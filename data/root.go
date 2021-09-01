package data

import (
	"context"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
)

// AlertReport defines an alert report
type AlertReport struct {
	Latitude               float64     `json:"latitude"`  // Latitude
	Longitude              float64     `json:"longitude"` // Longitude
	City                   string      `json:"city"`      // City name
	State                  string      `json:"state"`     // State name
	NWSZone                string      `json:"zone"`      // National weather service zone
	NWSZoneURL             string      `json:"zoneurl"`   // National weather service zone URL (for more information)
	ActiveAlertsForZoneURL string      `json:"alertsurl"` // URL to see active alerts on the NWS website for the current NWS zone
	Alerts                 []AlertItem `json:"alerts"`    // Active alerts
	Version                string      `json:"version"`   // Service version
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
	Effective       time.Time `json:"effective"`        // Event effective start time
	Expiration      time.Time `json:"expiration"`       // Event effective end time
}

// AlertService is the interface for all services that can fetch weather alerts
type AlertService interface {
	// GetWeatherReport gets the weather report
	GetWeatherAlerts(ctx context.Context, lat, long string) (AlertReport, error)
}

// GetWeatherAlerts calls all services in parallel and returns the first result
func GetWeatherAlerts(ctx context.Context, services []AlertService, lat, long string) AlertReport {

	ch := make(chan AlertReport, 1)

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "alert-report")
	defer seg.Close(nil)

	//	For each passed service ...
	for _, service := range services {

		//	Launch a goroutine for each service...
		go func(c context.Context, s AlertService, la, lo string) {

			//	Get its pollen report ...
			result, err := s.GetWeatherAlerts(c, la, lo)

			//	As long as we don't have an error, return what we found on the result channel
			if err == nil {
				select {
				case ch <- result:
				default:
				}
			}
		}(ctx, service, lat, long)

	}

	//	Return the first result passed on the channel
	return <-ch
}
