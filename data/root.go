package data

import (
	"time"
)

// AlertReport defines an alert report
type AlertReport struct {
	Latitude                 float64     `json:"latitude"`  // Latitude
	Longitude                float64     `json:"longitude"` // Longitude
	City                     string      `json:"city"`      // City name
	State                    string      `json:"state"`     // State name
	NWSCounty                string      `json:"county"`    // National weather service county
	ActiveAlertsForCountyURL string      `json:"alertsurl"` // URL to see active alerts on the NWS website for the current NWS zone
	Alerts                   []AlertItem `json:"alerts"`    // Active alerts
	Version                  string      `json:"version"`   // Service version
}

type AlertItem struct {
	Event           string    `json:"event"`            // Short event summary
	Headline        string    `json:"headline"`         // Full headline description
	Description     string    `json:"description"`      // Long description of the event
	Severity        string    `json:"severity"`         // Severity of the event
	Urgency         string    `json:"urgency"`          // Urgency of the event
	AreaDescription string    `json:"area_description"` // Affected Area description
	Sender          string    `json:"sender"`           // Sender (email) of the event
	SenderName      string    `json:"sendername"`       // Sender name of the event
	Start           time.Time `json:"start"`            // Event start time
	End             time.Time `json:"end"`              // Event end time
}
