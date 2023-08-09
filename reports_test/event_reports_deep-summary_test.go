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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByStatus(event.Key)
	require.NoError(t, err)
	require.Equal(t, 1, report.Items["booked"].Count)
	require.Equal(t, 1, report.Items["booked"].BySection["NO_SECTION"].Count)
	require.Equal(t, 1, report.Items["booked"].BySection["NO_SECTION"].ByAvailability["not_available"])

}

func TestDeepSummaryByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.DeepSummaryByObjectType(event.Key)
	require.NoError(t, err)
	require.Equal(t, 32, report.Items["seat"].Count)
	require.Equal(t, 32, report.Items["seat"].BySection["NO_SECTION"].Count)
	require.Equal(t, 32, report.Items["seat"].BySection["NO_SECTION"].ByAvailability["available"])
}

func TestDeepSummaryByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByCategoryKey(event.Key)
	require.NoError(t, err)
	require.Equal(t, 116, report.Items["9"].Count)
	require.Equal(t, 116, report.Items["9"].BySection["NO_SECTION"].Count)
	require.Equal(t, 1, report.Items["9"].BySection["NO_SECTION"].ByAvailability["not_available"])
}

func TestDeepSummaryByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByCategoryLabel(event.Key)
	require.NoError(t, err)
	require.Equal(t, 116, report.Items["Cat1"].Count)
	require.Equal(t, 116, report.Items["Cat1"].BySection["NO_SECTION"].Count)
	require.Equal(t, 1, report.Items["Cat1"].BySection["NO_SECTION"].ByAvailability["not_available"])
}

func TestDeepSummaryBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryBySection(event.Key)
	require.NoError(t, err)
	require.Equal(t, 232, report.Items["NO_SECTION"].Count)
	require.Equal(t, 116, report.Items["NO_SECTION"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 1, report.Items["NO_SECTION"].ByCategoryLabel["Cat1"].ByAvailability["not_available"])
}

func TestDeepSummaryByAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByAvailability(event.Key)
	require.NoError(t, err)
	require.Equal(t, 1, report.Items["not_available"].Count)
	require.Equal(t, 1, report.Items["not_available"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 1, report.Items["not_available"].ByCategoryLabel["Cat1"].BySection["NO_SECTION"])
}

func TestDeepSummaryByAvailabilityReason(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByAvailabilityReason(event.Key)
	require.NoError(t, err)
	require.Equal(t, 1, report.Items["booked"].Count)
	require.Equal(t, 1, report.Items["booked"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 1, report.Items["booked"].ByCategoryLabel["Cat1"].BySection["NO_SECTION"])
}

func TestDeepSummaryByChannel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus([]string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByChannel(event.Key)
	require.NoError(t, err)
	require.Equal(t, 232, report.Items["NO_CHANNEL"].Count)
	require.Equal(t, 116, report.Items["NO_CHANNEL"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 116, report.Items["NO_CHANNEL"].ByCategoryLabel["Cat1"].BySection["NO_SECTION"])
}
