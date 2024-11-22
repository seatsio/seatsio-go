package events_test

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/seasons"
	"github.com/seatsio/seatsio-go/v9/shared"
	"github.com/seatsio/seatsio-go/v9/test_util"
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
				Type:    events.CHANGE_STATUS_TO,
				Status:  "lolzor",
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
			},
		},
		events.StatusChangeInBatchParams{
			Event: event2.Key,
			StatusChanges: events.StatusChanges{
				Type:    events.CHANGE_STATUS_TO,
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

	require.Equal(t, "lolzor", result.Results[0].Objects["A-1"].Status)
	require.Equal(t, "lolzor", event1Info["A-1"].Status)
	require.Equal(t, "lolzor", result.Results[0].Objects["A-1"].Status)
	require.Equal(t, "lolzor", event2Info["A-2"].Status)
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
	require.Equal(t, "lolzor", result.Results[0].Objects["A-1"].Status)

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

func TestReleaseInBatch(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	_, err = client.Events.Book(event.Key, "A-1")
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{
			Event: event.Key,
			StatusChanges: events.StatusChanges{
				Type:    events.RELEASE,
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
			},
		},
	)
	require.NoError(t, err)

	eventInfo, err := client.Events.RetrieveObjectInfo(event.Key, "A-1")
	require.NoError(t, err)
	require.Equal(t, events.FREE, result.Results[0].Objects["A-1"].Status)
	require.Equal(t, events.FREE, eventInfo["A-1"].Status)
}

func TestOverrideSeasonStatusInBatch(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.CreateSeasonWithOptions(chartKey, &seasons.CreateSeasonParams{EventKeys: []string{"event1"}})
	require.NoError(t, err)
	_, err = client.Events.Book(season.Key, "A-1")
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{
			Event: "event1",
			StatusChanges: events.StatusChanges{
				Type:    events.OVERRIDE_SEASON_STATUS,
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, events.FREE, result.Results[0].Objects["A-1"].Status)
	info, _ := client.Events.RetrieveObjectInfo("event1", "A-1")
	require.Equal(t, events.FREE, info["A-1"].Status)
}

func TestUseSeasonStatusInBatch(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.CreateSeasonWithOptions(chartKey, &seasons.CreateSeasonParams{EventKeys: []string{"event1"}})
	require.NoError(t, err)
	_, err = client.Events.Book(season.Key, "A-1")
	require.NoError(t, err)
	err = client.Events.OverrideSeasonObjectStatus("event1", "A-1")
	require.NoError(t, err)

	result, err := client.Events.ChangeObjectStatusInBatch(
		events.StatusChangeInBatchParams{
			Event: "event1",
			StatusChanges: events.StatusChanges{
				Type:    events.USE_SEASON_STATUS,
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, result.Results[0].Objects["A-1"].Status)
	info, _ := client.Events.RetrieveObjectInfo("event1", "A-1")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
}
