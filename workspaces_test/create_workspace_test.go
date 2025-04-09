package workspaces_test

import (
	"github.com/seatsio/seatsio-go/v10"
	"github.com/seatsio/seatsio-go/v10/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateProductionWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	require.Equal(t, "my workspace", workspace.Name)
	require.NotNil(t, workspace.Key)
	require.NotNil(t, workspace.SecretKey)
	require.False(t, workspace.IsTest)
	require.True(t, workspace.IsActive)
}

func TestCreateTestWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateTestWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	require.Equal(t, "my workspace", workspace.Name)
	require.NotNil(t, workspace.Key)
	require.NotNil(t, workspace.SecretKey)
	require.True(t, workspace.IsTest)
	require.True(t, workspace.IsActive)
}
