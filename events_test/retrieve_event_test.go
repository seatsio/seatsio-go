package events_test

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRetrieveEvent(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	retrievedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)

	require.NoError(t, err)
	require.NotZero(t, retrievedEvent.Id)
	require.NotNil(t, retrievedEvent.Key)
	require.Equal(t, chartKey, retrievedEvent.ChartKey)
	require.Equal(t, events.TableBookingConfig{Mode: "INHERIT"}, retrievedEvent.TableBookingConfig)
	require.True(t, retrievedEvent.SupportsBestAvailable)
	require.Nil(t, retrievedEvent.ForSaleConfig)
	require.True(t, retrievedEvent.CreatedOn.After(start))
	require.Nil(t, retrievedEvent.UpdatedOn)
	require.Equal(t, []events.Category{
		{Key: events.CategoryKey{Key: 9}, Label: "Cat1", Color: "#87A9CD", Accessible: false},
		{Key: events.CategoryKey{Key: 10}, Label: "Cat2", Color: "#5E42ED", Accessible: false},
		{Key: events.CategoryKey{Key: "string11"}, Label: "Cat3", Color: "#5E42BB", Accessible: false},
	}, retrievedEvent.Categories)
	require.Nil(t, retrievedEvent.PartialSeasonKeysForEvent)
}
