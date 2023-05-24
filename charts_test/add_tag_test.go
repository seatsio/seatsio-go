package charts_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	err := client.Charts.AddTag(chartKey, "tag1")
	require.NoError(t, err)

	err2 := client.Charts.AddTag(chartKey, "tag2")
	require.NoError(t, err2)

	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.Equal(t, 2, len(retrievedChart.Tags))
	require.Equal(t, "tag1", retrievedChart.Tags[0])
	require.Equal(t, "tag2", retrievedChart.Tags[1])
}

func TestSpecialCharacters(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	err := client.Charts.AddTag(chartKey, "'tag1:-'<>")
	require.NoError(t, err)

	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.Equal(t, 1, len(retrievedChart.Tags))
	require.Equal(t, "'tag1:-'<>", retrievedChart.Tags[0])
}
