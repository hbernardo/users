package infra

import (
	"context"
	"fmt"

	"github.com/hbernardo/users/go-src/lib"
)

type (
	usersRepo struct {
		usersData []lib.User
		usersMap  map[string]*lib.User
	}
)

func NewUsersRepo(usersData []lib.User) *usersRepo {
	// creating a map for direct access when querying a single user
	usersMap := make(map[string]*lib.User, len(usersData))
	for _, user := range usersData {
		usersMap[user.ID] = &user
	}

	return &usersRepo{
		usersData,
		usersMap,
	}
}

func (r *usersRepo) GetUsers(ctx context.Context, limit int, offset int) ([]lib.User, error) {
	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("'limit' nor 'offset' cannot be negative: %w", lib.ErrPreconditionFailed)
	}

	if offset > len(r.usersData) {
		offset = len(r.usersData)
	}
	if (offset + limit) > len(r.usersData) {
		limit = len(r.usersData) - offset
	}

	return r.usersData[offset : offset+limit], nil
}

func (r *usersRepo) GetUser(ctx context.Context, userID string) (lib.User, error) {
	user, userExists := r.usersMap[userID]
	if !userExists || user == nil {
		return lib.User{}, lib.ErrNotFound
	}

	return *user, nil
}
