package workspaces

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/seatsio/seatsio-go/workspaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListInactiveWorkspaces(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_, err := client.Workspaces.CreateProductionWorkspace("workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace("workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace("workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(workspaces.Inactive)
	require.NoError(t, err)

	require.Equal(t, 1, len(retrievedWorkspaces))
	require.Equal(t, workspace2.Key, retrievedWorkspaces[0].Key)
}

func TestListInactiveWorkspacesWithFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	workspace1, err := client.Workspaces.CreateProductionWorkspace("workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace("workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace("workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(workspace1.Key)
	require.NoError(t, err)
	err = client.Workspaces.Deactivate(workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(workspaces.Inactive, workspaces.WorkspaceSupport.WithFilter("workspace 2"))
	require.NoError(t, err)

	require.Equal(t, 1, len(retrievedWorkspaces))
	require.Equal(t, workspace2.Key, retrievedWorkspaces[0].Key)
}

func TestListInactiveWorkspacesWithFilterNoResults(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	_, err := client.Workspaces.CreateProductionWorkspace("workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace("workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace("workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(workspaces.Inactive, workspaces.WorkspaceSupport.WithFilter("workspace 1"))
	require.NoError(t, err)

	require.Equal(t, 0, len(retrievedWorkspaces))
}
