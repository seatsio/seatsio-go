package events

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v11/shared"
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
	IsInThePast *bool `json:"isInThePast,omitempty"`
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

const (
	FREE   = "free"
	BOOKED = "booked"
	HELD   = "reservedByToken"
	RESALE = "resale"
)

type StatusChangeType string

const (
	CHANGE_STATUS_TO       StatusChangeType = "CHANGE_STATUS_TO"
	RELEASE                StatusChangeType = "RELEASE"
	OVERRIDE_SEASON_STATUS StatusChangeType = "OVERRIDE_SEASON_STATUS"
	USE_SEASON_STATUS      StatusChangeType = "USE_SEASON_STATUS"
)

type StatusChanges struct {
	Type                     StatusChangeType   `json:"type,omitempty"`
	Status                   string             `json:"status,omitempty"`
	Objects                  []ObjectProperties `json:"objects"`
	HoldToken                string             `json:"holdToken,omitempty"`
	OrderId                  string             `json:"orderId,omitempty"`
	KeepExtraData            bool               `json:"keepExtraData,omitempty"`
	IgnoreChannels           bool               `json:"ignoreChannels,omitempty"`
	ChannelKeys              []string           `json:"channelKeys,omitempty"`
	AllowedPreviousStatuses  []string           `json:"allowedPreviousStatuses,omitempty"`
	RejectedPreviousStatuses []string           `json:"rejectedPreviousStatuses,omitempty"`
	ResaleListingId          string             `json:"resaleListingId,omitempty"`
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
	Status         string              `json:"status"`
	BestAvailable  BestAvailableParams `json:"bestAvailable"`
	HoldToken      string              `json:"holdToken,omitempty"`
	OrderId        string              `json:"orderId,omitempty"`
	KeepExtraData  bool                `json:"keepExtraData"`
	IgnoreChannels bool                `json:"ignoreChannels"`
	ChannelKeys    []string            `json:"channelKeys,omitempty"`
}

type BestAvailableParams struct {
	Number          int           `json:"number"`
	Zone            string        `json:"zone,omitempty"`
	Sections        []string      `json:"sections,omitempty"`
	Categories      []CategoryKey `json:"categories,omitempty"`
	ExtraData       []ExtraData   `json:"extraData,omitempty"`
	TicketTypes     []string      `json:"ticketTypes,omitempty"`
	AccessibleSeats int           `json:"accessibleSeats,omitempty"`
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

func (events *Events) Create(context context.Context, params *CreateEventParams) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetContext(context).
		SetBody(params).
		SetSuccessResult(&event).
		Post("/events")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) CreateMultiple(context context.Context, chartKey string, params ...CreateMultipleEventParams) (*CreateEventResult, error) {
	var eventCreationResult CreateEventResult
	result, err := events.Client.R().
		SetContext(context).
		SetBody(&CreateMultipleEventsRequest{
			ChartKey: chartKey,
			Events:   params,
		}).
		SetSuccessResult(&eventCreationResult).
		Post("/events/actions/create-multiple")
	return shared.AssertOk(result, err, &eventCreationResult)
}

func (events *Events) Update(context context.Context, eventKey string, params *UpdateEventParams) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(params).
		SetPathParam("event", eventKey).
		Post("/events/{event}")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) ChangeObjectStatus(context context.Context, eventKeys []string, objects []string, status string) (*ChangeObjectStatusResult, error) {
	objectProperties := make([]ObjectProperties, len(objects))
	for i, object := range objects {
		objectProperties[i] = ObjectProperties{ObjectId: object}
	}
	return events.ChangeObjectStatusWithOptions(context, &StatusChangeParams{
		Events: eventKeys,
		StatusChanges: StatusChanges{
			Status:  status,
			Objects: objectProperties,
		},
	})
}

func (events *Events) ChangeObjectStatusWithOptions(context context.Context, statusChangeparams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	var changeObjectStatusResult ChangeObjectStatusResult
	result, err := events.Client.R().
		SetContext(context).
		SetBody(statusChangeparams).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusResult).
		Post("/events/groups/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusResult)
}

func (events *Events) ChangeObjectStatusInBatch(context context.Context, statusChangeInBatchParams ...StatusChangeInBatchParams) (*ChangeObjectStatusInBatchResult, error) {
	var changeObjectStatusInBatchResult ChangeObjectStatusInBatchResult
	result, err := events.Client.R().
		SetContext(context).
		SetBody(&StatusChangeInBatchRequest{
			StatusChanges: statusChangeInBatchParams,
		}).
		SetQueryParam("expand", "objects").
		SetSuccessResult(&changeObjectStatusInBatchResult).
		Post("/events/actions/change-object-status")
	return shared.AssertOk(result, err, &changeObjectStatusInBatchResult)
}

