# CP7 — Alert policy (burn rate)

Rules: `06-alerts.yml`. SLO: 99 % availability of `POST /api/v1/loans` in 30 days
(`05-slo.md`), budget = 1 % = **45 failed requests** or **432 minutes**.
Traffic: 12.5 loan requests per open hour, 4 500 per 30 days.

`burn rate = observed error ratio / budgeted error ratio (0.01)`. A burn rate of 1 spends
exactly the whole budget in 30 days (720 h).

## How much budget each alert has spent when it fires

| Alert | Condition | Error ratio needed | Budget spent before firing | In requests | In minutes of total outage |
|---|---|---|---|---|---|
| `LoanRegistrationAvailabilityFastBurn` (page) | burn 14.4 over 1 h **and** over 5 min | > 14.4 % | 14.4 x 1 h / 720 h = **2 %** | 2 % x 45 = **0.9 → 1 to 2 failed loans in one hour** | 2 % x 432 = **8.6 min** |
| `LoanRegistrationAvailabilitySlowBurn` (ticket) | burn 6 over 6 h **and** over 30 min | > 6 % | 6 x 6 h / 720 h = **5 %** | 5 % x 45 = **2.25 → 3 or more failed loans in 6 h** (out of ~75) | 5 % x 432 = **21.6 min** |
| `LoanRegistrationLatencyFastBurn` (page) | 14.4 x on the latency SLO (> 600 ms) | > 14.4 % slow | **2 %** of the latency budget | 1 to 2 slow loans in one hour | 8.6 min |

Reading the numbers for our traffic:

- With 12.5 requests per hour, **2 failed loans in one hour** (16 %) are enough for the
  fast-burn page, as long as the last 5 minutes also had a failure. That is the right
  sensitivity for a desk with a student waiting: two failures in a row mean the
  Administrator is already working around the system.
- The short windows (5 min and 30 min) stop the alert as soon as the failure stops, so a
  problem that was fixed 20 minutes ago does not keep paging.
- At night (19:00–07:00) there is no traffic, the ratios are `NaN` and nothing fires. This
  is correct: no loans, no user impact.
- Low traffic is also a risk: one failed request in a quiet hour can reach 14.4 %. If this
  produces false pages in `qa`, the fix is a minimum volume guard
  (`and sum(increase(lms_circulation_loan_requests_total[1h])) >= 5`), not a bigger
  threshold.

## Cause-based alerts we will NOT create

| Alert we will not create | Why not |
|---|---|
| CPU of the `lms-circulation-api` container > 80 % | The Go service can run hot without any loan failing, and loans can fail (catalog timeout) with the CPU at 5 %. It pages people without user impact and stays silent during real impact. CPU stays on the dashboard as a saturation signal, not as a page. |
| `lms-catalog-api` or `lms-membership-api` container restarted / MongoDB `loan_db` connections > N | A restart that `docker compose` heals in seconds spends almost no budget; if it does hurt users, the burn-rate alert already fires on the symptom. Alerting on every cause would create one alert per dependency and still miss the causes we did not think of. |

Both signals are still useful **after** the page, inside the runbook (`07-runbook.md`), to
find the cause.

## Routing

| Severity | Who | Channel | When |
|---|---|---|---|
| `page` | on-call developer of the week (rotation in the team) | team chat + phone call | desk opening hours |
| `ticket` | `lms-circulation-api` owner | GitHub issue in `code-corhuila/lms-circulation-api` with label `slo` | next working day |

## Validation

```text
$ promtool check rules 06-alerts.yml
Checking 06-alerts.yml
  SUCCESS: 9 rules found
```

(promtool 3.5.0 from the Prometheus Windows release, run on 2026-10-19.)
