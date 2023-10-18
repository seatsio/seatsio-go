package events_test

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/shared"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateMultipleEventsWithDefaultValues(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{},
		events.EventParams{},
	)
	require.NoError(t, err)

	require.Equal(t, chartKey, result.Events[0].ChartKey)
	require.Equal(t, chartKey, result.Events[1].ChartKey)
}

func TestCreateMultipleEventsWithEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{EventKey: "event1"},
		events.EventParams{EventKey: "event2"},
	)
	require.NoError(t, err)

	require.Equal(t, "event1", result.Events[0].Key)
	require.Equal(t, "event2", result.Events[1].Key)
}

func TestCreateMultipleEventsWithTableBookingConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]events.TableBookingMode{
		"T1": events.BY_TABLE, "T2": events.BY_SEAT,
	}}
	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{TableBookingConfig: &tableBookingConfig},
		events.EventParams{TableBookingConfig: &tableBookingConfig},
	)
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, result.Events[0].TableBookingConfig)
	require.Equal(t, tableBookingConfig, result.Events[1].TableBookingConfig)
}

func TestCreateMultipleEventsWithObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{ObjectCategories: &objectCategories},
		events.EventParams{ObjectCategories: &objectCategories},
	)
	require.NoError(t, err)

	require.Equal(t, objectCategories, result.Events[0].ObjectCategories)
	require.Equal(t, objectCategories, result.Events[1].ObjectCategories)
}

func TestCreateMultipleEventsWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{Categories: &categories},
		events.EventParams{Categories: &categories},
	)
	require.NoError(t, err)

	require.Contains(t, result.Events[0].Categories, category)
	require.Contains(t, result.Events[1].Categories, category)
}

func TestCreateMultipleEventsWithName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{Name: "Event 1"},
		events.EventParams{Name: "Event 2"},
	)
	require.NoError(t, err)

	require.Equal(t, "Event 1", result.Events[0].Name)
	require.Equal(t, "Event 2", result.Events[1].Name)
}

func TestCreateMultipleEventsWithDate(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{Date: "2023-07-18"},
		events.EventParams{Date: "2023-07-19"},
	)
	require.NoError(t, err)

	require.Equal(t, "2023-07-18", result.Events[0].Date)
	require.Equal(t, "2023-07-19", result.Events[1].Date)
}

func TestCreateMultipleEventsWithChannels(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	channels := []events.CreateChannelParams{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}

	result, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{Channels: &channels},
		events.EventParams{Channels: &channels},
	)
	require.NoError(t, err)

	expectedChannels := []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}
	require.Equal(t, expectedChannels, result.Events[0].Channels)
	require.Equal(t, expectedChannels, result.Events[1].Channels)
}

func TestCreateMultipleEventsWithDuplicateKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Events.CreateMultiple(chartKey,
		events.EventParams{EventKey: "event1"},
		events.EventParams{EventKey: "event1"},
	)
	seatsioError := err.(*shared.SeatsioError)
	require.Equal(t, "GENERAL_ERROR", seatsioError.Code)
	require.Equal(t, "Duplicate event keys are not allowed", seatsioError.Message)
}

func TestCreateSingleEventWithDefaultValues(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	start := time.Now()
	result, err := client.Events.CreateMultiple(chartKey, events.EventParams{})
	require.NoError(t, err)
	end := time.Now().Add(time.Second)

	require.NotEqual(t, int64(0), result.Events[0].Id)
	require.NotNil(t, result.Events[0].Key)
	require.Equal(t, chartKey, result.Events[0].ChartKey)
	require.Equal(t, events.TableBookingSupport.Inherit(), result.Events[0].TableBookingConfig)
	require.True(t, result.Events[0].SupportsBestAvailable)
	require.Equal(t, 3, len(result.Events[0].Categories))
	require.Nil(t, result.Events[0].ForSaleConfig)
	require.True(t, result.Events[0].CreatedOn.After(start))
	require.True(t, result.Events[0].CreatedOn.Before(end))
	require.Nil(t, result.Events[0].UpdatedOn)
}