func (events *Events) ChangeBestAvailableObjectStatus(context context.Context, eventKey string, bestAvailableStatusChangeParams *BestAvailableStatusChangeParams) (*BestAvailableResult, error) {
	var bestAvailableResult BestAvailableResult
	result, err := events.Client.R().
		SetContext(context).
		SetBody(bestAvailableStatusChangeParams).
		SetSuccessResult(&bestAvailableResult).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/change-object-status")
	return shared.AssertOk(result, err, &bestAvailableResult)
}

func (events *Events) OverrideSeasonObjectStatus(context context.Context, eventKey string, objects ...string) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(&overrideSeasonObjectStatusRequest{
			Objects: objects,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/override-season-status")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) UseSeasonObjectStatus(context context.Context, eventKey string, objects ...string) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(&overrideSeasonObjectStatusRequest{
			Objects: objects,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/use-season-status")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) UpdateExtraData(context context.Context, eventKey string, extraData map[string]ExtraData) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(&updateExtraDataRequest{
			ExtraData: extraData,
		}).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/update-extra-data")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) RetrieveObjectInfo(context context.Context, eventKey string, objectLabels ...string) (map[string]EventObjectInfo, error) {
	var eventObjectInfos map[string]EventObjectInfo
	request := events.Client.R().
		SetContext(context).
		SetSuccessResult(&eventObjectInfos).
		AddQueryParams("label", objectLabels...)
	result, err := request.Get("/events/" + eventKey + "/objects")
	return shared.AssertOkMap(result, err, eventObjectInfos)
}

func (events *Events) Delete(context context.Context, eventKey string) error {
	result, err := events.Client.R().
		SetContext(context).
		SetQueryParam("expand", "objects").
		SetPathParam("event", eventKey).
		Delete("/events/{event}")
	return shared.AssertOkNoBody(result, err)
}

func (events *Events) Retrieve(context context.Context, eventKey string) (*Event, error) {
	var event Event
	result, err := events.Client.R().
		SetContext(context).
		SetSuccessResult(&event).
		SetPathParam("event", eventKey).
		Get("/events/{event}")
	return shared.AssertOk(result, err, &event)
}

func (events *Events) MarkAsNotForSale(context context.Context, eventKey string, forSaleConfig *ForSaleConfigParams) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(forSaleConfig).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-as-not-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) MarkAsForSale(context context.Context, eventKey string, forSaleConfig *ForSaleConfigParams) error {
	result, err := events.Client.R().
		SetContext(context).
		SetBody(forSaleConfig).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-as-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) MarkEverythingAsForSale(context context.Context, eventKey string) error {
	result, err := events.Client.R().
		SetContext(context).
		SetPathParam("event", eventKey).
		Post("/events/{event}/actions/mark-everything-as-for-sale")
	return shared.AssertOkWithoutResult(result, err)
}

func (events *Events) StatusChanges(context context.Context, eventKey string, opts ...ListParamsOption) *shared.Lister[StatusChange] {
	pageFetcher := shared.PageFetcher[StatusChange]{
		Client:      events.Client,
		Url:         "/events/{eventKey}/status-changes",
		UrlParams:   map[string]string{"eventKey": eventKey},
		QueryParams: map[string]string{},
		Context:     &context,
	}
	for _, opt := range opts {
		opt(&pageFetcher)
	}
	return &shared.Lister[StatusChange]{PageFetcher: &pageFetcher}
}

func (events *Events) StatusChangesForObject(context context.Context, eventKey string, objectLabel string) *shared.Lister[StatusChange] {
	pageFetcher := shared.PageFetcher[StatusChange]{
		Client:    events.Client,
		Url:       "/events/{eventKey}/objects/{objectLabel}/status-changes",
		UrlParams: map[string]string{"eventKey": eventKey, "objectLabel": objectLabel},
		Context:   &context,
	}
	return &shared.Lister[StatusChange]{PageFetcher: &pageFetcher}
}

func (events *Events) ListAll(context context.Context, opts ...shared.PaginationParamsOption) ([]Event, error) {
	return events.lister(context).All(opts...)
}

func (events *Events) Book(context context.Context, eventKey string, objectIds ...string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, BOOKED, eventKey, events.toObjectProperties(objectIds), nil, nil)
}

func (events *Events) BookWithHoldToken(context context.Context, eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, BOOKED, eventKey, events.toObjectProperties(objectIds), holdToken, nil)
}

func (events *Events) BookWithObjectProperties(context context.Context, eventKey string, objectProperties ...ObjectProperties) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, BOOKED, eventKey, objectProperties, nil, nil)
}

func (events *Events) BookWithObjectPropertiesAndHoldToken(context context.Context, eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, BOOKED, eventKey, objectProperties, holdToken, nil)
}

