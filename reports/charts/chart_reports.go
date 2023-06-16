package reports

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
)

type ChartReports struct {
	Client *req.Client
}

type ChartSummaryReport struct {
	Items map[string]ChartSummaryReportItem
}

type ChartSummaryReportItem struct {
	Count           int                    `json:"count,omitempty"`
	BySection       map[string]interface{} `json:"bySection,omitempty"`
	ByCategoryKey   map[string]interface{} `json:"byCategoryKey,omitempty"`
	ByCategoryLabel map[string]interface{} `json:"byCategoryLabel,omitempty"`
	ByObjectType    map[string]interface{} `json:"byObjectType,omitempty"`
}

type ChartReport struct {
	Items map[string][]ChartReportItem
}

type ChartReportItem struct {
	Label                string        `json:"label,omitempty"`
	Labels               events.Labels `json:"labels,omitempty"`
	IDs                  events.IDs    `json:"ids,omitempty"`
	CategoryLabel        string        `json:"categoryLabel,omitempty"`
	CategoryKey          string        `json:"categoryKey,omitempty"`
	Section              string        `json:"section,omitempty"`
	Entrance             string        `json:"entrance,omitempty"`
	Capacity             int           `json:"capacity,omitempty"`
	ObjectType           string        `json:"objectType,omitempty"`
	LeftNeighbour        string        `json:"leftNeighbour,omitempty"`
	RightNeighbour       string        `json:"rightNeighbour,omitempty"`
	DistanceToFocalPoint float64       `json:"distanceToFocalPoint,omitempty"`
}

func (reports *ChartReports) SummaryByObjectType(chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport("byObjectType", chartKey, bookWholeTablesMode)
}

func (reports *ChartReports) SummaryByCategoryKey(chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport("byCategoryKey", chartKey, bookWholeTablesMode)
}

func (reports *ChartReports) SummaryByCategoryLabel(chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport("byCategoryLabel", chartKey, bookWholeTablesMode)
}

func (reports *ChartReports) SummaryBySection(chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport("bySection", chartKey, bookWholeTablesMode)
}

func (reports *ChartReports) ByLabel(chartKey string, bookWholeTablesMode string) (*ChartReport, error) {
	return reports.fetchChartReport("byLabel", chartKey, bookWholeTablesMode)
}

func (reports *ChartReports) fetchChartReport(reportType string, chartKey string, bookWholeTablesMode string) (*ChartReport, error) {
	var report map[string][]ChartReportItem
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", chartKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(toQueryParams(bookWholeTablesMode)).
		Get("/reports/{reportItemType}/{key}/{reportType}")
	return shared.AssertOk(result, err, &ChartReport{Items: report})
}

func (reports *ChartReports) fetchSummaryChartReport(reportType string, chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	var report map[string]ChartSummaryReportItem
	result, err := reports.Client.R().
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", chartKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(toQueryParams(bookWholeTablesMode)).
		Get("/reports/{reportItemType}/{key}/{reportType}/summary")
	return shared.AssertOk(result, err, &ChartSummaryReport{Items: report})
}

func toQueryParams(bookWholeTablesMode string) map[string]string {
	m := make(map[string]string)
	m["bookWholeTables"] = bookWholeTablesMode
	return m
}
