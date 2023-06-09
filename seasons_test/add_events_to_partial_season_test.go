package seasons

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/seasons"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddingEventToPartialSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	season, err := client.Seasons.CreateSeason(chartKey, seasons.SeasonSupport.WithKey("aSeason"), seasons.SeasonSupport.WithEventKeys("event1", "event2"))
	require.NoError(t, err)
	partialSeason, err := client.Seasons.CreatePartialSeason(season.Key, seasons.PartialSeasonSupport.WithKey("aPartialSeason"))
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.AddEventsToPartialSeason(season.Key, partialSeason.Key, "event1", "event2")
	require.NoError(t, err)
	require.Subset(t, []string{updatedSeason.Events[0].Key, updatedSeason.Events[1].Key}, []string{"event1", "event2"})
}
