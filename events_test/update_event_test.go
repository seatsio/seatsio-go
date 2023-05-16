package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateEventChartKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		ChartKey: chartKey2,
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, chartKey2, updatedEvent.ChartKey)
}

func TestUpdateEventEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		EventKey: "newKey",
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve("newKey")
	require.NoError(t, err)
	require.Equal(t, "newKey", updatedEvent.Key)
}

func TestUpdateEventTableBookingConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]string{
		"T1": "BY_TABLE", "T2": "BY_SEAT",
	}}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		TableBookingConfig: &tableBookingConfig,
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, tableBookingConfig, updatedEvent.TableBookingConfig)
}

func TestUpdateEventObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		ObjectCategories: &objectCategories,
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, objectCategories, updatedEvent.ObjectCategories)
}

func TestUpdateEventRemoveObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, ObjectCategories: &map[string]events.CategoryKey{
		"A-1": {10},
	}})
	require.NoError(t, err)

	objectCategories := map[string]events.CategoryKey{}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		ObjectCategories: &objectCategories,
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Empty(t, updatedEvent.ObjectCategories)
}

func TestUpdateEventCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		Categories: &categories,
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Contains(t, updatedEvent.Categories, category)
}

func TestUpdateEventRemoveCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	event, err := client.Events.Create(&events.CreateEventParams{
		ChartKey:   chartKey,
		Categories: &categories,
	})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{
		Categories: &[]events.Category{},
	})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.NotContains(t, updatedEvent.Categories, category)
}
