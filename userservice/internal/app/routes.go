package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http/pprof"
)

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	// group for gathering metrics and requests counting,
	// doesn't include the /metrics, /debug/pprof, and /swagger  endpoints
	r.Group(func(r chi.Router) {

		r.Use(a.ZapLogger)
		r.Use(a.RequestCounter)
		r.Use(a.RequestLatency)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", a.controller.GetUsers)
			r.Get("/{email}", a.controller.GetUser)
			r.Post("/", a.controller.CreateUser)
		})
	})

	r.Route("/debug/pprof/", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/{cmd}", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
	})

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	return r
}
