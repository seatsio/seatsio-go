package reports

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/reports"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDetailedReportItemProperties(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{{
				ObjectId:   "A-1",
				ExtraData:  events.ExtraData{"foo": "bar"},
				TicketType: "ticketType1",
			}},
			OrderId:        "order1",
			IgnoreChannels: true,
		},
	})
	require.NoError(t, err)

	report, err := client.EventReports.ByLabel(event.Key)
	require.NoError(t, err)

	reportItem := report.Items["A-1"][0]
	require.Equal(t, reports.Booked, reportItem.Status)
	require.Equal(t, "A-1", reportItem.Label)
	require.Equal(t, reports.Labels{
		Own: reports.LabelAndType{
			Label:     "1",
			LabelType: "seat",
		},
		Parent: reports.LabelAndType{
			Label:     "A",
			LabelType: "row",
		},
	}, reportItem.Labels)
	require.Equal(t, reports.IDs{
		Own:    "1",
		Parent: "A",
	}, reportItem.IDs)
	require.Equal(t, "Cat1", reportItem.CategoryLabel)
	require.Equal(t, "9", reportItem.CategoryKey)
	require.Equal(t, "ticketType1", reportItem.TicketType)
	require.Equal(t, "order1", reportItem.OrderId)
	require.True(t, reportItem.ForSale)
	require.Empty(t, reportItem.Section)
	require.Empty(t, reportItem.Entrance)
	require.Empty(t, reportItem.NumBooked)
	require.Empty(t, reportItem.Capacity)
	require.Equal(t, "seat", reportItem.ObjectType)
	require.Equal(t, "bar", reportItem.ExtraData["foo"])
	require.False(t, reportItem.HasRestrictedView)
	require.False(t, reportItem.IsAccessible)
	require.False(t, reportItem.IsCompanionSeat)
	require.Empty(t, reportItem.DisplayedObjectType)
	require.Empty(t, reportItem.LeftNeighbour)
	require.Equal(t, "A-2", reportItem.RightNeighbour)
	require.False(t, reportItem.IsAvailable)
	require.Equal(t, reports.Booked, reportItem.AvailabilityReason)
	require.Equal(t, "channel1", reportItem.Channel)
	require.NotNil(t, reportItem.DistanceToFocalPoint)

	gaItem := report.Items["GA1"][0]
	require.True(t, gaItem.VariableOccupancy)
	require.Equal(t, int32(1), gaItem.MinOccupancy)
	require.Equal(t, int32(100), gaItem.MaxOccupancy)
}

func TestHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	holdToken, err := client.HoldTokens.Create()
	_, err = client.Events.Hold(event.Key, []string{"A-1"}, &holdToken.HoldToken)

	report, err := client.EventReports.ByLabel(event.Key)
	require.NoError(t, err)
	require.Equal(t, holdToken.HoldToken, report.Items["A-1"][0].HoldToken)
}

func TestDetailedReportItemPropertiesForGA(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{{
				ObjectId: "GA1",
				Quantity: 5,
			}},
		},
	})
	holdToken, err := client.HoldTokens.Create()
	_, err = client.Events.HoldWithObjectProperties(
		event.Key,
		[]events.ObjectProperties{
			{
				ObjectId: "GA1",
				Quantity: 3,
			},
		},
		&holdToken.HoldToken)

	report, err := client.EventReports.ByLabel(event.Key)
	require.NoError(t, err)

	reportItem := report.Items["GA1"][0]
	require.Equal(t, int32(5), reportItem.NumBooked)
	require.Equal(t, int32(92), reportItem.NumFree)
	require.Equal(t, int32(3), reportItem.NumHeld)
	require.Equal(t, int32(100), reportItem.Capacity)
	require.Equal(t, "generalAdmission", reportItem.ObjectType)
	require.False(t, reportItem.BookAsAWhole)
	require.Empty(t, reportItem.HasRestrictedView)
	require.Empty(t, reportItem.IsAccessible)
	require.Empty(t, reportItem.IsCompanionSeat)
	require.Empty(t, reportItem.DisplayedObjectType)
}

func TestDetailedReportItemPropertiesForTable(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{
		ChartKey: chartKey,
		EventParams: &events.EventParams{
			TableBookingConfig: events.TableBookingSupport.AllByTables(),
		},
	})

	report, err := client.EventReports.ByLabel(event.Key)
	require.NoError(t, err)

	reportItem := report.Items["T1"][0]
	require.False(t, reportItem.BookAsAWhole)
	require.Equal(t, int32(6), reportItem.NumSeats)
}

func TestByStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatusWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: "lolzor",
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
		},
	})
	_, err = client.Events.Book(event.Key, "A-3")

	report, err := client.EventReports.ByStatus(event.Key)
	require.NoError(t, err)

	require.Equal(t, 2, len(report.Items["lolzor"]))
	require.Equal(t, 1, len(report.Items[reports.Booked]))
	require.Equal(t, 31, len(report.Items[reports.Free]))
}

func TestByStatusWithEmptyChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(&charts.CreateChartParams{})
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chart.Key})

	report, err := client.EventReports.ByStatus(event.Key)
	require.NoError(t, err)

	require.Empty(t, report.Items)
}

func TestBySpecificStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatusWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: "lolzor",
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
		},
	})
	_, err = client.Events.Book(event.Key, "A-3")

	items, err := client.EventReports.BySpecificStatus(event.Key, "lolzor")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestBySpecificNonExistingStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificStatus(event.Key, "lolzor")
	require.NoError(t, err)

	require.Empty(t, len(items))
}

func TestDetailedReportByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByObjectType(event.Key)
	require.NoError(t, err)

	require.Equal(t, 32, len(report.Items["seat"]))
	require.Equal(t, 2, len(report.Items["generalAdmission"]))
	require.Equal(t, 0, len(report.Items["booth"]))
}

func TestDetailedReportBySpecificObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificObjectType(event.Key, "seat")
	require.NoError(t, err)

	require.Equal(t, 32, len(items))
}

func TestDetailedReportByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByCategoryLabel(event.Key)
	require.NoError(t, err)

	require.Equal(t, 17, len(report.Items["Cat1"]))
	require.Equal(t, 17, len(report.Items["Cat2"]))
}

func TestDetailedReportBySpecificCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificCategoryLabel(event.Key, "Cat1")
	require.NoError(t, err)

	require.Equal(t, 17, len(items))
}

func TestDetailedReportByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByCategoryKey(event.Key)
	require.NoError(t, err)

	require.Equal(t, 17, len(report.Items["9"]))
	require.Equal(t, 17, len(report.Items["10"]))
}

func TestDetailedReportBySpecificCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificCategoryKey(event.Key, "9")
	require.NoError(t, err)

	require.Equal(t, 17, len(items))
}

func TestDetailedReportByLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByLabel(event.Key)
	require.NoError(t, err)

	require.Equal(t, 1, len(report.Items["A-1"]))
	require.Equal(t, 1, len(report.Items["A-2"]))
}

func TestDetailedReportBySpecificLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificLabel(event.Key, "A-1")
	require.NoError(t, err)

	require.Equal(t, 1, len(items))
}

func TestDetailedReportByOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
			OrderId: "order1",
		},
	})
	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-3"},
			},
			OrderId: "order2",
		},
	})

	report, err := client.EventReports.ByOrderId(event.Key)
	require.NoError(t, err)

	require.Equal(t, 2, len(report.Items["order1"]))
	require.Equal(t, 1, len(report.Items["order2"]))
	require.Equal(t, 31, len(report.Items[reports.NoOrderId]))
}

func TestDetailedReportBySpecificOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
			OrderId: "order1",
		},
	})
	_, err = client.Events.BookWithOptions(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-3"},
			},
			OrderId: "order2",
		},
	})

	items, err := client.EventReports.BySpecificOrderId(event.Key, "order1")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.BySection(event.Key)
	require.NoError(t, err)

	require.Equal(t, 34, len(report.Items[reports.NoSection]))
}

func TestDetailedReportBySpecificSection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificSection(event.Key, reports.NoSection)
	require.NoError(t, err)

	require.Equal(t, 34, len(items))
}

func TestDetailedReportByAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	report, err := client.EventReports.ByAvailability(event.Key)
	require.NoError(t, err)

	require.Equal(t, 32, len(report.Items[reports.Available]))
	require.Equal(t, 2, len(report.Items[reports.NotAvailable]))
}

func TestDetailedReportBySpecificAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	items, err := client.EventReports.BySpecificAvailability(event.Key, reports.NotAvailable)
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportByAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	report, err := client.EventReports.ByAvailabilityReason(event.Key)
	require.NoError(t, err)

	require.Equal(t, 32, len(report.Items[reports.Available]))
	require.Equal(t, 2, len(report.Items["lolzor"]))
}

func TestDetailedReportBySpecificAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	items, err := client.EventReports.BySpecificAvailabilityReason(event.Key, "lolzor")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportByChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1"}},
		},
	}})
	require.NoError(t, err)

	report, err := client.EventReports.ByChannel(event.Key)
	require.NoError(t, err)

	require.Equal(t, 2, len(report.Items["channel1"]))
	require.Equal(t, 32, len(report.Items[reports.NoChannel]))
}

func TestDetailedReportBySpecificChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	items, err := client.EventReports.BySpecificChannel(event.Key, "channel1")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}
