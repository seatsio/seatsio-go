package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateEventWithChartKey(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	require.NotZero(t, event.Id)
	require.NotNil(t, event.Key)
	require.Equal(t, chartKey, event.ChartKey)
	require.Equal(t, events.TableBookingConfig{Mode: events.INHERIT}, event.TableBookingConfig)
	require.True(t, event.SupportsBestAvailable)
	require.Nil(t, event.ForSaleConfig)
	require.True(t, event.CreatedOn.After(start))
	require.Nil(t, event.UpdatedOn)
	require.Equal(t, []events.Category{
		{Key: events.CategoryKey{Key: 9}, Label: "Cat1", Color: "#87A9CD", Accessible: false},
		{Key: events.CategoryKey{Key: 10}, Label: "Cat2", Color: "#5E42ED", Accessible: false},
		{Key: events.CategoryKey{Key: "string11"}, Label: "Cat3", Color: "#5E42BB", Accessible: false},
	}, event.Categories)
}

func TestCreateEventWithEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventKey: "anEvent"})
	require.NoError(t, err)

	require.Equal(t, "anEvent", event.Key)
}

func TestCreateEventWithTableBookingConfigCustom(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]events.TableBookingMode{
		"T1": events.BY_TABLE, "T2": events.BY_SEAT,
	}}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, TableBookingConfig: &tableBookingConfig})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestCreateEventWithTableBookingConfigInherit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	tableBookingConfig := events.TableBookingSupport.Inherit()
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, TableBookingConfig: &tableBookingConfig})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestCreateEventWithObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, ObjectCategories: &objectCategories})
	require.NoError(t, err)

	require.Equal(t, objectCategories, event.ObjectCategories)
}

func TestCreateEventWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, Categories: &categories})
	require.NoError(t, err)

	require.Contains(t, event.Categories, category)
}
