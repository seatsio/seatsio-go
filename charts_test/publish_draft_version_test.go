package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublishDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "oldname"})
	client.Events.Create(&events.CreateEventParams{ChartKey: chart.Key})
	client.Charts.Update(chart.Key, &charts.UpdateChartParams{Name: "newname"})

	err = client.Charts.PublishDraftVersion(chart.Key)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chart.Key)
	require.NoError(t, err)
	require.Equal(t, "newname", retrievedChart.Name)
	require.Equal(t, "PUBLISHED", retrievedChart.Status)
}
