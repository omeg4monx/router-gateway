# Wave 3B — Media Gateway (Audio/Images/Videos/Search/Fetch) & Tunnel Manager Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 3B by adding two isolated vertical-slice modules — `media` (gateway endpoints `/v1/audio/*`, `/v1/images/*`, `/v1/videos/*`, `/v1/search`, `/v1/web/fetch` with encrypted media provider management, ordered media combos, automatic fallback, and SSRF-safe built-in web fetch) and `tunnel` (Cloudflare Quick Tunnel and Tailscale Serve/Funnel lifecycle manager with binary detection, loopback-only target enforcement, and lifecycle event audit log) — complete with OpenAPI 3.1 contracts, Next.js 16 static-export studios, E2E tests, and single-binary embedding.

**Architecture:** `media` owns its own SQLite tables (`media_providers`, `media_combos`), encrypts upstream provider API keys via `crypto.Vault`, authenticates `/v1/*` requests via `contracts.KeyValidator`, records telemetry via `contracts.UsageRecorder`, and mounts both AdminAPI (`/api/media/*`) and Gateway (`/v1/audio/*`, `/v1/images/*`, `/v1/videos/*`, `/v1/search`, `/v1/web/fetch`) routes with zero cross-module Go imports. `tunnel` owns its own SQLite tables (`tunnel_config`, `tunnel_events`), abstracts process execution behind an injectable `ProcessRunner` interface (defaulting to `os/exec` for `cloudflared` and `tailscale`), enforces loopback-only tunnel targets, parses live public URLs from process output, and terminates child processes cleanly on `Stop()`.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `crypto.Vault`, `eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded request/response body limits.
- **Security Posture:** Media provider API keys are encrypted at rest with AES-256-GCM (`crypto.Vault`) and never returned in API responses; `/v1/web/fetch` enforces strict SSRF blocking (loopback, RFC1918, link-local, CGNAT, multicast, non-HTTP schemes, and redirect hops); `tunnel` strictly restricts `local_url` to loopback hosts (`127.0.0.1`, `localhost`, `::1`).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task passes `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and final verification runs `make check`, `make build`, and `graphify update .`.

## Review Focus

1. **Media gateway fallback distinguishes client vs upstream errors:** When a media combo targets `[primary, secondary]`, upstream `429`/`5xx`/transport failures fail over to `secondary` and return `200`, whereas client `400`/`422` errors return immediately without dialing `secondary`. _(Pinned in Task 1 & Task 3)_
2. **Built-in `/v1/web/fetch` blocks SSRF on initial URL and redirect hops:** Requests targeting `127.0.0.1`, `localhost`, `169.254.169.254`, `10.0.0.1`, `file://`, or a public URL redirecting to `127.0.0.1` fail closed with `403` / `400`. _(Pinned in Task 1 & Task 3)_
3. **Tunnel `local_url` rejects non-loopback targets & cleans up child processes on Stop:** Configuring or starting a tunnel with a non-loopback host (e.g. `http://10.0.0.5:8080` or `http://169.254.169.254`) fails with `400`; calling `Stop()` terminates any running tunnel process cleanly without goroutine leaks. _(Pinned in Task 2)_
4. **Disabled optional modules cleanly detach:** Disabling `media` or `tunnel` via `OMEG4MON_DISABLED_MODULES` omits their routes without affecting core gateway traffic. _(Pinned in Task 3)_
5. **Static-export UI compliance & error surfacing:** `/dashboard/media/` and `/dashboard/tunnel/` export cleanly to static HTML (`web/out/`), surface all query/mutation errors via `Alert`, and contain zero `backdrop-blur` or UI emojis. _(Pinned in Task 4)_

---

### Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch)

