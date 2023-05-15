package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeObjectStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:  events.ObjectStatusBooked,
		Events:  []string{event.Key},
		Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
	})
	require.NoError(t, err)

	require.Equal(t, events.ObjectStatusBooked, objects.Objects["A-1"].Status)
	require.Equal(t, events.ObjectStatusBooked, objects.Objects["A-2"].Status)
}

func TestObjectDetails(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:  "foo",
		Events:  []string{event.Key},
		Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
	})
	require.NoError(t, err)

	require.Len(t, objects.Objects, 1)
	eventObjectInfo := objects.Objects["A-1"]
	require.Equal(t, "foo", eventObjectInfo.Status)
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
}

func TestHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:    events.ObjectStatusHeld,
		Events:    []string{event.Key},
		Objects:   []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		HoldToken: holdToken.HoldToken,
	})
	require.NoError(t, err)

	require.Equal(t, events.ObjectStatusHeld, objects.Objects["A-1"].Status)
	require.Equal(t, holdToken.HoldToken, objects.Objects["A-1"].HoldToken)
}

func TestExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1", ExtraData: map[string]string{"foo": "bar"}},
		},
	})
	require.NoError(t, err)

	require.Equal(t, map[string]string{"foo": "bar"}, objects.Objects["A-1"].ExtraData)
}

func TestKeepExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, "A-1", map[string]string{"foo": "bar"})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1"},
		},
		KeepExtraData: true,
	})
	require.NoError(t, err)

	require.Equal(t, map[string]string{"foo": "bar"}, objects.Objects["A-1"].ExtraData)
}

func TestKeepExtraDataFalse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, "A-1", map[string]string{"foo": "bar"})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1"},
		},
		KeepExtraData: false,
	})
	require.NoError(t, err)

	require.Nil(t, objects.Objects["A-1"].ExtraData)
}

func TestTicketType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1", TicketType: "adult"},
		},
	})
	require.NoError(t, err)

	require.Equal(t, "adult", objects.Objects["A-1"].TicketType)
}

func TestQuantity(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "GA1", Quantity: 5},
		},
	})
	require.NoError(t, err)

	require.Equal(t, 5, objects.Objects["GA1"].NumBooked)
}

func TestChannelsKeys(t *testing.T) {
	// TODO
}

func TestIgnoreChannels(t *testing.T) {
	// TODO
}

func TestAllowedPreviousStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1"},
		},
		AllowedPreviousStatuses: []string{"onlyAllowedPreviousStatus"},
	})

	require.Equal(t, "ILLEGAL_STATUS_CHANGE", err.(*shared.SeatsioError).Code)
}

func TestRejectedPreviousStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status: events.ObjectStatusBooked,
		Events: []string{event.Key},
		Objects: []events.ObjectProperties{
			{ObjectId: "A-1"},
		},
		RejectedPreviousStatuses: []string{"free"},
	})

	require.Equal(t, "ILLEGAL_STATUS_CHANGE", err.(*shared.SeatsioError).Code)
}
