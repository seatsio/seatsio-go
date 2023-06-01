package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListChartsInArchive(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	err1 := client.Charts.MoveToArchive(chartKey1)
	require.NoError(t, err1)
	err2 := client.Charts.MoveToArchive(chartKey2)
	require.NoError(t, err2)

	charts, err := client.Charts.Archive.All(20)

	require.NoError(t, err)
	require.Equal(t, 2, len(charts))
	require.Equal(t, chartKey2, charts[0].Key)
	require.Equal(t, chartKey1, charts[1].Key)
}
