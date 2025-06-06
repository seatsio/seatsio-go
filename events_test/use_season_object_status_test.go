package events

import (
	"github.com/seatsio/seatsio-go/v11"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/seasons"
	"github.com/seatsio/seatsio-go/v11/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUseSeasonObjectStatus(t *testing.T) {
	t.Parallel()
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	client := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	season, err := client.Seasons.CreateSeasonWithOptions(test_util.RequestContext(), chartKey, &seasons.CreateSeasonParams{EventKeys: []string{"event1"}})
	require.NoError(t, err)
	_, err = client.Events.Book(test_util.RequestContext(), season.Key, "A-1", "A-2")
	require.NoError(t, err)
	err = client.Events.OverrideSeasonObjectStatus(test_util.RequestContext(), "event1", "A-1", "A-2")
	require.NoError(t, err)

	err = client.Events.UseSeasonObjectStatus(test_util.RequestContext(), "event1", "A-1", "A-2")
	require.NoError(t, err)

	info, _ := client.Events.RetrieveObjectInfo(test_util.RequestContext(), "event1", "A-1", "A-2")
	require.Equal(t, events.BOOKED, info["A-1"].Status)
	require.Equal(t, events.BOOKED, info["A-2"].Status)
}
