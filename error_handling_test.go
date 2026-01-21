package seatsio

import (
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/shared"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test300(t *testing.T) {
	t.Parallel()
	var event *events.Event
	response, err := shared.ApiClient("someSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/300")

	_, e := shared.AssertOk(response, err, &event)

	require.EqualError(t, e, "server returned error 300. Body: ")
}

func Test301(t *testing.T) {
	t.Parallel()
	var event *events.Event
	response, err := shared.ApiClient("someSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/301")

	_, e := shared.AssertOk(response, err, &event)

	require.NoError(t, e)
}

func Test400(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: "foo"})

	require.ErrorContains(t, err, "Chart not found: foo was not found in workspace")
}

func Test500(t *testing.T) {
	t.Parallel()
	var event *events.Event
	response, err := shared.ApiClient("someSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/500")

	_, e := shared.AssertOk(response, err, &event)

	require.EqualError(t, e, "server returned error 500. Body: ")
}

func TestWeirdError(t *testing.T) {
	t.Parallel()
	client := NewSeatsioClient("unknownProtocol://", "someSecretKey")

	_, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: "foo"})

	require.EqualError(t, err, "Post \"unknownprotocol:/events\": unsupported protocol scheme \"unknownprotocol\"")
}

func TestTimeout(t *testing.T) {
	t.Parallel()
	var event *events.Event
	response, err := shared.ApiClient("someSecretKey", "https://httpbin.seatsio.net").
		SetTimeout(10 * time.Millisecond).
		R().
		Get("/delay/5")

	_, e := shared.AssertOk(response, err, &event)

	require.ErrorContains(t, e, "context deadline exceeded")
}
