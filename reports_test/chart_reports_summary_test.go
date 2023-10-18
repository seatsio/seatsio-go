package reports_test

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/charts"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/reports"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSummaryByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByObjectType(chartKey)

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
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByObjectType(chartKey, reports.ChartReportOptions.BookWholeTablesTrue())

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

func TestSummaryByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByCategoryKey(chartKey)

	require.NoError(t, err)
	cat9Report := reports.ChartSummaryReportItem{
		Count:     116,
		BySection: map[string]interface{}{"NO_SECTION": float64(116)},
		ByObjectType: map[string]interface{}{
			"seat":             float64(16),
			"generalAdmission": float64(100),
		},
	}
	cat10Report := reports.ChartSummaryReportItem{
		Count:     116,
		BySection: map[string]interface{}{"NO_SECTION": float64(116)},
		ByObjectType: map[string]interface{}{
			"seat":             float64(16),
			"generalAdmission": float64(100),
		},
	}
	cat11Report := reports.ChartSummaryReportItem{
		Count:        0,
		BySection:    map[string]interface{}{},
		ByObjectType: map[string]interface{}{},
	}
	noCategoryReport := reports.ChartSummaryReportItem{
		Count:        0,
		BySection:    map[string]interface{}{},
		ByObjectType: map[string]interface{}{},
	}
	require.Equal(t, cat9Report, summaryChartReport.Items["9"])
	require.Equal(t, cat10Report, summaryChartReport.Items["10"])
	require.Equal(t, cat11Report, summaryChartReport.Items["string11"])
	require.Equal(t, noCategoryReport, summaryChartReport.Items["NO_CATEGORY"])
}

func TestSummaryByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryByCategoryLabel(chartKey)

	require.NoError(t, err)
	cat1Report := reports.ChartSummaryReportItem{
		Count:     116,
		BySection: map[string]interface{}{"NO_SECTION": float64(116)},
		ByObjectType: map[string]interface{}{
			"seat":             float64(16),
			"generalAdmission": float64(100),
		},
	}
	cat2Report := reports.ChartSummaryReportItem{
		Count:     116,
		BySection: map[string]interface{}{"NO_SECTION": float64(116)},
		ByObjectType: map[string]interface{}{
			"seat":             float64(16),
			"generalAdmission": float64(100),
		},
	}
	cat3Report := reports.ChartSummaryReportItem{
		Count:        0,
		BySection:    map[string]interface{}{},
		ByObjectType: map[string]interface{}{},
	}
	noCategoryReport := reports.ChartSummaryReportItem{
		Count:        0,
		BySection:    map[string]interface{}{},
		ByObjectType: map[string]interface{}{},
	}
	require.Equal(t, cat1Report, summaryChartReport.Items["Cat1"])
	require.Equal(t, cat2Report, summaryChartReport.Items["Cat2"])
	require.Equal(t, cat3Report, summaryChartReport.Items["Cat3"])
	require.Equal(t, noCategoryReport, summaryChartReport.Items["NO_CATEGORY"])
}

func TestSummaryBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	summaryChartReport, err := client.ChartReports.SummaryBySection(chartKey)

	require.NoError(t, err)
	noSectionReport := reports.ChartSummaryReportItem{
		Count: 232,
		ByCategoryKey: map[string]interface{}{
			"9":  float64(116),
			"10": float64(116),
		},
		ByCategoryLabel: map[string]interface{}{
			"Cat1": float64(116),
			"Cat2": float64(116),
		},
		ByObjectType: map[string]interface{}{
			"seat":             float64(32),
			"generalAdmission": float64(200),
		},
	}
	require.Equal(t, noSectionReport, summaryChartReport.Items["NO_SECTION"])
}

func TestSummaryDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "Foo"})

	summaryChartReport, err := client.ChartReports.SummaryBySection(chartKey, reports.ChartReportOptionsNS{}.UseDraftVersion())

	require.NoError(t, err)
	noSectionReport := reports.ChartSummaryReportItem{
		Count: 232,
		ByCategoryKey: map[string]interface{}{
			"9":  float64(116),
			"10": float64(116),
		},
		ByCategoryLabel: map[string]interface{}{
			"Cat1": float64(116),
			"Cat2": float64(116),
		},
		ByObjectType: map[string]interface{}{
			"seat":             float64(32),
			"generalAdmission": float64(200),
		},
	}
	require.Equal(t, noSectionReport, summaryChartReport.Items["NO_SECTION"])
}
