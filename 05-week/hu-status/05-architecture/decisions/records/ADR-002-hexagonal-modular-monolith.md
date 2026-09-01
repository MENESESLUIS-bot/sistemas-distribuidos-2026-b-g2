# ADR-002 — Architectural Style: Hexagonal Modular Monolith (Microservices-Ready)

| Field | Value |
|-------|-------|
| **ID** | ADR-002 |
| **Date** | 2026-08-26 |
| **Status** | Accepted |
| **Authors** | Oscar Areiza — Tech Lead |
| **Reviewers** | Hermes Pascuas, Luis Alejandro Meneses — Development team |

---

## Context

`02-domain/domain-map.md` already identifies four bounded contexts — Circulation (Core),
Catalog (Supporting), Membership (Supporting), Access (Generic) — and records the decision to
ship them inside a single deployable service, `library-api`, rather than as separate
microservices from day one. That decision was made at the domain-modeling level; this ADR
formalizes the corresponding **code-level architectural style**, since it directly shapes how
`library-api` is built and how easily it can be split later.

**Known constraints:**
- 3-person academic team (Oscar Areiza, Hermes Pascuas, Luis Meneses), one academic term
- No confirmed need for independent scaling per bounded context (`01-context/scope.md`,
  assumption #4: low hundreds/thousands of books and students)
- Must use Go + PostgreSQL + Docker (product brief); no message broker adopted in v1
  (`02-domain/domain-events.md` — events are in-process)
- The domain work in `02-domain/` must not be wasted by an incompatible code structure

---

## Decision

**We decided:** build `library-api` as a single deployable Go service, internally organized
using **Hexagonal Architecture (Ports & Adapters)** with one internal module per bounded
context (`access`, `circulation`, `catalog`, `membership`), each following the
domain/application/infrastructure split documented in `_stacks/go.md` and
`05-architecture/hexagonal-architecture.md`. Modules communicate in-process (direct calls /
in-process domain events), never through another module's database tables.

**Justification:** this keeps each bounded context's business logic isolated and unit-testable
without infrastructure, while avoiding the deployment/operations overhead (4 deployables, a
message broker, database-per-service) that a real microservices split would require for a team
and timeline this size. Because module boundaries already mirror bounded contexts, any module
can be promoted into its own microservice later — given its own datastore and a broker for
cross-context events — without rewriting its business logic.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| Full microservices from day 1 (one deployable + DB + broker per context) | Textbook microservices; independent scaling and deployment per context | Requires operating 4 services, a message broker (Kafka/RabbitMQ), and 4 databases with a 3-person academic team; no scaling need yet | No scaling need in v1 (`01-context/scope.md` assumption #4); team size/timeline can't operate this |
| Traditional layered monolith (Controller → Service → Repository, no ports/adapters) | Fastest to start, minimal ceremony | Business logic couples to the framework/DB; hard to unit-test the domain; hard to later extract a context into its own service without a rewrite | Would discard the DDD work already done in `02-domain/`; contradicts the TDD goals in `11-quality/` |
| **Hexagonal Modular Monolith (CHOSEN)** | Domain testable without infra; clear seams per bounded context; low-risk path to real microservices later | Slightly more ceremony (explicit ports/interfaces) than a plain layered monolith | — (chosen) |

---

## Consequences

**Positive:**
- Domain logic (entities, invariants from `02-domain/entities-and-rules.md`) is testable in
  isolation, supporting the TDD workflow planned for `11-quality/`
- The Core/Supporting/Generic classification in `02-domain/domain-map.md` maps directly onto
  code modules, keeping domain and code vocabulary identical
- A well-defined, low-risk migration path exists if the project ever needs to split a context
  into its own microservice (see "Planned evolution" in `05-architecture/overview.md`)

**Negative / Trade-offs:**
- All modules currently share one PostgreSQL database (see architectural debt `AT-001` in
  `05-architecture/overview.md`) — the team must not let one module query another module's
  tables directly, or the seams needed for a future split are lost
- More interfaces (ports) to write than a plain layered monolith would need

**Impact on the system:**
- Affected services: `library-api` internal package layout (`internal/domain/{access,circulation,catalog,membership}/...`)
- Documents that must be updated: `09-microservices/service-catalog.md` once that section is filled in

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| A module ends up calling another module's repository/table directly, silently coupling them | Medium | High (blocks any future split) | Enforce the Hexagonal Architecture Checklist in `05-architecture/hexagonal-architecture.md` during code review |
| Team over-invests in ports/interfaces for trivial operations, slowing delivery | Low | Low | Keep ports at the module boundary only — no ports needed for calls within the same module |

---

## References

- Bounded contexts and the "one deployable service" domain decision → `02-domain/domain-map.md`
- Hexagonal architecture concepts and checklist → `05-architecture/hexagonal-architecture.md`
- Go folder structure per module → `_stacks/go.md`
- Related to: `ADR-003-nginx-reverse-proxy.md` (the edge layer in front of this service)
