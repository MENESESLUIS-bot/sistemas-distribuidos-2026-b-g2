# CP4 — Golden signals and metric catalog

Scope: the three services of the HU-06 journey (`01-journey.md`). Naming follows the
Prometheus conventions: `snake_case`, a `lms_<service>_` prefix, the base unit as suffix
(`_seconds`), `_total` on counters. Routes are always the **template**
(`/api/v1/books/{id}/loan-copy`), never the raw path.

## lms-circulation-api (entry point of HU-06)

| Signal | Metric | Type | Labels | PromQL (dashboard panel) |
|---|---|---|---|---|
| Latency | `lms_circulation_loan_request_duration_seconds` | histogram | `route` | `histogram_quantile(0.95, sum by (le) (rate(lms_circulation_loan_request_duration_seconds_bucket{route="/api/v1/loans"}[5m])))` |
| Traffic | `lms_circulation_loan_requests_total` | counter | `route`, `outcome` | `sum(rate(lms_circulation_loan_requests_total{route="/api/v1/loans"}[5m]))` |
| Errors | `lms_circulation_loan_requests_total` | counter | `route`, `outcome` | `sum(rate(lms_circulation_loan_requests_total{route="/api/v1/loans",outcome="error"}[5m])) / sum(rate(lms_circulation_loan_requests_total{route="/api/v1/loans"}[5m]))` |
| Saturation | `lms_circulation_mongo_pool_in_use_connections` / `lms_circulation_mongo_pool_max_connections` | gauge | none | `max(lms_circulation_mongo_pool_in_use_connections) / max(lms_circulation_mongo_pool_max_connections)` |
| Dependency latency | `lms_circulation_downstream_request_duration_seconds` | histogram | `target`, `operation`, `outcome` | `histogram_quantile(0.95, sum by (le, target) (rate(lms_circulation_downstream_request_duration_seconds_bucket[5m])))` |

`outcome` values: `success` (201), `rejected` (422 `STUDENT_SUSPENDED`,
`LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`), `client_error` (400/401) and `error`
(5xx, dependency timeout). A business rejection is **not** an error: the system worked.
The first three metrics and the dependency histogram are already implemented in
`03-instrumentation/lms-circulation-api/internal/metrics`.

## lms-membership-api (hop 2: eligibility)

| Signal | Metric | Type | Labels | PromQL |
|---|---|---|---|---|
| Latency | `lms_membership_http_request_duration_seconds` | histogram | `route`, `method` | `histogram_quantile(0.95, sum by (le) (rate(lms_membership_http_request_duration_seconds_bucket{route="/api/v1/students/{id}",method="GET"}[5m])))` |
| Traffic | `lms_membership_http_requests_total` | counter | `route`, `method`, `status_class` | `sum by (route) (rate(lms_membership_http_requests_total[5m]))` |
| Errors | `lms_membership_http_requests_total` | counter | `route`, `method`, `status_class` | `sum(rate(lms_membership_http_requests_total{status_class="5xx"}[5m])) / sum(rate(lms_membership_http_requests_total[5m]))` |
| Saturation | `lms_membership_db_pool_acquired_connections` / `lms_membership_db_pool_max_connections` (from `pgxpool.Stat()`) | gauge | none | `max(lms_membership_db_pool_acquired_connections) / max(lms_membership_db_pool_max_connections)` |

## lms-catalog-api (hop 4: loan-copy, critical hop)

| Signal | Metric | Type | Labels | PromQL |
|---|---|---|---|---|
| Latency | `lms_catalog_http_request_duration_seconds` | histogram | `route`, `method` | `histogram_quantile(0.95, sum by (le) (rate(lms_catalog_http_request_duration_seconds_bucket{route="/api/v1/books/{id}/loan-copy"}[5m])))` |
| Traffic | `lms_catalog_http_requests_total` | counter | `route`, `method`, `status_class` | `sum(rate(lms_catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy"}[5m]))` |
| Errors | `lms_catalog_http_requests_total` | counter | `route`, `method`, `status_class` | `sum(rate(lms_catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy",status_class="5xx"}[5m])) / sum(rate(lms_catalog_http_requests_total{route="/api/v1/books/{id}/loan-copy"}[5m]))` |
| Saturation | `lms_catalog_db_pool_acquired_connections` / `lms_catalog_db_pool_max_connections` | gauge | none | `max(lms_catalog_db_pool_acquired_connections) / max(lms_catalog_db_pool_max_connections)` |
| Domain | `lms_catalog_loan_copy_total` | counter | `result` (`decremented`, `no_copies`, `not_found`) | `sum by (result) (increase(lms_catalog_loan_copy_total[1h]))` |

## Label cardinality

| Label | Values | Cardinality |
|---|---|---|
| `route` | route templates of each service (`/api/v1/loans`, `/api/v1/loans/{id}/return`, `/api/v1/students/{id}`, `/api/v1/books/{id}/loan-copy`, ...) | 4 to 7 per service |
| `method` | GET, POST, PATCH, DELETE | 4 |
| `outcome` | success, rejected, client_error, error | 4 |
| `status_class` | 2xx, 3xx, 4xx, 5xx | 4 |
| `target` | lms-membership-api, lms-catalog-api | 2 |
| `operation` | get_student, suspend, loan_copy, return_copy | 4 |
| `result` | decremented, no_copies, not_found | 3 |

Worst case, `lms_catalog_http_request_duration_seconds` has 7 routes x 4 methods x
(9 buckets + `+Inf` + `_sum` + `_count`) = 336 series. The three services together stay
under 1 500 series.

**Forbidden labels** (unbounded) and what replaces them:

| Forbidden | Why | Replacement |
|---|---|---|
| `student_id`, `book_id`, `loan_id` | one series per student/book/loan, grows with the library | `student_id`/`book_id`/`loan_id` fields in the log line (CP2) |
| raw URL (`/api/v1/books/c84e.../loan-copy`) | one series per book | `route` template from the router |
| student `email`, `document_id` | unbounded and personal data | nothing; never leaves `lms-membership-api` |
| `trace_id`, `Idempotency-Key` | unique per request | logs and, later, histogram exemplars |
| `error_code` free text | could grow with new messages | closed `outcome` enum |

## Latency buckets

`0.05, 0.1, 0.2, 0.3, 0.5, 0.6, 1, 2.5, 5` seconds.

- The user is the Administrator at the circulation desk with a student waiting in front of
  them. NFR-001 asks p95 < 300 ms and p99 < 600 ms for `POST /loans`, so `0.3` and `0.6`
  are exact bucket edges and the SLO in `05-slo.md` can be computed without interpolation.
- `0.05`, `0.1`, `0.2` split the healthy range: the loan does 2 HTTP hops and 2 MongoDB
  operations, and a healthy request is around 50 to 150 ms.
- `1` and `2.5` show degradation that a person notices but still waits for.
- `5` is the HTTP client timeout of `catalogclient` and `membershipclient`; anything slower
  ends in `+Inf` and is already a failed request.