func (events *Events) BookWithOptions(context context.Context, statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Status = BOOKED
	return events.ChangeObjectStatusWithOptions(context, statusChangeParams)
}

func (events *Events) BookBestAvailable(context context.Context, eventKey string, params BestAvailableParams) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(context, eventKey, &BestAvailableStatusChangeParams{
		Status:        BOOKED,
		BestAvailable: params,
	})
}

func (events *Events) BookBestAvailableWithHoldToken(context context.Context, eventKey string, params BestAvailableParams, holdToken string) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(context, eventKey, &BestAvailableStatusChangeParams{
		Status:        BOOKED,
		BestAvailable: params,
		HoldToken:     holdToken,
	})
}

func (events *Events) BookBestAvailableWithOptions(context context.Context, eventKey string, params BestAvailableStatusChangeParams) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(context, eventKey, &params)
}

func (events *Events) Hold(context context.Context, eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, HELD, eventKey, events.toObjectProperties(objectIds), holdToken, nil)
}

func (events *Events) HoldWithObjectProperties(context context.Context, eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, HELD, eventKey, objectProperties, holdToken, nil)
}

func (events *Events) HoldWithOptions(context context.Context, statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Status = HELD
	return events.ChangeObjectStatusWithOptions(context, statusChangeParams)
}

func (events *Events) HoldBestAvailable(context context.Context, eventKey string, params BestAvailableParams, holdToken string) (*BestAvailableResult, error) {
	return events.ChangeBestAvailableObjectStatus(context, eventKey, &BestAvailableStatusChangeParams{
		Status:        HELD,
		BestAvailable: params,
		HoldToken:     holdToken,
	})
}

func (events *Events) PutUpForResale(context context.Context, eventKey string, objectIds []string, resaleListingId *string) (*ChangeObjectStatusResult, error) {
	return events.changeStatus(context, RESALE, eventKey, events.toObjectProperties(objectIds), nil, resaleListingId)
}

func (events *Events) Release(context context.Context, eventKey string, objectIds ...string) (*ChangeObjectStatusResult, error) {
	return events.releaseObjects(context, eventKey, events.toObjectProperties(objectIds), nil)
}

func (events *Events) ReleaseWithHoldToken(context context.Context, eventKey string, objectIds []string, holdToken *string) (*ChangeObjectStatusResult, error) {
	return events.releaseObjects(context, eventKey, events.toObjectProperties(objectIds), holdToken)
}

func (events *Events) ReleaseWithOptions(context context.Context, statusChangeParams *StatusChangeParams) (*ChangeObjectStatusResult, error) {
	statusChangeParams.Type = "RELEASE"
	return events.ChangeObjectStatusWithOptions(context, statusChangeParams)
}

func (events *Events) releaseObjects(context context.Context, eventKey string, objectProperties []ObjectProperties, holdToken *string) (*ChangeObjectStatusResult, error) {
	params := StatusChangeParams{
		Events: []string{eventKey},
		StatusChanges: StatusChanges{
			Type:    "RELEASE",
			Objects: objectProperties,
		},
	}
	if holdToken != nil {
		params.HoldToken = *holdToken
	}
	return events.ChangeObjectStatusWithOptions(context, &params)
}

func (events *Events) changeStatus(context context.Context, status string, eventKey string, objectProperties []ObjectProperties, holdToken *string, resaleListingId *string) (*ChangeObjectStatusResult, error) {
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
	if resaleListingId != nil {
		params.ResaleListingId = *resaleListingId
	}
	return events.ChangeObjectStatusWithOptions(context, &params)
}

func (events *Events) toObjectProperties(objects []string) []ObjectProperties {
	objectProperties := make([]ObjectProperties, len(objects))
	for i, object := range objects {
		objectProperties[i] = ObjectProperties{ObjectId: object}
	}
	return objectProperties
}

func (events *Events) RemoveCategories(context context.Context, eventKey string) error {
	return events.Update(context, eventKey, &UpdateEventParams{
		EventParams: &EventParams{
			Categories: &[]Category{},
		},
	})
}

func (events *Events) RemoveObjectCategories(context context.Context, eventKey string) error {
	return events.Update(context, eventKey, &UpdateEventParams{
		EventParams: &EventParams{
			ObjectCategories: &map[string]CategoryKey{},
		},
	})
}

func (events *Events) lister(context context.Context) *shared.Lister[Event] {
	pageFetcher := shared.PageFetcher[Event]{
		Client:    events.Client,
		Url:       "/events",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[Event]{PageFetcher: &pageFetcher}
}

func (events *Events) ListFirstPage(context context.Context, opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister(context).ListFirstPage(opts...)
}

func (events *Events) ListPageAfter(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister(context).ListPageAfter(id, opts...)
}

func (events *Events) ListPageBefore(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[Event], error) {
	return events.lister(context).ListPageBefore(id, opts...)
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
