# W13 Workshop — Threat modeling and security backlog on LMS-Library

<!-- CONFIG-START -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- DOCS_REPO: https://github.com/code-corhuila/library-docs
- REFLECTION_DOC: (pending)
- JOURNEY: HU-06 Loan Registration (`POST /api/v1/loans`)
- BOUNDARY: lms-circulation-api → lms-catalog-api (internal `POST /api/v1/books/{id}/loan-copy` and `/return-copy`)
- CLAIM_ISSUE: (pending)
- LIBRARY_BOOK: (pending — filled after the library visit)
- LIBRARY_RECORD_TIME: (pending)
- PR: (written in CP2)
- STORIES: (written in CP3)
<!-- CONFIG-END -->

## Why this boundary

HU-06 is the only flow of LMS-Library where one service changes data owned by another:
to register a `Loan`, `lms-circulation-api` asks `lms-catalog-api` to decrement
`available_copies` of a row in `catalog.book`. The call crosses from the Circulation
context (MongoDB `loan_db`) into the Catalog context (PostgreSQL `catalog_db`) with a
service token, so it is where a forged or misused token turns directly into wrong
inventory at the circulation desk.

## Checkpoints

| CP | Where | Status |
|---|---|---|
| CP0 | this README + claim issue | done |
| CP1 | `01-attack-surface.md` | todo |
| CP2 | PR in `library-docs` (`05-architecture/security/stride-MENESESLUIS-bot.md`) | todo |
| CP3 | three `[security]` issues in `library-docs` | todo |
| CP4 | `04-security-ci/` | todo |
| CP5 | peer review + `13-week/hu-status/README.md` | todo |
