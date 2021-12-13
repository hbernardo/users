package srv

import (
	"context"
	"net/http"

	"github.com/hbernardo/users/go-src/lib"
)

type (
	usersService interface {
		GetUsers(ctx context.Context, limit int, offset int) ([]lib.User, error)
		GetUser(ctx context.Context, userID string) (lib.User, error)
	}

	usersHandler struct {
		http.Handler
		usersService
	}
)

const (
	// maxUsersLimit sets the maximum number of users
	// that the client can request to the server
	maxUsersLimit = 1000
)

// NewUsersHandler creates a new users handler, receives the users service as parameter
func NewUsersHandler(usersSvc usersService) *usersHandler {
	handler := http.NewServeMux()

	h := &usersHandler{
		handler,
		usersSvc,
	}

	// route for multiple users fetching, receiving pagination parameters
	handler.HandleFunc("/v1/users", h.handleGetUsers)
	// route for single user fetching, receiving the user id as URL parameter
	handler.HandleFunc("/v1/users/", h.handleGetUser)

	return h
}

// handleGetUsers is the HTTP handler function for getting multiples based on pagination querystrings (limit and offset)
func (h *usersHandler) handleGetUsers(w http.ResponseWriter, req *http.Request) {
	// validating GET method
	if req.Method != http.MethodGet {
		writeError(w, &httpError{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "method not allowed",
		})
		return
	}

	// getting and validating pagination parameters
	limit, offset, err := getAndValidatePaginationParams(req.URL.Query(), maxUsersLimit)
	if err != nil {
		writeError(w, err)
		return
	}

	users, err := h.usersService.GetUsers(req.Context(), limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// handleGetUser is the HTTP handler function for getting a single user by its ID (got from URL parameter)
func (h *usersHandler) handleGetUser(w http.ResponseWriter, req *http.Request) {
	// validating GET method
	if req.Method != http.MethodGet {
		writeError(w, &httpError{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "method not allowed",
		})
		return
	}

	// getting user id from URL parameter
	userID, err := getURLPathParam(req.URL.Path, "users")
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.usersService.GetUser(req.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
