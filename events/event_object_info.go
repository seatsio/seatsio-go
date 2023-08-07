package events

type EventObjectInfo struct {
	Status         ObjectStatus `json:"status"`
	OrderId        string       `json:"orderId"`
	ExtraData      ExtraData    `json:"extraData"`
	Label          string       `json:"label"`
	Labels         Labels       `json:"labels"`
	IDs            IDs          `json:"ids"`
	CategoryLabel  string       `json:"categoryLabel"`
	CategoryKey    CategoryKey  `json:"categoryKey"`
	TicketType     string       `json:"ticketType"`
	ForSale        bool         `json:"forSale"`
	Section        string       `json:"section"`
	Entrance       string       `json:"entrance"`
	NumBooked      int          `json:"numBooked"`
	Capacity       int          `json:"capacity"`
	ObjectType     string       `json:"objectType"`
	LeftNeighbour  string       `json:"leftNeighbour"`
	RightNeighbour string       `json:"rightNeighbour"`
	HoldToken      *string      `json:"holdToken"`
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
