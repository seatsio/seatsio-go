package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAll(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event1, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	events, err := client.Events.ListAll(2)
	require.NoError(t, err)

	require.Equal(t, event3.Key, events[0].Key)
	require.Equal(t, event2.Key, events[1].Key)
	require.Equal(t, event1.Key, events[2].Key)
}

func TestListFirstPage(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	_, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListFirstPage(2)
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event2.Id, page.NextPageStartsAfter)
	require.Equal(t, int64(0), page.PreviousPageEndsBefore)
}

func TestListPageAfter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	_, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event4, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageAfter(event4.Id, 2)
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event2.Id, page.NextPageStartsAfter)
	require.Equal(t, event3.Id, page.PreviousPageEndsBefore)
}

func TestListPageBefore(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event1, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	_, err = client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageBefore(event1.Id, 2)
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event2.Id, page.NextPageStartsAfter)
	require.Equal(t, event3.Id, page.PreviousPageEndsBefore)
}
