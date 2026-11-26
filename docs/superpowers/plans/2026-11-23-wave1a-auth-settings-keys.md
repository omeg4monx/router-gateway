# Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the first vertical slice of Wave 1 — single-admin cookie session authentication (`auth`), typed settings with revision history & rollback (`settings`), and gateway API key management (`keys`) — wired cleanly into the Go modular monolith, OpenAPI 3.1 contract pipeline, and Next.js 16 static-export dashboard UI.

**Architecture:** Three isolated vertical-slice Go modules (`internal/modules/auth`, `internal/modules/settings`, `internal/modules/keys`) implement `lifecycle.Module`, own their own SQLite tables/migrations, publish invalidation events over `eventbus`, and communicate with `system` and `httpserver` exclusively via narrow interfaces in `internal/contracts/`. Modular OpenAPI 3.1 fragments (`contracts/openapi/modules/{auth,settings,keys}.yaml`) generate strict TypeScript types consumed by matching feature slices in `web/src/features/{auth,settings,keys}/`.

**Tech Stack:** Go 1.24+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `crypto/sha256`, `crypto/rand`, HMAC-SHA256 `crypto.SessionSigner`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, Turbopack, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), `zustand`, TanStack Query, Bun `1.4+`.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying of source code, SQL, or assets from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER read or write another module's SQLite tables directly. Cross-module calls use narrow interfaces in `internal/contracts/` wired in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL:** Pure-Go `modernc.org/sqlite` using `db.Writer()` (max 1 conn) for writes and `db.Reader()` (max 8 conns) for reads; per-module embedded SQL migrations tracked in `schema_migrations`.
- **Single-admin auth & remote-bind posture:** Fresh installs default to `require_login=true`; session cookie is `HttpOnly`, `SameSite=Lax`, signed via `crypto.SessionSigner` (HMAC-SHA256 backed by vault master key). Non-loopback binds enforce `require_login=true` unless `OMEG4MON_ALLOW_INSECURE_REMOTE=true`.
- **Static-export frontend (`web/`):** Next.js `16.4.6` with `output: "export"` and `trailingSlash: true`; active usage of `react-hook-form`, `@hookform/resolvers` (`zod`), `shadcn/ui`, and `zustand`; strict adherence to `docs/DESIGN.md` (`#E56A4A` ember orange, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces without `backdrop-blur`, zero UI emojis, WCAG AA contrast).
- **Formatter & Linter Discipline:** Every task must pass `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, and `cd web && bun run check` (Biome v2) — verified via `make check` before commit.
- **Graphify & Branch Discipline:** Run `graphify update .` after code changes before committing; commit all Wave 1A work on `wave1a/auth-settings-keys` and open a PR targeting `migration/ts-go` (never `main`).

## Review Focus

1. **Non-loopback login bypass prevention:** Attempting `PUT /api/auth/policy` with `require_login: false` while bound to a non-loopback host (`0.0.0.0`) without `AllowInsecureRemote` must fail with `403 Forbidden`, and `auth.Module.RequireLogin(ctx)` must always return `true` on non-loopback binds regardless of stored DB value. _(Pinned in Task 2)_
2. **Auth login & cookie edge cases:** `POST /api/auth/login` against an uninitialized DB must return `409 Conflict` (`not_initialized`); invalid credentials return `401 Unauthorized` with constant-time/bcrypt verification; tampered or expired cookies are rejected with `401`; `Secure` cookie attribute is set automatically when `r.TLS != nil` or `X-Forwarded-Proto: https`. _(Pinned in Task 2)_
3. **Settings validation & rollback atomicity:** Updating an unknown setting key returns `400 invalid_setting_key`, invalid typed value returns `400 invalid_setting_value`, and rolling back to a nonexistent revision ID returns `404 revision_not_found` without mutating existing settings or appending corrupt history. _(Pinned in Task 3)_
4. **API key secret hygiene & cache invalidation:** Raw API key tokens (`sk-om4-...`) are returned ONLY in `POST /api/keys` (and `POST /api/keys/{id}/rotate`), stored exclusively as SHA-256 hex digest + display prefix (never plaintext in SQLite or `GET /api/keys`), and disabling/rotating/deleting a key immediately invalidates the in-memory `KeyValidator` cache so subsequent `ValidateKey` calls fail immediately. _(Pinned in Task 4)_
5. **Static-export routing & auth redirect loop safety:** Unauthenticated visits to `/dashboard/`, `/dashboard/keys/`, or `/dashboard/settings/` when `require_login=true` redirect cleanly to `/login/` without an infinite redirect loop across `/`, `/init/`, `/login/`, and `/dashboard/`, and direct refresh works on all static-exported trailing-slash routes. _(Pinned in Tasks 6 & 7)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Contracts | `internal/contracts/contracts.go`, `internal/contracts/contracts_test.go` | Add `AdminCredentialStore` & `ErrNotInitialized` for `system` <-> `auth` boundary |
| System adapter | `internal/modules/system/repo.go`, `internal/modules/system/module.go`, `internal/modules/system/module_test.go` | Implement `contracts.AdminCredentialStore` cleanly on `system.Module` |
| Auth module | `internal/modules/auth/{module.go,handlers.go,cookie.go,module_test.go}`, `contracts/openapi/modules/auth.yaml` | Admin login/logout/me/password/policy endpoints + `contracts.AuthVerifier` |
| Settings module | `internal/modules/settings/{migrations/0001_init.sql,schema.go,repo.go,module.go,module_test.go}`, `contracts/openapi/modules/settings.yaml` | Typed KV settings store, validation, revision history, atomic rollback, `eventbus` invalidation |
| Keys module | `internal/modules/keys/{migrations/0001_init.sql,token.go,repo.go,cache.go,module.go,module_test.go}`, `contracts/openapi/modules/keys.yaml` | API keys CRUD + rotate, SHA-256 lookup + cache, `contracts.KeyValidator`, `allowed_models` matching |
| Wiring & OpenAPI | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Register modules in dependency order, wire `AuthVerifier` into `httpserver`, run `make codegen` |
| Web Auth UI | `web/src/features/auth/{schemas.ts,hooks.ts,LoginForm.tsx,AuthGuard.tsx,schemas.test.ts}`, `web/src/app/login/page.tsx`, `web/src/app/page.tsx`, `web/src/shared/ui/AppShell.tsx` | Live cookie login form (`react-hook-form` + `zod`), dashboard auth guard, sign-out control |
| Web Keys & Settings UI | `web/src/features/keys/*`, `web/src/features/settings/*`, `web/src/app/(dashboard)/dashboard/keys/page.tsx`, `web/src/app/(dashboard)/dashboard/settings/page.tsx`, `web/src/shared/ui/{dialog.tsx,switch.tsx,select.tsx}` | API Keys management studio (create, reveal-once copy, edit scopes/models/limits, disable/rotate/delete) + Typed Settings & History Rollback studio |

---

### Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter

**Files:**
- Modify: `internal/contracts/contracts.go`
- Modify: `internal/contracts/contracts_test.go`
- Modify: `internal/modules/system/repo.go`
- Modify: `internal/modules/system/module.go`
- Test: `internal/modules/system/module_test.go`

**Interfaces:**
- Produces in `internal/contracts/contracts.go`:
  - `var ErrNotInitialized = errors.New("contracts: system is not initialized")`
  - `type AdminCredentialStore interface { AdminPasswordHash(ctx context.Context) (string, error); UpdateAdminPasswordHash(ctx context.Context, hash string) error; RequireLogin(ctx context.Context) (bool, error); SetRequireLogin(ctx context.Context, required bool) error }`
- Produces in `internal/modules/system/module.go`:
  - `(*Module)` satisfies `contracts.AdminCredentialStore` (compile-time check `var _ contracts.AdminCredentialStore = (*Module)(nil)`).

- [ ] **Step 1: Write failing tests for `AdminCredentialStore` on `system.Module`**
  - In `internal/modules/system/module_test.go`, add `TestSystemAdminCredentialStore`:
    - Before init: `mod.AdminPasswordHash(ctx)` returns `contracts.ErrNotInitialized`; `mod.UpdateAdminPasswordHash(ctx, hash)` returns `contracts.ErrNotInitialized`; `mod.RequireLogin(ctx)` returns `true, nil`.
    - After `POST /api/init`: `mod.AdminPasswordHash(ctx)` returns stored bcrypt hash; `mod.UpdateAdminPasswordHash(ctx, newHash)` updates hash atomically and subsequent `AdminPasswordHash` returns `newHash`.
    - `mod.SetRequireLogin(ctx, false)` updates `require_login` to `"false"` and `mod.RequireLogin(ctx)` returns `false, nil`; `mod.SetRequireLogin(ctx, true)` restores `true, nil`.
  - In `internal/contracts/contracts_test.go`, add a compile-time interface assertion for `AdminCredentialStore`.

- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/contracts/... ./internal/modules/system/...`
  - Expected: FAIL (`contracts.AdminCredentialStore` / `AdminPasswordHash` undefined).

- [ ] **Step 3: Implement contract and `system` repository/module methods**
  - Add `ErrNotInitialized` and `AdminCredentialStore` to `internal/contracts/contracts.go`.
  - Add `AdminPasswordHash`, `UpdateAdminPasswordHash`, `RequireLogin`, `SetRequireLogin` on `repository` in `internal/modules/system/repo.go` (using `r.db.Reader()` for reads and `r.db.Writer()` for writes; `UpdateAdminPasswordHash` updates only when `key = 'admin_password_hash'` exists, returning `contracts.ErrNotInitialized` if 0 rows affected).
  - Expose forwarding methods on `*Module` in `internal/modules/system/module.go` with compile-time assertion `var _ contracts.AdminCredentialStore = (*Module)(nil)`.

- [ ] **Step 4: Run tests to verify pass**
  - Run: `go test -count=1 -race ./internal/contracts/... ./internal/modules/system/...`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(system): expose AdminCredentialStore contract for auth module`.

---

### Task 2: `auth` Module & OpenAPI Fragment

**Files:**
- Create: `internal/modules/auth/cookie.go`
- Create: `internal/modules/auth/module.go`
- Create: `internal/modules/auth/handlers.go`
- Create: `internal/modules/auth/module_test.go`
- Create: `contracts/openapi/modules/auth.yaml`

**Interfaces:**
- Consumes: `contracts.AdminCredentialStore`, `*crypto.SessionSigner`, `config.Config`, `*eventbus.Bus`, `crypto.VerifyPassword`, `crypto.HashPassword`
- Produces:
  - `auth.Options{Store contracts.AdminCredentialStore, Sessions *crypto.SessionSigner, Config config.Config, Bus *eventbus.Bus, SessionTTL time.Duration}`
  - `auth.New(opts Options) *Module`
  - `(*Module)` implements `lifecycle.Module` (`Name() = "auth"`, `Dependencies() = []string{"system"}`, `Migrations() = nil`) AND `contracts.AuthVerifier`.
  - Cookie name constant: `SessionCookieName = "om4_session"`, default TTL `24 * time.Hour`.
  - PublicAPI routes (under `/api`):
    - `POST /auth/login` — body `{ "password": string }`, sets `HttpOnly; SameSite=Lax; Path=/` signed cookie, returns `200 { "authenticated": true, "subject": "admin", "expires_at": RFC3339 }`.
    - `POST /auth/logout` — clears cookie (`MaxAge: -1`), returns `200 { "authenticated": false }`.
    - `GET /auth/me` — returns `200 { "authenticated": bool, "subject": string, "require_login": bool, "remote_enforced": bool, "initialized": bool }` (never 401 so frontend can inspect session state cleanly).
  - AdminAPI routes (under `/api`, protected by `AdminAPI` middleware):
    - `PUT /auth/password` — body `{ "current_password": string, "new_password": string }`, verifies current password, enforces `8..72` rune/UTF-8 byte password rules matching `system`, hashes & updates via `Store.UpdateAdminPasswordHash`, re-issues fresh session cookie, returns `200 { "updated": true }`.
    - `PUT /auth/policy` — body `{ "require_login": bool }`, rejects `require_login: false` with `403 remote_login_required` when `!cfg.IsLoopback() && !cfg.AllowInsecureRemote`, persists via `Store.SetRequireLogin`, publishes `auth.policy_changed` on `eventbus`, returns `200 { "require_login": bool, "remote_enforced": bool }`.

- [ ] **Step 1: Write failing tests in `internal/modules/auth/module_test.go`**
  - Cover all happy paths + Review Focus #1 and #2:
    - `TestAuthLoginLifecycle`: uninitialized -> `409 not_initialized`; wrong password -> `401 invalid_credentials`; valid password -> `200` + `Set-Cookie` (`HttpOnly`, `SameSite=Lax`, `Secure` when TLS / `X-Forwarded-Proto: https`); `GET /auth/me` with cookie -> `authenticated: true`; `POST /auth/logout` clears cookie -> `GET /auth/me` returns `authenticated: false`.
    - `TestAuthVerifierCookieValidation`: valid cookie -> `("admin", nil)`; missing/malformed/tampered/expired cookie -> error.
    - `TestAuthRemoteBindEnforcement` _(Review Focus #1)_: on `Host: "0.0.0.0"` with `AllowInsecureRemote: false`, `PUT /api/auth/policy` with `{"require_login":false}` returns `403 remote_login_required`, AND `mod.RequireLogin(ctx)` returns `true` even if store returns `false`. On `Host: "127.0.0.1"`, toggling `require_login: false` succeeds (`200`) and unlocks AdminAPI without cookie.
    - `TestAuthPasswordChange`: wrong `current_password` -> `401`; short/too-long `new_password` -> `400`; valid rotation -> old password rejected on login, new password accepted.

- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/auth/...`
  - Expected: FAIL (package does not exist).

- [ ] **Step 3: Implement `auth` module and `contracts/openapi/modules/auth.yaml`**
  - Implement `cookie.go`, `module.go`, `handlers.go`, and the OpenAPI 3.1 fragment documenting `/api/auth/login`, `/api/auth/logout`, `/api/auth/me`, `/api/auth/password`, `/api/auth/policy` and their schemas.

- [ ] **Step 4: Run tests to verify pass**
  - Run: `go test -count=1 -race ./internal/modules/auth/...`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(auth): implement single-admin cookie auth module and OpenAPI contract`.

---

### Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/settings/migrations/0001_init.sql`
- Create: `internal/modules/settings/schema.go`
- Create: `internal/modules/settings/repo.go`
- Create: `internal/modules/settings/module.go`
- Create: `internal/modules/settings/module_test.go`
- Create: `contracts/openapi/modules/settings.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`
- Produces:
  - `settings.Options{DB *db.DB, Bus *eventbus.Bus}`
  - `settings.New(opts Options) *Module` (`Name() = "settings"`, `Dependencies() = []string{"system"}`)
  - SQLite tables (`0001_init.sql`):
    - `settings (key TEXT PRIMARY KEY, value_json TEXT NOT NULL, revision INTEGER NOT NULL DEFAULT 1, updated_at TEXT NOT NULL, updated_by TEXT NOT NULL)`
    - `settings_history (id INTEGER PRIMARY KEY AUTOINCREMENT, setting_key TEXT NOT NULL, old_value_json TEXT, new_value_json TEXT NOT NULL, revision INTEGER NOT NULL, changed_by TEXT NOT NULL, reason TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)` + index `idx_settings_history_key_created ON settings_history(setting_key, created_at DESC)`.
  - Built-in typed setting definitions (`schema.go`) with defaults & strict validators:
    - `ui.theme` (`"dark"` | `"light"` | `"system"`, default `"dark"`)
    - `ui.locale` (`"en"` | `"id"`, default `"en"`)
    - `usage.payload_logging_enabled` (`bool`, default `false`)
    - `usage.payload_retention_days` (`int`, `1..90`, default `7`)
    - `usage.metadata_retention_days` (`int`, `1..3650`, default `90`)
    - `gateway.default_timeout_seconds` (`int`, `5..600`, default `120`)
    - `gateway.allow_response_chain_fallback` (`bool`, default `true`)
  - AdminAPI routes (under `/api`):
    - `GET /settings` — returns all typed settings (merged defaults + persisted overrides) with metadata (`key`, `type`, `value`, `default_value`, `revision`, `updated_at`, `description`).
    - `PATCH /settings` — body `{ "updates": { "<key>": <typed-value>, ... }, "reason": string? }`, validates every key/value upfront inside a single atomic writer transaction, appends `settings_history` rows, publishes `settings.updated` on `eventbus`, returns updated settings snapshot.
    - `GET /settings/history` — optional query params `key`, `limit` (default `50`, max `200`), returns newest-first revision history items.
    - `POST /settings/rollback` — body `{ "history_id": int64, "reason": string? }`, restores `old_value_json` (or target snapshot value) atomically in a transaction, appends a new history entry recording the rollback, publishes `settings.updated`, returns the restored setting item.

- [ ] **Step 1: Write failing tests in `internal/modules/settings/module_test.go`**
  - Cover defaults, partial `PATCH`, history ordering, atomic rollback, `eventbus` publication, and Review Focus #3:
    - `TestSettingsDefaultsAndPatch`: `GET /api/settings` returns all 7 defaults; `PATCH /api/settings` updates `ui.theme` to `"light"` and `usage.payload_retention_days` to `14`; subsequent `GET` reflects updates and incremented revision.
    - `TestSettingsValidationRejectsUnknownAndInvalid` _(Review Focus #3)_: unknown key -> `400 invalid_setting_key`; out-of-range int (`usage.payload_retention_days: 999`) or wrong type (`ui.theme: 123`) -> `400 invalid_setting_value`; entire multi-key `PATCH` rolls back atomically when one entry is invalid.
    - `TestSettingsHistoryAndRollback` _(Review Focus #3)_: `GET /api/settings/history` records old/new values; `POST /api/settings/rollback` restores previous value and records rollback audit entry; rolling back nonexistent `history_id` returns `404 revision_not_found`.
    - `TestSettingsPublishesEventBus`: updating settings emits `settings.updated` event on `eventbus`.

- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/settings/...`
  - Expected: FAIL.

- [ ] **Step 3: Implement `settings` module, SQL migration, and `contracts/openapi/modules/settings.yaml`**
  - Implement `migrations/0001_init.sql`, `schema.go`, `repo.go`, `module.go`, and OpenAPI fragment.

- [ ] **Step 4: Run tests to verify pass**
  - Run: `go test -count=1 -race ./internal/modules/settings/...`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(settings): implement typed settings store with history, rollback, and OpenAPI contract`.

---

### Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/keys/migrations/0001_init.sql`
- Create: `internal/modules/keys/token.go`
- Create: `internal/modules/keys/repo.go`
- Create: `internal/modules/keys/cache.go`
- Create: `internal/modules/keys/module.go`
- Create: `internal/modules/keys/module_test.go`
- Create: `contracts/openapi/modules/keys.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`
- Produces:
  - `keys.Options{DB *db.DB, Bus *eventbus.Bus}`
  - `keys.New(opts Options) *Module` (`Name() = "keys"`, `Dependencies() = []string{"system"}`)
  - `(*Module)` implements `lifecycle.Module` AND `contracts.KeyValidator` (`var _ contracts.KeyValidator = (*Module)(nil)`).
  - Token format (`token.go`): `sk-om4-<48-char base64url CSPRNG>` (~288 bits entropy); stored in SQLite as `key_hash TEXT NOT NULL UNIQUE` (lowercase hex SHA-256) and `key_prefix TEXT NOT NULL` (first 14 chars e.g. `sk-om4-AbCdEf1`).
  - SQLite table (`0001_init.sql`):
    - `api_keys (id TEXT PRIMARY KEY, name TEXT NOT NULL, key_hash TEXT NOT NULL UNIQUE, key_prefix TEXT NOT NULL, allowed_models_json TEXT NOT NULL DEFAULT '[]', scopes_json TEXT NOT NULL DEFAULT '["gateway"]', tags_json TEXT NOT NULL DEFAULT '[]', rate_limit_rpm INTEGER NOT NULL DEFAULT 0, quota_usd REAL NOT NULL DEFAULT 0, disabled INTEGER NOT NULL DEFAULT 0, expires_at TEXT, last_used_at TEXT, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)` + index `idx_api_keys_hash ON api_keys(key_hash)`.
  - AdminAPI routes (under `/api`):
    - `GET /keys` — list all keys (metadata only — NEVER includes `key_hash` or raw token).
    - `POST /keys` — body `{ "name", "allowed_models"?, "scopes"?, "tags"?, "rate_limit_rpm"?, "quota_usd"?, "expires_at"? }`, generates token, persists hash, returns `201 { "key": ApiKeySummary, "secret": "sk-om4-..." }`.
    - `GET /keys/{id}` — single key metadata (`404` if missing).
    - `PATCH /keys/{id}` — update `name`, `allowed_models`, `scopes`, `tags`, `rate_limit_rpm`, `quota_usd`, `disabled`, `expires_at`; invalidates cache & publishes `keys.changed`.
    - `POST /keys/{id}/rotate` — generates new token + replaces `key_hash` & `key_prefix` atomically, invalidates old & new hash in cache, publishes `keys.changed`, returns `200 { "key": ApiKeySummary, "secret": "sk-om4-..." }`.
    - `DELETE /keys/{id}` — deletes key, invalidates cache & publishes `keys.changed`, returns `204 No Content`.
  - `ValidateKey(ctx, rawKey)`: computes SHA-256 hex digest, checks in-memory read-through cache (falling back to `db.Reader()`), rejects disabled or expired keys with explicit errors (`ErrKeyInvalid`, `ErrKeyDisabled`, `ErrKeyExpired`), returns `contracts.KeyInfo`.

- [ ] **Step 1: Write failing tests in `internal/modules/keys/module_test.go`**
  - Cover CRUD, validation, expiry, rotation, `allowed_models`, and Review Focus #4:
    - `TestKeysCRUDAndSecretHygiene` _(Review Focus #4)_: `POST /api/keys` returns `201` with `secret` starting with `sk-om4-`; `GET /api/keys` and `GET /api/keys/{id}` NEVER contain the raw secret or hash; SQLite `api_keys` table does not contain raw secret.
    - `TestKeysValidateAndImmediateCacheInvalidation` _(Review Focus #4)_: `ValidateKey` succeeds and populates cache; `PATCH /api/keys/{id}` with `disabled: true` makes next `ValidateKey` immediately fail with `ErrKeyDisabled`; re-enabling works; `POST /api/keys/{id}/rotate` immediately invalidates old secret and accepts new secret; `DELETE /api/keys/{id}` immediately rejects secret.
    - `TestKeysExpirationAndValidationErrors`: expired key fails `ValidateKey` with `ErrKeyExpired`; empty name or negative `rate_limit_rpm`/`quota_usd` returns `400`.

- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/keys/...`
  - Expected: FAIL.

- [ ] **Step 3: Implement `keys` module, SQL migration, cache, and `contracts/openapi/modules/keys.yaml`**
  - Implement `migrations/0001_init.sql`, `token.go`, `repo.go`, `cache.go`, `module.go`, and OpenAPI fragment.

- [ ] **Step 4: Run tests to verify pass**
  - Run: `go test -count=1 -race ./internal/modules/keys/...`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(keys): implement gateway API keys module, KeyValidator cache, and OpenAPI contract`.

---

### Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Modify: `web/src/shared/api/schema.d.ts` (generated via `make codegen`)

**Interfaces:**
- Consumes: `system.New`, `auth.New`, `settings.New`, `keys.New`, `httpserver.New`
- Produces:
  - Full Wave 1A module registration (`system` -> `auth`, `settings`, `keys`) with `httpserver.Options{Auth: authMod}` enforcing session auth on all AdminAPI routes (`/api/shutdown`, `/api/auth/password`, `/api/auth/policy`, `/api/settings*`, `/api/keys*`).
  - Regenerated `web/src/shared/api/schema.d.ts` with all Wave 1A paths and component schemas.

- [ ] **Step 1: Write failing integration tests in `cmd/omeg4mon-router/main_test.go`**
  - Add `TestAppWave1AAuthSettingsKeysEndToEnd`:
    - Boot `buildApp` with temp `DataDir` and `AdminPassword: "test-admin-pass-2026"`.
    - Verify `GET /api/health` includes `"system"`, `"auth"`, `"settings"`, `"keys"` in `modules`.
    - Verify unauthenticated `GET /api/keys` and `GET /api/settings` return `401 Unauthorized` when `require_login=true`.
    - Authenticate via `POST /api/auth/login` -> capture `om4_session` cookie.
    - With cookie: `GET /api/settings` returns `200`, `PATCH /api/settings` updates `ui.theme`, `POST /api/keys` creates key, `POST /api/auth/logout` clears session -> subsequent `GET /api/keys` returns `401`.

- [ ] **Step 2: Run test to verify failure**
  - Run: `go test -race -run TestAppWave1AAuthSettingsKeysEndToEnd ./cmd/omeg4mon-router/...`
  - Expected: FAIL (`404` on `/api/auth/login` or unauthenticated `GET /api/keys` not mounted).

- [ ] **Step 3: Wire modules in `cmd/omeg4mon-router/main.go`, update `contracts/openapi/openapi.yaml`, and run `make codegen`**
  - Construct `sysMod` and `authMod` before `httpserver.New`, pass `Auth: authMod` to `httpserver.Options`, register `sysMod`, `authMod`, `settingsMod`, `keysMod` in `registry`.
  - Add all `auth`, `settings`, and `keys` paths, tags, and schemas to `contracts/openapi/openapi.yaml`.
  - Run `make codegen` to regenerate `web/src/shared/api/schema.d.ts`.

- [ ] **Step 4: Run Go tests & Web typecheck**
  - Run: `go test -count=1 -race ./cmd/omeg4mon-router/... && cd web && bun run typecheck`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(cmd): wire auth, settings, and keys modules with OpenAPI TypeScript codegen`.

---

### Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)

**Files:**
- Create: `web/src/features/auth/schemas.ts`
- Create: `web/src/features/auth/schemas.test.ts`
- Create: `web/src/features/auth/hooks.ts`
- Create: `web/src/features/auth/LoginForm.tsx`
- Create: `web/src/features/auth/AuthGuard.tsx`
- Modify: `web/src/app/login/page.tsx`
- Modify: `web/src/app/page.tsx`
- Modify: `web/src/app/(dashboard)/dashboard/layout.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`

**Interfaces:**
- Consumes: typed `apiClient` (`GET /api/auth/me`, `POST /api/auth/login`, `POST /api/auth/logout`), `useInitStatus`, `useUIStore`, `shadcn/ui` primitives, `react-hook-form` + `zodResolver`
- Produces:
  - `useAuthSession()`, `useLoginMutation()`, `useLogoutMutation()` in `web/src/features/auth/hooks.ts`
  - `<AuthGuard>` wrapping dashboard routes: checks `GET /api/auth/me`; if `!initialized` -> redirects to `/init/`; if `(require_login || remote_enforced) && !authenticated` -> redirects to `/login/`; otherwise renders children.
  - Active cookie-backed `<LoginForm />` on `/login/` replacing the Wave 0 standby placeholder; sign-out button in `AppShell` header/sidebar.

- [ ] **Step 1: Write failing unit tests in `web/src/features/auth/schemas.test.ts`**
  - Test `loginSchema` (`zod`) rejects empty password, accepts valid password, and `resolvePostAuthRedirect` helper _(Review Focus #5)_ produces deterministic non-looping routes for `(initialized, authenticated, requireLogin, remoteEnforced)`.

- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/auth/schemas.test.ts`
  - Expected: FAIL (module not found).

- [ ] **Step 3: Implement `web/src/features/auth/*` and integrate into `/login/`, `/`, dashboard layout, and `AppShell`**
  - Build `LoginForm` with `react-hook-form` + `zodResolver(loginSchema)` + `shadcn/ui` (`Input`, `Label`, `Button`, `Alert`), displaying structured Go `ErrorEnvelope` messages on `401`/`409`.
  - Wrap `/dashboard/*` layout with `<AuthGuard>` and add working Sign Out control in `AppShell`.
  - Update `/` (`web/src/app/page.tsx`) to route via `useAuthSession()` (`!initialized` -> `/init/`, needs login -> `/login/`, otherwise -> `/dashboard/`).

- [ ] **Step 4: Run Web tests & typecheck**
  - Run: `cd web && bun test && bun run typecheck`
  - Expected: PASS.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement live operator login, auth guard, and sign-out flow`.

---

### Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification

**Files:**
- Create: `web/src/features/keys/schemas.ts`
- Create: `web/src/features/keys/schemas.test.ts`
- Create: `web/src/features/keys/hooks.ts`
- Create: `web/src/features/keys/KeysConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/keys/page.tsx`
- Create: `web/src/features/settings/schemas.ts`
- Create: `web/src/features/settings/schemas.test.ts`
- Create: `web/src/features/settings/hooks.ts`
- Create: `web/src/features/settings/SettingsConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/settings/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx` (promote `API Keys` -> `/dashboard/keys/` and `Settings` -> `/dashboard/settings/` into `ACTIVE_NAV_ITEMS`)
- Modify: `web/src/features/system/SystemOverview.tsx` (update active vs planned Wave 1 module indicators)

**Interfaces:**
- Consumes: typed `apiClient` (`/api/keys*`, `/api/settings*`, `/api/auth/password`, `/api/auth/policy`), `useUIStore`, `react-hook-form` + `zodResolver`, `shadcn/ui` (`Card`, `Button`, `Input`, `Label`, `Badge`, `Alert`, `EmberRule`)
- Produces:
  - `/dashboard/keys/` — create key form (`react-hook-form` + `zod`), one-time secret reveal panel with copy-to-clipboard button, table of active/disabled keys with prefix, allowed models, RPM/quota badges, inline toggle enable/disable, rotate (with one-time secret reveal), and delete.
  - `/dashboard/settings/` — typed settings form (`react-hook-form` + `zod`, syncing `ui.theme` with Zustand `useUIStore`), operator security section (change admin password form + `require_login` toggle showing remote-bind lock status), and revision history table with one-click Rollback action.

- [ ] **Step 1: Write failing schema & helper tests in `web/src/features/keys/schemas.test.ts` and `web/src/features/settings/schemas.test.ts`**
  - Verify `createKeySchema` parses comma-separated `allowed_models`/`tags`, rejects blank name or negative RPM/quota.
  - Verify `settingsFormSchema` enforces `payload_retention_days` (`1..90`), `metadata_retention_days` (`1..3650`), `default_timeout_seconds` (`5..600`), and `changePasswordSchema` enforces `8..72` rune/UTF-8 byte parity with Go.

- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/keys/schemas.test.ts src/features/settings/schemas.test.ts`
  - Expected: FAIL.

- [ ] **Step 3: Implement `KeysConsole`, `SettingsConsole`, dashboard route pages, and update `AppShell` navigation**
  - Adhere strictly to `docs/DESIGN.md` (solid opaque surfaces, `#E56A4A` accent, Geist Sans/Mono, zero emojis, WCAG AA contrast, no `backdrop-blur`).

- [ ] **Step 4: Run full verification suite (`make test` + `make build` + binary smoke check)**
  - Run: `make test && make build`
  - Verify static export creates `web/out/dashboard/keys/index.html` and `web/out/dashboard/settings/index.html`, and the compiled single binary `bin/omeg4mon-router` serves `/api/health`, `/login/`, `/dashboard/keys/`, and `/dashboard/settings/`.

- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement API keys console, typed settings rollback UI, and Wave 1A navigation`.
