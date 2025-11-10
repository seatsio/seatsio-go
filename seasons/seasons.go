package seasons

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v11/events"
	"github.com/seatsio/seatsio-go/v11/shared"
)

type Seasons struct {
	Client *req.Client
}

type CreateSeasonParams struct {
	ChartKey           string                         `json:"chartKey"`
	Key                string                         `json:"key,omitempty"`
	Name               string                         `json:"name,omitempty"`
	TableBookingConfig *events.TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]events.CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]events.Category             `json:"categories,omitempty"`
	EventKeys          []string                       `json:"eventKeys,omitempty"`
	NumberOfEvents     int                            `json:"numberOfEvents,omitempty"`
	Channels           *[]events.CreateChannelParams  `json:"channels,omitempty"`
	ForSaleConfig      *events.ForSaleConfig          `json:"forSaleConfig,omitempty"`
	ForSalePropagated  *bool                          `json:"forSalePropagated,omitempty"`
}

type UpdateSeasonParams struct {
	EventKey           string                         `json:"eventKey,omitempty"`
	Name               string                         `json:"name,omitempty"`
	TableBookingConfig *events.TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]events.CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]events.Category             `json:"categories,omitempty"`
	Channels           *[]events.CreateChannelParams  `json:"channels,omitempty"`
	ForSalePropagated  *bool                          `json:"forSalePropagated,omitempty"`
}

type CreatePartialSeasonParams struct {
	Key       string   `json:"key,omitempty"`
	Name      string   `json:"name,omitempty"`
	EventKeys []string `json:"eventKeys,omitempty"`
}

type eventsCreationResponse struct {
	Events []*events.Event `json:"events"`
}

type createEventsParams struct {
	EventKeys      []string `json:"eventKeys,omitempty"`
	NumberOfEvents int      `json:"numberOfEvents,omitempty"`
}

type addEventsToPartialSeasonParams struct {
	EventKeys []string `json:"eventKeys"`
}

func (seasons *Seasons) Create(context context.Context, chartKey string) (*Season, error) {
	return seasons.CreateWithOptions(context, chartKey, &CreateSeasonParams{})
}

func (seasons *Seasons) CreateWithOptions(context context.Context, chartKey string, params *CreateSeasonParams) (*Season, error) {
	params.ChartKey = chartKey
	var season Season
	result, err := seasons.Client.R().
		SetContext(context).
		SetBody(params).
		SetSuccessResult(&season).
		Post("/seasons")
	return shared.AssertOk(result, err, &season)
}

func (events *Seasons) Update(context context.Context, eventKey string, params *UpdateSeasonParams) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(params).
		SetPathParam("event", eventKey).
		Post("/events/{event}")
	return shared.AssertOkWithoutResult(result, err)
}

func (seasons *Seasons) CreateEventsWithEventKeys(context context.Context, seasonKey string, eventKeys ...string) ([]*events.Event, error) {
	return seasons.createEvents(context, seasonKey, createEventsParams{EventKeys: eventKeys})
}

func (seasons *Seasons) CreateNumberOfEvents(context context.Context, seasonKey string, numberOfEvents int) ([]*events.Event, error) {
	return seasons.createEvents(context, seasonKey, createEventsParams{NumberOfEvents: numberOfEvents})
}

func (seasons *Seasons) createEvents(context context.Context, seasonKey string, params createEventsParams) ([]*events.Event, error) {
	var response eventsCreationResponse
	result, err := seasons.Client.R().
		SetContext(context).
		SetPathParam("seasonKey", seasonKey).
		SetBody(params).
		SetSuccessResult(&response).
		Post("/seasons/{seasonKey}/actions/create-events")
	ok, err := shared.AssertOk(result, err, &response)
	if err == nil {
		return ok.Events, nil
	} else {
		return nil, err
	}
}

func (seasons *Seasons) CreatePartialSeason(context context.Context, topLevelSeasonKey string) (*Season, error) {
	return seasons.CreatePartialSeasonWithOptions(context, topLevelSeasonKey, &CreatePartialSeasonParams{})
}

func (seasons *Seasons) CreatePartialSeasonWithOptions(context context.Context, topLevelSeasonKey string, params *CreatePartialSeasonParams) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetContext(context).
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetBody(params).
		SetSuccessResult(&season).
		Post("/seasons/{topLevelSeasonKey}/partial-seasons")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) RemoveEventFromPartialSeason(context context.Context, topLevelSeasonKey string, partialSeasonKey string, eventKey string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetContext(context).
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetPathParam("partialSeasonKey", partialSeasonKey).
		SetPathParam("eventKey", eventKey).
		SetSuccessResult(&season).
		Delete("/seasons/{topLevelSeasonKey}/partial-seasons/{partialSeasonKey}/events/{eventKey}")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) AddEventsToPartialSeason(context context.Context, topLevelSeasonKey string, partialSeasonKey string, eventKeys ...string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetContext(context).
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetPathParam("partialSeasonKey", partialSeasonKey).
		SetBody(addEventsToPartialSeasonParams{EventKeys: eventKeys}).
		SetSuccessResult(&season).
		Post("/seasons/{topLevelSeasonKey}/partial-seasons/{partialSeasonKey}/actions/add-events")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) Retrieve(context context.Context, key string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetContext(context).
		SetPathParam("key", key).
		SetSuccessResult(&season).
		Get("/events/{key}")
	return shared.AssertOk(result, err, &season)
}
