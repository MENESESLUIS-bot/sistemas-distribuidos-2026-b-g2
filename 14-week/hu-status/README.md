<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       14-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 14

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Continue `lms-access-portal`, the Vue login remote that `lms-front` loads through Module Federation: add the standalone entry point (`src/main.ts`) and `App.vue` wrapper so the portal can run outside the shell, the Vite/Vue type declarations, the ambient types of the `shell` remotes and the HU-01 login view (`LoginView.vue`), the `useSession` composable, the federated `mountAccessLogin` (`mount.ts`) and the standalone router, and a first production build (`dist/`) that publishes `remoteEntry.js` with the `mountAccessLogin` export; plus the nginx config and compose file that serve that build (committed under `13-week/hu-status/deploy/`).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W14-01 | Standalone entry (`src/main.ts`), `App.vue` wrapper and `src/env.d.ts` | done | [d4b4a53](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d4b4a53) |
| W14-02 | nginx config (no-cache `remoteEntry.js`, immutable `/assets/`) and `compose.yml` for the portal | done | [deb0b54](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/deb0b54) |
| W14-03 | Production build of the federated remote (`dist/`, not committed) | done | [.gitignore](.gitignore) |
| W14-04 | Ambient types for `shell/apiClient` and `shell/session` (`src/shims/shell.d.ts`) | done | [ddc8b74](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ddc8b74) |
| W14-05 | HU-01 login view (`src/views/LoginView.vue`): form, submit to `/auth/login`, generic error | done | [8b96202](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8b96202), [edad5ef](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/edad5ef) |
| W14-06 | `useSession` composable that wraps `shell/session` (`src/composables/useSession.ts`) | done | [c0db426](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/c0db426) |
| W14-07 | Federated `mountAccessLogin` (`src/mount.ts`) and standalone router (`src/router.ts`) | done | [e309165](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e309165) |
| W14-08 | Portal manifest, `vite.config.ts` and tsconfig in this folder so it builds here | todo |  |

## 2. My individual contribution
- `src/main.ts`: standalone-dev entry; creates the Vue app with `App` and `router` and mounts it on `#app`. It is only used by `npm run dev` outside the shell; when `lms-front` composes the portal it loads `mount.ts` instead.
- `src/App.vue`: minimal wrapper with a single `<router-view />`, so the portal can be previewed in isolation.
- `src/env.d.ts`: `vite/client` types and a `*.vue` module declaration, so `vue-tsc` type-checks the `.vue` imports.
- `src/shims/shell.d.ts`: ambient modules for what `lms-front` exposes at runtime: `shell/apiClient` (a minimal `post<T>()` interface, so the portal does not need `axios` just for a type, and the `ShellError` envelope `{error, message, details?, traceId?}`) and `shell/session` (`getToken`, `setToken`, `isAuthenticated`, `clearSession`).
- `src/views/LoginView.vue` (HU-01, Vue port of the former React `LoginPage.tsx`, ADR-011):
  - Receives `onAuthenticated` as a prop and calls it right away if a session already exists.
  - Username and password fields (`required`, proper `autocomplete`); on submit it posts to `/auth/login` through the shell `apiClient`, stores `accessToken` with `setToken` and calls `onAuthenticated`.
  - The submit button is disabled while the request runs; on failure it shows the API `message` in a `role="alert"` box, without guessing whether the username or the password was wrong (INV-002).
  - Scoped BEM styles (`access-login__*`) so the portal's CSS does not leak into the React shell.
- `src/composables/useSession.ts`: the only file that imports `shell/session`; it returns `getToken`, `setToken`, `isAuthenticated` and `clearSession`, so the rest of the portal never depends on the federated module directly.
- `src/mount.ts`: exposed to `lms-front` through Module Federation. A Vue remote cannot export a React `RouteObject[]` like the other three portals, so it exports `mountAccessLogin(el, onAuthenticated)`, which creates a Vue app with `LoginView`, mounts it on `el` and returns an unmount function (ADR-011).
- `src/router.ts`: standalone-dev routing only, used by `main.ts`; route `/` renders `LoginView` with an `onAuthenticated` that just logs `[standalone] authenticated`. Without `lms-front` running on port 3000 the `shell/*` imports cannot resolve, so the login call fails in this mode.
- Build output (`dist/`, generated locally, ignored by git):
  - `remoteEntry.js` is the federated entry of `access_portal`; `@mf-types/compiled-types/src/mount.d.ts` publishes `mountAccessLogin(el, onAuthenticated): () => void` for the shell.
  - `assets/` holds content-hashed chunks (`LoginView-*.js/.css`, shared-import map, `loadRemote` stubs for `shell/apiClient` and `shell/session`).
  - `index.html` is titled `LMS — Access Portal` and preloads the shell remote from `http://localhost:3000/remoteEntry.js`.
