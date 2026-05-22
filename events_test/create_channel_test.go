package events

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(test_util.RequestContext(), event.Key, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, "foo", ch.Key)
	require.Equal(t, "bar", ch.Name)
	require.Equal(t, "#ED303D", ch.Color)
	require.Equal(t, 1, ch.Index)
	require.Equal(t, []string{"A-1", "A-2"}, ch.Objects)
	require.Equal(t, map[string]int{"GA1": 5}, ch.AreaPlaces)
}

func TestCreateChannels(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(test_util.RequestContext(), event.Key,
		&events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}, AreaPlaces: map[string]int{"GA1": 5}},
		&events.CreateChannelParams{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
	)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Len(t, retrievedEvent.Channels, 2)
	ch1 := retrievedEvent.Channels[0]
	require.Equal(t, "foo", ch1.Key)
	require.Equal(t, "bar", ch1.Name)
	require.Equal(t, "#ED303D", ch1.Color)
	require.Equal(t, 1, ch1.Index)
	require.Equal(t, []string{"A-1", "A-2"}, ch1.Objects)
	require.Equal(t, map[string]int{"GA1": 5}, ch1.AreaPlaces)
	ch2 := retrievedEvent.Channels[1]
	require.Equal(t, "hurdy", ch2.Key)
	require.Equal(t, "gurdy", ch2.Name)
	require.Equal(t, "#DFDFDF", ch2.Color)
	require.Equal(t, 2, ch2.Index)
	require.Equal(t, []string{"A-3", "A-4"}, ch2.Objects)
	require.Empty(t, ch2.AreaPlaces)
}

func TestIndexIsOptional(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(test_util.RequestContext(), event.Key, &events.CreateChannelParams{
		Key:     "foo",
		Name:    "bar",
		Color:   "#ED303D",
		Objects: []string{"A-1"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, 1, len(retrievedEvent.Channels))
	require.Equal(t, "foo", retrievedEvent.Channels[0].Key)
	require.Equal(t, 0, retrievedEvent.Channels[0].Index)
	require.Equal(t, "A-1", retrievedEvent.Channels[0].Objects[0])
}

func TestObjectsIsOptional(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(test_util.RequestContext(), event.Key, &events.CreateChannelParams{
		Key:   "foo",
		Name:  "bar",
		Color: "#ED303D",
		Index: 1,
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, 1, len(retrievedEvent.Channels))
	require.Equal(t, "foo", retrievedEvent.Channels[0].Key)
	require.Equal(t, 1, retrievedEvent.Channels[0].Index)
	require.Empty(t, retrievedEvent.Channels[0].Objects)
}

func TestCreateChannelWithAreaPlaces(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})

	err := client.Channels.Create(test_util.RequestContext(), event.Key, &events.CreateChannelParams{
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      1,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{"GA1": 5},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, map[string]int{"GA1": 5}, retrievedEvent.Channels[0].AreaPlaces)
}

func TestChannelHasID(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})

	err := client.Channels.Create(test_util.RequestContext(), event.Key, &events.CreateChannelParams{
		Key:   "foo",
		Name:  "bar",
		Color: "#ED303D",
		Index: 1,
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Len(t, retrievedEvent.Channels, 1)
	require.NotEmpty(t, retrievedEvent.Channels[0].Id)
}

func TestAreaPartitionLabel(t *testing.T) {
	ch := events.Channel{Id: "abc123"}
	require.Equal(t, "myArea##abc123", ch.AreaPartitionLabel("myArea"))
}
