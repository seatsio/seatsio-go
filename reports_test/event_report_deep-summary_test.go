package reports

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeepSummaryByStatus(t *testing.T) {
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

	report, err := client.EventReports.DeepSummaryByStatus(event.Key)
	require.NoError(t, err)
	require.Equal(t, 1, report.Items[events.ObjectStatusBooked].Count)
	// require.Equal(t, 1, report.Items[events.ObjectStatusBooked].Count) TODO bver

}
