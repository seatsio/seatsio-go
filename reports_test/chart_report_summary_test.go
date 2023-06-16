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

	summaryChartReport, err := client.ChartReports.SummaryByObjectType(chartKey, "false")

	require.NoError(t, err)
	seatReportItem := reports.ChartSummaryReportItem{
		Count:     32,
		BySection: map[string]interface{}{"NO_SECTION": float64(32)},
		ByCategoryKey: map[string]interface{}{
			"9":  float64(16),
			"10": float64(16),
		},
		ByCategoryLabel: map[string]interface{}{
			"Cat1": float64(16),
			"Cat2": float64(16),
		},
	}
	gaReportItem := reports.ChartSummaryReportItem{
		Count:     200,
		BySection: map[string]interface{}{"NO_SECTION": float64(200)},
		ByCategoryKey: map[string]interface{}{
			"9":  float64(100),
			"10": float64(100),
		},
		ByCategoryLabel: map[string]interface{}{
			"Cat1": float64(100),
			"Cat2": float64(100),
		},
	}
	emptyReportItem := reports.ChartSummaryReportItem{
		Count:           0,
		BySection:       map[string]interface{}{},
		ByCategoryKey:   map[string]interface{}{},
		ByCategoryLabel: map[string]interface{}{},
	}
	require.Equal(t, seatReportItem, summaryChartReport.Items["seat"])
	require.Equal(t, gaReportItem, summaryChartReport.Items["generalAdmission"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["booth"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["table"])
}

func TestSummaryByObjectType_BookWholeTablesTrue(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByObjectType(chartKey, "true")

	require.NoError(t, err)
	tableReportItem := reports.ChartSummaryReportItem{
		Count:     2,
		BySection: map[string]interface{}{"NO_SECTION": float64(2)},
		ByCategoryKey: map[string]interface{}{
			"9": float64(2),
		},
		ByCategoryLabel: map[string]interface{}{
			"Cat1": float64(2),
		},
	}
	emptyReportItem := reports.ChartSummaryReportItem{
		Count:           0,
		BySection:       map[string]interface{}{},
		ByCategoryKey:   map[string]interface{}{},
		ByCategoryLabel: map[string]interface{}{},
	}
	require.Equal(t, emptyReportItem, summaryChartReport.Items["seat"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["generalAdmission"])
	require.Equal(t, emptyReportItem, summaryChartReport.Items["booth"])
	require.Equal(t, tableReportItem, summaryChartReport.Items["table"])
}
