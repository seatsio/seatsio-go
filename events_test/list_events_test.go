package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAll(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	retrievedEvents, err := client.Events.ListAll()
	require.NoError(t, err)

	require.Equal(t, event3.Key, retrievedEvents[0].Key)
	require.Equal(t, event2.Key, retrievedEvents[1].Key)
	require.Equal(t, event1.Key, retrievedEvents[2].Key)
}

func TestListAllWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	retrievedEvents, err := client.Events.ListAll(shared.Pagination.PageSize(1))
	require.NoError(t, err)

	require.Equal(t, event3.Key, retrievedEvents[0].Key)
	require.Equal(t, event2.Key, retrievedEvents[1].Key)
	require.Equal(t, event1.Key, retrievedEvents[2].Key)
}

func TestListFirstPage(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListFirstPage()
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event1.Key, page.Items[2].Key)
}

func TestListFirstPageWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListFirstPage(shared.Pagination.PageSize(2))
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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event4, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageAfter(event4.Id)
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event1.Key, page.Items[2].Key)
	require.Equal(t, event3.Id, page.PreviousPageEndsBefore)
}

func TestListPageAfterWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event4, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageAfter(event4.Id, shared.Pagination.PageSize(2))
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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event4, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageBefore(event1.Id)
	require.NoError(t, err)

	require.Equal(t, event4.Key, page.Items[0].Key)
	require.Equal(t, event3.Key, page.Items[1].Key)
	require.Equal(t, event2.Key, page.Items[2].Key)
	require.Equal(t, event2.Id, page.NextPageStartsAfter)
	require.Equal(t, int64(0), page.PreviousPageEndsBefore)
}

func TestListPageBeforeWithLimit(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	_, err = client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	page, err := client.Events.ListPageBefore(event1.Id, shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, event3.Key, page.Items[0].Key)
	require.Equal(t, event2.Key, page.Items[1].Key)
	require.Equal(t, event2.Id, page.NextPageStartsAfter)
	require.Equal(t, event3.Id, page.PreviousPageEndsBefore)
}
