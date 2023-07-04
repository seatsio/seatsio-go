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

type CreateEventParams struct {
	ChartKey           string                  `json:"chartKey"`
	EventKey           string                  `json:"eventKey"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]Category             `json:"categories,omitempty"`
}

type UpdateEventParams struct {
	ChartKey           string                  `json:"chartKey,omitempty"`
	EventKey           string                  `json:"eventKey,omitempty"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]Category             `json:"categories,omitempty"`
}

type CreateMultipleEventsParams struct {
	EventKey           string                  `json:"eventKey"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]Category             `json:"categories,omitempty"`
}

type CreateMultipleEventsRequest struct {
	ChartKey string                       `json:"chartKey"`
	Events   []CreateMultipleEventsParams `json:"events"`
}

type CreateEventResult struct {
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

type updateExtraDatasRequest struct {
	ExtraData map[string]ExtraData `json:"extraData"`
}

func (events *Events) Create(params *CreateEventParams) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetBody(params).
		SetSuccessResult(&event).
		Post("/events")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) CreateMultiple(chartKey string, params []CreateMultipleEventsParams) (*CreateEventResult, error) {
	var eventCreationResult CreateEventResult
	result, err := events.Client.R().
		SetBody(&CreateMultipleEventsRequest{
			ChartKey: chartKey,
			Events:   params,
		}).
		SetSuccessResult(&eventCreationResult).
		Post("/events/actions/create-multiple")
	return shared.AssertOk(result, err, &eventCreationResult)
}

func (events *Events) Update(eventKey string, params *UpdateEventParams) error {
	result, err := events.Client.R().
		SetBody(params).
		SetPathParam("event", eventKey).
		Post("/events/{event}")
	return shared.AssertOkWithoutResult(result, err)
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

func (events *Events) UpdateExtraDatas(eventKey string, extraData map[string]ExtraData) error {
	result, err := events.Client.R().
		SetBody(&updateExtraDatasRequest{
			ExtraData: extraData,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/update-extra-data")
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

func (events *Events) ListAll(opts ...shared.PaginationParamsOption) ([]Event, error) {
	return events.lister().All(opts...)
}

func (events *Events) lister() *shared.Lister[Event] {
	pageFetcher := shared.PageFetcher[Event]{
		Client:    events.Client,
		Url:       "/events",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[Event]{PageFetcher: &pageFetcher}
}

func (events *Events) ListFirstPage(opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister().ListFirstPage(opts...)
}

func (events *Events) ListPageAfter(id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister().ListPageAfter(id, opts...)
}

func (events *Events) ListPageBefore(id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister().ListPageBefore(id, opts...)
}
