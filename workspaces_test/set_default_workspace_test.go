package workspaces

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateProductionWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)
	require.False(t, workspace.IsDefault)

	err = client.Workspaces.SetDefaultWorkspace(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.True(t, retrievedWorkspace.IsDefault)

	originalDefaultWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), company.Workspace.Key)
	require.NoError(t, err)
	require.False(t, originalDefaultWorkspace.IsDefault)
}
