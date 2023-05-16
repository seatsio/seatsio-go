package events

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

const ObjectStatusBooked = "booked"
const ObjectStatusHeld = "reservedByToken"
const ObjectStatusFree = "free"

type Events struct {
	Client *req.Client
}

type EventCreationParams struct {
	ChartKey           string                  `json:"chartKey"`
	EventKey           string                  `json:"eventKey"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         []Category              `json:"categories,omitempty"`
}

type MultipleEventCreationParams struct {
	EventKey           string                  `json:"eventKey"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         []Category              `json:"categories,omitempty"`
}

type CreateMultipleEventsRequest struct {
	ChartKey string                        `json:"chartKey"`
	Events   []MultipleEventCreationParams `json:"events"`
}

type EventCreationResult struct {
	Events []Event `json:"events"`
}

type ChangeObjectStatusResult struct {
	Objects map[string]EventObjectInfo `json:"objects"`
}

type ChangeObjectStatusInBatchResult struct {
	Results []ChangeObjectStatusResult `json:"results"`
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

type StatusChangeInBatchRequest struct {
	StatusChanges []StatusChangeInBatchParams `json:"statusChanges"`
}

type StatusChangeInBatchParams struct {
	Event                    string             `json:"event"`
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

type ForSaleConfigParams struct {
	Objects    []string       `json:"objects,omitempty"`
	AreaPlaces map[string]int `json:"areaPlaces,omitempty"`
	Categories []string       `json:"categories,omitempty"`
}

func (events *Events) Create(params *EventCreationParams) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetBody(params).
		SetSuccessResult(&event).
		Post("/events")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) CreateMultiple(chartKey string, params []MultipleEventCreationParams) (*EventCreationResult, error) {
	var eventCreationResult EventCreationResult
	result, err := events.Client.R().
		SetBody(&CreateMultipleEventsRequest{
			ChartKey: chartKey,
			Events:   params,
		}).
		SetSuccessResult(&eventCreationResult).
		Post("/events/actions/create-multiple")
	return shared.AssertOk(result, err, &eventCreationResult)
}

func (events *Events) ChangeObjectStatus(statusChangeparams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	var changeObjectStatusResult ChangeObjectStatusResult
	result, err := events.Client.R().
		SetBody(statusChangeparams).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusResult).
		Post("/events/groups/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusResult)
}

func (events *Events) ChangeObjectStatusInBatch(statusChangeInBatchParams []StatusChangeInBatchParams) (*ChangeObjectStatusInBatchResult, error) {
	var changeObjectStatusInBatchResult ChangeObjectStatusInBatchResult
	result, err := events.Client.R().
		SetBody(&StatusChangeInBatchRequest{
			StatusChanges: statusChangeInBatchParams,
		}).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusInBatchResult).
		Post("/events/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusInBatchResult)
}

func (events *Events) ChangeBestAvailableObjectStatus(eventKey string, bestAvailableStatusChangeParams *BestAvailableStatusChangeParams) (*BestAvailableResult, error) {
	var bestAvailableResult BestAvailableResult
	result, err := events.Client.R().
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
	result, err := events.Client.R().
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
	request := events.Client.R().
		SetSuccessResult(&eventObjectInfos)
	for _, objectLabel := range objectLabels {
		request.AddQueryParam("label", objectLabel)
	}
	result, err := request.Get("/events/" + eventKey + "/objects")
	return shared.AssertOkMap(result, err, eventObjectInfos)
}

func (events *Events) Delete(eventKey string) error {
	result, err := events.Client.R().
		SetQueryParam("expand", "objects").
		SetPathParam("event", eventKey).
		Delete("/events/{event}")
	return shared.AssertOkNoBody(result, err)
}

func (events *Events) Retrieve(eventKey string) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetSuccessResult(&event).
		SetPathParam("event", eventKey).
		Get("/events/{event}")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) MarkAsNotForSale(eventKey string, forSaleConfig *ForSaleConfigParams) error {
	result, err := events.Client.R().
		SetBody(forSaleConfig).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-as-not-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) MarkAsForSale(eventKey string, forSaleConfig *ForSaleConfigParams) error {
	result, err := events.Client.R().
		SetBody(forSaleConfig).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-as-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) MarkEverythingAsForSale(eventKey string) error {
	result, err := events.Client.R().
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-everything-as-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) StatusChanges(eventKey string, filter string, sortField string, sortDirection string) *shared.Lister[StatusChange] {
	pageFetcher := shared.PageFetcher[StatusChange]{
		Client:      events.Client,
		Url:         "/events/{eventKey}/status-changes",
		UrlParams:   map[string]string{"eventKey": eventKey},
		QueryParams: map[string]string{"filter": filter, "sort": shared.ToSort(sortField, sortDirection)},
	}
	return &shared.Lister[StatusChange]{PageFetcher: &pageFetcher}
}

func (events *Events) StatusChangesForObject(eventKey string, objectLabel string) *shared.Lister[StatusChange] {
	pageFetcher := shared.PageFetcher[StatusChange]{
		Client:    events.Client,
		Url:       "/events/{eventKey}/objects/{objectLabel}/status-changes",
		UrlParams: map[string]string{"eventKey": eventKey, "objectLabel": objectLabel},
	}
	return &shared.Lister[StatusChange]{PageFetcher: &pageFetcher}
}

func (events *Events) ListAll(pageSize int) ([]Event, error) {
	return events.lister().All(pageSize)
}

func (events *Events) lister() *shared.Lister[Event] {
	pageFetcher := shared.PageFetcher[Event]{
		Client:    events.Client,
		Url:       "/events",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[Event]{PageFetcher: &pageFetcher}
}

func (events *Events) ListFirstPage(pageSize int) (*shared.Page[Event], error) {
	return events.lister().ListFirstPage(pageSize)
}

func (events *Events) ListPageAfter(id int64, pageSize int) (*shared.Page[Event], error) {
	return events.lister().ListPageAfter(id, pageSize)
}

func (events *Events) ListPageBefore(id int64, pageSize int) (*shared.Page[Event], error) {
	return events.lister().ListPageBefore(id, pageSize)
}
