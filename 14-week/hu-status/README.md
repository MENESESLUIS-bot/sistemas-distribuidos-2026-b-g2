<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       14-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 14

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Continue `lms-access-portal`, the Vue login remote that `lms-front` loads through Module Federation: add the standalone entry point (`src/main.ts`) and `App.vue` wrapper so the portal can run outside the shell, the Vite/Vue type declarations, the ambient types of the `shell` remotes and the HU-01 login view (`LoginView.vue`), the `useSession` composable, the federated `mountAccessLogin` (`mount.ts`) and the standalone router, and a first production build (`dist/`) that publishes `remoteEntry.js` with the `mountAccessLogin` export; plus the nginx config and compose file that serve that build (committed under `13-week/hu-status/deploy/`). Start `lms-workflow`, the Go worker for the overdue → penalty saga (ADR-008): module, configuration, README with the saga scope, the hexagonal source (domain, `ApplyPenalty`, HTTP clients, PostgreSQL repository, ticker), the `workflow` schema migrations and the Docker/Liquibase deploy files.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W14-01 | Standalone entry (`src/main.ts`), `App.vue` wrapper and `src/env.d.ts` | done | [d4b4a53](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d4b4a53) |
| W14-02 | nginx config (no-cache `remoteEntry.js`, immutable `/assets/`) and `compose.yml` for the portal | done | [deb0b54](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/deb0b54) |
| W14-03 | Production build of the federated remote (`dist/`, not committed) | done | [.gitignore](.gitignore) |
| W14-04 | Ambient types for `shell/apiClient` and `shell/session` (`src/shims/shell.d.ts`) | done | [ddc8b74](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ddc8b74) |
| W14-05 | HU-01 login view (`src/views/LoginView.vue`): form, submit to `/auth/login`, generic error | done | [8b96202](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8b96202), [edad5ef](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/edad5ef) |
| W14-06 | `useSession` composable that wraps `shell/session` (`src/composables/useSession.ts`) | done | [c0db426](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/c0db426) |
| W14-07 | Federated `mountAccessLogin` (`src/mount.ts`) and standalone router (`src/router.ts`) | done | [e309165](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e309165) |
| W14-08 | Portal manifest, `vite.config.ts` and tsconfig in this folder so it builds here | todo |  |
| W14-09 | `lms-workflow` Go module, `.env.example`, `.gitignore` and `Makefile` | done | [de467c9](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/de467c9) |
| W14-10 | `lms-workflow` README: two-step saga scope, `workflow.saga_run` state, migration scope | done | [3a35b8c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3a35b8c) |
| W14-11 | Saga `Run` aggregate (`PENDING → MARKED_LATE → COMPLETED` / `FAILED`) and driven ports | done | [aaf73b1](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/aaf73b1) |
| W14-12 | `ApplyPenalty` use case: idempotent by `loan_id`, resume at `MARKED_LATE`, forward recovery | done | [7b31c3f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7b31c3f) |
| W14-13 | Config from env vars (fails fast without `INTERNAL_JWT_SECRET`) | done | [76c9d44](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/76c9d44) |
| W14-14 | HTTP clients for circulation and membership with a 1-minute HS256 service token; zap logger | done | [aecb1e4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/aecb1e4), [12e2a6a](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/12e2a6a) |
| W14-15 | PostgreSQL pool and `RunRepository` (upsert on `workflow.saga_run`) | done | [794ee65](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/794ee65) |
| W14-16 | Ticker scheduler and `cmd/worker/main.go` with graceful shutdown | done | [0e8cdfa](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0e8cdfa) |
| W14-17 | Liquibase migrations: `workflow` schema, `saga_run`, status index, NOLOGIN roles | done | [b9c3091](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/b9c3091) |
| W14-18 | Multi-stage Go `Dockerfile` and `workflow-migrate` compose service | done | [e9f108c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e9f108c) |
| W14-19 | `lms-workflow` unit tests (`ApplyPenalty` with fake ports) | todo |  |

