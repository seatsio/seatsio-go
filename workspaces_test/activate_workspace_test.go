package workspaces

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestActivateWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.False(t, retrievedWorkspace.IsActive)

	err = client.Workspaces.Activate(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	activatedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.True(t, activatedWorkspace.IsActive)
}
