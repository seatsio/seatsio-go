package eventlog

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v12/shared"
)

type EventLog struct {
	Client *req.Client
}

func (eventLog EventLog) ListAll(context context.Context, opts ...shared.PaginationParamsOption) ([]EventLogItem, error) {
	return eventLog.lister(context).All(opts...)
}

func (eventLog EventLog) ListFirstPage(context context.Context, opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister(context).ListFirstPage(opts...)
}

func (eventLog EventLog) ListPageAfter(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister(context).ListPageAfter(id, opts...)
}

func (eventLog EventLog) ListPageBefore(context context.Context, id int64, opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister(context).ListPageBefore(id, opts...)
}

func (eventLog EventLog) lister(context context.Context) *shared.Lister[EventLogItem] {
	pageFetcher := shared.PageFetcher[EventLogItem]{
		Client:    eventLog.Client,
		Url:       "/event-log",
		UrlParams: map[string]string{},
		Context:   &context,
	}
	return &shared.Lister[EventLogItem]{PageFetcher: &pageFetcher}
}
