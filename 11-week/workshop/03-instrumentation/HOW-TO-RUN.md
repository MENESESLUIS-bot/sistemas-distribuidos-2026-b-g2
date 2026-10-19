# CP3 — How to run the instrumentation module

Module: `lms-circulation-api/` (Go 1.25, same stack as the team service). It keeps only
hop 4 of HU-06 (`POST /api/v1/loans` → `lms-catalog-api` `loan-copy`) to show:

- `internal/tracecontext` — W3C `traceparent` parse/format, child spans
- `internal/httpx` — middleware that continues or starts the trace and returns it
- `internal/logging` — `log/slog` JSON in the format of `02-log-contract.md`
- `internal/catalogclient` — sends `traceparent` + `X-Correlation-Id` downstream
- `internal/metrics` — `/metrics` with the request counter and latency histograms

## Install

```bash
# Go 1.25.1 (https://go.dev/dl/), then:
cd 11-week/workshop/03-instrumentation/lms-circulation-api
go mod download
```

## Test

```bash
go vet ./...
go test ./... -v
```

Real output on my machine (Windows 11, go1.25.1 windows/amd64, 2026-10-19):

```text
?   	github.com/MENESESLUIS-bot/w11/lms-circulation-api/cmd/api	[no test files]
=== RUN   TestLoanCopyReceivesSameTraceID
--- PASS: TestLoanCopyReceivesSameTraceID (0.01s)
=== RUN   TestOutcomesAreCountedInMetrics
--- PASS: TestOutcomesAreCountedInMetrics (0.31s)
PASS
ok  	github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/api	1.369s
?   	github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/catalogclient	[no test files]
=== RUN   TestTraceContextContinuesIncomingTrace
--- PASS: TestTraceContextContinuesIncomingTrace (0.00s)
=== RUN   TestTraceContextStartsNewTraceOnInvalidHeader
--- PASS: TestTraceContextStartsNewTraceOnInvalidHeader (0.00s)
PASS
ok  	github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/httpx	(cached)
```

`TestLoanCopyReceivesSameTraceID` starts a fake `lms-catalog-api` with `httptest`, sends
`POST /api/v1/loans` with a known `traceparent` and checks that the fake catalog received
the **same `trace_id`** with a **new span id**, the same `X-Correlation-Id`, and that every
log line has the 11 required fields of the contract.

## Run

Terminal 1 — a fake `lms-catalog-api` on :8082 that prints what it receives:

```bash
node -e "require('http').createServer((q,s)=>{console.log('lms-catalog-api fake <-', q.method, q.url, 'traceparent=', q.headers.traceparent, 'x-correlation-id=', q.headers['x-correlation-id']); s.end()}).listen(8082)"
```

Terminal 2 — the service on :8083 (`CATALOG_SERVICE_URL`, `PORT`, `APP_VERSION`,
`APP_ENV` are optional env vars):

```bash
go run ./cmd/api
```

Terminal 3 — one loan with a known trace, then the metrics:

```bash
curl -i -X POST localhost:8083/api/v1/loans \
  -H "Content-Type: application/json" \
  -H "traceparent: 00-7c3e9a41d2b84f06a1e5c0b9d8f27e13-a01f5c3e9b72d846-01" \
  -H "X-Correlation-Id: 5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5" \
  -d '{"studentId":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","bookId":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82"}'
curl -s localhost:8083/metrics | grep ^lms_circulation
```

Real output:

```text
HTTP/1.1 201 Created
Content-Type: application/json
Traceparent: 00-7c3e9a41d2b84f06a1e5c0b9d8f27e13-9539ec6c6b5c9a0e-01
X-Correlation-Id: 5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5

{"id":"b530beba-23ff-44c1-8d80-c7a9eabd8d94","status":"ACTIVE","traceId":"7c3e9a41d2b84f06a1e5c0b9d8f27e13"}
```

Service logs (terminal 2):

```text
{"timestamp":"2026-10-19T14:44:42.639Z","level":"INFO","message":"listening on :8083","service":"lms-circulation-api","version":"2.0.1","environment":"dev","event":"server.started","duration_ms":0}
{"timestamp":"2026-10-19T14:44:43.153Z","level":"INFO","message":"register loan requested","service":"lms-circulation-api","version":"2.0.1","environment":"dev","event":"loan.request.received","duration_ms":0,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"9539ec6c6b5c9a0e","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5"}
{"timestamp":"2026-10-19T14:44:43.169Z","level":"INFO","message":"loan registered","service":"lms-circulation-api","version":"2.0.1","environment":"dev","event":"loan.registered","duration_ms":15,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","loan_id":"b530beba-23ff-44c1-8d80-c7a9eabd8d94","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"9539ec6c6b5c9a0e","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5"}
```

Fake catalog (terminal 1) — same `trace_id`, new client span `4572776df96bb9d4`:

```text
lms-catalog-api fake <- POST /api/v1/books/c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82/loan-copy traceparent= 00-7c3e9a41d2b84f06a1e5c0b9d8f27e13-4572776df96bb9d4-01 x-correlation-id= 5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5
```

Metrics:

```text
lms_circulation_downstream_request_duration_seconds_count{operation="loan_copy",outcome="success",target="lms-catalog-api"} 1
lms_circulation_loan_request_duration_seconds_count{route="/api/v1/loans"} 1
lms_circulation_loan_requests_total{outcome="success",route="/api/v1/loans"} 1
```

## What broke while building it

- The first run of the server failed with "No such file or directory" because the binary
  was started before the build had finished writing it; building first and then starting
  it fixed it.
- In the second run a fake catalog from the first run was still listening on :8082
  (`EADDRINUSE`); the request still reached that process, which is the line shown above.
