package reports_test

import (
	"github.com/seatsio/seatsio-go"
	reports "github.com/seatsio/seatsio-go/reports/charts"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSummaryByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByObjectType(chartKey)

	require.NoError(t, err)
	seatReportItem := reports.ChartSummaryReportItem{Count: 32}
	gaReportItem := reports.ChartSummaryReportItem{Count: 200}
	emptyReportItem := reports.ChartSummaryReportItem{Count: 0}
	require.Equal(t, seatReportItem, summaryChartReport.Items["seat"])
	require.Equal(t, gaReportItem, summaryChartReport.Items["generalAdmission"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["booth"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["table"])
}
