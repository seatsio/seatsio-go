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

	charts, err := client.Charts.ListAll()
	require.NoError(t, err)

	require.Equal(t, chartKey3, charts[0].Key)
	require.Equal(t, chartKey2, charts[1].Key)
	require.Equal(t, chartKey1, charts[2].Key)
}

func TestFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	client.Charts.Update(chartKey1, &charts.UpdateChartParams{Name: "foo"})
	client.Charts.Update(chartKey2, &charts.UpdateChartParams{Name: "bar"})
	client.Charts.Update(chartKey3, &charts.UpdateChartParams{Name: "foofoo"})

	charts, err := client.Charts.List(&charts.ListChartParams{Filter: "foo"}).All(5)
	require.NoError(t, err)

	require.Equal(t, chartKey3, charts[0].Key)
	require.Equal(t, chartKey1, charts[1].Key)
	require.Equal(t, 2, len(charts))
}

func TestTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)

	client.Charts.AddTag(chartKey1, "aTag")
	client.Charts.AddTag(chartKey2, "anotherTag")
	client.Charts.AddTag(chartKey3, "aTag")

	charts, err := client.Charts.List(&charts.ListChartParams{Tag: "aTag"}).All(5)
	require.NoError(t, err)
	require.Equal(t, chartKey3, charts[0].Key)
	require.Equal(t, chartKey1, charts[1].Key)
	require.Equal(t, 2, len(charts))
}

func TestTagAndFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.Update(chartKey1, &charts.UpdateChartParams{Name: "bar"})
	client.Charts.AddTag(chartKey1, "foo")

	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.Update(chartKey2, &charts.UpdateChartParams{Name: "someOtherName"})
	client.Charts.AddTag(chartKey2, "foo")

	chartKey3 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.Update(chartKey3, &charts.UpdateChartParams{Name: "bar"})
	client.Charts.AddTag(chartKey3, "foo")

	chartKey4 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.Update(chartKey4, &charts.UpdateChartParams{Name: "bar"})

	charts, err := client.Charts.List(&charts.ListChartParams{Filter: "bar", Tag: "foo"}).All(5)
	require.NoError(t, err)
	require.Equal(t, chartKey3, charts[0].Key)
	require.Equal(t, chartKey1, charts[1].Key)
	require.Equal(t, 2, len(charts))
}

func TestExpand(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	event1, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})
	event2, err := client.Events.Create(&events.CreateEventParams{ChartKey: chartKey})

	charts, err := client.Charts.List(&charts.ListChartParams{Expand: true}).All(5)
	require.NoError(t, err)
	require.Equal(t, 1, len(charts))
	require.Equal(t, event1.Key, charts[0].Events[1].Key)
	require.Equal(t, event2.Key, charts[0].Events[0].Key)
}

func TestPageSize(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(nil, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(chartsPage.Items))

}

func TestListChartsWithValidation(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(&charts.ListChartParams{Validation: true}, 10)
	require.NoError(t, err)
	charts := chartsPage.Items
	require.Equal(t, 1, len(charts))
	require.Empty(t, charts[0].Validation.Errors)
	require.Empty(t, charts[0].Validation.Warnings)
}

func TestListChartsWithoutValidation(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	chartsPage, err := client.Charts.ListFirstPage(&charts.ListChartParams{Validation: true}, 10)
	require.NoError(t, err)
	charts := chartsPage.Items
	require.Equal(t, 1, len(charts))
	require.Empty(t, charts[0].Validation.Errors)
	require.Empty(t, charts[0].Validation.Warnings)
}
