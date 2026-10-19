// Package logging writes the JSON line format of 02-log-contract.md:
// timestamp, level, service, version, environment, trace_id, span_id,
// correlation_id, event, message, duration_ms plus optional domain fields.
package logging

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

// New returns a logger whose every line carries the service identity and,
// when the context has one, the trace and correlation ids.
func New(w io.Writer, service, version, environment string, level slog.Level) *slog.Logger {
	json := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level, ReplaceAttr: rename})
	return slog.New(&traceHandler{Handler: json}).With(
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
	)
}

// Event writes one contract line. event is the stable "noun.verb" key and
// duration the time spent in the operation the line closes.
func Event(ctx context.Context, log *slog.Logger, level slog.Level, event, message string, duration time.Duration, attrs ...slog.Attr) {
	base := []slog.Attr{
		slog.String("event", event),
		slog.Int64("duration_ms", duration.Milliseconds()),
	}
	log.LogAttrs(ctx, level, message, append(base, attrs...)...)
}

// rename maps slog's default keys to the contract: time -> timestamp in
// ISO-8601 UTC with milliseconds, msg -> message, level in upper case.
func rename(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	switch a.Key {
	case slog.TimeKey:
		return slog.String("timestamp", a.Value.Time().UTC().Format("2006-01-02T15:04:05.000Z"))
	case slog.MessageKey:
		a.Key = "message"
	case slog.LevelKey:
		a.Value = slog.StringValue(a.Value.Any().(slog.Level).String())
	}
	return a
}

// traceHandler injects trace_id, span_id and correlation_id from the request
// context, so handlers never have to pass them by hand.
type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc, ok := tracecontext.From(ctx); ok {
		r.AddAttrs(slog.String("trace_id", sc.TraceID), slog.String("span_id", sc.SpanID))
	}
	if id := tracecontext.CorrelationID(ctx); id != "" {
		r.AddAttrs(slog.String("correlation_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name)}
}
