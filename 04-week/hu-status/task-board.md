# MVP 1 Task Board

| Task | HU | Owner | Status | Acceptance check |
|---|---|---|---|---|
| Authentication endpoint and JWT middleware | HU-01 | Access | To Do | Valid login, generic invalid-login error, protected route |
| Student registration and duplicate validation | HU-02 | Membership | To Do | Valid student persists; duplicate document rejected |
| Book registration and availability initialization | HU-04 | Catalog | To Do | Valid book persists; duplicate ISBN rejected |
| Loan use case and transaction | HU-06 | Circulation | To Do | Due date +7 days; availability decremented atomically |
| Return use case and history | HU-07 | Circulation | To Do | Active loan becomes returned; availability incremented |
| Health endpoint and database readiness | Cross-cutting | Platform | In Progress | `/health` and `/health/ready` respond correctly |
| Unit and integration test suite | All MVP HUs | QA | To Do | Domain tests and real PostgreSQL integration tests pass |
