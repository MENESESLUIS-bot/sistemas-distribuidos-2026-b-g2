// Package api is a minimal POST /api/v1/loans (HU-06) that keeps only the
// hop to lms-catalog-api; the Membership check and the MongoDB write of the
// real lms-circulation-api are left out to keep the instrumentation visible.
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/catalogclient"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/httpx"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/logging"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/metrics"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

const loansRoute = "/api/v1/loans"

type Server struct {
	Catalog *catalogclient.Client
	Log     *slog.Logger
	Metrics *metrics.Metrics
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", s.Metrics.Handler())
	mux.Handle("POST "+loansRoute, httpx.TraceContext(http.HandlerFunc(s.createLoan)))
	return mux
}

type createLoanRequest struct {
	StudentID string `json:"studentId"`
	BookID    string `json:"bookId"`
}

func (s *Server) createLoan(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()

	var req createLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookID == "" || req.StudentID == "" {
		s.finish(w, r, start, "client_error", http.StatusBadRequest, map[string]string{"code": "VALIDATION_ERROR"})
		return
	}
	ids := []slog.Attr{slog.String("student_id", req.StudentID), slog.String("book_id", req.BookID)}
	logging.Event(ctx, s.Log, slog.LevelInfo, "loan.request.received", "register loan requested", 0, ids...)

	callStart := time.Now()
	err := s.Catalog.LoanCopy(ctx, req.BookID)
	switch {
	case err == nil:
		loanID := uuid.NewString()
		logging.Event(ctx, s.Log, slog.LevelInfo, "loan.registered", "loan registered",
			time.Since(start), append(ids, slog.String("loan_id", loanID))...)
		s.finish(w, r, start, "success", http.StatusCreated, map[string]string{"id": loanID, "status": "ACTIVE"})
	case errors.Is(err, catalogclient.ErrNoCopiesAvailable):
		logging.Event(ctx, s.Log, slog.LevelWarn, "loan.rejected", "loan rejected by business rule",
			time.Since(start), append(ids, slog.String("error_code", "NO_COPIES_AVAILABLE"))...)
		s.finish(w, r, start, "rejected", http.StatusUnprocessableEntity, map[string]string{"code": "NO_COPIES_AVAILABLE"})
	case errors.Is(err, catalogclient.ErrBookNotFound):
		s.finish(w, r, start, "client_error", http.StatusBadRequest, map[string]string{"code": "VALIDATION_ERROR"})
	default:
		logging.Event(ctx, s.Log, slog.LevelError, "catalog.loan_copy.failed", err.Error(), time.Since(callStart), ids...)
		s.finish(w, r, start, "error", http.StatusBadGateway, map[string]string{"code": "DEPENDENCY_UNAVAILABLE"})
	}
}

// finish records the request metrics and writes the JSON body with the traceId
// the team's error envelope already carries.
func (s *Server) finish(w http.ResponseWriter, r *http.Request, start time.Time, outcome string, status int, body map[string]string) {
	s.Metrics.ObserveRequest(loansRoute, outcome, time.Since(start))
	if sc, ok := tracecontext.From(r.Context()); ok {
		body["traceId"] = sc.TraceID
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
