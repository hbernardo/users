package srv

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type (
	httpServer struct {
		srv *http.Server
	}
)

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(port int) *httpServer {
	return &httpServer{
		srv: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
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
	return h.srv.Shutdown(ctx)
}
