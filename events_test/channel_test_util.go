package events

import (
	"github.com/seatsio/seatsio-go/v2"
	"github.com/seatsio/seatsio-go/v2/events"
	"github.com/seatsio/seatsio-go/v2/test_util"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateChannel(t *testing.T, params ...*events.CreateChannelParams) (*events.Event, *seatsio.SeatsioClient) {
	company := test_util.CreateTestCompany(t)
	chartKey := test_util.CreateTestChart(t, company.Admin.SecretKey)
	c := seatsio.NewSeatsioClient(test_util.BaseUrl, company.Admin.SecretKey)
	event, _ := c.Events.Create(&events.CreateEventParams{ChartKey: chartKey, EventParams: &events.EventParams{
		EventKey: "anEvent",
	}})
	err := c.Channels.Create(event.Key, params...)
	require.NoError(t, err)
	return event, c
}
