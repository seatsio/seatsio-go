package events_test

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/events"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarkAsNotForSaleObjectsAndAreaPlacesAndCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(event.Key, &events.ForSaleConfigParams{
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	}, retrievedEvent.ForSaleConfig)
}

func TestMarkAsNotForSaleObjects(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(event.Key, &events.ForSaleConfigParams{
		Objects: []string{"o1", "o2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}

func TestMarkAsNotForSaleAreaPlaces(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(event.Key, &events.ForSaleConfigParams{
		AreaPlaces: map[string]int{"GA1": 3},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}

func TestMarkAsNotForSaleCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(event.Key, &events.ForSaleConfigParams{
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{},
		AreaPlaces: map[string]int{},
		Categories: []string{"cat1", "cat2"},
	}, retrievedEvent.ForSaleConfig)
}

func TestNumNotForSaleIsCorrectlyExposed(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(event.Key, &events.ForSaleConfigParams{
		AreaPlaces: map[string]int{"GA1": 3},
	})
	require.NoError(t, err)

	ga1Info, err := client.Events.RetrieveObjectInfo(event.Key, "GA1")
	require.NoError(t, err)

	require.Equal(t, 3, ga1Info["GA1"].NumNotForSale)
}
