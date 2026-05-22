package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestUpdateName(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	updateParams := events.UpdateChannelParams{Name: "hurdy"}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Len(t, postUpdateEvent.Channels, 1)
	require.Equal(t, "channelKey1", postUpdateEvent.Channels[0].Key)
	require.Equal(t, "hurdy", postUpdateEvent.Channels[0].Name)
	require.Equal(t, "#ED303D", postUpdateEvent.Channels[0].Color)
	require.Equal(t, []string{"A-1", "A-2"}, postUpdateEvent.Channels[0].Objects)
	require.Equal(t, map[string]int{"GA1": 5}, postUpdateEvent.Channels[0].AreaPlaces)
}

func TestUpdateColor(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	updateParams := events.UpdateChannelParams{Color: "#1E1E1E"}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Len(t, postUpdateEvent.Channels, 1)
	require.Equal(t, "channelKey1", postUpdateEvent.Channels[0].Key)
	require.Equal(t, "bar", postUpdateEvent.Channels[0].Name)
	require.Equal(t, "#1E1E1E", postUpdateEvent.Channels[0].Color)
	require.Equal(t, []string{"A-1", "A-2"}, postUpdateEvent.Channels[0].Objects)
	require.Equal(t, map[string]int{"GA1": 5}, postUpdateEvent.Channels[0].AreaPlaces)
}

func TestUpdateObjects(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	updateParams := events.UpdateChannelParams{Objects: []string{"A-3", "A-4"}}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Len(t, postUpdateEvent.Channels, 1)
	require.Equal(t, "channelKey1", postUpdateEvent.Channels[0].Key)
	require.Equal(t, "bar", postUpdateEvent.Channels[0].Name)
	require.Equal(t, "#ED303D", postUpdateEvent.Channels[0].Color)
	require.Equal(t, []string{"A-3", "A-4"}, postUpdateEvent.Channels[0].Objects)
	require.Equal(t, map[string]int{"GA1": 5}, postUpdateEvent.Channels[0].AreaPlaces)
}

func TestUpdateAreaPlaces(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1"}})
	updateParams := events.UpdateChannelParams{AreaPlaces: map[string]int{"GA1": 3}}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, map[string]int{"GA1": 3}, postUpdateEvent.Channels[0].AreaPlaces)
}
