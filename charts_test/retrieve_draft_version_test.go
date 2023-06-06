package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "newName"})

	drawing, err := client.Charts.RetrieveDraftVersion(chartKey)
	require.NoError(t, err)
	require.Equal(t, "newName", drawing["name"])
}
