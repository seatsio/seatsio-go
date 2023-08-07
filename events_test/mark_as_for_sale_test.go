package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarkAsForSale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.MarkAsForSale(event.Key, &events.ForSaleConfigParams{
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)

	require.Equal(t, &events.ForSaleConfig{
		ForSale:    true,
		Objects:    []string{"o1", "o2"},
		AreaPlaces: map[string]int{"GA1": 3},
		Categories: []string{"cat1", "cat2"},
	}, retrievedEvent.ForSaleConfig)
}
