package events

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReplaceChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
		Channels: &[]events.CreateChannelParams{
			{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}},
			{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
		},
	}})
	require.NoError(t, err)

	replaceError := client.Channels.Replace(event.Key,
		events.CreateChannelParams{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-5", "A-6"}},
		events.CreateChannelParams{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{"A-7", "A-8"}},
	)
	require.NoError(t, replaceError)

	postReplacementEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-5", "A-6"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{"A-7", "A-8"}},
	}, postReplacementEvent.Channels)
}
