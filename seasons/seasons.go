package seasons

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v9/events"
	"github.com/seatsio/seatsio-go/v9/shared"
)

type Seasons struct {
	Client *req.Client
}

type CreateSeasonParams struct {
	ChartKey           string                        `json:"chartKey"`
	Key                string                        `json:"key,omitempty"`
	Name               string                        `json:"name,omitempty"`
	TableBookingConfig *events.TableBookingConfig    `json:"tableBookingConfig,omitempty"`
	EventKeys          []string                      `json:"eventKeys,omitempty"`
	NumberOfEvents     int                           `json:"numberOfEvents,omitempty"`
	Channels           *[]events.CreateChannelParams `json:"channels,omitempty"`
	ForSaleConfig      *events.ForSaleConfig         `json:"forSaleConfig,omitempty"`
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

func (seasons *Seasons) CreateSeason(chartKey string) (*Season, error) {
	return seasons.CreateSeasonWithOptions(chartKey, &CreateSeasonParams{})
}

func (seasons *Seasons) CreateSeasonWithOptions(chartKey string, params *CreateSeasonParams) (*Season, error) {
	params.ChartKey = chartKey
	var season Season
	result, err := seasons.Client.R().
		SetBody(params).
		SetSuccessResult(&season).
		Post("/seasons")
	return shared.AssertOk(result, err, &season)
}

func (seasons *Seasons) CreateEventsWithEventKeys(seasonKey string, eventKeys ...string) ([]*events.Event, error) {
	return seasons.createEvents(seasonKey, createEventsParams{EventKeys: eventKeys})
}

func (seasons *Seasons) CreateNumberOfEvents(seasonKey string, numberOfEvents int) ([]*events.Event, error) {
	return seasons.createEvents(seasonKey, createEventsParams{NumberOfEvents: numberOfEvents})
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

func (seasons *Seasons) CreatePartialSeason(topLevelSeasonKey string) (*Season, error) {
	return seasons.CreatePartialSeasonWithOptions(topLevelSeasonKey, &CreatePartialSeasonParams{})
}

func (seasons *Seasons) CreatePartialSeasonWithOptions(topLevelSeasonKey string, params *CreatePartialSeasonParams) (*Season, error) {
	var season Season
	result, err := seasons.Client.R().
		SetPathParam("topLevelSeasonKey", topLevelSeasonKey).
		SetBody(params).
		SetSuccessResult(&season).
		Post("/seasons/{topLevelSeasonKey}/partial-seasons")
	return shared.AssertOk(result, err, &season)
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
