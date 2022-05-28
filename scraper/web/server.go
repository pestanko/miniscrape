package web

import (
	"log"
	"net/http"

	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
	middlewares "github.com/pestanko/miniscrape/scraper/web/middlewares"
)

type Server struct {
	cfg     *config.AppConfig
	service *scraper.Service
}

func MakeServer(cfg *config.AppConfig) Server {
	return Server{
		cfg:     cfg,
		service: scraper.NewService(cfg),
	}
}

func (s *Server) Serve() {
	mux := s.routes()

	addr := s.cfg.Web.Addr
	if addr == "" {
		addr = "127.0.01:8080"
	}

	log.Printf("Running server at %s", addr)

	mds := []middlewares.Middleware{
		middlewares.RequestLogger,
		middlewares.VisitorCookie,
		middlewares.SetupCors,
	}

	if err := http.ListenAndServe(addr, middlewares.ApplyMiddlewares(mux, mds)); err != nil {
		log.Fatalf("Unable to serve: %v", err)
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

	return mux
}
