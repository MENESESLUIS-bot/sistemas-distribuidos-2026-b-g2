<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Write the first `node --test` tests for `lms-access-api`: the `Administrator` aggregate of the Access domain (password hashing and authentication) and the `Login` use case, driven through fake ports so the domain and application layers are tested with no database or JWT I/O.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W13-01 | `Administrator` domain tests: no plaintext password, authenticate ok / wrong password | done | [83d426c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/83d426c) |
| W13-02 | `Login` use case tests with fake `administrators` and `tokens` ports | done | [dfe630b](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/dfe630b) |
| W13-03 | Source under test (`src/domain/access/administrator.js`, `errors.js`, `src/application/usecase/login.js`) | todo |  |

## 2. My individual contribution
- `test/domain/administrator.test.js` (3 tests): `newAdministrator('alice', ...)` never keeps the plaintext password and stores a bcrypt hash (`$2...` prefix); `authenticate` does not throw on the correct password; `authenticate` throws `ErrInvalidCredentials` on a wrong one.
- `test/application/login.test.js` (3 tests): the `Login` use case receives its ports by constructor (`administrators.findByUsername`, `tokens.issue`), replaced here by in-memory fakes, so no PostgreSQL or JWT key is needed.
  - Happy path returns a non-empty `accessToken` and `expiresIn = 3600`.
  - Unknown username (`ErrAdministratorNotFound` from the repository) and wrong password both reject with the same `ErrInvalidCredentials`, so the API does not reveal which usernames exist.
- Both files use only the built-in `node:test` and `node:assert/strict`, matching the `test` script (`node --test`) from the week 12 `package.json`.

## 3. Blockers and risks
- The tests import `src/domain/access/administrator.js`, `src/domain/access/errors.js` and `src/application/usecase/login.js`, which are not in the folder yet, so `npm test` fails with module-not-found until the source is added.
- `13-week/01-session` and `13-week/02-session` are still empty.

## 4. Plan for next week
- Implement `newAdministrator` / `authenticate` with `bcryptjs` and the `ErrInvalidCredentials` / `ErrAdministratorNotFound` errors so the domain tests pass.
- Implement the `Login` use case that maps "not found" to "invalid credentials", plus a PostgreSQL adapter for `administrators` and a JWT adapter for `tokens`.
- Run `npm test` and `npm run lint` green and add an integration test for the HTTP login endpoint.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Folder: [13-week/hu-status](./)
- Files: [test/domain/administrator.test.js](test/domain/administrator.test.js), [test/application/login.test.js](test/application/login.test.js)
- Commits: [83d426c](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/83d426c), [dfe630b](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/dfe630b)
