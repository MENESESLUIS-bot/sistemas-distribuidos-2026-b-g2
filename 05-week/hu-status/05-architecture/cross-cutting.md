# Cross-Cutting Concerns

> Concerns that apply to every module of `library-api` (Access, Circulation, Catalog,
> Membership), not just one. In this architecture (`ADR-002`) there is only one deployable
> service, so "cross-cutting" here means "shared across internal modules," not "shared across
> microservices."

---

## Logging

- **Format:** structured JSON, one line per event.
- **Library:** `zap` (Go) — see `_stacks/go.md`.
- **Correlation ID:** every incoming HTTP request gets a `correlationId` (UUID), generated at
  the `nginx` → `library-api` boundary if not already present, and propagated through every log
  line produced while handling that request.
- **What gets logged:** every state-changing operation (loan registered, return registered,
  suspension applied), every rejected operation with its reason (INV-00X code from
  `02-domain/entities-and-rules.md`), and every authentication attempt (success/failure, never
  the password).
- **What never gets logged:** plaintext passwords, JWTs, or full request bodies containing PII.

---

## Distributed tracing

- Not required in v1 — there is a single service, so there is nothing to trace *across*.
- `correlationId` (see Logging above) is enough to follow one request through
  `nginx → library-api → postgres` at academic scale.
- If a module is ever split into its own microservice (`ADR-002`, planned evolution), adopt
  OpenTelemetry + Jaeger at that point, per `04-requirements/non-functional.md` (NFR-005).

---

## Centralized configuration

- All environment-specific values (DB connection string, JWT secret, JWT expiry, log level)
  come from environment variables, injected via `.env` files per environment
  (`05-architecture/deployment.md`).
- No service reads configuration from a file baked into its image, and no config is hardcoded —
  this is required for the same image to run unmodified in Local, Development, and Production.
- No centralized config server (e.g., Consul, Spring Cloud Config) is used — unnecessary for a
  single deployable service.

---

## Feature flags

- Not used in v1. With a single service and a 3-person team working from the same backlog
  (`04-requirements/user-stories.md`), a feature flag system adds operational complexity with
  no corresponding benefit at this scale.

---

## Error handling

- Every domain rule violation (an `INV-00X` from `02-domain/entities-and-rules.md`) is
  translated at the module boundary into an application-level error, then into a standard HTTP
  error response by the primary (HTTP) adapter — the domain layer never returns an HTTP status
  code itself (hexagonal dependency rule, see `05-architecture/hexagonal-architecture.md`).
- **Standard error response shape** (finalized in `07-api/contracts/openapi/_shared.yaml` once
  written):
  ```json
  {
    "error": {
      "code": "STUDENT_SUSPENDED",
      "message": "This student cannot receive a new loan until their suspension ends.",
      "correlationId": "..."
    }
  }
  ```
- Authentication failures always return the same generic error regardless of whether the
  username or the password was wrong (`02-domain/entities-and-rules.md`, INV-002 on
  Administrator).

---

## Retry policies

- **`library-api` → `postgres`:** the DB connection pool (pgx) retries transient connection
  failures with a short bounded backoff; a query that fails due to a genuine constraint
  violation (e.g., duplicate ISBN) is never retried — it is a business rejection, not a
  transient fault.
- **`nginx` → `library-api`:** no retry — if `library-api` is down, NGINX returns a 502/503
  immediately; there is no second instance to fail over to in v1 (see NFR-003, up to 3
  instances is a *possible* future scaling step, not a v1 guarantee).
- No retries exist between bounded-context modules, because they communicate in-process
  (direct function calls) — there is nothing transient to retry.

---

## Correlations

- Architectural style this builds on → `05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`
- NGINX's role (and what it does *not* do — no auth, no retries) → `05-architecture/decisions/records/ADR-003-nginx-reverse-proxy.md`
- Observability requirements (NFR-005) → `04-requirements/non-functional.md`
- Domain invariants that drive error handling → `02-domain/entities-and-rules.md`
