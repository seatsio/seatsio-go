package charts

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/charts"
	"github.com/seatsio/seatsio-go/v6/events"
	"github.com/seatsio/seatsio-go/v6/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublishDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "oldname"})
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chart.Key})
	_ = client.Charts.Update(chart.Key, &charts.UpdateChartParams{Name: "newname"})

	err = client.Charts.PublishDraftVersion(chart.Key)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chart.Key)
	require.NoError(t, err)
	require.Equal(t, "newname", retrievedChart.Name)
	require.Equal(t, "PUBLISHED", retrievedChart.Status)
}
