package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopyDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "newname"})

	copiedChart, err := client.Charts.CopyDraftVersion(chartKey)
	require.NoError(t, err)
	require.Equal(t, "newname (copy)", copiedChart.Name)
}
