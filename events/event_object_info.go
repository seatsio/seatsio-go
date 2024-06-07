package events

type EventObjectInfo struct {
	Status                         string                    `json:"status,omitempty"`
	Label                          string                    `json:"label,omitempty"`
	Labels                         Labels                    `json:"labels,omitempty"`
	IDs                            IDs                       `json:"ids,omitempty"`
	CategoryLabel                  string                    `json:"categoryLabel,omitempty"`
	CategoryKey                    CategoryKey               `json:"categoryKey,omitempty"`
	TicketType                     string                    `json:"ticketType,omitempty"`
	HoldToken                      string                    `json:"holdToken,omitempty"`
	ObjectType                     string                    `json:"objectType,omitempty"`
	BookAsAWhole                   bool                      `json:"bookAsAWhole"`
	OrderId                        string                    `json:"orderId,omitempty"`
	ForSale                        bool                      `json:"forSale"`
	Section                        string                    `json:"section,omitempty"`
	Entrance                       string                    `json:"entrance,omitempty"`
	Capacity                       int                       `json:"capacity"`
	NumBooked                      int                       `json:"numBooked"`
	NumFree                        int                       `json:"numFree"`
	NumHeld                        int                       `json:"numHeld"`
	NumSeats                       int                       `json:"numSeats"`
	ExtraData                      ExtraData                 `json:"extraData,omitempty"`
	IsAccessible                   bool                      `json:"isAccessible"`
	IsCompanionSeat                bool                      `json:"isCompanionSeat"`
	HasRestrictedView              bool                      `json:"hasRestrictedView"`
	DisplayedObjectType            string                    `json:"displayedObjectType,omitempty"`
	LeftNeighbour                  string                    `json:"leftNeighbour,omitempty"`
	RightNeighbour                 string                    `json:"rightNeighbour,omitempty"`
	IsAvailable                    bool                      `json:"isAvailable"`
	AvailabilityReason             string                    `json:"availabilityReason,omitempty"`
	Channel                        string                    `json:"channel,omitempty"`
	DistanceToFocalPoint           float64                   `json:"distanceToFocalPoint"`
	Holds                          map[string]map[string]int `json:"holds,omitempty"`
	VariableOccupancy              bool                      `json:"variableOccupancy"`
	MinOccupancy                   int                       `json:"minOccupancy"`
	MaxOccupancy                   int                       `json:"maxOccupancy"`
	SeasonStatusOverriddenQuantity int                       `json:"seasonStatusOverriddenQuantity"`
	NumNotForSale                  int                       `json:"numNotForSale"`
}

type IDs struct {
	Own     string `json:"own"`
	Parent  string `json:"parent"`
	Section string `json:"section"`
}

type Labels struct {
	Own     LabelAndType `json:"own"`
	Parent  LabelAndType `json:"parent"`
	Section string       `json:"section"`
}

type LabelAndType struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}
