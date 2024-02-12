package eventlog

import "time"

type EventLogItem struct {
	Id        int64          `json:"id"`
	Type      string         `json:"type"`
	Timestamp *time.Time     `json:"timestamp"`
	Data      map[string]any `json:"data"`
}
