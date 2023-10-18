package events

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(event.Key, &events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}}}, retrievedEvent.Channels)
}

func TestCreateChannels(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(event.Key,
		&events.CreateChannelParams{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}},
		&events.CreateChannelParams{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
	)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.Equal(t, []events.Channel{
		{Key: "foo", Name: "bar", Color: "#ED303D", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "hurdy", Name: "gurdy", Color: "#DFDFDF", Index: 2, Objects: []string{"A-3", "A-4"}},
	}, retrievedEvent.Channels)
}

func TestIndexIsOptional(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(event.Key, &events.CreateChannelParams{
		Key:     "foo",
		Name:    "bar",
		Color:   "#ED303D",
		Objects: []string{"A-1"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.Equal(t, 1, len(retrievedEvent.Channels))
	require.Equal(t, "foo", retrievedEvent.Channels[0].Key)
	require.Equal(t, int32(0), retrievedEvent.Channels[0].Index)
	require.Equal(t, "A-1", retrievedEvent.Channels[0].Objects[0])
}

func TestObjectsIsOptional(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.Equal(t, 0, len(event.Channels))

	err := client.Channels.Create(event.Key, &events.CreateChannelParams{
		Key:   "foo",
		Name:  "bar",
		Color: "#ED303D",
		Index: 1,
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.Equal(t, 1, len(retrievedEvent.Channels))
	require.Equal(t, "foo", retrievedEvent.Channels[0].Key)
	require.Equal(t, int32(1), retrievedEvent.Channels[0].Index)
	require.Empty(t, retrievedEvent.Channels[0].Objects)
}
