package events

import (
	"github.com/seatsio/seatsio-go/v8/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObjectsAreAddedToChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	_ = client.Channels.AddObjects(event.Key, "foo", []string{"A-3", "A-4"})

	retrievedEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2", "A-3", "A-4"}}}, retrievedEvent.Channels)
}
