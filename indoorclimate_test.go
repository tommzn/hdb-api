package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type IndoorClimateRequestHandlerTestSuite struct {
	suite.Suite
}

func TestIndoorClimateRequestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(IndoorClimateRequestHandlerTestSuite))
}

func (suite *IndoorClimateRequestHandlerTestSuite) TestListRequest() {

	handler := suite.bootstrapHandler(indoorClimateRequestHandlerForTest())
	time.Sleep(3 * time.Second)

	req, err := http.NewRequest("GET", "/api/v1/indoorclimate", nil)
	suite.Nil(err)
	respRecorder := httptest.NewRecorder()

	handler.listIndoorClimate(respRecorder, req)
	suite.Equal(http.StatusOK, respRecorder.Code)
}

func (suite *IndoorClimateRequestHandlerTestSuite) TestGetLocationForDeviceId() {

	devideId := "device01"
	location := "location01"
	locations := make(map[string]string)
	locations[devideId] = location

	location1 := locationForDevice(locations, devideId)
	suite.NotNil(location1)
	suite.Equal(location, *location1)

	suite.Nil(locationForDevice(locations, "xxx"))
}

func (suite *IndoorClimateRequestHandlerTestSuite) TestGetLocationFromConfig() {

	conf := loadConfigForTest(nil)

	locations := locationsFromConfig(conf, "hdb.locations")
	suite.Len(locations, 1)
}

func (suite *IndoorClimateRequestHandlerTestSuite) TestFormatValues() {

	suite.Equal("32.2", format("32.16", 1))
	suite.Equal("32.1", format("32.11", 1))
	suite.Equal("58", format("58.11", 0))
	suite.Equal("60", format("59.89", 0))
	suite.Equal("ccc", format("ccc", 1))
}

func (suite *IndoorClimateRequestHandlerTestSuite) bootstrapHandler(handler *IndoorClimateRequestHandler) *IndoorClimateRequestHandler {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	handler.bootstrap(context.Background(), wg)
	return handler
}
