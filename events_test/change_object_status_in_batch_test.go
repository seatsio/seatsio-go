package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeObjectStatusInBatch(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event1, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: events.ObjectStatusBooked, Event: event1.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: events.ObjectStatusBooked, Event: event2.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
	})
	require.NoError(t, err)

	require.Equal(t, events.ObjectStatusBooked, result.Results[0].Objects["A-1"].Status)
	require.Equal(t, events.ObjectStatusBooked, result.Results[1].Objects["A-2"].Status)
}
