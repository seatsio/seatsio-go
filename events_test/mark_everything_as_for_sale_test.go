package events_test

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarkEverythingAsForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsNotForSale(test_util.RequestContext(), event.Key, &events.ForSaleConfigParams{
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	err = client.Events.MarkEverythingAsForSale(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Empty(t, retrievedEvent.ForSaleConfig)
}
