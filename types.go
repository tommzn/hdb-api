package main

import (
	"net/http"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// WebServer to handleAPI requests.
type webServer struct {
	port        string
	conf        config.Config
	logger      log.Logger
	httpServer  *http.Server
	handlerList []RequestHandler
}

// InddorClimateData contains a single indoor climate measurement.
// Temperature, humidity and battery.
type InddorClimateData struct {
	Timestamp       *string `json:"timestamp"`
	DeviceId        string  `json:"deviceid"`
	MeasurementType string  `json:"measurementtype"`
	Value           string  `json:"value"`
}

// IndoorClimateRequestHandler to handle all indoor climate requests.
type IndoorClimateRequestHandler struct {
	logger log.Logger
}

// HealthRequestHandler to handle health requests.
type HealthRequestHandler struct {
}
