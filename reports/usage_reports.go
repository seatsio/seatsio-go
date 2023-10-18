package reports

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v2/shared"
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
	NumFirstBookings             int32     `json:"numFirstBookings"`
	FirstBookingDate             time.Time `json:"firstBookingDate"`
	NumFirstSelections           int32     `json:"numFirstSelections"`
	NumFirstBookingsOrSelections int32     `json:"numFirstBookingsOrSelections"`
}

type UsageForObjectV2 struct {
	Object         string                `json:"object"`
	NumUsedObjects int32                 `json:"numUsedObjects"`
	UsageByReason  map[UsageReason]int32 `json:"usageByReason"`
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
	NumUsedObjects int32      `json:"numUsedObjects"`
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
	Month int32 `json:"month"`
	Year  int32 `json:"year"`
}

type UsageSummaryForMonth struct {
	Month          Month `json:"month"`
	NumUsedObjects int32 `json:"numUsedObjects"`
}

func (usageReports *UsageReports) SummaryForAllMonths() ([]UsageSummaryForMonth, error) {
	var summaries []UsageSummaryForMonth
	result, err := usageReports.Client.R().
		SetSuccessResult(&summaries).
		Get("/reports/usage")
	return shared.AssertOkArray(result, err, &summaries)
}

func (usageReports *UsageReports) DetailsForMonth(year int, month int) ([]UsageDetails, error) {
	var details []UsageDetails
	result, err := usageReports.Client.R().
		SetSuccessResult(&details).
		SetPathParam("month", formatMonth(year, month)).
		Get("/reports/usage/month/{month}")
	return shared.AssertOkArray(result, err, &details)
}

func (usageReports *UsageReports) DetailsForEventInMonth(year int, month int) ([]UsageForObjectV1, []UsageForObjectV2, error) {
	result, err := usageReports.Client.R().
		SetPathParam("month", formatMonth(year, month)).
		Get("/reports/usage/month/{month}")
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
