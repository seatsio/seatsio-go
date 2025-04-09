package charts

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/charts"
	"github.com/seatsio/seatsio-go/v10/events"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.AddTag(test_util.RequestContext(), chartKey, "tag1")
	_ = client.Charts.AddTag(test_util.RequestContext(), chartKey, "tag2")

	retrievedChart, err := client.Charts.Retrieve(test_util.RequestContext(), chartKey)

	require.NoError(t, err)
	require.NotEqual(t, 0, retrievedChart.Id)
	require.NotEmpty(t, retrievedChart.Key)
	require.Equal(t, "NOT_USED", retrievedChart.Status)
	require.Equal(t, "Sample chart", retrievedChart.Name)
	require.NotEmpty(t, retrievedChart.PublishedVersionThumbnailUrl)
	require.Empty(t, retrievedChart.DraftVersionThumbnailUrl)
	require.Nil(t, retrievedChart.Events)
	require.Contains(t, retrievedChart.Tags, "tag1")
	require.Contains(t, retrievedChart.Tags, "tag2")
	require.False(t, retrievedChart.Archived)
	require.Nil(t, retrievedChart.Zones)
}

func TestRetrieveChartZones(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)

	retrievedChart, err := client.Charts.Retrieve(test_util.RequestContext(), chartKey)

	require.NoError(t, err)
	require.Equal(t, []charts.Zone{{"finishline", "Finish Line"}, {"midtrack", "Mid Track"}}, retrievedChart.Zones)
}

func TestRetrieveChartWithEvents(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	retrievedChart, err := client.Charts.RetrieveWithEvents(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, 2, len(retrievedChart.Events))

}
