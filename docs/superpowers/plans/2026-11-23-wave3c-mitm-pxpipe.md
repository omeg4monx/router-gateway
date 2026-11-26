# Wave 3C — Local MITM Proxy & PXPIPE Transform Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 3C by adding two isolated vertical-slice modules — `mitm` (Go-native loopback TLS MITM proxy, X.509 Root CA & on-the-fly per-host leaf certificate engine, IDE tool host registry for Antigravity/Copilot/Kiro/Cursor, model alias rewriting, OS trust/hosts setup guide, and intercepted traffic audit log) and `pxpipe` (Go-native in-process request transform pipeline implementing `contracts.GatewayInterceptor` with runtime start/stop/restart controls, configurable transform stages, 3-step self-test health probe, live metrics, dry-run preview, and ring-buffered event logs) — complete with OpenAPI 3.1 contracts, Next.js 16 static-export studios, E2E tests, and single-binary embedding.

**Architecture:** `mitm` owns its own SQLite tables (`mitm_config`, `mitm_aliases`, `mitm_events`), stores its ECDSA P-256 Root CA under `<DataDir>/mitm/`, enforces loopback-only listen and gateway target addresses, rewrites matched IDE requests' JSON `model` fields according to active aliases, and shuts down its TLS listener cleanly on `Stop()`. `pxpipe` owns its own SQLite table (`pxpipe_config`), implements `contracts.GatewayInterceptor` at priority `35` (after `tokensaver` at `30`), applies deterministic JSON/message transforms in-process when active, exposes live self-test health and dry-run preview endpoints, and integrates cleanly into `cmd/omeg4mon-router/main.go` with zero cross-module Go imports.

