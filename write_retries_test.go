package seatsio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v13/shared"
	"github.com/stretchr/testify/require"
)

type fakeClock struct {
	mu            sync.Mutex
	now           time.Time
	sleeps        int
	extraPerSleep time.Duration
	onSleep       func(sleeps int)
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	c.now = c.now.Add(d + c.extraPerSleep)
	c.sleeps++
	sleeps, onSleep := c.sleeps, c.onSleep
	c.mu.Unlock()
	if onSleep != nil {
		onSleep(sleeps)
	}
	return nil
}

func (c *fakeClock) sleepCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sleeps
}

type testClient struct {
	*req.Client
	retries *shared.Retries
	clock   *fakeClock
}

func newTestClient(baseUrl string, headers ...shared.AdditionalHeader) *testClient {
	clock := &fakeClock{now: time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC)}
	retries := shared.NewRetries(clock)
	return &testClient{Client: shared.ApiClientWithRetries("aSecretKey", baseUrl, retries, headers...), retries: retries, clock: clock}
}

func (c *testClient) withMaxRetries(count int) *testClient {
	c.retries.SetMaxRetries(c.Client, count)
	return c
}

func (c *testClient) post(path string) error {
	response, err := c.R().Post(path)
	return shared.AssertOkWithoutResult(response, err)
}

type recordedRequest struct {
	method string
	path   string
	query  string
	body   string
	header http.Header
}

type requestRecorder struct {
	mu       sync.Mutex
	requests []recordedRequest
}

func (r *requestRecorder) record(request *http.Request) int {
	body, _ := io.ReadAll(request.Body)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests = append(r.requests, recordedRequest{
		method: request.Method,
		path:   request.URL.Path,
		query:  request.URL.RawQuery,
		body:   string(body),
		header: request.Header.Clone(),
	})
	return len(r.requests) - 1
}

func (r *requestRecorder) all() []recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedRequest(nil), r.requests...)
}

func (r *requestRecorder) count() int {
	return len(r.all())
}

func (r *requestRecorder) keys() []string {
	var keys []string
	for _, request := range r.all() {
		keys = append(keys, request.header.Get("Idempotency-Key"))
	}
	return keys
}

type step func(w http.ResponseWriter)

func reply(status int, keyStatus string, body string) step {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		if keyStatus != "" {
			w.Header().Set(shared.IdempotencyKeyStatusHeader, keyStatus)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

func lost() step {
	return dropConnection
}

func apiError(code string) string {
	return `{"errors":[{"code":"` + code + `","message":"an error"}]}`
}

const eventBody = `{"key":"anEvent"}`

var (
	executed = reply(http.StatusOK, "executed", eventBody)
	ignored  = reply(http.StatusOK, "ignored", eventBody)
	noStatus = reply(http.StatusOK, "", eventBody)
	keyInUse = reply(http.StatusConflict, "in-use", apiError("IDEMPOTENCY_KEY_IN_USE"))
	tooMany  = reply(http.StatusTooManyRequests, "", "")
)

func scriptedHandler(steps []step, requests *requestRecorder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := requests.record(r)
		steps[min(i, len(steps)-1)](w)
	})
}

func scriptedServer(t *testing.T, steps ...step) (*httptest.Server, *requestRecorder) {
	t.Helper()
	requests := &requestRecorder{}
	server := httptest.NewServer(scriptedHandler(steps, requests))
	t.Cleanup(server.Close)
	return server, requests
}

func reservedAddress(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	return address
}

func startScriptedServerOn(t *testing.T, address string, requests *requestRecorder, steps ...step) {
	t.Helper()
	listener, err := net.Listen("tcp", address)
	require.NoError(t, err)
	server := httptest.NewUnstartedServer(scriptedHandler(steps, requests))
	_ = server.Listener.Close()
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
}

func closedServerURL(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.NotFoundHandler())
	server.Close()
	return server.URL
}

func requireOutcome(t *testing.T, err error, outcomeUnknown bool, safeToRetryWithSameKey bool) {
	t.Helper()
	require.Error(t, err)
	require.Equal(t, outcomeUnknown, ClientSupport.IsOutcomeUnknown(err), "outcome unknown")
	require.Equal(t, safeToRetryWithSameKey, ClientSupport.IsSafeToRetryWithSameKey(err), "safe to retry with the same key")
}

func requireSameKey(t *testing.T, keys []string) {
	t.Helper()
	require.NotEmpty(t, keys)
	_, err := shared.ParseIdempotencyKey(keys[0])
	require.NoError(t, err)
	for _, key := range keys[1:] {
		require.Equal(t, keys[0], key)
	}
}

