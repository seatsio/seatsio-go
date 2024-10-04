package charts

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/charts"
	"github.com/seatsio/seatsio-go/v8/events"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

var sup = charts.ChartSupport

func TestListAll(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	retrievedCharts, err := client.Charts.ListAll()
	require.NoError(t, err)

	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey2, retrievedCharts[1].Key)
	require.Equal(t, chartKey1, retrievedCharts[2].Key)
}

func TestFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	_ = client.Charts.Update(chartKey1, &charts.UpdateChartParams{Name: "foo"})
	_ = client.Charts.Update(chartKey2, &charts.UpdateChartParams{Name: "bar"})
	_ = client.Charts.Update(chartKey3, &charts.UpdateChartParams{Name: "foofoo"})

	retrievedCharts, err := client.Charts.List().All(sup.WithFilter("foo"))
	require.NoError(t, err)

	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	_ = client.Charts.AddTag(chartKey1, "aTag")
	_ = client.Charts.AddTag(chartKey2, "anotherTag")
	_ = client.Charts.AddTag(chartKey3, "aTag")

	retrievedCharts, err := client.Charts.List().All(sup.WithTag("aTag"))
	require.NoError(t, err)
	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestTagAndFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.Update(chartKey1, &charts.UpdateChartParams{Name: "bar"})
	_ = client.Charts.AddTag(chartKey1, "foo")

	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.Update(chartKey2, &charts.UpdateChartParams{Name: "someOtherName"})
	_ = client.Charts.AddTag(chartKey2, "foo")

	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.Update(chartKey3, &charts.UpdateChartParams{Name: "bar"})
	_ = client.Charts.AddTag(chartKey3, "foo")

	chartKey4 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.Update(chartKey4, &charts.UpdateChartParams{Name: "bar"})

	retrievedCharts, err := client.Charts.List().All(sup.WithFilter("bar"), sup.WithTag("foo"))
	require.NoError(t, err)
	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestExpandAll(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey := test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	retrievedCharts, err := client.Charts.List().All(sup.WithExpandEvents(), sup.WithExpandValidation(), sup.WithExpandVenueType(), sup.WithExpandZones())
	require.NoError(t, err)
	require.Equal(t, 1, len(retrievedCharts))
	require.Equal(t, event1.Key, retrievedCharts[0].Events[1].Key)
	require.Equal(t, event2.Key, retrievedCharts[0].Events[0].Key)
	require.Empty(t, retrievedCharts[0].Validation.Errors)
	require.Empty(t, retrievedCharts[0].Validation.Warnings)
	require.Equal(t, "WITH_ZONES", retrievedCharts[0].VenueType)
	require.Equal(t, []charts.Zone{{"finishline", "Finish Line"}, {"midtrack", "Mid Track"}}, retrievedCharts[0].Zones)
}

func TestExpandNone(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	test_util.CreateTestChartWithZones(t, company.Admin.SecretKey)

	retrievedCharts, err := client.Charts.List().All()
	require.NoError(t, err)
	require.Equal(t, 1, len(retrievedCharts))
	require.Nil(t, retrievedCharts[0].Events)
	require.Nil(t, retrievedCharts[0].Validation)
	require.Empty(t, retrievedCharts[0].VenueType)
	require.Nil(t, retrievedCharts[0].Zones)
}

func TestPageSize(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage()
	require.NoError(t, err)
	require.Equal(t, 3, len(chartsPage.Items))
}
