package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestDeleteChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	retrievedEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, "foo", ch.Key)
	require.Equal(t, "bar", ch.Name)
	require.Equal(t, "#ED303D", ch.Color)
	require.Equal(t, 1, ch.Index)
	require.Equal(t, []string{"A-1", "A-2"}, ch.Objects)
	require.Equal(t, map[string]int{"GA1": 5}, ch.AreaPlaces)

	deleteErr := client.Channels.Delete(test_util.RequestContext(), event.Key, "foo")
	require.NoError(t, deleteErr)

	postDeleteEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Empty(t, postDeleteEvent.Channels)
}
