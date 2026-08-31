# Traceability Matrix — LMS-LIBRARY-V1

---

Traceability connects every line of code to its business justification. It allows answering: "Why does this function exist?" and "Which HU covers this part of the system?" It also identifies unimplemented requirements and code without a requirement (possible technical debt).

> **Status note:** LMS-LIBRARY-V1 has not started implementation yet — every row below is 🔴 **Pending** by design. Test file names are proposed (Go convention), not yet written.

---

## FR → HU → Test → Service matrix

| FR ID | FR Description | HU | Tests that verify it | Service¹ | Status |
|---|---|---|---|---|---|
| FR-001 | Authenticate the administrator (JWT session) | HU-01 | `auth_authenticate_test.go` | Auth | 🔴 Pending |
| FR-002 | Reject invalid credentials with a generic error | HU-01 | `auth_authenticate_test.go` | Auth | 🔴 Pending |
| FR-003 | Register a new student | HU-02 | `students_register_test.go` | Students | 🔴 Pending |
| FR-004 | Reject duplicate student document ID | HU-02 | `students_register_test.go` | Students | 🔴 Pending |
| FR-005 | Search and edit a student's contact information | HU-03 | `students_update_test.go` | Students | 🔴 Pending |
| FR-006 | Block deactivation of a student with active loans/suspension | HU-03 | `students_deactivate_test.go` | Students | 🔴 Pending |
| FR-007 | Register a new book | HU-04 | `catalog_register_book_test.go` | Catalog | 🔴 Pending |
| FR-008 | Reject duplicate ISBN on registration | HU-04 | `catalog_register_book_test.go` | Catalog | 🔴 Pending |
| FR-009 | Search catalog by title/author/ISBN/category | HU-05 | `catalog_search_books_test.go` | Catalog | 🔴 Pending |
| FR-010 | Return empty result set on no matches | HU-05 | `catalog_search_books_test.go` | Catalog | 🔴 Pending |
| FR-011 | Edit a registered book's information | HU-09 | `catalog_update_book_test.go` | Catalog | 🔴 Pending |
| FR-012 | Keep ISBN immutable through the edit action | HU-09 | `catalog_update_book_test.go` | Catalog | 🔴 Pending |
| FR-013 | Register a loan (eligible student + available book, due date +7 days) | HU-06 | `loans_register_test.go` | Loans | 🔴 Pending |
| FR-014 | Decrement availability on loan registration | HU-06 | `loans_register_test.go` | Loans | 🔴 Pending |
| FR-015 | Reject loan for suspended/over-limit student | HU-06 | `loans_register_test.go` | Loans | 🔴 Pending |
| FR-016 | Register a return, close the loan | HU-07 | `returns_register_test.go` | Returns | 🔴 Pending |
| FR-017 | Increment availability on return | HU-07 | `returns_register_test.go` | Returns | 🔴 Pending |
| FR-018 | Detect and record a late return | HU-07 | `returns_register_test.go` | Returns | 🔴 Pending |
| FR-019 | Apply a fixed 7-day suspension on late return | HU-08 | `penalties_suspend_test.go` | Penalties | 🔴 Pending |
| FR-020 | Query all currently overdue loans | HU-08 | `penalties_overdue_test.go` | Penalties | 🔴 Pending |

¹ "Service" reflects the **logical module** grouping used in `functional.md`, not a confirmed microservice boundary. Final service decomposition is pending the architecture stage (`05-architecture`, if applicable) — see `02-domain/domain-map.md` for the current single-service (`library-api`) decision.

**Coverage check:** all 20 FRs have an HU and a proposed test file; all 9 HUs (HU-01..HU-09) are represented.

---

## NFR → Validation matrix

| NFR ID | Description | How it is validated | Tool | Status |
|---|---|---|---|---|
| NFR-001 | Performance (P95 < 300ms on critical endpoints) | Load test in staging | k6 | 🔴 Pending |
| NFR-002 | Availability (95% monthly SLO) | Health-check + uptime monitoring | Uptime monitor / manual | 🔴 Pending |
| NFR-003 | Scalability (graceful behavior under 2x burst) | Load test at target concurrency | k6 / Locust | 🔴 Pending |
| NFR-004 | Security (auth, HTTPS, hashing, OWASP Top 10) | Access-control tests + SAST scan | SonarQube/Snyk, manual review | 🔴 Pending |
| NFR-005 | Observability (structured logs, metrics, traces) | Log/metric inspection during test execution | Prometheus + Grafana | 🔴 Pending |
| NFR-006 | Maintainability (≥ 80% coverage, complexity ≤ 10) | Coverage report + static analysis in CI | CI pipeline (Go tooling) | 🔴 Pending |
| NFR-007 | Portability (Docker images, no host OS dependency) | Clean-environment container build/run | Docker | 🔴 Pending |
| NFR-008 | Disaster Recovery (RTO/RPO targets) | Recovery drill | Manual | 🔴 Pending |

**Coverage check:** all 8 NFR categories have a defined validation path.

---

## Inverse traceability: HU → FR

| HU | Title | FR(s) it implements | Sprint |
|---|---|---|---|
| HU-01 | Administrator Authentication | FR-001, FR-002 | Sprint 1 |
| HU-02 | Student Registration | FR-003, FR-004 | Sprint 1 |
| HU-03 | Student Search, Editing & Deactivation | FR-005, FR-006 | Sprint 3 |
| HU-04 | Book Registration | FR-007, FR-008 | Sprint 1 |
| HU-05 | Inventory Control & Book Search | FR-009, FR-010 | Sprint 3 |
| HU-06 | Loan Registration | FR-013, FR-014, FR-015 | Sprint 2 |
| HU-07 | Return Registration & History Tracking | FR-016, FR-017, FR-018 | Sprint 2 |
| HU-08 | Late-Return Penalty System | FR-019, FR-020 | Sprint 4 |
| HU-09 | Book Editing | FR-011, FR-012 | Sprint 3 |

> Sprint assignments here match `user-stories.md` — see that file for the reasoning (dependency order + priority).

---

## Status legend

| Status | Meaning |
|---|---|
| ✅ Done | Implemented, tested, and in production |
| 🟡 In progress | Under development in the current sprint |
| 🔴 Pending | In the backlog, not started |
| ⏸ Blocked | Has an external blocker |
| ❌ Cancelled | Removed from scope |

---

## Identified gaps (requirements without coverage)

No structural gaps at this stage: every FR has an HU, every HU has a proposed test file, and every NFR has a defined validation path. All items are 🔴 Pending because implementation has not started — this is expected, not a gap. Re-run this check once development begins; a genuine gap would look like an FR without an HU, an HU without a test, or a test without an implementation file.

| Gap type | Description | Required action | Owner | Date |
|---|---|---|---|---|
| — | None identified as of this document's last update | — | — | 2026-08-26 |

---

## How to maintain this matrix

- When an HU is created: add the row in the FR → HU → Test → Service section.
- When a test is written: note the file in the "Tests that verify it" column.
- When an HU is completed: change the status to ✅.
- At each Sprint Planning: review gaps and assign actions.

## Correlations

- Functional Requirements → `04-requirements/functional.md`
- User Stories → `04-requirements/user-stories.md`
- Non-Functional Requirements → `04-requirements/non-functional.md`
- Testing strategy → `11-quality/testing-strategy.md`
- DoD that determines when an HU is Done → `00-governance/definition-of-done.md`
