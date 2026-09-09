package reports

import (
	"github.com/seatsio/seatsio-go/v13/events"
)

func findByLabel(report []events.EventObjectInfo, label string) events.EventObjectInfo {
	for _, item := range report {
		if item.Label == label {
			return item
		}
	}
	panic("no report item found with label " + label)
}

func findByLabelInDetailedReport(items map[string][]events.EventObjectInfo, label string) events.EventObjectInfo {
	var flat []events.EventObjectInfo
	for _, v := range items {
		flat = append(flat, v...)
	}
	return findByLabel(flat, label)
}
