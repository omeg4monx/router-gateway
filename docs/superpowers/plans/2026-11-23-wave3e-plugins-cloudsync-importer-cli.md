# Wave 3E — Plugins & Skills, Cloud Sync, Legacy Importer & CLI Launcher Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete Wave 3 of the modular Next.js (TypeScript) × Go rewrite by delivering three isolated backend modules (`plugins` with request-enhancement plugins + `GatewayInterceptor` + Skills catalog, `cloudsync` with vaulted SSRF-safe remote state sync + portable JSON snapshots, and `importer` with read-only legacy `data.sqlite` inspection + one-way import via narrow contracts), a thin zero-dependency multi-platform npm binary launcher (`cli/`), OpenAPI 3.1 contracts, Next.js 16 static-export studios, E2E tests, and single-binary verification.

**Architecture:** `internal/contracts` gains narrow `SettingsSnapshotter` and `ComboSnapshotter` interfaces implemented by `settings.Module` and `combos.Module` so `cloudsync` and `importer` never import sibling modules or touch foreign SQLite tables. `plugins` owns `plugin_configs` and `plugin_skills`, implements `contracts.GatewayInterceptor` (priority 60) to apply enabled model-scoped request enhancements (`imageVision`, `thinkDeeper`, `speedMode`, `jsonGuard`, `contextSqueezer`, `openaiToolBridge`, `antiSlop`) with a dry-run preview endpoint and CRUD + `SKILL.md` export for skills. `cloudsync` owns `cloudsync_config` and `cloudsync_history`, stores bearer tokens encrypted in `*crypto.Vault`, enforces SSRF validation on remote endpoints (rejecting loopback/private/link-local/metadata IPs unless overridden in tests), and supports `remote_wins` / `local_wins` / `merge` conflict policies plus portable JSON export/restore. `importer` opens legacy `data.sqlite` files strictly read-only (`mode=ro&_query_only=1`) after validating regular-file + SQLite3 magic header, inspects legacy tables gracefully, and imports settings/combos/API keys through contracts while recording every run in `importer_runs`. `cli/` provides a pure Node/Bun launcher resolving the platform single binary (`omeg4mon-router-<os>-<arch>[.exe]`).

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `*crypto.Vault`, `*eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, TanStack Query, Node/Bun (`cli/`), `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded request/response body limits, bounded timeouts.
- **Security Posture:** `cloudsync` encrypts remote auth tokens via `*crypto.Vault`, never returns raw tokens in HTTP responses, validates remote endpoint URLs against SSRF (private/loopback/link-local/metadata IPs blocked by default), and bounds remote response sizes; `importer` validates regular file + SQLite header magic and opens source DBs strictly read-only (`mode=ro&_query_only=1`) without mutating source files.
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task passes `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and final verification runs `make check`, `make build`, `cd cli && bun test`, and `graphify update .`.

## Review Focus

1. **Plugins interceptor model scoping & invalid-JSON resilience:** Enabling a plugin with a `models` allowlist mutates only matching models; non-matching models pass through untouched; malformed JSON bodies fail gracefully without panicking. _(Pinned in Task 1)_
2. **CloudSync SSRF guard & vaulted token isolation:** Configuring a loopback/private/metadata endpoint (e.g. `http://127.0.0.1:8080`, `http://169.254.169.254`) is rejected when private endpoints are disallowed; `GET/PATCH /api/cloudsync/config` never returns plaintext `auth_token`. _(Pinned in Task 2 & Task 4)_
3. **CloudSync conflict policy & snapshot round-trip:** `local_wins` preserves existing local keys on restore/pull, `remote_wins` overwrites with remote values, and `merge` unions cleanly while recording checksum + history. _(Pinned in Task 2)_
4. **Importer read-only safety & invalid-file rejection:** Passing a non-existent path, directory, or non-SQLite file returns `400`; inspecting/importing a valid legacy `data.sqlite` leaves the source file's modification time and contents untouched, and `dry_run: true` performs zero mutations. _(Pinned in Task 3 & Task 4)_
5. **Disabled optional modules, CLI resolver & static-export UI compliance:** Disabling `plugins,cloudsync,importer` via `OMEG4MON_DISABLED_MODULES` cleanly detaches all routes and interceptors; `cli/lib/resolve.js` maps Linux/macOS/Windows `x64`/`arm64` targets deterministically; `/dashboard/plugins/`, `/dashboard/skills/`, `/dashboard/cloud-sync/`, and `/dashboard/importer/` export cleanly to `web/out/` with zero `backdrop-blur` or UI emojis. _(Pinned in Task 4, Task 5 & Task 6)_

---

### Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module

**Files:**
- Modify: `internal/contracts/contracts.go`
- Modify: `internal/modules/settings/module.go`
- Modify: `internal/modules/combos/module.go`
- Create: `internal/modules/plugins/migrations/0001_init.sql`
- Create: `internal/modules/plugins/repo.go`
- Create: `internal/modules/plugins/interceptor.go`
- Create: `internal/modules/plugins/module.go`
- Create: `internal/modules/plugins/module_test.go`
- Create: `contracts/openapi/modules/plugins.yaml`

**Interfaces:**
- Produces in `internal/contracts`:
  - `SettingsSnapshotter`: `ExportSettings(ctx context.Context) (map[string]any, error)`, `ImportSettings(ctx context.Context, values map[string]any, actor, reason string) (int, error)`
  - `ComboSnapshot`: `{Name, Description, Strategy string, Targets []string, Enabled bool}`
  - `AliasSnapshot`: `{Alias, Target string}`
  - `ComboSnapshotter`: `ExportCombos(ctx context.Context) ([]ComboSnapshot, []AliasSnapshot, error)`, `ImportCombos(ctx context.Context, combos []ComboSnapshot, aliases []AliasSnapshot, overwrite bool, actor string) (int, int, error)`
- Produces in `internal/modules/plugins`:
  - `plugins.New(opts Options) *Module` (`Name() = "plugins"`, `Dependencies() = []string{"system"}`)
  - Implements `lifecycle.Module`, `contracts.GatewayInterceptor` (`Priority() = 60`)
  - AdminAPI routes (`/api`):
    - `GET /plugins`, `PUT /plugins/{key}`, `POST /plugins/preview`
    - `GET /skills`, `POST /skills`, `GET /skills/{id}`, `PATCH /skills/{id}`, `DELETE /skills/{id}`, `GET /skills/{id}/export`

- [ ] **Step 1:** Add `SettingsSnapshotter` and `ComboSnapshotter` to `internal/contracts/contracts.go` and implement them on `settings.Module` and `combos.Module`.
- [ ] **Step 2:** Write failing tests in `internal/modules/plugins/module_test.go` covering the 7 built-in plugins, `PUT /api/plugins/{key}`, `OnRequest` model scoping + `antiSlop`/`jsonGuard`/`thinkDeeper`/`speedMode`/`contextSqueezer`, `POST /api/plugins/preview`, and Skills CRUD + `SKILL.md` export.
- [ ] **Step 3:** Run `go test -count=1 ./internal/modules/plugins/...` to verify failure.
- [ ] **Step 4:** Implement `plugins` migrations, repo, interceptor, module handlers, and `contracts/openapi/modules/plugins.yaml`.
- [ ] **Step 5:** Run `golangci-lint` and `go test -count=1 -race ./internal/modules/settings/... ./internal/modules/combos/... ./internal/modules/plugins/...`, then commit.

---

### Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI)

**Files:**
- Create: `internal/modules/cloudsync/migrations/0001_init.sql`
- Create: `internal/modules/cloudsync/repo.go`
- Create: `internal/modules/cloudsync/ssrf.go`
- Create: `internal/modules/cloudsync/module.go`
- Create: `internal/modules/cloudsync/module_test.go`
- Create: `contracts/openapi/modules/cloudsync.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*crypto.Vault`, `*http.Client`, `*eventbus.Bus`, `*slog.Logger`, `contracts.SettingsSnapshotter`, `contracts.ComboSnapshotter`, `AllowPrivateEndpoints bool` (test-only override)
- Produces:
  - `cloudsync.New(opts Options) *Module` (`Name() = "cloudsync"`, `Dependencies() = []string{"system", "settings", "combos"}`)
  - AdminAPI routes (`/api`):
    - `GET /cloudsync/config`, `PATCH /cloudsync/config`
    - `POST /cloudsync/push`, `POST /cloudsync/pull`
    - `GET /cloudsync/snapshot`, `POST /cloudsync/restore`
    - `GET /cloudsync/history`, `DELETE /cloudsync/history`

- [ ] **Step 1:** Write failing tests in `internal/modules/cloudsync/module_test.go` covering vaulted token isolation, SSRF rejection of private/loopback endpoints, push/pull against `httptest.Server`, conflict policies (`remote_wins`/`local_wins`/`merge`), portable snapshot export/restore, and history CRUD.
- [ ] **Step 2:** Run `go test -count=1 ./internal/modules/cloudsync/...` to verify failure.
- [ ] **Step 3:** Implement `cloudsync` migrations, repo, SSRF validator, module handlers, and `contracts/openapi/modules/cloudsync.yaml`.
- [ ] **Step 4:** Run `golangci-lint` and `go test -count=1 -race ./internal/modules/cloudsync/...`, then commit.

