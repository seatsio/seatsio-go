package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/test_util"
	"github.com/stretchr/testify/require"
)

func TestMarkEverythingAsNotForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkEverythingAsNotForSale(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, &events.ForSaleConfig{
		ForSale:    true,
		Objects:    []string{},
		AreaPlaces: map[string]int{},
		Categories: []string{},
	}, retrievedEvent.ForSaleConfig)
}


