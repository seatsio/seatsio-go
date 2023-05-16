package events

import "time"

type Event struct {
	Id                    int64                  `json:"id"`
	Key                   string                 `json:"key"`
	ChartKey              string                 `json:"chartKey"`
	HoldToken             string                 `json:"holdToken"`
	TableBookingConfig    TableBookingConfig     `json:"tableBookingConfig"`
	SupportsBestAvailable bool                   `json:"supportsBestAvailable"`
	ForSaleConfig         *ForSaleConfig         `json:"forSaleConfig"`
	CreatedOn             *time.Time             `json:"createdOn"`
	UpdatedOn             *time.Time             `json:"updatedOn"`
	Categories            []Category             `json:"categories"`
	ObjectCategories      map[string]CategoryKey `json:"objectCategories"`
}
