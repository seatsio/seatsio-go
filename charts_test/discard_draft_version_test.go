package charts

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/charts"
	"github.com/seatsio/seatsio-go/v10/events"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDiscardDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(test_util.RequestContext(), chartKey, &charts.UpdateChartParams{Name: "newname"})

	err := client.Charts.DiscardDraftVersion(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	retrievedChart, err := client.Charts.Retrieve(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "Sample chart", retrievedChart.Name)
	require.Equal(t, "PUBLISHED", retrievedChart.Status)
}
