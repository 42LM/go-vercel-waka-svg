// Package service implements the actual service functions that deliver svg files
package service

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
)

// Service describes the Service interface.
type Service interface {
	Error(ctx context.Context) error
	Wakatime(ctx context.Context) error
}

// Check if service implements Service interface explicitly.
var _ Service = (*service)(nil)

// service implements the Service interface.
type service struct {
	queryParams    map[string]string
	responseWriter http.ResponseWriter
	templates      *template.Template
}

// ServiceConfig contains the configuration params of the service.
type ServiceConfig struct {
	Logger         *slog.Logger
	QueryParams    map[string]string
	ResponseWriter http.ResponseWriter
	Templates      *template.Template
}

// New returns a service with middleware wired in.
func New(config *ServiceConfig) Service {
	var svc Service
	svc = &service{
		queryParams:    config.QueryParams,
		responseWriter: config.ResponseWriter,
		templates:      config.Templates,
	}
	svc = LoggingMiddleware(config.Logger)(svc)
	return svc
}