## 2. My individual contribution
- `src/main.ts`: standalone-dev entry; creates the Vue app with `App` and `router` and mounts it on `#app`. It is only used by `npm run dev` outside the shell; when `lms-front` composes the portal it loads `mount.ts` instead.
- `src/App.vue`: minimal wrapper with a single `<router-view />`, so the portal can be previewed in isolation.
- `src/env.d.ts`: `vite/client` types and a `*.vue` module declaration, so `vue-tsc` type-checks the `.vue` imports.
- `src/shims/shell.d.ts`: ambient modules for what `lms-front` exposes at runtime: `shell/apiClient` (a minimal `post<T>()` interface, so the portal does not need `axios` just for a type, and the `ShellError` envelope `{error, message, details?, traceId?}`) and `shell/session` (`getToken`, `setToken`, `isAuthenticated`, `clearSession`).
- `src/views/LoginView.vue` (HU-01, Vue port of the former React `LoginPage.tsx`, ADR-011):
  - Receives `onAuthenticated` as a prop and calls it right away if a session already exists.
  - Username and password fields (`required`, proper `autocomplete`); on submit it posts to `/auth/login` through the shell `apiClient`, stores `accessToken` with `setToken` and calls `onAuthenticated`.
  - The submit button is disabled while the request runs; on failure it shows the API `message` in a `role="alert"` box, without guessing whether the username or the password was wrong (INV-002).
  - Scoped BEM styles (`access-login__*`) so the portal's CSS does not leak into the React shell.
- `src/composables/useSession.ts`: the only file that imports `shell/session`; it returns `getToken`, `setToken`, `isAuthenticated` and `clearSession`, so the rest of the portal never depends on the federated module directly.
- `src/mount.ts`: exposed to `lms-front` through Module Federation. A Vue remote cannot export a React `RouteObject[]` like the other three portals, so it exports `mountAccessLogin(el, onAuthenticated)`, which creates a Vue app with `LoginView`, mounts it on `el` and returns an unmount function (ADR-011).
- `src/router.ts`: standalone-dev routing only, used by `main.ts`; route `/` renders `LoginView` with an `onAuthenticated` that just logs `[standalone] authenticated`. Without `lms-front` running on port 3000 the `shell/*` imports cannot resolve, so the login call fails in this mode.
- Build output (`dist/`, generated locally, ignored by git):
  - `remoteEntry.js` is the federated entry of `access_portal`; `@mf-types/compiled-types/src/mount.d.ts` publishes `mountAccessLogin(el, onAuthenticated): () => void` for the shell.
  - `assets/` holds content-hashed chunks (`LoginView-*.js/.css`, shared-import map, `loadRemote` stubs for `shell/apiClient` and `shell/session`).
  - `index.html` is titled `LMS — Access Portal` and preloads the shell remote from `http://localhost:3000/remoteEntry.js`.
