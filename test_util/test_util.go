package test_util

import (
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"log"
	"os"
	"testing"
)

var client = req.C()

const SecretKey = "demoKey"
const BaseUrl = "https://api-staging-eu.seatsio.net"

func CreateTestChart(t *testing.T) string {
	chartJson, err := os.ReadFile("sampleChart.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	chartKey := uuid.New().String()
	result, err := client.R().
		SetBasicAuth(SecretKey, "").
		SetBody(string(chartJson)).
		Post(BaseUrl + "/system/public/charts/" + chartKey)
	if err != nil {
		t.Fatalf("unable to create test chart: %v", err)
	}
	if result.IsErrorState() {
		t.Fatalf("unable to create test chart: %v", result.String())
	}
	return chartKey
}
