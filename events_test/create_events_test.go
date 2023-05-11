package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChartKeyIsRequired(t *testing.T) {
	chartKey := test_util.CreateTestChart(t)
	client := seatsio.NewSeatsioClient(test_util.SecretKey, test_util.BaseUrl)

	event, _ := client.Events.Create(chartKey)

	assert.NotZero(t, event.Id)
}
