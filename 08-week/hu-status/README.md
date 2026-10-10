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

### My commits and PRs this week (all repos)
Every commit and pull request authored by `MENESESLUIS-bot` between 2026-09-21 and 2026-09-27 (America/Bogota). Cherry-picks of the same commit into develop/qa/main are listed once, with the original SHA.

**Pull requests (6)**

| Opened | Repo | PR | Title | State |
|---|---|---|---|---|
| 2026-09-24 | code-corhuila/library-docs | [#15](https://github.com/code-corhuila/library-docs/pull/15) | docs(api): close Catalog & Circulation contract decisions and declare gaps | closed |
| 2026-09-24 | code-corhuila/library-docs | [#16](https://github.com/code-corhuila/library-docs/pull/16) | docs(api): align OpenAPI contracts with lms-library v1.0.0 | closed |
| 2026-09-24 | code-corhuila/library-docs | [#21](https://github.com/code-corhuila/library-docs/pull/21) | docs(07-api): add catalog-service OpenAPI contract and rationale | merged |
| 2026-09-24 | code-corhuila/library-docs | [#22](https://github.com/code-corhuila/library-docs/pull/22) | docs(07-api): add catalog-service contract rationale | closed |
| 2026-09-24 | code-corhuila/lms-catalog-db | [#2](https://github.com/code-corhuila/lms-catalog-db/pull/2) | chore(migrations): migrate catalog-db schema | merged |
| 2026-09-26 | code-corhuila/lms-catalog-api | [#2](https://github.com/code-corhuila/lms-catalog-api/pull/2) | chore(catalog): bootstrap catalog-service skeleton | merged |

**Commits (47)**

| Date | Repo | Commit | Message |
|---|---|---|---|
| 2026-09-22 | code-corhuila/lms-catalog-api | [`901da3b`](https://github.com/code-corhuila/lms-catalog-api/commit/901da3b91ac0bdaf036b84496191b24df6b3fe07) | chore(docker): add multi-stage Dockerfile for catalog-service |
| 2026-09-22 | code-corhuila/lms-catalog-api | [`b90b26c`](https://github.com/code-corhuila/lms-catalog-api/commit/b90b26c56f6db61322fd3a67edb45428bb54c7c0) | chore(catalog): add API entrypoint skeleton |
| 2026-09-22 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`3616550`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3616550eb98cf65614b0c3035e37685448e949bf) | Add Dockerfile and API entrypoint |
| 2026-09-22 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`461d71d`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/461d71d48225e2ed49afda98033066e2983ea6c5) | docs(hu-status): update week 08 status with catalog-service entrypoint and Dockerfile |
| 2026-09-23 | code-corhuila/lms-catalog-api | [`c523df8`](https://github.com/code-corhuila/lms-catalog-api/commit/c523df8a571af753046a0f327eb3065789f32f3d) | chore(catalog): add go.sum for catalog-service |
| 2026-09-23 | code-corhuila/lms-catalog-api | [`6036fc0`](https://github.com/code-corhuila/lms-catalog-api/commit/6036fc0c94e8e517c9d56ee2526fd049122b74cb) | chore(catalog): add go.mod for catalog-service |
| 2026-09-23 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`d060fd5`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d060fd5c1959d50fbdea732d50a39a0f193f1e27) | Add Go module files for catalog-service |
| 2026-09-23 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`1ad7e87`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/1ad7e87b37e6cf0464aaa798c58b7d8b50bb4af0) | docs(hu-status): update week 08 status with go.mod/go.sum |
| 2026-09-24 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`0f09a11`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0f09a11f9d506a7388f5c46f69ea857712a0e514) | Merge branch 'code-corhuila:main' into main |
| 2026-09-24 | code-corhuila/library-docs | [`5e077e8`](https://github.com/code-corhuila/library-docs/commit/5e077e8eae74212e1e07a5529ee448d15aa2cf08) | docs(api): close Catalog & Circulation contract decisions and declare gaps |
| 2026-09-24 | code-corhuila/library-docs | [`39eb906`](https://github.com/code-corhuila/library-docs/commit/39eb90687bbee34cf3c7aef5325608d1f807996a) | docs(api): align OpenAPI contracts with lms-library v1.0.0 |
| 2026-09-24 | code-corhuila/library-docs | [`6b52dce`](https://github.com/code-corhuila/library-docs/commit/6b52dce5d05fc8f24b747a44b9761ba52a806e12) | docs(uml): add request path (read) sequence diagram |
| 2026-09-24 | code-corhuila/library-docs | [`6b3946d`](https://github.com/code-corhuila/library-docs/commit/6b3946dd549bef43ad0363c9b16824ae061e2c5e) | docs(uml): add hexagonal architecture diagram |
| 2026-09-24 | code-corhuila/library-docs | [`7308163`](https://github.com/code-corhuila/library-docs/commit/73081630da3dc6a2380423016ccf44e9615b3422) | docs(uml): add dependency rule diagram |
| 2026-09-24 | code-corhuila/lms-catalog-api | [`2ef7a67`](https://github.com/code-corhuila/lms-catalog-api/commit/2ef7a6746b4c9928cceb85480de1e025bd99506f) | feat(catalog): add adjust book availability use case (HU-06/HU-07) |
| 2026-09-24 | code-corhuila/lms-catalog-api | [`e7f22ed`](https://github.com/code-corhuila/lms-catalog-api/commit/e7f22edd4a763f526e2f4c850f3ff9c241da82ad) | feat(catalog): add create book use case (HU-04) |
| 2026-09-24 | code-corhuila/lms-catalog-api | [`fbb812f`](https://github.com/code-corhuila/lms-catalog-api/commit/fbb812f188cd770d0883728af9cea1609f712846) | chore(catalog): add Makefile for catalog-service |
| 2026-09-24 | code-corhuila/lms-catalog-api | [`bd58056`](https://github.com/code-corhuila/lms-catalog-api/commit/bd58056607fde28f6549c9067a886cb0ca32ee4b) | test(catalog): add tests for create book use case (HU-04) |
| 2026-09-24 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`37ec6e8`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/37ec6e80b39533a6124e734463e2c6df376c903b) | Add catalog-service use cases and build Makefile |
| 2026-09-24 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`570a824`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/570a824baa65697a88ca6f650ec7883762663273) | docs(hu-status): update week 08 status with catalog-service use cases |
| 2026-09-24 | code-corhuila/library-docs | [`41b8439`](https://github.com/code-corhuila/library-docs/commit/41b8439b7a3743b1f003b362390d106bffba888a) | Merge branch 'docs/uml-diagrams-svg' of https://github.com/code-corhuila/library-docs into docs/uml-diagrams-svg |
| 2026-09-24 | code-corhuila/library-docs | [`676089d`](https://github.com/code-corhuila/library-docs/commit/676089dfdc2820f0be88b881fd7a63f49093d570) | docs(07-api): add catalog-service OpenAPI contract |
| 2026-09-24 | code-corhuila/library-docs | [`7c5c458`](https://github.com/code-corhuila/library-docs/commit/7c5c458c20fe31c97293790717f20ed26f4dd986) | docs(07-api): explain catalog-service contract decisions and consumers |
| 2026-09-24 | code-corhuila/library-docs | [`f15d507`](https://github.com/code-corhuila/library-docs/commit/f15d50757fd41d8eb07a19a0e89758856c7c2266) | docs(07-api): apply review — split rationale out, mark proposed decisions |
| 2026-09-24 | code-corhuila/library-docs | [`cce2934`](https://github.com/code-corhuila/library-docs/commit/cce29348570d6c070defd6cfa3aeb47cced0ef58) | docs(07-api): keep contract and rationale in one PR under 400 lines |
| 2026-09-24 | code-corhuila/lms-catalog-db | [`75a7369`](https://github.com/code-corhuila/lms-catalog-db/commit/75a736932c7ff7092a8803481ca35a199b19a4bd) | feat(migrations): add books table up migration |
| 2026-09-24 | code-corhuila/lms-catalog-db | [`73d6539`](https://github.com/code-corhuila/lms-catalog-db/commit/73d65395d4af3b7e9b52cef579f6f3fbf24143f7) | chore(docker): add standalone compose for catalog-db and migrate |
| 2026-09-24 | code-corhuila/lms-catalog-db | [`a09d516`](https://github.com/code-corhuila/lms-catalog-db/commit/a09d516e8596cf5bc5a995896d22475edadd0aef) | feat(migrations): add books table down migration |
| 2026-09-24 | code-corhuila/lms-catalog-db | [`03c45d4`](https://github.com/code-corhuila/lms-catalog-db/commit/03c45d42d34fc552de426a0a44b282bf71ea15a2) | Merge pull request #2 from code-corhuila/chore/migrate-catalog-db-schema |
| 2026-09-24 | code-corhuila/library-docs | [`f219ac6`](https://github.com/code-corhuila/library-docs/commit/f219ac6a6647b69957e767f98179d63d9f4be132) | Merge pull request #21 from code-corhuila/docs/07-api-catalog-service-contract |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`0c07904`](https://github.com/code-corhuila/lms-catalog-api/commit/0c0790400d7c96303cb9d5af0e80cd7de1bbd467) | feat(catalog): add Book domain entity (HU-04/HU-06/HU-07/HU-09) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`9e30544`](https://github.com/code-corhuila/lms-catalog-api/commit/9e30544d34ba93114d7363bf9ee61e9dd10ed872) | chore(catalog): add service configuration loader |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`a7f7ea4`](https://github.com/code-corhuila/lms-catalog-api/commit/a7f7ea4fafb03f1929a34d8728cdfc53d362bac1) | test(catalog): add tests for Book domain entity (HU-04/HU-06/HU-07) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`17a64e8`](https://github.com/code-corhuila/lms-catalog-api/commit/17a64e8cc960e94590e8659f9c7405183edcb55b) | feat(catalog): add BookRepository port |
| 2026-09-26 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`d807d9d`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d807d9de5f20797928939010637c174cf648a7a0) | Add catalog domain entity and env-based config loader |
| 2026-09-26 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`a9c65ab`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/a9c65abda2387156fd4e071f7aa9a2c01a4ed173) | docs(hu-status): update week 08 status with catalog domain and config loader |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`0f8db00`](https://github.com/code-corhuila/lms-catalog-api/commit/0f8db009c903d0912a4d0daecd6627bc386774b6) | chore(catalog): add postgres connection pool setup |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`928dcfd`](https://github.com/code-corhuila/lms-catalog-api/commit/928dcfd076f3ce8193571abfb4b178553a3c102a) | chore(catalog): add zap logger setup |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`f9f4e0b`](https://github.com/code-corhuila/lms-catalog-api/commit/f9f4e0be5ae5fddc785670c5243eec5a35c715ab) | feat(catalog): add postgres BookRepository adapter (HU-04) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`6b43882`](https://github.com/code-corhuila/lms-catalog-api/commit/6b43882d1282df23ee3f69f0e7143e875f51ea38) | chore(catalog): add CORS middleware |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`5149145`](https://github.com/code-corhuila/lms-catalog-api/commit/51491459d5290401517a8fcdfb2c09a02d242ae5) | chore(catalog): add HTTP JSON response helpers |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`393dcb8`](https://github.com/code-corhuila/lms-catalog-api/commit/393dcb861d9892914780d92ccafe4616072f5d86) | feat(catalog): add JWT auth middleware (HU-01) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`dd61623`](https://github.com/code-corhuila/lms-catalog-api/commit/dd61623091178f88313277381b91bbe1401d3519) | chore(catalog): add correlation ID middleware |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`0a88c84`](https://github.com/code-corhuila/lms-catalog-api/commit/0a88c84068e27aa6c8fc2cf0ff3b49c4c38e3ccc) | feat(catalog): add book HTTP handler (HU-04) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`2220005`](https://github.com/code-corhuila/lms-catalog-api/commit/222000532c5d8d52944e01b63bbb0c37b8c92ce3) | chore(catalog): add liveness/readiness health handler |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`c0b2a39`](https://github.com/code-corhuila/lms-catalog-api/commit/c0b2a395edd3a1c31f9e1463c368fb9d3018c6dc) | feat(catalog): add HTTP router wiring (HU-04) |
| 2026-09-26 | code-corhuila/lms-catalog-api | [`30172a7`](https://github.com/code-corhuila/lms-catalog-api/commit/30172a76ff97f6f057578a1b72a44564932a63c8) | Merge pull request #2 from code-corhuila/chore/bootstrap-catalog-service-skeleton |
