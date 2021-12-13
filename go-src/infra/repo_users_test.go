package infra

import (
	"context"
	"fmt"
	"testing"

	"github.com/hbernardo/users/go-src/lib"
	"github.com/stretchr/testify/assert"
)

var testUsersData = []lib.User{
	{
		ID:           "144bf891-f161-4c9a-8d83-38a275e088a5",
		FirstName:    "Nicky",
		LastName:     "Blasio",
		Email:        "nblasio0@jiathis.com",
		Password:     "rKJKin",
		IPAddress:    "43.113.46.36",
		CreationDate: "06/06/2021",
	},
	{
		ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
		FirstName:    "Terrence",
		LastName:     "Trillow",
		Email:        "ttrillow1@feedburner.com",
		Password:     "5YLItbmdkfC1",
		IPAddress:    "63.119.6.98",
		CreationDate: "19/04/2021",
	},
	{
		ID:           "3e601207-0e80-4e7e-ae87-bb802b16a179",
		FirstName:    "Niels",
		LastName:     "MacPaik",
		Email:        "nmacpaik2@phoca.cz",
		Password:     "Vae1mnI",
		IPAddress:    "94.47.183.190",
		CreationDate: "19/01/2021",
	},
}

func TestGetUsers(t *testing.T) {
	testCases := []struct {
		name          string
		usersData     []lib.User
		limit         int
		offset        int
		expectedUsers []lib.User
		expectedError error
	}{
		{
			name:          "base case - all data",
			usersData:     testUsersData,
			limit:         3,
			offset:        0,
			expectedUsers: testUsersData,
			expectedError: nil,
		},
		{
			name:      "first 2 users",
			usersData: testUsersData,
			limit:     2,
			offset:    0,
			expectedUsers: []lib.User{
				{
					ID:           "144bf891-f161-4c9a-8d83-38a275e088a5",
					FirstName:    "Nicky",
					LastName:     "Blasio",
					Email:        "nblasio0@jiathis.com",
					Password:     "rKJKin",
					IPAddress:    "43.113.46.36",
					CreationDate: "06/06/2021",
				},
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
			name:      "last 2 users",
			usersData: testUsersData,
			limit:     2,
			offset:    1,
			expectedUsers: []lib.User{
				{
					ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
					FirstName:    "Terrence",
					LastName:     "Trillow",
					Email:        "ttrillow1@feedburner.com",
					Password:     "5YLItbmdkfC1",
					IPAddress:    "63.119.6.98",
					CreationDate: "19/04/2021",
				},
				{
					ID:           "3e601207-0e80-4e7e-ae87-bb802b16a179",
					FirstName:    "Niels",
					LastName:     "MacPaik",
					Email:        "nmacpaik2@phoca.cz",
					Password:     "Vae1mnI",
					IPAddress:    "94.47.183.190",
					CreationDate: "19/01/2021",
				},
			},
			expectedError: nil,
		},
		{
			name:      "limiting the limit",
			usersData: testUsersData,
			limit:     900,
			offset:    2,
			expectedUsers: []lib.User{
				{
					ID:           "3e601207-0e80-4e7e-ae87-bb802b16a179",
					FirstName:    "Niels",
					LastName:     "MacPaik",
					Email:        "nmacpaik2@phoca.cz",
					Password:     "Vae1mnI",
					IPAddress:    "94.47.183.190",
					CreationDate: "19/01/2021",
				},
			},
			expectedError: nil,
		},
		{
			name:          "offset greater than the data size",
			usersData:     testUsersData,
			limit:         10,
			offset:        4,
			expectedUsers: []lib.User{},
			expectedError: nil,
		},
		{
			name:          "error - invalid limit",
			usersData:     testUsersData,
			limit:         -2,
			offset:        0,
			expectedUsers: nil,
			expectedError: fmt.Errorf("'limit' nor 'offset' cannot be negative: %w", lib.ErrPreconditionFailed),
		},
		{
			name:          "error - invalid offset",
			usersData:     testUsersData,
			limit:         1,
			offset:        -2,
			expectedUsers: nil,
			expectedError: fmt.Errorf("'limit' nor 'offset' cannot be negative: %w", lib.ErrPreconditionFailed),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewUsersRepo(tc.usersData)

			users, err := repo.GetUsers(context.Background(), tc.limit, tc.offset)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUsers, users)
		})
	}
}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name          string
		usersData     []lib.User
		userID        string
		expectedUser  lib.User
		expectedError error
	}{
		{
			name:      "base case",
			usersData: testUsersData,
			userID:    "1311f914-1d4f-40b6-8886-80193265d5a4",
			expectedUser: lib.User{
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
			name:          "not found",
			usersData:     testUsersData,
			userID:        "unknown_id",
			expectedUser:  lib.User{},
			expectedError: lib.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewUsersRepo(tc.usersData)

			user, err := repo.GetUser(context.Background(), tc.userID)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
