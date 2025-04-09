package workspaces

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/seatsio/seatsio-go/v10/workspaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListInactiveWorkspaces(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(test_util.RequestContext(), workspaces.Inactive)
	require.NoError(t, err)

	require.Equal(t, 1, len(retrievedWorkspaces))
	require.Equal(t, workspace2.Key, retrievedWorkspaces[0].Key)
}

func TestListInactiveWorkspacesWithFilter(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace1, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace1.Key)
	require.NoError(t, err)
	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(test_util.RequestContext(), workspaces.Inactive, workspaces.WorkspaceSupport.WithFilter("workspace 2"))
	require.NoError(t, err)

	require.Equal(t, 1, len(retrievedWorkspaces))
	require.Equal(t, workspace2.Key, retrievedWorkspaces[0].Key)
}

func TestListInactiveWorkspacesWithFilterNoResults(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	_, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 1")
	require.NoError(t, err)
	workspace2, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 2")
	require.NoError(t, err)
	_, err = client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "workspace 3")
	require.NoError(t, err)

	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace2.Key)
	require.NoError(t, err)

	retrievedWorkspaces, err := client.Workspaces.ListAll(test_util.RequestContext(), workspaces.Inactive, workspaces.WorkspaceSupport.WithFilter("workspace 1"))
	require.NoError(t, err)

	require.Equal(t, 0, len(retrievedWorkspaces))
}
