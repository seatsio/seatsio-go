package seatsio

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v13/shared"
	"github.com/stretchr/testify/require"
)

type keyRecorder struct {
	mu   sync.Mutex
	keys []string
}

func (k *keyRecorder) record(r *http.Request) {
	if k == nil {
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.keys = append(k.keys, r.Header.Get("Idempotency-Key"))
}

func (k *keyRecorder) all() []string {
	k.mu.Lock()
	defer k.mu.Unlock()
	return append([]string(nil), k.keys...)
}

var writeMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

func respondingServer(t *testing.T, status int, body string) (*httptest.Server, *keyRecorder) {
	t.Helper()
	return respondingServerWithContentType(t, status, "application/json", body)
}

func respondingServerWithContentType(t *testing.T, status int, contentType string, body string) (*httptest.Server, *keyRecorder) {
	t.Helper()
	var keys keyRecorder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys.record(r)
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server, &keys
}

func connectionDroppingServer(t *testing.T) (*httptest.Server, *keyRecorder) {
	t.Helper()
	var keys keyRecorder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys.record(r)
		dropConnection(w)
	}))
	t.Cleanup(server.Close)
	return server, &keys
}

func dropConnection(w http.ResponseWriter) {
	conn, _, err := w.(http.Hijacker).Hijack()
	if err == nil {
		_ = conn.Close()
	}
}

func TestDoesNotSendIdempotencyKeyOnGet(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusNoContent, "")

	_, err := shared.ApiClient("aSecretKey", server.URL).R().Get("/events/anEvent")

	require.NoError(t, err)
	require.Equal(t, []string{""}, keys.all())
}

func TestUsesIdempotencyKeyFromContext(t *testing.T) {
	t.Parallel()
	for _, method := range writeMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			server, keys := respondingServer(t, http.StatusNoContent, "")
			key := ClientSupport.NewIdempotencyKey()

			_, err := shared.ApiClient("aSecretKey", server.URL).
				R().
				SetContext(ClientSupport.WithIdempotencyKey(context.Background(), key)).
				Send(method, "/events/anEvent")

			require.NoError(t, err)
			require.Equal(t, []string{key.String()}, keys.all())
		})
	}
}

func TestParseIdempotencyKey(t *testing.T) {
	t.Parallel()
	valid := uuid.NewString()
	key, err := ClientSupport.ParseIdempotencyKey(valid)
	require.NoError(t, err)
	require.Equal(t, valid, key.String())
	key, err = ClientSupport.ParseIdempotencyKey(strings.ToUpper(valid))
	require.NoError(t, err)
	require.Equal(t, valid, key.String())

	for name, invalid := range map[string]string{
		"empty":                       "",
		"notAUuid":                    "not-a-uuid",
		"v1":                          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"nil":                         uuid.Nil.String(),
		"v4WithReservedNcsVariant":    "6ba7b810-9dad-41d1-00b4-00c04fd430c8",
		"v4WithMicrosoftVariant":      "6ba7b810-9dad-41d1-c0b4-00c04fd430c8",
		"v4WithReservedFutureVariant": "6ba7b810-9dad-41d1-e0b4-00c04fd430c8",
		"tooFewDigits":                "1-1-1-1-1",
		"braces":                      "{" + valid + "}",
		"urnPrefix":                   "urn:uuid:" + valid,
		"withoutDashes":               strings.ReplaceAll(valid, "-", ""),
		"leadingSpace":                " " + valid,
	} {
		_, err := ClientSupport.ParseIdempotencyKey(invalid)
		require.ErrorContains(t, err, "is not a v4 UUID", name)
	}
}

func TestIdempotencyKeyIsAvailableOnConnectionError(t *testing.T) {
	t.Parallel()
	server, keys := connectionDroppingServer(t)

	response, err := shared.ApiClient("aSecretKey", server.URL).R().Post("/events")
	e := shared.AssertOkWithoutResult(response, err)

	key, ok := shared.IdempotencyKeyOf(e)
	require.True(t, ok)
	require.Equal(t, keys.all()[0], key.String())
}

