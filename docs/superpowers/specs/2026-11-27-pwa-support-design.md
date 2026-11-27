# PWA Support Design — `omeg4mon-router`

**Date:** 2026-11-27  
**Scope:** Small–medium feature (frontend Vite 8 SPA + Go embedded static handler)

## Intent & Why

Make `omeg4mon-router` a first-class installable Progressive Web App (PWA) across Desktop (Chrome/Edge/Safari) and Mobile (Android WebAPK & iOS Add to Home Screen), delivering fast standalone app-shell startup, safe single-binary upgrades, and clear offline/update feedback without introducing external bundler dependencies or caching sensitive API/streaming responses.

## Goals, Non-Goals & Anti-Goals

- **Goals:**
  - Full Web App Manifest (`manifest.webmanifest`) with standalone display, Omegamon Zwart theme colors (`#050608` dark / `#f5f6f8` light), and app shortcuts.
  - Complete icon set generated from `web/public/icon.svg`: `192x192`, `512x512`, `maskable 512x512` (with Android adaptive-icon safe zone), and `180x180` `apple-touch-icon.png`.
  - Dynamic `<meta name="theme-color">` synced with the active light/dark theme (before first paint and on theme toggle), plus `viewport-fit=cover` for notched mobile screens.
  - Build-time Service Worker (`sw.js`) with hashed asset precaching (via a small local Vite 8 plugin), automatic old-cache cleanup on binary upgrades, and SPA navigation fallback to cached `index.html` when offline.
  - Strict **Network-Only** policy for `/api/*`, `/v1/*`, `/v1beta/*`, and all non-GET requests so auth cookies/sessions, SSE streams, and live telemetry are never cached or stale.
  - Bilingual (EN/ID) PWA UX in `AppShell`: offline status indicator, user-controlled update prompt (*"Update available — Reload"* via `SKIP_WAITING`), and Install App button when `beforeinstallprompt` fires.
  - Correct Go HTTP static headers in `internal/platform/httpserver/static_ui.go`: `Cache-Control: no-store, no-cache, must-revalidate` + `Service-Worker-Allowed: /` for `sw.js`, and `Content-Type: application/manifest+json` + `Cache-Control: no-cache` for `manifest.webmanifest`.
- **Non-Goals:**
  - Offline read-only API data persistence (IndexedDB / TanStack Query persister).
  - Web Push notifications / VAPID key management (Telegram & webhooks already cover alerts).
- **Anti-Goals:**
  - Never cache API responses, auth states, or SSE/LLM streams in the Service Worker.
  - Never force-reload an active user session mid-chat/mid-monitoring on SW update.
  - Never add `vite-plugin-pwa` / Workbox dependency bloat that risks Vite 8 Rolldown breakage.

## Constraints

- Must work seamlessly with `web/` Vite 8 (`rolldownOptions`), React 19, Tailwind CSS v4, and `bun`.
- Must work identically when embedded into the single Go binary via `make build` (`embedui` tag -> `cmd/omeg4mon-router/ui_dist`).
- Android/iOS/Desktop PWA installation and Service Worker registration require a Secure Context (`https://` or `http://localhost`). Accessing from Android over plain LAN HTTP (`http://192.168.x.x:20128`) requires HTTPS (reverse proxy, Tailscale HTTPS, or the built-in tunnel module).

## Approaches Considered

1. **Native Lightweight PWA (Custom SW + Local Vite Precache Plugin + Go Static Headers) — Selected**
   - Zero new npm/Go dependencies, 100% compatible with Vite 8 Rolldown, transparent ~100-line SW easy to unit-test, strict network-only API boundary.
2. **`vite-plugin-pwa` (Workbox)**
   - Rejected: heavy dependency tree, Vite 8/Rolldown compatibility risk, still requires manual API exclusions and Go static header updates.
3. **Offline-First Data + Web Push (Approach 1 + IndexedDB + VAPID)**
   - Rejected: high complexity and stale/post-logout data leakage risk for a live self-hosted router control plane.

## Builder Decisions & Deliverables

- **Icons & Manifest:** Static PNG assets under `web/public/icons/` (generated deterministically from `web/public/icon.svg`) + `web/public/manifest.webmanifest` + `<head>` links/meta in `web/index.html` (and dynamic `theme-color` sync in `ThemeToggle`/theme script).
- **Service Worker:** Source template + small Vite plugin in `web/vite.config.ts` that emits `sw.js` into `web/out/sw.js` with a build-hash cache version and the exact list of emitted static assets to precache.
- **Frontend PWA Module:** `web/src/shared/pwa/` (`usePwaStatus.ts` / registration helper + `PwaStatusBanner.tsx` / install action) wired into `AppShell` with EN/ID copy and clean unmount/event listener cleanup.
- **Backend Static Handler:** Update `internal/platform/httpserver/static_ui.go` and unit tests in `internal/platform/httpserver/server_test.go` to enforce `sw.js` no-cache + `Service-Worker-Allowed: /` and manifest MIME/cache headers.
- **Verification:** Go unit tests (`go test -race ./internal/platform/httpserver/...`), frontend unit tests & typecheck (`bun test && bun run typecheck && bun run check`), and full embedded build (`make build`).
