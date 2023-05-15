package events

import "github.com/seatsio/seatsio-go/shared"

const ObjectStatusBooked = "booked"
const ObjectStatusHeld = "reservedByToken"
const ObjectStatusFree = "free"

type Events struct {
	secretKey string
	baseUrl   string
}

type EventCreationParams struct {
	ChartKey           string                  `json:"chartKey"`
	EventKey           string                  `json:"eventKey"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         []Category              `json:"categories,omitempty"`
}

func (events *Events) Create(params *EventCreationParams) (*Event, error) {
	var event Event
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(params).
		SetSuccessResult(&event).
		Post("/events")
	return shared.AssertOk(result, err, &event)
}

type IDs struct {
	Own     string `json:"own"`
	Parent  string `json:"parent"`
	Section string `json:"section"`
}

type Labels struct {
	Own     LabelAndType `json:"own"`
	Parent  LabelAndType `json:"parent"`
	Section string       `json:"section"`
}

type LabelAndType struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type EventObjectInfo struct {
	Status         string            `json:"status"`
	OrderId        string            `json:"orderId"`
	ExtraData      map[string]string `json:"extraData"`
	Label          string            `json:"label"`
	Labels         Labels            `json:"labels"`
	IDs            IDs               `json:"ids"`
	CategoryLabel  string            `json:"categoryLabel"`
	CategoryKey    CategoryKey       `json:"categoryKey"`
	TicketType     string            `json:"ticketType"`
	ForSale        bool              `json:"forSale"`
	Section        string            `json:"section"`
	Entrance       string            `json:"entrance"`
	NumBooked      int               `json:"numBooked"`
	Capacity       int               `json:"capacity"`
	ObjectType     string            `json:"objectType"`
	LeftNeighbour  string            `json:"leftNeighbour"`
	RightNeighbour string            `json:"rightNeighbour"`
	HoldToken      string            `json:"holdToken"`
}

type ChangeObjectStatusResult struct {
	Objects map[string]EventObjectInfo `json:"objects"`
}

type StatusChangeParams struct {
	Events                   []string           `json:"events"`
	Status                   string             `json:"status"`
	Objects                  []ObjectProperties `json:"objects"`
	HoldToken                string             `json:"holdToken,omitempty"`
	OrderId                  string             `json:"orderId,omitempty"`
	KeepExtraData            bool               `json:"keepExtraData"`
	AllowedPreviousStatuses  []string           `json:"allowedPreviousStatuses,omitempty"`
	RejectedPreviousStatuses []string           `json:"rejectedPreviousStatuses,omitempty"`
}

func (events *Events) ChangeObjectStatus(statusChangeparams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	var changeObjectStatusResult ChangeObjectStatusResult
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(statusChangeparams).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusResult).
		Post("/events/groups/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusResult)
}

type updateExtraDataRequest struct {
	ExtraData map[string]string `json:"extraData"`
}

func (events *Events) UpdateExtraData(eventKey string, objectLabel string, extraData map[string]string) error {
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(&updateExtraDataRequest{
			ExtraData: extraData,
		}).
		SetQueryParam("expand", "objects").
		// TODO: encode url args
		Post("/events/" + eventKey + "/objects/" + objectLabel + "/actions/update-extra-data")
	return shared.AssertOkWithoutResult(result, err)
}

func NewEvents(secretKey string, baseUrl string) *Events {
	return &Events{secretKey, baseUrl}
}
