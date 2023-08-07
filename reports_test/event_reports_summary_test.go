package reports

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/reports"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSummaryByStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryByStatus(event.Key)
	require.NoError(t, err)
	bookedReportItem := reports.EventSummaryReportItem{
		Count:                1,
		BySection:            map[string]int{"NO_SECTION": 1},
		ByCategoryKey:        map[string]int{"9": 1},
		ByCategoryLabel:      map[string]int{"Cat1": 1},
		ByAvailability:       map[string]int{"not_available": 1},
		ByAvailabilityReason: map[string]int{"booked": 1},
		ByChannel:            map[string]int{"NO_CHANNEL": 1},
	}
	freeReportItem := reports.EventSummaryReportItem{
		Count:     231,
		BySection: map[string]int{"NO_SECTION": 231},
		ByCategoryKey: map[string]int{
			"9":  115,
			"10": 116,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 115,
			"Cat2": 116,
		},
		ByAvailability:       map[string]int{"available": 231},
		ByAvailabilityReason: map[string]int{"available": 231},
		ByChannel:            map[string]int{"NO_CHANNEL": 231},
	}
	require.Equal(t, bookedReportItem, report.Items["booked"])
	require.Equal(t, freeReportItem, report.Items["free"])
}

func TestSummaryByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.SummaryByObjectType(event.Key)
	require.NoError(t, err)

	seatReport := reports.EventSummaryReportItem{
		Count:     32,
		BySection: map[string]int{"NO_SECTION": 32},
		ByCategoryKey: map[string]int{
			"9":  16,
			"10": 16,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 16,
			"Cat2": 16,
		},
		ByAvailability:       map[string]int{"available": 32},
		ByAvailabilityReason: map[string]int{"available": 32},
		ByChannel:            map[string]int{"NO_CHANNEL": 32},
		ByStatus:             map[string]int{"free": 32},
	}
	gaReport := reports.EventSummaryReportItem{
		Count:     200,
		BySection: map[string]int{"NO_SECTION": 200},
		ByCategoryKey: map[string]int{
			"9":  100,
			"10": 100,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 100,
			"Cat2": 100,
		},
		ByAvailability:       map[string]int{"available": 200},
		ByAvailabilityReason: map[string]int{"available": 200},
		ByChannel:            map[string]int{"NO_CHANNEL": 200},
		ByStatus:             map[string]int{"free": 200},
	}
	emptyReport := reports.EventSummaryReportItem{
		Count:                0,
		BySection:            map[string]int{},
		ByCategoryKey:        map[string]int{},
		ByCategoryLabel:      map[string]int{},
		ByAvailability:       map[string]int{},
		ByAvailabilityReason: map[string]int{},
		ByChannel:            map[string]int{},
		ByStatus:             map[string]int{},
	}
	require.Equal(t, seatReport, report.Items["seat"])
	require.Equal(t, gaReport, report.Items["generalAdmission"])
	require.Equal(t, emptyReport, report.Items["booth"])
	require.Equal(t, emptyReport, report.Items["table"])
}

func TestSummaryByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryByCategoryKey(event.Key)
	require.NoError(t, err)

	cat9Report := reports.EventSummaryReportItem{
		Count:                116,
		BySection:            map[string]int{"NO_SECTION": 116},
		ByStatus:             map[string]int{"booked": 1, "free": 115},
		ByAvailability:       map[string]int{"available": 115, "not_available": 1},
		ByAvailabilityReason: map[string]int{"available": 115, "booked": 1},
		ByChannel:            map[string]int{"NO_CHANNEL": 116},
	}
	cat10Report := reports.EventSummaryReportItem{
		Count:                116,
		BySection:            map[string]int{"NO_SECTION": 116},
		ByStatus:             map[string]int{"free": 116},
		ByAvailability:       map[string]int{"available": 116},
		ByAvailabilityReason: map[string]int{"available": 116},
		ByChannel:            map[string]int{"NO_CHANNEL": 116},
	}
	cat11Report := reports.EventSummaryReportItem{
		Count:                0,
		BySection:            map[string]int{},
		ByStatus:             map[string]int{},
		ByAvailability:       map[string]int{},
		ByAvailabilityReason: map[string]int{},
		ByChannel:            map[string]int{},
	}
	noCategoryReport := reports.EventSummaryReportItem{
		Count:                0,
		BySection:            map[string]int{},
		ByStatus:             map[string]int{},
		ByAvailability:       map[string]int{},
		ByAvailabilityReason: map[string]int{},
		ByChannel:            map[string]int{},
	}
	require.Equal(t, cat9Report, report.Items["9"])
	require.Equal(t, cat10Report, report.Items["10"])
	require.Equal(t, cat11Report, report.Items["string11"])
	require.Equal(t, noCategoryReport, report.Items["NO_CATEGORY"])
}

func TestSummaryByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryByCategoryLabel(event.Key)
	require.NoError(t, err)

	cat1Report := reports.EventSummaryReportItem{
		Count:                116,
		BySection:            map[string]int{"NO_SECTION": 116},
		ByStatus:             map[string]int{"booked": 1, "free": 115},
		ByAvailability:       map[string]int{"available": 115, "not_available": 1},
		ByAvailabilityReason: map[string]int{"available": 115, "booked": 1},
		ByChannel:            map[string]int{"NO_CHANNEL": 116},
	}
	cat2Report := reports.EventSummaryReportItem{
		Count:                116,
		BySection:            map[string]int{"NO_SECTION": 116},
		ByStatus:             map[string]int{"free": 116},
		ByAvailability:       map[string]int{"available": 116},
		ByAvailabilityReason: map[string]int{"available": 116},
		ByChannel:            map[string]int{"NO_CHANNEL": 116},
	}
	cat3Report := reports.EventSummaryReportItem{
		Count:                0,
		BySection:            map[string]int{},
		ByStatus:             map[string]int{},
		ByAvailability:       map[string]int{},
		ByAvailabilityReason: map[string]int{},
		ByChannel:            map[string]int{},
	}
	noCategoryReport := reports.EventSummaryReportItem{
		Count:                0,
		BySection:            map[string]int{},
		ByStatus:             map[string]int{},
		ByAvailability:       map[string]int{},
		ByAvailabilityReason: map[string]int{},
		ByChannel:            map[string]int{},
	}
	require.Equal(t, cat1Report, report.Items["Cat1"])
	require.Equal(t, cat2Report, report.Items["Cat2"])
	require.Equal(t, cat3Report, report.Items["Cat3"])
	require.Equal(t, noCategoryReport, report.Items["NO_CATEGORY"])
}

func TestSummaryBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryBySection(event.Key)
	require.NoError(t, err)

	noSectionReport := reports.EventSummaryReportItem{
		Count:    232,
		ByStatus: map[string]int{"booked": 1, "free": 231},
		ByCategoryKey: map[string]int{
			"9":  116,
			"10": 116,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 116,
			"Cat2": 116,
		},
		ByAvailability:       map[string]int{"available": 231, "not_available": 1},
		ByAvailabilityReason: map[string]int{"available": 231, "booked": 1},
		ByChannel:            map[string]int{"NO_CHANNEL": 232},
	}
	require.Equal(t, noSectionReport, report.Items["NO_SECTION"])
}

func TestSummaryByAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryByAvailability(event.Key)
	require.NoError(t, err)

	availableReport := reports.EventSummaryReportItem{
		Count:     231,
		BySection: map[string]int{"NO_SECTION": 231},
		ByStatus:  map[string]int{"free": 231},
		ByCategoryKey: map[string]int{
			"9":  115,
			"10": 116,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 115,
			"Cat2": 116,
		},
		ByChannel:            map[string]int{"NO_CHANNEL": 231},
		ByAvailabilityReason: map[string]int{"available": 231},
	}

	notavailableReport := reports.EventSummaryReportItem{
		Count:                1,
		BySection:            map[string]int{"NO_SECTION": 1},
		ByStatus:             map[string]int{"booked": 1},
		ByCategoryKey:        map[string]int{"9": 1},
		ByCategoryLabel:      map[string]int{"Cat1": 1},
		ByChannel:            map[string]int{"NO_CHANNEL": 1},
		ByAvailabilityReason: map[string]int{"booked": 1},
	}
	require.Equal(t, availableReport, report.Items["available"])
	require.Equal(t, notavailableReport, report.Items["not_available"])
}

func TestSummaryByAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Events: []string{event.Key},
		StatusChanges: events.StatusChanges{
			Status:  events.BOOKED,
			Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
		},
	})

	report, err := client.EventReports.SummaryByAvailabilityReason(event.Key)
	require.NoError(t, err)

	availableReport := reports.EventSummaryReportItem{
		Count:     231,
		BySection: map[string]int{"NO_SECTION": 231},
		ByStatus:  map[string]int{"free": 231},
		ByCategoryKey: map[string]int{
			"9":  115,
			"10": 116,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 115,
			"Cat2": 116,
		},
		ByChannel:      map[string]int{"NO_CHANNEL": 231},
		ByAvailability: map[string]int{"available": 231},
	}

	bookedReport := reports.EventSummaryReportItem{
		Count:           1,
		BySection:       map[string]int{"NO_SECTION": 1},
		ByStatus:        map[string]int{"booked": 1},
		ByCategoryKey:   map[string]int{"9": 1},
		ByCategoryLabel: map[string]int{"Cat1": 1},
		ByChannel:       map[string]int{"NO_CHANNEL": 1},
		ByAvailability:  map[string]int{"not_available": 1},
	}
	emptyReport := reports.EventSummaryReportItem{
		BySection:       map[string]int{},
		ByStatus:        map[string]int{},
		ByCategoryKey:   map[string]int{},
		ByCategoryLabel: map[string]int{},
		ByChannel:       map[string]int{},
		ByAvailability:  map[string]int{},
	}
	require.Equal(t, availableReport, report.Items["available"])
	require.Equal(t, bookedReport, report.Items["booked"])
	require.Equal(t, emptyReport, report.Items["reservedByToken"])
	require.Equal(t, emptyReport, report.Items["disabled_by_social_distancing"])
	require.Equal(t, emptyReport, report.Items["not_for_sale"])
}

func TestSummaryByChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	err := client.Channels.Replace(event.Key, events.Channel{Key: "channel1", Name: "channel 1", Color: "#FFFF99", Index: 1, Objects: []string{"A-1", "A-2"}})

	report, err := client.EventReports.SummaryByChannel(event.Key)
	require.NoError(t, err)

	channelReport := reports.EventSummaryReportItem{
		Count:                2,
		BySection:            map[string]int{"NO_SECTION": 2},
		ByStatus:             map[string]int{string(events.FREE): 2},
		ByCategoryKey:        map[string]int{"9": 2},
		ByCategoryLabel:      map[string]int{"Cat1": 2},
		ByAvailability:       map[string]int{"available": 2},
		ByAvailabilityReason: map[string]int{"available": 2},
	}
	noChannelReport := reports.EventSummaryReportItem{
		Count:     230,
		BySection: map[string]int{"NO_SECTION": 230},
		ByStatus:  map[string]int{string(events.FREE): 230},
		ByCategoryKey: map[string]int{
			"9":  114,
			"10": 116,
		},
		ByCategoryLabel: map[string]int{
			"Cat1": 114,
			"Cat2": 116,
		},
		ByAvailability:       map[string]int{"available": 230},
		ByAvailabilityReason: map[string]int{"available": 230},
	}

	require.Equal(t, map[string]reports.EventSummaryReportItem{
		"channel1":   channelReport,
		"NO_CHANNEL": noChannelReport,
	}, report.Items)
}
