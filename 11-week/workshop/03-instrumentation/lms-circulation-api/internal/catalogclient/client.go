// Package catalogclient calls lms-catalog-api's POST /api/v1/books/{id}/loan-copy
// (hop 4 of 01-journey.md) and propagates the trace context, which the team's
// current client does not do.
package catalogclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/tracecontext"
)

// Errors mirror lms-catalog-api's 404 and 409 on loan-copy.
var (
	ErrBookNotFound      = errors.New("book not found")
	ErrNoCopiesAvailable = errors.New("no copies available")
)

// Observer receives the duration and outcome of every call, for metrics.
type Observer func(operation, outcome string, d time.Duration)

type Client struct {
	baseURL string
	http    *http.Client
	observe Observer
}

func New(baseURL string, timeout time.Duration, observe Observer) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: timeout}, observe: observe}
}

// LoanCopy asks lms-catalog-api to decrement available_copies of a book. It
// opens a client span as a child of the request span stored in ctx and sends
// it as traceparent, together with X-Correlation-Id.
func (c *Client) LoanCopy(ctx context.Context, bookID string) (err error) {
	start := time.Now()
	outcome := "error"
	defer func() { c.observe("loan_copy", outcome, time.Since(start)) }()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/api/v1/books/%s/loan-copy", c.baseURL, bookID), nil)
	if err != nil {
		return err
	}
	if parent, ok := tracecontext.From(ctx); ok {
		req.Header.Set(tracecontext.Header, parent.Child().String())
	}
	if id := tracecontext.CorrelationID(ctx); id != "" {
		req.Header.Set("X-Correlation-Id", id)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		outcome = "timeout_or_network"
		return fmt.Errorf("calling lms-catalog-api: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		outcome = "success"
		return nil
	case http.StatusNotFound:
		outcome = "rejected"
		return ErrBookNotFound
	case http.StatusConflict:
		outcome = "rejected"
		return ErrNoCopiesAvailable
	default:
		return fmt.Errorf("lms-catalog-api returned %d", resp.StatusCode)
	}
}