**Tech Stack:** Go 1.26+ (`net/http`, `crypto/ecdsa`, `crypto/x509`, `crypto/tls`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded request/response body limits and bounded in-memory caches/rings.
- **Security Posture:** `mitm` strictly restricts both `listen_addr` and `target_gateway_url` to loopback hosts (`127.0.0.1`, `localhost`, `::1`), stores Root CA private keys with `0600` permissions and never exposes private keys via HTTP APIs; `pxpipe` redacts auth headers and never leaks raw prompts in logs.
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task passes `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and final verification runs `make check`, `make build`, and `graphify update .`.

## Review Focus

1. **MITM loopback enforcement & private key isolation:** Configuring `listen_addr` or `target_gateway_url` with a non-loopback address (`0.0.0.0:8443`, `http://10.0.0.5:20128`) fails with `400`; `GET /api/mitm/ca/download` returns only the public Root CA PEM (`BEGIN CERTIFICATE`) and never the private key (`0600` on disk). _(Pinned in Task 1 & Task 3)_
2. **Live TLS MITM interception & alias rewriting:** Starting `mitm` on `127.0.0.1:0` with a test gateway target intercepts a TLS request verified against the generated Root CA, rewrites the JSON `model` via an active alias for the matched IDE tool host, forwards to the target gateway, and records an audit event in `mitm_events`; `Stop()` cleanly closes the listener with zero goroutine leaks. _(Pinned in Task 1 & Task 3)_
3. **PXPIPE interceptor runtime gating & transform fidelity:** Stopping `pxpipe` via `POST /api/pxpipe/stop` causes `OnRequest` to pass payloads through untouched; starting it applies configured transforms (null field stripping, whitespace collapse, system prompt deduplication, history windowing) without corrupting valid JSON or non-JSON bodies. _(Pinned in Task 2 & Task 3)_
4. **Disabled optional modules cleanly detach:** Disabling `mitm` or `pxpipe` via `OMEG4MON_DISABLED_MODULES` omits their routes and excludes `pxpipe` from the gateway interceptor chain without affecting core gateway traffic. _(Pinned in Task 3)_
5. **Static-export UI compliance & error surfacing:** `/dashboard/mitm/` and `/dashboard/pxpipe/` export cleanly to static HTML (`web/out/`), surface all query/mutation errors via `Alert`, and contain zero `backdrop-blur` or UI emojis. _(Pinned in Task 4)_

---

### Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases)

**Files:**
- Create: `internal/modules/mitm/migrations/0001_init.sql`
- Create: `internal/modules/mitm/certs.go`
- Create: `internal/modules/mitm/repo.go`
- Create: `internal/modules/mitm/proxy.go`
- Create: `internal/modules/mitm/module.go`
- Create: `internal/modules/mitm/module_test.go`
- Create: `contracts/openapi/modules/mitm.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `config.Config`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `mitm.New(opts Options) *Module` (`Name() = "mitm"`, `Dependencies() = []string{"system"}`)
  - Implements `lifecycle.Module`
  - SQLite tables (`0001_init.sql`):
    - `mitm_config(id INTEGER PRIMARY KEY CHECK (id = 1), listen_addr TEXT NOT NULL DEFAULT '127.0.0.1:8443', target_gateway_url TEXT NOT NULL DEFAULT '', enabled_tools_json TEXT NOT NULL DEFAULT '["antigravity","copilot","kiro","cursor"]', auto_start INTEGER NOT NULL DEFAULT 0, updated_at TEXT NOT NULL)`
    - `mitm_aliases(id TEXT PRIMARY KEY, tool TEXT NOT NULL, source_model TEXT NOT NULL, target_model TEXT NOT NULL, enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, UNIQUE(tool, source_model))`
    - `mitm_events(id TEXT PRIMARY KEY, tool TEXT NOT NULL, host TEXT NOT NULL, method TEXT NOT NULL, path TEXT NOT NULL, source_model TEXT NOT NULL DEFAULT '', target_model TEXT NOT NULL DEFAULT '', status_code INTEGER NOT NULL DEFAULT 0, latency_ms INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL)`
  - AdminAPI routes (`/api`):
    - `GET /mitm/status`, `PATCH /mitm/config`, `POST /mitm/start`, `POST /mitm/stop`
    - `POST /mitm/ca/generate`, `GET /mitm/ca/download`
    - `GET /mitm/aliases`, `POST /mitm/aliases`, `DELETE /mitm/aliases/{id}`
    - `GET /mitm/events`, `DELETE /mitm/events`

- [ ] **Step 1: Write failing tests in `internal/modules/mitm/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `mitm` migrations, cert manager, repo, loopback TLS proxy, module handlers, and `contracts/openapi/modules/mitm.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs)

**Files:**
- Create: `internal/modules/pxpipe/migrations/0001_init.sql`
- Create: `internal/modules/pxpipe/repo.go`
- Create: `internal/modules/pxpipe/transform.go`
- Create: `internal/modules/pxpipe/module.go`
- Create: `internal/modules/pxpipe/module_test.go`
- Create: `contracts/openapi/modules/pxpipe.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `pxpipe.New(opts Options) *Module` (`Name() = "pxpipe"`, `Dependencies() = []string{"system"}`, `Priority() = 35`)
  - Implements `lifecycle.Module` and `contracts.GatewayInterceptor`
  - SQLite table (`0001_init.sql`):
    - `pxpipe_config(id INTEGER PRIMARY KEY CHECK (id = 1), enabled INTEGER NOT NULL DEFAULT 1, auto_start INTEGER NOT NULL DEFAULT 1, strip_null_fields INTEGER NOT NULL DEFAULT 1, collapse_whitespace INTEGER NOT NULL DEFAULT 1, dedupe_system_prompts INTEGER NOT NULL DEFAULT 1, max_history_messages INTEGER NOT NULL DEFAULT 0, updated_at TEXT NOT NULL)`
  - AdminAPI routes (`/api`):
    - `GET /pxpipe/status`, `PATCH /pxpipe/config`, `POST /pxpipe/start`, `POST /pxpipe/stop`, `POST /pxpipe/restart`
    - `GET /pxpipe/health`, `GET /pxpipe/stats`, `POST /pxpipe/stats/reset`
    - `GET /pxpipe/logs`, `DELETE /pxpipe/logs`, `POST /pxpipe/preview`

- [ ] **Step 1: Write failing tests in `internal/modules/pxpipe/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `pxpipe` migrations, repo, transform engine, module handlers, and `contracts/openapi/modules/pxpipe.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 3: Wire Wave 3C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go` (`TestWave3CMitmAndPxpipeE2E`, `TestWave3CDisabledModulesDetachCleanly`)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Wire `mitm` and `pxpipe` in `main.go`, update `openapi.yaml`, run `make codegen`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/mitm/{schemas.ts,schemas.test.ts,hooks.ts,MitmConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/mitm/page.tsx`
- Create: `web/src/features/pxpipe/{schemas.ts,schemas.test.ts,hooks.ts,PxpipeConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/pxpipe/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

- [ ] **Step 1: Write failing Zod schema tests**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement schemas, hooks, `MitmConsole`, `PxpipeConsole`, pages, `AppShell`, `SystemOverview`**
- [ ] **Step 4: Run `make check`, `make build`, and `graphify update .`**
- [ ] **Step 5: Commit**
