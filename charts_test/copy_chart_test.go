package charts

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopyChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	copiedChart, err := client.Charts.Copy(test_util.RequestContext(), chartKey)

	require.NoError(t, err)
	require.Equal(t, "Sample chart (copy)", copiedChart.Name)
}