---

### Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI)

**Files:**
- Create: `internal/modules/importer/migrations/0001_init.sql`
- Create: `internal/modules/importer/repo.go`
- Create: `internal/modules/importer/legacy.go`
- Create: `internal/modules/importer/module.go`
- Create: `internal/modules/importer/module_test.go`
- Create: `contracts/openapi/modules/importer.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.SettingsSnapshotter`, `contracts.ComboSnapshotter`, `contracts.KeyProvisioner`
- Produces:
  - `importer.New(opts Options) *Module` (`Name() = "importer"`, `Dependencies() = []string{"system", "settings", "combos", "keys"}`)
  - AdminAPI routes (`/api`):
    - `POST /importer/inspect`
    - `POST /importer/run`
    - `GET /importer/history`, `DELETE /importer/history`

- [ ] **Step 1:** Write failing tests in `internal/modules/importer/module_test.go` creating a synthetic legacy `data.sqlite`, testing invalid path/non-SQLite rejection, read-only inspection (verifying source file untouched), `dry_run: true`, live import of settings/combos/keys, and run history.
- [ ] **Step 2:** Run `go test -count=1 ./internal/modules/importer/...` to verify failure.
- [ ] **Step 3:** Implement `importer` migrations, repo, read-only legacy reader, module handlers, and `contracts/openapi/modules/importer.yaml`.
- [ ] **Step 4:** Run `golangci-lint` and `go test -count=1 -race ./internal/modules/importer/...`, then commit.

---

### Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1:** Wire `plugins`, `cloudsync`, and `importer` in `cmd/omeg4mon-router/main.go` (including `plugins` in `GatewayInterceptor` chain when enabled).
- [ ] **Step 2:** Register `plugins.yaml`, `cloudsync.yaml`, and `importer.yaml` in `contracts/openapi/openapi.yaml` and run `make codegen`.
- [ ] **Step 3:** Add E2E tests in `cmd/omeg4mon-router/main_test.go` exercising `plugins`, `cloudsync`, `importer`, and `OMEG4MON_DISABLED_MODULES=plugins,cloudsync,importer`.
- [ ] **Step 4:** Run `golangci-lint` and `go test -count=1 -race ./cmd/omeg4mon-router/...`, then commit.

---

### Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`)

**Files:**
- Create: `cli/package.json`
- Create: `cli/lib/resolve.js`
- Create: `cli/bin/omeg4mon-router.js`
- Create: `cli/test/resolve.test.mjs`
- Create: `cli/README.md`

- [ ] **Step 1:** Write failing unit tests in `cli/test/resolve.test.mjs` for platform/arch binary naming (`linux`/`darwin`/`win32` × `x64`/`arm64`), unsupported platform errors, `OMEG4MON_ROUTER_BIN` override, and candidate search order.
- [ ] **Step 2:** Run `cd cli && bun test` to verify failure.
- [ ] **Step 3:** Implement `cli/package.json`, `cli/lib/resolve.js`, `cli/bin/omeg4mon-router.js`, and `cli/README.md`.
- [ ] **Step 4:** Run `cd cli && bun test`, then commit.

---

### Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification

**Files:**
- Create: `web/src/features/plugins/hooks.ts`
- Create: `web/src/features/plugins/PluginsStudio.tsx`
- Create: `web/src/features/plugins/SkillsStudio.tsx`
- Create: `web/src/app/(dashboard)/dashboard/plugins/page.tsx`
- Create: `web/src/app/(dashboard)/dashboard/skills/page.tsx`
- Create: `web/src/features/cloudsync/hooks.ts`
- Create: `web/src/features/cloudsync/CloudSyncStudio.tsx`
- Create: `web/src/app/(dashboard)/dashboard/cloud-sync/page.tsx`
- Create: `web/src/features/importer/hooks.ts`
- Create: `web/src/features/importer/ImporterStudio.tsx`
- Create: `web/src/app/(dashboard)/dashboard/importer/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`

- [ ] **Step 1:** Implement `plugins`, `cloudsync`, and `importer` TanStack Query hooks and `DESIGN.md`-compliant studios + static-export App Router pages.
- [ ] **Step 2:** Add `Plugins`, `Skills`, `Cloud Sync`, and `Legacy Importer` entries to `ACTIVE_NAV_ITEMS` in `web/src/shared/ui/AppShell.tsx`.
- [ ] **Step 3:** Run full verification (`make check`, `make build`, `cd cli && bun test`, single-binary smoke check, and `graphify update .`), then commit.
