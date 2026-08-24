package reports

import (
	"strings"
	"testing"

	seatsio "github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/test_util"
	"github.com/stretchr/testify/require"
)

func TestFlatList(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	report, err := client.EventReports.FlatList(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 34, len(report))
	require.Equal(t, "A-1", report[0].Label)
}

func TestFlatListCsv(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	csv, err := client.EventReports.FlatListCsv(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.True(t, strings.Contains(csv, "A-1"))
}
