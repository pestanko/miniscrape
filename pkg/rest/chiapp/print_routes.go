package chiapp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// LogChiRoutes logs all registered routes in the chi.Mux
func LogChiRoutes(app chi.Routes) {
	if err := chi.Walk(app, func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		log.Info().Fields(map[string]any{
			"method": method,
			"route":  route,
		}).Msg("registered route")
		return nil
	}); err != nil {
		log.Warn().Err(err).Msg("error logging chi routes")
	}
}
