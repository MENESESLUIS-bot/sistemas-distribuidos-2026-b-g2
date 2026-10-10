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

### My commits and PRs this week (all repos)
Every commit and pull request authored by `MENESESLUIS-bot` between 2026-09-28 and 2026-10-04 (America/Bogota). Cherry-picks of the same commit into develop/qa/main are listed once, with the original SHA.

**Pull requests (7)**

| Opened | Repo | PR | Title | State |
|---|---|---|---|---|
| 2026-09-28 | code-corhuila/lms-catalog-portal | [#2](https://github.com/code-corhuila/lms-catalog-portal/pull/2) | chore: scaffold catalog portal with book pages and deploy files | merged |
| 2026-10-02 | code-corhuila/lms-catalog-portal | [#3](https://github.com/code-corhuila/lms-catalog-portal/pull/3) | feat(catalog): add book pages, ui components and deploy files | merged |
| 2026-10-03 | code-corhuila/lms-catalog-api | [#6](https://github.com/code-corhuila/lms-catalog-api/pull/6) | chore(qa): promote catalog service skeleton and Annex C fixes from develop | closed |
| 2026-10-03 | code-corhuila/lms-catalog-api | [#7](https://github.com/code-corhuila/lms-catalog-api/pull/7) | chore(release): promote develop to qa | closed |
| 2026-10-03 | code-corhuila/lms-catalog-api | [#8](https://github.com/code-corhuila/lms-catalog-api/pull/8) | qa: promote catalog service skeleton from develop | merged |
| 2026-10-03 | code-corhuila/lms-catalog-db | [#5](https://github.com/code-corhuila/lms-catalog-db/pull/5) | chore(qa): promote catalog-db Liquibase schema from develop | merged |
| 2026-10-03 | code-corhuila/lms-catalog-portal | [#4](https://github.com/code-corhuila/lms-catalog-portal/pull/4) | chore(qa): promote project scaffold and catalog pages from develop | merged |

**Commits (45)**

| Date | Repo | Commit | Message |
|---|---|---|---|
| 2026-09-28 | code-corhuila/lms-catalog-portal | [`d7bd273`](https://github.com/code-corhuila/lms-catalog-portal/commit/d7bd273a947281fb0dd4fa10baa854e284020ac1) | chore(deps): add npm lockfile |
| 2026-09-28 | code-corhuila/lms-catalog-portal | [`b78cbc2`](https://github.com/code-corhuila/lms-catalog-portal/commit/b78cbc20e8c78b460f1be6fecfd748992db91408) | chore(deps): add package manifest with vite, react and module federation |
| 2026-09-28 | code-corhuila/lms-catalog-portal | [`0f268d5`](https://github.com/code-corhuila/lms-catalog-portal/commit/0f268d54457bea3c15ee54423ce076949866eb12) | chore(config): add env example documenting no local variables |
| 2026-09-28 | code-corhuila/lms-catalog-portal | [`a462168`](https://github.com/code-corhuila/lms-catalog-portal/commit/a462168a776e96ec12462d1ac98e735a4b0136d3) | chore(config): add gitignore for node, vite and secrets |
| 2026-09-28 | code-corhuila/lms-catalog-portal | [`2d36587`](https://github.com/code-corhuila/lms-catalog-portal/commit/2d36587435d34311fb25e340753230e112a00ba9) | feat(portal): add html entry point for the catalog portal |
| 2026-09-29 | code-corhuila/lms-catalog-portal | [`32b5afb`](https://github.com/code-corhuila/lms-catalog-portal/commit/32b5afb93c047b1aed97b90974ddef4abf9502a7) | chore(config): add tsconfig for vite config |
| 2026-09-29 | code-corhuila/lms-catalog-portal | [`abd4dca`](https://github.com/code-corhuila/lms-catalog-portal/commit/abd4dca1c05927551810a92b5b0bcc5265eac914) | chore(config): add strict tsconfig for app sources |
| 2026-09-29 | code-corhuila/lms-catalog-portal | [`ad7a0a9`](https://github.com/code-corhuila/lms-catalog-portal/commit/ad7a0a9819579157c6d8c2b015e40d63bd60bc35) | chore(config): add root tsconfig with project references |
| 2026-09-29 | code-corhuila/lms-catalog-portal | [`8216347`](https://github.com/code-corhuila/lms-catalog-portal/commit/8216347a89d14acdedf4c30f27c0aa5c018f1e08) | chore(config): add vite config exposing routes via module federation |
| 2026-09-29 | code-corhuila/lms-catalog-portal | [`5c973fa`](https://github.com/code-corhuila/lms-catalog-portal/commit/5c973fa92918973a39610c78b46523ed9e0b20fe) | Merge pull request #2 from code-corhuila/chore/project-scaffold |
| 2026-09-29 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`60d220d`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/60d220d4d61ec7c55e129020319d11ead7e062c2) | chore(catalog-portal): scaffold Vite + React micro-frontend with Module Federation |
| 2026-09-29 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`0f77968`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0f779681b2c0f0b62fde372620df2164746c1e3c) | docs(hu-status): update week 09 status with catalog-portal scaffold |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`c12f7aa`](https://github.com/code-corhuila/lms-catalog-portal/commit/c12f7aa62178f15857d7285ec897605706d14bd5) | chore(config): ignore module federation diagnostics folder |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`a69aa34`](https://github.com/code-corhuila/lms-catalog-portal/commit/a69aa3438e2665ab4fedc059fe9e3f983c975930) | feat(ui): add tailwind theme with lms design tokens |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`7e724e6`](https://github.com/code-corhuila/lms-catalog-portal/commit/7e724e6d07febe16e98faf1f7a6fedc761d4d714) | feat(ui): add button component with variants and loading state |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`061a8df`](https://github.com/code-corhuila/lms-catalog-portal/commit/061a8df8716d7271b9cfcf56637e011e6a11d3f6) | feat(catalog): add book and pagination types |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`08468e6`](https://github.com/code-corhuila/lms-catalog-portal/commit/08468e6144ac94594c8bdd3e5059a20d896ddc5d) | feat(shell): add ambient types for the lms-front api client and session |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`8b53476`](https://github.com/code-corhuila/lms-catalog-portal/commit/8b534763c24f4864cc57b06aa97969eb4694f0fe) | feat(catalog): add book registration form page (HU-04) |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`342727f`](https://github.com/code-corhuila/lms-catalog-portal/commit/342727f58493822e2f336c620bb9c5a4ca7fc7b7) | feat(ui): add empty state component |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`d9a48ed`](https://github.com/code-corhuila/lms-catalog-portal/commit/d9a48edb1762bf6721df06d7bd63b5d76257225a) | feat(ui): add card component |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`f31186b`](https://github.com/code-corhuila/lms-catalog-portal/commit/f31186b7288d1baec01a86a5f041db1aede972dd) | feat(portal): add standalone app route tree |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`609304b`](https://github.com/code-corhuila/lms-catalog-portal/commit/609304b374510177f669ca15039eb0619e5301d0) | feat(catalog): add catalog routes exposed to the shell |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`c3a0629`](https://github.com/code-corhuila/lms-catalog-portal/commit/c3a0629144aedbffb16bad6d54b424fbce2d5a80) | feat(catalog): add books list page with search and inline edit (HU-05, HU-09) |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`d377d2a`](https://github.com/code-corhuila/lms-catalog-portal/commit/d377d2a3d2473bb30387c83b0c311ed349061ea0) | chore(deploy): add nginx config with no-store for remote entry |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`ac06e14`](https://github.com/code-corhuila/lms-catalog-portal/commit/ac06e1491afd5d6a799ad78994a8c41cab08bb40) | feat(portal): add entry point that loads bootstrap dynamically |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`210e683`](https://github.com/code-corhuila/lms-catalog-portal/commit/210e6838117b4467d070d5f5369b2038b2070101) | feat(portal): add react bootstrap with browser router |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`a007b82`](https://github.com/code-corhuila/lms-catalog-portal/commit/a007b821022b51e937327ca503de837ffa4cabf7) | chore(deploy): add compose service on the lms network |
| 2026-10-01 | code-corhuila/lms-catalog-portal | [`5cd5b74`](https://github.com/code-corhuila/lms-catalog-portal/commit/5cd5b74cf221be79c350931d4d8c0faf758648a3) | chore(deploy): add multi-stage dockerfile serving the build with nginx |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`450a76e`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/450a76e2b57bd5efdacaf508653414ede2dd07ba) | feat(catalog-portal): add Book and Paginated contract types |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`3f53ab7`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3f53ab7ce3e44f567e86696a87f9d6cb05a4fdd4) | feat(catalog-portal): add ambient types for shell apiClient and session remotes |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`ec1f698`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ec1f6983dd083576685337ad60f6134a3624e7fc) | feat(catalog-portal): add Card ui component |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`7993999`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/79939990c0244d95944c7814c697910cb71f4d69) | feat(catalog-portal): add Button ui component with variants and loading state |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`8f2c721`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8f2c721f3af164bee3ac640d4f46c849db5d500d) | style(catalog-portal): add Tailwind entry with LMS design tokens |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`694f2a8`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/694f2a8d9ee2267a111bc0549d184736908c39f3) | feat(catalog-portal): add book registration form page (HU-04) |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`00d1057`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/00d105771458aec1d0598a81c9417694777cf6ba) | feat(catalog-portal): add EmptyState ui component |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`0b71eca`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/0b71eca3ea1f5557d51346ee2516b0263844a958) | feat(catalog-portal): expose catalog routes for the shell |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`82a3082`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/82a30826360ddfcaee98b1b5940c1ffb4253b7ba) | feat(catalog-portal): add catalog search and inline edit page (HU-05, HU-09) |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`8eb5136`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8eb51369c48e061e20cf7fbeb41179d56eae1f7e) | feat(catalog-portal): add React bootstrap with BrowserRouter |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`b128261`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/b128261e4917fd899beb06da63153f4966c3d5d6) | feat(catalog-portal): add standalone App route tree |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`177ce96`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/177ce96b860f2318a5db22bb24df94c5b764582a) | build(catalog-portal): add nginx config with remoteEntry no-store caching |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`290f1ea`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/290f1ea8a54b8a54d47f4db551fbf4ddcbb46fa4) | feat(catalog-portal): add async entry point for Module Federation |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`3e9d154`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/3e9d1549eec5cd40765a207b999e74578e76d1f3) | build(catalog-portal): add compose service on lms-network |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`ba98f99`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/ba98f997d651c8a521a465ddbf37ca3b04c16ad4) | build(catalog-portal): add multi-stage Node build and nginx runtime Dockerfile |
| 2026-10-01 | MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2 | [`8aa18db`](https://github.com/MENESESLUIS-bot/sistemas-distribuidos-2026-b-g2/commit/8aa18dbfba9ef8cb085496101252695f6dea11b7) | docs(hu-status): update week 09 status with catalog-portal pages and deploy |
| 2026-10-03 | code-corhuila/lms-catalog-portal | [`6a7797f`](https://github.com/code-corhuila/lms-catalog-portal/commit/6a7797fa4e659bd88a01772fc6179efb567ec76b) | Merge pull request #3 from code-corhuila/chore/project-scaffold |
