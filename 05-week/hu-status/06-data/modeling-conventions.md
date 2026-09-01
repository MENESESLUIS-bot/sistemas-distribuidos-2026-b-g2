# Data Modeling Conventions — LMS-LIBRARY-V1

> Naming and style rules applied consistently across `library-api`'s schema
> (`06-data/models.md`). Any new table or column must follow these.

---

## Naming

| Artifact | Convention | Example |
|----------|-----------|---------|
| Tables | `snake_case`, plural | `students`, `loans` |
| Columns | `snake_case`, descriptive | `document_id`, `available_copies` |
| Foreign keys | `[referenced_table_singular]_id` | `student_id`, `book_id` |
| Indexes | `idx_[table]_[column(s)]` | `idx_loans_student_id` |
| Booleans | prefixed to read as a question where possible | `was_late` |
| Timestamps | always `TIMESTAMPTZ`, never bare `TIMESTAMP` (per ADR-001, English + explicit timezone) | `created_at`, `due_date` |

---

## Identifiers

- Every table's primary key is a `UUID` (`gen_random_uuid()`), matching the `Identifier: id (UUID)`
  field recorded for every entity in `02-domain/entities-and-rules.md`.
- Business keys (`students.document_id`, `books.isbn`, `administrators.username`) are enforced
  as `UNIQUE` but are **not** the primary key — the surrogate `id` is, so a business key can be
  corrected (e.g., a mistyped ISBN) without cascading a primary-key change through `loans`.

## Standard audit fields

Every table includes:

```sql
created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
```

`updated_at` is refreshed by the application layer on every write (no DB trigger in v1 — one
more moving part the team doesn't need to maintain at this scale).

## Soft delete — only where the domain requires it

Unlike the general scaffold default of soft-deleting everything, LMS-LIBRARY-V1 only
soft-deletes where the domain explicitly calls for it:

| Table | Soft delete? | Field | Why |
|-------|-------------|-------|-----|
| `students` | Yes | `deactivated_at` | HU-03 requires "deactivation," and a hard delete would orphan loan history |
| `books` | No | — | No HU deletes a book in v1 (only register/edit) |
| `loans` | No | — | A loan is never deleted; its terminal state is `status = 'RETURNED'` |
| `administrators` | No | — | Provisioned directly in the DB, not through a user-facing flow |

**Rule going forward:** don't add `deleted_at` to a table unless a real HU requires deleting
that kind of record — an unused soft-delete column is dead weight, not a safety net.

## `created_by` / `updated_by` — not used in v1

The generic scaffold convention suggests `created_by`/`updated_by` audit columns. LMS-LIBRARY-V1
omits them: there is exactly one Administrator account operating the system in v1
(`01-context/scope.md`), so "who made this change" is never ambiguous. Revisit only if a future
version introduces multiple Administrator accounts (see `01-context/scope.md`, Candidates for
future versions).

---

## Correlations

- Applied schema → `06-data/models.md`
- Entities these conventions serialize → `02-domain/entities-and-rules.md`
- Documentation language rule (English, incl. DB naming) → `05-architecture/decisions/records/ADR-001-idioma-documentacion.md`
