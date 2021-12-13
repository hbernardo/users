package srv

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hbernardo/users/go-src/lib"
	log "github.com/sirupsen/logrus"
)

// httpError wraps errors and adds HTTP specific info
type httpError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
}

// Error formats the error in a descriptive format (required for the custom error)
func (e *httpError) Error() string {
	return fmt.Sprintf("%s (status code: %d)", e.Message, e.StatusCode)
}

// writeError handles, formats and writes error to the HTTP response
func writeError(w http.ResponseWriter, err error) {
	httpError := handleError(err)
	writeJSON(w, httpError.StatusCode, httpError)
}

// handleError handles the error properly to have the final HTTP error
func handleError(err error) *httpError {
	var httpErr *httpError
	// just return it's already a HTTP error
	if errors.As(err, &httpErr) {
		// log if internal server error
		if httpErr.StatusCode >= 500 {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("http server error")
		}
		return httpErr
	}

	// not found error converting to HTTP error
	if errors.Is(err, lib.ErrNotFound) {
		return &httpError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
	}

	// precondition failed error converting to HTTP error
	if errors.Is(err, lib.ErrPreconditionFailed) {
		return &httpError{
			StatusCode: http.StatusPreconditionFailed,
			Message:    err.Error(),
		}
	}

	// default error handling:
	// - log as internal error
	// - convert to HTTP internal server error

	log.WithFields(log.Fields{
		"error": err.Error(),
	}).Error("internal error")

	return &httpError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}
}
