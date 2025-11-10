package charts

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListChartsInArchive(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	chartKey1 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	chartKey2 := test_util.CreateTestChart(t, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	err1 := client.Charts.MoveToArchive(test_util.RequestContext(), chartKey1)
	require.NoError(t, err1)
	err2 := client.Charts.MoveToArchive(test_util.RequestContext(), chartKey2)
	require.NoError(t, err2)

	charts, err := client.Charts.Archive.All(test_util.RequestContext())

	require.NoError(t, err)
	require.Equal(t, 2, len(charts))
	require.Equal(t, chartKey2, charts[0].Key)
	require.Equal(t, chartKey1, charts[1].Key)
}
