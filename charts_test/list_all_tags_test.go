package charts

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAllTags(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.AddTag(chartKey1, "tag1")
	_ = client.Charts.AddTag(chartKey1, "tag2")

	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.AddTag(chartKey2, "tag3")

	list, err := client.Charts.ListAllTags()
	require.NoError(t, err)

	require.Len(t, list, 3)
	require.Subset(t, list, []string{"tag1", "tag2", "tag3"})
}