func TestALostResponseIsRetriedWithTheSameKeyWhenKeysAreHonored(t *testing.T) {
	t.Parallel()
	for _, method := range []string{http.MethodPost, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, lost(), reply(http.StatusOK, "replayed", `{"key":"replayed"}`))
			client := newTestClient(server.URL)
			_, err := client.R().Send(method, "/events")
			require.NoError(t, err)

			var event unparseableEvent
			response, err := client.R().SetSuccessResult(&event).Send(method, "/events")

			require.NoError(t, shared.AssertOkWithoutResult(response, err))
			require.Equal(t, "replayed", event.Key)
			require.Equal(t, 3, requests.count())
			requireSameKey(t, requests.keys()[1:])
			require.NotEqual(t, requests.keys()[0], requests.keys()[1])
			require.Equal(t, 1, client.clock.sleepCount())
		})
	}
}

func TestALostResponseIsNotRetriedWhenKeysAreIgnored(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, ignored, lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	err := client.post("/events")

	require.Equal(t, 2, requests.count())
	requireOutcome(t, err, true, false)
	key, ok := ClientSupport.IdempotencyKeyOf(err)
	require.True(t, ok)
	require.Equal(t, requests.keys()[1], key.String())
}

func TestALostResponseIsNotRetriedWhenNoResponseHadAStatus(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, noStatus, lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	err := client.post("/events")

	require.Equal(t, 2, requests.count())
	requireOutcome(t, err, true, false)
}

func TestTheFirstRequestOfAClientIsNotRetriedAfterALostResponse(t *testing.T) {
	t.Parallel()
	for _, method := range writeMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, lost(), executed)
			client := newTestClient(server.URL)

			response, err := client.R().Send(method, "/events")
			err = shared.AssertOkWithoutResult(response, err)

			require.Equal(t, 1, requests.count())
			requireOutcome(t, err, true, false)
			require.Zero(t, client.clock.sleepCount())
		})
	}
}

func TestARequestThatWasNeverSentIsRetriedWithTheSameKey(t *testing.T) {
	t.Parallel()
	address := reservedAddress(t)
	requests := &requestRecorder{}
	client := newTestClient("http://" + address)
	client.clock.onSleep = func(int) { startScriptedServerOn(t, address, requests, executed) }

	require.NoError(t, client.post("/events"))

	require.Equal(t, 1, requests.count())
	require.Equal(t, 1, client.clock.sleepCount())
	requireSameKey(t, requests.keys())
}

func TestARequestThatWasNeverSentIsSafeToRetryWithTheSameKey(t *testing.T) {
	t.Parallel()
	client := newTestClient(closedServerURL(t)).withMaxRetries(0)

	err := client.post("/events")

	requireOutcome(t, err, false, true)
	var requestErr *shared.RequestError
	require.ErrorAs(t, err, &requestErr)
	require.Nil(t, requestErr.IdempotencyKeyHonored)
	require.False(t, requestErr.Replayed)
	_, ok := ClientSupport.IdempotencyKeyOf(err)
	require.True(t, ok)
}

func TestEveryKindOfRequestThatWasNeverSentIsSafeToRetryWithTheSameKey(t *testing.T) {
	t.Parallel()
	tests := map[string]func(t *testing.T) (client *testClient, ctx context.Context, requests *requestRecorder){
		"unknown host": func(t *testing.T) (*testClient, context.Context, *requestRecorder) {
			return newTestClient("http://seatsio-go-test.invalid"), context.Background(), nil
		},
		"untrusted certificate": func(t *testing.T) (*testClient, context.Context, *requestRecorder) {
			requests := &requestRecorder{}
			server := httptest.NewTLSServer(scriptedHandler([]step{executed}, requests))
			t.Cleanup(server.Close)
			return newTestClient(server.URL), context.Background(), requests
		},
		"connect timeout": func(t *testing.T) (*testClient, context.Context, *requestRecorder) {
			server, requests := scriptedServer(t, executed)
			client := newTestClient(server.URL)
			client.SetDial((&net.Dialer{Timeout: time.Nanosecond}).DialContext)
			return client, context.Background(), requests
		},
		"context cancelled before connecting": func(t *testing.T) (*testClient, context.Context, *requestRecorder) {
			server, requests := scriptedServer(t, executed)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			return newTestClient(server.URL), ctx, requests
		},
	}
	for name, setUp := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client, ctx, requests := setUp(t)
			client.withMaxRetries(0)

			response, err := client.R().SetContext(ctx).Post("/events")
			err = shared.AssertOkWithoutResult(response, err)

			requireOutcome(t, err, false, true)
			key, ok := ClientSupport.IdempotencyKeyOf(err)
			require.True(t, ok)
			require.Equal(t, sentIdempotencyKey(t, response), key.String())
			if requests != nil {
				require.Zero(t, requests.count())
			}
		})
	}
}

func sentIdempotencyKey(t *testing.T, response *req.Response) string {
	t.Helper()
	require.NotNil(t, response)
	require.NotNil(t, response.Request)
	keys := response.Request.Headers["idempotency-key"]
	require.Len(t, keys, 1)
	return keys[0]
}

