# Migration Strategy — LMS-LIBRARY-V1

**Tool:** `golang-migrate` (`github.com/golang-migrate/migrate/v4`), per `_stacks/go.md`.

## File naming convention

```
migrations/
├── 000001_create_administrators_table.up.sql
├── 000001_create_administrators_table.down.sql
├── 000002_create_students_table.up.sql
├── 000002_create_students_table.down.sql
├── 000003_create_books_table.up.sql
├── 000003_create_books_table.down.sql
├── 000004_create_loans_table.up.sql
└── 000004_create_loans_table.down.sql
```

Sequential numeric prefix (`000001`, `000002`, ...), snake_case description, `.up.sql` /
`.down.sql` pair per change — matching `golang-migrate`'s expected format.

## Rules

```
✓ Migrations are ALWAYS forward-only in shared environments (dev/production) — a mistake is
  fixed with a new migration, not by editing an already-applied one
✓ One migration per logical change (one table, or one focused alter)
✓ Every .up.sql has a matching .down.sql that reverses it exactly
✗ Never edit a migration file that has already run in dev or production
✗ Never DROP COLUMN / DROP TABLE while code in production still reads it
    (process: 1 — stop reading/writing the column in code → 2 — drop it in the next migration)
```

## Compatible vs. incompatible changes

Same rules as the general scaffold guidance in `06-data/models.md`'s template section apply
here — additive changes (new nullable column, new index) are safe; renames/type changes need a
two-phase migration (add new → backfill → switch reads → drop old).

## Applying migrations

| Environment | How |
|-------------|-----|
| Local | `make migrate-up` (per the Makefile in `_stacks/go.md`) |
| Development / Production | Run as a one-off step in the deploy pipeline, **before** the new `library-api` image starts serving traffic (see `10-devops/ci-cd.md`) |

## Seed data

For local development, a `S001__seed_dev_data.sql` (or equivalent Go seed script) creates one
Administrator account and a handful of sample books/students — never real student PII, and
never applied outside `local`/`development`.

---

## Correlations

- Schema these migrations build → `06-data/models.md`
- Local setup command → `10-devops/local-setup.md`
- Deploy pipeline step that runs migrations → `10-devops/ci-cd.md`
- Go tooling reference → `_stacks/go.md`
