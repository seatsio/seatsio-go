package reports

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v8/shared"
	"strconv"
	"time"
)

type UsageReports struct {
	Client *req.Client
}

type UsageReason string

const (
	StatusChanged     UsageReason = "STATUS_CHANGED"
	Selected          UsageReason = "SELECTED"
	AssignedToChannel UsageReason = "ASSIGNED_TO_CHANNEL"
)

type UsageForObjectV1 struct {
	Object                       string    `json:"object"`
	NumFirstBookings             int       `json:"numFirstBookings"`
	FirstBookingDate             time.Time `json:"firstBookingDate"`
	NumFirstSelections           int       `json:"numFirstSelections"`
	NumFirstBookingsOrSelections int       `json:"numFirstBookingsOrSelections"`
}

type UsageForObjectV2 struct {
	Object         string              `json:"object"`
	NumUsedObjects int                 `json:"numUsedObjects"`
	UsageByReason  map[UsageReason]int `json:"usageByReason"`
}

type UsageChart struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type UsageEvent struct {
	Id  int64  `json:"id"`
	Key string `json:"key"`
}

type UsageByEvent struct {
	Event          UsageEvent `json:"event"`
	NumUsedObjects int        `json:"numUsedObjects"`
}

type UsageByChart struct {
	Chart        UsageChart     `json:"chart"`
	UsageByEvent []UsageByEvent `json:"usageByEvent"`
}

type UsageDetails struct {
	Workspace    int64          `json:"workspace"`
	UsageByChart []UsageByChart `json:"usageByChart"`
}

type Month struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type UsageSummaryForAllMonths struct {
	Usage           []UsageSummaryForMonth `json:"usage"`
	UsageCutoffDate *time.Time             `json:"usageCutoffDate"`
}

type UsageSummaryForMonth struct {
	Month          Month `json:"month"`
	NumUsedObjects int   `json:"numUsedObjects"`
}

func (usageReports *UsageReports) SummaryForAllMonths() (*UsageSummaryForAllMonths, error) {
	var report UsageSummaryForAllMonths
	result, err := usageReports.Client.R().
		SetSuccessResult(&report).
		Get("/reports/usage?version=2")
	return shared.AssertOk(result, err, &report)
}

func (usageReports *UsageReports) DetailsForMonth(year int, month int) ([]UsageDetails, error) {
	var details []UsageDetails
	result, err := usageReports.Client.R().
		SetSuccessResult(&details).
		SetPathParam("month", formatMonth(year, month)).
		Get("/reports/usage/month/{month}")
	return shared.AssertOkArray(result, err, &details)
}

func (usageReports *UsageReports) DetailsForEventInMonth(eventId int, year int, month int) ([]UsageForObjectV1, []UsageForObjectV2, error) {
	result, err := usageReports.Client.R().
		SetPathParam("month", formatMonth(year, month)).
		SetPathParam("eventId", strconv.Itoa(eventId)).
		Get("/reports/usage/month/{month}/event/{eventId}")
	err = shared.AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, nil, err
	}

	var jsonContent []interface{}
	bodyAsBytes := result.Bytes()
	_ = json.Unmarshal(bodyAsBytes, &jsonContent)
	if len(jsonContent) == 0 {
		return []UsageForObjectV1{}, []UsageForObjectV2{}, nil
	}

	var sampleUsageMap = jsonContent[0].(map[string]interface{})
	_, containsKey := sampleUsageMap["usageByReason"]
	if containsKey {
		var results []UsageForObjectV2
		err = json.Unmarshal(bodyAsBytes, &results)
		return nil, results, nil
	} else {
		var results []UsageForObjectV1
		err = json.Unmarshal(bodyAsBytes, &results)
		return results, nil, nil
	}
}

func formatMonth(year int, month int) string {
	return fmt.Sprintf("%d-%02d", year, month)
}
