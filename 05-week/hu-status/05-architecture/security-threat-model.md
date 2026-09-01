# Security Threat Model (STRIDE)

> STRIDE analysis of LMS-LIBRARY-V1: Spoofing, Tampering, Repudiation, Information Disclosure,
> Denial of Service, Elevation of Privilege. Scope: the single Administrator-facing system
> described in `05-architecture/overview.md` (NGINX → `library-api` → PostgreSQL).

---

## Spoofing (pretending to be someone/something else)

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| An attacker guesses/brute-forces the Administrator's password | Administrator account | Passwords hashed with bcrypt (cost ≥ 12); generic error on failed login that doesn't reveal which field was wrong (`02-domain/entities-and-rules.md`, INV-002) | Confirmed in domain model |
| An attacker forges or steals a JWT to impersonate the Administrator | Every write operation in the system | Short-lived tokens (1h expiry, NFR-004), HTTPS-only transmission so tokens can't be sniffed in transit | Confirmed in NFR-004 |
| A client pretends to be `library-api` to a browser (man-in-the-middle) | All data in transit | TLS terminated at `nginx` (`ADR-003`), HTTP disallowed outside local dev | Confirmed in NFR-004 / ADR-003 |

---

## Tampering (unauthorized modification of data)

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| Direct manipulation of `availableCopies` bypassing loan/return validation | Book aggregate | `availableCopies` only changes through `Book.loanOneCopy()` / `Book.returnOneCopy()`, never written directly (INV-001, `02-domain/entities-and-rules.md`) | Confirmed in domain model |
| A tampered request tries to set a loan's `dueDate` directly, bypassing the fixed 7-day rule | Loan aggregate | `dueDate` is computed in the Loan factory, never exposed as a settable field (INV-001 on Loan) | Confirmed in domain model |
| SQL injection via unsanitized input (e.g., book search, student registration) | PostgreSQL database | Parameterized queries only (pgx), no string-concatenated SQL; input validation at the HTTP adapter | Design principle — enforce in code review |

---

## Repudiation (denying having performed an action)

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| Administrator denies having registered/edited a loan, book, or student | Audit trail | Every entity has immutable `createdAt` and auto-updated `updatedAt` fields; every domain event (`02-domain/domain-events.md`) carries `metadata.userId` and `occurredAt` | Confirmed in domain model |
| No record of who changed what, if the system is ever opened to more than one Administrator account | Accountability | Out of scope in v1 (single Administrator account, `01-context/scope.md`); revisit if multi-admin is ever added (see "Candidates for future versions") | Accepted risk for v1 |

---

## Information Disclosure (exposing data to unauthorized parties)

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| Login error reveals whether the username or password was wrong | Administrator credentials | Single generic error message for both cases (INV-002, `02-domain/entities-and-rules.md`) | Confirmed in domain model |
| Student PII (email, phone, document ID) exposed via an unauthenticated endpoint | Student records | Every operation requires a valid Administrator session (`02-domain/domain-map.md`, Access context is OHS in front of all others) | Confirmed in domain model |
| Sensitive data intercepted in transit | All data | HTTPS mandatory in production (NFR-004), TLS terminated at `nginx` | Confirmed in ADR-003 |
| Secrets (JWT signing key, DB credentials) leaked via source control or logs | Backend secrets | Secrets only in environment variables, never in code or logs (NFR-004, `05-architecture/cross-cutting.md`) | Confirmed in NFR-004 |

---

## Denial of Service

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| Flood of requests to a single endpoint (e.g., login) exhausts `library-api` | System availability | Basic rate limiting configured at `nginx` (`limit_req`), per `ADR-003` and `05-architecture/deployment.md` | Confirmed in ADR-003 |
| Unbounded result sets (e.g., catalog search with no pagination) degrade performance | System availability | Pagination on list/search endpoints, finalized in `07-api/contracts/openapi/` | Design principle — enforce in API contract |
| A single `library-api` instance failing takes the whole system down | System availability | Accepted risk for v1 — no load balancing across multiple instances yet; NFR-003 allows scaling to up to 3 instances if needed later | Accepted risk for v1 |

---

## Elevation of Privilege

| Threat | Affected asset | Mitigation | Status |
|--------|----------------|-----------|--------|
| A forged/modified JWT claims Administrator privileges without a valid login | Every operation in the system | JWT signature verified server-side inside `library-api` on every request (NGINX does not perform this check, per `ADR-003`) | Confirmed in ADR-003 |
| Privilege escalation between roles | N/A | Not applicable in v1 — there is only one role, Administrator, with no role hierarchy (`00-governance/security-policy.md`, `02-domain/domain-map.md`); revisit only if a future version introduces RBAC | Not applicable in v1 |

---

## Out of scope for this threat model

- Multi-tenant isolation — v1 serves a single university library, not multiple tenants.
- Physical security of the hosting infrastructure — outside the team's control for a
  cloud/virtualized deployment (`01-context/scope.md`).
- Supply-chain attacks on third-party Go/npm dependencies — tracked instead under
  `04-requirements/non-functional.md` (NFR-004, dependency scanning via SAST tools).

---

## Correlations

- Domain invariants referenced above → `02-domain/entities-and-rules.md`
- Security NFRs (HTTPS, hashing, OWASP Top 10) → `04-requirements/non-functional.md`
- RBAC status (not applicable in v1) → `00-governance/security-policy.md`
- Reverse proxy / edge layer → `05-architecture/decisions/records/ADR-003-nginx-reverse-proxy.md`
- Security rules and thresholds → `00-governance/security-rules.md`
