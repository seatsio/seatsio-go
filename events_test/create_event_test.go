package events_test

import (
	"github.com/seatsio/seatsio-go"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/test_util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChartKeyIsRequired(t *testing.T) {
	start := time.Now()
	chartKey := test_util.CreateTestChart(t)
	client := seatsio.NewSeatsioClient(test_util.SecretKey, test_util.BaseUrl)

	event, _ := client.Events.Create(chartKey)

	assert.NotZero(t, event.Id)
	assert.NotNil(t, event.Key)
	assert.Equal(t, event.ChartKey, chartKey)
	assert.Equal(t, event.TableBookingConfig, events.TableBookingConfig{Mode: "INHERIT"})
	assert.True(t, event.SupportsBestAvailable)
	assert.Nil(t, event.ForSaleConfig)
	assert.True(t, event.CreatedOn.After(start))
	assert.Nil(t, event.UpdatedOn)
	assert.Equal(t, event.Categories, []events.Category{
		{Key: events.CategoryKey{Key: 9}, Label: "Cat1", Color: "#87A9CD", Accessible: false},
		{Key: events.CategoryKey{Key: 10}, Label: "Cat2", Color: "#5E42ED", Accessible: false},
		{Key: events.CategoryKey{Key: "string11"}, Label: "Cat3", Color: "#5E42BB", Accessible: false},
	})
}
