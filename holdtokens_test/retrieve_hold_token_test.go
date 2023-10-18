package holdtokens_test

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	retrievedHoldToken, err := client.HoldTokens.Retrieve(holdToken.HoldToken)
	require.NoError(t, err)

	require.Equal(t, holdToken, retrievedHoldToken)
}
