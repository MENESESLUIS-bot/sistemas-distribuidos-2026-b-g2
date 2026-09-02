# ADR-001 — Hexagonal Modular Monolith for MVP 1

| Field | Value |
| --- | --- |
| Status | Accepted |
| Date | 2026-08-22 |
| Authors | Oscar Areiza, Hermes Pascuas, Luis Alejandro Meneses |

## Context

LMS Library has four bounded contexts: Access, Membership, Catalog, and Circulation. The MVP must be delivered by a three-person academic team using Go, PostgreSQL, Docker, and a React frontend. Expected traffic and data volume do not justify independent scaling in the first release.

## Decision

Build one deployable Go service named `library-api`, organized as a hexagonal modular monolith. Each bounded context is an internal module with its own domain, application, and infrastructure code. The modules use explicit ports and adapters and do not access another module's persistence tables directly. PostgreSQL is shared in MVP 1, with logical ownership separated by module.

## Alternatives considered

| Alternative | Reason not selected |
| --- | --- |
| Four independent microservices with four databases and a broker | Operational complexity and distributed failure modes exceed the MVP team's capacity and current scaling needs. |
| Traditional layered monolith | Makes domain rules dependent on framework and persistence concerns and weakens future extraction boundaries. |
| Hexagonal modular monolith | Selected: preserves bounded-context boundaries, supports isolated domain tests, and leaves a path to future extraction. |

## Consequences

- MVP deployment has one API process, one PostgreSQL instance, and a simpler local environment.
- Domain logic remains testable without I/O.
- Shared PostgreSQL is technical debt; table ownership and repository boundaries must be enforced.
- A context can later become a service with its own database and event transport without rewriting its domain model.
