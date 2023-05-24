package charts

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type Charts struct {
	Client *req.Client
}

type UpdateChartParams struct {
	Name string `json:"name,omitempty"`
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
