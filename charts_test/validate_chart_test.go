package charts

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatePublishedChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	result, err := client.Charts.ValidatePublishedVersion(test_util.RequestContext(), chartKey)
	require.NoError(t, err)
	require.Empty(t, result.Errors)
}

func TestValidateDraftChart(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	chartKey := test_util.CreateTestChartWithErrors(t, company.Admin.SecretKey)

	_, err := client.Charts.ValidateDraftVersion(test_util.RequestContext(), chartKey)
	require.Error(t, err)
}
