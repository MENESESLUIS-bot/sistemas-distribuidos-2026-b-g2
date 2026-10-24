<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       12-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 12

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Bootstrap `lms-access-api`, the Access bounded context (identity, roles and policies, ADR-011) on Node.js/Express, with its dependency manifest, every environment variable it reads documented in `.env.example`, a Docker build context that leaves secrets and tests out, and a gitleaks fingerprint for the throwaway dev key.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W12-01 | `lms-access-api` package manifest (Express, pg, JWT, bcrypt, pino, oxlint) | done | [708c135](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/708c135) |
| W12-02 | Config via environment variables, Docker ignore list and secret-scan fingerprint | done | [4504bb4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/4504bb4) |
| W12-03 | Service source (`src/index.js`, `src/config/config.js`) and tests (`test/**/*.test.js`) | todo |  |

## 2. My individual contribution
- `package.json`: `lms-access-api` 1.0.0 as an ES module (`"type": "module"`) on Node `>=22`; runtime dependencies `express` (HTTP), `pg` (PostgreSQL), `jsonwebtoken` (JWT signing), `bcryptjs` (password hashing) and `pino` (structured JSON logs); `oxlint` as the linter; scripts `dev` (`node --watch`), `start`, `test` (built-in `node --test`) and `lint`.
- `.env.example`: every variable the service reads from `src/config/config.js` - `PORT`, `DB_HOST`/`DB_PORT`/`DB_USER`/`DB_PASSWORD`/`DB_NAME` for `lms_db`, `JWT_PRIVATE_KEY` + `JWT_EXPIRY=1h`, `LOG_LEVEL` and `CORS_ORIGIN`; header with the `openssl genrsa` / `openssl rsa -pubout` commands to generate your own key pair and the one-line `\n` format (anexo C, numeral 5.3.7).
- `.dockerignore`: keeps `node_modules`, `test`, `.env` and `.git` out of the image build context, so local secrets and test files never reach the image.
- `.gitleaksignore`: a single fingerprint (`.env.example:private-key:16`) instead of a whole-file ignore, so the documented throwaway dev key passes the pre-commit secret scan while any new secret in that file is still caught.

## 3. Blockers and risks
- `src/` and `test/` are not in the folder yet, so `npm start`, `npm test` and `npm run lint` cannot run.
- `package-lock.json` exists locally but the root `.gitignore` excludes it, so installs of `lms-access-api` are not reproducible from the repo.
- `.env.example` carries a real (throwaway) RSA private key; GitHub secret scanning may flag it even though gitleaks has its fingerprint.

## 4. Plan for next week
- Add `src/config/config.js` that reads and validates the variables of `.env.example`, and `src/index.js` with the Express app and a health endpoint.
- Add the first `node --test` tests for config loading and JWT signing with the dev key.
- Add a `Dockerfile` that uses the `.dockerignore` above.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Folder: [12-week/hu-status](./)
- Files: [package.json](package.json), [.env.example](.env.example), [.dockerignore](.dockerignore), [.gitleaksignore](.gitleaksignore)
- Commits: [708c135](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/708c135), [4504bb4](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/4504bb4)
