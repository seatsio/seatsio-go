package reports

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/charts"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/reports"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestReportItemProperties(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey)
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
	require.False(t, item.IsAccessible)
	require.False(t, item.HasLiftUpArmrests)
	require.False(t, item.IsHearingImpaired)
	require.False(t, item.IsSemiAmbulatorySeat)
	require.False(t, item.HasSignLanguageInterpretation)
	require.False(t, item.IsPlusSize)
	require.False(t, item.IsCompanionSeat)
	require.False(t, item.HasRestrictedView)
	require.NotNil(t, item.Floor)
}

func TestReportItemPropertiesForGA(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey)
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

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey, reports.ChartReportOptions.BookWholeTablesTrue())
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

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["A-1"], 1)
	require.Len(t, chartReport.Items["A-2"], 1)
}

func TestByLabel_BookWholeTablesTrue(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithTables(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey, reports.ChartReportOptions.BookWholeTablesTrue())
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

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey, reports.ChartReportOptions.BookWholeTablesChart())
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

func TestByLabel_WithFloors(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithFloors(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByLabel(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "1", chartReport.Items["S1-A-1"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["S1-A-1"][0].Floor.DisplayName)
	require.Equal(t, "1", chartReport.Items["S1-A-2"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["S1-A-2"][0].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["S2-B-1"][0].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["S2-B-1"][0].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["S2-B-2"][0].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["S2-B-2"][0].Floor.DisplayName)
}

func TestByObjectType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByObjectType(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["seat"], 32)
	require.Len(t, chartReport.Items["generalAdmission"], 2)
}

func TestByObjectType_WithFloors(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithFloors(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByObjectType(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "1", chartReport.Items["seat"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["seat"][0].Floor.DisplayName)
	require.Equal(t, "1", chartReport.Items["seat"][1].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["seat"][1].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["seat"][2].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["seat"][2].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["seat"][3].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["seat"][3].Floor.DisplayName)
}

func TestByCategoryKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryKey(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["9"], 17)
	require.Len(t, chartReport.Items["10"], 17)
}

func TestByCategoryKey_WithFloors(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithFloors(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryKey(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "1", chartReport.Items["1"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["1"][0].Floor.DisplayName)
	require.Equal(t, "1", chartReport.Items["1"][1].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["1"][1].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["2"][0].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["2"][0].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["2"][1].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["2"][1].Floor.DisplayName)
}

func TestByCategoryLabel(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryLabel(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Cat1"], 17)
	require.Len(t, chartReport.Items["Cat2"], 17)
}

func TestByCategoryLabel_WithFloors(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithFloors(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByCategoryLabel(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "1", chartReport.Items["CatA"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["CatA"][0].Floor.DisplayName)
	require.Equal(t, "1", chartReport.Items["CatA"][1].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["CatA"][1].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["CatB"][0].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["CatB"][0].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["CatB"][1].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["CatB"][1].Floor.DisplayName)
}

func TestBySection(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithSections(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.BySection(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Section A"], 36)
	require.Len(t, chartReport.Items["Section B"], 35)
}

func TestBySection_WithFloors(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithFloors(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.BySection(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Equal(t, "1", chartReport.Items["S1"][0].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["S1"][0].Floor.DisplayName)
	require.Equal(t, "1", chartReport.Items["S1"][1].Floor.Name)
	require.Equal(t, "Floor 1", chartReport.Items["S1"][1].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["S2"][0].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["S2"][0].Floor.DisplayName)
	require.Equal(t, "2", chartReport.Items["S2"][1].Floor.Name)
	require.Equal(t, "Floor 2", chartReport.Items["S2"][1].Floor.DisplayName)
}

func TestByZone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)

	chartReport, err := client.ChartReports.ByZone(test_util.RequestContext(), chartKey)
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
	_, _ = client.Events.Create(test_util.RequestContext(), &events.CreateEventParams{ChartKey: chartKey})
	_ = client.Charts.Update(test_util.RequestContext(), chartKey, &charts.UpdateChartParams{Name: "Foo"})

	chartReport, err := client.ChartReports.BySection(test_util.RequestContext(), chartKey, reports.ChartReportOptions.UseDraftVersion())
	require.NoError(t, err)
	require.Len(t, chartReport.Items["Section A"], 36)
	require.Len(t, chartReport.Items["Section B"], 35)
}