func TestThe30SecondLimitStopsApplyingOnceTheServerAnswers(t *testing.T) {
	t.Parallel()
	address := reservedAddress(t)
	requests := &requestRecorder{}
	client := newTestClient("http://" + address).withMaxRetries(5)
	client.clock.extraPerSleep = 10 * time.Second
	client.clock.onSleep = func(sleeps int) {
		if sleeps == 1 {
			startScriptedServerOn(t, address, requests, tooMany)
		}
	}

	err := client.post("/events")

	require.Equal(t, 5, client.clock.sleepCount())
	require.Equal(t, 5, requests.count())
	require.Greater(t, client.clock.Now().Sub(time.Date(2026, 11, 24, 12, 0, 0, 0, time.UTC)), 30*time.Second)
	requireOutcome(t, err, false, false)
}

func TestAnIncompleteBodyKeepsTheKeyStatusOfItsHeaders(t *testing.T) {
	t.Parallel()
	for name, incomplete := range incompleteBodyResponses() {
		for _, keyStatus := range []string{"executed", "replayed"} {
			t.Run(name+", "+keyStatus, func(t *testing.T) {
				t.Parallel()
				server, _ := scriptedServer(t, func(w http.ResponseWriter) {
					w.Header().Set(shared.IdempotencyKeyStatusHeader, keyStatus)
					incomplete(w)
				})
				client := newTestClient(server.URL).withMaxRetries(0)
				client.SetTimeout(150 * time.Millisecond)

				var event unparseableEvent
				response, err := client.R().SetSuccessResult(&event).Post("/events")
				_, err = shared.AssertOk(response, err, &event)

				var requestErr *shared.RequestError
				require.ErrorAs(t, err, &requestErr)
				require.NotNil(t, requestErr.IdempotencyKeyHonored)
				require.True(t, *requestErr.IdempotencyKeyHonored)
				require.Equal(t, keyStatus == "replayed", requestErr.Replayed)
				requireOutcome(t, err, true, true)
			})
		}
	}
}

func TestALostResponseIsSafeToRetryWithTheSameKeyWhenKeysAreHonored(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost())
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	client.withMaxRetries(0)

	err := client.post("/events")

	require.Equal(t, 2, requests.count())
	requireOutcome(t, err, true, true)
}

func TestTheServersAnswerForOneOperationAppliesToAllOperations(t *testing.T) {
	t.Parallel()
	for name, second := range map[string]func(client *testClient) error{
		"another path": func(client *testClient) error { return client.post("/charts") },
		"another method": func(client *testClient) error {
			response, err := client.R().Delete("/events/anEvent")
			return shared.AssertOkWithoutResult(response, err)
		},
		"a PUT": func(client *testClient) error {
			response, err := client.R().Put("/events/anEvent")
			return shared.AssertOkWithoutResult(response, err)
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, lost(), executed)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))

			require.NoError(t, second(client))

			require.Equal(t, 3, requests.count())
			requireSameKey(t, requests.keys()[1:])
		})
	}
}

func TestAnOperationThatIgnoresKeysStopsRetriesOfAllOperations(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, ignored, lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	require.NoError(t, client.post("/charts"))

	err := client.post("/events")

	require.Equal(t, 3, requests.count())
	requireOutcome(t, err, true, false)
}

func TestTheLatestKnownStatusWins(t *testing.T) {
	t.Parallel()
	for name, second := range map[string]step{
		"200": reply(http.StatusOK, "ignored", eventBody),
		"400": reply(http.StatusBadRequest, "ignored", apiError("BAD")),
		"500": reply(http.StatusInternalServerError, "ignored", apiError("BOOM")),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, second, lost(), executed)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))
			_ = client.post("/events")

			err := client.post("/events")

			require.Equal(t, 3, requests.count())
			requireOutcome(t, err, true, false)
		})
	}
}

func TestAResponseWithoutAStatusLeavesTheMemoryAlone(t *testing.T) {
	t.Parallel()
	for name, second := range map[string]step{
		"401": reply(http.StatusUnauthorized, "", apiError("UNAUTHORIZED")),
		"429": tooMany,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, second, lost(), executed)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))
			client.withMaxRetries(0)
			_ = client.post("/events")
			client.withMaxRetries(5)

			require.NoError(t, client.post("/events"))

			require.Equal(t, 4, requests.count())
		})
	}
}

func TestSetDefaultWorkspaceIsRetriedLikeAnyOtherOperation(t *testing.T) {
	t.Parallel()
	noContent := reply(http.StatusNoContent, "executed", "")
	server, requests := scriptedServer(t, noContent, lost(), noContent)
	client := newTestClient(server.URL)
	_, err := client.R().SetPathParam("key", "workspace1").Post("/workspaces/actions/set-default/{key}")
	require.NoError(t, err)

	_, err = client.R().SetPathParam("key", "workspace2").Post("/workspaces/actions/set-default/{key}")

	require.NoError(t, err)
	require.Equal(t, 3, requests.count())
	all := requests.all()
	require.Equal(t, "/workspaces/actions/set-default/workspace2", all[2].path)
	require.NotEmpty(t, all[2].header.Get("Idempotency-Key"))
	require.Equal(t, all[1].header.Get("Idempotency-Key"), all[2].header.Get("Idempotency-Key"))
}

