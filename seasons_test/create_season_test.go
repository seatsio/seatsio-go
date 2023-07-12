package seasons_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/seasons"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestChartKeyIsRequired(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	startTime := time.Now()
	season, err := client.Seasons.CreateSeason(chartKey)
	require.NoError(t, err)

	require.NotEqual(t, 0, season.Id)
	require.True(t, season.IsTopLevelSeason)
	require.Nil(t, season.TopLevelSeasonKey)
	require.NotNil(t, season.Key)
	require.Empty(t, season.PartialSeasonKeys)
	require.NotEqual(t, 0, season.Id)
	require.Equal(t, chartKey, season.ChartKey)
	require.Equal(t, events.TableBookingConfig{Mode: events.INHERIT, Tables: nil}, season.TableBookingConfig)
	require.True(t, season.SupportsBestAvailable)
	require.Nil(t, season.ForSaleConfig)
	require.True(t, season.CreatedOn.After(startTime))
	require.Nil(t, season.UpdatedOn)
	require.Empty(t, season.Events)
}

func TestKeyCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	season, err := client.Seasons.CreateSeason(chartKey, seasons.SeasonSupport.WithKey("aSeason"))
	require.NoError(t, err)
	require.Equal(t, "aSeason", season.Key)
}

func TestNumberOfEventsCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	support := seasons.SeasonSupport
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	season, err := client.Seasons.CreateSeason(chartKey, support.WithKey("aSeason"), support.WithNumberOfEvents(2))
	require.NoError(t, err)
	require.Equal(t, 2, len(season.Events))
}

func TestEventKeysCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	support := seasons.SeasonSupport
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	season, err := client.Seasons.CreateSeason(chartKey, support.WithKey("aSeason"), support.WithEventKeys("event1", "event2"))
	require.NoError(t, err)
	require.Subset(t, []string{season.Events[0].Key, season.Events[1].Key}, []string{"event1", "event2"})
}
