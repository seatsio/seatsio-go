package charts

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/charts"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	category := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa", Accessible: true}
	err := client.Charts.AddCategory(test_util.RequestContext(), chartKey1, category)
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(test_util.RequestContext(), chartKey1)
	require.NoError(t, err)
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": true})
}

func TestRemoveCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	categories := []events.Category{category1, category2}
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{Categories: categories})

	err = client.Charts.RemoveCategory(test_util.RequestContext(), chart.Key, events.CategoryKey{Key: 1})
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(test_util.RequestContext(), chart.Key)
	require.NoError(t, err)
	require.NotContains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": true})
}

func TestListCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{Categories: []events.Category{category1, category2}})

	categories, err := client.Charts.ListCategories(test_util.RequestContext(), chart.Key)
	require.NoError(t, err)
	require.Contains(t, categories, category1)
	require.Contains(t, categories, category2)
	require.Equal(t, 2, len(categories))
}

func TestListCategories_unknownChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Charts.ListCategories(test_util.RequestContext(), "someUnknownChart")
	require.Error(t, err)
}

func TestUpdateCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	err := client.Charts.UpdateCategory(test_util.RequestContext(), chartKey, events.CategoryKey{Key: "string11"}, charts.UpdateCategoryParams{
		Label:      "New label",
		Color:      "#bbbbbb",
		Accessible: false,
	})
	require.NoError(t, err)

	categories, err := client.Charts.ListCategories(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Contains(t, categories, events.Category{Key: events.CategoryKey{Key: "string11"}, Label: "New label", Color: "#bbbbbb", Accessible: false})
}

func TestUpdateCategory_unknownChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	err := client.Charts.UpdateCategory(test_util.RequestContext(), "someUnknownChart", events.CategoryKey{Key: "string11"}, charts.UpdateCategoryParams{
		Label:      "New label",
		Color:      "#bbbbbb",
		Accessible: false,
	})
	require.Error(t, err)
}

func TestUpdateCategory_unknownCategory(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	err := client.Charts.UpdateCategory(test_util.RequestContext(), chartKey, events.CategoryKey{Key: "unknownCategory"}, charts.UpdateCategoryParams{
		Label:      "New label",
		Color:      "#bbbbbb",
		Accessible: false,
	})
	require.Error(t, err)
}
