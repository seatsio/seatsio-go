package charts

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemoveTag(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.AddTag(test_util.RequestContext(), chartKey1, "tag1")
	_ = client.Charts.AddTag(test_util.RequestContext(), chartKey1, "tag2")

	err := client.Charts.RemoveTag(test_util.RequestContext(), chartKey1, "tag2")

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(test_util.RequestContext(), chartKey1)
	require.NoError(t, err)
	require.Contains(t, retrievedChart.Tags, "tag1")
	require.NotContains(t, retrievedChart.Tags, "tag2")
}
