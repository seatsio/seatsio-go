package charts

import "github.com/seatsio/seatsio-go/events"

type Chart struct {
	Id                           int64          `json:"id"`
	Key                          string         `json:"key"`
	Name                         string         `json:"name"`
	Status                       string         `json:"status"`
	Tags                         []string       `json:"tags"`
	PublishedVersionThumbnailUrl string         `json:"publishedVersionThumbnailUrl"`
	DraftVersionThumbnailUrl     string         `json:"draftVersionThumbnailUrl"`
	Events                       []events.Event `json:"events"`
	Archived                     bool           `json:"archived"`
}

type Tags struct {
	Tags []string `json:"tags"`
}