func TestClientsDontShareTheMemory(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost(), executed)
	require.NoError(t, newTestClient(server.URL).post("/events"))

	err := newTestClient(server.URL).post("/events")

	require.Equal(t, 2, requests.count())
	requireOutcome(t, err, true, false)
}

func TestARetrySendsTheSameRequest(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost(), executed)
	client := newTestClient(server.URL, ClientSupport.WorkspaceKey("aWorkspace"))
	require.NoError(t, client.post("/events/{event}/actions/book"))

	_, err := client.R().
		SetPathParam("event", "anEvent").
		SetQueryParam("expand", "x").
		SetBodyJsonString(`{"objects":["A-1"]}`).
		Post("/events/{event}/actions/book")

	require.NoError(t, err)
	all := requests.all()
	require.Len(t, all, 3)
	for _, request := range all[1:] {
		require.Equal(t, http.MethodPost, request.method)
		require.Equal(t, "/events/anEvent/actions/book", request.path)
		require.Equal(t, "expand=x", request.query)
		require.Equal(t, `{"objects":["A-1"]}`, request.body)
		require.Equal(t, []string{"aWorkspace"}, request.header.Values("X-Workspace-Key"))
		require.Len(t, request.header.Values("Authorization"), 1)
		require.Len(t, request.header.Values("Idempotency-Key"), 1)
	}
	require.Equal(t, all[1].header.Get("Idempotency-Key"), all[2].header.Get("Idempotency-Key"))
}

func TestAnUnknownStatusLeavesTheMemoryAlone(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, reply(http.StatusOK, "maybe", eventBody), lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	require.NoError(t, client.post("/events"))

	require.NoError(t, client.post("/events"))

	require.Equal(t, 4, requests.count())
}

func TestAnUnknownStatusDoesNotMarkKeysAsHonored(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, reply(http.StatusOK, "maybe", eventBody), lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	require.Error(t, client.post("/events"))

	require.Equal(t, 2, requests.count())
}

func TestTheStatusIsCaseInsensitive(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, reply(http.StatusOK, " Executed ", eventBody), lost(), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	require.NoError(t, client.post("/events"))

	require.Equal(t, 3, requests.count())
}

func TestStatusesThatMarkKeysAsHonored(t *testing.T) {
	t.Parallel()
	for name, first := range map[string]step{
		"in-use":   keyInUse,
		"reused":   reply(http.StatusUnprocessableEntity, "reused", apiError("IDEMPOTENCY_KEY_REUSED")),
		"replayed": reply(http.StatusOK, "replayed", eventBody),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, first, lost(), executed)
			client := newTestClient(server.URL).withMaxRetries(0)
			_ = client.post("/events")
			client.withMaxRetries(1)

			require.NoError(t, client.post("/events"))

			require.Equal(t, 3, requests.count())
		})
	}
}

func TestAReplayedErrorIsReportedAsReplayed(t *testing.T) {
	t.Parallel()
	server, _ := scriptedServer(t, reply(http.StatusBadRequest, "replayed", apiError("BAD")))

	err := newTestClient(server.URL).post("/events")

	var apiErr *shared.SeatsioError
	require.ErrorAs(t, err, &apiErr)
	require.True(t, apiErr.Replayed)
	require.NotNil(t, apiErr.IdempotencyKeyHonored)
	require.True(t, *apiErr.IdempotencyKeyHonored)
	requireOutcome(t, err, false, false)
}

func TestAnErrorWithoutAStatusSaysNothingAboutTheKey(t *testing.T) {
	t.Parallel()
	server, _ := scriptedServer(t, reply(http.StatusBadRequest, "", apiError("BAD")))

	err := newTestClient(server.URL).post("/events")

	var apiErr *shared.SeatsioError
	require.ErrorAs(t, err, &apiErr)
	require.False(t, apiErr.Replayed)
	require.Nil(t, apiErr.IdempotencyKeyHonored)
	require.NotEmpty(t, apiErr.IdempotencyKey)
}

func TestAnIncompleteBodyIsRetriedWhenKeysAreHonored(t *testing.T) {
	t.Parallel()
	for name, incomplete := range incompleteBodyResponses() {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, func(w http.ResponseWriter) { incomplete(w) }, executed)
			client := newTestClient(server.URL)
			client.SetTimeout(150 * time.Millisecond)
			require.NoError(t, client.post("/events"))

			var event unparseableEvent
			response, err := client.R().SetSuccessResult(&event).Post("/events")

			require.NoError(t, shared.AssertOkWithoutResult(response, err))
			require.Equal(t, "anEvent", event.Key)
			require.Equal(t, 3, requests.count())
			requireSameKey(t, requests.keys()[1:])
		})
	}
}

