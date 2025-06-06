package charts

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/charts"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateName(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	categories := []events.Category{category1}
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{VenueType: "SIMPLE", Categories: categories})

	err = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Name: "aChart"})
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(test_util.RequestContext(), chart.Key)
	require.Equal(t, "aChart", drawing["name"])
	require.Equal(t, "SIMPLE", drawing["venueType"])
}

func TestUpdateCategories(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	category1 := events.Category{Key: events.CategoryKey{Key: 1}, Label: "Category 1", Color: "#aaaaaa"}
	category2 := events.Category{Key: events.CategoryKey{Key: "anotherCat"}, Label: "Category 2", Color: "#bbbbbb", Accessible: true}
	categories := []events.Category{category1, category2}
	chart, err := client.Charts.Create(test_util.RequestContext(), &charts.CreateChartParams{})
	require.NoError(t, err)

	err = client.Charts.Update(test_util.RequestContext(), chart.Key, &charts.UpdateChartParams{Categories: categories})
	require.NoError(t, err)

	drawing, err := client.Charts.RetrievePublishedVersion(test_util.RequestContext(), chart.Key)
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": float64(1), "label": "Category 1", "color": "#aaaaaa", "accessible": false})
	require.Contains(t,
		getCategories(drawing),
		map[string]interface{}{"key": "anotherCat", "label": "Category 2", "color": "#bbbbbb", "accessible": true})
	require.Equal(t, 2, len(getCategories(drawing)))

}
