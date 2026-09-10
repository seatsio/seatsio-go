package reports

import (
	"strings"
	"testing"

	seatsio "github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/seasons"
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

func TestFlatListWithSeasonBookingsNotPropagated(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{NumberOfEvents: 1})
	require.NoError(t, err)
	event := season.Events[0]
	_, err = client.Events.Book(test_util.RequestContext(), season.Key, "A-1", "A-2")
	require.NoError(t, err)
	_, err = client.Events.Book(test_util.RequestContext(), event.Key, "A-3")
	require.NoError(t, err)

	reportWithPropagation, err := client.EventReports.FlatList(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	reportWithoutPropagation, err := client.EventReports.WithSeasonBookingsNotPropagated().FlatList(test_util.RequestContext(), season.Key)
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, findByLabel(reportWithPropagation, "A-3").Status)
	require.NotEqual(t, events.BOOKED, findByLabel(reportWithoutPropagation, "A-3").Status)
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
