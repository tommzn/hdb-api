package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	hdbcore "github.com/tommzn/hdb-core"
	events "github.com/tommzn/hdb-events-go"
)

// ApplyRoutes registers all methods sipported by this handler.
func (handler *IndoorClimateRequestHandler) applyRoutes(router *mux.Router) {
	router.HandleFunc("/v1/indoorclimate", handler.listIndoorClimate).Methods("GET")
}

// ListIndoorClimate returns list of availab inddor climate data records.
func (handler *IndoorClimateRequestHandler) listIndoorClimate(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(handler.climateData)
}

// Bootstrap will listen for new indoor climate data provided by used datasource.
func (handler *IndoorClimateRequestHandler) bootstrap(ctx context.Context, waitGroup *sync.WaitGroup) {

	waitGroup.Add(1)
	go handler.datasource.Run(ctx, waitGroup)
	handler.mutex = &sync.Mutex{}

	filter := []hdbcore.DataSource{hdbcore.DATASOURCE_INDOORCLIMATE}
	handler.dataSourceChan = handler.datasource.Observe(&filter)
	go func() {
		for {
			select {
			case message := <-handler.dataSourceChan:
				if indoorClimate, ok := message.(*events.IndoorClimate); ok {
					handler.addAsClimateData(indoorClimate)
				}

			case <-ctx.Done():
				handler.logger.Info("Camceled, stop observing.")
				handler.logger.Flush()
				waitGroup.Done()
				return
			}
		}
	}()
}

// addAsClimateData will try to add passed message to local indoor climate data.
func (handler *IndoorClimateRequestHandler) addAsClimateData(indoorClimate *events.IndoorClimate) {

	deviceId := strings.ToUpper(indoorClimate.DeviceId)
	climateData := handler.climateDataForDevice(deviceId)

	handler.logger.Debugf("Receive new indoor climate data, %s, %s", indoorClimate.Type, indoorClimate.Value)
	switch indoorClimate.Type {

	case events.MeasurementType_TEMPERATURE:
		if climateData.Temperature == nil ||
			climateData.Temperature.Timestamp.Before(indoorClimate.Timestamp.AsTime()) {
			climateData.Temperature = &Measurement{
				Value:     format(indoorClimate.Value, 1),
				Unit:      UNIT_CELSIUS,
				Timestamp: indoorClimate.Timestamp.AsTime(),
			}
		}

	case events.MeasurementType_HUMIDITY:
		if climateData.Humidity == nil ||
			climateData.Humidity.Timestamp.Before(indoorClimate.Timestamp.AsTime()) {
			climateData.Humidity = &Measurement{
				Value:     format(indoorClimate.Value, 0),
				Unit:      UNIT_PERCENT,
				Timestamp: indoorClimate.Timestamp.AsTime(),
			}
		}

	case events.MeasurementType_BATTERY:
		if climateData.Battery == nil ||
			climateData.Battery.Timestamp.Before(indoorClimate.Timestamp.AsTime()) {
			climateData.Battery = &Measurement{
				Value:     format(indoorClimate.Value, 0),
				Unit:      UNIT_PERCENT,
				Timestamp: indoorClimate.Timestamp.AsTime(),
			}
		}
	}
	handler.putClimateData(climateData, deviceId)
}

// PutClimateData write given climate data to internal map, concurrency save.
func (handler *IndoorClimateRequestHandler) putClimateData(climateData ClimateData, deviceId string) {
	handler.mutex.Lock()
	handler.climateData[deviceId] = climateData
	handler.mutex.Unlock()
}

// ClimateDataForDevice will try to get climate data for passed device from local storage.
// If none exists a new, empty cliemate data wil be returned.
func (handler *IndoorClimateRequestHandler) climateDataForDevice(deviceId string) ClimateData {

	if climateData, ok := handler.climateData[deviceId]; ok {
		return climateData
	}

	climateData := ClimateData{
		DeviceId: deviceId,
		Location: locationForDevice(handler.locations, deviceId),
	}
	handler.putClimateData(climateData, deviceId)
	return climateData
}

// LocationForDevice lookup if there's a location for given device.
// Returns null if no location has been defined for passed device id.
func locationForDevice(locations map[string]string, deviceId string) *string {
	if location, ok := locations[deviceId]; ok {
		return &location
	}
	return nil
}

// Format passed value to float with given decimal places.
func format(value string, decimals int) string {
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return fmt.Sprintf("%."+fmt.Sprintf("%d", decimals)+"f", floatValue)
	}
	return value
}
