package charts

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
	"strconv"
)

type Charts struct {
	Client  *req.Client
	Archive *Archive
}

type Archive struct {
	Client *req.Client
}

type CreateChartParams struct {
	Name       string            `json:"name,omitempty"`
	VenueType  string            `json:"venueType,omitempty"`
	Categories []events.Category `json:"categories,omitempty"`
}

type UpdateChartParams struct {
	Name       string            `json:"name,omitempty"`
	Categories []events.Category `json:"categories,omitempty"`
}

type chartSupportNS struct{}

var ChartSupport chartSupportNS

func (charts *Charts) Create(params *CreateChartParams) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetBody(params).
		Post("/charts")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) Update(chartKey string, params *UpdateChartParams) error {
	result, err := charts.Client.R().
		SetBody(params).
		SetPathParam("key", chartKey).
		Post("/charts/{key}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) Retrieve(chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Get("/charts/{key}")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) RetrieveWithEvents(chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Get("/charts/{key}?expand=events")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) AddTag(chartKey string, tag string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		SetPathParam("tag", tag).
		Post("/charts/{key}/tags/{tag}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RemoveTag(chartKey string, tag string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		SetPathParam("tag", tag).
		Delete("/charts/{key}/tags/{tag}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) Copy(chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/published/actions/copy")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) CopyToWorkspace(chartKey string, targetWorkspaceKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetPathParam("chartKey", chartKey).
		SetPathParam("targetWorkspaceKey", targetWorkspaceKey).
		Post("/charts/{chartKey}/version/published/actions/copy-to-workspace/{targetWorkspaceKey}")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) CopyDraftVersion(chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/copy")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) DiscardDraftVersion(chartKey string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/discard")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) MoveToArchive(chartKey string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		Post("/charts/{key}/actions/move-to-archive")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) MoveOutOfArchive(chartKey string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		Post("/charts/{key}/actions/move-out-of-archive")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) lister() *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    charts.Client,
		Url:       "/charts",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (charts *Charts) ListAll() ([]Chart, error) {
	return charts.lister().All()
}

func (charts *Charts) List() *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    charts.Client,
		Url:       "/charts",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (charts *Charts) ListFirstPage(opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return charts.List().ListFirstPage(opts...)
}

func (charts *Charts) AddCategory(chartKey string, category events.Category) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		SetBody(category).
		Post("/charts/{key}/categories")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RemoveCategory(chartKey string, categoryKey events.CategoryKey) error {
	result, err := charts.Client.R().
		SetPathParam("chartKey", chartKey).
		SetPathParam("categoryKey", categoryKey.KeyAsString()).
		Delete("/charts/{chartKey}/categories/{categoryKey}")
	return shared.AssertOkWithoutResult(result, err)
}

type ListCategoriesResponse struct {
	Categories []events.Category `json:"categories"`
}

func (charts *Charts) ListCategories(chartKey string) (*ListCategoriesResponse, error) {
	var response ListCategoriesResponse
	result, err := charts.Client.R().
		SetSuccessResult(&response).
		SetPathParam("chartKey", chartKey).
		Get("/charts/{chartKey}/categories")
	return shared.AssertOk(result, err, &response)
}

func (charts *Charts) PublishDraftVersion(chartKey string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/publish")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RetrievePublishedVersion(chartKey string) (map[string]interface{}, error) {
	var drawing map[string]interface{}
	result, err := charts.Client.R().
		SetSuccessResult(&drawing).
		SetPathParam("key", chartKey).
		Get("/charts/{key}/version/published")
	return shared.AssertOkMap(result, err, drawing)
}

func (charts *Charts) RetrieveDraftVersion(chartKey string) (map[string]interface{}, error) {
	var drawing map[string]interface{}
	result, err := charts.Client.R().
		SetSuccessResult(&drawing).
		SetPathParam("key", chartKey).
		Get("/charts/{key}/version/draft")
	return shared.AssertOkMap(result, err, drawing)
}

func (charts *Charts) ValidatePublishedVersion(key string) (*ChartValidationResult, error) {
	var response ChartValidationResult
	result, err := charts.Client.R().
		SetSuccessResult(&response).
		SetPathParam("key", key).
		Post("/charts/{key}/version/published/actions/validate")
	return shared.AssertOk(result, err, &response)
}

func (charts *Charts) ValidateDraftVersion(key string) (*ChartValidationResult, error) {
	var response ChartValidationResult
	result, err := charts.Client.R().
		SetSuccessResult(&response).
		SetPathParam("chartKey", key).
		Post("/charts/{chartKey}/version/draft/actions/validate")
	return shared.AssertOk(result, err, &response)
}

func (archive *Archive) lister() *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    archive.Client,
		Url:       "/charts/archive",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (archive *Archive) All(opts ...shared.PaginationParamsOption) ([]Chart, error) {
	return archive.lister().All(opts...)
}

func (charts *Charts) ListAllTags() ([]string, error) {
	var tags Tags
	result, err := charts.Client.R().
		SetSuccessResult(&tags).
		Get("/charts/tags")
	ok, err := shared.AssertOk(result, err, &tags)
	if err != nil {
		return nil, err
	} else {
		return ok.Tags, nil
	}
}

func (chartSupportNS) WithFilter(filterValue string) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["filter"] = filterValue
	}
}

func (chartSupportNS) WithTag(tag string) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["tag"] = tag
	}
}

func (chartSupportNS) WithValidation(validate bool) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["validation"] = strconv.FormatBool(validate)
	}
}

func (chartSupportNS) WithExpand() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["expand"] = "events"
	}
}

func (chartSupportNS) WithEventsLimit(limit int) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["eventsLimit"] = strconv.Itoa(limit)
	}
}
