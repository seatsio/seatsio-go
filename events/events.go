package events

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v6/shared"
	"time"
)

type Events struct {
	Client *req.Client
}

type EventParams struct {
	EventKey           string                  `json:"eventKey,omitempty"`
	Name               string                  `json:"name,omitempty"`
	Date               string                  `json:"date,omitempty"`
	TableBookingConfig *TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	ObjectCategories   *map[string]CategoryKey `json:"objectCategories,omitempty"`
	Categories         *[]Category             `json:"categories,omitempty"`
	Channels           *[]CreateChannelParams  `json:"channels,omitempty"`
}

type CreateEventParams struct {
	ChartKey      string         `json:"chartKey"`
	ForSaleConfig *ForSaleConfig `json:"forSaleConfig,omitempty"`
	*EventParams
}

type CreateMultipleEventParams struct {
	ForSaleConfig *ForSaleConfig `json:"forSaleConfig,omitempty"`
	*EventParams
}

type UpdateEventParams struct {
	ChartKey    *string `json:"chartKey,omitempty"`
	IsInThePast *bool   `json:"isInThePast,omitempty"`
	*EventParams
}

type CreateMultipleEventsRequest struct {
	ChartKey string                      `json:"chartKey"`
	Events   []CreateMultipleEventParams `json:"events"`
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

type ObjectStatus string

const (
	FREE   ObjectStatus = "free"
	BOOKED ObjectStatus = "booked"
	HELD   ObjectStatus = "reservedByToken"
)

type StatusChanges struct {
	Status                   ObjectStatus       `json:"status"`
	Objects                  []ObjectProperties `json:"objects"`
	HoldToken                string             `json:"holdToken,omitempty"`
	OrderId                  string             `json:"orderId,omitempty"`
	KeepExtraData            bool               `json:"keepExtraData"`
	IgnoreChannels           bool               `json:"ignoreChannels"`
	ChannelKeys              []string           `json:"channelKeys,omitempty"`
	AllowedPreviousStatuses  []ObjectStatus     `json:"allowedPreviousStatuses,omitempty"`
	RejectedPreviousStatuses []ObjectStatus     `json:"rejectedPreviousStatuses,omitempty"`
}

type StatusChangeParams struct {
	Events []string `json:"events"`
	StatusChanges
}

type StatusChangeInBatchRequest struct {
	StatusChanges []StatusChangeInBatchParams `json:"statusChanges"`
}

type StatusChangeInBatchParams struct {
	Event string `json:"event"`
	StatusChanges
}

type BestAvailableStatusChangeParams struct {
	Status         ObjectStatus        `json:"status"`
	BestAvailable  BestAvailableParams `json:"bestAvailable"`
	HoldToken      string              `json:"holdToken,omitempty"`
	OrderId        string              `json:"orderId,omitempty"`
	KeepExtraData  bool                `json:"keepExtraData"`
	IgnoreChannels bool                `json:"ignoreChannels"`
	ChannelKeys    []string            `json:"channelKeys,omitempty"`
}

type BestAvailableParams struct {
	Number      int           `json:"number"`
	Categories  []CategoryKey `json:"categories,omitempty"`
	ExtraData   []ExtraData   `json:"extraData,omitempty"`
	TicketTypes []string      `json:"ticketTypes,omitempty"`
}

type ExtraData = map[string]any

type ForSaleConfigParams struct {
	Objects    []string       `json:"objects,omitempty"`
	AreaPlaces map[string]int `json:"areaPlaces,omitempty"`
	Categories []string       `json:"categories,omitempty"`
}

type overrideSeasonObjectStatusRequest struct {
	Objects []string `json:"objects"`
}

type updateExtraDataRequest struct {
	ExtraData map[string]ExtraData `json:"extraData"`
}

type ListParamsOption func(Params *shared.PageFetcher[StatusChange])

type eventSupportNS struct{}

var EventSupport eventSupportNS

func (events *Events) Create(params *CreateEventParams) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetBody(params).
		SetSuccessResult(&event).
		Post("/events")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) CreateMultiple(chartKey string, params ...CreateMultipleEventParams) (*CreateEventResult, error) {
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

func (events *Events) ChangeObjectStatus(eventKeys []string, objects []string, status ObjectStatus) (*ChangeObjectStatusResult, error) {
	objectProperties := make([]ObjectProperties, len(objects))
	for i, object := range objects {
		objectProperties[i] = ObjectProperties{ObjectId: object}
	}
	return events.ChangeObjectStatusWithOptions(&StatusChangeParams{
		Events: eventKeys,
		StatusChanges: StatusChanges{
			Status:  status,
			Objects: objectProperties,
		},
	})
}

func (events *Events) ChangeObjectStatusWithOptions(statusChangeparams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	var changeObjectStatusResult ChangeObjectStatusResult
	result, err := events.Client.R().
		SetBody(statusChangeparams).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusResult).
		Post("/events/groups/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusResult)
}

func (events *Events) ChangeObjectStatusInBatch(statusChangeInBatchParams ...StatusChangeInBatchParams) (*ChangeObjectStatusInBatchResult, error) {
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

func (events *Events) OverrideSeasonObjectStatus(eventKey string, objects ...string) error {
	result, err := events.Client.R().
		SetBody(&overrideSeasonObjectStatusRequest{
			Objects: objects,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/override-season-status")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) UseSeasonObjectStatus(eventKey string, objects ...string) error {
	result, err := events.Client.R().
		SetBody(&overrideSeasonObjectStatusRequest{
			Objects: objects,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/use-season-status")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) UpdateExtraData(eventKey string, extraData map[string]ExtraData) error {
	result, err := events.Client.R().
		SetBody(&updateExtraDataRequest{
			ExtraData: extraData,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/update-extra-data")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) RetrieveObjectInfo(eventKey string, objectLabels ...string) (map[string]EventObjectInfo, error) {
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

func (events *Events) StatusChanges(eventKey string, opts ...ListParamsOption) *shared.Lister[StatusChange] {
	pageFetcher := shared.PageFetcher[StatusChange]{
		Client:      events.Client,
		Url:         "/events/{eventKey}/status-changes",
		UrlParams:   map[string]string{"eventKey": eventKey},
		QueryParams: map[string]string{},
	}
	for _, opt := range opts {
		opt(&pageFetcher)
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

func (events *Events) Book(eventKey string, objectIds ...string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(BOOKED, eventKey, events.toObjectProperties(objectIds), nil)
}

func (events *Events) BookWithHoldToken(eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(BOOKED, eventKey, events.toObjectProperties(objectIds), holdToken)
}

func (events *Events) BookWithObjectProperties(eventKey string, objectProperties ...ObjectProperties) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(BOOKED, eventKey, objectProperties, nil)
}

func (events *Events) BookWithObjectPropertiesAndHoldToken(eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(BOOKED, eventKey, objectProperties, holdToken)
}

func (events *Events) BookWithOptions(statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Status = BOOKED
	return events.ChangeObjectStatusWithOptions(statusChangeParams)
}

func (events *Events) BookBestAvailable(eventKey string, params BestAvailableParams) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(eventKey, &BestAvailableStatusChangeParams{
		Status:        BOOKED,
		BestAvailable: params,
	})
}

func (events *Events) BookBestAvailableWithHoldToken(eventKey string, params BestAvailableParams, holdToken string) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(eventKey, &BestAvailableStatusChangeParams{
		Status:        BOOKED,
		BestAvailable: params,
		HoldToken:     holdToken,
	})
}

func (events *Events) BookBestAvailableWithOptions(eventKey string, params BestAvailableStatusChangeParams) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(eventKey, &params)
}

func (events *Events) Hold(eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(HELD, eventKey, events.toObjectProperties(objectIds), holdToken)
}

func (events *Events) HoldWithObjectProperties(eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(HELD, eventKey, objectProperties, holdToken)
}

func (events *Events) HoldWithOptions(statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Status = HELD
	return events.ChangeObjectStatusWithOptions(statusChangeParams)
}

func (events *Events) HoldBestAvailable(eventKey string, params BestAvailableParams, holdToken string) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(eventKey, &BestAvailableStatusChangeParams{
		Status:        HELD,
		BestAvailable: params,
		HoldToken:     holdToken,
	})
}

func (events *Events) Release(eventKey string, objectIds ...string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(FREE, eventKey, events.toObjectProperties(objectIds), nil)
}

func (events *Events) ReleaseWithHoldToken(eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(FREE, eventKey, events.toObjectProperties(objectIds), holdToken)
}

func (events *Events) ReleaseWithOptions(statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Status = FREE
	return events.ChangeObjectStatusWithOptions(statusChangeParams)
}

func (events *Events) changeStatus(status ObjectStatus, eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	params := StatusChangeParams{
		Events: []string{eventKey},
		StatusChanges: StatusChanges{
			Status:  status,
			Objects: objectProperties,
		},
	}
	if holdToken != nil {
		params.HoldToken = *holdToken
	}
	return events.ChangeObjectStatusWithOptions(&params)
}

func (events *Events) toObjectProperties(objects []string) []ObjectProperties {
	objectProperties := make([]ObjectProperties, len(objects))
	for i, object := range objects {
		objectProperties[i] = ObjectProperties{ObjectId: object}
	}
	return objectProperties
}

func (events *Events) RemoveCategories(eventKey string) error {
	return events.Update(eventKey, &UpdateEventParams{
		EventParams: &EventParams{
			Categories: &[]Category{},
		},
	})
}

func (events *Events) RemoveObjectCategories(eventKey string) error {
	return events.Update(eventKey, &UpdateEventParams{
		EventParams: &EventParams{
			ObjectCategories: &map[string]CategoryKey{},
		},
	})
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

func DateFormat(date *time.Time) string {
	return date.Format(time.DateOnly)
}

func (eventSupportNS) WithFilter(filterValue string) ListParamsOption {
	return func(pageFetcher *shared.PageFetcher[StatusChange]) {
		pageFetcher.QueryParams["filter"] = filterValue
	}
}

func (eventSupportNS) WithSortAsc(sortField string) ListParamsOption {
	return func(pageFetcher *shared.PageFetcher[StatusChange]) {
		pageFetcher.QueryParams["sort"] = shared.ToSort(sortField, "asc")
	}
}

func (eventSupportNS) WithSortDesc(sortField string) ListParamsOption {
	return func(pageFetcher *shared.PageFetcher[StatusChange]) {
		pageFetcher.QueryParams["sort"] = shared.ToSort(sortField, "desc")
	}
}
