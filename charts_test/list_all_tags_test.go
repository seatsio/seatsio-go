package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAllTags(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey1, "tag1")
	client.Charts.AddTag(chartKey1, "tag2")

	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client.Charts.AddTag(chartKey2, "tag3")

	list, err := client.Charts.ListAllTags()
	require.NoError(t, err)

	require.Len(t, list, 3)
	require.Subset(t, list, []string{"tag1", "tag2", "tag3"})
}
