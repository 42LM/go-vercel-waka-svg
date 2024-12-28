// Package query implements utility routines for extracting URL Query Parameters.
package query

import (
	"net/http"
	"strings"
)

// GetQueryParams retrieves all URL query parameters.
func GetQueryParams(r *http.Request) map[string]string {
	urlValues := r.URL.Query()
	res := make(map[string]string, len(urlValues))

	for k, v := range urlValues {
		res[k] = strings.Join(v, ",")
	}

	return res
}

// GetQueryParam retrieves a single URL query parameter from a request by given param name.
func GetQueryParam(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}
