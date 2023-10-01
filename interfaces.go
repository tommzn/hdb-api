package main

import "github.com/gorilla/mux"

// RequestHandler is a generic interface for each kind of requests.
type RequestHandler interface {

	// ApplyRoutes registers all methods a request hadnler suppors.
	applyRoutes(*mux.Router)
}
