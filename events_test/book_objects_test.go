package events

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/events"
	"github.com/seatsio/seatsio-go/v6/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBook(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.Book(event.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, objects.Objects["A-1"].Status)
	require.Equal(t, events.BOOKED, objects.Objects["A-2"].Status)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "A-1", "A-2", "A-3")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
	require.Equal(t, events.BOOKED, info["A-2"].Status)
	require.Equal(t, events.FREE, info["A-3"].Status)
}

func TestBookSections(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithSections(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.Book(event.Key, "Section A-A-1", "Section A-A-2")
	require.NoError(t, err)
	keys := make([]string, len(objects.Objects))
	i := 0
	for k := range objects.Objects {
		keys[i] = k
		i++
	}
	require.Subset(t, keys, []string{"Section A-A-1", "Section A-A-2"})
	require.Equal(t, "Entrance 1", objects.Objects["Section A-A-1"].Entrance)
	require.Equal(t, "Section A", objects.Objects["Section A-A-1"].Section)
	require.Equal(t, events.Labels{
		Own:     events.LabelAndType{Label: "1", Type: "seat"},
		Parent:  events.LabelAndType{Label: "A", Type: "row"},
		Section: "Section A",
	}, objects.Objects["Section A-A-1"].Labels)
	require.Equal(t, events.IDs{Own: "1", Parent: "A", Section: "Section A"}, objects.Objects["Section A-A-1"].IDs)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "Section A-A-1", "Section A-A-2", "Section A-A-3")
	require.Equal(t, events.BOOKED, info["Section A-A-1"].Status)
	require.Equal(t, events.BOOKED, info["Section A-A-2"].Status)
	require.Equal(t, events.FREE, info["Section A-A-3"].Status)
}

func TestHoldTokens(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)
	_, _ = client.Events.Hold(event.Key, []string{"A-1", "A-2"}, &holdToken.HoldToken)

	_, err = client.Events.BookWithHoldToken(event.Key, []string{"A-1", "A-2"}, &holdToken.HoldToken)
	require.NoError(t, err)

	objects, err := client.Events.RetrieveObjectInfo(event.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, objects["A-1"].Status)
	require.Empty(t, objects["A-1"].HoldToken)
	require.Equal(t, events.BOOKED, objects["A-2"].Status)
	require.Empty(t, objects["A-1"].HoldToken)
}

func TestOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
			OrderId: "order1",
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "A-1", "A-2")
	require.Equal(t, "order1", info["A-1"].OrderId)
	require.Equal(t, "order1", info["A-2"].OrderId)
}

func TestKeepExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, map[string]events.ExtraData{
		"A-1": {"foo": "bar"},
	})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:       []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData: true,
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "A-1")
	require.Equal(t, events.ExtraData{"foo": "bar"}, info["A-1"].ExtraData)
}

func TestChannelKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channelKey1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:     []events.ObjectProperties{{ObjectId: "A-1"}},
			ChannelKeys: []string{"channelKey1"},
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "A-1")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
}

func TestIgnoreChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channelKey1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:        []events.ObjectProperties{{ObjectId: "A-1"}},
			IgnoreChannels: true,
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(event.Key, "A-1")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
}

func TestBookBestAvailable(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.BookBestAvailable(event.Key, events.BestAvailableParams{
		Number:     3,
		Categories: []events.CategoryKey{{Key: "cat2"}},
	})
	require.NoError(t, err)

	require.Equal(t, []string{"C-4", "C-5", "C-6"}, bestAvailableResult.Objects)
}
