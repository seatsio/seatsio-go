package seasons

import (
	"github.com/seatsio/seatsio-go/v12"
	"github.com/seatsio/seatsio-go/v12/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateEventsWithEventKeys(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	events, err := client.Seasons.CreateEventsWithEventKeys(test_util.RequestContext(), season.Key, "event1", "event2")
	require.NoError(t, err)
	require.Subset(t, []string{events[0].Key, events[1].Key}, []string{"event1", "event2"})
}

func TestCreateEventsWithNumberOfEvents(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)

	season, err := client.Seasons.Create(test_util.RequestContext(), chartKey)
	require.NoError(t, err)

	events, err := client.Seasons.CreateNumberOfEvents(test_util.RequestContext(), season.Key, 2)
	require.NoError(t, err)
	require.Len(t, events, 2)
}
