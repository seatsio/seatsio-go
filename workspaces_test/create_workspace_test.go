package workspaces_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateProductionWorkspace(t *testing.T) {
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)

	require.Equal(t, "my workspace", workspace.Name)
	require.NotNil(t, workspace.Key)
	require.NotNil(t, workspace.SecretKey)
	require.False(t, workspace.IsTest)
	require.True(t, workspace.IsActive)
}

func TestCreateTestWorkspace(t *testing.T) {
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	workspace, err := client.Workspaces.CreateTestWorkspace("my workspace")
	require.NoError(t, err)

	require.Equal(t, "my workspace", workspace.Name)
	require.NotNil(t, workspace.Key)
	require.NotNil(t, workspace.SecretKey)
	require.True(t, workspace.IsTest)
	require.True(t, workspace.IsActive)
}
