package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	registry := prometheus.NewRegistry()

	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Duration of HTTP requests in seconds.",
		Buckets:   prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})

	registry.MustRegister(histogram)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(startTime).Seconds()
		histogram.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
