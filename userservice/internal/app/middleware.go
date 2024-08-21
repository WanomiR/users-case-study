package app

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "total_requests_count",
	Help: "Number of requests received.",
}, []string{"type"})

var requestLatencyHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_latency_seconds",
	Help:    "Request latency in seconds",
	Buckets: []float64{0.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
}, []string{"method"})

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestLatencyHist)
}

func (a *App) ZapLogger(next http.Handler) http.Handler {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger := zap.Must(cfg.Build())
	defer logger.Sync()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logger.Info("new request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("addr", r.RemoteAddr),
		)
	})
}

func (a *App) RequestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		requestsTotal.With(prometheus.Labels{"type": "requests"}).Inc()
	})
}

func (a *App) RequestLatency(next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		requestLatencyHist.With(prometheus.Labels{"method": r.Method}).Observe(time.Since(start).Seconds())
	})
}
