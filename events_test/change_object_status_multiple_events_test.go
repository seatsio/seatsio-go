package events_test

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/shared"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatusChange(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	results, err := client.Events.ChangeObjectStatusWithOptions(
		test_util.RequestContext(),
		&events.StatusChangeParams{
			Events: []string{event1.Key, event2.Key},
			StatusChanges: events.StatusChanges{
				Status:  "foo",
				Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, "foo", string(results.Objects["A-1"].Status))
	require.Equal(t, "foo", string(results.Objects["A-2"].Status))

	event1Data, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2Data, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, "foo", string(event1Data["A-1"].Status))
	require.Equal(t, "foo", string(event1Data["A-2"].Status))
	require.Equal(t, "foo", string(event2Data["A-1"].Status))
	require.Equal(t, "foo", string(event2Data["A-2"].Status))
}

func TestBook(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		},
	})
	require.NoError(t, err)

	event1ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, objects.Objects["A-1"].Status)
	require.Equal(t, events.BOOKED, objects.Objects["A-2"].Status)
	require.Equal(t, events.BOOKED, event1ObjectInfos["A-1"].Status)
	require.Equal(t, events.BOOKED, event1ObjectInfos["A-2"].Status)
	require.Equal(t, events.BOOKED, event2ObjectInfos["A-1"].Status)
	require.Equal(t, events.BOOKED, event2ObjectInfos["A-2"].Status)
}

func TestHold(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:    events.HELD,
			HoldToken: holdToken.HoldToken,
			Objects:   []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		},
	})
	require.NoError(t, err)

	event1ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.HELD, objects.Objects["A-1"].Status)
	require.Equal(t, events.HELD, objects.Objects["A-2"].Status)
	require.Equal(t, events.HELD, event1ObjectInfos["A-1"].Status)
	require.Equal(t, events.HELD, event1ObjectInfos["A-2"].Status)
	require.Equal(t, events.HELD, event2ObjectInfos["A-1"].Status)
	require.Equal(t, events.HELD, event2ObjectInfos["A-2"].Status)
}

func TestPutUpForResale(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.RESALE,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		},
	})
	require.NoError(t, err)

	event1ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2ObjectInfos, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.RESALE, objects.Objects["A-1"].Status)
	require.Equal(t, events.RESALE, objects.Objects["A-2"].Status)
	require.Equal(t, events.RESALE, event1ObjectInfos["A-1"].Status)
	require.Equal(t, events.RESALE, event1ObjectInfos["A-2"].Status)
	require.Equal(t, events.RESALE, event2ObjectInfos["A-1"].Status)
	require.Equal(t, events.RESALE, event2ObjectInfos["A-2"].Status)
}

func TestRelease(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	objects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		},
	})
	require.NoError(t, err)

	event1ObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2ObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.BOOKED, objects.Objects["A-1"].Status)
	require.Equal(t, events.BOOKED, objects.Objects["A-2"].Status)
	require.Equal(t, events.BOOKED, event1ObjectInfo["A-1"].Status)
	require.Equal(t, events.BOOKED, event1ObjectInfo["A-2"].Status)
	require.Equal(t, events.BOOKED, event2ObjectInfo["A-1"].Status)
	require.Equal(t, events.BOOKED, event2ObjectInfo["A-2"].Status)

	releasedObjects, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.FREE,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}, {ObjectId: "A-2"}},
		},
	})
	require.NoError(t, err)

	event1ReleasedObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1", "A-2")
	require.NoError(t, err)
	event2ReleasedObjectInfo, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1", "A-2")
	require.NoError(t, err)

	require.Equal(t, events.FREE, releasedObjects.Objects["A-1"].Status)
	require.Equal(t, events.FREE, releasedObjects.Objects["A-2"].Status)
	require.Equal(t, events.FREE, event1ReleasedObjectInfo["A-1"].Status)
	require.Equal(t, events.FREE, event1ReleasedObjectInfo["A-2"].Status)
	require.Equal(t, events.FREE, event2ReleasedObjectInfo["A-1"].Status)
	require.Equal(t, events.FREE, event2ReleasedObjectInfo["A-2"].Status)
}

func TestAllowedPreviousStatuses(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:                  events.BOOKED,
			Objects:                 []events.ObjectProperties{{ObjectId: "A-1"}},
			AllowedPreviousStatuses: []string{"MustBeThisStatus"},
		},
	})
	seatsioErr := err.(*shared.SeatsioError)
	require.Equal(t, "ILLEGAL_STATUS_CHANGE", seatsioErr.Code)
}

func TestRejectedPreviousStatuses(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	_, err = client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event1.Key, event2.Key},
		StatusChanges: events.StatusChanges{
			Status:                   events.BOOKED,
			Objects:                  []events.ObjectProperties{{ObjectId: "A-1"}},
			RejectedPreviousStatuses: []string{events.FREE},
		},
	})
	seatsioErr := err.(*shared.SeatsioError)
	require.Equal(t, "ILLEGAL_STATUS_CHANGE", seatsioErr.Code)
}

func TestResaleListingId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event1, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)
	event2, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	results, err := client.Events.ChangeObjectStatusWithOptions(
		test_util.RequestContext(),
		&events.StatusChangeParams{
			Events: []string{event1.Key, event2.Key},
			StatusChanges: events.StatusChanges{
				Status:          events.RESALE,
				Objects:         []events.ObjectProperties{{ObjectId: "A-1"}},
				ResaleListingId: "listing1",
			},
		},
	)
	require.NoError(t, err)

	require.Equal(t, "listing1", string(results.Objects["A-1"].ResaleListingId))

	event1Data, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event1.Key, "A-1")
	require.NoError(t, err)
	event2Data, err := client.Events.RetrieveObjectInfo(test_util.RequestContext(), event2.Key, "A-1")
	require.NoError(t, err)

	require.Equal(t, "listing1", string(event1Data["A-1"].ResaleListingId))
	require.Equal(t, "listing1", string(event2Data["A-1"].ResaleListingId))
}
