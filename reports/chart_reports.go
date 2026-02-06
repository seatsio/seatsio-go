package reports

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/shared"
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

type Floor struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type ChartReportItem struct {
	Label                         string        `json:"label,omitempty"`
	Labels                        events.Labels `json:"labels,omitempty"`
	IDs                           events.IDs    `json:"ids,omitempty"`
	CategoryLabel                 string        `json:"categoryLabel,omitempty"`
	CategoryKey                   string        `json:"categoryKey,omitempty"`
	Section                       string        `json:"section,omitempty"`
	Entrance                      string        `json:"entrance,omitempty"`
	Capacity                      int           `json:"capacity,omitempty"`
	ObjectType                    string        `json:"objectType,omitempty"`
	LeftNeighbour                 string        `json:"leftNeighbour,omitempty"`
	RightNeighbour                string        `json:"rightNeighbour,omitempty"`
	DistanceToFocalPoint          float64       `json:"distanceToFocalPoint,omitempty"`
	BookAsAWhole                  bool          `json:"bookAsAWhole,omitempty"`
	NumSeats                      int           `json:"numSeats,omitempty"`
	IsAccessible                  bool          `json:"isAccessible,omitempty"`
	IsCompanionSeat               bool          `json:"isCompanionSeat,omitempty"`
	HasLiftUpArmrests             bool          `json:"hasLiftUpArmrests,omitempty"`
	IsHearingImpaired             bool          `json:"isHearingImpaired,omitempty"`
	IsSemiAmbulatorySeat          bool          `json:"isSemiAmbulatorySeat,omitempty"`
	HasSignLanguageInterpretation bool          `json:"hasSignLanguageInterpretation,omitempty"`
	IsPlusSize                    bool          `json:"isPlusSize,omitempty"`
	HasRestrictedView             bool          `json:"hasRestrictedView,omitempty"`
	Zone                          string        `json:"zone,omitempty"`
	Floor                         Floor         `json:"floor,omitempty"`
}

type chartReportOptions struct {
	options map[string]string
}

type chartReportOptionsOption func(params *chartReportOptions)

type ChartReportOptionsNS struct{}

var ChartReportOptions ChartReportOptionsNS

func (reports *ChartReports) SummaryByObjectType(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport(context, "byObjectType", chartKey, chartReportOptions...)
}

func (reports *ChartReports) SummaryByCategoryKey(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport(context, "byCategoryKey", chartKey, chartReportOptions...)
}

func (reports *ChartReports) SummaryByCategoryLabel(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport(context, "byCategoryLabel", chartKey, chartReportOptions...)
}

func (reports *ChartReports) SummaryBySection(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport(context, "bySection", chartKey, chartReportOptions...)
}

func (reports *ChartReports) SummaryByZone(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	return reports.fetchSummaryChartReport(context, "byZone", chartKey, chartReportOptions...)
}

func (reports *ChartReports) ByLabel(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "byLabel", chartKey, chartReportOptions...)
}

func (reports *ChartReports) ByObjectType(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "byObjectType", chartKey, chartReportOptions...)
}

func (reports *ChartReports) ByCategoryKey(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "byCategoryKey", chartKey, chartReportOptions...)
}

func (reports *ChartReports) ByCategoryLabel(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "byCategoryLabel", chartKey, chartReportOptions...)
}

func (reports *ChartReports) BySection(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "bySection", chartKey, chartReportOptions...)
}

func (reports *ChartReports) ByZone(context context.Context, chartKey string, chartReportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	return reports.fetchChartReport(context, "byZone", chartKey, chartReportOptions...)
}

func (reports *ChartReports) fetchChartReport(context context.Context, reportType string, chartKey string, reportOptions ...chartReportOptionsOption) (*ChartReport, error) {
	options := chartReportOptions{options: map[string]string{}}
	for _, opt := range reportOptions {
		opt(&options)
	}
	var report map[string][]ChartReportItem
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", chartKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(options.options).
		Get("/reports/{reportItemType}/{key}/{reportType}")
	return shared.AssertOk(result, err, &ChartReport{Items: report})
}

func (reports *ChartReports) fetchSummaryChartReport(context context.Context, reportType string, chartKey string, reportOptions ...chartReportOptionsOption) (*ChartSummaryReport, error) {
	options := chartReportOptions{options: map[string]string{}}
	for _, opt := range reportOptions {
		opt(&options)
	}
	var report map[string]ChartSummaryReportItem
	result, err := reports.Client.R().
		SetContext(context).
		SetSuccessResult(&report).
		SetPathParam("reportItemType", "charts").
		SetPathParam("key", chartKey).
		SetPathParam("reportType", reportType).
		SetQueryParams(options.options).
		Get("/reports/{reportItemType}/{key}/{reportType}/summary")
	return shared.AssertOk(result, err, &ChartSummaryReport{Items: report})
}

func (ChartReportOptionsNS) BookWholeTablesChart() chartReportOptionsOption {
	return func(options *chartReportOptions) {
		options.options["bookWholeTables"] = "chart"
	}
}

func (ChartReportOptionsNS) BookWholeTablesTrue() chartReportOptionsOption {
	return func(options *chartReportOptions) {
		options.options["bookWholeTables"] = "true"
	}
}

func (ChartReportOptionsNS) BookWholeTablesFalse() chartReportOptionsOption {
	return func(options *chartReportOptions) {
		options.options["bookWholeTables"] = "false"
	}
}

func (ChartReportOptionsNS) UseDraftVersion() chartReportOptionsOption {
	return func(options *chartReportOptions) {
		options.options["version"] = "draft"
	}
}
