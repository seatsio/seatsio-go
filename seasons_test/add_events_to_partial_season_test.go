package seasons

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/seasons"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddingEventToPartialSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateSeasonWithOptions(chartKey, &seasons.CreateSeasonParams{
		Key:       "aSeason",
		EventKeys: []string{"event1", "event2"},
	})
	require.NoError(t, err)
	partialSeason, err := client.Seasons.CreatePartialSeasonWithOptions(season.Key, &seasons.CreatePartialSeasonParams{Key: "aPartialSeason"})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.AddEventsToPartialSeason(season.Key, partialSeason.Key, "event1", "event2")
	require.NoError(t, err)
	require.Subset(t, []string{updatedSeason.Events[0].Key, updatedSeason.Events[1].Key}, []string{"event1", "event2"})
}
