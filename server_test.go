package main

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	conf       config.Config
	logger     log.Logger
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
	mock       *publisherMock
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.conf = loadConfigForTest(nil)
	suite.logger = loggerForTest()
	suite.mock = newPublisherMock()
}

func (suite *ServerTestSuite) SetupTest() {
	suite.ctx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *ServerTestSuite) TestHealthRequest() {

	server := suite.serverForTest()
	suite.startServer(server)

	resp, err := http.Get("http://localhost:8080/health")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	suite.stopServer()
}

func (suite *ServerTestSuite) TestPublishIndoorClimateData() {

	server := suite.serverForTest()
	suite.startServer(server)

	indoorClimateDataContent, err := os.ReadFile("fixtures/indoorcliamtedata.json")
	suite.Nil(err)

	resp, err := http.Post("http://localhost:8080/api/v1/indoorclimate", "application/json", bytes.NewReader(indoorClimateDataContent))
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	suite.Equal(2, suite.mock.messageCount)
	suite.stopServer()
}

func (suite *ServerTestSuite) startServer(server *webServer) {
	suite.wg = &sync.WaitGroup{}
	go func() {
		suite.wg.Add(1)
		suite.Nil(server.Run(suite.ctx, suite.wg))
	}()
	time.Sleep(1 * time.Second)
}

func (suite *ServerTestSuite) serverForTest() *webServer {
	return newServer(suite.conf, suite.logger, suite.mock)
}

func (suite *ServerTestSuite) stopServer() {

	waitChan := make(chan bool, 1)
	go func() {
		suite.wg.Wait()
		waitChan <- true
	}()

	suite.cancelFunc()
	select {
	case <-time.After(1 * time.Second):
		suite.T().Error("Server stop timeput!")
	case ok := <-waitChan:
		suite.True(ok)
	}
}
