package holdtokens_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrieveHoldToken(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(company.Admin.SecretKey, test_util.BaseUrl)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	retrievedHoldToken, err := client.HoldTokens.Retrieve(holdToken.HoldToken)
	require.NoError(t, err)

	require.Equal(t, holdToken, retrievedHoldToken)
}
