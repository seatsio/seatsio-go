package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
)

func TestEditForSaleConfigMakeForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	_, err = client.Events.EditForSaleConfig(test_util.RequestContext(), event.Key, []events.ObjectAndQuantity{{Object: "A-1"}, {Object: "A-2"}}, nil)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	expectedForSaleConfig := &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}

	require.Equal(t, expectedForSaleConfig, retrievedEvent.ForSaleConfig)
}

func TestEditForSaleConfigReturnsResponse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	r, err := client.Events.EditForSaleConfig(test_util.RequestContext(), event.Key, []events.ObjectAndQuantity{{Object: "A-1"}, {Object: "A-2"}}, nil)
	require.NoError(t, err)

	expectedForSaleConfig := &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}

	require.Equal(t, expectedForSaleConfig, r.ForSaleConfig)
	require.Equal(t, 9, r.ForSaleRateLimitInfo.RateLimitRemainingCalls)
	require.NotNil(t, r.ForSaleRateLimitInfo.RateLimitResetDate)
}

func TestEditForSaleConfigMakeNotForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.EditForSaleConfig(test_util.RequestContext(), event.Key, nil, []events.ObjectAndQuantity{{Object: "A-1"}, {Object: "A-2"}})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1", "A-2"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}

func TestEditForSaleConfigMakeAreaPlacesNotForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.EditForSaleConfig(test_util.RequestContext(), event.Key, nil, []events.ObjectAndQuantity{{Object: "GA1", Quantity: 5}})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{},
		AreaPlaces: map[string]int{"GA1": 5},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}
