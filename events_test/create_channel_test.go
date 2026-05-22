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
	require.Equal(t, events.Channel{
		Id:         ch.Id,
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      1,
		Objects:    []string{"A-1", "A-2"},
		AreaPlaces: map[string]int{"GA1": 5},
	}, ch)
	require.NotEmpty(t, ch.Id)
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
	require.Equal(t, events.Channel{
		Id:         ch1.Id,
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      1,
		Objects:    []string{"A-1", "A-2"},
		AreaPlaces: map[string]int{"GA1": 5},
	}, ch1)
	ch2 := retrievedEvent.Channels[1]
	require.Equal(t, events.Channel{
		Id:         ch2.Id,
		Key:        "hurdy",
		Name:       "gurdy",
		Color:      "#DFDFDF",
		Index:      2,
		Objects:    []string{"A-3", "A-4"},
		AreaPlaces: map[string]int{},
	}, ch2)
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
	require.NoError(t, err)
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, events.Channel{
		Id:         ch.Id,
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      0,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{},
	}, ch)
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
	require.NoError(t, err)
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, events.Channel{
		Id:         ch.Id,
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      1,
		Objects:    []string{},
		AreaPlaces: map[string]int{},
	}, ch)
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
	require.Len(t, retrievedEvent.Channels, 1)
	ch := retrievedEvent.Channels[0]
	require.Equal(t, events.Channel{
		Id:         ch.Id,
		Key:        "foo",
		Name:       "bar",
		Color:      "#ED303D",
		Index:      1,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{"GA1": 5},
	}, ch)
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
