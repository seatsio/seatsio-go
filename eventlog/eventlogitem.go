package eventlog

import "time"

type EventLogItem struct {
	Id           int64          `json:"id"`
	Type         string         `json:"type"`
	WorkspaceKey string         `json:"workspaceKey"`
	Date         *time.Time     `json:"date"`
	Data         map[string]any `json:"data"`
}
