package srv

import (
	"net/http"
)

func NewHealthHandler(livenessProbePath, readinessProbePath string) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc(livenessProbePath, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.HandleFunc(readinessProbePath, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return handler
}
