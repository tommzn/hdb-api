package main

import (
	"net/http"

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
