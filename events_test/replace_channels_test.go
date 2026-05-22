package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestReplaceChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
		Channels: &[]events.CreateChannelParams{
			{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}},
			{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}, AreaPlaces: map[string]int{"GA1": 3}},
		},
	}})
	require.NoError(t, err)

	replaceError := client.Channels.Replace(test_util.RequestContext(), event.Key,
		events.CreateChannelParams{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-5", "A-6"}, AreaPlaces: map[string]int{"GA1": 7}},
		events.CreateChannelParams{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{"A-7", "A-8"}},
	)
	require.NoError(t, replaceError)

	postReplacementEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Len(t, postReplacementEvent.Channels, 2)
	ch1 := postReplacementEvent.Channels[0]
	require.Equal(t, "aaa", ch1.Key)
	require.Equal(t, "bbb", ch1.Name)
	require.Equal(t, "#101010", ch1.Color)
	require.Equal(t, 1, ch1.Index)
	require.Equal(t, []string{"A-5", "A-6"}, ch1.Objects)
	require.Equal(t, map[string]int{"GA1": 7}, ch1.AreaPlaces)
	ch2 := postReplacementEvent.Channels[1]
	require.Equal(t, "ccc", ch2.Key)
	require.Equal(t, "ddd", ch2.Name)
	require.Equal(t, "#F2F2F2", ch2.Color)
	require.Equal(t, 2, ch2.Index)
	require.Equal(t, []string{"A-7", "A-8"}, ch2.Objects)
	require.Empty(t, ch2.AreaPlaces)
}
