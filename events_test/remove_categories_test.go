package events

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateEventRemoveObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Categories: &[]events.Category{
			{
				Key:        events.CategoryKey{Key: "CatA"},
				Label:      "cat-a",
				Color:      "#111111",
				Accessible: false,
			},
		},
	}})
	require.NoError(t, err)

	err = client.Events.RemoveCategories(event.Key)
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(event.Key)
	require.NoError(t, err)
	require.Empty(t, updatedEvent.ObjectCategories)
}
