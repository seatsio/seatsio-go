package events

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		ObjectCategories: &map[string]events.CategoryKey{
			"A-1": {10},
		},
	}})
	require.NoError(t, err)

	err = client.Events.RemoveObjectCategories(event.Key)
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Empty(t, updatedEvent.ObjectCategories)
}
