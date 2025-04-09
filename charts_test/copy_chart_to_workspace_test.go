package charts

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCopyChartToWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my ws")
	require.NoError(t, err)

	copiedChart, err := client.Charts.CopyToWorkspace(test_util.RequestContext(), chartKey, workspace.Key)

	require.Equal(t, "Sample chart", copiedChart.Name)

	workspaceClient := seatsio.NewSeatsioClient(test_util.BaseUrl, workspace.SecretKey)
	retrievedChart, err := workspaceClient.Charts.Retrieve(test_util.RequestContext(), copiedChart.Key)
	require.NoError(t, err)
	require.Equal(t, "Sample chart", retrievedChart.Name)
}
