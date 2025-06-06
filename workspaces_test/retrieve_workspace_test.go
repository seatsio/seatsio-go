package workspaces

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	require.Equal(t, workspace.Key, retrievedWorkspace.Key)
	require.Equal(t, workspace.Name, retrievedWorkspace.Name)
	require.Equal(t, workspace.SecretKey, retrievedWorkspace.SecretKey)
	require.Equal(t, workspace.IsTest, retrievedWorkspace.IsTest)
	require.Equal(t, workspace.IsActive, retrievedWorkspace.IsActive)
	require.Equal(t, workspace.IsDefault, retrievedWorkspace.IsDefault)
}
