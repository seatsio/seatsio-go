package charts

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMoveOutOfArchive(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = client.Charts.MoveToArchive(chartKey)

	err := client.Charts.MoveOutOfArchive(chartKey)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.False(t, retrievedChart.Archived)
}
