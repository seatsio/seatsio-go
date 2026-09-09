package reports

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v13/events"
	"github.com/seatsio/seatsio-go/v13/shared"
)

type EventReports struct {
	Client                      *req.Client
	seasonBookingsNotPropagated bool
}

// WithSeasonBookingsNotPropagated returns a new EventReports that, when fetching a report for an
// event that's part of a season, only takes into account bookings made directly on that event -
// not bookings made on the season (or on other events in the season) that are normally propagated
// to it. The original EventReports is left unchanged.
func (reports *EventReports) WithSeasonBookingsNotPropagated() *EventReports {
	return &EventReports{Client: reports.Client, seasonBookingsNotPropagated: true}
}

func (reports *EventReports) queryParams() map[string]string {
	if reports.seasonBookingsNotPropagated {
		return map[string]string{"seasonBookingsPropagated": "false"}
	}
	return map[string]string{}
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

func (reports *EventReports) FlatList(context context.Context, eventKey string) ([]events.EventObjectInfo, error) {
	var report []events.EventObjectInfo
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("eventKey", eventKey).
		SetQueryParams(reports.queryParams()).
		Get("/reports/events/{eventKey}")
	return shared.AssertOkArray(result, err, &report)
}

func (reports *EventReports) FlatListCsv(context context.Context, eventKey string) (string, error) {
	result, err := reports.Client.R().
		SetContext(context).
		SetPathParam("eventKey", eventKey).
		SetQueryParams(reports.queryParams()).
		Get("/reports/events/{eventKey}.csv")
	err = shared.AssertOkWithoutResult(result, err)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (reports *EventReports) fetchReport(context context.Context, eventKey string, reportType string) (*DetailedEventReport, error) {
	var report map[string][]events.EventObjectInfo
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("eventKey", eventKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(reports.queryParams()).
		Get("/reports/events/{eventKey}/{reportType}")
	return shared.AssertOk(result, err, &DetailedEventReport{Items: report})
}

func (reports *EventReports) fetchReportWithFilter(context context.Context, eventKey string, reportType string, filter string) ([]events.EventObjectInfo, error) {
	var report map[string][]events.EventObjectInfo
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("eventKey", eventKey).
		SetPathParam("reportType", reportType).
		SetPathParam("filter", filter).
		SetQueryParams(reports.queryParams()).
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

func (reports *EventReports) ByAvailabilityReason(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byAvailabilityReason")
}

func (reports *EventReports) BySpecificAvailabilityReason(context context.Context, eventKey string, reason string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byAvailabilityReason", reason)
}

func (reports *EventReports) ByAvailability(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byAvailability")
}

func (reports *EventReports) BySpecificAvailability(context context.Context, eventKey string, availability string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byAvailability", availability)
}

func (reports *EventReports) ByStatus(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byStatus")
}

func (reports *EventReports) BySpecificStatus(context context.Context, eventKey string, status string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byStatus", status)
}

func (reports *EventReports) ByCategoryLabel(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byCategoryLabel")
}

func (reports *EventReports) BySpecificCategoryLabel(context context.Context, eventKey string, label string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byCategoryLabel", label)
}

func (reports *EventReports) ByCategoryKey(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byCategoryKey")
}

func (reports *EventReports) BySpecificCategoryKey(context context.Context, eventKey string, key string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byCategoryKey", key)
}

func (reports *EventReports) ByLabel(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byLabel")
}

func (reports *EventReports) BySpecificLabel(context context.Context, eventKey string, label string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byLabel", label)
}

func (reports *EventReports) ByOrderId(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byOrderId")
}

func (reports *EventReports) BySpecificOrderId(context context.Context, eventKey string, orderId string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byOrderId", orderId)
}

func (reports *EventReports) BySection(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "bySection")
}

func (reports *EventReports) BySpecificSection(context context.Context, eventKey string, section string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "bySection", section)
}

func (reports *EventReports) ByZone(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byZone")
}

func (reports *EventReports) BySpecificZone(context context.Context, eventKey string, zone string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byZone", zone)
}

func (reports *EventReports) ByChannel(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byChannel")
}

func (reports *EventReports) BySpecificChannel(context context.Context, eventKey string, channel string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byChannel", channel)
}

func (reports *EventReports) ByObjectType(context context.Context, eventKey string) (*DetailedEventReport, error) {
	return reports.fetchReport(context, eventKey, "byObjectType")
}

func (reports *EventReports) BySpecificObjectType(context context.Context, eventKey string, objectType string) ([]events.EventObjectInfo, error) {
	return reports.fetchReportWithFilter(context, eventKey, "byObjectType", objectType)
}

func (reports *EventReports) SummaryByStatus(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byStatus", eventKey)
}

func (reports *EventReports) SummaryByObjectType(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byObjectType", eventKey)
}

func (reports *EventReports) SummaryByCategoryKey(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byCategoryKey", eventKey)
}

func (reports *EventReports) SummaryByCategoryLabel(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byCategoryLabel", eventKey)
}

func (reports *EventReports) SummaryBySection(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "bySection", eventKey)
}

func (reports *EventReports) SummaryByZone(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byZone", eventKey)
}

func (reports *EventReports) SummaryByAvailability(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byAvailability", eventKey)
}

func (reports *EventReports) SummaryByAvailabilityReason(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byAvailabilityReason", eventKey)
}

func (reports *EventReports) SummaryByChannel(context context.Context, eventKey string) (*EventSummaryReport, error) {
	return reports.fetchEventSummaryReport(context, "byChannel", eventKey)
}

func (reports *EventReports) DeepSummaryByStatus(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byStatus", eventKey)
}

func (reports *EventReports) DeepSummaryByObjectType(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byObjectType", eventKey)
}

func (reports *EventReports) DeepSummaryByCategoryKey(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byCategoryKey", eventKey)
}

func (reports *EventReports) DeepSummaryByCategoryLabel(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byCategoryLabel", eventKey)
}

func (reports *EventReports) DeepSummaryBySection(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "bySection", eventKey)
}

func (reports *EventReports) DeepSummaryByZone(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byZone", eventKey)
}

func (reports *EventReports) DeepSummaryByAvailability(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byAvailability", eventKey)
}

func (reports *EventReports) DeepSummaryByAvailabilityReason(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byAvailabilityReason", eventKey)
}

func (reports *EventReports) DeepSummaryByChannel(context context.Context, eventKey string) (*EventDeepSummaryReport, error) {
	return reports.fetchEventDeepSummaryReport(context, "byChannel", eventKey)
}

func (reports *EventReports) fetchEventDeepSummaryReport(context context.Context, reportType string, eventKey string) (*EventDeepSummaryReport, error) {
	var report map[string]EventDeepSummaryReportItem
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", eventKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(reports.queryParams()).
		Get("/reports/events/{key}/{reportType}/summary/deep")
	return shared.AssertOk(result, err, &EventDeepSummaryReport{Items: report})
}

func (reports *EventReports) fetchEventSummaryReport(context context.Context, reportType string, eventKey string) (*EventSummaryReport, error) {
	var report map[string]EventSummaryReportItem
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", eventKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(reports.queryParams()).
		Get("/reports/events/{key}/{reportType}/summary")
	return shared.AssertOk(result, err, &EventSummaryReport{Items: report})
}
