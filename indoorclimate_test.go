package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IndoorClimateRequestHandlerTestSuite struct {
	suite.Suite
}

func TestIndoorClimateRequestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthRequestHandlerTestSuite))
}

func (suite *IndoorClimateRequestHandlerTestSuite) TestListRequest() {

	handler := indoorClimateRequestHandlerForTest()
	req, err := http.NewRequest("GET", "/api/v1/indoorclimate", nil)
	suite.Nil(err)
	respRecorder := httptest.NewRecorder()

	handler.listIndoorClimate(respRecorder, req)
	suite.Equal(http.StatusNoContent, respRecorder.Code)
}
