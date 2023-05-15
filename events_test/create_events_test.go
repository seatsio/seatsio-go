package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateMultipleEventsWithDefaultValues(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	result, err := client.Events.CreateMultiple(chartKey, []events.MultipleEventCreationParams{
		{},
		{},
	})
	require.NoError(t, err)

	require.Equal(t, chartKey, result.Events[0].ChartKey)
	require.Equal(t, chartKey, result.Events[1].ChartKey)
}

func TestCreateMultipleEventsWithEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	result, err := client.Events.CreateMultiple(chartKey, []events.MultipleEventCreationParams{
		{EventKey: "event1"},
		{EventKey: "event2"},
	})
	require.NoError(t, err)

	require.Equal(t, "event1", result.Events[0].Key)
	require.Equal(t, "event2", result.Events[1].Key)
}

func TestCreateMultipleEventsWithTableBookingConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]string{
		"T1": "BY_TABLE", "T2": "BY_SEAT",
	}}
	result, err := client.Events.CreateMultiple(chartKey, []events.MultipleEventCreationParams{
		{TableBookingConfig: &tableBookingConfig},
		{TableBookingConfig: &tableBookingConfig},
	})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, result.Events[0].TableBookingConfig)
	require.Equal(t, tableBookingConfig, result.Events[1].TableBookingConfig)
}

func TestCreateMultipleEventsWithObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	result, err := client.Events.CreateMultiple(chartKey, []events.MultipleEventCreationParams{
		{ObjectCategories: &objectCategories},
		{ObjectCategories: &objectCategories},
	})
	require.NoError(t, err)

	require.Equal(t, objectCategories, result.Events[0].ObjectCategories)
	require.Equal(t, objectCategories, result.Events[1].ObjectCategories)
}

func TestCreateMultipleEventsWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	result, err := client.Events.CreateMultiple(chartKey, []events.MultipleEventCreationParams{
		{Categories: categories},
		{Categories: categories},
	})
	require.NoError(t, err)

	require.Contains(t, result.Events[0].Categories, category)
	require.Contains(t, result.Events[1].Categories, category)
}
