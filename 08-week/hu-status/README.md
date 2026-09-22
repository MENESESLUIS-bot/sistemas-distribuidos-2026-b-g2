<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       08-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 08

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Wire the catalog-service HTTP entry point and containerize it with a multi-stage Dockerfile.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| catalog-service | API entry point (`cmd/api/main.go`) + Docker build | done | [Dockerfile](Dockerfile); [cmd/api/main.go](cmd/api/main.go) (commit 3616550) |

## 2. My individual contribution
- Authored `cmd/api/main.go`: the catalog-service process entry point — loads config, sets up a JSON zap logger, opens a Postgres pool, wires the book repository/use cases (create, loan copy, return copy) into the HTTP handler and router, and runs the server with graceful shutdown on SIGINT/SIGTERM.
- Authored `Dockerfile`: a two-stage build — `golang:1.25-alpine` compiles a static (`CGO_ENABLED=0`) binary from `./cmd/api`, and the runtime stage copies just the binary onto `alpine:3.21` with CA certificates, exposing port 8080.

## 3. Blockers and risks
- No unit/integration tests included yet for `main.go`'s wiring.
- `go.mod`/`go.sum` and the rest of the service's internal packages are not present in this docs folder, so the Dockerfile build context assumes they exist in the actual service repo.

## 4. Plan for next week
- Add a `docker-compose` service definition (catalog-service + Postgres) and verify the container builds/runs end to end.
- Add tests for the router wiring and open the corresponding HU PR.

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
- Commit: `3616550` — "Add Dockerfile and API entrypoint"
