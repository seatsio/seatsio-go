package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	category := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa", Accessible: true}
	err := client.Charts.AddCategory(chartKey1, category)
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(chartKey1)
	require.NoError(t, err)
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": true})
}

func TestRemoveCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	categories := []events.Category{category1, category2}
	chart, err := client.Charts.Create(&charts.CreateChartParams{Categories: categories})

	err = client.Charts.RemoveCategory(chart.Key, events.CategoryKey{Key: 1})
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(chart.Key)
	require.NoError(t, err)
	require.NotContains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": true})
}

func TestListCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	chart, err := client.Charts.Create(&charts.CreateChartParams{Categories: []events.Category{category1, category2}})

	categories, err := client.Charts.ListCategories(chart.Key)
	require.NoError(t, err)
	require.Contains(t, categories.Categories, category1)
	require.Contains(t, categories.Categories, category2)
	require.Equal(t, 2, len(categories.Categories))
}

func TestListCategories_unknownChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_, err := client.Charts.ListCategories("someUnknownChart")
	require.Error(t, err)
}
