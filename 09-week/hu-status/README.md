<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       09-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 09

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Luis Alejandro Meneses
- GITHUB_USER: MENESESLUIS-bot
- TEAM: LMS-Library
- SPRINT_GOAL: Build the lms-catalog-portal micro-frontend (Vite + React + TypeScript): expose its routes as a Module Federation remote for the lms-front shell, implement the HU-04/HU-05/HU-09 catalog screens, and containerize it behind nginx.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| catalog-portal | Project scaffold: Vite + React 19 + TypeScript (strict) + Tailwind | done | [package.json](package.json); [index.html](index.html); [tsconfig.json](tsconfig.json) (commit 60d220d) |
| catalog-portal | Module Federation remote exposing `./routes` to the lms-front shell | done | [vite.config.ts](vite.config.ts); [routes.tsx](src/routes.tsx); [shell.d.ts](src/shell.d.ts); [main.tsx](src/main.tsx) (commits 60d220d, 0b71eca, 3f53ab7, 290f1ea) |
| HU-04 | Register book form | done | [BookFormPage.tsx](src/pages/books/BookFormPage.tsx) (commit 694f2a8) |
| HU-05 / HU-09 | Search the catalog and edit a book's details | doing | [BooksListPage.tsx](src/pages/books/BooksListPage.tsx); [types/index.ts](src/types/index.ts) (commits 82a3082, 450a76e) |
| catalog-portal | Shared UI components and design tokens | done | [components/ui/](src/components/ui/); [index.css](src/index.css) (commits 7993999, ec1f698, 00d1057, 8f2c721) |
| catalog-portal | Containerization: Docker build + nginx + compose | done | [Dockerfile](deploy/Dockerfile); [nginx.conf](deploy/nginx.conf); [compose.yml](deploy/compose.yml) (commits ba98f99, 177ce96, 3e9d154) |
| catalog-portal | Repo hygiene: `.gitignore` (no `.env`/`*.pem`) and `.env.example` | done | [.gitignore](.gitignore); [.env.example](.env.example) (commit 60d220d) |

## 2. My individual contribution
- Authored `package.json` for `lms-catalog-portal`: React 19, `react-router-dom` 7, Tailwind CSS 4 (via `@tailwindcss/vite`), `clsx` and the Inter font as runtime dependencies; Vite 8, TypeScript 6, `@module-federation/vite` and `oxlint` as dev dependencies, with `dev`/`build`/`lint`/`preview` scripts (`build` runs `tsc -b` before `vite build`).
- Authored `vite.config.ts`: registers the React and Tailwind plugins and configures Module Federation as the `catalog_portal` remote — it serves `remoteEntry.js`, exposes `./routes` (`src/routes.tsx`) so the lms-front shell can lazy-load the portal, consumes the `shell` remote (`http://localhost:3000/remoteEntry.js`) for `shell/apiClient` and `shell/session`, and shares `react`/`react-dom`/`react-router-dom` so they are not bundled twice. Dev server runs on port 3002.
- Authored the TypeScript project references: `tsconfig.json` points to `tsconfig.app.json` (browser code in `src`, `strict` mode, bundler resolution, `react-jsx`, unused-locals/params checks) and `tsconfig.node.json` (type-checks `vite.config.ts` against Node types).
- Authored `index.html`: the Vite entry page ("LMS — Catalog Portal") mounting the app into `#root` from `/src/main.tsx`.
- Authored `.gitignore` (logs, `node_modules`, `dist`, editor files, and always `.env`/`*.pem`) and `.env.example`, which documents that the portal reads no environment variables of its own — the gateway URL, token and correlation ID come from the shell's `apiClient` via Module Federation.
- Authored `src/routes.tsx`: `catalogRoutes`, the route table exposed to the shell — the index route renders `BooksListPage` and `new` renders `BookFormPage`. Paths are relative, so the portal does not depend on the prefix where the shell mounts it.
- Authored `src/main.tsx`, `src/bootstrap.tsx` and `src/App.tsx`: `main.tsx` dynamically imports `bootstrap` (required so the shared Module Federation dependencies load before React is used), `bootstrap.tsx` mounts the app in `StrictMode` inside a `BrowserRouter`, and `App.tsx` builds a standalone route tree from `catalogRoutes` so the portal can be previewed without the shell.
- Authored `src/shell.d.ts`: ambient type declarations for the two modules the shell exposes — `shell/apiClient` (`get`/`post`/`patch` plus the `ShellError` shape with per-field `details`) and `shell/session` (`getToken`, `isAuthenticated`, `clearSession`).
- Authored `src/types/index.ts`: the `Book` type mirroring the catalog OpenAPI schema, and the `Paginated<T>`/`PaginatedMeta` envelope the list page expects from `GET /books`.
- Authored `src/pages/books/BookFormPage.tsx` (HU-04): the register-book form (title, author, ISBN, category, year, total copies). It posts to `/books` with an `Idempotency-Key` generated once per registration intent and reused on retries, maps the API's per-field `details` to accessible inline errors (`aria-invalid`/`aria-describedby`), and returns to the list on success.
- Authored `src/pages/books/BooksListPage.tsx` (HU-05, HU-09): catalog table with search by title/author/ISBN, an availability badge per book, and inline editing of title/author/category/year through `PATCH /books/{id}` (ISBN stays read-only, matching the domain rule). It handles the four view states (loading, error with retry, empty, data) and discards stale responses so a newer search always replaces an older one.
- Authored `src/components/ui/`: `Button` (primary/secondary/danger/ghost variants with a loading spinner), `Card` and `EmptyState`. They are duplicated here temporarily until lms-front exposes shared UI components (the risk flagged in ADR-006).
- Authored `src/index.css`: Tailwind entry point with the LMS design tokens (indigo primary palette, success/warning/error colors, Inter typeface).
- Authored `deploy/Dockerfile`: two-stage build — `node:22-alpine` runs `npm ci` and `npm run build`, then `nginx:1.27-alpine` serves the `dist` output on port 80.
- Authored `deploy/nginx.conf`: SPA fallback to `index.html`, long-lived immutable caching for hashed assets, and `no-store` for `remoteEntry.js` and the HTML so the shell always loads the current portal version.
- Authored `deploy/compose.yml`: the `catalog-portal` service built from the Dockerfile, published on host port 3002 and attached to the external `lms-network`.

