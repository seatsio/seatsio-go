package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestMoveEventToNewChartCopy(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	movedEvent, err := client.Events.MoveEventToNewChartCopy(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.NotEqual(t, event.ChartKey, movedEvent.ChartKey)
	require.Equal(t, event.Key, movedEvent.Key)
}
