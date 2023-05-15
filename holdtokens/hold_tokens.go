package holdtokens

import (
	"github.com/imroc/req/v3"
	"github.com/seatsio/seatsio-go/shared"
)

type HoldTokens struct {
	Client *req.Client
}

func (holdTokens *HoldTokens) Create() (*HoldToken, error) {
	var holdToken HoldToken
	result, err := holdTokens.Client.R().
		SetSuccessResult(&holdToken).
		Post("/hold-tokens")
	return shared.AssertOk(result, err, &holdToken)
}
