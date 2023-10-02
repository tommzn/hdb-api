package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

// NewServer returns a new HTTP server.
func newServer(conf config.Config, logger log.Logger, handlerList []RequestHandler) *webServer {
	port := conf.Get("hdb.api.port", config.AsStringPtr("8080"))
	return &webServer{
		port:        *port,
		conf:        conf,
		logger:      logger,
		handlerList: handlerList,
	}
}

// Run starts a HTTP server to listen for rendering requests.
func (server *webServer) Run(ctx context.Context, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	defer server.logger.Flush()

	router := mux.NewRouter()
	router.Use(server.logMiddleware)

	for _, handler := range server.handlerList {
		waitGroup.Add(1)
		handler.bootstrap(ctx, waitGroup)
		handler.applyRoutes(router)
	}

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
	server.httpServer.Shutdown(ctx)
}

// LogMiddleware adds a logger for all requests. Used log level if debug.
func (server *webServer) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.logger.Debugf("Method: %s, URL: %+v, Header: %+v, URI: %s", r.Method, r.URL, r.Header, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
