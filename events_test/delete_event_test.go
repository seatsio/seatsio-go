package events_test

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/shared"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Delete(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	_, err = client.Events.Retrieve(test_util.RequestContext(), event.Key)
	require.Equal(t, "EVENT_NOT_FOUND", err.(*shared.SeatsioError).Code)
}

func TestDeleteSeason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	err = client.Events.Delete(test_util.RequestContext(), season.Key)
	require.NoError(t, err)

	_, err = client.Seasons.Retrieve(test_util.RequestContext(), season.Key)
	seatsioError := err.(*shared.SeatsioError)
	require.Equal(t, "EVENT_NOT_FOUND", seatsioError.Code)
}
