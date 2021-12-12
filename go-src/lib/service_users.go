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

func NewUsersService(usersRepo usersRepo) *usersService {
	return &usersService{
		usersRepo,
	}
}

func (s *usersService) GetUsers(ctx context.Context, limit int, offset int) ([]User, error) {
	return s.usersRepo.GetUsers(ctx, limit, offset)
}

func (s *usersService) GetUser(ctx context.Context, userID string) (User, error) {
	return s.usersRepo.GetUser(ctx, userID)
}
