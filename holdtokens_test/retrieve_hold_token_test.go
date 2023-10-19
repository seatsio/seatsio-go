package holdtokens_test

import (
	"github.com/seatsio/seatsio-go/v6"
	"github.com/seatsio/seatsio-go/v6/test_util"
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

	require.NotEmpty(t, retrievedHoldToken.HoldToken)
	require.Equal(t, holdToken.ExpiresAt, retrievedHoldToken.ExpiresAt)
	require.True(t, retrievedHoldToken.ExpiresInSeconds > 14*60)
	require.True(t, retrievedHoldToken.ExpiresInSeconds <= 15*60)
}
