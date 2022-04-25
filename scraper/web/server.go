package web

import (
	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
	"log"
	"net/http"
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
		addr = ":8080"
	}

	log.Printf("Running server at %s", addr)

	if err := http.ListenAndServe(addr, applyMiddlewares(mux)); err != nil {
		log.Fatalf("Unable to serve: %v", err)
	}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/pages", func(w http.ResponseWriter, req *http.Request) {
		HandlePages(s.service, w, req)

	})

	mux.HandleFunc("/api/v1/categories", func(w http.ResponseWriter, req *http.Request) {
		HandleCategories(s.service, w, req)
	})

	mux.HandleFunc("/api/v1/content", func(w http.ResponseWriter, req *http.Request) {

	})

	mux.HandleFunc("/api/v1/health", HandleHealthStatus)

	return mux
}
