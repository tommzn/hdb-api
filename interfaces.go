package main

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	core "github.com/tommzn/hdb-core"
)

// RequestHandler is a generic interface for each kind of requests.
type RequestHandler interface {

	// Bootstrap runs request handler initialization.
	bootstrap(context.Context, *sync.WaitGroup)

	// ApplyRoutes registers all methods a request hadnler suppors.
	applyRoutes(*mux.Router)
}

// DataSource is used to obtain events from different sources.
type DataSource interface {

	// Runable core interface. Used to run message fetch n background.
	core.Runable

	// Observe returns a channel clients can subsribe to get new messages.
	Observe(*[]core.DataSource) <-chan proto.Message
}
