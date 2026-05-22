package events

import "fmt"

type Channel struct {
	Id         string         `json:"id"`
	Key        string         `json:"key"`
	Name       string         `json:"name"`
	Color      string         `json:"color"`
	Index      int            `json:"index"`
	Objects    []string       `json:"objects"`
	AreaPlaces map[string]int `json:"areaPlaces"`
}

func (c Channel) AreaPartitionLabel(areaLabel string) string {
	return fmt.Sprintf("%s##%s", areaLabel, c.Id)
}
