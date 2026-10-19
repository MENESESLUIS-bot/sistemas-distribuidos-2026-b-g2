# W11 Workshop — Observability, SLOs and incident response on LMS-Library

<!-- CONFIG-START -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- DOCS_REPO: https://github.com/code-corhuila/library-docs
- REFLECTION_DOC: https://docs.google.com/document/d/1kZ2divgrQ_O7ob9PucDD8yVcCxuL6HNsYG2eHIhDSlc/edit?usp=sharing
- JOURNEY: HU-06 Loan Registration — `POST /api/v1/loans` across lms-circulation-api, lms-membership-api and lms-catalog-api
<!-- CONFIG-END -->

## Scope

One critical user journey of the LMS-Library system is made observable end to end.
The services involved live in the `code-corhuila` organization: `lms-circulation-api`
(Go, MongoDB `loan_db`), `lms-membership-api` (Go, PostgreSQL `membership_db`) and
`lms-catalog-api` (Go, PostgreSQL `catalog_db`).

## Checkpoints

| CP | File | Status |
|---|---|---|
| CP0 | `README.md` | done |
| CP1 | `01-journey.md` | done |
| CP2 | `02-log-contract.md` | done |
| CP3 | `03-instrumentation/` | todo |
| CP4 | `04-metrics.md` | todo |
| CP6 | `05-slo.md` | todo |
| CP7 | `06-alerts.yml`, `06-alert-policy.md` | todo |
| CP8 | `07-runbook.md` | todo |
| CP9 | `08-incident.md` | todo |
