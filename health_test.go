package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HealthRequestHandlerTestSuite struct {
	suite.Suite
}

func TestHealthRequestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HealthRequestHandlerTestSuite))
}

func (suite *HealthRequestHandlerTestSuite) TestHealthGetRequest() {

	handler := HealthRequestHandler{}
	req, err := http.NewRequest("GET", "/health", nil)
	suite.Nil(err)
	respRecorder := httptest.NewRecorder()

	handler.handleHealthCheckRequest(respRecorder, req)
	suite.Equal(http.StatusNoContent, respRecorder.Code)
}
