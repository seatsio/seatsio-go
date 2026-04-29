package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestUpdateName(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Name: "hurdy"}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "hurdy", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{}}}, postUpdateEvent.Channels)
}

func TestUpdateColor(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Color: "#1E1E1E"}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "bar", Color: "#1E1E1E", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{}}}, postUpdateEvent.Channels)
}

func TestUpdateObjects(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	updateParams := events.UpdateChannelParams{Objects: []string{"A-3", "A-4"}}
	err := client.Channels.Update(test_util.RequestContext(), event.Key, "channelKey1", updateParams)
	require.NoError(t, err)

	postUpdateEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, []events.Channel{{Key: "channelKey1", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-3", "A-4"}, AreaPlaces: map[string]int{}}}, postUpdateEvent.Channels)
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
