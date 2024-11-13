package reports

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/shared"
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
	ByZone          map[string]EventSummaryReportItem `json:"byZone,omitempty"`
	ByAvailability  map[string]EventSummaryReportItem `json:"byAvailability,omitempty"`
	ByChannel       map[string]EventSummaryReportItem `json:"byChannel,omitempty"`
}

type EventSummaryReportItem struct {
	Count                int            `json:"count,omitempty"`
	ByStatus             map[string]int `json:"byStatus,omitempty"`
	ByCategoryKey        map[string]int `json:"byCategoryKey,omitempty"`
	ByCategoryLabel      map[string]int `json:"byCategoryLabel,omitempty"`
	BySection            map[string]int `json:"bySection,omitempty"`
	ByZone               map[string]int `json:"byZone,omitempty"`
	ByAvailability       map[string]int `json:"byAvailability,omitempty"`
	ByAvailabilityReason map[string]int `json:"byAvailabilityReason,omitempty"`
	ByChannel            map[string]int `json:"byChannel,omitempty"`
}

const (
	NoOrderId    string = "NO_ORDER_ID"
	NoSection    string = "NO_SECTION"
	Available    string = "available"
	NoChannel    string = "NO_CHANNEL"
	NoCategory   string = "NO_CATEGORY"
	NotAvailable string = "not_available"
	NotForSale   string = "not_for_sale"
)

type DetailedEventReport struct {
	Items map[string][]events.EventObjectInfo
}

func (reports *EventReports) fetchReport(eventKey string, reportType string) (*DetailedEventReport, error) {
	var report map[string][]events.EventObjectInfo
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("eventKey", eventKey).
		SetPathParam("reportType", reportType).
		Get("/reports/events/{eventKey}/{reportType}")
	return shared.AssertOk(result, err, &DetailedEventReport{Items: report})
}

func (reports *EventReports) fetchReportWithFilter(eventKey string, reportType string, filter string) ([]events.EventObjectInfo, error) {
	var report map[string][]events.EventObjectInfo
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("eventKey", eventKey).
		SetPathParam("reportType", reportType).
		SetPathParam("filter", filter).
		Get("/reports/events/{eventKey}/{reportType}/{filter}")
	ok, err := shared.AssertOk(result, err, &DetailedEventReport{Items: report})
	if err == nil {
		return reports.doCast(ok).Items[filter], nil
	} else {
		return nil, err
	}
}

func (reports *EventReports) doCast(report *DetailedEventReport) *DetailedEventReport {
	return report
}

func (reports *EventReports) ByAvailabilityReason(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byAvailabilityReason")
}

func (reports *EventReports) BySpecificAvailabilityReason(eventKey string, reason string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byAvailabilityReason", reason)
}

func (reports *EventReports) ByAvailability(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byAvailability")
}

func (reports *EventReports) BySpecificAvailability(eventKey string, availability string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byAvailability", availability)
}

func (reports *EventReports) ByStatus(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byStatus")
}

func (reports *EventReports) BySpecificStatus(eventKey string, status string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byStatus", status)
}

func (reports *EventReports) ByCategoryLabel(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byCategoryLabel")
}

func (reports *EventReports) BySpecificCategoryLabel(eventKey string, label string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byCategoryLabel", label)
}

func (reports *EventReports) ByCategoryKey(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byCategoryKey")
}

func (reports *EventReports) BySpecificCategoryKey(eventKey string, key string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byCategoryKey", key)
}

func (reports *EventReports) ByLabel(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byLabel")
}

func (reports *EventReports) BySpecificLabel(eventKey string, label string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byLabel", label)
}

func (reports *EventReports) ByOrderId(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byOrderId")
}

func (reports *EventReports) BySpecificOrderId(eventKey string, orderId string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byOrderId", orderId)
}

func (reports *EventReports) BySection(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "bySection")
}

func (reports *EventReports) BySpecificSection(eventKey string, section string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "bySection", section)
}

func (reports *EventReports) ByZone(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byZone")
}

func (reports *EventReports) BySpecificZone(eventKey string, zone string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byZone", zone)
}

func (reports *EventReports) ByChannel(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byChannel")
}

func (reports *EventReports) BySpecificChannel(eventKey string, channel string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byChannel", channel)
}

func (reports *EventReports) ByObjectType(eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(eventKey, "byObjectType")
}

func (reports *EventReports) BySpecificObjectType(eventKey string, objectType string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(eventKey, "byObjectType", objectType)
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

func (reports *EventReports) SummaryByZone(eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport("byZone", eventKey)
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

func (reports *EventReports) DeepSummaryByZone(eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport("byZone", eventKey)
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
