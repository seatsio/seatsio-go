package eventlog

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v6/shared"
)

type EventLog struct {
	Client *req.Client
}

func (eventLog EventLog) ListAll(opts ...shared.PaginationParamsOption) ([]EventLogItem, error) {
	return eventLog.lister().All(opts...)
}

func (eventLog EventLog) ListFirstPage(opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister().ListFirstPage(opts...)
}

func (eventLog EventLog) ListPageAfter(id int64, opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister().ListPageAfter(id, opts...)
}

func (eventLog EventLog) ListPageBefore(id int64, opts ...shared.PaginationParamsOption) (*shared.Page[EventLogItem], error) {
	return eventLog.lister().ListPageBefore(id, opts...)
}

func (eventLog EventLog) lister() *shared.Lister[EventLogItem] {
	pageFetcher := shared.PageFetcher[EventLogItem]{
		Client:    eventLog.Client,
		Url:       "/event-log",
		UrlParams: map[string]string{},
	}
	return &shared.Lister[EventLogItem]{PageFetcher: &pageFetcher}
}
