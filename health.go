package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// ApplyRoutes registers all methods sipported by this handler.
func (handler *HealthRequestHandler) applyRoutes(router *mux.Router) {
	router.HandleFunc("/health", handler.handleHealthCheckRequest).Methods("GET")
}

// HandleHealthCheckRequest always returns a 204 status code.
func (handler *HealthRequestHandler) handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (handler *HealthRequestHandler) bootstrap(ctx context.Context, waitGroup *sync.WaitGroup) {
	waitGroup.Done()
}