**Files:**
- Create: `internal/modules/media/migrations/0001_init.sql`
- Create: `internal/modules/media/repo.go`
- Create: `internal/modules/media/fetch.go`
- Create: `internal/modules/media/module.go`
- Create: `internal/modules/media/module_test.go`
- Create: `contracts/openapi/modules/media.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*crypto.Vault`, `*http.Client`, `contracts.KeyValidator`, `contracts.UsageRecorder`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `media.New(opts Options) *Module` (`Name() = "media"`, `Dependencies() = []string{"system", "keys"}`)
  - Implements `lifecycle.Module`
  - SQLite tables (`0001_init.sql`):
    - `media_providers(id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, capability TEXT NOT NULL, driver TEXT NOT NULL, base_url TEXT NOT NULL DEFAULT '', api_key_enc BLOB, default_model TEXT NOT NULL DEFAULT '', extra_headers_json TEXT NOT NULL DEFAULT '{}', priority INTEGER NOT NULL DEFAULT 100, enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `media_combos(id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, capability TEXT NOT NULL, targets_json TEXT NOT NULL DEFAULT '[]', enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
  - AdminAPI routes (`/api`):
    - `GET /media/providers`, `POST /media/providers`, `PATCH /media/providers/{id}`, `DELETE /media/providers/{id}`
    - `GET /media/combos`, `POST /media/combos`, `PATCH /media/combos/{id}`, `DELETE /media/combos/{id}`
  - Gateway routes:
    - `POST /v1/audio/transcriptions`, `POST /v1/audio/translations`, `POST /v1/audio/speech`
    - `POST /v1/images/generations`, `POST /v1/images/edits`, `POST /v1/images/variations`
    - `POST /v1/videos/generations`, `GET /v1/videos/{id}`
    - `POST /v1/search`, `POST /v1/web/fetch`

- [ ] **Step 1: Write failing tests in `internal/modules/media/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `media` migrations, repo, SSRF-safe web fetcher, module handlers, and `contracts/openapi/modules/media.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager)

**Files:**
- Create: `internal/modules/tunnel/migrations/0001_init.sql`
- Create: `internal/modules/tunnel/repo.go`
- Create: `internal/modules/tunnel/runner.go`
- Create: `internal/modules/tunnel/module.go`
- Create: `internal/modules/tunnel/module_test.go`
- Create: `contracts/openapi/modules/tunnel.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `config.Config`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `tunnel.New(opts Options) *Module` (`Name() = "tunnel"`, `Dependencies() = []string{"system"}`)
  - Implements `lifecycle.Module`
  - SQLite tables (`0001_init.sql`):
    - `tunnel_config(id INTEGER PRIMARY KEY CHECK (id = 1), provider TEXT NOT NULL DEFAULT 'cloudflare', local_url TEXT NOT NULL DEFAULT '', tailscale_mode TEXT NOT NULL DEFAULT 'funnel', auto_start INTEGER NOT NULL DEFAULT 0, updated_at TEXT NOT NULL)`
    - `tunnel_events(id TEXT PRIMARY KEY, provider TEXT NOT NULL, event_type TEXT NOT NULL, message TEXT NOT NULL, public_url TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
  - AdminAPI routes (`/api`):
    - `GET /tunnel/status`, `PATCH /tunnel/config`, `POST /tunnel/start`, `POST /tunnel/stop`, `GET /tunnel/events`

- [ ] **Step 1: Write failing tests in `internal/modules/tunnel/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `tunnel` migrations, repo, process runner, module handlers, and `contracts/openapi/modules/tunnel.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 3: Wire Wave 3B in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go` (`TestWave3BMediaAndTunnelE2E`, `TestWave3BDisabledModulesDetachCleanly`)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Wire `media` and `tunnel` in `main.go`, update `openapi.yaml`, run `make codegen`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/media/{schemas.ts,schemas.test.ts,hooks.ts,MediaConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/media/page.tsx`
- Create: `web/src/features/tunnel/{schemas.ts,schemas.test.ts,hooks.ts,TunnelConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/tunnel/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

- [ ] **Step 1: Write failing Zod schema tests**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement schemas, hooks, `MediaConsole`, `TunnelConsole`, pages, `AppShell`, `SystemOverview`**
- [ ] **Step 4: Run `make check`, `make build`, and `graphify update .`**
- [ ] **Step 5: Commit**
