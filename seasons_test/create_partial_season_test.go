package seasons

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/seasons"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	topLevelSeason, err := client.Seasons.CreateSeasonWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aTopLevelSeason"})
	require.NoError(t, err)

	partialSeason, err := client.Seasons.CreatePartialSeasonWithOptions(test_util.RequestContext(), topLevelSeason.Key, &seasons.CreatePartialSeasonParams{Key: "aPartialSeason"})
	require.NoError(t, err)

	require.Equal(t, "aPartialSeason", partialSeason.Key)
	require.True(t, partialSeason.IsPartialSeason)
	require.Equal(t, topLevelSeason.Key, *partialSeason.TopLevelSeasonKey)
}

func TestNameCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	topLevelSeason, err := client.Seasons.CreateSeasonWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aTopLevelSeason"})
	require.NoError(t, err)

	partialSeason, err := client.Seasons.CreatePartialSeasonWithOptions(test_util.RequestContext(), topLevelSeason.Key, &seasons.CreatePartialSeasonParams{Name: "aPartialSeason"})
	require.NoError(t, err)

	require.Equal(t, "aPartialSeason", partialSeason.Name)
}

func TestEventKeysCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	topLevelSeason, err := client.Seasons.CreateSeasonWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aTopLevelSeason", EventKeys: []string{"event1", "event2", "event3"}})
	require.NoError(t, err)

	partialSeason, err := client.Seasons.CreatePartialSeasonWithOptions(test_util.RequestContext(), topLevelSeason.Key, &seasons.CreatePartialSeasonParams{EventKeys: []string{"event1", "event3"}})
	require.NoError(t, err)
	require.Subset(t, []string{partialSeason.Events[0].Key, partialSeason.Events[1].Key}, []string{"event1", "event3"})
	require.Equal(t, []string{partialSeason.Key}, partialSeason.Events[0].PartialSeasonKeysForEvent)
}
