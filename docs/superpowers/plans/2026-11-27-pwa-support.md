# PWA Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `omeg4mon-router` a full installable Progressive Web App (Desktop, Android WebAPK, and iOS Add to Home Screen) with build-time Service Worker app-shell precaching, strict network-only API boundaries, bilingual offline/update/install UX, and proper Go embedded static HTTP headers.

**Architecture:** Zero-dependency native PWA integrated into the existing Vite 8 + React 19 SPA (`web/`) and Go embedded static server (`internal/platform/httpserver/static_ui.go`). A small local Vite plugin generates `web/out/sw.js` with a deterministic build-hash cache key and precache manifest from the emitted bundle, while `web/src/shared/pwa/` manages SW lifecycle, install prompts, and offline/update banners inside `AppShell`.

**Tech Stack:** Go 1.24 (`net/http`, `io/fs`), Vite 8 (Rolldown), React 19, TypeScript 5.8, Tailwind CSS v4, Zustand, Bun test, Biome.

**Spec:** [docs/superpowers/specs/2026-11-27-pwa-support-design.md](../specs/2026-11-27-pwa-support-design.md)

## Global Constraints

- Zero new npm or Go dependencies (`package.json` and `go.mod` remain dependency-clean).
- Must work identically in `vite build` (`web/out`) and Go single-binary `make build` (`-tags embedui` -> `cmd/omeg4mon-router/ui_dist`).
- Strict **Network-Only** for `/api/*`, `/v1/*`, `/v1beta/*`, and all non-GET requests in the Service Worker.
- Never force-reload active sessions on SW update; prompt user and reload only after explicit `SKIP_WAITING`.
- All UI copy must support both `en` and `id` locales and use semantic CSS variables from `web/src/app/globals.css` (Light `#f5f6f8` / Dark `#050608`).
- Every `useEffect` listener (`online`, `offline`, `beforeinstallprompt`, `appinstalled`, `controllerchange`, `updatefound`, `statechange`) must clean up on unmount.

## Review Focus

1. **Stale Service Worker across single-binary upgrades:** Browser caching `sw.js` would keep serving old hashed `/assets/*` paths after upgrading the Go binary. Pinned in Task 1 (`sw.js` serves `Cache-Control: no-store, no-cache, must-revalidate` + `Service-Worker-Allowed: /`) and Task 3 (SW `activate` purges old cache prefixes).
2. **API / SSE / auth requests intercepted by SW:** Caching `/api/*`, `/v1/*`, `/v1beta/*`, or non-GET requests would break auth cookies and LLM streaming. Pinned in Task 3 (`shouldHandleFetch` / routing unit tests).
3. **Offline navigation on nested SPA routes:** Reloading `/dashboard/providers/detail` offline must fall back to cached `/index.html`, not a 404 or network error. Pinned in Task 3 (navigation fallback strategy test).
4. **Android WebAPK installability requirements:** Missing `application/manifest+json` MIME type, missing 192/512/maskable PNG icons, or broken `start_url`/`scope` silently disables Android Chrome install. Pinned in Task 1 (`manifest.webmanifest` header test) and Task 2 (manifest & PNG asset validation test).
5. **Android status bar `theme-color` drift:** Toggling between Light (`#f5f6f8`) and Dark (`#050608`) themes must update `<meta name="theme-color">` immediately so the standalone Android status bar matches the active theme. Pinned in Task 2 (`uiStore` theme-color sync test).

---

### Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`)

**Files:**
- Modify: `internal/platform/httpserver/static_ui.go:12-135`
- Test: `internal/platform/httpserver/server_test.go`

**Interfaces:**
- Consumes: Existing `NewStaticUIHandler(fsys fs.FS) http.Handler` in `internal/platform/httpserver/static_ui.go`.
- Produces: Explicit cache & MIME headers for `sw.js` (`Cache-Control: no-store, no-cache, must-revalidate`, `Service-Worker-Allowed: /`, `Content-Type: text/javascript; charset=utf-8`) and `manifest.webmanifest` (`Cache-Control: no-cache`, `Content-Type: application/manifest+json`).

- [ ] **Step 1: Write the failing test**

