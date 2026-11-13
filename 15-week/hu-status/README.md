<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       15-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 15

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Promote my services through `develop → qa → main` with `cherry-pick -x` and cut their first releases: the HU-01 Vue login remote (`lms-access-portal` v3.0.0), the Node.js/Express rewrite of `lms-access-api` (v3.0.0, plus the `body-parser` DoS fix) and the overdue penalty saga (`lms-workflow` v3.0.0, with its first test suite and a retry fix); close the STRIDE threat model of the `lms-circulation-api → lms-catalog-api` internal calls in `library-docs`; track the HU-09 book-editing endpoint (`PATCH /books/{id}`) in `lms-catalog-api`; and review the team's pending PRs, approving the `lms-api-gateway` v1.0.0 release (ADR-007).
<!-- CONFIG-END -->

## 1. User stories worked this week

| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01 | Administrator authentication: Vue login remote in `lms-access-portal`, promoted to `qa` and released as v3.0.0 | done | [lms-access-portal#3](https://github.com/code-corhuila/lms-access-portal/pull/3), [lms-access-portal#5](https://github.com/code-corhuila/lms-access-portal/pull/5) (replaces [#4](https://github.com/code-corhuila/lms-access-portal/pull/4)) |
| W15-04 | `lms-access-api` Node.js/Express rewrite (ADR-011), promoted to `qa` and released as v3.0.0 | done | [lms-access-api#9](https://github.com/code-corhuila/lms-access-api/pull/9), [lms-access-api#11](https://github.com/code-corhuila/lms-access-api/pull/11) |
| W15-05 | Patch the `body-parser` DoS vulnerability (GHSA-27r3-jxmh-p828) in `lms-access-api` | done | [lms-access-api#10](https://github.com/code-corhuila/lms-access-api/pull/10) |
| W15-06 | Overdue penalty saga (ADR-008) in `lms-workflow`, promoted to `qa` and released as v3.0.0 | done | [lms-workflow#3](https://github.com/code-corhuila/lms-workflow/pull/3), [lms-workflow#7](https://github.com/code-corhuila/lms-workflow/pull/7) (replaces [#6](https://github.com/code-corhuila/lms-workflow/pull/6)) |
| W15-07 | Saga test suite and fix for the retry after a failed suspend | done | [lms-workflow#4](https://github.com/code-corhuila/lms-workflow/pull/4), [lms-workflow#5](https://github.com/code-corhuila/lms-workflow/pull/5) |
| HU-09 | Book editing: `PATCH /books/{id}` in `lms-catalog-api` (FR-011, FR-012) | doing | [lms-catalog-api#11](https://github.com/code-corhuila/lms-catalog-api/pull/11) |
| W15-01 | STRIDE threat model `lms-circulation-api → lms-catalog-api`: address review of #36 (merged 2026-11-13) | done | [library-docs#45](https://github.com/code-corhuila/library-docs/pull/45) |
| W15-02 | Revert of the threat model merged in #36 (closed without merging, #45 was kept) | done | [library-docs#44](https://github.com/code-corhuila/library-docs/pull/44) |
| W15-03 | Code review: `lms-api-gateway` release v1.0.0 (ADR-007) | done | [lms-api-gateway#7 (review)](https://github.com/code-corhuila/lms-api-gateway/pull/7#pullrequestreview-5896143786) |

## 2. My individual contribution

### `lms-access-portal` — HU-01 Vue login remote

- **[lms-access-portal#3](https://github.com/code-corhuila/lms-access-portal/pull/3) — `chore(qa): promote HU-01 Vue login remote to qa`** (`qa-hu-01-vue-login-portal` → `qa`, +2410/−25, 22 files, approved by `Hermesss29` and `OscarAreiza`, merged 2026-11-09):
  - Promotes HU-01 (the Vue login remote merged into `develop` in #2) by cutting the branch from `qa` and re-applying all 23 `develop` commits with `git cherry-pick -x`; no conflicts, and the resulting tree is identical to `develop` at `3673c95`.
  - The branch is `qa-hu-01-…` instead of `qa/hu-01-…` because Git cannot create `refs/heads/qa/...` while the `qa` branch exists.
  - Documents the `qa` gate gaps: no Pact tests in the repo yet, and GitHub Actions do not start while the organization account is locked over billing (`npm ci` and `npm run build` pass locally).
- **[lms-access-portal#4](https://github.com/code-corhuila/lms-access-portal/pull/4) — `chore(release): 0.1.0 — HU-01 Vue login remote`** (`release/0.1.0` → `main`, closed without merging): GitHub closed it when the branch was renamed to `release/v3.0.0`; replaced by #5 with the same commits.
- **[lms-access-portal#5](https://github.com/code-corhuila/lms-access-portal/pull/5) — `chore(release): v3.0.0 — HU-01 Vue login remote`** (`release/v3.0.0` → `main`, +2410/−25, approved by `ariel5253`, merged 2026-11-12, tagged `v3.0.0`):
  - First release of the portal: the login remote that `lms-front` mounts through Module Federation as `./mount` (`mountAccessLogin(el, onAuthenticated)`). Cut from `main`, each commit validated in `qa` re-applied with `cherry-pick -x`, with a trail to both `develop` (#2) and `qa` (#3).
  - Validation on `qa` (`2efe65b`): `npm ci`, `oxlint` (0 findings), `vue-tsc -b && vite build`, `npm audit --omit=dev` (0 runtime vulnerabilities), Docker image build, and **8 of 8 Vitest behaviour tests** covering HU-01's acceptance criteria (required/masked fields, existing session, successful login stores `accessToken`, loading state, generic `401` message per INV-002, fallback error, error reset, mount/unmount).
  - Contract check against `access-service.yaml`: the request matches `LoginRequest` and the view reads `accessToken` from `LoginResponse`.

### `lms-access-api` — Node.js/Express rewrite

- **[lms-access-api#9](https://github.com/code-corhuila/lms-access-api/pull/9) — `chore: promote access-service Node.js/Express rewrite to qa`** (`promote-qa/access-service-node-express-rewrite` → `qa`, +2342/−0, 25 files, approved by `OscarAreiza`, merged 2026-11-09):
  - Re-applies the 50 commits of the rewrite (#8) onto a branch cut from `qa` with `cherry-pick -x`.
  - The Go-only RS256/gitleaks fixes from #6/#7 were intentionally not re-applied, because the rewrite removes the Go code they touched. The tree matches `develop` except for 9 dead `.gitleaksignore` lines.
- **[lms-access-api#10](https://github.com/code-corhuila/lms-access-api/pull/10) — `fix: patch body-parser DoS vulnerability`** (`promote-qa/fix-body-parser-dos-vulnerability` → `qa`, +3/−27, merged 2026-11-09):
  - `npm audit` flagged a high-severity DoS in `body-parser` < 1.20.9 (transitive via `express`, GHSA-27r3-jxmh-p828). `npm audit fix` resolved it with no breaking changes.
  - `npm test` 6/6, `npm run lint` clean, `npm audit` 0 vulnerabilities.
- **[lms-access-api#11](https://github.com/code-corhuila/lms-access-api/pull/11) — `release: v3.0.0 — access-service Node.js/Express rewrite`** (`release/v3.0.0` → `main`, +2443/−0, 26 files, approved by `ariel5253`, merged 2026-11-12):
  - Brings #9 and #10 to `main` by re-applying the 28 `qa`-validated commits with `cherry-pick -x`; the tree is identical to `qa`.
  - `npm install`, `npm test` (6/6), `npm run lint` and `npm audit` (0 vulnerabilities) pass.

### `lms-workflow` — overdue penalty saga (ADR-008)

- **[lms-workflow#3](https://github.com/code-corhuila/lms-workflow/pull/3) — `qa: promote overdue penalty saga from develop`** (`promote-qa/overdue-penalty-saga` → `qa`, +1027/−1, 26 files, approved by `Hermesss29` and `OscarAreiza`, merged 2026-11-09):
  - Promotes the saga worker: poll overdue loans → mark late in `lms-circulation-api` → suspend the student in `lms-membership-api`, with durable, idempotent state in `workflow.saga_run`, plus the Liquibase migrations and deploy files. The tree is identical to `develop`.
- **[lms-workflow#4](https://github.com/code-corhuila/lms-workflow/pull/4) — `test(saga): add test suite and fix retry after suspend failure`** (`chore/saga-test-suite` → `develop`, +900/−3, 10 files, approved by `OscarAreiza`, merged 2026-11-09):
  - **Bug fixed:** when the suspend step failed, `run.Failed()` moved the run from `MARKED_LATE` to `FAILED`, so the next tick called `mark-late` again instead of resuming at the suspend step (against ADR-008's forward-only recovery). New `Run.RecordFailure(err)` counts the attempt and stores `LastError` without changing the status.
  - **Tests added:** `domain/saga` 100%, `application/usecase` 95.8% (including a regression test for the bug), `config` 100%, `infrastructure/httpclient` 87.0% (HS256 token with `sub=lms-workflow`, TTL ≤ 1 min), `infrastructure/scheduler` 100%, `infrastructure/postgres` 83.3% (integration tests, skipped without `WORKFLOW_TEST_DATABASE_URL`).
  - Verified with `go test -race` against PostgreSQL 17 as `workflow_app`, Liquibase `update` / no-op `update` / `rollback-count 4` / re-apply, Docker image build, and an end-to-end smoke test where the first suspend returns `500` and the run ends `COMPLETED` with `mark-late` called only once.
- **[lms-workflow#5](https://github.com/code-corhuila/lms-workflow/pull/5) — `qa: promote saga test suite and retry fix from develop`** (`promote-qa/saga-test-suite` → `qa`, +900/−3, approved by `OscarAreiza`, merged 2026-11-09): re-applies the merge of #4 with `cherry-pick -x -m 1`, so the retry fix reaches `qa` before any release. The tree is identical to `develop`.
- **[lms-workflow#6](https://github.com/code-corhuila/lms-workflow/pull/6) — `release: 0.1.0 overdue penalty saga`** (`release/0.1.0` → `main`, closed without merging): replaced by #7 after the release was renumbered to v3.0.0.
- **[lms-workflow#7](https://github.com/code-corhuila/lms-workflow/pull/7) — `release: v3.0.0 overdue penalty saga`** (`release/v3.0.0` → `main`, +1924/−1, 34 files, approved by `ariel5253`, merged 2026-11-12, tagged `v3.0.0`):
  - First release of `lms-workflow`: one commit per change validated in `qa` (migration scope, saga worker, test suite and retry fix), each with its trail from `develop` through `qa`. The tree is identical to `qa`.
  - On a 15-minute polling ticker: `GET /api/v1/loans/overdue` → `POST /api/v1/loans/{id}/mark-late` → `POST /api/v1/students/{id}/suspend`. Completed runs are skipped and a failed suspend resumes without repeating `mark-late`.

### `library-docs` — STRIDE threat model

- **[library-docs#45](https://github.com/code-corhuila/library-docs/pull/45) — Threat model: `lms-circulation-api → lms-catalog-api` — address review** (`security/stride-MENESESLUIS-bot` → `main`, +16/−3 on `05-architecture/security/stride-MENESESLUIS-bot.md`, approved by `Hermesss29` and `ariel5253`, merged 2026-11-13, refs #34, follow-up of #36):
  - Answers the review `Hermesss29` left on #36, which arrived after #36 was merged.
  - Row **T (b)** rewritten: `catalogclient.LoanCopy` does not retry, and `RegisterLoan` creates the `Loan` *after* `loan-copy`, so `loanId` cannot be the idempotency key. The real double decrement is an orphan decrement followed by a new `POST /api/v1/loans` with a new `Idempotency-Key`. Mitigation: forward the original `POST /loans` key to `loan-copy`, plus a `return-copy` compensation when `Save` fails.
  - Row **R**: daily orphan-copy check in `lms-worker` (`total_copies − available_copies` vs. active loans in `loan_db.loans`).
  - New residual-risk paragraph for the new-key case until the compensation exists, and a **Review log** section with each comment and what was done. Risk values unchanged (T = 9, R = 6).
- **[library-docs#44](https://github.com/code-corhuila/library-docs/pull/44) — Revert "Threat model: `lms-circulation-api → lms-catalog-api` (internal loan-copy / return-copy)"** (`revert-36-security/stride-MENESESLUIS-bot` → `main`, +0/−98, approved by `ariel5253`): would have reverted #36 by removing `stride-MENESESLUIS-bot.md`. Closed without merging on 2026-11-13, because #45 was kept.

### `lms-catalog-api` — HU-09 book editing

- **[lms-catalog-api#11](https://github.com/code-corhuila/lms-catalog-api/pull/11) — `fix(catalog): add PATCH /books/{id} for HU-09 book editing`** (`fix/book-update-endpoint` → `develop`, +188/−4, approved by `ariel5253`, `OscarAreiza` and `Hermesss29`, not merged yet):
  - Problem: editing a book in `lms-catalog-portal` always failed with *"Something went wrong. Please try again."* The portal sends `PATCH /books/{id}`, but the service never registered that route, so chi answered `404`. The domain already had `Book.Update` (HU-09), but nothing called it.
  - Ports HU-09 from the monolith's `feat/HU-09-book-editing` to the hexagonal layout: `UpdateBookUseCase` input port, `usecase.UpdateBook` (`FindByID` → `Book.Update` → `Save`), `BookHandler.Update`, the `books.Patch("/{id}", …)` route and its wiring in `cmd/api/main.go`.
  - Title and author are validated at the boundary; the `400 VALIDATION_ERROR` lists each invalid field in `details`, which the portal's edit row reads to highlight that input. Other responses: `200` with the updated book, `404 NOT_FOUND`, `500 INTERNAL_ERROR`.
  - The ISBN is not part of the request (FR-012) and copy counts are not touched, so existing loans are unaffected (FR-011).
  - Tests: new `update_book_test.go` (successful edit is saved; ISBN and copies 2/3 are kept; unknown id returns `ErrBookNotFound`; empty title is rejected and leaves the stored book untouched). `gofmt -l .`, `go vet ./...` and `go test ./...` pass in `golang:1.25-alpine`.
  - End to end against a local stack (catalog-db with Liquibase migrations, access-api for the RS256 token): valid `PATCH` → `200` (ISBN unchanged even when the body sends `"isbn"`, copies stay 2/2); empty title/author → `400` with `details`; unknown id → `404`.

### Code review

- **Review of [lms-api-gateway#7](https://github.com/code-corhuila/lms-api-gateway/pull/7)** — `release: v1.0.0 — ADR-007 gateway implementation` (author `OscarAreiza`, `release/v1.0.0` → `main`): checked that it is a single `cherry-pick -x` of `qa`, and that `deploy/nginx.conf` matches ADR-007: path-prefix routing to access, membership, catalog and circulation, the `Authorization`-header presence check with a fast `401` (no JWT parsing at the gateway, ADR-003), and the two `limit_req` zones (5 r/m on login, 60 r/m on the rest of the API, `429`). Approved, with a non-blocking note: `add_header Retry-After 60 always;` is sent on every response of those locations, not only on `429`s.

## 3. Blockers and risks

- `POST /api/v1/loans/{id}/mark-late` does not exist in `lms-circulation-api` yet, so the saga released in `lms-workflow` v3.0.0 fails at step one and its runs end `FAILED`; it will not apply penalties end to end until that route is merged (open in [lms-circulation-api#18](https://github.com/code-corhuila/lms-circulation-api/pull/18)).
- No Pact contract tests exist for `lms-access-portal` or for the saga's calls to `lms-circulation-api` and `lms-membership-api`, so the `qa` gate in `branching-policy.md` was covered by local contract checks, not by Pact.
- GitHub Actions do not start while the organization account is locked over billing; every check reported in this week's PRs was run locally.
- `lms-front` still renders its own `LoginPage` at `/login` and does not load the `access_portal` remote yet, so HU-01's portal is released but not wired into the shell.
- The saga's service token has no `aud` claim; once `lms-membership-api` enforces `aud` (library-docs #39 / #42), the suspend calls will be rejected.
- `lms-access-api` v3.0.0 is merged to `main` but the `v3.0.0` tag has not been created yet (`lms-access-portal` and `lms-workflow` are already tagged).
- [lms-catalog-api#11](https://github.com/code-corhuila/lms-catalog-api/pull/11) has been open since 2026-10-09 with three approvals but is not merged; until it is, editing a book in `lms-catalog-portal` keeps failing with `404`.
- [lms-api-gateway#7](https://github.com/code-corhuila/lms-api-gateway/pull/7) has no CI checks and still needs the mandatory approval of `ariel5253` before it can be merged to `main`.
- The `return-copy` compensation proposed in #45 does not exist yet, so the orphan-decrement risk (row T) is still open in the code.

## 4. Plan for next week

- Create the `v3.0.0` tag on `lms-access-api`'s `main`, as `git-conventions.md` requires.
- Follow [lms-circulation-api#18](https://github.com/code-corhuila/lms-circulation-api/pull/18) (`mark-late`) and run the saga end to end once it is merged.
- Add Pact contract tests for the saga's calls and for the `lms-access-portal` login.
- Wire the `access_portal` remote into `lms-front`.
- Get #11 merged into `develop` and promote it to `qa` with `cherry-pick -x`.
- Implement the `return-copy` compensation and the forwarding of the `Idempotency-Key` to `loan-copy` in `lms-circulation-api`, as described in #45.
- Follow up on the `Retry-After` note in `lms-api-gateway` once v1.0.0 is on `main`.

## 5. Compliance self-check

- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links

- Folder: [15-week/hu-status](./)
- My Pull Requests this week:
  - `lms-access-portal`: [#3](https://github.com/code-corhuila/lms-access-portal/pull/3), [#4](https://github.com/code-corhuila/lms-access-portal/pull/4) (closed), [#5](https://github.com/code-corhuila/lms-access-portal/pull/5)
  - `lms-access-api`: [#9](https://github.com/code-corhuila/lms-access-api/pull/9), [#10](https://github.com/code-corhuila/lms-access-api/pull/10), [#11](https://github.com/code-corhuila/lms-access-api/pull/11)
  - `lms-workflow`: [#3](https://github.com/code-corhuila/lms-workflow/pull/3), [#4](https://github.com/code-corhuila/lms-workflow/pull/4), [#5](https://github.com/code-corhuila/lms-workflow/pull/5), [#6](https://github.com/code-corhuila/lms-workflow/pull/6) (closed), [#7](https://github.com/code-corhuila/lms-workflow/pull/7)
- Other Pull Requests of mine tracked this week: [lms-catalog-api#11](https://github.com/code-corhuila/lms-catalog-api/pull/11), [library-docs#45](https://github.com/code-corhuila/library-docs/pull/45), [library-docs#44](https://github.com/code-corhuila/library-docs/pull/44)
- Review: [lms-api-gateway#7](https://github.com/code-corhuila/lms-api-gateway/pull/7#pullrequestreview-5896143786)
- Related: [library-docs#36](https://github.com/code-corhuila/library-docs/pull/36) (original threat model), [library-docs#34](https://github.com/code-corhuila/library-docs/issues/34)
