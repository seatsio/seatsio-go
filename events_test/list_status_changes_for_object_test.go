package events_test

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/shared"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListStatusChangesForObject(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChangesForObject(event.Key, "A-1").All()
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesForObjectWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChangesForObject(event.Key, "A-1").All(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}
