package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateExtraDatas(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	err = client.Events.UpdateExtraData(event.Key, map[string]events.ExtraData{
		"A-1": {"foo": "bar"},
		"A-2": {"foo": "baz"},
	})
	require.NoError(t, err)

	eventObjectInfos, err := client.Events.RetrieveObjectInfo(event.Key, "A-1", "A-2")
	require.NoError(t, err)
	require.Equal(t, events.ExtraData{"foo": "bar"}, eventObjectInfos["A-1"].ExtraData)
	require.Equal(t, events.ExtraData{"foo": "baz"}, eventObjectInfos["A-2"].ExtraData)
}
