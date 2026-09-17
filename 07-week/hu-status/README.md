<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       07-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 07

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Define the shared OpenAPI schema and the service contract template used to standardize API design across LMS services.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| API contracts | Shared OpenAPI components + per-service contract template | done | [_shared.yaml](_shared.yaml); [_template-service.yaml](_template-service.yaml) (commit 8db6b5c) |
| HU-01 | Access-service infrastructure adapters: HTTP router, login handler, JWT issuer/validation, health checks, Postgres repository | doing | [infrastructure/](infrastructure/) (untracked, not yet committed) |

## 2. My individual contribution
- Authored `_shared.yaml`: reusable OpenAPI components shared by every service contract, including `PaginatedMeta`, `ErrorResponse`, `UUID`, `Timestamp`, and `AuditFields` schemas, the `IdParam` parameter, standard error responses (`BadRequest`, `Unauthorized`, `NotFound`, `InternalError`), and the `bearerAuth` JWT security scheme.
- Authored `_template-service.yaml`: a base OpenAPI 3.0.3 template for per-service contracts, covering server variables, default bearer security, a `/health` endpoint, and full CRUD paths (list/create/get/update/delete) with request/response schemas referencing `_shared.yaml`.
- Implemented the access-service's infrastructure layer in Go (hexagonal architecture, secondary/primary adapters):
  - `infrastructure/auth/jwt_issuer.go` — `JWTIssuer`, a driven adapter for the `access.TokenIssuer` port; signs HS256 JWTs carrying only `sub`/`iat`/`exp` (no RBAC claims, per the single-role v1 domain model).
  - `infrastructure/http/handler/auth_handler.go` — `AuthHandler.Login`, implementing HU-01's `/auth/login`; decodes the request, delegates to the `usecase.Login` application service, and maps `access.ErrInvalidCredentials` to a generic `401 INVALID_CREDENTIALS`.
  - `infrastructure/http/handler/health_handler.go` — liveness (`/health`) and readiness (`/health/ready`, pings the DB with a 2s timeout) endpoints.
  - `infrastructure/http/middleware/auth.go` — `RequireAuth`, validates the Bearer JWT's signature/expiry and injects the administrator ID into the request context.
  - `infrastructure/http/middleware/correlation_id.go` — assigns/propagates an `X-Correlation-Id` per request via context and response header for structured logging.
  - `infrastructure/http/middleware/cors.go` — permissive CORS for local dev (production traffic is fronted by nginx).
  - `infrastructure/http/response/response.go` — shared `JSON`/`Error` envelope helpers matching the `ErrorResponse` schema in `_shared.yaml`.
  - `infrastructure/http/router.go` — `NewRouter`, wires chi with the middleware stack, health routes, and the public `POST /api/v1/auth/login` route.
  - `infrastructure/logger/logger.go` — JSON `zap.Logger` factory with configurable level.
  - `infrastructure/postgres/db.go` — `NewPool`, opens a pgx connection pool and verifies connectivity with a bounded ping timeout.
  - `infrastructure/postgres/administrator_repository.go` — `AdministratorRepository`, the sole adapter allowed to touch the `administrators` table; implements `FindByUsername` and `Save`.

## 3. Blockers and risks
- `_template-service.yaml` still contains placeholders (`[NombreServicio]`, `[recurso]`, `[campo1]`, etc.) — no concrete per-service contract has been instantiated from it yet in this folder.
- The `infrastructure/` code is untracked (not yet committed/pushed) and depends on `internal/application/usecase.Login` and `internal/domain/access` types that live in the actual service repo, not in this docs repo — it is included here only as this week's status evidence.
- No unit/integration tests included yet for the new adapters (JWT issuer, auth middleware, handlers, repository).

## 4. Plan for next week
- Instantiate the template into real per-service OpenAPI contracts (e.g. access-service, catalog-service, membership-service), replacing placeholders with the actual resources and fields for each service.
- Commit and push the `infrastructure/` adapters to the access-service repository, add unit tests, and open the corresponding HU-01 PR.

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- [Shared OpenAPI components](_shared.yaml)
- [Service contract template](_template-service.yaml)
- [Access-service infrastructure adapters](infrastructure/)
- Commit: `8db6b5c` — "Document drafts 07"
