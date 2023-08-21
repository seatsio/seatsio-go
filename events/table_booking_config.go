package events

type TableBookingModeSupportNS struct{}

var TableBookingSupport TableBookingModeSupportNS

type TableBookingMode string

const (
	BY_TABLE TableBookingMode = "BY_TABLE"
	BY_SEAT  TableBookingMode = "BY_SEAT"
)

type Mode string

const (
	INHERIT      Mode = "INHERIT"
	ALL_BY_SEAT  Mode = "ALL_BY_SEAT"
	ALL_BY_TABLE Mode = "ALL_BY_TABLE"
	CUSTOM       Mode = "CUSTOM"
)

type TableBookingConfig struct {
	Mode   Mode                        `json:"mode"`
	Tables map[string]TableBookingMode `json:"tables,omitempty"`
}

func (TableBookingModeSupportNS) Inherit() TableBookingConfig {
	return TableBookingConfig{Mode: INHERIT, Tables: nil}
}

func (TableBookingModeSupportNS) AllBySeat() *TableBookingConfig {
	return &TableBookingConfig{Mode: ALL_BY_SEAT, Tables: nil}
}

func (TableBookingModeSupportNS) AllByTables() *TableBookingConfig {
	return &TableBookingConfig{Mode: ALL_BY_TABLE, Tables: nil}
}

func (TableBookingModeSupportNS) Custom() *TableBookingConfig {
	return &TableBookingConfig{Mode: CUSTOM, Tables: nil}
}
