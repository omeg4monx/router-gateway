# Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 3A by adding two isolated vertical-slice modules under the single-admin dashboard model — `teams` (API-key grouping, shared USD budgets, shared RPM limits, spend tracking via `usage.recorded`, and `GatewayInterceptor` enforcement at `Priority=5`) and `vouchers` (voucher creation, atomic claim flow with API-key provisioning, optional team attachment, and per-code/per-actor/per-IP anti-abuse enforcement) — complete with OpenAPI 3.1 contracts, Next.js 16 static-export studios, E2E tests, and single-binary embedding.

**Architecture:** `internal/contracts` defines narrow `KeyProvisioner` (implemented by `keys.Module`) and `TeamAssigner` (implemented by `teams.Module`) so `vouchers` can mint gateway API keys and bind them to a team during claim with zero cross-module Go imports. `teams` implements `contracts.GatewayInterceptor` (`Priority=5`) to enforce shared team USD budgets and shared sliding-window RPM limits before upstream dialing, and subscribes to `usage.recorded` on `eventbus.Bus` to accumulate team spend asynchronously. `vouchers` mounts admin CRUD on `AdminAPI` and self-service `POST /api/vouchers/claim` on `PublicAPI` (also usable from the admin studio) protected by an atomic SQLite transaction and an in-memory per-IP sliding-window anti-abuse guard.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Single-admin dashboard model (Owner Decision 2A):** `teams` operates strictly as API-key, budget, and quota grouping — no multi-user dashboard RBAC or per-tenant auth.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded in-memory rate-limit windows.
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task passes `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and final verification runs `make check`, `make build`, and `graphify update .`.

## Review Focus

1. **Team budget & RPM enforcement fails closed with 429 before dialing upstream:** When a key belongs to a disabled team, a team whose `spent_usd >= budget_usd` (with `budget_usd > 0`), or a team exceeding `rate_limit_rpm`, `teams.OnRequest` rejects with HTTP `429` (or `403` when disabled) without dialing upstream. _(Pinned in Tasks 1 & 3)_
2. **Voucher claim atomicity prevents double-redemption races:** Concurrent `POST /api/vouchers/claim` calls for a single-use voucher (`max_redemptions = 1`) must allow exactly one claim (`201`/`200`) and reject all concurrent callers with `409 Conflict` (`voucher_exhausted`), never minting duplicate API keys. _(Pinned in Task 2)_
3. **Voucher anti-abuse enforces per-actor uniqueness & per-IP rate limiting:** The same `claimer_id` (or normalized IP when omitted) cannot claim the same voucher code twice (`409`), and rapid claim attempts from the same IP exceed the anti-abuse window and receive `429 Too Many Requests`. _(Pinned in Task 2)_
4. **Disabled optional modules cleanly detach:** Disabling `teams` or `vouchers` via `OMEG4MON_DISABLED_MODULES` omits their routes, interceptor, and eventbus listeners without affecting core gateway traffic. _(Pinned in Task 3)_
5. **Static-export UI compliance & error surfacing:** `/dashboard/teams/` and `/dashboard/vouchers/` export cleanly to static HTML (`web/out/`), surface all query/mutation errors via `Alert`, and contain zero `backdrop-blur` or UI emojis. _(Pinned in Task 4)_

---

### Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module

**Files:**
- Modify: `internal/contracts/contracts.go`
- Modify: `internal/modules/keys/module.go`
- Modify: `internal/modules/keys/module_test.go`
- Create: `internal/modules/teams/migrations/0001_init.sql`
- Create: `internal/modules/teams/repo.go`
- Create: `internal/modules/teams/module.go`
- Create: `internal/modules/teams/module_test.go`
- Create: `contracts/openapi/modules/teams.yaml`

**Interfaces:**
- Produces:
  - `contracts.ProvisionKeyRequest`, `contracts.ProvisionedKey`, `contracts.KeyProvisioner`, `contracts.TeamAssigner`
  - `(*keys.Module).ProvisionKey(ctx, req)` implementing `contracts.KeyProvisioner`
  - `teams.New(opts Options) *Module` (`Name() = "teams"`, `Dependencies() = []string{"system", "keys"}`, `Priority() = 5`)
  - Implements `lifecycle.Module`, `contracts.GatewayInterceptor`, and `contracts.TeamAssigner`
  - SQLite tables (`0001_init.sql`):
    - `teams(id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, description TEXT NOT NULL DEFAULT '', budget_usd REAL NOT NULL DEFAULT 0, spent_usd REAL NOT NULL DEFAULT 0, rate_limit_rpm INTEGER NOT NULL DEFAULT 0, enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `team_keys(team_id TEXT NOT NULL, key_id TEXT NOT NULL PRIMARY KEY, added_at TEXT NOT NULL, FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE)`
  - AdminAPI routes (`/api`):
    - `GET /teams`, `POST /teams`, `GET /teams/{id}`, `PATCH /teams/{id}`, `DELETE /teams/{id}`
    - `POST /teams/{id}/keys`, `DELETE /teams/{id}/keys/{keyId}`
    - `POST /teams/{id}/reset-spend`

- [ ] **Step 1: Write failing tests in `internal/modules/keys/module_test.go` and `internal/modules/teams/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement contracts, `keys.ProvisionKey`, `teams` module, and `contracts/openapi/modules/teams.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse)

