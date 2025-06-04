package ticketbuyers

import (
	"context"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v11/shared"
)

type TicketBuyers struct {
	Client *req.Client
}

type TicketBuyerParams struct {
	Ids []uuid.UUID `json:"ids"`
}

type AddTicketBuyerIdsResponse struct {
	Added          []uuid.UUID `json:"added"`
	AlreadyPresent []uuid.UUID `json:"alreadyPresent"`
}

type RemoveTicketBuyerIdsResponse struct {
	Removed    []uuid.UUID `json:"removed"`
	NotPresent []uuid.UUID `json:"notPresent"`
}

func (ticketBuyers *TicketBuyers) Add(context context.Context, params *TicketBuyerParams) (*AddTicketBuyerIdsResponse, error) {
	var response AddTicketBuyerIdsResponse
	result, err := ticketBuyers.Client.R().
		SetContext(context).
		SetBody(params).
		SetSuccessResult(&response).
		Post("/ticket-buyers")
	return shared.AssertOk(result, err, &response)
}

func (ticketBuyers *TicketBuyers) Remove(context context.Context, params *TicketBuyerParams) (*RemoveTicketBuyerIdsResponse, error) {
	var response RemoveTicketBuyerIdsResponse
	result, err := ticketBuyers.Client.R().
		SetContext(context).
		SetBody(params).
		SetSuccessResult(&response).
		Delete("/ticket-buyers")
	return shared.AssertOk(result, err, &response)
}

func (ticketBuyers *TicketBuyers) lister(context context.Context) *shared.Lister[uuid.UUID] {
	pageFetcher := shared.PageFetcher[uuid.UUID]{
		Client:    ticketBuyers.Client,
		Url:       "/ticket-buyers",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[uuid.UUID]{PageFetcher: &pageFetcher}
}

func (ticketBuyers *TicketBuyers) ListAll(context context.Context) ([]uuid.UUID, error) {
	return ticketBuyers.lister(context).All()
}

func (ticketBuyers *TicketBuyers) ListFirstPage(context context.Context, pageSize ...int) (*shared.Page[uuid.UUID], error) {
	if pageSize != nil {
		return ticketBuyers.lister(context).ListFirstPage(shared.Pagination.PageSize(pageSize[0]))
	}
	return ticketBuyers.lister(context).ListFirstPage()
}

func (ticketBuyers *TicketBuyers) ListPageAfter(context context.Context, id int64, pageSize ...int) (*shared.Page[uuid.UUID], error) {
	if pageSize != nil {
		return ticketBuyers.lister(context).ListPageAfter(id, shared.Pagination.PageSize(pageSize[0]))
	}
	return ticketBuyers.lister(context).ListPageAfter(id)
}

func (ticketBuyers *TicketBuyers) ListPageBefore(context context.Context, id int64, pageSize ...int) (*shared.Page[uuid.UUID], error) {
	if pageSize != nil {
		return ticketBuyers.lister(context).ListPageBefore(id, shared.Pagination.PageSize(pageSize[0]))
	}
	return ticketBuyers.lister(context).ListPageBefore(id)
}
