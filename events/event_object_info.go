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
	BookAsAWhole                   bool                      `json:"bookAsAWhole,omitempty"`
	OrderId                        string                    `json:"orderId,omitempty"`
	ForSale                        bool                      `json:"forSale,omitempty"`
	Section                        string                    `json:"section,omitempty"`
	Entrance                       string                    `json:"entrance,omitempty"`
	Capacity                       int                       `json:"capacity,omitempty"`
	NumBooked                      int                       `json:"numBooked,omitempty"`
	NumFree                        int                       `json:"numFree,omitempty"`
	NumHeld                        int                       `json:"numHeld,omitempty"`
	NumSeats                       int                       `json:"numSeats,omitempty"`
	ExtraData                      ExtraData                 `json:"extraData,omitempty"`
	IsAccessible                   bool                      `json:"isAccessible,omitempty"`
	IsCompanionSeat                bool                      `json:"isCompanionSeat,omitempty"`
	HasLiftUpArmrests              bool                      `json:"hasLiftUpArmrests,omitempty"`
	IsHearingImpaired              bool                      `json:"isHearingImpaired,omitempty"`
	IsSemiAmbulatorySeat           bool                      `json:"isSemiAmbulatorySeat,omitempty"`
	HasSignLanguageInterpretation  bool                      `json:"hasSignLanguageInterpretation,omitempty"`
	IsPlusSize                     bool                      `json:"isPlusSize,omitempty"`
	HasRestrictedView              bool                      `json:"hasRestrictedView,omitempty"`
	DisplayedObjectType            string                    `json:"displayedObjectType,omitempty"`
	ParentDisplayedObjectType      string                    `json:"parentDisplayedObjectType,omitempty"`
	LeftNeighbour                  string                    `json:"leftNeighbour,omitempty"`
	RightNeighbour                 string                    `json:"rightNeighbour,omitempty"`
	IsAvailable                    bool                      `json:"isAvailable,omitempty"`
	AvailabilityReason             string                    `json:"availabilityReason,omitempty"`
	Channel                        string                    `json:"channel,omitempty"`
	DistanceToFocalPoint           float64                   `json:"distanceToFocalPoint,omitempty"`
	Holds                          map[string]map[string]int `json:"holds,omitempty"`
	VariableOccupancy              bool                      `json:"variableOccupancy,omitempty"`
	MinOccupancy                   int                       `json:"minOccupancy,omitempty"`
	MaxOccupancy                   int                       `json:"maxOccupancy,omitempty"`
	SeasonStatusOverriddenQuantity int                       `json:"seasonStatusOverriddenQuantity,omitempty"`
	NumNotForSale                  int                       `json:"numNotForSale,omitempty"`
	Zone                           string                    `json:"zone,omitempty"`
	Floor                          Floor                     `json:"floor,omitempty"`
	ResaleListingId                string                    `json:"resaleListingId,omitempty"`
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

type Floor struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
