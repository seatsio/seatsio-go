package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChartKeyIsRequired(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	require.NotZero(t, event.Id)
	require.NotNil(t, event.Key)
	require.Equal(t, chartKey, event.ChartKey)
	require.Equal(t, events.TableBookingConfig{Mode: "INHERIT"}, event.TableBookingConfig)
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

func TestEventKeyCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey, EventKey: "anEvent"})
	require.NoError(t, err)

	require.Equal(t, "anEvent", event.Key)
}

func TestTableBookingConfigCustomCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]string{
		"T1": "BY_TABLE", "T2": "BY_SEAT",
	}}
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey, TableBookingConfig: &tableBookingConfig})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestTableBookingConfigInheritCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	tableBookingConfig := events.TableBookingConfig{Mode: "INHERIT"}
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey, TableBookingConfig: &tableBookingConfig})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestObjectCategoriesCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey, ObjectCategories: &objectCategories})
	require.NoError(t, err)

	require.Equal(t, objectCategories, event.ObjectCategories)
}

func TestCategoriesCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey, Categories: categories})
	require.NoError(t, err)

	require.Contains(t, event.Categories, category)
}