func TestOutcomeIsUnknownWhenASuccessfulResponseBodyIsIncomplete(t *testing.T) {
	t.Parallel()
	for name, incomplete := range incompleteBodyResponses() {
		for _, honored := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s, keys honored %v", name, honored), func(t *testing.T) {
				t.Parallel()
				first := ignored
				if honored {
					first = executed
				}
				server, requests := scriptedServer(t, first, func(w http.ResponseWriter) { incomplete(w) })
				client := newTestClient(server.URL).withMaxRetries(1)
				client.SetTimeout(150 * time.Millisecond)
				require.NoError(t, client.post("/events"))

				var event unparseableEvent
				response, err := client.R().SetSuccessResult(&event).Post("/events")
				result, err := shared.AssertOk(response, err, &event)

				require.Nil(t, result)
				requireOutcome(t, err, true, honored)
				key, ok := ClientSupport.IdempotencyKeyOf(err)
				require.True(t, ok)
				require.Equal(t, requests.keys()[len(requests.keys())-1], key.String())
			})
		}
	}
}

func TestLostResponsesCountAgainstMaxRetries(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost(), lost(), lost(), executed)
	client := newTestClient(server.URL).withMaxRetries(2)
	require.NoError(t, client.post("/events"))

	err := client.post("/events")

	require.Equal(t, 4, requests.count())
	require.Equal(t, 2, client.clock.sleepCount())
	requireOutcome(t, err, true, true)
}

func TestKeyInUseErrorIsReturnedWhenRetriesRunOut(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, keyInUse)
	client := newTestClient(server.URL).withMaxRetries(2)

	err := client.post("/events")

	require.Equal(t, 3, requests.count())
	requireSameKey(t, requests.keys())
	var apiErr *shared.SeatsioError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, "IDEMPOTENCY_KEY_IN_USE", apiErr.Code)
	require.Equal(t, requests.keys()[0], apiErr.IdempotencyKey)
	requireOutcome(t, err, true, true)
}

func TestARequestThatWasNeverSentStopsBeingRetried30SecondsAfterTheFirstAttempt(t *testing.T) {
	t.Parallel()
	client := newTestClient(closedServerURL(t)).withMaxRetries(10)
	client.clock.extraPerSleep = 10 * time.Second

	err := client.post("/events")

	require.Equal(t, 3, client.clock.sleepCount())
	requireOutcome(t, err, false, true)
}

func TestARequestThatWasNeverSentIsNotRetriedWhenTheWaitWouldEndPast30Seconds(t *testing.T) {
	t.Parallel()
	client := newTestClient(closedServerURL(t)).withMaxRetries(10)
	start := client.clock.Now()
	client.clock.onSleep = func(int) {
		client.clock.mu.Lock()
		client.clock.now = start.Add(29500 * time.Millisecond)
		client.clock.mu.Unlock()
	}

	err := client.post("/events")

	require.Equal(t, 1, client.clock.sleepCount())
	requireOutcome(t, err, false, true)
}

func TestMaxRetriesZeroGivesOneAttemptForEveryKindOfFailure(t *testing.T) {
	t.Parallel()
	for name, failure := range map[string]step{"429": tooMany, "key in use": keyInUse} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, failure, executed)
			client := newTestClient(server.URL).withMaxRetries(0)

			require.Error(t, client.post("/events"))

			require.Equal(t, 1, requests.count())
			require.Zero(t, client.clock.sleepCount())
		})
	}
	t.Run("never sent", func(t *testing.T) {
		t.Parallel()
		client := newTestClient(closedServerURL(t)).withMaxRetries(0)

		require.Error(t, client.post("/events"))

		require.Zero(t, client.clock.sleepCount())
	})
	t.Run("lost when keys are honored", func(t *testing.T) {
		t.Parallel()
		server, requests := scriptedServer(t, executed, lost(), executed)
		client := newTestClient(server.URL).withMaxRetries(0)
		require.NoError(t, client.post("/events"))

		require.Error(t, client.post("/events"))

		require.Equal(t, 2, requests.count())
		require.Zero(t, client.clock.sleepCount())
	})
}

func TestAllKindsOfFailureShareOneRetryCount(t *testing.T) {
	t.Parallel()
	address := reservedAddress(t)
	requests := &requestRecorder{}
	client := newTestClient("http://" + address).withMaxRetries(3)
	client.clock.onSleep = func(sleeps int) {
		if sleeps == 1 {
			startScriptedServerOn(t, address, requests, tooMany, keyInUse, lost(), executed)
		}
	}

	err := client.post("/events")

	require.Equal(t, 3, requests.count())
	require.Equal(t, 3, client.clock.sleepCount())
	requireSameKey(t, requests.keys())
	requireOutcome(t, err, true, true)
}

func TestWritesAreRetriedFiveTimesByDefault(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, tooMany)
	client := newTestClient(server.URL)

	require.Error(t, client.post("/events"))

	require.Equal(t, 1+5, requests.count())
}

