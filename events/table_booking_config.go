package events

type TableBookingConfig struct {
	Mode   string            `json:"mode"`
	Tables map[string]string `json:"tables"`
}