- `.gitignore`: `node_modules`, `dist`, the `.mf` cache, `.env` and `*.pem`, so the build is not committed.
- `13-week/hu-status/deploy/nginx.conf`: `remoteEntry.js` with `Cache-Control: no-store` (H-front rule 4, so a deployed fix always reaches the browser), `/assets/` cached for one year as `immutable`, CORS `*` on both, SPA fallback to `index.html`, and a local-only `/api/v1/` proxy to `access-service:8080` while `lms-api-gateway` does not exist.
- `13-week/hu-status/deploy/compose.yml`: service `access-portal` on port `3004:80`, on the external `lms-network`; it runs the portal alone and is not the `lms-infra` composition.
- `lms-workflow/` (new Go service, `github.com/code-corhuila/lms-workflow`):
  - `go.mod` on Go 1.25.1 with `golang-jwt/jwt/v5` (internal HS256 token), `jackc/pgx/v5` (PostgreSQL) and `go.uber.org/zap` (structured logs); `go.sum` pins the versions.
  - `.env.example`: DB connection, `INTERNAL_JWT_SECRET` (service-to-service, not the end-user RS256 token), `CIRCULATION_SERVICE_URL`, `MEMBERSHIP_SERVICE_URL`, `POLL_INTERVAL=15m` (same cadence as ADR-009) and `LOG_LEVEL`.
  - `Makefile`: `dev` (`go run ./cmd/worker/...`), `test`, `test-cover` (HTML coverage), `build` (`bin/worker`) and `lint` (`golangci-lint`).
  - `.gitignore`: binaries, `coverage.out`, `/vendor/`, `.env`, `*.pem`, Liquibase local files and IDE folders.
  - `README.md`: per ADR-008 the saga only has two steps once a loan is overdue: mark the loan late in `lms-circulation-api`, then suspend the student in `lms-membership-api`. It is orchestrated with no broker (AT-002), driven by a ticker instead of an HTTP server. State lives in `workflow.saga_run` (`PENDING → MARKED_LATE → COMPLETED` or `FAILED`), with `loan_id` as the idempotency key and forward recovery only. The `workflow` schema lives on the shared PostgreSQL instance, applied by a Liquibase executor in this repo. It also documents the branching policy (`develop` / `qa` / `main`, promotion by `cherry-pick -x`).
  - Domain (`internal/domain/saga/`): `Run` aggregate with `NewRun`, `MarkedLate`, `Completed` and `Failed` (increments `AttemptCount`, keeps `LastError`); ports `RunRepository`, `CirculationClient` (`FindOverdueLoans`, `MarkLate`), `MembershipClient` (`Suspend`) and `ErrRunNotFound`. No I/O.
  - Application (`internal/application/usecase/apply_penalty.go`): loads the run by `loan_id` or creates it; skips it if `COMPLETED`; runs `MarkLate` only if not already `MARKED_LATE`; then `Suspend`; never rolls back step one.
  - `internal/config/config.go`: reads every variable of `.env.example`, parses `POLL_INTERVAL` as a Go duration and refuses to start without `INTERNAL_JWT_SECRET`.
  - `internal/infrastructure/httpclient/`: `GET /api/v1/loans/overdue` (maps `id` → `LoanID`), `POST /api/v1/loans/{id}/mark-late` and `POST /api/v1/students/{id}/suspend`, each with a 5 s timeout and a fresh HS256 token (`sub=lms-workflow`, 1 min expiry).
  - `internal/infrastructure/postgres/`: pgx pool checked with a 5 s ping; `RunRepository` with parameterized queries and `INSERT ... ON CONFLICT (loan_id) DO UPDATE`.
  - `internal/infrastructure/scheduler/ticker.go` runs one tick at start and then every `POLL_INTERVAL`; `cmd/worker/main.go` wires everything and stops on `SIGINT`/`SIGTERM`.
  - `migrations/`: flat `changelog.yaml` with inline rollbacks (changesets carry `author: OscarAreiza`); `workflow` schema, `saga_run` with `loan_id text` (MongoDB ObjectId, ADR-005) and a `CHECK` on status, `idx_saga_run_status`, and `workflow_reader` / `workflow_writer` NOLOGIN roles with the writer granted to `workflow_app`.
  - `deploy/`: two-stage `Dockerfile` (`golang:1.25-alpine` → `alpine:3.21`, static binary, no `EXPOSE`); `compose.yml` with only the `workflow-migrate` Liquibase 4.29 service on the external `platform` network; `liquibase.properties` with per-domain changelog table names.
  - Checked locally: `go build ./...` and `go vet ./...` pass.

## 3. Blockers and risks
- The portal source is now complete in `src/`, but it still cannot be built from this folder (see the next point), so `dist/` was produced with the week 13 configuration and has not been rebuilt from this source.
- The standalone mode (`npm run dev` without the shell) cannot really log in, because `shell/apiClient` and `shell/session` only resolve when `lms-front` runs on port 3000.
- If the shell `apiClient` rejects with a raw Axios error instead of the unwrapped `ShellError`, the view shows Axios's generic message instead of the API one.
- This folder has no `package.json`, `vite.config.ts`, `index.html` or tsconfig files; they are in `13-week/hu-status`, so `npm run dev` / `npm run build` cannot run from `14-week/hu-status`.
- `compose.yml` builds `deploy/Dockerfile`, but that Dockerfile is the Node image of `lms-access-api` (`EXPOSE 8080`, `node src/index.js`) and does not copy `dist/` into nginx, while compose maps port `80`. The portal needs its own nginx Dockerfile.
- The shell URL is still fixed to `http://localhost:3000/remoteEntry.js` in the build, so it needs a per-environment value before any deploy outside local development.
- CORS `*` on `remoteEntry.js` and `/assets/` is acceptable for static files but should be narrowed to the `lms-front` origin once it is known.
- `lms-workflow` has no tests: `go test ./...` finds no test files.
- `Run.Failed()` overwrites the status with `FAILED`, so a run that failed at `Suspend` loses `MARKED_LATE`; on the next tick `ApplyPenalty` calls `MarkLate` again, against the "resume at the suspend step" rule of ADR-008 (it is only safe if `mark-late` is idempotent).
- `FAILED` runs are retried on every tick with no limit on `AttemptCount`.
- A new run is saved as `PENDING` before step one, while the comment in `run.go` says a row is only inserted after step one.
- `FindOverdueLoans` reads only the first page of `/loans/overdue`; `MarkLate` and `Suspend` accept only `200`, so a `204` would be treated as a failure.
- `004_create_roles_and_grants.sql` grants to `workflow_app`, which must already exist (created by `lms-infra-postgres`); `compose.yml` defaults `PG_ADMIN_PASSWORD` to `change-me-admin-password`, and `config.go` builds the DSN with `sslmode=disable`.
- The Dockerfile has no `USER`, so the worker runs as root.
- The saga needs `POST /api/v1/loans/{id}/mark-late` on `lms-circulation-api`, which is not built yet; only `GET /api/v1/loans/overdue` (HU-08) exists.
- `go.sum` still lists `go-chi/chi/v5`, which `go.mod` does not require; `go mod tidy` should clean it.
- `lms-workflow` is on Go 1.25.1, the same version where `govulncheck` found 8 reachable standard library vulnerabilities in week 13; it needs Go 1.26.9.
- `.env.example` ships `INTERNAL_JWT_SECRET=change-me-internal-secret` and `lms_password`, the same example values flagged in the W13 STRIDE (row S).
- ADR-006 still says workflow is stateless, which contradicts ADR-008; the README flags it but the ADR text is not corrected.
- `14-week/01-session` and `14-week/02-session` are empty.

