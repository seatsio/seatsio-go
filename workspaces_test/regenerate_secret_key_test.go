package workspaces

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegenerateSecretKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)

	newKey, err := client.Workspaces.RegenerateSecretKey(workspace.Key)
	require.NoError(t, err)
	require.NotEqual(t, newKey, workspace.SecretKey)

	retrievedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.Equal(t, *newKey, retrievedWorkspace.SecretKey)
}
