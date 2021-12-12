package srv

import (
	"fmt"
	"net/http"
)

func NewHealthServerHandler(livenessProbePath, readinessProbePath string) http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc(livenessProbePath, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "live")
	})

	handler.HandleFunc(readinessProbePath, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ready")
	})

	return handler
}