In `internal/platform/httpserver/server_test.go`, add a subtest inside `TestServerIsolatesJSON404FromSPAFallback` (or a dedicated `TestStaticUIHandlerPWAHeaders(t *testing.T)`) with `fstest.MapFS` containing `sw.js` and `manifest.webmanifest`:
- `GET /sw.js` -> status `200`, `Cache-Control` contains `no-store` and `no-cache`, `Service-Worker-Allowed == "/"`, `Content-Type` contains `javascript`.
- `GET /manifest.webmanifest` -> status `200`, `Cache-Control` contains `no-cache`, `Content-Type` contains `application/manifest+json`.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -count=1 -race -run TestStaticUIHandlerPWAHeaders ./internal/platform/httpserver/...`
Expected: FAIL (missing `Cache-Control`, `Service-Worker-Allowed`, or `application/manifest+json`).

- [ ] **Step 3: Implement PWA header rules in `internal/platform/httpserver/static_ui.go`**

Update `applyStaticCacheHeaders(h http.Header, servedName string)` (or add `applyStaticResponseHeaders`) so:
- `base := strings.ToLower(path.Base(servedName))`
- When `servedName == "sw.js"` (case-insensitive): set `Cache-Control` to `htmlCacheControl` (`no-store, no-cache, must-revalidate`), `Service-Worker-Allowed` to `/`, and `Content-Type` to `text/javascript; charset=utf-8`.
- When `strings.HasSuffix(base, ".webmanifest")`: set `Cache-Control` to `"no-cache"` and `Content-Type` to `"application/manifest+json"`.
- Preserve existing rules for `assets/*`, `_next/static/*`, and `*.html`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -count=1 -race ./internal/platform/httpserver/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/platform/httpserver/static_ui.go internal/platform/httpserver/server_test.go
git commit -m "feat(httpserver): add PWA cache and MIME headers for sw.js and manifest"
```

---

### Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`

**Files:**
- Create: `web/public/icons/pwa-192x192.png`
- Create: `web/public/icons/pwa-512x512.png`
- Create: `web/public/icons/pwa-maskable-512x512.png`
- Create: `web/public/icons/apple-touch-icon.png`
- Create: `web/public/manifest.webmanifest`
- Modify: `web/index.html`
- Modify: `web/src/shared/store/uiStore.ts`
- Test: `web/src/shared/pwa/manifest.test.ts`

**Interfaces:**
- Consumes: `web/public/icon.svg` visual geometry (`#050608` dark canvas, `#eb6583` Caduceus emblem) and `web/src/shared/store/uiStore.ts` (`ThemeMode`, `THEME_COLORS`).
- Produces:
  - `web/public/manifest.webmanifest` with `id: "/"`, `name: "OMEG4MON ROUTER"`, `short_name: "OMEG4MON"`, `start_url: "/dashboard/"`, `scope: "/"`, `display: "standalone"`, `theme_color: "#050608"`, `background_color: "#050608"`, `icons` (192 any, 512 any, 512 maskable), and `shortcuts` (Overview, Chat, Providers, Usage).
  - `THEME_COLOR_BY_MODE: Record<ThemeMode, string>` (`light: "#f5f6f8"`, `dark: "#050608"`) exported from `web/src/shared/store/uiStore.ts` and synced to `<meta name="theme-color">` on boot and on every `setTheme` / `toggleTheme` / `hydrateTheme`.

- [ ] **Step 1: Write the failing test**

Create `web/src/shared/pwa/manifest.test.ts` testing:
- `web/public/manifest.webmanifest` parses as valid JSON with `display === "standalone"`, `start_url === "/dashboard/"`, `scope === "/"`, `name === "OMEG4MON ROUTER"`, and at least 3 icon entries (`192x192` any, `512x512` any, `512x512` maskable) whose referenced files exist on disk under `web/public/` with valid PNG magic bytes (`89 50 4E 47`).
- `web/public/icons/apple-touch-icon.png` exists with valid PNG magic bytes.
- `web/index.html` contains `rel="manifest" href="/manifest.webmanifest"`, `rel="apple-touch-icon"`, `name="theme-color"`, `viewport-fit=cover`, and `apple-mobile-web-app-capable`.
- `THEME_COLOR_BY_MODE` maps `light -> "#f5f6f8"` and `dark -> "#050608"`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && bun test src/shared/pwa/manifest.test.ts`
Expected: FAIL (manifest and icons do not exist yet).

- [ ] **Step 3: Generate PNG icons, create manifest, update `web/index.html` and `uiStore.ts`**

- Generate crisp 4x-supersampled PNG icons from the exact `icon.svg` vector geometry using Python 3 + Pillow:
  - `web/public/icons/pwa-192x192.png` (192x192)
  - `web/public/icons/pwa-512x512.png` (512x512)
  - `web/public/icons/pwa-maskable-512x512.png` (512x512, full-bleed `#050608` background with logo scaled to 72% centered safe zone for Android adaptive icons)
  - `web/public/icons/apple-touch-icon.png` (180x180, full-bleed `#050608` background for iOS)
- Create `web/public/manifest.webmanifest`.
- Update `web/index.html`:
  - `viewport` meta adds `viewport-fit=cover`.
  - Add `<meta name="theme-color" content="#f5f6f8" />`, `<meta name="mobile-web-app-capable" content="yes" />`, `<meta name="apple-mobile-web-app-capable" content="yes" />`, `<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />`, `<meta name="apple-mobile-web-app-title" content="OMEG4MON" />`.
  - Add `<link rel="manifest" href="/manifest.webmanifest" />` and `<link rel="apple-touch-icon" sizes="180x180" href="/icons/apple-touch-icon.png" />`.
  - In the inline pre-paint theme script, also update `meta[name="theme-color"]` to `#050608` when `t === "dark"` and `#f5f6f8` otherwise.
- Update `web/src/shared/store/uiStore.ts`: export `THEME_COLOR_BY_MODE` and update `meta[name="theme-color"]` inside `applyThemeToDocument` and `hydrateTheme`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd web && bun test src/shared/pwa/manifest.test.ts && bun run typecheck`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add web/public/icons web/public/manifest.webmanifest web/index.html web/src/shared/store/uiStore.ts web/src/shared/pwa/manifest.test.ts
git commit -m "feat(web): add PWA manifest, adaptive icons, and dynamic theme-color sync"
```

---

### Task 3: Build-Time Service Worker & Vite 8 Precache Plugin

**Files:**
- Create: `web/src/shared/pwa/swStrategy.ts`
- Create: `web/pwaPlugin.ts`
- Modify: `web/vite.config.ts`
- Test: `web/src/shared/pwa/swStrategy.test.ts`

**Interfaces:**
- Consumes: Vite 8 `generateBundle` output bundle filenames.
- Produces:
  - `classifyFetchRequest(url: URL, method: string, mode: string, origin: string): "ignore" | "network-only" | "navigate" | "asset"` in `web/src/shared/pwa/swStrategy.ts`.
  - `buildServiceWorkerSource(cacheName: string, precacheUrls: readonly string[]): string` in `web/pwaPlugin.ts` and `omeg4monPwaPlugin(): Plugin` wired into `web/vite.config.ts` emitting `web/out/sw.js`.

- [ ] **Step 1: Write the failing test**

Create `web/src/shared/pwa/swStrategy.test.ts` asserting:
- Cross-origin requests -> `"ignore"`.
- Non-GET requests (`POST`, `PUT`, `DELETE`) -> `"network-only"`.
- Same-origin `/api/auth/session`, `/api/system/health`, `/v1/chat/completions`, `/v1beta/models`, `/codex/responses`, `/responses`, `/models` -> `"network-only"`.
- Same-origin `/sw.js` -> `"network-only"`.
- Same-origin navigation (`mode === "navigate"`) to `/dashboard/`, `/dashboard/providers/detail`, `/login/` -> `"navigate"`.
- Same-origin static assets (`/assets/index-abc.js`, `/icons/pwa-192x192.png`, `/manifest.webmanifest`, `/icon.svg`) -> `"asset"`.
- `buildServiceWorkerSource("omeg4mon-router-v1-abc123", ["/index.html", "/assets/app.js"])` produces valid JS containing the cache name, precache URLs, `SKIP_WAITING` message listener, old-cache cleanup (`omeg4mon-router-` prefix), and strict network-only bypass for `/api/`, `/v1/`, `/v1beta/`.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && bun test src/shared/pwa/swStrategy.test.ts`
Expected: FAIL (module not found).

- [ ] **Step 3: Implement `swStrategy.ts`, `web/pwaPlugin.ts`, and wire into `web/vite.config.ts`**

- `web/src/shared/pwa/swStrategy.ts`: pure `classifyFetchRequest` function with explicit `NETWORK_ONLY_PREFIXES` (`/api`, `/v1`, `/v1beta`, `/codex`, `/responses`, `/models`).
- `web/pwaPlugin.ts`:
  - `buildServiceWorkerSource(cacheName, precacheUrls)` generates a clean, self-contained `sw.js`:
    - `install`: opens `cacheName`, precaches `PRECACHE_URLS` (deduped). Does NOT call `skipWaiting()` automatically (waits for user confirmation via `{ type: "SKIP_WAITING" }`).
    - `activate`: deletes any cache starting with `"omeg4mon-router-"` except `cacheName`, then calls `self.clients.claim()`.
    - `message`: on `{ type: "SKIP_WAITING" }`, calls `self.skipWaiting()`.
    - `fetch`:
      - ignores non-GET, cross-origin, `/sw.js`, and API/gateway prefixes (`/api`, `/v1`, `/v1beta`, `/codex`, `/responses`, `/models`).
      - `navigate`: network-first; on network failure falls back to cached `/index.html` (or `/`).
      - static assets (`/assets/*`, `/icons/*`, `/providers/*`, `/icon.svg`, `/favicon.ico`, `/manifest.webmanifest`): cache-first, populating cache on 200 basic responses.
  - `omeg4monPwaPlugin()`: Vite `Plugin` (`apply: "build"`) using `generateBundle` to collect `/index.html`, `/manifest.webmanifest`, `/icon.svg`, `/icons/pwa-192x192.png`, `/icons/pwa-512x512.png`, `/icons/pwa-maskable-512x512.png`, plus all emitted `assets/*` chunk/asset paths, compute a short SHA-256 hash over the sorted list, and emit `sw.js` via `this.emitFile({ type: "asset", fileName: "sw.js", source })`.
- Register `omeg4monPwaPlugin()` in `web/vite.config.ts`.

- [ ] **Step 4: Run test and build to verify `web/out/sw.js` is emitted cleanly**

Run: `cd web && bun test src/shared/pwa/swStrategy.test.ts && bun run build && ls -lh out/sw.js out/manifest.webmanifest out/icons/`
Expected: PASS and `out/sw.js` exists with precached `/assets/*` entries.

- [ ] **Step 5: Commit**

```bash
git add web/src/shared/pwa/swStrategy.ts web/src/shared/pwa/swStrategy.test.ts web/pwaPlugin.ts web/vite.config.ts
git commit -m "feat(web): add Vite 8 Service Worker precache plugin with network-only API boundary"
```

---

### Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification

**Files:**
- Create: `web/src/shared/pwa/usePwaStatus.ts`
- Create: `web/src/shared/pwa/PwaControls.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Test: `web/src/shared/pwa/pwaCopy.test.ts`

**Interfaces:**
- Consumes: `useUIStore` (`locale`), `Button` (`@/shared/ui/button`), `Alert` (`@/shared/ui/alert`), `navigator.serviceWorker`, `BeforeInstallPromptEvent`.
- Produces:
  - `usePwaStatus(): { isOffline: boolean; updateAvailable: boolean; canInstall: boolean; isInstalling: boolean; applyUpdate: () => void; dismissUpdate: () => void; promptInstall: () => Promise<void> }`
  - `getPwaCopy(locale: UILocale)` (pure bilingual EN/ID copy helper).
  - `PwaInstallButton` (compact header action shown only when `canInstall` is true) and `PwaStatusBanner` (non-intrusive banner above `<main>` content when offline or when a new version is ready to reload).

- [ ] **Step 1: Write the failing test**

Create `web/src/shared/pwa/pwaCopy.test.ts` testing:
- `getPwaCopy("en")` and `getPwaCopy("id")` return distinct, non-empty strings for `installLabel`, `installingLabel`, `offlineTitle`, `offlineDescription`, `updateTitle`, `updateAction`, `dismissAction`.
- `isStandaloneDisplayMode()` helper returns boolean safely when `window` is undefined (SSR/test environment).

- [ ] **Step 2: Run test to verify it fails**

Run: `cd web && bun test src/shared/pwa/pwaCopy.test.ts`
Expected: FAIL (module not found).

- [ ] **Step 3: Implement `usePwaStatus.ts`, `PwaControls.tsx`, and wire into `AppShell.tsx`**

- `web/src/shared/pwa/usePwaStatus.ts`:
  - Registers `/sw.js` with `{ scope: "/" }` only in production (`import.meta.env.PROD`) when `"serviceWorker" in navigator` and `window.isSecureContext` (prevents dev-server SW caching interference while keeping prod + localhost/HTTPS working).
  - Tracks `waitingWorker` when `registration.waiting` exists or `updatefound` -> `installed` with existing `navigator.serviceWorker.controller`.
  - `applyUpdate()` posts `{ type: "SKIP_WAITING" }` to the waiting worker and reloads once on `controllerchange`.
  - Tracks `online` / `offline` window events.
  - Captures `beforeinstallprompt` (calling `event.preventDefault()`) and clears prompt on `appinstalled` or after `promptInstall()`.
  - Cleans up every event listener on unmount.
- `web/src/shared/pwa/PwaControls.tsx`:
  - `PwaInstallButton`: renders compact `Button` (`variant="secondary"`, `size="sm"`, `Download` icon from `lucide-react`) only when `canInstall` is true.
  - `PwaStatusBanner`: renders crisp semantic-token banner (offline warning or update-available action bar with Reload + Dismiss) only when `isOffline || updateAvailable`.
- `web/src/shared/ui/AppShell.tsx`:
  - Call `usePwaStatus()` once in `AppShell`, render `<PwaInstallButton />` beside `<LocaleToggle />` / `<ThemeToggle />` in the header, and `<PwaStatusBanner />` at the top of `<main>`.

- [ ] **Step 4: Run full project verification gate (`make check` and `make build`)**

Run: `make check && make build`
Expected: Exit code 0 across `gofmt`, `golangci-lint`, `biome check`, `go test -race`, `CGO_ENABLED=0` db tests, `bun test`, `tsc --noEmit`, `cli` tests, `web` production build, and `embedui` Go binary build.

- [ ] **Step 5: Update knowledge graph & commit**

Run: `graphify update .` (per `GEMINI.md`) and commit Task 4 changes.
