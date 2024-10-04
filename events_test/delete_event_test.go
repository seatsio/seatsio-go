package events_test

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/events"
	"github.com/seatsio/seatsio-go/v8/shared"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Delete(event.Key)
	require.NoError(t, err)

	_, err = client.Events.Retrieve(event.Key)
	require.Equal(t, "EVENT_NOT_FOUND", err.(*shared.SeatsioError).Code)
}

func TestDeleteSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.CreateSeason(chartKey)
	require.NoError(t, err)

	err = client.Events.Delete(season.Key)
	require.NoError(t, err)

	_, err = client.Seasons.Retrieve(season.Key)
	seatsioError := err.(*shared.SeatsioError)
	require.Equal(t, "EVENT_NOT_FOUND", seatsioError.Code)
}
