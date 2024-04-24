package events_test

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/events"
	"github.com/seatsio/seatsio-go/v7/shared"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeObjectStatusInBatch(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{
			Event: event1.Key,
			StatusChanges: events.StatusChanges{
				Status:  "lolzor",
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
			},
		},
		events.StatusChangeInBatchParams{
			Event: event2.Key,
			StatusChanges: events.StatusChanges{
				Status:  "lolzor",
				Objects: []events.ObjectProperties{{ObjectId: "A-2"}},
			},
		},
	)
	require.NoError(t, err)

	event1Info, err := client.Events.RetrieveObjectInfo(event1.Key, "A-1")
	require.NoError(t, err)
	event2Info, err := client.Events.RetrieveObjectInfo(event2.Key, "A-2")
	require.NoError(t, err)

	require.Equal(t, "lolzor", string(result.Results[0].Objects["A-1"].Status))
	require.Equal(t, "lolzor", string(event1Info["A-1"].Status))
	require.Equal(t, "lolzor", string(result.Results[0].Objects["A-1"].Status))
	require.Equal(t, "lolzor", string(event2Info["A-2"].Status))
}

func TestChannelKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channelKey1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1"}},
		},
	}})
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "lolzor", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}, IgnoreChannels: true}},
	)
	require.NoError(t, err)
	require.Equal(t, "lolzor", string(result.Results[0].Objects["A-1"].Status))

}

func TestBatchAllowedPreviousStatuses(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{
			Event: event.Key,
			StatusChanges: events.StatusChanges{
				Status:                  "lolzor",
				Objects:                 []events.ObjectProperties{{ObjectId: "A-1"}},
				AllowedPreviousStatuses: []string{"MustBeThisStatus"}},
		},
	)
	seatsioError := err.(*shared.SeatsioError)
	require.Equal(t, "ILLEGAL_STATUS_CHANGE", seatsioError.Code)

}

func TestBatchRejectedPreviousStatuses(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{Event: event.Key, StatusChanges: events.StatusChanges{Status: "lolzor", Objects: []events.ObjectProperties{{ObjectId: "A-1"}}, RejectedPreviousStatuses: []string{events.FREE}}},
	)
	seatsioError := err.(*shared.SeatsioError)
	require.Equal(t, "ILLEGAL_STATUS_CHANGE", seatsioError.Code)

}
