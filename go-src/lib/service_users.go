package lib

import (
	"context"
)

type (
	usersRepo interface {
		GetUsers(ctx context.Context, limit int, offset int) ([]User, error)
		GetUser(ctx context.Context, userID string) (User, error)
	}

	usersService struct {
		usersRepo
	}
)

// NewUsersRepo creates a new users service, receives the users repo as parameter
func NewUsersService(usersRepo usersRepo) *usersService {
	return &usersService{
		usersRepo,
	}
}

// GetUsers gets users based on pagination (limit and offset)
func (s *usersService) GetUsers(ctx context.Context, limit int, offset int) ([]User, error) {
	// only forwarding request to repo, no extra logic required for now
	return s.usersRepo.GetUsers(ctx, limit, offset)
}

// GetUser gets user based on its ID
func (s *usersService) GetUser(ctx context.Context, userID string) (User, error) {
	// only forwarding request to repo, no extra logic required for now
	return s.usersRepo.GetUser(ctx, userID)
}
