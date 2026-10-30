<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Write the first `node --test` tests for `lms-access-api`: the `Administrator` aggregate of the Access domain (password hashing and authentication) and the `Login` use case, driven through fake ports so the domain and application layers are tested with no database or JWT I/O; and a multi-stage `Dockerfile` to package the service; then the hexagonal source for HU-01 (`POST /api/v1/auth/login`) that makes those tests pass. W13 workshop: threat model of the trust boundary `lms-circulation-api` → `lms-catalog-api` (HU-06 `loan-copy`, HU-07 `return-copy`) and turn it into security backlog and CI gates.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W13-01 | `Administrator` domain tests: no plaintext password, authenticate ok / wrong password | done | [83d426c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/83d426c) |
| W13-02 | `Login` use case tests with fake `administrators` and `tokens` ports | done | [dfe630b](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/dfe630b) |
| W13-03 | Multi-stage `Dockerfile` (`node:22-alpine`, production deps only, port 8080) | done | [d22c456](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d22c456) |
| W13-04 | Access domain (`Administrator`, errors) and `Login` use case | done | [0bb0a25](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0bb0a25) |
| W13-05 | Config from env vars, pino logger and PostgreSQL adapter (`access.administrator`) | done | [483a367](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/483a367) |
| W13-06 | RS256 JWT issuer, `/auth/login` and `/health` handlers | done | [577b783](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/577b783) |
| W13-07 | Middleware (correlation ID, CORS, error envelope), router and `src/index.js` with graceful shutdown | done | [ba33b5f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ba33b5f) |
| W13-TM-CP0 | Claim trust boundary `lms-circulation-api` → `lms-catalog-api` | done | [library-docs#34](https://github.com/code-corhuila/library-docs/issues/34), [0f6d072](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0f6d072) |
| W13-TM-CP1 | Attack surface map and trust boundaries | done | [be13d37](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/be13d37) |
| W13-TM-CP2 | STRIDE on the boundary (PR in `library-docs`) | done | [library-docs#36](https://github.com/code-corhuila/library-docs/pull/36) |
| W13-TM-CP3 | Three `[security]` stories in the team backlog | done | [#37](https://github.com/code-corhuila/library-docs/issues/37), [#38](https://github.com/code-corhuila/library-docs/issues/38), [#39](https://github.com/code-corhuila/library-docs/issues/39), [7567c76](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7567c76) |
| W13-TM-CP4 | Security gates workflow and negative authorization test | done | [8f4eb57](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8f4eb57) |
| W13-TM-CP5 | Peer review of the circulation → membership STRIDE | done | [library-docs#40](https://github.com/code-corhuila/library-docs/pull/40), [7912ed7](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7912ed7) |
| W13-TM-CP6 | Answer the review and finish the STRIDE | done | [library-docs#45](https://github.com/code-corhuila/library-docs/pull/45), [f043c90](https://github.com/code-corhuila/library-docs/commit/f043c90) |
| W13-TM-CP7 | Refine the security backlog | done | [318ac3f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/318ac3f) |
| W13-TM-CP8 | Run the negative test and the workflow | done | [e9e7f04](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e9e7f04), [Actions run](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/actions/runs/41092565678) |

## 2. My individual contribution
- `test/domain/administrator.test.js` (3 tests): `newAdministrator('alice', ...)` never keeps the plaintext password and stores a bcrypt hash (`$2...` prefix); `authenticate` does not throw on the correct password; `authenticate` throws `ErrInvalidCredentials` on a wrong one.
- `test/application/login.test.js` (3 tests): the `Login` use case receives its ports by constructor (`administrators.findByUsername`, `tokens.issue`), replaced here by in-memory fakes, so no PostgreSQL or JWT key is needed.
  - Happy path returns a non-empty `accessToken` and `expiresIn = 3600`.
  - Unknown username (`ErrAdministratorNotFound` from the repository) and wrong password both reject with the same `ErrInvalidCredentials`, so the API does not reveal which usernames exist.
- Both files use only the built-in `node:test` and `node:assert/strict`, matching the `test` script (`node --test`) from the week 12 `package.json`.
- `deploy/Dockerfile`: two stages on `node:22-alpine` (same Node major as `engines` in `package.json`).
  - Build stage copies `package.json` (and `package-lock.json` if present) first and runs `npm install --omit=dev`, so the dependency layer is cached while only `src/` changes; then copies `src/`.
  - Run stage adds `ca-certificates` (TLS to PostgreSQL and other services) and copies only `node_modules`, `package.json` and `src/` from the builder, so tests, dev tools and the npm cache stay out of the image.
  - `EXPOSE 8080` matches `PORT=8080` in `.env.example`; `CMD ["node", "src/index.js"]` in exec form so Node receives `SIGTERM` directly on `docker stop`.
- `src/` (16 files, hexagonal layout; the 6 tests above pass with `node --test`):
  - Domain (`src/domain/access/`): `Administrator` aggregate with `authenticate()`; `newAdministrator()` hashes with bcrypt cost 12 (NFR-004) and never keeps the plaintext (INV-001); `ErrInvalidCredentials` and `ErrAdministratorNotFound`. No I/O besides `bcryptjs`.
  - Application (`src/application/usecase/login.js`): `Login` gets its ports by constructor and maps "not found" and "wrong password" to the same `ErrInvalidCredentials` (INV-002).
  - Driven adapters: `config/config.js` reads every variable of `.env.example`, fails fast without `JWT_PRIVATE_KEY` and turns the literal `\n` back into a PEM; `postgres/db.js` opens a `pg` pool and checks it with `SELECT 1`; `postgres/administratorRepository.js` uses parameterized queries on `access.administrator` (schema per domain, singular table); `auth/jwtIssuer.js` signs RS256 tokens with only `sub` (no role claim) and parses `JWT_EXPIRY` (`1h`, `30m`, `15s`); `logger/logger.js` is pino.
  - Driving adapters (`src/infrastructure/http/`): `POST /api/v1/auth/login` (400 on bad body, 401 `INVALID_CREDENTIALS`, 200 with `accessToken` + `expiresIn`), `GET /health` and `GET /health/ready` (DB ping with a 2 s bound, 503 when down), `X-Correlation-Id` middleware, CORS, and an error handler that returns the standard `{error, message, correlationId}` envelope.
  - `src/index.js` wires everything and closes the HTTP server and the pool on `SIGINT`/`SIGTERM`, with a 10 s forced exit.
- **W13 threat modeling workshop** (`13-week/workshop/`):
  - `01-attack-surface.md`: Mermaid data flow of all LMS services (`lms-front`, `lms-access-api`, `lms-membership-api`, `lms-catalog-api`, `lms-circulation-api`, their databases, planned `lms-api-gateway`, `lms-broker`, `lms-iam`), trust boundaries TB1–TB6 with mine (TB4) highlighted, asset inventory and the entry points of TB4.
  - PR library-docs#36 (`05-architecture/security/stride-MENESESLUIS-bot.md`): six STRIDE rows for the internal `loan-copy` / `return-copy` calls; top risks are E (no per-route scope, 9), T (lost/double decrement of `available_copies`, 9) and S (example `INTERNAL_JWT_SECRET`, 6).
  - Issues #37 (P0), #38 (P0), #39 (P1) with user story, testable criteria (one automated test each) and the link to the PR.
  - `04-security-ci/`: `security.yml` (gitleaks, govulncheck, authz tests; `actionlint` clean) and a minimal `lms-catalog-api` module with `RequireScope` and 4 negative tests (internal token on `DELETE /api/v1/books/{id}` → 403, Administrator token on `loan-copy` → 403), all passing.
  - Review answered (PR library-docs#45): row T (b) rewritten after @Hermesss29's comment, because in `RegisterLoan` the `Loan` is created after `loan-copy`, so the key is the `Idempotency-Key` of the original `POST /api/v1/loans`, plus a `return-copy` compensation; row R gets a daily orphan-copy check; new residual-risk paragraph and a review log.
  - Issue #38 rewritten with those criteria; #37 linked to #41 and #39 linked to #42 (Hermes's stories).
  - Workflow executed on GitHub Actions from `.github/workflows/security.yml` of this fork: negative tests and `govulncheck` pass, `gitleaks` fails on one true positive. Removing `RequireScope` in a temporary copy makes both negative tests fail with 200 instead of 403, so the tests really protect the rule.

## 3. Blockers and risks
- `package.json`, `.env.example` and `.dockerignore` live in `12-week/hu-status`, so `npm test` and `docker build` cannot run from this folder; the tests were checked by copying `src/`, `test/` and that `package.json` together.
- `config.js` falls back to `lms_user` / `lms_password` and `CORS_ORIGIN=*` when the variables are missing, so a misconfigured deploy starts with weak defaults instead of failing.
- An unknown username returns without running bcrypt, while a wrong password takes about 250 ms (cost 12), so response time can still reveal which usernames exist despite INV-002; `compareSync` also blocks the event loop.
- No HTTP or repository tests yet, and there is no migration that creates `access.administrator` or seeds the first Administrator.
- The run stage has no `USER node`, so the container runs as root.
- `package-lock.json` is gitignored, so the image uses `npm install` instead of `npm ci` and dependency versions are not pinned between builds.
- `13-week/01-session` and `13-week/02-session` only contain `.gitkeep`.
- The real `RequireAuth` of `lms-catalog-api` accepts the internal HS256 token on every `/books` route, and both `.env.example` files ship `INTERNAL_JWT_SECRET=change-me-internal-secret` (STRIDE rows E and S).
- The library book (García Serrano, *Inteligencia artificial*, 2012, 006.3 G17i EJ.1, p. 171) is cited in row E of the STRIDE ([57c7019](https://github.com/code-corhuila/library-docs/commit/57c7019)) and in issue #37; it is not a security book, so its use is limited to the rule-based reasoning idea.
- `gitleaks` found a true positive: `JWT_PRIVATE_KEY` in `12-week/hu-status/.env.example` of this public fork is the pair of the `JWT_PUBLIC_KEY` in the `.env.example` of `lms-catalog-api`, `lms-circulation-api` and `lms-membership-api`. It stays in the history (no rewrites), so a new pair must be generated for every shared environment. I could not verify which key `lms-pod` uses.
- `govulncheck` shows 8 reachable standard library vulnerabilities on Go 1.25.1, the version of the team's Go services; they need Go 1.26.9.
- PR library-docs#36 was merged by its own author, against rule 8 of the workshop. Revert PR library-docs#44 is open and undecided; merging it would remove the STRIDE from `main`.
- My review on PR library-docs#40 was submitted as an approval, while CP5 asked for a comment only.

## 4. Plan for next week
- Put `package.json`, `.env.example` and `.dockerignore` next to `src/` so `npm test`, `npm run lint` and `docker build` run from one folder.
- Compare against a dummy hash when the username is unknown and switch to async `bcrypt.compare`; remove the DB credential fallbacks.
- Add an integration test for `POST /api/v1/auth/login` (200 / 400 / 401) and the `access.administrator` migration with a seed Administrator.
- Add `USER node` to the run stage, a `HEALTHCHECK`, and switch to `npm ci` once `package-lock.json` is committed; then build and run the image with `--env-file .env`.
- Rotate the RS256 key pair and replace the example values in the team `.env.example` files.
- Start with story #37 (per-route scopes) in `lms-catalog-api` and move `security.yml` to the `.github/workflows/` of that repository, on Go 1.26.9.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [ ] No secrets; config via environment variables (`gitleaks` finding, see section 3)

## 6. Evidence links
- Folder: [13-week/hu-status](./)
- Files: [test/domain/administrator.test.js](test/domain/administrator.test.js), [test/application/login.test.js](test/application/login.test.js), [deploy/Dockerfile](deploy/Dockerfile), [src/](src/) ([index.js](src/index.js), [login.js](src/application/usecase/login.js), [administrator.js](src/domain/access/administrator.js), [router.js](src/infrastructure/http/router.js))
- Commits: [83d426c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/83d426c), [dfe630b](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/dfe630b), [68bce61](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/68bce61), [d22c456](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d22c456), [0bb0a25](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0bb0a25), [483a367](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/483a367), [577b783](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/577b783), [ba33b5f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ba33b5f)
- Workshop: [13-week/workshop](../workshop/) · claim [library-docs#34](https://github.com/code-corhuila/library-docs/issues/34) · PR [library-docs#36](https://github.com/code-corhuila/library-docs/pull/36) · stories [#37](https://github.com/code-corhuila/library-docs/issues/37), [#38](https://github.com/code-corhuila/library-docs/issues/38), [#39](https://github.com/code-corhuila/library-docs/issues/39)
- Workshop commits: [0f6d072](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0f6d072) (CP0), [be13d37](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/be13d37) (CP1), library-docs [a953716](https://github.com/code-corhuila/library-docs/commit/a953716) / [030f0d6](https://github.com/code-corhuila/library-docs/commit/030f0d6) (CP2), [7567c76](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7567c76) (CP3), [8f4eb57](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8f4eb57) (CP4), [7912ed7](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7912ed7) (CP5)
- Delivery 2 commits: library-docs [f043c90](https://github.com/code-corhuila/library-docs/commit/f043c90) (CP6), [318ac3f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/318ac3f) (CP7), [78dd941](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/78dd941) / [2127a4c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/2127a4c) / [e9e7f04](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e9e7f04) (CP8)
- Peer review: comment on Hermes's STRIDE PR [library-docs#40](https://github.com/code-corhuila/library-docs/pull/40) (missing Tampering threat: `Student.Suspend` overwrites `suspended_until` and `/suspend` has no idempotency key; row R `actor` is always `circulation-service`)
