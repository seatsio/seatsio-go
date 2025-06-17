package workspaces

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/shared"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteInactiveWorkspace(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	workspace, err := client.Workspaces.CreateProductionWorkspace(test_util.RequestContext(), "my workspace")
	require.NoError(t, err)
	require.True(t, workspace.IsActive)

	err = client.Workspaces.Deactivate(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	err = client.Workspaces.Delete(test_util.RequestContext(), workspace.Key)
	require.NoError(t, err)

	_, err = client.Workspaces.Retrieve(test_util.RequestContext(), workspace.Key)
	require.Equal(t, "WORKSPACE_NOT_FOUND", err.(*shared.SeatsioError).Code)
}