- `.gitignore`: `node_modules`, `dist`, the `.mf` cache, `.env` and `*.pem`, so the build is not committed.
- `13-week/hu-status/deploy/nginx.conf`: `remoteEntry.js` with `Cache-Control: no-store` (H-front rule 4, so a deployed fix always reaches the browser), `/assets/` cached for one year as `immutable`, CORS `*` on both, SPA fallback to `index.html`, and a local-only `/api/v1/` proxy to `access-service:8080` while `lms-api-gateway` does not exist.
- `13-week/hu-status/deploy/compose.yml`: service `access-portal` on port `3004:80`, on the external `lms-network`; it runs the portal alone and is not the `lms-infra` composition.

## 3. Blockers and risks
- The portal source is now complete in `src/`, but it still cannot be built from this folder (see the next point), so `dist/` was produced with the week 13 configuration and has not been rebuilt from this source.
- The standalone mode (`npm run dev` without the shell) cannot really log in, because `shell/apiClient` and `shell/session` only resolve when `lms-front` runs on port 3000.
- If the shell `apiClient` rejects with a raw Axios error instead of the unwrapped `ShellError`, the view shows Axios's generic message instead of the API one.
- This folder has no `package.json`, `vite.config.ts`, `index.html` or tsconfig files; they are in `13-week/hu-status`, so `npm run dev` / `npm run build` cannot run from `14-week/hu-status`.
- `compose.yml` builds `deploy/Dockerfile`, but that Dockerfile is the Node image of `lms-access-api` (`EXPOSE 8080`, `node src/index.js`) and does not copy `dist/` into nginx, while compose maps port `80`. The portal needs its own nginx Dockerfile.
- The shell URL is still fixed to `http://localhost:3000/remoteEntry.js` in the build, so it needs a per-environment value before any deploy outside local development.
- CORS `*` on `remoteEntry.js` and `/assets/` is acceptable for static files but should be narrowed to the `lms-front` origin once it is known.
- `14-week/01-session` and `14-week/02-session` are empty.

## 4. Plan for next week
- Bring the portal manifest, `vite.config.ts`, `index.html` and tsconfig files next to `src/` so `npm run build` reproduces `dist/` from this source.
- Write an nginx multi-stage Dockerfile for the portal (`node:22-alpine` build → `nginx:alpine` with `nginx.conf` and `dist/`) and point `compose.yml` to it.
- Make the shell `remoteEntry.js` URL configurable per environment.
- Run the login view end to end against `POST /api/v1/auth/login` of `lms-access-api` and add a test for `mountAccessLogin` (mount, `onAuthenticated` callback, unmount).

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables (`.env` and `*.pem` ignored)

## 6. Evidence links
- Folder: [14-week/hu-status](./)
- Files: [src/main.ts](src/main.ts), [src/App.vue](src/App.vue), [src/env.d.ts](src/env.d.ts), [src/shims/shell.d.ts](src/shims/shell.d.ts), [src/views/LoginView.vue](src/views/LoginView.vue), [src/composables/useSession.ts](src/composables/useSession.ts), [src/mount.ts](src/mount.ts), [src/router.ts](src/router.ts), [.gitignore](.gitignore), [nginx.conf](../../13-week/hu-status/deploy/nginx.conf), [compose.yml](../../13-week/hu-status/deploy/compose.yml)
- Commits: [deb0b54](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/deb0b54), [d4b4a53](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/d4b4a53), [ddc8b74](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ddc8b74), [8b96202](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8b96202), [edad5ef](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/edad5ef), [c0db426](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/c0db426), [e309165](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/e309165)
