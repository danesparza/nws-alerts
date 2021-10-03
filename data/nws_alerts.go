package data

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
)

// NWSPointsResponse defines the expected response format from the NWS points service
type NWSPointsResponse struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		ID                  string `json:"@id"`
		Type                string `json:"@type"`
		Cwa                 string `json:"cwa"`
		ForecastOffice      string `json:"forecastOffice"`
		GridID              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		Forecast            string `json:"forecast"`
		ForecastHourly      string `json:"forecastHourly"`
		ForecastGridData    string `json:"forecastGridData"`
		ObservationStations string `json:"observationStations"`
		RelativeLocation    struct {
			Type     string `json:"type"`
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				City     string `json:"city"`  // City name
				State    string `json:"state"` // State name
				Distance struct {
					Value    float64 `json:"value"`
					UnitCode string  `json:"unitCode"`
				} `json:"distance"`
				Bearing struct {
					Value    int    `json:"value"`
					UnitCode string `json:"unitCode"`
				} `json:"bearing"`
			} `json:"properties"`
		} `json:"relativeLocation"`
		ForecastZone    string `json:"forecastZone"` // Forecast zone
		County          string `json:"county"`
		FireWeatherZone string `json:"fireWeatherZone"`
		TimeZone        string `json:"timeZone"`
		RadarStation    string `json:"radarStation"`
	} `json:"properties"`
}

// NWSAlertsResponse defines the expected response from the NWS alerts service (for a specific zone)
type NWSAlertsResponse struct {
	Type     string `json:"type"`
	Features []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Properties struct {
			ID       string `json:"@id"`
			Type     string `json:"@type"`
			AreaDesc string `json:"areaDesc"`
			Geocode  struct {
				SAME []string `json:"SAME"`
				UGC  []string `json:"UGC"`
			} `json:"geocode"`
			AffectedZones []string `json:"affectedZones"`
			References    []struct {
				ID         string `json:"@id"`
				Identifier string `json:"identifier"`
				Sender     string `json:"sender"`
				Sent       string `json:"sent"`
			} `json:"references"`
			Sent        time.Time `json:"sent"`
			Effective   time.Time `json:"effective"`
			Onset       time.Time `json:"onset"`
			Expires     time.Time `json:"expires"`
			Ends        time.Time `json:"ends"`
			Status      string    `json:"status"`
			MessageType string    `json:"messageType"`
			Category    string    `json:"category"`
			Severity    string    `json:"severity"`
			Certainty   string    `json:"certainty"`
			Urgency     string    `json:"urgency"`
			Event       string    `json:"event"`
			Sender      string    `json:"sender"`
			SenderName  string    `json:"senderName"`
			Headline    string    `json:"headline"`
			Description string    `json:"description"`
			Instruction string    `json:"instruction"`
			Response    string    `json:"response"`
			Parameters  struct {
				PIL               []string    `json:"PIL"`
				NWSheadline       []string    `json:"NWSheadline"`
				BLOCKCHANNEL      []string    `json:"BLOCKCHANNEL"`
				EASORG            []string    `json:"EAS-ORG"`
				VTEC              []string    `json:"VTEC"`
				EventEndingTime   []time.Time `json:"eventEndingTime"`
				ExpiredReferences []string    `json:"expiredReferences"`
			} `json:"parameters"`
		} `json:"properties"`
	} `json:"features"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
}

// OpenWeatherService is a weather service for OpenWeather formatted data
type NWSAlertsService struct{}

// GetWeatherReport gets the weather report
func (s NWSAlertsService) GetWeatherAlerts(ctx context.Context, lat, long string) (AlertReport, error) {
	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "nwsalert-service")

	//	Our return value
	retval := AlertReport{}
	retval.Alerts = []AlertItem{} // Initialize the array

	//	Call the alerts service for the lat/long specified
	alertsServiceUrl := fmt.Sprintf("https://api.weather.gov/alerts?point=%s,%s", lat, long)
	alertClientRequest, err := http.NewRequest("GET", alertsServiceUrl, nil)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("problem creating request to the NWS alerts service: %v", err)
	}

	//	Set our headers
	alertClientRequest.Header.Set("Content-Type", "application/geo+json; charset=UTF-8")

	//	Execute the request
	alertClient := &http.Client{}
	alertClientResponse, err := ctxhttp.Do(ctx, xray.Client(alertClient), alertClientRequest)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("error when sending request to the NWS alerts service: %v", err)
	}
	defer alertClientResponse.Body.Close()

	//	Decode the response:
	alertsResponse := NWSAlertsResponse{}
	err = json.NewDecoder(alertClientResponse.Body).Decode(&alertsResponse)
	if err != nil {
		seg.AddError(err)
		return retval, fmt.Errorf("problem decoding the response from the NWS alerts service: %v", err)
	}

	seg.AddMetadata("AlertsResponse", alertsResponse)

	//	Compile our report
	for _, item := range alertsResponse.Features {

		alertItem := AlertItem{
			Event:           item.Properties.Event,
			Headline:        item.Properties.Headline,
			Description:     item.Properties.Description,
			Severity:        item.Properties.Severity,
			Urgency:         item.Properties.Urgency,
			AreaDescription: item.Properties.AreaDesc,
			Sender:          item.Properties.Sender,
			SenderName:      item.Properties.SenderName,
			Start:           item.Properties.Effective,
			End:             item.Properties.Ends,
		}

		retval.Alerts = append(retval.Alerts, alertItem)
	}

	//	Add the report to the request metadata
	xray.AddMetadata(ctx, "AlertResult", retval)

	// Close the segment
	seg.Close(nil)

	//	Return the report
	return retval, nil
}
