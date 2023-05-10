package events_test

import (
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var client = req.C()

const secretKey = "demoKey"
const baseUrl = "https://api-staging-eu.seatsio.net"

func testClient() *seatsio.SeatsioClient {
	return seatsio.NewSeatsioClient(secretKey, baseUrl)
}

func TestChartKeyIsRequired(t *testing.T) {
	chartKey := createTestChart(t)
	client := testClient()

	event := client.Events.Create(chartKey)

	assert.NotZero(t, event.Id)
}

func createTestChart(t *testing.T) string {
	chartJson, err := os.ReadFile("sampleChart.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	chartKey := uuid.New().String()
	result, err := client.R().
		SetBasicAuth(secretKey, "").
		SetBody(string(chartJson)).
		Post(baseUrl + "/system/public/charts/" + chartKey)
	if err != nil {
		t.Fatalf("unable to create test chart: %v", err)
	}
	if result.IsErrorState() {
		t.Fatalf("unable to create test chart: %v", result.String())
	}
	return chartKey
}
