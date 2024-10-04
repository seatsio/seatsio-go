package reports_test

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUsageReportForAllMonths(t *testing.T) {
	t.Parallel()
	test_util.AssertDemoCompanySecretKeySet(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, test_util.DemoCompanySecretKey())

	report, err := client.UsageReports.SummaryForAllMonths()

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(report.Usage), 0)
	require.NotNil(t, report.UsageCutoffDate)
	require.Equal(t, report.Usage[0].Month.Month, 2)
	require.Equal(t, report.Usage[0].Month.Year, 2014)
}

func TestUsageReportForMonth(t *testing.T) {
	t.Parallel()
	test_util.AssertDemoCompanySecretKeySet(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, test_util.DemoCompanySecretKey())

	report, err := client.UsageReports.DetailsForMonth(2021, 11)

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(report), 0)
	require.GreaterOrEqual(t, len(report[0].UsageByChart), 0)
	require.Equal(t, report[0].UsageByChart[0].UsageByEvent[0].NumUsedObjects, 143)
}

func TestUsageReportForEventInMonth(t *testing.T) {
	t.Parallel()
	test_util.AssertDemoCompanySecretKeySet(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, test_util.DemoCompanySecretKey())

	report1, report2, err := client.UsageReports.DetailsForEventInMonth(580293, 2021, 11)

	require.NoError(t, err)
	require.GreaterOrEqual(t, len(report1), 0)
	require.Equal(t, report1[0].NumFirstSelections, 1)
	require.Nil(t, report2)
}
