package srv

import (
	"net/http"
)

// NewHealthHandler creates a new HTTP handler for heath requests (liveness and readiness)
func NewHealthHandler(livenessProbePath, readinessProbePath string) http.Handler {
	handler := http.NewServeMux()

	// for now, it's only necessary to return success status code (200)
	handler.HandleFunc(livenessProbePath, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// for now, it's only necessary to return success status code (200)
	handler.HandleFunc(readinessProbePath, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return handler
}
