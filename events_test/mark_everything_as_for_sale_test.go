package events_test

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarkEverythingAsForSale(t *testing.T) {
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

	err = client.Events.MarkEverythingAsForSale(event.Key)
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Empty(t, retrievedEvent.ForSaleConfig)
}
