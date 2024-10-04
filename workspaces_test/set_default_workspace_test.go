package workspaces

import (
	"github.com/seatsio/seatsio-go/v8"
	"github.com/seatsio/seatsio-go/v8/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateProductionWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)
	require.False(t, workspace.IsDefault)

	err = client.Workspaces.SetDefaultWorkspace(workspace.Key)
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.True(t, retrievedWorkspace.IsDefault)

	originalDefaultWorkspace, err := client.Workspaces.Retrieve(company.Workspace.Key)
	require.NoError(t, err)
	require.False(t, originalDefaultWorkspace.IsDefault)
}
