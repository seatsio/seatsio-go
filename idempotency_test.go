package seatsio

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestDoesNotRetryWritesByDefaultWhenConnectionFails(t *testing.T) {
	t.Parallel()
	for _, method := range writeMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			var attempts atomic.Int32
			server := serverDroppingFirstRequest(t, &attempts)

			_, err := shared.ApiClient("aSecretKey", server.URL).R().Send(method, "/events")

			require.Error(t, err)
			require.Equal(t, int32(1), attempts.Load())
		})
	}
}

func TestRetriesWritesWithSameIdempotencyKeyWhenEnabled(t *testing.T) {
	t.Parallel()
	for _, method := range writeMethods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()
			var attempts atomic.Int32
			var keys keyRecorder
			server := serverDroppingFirstRequestRecordingKeys(t, &attempts, &keys)
			client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryFixedInterval(10 * time.Millisecond)
			shared.SetRetryWritesOnConnectionFailure(client, true)

			_, err := client.R().Send(method, "/events")

			require.NoError(t, err)
			require.Equal(t, int32(2), attempts.Load())
			require.Len(t, keys.all(), 2)
			require.NotEmpty(t, keys.all()[0])
			require.Equal(t, keys.all()[0], keys.all()[1])
		})
	}
}

func TestDoesNotRetryWritesAfterRetriesAreDisabledAgain(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	server := serverDroppingFirstRequest(t, &attempts)
	client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryFixedInterval(10 * time.Millisecond)
	shared.SetRetryWritesOnConnectionFailure(client, true)
	shared.SetRetryWritesOnConnectionFailure(client, false)

	_, err := client.R().Post("/events")

	require.Error(t, err)
	require.Equal(t, int32(1), attempts.Load())
}

