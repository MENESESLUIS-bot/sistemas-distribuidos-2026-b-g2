# ADR-003 — NGINX as Reverse Proxy / Edge Layer

| Field | Value |
|-------|-------|
| **ID** | ADR-003 |
| **Date** | 2026-08-26 |
| **Status** | Accepted |
| **Authors** | Oscar Areiza — Tech Lead |
| **Reviewers** | Hermes Pascuas, Luis Alejandro Meneses — Development team |

---

## Context

The Administrator reaches the system through a browser running the React SPA. The system
needs a single entry point that can: terminate TLS, serve the compiled React build, and
reverse-proxy API calls to `library-api` (the Go backend) — while staying within a Docker
Compose deployment a 3-person academic team can run and maintain (`01-context/overview.md`,
`01-context/scope.md`).

**Known constraints:**
- Single backend service in v1 (`ADR-002`) — no multi-service routing need yet
- `04-requirements/non-functional.md` (NFR-004) requires HTTPS in production and (NFR-001)
  defines latency/throughput targets that benefit from a place to apply basic rate limiting
- No budget for a managed API Gateway product; infrastructure limited to free/local tooling
  (`01-context/scope.md`, Constraints)

---

## Decision

**We decided:** use **NGINX** as a lightweight reverse proxy and static file server in front of
`library-api`, instead of adopting a full API Gateway product (Kong, Traefik).

**Justification:** NGINX already covers everything v1 needs — TLS termination, reverse-proxying
`/api/*` to `library-api`, serving the built React SPA, and basic rate limiting/security
headers — using a single, well-documented, minimal-configuration tool. A full API Gateway adds
capabilities (service discovery, per-route auth plugins, multi-service request aggregation)
that a single-backend v1 system does not need.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| Kong / Traefik (full API Gateway) | Built-in auth plugins, rate limiting, service discovery, ready for multi-service routing | Extra moving part to deploy and operate (Kong needs its own datastore or DB-less mode; Traefik needs dynamic config); overkill for one backend service | No multi-service routing need in v1; no spare team time to operate a gateway product |
| No reverse proxy — expose `library-api` directly | Simplest possible setup | No single TLS termination point; no place to serve the static SPA build; no central point for rate limiting/security headers; harder to insert a gateway later without a client-facing URL change | Would leave `library-api` responsible for TLS/static files/rate limiting itself, against NFR-001/NFR-004 |
| **NGINX (CHOSEN)** | Lightweight, battle-tested, doubles as static file server + reverse proxy + TLS termination; trivial to swap for Kong/Traefik later behind the same DNS name if the system grows | NGINX config becomes another artifact to version and document | — (chosen) |

---

## Consequences

**Positive:**
- One well-known component handles TLS, static asset serving, and reverse proxying, keeping
  `library-api` free of those infrastructure concerns (consistent with the Hexagonal Architecture
  principle of keeping infra out of the domain/application core — see `ADR-002`)
- Easy migration path: if the system later needs multi-service routing, NGINX can be replaced
  by Kong/Traefik behind the same public DNS name without changing `library-api`

**Negative / Trade-offs:**
- NGINX configuration must be written, versioned, and kept in sync with `library-api`'s routes
  (tracked in `05-architecture/deployment.md`)
- **NGINX does not perform authentication or authorization** — it only handles TLS, routing,
  and rate limiting. JWT validation stays inside `library-api`. This must not be confused with
  an API Gateway's auth-plugin capabilities.

**Impact on the system:**
- Affected services: adds one container (`nginx`) to the Docker Compose topology
- Documents that must be updated: `05-architecture/deployment.md`, `10-devops/local-setup.md`
  (once written)

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Team assumes NGINX validates the JWT / enforces authorization | Medium | High (would leave endpoints unprotected) | Explicit note in `05-architecture/cross-cutting.md` and this ADR: authn/authz always happens inside `library-api` |
| NGINX config drifts from `library-api`'s actual routes as the API evolves | Medium | Medium | Review NGINX config as part of the PR whenever a route changes (per `00-governance/git-conventions.md`) |

---

## References

- API Gateway pattern description → `05-architecture/pattern-guide.md`
- Deployment topology → `05-architecture/deployment.md`
- Related to: `ADR-002-hexagonal-modular-monolith.md`
