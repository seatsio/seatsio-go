package seatsio

import (
	"github.com/seatsio/seatsio-go/events"
	"github.com/seatsio/seatsio-go/holdtokens"
	"github.com/seatsio/seatsio-go/shared"
)

type SeatsioClient struct {
	baseUrl    string
	secretKey  string
	Events     *events.Events
	HoldTokens *holdtokens.HoldTokens
}

func NewSeatsioClient(secretKey string, baseUrl string) *SeatsioClient {
	apiClient := shared.ApiClient(secretKey, baseUrl)
	return &SeatsioClient{
		baseUrl:    baseUrl,
		secretKey:  secretKey,
		Events:     &events.Events{Client: apiClient},
		HoldTokens: &holdtokens.HoldTokens{Client: apiClient},
	}
}