func TestRetriesWritesOn429AndKeyInUseWhenWriteRetriesAreEnabled(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		status int
		body   string
	}{
		"429":        {http.StatusTooManyRequests, ``},
		"key in use": {http.StatusConflict, `{"errors":[{"code":"IDEMPOTENCY_KEY_IN_USE","message":"in use"}]}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var attempts atomic.Int32
			var keys keyRecorder
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				keys.record(r)
				if attempts.Add(1) == 1 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(test.status)
					_, _ = w.Write([]byte(test.body))
					return
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			t.Cleanup(server.Close)
			client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryFixedInterval(10 * time.Millisecond)
			shared.SetRetryWritesOnConnectionFailure(client, true)

			response, err := client.R().Post("/events")

			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode)
			require.Len(t, keys.all(), 2)
			require.Equal(t, keys.all()[0], keys.all()[1])
		})
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

func TestKeyFromContextSurvivesRetry(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	var keys keyRecorder
	server := serverDroppingFirstRequestRecordingKeys(t, &attempts, &keys)
	client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryFixedInterval(10 * time.Millisecond)
	shared.SetRetryWritesOnConnectionFailure(client, true)
	key := shared.NewIdempotencyKey()

	_, err := client.R().SetContext(shared.WithIdempotencyKey(context.Background(), key)).Post("/events")

	require.NoError(t, err)
	require.Equal(t, []string{key.String(), key.String()}, keys.all())
}

func TestParseIdempotencyKey(t *testing.T) {
	t.Parallel()
	valid := uuid.NewString()
	key, err := ClientSupport.ParseIdempotencyKey(valid)
	require.NoError(t, err)
	require.Equal(t, valid, key.String())

	for _, invalid := range []string{
		"",
		"not-a-uuid",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8", // v1
		uuid.Nil.String(),
		"6ba7b810-9dad-41d1-00b4-00c04fd430c8", // v4 with the reserved NCS variant
		"6ba7b810-9dad-41d1-c0b4-00c04fd430c8", // v4 with the Microsoft variant
		"6ba7b810-9dad-41d1-e0b4-00c04fd430c8", // v4 with the reserved future variant
	} {
		_, err := ClientSupport.ParseIdempotencyKey(invalid)
		require.ErrorContains(t, err, "is not a v4 UUID", invalid)
	}
}

func TestZeroIdempotencyKeyIsIgnored(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusNoContent, "")

	_, err := shared.ApiClient("aSecretKey", server.URL).
		R().
		SetContext(shared.WithIdempotencyKey(context.Background(), shared.IdempotencyKey{})).
		Post("/events")

	require.NoError(t, err)
	require.Len(t, keys.all(), 1)
	_, parseErr := shared.ParseIdempotencyKey(keys.all()[0])
	require.NoError(t, parseErr)
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
	client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryCount(0)

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

func TestRetriesWithSameKeyWhileKeyIsInUse(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	var keys keyRecorder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys.record(r)
		if attempts.Add(1) == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"errors":[{"code":"IDEMPOTENCY_KEY_IN_USE","message":"in use"}]}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	response, err := shared.ApiClient("aSecretKey", server.URL).
		SetCommonRetryFixedInterval(10 * time.Millisecond).
		R().
		Post("/events")

	require.NoError(t, err)
	require.Equal(t, 204, response.StatusCode)
	require.Len(t, keys.all(), 2)
	require.Equal(t, keys.all()[0], keys.all()[1])
}

func TestDoesNotRetryOtherConflicts(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusConflict, `{"errors":[{"code":"SOMETHING_ELSE","message":"conflict"}]}`)

	_, err := shared.ApiClient("aSecretKey", server.URL).
		SetCommonRetryFixedInterval(10 * time.Millisecond).
		R().
		Post("/events")

	require.NoError(t, err)
	require.Len(t, keys.all(), 1)
}

func TestDoesNotRetryConflictWithoutAnErrorCode(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"no errors":      `{"errors":[]}`,
		"empty object":   `{}`,
		"malformed json": `{"errors":`,
		"empty body":     ``,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			server, keys := respondingServer(t, http.StatusConflict, body)

			_, err := shared.ApiClient("aSecretKey", server.URL).
				SetCommonRetryFixedInterval(10 * time.Millisecond).
				R().
				Post("/events")

			require.NoError(t, err)
			require.Len(t, keys.all(), 1)
		})
	}
}

func TestKeyInUseErrorIsReturnedWhenRetriesRunOut(t *testing.T) {
	t.Parallel()
	server, keys := respondingServer(t, http.StatusConflict, `{"errors":[{"code":"IDEMPOTENCY_KEY_IN_USE","message":"in use"}]}`)

	response, err := shared.ApiClient("aSecretKey", server.URL).
		SetCommonRetryFixedInterval(10 * time.Millisecond).
		SetCommonRetryCount(2).
		R().
		Post("/events")
	e := shared.AssertOkWithoutResult(response, err)

	require.Len(t, keys.all(), 3)
	require.Equal(t, keys.all()[0], keys.all()[1])
	require.Equal(t, keys.all()[0], keys.all()[2])
	var apiErr *shared.SeatsioError
	require.ErrorAs(t, e, &apiErr)
	require.Equal(t, "IDEMPOTENCY_KEY_IN_USE", apiErr.Code)
	require.Equal(t, keys.all()[0], apiErr.IdempotencyKey)
}

func TestErrorResponseWithoutErrorsIsAnUnexpectedResponse(t *testing.T) {
	t.Parallel()
	server, _ := respondingServer(t, http.StatusBadRequest, `{"errors":[]}`)

	response, err := shared.ApiClient("aSecretKey", server.URL).R().Post("/events")
	e := shared.AssertOkWithoutResult(response, err)

	var apiErr *shared.SeatsioError
	require.False(t, errors.As(e, &apiErr))
	require.ErrorContains(t, e, "server returned error 400")
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

			response, err := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryCount(0).R().Post("/events")
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

	response, err := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryCount(0).R().Get("/events/anEvent")
	e := shared.AssertOkWithoutResult(response, err)

	require.Error(t, e)
	require.False(t, ClientSupport.IsOutcomeUnknown(e))
}

// The transport resends a request with a canonical Idempotency-Key header by itself when a reused connection breaks
// before the first response byte, regardless of the client's retry settings.
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
			client := shared.ApiClient("aSecretKey", server.URL).SetCommonRetryCount(0)

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

func TestRetriesWritesWhenTheResponseBodyIsIncomplete(t *testing.T) {
	t.Parallel()
	tests := map[string]func(w http.ResponseWriter){
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
	for name, firstResponse := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var attempts atomic.Int32
			var keys keyRecorder
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				keys.record(r)
				if attempts.Add(1) == 1 {
					firstResponse(w)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"key":"anEvent"}`))
			}))
			t.Cleanup(server.Close)
			client := shared.ApiClient("aSecretKey", server.URL).
				SetTimeout(150 * time.Millisecond).
				SetCommonRetryFixedInterval(10 * time.Millisecond)
			shared.SetRetryWritesOnConnectionFailure(client, true)

			var event unparseableEvent
			_, err := client.R().SetSuccessResult(&event).Post("/events")

			require.NoError(t, err)
			require.Equal(t, "anEvent", event.Key)
			require.Equal(t, int32(2), attempts.Load())
			require.Equal(t, keys.all()[0], keys.all()[1])
		})
	}
}

func TestStopsWaitingForARetryWhenTheContextEnds(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		status int
		body   string
	}{
		"429":        {http.StatusTooManyRequests, ``},
		"key in use": {http.StatusConflict, `{"errors":[{"code":"IDEMPOTENCY_KEY_IN_USE","message":"in use"}]}`},
	}
	contexts := map[string]func() (context.Context, context.CancelFunc){
		"cancelled": func() (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(context.Background())
			time.AfterFunc(50*time.Millisecond, cancel)
			return ctx, cancel
		},
		"expired": func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(context.Background(), 50*time.Millisecond)
		},
	}
	for name, test := range tests {
		for contextName, newContext := range contexts {
			t.Run(name+" "+contextName, func(t *testing.T) {
				t.Parallel()
				server, keys := respondingServer(t, test.status, test.body)
				ctx, cancel := newContext()
				defer cancel()

				start := time.Now()
				response, err := shared.ApiClient("aSecretKey", server.URL).R().SetContext(ctx).Post("/events")
				e := shared.AssertOkWithoutResult(response, err)

				// the first retry waits at least 400ms
				require.Less(t, time.Since(start), 300*time.Millisecond)
				require.Error(t, e)
				require.Len(t, keys.all(), 1)
			})
		}
	}
}
