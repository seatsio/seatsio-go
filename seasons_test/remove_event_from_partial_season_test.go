package seasons

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/seasons"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveEventFromPartialSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	season, err := client.Seasons.CreateSeason(chartKey, seasons.SeasonSupport.WithKey("aSeason"), seasons.SeasonSupport.WithEventKeys("event1", "event2"))
	require.NoError(t, err)
	support := seasons.PartialSeasonSupport
	partialSeason, err := client.Seasons.CreatePartialSeason(season.Key, support.WithKey("aPartialSeason"), support.WithEventKeys("event1", "event2"))
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.RemoveEventFromPartialSeason(season.Key, partialSeason.Key, "event2")
	require.NoError(t, err)
	require.Len(t, updatedSeason.Events, 1)
	require.Contains(t, updatedSeason.Events[0].Key, "event1")
}
