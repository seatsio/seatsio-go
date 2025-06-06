package seasons

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/seasons"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRetrieveSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	startTime := time.Now()
	season, err := client.Seasons.CreateSeasonWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{EventKeys: []string{"event1", "event2"}})
	require.NoError(t, err)
	partialSeason1, err := client.Seasons.CreatePartialSeason(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	partialSeason2, err := client.Seasons.CreatePartialSeason(test_util.RequestContext(), season.Key)
	require.NoError(t, err)

	retrievedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)

	require.NotEqual(t, 0, retrievedSeason.Id)
	require.Equal(t, season.Key, retrievedSeason.Key)
	require.Equal(t, season.ChartKey, retrievedSeason.ChartKey)
	require.Equal(t, events.TableBookingSupport.Inherit(), retrievedSeason.TableBookingConfig)
	require.Nil(t, retrievedSeason.ForSaleConfig)
	require.True(t, retrievedSeason.SupportsBestAvailable)
	require.True(t, retrievedSeason.CreatedOn.After(startTime))
	require.Nil(t, retrievedSeason.UpdatedOn)
	topLevelSeasonEventKeys := []string{retrievedSeason.Events[0].Key, retrievedSeason.Events[1].Key}
	require.Subset(t, topLevelSeasonEventKeys, []string{"event1", "event2"})
	require.Subset(t, retrievedSeason.PartialSeasonKeys, []string{partialSeason1.Key, partialSeason2.Key})
}
