package charts

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
)

type Charts struct {
	Client *req.Client
}

type CreateChartParams struct {
	Name       string            `json:"name,omitempty"`
	VenueType  string            `json:"venueType,omitempty"`
	Categories []events.Category `json:"categories,omitempty"`
}

type UpdateChartParams struct {
	Name string `json:"name,omitempty"`
}

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

func (charts *Charts) AddTag(chartKey string, tag string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		SetPathParam("tag", tag).
		Post("/charts/{key}/tags/{tag}")
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

func (charts *Charts) RetrievePublishedVersion(chartKey string) (map[string]interface{}, error) {
	var drawing map[string]interface{}
	result, err := charts.Client.R().
		SetSuccessResult(&drawing).
		SetPathParam("key", chartKey).
		Get("/charts/{key}/version/published")
	return shared.AssertOkMap(result, err, drawing)
}

func (charts *Charts) DiscardDraftVersion(chartKey string) error {
	result, err := charts.Client.R().
		SetPathParam("key", chartKey).
		Post("/charts/{key}/version/draft/actions/discard")
	return shared.AssertOkWithoutResult(result, err)
}

/*  TODO
func (charts *Charts) ListAllTags() (*[]string, error) {
	var tags Tags
	result, err := charts.Client.R().
		SetSuccessResult(&tags).
		Get("/charts/tags")
	return shared.AssertOk(result, err, &tags.Tags)
}
*/
