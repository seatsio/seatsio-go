package seasons

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateEventsWithEventKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateSeason(chartKey)
	require.NoError(t, err)

	events, err := client.Seasons.CreateEventsWithEventKeys(season.Key, "event1", "event2")
	require.NoError(t, err)
	require.Subset(t, []string{events[0].Key, events[1].Key}, []string{"event1", "event2"})
}

func TestCreateEventsWithNumberOfEvents(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateSeason(chartKey)
	require.NoError(t, err)

	events, err := client.Seasons.CreateNumberOfEvents(season.Key, 2)
	require.NoError(t, err)
	require.Len(t, events, 2)
}
