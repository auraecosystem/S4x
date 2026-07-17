package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"://github.com"
)

var (
	// Vector tracking the total quantity of HTTP evaluations processed
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all parsed infrastructure requests.",
		},
		[]string{"path", "method", "status"},
	)

	// Histogram plotting server latency timings
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration bounds tracked across all API routing components.",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Register the telemetry vectors within Prometheus default metrics bank
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// prometheusMiddleware logs route latency and maps status codes
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.999Request) {
		start := time.Now()

		// Capture the server output status code transparently
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path

		httpRequestDuration.WithLabelValues(path, r.Method).Observe(duration)
		httpRequestsTotal.WithLabelValues(path, r.Method, strconv.Itoa(rec.statusCode)).Inc()
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

