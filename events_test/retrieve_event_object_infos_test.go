package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveEventObjectInfos(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	eventObjectInfos, err := client.Events.RetrieveObjectInfos(event.Key, []string{"A-1", "A-2"})

	require.NoError(t, err)
	require.Equal(t, "A-1", eventObjectInfos["A-1"].Label)
	require.Equal(t, "A-2", eventObjectInfos["A-2"].Label)
}
