// Package handler implements a single HTTP handler
// that is being used by the Go Runtime of Vercel.
package handler

import (
	"log/slog"
	"net/http"
	"os"

	"go-vercel-waka-svg/internal/query"
	"go-vercel-waka-svg/internal/service"
	"go-vercel-waka-svg/internal/svgtemplate"
)

// GenerateSVG is the single handler being triggered by vercel.
// It generates an svg according to the given `type`.
func GenerateSVG(w http.ResponseWriter, r *http.Request) {
	// fetch query parameter
	qp := query.GetQueryParams(r)
	// load templates
	svgTemplates, err := svgtemplate.GetSVGTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// set the headers content type to svg
	w.Header().Set("Content-Type", "image/svg+xml")
	// build service
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	svc := service.New(&service.ServiceConfig{
		Logger:         slogger.With("layer", "service"),
		QueryParams:    qp,
		ResponseWriter: w,
		Templates:      svgTemplates,
	})

	ctx := r.Context()

	// type switch for svg type
	// no type given leads to creation of a default error svg
	switch qp["type"] {
	case "waka":
		err := svc.Wakatime(ctx)
		if err != nil {
			return
		}
	default:
		err := svc.Error(ctx)
		if err != nil {
			return
		}
	}
}
