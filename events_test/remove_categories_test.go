package events

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/events"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateEventRemoveObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
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

	err = client.Events.RemoveCategories(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	updatedEvent, err := client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Empty(t, updatedEvent.ObjectCategories)
}
