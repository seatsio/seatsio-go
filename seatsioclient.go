package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/holdtokens"
)

type SeatsioClient struct {
	baseUrl    string
	secretKey  string
	Events     *events.Events
	HoldTokens *holdtokens.HoldTokens
}

func NewSeatsioClient(secretKey string, baseUrl string) *SeatsioClient {
	return &SeatsioClient{
		baseUrl:    baseUrl,
		secretKey:  secretKey,
		Events:     events.NewEvents(secretKey, baseUrl),
		HoldTokens: holdtokens.NewHoldTokens(secretKey, baseUrl),
	}
}
