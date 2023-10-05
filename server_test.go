package main

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

type ServerTestSuite struct {
	suite.Suite
	conf       config.Config
	logger     log.Logger
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.conf = loadConfigForTest(nil)
	suite.logger = loggerForTest()
}

func (suite *ServerTestSuite) SetupTest() {
	suite.ctx, suite.cancelFunc = context.WithCancel(context.Background())
}

func (suite *ServerTestSuite) TestHealthRequest() {

	server := suite.serverForTest()
	suite.startServer(server)

	resp, err := http.Get("http://localhost/v1/health")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusNoContent, resp.StatusCode)

	suite.stopServer()
}

func (suite *ServerTestSuite) TestGetIndoorClimateData() {

	server := suite.serverForTest()
	suite.startServer(server)

	resp, err := http.Get("http://localhost/v1/indoorclimate")
	suite.Nil(err)
	suite.NotNil(resp)
	suite.Equal(http.StatusOK, resp.StatusCode)

	suite.stopServer()
}

func (suite *ServerTestSuite) startServer(server *webServer) {
	suite.wg = &sync.WaitGroup{}
	go func() {
		suite.wg.Add(1)
		suite.Nil(server.Run(suite.ctx, suite.wg))
	}()
	time.Sleep(2 * time.Second)
}

func (suite *ServerTestSuite) serverForTest() *webServer {
	handlerList := []RequestHandler{
		indoorClimateRequestHandlerForTest(),
		&HealthRequestHandler{},
	}
	return newServer(suite.conf, suite.logger, handlerList)
}

func (suite *ServerTestSuite) stopServer() {

	waitChan := make(chan bool, 1)
	go func() {
		suite.wg.Wait()
		waitChan <- true
	}()

	suite.cancelFunc()
	time.Sleep(1 * time.Second)
	select {
	case <-time.After(1 * time.Second):
		suite.T().Error("Server stop timeput!")
	case ok := <-waitChan:
		suite.True(ok)
	}
}
