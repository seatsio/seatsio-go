package events

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReplaceChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventKey: "anEvent"})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.CreateMultiple(event.Key, &[]events.CreateChannelParams{
		*&events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}},
		*&events.CreateChannelParams{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{
		{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
	}, retrievedEvent.Channels)

	replaceError := client.Channels.Replace(event.Key, []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-5", "A-6"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{"A-7", "A-8"}},
	})
	require.NoError(t, replaceError)

	postReplacementEvent, _ := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-5", "A-6"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{"A-7", "A-8"}},
	}, postReplacementEvent.Channels)
}
