package charts

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/charts"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrievePublishedVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{Name: "chartName"})
	_ = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Name: "chartName"})

	drawing, err := client.Charts.RetrievePublishedVersion(test_util.RequestContext(), chart.Key)

	require.NoError(t, err)
	require.Equal(t, "chartName", drawing["name"])
}