**Files:**
- Create: `internal/modules/vouchers/migrations/0001_init.sql`
- Create: `internal/modules/vouchers/repo.go`
- Create: `internal/modules/vouchers/module.go`
- Create: `internal/modules/vouchers/module_test.go`
- Create: `contracts/openapi/modules/vouchers.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.KeyProvisioner`, `contracts.TeamAssigner`
- Produces:
  - `vouchers.New(opts Options) *Module` (`Name() = "vouchers"`, `Dependencies() = []string{"system", "keys"}`)
  - SQLite tables (`0001_init.sql`):
    - `vouchers(id TEXT PRIMARY KEY, code TEXT NOT NULL UNIQUE, name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', quota_usd REAL NOT NULL DEFAULT 0, rate_limit_rpm INTEGER NOT NULL DEFAULT 0, allowed_models_json TEXT NOT NULL DEFAULT '[]', team_id TEXT NOT NULL DEFAULT '', max_redemptions INTEGER NOT NULL DEFAULT 1, redeemed_count INTEGER NOT NULL DEFAULT 0, enabled INTEGER NOT NULL DEFAULT 1, expires_at TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `voucher_redemptions(id TEXT PRIMARY KEY, voucher_id TEXT NOT NULL, code TEXT NOT NULL, claimer_id TEXT NOT NULL, key_id TEXT NOT NULL, ip_address TEXT NOT NULL DEFAULT '', redeemed_at TEXT NOT NULL, FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE, UNIQUE(voucher_id, claimer_id))`
  - AdminAPI routes (`/api`):
    - `GET /vouchers`, `POST /vouchers`, `PATCH /vouchers/{id}`, `DELETE /vouchers/{id}`
    - `GET /vouchers/{id}/redemptions`
  - PublicAPI routes (`/api`):
    - `POST /vouchers/claim` (atomic single-tx claim + anti-abuse per-IP sliding window + per-actor uniqueness + API key minting + optional team binding)

- [ ] **Step 1: Write failing tests in `internal/modules/vouchers/module_test.go` (including concurrent race test pinning Review Focus #2 & #3)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `vouchers` module and `contracts/openapi/modules/vouchers.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go` (`TestWave3ATeamsBudgetAndVoucherClaimE2E`, `TestWave3ADisabledModulesDetachCleanly`)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Wire `teams` (with `GatewayInterceptor`) and `vouchers` in `main.go`, update `openapi.yaml`, run `make codegen`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/teams/{schemas.ts,schemas.test.ts,hooks.ts,TeamsConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/teams/page.tsx`
- Create: `web/src/features/vouchers/{schemas.ts,schemas.test.ts,hooks.ts,VouchersConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/vouchers/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

- [ ] **Step 1: Write failing Zod schema tests**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement schemas, hooks, `TeamsConsole`, `VouchersConsole`, pages, `AppShell`, `SystemOverview`**
- [ ] **Step 4: Run `make check`, `make build`, and `graphify update .`**
- [ ] **Step 5: Commit**
