package srv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func getURLPathParam(urlPath, group string) (string, error) {
	urlPathParts := strings.Split(urlPath, group)
	if len(urlPathParts) < 2 {
		return "", fmt.Errorf("invalid group %s for url path %s", group, urlPath)
	}

	groupPath := urlPathParts[1]

	groupPathParts := strings.Split(groupPath, "/")
	if len(groupPathParts) < 2 {
		return "", fmt.Errorf("invalid url path %s", urlPath)
	}

	return groupPathParts[1], nil
}

func getURLQueryParam(urlValues url.Values, key string) string {
	values, ok := urlValues[key]

	if !ok || len(values[0]) < 1 {
		return ""
	}

	return values[0]
}

func getAndValidatePaginationParams(req *http.Request) (limit, offset int, err error) {
	urlQuery := req.URL.Query()

	limitStr := getURLQueryParam(urlQuery, "limit")
	if limitStr == "" { // required param
		return 0, 0, &httpError{
			StatusCode: http.StatusBadRequest,
			Message:    "missing required query param 'limit'",
		}
	}
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return 0, 0, &httpError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid integer param 'limit'",
		}
	}
	if limit > maxNumberOfUsersInResponse {
		return 0, 0, &httpError{
			StatusCode: http.StatusPreconditionFailed,
			Message:    fmt.Sprintf("'limit' is greater than %d", maxNumberOfUsersInResponse),
		}
	}

	offsetStr := getURLQueryParam(urlQuery, "offset")
	if offsetStr != "" { // optional param
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

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}
