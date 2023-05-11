package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test400(t *testing.T) {
	company := test_util.CreateTestCompany(t)
	client := NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_, err := client.Events.Create(&events.EventCreationParams{ChartKey: "foo"})

	require.EqualError(t, err, "Chart not found: foo")
}

func Test500(t *testing.T) {
	var event *events.Event
	response, err := events.ApiClient("someSecretKey", "https://httpbin.seatsio.net").
		R().
		Get("/status/500")

	_, e := events.AssertOk(response, err, &event)

	require.EqualError(t, e, "server returned error 500. Body: ")
}

func TestWeirdError(t *testing.T) {
	client := NewSeatsioClient("someSecretKey", "unknownProtocol://")

	_, err := client.Events.Create(&events.EventCreationParams{ChartKey: "foo"})

	require.EqualError(t, err, "Post \"unknownprotocol:/events\": unsupported protocol scheme \"unknownprotocol\"")
}
