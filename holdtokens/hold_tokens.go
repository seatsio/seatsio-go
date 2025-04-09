package holdtokens

import (
	"context"
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/v10/shared"
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

func (holdTokens *HoldTokens) Create(context context.Context) (*HoldToken, error) {
	var holdToken HoldToken
	result, err := holdTokens.Client.R().
		SetContext(context).
		SetSuccessResult(&holdToken).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) CreateWithExpiration(context context.Context, expiresInMinutes int) (*HoldToken, error) {
	var holdToken HoldToken
	request := &CreateHoldTokenRequest{ExpiresInMinutes: expiresInMinutes}
	result, err := holdTokens.Client.R().
		SetContext(context).
		SetSuccessResult(&holdToken).
		SetBody(&request).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) Retrieve(context context.Context, token string) (*HoldToken, error) {
	var holdToken HoldToken
	result, err := holdTokens.Client.R().
		SetContext(context).
		SetSuccessResult(&holdToken).
		SetPathParam("token", token).
		Get("/hold-tokens/{token}")
	return shared.AssertOk(result, err, &holdToken)
}

func (holdTokens *HoldTokens) ExpireInMinutes(context context.Context, token string, expiresInMinutes int) (*HoldToken, error) {
	var holdToken HoldToken
	request := &SetExpirationDateOfHoldTokenRequest{ExpiresInMinutes: expiresInMinutes}
	result, err := holdTokens.Client.R().
		SetContext(context).
		SetSuccessResult(&holdToken).
		SetBody(&request).
		SetPathParam("token", token).
		Post("/hold-tokens/{token}")
	return shared.AssertOk(result, err, &holdToken)
}
