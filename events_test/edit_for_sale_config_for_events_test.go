package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestEditForSaleConfigForEventsMakeForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	_, err = client.Events.EditForSaleConfigForEvents(test_util.RequestContext(), map[string]events.EditForSaleConfigRequest{
		event1.Key: {ForSale: []events.ObjectAndQuantity{{Object: "A-1"}}},
		event2.Key: {ForSale: []events.ObjectAndQuantity{{Object: "A-2"}}},
	})
	require.NoError(t, err)

	retrievedEvent1, err := client.Events.Retrieve(test_util.RequestContext(), event1.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-2", "A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent1.ForSaleConfig)

	retrievedEvent2, err := client.Events.Retrieve(test_util.RequestContext(), event2.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1", "A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent2.ForSaleConfig)
}

func TestEditForSaleConfigForEventsReturnsResponse(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey:      chartKey,
		ForSaleConfig: &events.ForSaleConfig{ForSale: false, Objects: []string{"A-1", "A-2", "A-3"}}})
	require.NoError(t, err)

	r, err := client.Events.EditForSaleConfigForEvents(test_util.RequestContext(), map[string]events.EditForSaleConfigRequest{
		event1.Key: {ForSale: []events.ObjectAndQuantity{{Object: "A-1"}}},
		event2.Key: {ForSale: []events.ObjectAndQuantity{{Object: "A-2"}}},
	})
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-2", "A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, r[event1.Key].ForSaleConfig)
	require.Equal(t, 9, r[event1.Key].ForSaleRateLimitInfo.RateLimitRemainingCalls)
	require.NotNil(t, 9, r[event1.Key].ForSaleRateLimitInfo.RateLimitResetDate)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1", "A-3"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, r[event2.Key].ForSaleConfig)
	require.Equal(t, 9, r[event2.Key].ForSaleRateLimitInfo.RateLimitRemainingCalls)
	require.NotNil(t, 9, r[event2.Key].ForSaleRateLimitInfo.RateLimitResetDate)
}

func TestEditForSaleConfigForEventsMakeNotForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.EditForSaleConfigForEvents(test_util.RequestContext(), map[string]events.EditForSaleConfigRequest{
		event1.Key: {NotForSale: []events.ObjectAndQuantity{{Object: "A-1"}}},
		event2.Key: {NotForSale: []events.ObjectAndQuantity{{Object: "A-2"}}},
	})
	require.NoError(t, err)

	retrievedEvent1, err := client.Events.Retrieve(test_util.RequestContext(), event1.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent1.ForSaleConfig)

	retrievedEvent2, err := client.Events.Retrieve(test_util.RequestContext(), event2.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-2"},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent2.ForSaleConfig)
}
