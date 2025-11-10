package events

type ForSaleConfig struct {
	ForSale    bool           `json:"forSale"`
	Objects    []string       `json:"objects,omitempty"`
	AreaPlaces map[string]int `json:"areaPlaces,omitempty"`
	Categories []string       `json:"categories,omitempty"`
}
