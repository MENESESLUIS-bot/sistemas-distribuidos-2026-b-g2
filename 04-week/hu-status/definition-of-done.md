# Definition of Done — MVP 1

A story is Done only when all applicable conditions are satisfied:

- Acceptance criteria are implemented and verified by automated or documented manual tests.
- Domain invariants are covered by unit tests; persistence behavior is covered by integration tests where applicable.
- The code follows the hexagonal boundary: domain code has no HTTP, database, framework, or environment I/O.
- API changes are reflected in `api-contract-openapi.yaml`.
- Database migrations are repeatable and tested against PostgreSQL in a container.
- Error responses are consistent, actionable, and do not expose secrets or sensitive internals.
- Code is reviewed and approved by at least one teammate.
- Conventional Commits and the branch/PR policy are followed.
- No credentials or private data are committed; configuration comes from environment variables.
- The feature runs through Docker Compose and does not break `/health` or `/health/ready`.
- Traceability is updated from HU to functional requirement and test evidence.
