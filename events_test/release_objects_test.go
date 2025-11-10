package events

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReleaseObjects(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.Book(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)

	_, err = client.Events.Release(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)
	require.Equal(t, events.FREE, retrieveObjectInfo["A-1"].Status)
	require.Equal(t, events.FREE, retrieveObjectInfo["A-2"].Status)
}

func TestWithHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	require.NoError(t, err)

	_, err = client.Events.Hold(test_util.RequestContext(), event.Key, []string{"A-1", "A-2"}, &holdToken.HoldToken)
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithHoldToken(test_util.RequestContext(), event.Key, []string{"A-1", "A-2"}, &holdToken.HoldToken)
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)
	require.Equal(t, events.FREE, retrieveObjectInfo["A-1"].Status)
	require.Empty(t, retrieveObjectInfo["A-2"].HoldToken)
}

func TestWithOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
			OrderId: "order1",
		},
	})
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)
	require.Equal(t, "order1", retrieveObjectInfo["A-1"].OrderId)
	require.Equal(t, "order1", retrieveObjectInfo["A-2"].OrderId)
}

func TestWithExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(test_util.RequestContext(), event.Key, map[string]events.ExtraData{
		"A-1": {"foo": "bar"},
	})
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:       []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData: true,
		},
	})
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1")
	require.NoError(t, err)
	require.Equal(t, events.ExtraData{"foo": "bar"}, retrieveObjectInfo["A-1"].ExtraData)
}

func TestWithChannelKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channelKey1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:       []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData: true,
			ChannelKeys:   []string{"channelKey1"},
		},
	})
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:       []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData: true,
			ChannelKeys:   []string{"channelKey1"},
		},
	})
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1")
	require.NoError(t, err)
	require.Equal(t, events.FREE, retrieveObjectInfo["A-1"].Status)
}

func TestIgnoreChannelKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channelKey1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:        []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData:  true,
			IgnoreChannels: true,
		},
	})
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects:        []events.ObjectProperties{{ObjectId: "A-1"}},
			KeepExtraData:  true,
			IgnoreChannels: true,
		},
	})
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1")
	require.NoError(t, err)
	require.Equal(t, events.FREE, retrieveObjectInfo["A-1"].Status)
}

func TestBestAvailable(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	require.NoError(t, err)

	_, err = client.Events.HoldBestAvailable(test_util.RequestContext(), event.Key, events.BestAvailableParams{
		Number: 2,
		Categories: []events.CategoryKey{
			{Key: "cat2"},
		},
	}, holdToken.HoldToken)
	require.NoError(t, err)

	_, err = client.Events.ReleaseWithHoldToken(test_util.RequestContext(), event.Key, []string{"C-4", "C-5"}, &holdToken.HoldToken)
	require.NoError(t, err)

	retrieveObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "C-4", "C-5")
	require.NoError(t, err)
	require.Equal(t, events.FREE, retrieveObjectInfo["C-4"].Status)
	require.Empty(t, retrieveObjectInfo["C-5"].HoldToken)
}
