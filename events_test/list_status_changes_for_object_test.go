package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListStatusChangesForObject(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChangesForObject(event.Key, "A-1").All(2)
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}
