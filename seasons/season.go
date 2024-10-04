package seasons

import "github.com/seatsio/seatsio-go/v8/events"

type Season struct {
	events.Event
	Events            []events.Event `json:"events"`
	PartialSeasonKeys []string       `json:"partialSeasonKeys"`
	IsSeason          bool           `json:"isSeason"`
	IsTopLevelSeason  bool           `json:"isTopLevelSeason"`
	IsPartialSeason   bool           `json:"isPartialSeason"`
	IsEventInSeason   bool           `json:"isEventInSeason"`
	TopLevelSeasonKey *string        `json:"topLevelSeasonKey"`
}
