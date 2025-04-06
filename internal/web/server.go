// Package web provides a web server for the miniscrape application.
package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/scraper"
	"github.com/pestanko/miniscrape/internal/web/handlers"
	"github.com/pestanko/miniscrape/internal/web/middlewares"
	"github.com/pestanko/miniscrape/pkg/rest/chiapp"
)

// NewServer creates a new chi multiplexer instance
func NewServer(cfg *config.AppConfig) *chi.Mux {
	service := scraper.NewService(cfg)

	app := chiapp.CreateChiApp(
		chiapp.WithServiceName(cfg.ServiceInfo.Name),
		chiapp.WithPublicHealthEndpoints("/api/health"),
		chiapp.WithPrometheus(true),
	)

	registerRoutes(app, service)

	return app
}

func registerRoutes(mux chi.Router, service *scraper.Service) {
	// Health check routes don't need to be instrumented
	registerHealthRoutes(mux)

	// Register API routes
	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/categories", handlers.HandleCategories(service))
		r.Get("/pages", handlers.HandlePages(service))
		r.Get("/content", handlers.HandlePagesContent(service))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handlers.HandleAuthLogin(service))
			r.Post("/logout", handlers.HandleAuthLogout(service))
			r.Get("/sessionstatus", handlers.HandleSessionStatus(service))
		})

		r.Route("/cache", func(r chi.Router) {
			r.Use(middlewares.AuthRequired(service))
			r.Post("/", handlers.HandleCacheInvalidation(service))
		})
	})

	chiapp.LogChiRoutes(mux)
}

func registerHealthRoutes(mux chi.Router) {
	// Register health routes
	mux.Get("/health/live", handlers.HandleHealthStatus())
	mux.Get("/health/ready", handlers.HandleHealthStatus())
	mux.Get("/health", handlers.HandleHealthStatus())
}
