<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       07-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 07

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Define the shared OpenAPI schema and the service contract template used to standardize API design across LMS services.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| API contracts | Shared OpenAPI components + per-service contract template | doing | [_shared.yaml](_shared.yaml); [_template-service.yaml](_template-service.yaml) (commit 8db6b5c) |

## 2. My individual contribution
- Authored `_shared.yaml`: reusable OpenAPI components shared by every service contract, including `PaginatedMeta`, `ErrorResponse`, `UUID`, `Timestamp`, and `AuditFields` schemas, the `IdParam` parameter, standard error responses (`BadRequest`, `Unauthorized`, `NotFound`, `InternalError`), and the `bearerAuth` JWT security scheme.
- Authored `_template-service.yaml`: a base OpenAPI 3.0.3 template for per-service contracts, covering server variables, default bearer security, a `/health` endpoint, and full CRUD paths (list/create/get/update/delete) with request/response schemas referencing `_shared.yaml`.

## 3. Blockers and risks
- `_template-service.yaml` still contains placeholders (`[NombreServicio]`, `[recurso]`, `[campo1]`, etc.) — no concrete per-service contract has been instantiated from it yet in this folder.

## 4. Plan for next week
- Instantiate the template into real per-service OpenAPI contracts (e.g. access-service, catalog-service, membership-service), replacing placeholders with the actual resources and fields for each service.

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- [Shared OpenAPI components](_shared.yaml)
- [Service contract template](_template-service.yaml)
- Commit: `8db6b5c` — "Document drafts 07"
