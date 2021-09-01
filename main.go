package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/danesparza/nws-alerts/data"
)

var (
	// BuildVersion contains the version information for the app
	BuildVersion = "Unknown"

	// CommitID is the git commitId for the app.  It's filled in as
	// part of the automated build
	CommitID string
)

// Message is a custom struct event type to handle the Lambda input
type Message struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"long"`
}

// HandleRequest handles the AWS lambda request
func HandleRequest(ctx context.Context, msg Message) (data.AlertReport, error) {
	xray.Configure(xray.Config{LogLevel: "trace"})
	ctx, seg := xray.BeginSegment(ctx, "nws-alerts-lambda-handler")

	service := data.NWSAlertsService{}
	response, err := service.GetWeatherAlerts(ctx, msg.Latitude, msg.Longitude)
	if err != nil {
		seg.Close(err)
		log.Fatalf("problem getting weather alerts: %v", err)
	}

	//	Set the service version information:
	response.Version = fmt.Sprintf("%s.%s", BuildVersion, CommitID)

	//	Close the segment
	seg.Close(nil)

	//	Return our response
	return response, nil
}

func main() {
	//	Immediately forward to Lambda
	lambda.Start(HandleRequest)
}
