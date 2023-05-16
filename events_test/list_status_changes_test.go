package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListStatusChanges(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "", "", "").All(2)
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestPropertiesOfStatusChange(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatus(&events.StatusChangeParams{Status: "s1", Events: []string{event.Key}, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "", "", "").All(1)
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
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:    events.ObjectStatusHeld,
		Events:    []string{event.Key},
		Objects:   []events.ObjectProperties{{ObjectId: "A-1"}},
		HoldToken: holdToken.HoldToken,
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "", "", "").All(1)
	require.NoError(t, err)

	statusChange := statusChanges[0]
	require.Equal(t, holdToken.HoldToken, statusChange.HoldToken)
}

func TestListStatusChangesWithFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "B-1"}}},
		{Status: "s4", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "A", "", "").All(2)
	require.NoError(t, err)

	require.Equal(t, "s4", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}

func TestListStatusChangesSortAsc(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "", "objectLabel", "asc").All(2)
	require.NoError(t, err)

	require.Equal(t, "s1", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s3", statusChanges[2].Status)
}

func TestListStatusChangesSortAscPageBefore(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, "", "objectLabel", "asc")
	statusChanges, err := statusChangeLister.All(10)
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageBefore(statusChanges[2].Id, 10)
	require.NoError(t, err)

	require.Equal(t, "s1", statusChangesPage.Items[0].Status)
	require.Equal(t, "s2", statusChangesPage.Items[1].Status)
}

func TestListStatusChangesSortAscPageAfter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChangeLister := client.Events.StatusChanges(event.Key, "", "objectLabel", "asc")
	statusChanges, err := statusChangeLister.All(10)
	require.NoError(t, err)

	statusChangesPage, err := statusChangeLister.ListPageAfter(statusChanges[0].Id, 10)
	require.NoError(t, err)

	require.Equal(t, "s2", statusChangesPage.Items[0].Status)
	require.Equal(t, "s3", statusChangesPage.Items[1].Status)
}

func TestListStatusChangesSortDesc(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch([]events.StatusChangeInBatchParams{
		{Status: "s1", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-1"}}},
		{Status: "s2", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-2"}}},
		{Status: "s3", Event: event.Key, Objects: []events.ObjectProperties{{ObjectId: "A-3"}}},
	})
	require.NoError(t, err)

	statusChanges, err := client.Events.StatusChanges(event.Key, "", "objectLabel", "desc").All(2)
	require.NoError(t, err)

	require.Equal(t, "s3", statusChanges[0].Status)
	require.Equal(t, "s2", statusChanges[1].Status)
	require.Equal(t, "s1", statusChanges[2].Status)
}
