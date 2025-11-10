package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/seasons"
	"github.com/seatsio/seatsio-go/v12/shared"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestUpdateSeasonEventKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{EventKey: "newKey"})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), "newKey")
	require.NoError(t, err)
	require.Equal(t, "newKey", updatedSeason.Key)
}

func TestUpdateSeasonName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{
		Name: "foo",
	})
	require.NoError(t, err)

	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{
		Name: "bar",
	})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	require.Equal(t, "bar", updatedSeason.Name)
}

func TestUpdateSeasonTableBookingConfig(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	tableBookingConfig := events.TableBookingConfig{Mode: events.CUSTOM, Tables: map[string]events.TableBookingMode{
		"T1": events.BY_TABLE, "T2": events.BY_SEAT,
	}}
	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{
		TableBookingConfig: &tableBookingConfig,
	})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	require.Equal(t, tableBookingConfig, updatedSeason.TableBookingConfig)
}

func TestUpdateSeasonObjectCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	objectCategories := map[string]events.CategoryKey{
		"A-1": {10},
	}
	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{
		ObjectCategories: &objectCategories,
	})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	require.Equal(t, objectCategories, updatedSeason.ObjectCategories)
}

func TestUpdateSeasonCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	category := events.Category{Key: events.CategoryKey{Key: "eventCategory"}, Label: "event-level category", Color: "#AAABBB"}
	categories := []events.Category{
		category,
	}
	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{
		Categories: &categories,
	})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	require.Contains(t, updatedSeason.Categories, category)
}

func TestUpdateSeasonForSalePropagated(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	err = client.Seasons.Update(test_util.RequestContext(), season.Key, &seasons.UpdateSeasonParams{
		ForSalePropagated: shared.OptionalBool(false),
	})
	require.NoError(t, err)

	updatedSeason, err := client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	require.NoError(t, err)
	require.False(t, updatedSeason.ForSalePropagated)
}
