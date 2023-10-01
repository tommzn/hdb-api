package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ApplyRoutes registers all methods sipported by this handler.
func (handler *IndoorClimateRequestHandler) applyRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/indoorclimate", handler.listIndoorClimate).Methods("GET")
}

// ListIndoorClimate returns list of availab inddor climate data records.
func (handler *IndoorClimateRequestHandler) listIndoorClimate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

/**

// IndoorClimatePostNodeRequest will publish given indoor climate data.
func (server *webServer) indoorClimatePostNodeRequest(w http.ResponseWriter, r *http.Request) {

	defer server.logger.Flush()

	var inddorClimateData InddorClimateData
	if err := json.NewDecoder(r.Body).Decode(&inddorClimateData); err != nil {
		server.logger.Errorf("Unable to parse request, reason: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := toIndoorCliemateDataEvent(inddorClimateData)
	if err != nil {
		server.logger.Errorf("Unable to convert to indoorclimate event, reason: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = server.messagePublisher.Send(event)
	if err != nil {
		server.logger.Errorf("Unable to publish event, reason: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func toIndoorCliemateDataEvent(inddorClimateData InddorClimateData) (*events.IndoorClimate, error) {

	var timestamp time.Time
	if inddorClimateData.Timestamp != nil {
		parsedTime, err := time.Parse(time.RFC3339, *inddorClimateData.Timestamp)
		if err != nil {
			return nil, err
		}
		timestamp = parsedTime
	} else {
		timestamp = time.Now()
	}

	return &events.IndoorClimate{
		DeviceId:  inddorClimateData.DeviceId,
		Timestamp: timestamppb.New(timestamp),
		Type:      toMeasurementType(inddorClimateData.MeasurementType),
		Value:     inddorClimateData.Value,
	}, nil
}

func toMeasurementType(measurementType string) events.MeasurementType {
	switch measurementType {
	case "humidity":
		return events.MeasurementType_HUMIDITY
	case "battery":
		return events.MeasurementType_BATTERY
	default:
		return events.MeasurementType_TEMPERATURE
	}
}
*/
