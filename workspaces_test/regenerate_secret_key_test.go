package workspaces

import (
	"github.com/seatsio/seatsio-go/v9"
	"github.com/seatsio/seatsio-go/v9/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegenerateSecretKey(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	newKey, err := client.Workspaces.RegenerateSecretKey(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.NotEqual(t, newKey, workspace.SecretKey)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.Equal(t, *newKey, retrievedWorkspace.SecretKey)
}
