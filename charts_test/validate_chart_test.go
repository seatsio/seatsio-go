package charts

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatePublishedChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	result, err := client.Charts.ValidatePublishedVersion(chartKey)
	require.NoError(t, err)
	require.Empty(t, result.Errors)
}

func TestValidateDraftChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	chartKey := test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	_, err := client.Charts.ValidateDraftVersion(chartKey)
	require.Error(t, err)
}
