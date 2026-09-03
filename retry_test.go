package seatsio

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/seatsio/seatsio-go/v12/shared"
	"github.com/stretchr/testify/require"
)

func serverDroppingFirstRequest(t *testing.T, attempts *atomic.Int32) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				_ = conn.Close()
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestRetriesGetWhenConnectionFails(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	server := serverDroppingFirstRequest(t, &attempts)

	response, err := shared.ApiClient("aSecretKey", server.URL).R().Get("/events/anEvent/objects")

	require.NoError(t, err)
	require.Equal(t, 204, response.StatusCode)
	require.Equal(t, int32(2), attempts.Load())
}

func TestDoesNotRetryPostWhenConnectionFails(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	server := serverDroppingFirstRequest(t, &attempts)

	_, err := shared.ApiClient("aSecretKey", server.URL).R().Post("/events")

	require.Error(t, err)
	require.Equal(t, int32(1), attempts.Load())
}

func TestRetriesPostOn429(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	response, err := shared.ApiClient("aSecretKey", server.URL).
		SetCommonRetryFixedInterval(10*time.Millisecond).
		R().
		Post("/events")

	require.NoError(t, err)
	require.Equal(t, 204, response.StatusCode)
	require.Equal(t, int32(2), attempts.Load())
}

type unparseableEvent struct {
	Key string `json:"key"`
}

func TestDoesNotRetryWhenResponseBodyCannotBeParsed(t *testing.T) {
	t.Parallel()
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"truncated"))
	}))
	defer server.Close()

	var event unparseableEvent
	_, err := shared.ApiClient("aSecretKey", server.URL).
		SetCommonRetryFixedInterval(10*time.Millisecond).
		R().
		SetSuccessResult(&event).
		Get("/events/anEvent")

	require.Error(t, err)
	require.Equal(t, int32(1), attempts.Load())
}

