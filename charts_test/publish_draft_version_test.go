package charts

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/charts"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublishDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{Name: "oldname"})
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chart.Key})
	_ = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Name: "newname"})

	err = client.Charts.PublishDraftVersion(test_util.RequestContext(), chart.Key)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(test_util.RequestContext(), chart.Key)
	require.NoError(t, err)
	require.Equal(t, "newname", retrievedChart.Name)
	require.Equal(t, "PUBLISHED", retrievedChart.Status)
}
