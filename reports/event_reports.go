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

type EventSummaryReport struct {
	Items map[string]EventSummaryReportItem
}

type EventDeepSummaryReportItem struct {
	Count           int                               `json:"count,omitempty"`
	ByStatus        map[string]EventSummaryReportItem `json:"byStatus,omitempty"`
	ByCategoryKey   map[string]EventSummaryReportItem `json:"byCategoryKey,omitempty"`
	ByCategoryLabel map[string]EventSummaryReportItem `json:"byCategoryLabel,omitempty"`
	BySection       map[string]EventSummaryReportItem `json:"bySection,omitempty"`
	ByAvailability  map[string]EventSummaryReportItem `json:"byAvailability,omitempty"`
	ByChannel       map[string]EventSummaryReportItem `json:"byChannel,omitempty"`
}

type EventSummaryReportItem struct {
	Count                int            `json:"count,omitempty"`
	ByStatus             map[string]int `json:"byStatus,omitempty"`
	ByCategoryKey        map[string]int `json:"byCategoryKey,omitempty"`
	ByCategoryLabel      map[string]int `json:"byCategoryLabel,omitempty"`
	BySection            map[string]int `json:"bySection,omitempty"`
	ByAvailability       map[string]int `json:"byAvailability,omitempty"`
	ByAvailabilityReason map[string]int `json:"byAvailabilityReason,omitempty"`
	ByChannel            map[string]int `json:"byChannel,omitempty"`
}

func (reports *EventReports) SummaryByStatus(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byStatus", eventKey)
}

func (reports *EventReports) SummaryByObjectType(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byObjectType", eventKey)
}

func (reports *EventReports) SummaryByCategoryKey(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byCategoryKey", eventKey)
}

func (reports *EventReports) SummaryByCategoryLabel(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byCategoryLabel", eventKey)
}

func (reports *EventReports) SummaryBySection(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("bySection", eventKey)
}

func (reports *EventReports) SummaryByAvailability(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byAvailability", eventKey)
}

func (reports *EventReports) SummaryByAvailabilityReason(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byAvailabilityReason", eventKey)
}

func (reports *EventReports) SummaryByChannel(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byChannel", eventKey)
}

func (reports *EventReports) DeepSummaryByStatus(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byStatus", eventKey)
}

func (reports *EventReports) DeepSummaryByObjectType(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byObjectType", eventKey)
}

func (reports *EventReports) DeepSummaryByCategoryKey(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byCategoryKey", eventKey)
}

func (reports *EventReports) DeepSummaryByCategoryLabel(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byCategoryLabel", eventKey)
}

func (reports *EventReports) DeepSummaryBySection(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("bySection", eventKey)
}

func (reports *EventReports) DeepSummaryByAvailability(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byAvailability", eventKey)
}

func (reports *EventReports) DeepSummaryByAvailabilityReason(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byAvailabilityReason", eventKey)
}

func (reports *EventReports) DeepSummaryByChannel(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byChannel", eventKey)
}

func (reports *EventReports) fetchEventDeepSummaryReport(reportType string, eventKey string) (*EventDeepSummaryReport, error) {
	var report map[string]EventDeepSummaryReportItem
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", eventKey).
		SetPathParam("reportType", reportType).
		Get("/reports/events/{key}/{reportType}/summary/deep")
	return shared.AssertOk(result, err, &EventDeepSummaryReport{Items: report})
}

func (reports *EventReports) fetchEventSummaryReport(reportType string, eventKey string) (*EventSummaryReport, error) {
	var report map[string]EventSummaryReportItem
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", eventKey).
		SetPathParam("reportType", reportType).
		Get("/reports/events/{key}/{reportType}/summary")
	return shared.AssertOk(result, err, &EventSummaryReport{Items: report})
}
