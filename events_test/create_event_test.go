package events_test

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/events"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateEventWithChartKey(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	require.NoError(t, err)

	require.Equal(t, "anEvent", event.Key)
}

func TestCreateEventWithTableBookingConfigCustom(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	tableBookingConfig := events.TableBookingConfig{Mode: "CUSTOM", Tables: map[string]events.TableBookingMode{
		"T1": events.BY_TABLE, "T2": events.BY_SEAT,
	}}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		TableBookingConfig: &tableBookingConfig,
	}})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestCreateEventWithTableBookingConfigInherit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	tableBookingConfig := events.TableBookingSupport.Inherit()
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		TableBookingConfig: &tableBookingConfig,
	}})
	require.NoError(t, err)

	require.Equal(t, tableBookingConfig, event.TableBookingConfig)
}

func TestCreateEventWithObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		ObjectCategories: &objectCategories,
	}})
	require.NoError(t, err)

	require.Equal(t, objectCategories, event.ObjectCategories)
}

func TestCreateEventWithCategories(t *testing.T) {
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

	require.Contains(t, event.Categories, category)
}

func TestCreateEventWithName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Name: "My event",
	}})
	require.NoError(t, err)

	require.Equal(t, "My event", event.Name)
}

func TestCreateEventWithDate(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	now, _ := time.Parse(time.DateOnly, "2023-07-18")
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Date: events.DateFormat(&now),
	}})
	require.NoError(t, err)

	require.Equal(t, "2023-07-18", event.Date)
}

func TestCreateEventWithChannels(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	channels := []events.CreateChannelParams{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &channels,
	}})
	require.NoError(t, err)

	expectedChannels := []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}
	require.Equal(t, expectedChannels, event.Channels)
}

func TestCreateEventWithForSaleConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	forSaleConfig := &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{"GA1": 5},
		Categories: []string{"Cat1"},
	}

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, ForSaleConfig: forSaleConfig})
	require.NoError(t, err)

	require.Equal(t, forSaleConfig, event.ForSaleConfig)
}
