package workspaces

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)

	err = client.Workspaces.Update(test_util.RequestContext(), workspace.Key, "New name")
	require.NoError(t, err)

	retrievedWorkspace, err := client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)
	require.Equal(t, "New name", retrievedWorkspace.Name)
}
