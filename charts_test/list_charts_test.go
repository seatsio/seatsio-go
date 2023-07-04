package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/charts"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAll(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

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
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	_ = client.Charts.Update(chartKey1, &charts.UpdateChartParams{Name: "foo"})
	_ = client.Charts.Update(chartKey2, &charts.UpdateChartParams{Name: "bar"})
	_ = client.Charts.Update(chartKey3, &charts.UpdateChartParams{Name: "foofoo"})

	retrievedCharts, err := client.Charts.List(&charts.ListChartParams{Filter: "foo"}).All()
	require.NoError(t, err)

	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	_ = client.Charts.AddTag(chartKey1, "aTag")
	_ = client.Charts.AddTag(chartKey2, "anotherTag")
	_ = client.Charts.AddTag(chartKey3, "aTag")

	retrievedCharts, err := client.Charts.List(&charts.ListChartParams{Tag: "aTag"}).All()
	require.NoError(t, err)
	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestTagAndFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

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

	retrievedCharts, err := client.Charts.List(&charts.ListChartParams{Filter: "bar", Tag: "foo"}).All()
	require.NoError(t, err)
	require.Equal(t, chartKey3, retrievedCharts[0].Key)
	require.Equal(t, chartKey1, retrievedCharts[1].Key)
	require.Equal(t, 2, len(retrievedCharts))
}

func TestExpand(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	retrievedCharts, err := client.Charts.List(&charts.ListChartParams{Expand: true}).All()
	require.NoError(t, err)
	require.Equal(t, 1, len(retrievedCharts))
	require.Equal(t, event1.Key, retrievedCharts[0].Events[1].Key)
	require.Equal(t, event2.Key, retrievedCharts[0].Events[0].Key)
}

func TestPageSize(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(nil)
	require.NoError(t, err)
	require.Equal(t, 3, len(chartsPage.Items))

}

func TestListChartsWithValidation(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(&charts.ListChartParams{Validation: true})
	require.NoError(t, err)
	retrievedCharts := chartsPage.Items
	require.Equal(t, 1, len(retrievedCharts))
	require.Empty(t, retrievedCharts[0].Validation.Errors)
	require.Empty(t, retrievedCharts[0].Validation.Warnings)
}

func TestListChartsWithoutValidation(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(&charts.ListChartParams{Validation: true})
	require.NoError(t, err)
	retrievedCharts := chartsPage.Items
	require.Equal(t, 1, len(retrievedCharts))
	require.Empty(t, retrievedCharts[0].Validation.Errors)
	require.Empty(t, retrievedCharts[0].Validation.Warnings)
}
