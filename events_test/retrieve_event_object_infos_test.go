package events_test

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/events"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveEventObjectInfos(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	eventObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event.Key, "A-1", "A-2")

	require.NoError(t, err)
	require.Equal(t, "A-1", eventObjectInfos["A-1"].Label)
	require.Equal(t, "A-2", eventObjectInfos["A-2"].Label)
}
