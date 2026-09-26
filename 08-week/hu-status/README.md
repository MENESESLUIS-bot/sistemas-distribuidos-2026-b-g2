<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       08-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 08

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Wire the catalog-service HTTP entry point, containerize it, and implement the HU-04/HU-06/HU-07 application use cases on top of the Catalog domain model.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| catalog-service | API entry point (`cmd/api/main.go`) + Docker build | done | [Dockerfile](Dockerfile); [cmd/api/main.go](cmd/api/main.go) (commit 3616550) |
| catalog-service | Go module dependencies (`go.mod`/`go.sum`) | done | [go.mod](go.mod); [go.sum](go.sum) (commit d060fd5) |
| HU-04 | Create book use case (reject duplicate ISBN) | done | [create_book.go](internal/application/usecase/create_book.go); [create_book_test.go](internal/application/usecase/create_book_test.go); [book.go](internal/domain/catalog/book.go) (commits 37ec6e8, d807d9d) |
| HU-06 / HU-07 | Loan / return book copy use cases | done | [adjust_book_availability.go](internal/application/usecase/adjust_book_availability.go); [book.go](internal/domain/catalog/book.go) (commits 37ec6e8, d807d9d) |

## 2. My individual contribution
- Authored `cmd/api/main.go`: the catalog-service process entry point — loads config, sets up a JSON zap logger, opens a Postgres pool, wires the book repository/use cases (create, loan copy, return copy) into the HTTP handler and router, and runs the server with graceful shutdown on SIGINT/SIGTERM.
- Authored `Dockerfile`: a two-stage build — `golang:1.25-alpine` compiles a static (`CGO_ENABLED=0`) binary from `./cmd/api`, and the runtime stage copies just the binary onto `alpine:3.21` with CA certificates, exposing port 8080.
- Added `go.mod`/`go.sum` (module `github.com/code-corhuila/lms-catalog-api`, Go 1.25.1), pinning the service's direct dependencies (`chi`, `golang-jwt`, `google/uuid`, `pgx`, `zap`, `golang.org/x/crypto`) and their transitive/indirect requirements, so `go mod download` and the Dockerfile build now resolve.
- Authored `internal/application/usecase/create_book.go`: `CreateBook` (HU-04) — rejects registration when the ISBN already exists (`ErrISBNAlreadyExists`), otherwise builds a `catalog.Book` via the domain constructor and persists it through `catalog.BookRepository`; covered by `create_book_test.go`.
- Authored `internal/application/usecase/adjust_book_availability.go`: `LoanBookCopy` (HU-06) and `ReturnBookCopy` (HU-07) — look up the book by ID, delegate the copy-count change to `catalog.Book`'s domain methods, and persist the result. These exist as HTTP-facing use cases (rather than direct repository calls) because `Book` now lives in a separate database from circulation-service, per the service-boundary decomposition.
- Added `Makefile` with `dev`/`test`/`test-cover`/`build`/`lint` targets; `migrate-up`/`migrate-down` were intentionally dropped since schema ownership moved to the `lms-catalog-db` repo (ADR-006).
- Authored `internal/domain/catalog/book.go`: the `Book` aggregate root (Catalog bounded context) — `NewBook` (INV-003: at least one copy at registration), `LoanOneCopy`/`ReturnOneCopy` (INV-001: availability never goes negative or exceeds total), and `Update` (HU-09, ISBN intentionally non-editable); covered by `book_test.go`.
- Authored `internal/domain/catalog/port.go`: the `BookRepository` driven port (`FindByID`, `FindByISBN`, `Search`, `Save`) that the use cases and future Postgres adapter depend on.
- Authored `internal/config/config.go`: env-var-only configuration loader (`Config.Load`) for port, DB connection, JWT secret/expiry, log level, and CORS origin; fails fast if `JWT_SECRET` is unset, and builds the Postgres DSN via `Config.DSN()`.

## 3. Blockers and risks
- Only `create_book.go` and `book.go` have tests so far; `adjust_book_availability.go` (loan/return use cases) still needs unit test coverage.
- No Postgres adapter implementing `catalog.BookRepository` yet, so the service cannot run end to end — `main.go`'s wiring is still unverified against a real database.

## 4. Plan for next week
- Implement the Postgres adapter for `catalog.BookRepository` and verify the service builds/runs end to end via the Dockerfile.
- Add tests for `LoanBookCopy`/`ReturnBookCopy`, wire up `golangci-lint`, and open the corresponding HU PRs.

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- [API entry point](cmd/api/main.go)
- [Dockerfile](Dockerfile)
- [go.mod](go.mod) / [go.sum](go.sum)
- [Use cases](internal/application/usecase/) / [Makefile](Makefile)
- [Catalog domain](internal/domain/catalog/) / [Config loader](internal/config/config.go)
- Commit: `3616550` — "Add Dockerfile and API entrypoint"
- Commit: `d060fd5` — "Add Go module files for catalog-service"
- Commit: `37ec6e8` — "Add catalog-service use cases and build Makefile"
- Commit: `d807d9d` — "Add catalog domain entity and env-based config loader"
