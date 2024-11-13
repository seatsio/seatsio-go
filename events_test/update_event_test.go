package events_test

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/seasons"
	"github.com/seatsio/seatsio-go/v9/shared"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateEventEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{EventKey: "newKey"}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve("newKey")
	require.NoError(t, err)
	require.Equal(t, "newKey", updatedEvent.Key)
}

func TestUpdateEventName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Name: "foo",
	}})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		Name: "bar",
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, "bar", updatedEvent.Name)
}

func TestUpdateEventDate(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	now, _ := time.Parse(time.DateOnly, "2023-07-18")
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Date: events.DateFormat(&now),
	}})
	require.NoError(t, err)

	updatedDate := "2023-08-03"
	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		Date: updatedDate,
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, updatedDate, updatedEvent.Date)
}

func TestUpdateEventTableBookingConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	tableBookingConfig := events.TableBookingConfig{Mode: events.CUSTOM, Tables: map[string]events.TableBookingMode{
		"T1": events.BY_TABLE, "T2": events.BY_SEAT,
	}}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		TableBookingConfig: &tableBookingConfig,
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, tableBookingConfig, updatedEvent.TableBookingConfig)
}

func TestUpdateEventObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		ObjectCategories: &objectCategories,
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Equal(t, objectCategories, updatedEvent.ObjectCategories)
}

func TestUpdateEventRemoveObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		ObjectCategories: &map[string]events.CategoryKey{
			"A-1": {10},
		},
	}})
	require.NoError(t, err)

	objectCategories := map[string]events.CategoryKey{}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		ObjectCategories: &objectCategories,
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Empty(t, updatedEvent.ObjectCategories)
}

func TestUpdateEventCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		Categories: &categories,
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Contains(t, updatedEvent.Categories, category)
}

func TestUpdateEventRemoveCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Categories: &categories,
	}})
	require.NoError(t, err)

	err = client.Events.Update(event.Key, &events.UpdateEventParams{EventParams: &events.EventParams{
		Categories: &[]events.Category{},
	}})
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.NotContains(t, updatedEvent.Categories, category)
}

func TestUpdateEventIsInThePast(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_, err := client.Seasons.CreateSeasonWithOptions(chartKey, &seasons.CreateSeasonParams{EventKeys: []string{"event1"}})
	require.NoError(t, err)

	err = client.Events.Update("event1", &events.UpdateEventParams{IsInThePast: shared.OptionalBool(true)})
	require.NoError(t, err)
	updatedEvent, err := client.Events.Retrieve("event1")
	require.NoError(t, err)
	require.True(t, updatedEvent.IsInThePast)
}
