# CP4 — Security gates and negative authorization test

## What is here

| File | Purpose |
|---|---|
| `security.yml` | GitHub Actions workflow proposed for the team pipeline: `gitleaks`, `govulncheck`, negative authorization tests, on every `pull_request` to `develop`, `qa`, `main` |
| `lms-catalog-api/` | minimal module of `lms-catalog-api` (Go 1.25, chi, golang-jwt): the team's `RequireAuth` plus the proposed `RequireScope`, and the `/api/v1/books` routes |
| `lms-catalog-api/internal/httpapi/authz_negative_test.go` | the negative tests |

The workflow lives in this folder as a proposal; it is not copied into the code repository
of the team today. A copy in `.github/workflows/security.yml` of this fork lets it run for
real on GitHub Actions (see "CP8 — Recorded results").

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

## CP8 — Recorded results (2026-10-30)

### Negative authorization test

Run on my machine against the minimal `lms-catalog-api` module (Windows 11, go1.25.1):

```text
$ go vet ./...
$ go test -count=1 ./internal/httpapi/ -run 'Cannot|401|RightToken' -v
=== RUN   TestCirculationTokenCannotDeleteBook
--- PASS: TestCirculationTokenCannotDeleteBook (0.02s)
=== RUN   TestAdministratorTokenCannotTakeCopyWithoutLoan
--- PASS: TestAdministratorTokenCannotTakeCopyWithoutLoan (0.05s)
=== RUN   TestMissingTokenIsStill401
--- PASS: TestMissingTokenIsStill401 (0.01s)
=== RUN   TestRightTokenOnRightRoute
--- PASS: TestRightTokenOnRightRoute (0.02s)
PASS
ok  	github.com/MENESESLUIS-bot/w13/lms-catalog-api/internal/httpapi	0.513s
```

To check that the test can really fail, I removed the `RequireScope` lines from `router.go`
in a temporary copy. Both negative tests went red with the exact threat of row E:

```text
--- FAIL: TestCirculationTokenCannotDeleteBook (0.03s)
    authz_negative_test.go:76: DELETE /api/v1/books/{id} with the internal token = 200, want 403
--- FAIL: TestAdministratorTokenCannotTakeCopyWithoutLoan (0.09s)
    authz_negative_test.go:86: POST /api/v1/books/{id}/loan-copy with the Administrator token = 200, want 403
FAIL
```

### Workflow validation

```text
$ actionlint -version
v1.7.7
$ actionlint security.yml
(no output, exit code 0)
```

### Real run on GitHub Actions

Run: <https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/actions/runs/41092565678>
(manual `workflow_dispatch`, commit `2127a4c`). **The run is red because one gate found a
real problem.**

| Gate | Result | Detail |
|---|---|---|
| `negative authorization tests` | pass | the 4 tests above, on `ubuntu-latest` |
| `govulncheck` | pass on Go 1.26.9 | on Go 1.25.1 it exits with code 3: 8 reachable vulnerabilities in the standard library (`net/http`, `crypto/x509`, `net/url`). The `net/http` ones are only fixed in 1.26.9, so the workflow uses that version |
| `gitleaks` | **fail** | 1 finding, see below. On `workflow_dispatch` it scans the whole history; on `pull_request` it scans only the commits of the PR |

`gitleaks` first reported 3 findings:

| Finding | Decision |
|---|---|
| `authz_negative_test.go:19`, fixed test secret | false positive, listed in `.gitleaksignore` |
| `06-week/hu-status/access-service.yaml:80`, example `accessToken` in the OpenAPI contract | false positive, listed in `.gitleaksignore` |
| `12-week/hu-status/.env.example:16`, `JWT_PRIVATE_KEY` | **true positive, not ignored** |

The third finding is a local development key, but it is the pair of the `JWT_PUBLIC_KEY` in the
`.env.example` of `lms-catalog-api`, `lms-circulation-api` and `lms-membership-api` (same
SHA-256 of the public key, `07e222221f90dd63`), and this fork is public. Any environment
that kept that example public key accepts Administrator tokens signed by anyone who reads
the file. It cannot be removed from the history without rewriting it, so the fix is to
generate a new pair for every shared environment and replace the example values.
