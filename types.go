package main

import (
	"net/http"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-datasource-core"
)

type webServer struct {
	port             string
	conf             config.Config
	logger           log.Logger
	messagePublisher core.Publisher
	httpServer       *http.Server
}

type InddorClimateData struct {
	Timestamp       *string `json:"timestamp"`
	DeviceId        string  `json:"deviceid"`
	MeasurementType string  `json:"measurementtype"`
	Value           string  `json:"value"`
}
