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
	maxNumberOfUsersInResponse = 1000
)

func NewUsersHandler(usersSvc usersService) *usersHandler {
	handler := http.NewServeMux()

	h := &usersHandler{
		handler,
		usersSvc,
	}

	handler.HandleFunc("/users", h.handleGetUsers)
	handler.HandleFunc("/users/", h.handleGetUser)

	return h
}

func (h *usersHandler) handleGetUsers(w http.ResponseWriter, req *http.Request) {
	limit, offset, err := getAndValidatePaginationParams(req)
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

func (h *usersHandler) handleGetUser(w http.ResponseWriter, req *http.Request) {
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
