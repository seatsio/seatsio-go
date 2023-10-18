package charts

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/charts"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "newName"})

	drawing, err := client.Charts.RetrieveDraftVersion(chartKey)
	require.NoError(t, err)
	require.Equal(t, "newName", drawing["name"])
}
