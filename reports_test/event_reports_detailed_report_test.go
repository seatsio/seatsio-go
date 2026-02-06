package reports

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/charts"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/reports"
	"github.com/seatsio/seatsio-go/v12/seasons"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestDetailedReportItemProperties(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1"}},
		},
	}})
	require.NoError(t, err)

	_, err = client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
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

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	reportItem := report.Items["A-1"][0]
	require.Equal(t, events.BOOKED, reportItem.Status)
	require.Equal(t, "A-1", reportItem.Label)
	require.Equal(t, events.Labels{
		Own: events.LabelAndType{
			Label: "1",
			Type:  "seat",
		},
		Parent: events.LabelAndType{
			Label: "A",
			Type:  "row",
		},
	}, reportItem.Labels)
	require.Equal(t, events.IDs{
		Own:    "1",
		Parent: "A",
	}, reportItem.IDs)
	require.Equal(t, "Cat1", reportItem.CategoryLabel)
	require.Equal(t, events.CategoryKey{"9"}, reportItem.CategoryKey)
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
	require.False(t, reportItem.HasLiftUpArmrests)
	require.False(t, reportItem.IsHearingImpaired)
	require.False(t, reportItem.IsSemiAmbulatorySeat)
	require.False(t, reportItem.HasSignLanguageInterpretation)
	require.False(t, reportItem.IsPlusSize)
	require.Empty(t, reportItem.DisplayedObjectType)
	require.Empty(t, reportItem.ParentDisplayedObjectType)
	require.Empty(t, reportItem.LeftNeighbour)
	require.Equal(t, "A-2", reportItem.RightNeighbour)
	require.False(t, reportItem.IsAvailable)
	require.Equal(t, events.BOOKED, reportItem.AvailabilityReason)
	require.Equal(t, "channel1", reportItem.Channel)
	require.NotNil(t, reportItem.DistanceToFocalPoint)
	require.Equal(t, 0, reportItem.SeasonStatusOverriddenQuantity)
	require.Empty(t, reportItem.ResaleListingId)

	gaItem := report.Items["GA1"][0]
	require.True(t, gaItem.VariableOccupancy)
	require.Equal(t, 1, gaItem.MinOccupancy)
	require.Equal(t, 100, gaItem.MaxOccupancy)
}

func TestHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	_, err = client.Events.Hold(test_util.RequestContext(), event.Key, []string{"A-1"}, &holdToken.HoldToken)

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, holdToken.HoldToken, report.Items["A-1"][0].HoldToken)
}

func TestSeasonStatusOverriddenQuantity(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	season, err := client.Seasons.CreateWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{NumberOfEvents: 1})
	require.NoError(t, err)

	event := season.Events[0]
	err = client.Events.OverrideSeasonObjectStatus(test_util.RequestContext(), event.Key, "A-1")
	require.NoError(t, err)

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 1, report.Items["A-1"][0].SeasonStatusOverriddenQuantity)
}

func TestDetailedReportItemPropertiesForGA(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{{
				ObjectId: "GA1",
				Quantity: 5,
			}},
		},
	})
	holdToken, err := client.HoldTokens.Create(test_util.RequestContext())
	_, err = client.Events.HoldWithObjectProperties(
		test_util.RequestContext(),
		event.Key,
		[]events.ObjectProperties{
			{
				ObjectId: "GA1",
				Quantity: 3,
			},
		},
		&holdToken.HoldToken)

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	reportItem := report.Items["GA1"][0]
	require.Equal(t, 5, reportItem.NumBooked)
	require.Equal(t, 92, reportItem.NumFree)
	require.Equal(t, 3, reportItem.NumHeld)
	require.Equal(t, 100, reportItem.Capacity)
	require.Equal(t, "generalAdmission", reportItem.ObjectType)
	require.False(t, reportItem.BookAsAWhole)
	require.Empty(t, reportItem.HasRestrictedView)
	require.Empty(t, reportItem.IsAccessible)
	require.Empty(t, reportItem.IsCompanionSeat)
	require.Empty(t, reportItem.HasLiftUpArmrests)
	require.Empty(t, reportItem.IsHearingImpaired)
	require.Empty(t, reportItem.IsSemiAmbulatorySeat)
	require.Empty(t, reportItem.HasSignLanguageInterpretation)
	require.Empty(t, reportItem.IsPlusSize)
	require.Empty(t, reportItem.DisplayedObjectType)
	require.Empty(t, reportItem.ParentDisplayedObjectType)
}

func TestDetailedReportItemPropertiesForTable(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{
		ChartKey: chartKey,
		EventParams: &events.EventParams{
			TableBookingConfig: events.TableBookingSupport.AllByTables(),
		},
	})

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	reportItem := report.Items["T1"][0]
	require.False(t, reportItem.BookAsAWhole)
	require.Equal(t, 6, reportItem.NumSeats)
}

func TestByStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: "lolzor",
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
		},
	})
	_, err = client.Events.Book(test_util.RequestContext(), event.Key, "A-3")

	report, err := client.EventReports.ByStatus(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 2, len(report.Items["lolzor"]))
	require.Equal(t, 1, len(report.Items[events.BOOKED]))
	require.Equal(t, 31, len(report.Items[events.FREE]))
}

func TestByStatusWithEmptyChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{})
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chart.Key})

	report, err := client.EventReports.ByStatus(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Empty(t, report.Items)
}

func TestBySpecificStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatusWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status: "lolzor",
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
		},
	})
	_, err = client.Events.Book(test_util.RequestContext(), event.Key, "A-3")

	items, err := client.EventReports.BySpecificStatus(test_util.RequestContext(), event.Key, "lolzor")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestBySpecificNonExistingStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificStatus(test_util.RequestContext(), event.Key, "lolzor")
	require.NoError(t, err)

	require.Empty(t, len(items))
}

func TestDetailedReportByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByObjectType(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificObjectType(test_util.RequestContext(), event.Key, "seat")
	require.NoError(t, err)

	require.Equal(t, 32, len(items))
}

func TestDetailedReportByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByCategoryLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 17, len(report.Items["Cat1"]))
	require.Equal(t, 17, len(report.Items["Cat2"]))
}

func TestDetailedReportBySpecificCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificCategoryLabel(test_util.RequestContext(), event.Key, "Cat1")
	require.NoError(t, err)

	require.Equal(t, 17, len(items))
}

func TestDetailedReportByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByCategoryKey(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 17, len(report.Items["9"]))
	require.Equal(t, 17, len(report.Items["10"]))
}

func TestDetailedReportBySpecificCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificCategoryKey(test_util.RequestContext(), event.Key, "9")
	require.NoError(t, err)

	require.Equal(t, 17, len(items))
}

func TestDetailedReportByLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 1, len(report.Items["A-1"]))
	require.Equal(t, 1, len(report.Items["A-2"]))
}

func TestDetailedReportBySpecificLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificLabel(test_util.RequestContext(), event.Key, "A-1")
	require.NoError(t, err)

	require.Equal(t, 1, len(items))
}

func TestDetailedReportByOrderId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
			OrderId: "order1",
		},
	})
	_, err = client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-3"},
			},
			OrderId: "order2",
		},
	})

	report, err := client.EventReports.ByOrderId(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-1"},
				{ObjectId: "A-2"},
			},
			OrderId: "order1",
		},
	})
	_, err = client.Events.BookWithOptions(test_util.RequestContext(), &events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Objects: []events.ObjectProperties{
				{ObjectId: "A-3"},
			},
			OrderId: "order2",
		},
	})

	items, err := client.EventReports.BySpecificOrderId(test_util.RequestContext(), event.Key, "order1")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.BySection(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 34, len(report.Items[reports.NoSection]))
}

func TestDetailedReportBySpecificSection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificSection(test_util.RequestContext(), event.Key, reports.NoSection)
	require.NoError(t, err)

	require.Equal(t, 34, len(items))
}

func TestDetailedReportByZone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.ByZone(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 6032, len(report.Items["midtrack"]))
	require.Equal(t, "midtrack", report.Items["midtrack"][0].Zone)
}

func TestDetailedReportBySpecificZone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	items, err := client.EventReports.BySpecificZone(test_util.RequestContext(), event.Key, "midtrack")
	require.NoError(t, err)

	require.Equal(t, 6032, len(items))
}

func TestDetailedReportByAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	report, err := client.EventReports.ByAvailability(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 32, len(report.Items[reports.Available]))
	require.Equal(t, 2, len(report.Items[reports.NotAvailable]))
}

func TestDetailedReportBySpecificAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	items, err := client.EventReports.BySpecificAvailability(test_util.RequestContext(), event.Key, reports.NotAvailable)
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportByAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	report, err := client.EventReports.ByAvailabilityReason(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 32, len(report.Items[reports.Available]))
	require.Equal(t, 2, len(report.Items["lolzor"]))
}

func TestDetailedReportBySpecificAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, err := client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1", "A-2"}, "lolzor")

	items, err := client.EventReports.BySpecificAvailabilityReason(test_util.RequestContext(), event.Key, "lolzor")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestDetailedReportByChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	report, err := client.EventReports.ByChannel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)

	require.Equal(t, 2, len(report.Items["channel1"]))
	require.Equal(t, 32, len(report.Items[reports.NoChannel]))
}

func TestDetailedReportBySpecificChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		Channels: &[]events.CreateChannelParams{
			{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}},
		},
	}})
	require.NoError(t, err)

	items, err := client.EventReports.BySpecificChannel(test_util.RequestContext(), event.Key, "channel1")
	require.NoError(t, err)

	require.Equal(t, 2, len(items))
}

func TestResaleListingId(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, err := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	require.NoError(t, err)

	listingId := "listing1"
	_, err = client.Events.PutUpForResale(test_util.RequestContext(), event.Key, []string{"A-1"}, &listingId)
	require.NoError(t, err)

	report, err := client.EventReports.ByLabel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, "listing1", report.Items["A-1"][0].ResaleListingId)
}
