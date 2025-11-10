package charts

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v12/events"
	"github.com/seatsio/seatsio-go/v12/shared"
	"os"
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

type UpdateCategoryParams struct {
	Label      string `json:"label,omitempty"`
	Color      string `json:"color,omitempty"`
	Accessible bool   `json:"accessible,omitempty"`
}

type chartSupportNS struct{}

var ChartSupport chartSupportNS

func (charts *Charts) Create(context context.Context, params *CreateChartParams) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetBody(params).
		Post("/charts")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) Update(context context.Context, chartKey string, params *UpdateChartParams) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetBody(params).
		SetPathParam("key", chartKey).
		Post("/charts/{key}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) Retrieve(context context.Context, chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Get("/charts/{key}")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) RetrieveWithEvents(context context.Context, chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Get("/charts/{key}?expand=events")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) AddTag(context context.Context, chartKey string, tag string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		SetPathParam("tag", tag).
		Post("/charts/{key}/tags/{tag}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RemoveTag(context context.Context, chartKey string, tag string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		SetPathParam("tag", tag).
		Delete("/charts/{key}/tags/{tag}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) Copy(context context.Context, chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/published/actions/copy")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) CopyToWorkspace(context context.Context, chartKey string, toWorkspaceKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("chartKey", chartKey).
		SetPathParam("toWorkspaceKey", toWorkspaceKey).
		Post("/charts/{chartKey}/version/published/actions/copy-to-workspace/{toWorkspaceKey}")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) CopyFromWorkspaceTo(context context.Context, chartKey string, fromWorkspaceKey string, toWorkspaceKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("chartKey", chartKey).
		SetPathParam("fromWorkspaceKey", fromWorkspaceKey).
		SetPathParam("toWorkspaceKey", toWorkspaceKey).
		Post("/charts/{chartKey}/version/published/actions/copy/from/{fromWorkspaceKey}/to/{toWorkspaceKey}")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) CopyDraftVersion(context context.Context, chartKey string) (*Chart, error) {
	var chart Chart
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&chart).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/copy")
	return shared.AssertOk(result, err, &chart)
}

func (charts *Charts) DiscardDraftVersion(context context.Context, chartKey string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/discard")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) MoveToArchive(context context.Context, chartKey string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/actions/move-to-archive")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) MoveOutOfArchive(context context.Context, chartKey string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/actions/move-out-of-archive")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) lister(context context.Context) *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    charts.Client,
		Url:       "/charts",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (charts *Charts) ListAll(context context.Context) ([]Chart, error) {
	return charts.lister(context).All()
}

func (charts *Charts) List(context context.Context) *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    charts.Client,
		Url:       "/charts",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (charts *Charts) ListFirstPage(context context.Context, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return charts.List(context).ListFirstPage(opts...)
}

func (charts *Charts) ListPageAfter(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return charts.List(context).ListPageAfter(id, opts...)
}

func (charts *Charts) ListPageBefore(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return charts.List(context).ListPageBefore(id, opts...)
}

func (charts *Charts) AddCategory(context context.Context, chartKey string, category events.Category) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		SetBody(category).
		Post("/charts/{key}/categories")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RemoveCategory(context context.Context, chartKey string, categoryKey events.CategoryKey) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("chartKey", chartKey).
		SetPathParam("categoryKey", categoryKey.KeyAsString()).
		Delete("/charts/{chartKey}/categories/{categoryKey}")
	return shared.AssertOkWithoutResult(result, err)
}

type listCategoriesResponse struct {
	Categories []events.Category `json:"categories"`
}

func (charts *Charts) ListCategories(context context.Context, chartKey string) ([]events.Category, error) {
	var response listCategoriesResponse
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&response).
		SetPathParam("chartKey", chartKey).
		Get("/charts/{chartKey}/categories")
	return shared.AssertOkArray(result, err, &response.Categories)
}

func (charts *Charts) UpdateCategory(context context.Context, chartKey string, categoryKey events.CategoryKey, params UpdateCategoryParams) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("chartKey", chartKey).
		SetPathParam("categoryKey", categoryKey.KeyAsString()).
		SetBody(params).
		Post("/charts/{chartKey}/categories/{categoryKey}")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) PublishDraftVersion(context context.Context, chartKey string) error {
	result, err := charts.Client.R().
		SetContext(context).
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/publish")
	return shared.AssertOkWithoutResult(result, err)
}

