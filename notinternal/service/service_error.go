// Package service implements the actual service functions that deliver svg files
package service

import (
	"context"
	"time"
)

// Error returns a default error svg.
// Templates can be found in `../svgtemplate/templates/...`
func (s *service) Error(ctx context.Context) error {
	return s.templates.ExecuteTemplate(s.responseWriter, "error.gosvg", nil)
}

func (mw loggingMiddleware) Error(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service invocation",
			"method", "Error",
			"result.err", err,
			"took", (time.Since(begin) / 1e6).String(),
		)
	}(time.Now())
	return mw.next.Error(ctx)
}
