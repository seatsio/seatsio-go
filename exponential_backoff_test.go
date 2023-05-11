package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAbortsEventuallyIfServerKeepsReturning429(t *testing.T) {
	start := time.Now()

	response, _ := events.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/429")

	elapsed := time.Now().Sub(start)
	require.Greater(t, int(elapsed.Seconds()), 10)
	require.Less(t, int(elapsed.Seconds()), 25)
	require.Equal(t, response.StatusCode, 429)
}

func TestAbortsDirectlyIfServerReturnsOtherErrorThan429(t *testing.T) {
	start := time.Now()

	response, _ := events.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/400")

	elapsed := time.Now().Sub(start)
	require.Less(t, int(elapsed.Seconds()), 2)
	require.Equal(t, response.StatusCode, 400)
}

func TestReturnsSuccessfullyWhenServerSends429FirstAndThenSuccess(t *testing.T) {
	for i := 0; i < 20; i++ {
		response, _ := events.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
			R().
			Get("/status/429:0.25,204:0.75")

		require.Equal(t, response.StatusCode, 204)
	}
}
