<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       11-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 11

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Make the HU-06 loan registration journey (lms-circulation-api -> lms-membership-api -> lms-catalog-api) observable end to end: trace map, structured log contract, trace-context instrumentation with tests, metric catalog, and then SLOs, burn-rate alerts, runbook and an incident drill.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W11-CP1 | Critical journey HU-06 and trace map | done | [fb082f4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/fb082f4) |
| W11-CP2 | Structured log contract | done | [01340eb](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/01340eb) |
| W11-CP3 | Trace-context middleware, logger, propagation to lms-catalog-api, metrics and test | done | [4927fdc](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/4927fdc), [fad81a8](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/fad81a8), [516f394](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/516f394) |
| W11-CP4 | Golden signals and metric catalog | done | [46ee418](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/46ee418) |
| W11-CP6..CP9 | SLOs, burn-rate alerts, runbook, incident drill (delivery 2) | todo | — |

## 2. My individual contribution
- `11-week/workshop/01-journey.md`: hop table, Mermaid sequence diagram with the `traceparent` on every arrow, and the span tree of `POST /api/v1/loans` across `lms-circulation-api`, `lms-membership-api` (`GET /api/v1/students/{id}`), MongoDB `loan_db.loans` and `lms-catalog-api` (`POST /api/v1/books/{id}/loan-copy`); `LoanRegistered` mapped as the planned async boundary.
- `11-week/workshop/02-log-contract.md`: 11 required fields, 9 domain fields (`loan_id`, `student_id`, `book_id`, `error_code`, ...), never-log list for student PII, JWTs and secrets, and seven JSON lines of one trace showing a `loan-copy` timeout and a `NO_COPIES_AVAILABLE` rejection.
- `11-week/workshop/03-instrumentation/lms-circulation-api` (Go 1.25): W3C `traceparent` parsing, middleware that continues or starts the trace, `log/slog` JSON logger in the contract format, `catalogclient` that forwards `traceparent` and `X-Correlation-Id`, `/metrics` with `lms_circulation_loan_requests_total` and two latency histograms, and 4 passing tests (`HOW-TO-RUN.md` has the real console output).
- `11-week/workshop/04-metrics.md`: four golden signals for the three services with PromQL, label cardinality, forbidden labels and latency buckets based on NFR-001.

## 3. Blockers and risks
- The team services still do not propagate `traceparent`: `catalogclient` and `membershipclient` of `lms-circulation-api` only send `Authorization`.
- `LoanPeriodDays` is 3 in `lms-circulation-api` `develop` but 7 in `library-docs`.
- No Prometheus instance in the LMS compose yet, so the metrics are only scraped manually.

## 4. Plan for next week
- Delivery 2 of the workshop: SLOs and error budget, burn-rate alerts, runbook and incident drill with postmortem.
- Propose a PR to `lms-circulation-api` that moves the trace propagation of `03-instrumentation/` into both HTTP clients.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Workshop folder: [11-week/workshop](../workshop/)
- Commits of delivery 1: [50abc14](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/50abc14) (CP0), [fb082f4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/fb082f4) (CP1), [01340eb](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/01340eb) (CP2), [14e8b4c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/14e8b4c) / [4927fdc](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/4927fdc) / [3418981](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3418981) (CP3 push 1), [16009db](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/16009db) / [fad81a8](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/fad81a8) / [293699f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/293699f) / [516f394](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/516f394) / [3e7139a](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3e7139a) (CP3 push 2), [46ee418](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/46ee418) (CP4)
