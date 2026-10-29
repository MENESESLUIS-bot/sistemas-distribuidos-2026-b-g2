# W13 Workshop — Threat modeling and security backlog on LMS-Library

<!-- CONFIG-START -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- DOCS_REPO: https://github.com/code-corhuila/library-docs
- REFLECTION_DOC: https://docs.google.com/document/d/1yi0_c2tiVAOLDE_Lyxr8Rf8MZUbYPpeg8hm70jFbBdE/edit?usp=sharing
- JOURNEY: HU-06 Loan Registration (`POST /api/v1/loans`)
- BOUNDARY: lms-circulation-api → lms-catalog-api (internal `POST /api/v1/books/{id}/loan-copy` and `/return-copy`)
- CLAIM_ISSUE: https://github.com/code-corhuila/library-docs/issues/34
- LIBRARY_BOOK: (pending — filled after the library visit)
- LIBRARY_RECORD_TIME: (pending)
- PR: https://github.com/code-corhuila/library-docs/pull/36
- STORIES: https://github.com/code-corhuila/library-docs/issues/37 (E, risk 9, P0), https://github.com/code-corhuila/library-docs/issues/38 (T, risk 9, P0), https://github.com/code-corhuila/library-docs/issues/39 (S, risk 6, P1)
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
| CP1 | `01-attack-surface.md` | done |
| CP2 | PR in `library-docs` (`05-architecture/security/stride-MENESESLUIS-bot.md`) | done |
| CP3 | three `[security]` issues in `library-docs` | done |
| CP4 | `04-security-ci/` | done |
| CP5 | peer review + `13-week/hu-status/README.md` | done |
