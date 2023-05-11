package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test400(t *testing.T) {
	client := NewSeatsioClient(test_util.SecretKey, test_util.BaseUrl)

	_, err := client.Events.Create("foo")

	assert.EqualError(t, err, "Chart not found: foo")
}

func Test500(t *testing.T) {
	var event *events.Event
	response, err := events.ApiClient(test_util.SecretKey, "https://httpbin.seatsio.net").
		R().
		Get("/status/500")

	_, e := events.AssertOk(response, err, event)

	assert.EqualError(t, e, "server returned error 500. Body: ")
}

func TestWeirdError(t *testing.T) {
	client := NewSeatsioClient(test_util.SecretKey, "unknownProtocol://")

	_, err := client.Events.Create("foo")

	assert.EqualError(t, err, "Post \"unknownprotocol:/events\": unsupported protocol scheme \"unknownprotocol\"")
}
