package lib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsersRepo struct {
	mock.Mock
}

func (m *mockUsersRepo) GetUsers(ctx context.Context, limit int, offset int) ([]User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]User), args.Error(1)
}

func (m *mockUsersRepo) GetUser(ctx context.Context, userID string) (User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Error(1)
}

func TestGetUsers(t *testing.T) {
	testCases := []struct {
		name          string
		limit         int
		offset        int
		repoResponse  []User
		repoError     error
		expectedUsers []User
		expectedError error
	}{
		{
			name:   "base case",
			limit:  1,
			offset: 5,
			repoResponse: []User{
				{
					ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
					FirstName:    "Terrence",
					LastName:     "Trillow",
					Email:        "ttrillow1@feedburner.com",
					Password:     "5YLItbmdkfC1",
					IPAddress:    "63.119.6.98",
					CreationDate: "19/04/2021",
				},
			},
			repoError: nil,
			expectedUsers: []User{
				{
					ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
					FirstName:    "Terrence",
					LastName:     "Trillow",
					Email:        "ttrillow1@feedburner.com",
					Password:     "5YLItbmdkfC1",
					IPAddress:    "63.119.6.98",
					CreationDate: "19/04/2021",
				},
			},
			expectedError: nil,
		},
		{
			name:          "repo error",
			limit:         -1,
			offset:        -1,
			repoResponse:  nil,
			repoError:     ErrPreconditionFailed,
			expectedUsers: nil,
			expectedError: ErrPreconditionFailed,
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersRepo := new(mockUsersRepo)

			ctx := context.Background()

			mockUsersRepo.On("GetUsers", ctx, tc.limit, tc.offset).Return(tc.repoResponse, tc.repoError)

			svc := NewUsersService(mockUsersRepo)

			users, err := svc.GetUsers(ctx, tc.limit, tc.offset)

			mockUsersRepo.AssertExpectations(t)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUsers, users)
		})
	}
}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		repoResponse  User
		repoError     error
		expectedUser  User
		expectedError error
	}{
		{
			name:   "base case",
			userID: "1311f914-1d4f-40b6-8886-80193265d5a4",
			repoResponse: User{
				ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
				FirstName:    "Terrence",
				LastName:     "Trillow",
				Email:        "ttrillow1@feedburner.com",
				Password:     "5YLItbmdkfC1",
				IPAddress:    "63.119.6.98",
				CreationDate: "19/04/2021",
			},
			repoError: nil,
			expectedUser: User{
				ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
				FirstName:    "Terrence",
				LastName:     "Trillow",
				Email:        "ttrillow1@feedburner.com",
				Password:     "5YLItbmdkfC1",
				IPAddress:    "63.119.6.98",
				CreationDate: "19/04/2021",
			},
			expectedError: nil,
		},
		{
			name:          "repo error",
			userID:        "unknown_id",
			repoResponse:  User{},
			repoError:     ErrNotFound,
			expectedUser:  User{},
			expectedError: ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersRepo := new(mockUsersRepo)

			ctx := context.Background()

			mockUsersRepo.On("GetUser", ctx, tc.userID).Return(tc.repoResponse, tc.repoError)

			svc := NewUsersService(mockUsersRepo)

			user, err := svc.GetUser(ctx, tc.userID)

			mockUsersRepo.AssertExpectations(t)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
