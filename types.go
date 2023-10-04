package main

import (
	"net/http"
	"sync"
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	"google.golang.org/protobuf/proto"
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

// MeasurementUnit e.g. Celsius, Fahrenheit or percent.
type MeasurementUnit string

const (
	// UNIT_CELSIUS, temperature in Celsius.
	UNIT_CELSIUS MeasurementUnit = "C"
	// UNIT_FAHRENHEIT, temperature in Fahrenheit.
	UNIT_FAHRENHEIT MeasurementUnit = "F"
	// UNIT_PERCENT, measurement in percent, e.g. battery level.
	UNIT_PERCENT MeasurementUnit = "%"
)

// ClimateData contains temperate, humidity and battery level measured by a device.
type ClimateData struct {
	DeviceId    string       `json:"deviceid"`
	Location    *string      `json:"location,omitempty"`
	Temperature *Measurement `json:"temperature,omitempty"`
	Humidity    *Measurement `json:"humidity,omitempty"`
	Battery     *Measurement `json:"battery,omitempty"`
}

// Measurement, a single value with a unit.
type Measurement struct {
	Value     string          `json:"value"`
	Unit      MeasurementUnit `json:"unit"`
	Timestamp time.Time       `json:"-"`
}

// IndoorClimateRequestHandler to handle all indoor climate requests.
type IndoorClimateRequestHandler struct {
	logger         log.Logger
	climateData    map[string]ClimateData
	locations      map[string]string
	datasource     DataSource
	dataSourceChan <-chan proto.Message
	mutex          *sync.Mutex
}

// HealthRequestHandler to handle health requests.
type HealthRequestHandler struct {
}
