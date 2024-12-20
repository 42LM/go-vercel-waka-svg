package service_test

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"go-vercel-waka-svg/internal/service"
	"go-vercel-waka-svg/internal/svgtemplate"
)

func Test_Service_Error(t *testing.T) {
	testCases := map[string]struct {
		queryParams map[string]string

		expBody string
		expErr  error
	}{
		"all valid params": {
			expBody: error_svg,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// take a look at the general URL form:
			// https://pkg.go.dev/net/url#:~:text=%5Bscheme%3A%5D%5B//%5Buserinfo%40%5Dhost%5D%5B/%5Dpath%5B%3Fquery%5D%5B%23fragment%5D
			svgTemplates, _ := svgtemplate.GetSVGTemplates()

			rec := httptest.NewRecorder()
			s := service.New(&service.ServiceConfig{
				Logger:         slog.Default(),
				QueryParams:    nil,
				ResponseWriter: rec,
				Templates:      svgTemplates,
			})

			err := s.Error(context.Background())

			// Get the captured response from the ResponseRecorder
			resp := rec.Result()

			// Read the response body
			body, _ := io.ReadAll(resp.Body)

			if tc.expBody != string(body) {
				t.Errorf("\nexpected: %v\nactual: %v\n", tc.expBody, string(body))
			}

			if tc.expErr != err {
				t.Errorf("\nexpected: %v\nactual: %v\n", tc.expErr, err)
			}
		})
	}
}

const (
	error_svg = "<svg width=\"100\" height=\"100\" viewBox=\"0 0 1397 1296\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n  <title>error</title>\n  <g clip-path=\"url(#clip0_25_7)\">\n    <path d=\"M698.5 166L1207.29 1047.25H189.71L698.5 166Z\" fill=\"#F31559\"/>\n    <path d=\"M645 355V802H753V355H645Z\" fill=\"#D9D9D9\"/>\n    <path d=\"M645 882V990H753V882H645Z\" fill=\"#D9D9D9\"/>\n  </g>\n  <defs>\n    <clipPath id=\"clip0_25_7\">\n      <rect width=\"1397\" height=\"1296\" fill=\"white\"/>\n    </clipPath>\n  </defs>\n</svg>\n"
)
