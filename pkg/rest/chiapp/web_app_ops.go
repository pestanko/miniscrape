package chiapp

import (
	appmidlewares "github.com/pestanko/miniscrape/internal/web/middlewares"
	"github.com/pestanko/miniscrape/pkg/utils/applog"
	"github.com/riandyrn/otelchi"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pestanko/miniscrape/pkg/utils/collut"
	"github.com/rs/zerolog"
)

// envPrometheusEnabled - default is false
const envPrometheusEnabled = "PROMETHEUS_ENABLED"
const defaultReadTimeout = 10 * time.Second
const defaultReadBufSize = 16384

// AppOpsFn represents a modifier function for the AppOps
type AppOpsFn = func(op *AppOps)

// AppOps represents a fiber application options
type AppOps struct {
	Name                  string
	PrometheusEnabled     bool
	PrometheusSetupFunc   func(app chi.Router, appOps AppOps)
	PublicHealthEndpoints []string
	HealthFunc            http.HandlerFunc
	ReadTimeout           time.Duration
	ReadBufferSize        int
	DefaultMiddlewares    func(app chi.Router)
}

// WithServiceName override the service/app name default is empty
func WithServiceName(name string) AppOpsFn {
	return func(ops *AppOps) {
		ops.Name = name
	}
}

// WithPrometheus set prometheus either enabled or disabled where true = enabled
// Overrides default value that is based on the PROMETHEUS_ENABLED env variable
func WithPrometheus(enabled bool) AppOpsFn {
	return func(ops *AppOps) {
		ops.PrometheusEnabled = enabled
	}
}

// WithDefaultMiddlewares registers a default middlewares
func WithDefaultMiddlewares(middlewareFn func(r chi.Router)) AppOpsFn {
	return func(ops *AppOps) {
		ops.DefaultMiddlewares = middlewareFn
	}
}

// WithPublicHealthEndpoints add list of public health endpoints to be setup
func WithPublicHealthEndpoints(endpoints ...string) AppOpsFn {
	return func(ops *AppOps) {
		ops.PublicHealthEndpoints = append(ops.PublicHealthEndpoints, endpoints...)
	}
}

// CreateChiApp create a chi application
func CreateChiApp(ops ...AppOpsFn) *chi.Mux {
	isPrometheusEnabled, _ := strconv.ParseBool(os.Getenv(envPrometheusEnabled))

	appOps := AppOps{
		Name:                  "",
		ReadTimeout:           defaultReadTimeout,
		ReadBufferSize:        defaultReadBufSize,
		PrometheusEnabled:     isPrometheusEnabled,
		PrometheusSetupFunc:   nil,
		PublicHealthEndpoints: []string{},
		DefaultMiddlewares:    defaultMiddlewares,
	}

	appOps = collut.OpsApplyAll(appOps, ops...)

	app := chi.NewRouter()

	if appOps.PrometheusEnabled && appOps.PrometheusSetupFunc != nil {
		appOps.PrometheusSetupFunc(app, appOps)
	}

	if appOps.DefaultMiddlewares != nil {
		appOps.DefaultMiddlewares(app)
	}

	return app
}

func defaultMiddlewares(r chi.Router) {
	r.Use(otelchi.Middleware("miniscrape", otelchi.WithChiRoutes(r)))
	r.Use(appmidlewares.RealIP())
	r.Use(middleware.RequestID)
	r.Use(appmidlewares.SetupCors())
	r.Use(appmidlewares.VisitorCookie())
	r.Use(appmidlewares.Logger(appmidlewares.LogParams{
		LogCfg: applog.LogConfig{},
		Log:    zerolog.Logger{},
	}))
}
