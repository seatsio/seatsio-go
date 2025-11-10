package events

type ObjectAndQuantity struct {
	Object   string `json:"object"`
	Quantity int    `json:"quantity,omitempty"`
}
