package charts

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type Charts struct {
	Client *req.Client
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
