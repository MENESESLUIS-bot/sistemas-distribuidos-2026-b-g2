# CP6 — SLIs, SLOs and error budget

Journey: HU-06 `POST /api/v1/loans` (`01-journey.md`). Metrics: `04-metrics.md`.
Window for every SLO: **30 days, rolling**.

## Traffic estimate

| Input | Value | Source / reasoning |
|---|---|---|
| Desk opening hours | 12 h/day (07:00–19:00), every day of the window | academic library, one circulation desk |
| Loan requests per day | 150 (including 422 rejections and retries) | about 120 real loans plus ~25 % rejected or retried attempts |
| Requests per hour while open | 150 / 12 = **12.5** | used by the alert policy (CP7) |
| Requests in 30 days | 150 x 30 = **4 500** | denominator of every budget below |
| Minutes in 30 days | 30 x 24 x 60 = **43 200** | time-based budget |

The number is small on purpose: LMS-Library serves one university library, not a
national service. NFR-001 asks the system to support 20 RPS, which is capacity, not demand.

## SLO 1 — Availability of loan registration

- **SLI** = good events / valid events
  - valid = every `POST /api/v1/loans` except `client_error` (malformed body, bad token:
    the caller's fault)
  - good = `success` (201) **or** `rejected` (422 `STUDENT_SUSPENDED`,
    `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`), because the system answered correctly
- **Objective: 99 %** of valid requests in 30 days.

```promql
sum(increase(lms_circulation_loan_requests_total{route="/api/v1/loans",outcome=~"success|rejected"}[30d]))
/
sum(increase(lms_circulation_loan_requests_total{route="/api/v1/loans",outcome!="client_error"}[30d]))
```

- **Budget**: 1 % of 4 500 = **45 failed loan requests** per 30 days, or
  1 % of 43 200 = **432 minutes (7.2 h)** of total outage.
- Why 99 % and not the 95 % of NFR-002: 95 % (36 h/month) is the target for the whole
  deployment including maintenance. The loan desk has a student waiting in front of the
  Administrator, so this journey gets a stricter target; 99 % is still realistic for a
  single Docker host with `docker compose`. 99.9 % (43 min) is not: one `compose up` of a
  new image already takes a few minutes.

## SLO 2 — Latency of loan registration

- **SLI** = requests answered in ≤ 600 ms / all requests recorded by the latency
  histogram of `POST /api/v1/loans` (it is labelled only by `route`, so a slow error also
  counts as slow)
- **Objective: 99 %** under 600 ms in 30 days (NFR-001 p99 < 600 ms). The p95 < 300 ms
  of NFR-001 stays as a dashboard target, not as an SLO.

```promql
sum(increase(lms_circulation_loan_request_duration_seconds_bucket{route="/api/v1/loans",le="0.6"}[30d]))
/
sum(increase(lms_circulation_loan_request_duration_seconds_count{route="/api/v1/loans"}[30d]))
```

`0.6` is an exact bucket edge (`04-metrics.md`), so no interpolation is needed.

- **Budget**: 1 % of 4 500 = **45 slow requests**, or **432 minutes** fully slow.

## SLO 3 — Loan / copy consistency (product specific)

Registering a loan changes two databases: `catalog.book.available_copies` in
`catalog_db` and `loans` in `loan_db`. There is no saga compensation yet (week 10), so a
timeout after `loan-copy` can take a copy with no `Loan` behind it (see the trace in
`02-log-contract.md`). For the Administrator this is the worst failure: the system says
there are no copies while the book is on the shelf.

- **SLI** = loans registered / copies decremented by `loan-copy`
- **Objective: 99.5 %** in 30 days.

```promql
sum(increase(lms_circulation_loan_requests_total{route="/api/v1/loans",outcome="success"}[30d]))
/
sum(increase(lms_catalog_loan_copy_total{result="decremented"}[30d]))
```

- **Budget**: about 3 600 successful loans per month (4 500 minus rejections) x 0.5 % =
  **18 orphan copies** per month. A time budget does not apply here: one orphan copy stays
  wrong until somebody fixes it by hand, so this budget is counted in copies only.

## Error budget policy (team LMS-Library)

Measured on SLO 1 and SLO 2; SLO 3 has its own rule at the end.

| Budget consumed | Action |
|---|---|
| **50 %** (22 failed or slow requests, 216 min) | Warning in the team channel. The next PR to `develop` of the guilty service needs the trace of one failed request in its description. No change to the release plan. |
| **75 %** (34 requests, 324 min) | Feature freeze for `lms-circulation-api` and the service that caused it. Only fixes and reliability stories go `develop -> qa -> main`. Tech Lead (Oscar Areiza) reviews every promotion to `main`. |
| **100 %** (45 requests, 432 min) | Full freeze until the 30-day window recovers. The sprint goal switches to reliability: retries/timeouts in `catalogclient`, compensation for `loan-copy`, and a postmortem (`08-incident.md` format) in `library-docs/13-operations`. |

SLO 3: every orphan copy found is fixed with `POST /api/v1/books/{id}/return-copy` the same
day; at 9 orphan copies (50 %) the compensation story moves to the top of the backlog.
