package test_util

import (
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v8/shared"
	"log"
	"os"
	"testing"
)

var client = req.C()

const BaseUrl = "https://api-staging-eu.seatsio.net"

type User struct {
	SecretKey string `json:"secretKey"`
}

type Workspace struct {
	Key       string `json:"key"`
	SecretKey string `json:"secretKey"`
}

type TestCompany struct {
	Admin     User      `json:"admin"`
	Workspace Workspace `json:"workspace"`
}

func CreateTestCompany(t *testing.T) *TestCompany {
	var testCompany TestCompany
	result, err := req.C().
		SetBaseURL(BaseUrl).
		R().
		SetSuccessResult(&testCompany).
		Post("/system/public/users/actions/create-test-company")
	_, e := shared.AssertOk(result, err, &testCompany)
	if e != nil {
		t.Fatalf("unable to create test company: #{err}")
	}
	return &testCompany
}

func CreateTestChart(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChart.json")
}

func CreateTestChartWithTables(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChartWithTables.json")
}

func CreateTestChartWithSections(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChartWithSections.json")
}

func CreateTestChartWithFloors(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChartWithFloors.json")
}

func CreateTestChartWithZones(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChartWithZones.json")
}

func CreateTestChartWithErrors(t *testing.T, secretKey string) string {
	return createTestChart(t, secretKey, "sampleChartWithErrors.json")
}

func createTestChart(t *testing.T, secretKey string, fileName string) string {
	chartJson, err := os.ReadFile("../test_util/charts/" + fileName)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	chartKey := uuid.New().String()
	result, err := shared.ApiClient(secretKey, BaseUrl).R().
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

func DemoCompanySecretKey() string {
	return os.Getenv("DEMO_COMPANY_SECRET_KEY")
}

func AssertDemoCompanySecretKeySet(t *testing.T) {
	_, exists := os.LookupEnv("DEMO_COMPANY_SECRET_KEY")
	if !exists {
		t.Skip("DEMO_COMPANY_SECRET_KEY environment variable not set, skipping test")
	}
}
