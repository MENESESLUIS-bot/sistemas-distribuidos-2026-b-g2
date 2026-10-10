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

### My commits and PRs this week (all repos)
Every commit and pull request authored by `MENESESLUIS-bot` between 2026-10-05 and 2026-10-11 (America/Bogota). Cherry-picks of the same commit into develop/qa/main are listed once, with the original SHA.

**Pull requests (8)**

| Opened | Repo | PR | Title | State |
|---|---|---|---|---|
| 2026-10-06 | code-corhuila/lms-catalog-api | [#10](https://github.com/code-corhuila/lms-catalog-api/pull/10) | release: 1.0.0 — promote validated qa stories to main | merged |
| 2026-10-06 | code-corhuila/lms-catalog-portal | [#5](https://github.com/code-corhuila/lms-catalog-portal/pull/5) | release: 0.1.0 catalog portal mvp1 | closed |
| 2026-10-06 | code-corhuila/lms-catalog-db | [#6](https://github.com/code-corhuila/lms-catalog-db/pull/6) | release: 1.0.0 — catalog-db Liquibase schema | merged |
| 2026-10-06 | code-corhuila/lms-catalog-portal | [#6](https://github.com/code-corhuila/lms-catalog-portal/pull/6) | release: 1.0.0 catalog portal mvp1 | merged |
| 2026-10-08 | code-corhuila/lms-front | [#9](https://github.com/code-corhuila/lms-front/pull/9) | test(front): add unit and component test suite with Vitest | merged |
| 2026-10-08 | code-corhuila/lms-front | [#10](https://github.com/code-corhuila/lms-front/pull/10) | qa: promote shell + test suite (unit, connection, stress, perf) to qa | merged |
| 2026-10-08 | code-corhuila/lms-front | [#12](https://github.com/code-corhuila/lms-front/pull/12) | test(front): add connection, stress and performance tests | merged |
| 2026-10-09 | code-corhuila/lms-catalog-api | [#11](https://github.com/code-corhuila/lms-catalog-api/pull/11) | fix(catalog): add PATCH /books/{id} for HU-09 book editing | open |

**Commits (14)**

| Date | Repo | Commit | Message |
|---|---|---|---|
| 2026-10-05 | code-corhuila/lms-catalog-api | [`4dd992f`](https://github.com/code-corhuila/lms-catalog-api/commit/4dd992fa9383a140fbf444e8059357308d2757fb) | Merge pull request #8 from code-corhuila/qa-promote/develop-catalog-skeleton |
| 2026-10-06 | code-corhuila/lms-catalog-portal | [`8e36870`](https://github.com/code-corhuila/lms-catalog-portal/commit/8e368708b1e76afcf70bf934d42dead643be687b) | Merge pull request #4 from code-corhuila/promote-qa/project-scaffold |
| 2026-10-06 | code-corhuila/lms-catalog-db | [`f1f00c4`](https://github.com/code-corhuila/lms-catalog-db/commit/f1f00c434581df07e183c552ae073dad6c2b2381) | Merge pull request #5 from code-corhuila/qa-catalog-db-liquibase-schema |
| 2026-10-06 | code-corhuila/lms-catalog-api | [`5a2ccbb`](https://github.com/code-corhuila/lms-catalog-api/commit/5a2ccbbc77b150b6846b9dcfca4a04437dc039a5) | Merge pull request #10 from code-corhuila/release/1.0.0 |
| 2026-10-06 | code-corhuila/lms-catalog-portal | [`b3b578d`](https://github.com/code-corhuila/lms-catalog-portal/commit/b3b578dde029d72ae19e21c7335780dd02666cfc) | Merge pull request #6 from code-corhuila/release/1.0.0 |
| 2026-10-06 | code-corhuila/lms-catalog-db | [`8b05827`](https://github.com/code-corhuila/lms-catalog-db/commit/8b0582795a7f433ea2557f76c56a9f812a5c3ca8) | Merge pull request #6 from code-corhuila/release/1.0.0 |
| 2026-10-07 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`10b23d1`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/10b23d1b977d1e3ec3a967c14b6f9d101ab289da) | docs(hu-status): link Saga Outbox CQRS diagram |
| 2026-10-07 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`d64e1a5`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d64e1a5ede45bd00f0f5c8092195e2a7ca5ab7b6) | docs(hu-status): update week 10 status with Saga, Outbox and CQRS concept map |
| 2026-10-08 | code-corhuila/lms-front | [`38466bd`](https://github.com/code-corhuila/lms-front/commit/38466bd01fb2601bd7c0a314f7a0dca960d1c7f2) | test(front): add unit and component test suite with Vitest |
| 2026-10-08 | code-corhuila/lms-front | [`a76845d`](https://github.com/code-corhuila/lms-front/commit/a76845d90e894114f640179d1c8cff0cfea06fba) | test(front): add connection, stress and performance tests |
| 2026-10-08 | code-corhuila/lms-front | [`3e32e55`](https://github.com/code-corhuila/lms-front/commit/3e32e55b1faf587a73dccb537823449c96753474) | Merge pull request #9 from code-corhuila/test/shell-unit-tests |
| 2026-10-08 | code-corhuila/lms-front | [`e680a4b`](https://github.com/code-corhuila/lms-front/commit/e680a4b18416e79e013146c23c99ec9b6e88998b) | Merge pull request #12 from code-corhuila/chore/front-connection-stress-perf-tests |
| 2026-10-08 | code-corhuila/lms-front | [`2ed1d0b`](https://github.com/code-corhuila/lms-front/commit/2ed1d0b4e3c54bd2c173d0b5ef6578188b4ffed2) | Merge pull request #10 from code-corhuila/promote-qa/shell-unit-tests |
| 2026-10-09 | code-corhuila/lms-catalog-api | [`8bad842`](https://github.com/code-corhuila/lms-catalog-api/commit/8bad8429dd8e31f19a1dafa9dcf48bd0b3f81077) | fix(catalog): add PATCH /books/{id} for HU-09 book editing |
