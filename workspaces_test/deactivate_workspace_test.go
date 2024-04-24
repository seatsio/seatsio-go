package workspaces

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeactivateWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)
	require.True(t, workspace.IsActive)

	err = client.Workspaces.Deactivate(workspace.Key)
	require.NoError(t, err)

	deactivatedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.False(t, deactivatedWorkspace.IsActive)
}
