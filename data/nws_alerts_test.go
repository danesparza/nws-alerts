package data_test

import (
	"context"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/danesparza/nws-alerts/data"
)

func TestNWSAlerts_GetWeatherAlerts_ReturnsValidData(t *testing.T) {
	//	Arrange
	service := data.NWSAlertsService{}
	lat := "41.256989522114395"
	long := "-70.11151820621481"
	ctx := context.Background()
	ctx, seg := xray.BeginSegment(ctx, "unit-test")
	defer seg.Close(nil)

	//	Act
	response, err := service.GetWeatherAlerts(ctx, lat, long)

	//	Assert
	if err != nil {
		t.Errorf("Error calling GetAlerts: %v", err)
	}

	t.Logf("Returned object: %+v", response)

}
