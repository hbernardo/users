package srv

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/hbernardo/users/go-src/lib"
	"github.com/stretchr/testify/mock"
)

type mockUsersService struct {
	mock.Mock
}

func (m *mockUsersService) GetUsers(ctx context.Context, limit int, offset int) ([]lib.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]lib.User), args.Error(1)
}

func (m *mockUsersService) GetUser(ctx context.Context, userID string) (lib.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(lib.User), args.Error(1)
}

type mockHTTPResponseWriter struct {
	mock.Mock
}

func (m *mockHTTPResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *mockHTTPResponseWriter) WriteHeader(status int) {
	m.Called(status)
}

func (m *mockHTTPResponseWriter) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Get(0).(int), args.Error(1)
}

func TestHandleGetUsers(t *testing.T) {
	// NOTE: function "getAndValidatePaginationParams" is already being tested in "helper_test.go"
	// and function "handleError" is already being tested in "errors_test.go"
	// so all tests here will assume the success case scenario for them

	testCases := []struct {
		name               string
		httpRequest        *http.Request
		svcNotCalled       bool
		svcResponse        []lib.User
		svcError           error
		expectedLimit      int
		expectedOffset     int
		expectedHTTPStatus int
		expectedResponse   []byte
	}{
		{
			name: "base case",
			httpRequest: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path:     "/v1/users",
					RawQuery: "limit=1&offset=5",
				},
			},
			svcResponse: []lib.User{
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
			svcError:           nil,
			expectedLimit:      1,
			expectedOffset:     5,
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   []byte(`[{"id":"1311f914-1d4f-40b6-8886-80193265d5a4","first_name":"Terrence","last_name":"Trillow","email":"ttrillow1@feedburner.com","password":"5YLItbmdkfC1","ip_address":"63.119.6.98","creation_date":"19/04/2021"}]` + "\n"),
		},
		{
			name: "service error",
			httpRequest: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path:     "/v1/users",
					RawQuery: "limit=1&offset=5",
				},
			},
			svcResponse:        nil,
			svcError:           fmt.Errorf("svc error"),
			expectedLimit:      1,
			expectedOffset:     5,
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedResponse:   []byte(`{"error":"internal server error"}` + "\n"),
		},
		{
			name: "not allowed method",
			httpRequest: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path:     "/v1/users",
					RawQuery: "limit=1&offset=5",
				},
			},
			svcNotCalled:       true,
			expectedHTTPStatus: http.StatusMethodNotAllowed,
			expectedResponse:   []byte(`{"error":"method not allowed"}` + "\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := new(mockUsersService)
			mockUsersService.On("GetUsers", mock.Anything, tc.expectedLimit, tc.expectedOffset).Return(tc.svcResponse, tc.svcError)

			mockHTTPResponseWriter := new(mockHTTPResponseWriter)
			mockHTTPResponseWriter.On("Header").Return(make(http.Header))
			mockHTTPResponseWriter.On("WriteHeader", tc.expectedHTTPStatus)
			mockHTTPResponseWriter.On("Write", tc.expectedResponse).Return(len(tc.expectedResponse), nil)

			handler := NewUsersHandler(mockUsersService)
			handler.handleGetUsers(mockHTTPResponseWriter, tc.httpRequest)

			if tc.svcNotCalled == false {
				mockUsersService.AssertExpectations(t)
			}
			mockHTTPResponseWriter.AssertExpectations(t)
		})
	}
}

func TestHandleGetUser(t *testing.T) {
	// NOTE: function "getURLPathParam" is already being tested in "helper_test.go"
	// and function "handleError" is already being tested in "errors_test.go"
	// so all tests here will assume the success case scenario for them

	testCases := []struct {
		name               string
		httpRequest        *http.Request
		svcNotCalled       bool
		svcResponse        lib.User
		svcError           error
		expectedUserID     string
		expectedHTTPStatus int
		expectedResponse   []byte
	}{
		{
			name: "base case",
			httpRequest: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/v1/users/1311f914-1d4f-40b6-8886-80193265d5a4",
				},
			},
			svcResponse: lib.User{
				ID:           "1311f914-1d4f-40b6-8886-80193265d5a4",
				FirstName:    "Terrence",
				LastName:     "Trillow",
				Email:        "ttrillow1@feedburner.com",
				Password:     "5YLItbmdkfC1",
				IPAddress:    "63.119.6.98",
				CreationDate: "19/04/2021",
			},
			svcError:           nil,
			expectedUserID:     "1311f914-1d4f-40b6-8886-80193265d5a4",
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   []byte(`{"id":"1311f914-1d4f-40b6-8886-80193265d5a4","first_name":"Terrence","last_name":"Trillow","email":"ttrillow1@feedburner.com","password":"5YLItbmdkfC1","ip_address":"63.119.6.98","creation_date":"19/04/2021"}` + "\n"),
		},
		{
			name: "service error",
			httpRequest: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/v1/users/1311f914-1d4f-40b6-8886-80193265d5a4",
				},
			},
			svcResponse:        lib.User{},
			svcError:           lib.ErrNotFound,
			expectedUserID:     "1311f914-1d4f-40b6-8886-80193265d5a4",
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   []byte(`{"error":"not found"}` + "\n"),
		},
		{
			name: "not allowed method",
			httpRequest: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/v1/users/1311f914-1d4f-40b6-8886-80193265d5a4",
				},
			},
			svcNotCalled:       true,
			expectedHTTPStatus: http.StatusMethodNotAllowed,
			expectedResponse:   []byte(`{"error":"method not allowed"}` + "\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := new(mockUsersService)
			mockUsersService.On("GetUser", mock.Anything, tc.expectedUserID).Return(tc.svcResponse, tc.svcError)

			mockHTTPResponseWriter := new(mockHTTPResponseWriter)
			mockHTTPResponseWriter.On("Header").Return(make(http.Header))
			mockHTTPResponseWriter.On("WriteHeader", tc.expectedHTTPStatus)
			mockHTTPResponseWriter.On("Write", tc.expectedResponse).Return(len(tc.expectedResponse), nil)

			handler := NewUsersHandler(mockUsersService)
			handler.handleGetUser(mockHTTPResponseWriter, tc.httpRequest)

			if tc.svcNotCalled == false {
				mockUsersService.AssertExpectations(t)
			}
			mockHTTPResponseWriter.AssertExpectations(t)
		})
	}
}
