# Candidate Architecture — LMS Library

## Bounded contexts

| Bounded context | Business responsibility | Candidate service | Decision for MVP 1 |
| --- | --- | --- | --- |
| Access | Administrator authentication and session validation | `access-service` | Keep as an internal module of `library-api` |
| Membership | Student registration, search, editing, and deactivation | `membership-service` | Keep as an internal module of `library-api` |
| Catalog | Book registration, editing, search, and availability counters | `catalog-service` | Keep as an internal module of `library-api` |
| Circulation | Loans, returns, overdue detection, and suspension rules | `circulation-service` | Keep as an internal module of `library-api` |

## Candidate options

1. **Four independent microservices:** one service and database per bounded context, communicating through REST and domain events.
2. **Traditional layered monolith:** one deployable application with shared controllers, application services, and repositories.
3. **Hexagonal modular monolith:** one deployable Go API with one isolated module per bounded context and explicit ports and adapters.

## Evaluation

The four-service option gives independent deployment and scaling, but it adds four deployables, multiple databases, a broker, distributed failure modes, and operational overhead that is disproportionate for a three-person academic team and the expected low-volume MVP. The traditional monolith is simpler initially but couples business rules to infrastructure and makes later extraction harder.

## Candidate architecture

The recommended candidate is a **hexagonal modular monolith**. Each bounded context keeps its own domain model, application use cases, ports, and adapters. Modules share the deployable process but do not access another module's tables directly. PostgreSQL is shared in MVP 1 as an explicitly documented trade-off; ownership remains separated by module.

## Decision input

The team will formalize this recommendation in `decisions/ADR-001-architecture.md` during Session 2.
