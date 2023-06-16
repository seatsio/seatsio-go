package reports

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestByLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey, "false")
	require.NoError(t, err)
	require.Len(t, chartReport.Items["A-1"], 1)
	item := chartReport.Items["A-1"][0]
	require.Equal(t, "A-1", item.Label)
	require.Equal(t, events.Labels{
		Own:    events.LabelAndType{Label: "1", Type: "seat"},
		Parent: events.LabelAndType{Label: "A", Type: "row"},
	}, item.Labels)
	require.Equal(t, events.IDs{Own: "1", Parent: "A", Section: ""}, item.IDs)
	require.Equal(t, "Cat1", item.CategoryLabel)
	require.Equal(t, "9", item.CategoryKey)
	require.Empty(t, item.Section)
	require.Empty(t, item.Entrance)
	require.Empty(t, item.Capacity)
	require.Equal(t, "seat", item.ObjectType)
	require.Empty(t, item.LeftNeighbour)
	require.Equal(t, "A-2", item.RightNeighbour)
	require.NotEmpty(t, item.DistanceToFocalPoint)
}
