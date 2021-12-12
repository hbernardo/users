package srv

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	httpServer struct {
		srv *http.Server
	}
)

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(port int, handler http.Handler) *httpServer {
	return &httpServer{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}

// ListenAndServe starts the HTTP server in the background
func (h *httpServer) ListenAndServe() {
	go func(srv *http.Server) {
		log.Fatal(srv.ListenAndServe())
	}(h.srv)
}

// Close closes the HTTP server
func (h *httpServer) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return h.srv.Shutdown(ctx)
}
