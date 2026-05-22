package events

import "fmt"

type Channel struct {
	ID         string         `json:"id"`
	Key        string         `json:"key"`
	Name       string         `json:"name"`
	Color      string         `json:"color"`
	Index      int            `json:"index"`
	Objects    []string       `json:"objects,omitempty"`
	AreaPlaces map[string]int `json:"areaPlaces,omitempty"`
}

func (c Channel) AreaPartitionLabel(areaLabel string) string {
	return fmt.Sprintf("%s##%s", areaLabel, c.ID)
}
