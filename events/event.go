package events

import "time"

type Event struct {
	Id                    int64                  `json:"id,omitempty"`
	Key                   string                 `json:"key,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	Date                  string                 `json:"date,omitempty"`
	ChartKey              string                 `json:"chartKey"`
	HoldToken             string                 `json:"holdToken,omitempty"`
	TableBookingConfig    TableBookingConfig     `json:"tableBookingConfig,omitempty"`
	SupportsBestAvailable bool                   `json:"supportsBestAvailable,omitempty"`
	ForSaleConfig         *ForSaleConfig         `json:"forSaleConfig,omitempty"`
	CreatedOn             *time.Time             `json:"createdOn"`
	UpdatedOn             *time.Time             `json:"updatedOn,omitempty"`
	Categories            []Category             `json:"categories,omitempty"`
	ObjectCategories      map[string]CategoryKey `json:"objectCategories,omitempty"`
	Channels              []Channel              `json:"channels,omitempty"`
	IsInThePast           bool                   `json:"isInThePast"`
}
