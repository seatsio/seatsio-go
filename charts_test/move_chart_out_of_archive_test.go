package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMoveOutOfArchive(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.MoveToArchive(chartKey)

	err := client.Charts.MoveOutOfArchive(chartKey)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.False(t, retrievedChart.Archived)
}
