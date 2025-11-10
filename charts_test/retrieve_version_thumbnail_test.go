package charts

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/charts"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRetrievePublishedVersionThumbnail(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{VenueType: "SIMPLE"})
	require.NoError(t, err)

	file, err := client.Charts.RetrievePublishedVersionThumbnail(test_util.RequestContext(), chart.Key)
	require.NoError(t, err)
	require.Contains(t, file.Name(), chart.Key)
	_ = os.Remove(file.Name())
}

func TestRetrieveDraftVersionThumbnail(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{Name: "oldname"})
	require.NoError(t, err)
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chart.Key})
	_ = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Name: "newname"})

	file, err := client.Charts.RetrieveDraftVersionThumbnail(test_util.RequestContext(), chart.Key)
	require.NoError(t, err)
	require.Contains(t, file.Name(), chart.Key)
	_ = os.Remove(file.Name())
}
