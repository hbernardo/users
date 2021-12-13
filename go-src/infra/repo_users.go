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

// NewUsersRepo creates a new users repo, receives the users data as parameter
func NewUsersRepo(usersData []lib.User) *usersRepo {
	repo := &usersRepo{
		usersData: usersData,
		usersMap:  make(map[string]*lib.User, len(usersData)),
	}

	// map for direct/instant access when querying a single user
	for i, user := range repo.usersData {
		repo.usersMap[user.ID] = &(repo.usersData[i])
	}

	return repo
}

// GetUsers gets users based on pagination (limit and offset)
func (r *usersRepo) GetUsers(ctx context.Context, limit int, offset int) ([]lib.User, error) {
	// validating pagination parameters
	if limit < 0 || offset < 0 {
		return nil, fmt.Errorf("'limit' nor 'offset' cannot be negative: %w", lib.ErrPreconditionFailed)
	}

	// fixing out of bonds slice access
	// empty slice should be expected in that case
	if offset > len(r.usersData) {
		offset = len(r.usersData)
	}
	if (offset + limit) > len(r.usersData) {
		limit = len(r.usersData) - offset
	}

	return r.usersData[offset : offset+limit], nil
}

// GetUser gets user based on its ID
func (r *usersRepo) GetUser(ctx context.Context, userID string) (lib.User, error) {
	// direct access to queried user
	// returning "not found" error if user doesn't exists
	user, userExists := r.usersMap[userID]
	if !userExists || user == nil {
		return lib.User{}, lib.ErrNotFound
	}

	return *user, nil
}
