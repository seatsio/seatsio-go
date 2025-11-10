package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
)

func TestReplaceForSaleConfigAreaPlacesAndCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.ReplaceForSaleConfig(test_util.RequestContext(), event.Key, false, &events.ForSaleConfigParams{
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	}, retrievedEvent.ForSaleConfig)
}

func TestReplaceForSaleConfigObjects(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.ReplaceForSaleConfig(test_util.RequestContext(), event.Key, false, &events.ForSaleConfigParams{
		Objects: []string{"o1", "o2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}

func TestReplaceForSaleConfigAreaPlaces(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.ReplaceForSaleConfig(test_util.RequestContext(), event.Key, false, &events.ForSaleConfigParams{
		AreaPlaces: map[string]int{"GA1": 3},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}

func TestReplaceForSaleConfigCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.ReplaceForSaleConfig(test_util.RequestContext(), event.Key, false, &events.ForSaleConfigParams{
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{},
		AreaPlaces: map[string]int{},
		Categories: []string{"cat1", "cat2"},
	}, retrievedEvent.ForSaleConfig)
}

func TestReplaceForSaleConfigForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.ReplaceForSaleConfig(test_util.RequestContext(), event.Key, true, &events.ForSaleConfigParams{
		Objects: []string{"o1", "o2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    true,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}