func TestSetMaxRetriesAppliesToWritesOfItsOwnClientOnly(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, tooMany)
	clock := &fakeClock{}
	client := newSeatsioClient(server.URL, "aSecretKey", shared.NewRetries(clock))
	otherClient := newSeatsioClient(server.URL, "aSecretKey", shared.NewRetries(clock))

	require.NoError(t, client.SetMaxRetries(2))
	require.NoError(t, otherClient.SetMaxRetries(0))
	_, _ = client.apiClient.R().Post("/events")

	require.Equal(t, 3, requests.count())
}

func TestLostResponsesAreNotRetriedOnceTheKeyMayHaveExpired(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost())
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	client.clock.extraPerSleep = 2 * time.Minute

	err := client.post("/events")

	require.Equal(t, 1+3, requests.count())
	require.Equal(t, 3, client.clock.sleepCount())
	requireOutcome(t, err, true, false)
}

func TestKeyInUseIsNotSafeToRetryOnceTheRetryWaitEndsBeyondTheWindow(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, keyInUse)
	client := newTestClient(server.URL)
	client.clock.extraPerSleep = 11 * time.Minute

	err := client.post("/events")

	require.Equal(t, 1, requests.count())
	requireOutcome(t, err, true, false)
}

func TestA429OnlyWriteKeepsRetryingPastTheRetryWindow(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, tooMany)
	client := newTestClient(server.URL)
	client.clock.extraPerSleep = 2 * time.Minute

	err := client.post("/events")

	require.Equal(t, 1+5, requests.count())
	requireOutcome(t, err, false, false)
}

func TestARetryDoesNotStartWhenTheWaitEndsPastTheWindow(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, keyInUse, executed)
	client := newTestClient(server.URL)
	client.clock.extraPerSleep = 6 * time.Minute

	err := client.post("/events")

	require.Equal(t, 1, requests.count())
	require.Equal(t, 1, client.clock.sleepCount())
	requireOutcome(t, err, true, false)
}

func TestA429AfterAnAttemptThatMayHaveRunLeavesTheOutcomeUnknown(t *testing.T) {
	t.Parallel()
	for name, first := range map[string]step{"lost": lost(), "key in use": keyInUse} {
		t.Run(name+", within the window", func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, first, tooMany)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))
			client.withMaxRetries(2)

			err := client.post("/events")

			require.Equal(t, 1+3, requests.count())
			requireSameKey(t, requests.keys()[1:])
			requireOutcome(t, err, true, true)
		})
		t.Run(name+", past the window", func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, first, tooMany)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))
			client.clock.extraPerSleep = 2 * time.Minute

			err := client.post("/events")

			require.Equal(t, 1+3, requests.count())
			requireOutcome(t, err, true, false)
		})
	}
}

func TestAWriteThatMayHaveRunIsNotResentWhenAnotherOperationIgnoresKeysDuringTheWait(t *testing.T) {
	t.Parallel()
	for name, steps := range map[string]struct {
		script     []step
		otherAfter int
		safe       bool
	}{
		"after a lost response":    {[]step{executed, lost(), ignored, executed}, 1, false},
		"after a 429 that follows": {[]step{executed, lost(), tooMany, ignored, executed}, 2, false},
		"after a key in use":       {[]step{executed, keyInUse, ignored, executed}, 1, true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, steps.script...)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))
			client.clock.onSleep = func(sleeps int) {
				if sleeps == steps.otherAfter {
					require.NoError(t, client.post("/charts"))
				}
			}

			err := client.post("/events")

			require.Equal(t, len(steps.script)-1, requests.count())
			requireOutcome(t, err, true, steps.safe)
		})
	}
}

func TestAWriteThatMayHaveRunIsNotResentAfterA429ThatIgnoresKeys(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost(), reply(http.StatusTooManyRequests, "ignored", ""), executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	err := client.post("/events")

	require.Equal(t, 3, requests.count())
	requireOutcome(t, err, true, false)
}

func TestKeyInUseWithoutAStatusIsRetriedBeforeKeysAreKnownToBeHonored(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, reply(http.StatusConflict, "", apiError("IDEMPOTENCY_KEY_IN_USE")), executed)
	client := newTestClient(server.URL)

	require.NoError(t, client.post("/events"))

	require.Equal(t, 2, requests.count())
	requireSameKey(t, requests.keys())
}

func TestA429IsRetriedWhenKeysAreIgnored(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, ignored, tooMany, ignored)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	require.NoError(t, client.post("/events"))

	require.Equal(t, 3, requests.count())
}

