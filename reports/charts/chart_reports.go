package reports

import (
	"github.com/imroc/req/v3"
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
}

func (reports *ChartReports) SummaryByObjectType(chartKey string, bookWholeTablesMode string) (*ChartSummaryReport, error) {
	reportType := "byObjectType"
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
