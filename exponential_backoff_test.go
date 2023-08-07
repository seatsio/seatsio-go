package seatsio

import (
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	_, _ = shared.ApiClient("aSecretKey", "https://httpbin.seatsio.net").R().Get("/status/200")
	m.Run()
}

func TestAbortsEventuallyIfServerKeepsReturning429(t *testing.T) {
	t.Parallel()
	start := time.Now()

	response, _ := shared.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/429")

	elapsed := time.Now().Sub(start)
	require.Greater(t, int(elapsed.Seconds()), 10)
	require.Less(t, int(elapsed.Seconds()), 25)
	require.Equal(t, 429, response.StatusCode)
}

func TestAbortsDirectlyIfServerReturnsOtherErrorThan429(t *testing.T) {
	t.Parallel()
	start := time.Now()

	response, _ := shared.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/400")

	elapsed := time.Now().Sub(start)
	require.Less(t, int(elapsed.Seconds()), 2)
	require.Equal(t, 400, response.StatusCode)
}

func TestReturnsSuccessfullyWhenServerSends429FirstAndThenSuccess(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		response, _ := shared.ApiClient("aSecretKey", "https://httpbin.seatsio.net").
			R().
			Get("/status/429:0.25,204:0.75")

		require.Equal(t, 204, response.StatusCode)
	}
}

func TestMaxRecountMustNotBeNegative(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	err := client.SetMaxRetries(-1)
	require.Equal(t, "retry count must not be negative", err.Error())
}
