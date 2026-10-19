// Package tracecontext implements the W3C traceparent header used to join the
// logs of lms-circulation-api, lms-membership-api and lms-catalog-api into one
// trace (see 01-journey.md).
package tracecontext

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// Header is the W3C Trace Context header name.
const Header = "traceparent"

// SpanContext is the part of a span that travels between services.
type SpanContext struct {
	TraceID      string // 32 lowercase hex
	SpanID       string // 16 lowercase hex, the span that owns this context
	ParentSpanID string // span id received from the caller, empty on a root span
	Sampled      bool
}

// Parse reads "00-<trace_id>-<span_id>-<flags>". ok is false for any malformed
// or all-zero value, in which case the caller starts a new trace.
func Parse(value string) (sc SpanContext, ok bool) {
	parts := strings.Split(strings.TrimSpace(value), "-")
	if len(parts) != 4 || parts[0] != "00" {
		return SpanContext{}, false
	}
	traceID, spanID, flags := parts[1], parts[2], parts[3]
	if !isHex(traceID, 32) || !isHex(spanID, 16) || !isHex(flags, 2) {
		return SpanContext{}, false
	}
	if traceID == strings.Repeat("0", 32) || spanID == strings.Repeat("0", 16) {
		return SpanContext{}, false
	}
	return SpanContext{TraceID: traceID, SpanID: spanID, Sampled: flags == "01"}, true
}

// NewRoot starts a new trace when no valid traceparent arrived.
func NewRoot() SpanContext {
	return SpanContext{TraceID: randomHex(16), SpanID: randomHex(8), Sampled: true}
}

// Child opens a new span in the same trace, with sc as its parent.
func (sc SpanContext) Child() SpanContext {
	return SpanContext{TraceID: sc.TraceID, SpanID: randomHex(8), ParentSpanID: sc.SpanID, Sampled: sc.Sampled}
}

// String formats the context as a traceparent header value.
func (sc SpanContext) String() string {
	flags := "00"
	if sc.Sampled {
		flags = "01"
	}
	return fmt.Sprintf("00-%s-%s-%s", sc.TraceID, sc.SpanID, flags)
}

type ctxKey struct{}

type correlationKey struct{}

// WithCorrelationID stores the X-Correlation-Id value in ctx.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationKey{}, id)
}

// CorrelationID returns the id stored by WithCorrelationID, or "".
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationKey{}).(string)
	return id
}

// With stores sc in ctx.
func With(ctx context.Context, sc SpanContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, sc)
}

// From returns the span context stored in ctx, if any.
func From(ctx context.Context) (SpanContext, bool) {
	sc, ok := ctx.Value(ctxKey{}).(SpanContext)
	return sc, ok
}

func isHex(s string, n int) bool {
	if len(s) != n || strings.ToLower(s) != s {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
