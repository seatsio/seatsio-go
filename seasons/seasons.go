package seasons

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/shared"
)

type Seasons struct {
	Client *req.Client
}

type createSeasonParams struct {
	ChartKey           string                     `json:"chartKey"`
	Key                *string                    `json:"key,omitempty"`
	TableBookingConfig *events.TableBookingConfig `json:"tableBookingConfig,omitempty"`
	EventKeys          *[]string                  `json:"eventKeys,omitempty"`
	NumberOfEvents     *int32                     `json:"numberOfEvents,omitempty"`
}

type CreateSeasonParamsOption func(Params *createSeasonParams)

type createPartialSeasonParams struct {
	Key       *string   `json:"key,omitempty"`
	EventKeys *[]string `json:"eventKeys,omitempty"`
}

type CreatePartialSeasonParamsOption func(Params *createPartialSeasonParams)

type eventsCreationResponse struct {
	Events []*events.Event `json:"events"`
}

type createEventsParams struct {
	EventKeys      *[]string `json:"eventKeys,omitempty"`
	NumberOfEvents *int32    `json:"numberOfEvents,omitempty"`
}

type addEventsToPartialSeasonParams struct {
	EventKeys []string `json:"eventKeys"`
}

type seasonSupportNS struct{}

var SeasonSupport seasonSupportNS

type partialSeasonSupportNS struct{}

var PartialSeasonSupport partialSeasonSupportNS

func (seasons *Seasons) CreateSeason(chartKey string, opts ...CreateSeasonParamsOption) (*Season, error) {

	var season Season
	result, err := seasons.Client.R().
		SetBody(seasons.newCreateSeasonParams(chartKey, opts...)).
		SetSuccessResult(&season).
		Post("/seasons")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) newCreateSeasonParams(chartKey string, opts ...CreateSeasonParamsOption) *createSeasonParams {
	params := &createSeasonParams{
		ChartKey:           chartKey,
		Key:                nil,
		EventKeys:          nil,
		NumberOfEvents:     nil,
		TableBookingConfig: nil,
	}
	for _, opt := range opts {
		opt(params)
	}
	return params
}

func (seasons *Seasons) CreateEventsWithEventKeys(seasonKey string, eventKeys ...string) ([]*events.Event, error) {
	return seasons.createEvents(seasonKey, createEventsParams{EventKeys: &eventKeys, NumberOfEvents: nil})
}

func (seasons *Seasons) CreateNumberOfEvents(seasonKey string, numberOfEvents int32) ([]*events.Event, error) {
	return seasons.createEvents(seasonKey, createEventsParams{EventKeys: nil, NumberOfEvents: &numberOfEvents})
}

func (seasons *Seasons) createEvents(seasonKey string, params createEventsParams) ([]*events.Event, error) {
	var response eventsCreationResponse
	result, err := seasons.Client.R().
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

func (seasons *Seasons) CreatePartialSeason(topLevelSeasonKey string, opts ...CreatePartialSeasonParamsOption) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetBody(seasons.newCreatePartialSeasonParams(opts...)).
		SetSuccessResult(&season).
		Post("/seasons/{topLevelSeasonKey}/partial-seasons")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) newCreatePartialSeasonParams(opts ...CreatePartialSeasonParamsOption) *createPartialSeasonParams {
	params := &createPartialSeasonParams{
		Key:       nil,
		EventKeys: nil,
	}
	for _, opt := range opts {
		opt(params)
	}
	return params
}

func (seasons *Seasons) RemoveEventFromPartialSeason(topLevelSeasonKey string, partialSeasonKey string, eventKey string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetPathParam("partialSeasonKey", partialSeasonKey).
		SetPathParam("eventKey", eventKey).
		SetSuccessResult(&season).
		Delete("/seasons/{topLevelSeasonKey}/partial-seasons/{partialSeasonKey}/events/{eventKey}")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) AddEventsToPartialSeason(topLevelSeasonKey string, partialSeasonKey string, eventKeys ...string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetPathParam("partialSeasonKey", partialSeasonKey).
		SetBody(addEventsToPartialSeasonParams{EventKeys: eventKeys}).
		SetSuccessResult(&season).
		Post("/seasons/{topLevelSeasonKey}/partial-seasons/{partialSeasonKey}/actions/add-events")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) Retrieve(key string) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetPathParam("key", key).
		SetSuccessResult(&season).
		Get("/events/{key}")
	return shared.AssertOk(result, err, &season)
}

func (seasonSupportNS) WithKey(key string) CreateSeasonParamsOption {
	return func(params *createSeasonParams) {
		params.Key = &key
	}
}

func (seasonSupportNS) WithNumberOfEvents(numberOfEvents int32) CreateSeasonParamsOption {
	return func(params *createSeasonParams) {
		params.NumberOfEvents = &numberOfEvents
		params.EventKeys = nil
	}
}

func (seasonSupportNS) WithEventKeys(eventKeys ...string) CreateSeasonParamsOption {
	return func(params *createSeasonParams) {
		params.EventKeys = &eventKeys
		params.NumberOfEvents = nil
	}
}

func (partialSeasonSupportNS) WithKey(key string) CreatePartialSeasonParamsOption {
	return func(params *createPartialSeasonParams) {
		params.Key = &key
	}
}

func (partialSeasonSupportNS) WithEventKeys(eventKeys ...string) CreatePartialSeasonParamsOption {
	return func(params *createPartialSeasonParams) {
		params.EventKeys = &eventKeys
	}
}
