package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAbortsEventuallyIfServerKeepsReturning429(t *testing.T) {
	start := time.Now()

	response, _ := events.ApiClient(test_util.SecretKey, "https://httpbin.seatsio.net").
		R().
		Get("/status/429")

	elapsed := time.Now().Sub(start)
	assert.Greater(t, int(elapsed.Seconds()), 10)
	assert.Less(t, int(elapsed.Seconds()), 25)
	assert.Equal(t, response.StatusCode, 429)
}

func TestAbortsDirectlyIfServerReturnsOtherErrorThan429(t *testing.T) {
	start := time.Now()

	response, _ := events.ApiClient(test_util.SecretKey, "https://httpbin.seatsio.net").
		R().
		Get("/status/400")

	elapsed := time.Now().Sub(start)
	assert.Less(t, int(elapsed.Seconds()), 2)
	assert.Equal(t, response.StatusCode, 400)
}

func TestReturnsSuccessfullyWhenServerSends429FirstAndThenSuccess(t *testing.T) {
	for i := 0; i < 20; i++ {
		response, _ := events.ApiClient(test_util.SecretKey, "https://httpbin.seatsio.net").
			R().
			Get("/status/429:0.25,204:0.75")

		assert.Equal(t, response.StatusCode, 204)
	}
}
