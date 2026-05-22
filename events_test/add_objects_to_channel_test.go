package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestObjectsAreAddedToChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	_ = client.Channels.AddObjects(test_util.RequestContext(), event.Key, "foo", []string{"A-3", "A-4"})

	retrievedEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, "foo", ch.Key)
	require.Equal(t, "bar", ch.Name)
	require.Equal(t, "#ED303D", ch.Color)
	require.Equal(t, 1, ch.Index)
	require.Equal(t, []string{"A-1", "A-2", "A-3", "A-4"}, ch.Objects)
	require.Equal(t, map[string]int{"GA1": 5}, ch.AreaPlaces)
}

func TestAreaPlacesAreAddedToChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1"}})
	_ = client.Channels.AddObjects(test_util.RequestContext(), event.Key, "foo", []string{}, map[string]int{"GA1": 5})

	retrievedEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, map[string]int{"GA1": 5}, retrievedEvent.Channels[0].AreaPlaces)
}
