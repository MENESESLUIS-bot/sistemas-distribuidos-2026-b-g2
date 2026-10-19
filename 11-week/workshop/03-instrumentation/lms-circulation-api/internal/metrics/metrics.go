// Package metrics exposes the Prometheus metrics of lms-circulation-api used
// by 04-metrics.md and the SLOs of 05-slo.md.
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// LatencyBuckets follow NFR-001 (p95 < 300 ms, p99 < 600 ms) and end at the
// 5 s client timeout used towards lms-catalog-api and lms-membership-api.
var LatencyBuckets = []float64{0.05, 0.1, 0.2, 0.3, 0.5, 0.6, 1, 2.5, 5}

type Metrics struct {
	registry   *prometheus.Registry
	requests   *prometheus.CounterVec
	duration   *prometheus.HistogramVec
	downstream *prometheus.HistogramVec
}

func New() *Metrics {
	m := &Metrics{
		registry: prometheus.NewRegistry(),
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "lms_circulation_loan_requests_total",
			Help: "HU-06 loan requests by route and outcome (success, rejected, client_error, error).",
		}, []string{"route", "outcome"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "lms_circulation_loan_request_duration_seconds",
			Help:    "Server-side latency of HU-06 loan requests.",
			Buckets: LatencyBuckets,
		}, []string{"route"}),
		downstream: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "lms_circulation_downstream_request_duration_seconds",
			Help:    "Latency of calls from lms-circulation-api to other LMS services.",
			Buckets: LatencyBuckets,
		}, []string{"target", "operation", "outcome"}),
	}
	m.registry.MustRegister(m.requests, m.duration, m.downstream)
	return m
}

// ObserveRequest records one finished request of the given route template.
func (m *Metrics) ObserveRequest(route, outcome string, d time.Duration) {
	m.requests.WithLabelValues(route, outcome).Inc()
	m.duration.WithLabelValues(route).Observe(d.Seconds())
}

// CatalogObserver records calls to lms-catalog-api.
func (m *Metrics) CatalogObserver(operation, outcome string, d time.Duration) {
	m.downstream.WithLabelValues("lms-catalog-api", operation, outcome).Observe(d.Seconds())
}

// Handler serves the registry in the Prometheus text format on /metrics.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
