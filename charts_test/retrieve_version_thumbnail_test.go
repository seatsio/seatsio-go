package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRetrievePublishedVersionThumbnail(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(&charts.CreateChartParams{VenueType: "BOOTHS"})
	require.NoError(t, err)

	file, err := client.Charts.RetrievePublishedVersionThumbnail(chart.Key)
	require.NoError(t, err)
	require.Contains(t, file.Name(), chart.Key)
	_ = os.Remove(file.Name())
}

func TestRetrieveDraftVersionThumbnail(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "oldname"})
	require.NoError(t, err)
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chart.Key})
	_ = client.Charts.Update(chart.Key, &charts.UpdateChartParams{Name: "newname"})

	file, err := client.Charts.RetrieveDraftVersionThumbnail(chart.Key)
	require.NoError(t, err)
	require.Contains(t, file.Name(), chart.Key)
	_ = os.Remove(file.Name())
}
