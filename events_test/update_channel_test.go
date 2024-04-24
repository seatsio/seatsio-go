package events

import (
	"github.com/seatsio/seatsio-go/v7/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateName(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Name: "hurdy"}
	err := client.Channels.Update(event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "hurdy", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}}}, postUpdateEvent.Channels)
}

func TestUpdateColor(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Color: "#1E1E1E"}
	err := client.Channels.Update(event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "bar", Color: "#1E1E1E", Index: 1, Objects: []string{"A-1", "A-2"}}}, postUpdateEvent.Channels)
}

func TestUpdateObjects(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Objects: []string{"A-3", "A-4"}}
	err := client.Channels.Update(event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-3", "A-4"}}}, postUpdateEvent.Channels)
}
