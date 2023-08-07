package workspaces

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace("my workspace")
	require.NoError(t, err)

	err = client.Workspaces.Update(workspace.Key, "New name")
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(workspace.Key)
	require.NoError(t, err)
	require.Equal(t, "New name", retrievedWorkspace.Name)
}
