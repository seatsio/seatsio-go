package events_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/shared"
	"github.com/seatsio/seatsio-go/v13/test_util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSameIdempotencyKeyReplaysTheResponse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompanyWithIdempotencyEnabled(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	ctx := seatsio.ClientSupport.WithIdempotencyKey(test_util.RequestContext(), seatsio.ClientSupport.NewIdempotencyKey())

	first, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	second, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	require.Equal(t, first.Id, second.Id)
	require.Equal(t, first.Key, second.Key)
}

func TestSameIdempotencyKeyWithADifferentRequestIsRejected(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompanyWithIdempotencyEnabled(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	ctx := seatsio.ClientSupport.WithIdempotencyKey(test_util.RequestContext(), seatsio.ClientSupport.NewIdempotencyKey())
	_, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{EventKey: "anotherEvent"}})

	var seatsioErr *shared.SeatsioError
	require.ErrorAs(t, err, &seatsioErr)
	require.Equal(t, "IDEMPOTENCY_KEY_REUSED", seatsioErr.Code)
}

func TestIdempotencyKeyIsIgnoredWhenNotEnabledForTheCompany(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	ctx := seatsio.ClientSupport.WithIdempotencyKey(test_util.RequestContext(), seatsio.ClientSupport.NewIdempotencyKey())

	first, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	second, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	require.NotEqual(t, first.Id, second.Id)
}

func TestIdempotencyKeyIsAvailableOnAFailedRequest(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompanyWithIdempotencyEnabled(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	key := seatsio.ClientSupport.NewIdempotencyKey()
	ctx := seatsio.ClientSupport.WithIdempotencyKey(test_util.RequestContext(), key)

	_, err := client.Events.Create(ctx, &events.CreateEventParams{ChartKey: "unknownChart"})

	require.Error(t, err)
	actual, ok := seatsio.ClientSupport.IdempotencyKeyOf(err)
	require.True(t, ok)
	require.Equal(t, key, actual)
}

func TestGeneratedIdempotencyKeyIsAvailableOnAFailedRequest(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompanyWithIdempotencyEnabled(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: "unknownChart"})

	_, ok := seatsio.ClientSupport.IdempotencyKeyOf(err)
	require.True(t, ok)
}

func proxyThatForwardsButDropsTheFirstResponse(t *testing.T) string {
	t.Helper()
	var requests atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forwarded, err := http.NewRequest(r.Method, test_util.BaseUrl+r.RequestURI, r.Body)
		if !assert.NoError(t, err) {
			return
		}
		forwarded.Header = r.Header.Clone()
		response, err := http.DefaultClient.Do(forwarded)
		if !assert.NoError(t, err) {
			return
		}
		defer response.Body.Close()
		if requests.Add(1) == 1 {
			conn, _, err := w.(http.Hijacker).Hijack()
			if assert.NoError(t, err) {
				_ = conn.Close()
			}
			return
		}
		for key, values := range response.Header {
			w.Header()[key] = values
		}
		w.WriteHeader(response.StatusCode)
		_, _ = io.Copy(w, response.Body)
	}))
	t.Cleanup(proxy.Close)
	return proxy.URL
}

func TestRetryAfterALostResponseCreatesOneEvent(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompanyWithIdempotencyEnabled(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(proxyThatForwardsButDropsTheFirstResponse(t), company.Admin.SecretKey)
	client.SetRetryWritesOnConnectionFailure(true)
	realClient := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	all, err := realClient.Events.ListAll(test_util.RequestContext())
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, event.Id, all[0].Id)
}
