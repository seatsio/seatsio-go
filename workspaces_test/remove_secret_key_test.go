package workspaces

import (
	"testing"

	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
)

func TestRemoveSecretKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	newKey, err := client.Workspaces.AddSecretKey(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.NotEqual(t, newKey, workspace.SecretKey)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.Contains(t, retrievedWorkspace.SecretKeys, workspace.SecretKey, newKey)

	err = client.Workspaces.RemoveSecretKey(test_util.RequestContext(), workspace.Key, workspace.SecretKey)
	require.NoError(t, err)

	finalStateWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.NotContains(t, finalStateWorkspace.SecretKeys, workspace.Key)
	require.Contains(t, finalStateWorkspace.SecretKeys, *newKey)
}
