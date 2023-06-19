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
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	client.Events.ChangeObjectStatus(&events.StatusChangeParams{
		Status:  events.ObjectStatusBooked,
		Events:  []string{event.Key},
		Objects: []events.ObjectProperties{{ObjectId: "A-1"}},
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