func TestAnErrorAfterAnAttemptThatMayHaveRunLeavesTheOutcomeUnknown(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		final          step
		outcomeUnknown bool
		safe           bool
	}{
		"400 without status": {reply(http.StatusBadRequest, "", apiError("BAD")), true, true},
		"403 from a proxy":   {reply(http.StatusForbidden, "", "forbidden"), true, true},
		"400 executed":       {reply(http.StatusBadRequest, "executed", apiError("BAD")), true, true},
		"400 ignored":        {reply(http.StatusBadRequest, "ignored", apiError("BAD")), true, false},
		"400 replayed":       {reply(http.StatusBadRequest, "replayed", apiError("BAD")), false, false},
		"422 reused":         {reply(http.StatusUnprocessableEntity, "reused", apiError("IDEMPOTENCY_KEY_REUSED")), true, true},
		"400 invalid key":    {reply(http.StatusBadRequest, "", apiError("INVALID_IDEMPOTENCY_KEY")), true, true},
		"404":                {reply(http.StatusNotFound, "", apiError("NOT_FOUND")), true, true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, executed, lost(), test.final)
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))

			err := client.post("/events")

			require.Equal(t, 3, requests.count())
			requireOutcome(t, err, test.outcomeUnknown, test.safe)
		})
	}
}

func TestA429OnTheFirstAttemptLeavesTheOutcomeKnown(t *testing.T) {
	t.Parallel()
	server, _ := scriptedServer(t, executed, tooMany)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	client.withMaxRetries(2)

	err := client.post("/events")

	requireOutcome(t, err, false, false)
}

func TestWritesAreRetriedOn429AndKeyInUseWithTheSameKey(t *testing.T) {
	t.Parallel()
	for name, failure := range map[string]step{"429": tooMany, "key in use": keyInUse} {
		for _, method := range writeMethods {
			t.Run(name+" "+method, func(t *testing.T) {
				t.Parallel()
				server, requests := scriptedServer(t, failure, reply(http.StatusNoContent, "executed", ""))
				client := newTestClient(server.URL)

				response, err := client.R().Send(method, "/events")

				require.NoError(t, err)
				require.Equal(t, http.StatusNoContent, response.StatusCode)
				require.Equal(t, 2, requests.count())
				requireSameKey(t, requests.keys())
			})
		}
	}
}

func TestOtherConflictsAreNotRetried(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"other code":     apiError("SOMETHING_ELSE"),
		"no errors":      `{"errors":[]}`,
		"empty object":   `{}`,
		"malformed json": `{"errors":`,
		"empty body":     ``,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, requests := scriptedServer(t, reply(http.StatusConflict, "", body), executed)

			err := newTestClient(server.URL).post("/events")

			require.Equal(t, 1, requests.count())
			requireOutcome(t, err, false, false)
		})
	}
}

func TestReusedIsNotRetried(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, reply(http.StatusUnprocessableEntity, "reused", apiError("IDEMPOTENCY_KEY_REUSED")), executed)

	err := newTestClient(server.URL).post("/events")

	require.Equal(t, 1, requests.count())
	var apiErr *shared.SeatsioError
	require.ErrorAs(t, err, &apiErr)
	require.True(t, *apiErr.IdempotencyKeyHonored)
	requireOutcome(t, err, false, false)
}

func TestCallerKeySurvivesARetry(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, tooMany, executed)
	key := ClientSupport.NewIdempotencyKey()

	_, err := newTestClient(server.URL).R().SetContext(ClientSupport.WithIdempotencyKey(context.Background(), key)).Post("/events")

	require.NoError(t, err)
	require.Equal(t, []string{key.String(), key.String()}, requests.keys())
}

func TestErrorValuesOfErrorResponses(t *testing.T) {
	t.Parallel()
	yes, no := true, false
	tests := map[string]struct {
		learned        step
		response       step
		outcomeUnknown bool
		safe           bool
		honored        *bool
		replayed       bool
	}{
		"5xx executed":                  {nil, reply(500, "executed", apiError("BOOM")), true, true, &yes, false},
		"5xx in-use":                    {nil, reply(500, "in-use", apiError("BOOM")), true, true, &yes, false},
		"5xx ignored":                   {nil, reply(500, "ignored", apiError("BOOM")), true, false, &no, false},
		"5xx without status, unknown":   {nil, reply(502, "", "bad gateway"), true, false, nil, false},
		"5xx without status, honored":   {executed, reply(502, "", "bad gateway"), true, true, nil, false},
		"5xx without status, ignored":   {ignored, reply(502, "", "bad gateway"), true, false, nil, false},
		"5xx replayed":                  {nil, reply(500, "replayed", apiError("BOOM")), true, false, &yes, true},
		"400":                           {nil, reply(400, "executed", apiError("BAD")), false, false, &yes, false},
		"404":                           {nil, reply(404, "", apiError("NOT_FOUND")), false, false, nil, false},
		"409":                           {nil, reply(409, "executed", apiError("CONFLICT")), false, false, &yes, false},
		"422 reused":                    {nil, reply(422, "reused", apiError("IDEMPOTENCY_KEY_REUSED")), false, false, &yes, false},
		"key in use after the last try": {nil, keyInUse, true, true, &yes, false},
		"400 invalid idempotency key":   {nil, reply(400, "", apiError("INVALID_IDEMPOTENCY_KEY")), false, false, nil, false},
		"non-json 5xx, keys honored":    {executed, reply(503, "", "unavailable"), true, true, nil, false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			steps := []step{test.response}
			if test.learned != nil {
				steps = []step{test.learned, test.response}
			}
			server, _ := scriptedServer(t, steps...)
			client := newTestClient(server.URL).withMaxRetries(0)
			if test.learned != nil {
				require.NoError(t, client.post("/events"))
			}

			err := client.post("/events")

			requireOutcome(t, err, test.outcomeUnknown, test.safe)
			var honored *bool
			var replayed bool
			var apiErr *shared.SeatsioError
			var requestErr *shared.RequestError
			if errors.As(err, &apiErr) {
				honored, replayed = apiErr.IdempotencyKeyHonored, apiErr.Replayed
			} else {
				require.ErrorAs(t, err, &requestErr)
				honored, replayed = requestErr.IdempotencyKeyHonored, requestErr.Replayed
			}
			require.Equal(t, test.honored, honored)
			require.Equal(t, test.replayed, replayed)
			_, ok := ClientSupport.IdempotencyKeyOf(err)
			require.True(t, ok)
		})
	}
}

