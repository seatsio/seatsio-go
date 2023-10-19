package charts

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/charts"
	"github.com/seatsio/seatsio-go/v6/events"
	"github.com/seatsio/seatsio-go/v6/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDiscardDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "newname"})

	err := client.Charts.DiscardDraftVersion(chartKey)
	require.NoError(t, err)

	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.Equal(t, "Sample chart", retrievedChart.Name)
	require.Equal(t, "PUBLISHED", retrievedChart.Status)
}
