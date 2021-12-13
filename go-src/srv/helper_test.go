package srv

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLPathParam(t *testing.T) {
	testCases := []struct {
		name          string
		urlPath       string
		group         string
		expectedValue string
		expectedError error
	}{
		{
			name:          "base case",
			urlPath:       "/v1/users/144bf891-f161-4c9a-8d83-38a275e088a5",
			group:         "users",
			expectedValue: "144bf891-f161-4c9a-8d83-38a275e088a5",
			expectedError: nil,
		},
		{
			name:          "base case 2",
			urlPath:       "/v1/users/144bf891-f161-4c9a-8d83-38a275e088a5/test",
			group:         "users",
			expectedValue: "144bf891-f161-4c9a-8d83-38a275e088a5",
			expectedError: nil,
		},
		{
			name:          "invalid group",
			urlPath:       "/v1/users/144bf891-f161-4c9a-8d83-38a275e088a5",
			group:         "unknown",
			expectedValue: "",
			expectedError: fmt.Errorf("invalid group unknown for url path /v1/users/144bf891-f161-4c9a-8d83-38a275e088a5"),
		},
		{
			name:          "invalid url path",
			urlPath:       "/v1/users144bf891-f161-4c9a-8d83-38a275e088a5",
			group:         "users",
			expectedValue: "",
			expectedError: fmt.Errorf("invalid url path /v1/users144bf891-f161-4c9a-8d83-38a275e088a5"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := getURLPathParam(tc.urlPath, tc.group)

			assert.Equal(t, tc.expectedValue, value)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetURLQueryParam(t *testing.T) {
	testCases := []struct {
		name          string
		urlValues     url.Values
		key           string
		expectedValue string
	}{
		{
			name: "base case 1",
			urlValues: url.Values{
				"limit":  []string{"10"},
				"offset": []string{"5"},
			},
			key:           "limit",
			expectedValue: "10",
		},
		{
			name: "base case 2",
			urlValues: url.Values{
				"limit":  []string{"10", "1"},
				"offset": []string{"5", "8"},
			},
			key:           "offset",
			expectedValue: "5",
		},
		{
			name: "no key found 1",
			urlValues: url.Values{
				"limit":  []string{"10"},
				"offset": []string{"5"},
			},
			key:           "unkown",
			expectedValue: "",
		},
		{
			name:          "no key found 2",
			urlValues:     url.Values{},
			key:           "limit",
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := getURLQueryParam(tc.urlValues, tc.key)

			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestGetAndValidatePaginationParams(t *testing.T) {
	testCases := []struct {
		name           string
		urlValues      url.Values
		maxLimit       int
		expectedLimit  int
		expectedOffset int
		expectedError  error
	}{
		{
			name: "base case",
			urlValues: url.Values{
				"limit":  []string{"10"},
				"offset": []string{"5"},
			},
			maxLimit:       10,
			expectedLimit:  10,
			expectedOffset: 5,
			expectedError:  nil,
		},
		{
			name: "no offset - defaulting to 0",
			urlValues: url.Values{
				"limit": []string{"10"},
			},
			maxLimit:       20,
			expectedLimit:  10,
			expectedOffset: 0,
			expectedError:  nil,
		},
		{
			name:           "error - no limit query",
			urlValues:      url.Values{},
			maxLimit:       1000,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "missing required query param 'limit'",
			},
		},
		{
			name: "error - invalid limit 1",
			urlValues: url.Values{
				"limit": []string{"a"},
			},
			maxLimit:       1000,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid integer param 'limit'",
			},
		},
		{
			name: "error - invalid limit 2",
			urlValues: url.Values{
				"limit": []string{"-1"},
			},
			maxLimit:       1000,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid integer param 'limit'",
			},
		},
		{
			name: "error - limit above max permitted",
			urlValues: url.Values{
				"limit": []string{"1001"},
			},
			maxLimit:       1000,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusPreconditionFailed,
				Message:    "'limit' is greater than 1000",
			},
		},
		{
			name: "error - invalid offset 1",
			urlValues: url.Values{
				"limit":  []string{"5"},
				"offset": []string{"a"},
			},
			maxLimit:       10,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid integer param 'offset'",
			},
		},
		{
			name: "error - invalid offset 2",
			urlValues: url.Values{
				"limit":  []string{"5"},
				"offset": []string{"-1"},
			},
			maxLimit:       10,
			expectedLimit:  0,
			expectedOffset: 0,
			expectedError: &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid integer param 'offset'",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			limit, offset, err := getAndValidatePaginationParams(tc.urlValues, tc.maxLimit)

			assert.Equal(t, tc.expectedLimit, limit)
			assert.Equal(t, tc.expectedOffset, offset)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
