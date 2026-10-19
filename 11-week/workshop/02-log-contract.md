# CP2 — Structured log contract

One JSON object per line, written to stdout by `lms-circulation-api`,
`lms-membership-api` and `lms-catalog-api`. Field names are `snake_case` and identical
in the three services, so one query (`trace_id = X`) returns the whole HU-06 request.
This replaces the free-form fields of the current zap logger
(`internal/infrastructure/logger/logger.go`), which today only adds `timestamp`.

## Required fields

| Field | Type | Req. | Example | Why it exists |
|---|---|---|---|---|
| `timestamp` | string, ISO-8601 UTC with ms | yes | `2026-10-19T14:20:03.114Z` | Order events across services without clock-zone guessing |
| `level` | enum `DEBUG`/`INFO`/`WARN`/`ERROR` | yes | `ERROR` | Filter noise; alerting reads only `ERROR` |
| `service` | string | yes | `lms-circulation-api` | Which repo/container wrote the line |
| `version` | string, semver | yes | `2.0.1` | Tie a failure to a release tag (rollback decision) |
| `environment` | enum `dev`/`qa`/`main` | yes | `qa` | Same names as the team branches/environments |
| `trace_id` | string, 32 hex | yes | `7c3e9a41d2b84f06a1e5c0b9d8f27e13` | Join all lines of one request across services (W3C `traceparent`) |
| `span_id` | string, 16 hex | yes | `b41c07e9a2d35f68` | Which hop of the span tree (CP1) wrote the line |
| `correlation_id` | string, UUID | yes | `5f0c2b7e-...` | Business id from `X-Correlation-Id`, already returned to the UI |
| `event` | string, `noun.verb` dot case | yes | `catalog.loan_copy.timeout` | Stable key for counting and dashboards; `message` may change, `event` may not |
| `message` | string | yes | `lms-catalog-api did not answer in 5s` | Human explanation, never contains personal data |
| `duration_ms` | integer | yes (0 if n/a) | `5002` | Latency of the operation the line closes |

## Domain fields (optional, only when known)

| Field | Type | Example | Why it exists |
|---|---|---|---|
| `loan_id` | UUID | `9a7d…` | Id of the `Loan` aggregate in `loan_db.loans`; links logs with the stored document |
| `student_id` | UUID | `3f1b…` | Opaque Membership id (never name or document); needed to explain INV-003/INV-004 rejections |
| `book_id` | UUID | `c84e…` | Id of the row in `catalog.book`; explains `loan-copy` and `NO_COPIES_AVAILABLE` |
| `active_loans` | integer | `1` | Result of the INV-004 count (max 2) at decision time |
| `available_copies` | integer | `0` | Copies left in `catalog.book` after `loan-copy` |
| `error_code` | enum | `STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE` | Same code the API returns in the 422 body |
| `idempotency_key_hash` | string, sha256 prefix | `e3b0c442` | Detect retries of the same `Idempotency-Key` without storing the raw key |
| `attempt` | integer | `2` | Retry number inside the same trace |
| `http_status` | integer | `422` | Outcome of the hop |

## Never log (LMS-Library domain)

| Data | Where it appears | How it is handled |
|---|---|---|
| Student `FullName`, `DocumentID`, `Phone` | `lms-membership-api` `Student` | Not logged; use `student_id` |
| Student `Email` | `Student.Email` | Masked: `l***@corhuila.edu.co` |
| Administrator password | `POST /api/v1/auth/login` body (`lms-access-api`) | Request bodies of auth routes are never logged |
| RS256 session JWT / internal HS256 token | `Authorization` header | Header dropped; at most `token=***` |
| `INTERNAL_JWT_SECRET`, RS256 private key | env / config | Config is never printed at start-up |
| MongoDB URI with `circulation_app` password | `MONGO_URI` | Logged as `mongodb://***@host/loan_db` |
| Raw `Idempotency-Key` | request header | Only `idempotency_key_hash` (first 8 hex of sha256) |

## Seven lines of one failing request (same `trace_id`)

Story: the student has 1 active loan. `loan-copy` times out after 5 s, but `lms-catalog-api`
finishes late and takes the last copy. The portal retries with the same `Idempotency-Key`
inside the same trace, and the retry is rejected because no copy is left.

```json
{"timestamp":"2026-10-19T14:20:03.114Z","level":"INFO","service":"lms-circulation-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"a01f5c3e9b72d846","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"loan.request.received","message":"register loan requested","duration_ms":0,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","idempotency_key_hash":"e3b0c442","attempt":1}
{"timestamp":"2026-10-19T14:20:03.139Z","level":"INFO","service":"lms-membership-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"5d29e0b4c18a7f36","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"student.fetched","message":"student found, no active suspension","duration_ms":17,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","http_status":200}
{"timestamp":"2026-10-19T14:20:03.151Z","level":"DEBUG","service":"lms-circulation-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"e6a3b98f1c04d275","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"loan.active_count.checked","message":"student is under the loan limit","duration_ms":9,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","active_loans":1}
{"timestamp":"2026-10-19T14:20:08.155Z","level":"ERROR","service":"lms-circulation-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"3c7f12a8e5d90b64","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"catalog.loan_copy.timeout","message":"lms-catalog-api did not answer in 5s","duration_ms":5002,"book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","http_status":504,"attempt":1}
{"timestamp":"2026-10-19T14:20:08.571Z","level":"WARN","service":"lms-catalog-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"f8b05d61a3c2e947","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"book.copy.loaned","message":"copy decremented after the caller gave up","duration_ms":5418,"book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","available_copies":0,"http_status":200}
{"timestamp":"2026-10-19T14:20:10.204Z","level":"WARN","service":"lms-catalog-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"0a94c6e2b7d15f38","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"book.copy.unavailable","message":"no available copies to loan","duration_ms":6,"book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","available_copies":0,"http_status":409}
{"timestamp":"2026-10-19T14:20:10.212Z","level":"WARN","service":"lms-circulation-api","version":"2.0.1","environment":"qa","trace_id":"7c3e9a41d2b84f06a1e5c0b9d8f27e13","span_id":"d2e71b5a90c3f846","correlation_id":"5f0c2b7e-8d41-4a39-9e6b-21c7f0d3a8b5","event":"loan.rejected","message":"loan rejected by business rule","duration_ms":41,"student_id":"3f1b6a90-2c4e-4d7a-8f51-0b9e7c2d6a14","book_id":"c84e2d17-5a9b-4f3c-b6e0-7d1a9c3f5e82","error_code":"NO_COPIES_AVAILABLE","http_status":422,"attempt":2}
```

Filtered by `trace_id`, the lines show the full story, including the inconsistency: a copy
was taken in `catalog.book` but no `Loan` exists in `loan_db.loans`. That is the missing
saga compensation from week 10, now visible in the logs.

## Level policy

| Level | When the team uses it | LMS example |
|---|---|---|
| `ERROR` | The request could not be completed because of our system or a dependency; somebody may need to act | `catalog.loan_copy.timeout`, MongoDB `replaceOne` failure on `loans` |
| `WARN` | Expected but abnormal: a business rule said no, a retry, a dependency that was slow but answered | `loan.rejected` with `LOAN_LIMIT_REACHED`, `book.copy.loaned` after 5 s |
| `INFO` | Normal business milestones, one or two per request | `loan.request.received`, `loan.registered` with `loan_id` and `due_date` |
| `DEBUG` | Decision details for development; off in `main` | `loan.active_count.checked` with `active_loans` |
