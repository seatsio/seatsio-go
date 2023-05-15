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
	Status         string      `json:"status"`
	OrderId        string      `json:"orderId"`
	ExtraData      ExtraData   `json:"extraData"`
	Label          string      `json:"label"`
	Labels         Labels      `json:"labels"`
	IDs            IDs         `json:"ids"`
	CategoryLabel  string      `json:"categoryLabel"`
	CategoryKey    CategoryKey `json:"categoryKey"`
	TicketType     string      `json:"ticketType"`
	ForSale        bool        `json:"forSale"`
	Section        string      `json:"section"`
	Entrance       string      `json:"entrance"`
	NumBooked      int         `json:"numBooked"`
	Capacity       int         `json:"capacity"`
	ObjectType     string      `json:"objectType"`
	LeftNeighbour  string      `json:"leftNeighbour"`
	RightNeighbour string      `json:"rightNeighbour"`
	HoldToken      string      `json:"holdToken"`
}

type ChangeObjectStatusResult struct {
	Objects map[string]EventObjectInfo `json:"objects"`
}

type BestAvailableResult struct {
	NextToEachOther bool                       `json:"nextToEachOther"`
	Objects         []string                   `json:"objects"`
	ObjectDetails   map[string]EventObjectInfo `json:"objectDetails"`
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

type BestAvailableStatusChangeParams struct {
	Status        string              `json:"status"`
	BestAvailable BestAvailableParams `json:"bestAvailable"`
	HoldToken     string              `json:"holdToken,omitempty"`
	OrderId       string              `json:"orderId,omitempty"`
	KeepExtraData bool                `json:"keepExtraData"`
}

type BestAvailableParams struct {
	Number      int           `json:"number"`
	Categories  []CategoryKey `json:"categories,omitempty"`
	ExtraData   []ExtraData   `json:"extraData,omitempty"`
	TicketTypes []string      `json:"ticketTypes,omitempty"`
}

type ExtraData = map[string]string

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

func (events *Events) ChangeBestAvailableObjectStatus(eventKey string, bestAvailableStatusChangeParams *BestAvailableStatusChangeParams) (*BestAvailableResult, error) {
	var bestAvailableResult BestAvailableResult
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(bestAvailableStatusChangeParams).
		SetSuccessResult(&bestAvailableResult).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/change-object-status")
	return shared.AssertOk(result, err, &bestAvailableResult)
}

type updateExtraDataRequest struct {
	ExtraData ExtraData `json:"extraData"`
}

func (events *Events) UpdateExtraData(eventKey string, objectLabel string, extraData ExtraData) error {
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	result, err := client.R().
		SetBody(&updateExtraDataRequest{
			ExtraData: extraData,
		}).
		SetQueryParam("expand", "objects").
		SetPathParam("event", eventKey).
		SetPathParam("object", objectLabel).
		Post("/events/{event}/objects/{object}/actions/update-extra-data")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) RetrieveObjectInfos(eventKey string, objectLabels []string) (map[string]EventObjectInfo, error) {
	var eventObjectInfos map[string]EventObjectInfo
	client := shared.ApiClient(events.secretKey, events.baseUrl)
	request := client.R().
		SetSuccessResult(&eventObjectInfos)
	for _, objectLabel := range objectLabels {
		request.AddQueryParam("label", objectLabel)
	}
	result, err := request.Get("/events/" + eventKey + "/objects")
	return shared.AssertOkMap(result, err, eventObjectInfos)
}

func NewEvents(secretKey string, baseUrl string) *Events {
	return &Events{secretKey, baseUrl}
}
