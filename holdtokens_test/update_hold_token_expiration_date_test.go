package holdtokens_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUpdateHoldTokenExpirationDate(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	updatedHoldToken, err := client.HoldTokens.ExpireInMinutes(holdToken.HoldToken, 30)
	require.NoError(t, err)

	require.Equal(t, holdToken.HoldToken, updatedHoldToken.HoldToken)
	require.True(t, updatedHoldToken.ExpiresAt.After(start.Add(-29*time.Minute)))
	require.True(t, updatedHoldToken.ExpiresAt.Before(start.Add(31*time.Minute)))
	require.True(t, updatedHoldToken.ExpiresInSeconds > 29*60)
	require.True(t, updatedHoldToken.ExpiresInSeconds <= 30*60)
}
