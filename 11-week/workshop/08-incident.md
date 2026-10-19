# CP9 — Incident drill and blameless postmortem

## Scenario

Critical hop of `01-journey.md` is hop 4: `lms-circulation-api` →
`POST /api/v1/books/{id}/loan-copy` on `lms-catalog-api`. Drill: **for 40 minutes,
`loan-copy` answers in 3 s for 25 % of the requests**.

Drill assumptions:

- Date 2026-10-26 (first week of the new loan period), 10:00–10:40 local (15:00–15:40 UTC),
  the peak hour of the desk: **60 loan requests per hour** (1 per minute), five times the
  average of `05-slo.md`.
- Hypothetical cause used for the drill: while the Administrator types in the catalog
  search of `lms-catalog-portal`, `GET /api/v1/books?search=` runs an `ILIKE '%term%'` over
  `catalog.book` with no index, about 3 s per query. These queries hold connections of the
  `pgxpool` of `lms-catalog-api` (default `MaxConns` = max(4, CPUs), so 4 on a small host), and one `loan-copy` out of four
  waits ~3 s for a free connection.

## Trace of one slow request

```text
POST /api/v1/loans                          lms-circulation-api   3 112 ms  trace 4e81c07a3d9f42b6a5c1e8d0f7b2936c
├── RegisterLoan                            lms-circulation-api   3 104 ms
│   ├── GET /api/v1/students/{id}           -> lms-membership-api    21 ms
│   │   └── SELECT student                     membership_db          4 ms
│   ├── loans.countDocuments                   loan_db.loans          7 ms
│   ├── POST /api/v1/books/{id}/loan-copy   -> lms-catalog-api    3 046 ms   <- slow hop
│   │   ├── pgxpool.Acquire                    catalog_db         3 031 ms   <- waiting for a connection
│   │   └── UPDATE catalog.book                catalog_db             9 ms
│   └── loans.replaceOne                       loan_db.loans          12 ms
└── 201 Created
```

The request **succeeds**: 3.1 s is below the 5 s client timeout of `catalogclient`. No
orphan copy is created.

## Four log lines of that trace (format of `02-log-contract.md`)

```json
{"timestamp":"2026-10-26T15:12:40.208Z","level":"INFO","service":"lms-circulation-api","version":"2.0.1","environment":"main","trace_id":"4e81c07a3d9f42b6a5c1e8d0f7b2936c","span_id":"9b3e6d1f0a52c784","correlation_id":"b7d2e5a1-6c3f-4e90-8a14-3f9c2d7e0b61","event":"loan.request.received","message":"register loan requested","duration_ms":0,"student_id":"8d4a1f3e-7b2c-4c95-a061-5e2f9b8d3c70","book_id":"1c7e9b4d-3a8f-4d26-b5e1-9f0a6c2d7e43","attempt":1}
{"timestamp":"2026-10-26T15:12:43.257Z","level":"WARN","service":"lms-catalog-api","version":"2.0.1","environment":"main","trace_id":"4e81c07a3d9f42b6a5c1e8d0f7b2936c","span_id":"2f6a9c4e1d7b3058","correlation_id":"b7d2e5a1-6c3f-4e90-8a14-3f9c2d7e0b61","event":"book.copy.loaned","message":"copy decremented, waited for a db connection","duration_ms":3046,"book_id":"1c7e9b4d-3a8f-4d26-b5e1-9f0a6c2d7e43","available_copies":2,"http_status":200}
{"timestamp":"2026-10-26T15:12:43.261Z","level":"WARN","service":"lms-circulation-api","version":"2.0.1","environment":"main","trace_id":"4e81c07a3d9f42b6a5c1e8d0f7b2936c","span_id":"c5d81e3a7f9b2046","correlation_id":"b7d2e5a1-6c3f-4e90-8a14-3f9c2d7e0b61","event":"catalog.loan_copy.slow","message":"lms-catalog-api answered after 3s","duration_ms":3050,"book_id":"1c7e9b4d-3a8f-4d26-b5e1-9f0a6c2d7e43","http_status":200}
{"timestamp":"2026-10-26T15:12:43.320Z","level":"INFO","service":"lms-circulation-api","version":"2.0.1","environment":"main","trace_id":"4e81c07a3d9f42b6a5c1e8d0f7b2936c","span_id":"9b3e6d1f0a52c784","correlation_id":"b7d2e5a1-6c3f-4e90-8a14-3f9c2d7e0b61","event":"loan.registered","message":"loan registered","duration_ms":3112,"loan_id":"6a0f3c9e-2d7b-4e18-9c54-b1e8d2f7a036","student_id":"8d4a1f3e-7b2c-4c95-a061-5e2f9b8d3c70","book_id":"1c7e9b4d-3a8f-4d26-b5e1-9f0a6c2d7e43"}
```

## SLI values during the incident

| SLI (`05-slo.md`) | Normal | During the 40 min | Why |
|---|---|---|---|
| SLO 1 availability (target 99 %) | ~100 % | **100 %** | every slow request still ends in 201/422 before the 5 s timeout |
| SLO 2 latency ≤ 600 ms (target 99 %) | ~99.8 % | **75 %** | 1 of 4 loans takes ~3.1 s |
| SLO 3 consistency (target 99.5 %) | 100 % | **100 %** | no timeout, so no orphan copy |
| Latency burn rate | < 1 | **(1 − 0.75) / 0.01 = 25** | above the 14.4 page threshold |

## Which alert fires first, and when

