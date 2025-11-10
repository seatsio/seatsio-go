package seasons

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/seasons"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveEventFromPartialSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason", EventKeys: []string{"event1", "event2"}})
	require.NoError(t, err)
	partialSeason, err := client.Seasons.CreatePartialSeasonWithOptions(test_util.RequestContext(), season.Key, &seasons.CreatePartialSeasonParams{Key: "aPartialSeason", EventKeys: []string{"event1", "event2"}})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.RemoveEventFromPartialSeason(test_util.RequestContext(), season.Key, partialSeason.Key, "event2")
	require.NoError(t, err)
	require.Len(t, updatedSeason.Events, 1)
	require.Contains(t, updatedSeason.Events[0].Key, "event1")
}
