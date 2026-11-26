package events_test

import (
	"testing"

	"github.com/seatsio/seatsio-go/v13"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/test_util"
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

func TestMoveEventsToNewChartCopy(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	movedEvents, err := client.Events.MoveEventsToNewChartCopy(test_util.RequestContext(), []string{event2.Key, event1.Key})
	require.NoError(t, err)

	require.Equal(t, []string{event2.Key, event1.Key}, []string{movedEvents[0].Key, movedEvents[1].Key})
	require.NotEqual(t, chartKey, movedEvents[0].ChartKey)
	require.Equal(t, movedEvents[0].ChartKey, movedEvents[1].ChartKey)
}
