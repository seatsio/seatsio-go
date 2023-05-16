package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteEvent(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.Delete(event.Key)
	require.NoError(t, err)

	_, err = client.Events.Retrieve(event.Key)
	require.Equal(t, "EVENT_NOT_FOUND", err.(*shared.SeatsioError).Code)
}

func TestDeleteSeason(t *testing.T) {
	// TODO
}
