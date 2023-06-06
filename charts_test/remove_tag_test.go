package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey1, "tag1")
	client.Charts.AddTag(chartKey1, "tag2")

	err := client.Charts.RemoveTag(chartKey1, "tag2")

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chartKey1)
	require.NoError(t, err)
	require.Contains(t, retrievedChart.Tags, "tag1")
	require.NotContains(t, retrievedChart.Tags, "tag2")
}