- **Availability alerts never fire**: the error ratio stays at 0. A team with only the
  availability alerts of CP7 would not be paged at all.
- **`LoanRegistrationLatencyFastBurn` fires.** Its 5-minute window is at 25 % from minute 5.
  At minute `t` the 1-hour window still holds `60 − t` healthy minutes from before the
  incident, so its slow ratio is `0.25 x t / 60`. It passes 14.4 % when `t > 0.144 x 60 / 0.25 = 34.6 min`. With
  `for: 1m` and a 30 s evaluation interval, **the page arrives at minute 35 (10:35)**.
  Verified with `promtool test rules 08-incident-drill.test.yml`: not firing at minute 34,
  firing at minute 38, and the availability alerts stay silent.
- That is late: 35 of the 40 minutes were already slow. This is a finding of the drill.

## Which runbook step stops the damage

1. 10:35 page. Step 2.2 of `07-runbook.md` (dependency p95) shows `lms-catalog-api` at
   p95 ≈ 3 s and `lms-membership-api` at 20 ms: the slow hop is `loan-copy`.
2. 10:37 step 2.3: the trace above shows 3 031 ms in `pgxpool.Acquire`.
3. 10:38 **mitigation 1, `docker restart catalog-api`**: the long `ILIKE` searches are
   cancelled and the pool is free.
4. 10:40 step 5: `lms_circulation:loan_slow_ratio:rate5m` back to 0, alert resolves at 10:45.

## Postmortem (blameless)

**Summary.** From 10:00 to 10:40 on 2026-10-26, one in four loan registrations took about
3 seconds because `lms-catalog-api` ran out of database connections. No loan failed and no
data became inconsistent, but the circulation desk was slow during its busiest hour.

**Timeline (local time)**

| Time | Event |
|---|---|
| 10:00 | Catalog searches start holding `pgxpool` connections; 25 % of `loan-copy` calls wait 3 s |
| 10:05 | Latency 5-minute slow ratio reaches 25 % (no alert: the 1 h window is still low) |
| 10:20 | Library staff mention that "the system is slow"; no ticket is opened |
| 10:35 | `LoanRegistrationLatencyFastBurn` pages the on-call developer |
| 10:37 | Trace shows the wait in `pgxpool.Acquire` of `lms-catalog-api` |
| 10:38 | `docker restart catalog-api` (runbook mitigation 1) |
| 10:40 | Latency back to normal; incident over |
| 10:45 | Alert resolved |

**Impact**

- Requests: 40 loan requests in the incident, **10 slow** (≥ 600 ms), 0 failed.
- Latency budget (45 slow requests / 30 days): 10 / 45 = **22 %** spent in 40 minutes.
- In minutes: 40 min x 25 % = **10 minutes** of full-slowness equivalent = 10 / 432 =
  **2.3 %** of the time budget. The two numbers differ because the incident happened at
  peak hour (60 req/h, 5x the average): counting requests is the user's view.
- Availability and consistency budgets: 0 % spent.

**Root cause.** `GET /api/v1/books?search=` runs an unindexed `ILIKE '%term%'` on
`catalog.book` and shares the same `pgxpool` (4 connections) as `loan-copy`. Slow reads
starved a critical write.

**Contributing factors**

- The 5 s client timeout in `catalogclient` hides slowness: nothing fails, so the
  availability alerts stay silent.
- The latency fast-burn alert needs ~35 minutes to fire because the 1 h window dilutes a
  partial slowdown (25 %).
- Before W11 the services did not propagate `traceparent`; without it the 3 s could not be
  placed in `pgxpool.Acquire` in two minutes.
- `lms-catalog-api` exposes no pool saturation metric yet (`04-metrics.md` defines
  `lms_catalog_db_pool_acquired_connections`, not implemented).

**Action items (user stories)**

1. **Protect `loan-copy` from slow searches** — *As the Administrator, I want a loan to be
   registered fast even while the catalog is being searched, so that the student at the desk
   does not wait.*
   - Given 20 concurrent `GET /api/v1/books?search=` of 3 s in `qa`, when 10 `POST
     /api/v1/books/{id}/loan-copy` run, then their p99 is below 300 ms (k6 script in the repo).
   - `GET /books?search=` has a `statement_timeout` of 1 s and uses a trigram index on
     `catalog.book.title`/`author` (Liquibase changeset in `lms-catalog-db`).
2. **Propagate trace context in the three services** — *As the on-call developer, I want
   every log line of a loan to carry the same `trace_id`, so that I can find the slow hop in
   minutes.*
   - Given a `POST /api/v1/loans` with `traceparent` T, when it completes, then the log
     lines of `circulation-api`, `membership-api` and `catalog-api` all contain `trace_id`
     T (integration test in `lms-circulation-api` with both clients).
   - `catalogclient` and `membershipclient` send `traceparent` and `X-Correlation-Id`
     (the code of `03-instrumentation/` moved into the team repo).
3. **Detect partial slowdowns sooner** — *As the on-call developer, I want to be paged
   when the catalog becomes slow for part of the loans, before most of the hour is lost.*
   - `lms-catalog-api` exports `lms_catalog_db_pool_acquired_connections` and
     `lms_catalog_db_pool_max_connections` from `pgxpool.Stat()`.
   - Given the same drill in `qa` (25 % of `loan-copy` at 3 s), when it starts, then a
     page arrives in **10 minutes or less** (a 6x burn on a 15 min/5 min window pair is
     added to `06-alerts.yml` and tested with `promtool test rules`).
