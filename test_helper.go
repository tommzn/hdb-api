package main

import (
	"time"

	"github.com/golang/protobuf/proto"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

func loggerForTest() log.Logger {
	return log.NewLogger(log.Debug, nil, nil)
}

func indoorClimateRequestHandlerForTest() *IndoorClimateRequestHandler {

	handler := &IndoorClimateRequestHandler{
		logger:      loggerForTest(),
		climateData: make(map[string]ClimateData),
		locations:   map[string]string{"Device1": "Location1"},
		datasource:  datasourceMockForTest(),
	}
	return handler
}

func datasourceMockForTest() *datasourceMock {
	mock := &datasourceMock{
		messages:       []proto.Message{},
		offset:         0,
		delay:          500 * time.Millisecond,
		dataSourceChan: make(chan proto.Message, 10),
	}

	mock.messages = []proto.Message{
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device5",
			Type:      events.MeasurementType_BATTERY,
			Value:     "23",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_BATTERY,
			Value:     "23",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_TEMPERATURE,
			Value:     "23.5",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_HUMIDITY,
			Value:     "57",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_TEMPERATURE,
			Value:     "17.1",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device2",
			Type:      events.MeasurementType_HUMIDITY,
			Value:     "65",
		},
		&events.IndoorClimate{
			Timestamp: timestamppb.New(time.Now()),
			DeviceId:  "Device1",
			Type:      events.MeasurementType_BATTERY,
			Value:     "97",
		},
	}
	return mock
}
