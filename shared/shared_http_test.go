package shared

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimedOutRequestsAreNotRetried(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		time.Sleep(300 * time.Millisecond)
	}))
	defer server.Close()

	_, err := ApiClient("aSecretKey", server.URL).SetTimeout(50 * time.Millisecond).R().Get("/")

	require.Error(t, err)
	require.Equal(t, int32(1), requests.Load())
}

func TestTooManyRequestsAreStillRetried(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	response, err := ApiClient("aSecretKey", server.URL).
		SetCommonRetryBackoffInterval(time.Millisecond, 5*time.Millisecond).
		R().
		Get("/")

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Equal(t, int32(3), requests.Load())
}
