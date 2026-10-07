<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       10-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 10

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Study distributed persistence for Unit 3 · Corte 2 (database per service, Saga, Outbox, CQRS and eventual consistency) and map the MVP 2 release requirements (integrated system, promotion develop -> qa -> main, v2.0.0 tag, failure-path demo) onto the LMS services.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| MVP2-PERSISTENCE | Concept map: Saga, Outbox, CQRS and eventual consistency (Part A) | done | [Mapa_Saga_Outbox_CQRS_MVP2.png](./Mapa_Saga_Outbox_CQRS_MVP2.png) (commit 10b23d1) |
| MVP2-RELEASE | MVP 2 release checklist (DoD), promotion flow and failure demo (Part B) | done | [Mapa_Saga_Outbox_CQRS_MVP2.png](./Mapa_Saga_Outbox_CQRS_MVP2.png) (commit 10b23d1) |
| MVP2-OUTBOX | Apply Outbox + idempotent consumers to the LMS services | todo | — |
| MVP2-TAG | Promote develop -> qa -> main and tag `v2.0.0` with CHANGELOG and ADRs | todo | — |

## 2. My individual contribution
- Added `Mapa_Saga_Outbox_CQRS_MVP2.png`, a one-page concept map of Unit 3 · Corte 2 ("Distributed persistence + MVP 2 release"), split in two parts.
- Part A (persistence in distributed systems) summarizes:
  - **Database per service**: each service owns its database and is accessed only through its contract (API/events); polyglot persistence (e.g. Postgres for orders, MongoDB for the catalog, Redis for sessions).
  - **Saga**: there is no global ACID transaction across services; a saga is a sequence of local transactions with a compensation per step, either orchestrated (a coordinator drives reserve -> charge -> confirm and releases stock if the charge fails) or choreographed (services react to `OrderPlaced` -> `StockReserved` -> `PaymentDone`).
  - **Outbox pattern**: avoids the dual-write problem by writing the state change and the event to an `outbox` table in the same DB transaction; a relay publishes to the broker with at-least-once delivery, so consumers must be idempotent.
  - **CQRS + eventual consistency**: separate write model (commands, aggregate) from a denormalized read model built from events; the UI must show intermediate states such as "processing" instead of stale data.
  - Common mistakes: 2PC everywhere, DB-then-broker without outbox, non-idempotent consumers, expecting an instantly consistent read model.
- Part B (MVP 2 release) summarizes:
  - What changes: MVP 1 was one working service; MVP 2 is an integrated system (gateway + services + message broker) with real contracts, saga + outbox consistency and per-environment config/secrets.
  - Promotion: develop -> qa -> main, verify in qa with production-like config, tag `v2.0.0` on main.
  - Release DoD: acceptance criteria met, unit + integration + contract tests green, `docker compose up` with health checks, end-to-end flow, tested compensation path, outbox + idempotent consumers, no secrets in git or images, `v2.0.0` + CHANGELOG + updated ADRs.
  - Demo including a failure (stop a service or send a duplicate event and show the saga compensates and the consumer ignores the duplicate), retrospective feeding Corte 3, and grading criteria.
- Key idea captured: consistency between services = Sagas + Outbox + idempotent consumers; MVP 2 must demonstrate it by stopping a service mid-flow.

## 3. Blockers and risks
- The LMS services do not yet have an outbox table or relay; events (e.g. from lms-catalog-api) are still at risk of the dual-write problem.
- Consumers are not idempotent yet, so at-least-once delivery could cause duplicated effects.
- `GET /books` is still missing in lms-catalog-api (carried over from week 09), which blocks the end-to-end catalog flow required by the MVP 2 demo.
- No contract/integration tests yet, which the MVP 2 DoD requires.

## 4. Plan for next week
- Identify the cross-service flow for the LMS (e.g. loan: reserve copy -> register loan -> confirm) and decide orchestration vs choreography, documenting it in an ADR.
- Implement the Outbox pattern (table + relay) in the service that publishes events, and make the consumer idempotent (processed-event key).
- Add contract and integration tests, wire health checks into `docker compose`, and prepare the failure-path demo.
- Promote develop -> qa -> main and tag `v2.0.0` with CHANGELOG.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Diagram: [Mapa_Saga_Outbox_CQRS_MVP2.png](./Mapa_Saga_Outbox_CQRS_MVP2.png)
- Commit: `10b23d1` — "docs(hu-status): link Saga Outbox CQRS diagram"
