package charts

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/charts"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateChartWithDefaultParameters(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(&charts.CreateChartParams{})

	require.NoError(t, err)
	require.NotNil(t, chart.Id)
	require.NotNil(t, chart.Key)
	require.Equal(t, "NOT_USED", chart.Status)
	require.Equal(t, "Untitled chart", chart.Name)
	require.NotNil(t, chart.PublishedVersionThumbnailUrl)
	require.NotNil(t, chart.DraftVersionThumbnailUrl)
	require.Nil(t, chart.Events)
	require.Empty(t, chart.Tags)
	require.False(t, chart.Archived)
	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)
	require.Equal(t, "SIMPLE", drawing["venueType"])
	require.Empty(t, getCategories(drawing))
}

func TestCreateChartWithName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(&charts.CreateChartParams{Name: "aChart"})

	require.NoError(t, err)
	require.Equal(t, "aChart", chart.Name)
	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)
	require.Equal(t, "SIMPLE", drawing["venueType"])
	require.Empty(t, getCategories(drawing))
}

func TestCreateChartWithVenueType(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chart, err := client.Charts.Create(&charts.CreateChartParams{VenueType: "SIMPLE"})

	require.NoError(t, err)
	require.Equal(t, "Untitled chart", chart.Name)
	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)
	require.Equal(t, "SIMPLE", drawing["venueType"])
	require.Empty(t, getCategories(drawing))
}

func TestCreateChartWithCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	categories := []events.Category{category1, category2}

	chart, err := client.Charts.Create(&charts.CreateChartParams{Categories: categories})
	require.NoError(t, err)
	require.Equal(t, "Untitled chart", chart.Name)
	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)
	require.Equal(t, "SIMPLE", drawing["venueType"])
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": false})
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": "anotherCat", "label": "Category 2", "color": "#bbbbbb", "accessible": true})
}

func TestCreateChartInSpecificWorkspaceAsCompanyAdmin(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	defaultClient := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := defaultClient.Workspaces.CreateTestWorkspace("anotherWorkspace")
	require.NoError(t, err)

	workspaceClient := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey, seatsio.ClientSupport.WorkspaceKey(workspace.Key))
	chart, err := workspaceClient.Charts.Create(&charts.CreateChartParams{})
	require.NoError(t, err)

	retrievedDefaultWorkspaceCharts, err := defaultClient.Charts.ListAll()
	require.NoError(t, err)
	require.Len(t, retrievedDefaultWorkspaceCharts, 0)

	workspaceCharts, err := workspaceClient.Charts.ListAll()
	require.NoError(t, err)
	require.Len(t, workspaceCharts, 1)
	require.Equal(t, chart.Key, workspaceCharts[0].Key)
}

func TestCreateChartInSpecificWorkspaceAsWorkspaceAdmin(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	defaultClient := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := defaultClient.Workspaces.CreateTestWorkspace("anotherWorkspace")
	require.NoError(t, err)

	workspaceClient := seatsio.NewSeatsioClient(test_util.BaseUrl, workspace.SecretKey)
	chart, err := workspaceClient.Charts.Create(&charts.CreateChartParams{})
	require.NoError(t, err)

	retrievedDefaultWorkspaceCharts, err := defaultClient.Charts.ListAll()
	require.NoError(t, err)
	require.Len(t, retrievedDefaultWorkspaceCharts, 0)

	workspaceCharts, err := workspaceClient.Charts.ListAll()
	require.NoError(t, err)
	require.Len(t, workspaceCharts, 1)
	require.Equal(t, chart.Key, workspaceCharts[0].Key)
}
