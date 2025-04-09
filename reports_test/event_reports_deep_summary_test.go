package reports

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/events"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeepSummaryByStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByStatus(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.DeepSummaryByObjectType(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByCategoryKey(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByCategoryLabel(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryBySection(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, 232, report.Items["NO_SECTION"].Count)
	require.Equal(t, 116, report.Items["NO_SECTION"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 1, report.Items["NO_SECTION"].ByCategoryLabel["Cat1"].ByAvailability["not_available"])
}

func TestDeepSummaryByZone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})

	report, err := client.EventReports.DeepSummaryByZone(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, 6032, report.Items["midtrack"].Count)
	require.Equal(t, 6032, report.Items["midtrack"].ByCategoryLabel["Mid Track Stand"].Count)
}

func TestDeepSummaryByAvailability(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByAvailability(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByAvailabilityReason(test_util.RequestContext(), event.Key)
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
	event, _ := client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_, _ = client.Events.ChangeObjectStatus(test_util.RequestContext(), []string{event.Key}, []string{"A-1"}, events.BOOKED)

	report, err := client.EventReports.DeepSummaryByChannel(test_util.RequestContext(), event.Key)
	require.NoError(t, err)
	require.Equal(t, 232, report.Items["NO_CHANNEL"].Count)
	require.Equal(t, 116, report.Items["NO_CHANNEL"].ByCategoryLabel["Cat1"].Count)
	require.Equal(t, 116, report.Items["NO_CHANNEL"].ByCategoryLabel["Cat1"].BySection["NO_SECTION"])
}
