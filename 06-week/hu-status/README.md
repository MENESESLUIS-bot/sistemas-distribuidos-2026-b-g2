<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       06-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 06

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Draft the OpenAPI contracts for the access-service (authentication) and the API gateway routing behavior.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| API contracts | Access Service (auth) + API Gateway OpenAPI contracts | doing | [access-service.yaml](access-service.yaml); [api-gateway.yaml](api-gateway.yaml) (uncommitted, working tree) |

## 2. My individual contribution
- Authored `access-service.yaml`: the Access bounded-context contract with `GET /health` and `POST /auth/login`, including `LoginRequest`/`LoginResponse` schemas, the generic 401 on any credential failure, and explicit notes on what does not apply (no registration, refresh token, logout, JWKS, or role claims — single seeded Administrator, HS256 shared-secret JWT).
- Authored `api-gateway.yaml`: documents the NGINX gateway as a pure reverse proxy with no injected headers and no RBAC, its passthrough `/health`, and the real routing table (`/api/v1/auth` -> access-service, `/api/v1/students` -> membership-service, everything else -> library-api).

## 3. Blockers and risks
- Both files reference `_shared.yaml` and sibling contracts (`membership-service.yaml`, `library-api.yaml`) that are not present in this `06-week/hu-status` folder.
- `access-service.yaml` and `api-gateway.yaml` are still untracked/uncommitted in git — no PR or commit evidence yet.

## 4. Plan for next week
- Add the missing service contracts (`membership-service.yaml`, `library-api.yaml`), commit these drafts, and open the corresponding PR.

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- [Access Service contract](access-service.yaml)
- [API Gateway contract](api-gateway.yaml)
- Status: both files currently untracked (not yet committed)
