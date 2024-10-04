package charts

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMoveToArchive(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)

	err := client.Charts.MoveToArchive(chartKey)

	require.NoError(t, err)
	retrievedChart, err := client.Charts.Retrieve(chartKey)
	require.NoError(t, err)
	require.True(t, retrievedChart.Archived)
}
