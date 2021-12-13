package srv

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hbernardo/users/go-src/lib"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	testCases := []struct {
		name              string
		err               error
		expectedHTTPError *httpError
	}{
		{
			name: "already http error 1",
			err: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "msg",
			},
			expectedHTTPError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "msg",
			},
		},
		{
			name: "already http error 2",
			err: &httpError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			},
			expectedHTTPError: &httpError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			},
		},
		{
			name: "service error",
			err:  fmt.Errorf("svc error"),
			expectedHTTPError: &httpError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal server error",
			},
		},
		{
			name: "not found error",
			err:  fmt.Errorf("element: %w", lib.ErrNotFound),
			expectedHTTPError: &httpError{
				StatusCode: http.StatusNotFound,
				Message:    "element: not found",
			},
		},
		{
			name: "precondition failed error",
			err:  fmt.Errorf("invalid: %w", lib.ErrPreconditionFailed),
			expectedHTTPError: &httpError{
				StatusCode: http.StatusPreconditionFailed,
				Message:    "invalid: precondition failed",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpError := handleError(tc.err)

			assert.Equal(t, tc.expectedHTTPError, httpError)
		})
	}
}
