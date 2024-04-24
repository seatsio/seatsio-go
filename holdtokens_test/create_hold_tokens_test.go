package holdtokens_test

import (
	"github.com/seatsio/seatsio-go/v7"
	"github.com/seatsio/seatsio-go/v7/test_util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateHoldToken(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	holdToken, err := client.HoldTokens.Create()
	require.NoError(t, err)

	require.NotEmpty(t, holdToken.HoldToken)
	require.True(t, holdToken.ExpiresAt.After(start.Add(-14*time.Minute)))
	require.True(t, holdToken.ExpiresAt.Before(start.Add(16*time.Minute)))
	require.True(t, holdToken.ExpiresInSeconds > 14*60)
	require.True(t, holdToken.ExpiresInSeconds <= 15*60)
}

func TestCreateHoldTokenWithExpiresInMinutes(t *testing.T) {
	t.Parallel()
	start := time.Now()
	company := test_util.CreateTestCompany(t)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	holdToken, err := client.HoldTokens.CreateWithExpiration(5)
	require.NoError(t, err)

	require.NotEmpty(t, holdToken.HoldToken)
	require.True(t, holdToken.ExpiresAt.After(start.Add(-4*time.Minute)))
	require.True(t, holdToken.ExpiresAt.Before(start.Add(6*time.Minute)))
	require.True(t, holdToken.ExpiresInSeconds > 4*60)
	require.True(t, holdToken.ExpiresInSeconds <= 5*60)
}
