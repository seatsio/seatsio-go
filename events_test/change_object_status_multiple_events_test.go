package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	event1, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.EventCreationParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:  events.ObjectStatusBooked,
		Events:  []string{event1.Key, event2.Key},
		Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
	})
	require.NoError(t, err)

	event1ObjectInfos, err := client.Events.RetrieveObjectInfos(event1.Key, []string{"A-1", "A-2"})
	require.NoError(t, err)
	event2ObjectInfos, err := client.Events.RetrieveObjectInfos(event2.Key, []string{"A-1", "A-2"})
	require.NoError(t, err)

	require.Equal(t, events.ObjectStatusBooked, objects.Objects["A-1"].Status)
	require.Equal(t, events.ObjectStatusBooked, objects.Objects["A-2"].Status)
	require.Equal(t, events.ObjectStatusBooked, event1ObjectInfos["A-1"].Status)
	require.Equal(t, events.ObjectStatusBooked, event1ObjectInfos["A-2"].Status)
	require.Equal(t, events.ObjectStatusBooked, event2ObjectInfos["A-1"].Status)
	require.Equal(t, events.ObjectStatusBooked, event2ObjectInfos["A-2"].Status)
}
