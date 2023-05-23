package holdtokens

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type HoldTokens struct {
	Client *req.Client
}

type CreateHoldTokenRequest struct {
	ExpiresInMinutes int `json:"expiresInMinutes"`
}

type SetExpirationDateOfHoldTokenRequest struct {
	ExpiresInMinutes int `json:"expiresInMinutes"`
}

func (holdTokens *HoldTokens) Create() (*HoldToken, error) {
	var holdToken HoldToken
	result, err := holdTokens.Client.R().
		SetSuccessResult(&holdToken).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) CreateWithExpiration(expiresInMinutes int) (*HoldToken, error) {
	var holdToken HoldToken
	request := &CreateHoldTokenRequest{ExpiresInMinutes: expiresInMinutes}
	result, err := holdTokens.Client.R().
		SetSuccessResult(&holdToken).
		SetBody(&request).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) Retrieve(token string) (*HoldToken, error) {
	var holdToken HoldToken
	result, err := holdTokens.Client.R().
		SetSuccessResult(&holdToken).
		SetPathParam("token", token).
		Get("/hold-tokens/{token}")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) ExpireInMinutes(token string, expiresInMinutes int) (*HoldToken, error) {
	var holdToken HoldToken
	request := &SetExpirationDateOfHoldTokenRequest{ExpiresInMinutes: expiresInMinutes}
	result, err := holdTokens.Client.R().
		SetSuccessResult(&holdToken).
		SetBody(&request).
		SetPathParam("token", token).
		Post("/hold-tokens/{token}")
	return shared.AssertOk(result, err, &holdToken)
}
