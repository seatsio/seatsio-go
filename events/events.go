package events

import "github.com/imroc/req/v3"

type Events struct {
	secretKey string
	baseUrl   string
}

type createEventRequest struct {
	ChartKey string `json:"chartKey"`
}

func (events *Events) Create(chartKey string) *Event {
	client := req.C()
	var event Event
	// TODO: error handling
	client.R().
		SetBasicAuth(events.secretKey, "").
		SetBody(&createEventRequest{
			chartKey,
		}).
		SetSuccessResult(&event).
		Post(events.baseUrl + "/events")
	return &event
}

func NewEvents(secretKey string, baseUrl string) *Events {
	return &Events{secretKey, baseUrl}
}
