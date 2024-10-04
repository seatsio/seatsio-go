package events

import (
	"github.com/seatsio/seatsio-go/v8/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	retrievedEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}}}, retrievedEvent.Channels)

	deleteErr := client.Channels.Delete(event.Key, "foo")
	require.NoError(t, deleteErr)

	postDeleteEvent, _ := client.Events.Retrieve(event.Key)
	require.Empty(t, postDeleteEvent.Channels)
}
