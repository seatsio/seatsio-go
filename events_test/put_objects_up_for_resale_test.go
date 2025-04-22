package events

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPutUpForResaleSingleEvent(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.PutUpForResale(test_util.RequestContext(), event.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.RESALE, objects.Objects["A-1"].Status)
	require.Equal(t, events.RESALE, objects.Objects["A-2"].Status)

	info, _ := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1", "A-2", "A-3")
	require.Equal(t, events.RESALE, info["A-1"].Status)
	require.Equal(t, events.RESALE, info["A-2"].Status)
	require.Equal(t, events.FREE, info["A-3"].Status)
}
