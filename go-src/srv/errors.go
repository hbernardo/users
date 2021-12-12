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

func (e *httpError) Error() string {
	return fmt.Sprintf("%s (status code: %d)", e.Message, e.StatusCode)
}

func writeError(w http.ResponseWriter, err error) {
	httpError := handleError(err)
	writeJSON(w, httpError.StatusCode, httpError)
}

func handleError(err error) *httpError {
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode >= 500 {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("http server error")
		}
		return httpErr
	}

	if errors.Is(err, lib.ErrNotFound) {
		return &httpError{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
	}

	if errors.Is(err, lib.ErrPreconditionFailed) {
		return &httpError{
			StatusCode: http.StatusPreconditionFailed,
			Message:    err.Error(),
		}
	}

	log.WithFields(log.Fields{
		"error": err.Error(),
	}).Error("internal error")

	return &httpError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}
}
