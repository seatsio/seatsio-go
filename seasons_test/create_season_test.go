package seasons_test

import (
	"testing"
	"time"

	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/seasons"
	"github.com/seatsio/seatsio-go/v11/shared"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
)

func TestChartKeyIsRequired(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	startTime := time.Now()
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason"})
	require.NoError(t, err)
	require.Equal(t, "aSeason", season.Key)
}

func TestNameCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Name: "aSeason"})
	require.NoError(t, err)
	require.Equal(t, "aSeason", season.Name)
}

func TestNumberOfEventsCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason", NumberOfEvents: 2})
	require.NoError(t, err)
	require.Equal(t, 2, len(season.Events))
}

func TestEventKeysCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason", EventKeys: []string{"event1", "event2"}})
	require.NoError(t, err)
	require.Equal(t, []string{"event1", "event2"}, []string{season.Events[0].Key, season.Events[1].Key})
}

func TestChannelsCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	channels := []events.CreateChannelParams{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason", Channels: &channels})
	require.NoError(t, err)

	expectedChannels := []events.Channel{
		{Key: "aaa", Name: "bbb", Color: "#101010", Index: 1, Objects: []string{"A-1", "A-2"}},
		{Key: "ccc", Name: "ddd", Color: "#F2F2F2", Index: 2, Objects: []string{}},
	}
	require.Equal(t, expectedChannels, season.Channels)
}

func TestForSaleConfigCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	forSaleConfig := &events.ForSaleConfig{
		ForSale:    false,
		Objects:    []string{"A-1"},
		AreaPlaces: map[string]int{"GA1": 5},
		Categories: []string{"Cat1"},
	}

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{Key: "aSeason", ForSaleConfig: forSaleConfig})
	require.NoError(t, err)

	require.Equal(t, forSaleConfig, season.ForSaleConfig)
}

func TestCreateSeasonWithObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{
		ObjectCategories: &objectCategories,
	})
	require.NoError(t, err)

	require.Equal(t, objectCategories, season.ObjectCategories)
}

func TestCreateSeasonWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{
		Categories: &categories,
	})
	require.NoError(t, err)

	require.Contains(t, season.Categories, category)
}

func TestForSalePropagatedCanBePassedIn(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{ForSalePropagated: shared.OptionalBool(false)})
	require.NoError(t, err)

	require.False(t, season.ForSalePropagated)
}
