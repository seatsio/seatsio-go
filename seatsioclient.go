package seatsio

import "github.com/seatsio/seatsio-go/events"

type SeatsioClient struct {
	baseUrl   string
	secretKey string
	Events    *events.Events
}

func NewSeatsioClient(secretKey string, baseUrl string) *SeatsioClient {
	return &SeatsioClient{
		baseUrl:   baseUrl,
		secretKey: secretKey,
		Events:    events.NewEvents(secretKey, baseUrl),
	}
}
