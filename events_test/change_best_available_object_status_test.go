package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeBestAvailableObjectStatusWithNumber(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 3},
	})
	require.NoError(t, err)

	require.True(t, bestAvailableResult.NextToEachOther)
	require.Equal(t, []string{"A-4", "A-5", "A-6"}, bestAvailableResult.Objects)
	require.Equal(t, events.ObjectStatusBooked, bestAvailableResult.ObjectDetails["A-4"].Status)
	require.Equal(t, events.ObjectStatusBooked, bestAvailableResult.ObjectDetails["A-5"].Status)
	require.Equal(t, events.ObjectStatusBooked, bestAvailableResult.ObjectDetails["A-6"].Status)
}

func TestChangeBestAvailableObjectStatusWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 3, Categories: []events.CategoryKey{{Key: "cat2"}}},
	})
	require.NoError(t, err)

	require.Equal(t, []string{"C-4", "C-5", "C-6"}, bestAvailableResult.Objects)
}

func TestChangeBestAvailableObjectStatusWithExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 2, ExtraData: []events.ExtraData{{"foo": "bar"}, {"foo": "baz"}}},
	})
	require.NoError(t, err)

	require.Equal(t, events.ExtraData{"foo": "bar"}, bestAvailableResult.ObjectDetails["A-4"].ExtraData)
	require.Equal(t, events.ExtraData{"foo": "baz"}, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestChangeBestAvailableObjectStatusWithTicketTypes(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 2, TicketTypes: []string{"adult", "child"}},
	})
	require.NoError(t, err)

	require.Equal(t, "adult", bestAvailableResult.ObjectDetails["A-4"].TicketType)
	require.Equal(t, "child", bestAvailableResult.ObjectDetails["A-5"].TicketType)
}

func TestChangeBestAvailableObjectStatusWithKeepExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraDatas(event.Key, map[string]events.ExtraData{"A-5": {"foo": "bar"}})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 1},
		KeepExtraData: true,
	})
	require.NoError(t, err)

	require.Equal(t, events.ExtraData{"foo": "bar"}, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestChangeBestAvailableObjectStatusWithKeepExtraDataFalse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraDatas(event.Key, map[string]events.ExtraData{"A-5": {"foo": "bar"}})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 1},
		KeepExtraData: false,
	})
	require.NoError(t, err)

	require.Nil(t, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestChangeBestAvailableObjectStatusWithOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 2},
		OrderId:       "anOrder",
	})
	require.NoError(t, err)

	require.True(t, bestAvailableResult.NextToEachOther)
	require.Equal(t, []string{"A-4", "A-5"}, bestAvailableResult.Objects)
	require.Equal(t, "anOrder", bestAvailableResult.ObjectDetails["A-4"].OrderId)
	require.Equal(t, "anOrder", bestAvailableResult.ObjectDetails["A-5"].OrderId)
}

func TestChangeBestAvailableObjectStatusWithHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusHeld,
		BestAvailable: events.BestAvailableParams{Number: 2},
		HoldToken:     holdToken.HoldToken,
	})
	require.NoError(t, err)

	require.True(t, bestAvailableResult.NextToEachOther)
	require.Equal(t, []string{"A-4", "A-5"}, bestAvailableResult.Objects)
	require.Equal(t, holdToken.HoldToken, bestAvailableResult.ObjectDetails["A-4"].HoldToken)
	require.Equal(t, holdToken.HoldToken, bestAvailableResult.ObjectDetails["A-5"].HoldToken)
}

func TestChangeBestAvailableObjectStatusWithChannelsKeys(t *testing.T) {
	// TODO
}

func TestChangeBestAvailableObjectStatusWithIgnoreChannels(t *testing.T) {
	// TODO
}
