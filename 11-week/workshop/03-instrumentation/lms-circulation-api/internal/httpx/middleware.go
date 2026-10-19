// Package httpx holds the inbound HTTP middleware of lms-circulation-api.
package httpx

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

// CorrelationIDHeader is the business id the three LMS services already return.
const CorrelationIDHeader = "X-Correlation-Id"

// TraceContext reads the incoming traceparent (or starts a new trace), opens
// the server span for this request and returns it in the response headers,
// together with X-Correlation-Id.
func TraceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var server tracecontext.SpanContext
		if incoming, ok := tracecontext.Parse(r.Header.Get(tracecontext.Header)); ok {
			server = incoming.Child()
		} else {
			server = tracecontext.NewRoot()
		}

		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		w.Header().Set(tracecontext.Header, server.String())
		w.Header().Set(CorrelationIDHeader, correlationID)

		ctx := tracecontext.With(r.Context(), server)
		ctx = tracecontext.WithCorrelationID(ctx, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
