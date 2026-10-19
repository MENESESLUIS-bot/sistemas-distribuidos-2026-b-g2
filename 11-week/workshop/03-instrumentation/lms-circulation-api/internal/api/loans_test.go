package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/catalogclient"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/logging"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/metrics"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

const (
	traceID     = "7c3e9a41d2b84f06a1e5c0b9d8f27e13"
	traceparent = "00-" + traceID + "-a01f5c3e9b72d846-01"
	body        = `{"studentId":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","bookId":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82"}`
)

// fakeCatalog stands in for lms-catalog-api and remembers the headers it got.
func fakeCatalog(t *testing.T, status int, delay time.Duration, got *http.Header) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/loan-copy") {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		*got = r.Header.Clone()
		time.Sleep(delay)
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func newServer(catalogURL string, timeout time.Duration, logs io.Writer) *Server {
	m := metrics.New()
	return &Server{
		Catalog: catalogclient.New(catalogURL, timeout, m.CatalogObserver),
		Log:     logging.New(logs, "lms-circulation-api", "2.0.1", "test", slog.LevelDebug),
		Metrics: m,
	}
}

func post(t *testing.T, s *Server) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", strings.NewReader(body))
	req.Header.Set(tracecontext.Header, traceparent)
	req.Header.Set("X-Correlation-Id", "5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5")
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	return rec
}

func TestLoanCopyReceivesSameTraceID(t *testing.T) {
	var got http.Header
	catalog := fakeCatalog(t, http.StatusOK, 0, &got)
	var logs bytes.Buffer
	s := newServer(catalog.URL, time.Second, &logs)

	if rec := post(t, s); rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}

	down, ok := tracecontext.Parse(got.Get(tracecontext.Header))
	if !ok {
		t.Fatalf("lms-catalog-api got no valid traceparent: %q", got.Get(tracecontext.Header))
	}
	if down.TraceID != traceID {
		t.Fatalf("downstream trace_id = %s, want %s", down.TraceID, traceID)
	}
	if down.SpanID == "a01f5c3e9b72d846" {
		t.Fatal("downstream span must be a new client span, not the caller's span")
	}
	if got.Get("X-Correlation-Id") != "5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5" {
		t.Fatal("X-Correlation-Id was not forwarded")
	}

	for _, line := range strings.Split(strings.TrimSpace(logs.String()), "\n") {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("log line is not JSON: %s", line)
		}
		for _, k := range []string{"timestamp", "level", "service", "version", "environment", "trace_id", "span_id", "correlation_id", "event", "message", "duration_ms"} {
			if _, ok := entry[k]; !ok {
				t.Fatalf("log line misses %q: %s", k, line)
			}
		}
		if entry["trace_id"] != traceID {
			t.Fatalf("log trace_id = %v, want %s", entry["trace_id"], traceID)
		}
	}
}

func TestOutcomesAreCountedInMetrics(t *testing.T) {
	var got http.Header
	noCopies := newServer(fakeCatalog(t, http.StatusConflict, 0, &got).URL, time.Second, io.Discard)
	if rec := post(t, noCopies); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("no copies: status = %d, want 422", rec.Code)
	}
	assertMetric(t, noCopies, `lms_circulation_loan_requests_total{outcome="rejected",route="/api/v1/loans"} 1`)

	slow := newServer(fakeCatalog(t, http.StatusOK, 300*time.Millisecond, &got).URL, 100*time.Millisecond, io.Discard)
	if rec := post(t, slow); rec.Code != http.StatusBadGateway {
		t.Fatalf("timeout: status = %d, want 502", rec.Code)
	}
	assertMetric(t, slow, `lms_circulation_loan_requests_total{outcome="error",route="/api/v1/loans"} 1`)
	assertMetric(t, slow, `operation="loan_copy",outcome="timeout_or_network",target="lms-catalog-api"`)
}

func assertMetric(t *testing.T, s *Server, want string) {
	t.Helper()
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if !strings.Contains(rec.Body.String(), want) {
		t.Fatalf("/metrics does not contain %s", want)
	}
}
