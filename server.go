package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newServer(conf config.Config, logger log.Logger, messagePublisher core.Publisher) *webServer {
	port := conf.Get("hdb.api.port", config.AsStringPtr("8080"))
	return &webServer{
		port:             *port,
		conf:             conf,
		logger:           logger,
		messagePublisher: messagePublisher,
	}
}

// Run starts a HTTP server to listen for rendering requests.
func (server *webServer) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	defer server.logger.Flush()

	router := mux.NewRouter()
	router.Use(server.logMiddleware)

	router.HandleFunc("/api/v1/indoorclimate", server.indoorClimatePostNodeRequest).Methods("POST")
	router.HandleFunc("/health", server.handleHealthCheckRequest).Methods("GET")

	server.logger.Infof("Listen [%s]", server.port)
	server.logger.Flush()
	server.httpServer = &http.Server{Addr: ":" + server.port, Handler: router}

	endChan := make(chan error, 1)
	go func() {
		endChan <- server.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		server.stopHttpServer()
	case err := <-endChan:
		return err
	}
	return nil
}

// StopHttpServer will try to sop running HTTP server graceful. Timeout is 3s.
func (server *webServer) stopHttpServer() {
	server.logger.Info("Stopping HTTP server.")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("Unable to stop HTTP server, reason: ", err)
	}
}

// LogMiddleware adds a logger for all requests. Used log level if debug.
func (server *webServer) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.logger.Debugf("Method: %s, URL: %+v, Header: %+v, URI: %s", r.Method, r.URL, r.Header, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// HandleHealthCheckRequest always returns a 204 status code.
func (server *webServer) handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

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
