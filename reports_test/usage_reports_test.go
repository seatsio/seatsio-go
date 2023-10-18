package reports

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSummaryForAllMonths(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	summary, err := client.UsageReports.SummaryForAllMonths()
	require.NoError(t, err)

	require.NotEmpty(t, summary)
}

/*
// Further testing requires additional implementation of the event API (can't book right now!)
func TestDetailsForMonth(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)
	_ = test_util.CreateTestChart(t, company.Admin.SecretKey)

	_, err := client.UsageReports.DetailsForMonth(2020, 8)
	require.NoError(t, err)
}
*/
