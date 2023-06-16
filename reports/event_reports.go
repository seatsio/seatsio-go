package reports

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type EventReports struct {
	Client *req.Client
}

type EventDeepSummaryReport struct {
	Items map[string]EventDeepSummaryReportItem
}

type EventDeepSummaryReportItem struct {
	Count int `json:"count,omitempty"`
}

func (reports *EventReports) DeepSummaryByStatus(eventKey string) (*EventDeepSummaryReport, error) {
	reportType := "byStatus"
	var report map[string]EventDeepSummaryReportItem
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", eventKey).
		SetPathParam("reportType", reportType).
		Get("/reports/events/{key}/{reportType}/summary/deep")
	return shared.AssertOk(result, err, &EventDeepSummaryReport{Items: report})
}
