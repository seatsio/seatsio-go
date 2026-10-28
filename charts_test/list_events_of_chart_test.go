package charts

import (
	"github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/shared"
	"github.com/seatsio/seatsio-go/v13/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAllEventsOfChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	otherChartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	_, err = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: otherChartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	retrievedEvents, err := client.Charts.Events(test_util.RequestContext(), chartKey).All()
	require.NoError(t, err)

	require.Equal(t, 2, len(retrievedEvents))
	require.Equal(t, event2.Key, retrievedEvents[0].Key)
	require.Equal(t, event1.Key, retrievedEvents[1].Key)
}

func TestListEventsOfChartPageByPage(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event3, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	firstPage, err := client.Charts.Events(test_util.RequestContext(), chartKey).ListFirstPage(shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, 2, len(firstPage.Items))
	require.Equal(t, event3.Key, firstPage.Items[0].Key)
	require.Equal(t, event2.Key, firstPage.Items[1].Key)
	require.Equal(t, event2.Id, firstPage.NextPageStartsAfter)
	require.Equal(t, int64(0), firstPage.PreviousPageEndsBefore)

	pageAfter, err := client.Charts.Events(test_util.RequestContext(), chartKey).ListPageAfter(firstPage.NextPageStartsAfter, shared.Pagination.PageSize(2))
	require.NoError(t, err)

	require.Equal(t, 1, len(pageAfter.Items))
	require.Equal(t, event1.Key, pageAfter.Items[0].Key)
	require.Equal(t, int64(0), pageAfter.NextPageStartsAfter)
	require.Equal(t, event1.Id, pageAfter.PreviousPageEndsBefore)
}