func TestAnOperationWithUnknownOutcomeCanBeRetriedWithTheKeyOfTheError(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed, lost(), reply(http.StatusOK, "replayed", eventBody))
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))
	client.withMaxRetries(0)

	err := client.post("/events")
	require.True(t, ClientSupport.IsSafeToRetryWithSameKey(err))
	key, _ := ClientSupport.IdempotencyKeyOf(err)
	response, err := client.R().SetContext(ClientSupport.WithIdempotencyKey(context.Background(), key)).Post("/events")

	require.NoError(t, shared.AssertOkWithoutResult(response, err))
	require.Equal(t, 3, requests.count())
	require.Equal(t, requests.keys()[1], requests.keys()[2])
}

func TestGetsAreRetriedByReqOnly(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, lost(), executed)
	client := newTestClient(server.URL)

	response, err := client.R().Get("/events/anEvent")

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, 2, requests.count())
	require.Equal(t, 1, client.clock.sleepCount())
	require.Equal(t, []string{"", ""}, requests.keys())
}

func TestAnUnknownFailureIsNotRetriedBeforeKeysAreKnownToBeHonored(t *testing.T) {
	t.Parallel()
	client := newTestClient("unknown-scheme://localhost")

	err := client.post("/events")

	require.Zero(t, client.clock.sleepCount())
	requireOutcome(t, err, true, false)
}

func TestAnUnknownFailureIsRetriedWhenKeysAreHonored(t *testing.T) {
	t.Parallel()
	server, _ := scriptedServer(t, executed)
	client := newTestClient(server.URL).withMaxRetries(2)
	require.NoError(t, client.post("/events"))

	err := client.post("unknown-scheme://localhost/events")

	require.Equal(t, 2, client.clock.sleepCount())
	requireOutcome(t, err, true, true)
}

func TestDoesNotRetryALostResponseWhenTheContextIsCancelled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, requests := scriptedServer(t, executed, func(w http.ResponseWriter) {
		cancel()
		dropConnection(w)
	}, executed)
	client := newTestClient(server.URL)
	require.NoError(t, client.post("/events"))

	response, err := client.R().SetContext(ctx).Post("/events")
	err = shared.AssertOkWithoutResult(response, err)

	require.Error(t, err)
	require.Equal(t, 2, requests.count())
	require.Zero(t, client.clock.sleepCount())
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}

func TestOutcomeIsUnknownWhenTheContextIsCancelledWhileTheBodyIsRead(t *testing.T) {
	t.Parallel()
	for _, honored := range []bool{false, true} {
		t.Run(fmt.Sprintf("keys honored %v", honored), func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			first := ignored
			if honored {
				first = executed
			}
			server, requests := scriptedServer(t, first, func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"key":`))
				w.(http.Flusher).Flush()
				time.Sleep(50 * time.Millisecond)
				cancel()
				time.Sleep(300 * time.Millisecond)
			})
			client := newTestClient(server.URL)
			require.NoError(t, client.post("/events"))

			var event unparseableEvent
			response, err := client.R().SetContext(ctx).SetSuccessResult(&event).Post("/events")
			_, err = shared.AssertOk(response, err, &event)

			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, 2, requests.count())
			requireOutcome(t, err, true, honored)
		})
	}
}

func TestAZeroKeyInTheContextFailsBeforeSending(t *testing.T) {
	t.Parallel()
	server, requests := scriptedServer(t, executed)

	_, err := newTestClient(server.URL).R().
		SetContext(ClientSupport.WithIdempotencyKey(context.Background(), shared.IdempotencyKey{})).
		Post("/events")

	require.ErrorContains(t, err, "zero value")
	require.Zero(t, requests.count())
}
