# Data Ownership — LMS Library

| Entity / table | Owning context | Write authority | Consumers |
|---|---|---|---|
| `administrators` | Access | Access repositories and authentication use cases | API authorization middleware |
| `students` | Membership | Membership use cases | Circulation through eligibility port; reporting read models |
| `books` | Catalog | Catalog use cases and availability port | Circulation through reservation/release port; search API |
| `loans` | Circulation | Circulation use cases | Returns/history queries and penalty reporting |

## Shared database rule

MVP 1 uses one PostgreSQL database because the system is a modular monolith. This does not create shared ownership: each module may write only its own tables. Foreign keys preserve referential integrity, while application ports prevent business logic from bypassing module invariants.

## Future extraction

When independent services are needed, each owning context receives its own datastore. Cross-context foreign keys are replaced with identifiers and explicit contracts or events. The current logical ownership makes that migration possible without changing business vocabulary.
