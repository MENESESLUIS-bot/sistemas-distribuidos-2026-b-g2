package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

const incoming = "00-7c3e9a41d2b84f06a1e5c0b9d8f27e13-a01f5c3e9b72d846-01"

func TestTraceContextContinuesIncomingTrace(t *testing.T) {
	var seen tracecontext.SpanContext
	h := TraceContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen, _ = tracecontext.From(r.Context())
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", nil)
	req.Header.Set(tracecontext.Header, incoming)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if seen.TraceID != "7c3e9a41d2b84f06a1e5c0b9d8f27e13" {
		t.Fatalf("trace_id = %q, want the incoming one", seen.TraceID)
	}
	if seen.ParentSpanID != "a01f5c3e9b72d846" || seen.SpanID == seen.ParentSpanID {
		t.Fatalf("server span %q must be a new child of a01f5c3e9b72d846", seen.SpanID)
	}
	if got := rec.Header().Get(tracecontext.Header); got != seen.String() {
		t.Fatalf("response traceparent = %q, want %q", got, seen.String())
	}
	if rec.Header().Get(CorrelationIDHeader) == "" {
		t.Fatal("X-Correlation-Id missing in response")
	}
}

func TestTraceContextStartsNewTraceOnInvalidHeader(t *testing.T) {
	for _, bad := range []string{"", "garbage", "00-00000000000000000000000000000000-a01f5c3e9b72d846-01"} {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/loans", nil)
		req.Header.Set(tracecontext.Header, bad)
		rec := httptest.NewRecorder()
		TraceContext(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rec, req)

		if _, ok := tracecontext.Parse(rec.Header().Get(tracecontext.Header)); !ok {
			t.Fatalf("header %q: response traceparent is not valid", bad)
		}
	}
}
