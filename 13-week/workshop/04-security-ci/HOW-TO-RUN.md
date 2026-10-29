# CP4 — Security gates and negative authorization test

## What is here

| File | Purpose |
|---|---|
| `security.yml` | GitHub Actions workflow proposed for the team pipeline: `gitleaks`, `govulncheck`, negative authorization tests, on every `pull_request` to `develop`, `qa`, `main` |
| `lms-catalog-api/` | minimal module of `lms-catalog-api` (Go 1.25, chi, golang-jwt): the team's `RequireAuth` plus the proposed `RequireScope`, and the `/api/v1/books` routes |
| `lms-catalog-api/internal/httpapi/authz_negative_test.go` | the negative tests |

The workflow lives in this folder as a proposal; it is not copied into the code repository
of the team today.

## Story and criteria covered

Story: **[security] Per-route scopes on lms-catalog-api** (STRIDE row E, risk 9, P0 —
PR code-corhuila/library-docs#36). Boundary `lms-circulation-api` → `lms-catalog-api`.

| Test | Acceptance criterion |
|---|---|
| `TestCirculationTokenCannotDeleteBook` | `DELETE /api/v1/books/{id}` with the internal HS256 token of `lms-circulation-api` → **403**, the delete handler never runs |
| `TestAdministratorTokenCannotTakeCopyWithoutLoan` | `POST /api/v1/books/{id}/loan-copy` with the Administrator RS256 token → **403**, `available_copies` is not touched |
| `TestMissingTokenIsStill401` | no token → **401** (403 only for a valid token with the wrong scope) |
| `TestRightTokenOnRightRoute` | internal token on `loan-copy` → 200; Administrator token on `PATCH /books/{id}` → 200 |

## Run locally

```bash
cd 13-week/workshop/04-security-ci/lms-catalog-api
go mod download
go vet ./...
go test ./internal/httpapi/ -run 'Cannot|401|RightToken' -v
```

Requirements: Go 1.25.1. No database or running service is needed: the test generates an
RSA key pair for the Administrator token and signs the internal token with a test secret,
exactly like `catalogclient` does in `lms-circulation-api`.

The other two gates:

```bash
go install golang.org/x/vuln/cmd/govulncheck@v1.1.4 && govulncheck ./...
gitleaks detect --source ../../../.. --no-banner
```
