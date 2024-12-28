package query_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"go-vercel-waka-svg/nointeral/query"
)

func Test_GetQueryParams(t *testing.T) {
	testCases := map[string]struct {
		queryParams map[string]string

		expQueryParams map[string]string
	}{
		"all valid params": {
			queryParams: map[string]string{
				"color":       "red",
				"font_family": "arial",
				"font_size":   "20",
				"text":        "hello%20world",
			},
			expQueryParams: map[string]string{
				"color":       "red",
				"font_family": "arial",
				"font_size":   "20",
				"text":        "hello world",
			},
		},
		"mixed valid and invalid params - still get all params": {
			queryParams: map[string]string{"color": "blue", "a": "r2d2", "b": "20", "c": "hello%20world"},
			expQueryParams: map[string]string{
				"color": "blue",
				"a":     "r2d2",
				"b":     "20",
				"c":     "hello world",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// take a look at the general URL form:
			// https://pkg.go.dev/net/url#:~:text=%5Bscheme%3A%5D%5B//%5Buserinfo%40%5Dhost%5D%5B/%5Dpath%5B%3Fquery%5D%5B%23fragment%5D
			querySlice := make([]string, 0, len(tc.queryParams))
			for k, v := range tc.queryParams {
				querySlice = append(querySlice, k+"="+v)
			}

			q := "//?" + strings.Join(querySlice, "&")

			r := httptest.NewRequest("GET", q, nil)

			params := query.GetQueryParams(r)

			if !mapsEqual(tc.expQueryParams, params) {
				t.Errorf("\nexpected: %v\nactual: %v\n", tc.expQueryParams, params)
			}
		})
	}
}

func Test_GetQueryParam(t *testing.T) {
	testCases := map[string]struct {
		queryParamName string
		queryParams    map[string]string

		expQueryParam string
	}{
		"param exists in req": {
			queryParams:    map[string]string{"color": "123456"},
			queryParamName: "color",
			expQueryParam:  "123456",
		},
		"param does not exist in req": {
			queryParams:    map[string]string{},
			queryParamName: "color",
			expQueryParam:  "",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// take a look at the general URL form:
			// https://pkg.go.dev/net/url#:~:text=%5Bscheme%3A%5D%5B//%5Buserinfo%40%5Dhost%5D%5B/%5Dpath%5B%3Fquery%5D%5B%23fragment%5D
			querySlice := make([]string, 0, len(tc.queryParams))
			for k, v := range tc.queryParams {
				querySlice = append(querySlice, k+"="+v)
			}

			q := "//?" + strings.Join(querySlice, "&")

			r := httptest.NewRequest("GET", q, nil)

			param := query.GetQueryParam(r, tc.queryParamName)

			if tc.expQueryParam != param {
				t.Errorf("\nexpected: %v\nactual: %v\n", tc.expQueryParam, param)
			}
		})
	}
}

func mapsEqual(map1, map2 map[string]string) bool {
	// Check if the maps have the same length
	if len(map1) != len(map2) {
		return false
	}

	// Iterate over the key-value pairs of map1
	for key, value1 := range map1 {
		// Check if the key exists in map2
		value2, ok := map2[key]
		if !ok {
			return false
		}

		// Check if the values are equal
		if value1 != value2 {
			return false
		}
	}

	return true
}
