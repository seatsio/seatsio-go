package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrievePublishedVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "chartName"})
	client.Charts.Update(chart.Key, &charts.UpdateChartParams{Name: "chartName"})

	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)

	require.NoError(t, err)
	require.Equal(t, "chartName", drawing["name"])
}
