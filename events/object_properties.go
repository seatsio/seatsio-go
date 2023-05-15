package events

type ObjectProperties struct {
	ObjectId   string            `json:"objectId"`
	ExtraData  map[string]string `json:"extraData,omitempty"`
	TicketType string            `json:"ticketType,omitempty"`
	Quantity   int               `json:"quantity,omitempty"`
}
