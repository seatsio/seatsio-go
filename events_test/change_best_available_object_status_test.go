package events

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNumber(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
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

func TestCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 3, Categories: []events.CategoryKey{{Key: "cat2"}}},
	})
	require.NoError(t, err)

	require.Equal(t, []string{"C-4", "C-5", "C-6"}, bestAvailableResult.Objects)
}

func TestExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 2, ExtraData: []events.ExtraData{{"foo": "bar"}, {"foo": "baz"}}},
	})
	require.NoError(t, err)

	require.Equal(t, map[string]string{"foo": "bar"}, bestAvailableResult.ObjectDetails["A-4"].ExtraData)
	require.Equal(t, map[string]string{"foo": "baz"}, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestTicketTypes(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 2, TicketTypes: []string{"adult", "child"}},
	})
	require.NoError(t, err)

	require.Equal(t, "adult", bestAvailableResult.ObjectDetails["A-4"].TicketType)
	require.Equal(t, "child", bestAvailableResult.ObjectDetails["A-5"].TicketType)
}

func TestKeepExtraData(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, "A-5", map[string]string{"foo": "bar"})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 1},
		KeepExtraData: true,
	})
	require.NoError(t, err)

	require.Equal(t, map[string]string{"foo": "bar"}, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestKeepExtraDataFalse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, "A-5", map[string]string{"foo": "bar"})
	require.NoError(t, err)

	bestAvailableResult, err := client.Events.ChangeBestAvailableObjectStatus(event.Key, &events.BestAvailableStatusChangeParams{
		Status:        events.ObjectStatusBooked,
		BestAvailable: events.BestAvailableParams{Number: 1},
		KeepExtraData: false,
	})
	require.NoError(t, err)

	require.Nil(t, bestAvailableResult.ObjectDetails["A-5"].ExtraData)
}

func TestOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
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

func TestHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
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

func TestChannelsKeys(t *testing.T) {
	// TODO
}

func TestIgnoreChannels(t *testing.T) {
	// TODO
}
