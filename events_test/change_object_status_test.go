package events_test

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/shared"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeObjectStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1", "A-2"}, events.BOOKED)
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, objects.Objects["A-1"].Status)
	require.Equal(t, events.BOOKED, objects.Objects["A-2"].Status)
}

func TestChangeObjectStatusWithObjectDetails(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, "foo")
	require.NoError(t, err)

	var status = "foo"
	require.Len(t, objects.Objects, 1)
	eventObjectInfo := objects.Objects["A-1"]
	require.Equal(t, status, eventObjectInfo.Status)
	require.Equal(t, "A-1", eventObjectInfo.Label)
	require.Equal(t, events.Labels{Own: events.LabelAndType{Label: "1", Type: "seat"}, Parent: events.LabelAndType{Label: "A", Type: "row"}, Section: ""}, eventObjectInfo.Labels)
	require.Equal(t, events.IDs{Own: "1", Parent: "A", Section: ""}, eventObjectInfo.IDs)
	require.Equal(t, "Cat1", eventObjectInfo.CategoryLabel)
	require.Equal(t, events.CategoryKey{Key: "9"}, eventObjectInfo.CategoryKey)
	require.Empty(t, eventObjectInfo.TicketType)
	require.Empty(t, eventObjectInfo.OrderId)
	require.True(t, eventObjectInfo.ForSale)
	require.Empty(t, eventObjectInfo.Section)
	require.Empty(t, eventObjectInfo.Entrance)
	require.Equal(t, 0, eventObjectInfo.NumBooked)
	require.Empty(t, eventObjectInfo.Capacity)
	require.Equal(t, "seat", eventObjectInfo.ObjectType)
	require.Nil(t, eventObjectInfo.ExtraData)
	require.Empty(t, eventObjectInfo.LeftNeighbour)
	require.Equal(t, "A-2", eventObjectInfo.RightNeighbour)
	require.Empty(t, eventObjectInfo.HoldToken)
	require.Equal(t, 0, eventObjectInfo.SeasonStatusOverriddenQuantity)
}

func TestChangeObjectStatusWithHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:    events.HELD,
			Objects:   []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
			HoldToken: holdToken.HoldToken,
		},
	})
	require.NoError(t, err)

	require.Equal(t, events.HELD, objects.Objects["A-1"].Status)
	require.Equal(t, holdToken.HoldToken, objects.Objects["A-1"].HoldToken)
}

func TestChangeObjectStatusWithExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1", ExtraData: events.ExtraData{"foo": "bar"}},
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, events.ExtraData{"foo": "bar"}, objects.Objects["A-1"].ExtraData)
}

func TestChangeObjectStatusWithKeepExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(test_util.RequestContext(), event.Key, map[string]events.ExtraData{"A-1": {"foo": "bar"}})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
			},
			KeepExtraData: true,
		},
	})
	require.NoError(t, err)

	require.Equal(t, events.ExtraData{"foo": "bar"}, objects.Objects["A-1"].ExtraData)
}

func TestChangeObjectStatusWithKeepExtraDataFalse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(test_util.RequestContext(), event.Key, map[string]events.ExtraData{"A-1": {"foo": "bar"}})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
			},
			KeepExtraData: false,
		},
	})
	require.NoError(t, err)

	require.Nil(t, objects.Objects["A-1"].ExtraData)
}

func TestChangeObjectStatusWithTicketType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1", TicketType: "adult"},
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, "adult", objects.Objects["A-1"].TicketType)
}

func TestChangeObjectStatusWithQuantity(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "GA1", Quantity: 5},
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, 5, objects.Objects["GA1"].NumBooked)
}

func TestChangeObjectStatusWithChannelsKeys(t *testing.T) {
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

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:      events.BOOKED,
			Objects:     []events.ObjectProperties{{ObjectId: "A-1"}},
			ChannelKeys: []string{"channelKey1"},
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
}

func TestChangeObjectStatusWithIgnoreChannels(t *testing.T) {
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

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:         events.BOOKED,
			Objects:        []events.ObjectProperties{{ObjectId: "A-1"}},
			IgnoreChannels: true,
		},
	})
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
}

func TestChangeObjectStatusWithAllowedPreviousStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
			},
			AllowedPreviousStatuses: []string{"onlyAllowedPreviousStatus"},
		},
	})

	require.Equal(t, "ILLEGAL_STATUS_CHANGE", err.(*shared.SeatsioError).Code)
}

func TestChangeObjectStatusWithRejectedPreviousStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: events.BOOKED,
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
			},
			RejectedPreviousStatuses: []string{"free"},
		},
	})

	require.Equal(t, "ILLEGAL_STATUS_CHANGE", err.(*shared.SeatsioError).Code)
}

func TestChangeObjectStatusWithResaleListingId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:          events.RESALE,
			Objects:         []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
			ResaleListingId: "listing1",
		},
	})
	require.NoError(t, err)

	require.Equal(t, "listing1", objects.Objects["A-1"].ResaleListingId)
}
