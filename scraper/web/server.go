package web

import (
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
	middlewares "github.com/pestanko/miniscrape/scraper/web/middlewares"
	"github.com/rs/zerolog/log"
)

// Server representation
type Server struct {
	cfg     *config.AppConfig
	service *scraper.Service
}

// NewServer instance
func NewServer(cfg *config.AppConfig) Server {
	return Server{
		cfg:     cfg,
		service: scraper.NewService(cfg),
	}
}

// Serve the server
func (s *Server) Serve() {
	mux := s.routes()

	addr := s.cfg.Web.Addr
	if addr == "" {
		addr = "127.0.01:8080"
	}

	log.Info().
		Str("addr", addr).
		Msg("Running server")

	mds := []middlewares.Middleware{
		func(handler http.Handler, _ *config.AppConfig) http.Handler {
			return middlewares.RealIP(handler)
		},
		middlewares.RequestLogger,
		middlewares.VisitorCookie,
		middlewares.SetupCors,
	}

	if err := http.ListenAndServe(
		addr,
		middlewares.ApplyMiddlewares(mux, s.cfg, mds),
	); err != nil {
		log.Fatal().Err(err).Msg("Unable to serve")
	}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, req *http.Request) {
		HandleCategories(s.service, w, req)
	})

	mux.HandleFunc("/api/v1/pages", func(w http.ResponseWriter, req *http.Request) {
		HandlePages(s.service, w, req)
	})

	mux.HandleFunc("/api/v1/content", func(w http.ResponseWriter, req *http.Request) {
		HandlePagesContent(s.service, w, req)
	})

	mux.HandleFunc("/api/v1/health", HandleHealthStatus)

	// Auth related

	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, req *http.Request) {
		// sample way how to deal with annotations
		requireHTTPMethod(w, req, []string{http.MethodPost}, func() {
			HandleAuthLogin(s.service, w, req)
		})
	})

	mux.HandleFunc("/api/v1/auth/logout", func(w http.ResponseWriter, req *http.Request) {
		HandleAuthLogout(s.service, w, req)
	})

	mux.HandleFunc("/api/v1/auth/sessionstatus", func(w http.ResponseWriter, req *http.Request) {
		HandleSessionStatus(s.service, w, req)
	})

	// Auth Required

	mux.HandleFunc("/api/v1/cache", func(w http.ResponseWriter, req *http.Request) {
		requireAuthentication(s.service, w, req, func() {
			HandleCacheInvalidation(s.service, w, req)
		})
	})

	return mux
}
