package srv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// getURLPathParam gets URL path parameter value based on the group (e.g. "users")
func getURLPathParam(urlPath, group string) (string, error) {
	// splitting the URL path by the group string to get the group path (second part)
	urlPathParts := strings.Split(urlPath, group)
	if len(urlPathParts) < 2 {
		return "", fmt.Errorf("invalid group %s for url path %s", group, urlPath)
	}

	groupPath := urlPathParts[1]

	// splitting the group path by "/" to get the parameter value (second part)
	groupPathParts := strings.Split(groupPath, "/")
	if len(groupPathParts) < 2 {
		return "", fmt.Errorf("invalid url path %s", urlPath)
	}

	return groupPathParts[1], nil
}

// getURLQueryParam gets querystring value from request URL
func getURLQueryParam(urlValues url.Values, key string) string {
	values, ok := urlValues[key]

	if !ok || len(values[0]) < 1 {
		return ""
	}

	return values[0]
}

// getAndValidatePaginationParams gets and validates the pagination parameters (limit and offset) from the URL querystrings
func getAndValidatePaginationParams(urlQuery url.Values, maxLimit int) (limit, offset int, err error) {
	// getting the required "limit" parameter
	limitStr := getURLQueryParam(urlQuery, "limit")
	if limitStr == "" { // required param
		return 0, 0, &httpError{
			StatusCode: http.StatusBadRequest,
			Message:    "missing required query param 'limit'",
		}
	}

	// limit must be positive integer
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return 0, 0, &httpError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid integer param 'limit'",
		}
	}
	// limit cannot be greater than the maximum limit
	if limit > maxLimit {
		return 0, 0, &httpError{
			StatusCode: http.StatusPreconditionFailed,
			Message:    fmt.Sprintf("'limit' is greater than %d", maxLimit),
		}
	}

	// getting the "offset" parameter (optional, defaulting to 0)
	offsetStr := getURLQueryParam(urlQuery, "offset")
	if offsetStr != "" { // optional param
		// offset must be positive integer
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, &httpError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid integer param 'offset'",
			}
		}
	}

	return limit, offset, nil
}

// writeJSON writes the correct header, status code and JSON format to the response
func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
