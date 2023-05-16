package events

import "time"

type StatusChangeOrigin struct {
	Type string `json:"type"`
	Ip   string `json:"ip"`
}

type StatusChange struct {
	Id                      int64              `json:"id"`
	EventId                 int64              `json:"eventId"`
	Status                  string             `json:"status"`
	Date                    *time.Time         `json:"date"`
	OrderId                 string             `json:"orderId"`
	ObjectLabel             string             `json:"objectLabel"`
	ExtraData               ExtraData          `json:"extraData"`
	Origin                  StatusChangeOrigin `json:"origin"`
	IsPresentOnChart        bool               `json:"isPresentOnChart"`
	NotPresentOnChartReason string             `json:"notPresentOnChartReason"`
	HoldToken               string             `json:"holdToken"`
}
