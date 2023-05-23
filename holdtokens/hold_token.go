package holdtokens

import "time"

type HoldToken struct {
	HoldToken        string     `json:"holdToken"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	ExpiresInSeconds int64      `json:"expiresInSeconds"`
}
