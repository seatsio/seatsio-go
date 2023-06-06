package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey, "tag1")
	client.Charts.AddTag(chartKey, "tag2")

	retrievedChart, err := client.Charts.Retrieve(chartKey)

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
}

func TestRetrieveChartWithEvents(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	retrievedChart, err := client.Charts.RetrieveWithEvents(chartKey)
	require.NoError(t, err)
	require.Equal(t, 2, len(retrievedChart.Events))

}
