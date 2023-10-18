package charts

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/charts"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrievePublishedVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "chartName"})
	_ = client.Charts.Update(chart.Key, &charts.UpdateChartParams{Name: "chartName"})

	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)

	require.NoError(t, err)
	require.Equal(t, "chartName", drawing["name"])
}
