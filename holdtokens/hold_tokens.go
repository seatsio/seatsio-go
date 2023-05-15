package holdtokens

import (
	"github.com/seatsio/seatsio-go/shared"
)

type HoldTokens struct {
	secretKey string
	baseUrl   string
}

func (holdTokens *HoldTokens) Create() (*HoldToken, error) {
	var holdToken HoldToken
	client := shared.ApiClient(holdTokens.secretKey, holdTokens.baseUrl)
	result, err := client.R().
		SetSuccessResult(&holdToken).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}

func NewHoldTokens(secretKey string, baseUrl string) *HoldTokens {
	return &HoldTokens{secretKey, baseUrl}
}
