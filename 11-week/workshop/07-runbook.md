# CP8 — Runbook: `LoanRegistrationAvailabilityFastBurn`

| Field | Value |
|---|---|
| Alert | `LoanRegistrationAvailabilityFastBurn` (`06-alerts.yml`), severity `page`, team `lms-library` |
| Journey | HU-06 `POST /api/v1/loans` (`01-journey.md`) |
| Service paged | `lms-circulation-api` (container `circulation-api`, port 8083) |
| Dependencies | `lms-membership-api` (`membership-api`), `lms-catalog-api` (`catalog-api`), MongoDB `loan_db` |
| Last good release | `v2.0.1` (tag in every LMS repo) |
| Budget at firing time | ≥ 2 % of the monthly budget already spent (`06-alert-policy.md`) |

## 1. Symptoms

- The Administrator gets an error page when registering a loan; the API answers 502/504
  with `DEPENDENCY_UNAVAILABLE` or 500.
- `lms_circulation_loan_requests_total{outcome="error"}` grows; more than 14.4 % of the
  loans in the last hour failed.
- Students may be told "no copies available" for a book that is on the shelf (orphan copy
  after a `loan-copy` timeout, see `02-log-contract.md`).

## 2. Confirm (5 minutes)

1. Is it real and still happening?

   ```promql
   lms_circulation:loan_error_ratio:rate5m
   ```

2. Which hop fails? p95 of each dependency and its outcomes:

   ```promql
   histogram_quantile(0.95, sum by (le, target) (rate(lms_circulation_downstream_request_duration_seconds_bucket[5m])))
   sum by (target, outcome) (rate(lms_circulation_downstream_request_duration_seconds_count[5m]))
   ```

3. Take one failed request and read its whole story. Copy `traceId` from a failed response
   body (or from an `ERROR` line) and filter the three services by it:

   ```bash
   T=<trace_id>
   for c in circulation-api membership-api catalog-api; do
     docker logs --since 30m "$c" 2>&1 | grep "\"trace_id\":\"$T\""
   done | sort
   ```

   The `event` field says where it broke: `catalog.loan_copy.timeout`,
   `catalog.loan_copy.failed`, `student.fetch.failed` or a MongoDB error on `loans`.

4. Health of each container:

   ```bash
   curl -s localhost:8083/health/ready; curl -s localhost:8082/health/ready; curl -s localhost:8081/health/ready
   docker ps --filter name=api --format "table {{.Names}}\t{{.Status}}\t{{.Image}}"
   ```

5. Was there a deploy in the last hour? `docker ps` shows the image age; compare with the
   last merge to `main` in `code-corhuila/lms-circulation-api` and `lms-catalog-api`.

## 3. Mitigate — first three actions, lowest risk first

| # | Action | When | How | Risk |
|---|---|---|---|---|
| 1 | **Restart the failing dependency** | `catalog-api` or `membership-api` not ready, pool exhausted, stuck after a DB restart | `docker restart catalog-api` (or `membership-api`) | Low: stateless, ~5 s of errors |
| 2 | **Roll back to `v2.0.1`** | The burn started right after a deploy of any of the three services | In the service repo: `git checkout v2.0.1 && docker compose -f deploy/compose.yml up -d --build` | Medium: loses the features of the new release; no data migration in `loan_db` or `catalog_db` between releases |
| 3 | **Hot-fix: fail fast instead of waiting 5 s** | Dependency is slow but alive (p95 > 1 s) and 1 and 2 did not help | The 5 s timeout is hard-coded in `catalogclient.NewClient`; a one-line PR lowers it to 1 s, promoted `develop -> qa -> main` with the Tech Lead's approval, then `docker compose up -d --build` of `circulation-api` | Higher: a code change under pressure; more 502 for a while, but no 5 s waits and fewer orphan copies |

Not available today (do not try at 2 a.m.): feature flags, horizontal scaling (the compose
files use a fixed `container_name`, so `--scale` fails) and a circuit breaker. They are
action items in `08-incident.md`.

After mitigation, fix **orphan copies**: for each `catalog.loan_copy.timeout` trace whose
`book_id` has no `Loan` in `loan_db.loans`, call
`POST /api/v1/books/{id}/return-copy` once and record the trace id in the incident.

## 4. Escalation

| After | Who | Role |
|---|---|---|
| 0 min | On-call developer of the week | Runs this runbook |
| 15 min without mitigation | Owner of the failing hop: Hermes Pascuas (HU-06, circulation) or Luis Meneses (catalog, HU-07) | Knows the service code |
| 30 min or data inconsistency | Oscar Areiza, Tech Lead | Decides rollback of more than one service, approves manual data fixes |
| Desk impact > 1 h | Library staff | Switch the desk to the paper loan form and register the loans later |

## 5. Verify recovery

- `lms_circulation:loan_error_ratio:rate5m` is 0 for 15 minutes in a row.
- One real loan registered from the portal returns 201 and its `trace_id` shows
  `loan.registered` with no `ERROR` line.
- `lms_catalog_loan_copy_total{result="decremented"}` and successful loans grow together
  again (SLO 3 of `05-slo.md`).
- The fast-burn alert resolves by itself (its 5-minute window clears). Leave a note with the
  start, end and failed-request count for the postmortem.
