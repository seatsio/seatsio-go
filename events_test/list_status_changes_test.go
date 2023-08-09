package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListStatusChanges(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key).All()
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key).All(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestPropertiesOfStatusChange(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, "s1")
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key).All(shared.Pagination.PageSize(1))
	require.NoError(t, err)

	statusChange := statusChanges[0]
	require.NotEmpty(t, statusChange.Id)
	require.NotEmpty(t, statusChange.Date)
	require.Equal(t, "s1", statusChange.Status)
	require.Equal(t, "A-1", statusChange.ObjectLabel)
	require.Equal(t, event.Id, statusChange.EventId)
	require.Equal(t, "API_CALL", statusChange.Origin.Type)
	require.NotEmpty(t, statusChange.Origin.Ip)
	require.True(t, statusChange.IsPresentOnChart)
	require.Empty(t, statusChange.NotPresentOnChartReason)
}

func TestPropertiesOfStatusChangeHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:    events.HELD,
			Objects:   []events.ObjectProperties{{ObjectId: "A-1"}},
			HoldToken: holdToken.HoldToken,
		},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key).All(shared.Pagination.PageSize(1))
	require.NoError(t, err)

	statusChange := statusChanges[0]
	require.Equal(t, holdToken.HoldToken, statusChange.HoldToken)
}

func TestListStatusChangesWithFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "B-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s4", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithFilter("A")).All()
	require.NoError(t, err)

	require.Equal(t, "s4", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesWithFilterAndLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "B-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s4", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithFilter("A")).All(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, "s4", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesWithFilterAndSort(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "B-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s4", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithFilter("A"), events.EventSupport.WithSortAsc("objectLabel")).All()
	require.NoError(t, err)

	require.Equal(t, "s1", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s4", statusChanges[2].Status)
}

func TestListStatusChangesSortAsc(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel")).All()
	require.NoError(t, err)

	require.Equal(t, "s1", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s3", statusChanges[2].Status)
}

func TestListStatusChangesSortAscWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel")).All(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, "s1", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s3", statusChanges[2].Status)
}

func TestListStatusChangesSortAscPageBefore(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel"))
	statusChanges, err := statusChangeLister.All()
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageBefore(statusChanges[2].Id)
	require.NoError(t, err)

	require.Equal(t, "s1", statusChangesPage.Items[0].Status)
	require.Equal(t, "s2", statusChangesPage.Items[1].Status)
}

func TestListStatusChangesSortAscPageBeforeWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel"))
	statusChanges, err := statusChangeLister.All()
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageBefore(statusChanges[2].Id, shared.Pagination.PageSize(10))
	require.NoError(t, err)

	require.Equal(t, "s1", statusChangesPage.Items[0].Status)
	require.Equal(t, "s2", statusChangesPage.Items[1].Status)
}

func TestListStatusChangesSortAscPageAfter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel"))
	statusChanges, err := statusChangeLister.All()
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageAfter(statusChanges[0].Id)
	require.NoError(t, err)

	require.Equal(t, "s2", statusChangesPage.Items[0].Status)
	require.Equal(t, "s3", statusChangesPage.Items[1].Status)
}

func TestListStatusChangesSortAscPageAfterWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortAsc("objectLabel"))
	statusChanges, err := statusChangeLister.All()
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageAfter(statusChanges[0].Id, shared.Pagination.PageSize(1))
	require.NoError(t, err)

	require.Equal(t, "s2", statusChangesPage.Items[0].Status)
}

func TestListStatusChangesSortDesc(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortDesc("objectLabel")).All()
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesSortDescWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s1", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s2", Objects: []events.ObjectProperties{{ObjectId: "A-2"}}}},
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "s3", Objects: []events.ObjectProperties{{ObjectId: "A-3"}}}},
	)
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, events.EventSupport.WithSortDesc("objectLabel")).All(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}
