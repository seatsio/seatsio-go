package seasons

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/seasons"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	topLevelSeason, err := client.Seasons.CreateSeason(chartKey, seasons.SeasonSupport.WithKey("aTopLevelSeason"))
	require.NoError(t, err)

	partialSeason, err := client.Seasons.CreatePartialSeason(topLevelSeason.Key, seasons.PartialSeasonSupport.WithKey("aPartialSeason"))
	require.NoError(t, err)

	require.Equal(t, "aPartialSeason", partialSeason.Key)
	require.True(t, partialSeason.IsPartialSeason)
	require.Equal(t, topLevelSeason.Key, *partialSeason.TopLevelSeasonKey)
}

func TestEventKeysCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	support := seasons.SeasonSupport
	topLevelSeason, err := client.Seasons.CreateSeason(chartKey, support.WithKey("aTopLevelSeason"), support.WithEventKeys("event1", "event2", "event3"))
	require.NoError(t, err)

	partialSeason, err := client.Seasons.CreatePartialSeason(topLevelSeason.Key, seasons.PartialSeasonSupport.WithEventKeys("event1", "event3"))
	require.NoError(t, err)
	require.Subset(t, []string{partialSeason.Events[0].Key, partialSeason.Events[1].Key}, []string{"event1", "event3"})
}