func TestGetErrorsAreNotWrapped(t *testing.T) {
	t.Parallel()
	server, _ := connectionDroppingServer(t)
	client := newTestClient(server.URL).withMaxRetries(0)

	response, err := client.R().Get("/events/anEvent")
	e := shared.AssertOkWithoutResult(response, err)

	var requestErr *shared.RequestError
	require.Error(t, e)
	require.False(t, errors.As(e, &requestErr))
	_, ok := shared.IdempotencyKeyOf(e)
	require.False(t, ok)
}

func TestIdempotencyKeyIsAvailableOnErrorResponses(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		status      int
		contentType string
		body        string
	}{
		"api error":           {http.StatusBadRequest, "application/json", `{"errors":[{"code":"BAD","message":"bad"}]}`},
		"malformed json":      {http.StatusInternalServerError, "application/json", `{"errors":`},
		"no errors in json":   {http.StatusInternalServerError, "application/json", `{}`},
		"non-json error body": {http.StatusBadGateway, "text/html", `<html>bad gateway</html>`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, keys := respondingServerWithContentType(t, test.status, test.contentType, test.body)

			response, err := shared.ApiClient("aSecretKey", server.URL).R().Post("/events")
			e := shared.AssertOkWithoutResult(response, err)

			key, ok := shared.IdempotencyKeyOf(e)
			require.True(t, ok)
			require.Equal(t, keys.all()[0], key.String())
		})
	}
}

func TestIdempotencyKeyOf(t *testing.T) {
	t.Parallel()
	k := uuid.NewString()
	tests := map[string]struct {
		err error
		key string
		ok  bool
	}{
		"nil":           {nil, "", false},
		"foreign error": {errors.New("boom"), "", false},
		"request error": {&shared.RequestError{IdempotencyKey: k, Err: errors.New("boom")}, k, true},
		"api error":     {&shared.SeatsioError{Code: "c", IdempotencyKey: k}, k, true},
		"wrapped":       {fmt.Errorf("wrapped: %w", &shared.SeatsioError{IdempotencyKey: k}), k, true},
		"joined":        {errors.Join(errors.New("a"), &shared.RequestError{IdempotencyKey: k, Err: errors.New("b")}), k, true},
		"empty key":     {&shared.RequestError{Err: errors.New("boom")}, "", false},
		"invalid key":   {&shared.RequestError{IdempotencyKey: "k", Err: errors.New("boom")}, "", false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			key, ok := ClientSupport.IdempotencyKeyOf(test.err)
			require.Equal(t, test.ok, ok)
			if test.ok {
				require.Equal(t, test.key, key.String())
			}
		})
	}
}

func TestErrorResponseWithoutErrorsIsAnUnexpectedResponse(t *testing.T) {
	t.Parallel()
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			server, _ := respondingServer(t, http.StatusBadRequest, `{"errors":[]}`)

			response, err := shared.ApiClient("aSecretKey", server.URL).R().Send(method, "/events")
			e := shared.AssertOkWithoutResult(response, err)

			var apiErr *shared.SeatsioError
			require.False(t, errors.As(e, &apiErr))
			require.ErrorContains(t, e, "server returned error 400")
		})
	}
}

type temporaryError bool

func (e temporaryError) Error() string   { return "temporary?" }
func (e temporaryError) Temporary() bool { return bool(e) }

func TestRequestErrorIsTemporaryWhenItsCauseIs(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		cause     error
		temporary bool
	}{
		"temporary":         {temporaryError(true), true},
		"not temporary":     {temporaryError(false), false},
		"wrapped temporary": {fmt.Errorf("wrapped: %w", temporaryError(true)), true},
		"no Temporary":      {errors.New("boom"), false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			e := &shared.RequestError{Err: test.cause}

			require.Equal(t, test.temporary, e.Temporary())
		})
	}
}

func TestGeneratesADifferentKeyForEachRequest(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusNoContent, "")
	client := shared.ApiClient("aSecretKey", server.URL)

	_, err := client.R().Post("/events")
	require.NoError(t, err)
	_, err = client.R().Post("/events")
	require.NoError(t, err)

	require.Len(t, keys.all(), 2)
	require.NotEmpty(t, keys.all()[0])
	require.NotEmpty(t, keys.all()[1])
	require.NotEqual(t, keys.all()[0], keys.all()[1])
}

