package events

type ForSaleConfig struct {
	ForSale    bool           `json:"forSale"`
	Objects    []string       `json:"objects"`
	AreaPlaces map[string]int `json:"areaPlaces"`
	Categories []string       `json:"categories"`
}