## 4. Plan for next week
- Bring the portal manifest, `vite.config.ts`, `index.html` and tsconfig files next to `src/` so `npm run build` reproduces `dist/` from this source.
- Write an nginx multi-stage Dockerfile for the portal (`node:22-alpine` build → `nginx:alpine` with `nginx.conf` and `dist/`) and point `compose.yml` to it.
- Make the shell `remoteEntry.js` URL configurable per environment.
- Run the login view end to end against `POST /api/v1/auth/login` of `lms-access-api` and add a test for `mountAccessLogin` (mount, `onAuthenticated` callback, unmount).
- Add `ApplyPenalty` unit tests with fake ports: new run, skip if `COMPLETED`, resume at `MARKED_LATE`, failure in each step.
- Keep the last reached step when a run fails (for example a separate `FailedAt` field) so a retry after a failed `Suspend` does not call `MarkLate` again, and add a maximum `AttemptCount`.
- Handle pagination of `/loans/overdue`, add `USER` to the Dockerfile, and run the migrations once with `docker compose` against the shared PostgreSQL.
- Run `go mod tidy` and move the module to Go 1.26.9.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables (`.env` and `*.pem` ignored)

## 6. Evidence links
- Folder: [14-week/hu-status](./)
- Files: [src/main.ts](src/main.ts), [src/App.vue](src/App.vue), [src/env.d.ts](src/env.d.ts), [src/shims/shell.d.ts](src/shims/shell.d.ts), [src/views/LoginView.vue](src/views/LoginView.vue), [src/composables/useSession.ts](src/composables/useSession.ts), [src/mount.ts](src/mount.ts), [src/router.ts](src/router.ts), [.gitignore](.gitignore), [lms-workflow/README.md](lms-workflow/README.md), [lms-workflow/go.mod](lms-workflow/go.mod), [lms-workflow/.env.example](lms-workflow/.env.example), [lms-workflow/Makefile](lms-workflow/Makefile), [run.go](lms-workflow/internal/domain/saga/run.go), [apply_penalty.go](lms-workflow/internal/application/usecase/apply_penalty.go), [main.go](lms-workflow/cmd/worker/main.go), [changelog.yaml](lms-workflow/migrations/changelog.yaml), [Dockerfile](lms-workflow/deploy/Dockerfile), [nginx.conf](../../13-week/hu-status/deploy/nginx.conf), [compose.yml](../../13-week/hu-status/deploy/compose.yml)
- Commits: [deb0b54](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/deb0b54), [d4b4a53](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d4b4a53), [ddc8b74](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ddc8b74), [8b96202](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8b96202), [edad5ef](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/edad5ef), [c0db426](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/c0db426), [e309165](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e309165), [de467c9](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/de467c9), [3a35b8c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3a35b8c), [aaf73b1](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/aaf73b1), [7b31c3f](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/7b31c3f), [76c9d44](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/76c9d44), [aecb1e4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/aecb1e4), [12e2a6a](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/12e2a6a), [794ee65](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/794ee65), [0e8cdfa](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0e8cdfa), [b9c3091](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/b9c3091), [e9f108c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e9f108c)
