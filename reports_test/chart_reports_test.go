package reports

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/charts"
	"github.com/seatsio/seatsio-go/v7/events"
	"github.com/seatsio/seatsio-go/v7/reports"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReportItemProperties(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["A-1"], 1)
	item := chartReport.Items["A-1"][0]
	require.Equal(t, "A-1", item.Label)
	require.Equal(t, events.Labels{
		Own:    events.LabelAndType{Label: "1", Type: "seat"},
		Parent: events.LabelAndType{Label: "A", Type: "row"},
	}, item.Labels)
	require.Equal(t, events.IDs{Own: "1", Parent: "A", Section: ""}, item.IDs)
	require.Equal(t, "Cat1", item.CategoryLabel)
	require.Equal(t, "9", item.CategoryKey)
	require.Empty(t, item.Section)
	require.Empty(t, item.Entrance)
	require.Empty(t, item.Capacity)
	require.Equal(t, "seat", item.ObjectType)
	require.Empty(t, item.LeftNeighbour)
	require.Equal(t, "A-2", item.RightNeighbour)
	require.NotEmpty(t, item.DistanceToFocalPoint)
	require.NotNil(t, item.IsAccessible)
	require.NotNil(t, item.IsCompanionSeat)
	require.NotNil(t, item.HasRestrictedView)
}

func TestReportItemPropertiesForGA(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["GA1"], 1)
	item := chartReport.Items["GA1"][0]
	require.Equal(t, item.Capacity, 100)
	require.Equal(t, "generalAdmission", item.ObjectType)
	require.False(t, item.BookAsAWhole)
}

func TestReportItemPropertiesForTable(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey, reports.ChartReportOptions.BookWholeTablesTrue())
	require.NoError(t, err)
	require.Len(t, chartReport.Items["T1"], 1)
	item := chartReport.Items["T1"][0]
	require.False(t, item.BookAsAWhole)
	require.Equal(t, 6, item.NumSeats)
}

func TestByLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["A-1"], 1)
	require.Len(t, chartReport.Items["A-2"], 1)
}

func TestByLabel_BookWholeTablesTrue(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey, reports.ChartReportOptions.BookWholeTablesTrue())
	require.NoError(t, err)
	require.Len(t, chartReport.Items, 2)
	require.Nil(t, chartReport.Items["T1-1"])
	require.Nil(t, chartReport.Items["T1-2"])
	require.NotNil(t, chartReport.Items["T1"])
	require.NotNil(t, chartReport.Items["T2"])
}

func TestByLabel_BookWholeTablesChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(chartKey, reports.ChartReportOptions.BookWholeTablesChart())
	require.NoError(t, err)
	require.Len(t, chartReport.Items, 7)
	require.NotNil(t, chartReport.Items["T1-1"])
	require.NotNil(t, chartReport.Items["T1-2"])
	require.NotNil(t, chartReport.Items["T1-3"])
	require.NotNil(t, chartReport.Items["T1-4"])
	require.NotNil(t, chartReport.Items["T1-5"])
	require.NotNil(t, chartReport.Items["T1-6"])
	require.NotNil(t, chartReport.Items["T2"])
}

func TestByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByObjectType(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["seat"], 32)
	require.Len(t, chartReport.Items["generalAdmission"], 2)
}

func TestByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryKey(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["9"], 17)
	require.Len(t, chartReport.Items["10"], 17)
}

func TestByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryLabel(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Cat1"], 17)
	require.Len(t, chartReport.Items["Cat2"], 17)
}

func TestBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithSections(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.BySection(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Section A"], 36)
	require.Len(t, chartReport.Items["Section B"], 35)
}

func TestByZone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByZone(chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["midtrack"], 6032)
	require.Equal(t, "midtrack", chartReport.Items["midtrack"][0].Zone)
	require.Len(t, chartReport.Items["finishline"], 2865)
	require.Len(t, chartReport.Items["NO_ZONE"], 0)
}

func TestDraftVersion(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithSections(t, company.Admin.SecretKey)
	_, _ = client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(chartKey, &charts.UpdateChartParams{Name: "Foo"})

	chartReport, err := client.ChartReports.BySection(chartKey, reports.ChartReportOptions.UseDraftVersion())
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Section A"], 36)
	require.Len(t, chartReport.Items["Section B"], 35)
}