func TestKeepsIdempotencyKeyHeaderSetByTheCaller(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusNoContent, "")
	callerKey := uuid.NewString()

	_, err := shared.ApiClient("aSecretKey", server.URL).
		R().
		SetHeader("Idempotency-Key", callerKey).
		SetContext(shared.WithIdempotencyKey(context.Background(), shared.NewIdempotencyKey())).
		Post("/events")

	require.NoError(t, err)
	require.Equal(t, []string{callerKey}, keys.all())
}

func TestFailedWriteKeepsTheCauseOfTheError(t *testing.T) {
	t.Parallel()
	server, _ := respondingServer(t, http.StatusNoContent, "")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	response, err := shared.ApiClient("aSecretKey", server.URL).R().SetContext(ctx).Post("/events")
	e := shared.AssertOkWithoutResult(response, err)

	var requestErr *shared.RequestError
	require.ErrorAs(t, e, &requestErr)
	require.ErrorIs(t, e, context.Canceled)
}

func TestTimedOutWriteIsATimeout(t *testing.T) {
	t.Parallel()
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })

	response, err := shared.ApiClient("aSecretKey", server.URL).SetTimeout(50 * time.Millisecond).R().Post("/events")
	e := shared.AssertOkWithoutResult(response, err)

	var requestErr *shared.RequestError
	require.ErrorAs(t, e, &requestErr)
	require.True(t, os.IsTimeout(e))
	var netErr net.Error
	require.True(t, errors.As(e, &netErr))
	require.True(t, netErr.Timeout())
}

func TestIsOutcomeUnknown(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		server  func(t *testing.T) *httptest.Server
		unknown bool
	}{
		"connection dropped": {func(t *testing.T) *httptest.Server { s, _ := connectionDroppingServer(t); return s }, true},
		"5xx api error": {func(t *testing.T) *httptest.Server {
			s, _ := respondingServer(t, http.StatusInternalServerError, `{"errors":[{"code":"BOOM","message":"boom"}]}`)
			return s
		}, true},
		"5xx non-json": {func(t *testing.T) *httptest.Server {
			s, _ := respondingServerWithContentType(t, http.StatusBadGateway, "text/html", "bad gateway")
			return s
		}, true},
		"4xx api error": {func(t *testing.T) *httptest.Server {
			s, _ := respondingServer(t, http.StatusBadRequest, `{"errors":[{"code":"BAD","message":"bad"}]}`)
			return s
		}, false},
		"4xx non-json": {func(t *testing.T) *httptest.Server {
			s, _ := respondingServerWithContentType(t, http.StatusNotFound, "text/html", "not found")
			return s
		}, false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server := test.server(t)

			response, err := newTestClient(server.URL).withMaxRetries(0).R().Post("/events")
			e := shared.AssertOkWithoutResult(response, err)

			require.Error(t, e)
			require.Equal(t, test.unknown, ClientSupport.IsOutcomeUnknown(e))
		})
	}
	require.False(t, ClientSupport.IsOutcomeUnknown(nil))
	require.False(t, ClientSupport.IsOutcomeUnknown(errors.New("boom")))
}

func TestGetErrorsHaveNoUnknownOutcome(t *testing.T) {
	t.Parallel()
	server, _ := connectionDroppingServer(t)

	response, err := newTestClient(server.URL).withMaxRetries(0).R().Get("/events/anEvent")
	e := shared.AssertOkWithoutResult(response, err)

	require.Error(t, e)
	require.False(t, ClientSupport.IsOutcomeUnknown(e))
}

func TestTransportDoesNotResendAWriteWhenAReusedConnectionBreaks(t *testing.T) {
	t.Parallel()
	for name, setKey := range map[string]func(r *req.Request){
		"generated key": func(r *req.Request) {},
		"key in context": func(r *req.Request) {
			r.SetContext(shared.WithIdempotencyKey(context.Background(), shared.NewIdempotencyKey()))
		},
		"key in header": func(r *req.Request) { r.SetHeader("Idempotency-Key", uuid.NewString()) },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if attempts.Add(1) == 1 {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				dropConnection(w)
			}))
			t.Cleanup(server.Close)
			client := newTestClient(server.URL).withMaxRetries(0)

			_, err := client.R().Post("/events")
			require.NoError(t, err)
			second := client.R()
			setKey(second)
			_, err = second.Post("/events")

			require.Error(t, err)
			require.Equal(t, int32(2), attempts.Load())
		})
	}
}

