package charts_test

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	err := client.Charts.AddTag(chartKey, "tag1")
	require.NoError(t, err)

	err2 := client.Charts.AddTag(chartKey, "tag2")
	require.NoError(t, err2)

	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.Equal(t, 2, len(retrievedChart.Tags))
	require.Contains(t, []string{"tag1", "tag2"}, retrievedChart.Tags[0])
}

func TestSpecialCharacters(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	err := client.Charts.AddTag(chartKey, "'tag1:-'<>")
	require.NoError(t, err)

	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.Equal(t, 1, len(retrievedChart.Tags))
	require.Equal(t, "'tag1:-'<>", retrievedChart.Tags[0])
}
