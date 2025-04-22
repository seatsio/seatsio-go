package events

import (
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObjectsAreAddedToChannel(t *testing.T) {
	t.Parallel()

	event, client := CreateChannel(t, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	_ = client.Channels.AddObjects(test_util.RequestContext(), event.Key, "foo", []string{"A-3", "A-4"})

	retrievedEvent, _ := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, []events.Channel{{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2", "A-3", "A-4"}}}, retrievedEvent.Channels)
}