func TestSendsTheKeyOnTheWire(t *testing.T) {
	t.Parallel()
	var raw atomic.Value
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw.Store(r.Header.Get("Idempotency-Key"))
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	key := shared.NewIdempotencyKey()

	response, err := shared.ApiClient("aSecretKey", server.URL).R().
		SetContext(shared.WithIdempotencyKey(context.Background(), key)).
		Post("/events")

	require.NoError(t, err)
	require.Equal(t, key.String(), raw.Load())
	require.NoError(t, shared.AssertOkWithoutResult(response, err))
}

func incompleteBodyResponses() map[string]func(w http.ResponseWriter) {
	return map[string]func(w http.ResponseWriter){
		"truncated body": func(w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"key":`))
			w.(http.Flusher).Flush()
			dropConnection(w)
		},
		"stalled body": func(w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"key":`))
			w.(http.Flusher).Flush()
			time.Sleep(300 * time.Millisecond)
		},
	}
}

func TestOutcomeIsKnownWhenACompleteSuccessfulResponseBodyCannotBeParsed(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusOK, `{"key":`)

	var event unparseableEvent
	response, err := newTestClient(server.URL).withMaxRetries(0).R().SetSuccessResult(&event).Post("/events")
	_, err = shared.AssertOk(response, err, &event)

	var requestErr *shared.RequestError
	require.ErrorAs(t, err, &requestErr)
	requireOutcome(t, err, false, false)
	require.Len(t, keys.all(), 1)
	key, ok := shared.IdempotencyKeyOf(err)
	require.True(t, ok)
	require.Equal(t, keys.all()[0], key.String())
}

func endingContexts() map[string]struct {
	newContext func() (context.Context, context.CancelFunc)
	err        error
} {
	return map[string]struct {
		newContext func() (context.Context, context.CancelFunc)
		err        error
	}{
		"cancelled": {func() (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(context.Background())
			time.AfterFunc(50*time.Millisecond, cancel)
			return ctx, cancel
		}, context.Canceled},
		"expired": {func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), 50*time.Millisecond)
		}, context.DeadlineExceeded},
	}
}

func TestAWriteReturnsItsLastResponseWhenTheContextEndsDuringARetryWait(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		status int
		body   string
	}{
		"429":        {http.StatusTooManyRequests, ``},
		"key in use": {http.StatusConflict, `{"errors":[{"code":"IDEMPOTENCY_KEY_IN_USE","message":"in use"}]}`},
	}
	for name, test := range tests {
		for contextName, ending := range endingContexts() {
			t.Run(name+" "+contextName, func(t *testing.T) {
				t.Parallel()
				server, keys := respondingServer(t, test.status, test.body)
				ctx, cancel := ending.newContext()
				defer cancel()

				start := time.Now()
				response, err := shared.ApiClient("aSecretKey", server.URL).R().SetContext(ctx).Post("/events")
				e := shared.AssertOkWithoutResult(response, err)

				require.Less(t, time.Since(start), 300*time.Millisecond)
				require.Len(t, keys.all(), 1)
				require.NoError(t, err)
				require.Equal(t, test.status, response.StatusCode)
				require.Error(t, e)
				require.NotErrorIs(t, e, ending.err)
			})
		}
	}
}

func TestAGetReturnsTheContextErrorWhenTheContextEndsDuringARetryWait(t *testing.T) {
	t.Parallel()
	for contextName, ending := range endingContexts() {
		t.Run(contextName, func(t *testing.T) {
			t.Parallel()
			server, keys := respondingServer(t, http.StatusTooManyRequests, ``)
			ctx, cancel := ending.newContext()
			defer cancel()

			start := time.Now()
			response, err := shared.ApiClient("aSecretKey", server.URL).R().SetContext(ctx).Get("/events")
			e := shared.AssertOkWithoutResult(response, err)

			require.Less(t, time.Since(start), 300*time.Millisecond)
			require.Len(t, keys.all(), 1)
			require.ErrorIs(t, e, ending.err)
		})
	}
}
