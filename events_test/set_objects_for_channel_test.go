package events

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetObjects(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventKey: "anEvent"})
	require.Equal(t, 0, len(event.Channels))

	_ = client.Channels.Replace(event.Key, []events.Channel{
		{Key: "channelKey1", Name: "channel 1", Color: "#101010", Index: 1, Objects: []string{}},
		{Key: "channelKey2", Name: "channel 2", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	})

	_ = client.Channels.SetObjects(event.Key, map[string][]string{
		"channelKey1": {"A-1", "A-2"},
		"channelKey2": {"A-3"},
	})

	retrievedEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{
		{Key: "channelKey1", Name: "channel 1", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "channelKey2", Name: "channel 2", Color: "#F2F2F2", Index: 2, Objects: []string{"A-3"}},
	}, retrievedEvent.Channels)
}