func (charts *Charts) RetrievePublishedVersion(context context.Context, chartKey string) (map[string]interface{}, error) {
	var drawing map[string]interface{}
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&drawing).
		SetPathParam("key", chartKey).
		Get("/charts/{key}/version/published")
	return shared.AssertOkMap(result, err, drawing)
}

func (charts *Charts) RetrieveDraftVersion(context context.Context, chartKey string) (map[string]interface{}, error) {
	var drawing map[string]interface{}
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&drawing).
		SetPathParam("key", chartKey).
		Get("/charts/{key}/version/draft")
	return shared.AssertOkMap(result, err, drawing)
}

func (charts *Charts) ValidatePublishedVersion(context context.Context, key string) (*ChartValidationResult, error) {
	var response ChartValidationResult
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&response).
		SetPathParam("key", key).
		Post("/charts/{key}/version/published/actions/validate")
	return shared.AssertOk(result, err, &response)
}

func (charts *Charts) ValidateDraftVersion(context context.Context, key string) (*ChartValidationResult, error) {
	var response ChartValidationResult
	result, err := charts.Client.R().
		SetContext(context).
		SetSuccessResult(&response).
		SetPathParam("chartKey", key).
		Post("/charts/{chartKey}/version/draft/actions/validate")
	return shared.AssertOk(result, err, &response)
}

func (charts *Charts) RetrievePublishedVersionThumbnail(context context.Context, chartKey string) (*os.File, error) {
	return charts.retrieveThumbnail(context, "published", chartKey)
}

func (charts *Charts) RetrieveDraftVersionThumbnail(context context.Context, chartKey string) (*os.File, error) {
	return charts.retrieveThumbnail(context, "draft", chartKey)
}

func (charts *Charts) retrieveThumbnail(context context.Context, imageType string, chartKey string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", chartKey+".png.")
	result, err := charts.Client.R().
		SetContext(context).
		SetOutputFile(tempFile.Name()).
		SetPathParam("key", chartKey).
		SetPathParam("imageType", imageType).
		Get("/charts/{key}/version/{imageType}/thumbnail")
	return shared.AssertOk(result, err, tempFile)
}

func (archive *Archive) lister(context context.Context) *shared.Lister[Chart] {
	pageFetcher := shared.PageFetcher[Chart]{
		Client:    archive.Client,
		Url:       "/charts/archive",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[Chart]{PageFetcher: &pageFetcher}
}

func (archive *Archive) All(context context.Context, opts ...shared.PaginationParamsOption) ([]Chart, error) {
	return archive.lister(context).All(opts...)
}

func (archive *Archive) ListFirstPage(context context.Context, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return archive.lister(context).ListFirstPage(opts...)
}

func (archive *Archive) ListPageAfter(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return archive.lister(context).ListPageAfter(id, opts...)
}

func (archive *Archive) ListPageBefore(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Chart], error) {
	return archive.lister(context).ListPageBefore(id, opts...)
}

func (charts *Charts) ListAllTags(context context.Context) ([]string, error) {
	var tags Tags
	result, err := charts.Client.R().
		SetContext(context).
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

// WithValidation Deprecated
func (chartSupportNS) WithValidation(validate bool) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["validation"] = strconv.FormatBool(validate)
	}
}

// WithExpand Deprecated
func (chartSupportNS) WithExpand() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["expand"] = "events"
	}
}

func (chartSupportNS) WithExpandEvents() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.AddToArrayQueryParam("expand", "events")
	}
}

func (chartSupportNS) WithExpandValidation() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.AddToArrayQueryParam("expand", "validation")
	}
}

func (chartSupportNS) WithExpandVenueType() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.AddToArrayQueryParam("expand", "venueType")
	}
}

func (chartSupportNS) WithExpandZones() shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.AddToArrayQueryParam("expand", "zones")
	}
}

func (chartSupportNS) WithEventsLimit(limit int) shared.PaginationParamsOption {
	return func(Params *shared.PaginationParams) {
		Params.QueryParams["eventsLimit"] = strconv.Itoa(limit)
	}
}