## 3. Blockers and risks
- `GET /books` does not exist in lms-catalog-api yet, so the search screen (HU-05) cannot load data end to end; `Paginated<Book>` describes the contract the page expects once the endpoint is added.
- The portal depends on the lms-front shell exposing `remoteEntry.js` on port 3000; without it running, `shell/apiClient` and `shell/session` cannot resolve at runtime.
- `deploy/Dockerfile` copies `package-lock.json`, but the repository's root `.gitignore` excludes that file, so the image build fails on a clean clone until the lockfile is tracked.
- `deploy/nginx.conf` applies `no-store` to `/assets/remoteEntry.js`, while `vite.config.ts` names the file `remoteEntry.js`; the emitted path still has to be verified against a real build.
- No tests or test runner configured yet for the portal.

## 4. Plan for next week
- Add `GET /books` (search + pagination) to lms-catalog-api and verify HU-05/HU-09 end to end through the gateway.
- Load the portal as a remote inside lms-front, verify the Docker image builds and serves `remoteEntry.js` correctly, and track the lockfile.
- Add a test setup (e.g. Vitest + Testing Library) covering the form and list pages, and open the corresponding HU PRs.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- [package.json](package.json) / [index.html](index.html)
- [vite.config.ts](vite.config.ts) (Module Federation remote)
- [tsconfig.json](tsconfig.json) / [tsconfig.app.json](tsconfig.app.json) / [tsconfig.node.json](tsconfig.node.json)
- [.gitignore](.gitignore) / [.env.example](.env.example)
- [Routes](src/routes.tsx) / [Entry point](src/main.tsx) / [Bootstrap](src/bootstrap.tsx) / [App](src/App.tsx)
- [Pages](src/pages/books/) / [UI components](src/components/ui/) / [Types](src/types/index.ts) / [Shell types](src/shell.d.ts) / [Styles](src/index.css)
- [Deploy](deploy/) — [Dockerfile](deploy/Dockerfile), [nginx.conf](deploy/nginx.conf), [compose.yml](deploy/compose.yml)
- Commit: `60d220d` — "chore(catalog-portal): scaffold Vite + React micro-frontend with Module Federation"
- Commit: `3f53ab7` — "feat(catalog-portal): add ambient types for shell apiClient and session remotes"
- Commit: `450a76e` — "feat(catalog-portal): add Book and Paginated contract types"
- Commit: `8f2c721` — "style(catalog-portal): add Tailwind entry with LMS design tokens"
- Commit: `7993999` — "feat(catalog-portal): add Button ui component with variants and loading state"
- Commit: `ec1f698` — "feat(catalog-portal): add Card ui component"
- Commit: `00d1057` — "feat(catalog-portal): add EmptyState ui component"
- Commit: `694f2a8` — "feat(catalog-portal): add book registration form page (HU-04)"
- Commit: `82a3082` — "feat(catalog-portal): add catalog search and inline edit page (HU-05, HU-09)"
- Commit: `0b71eca` — "feat(catalog-portal): expose catalog routes for the shell"
- Commit: `b128261` — "feat(catalog-portal): add standalone App route tree"
- Commit: `8eb5136` — "feat(catalog-portal): add React bootstrap with BrowserRouter"
- Commit: `290f1ea` — "feat(catalog-portal): add async entry point for Module Federation"
- Commit: `177ce96` — "build(catalog-portal): add nginx config with remoteEntry no-store caching"
- Commit: `ba98f99` — "build(catalog-portal): add multi-stage Node build and nginx runtime Dockerfile"
- Commit: `3e9d154` — "build(catalog-portal): add compose service on lms-network"
