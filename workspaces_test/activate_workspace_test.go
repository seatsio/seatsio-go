package workspaces

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestActivateWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(workspace.Key)
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.False(t, retrievedWorkspace.IsActive)

	err = client.Workspaces.Activate(workspace.Key)
	require.NoError(t, err)

	activatedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.True(t, activatedWorkspace.IsActive)
}
