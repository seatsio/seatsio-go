package charts

import "github.com/seatsio/seatsio-go/v11/events"

type Chart struct {
	Id                           int64                  `json:"id"`
	Key                          string                 `json:"key"`
	Name                         string                 `json:"name"`
	Status                       string                 `json:"status"`
	Tags                         []string               `json:"tags"`
	PublishedVersionThumbnailUrl string                 `json:"publishedVersionThumbnailUrl"`
	DraftVersionThumbnailUrl     string                 `json:"draftVersionThumbnailUrl"`
	Events                       []events.Event         `json:"events"`
	Archived                     bool                   `json:"archived"`
	Validation                   *ChartValidationResult `json:"validation"`
	VenueType                    string                 `json:"venueType"`
	Zones                        []Zone                 `json:"zones"`
}

type Tags struct {
	Tags []string `json:"tags"`
}
