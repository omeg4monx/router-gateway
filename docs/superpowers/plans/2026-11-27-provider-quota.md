# Provider Quota Menu & Token Harbor Quota Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans or superpowers:subagent-driven-development to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a new **Quota** sidebar menu (`/dashboard/quota/`) backed by an extensible provider quota API in `internal/modules/providers`, initially supporting **Token Harbor** (`GET /api/cli/free-tier`).

**Architecture:** Register a built-in `tokenharbor` adapter and an extensible `quotaFetcher` registry with 15-second bounded in-memory caching and bounded parallel fan-out inside `internal/modules/providers`. Expose `GET /api/providers/quotas` and `POST /api/providers/connections/{id}/quota/refresh` via OpenAPI + `make codegen`, then render them in a responsive, localized `QuotaConsole` under `web/src/features/quota/`.

**Tech Stack:** Go 1.24 (`chi`, `net/http`, `sync`), OpenAPI 3.1 + `openapi-typescript`, React 19, TypeScript, TanStack Query, Tailwind CSS, Zustand, Bun test.

**Spec:** `docs/superpowers/specs/2026-11-27-provider-quota-design.md`

## Global Constraints

- Never expose or echo upstream API keys / access tokens in HTTP responses or error strings.
- All in-memory caches and concurrent fan-out workers must be strictly bounded and respect `ctx.Done()`.
- OpenAPI changes must be synced via `make codegen` (`web/src/shared/api/schema.d.ts`).
- UI components must reuse `web/src/shared/ui/` primitives and CSS custom properties from `web/src/app/globals.css` with `en`/`id` locale support.
- All checks in `make check` (`gofmt`, `golangci-lint`, `biome`, `go test -race`, `CGO_ENABLED=0`, `bun test`, `tsc --noEmit`, `cli` tests) and `cd web && bun run build` must pass.

## Review Focus

1. Custom `openai-compatible` / custom-node Token Harbor connection whose `provider_id` is NOT `"tokenharbor"` but `base_url` is `https://tokenharbor.ai/v1` — must still be matched and queried.
2. Token Harbor response with `on_plan: false`, `activated: false`, `plan_name: null` — must report `status: "not_activated"` and single `free_tier` bucket without nil-pointer panic.
3. Upstream returns HTTP 401/500 or invalid JSON containing the raw API key — must return per-connection `status: "error"` with sanitized error string while other accounts still succeed.
4. Cached quota entry must be invalidated when a connection's API key/URL is patched or deleted, and bypassed on `?refresh=true` or `POST .../quota/refresh`.
5. Navigation contract tests (`navigation.test.ts`, `App.test.ts`) must reflect the new `/dashboard/quota/` route with unique trailing-slash href and `providers` module registration gating.

---

### Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD)

**Files:**
- Modify: `internal/modules/providers/adapter.go`
- Create: `internal/modules/providers/quota.go`
- Modify: `internal/modules/providers/module.go`
- Modify: `internal/modules/providers/handlers.go` (invalidate quota cache on connection patch/delete)
- Create: `internal/modules/providers/quota_test.go`

**Interfaces:**
- Produces:
  - Built-in adapter `"tokenharbor"` in `defaultBuiltInAdapters()`
  - `GET /api/providers/quotas` (`ProviderQuotasResponse`)
  - `POST /api/providers/connections/{id}/quota/refresh` (`ProviderConnectionQuota`)

- [ ] **Step 1: Write failing Go tests in `internal/modules/providers/quota_test.go`** covering:
  - Built-in `tokenharbor` adapter presence in `/api/providers/catalog`
  - `GET /api/providers/quotas` for built-in `tokenharbor` connection (plan + free tier buckets, `remaining_pct`, `reset_at`, `models`)
  - Auto-detection of custom node / `openai-compatible` connection with `tokenharbor.ai` host (via custom transport/round-tripper or helper test)
  - `not_activated` (`on_plan: false, activated: false`), `low` (`used_pct: 85`), `exhausted` (`exhausted: true`) states
  - Per-connection error isolation + API key sanitization when upstream returns 401 with key echoed in body
  - 15s cache hit, `?refresh=true` bypass, `POST /api/providers/connections/{id}/quota/refresh`, and cache invalidation on connection `PATCH`
  - Unsupported connection on single-connection refresh returns HTTP 400; unknown connection returns HTTP 404
- [ ] **Step 2: Run `go test ./internal/modules/providers -run TestQuota -v`** and confirm RED failure.
- [ ] **Step 3: Implement `tokenharbor` adapter in `adapter.go`, quota engine in `quota.go`, route wiring in `module.go`, and cache invalidation in `handlers.go`.**
- [ ] **Step 4: Run `go test -count=1 -race ./internal/modules/providers/...`** and verify PASS.

---

### Task 2: OpenAPI Contract & TypeScript Schema Codegen

**Files:**
- Modify: `contracts/openapi/modules/providers.yaml`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1: Add `/api/providers/quotas` and `/api/providers/connections/{id}/quota/refresh` paths and `ProviderQuotaBucket`, `ProviderConnectionQuota`, `ProviderQuotasResponse` schemas to `contracts/openapi/modules/providers.yaml` and wire `$ref`s in `contracts/openapi/openapi.yaml`.**
- [ ] **Step 2: Run `make codegen`** and verify `web/src/shared/api/schema.d.ts` includes the new paths and schemas cleanly.

---

### Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD)

**Files:**
- Modify: `web/src/shared/ui/navigation.ts`
- Modify: `web/src/shared/ui/navigation.test.ts`
- Modify: `web/src/App.tsx`
- Modify: `web/src/features/providers/providerRegistry.ts`

- [ ] **Step 1: Update `web/src/shared/ui/navigation.test.ts`** to expect 24 nav items including `/dashboard/quota/` under `operations` (gated by `providers` module, searchable by `"quota"`, `"kuota"`, `"tokenharbor"`).
- [ ] **Step 2: Run `cd web && bun test src/shared/ui/navigation.test.ts`** and confirm RED failure.
- [ ] **Step 3: Add Quota nav item to `OPERATIONS_ITEMS` in `navigation.ts`, add `tokenharbor` spec in `providerRegistry.ts`, and add `"quota"` route in `App.tsx`.**
- [ ] **Step 4: Run `cd web && bun test src/shared/ui/navigation.test.ts src/App.test.ts`** (after Task 4 provides `QuotaConsole` or stub) and verify PASS.

---

### Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD)

**Files:**
- Create: `web/src/features/quota/schemas.ts`
- Create: `web/src/features/quota/schemas.test.ts`
- Create: `web/src/features/quota/hooks.ts`
- Create: `web/src/features/quota/QuotaConsole.tsx`

- [ ] **Step 1: Write failing unit tests in `web/src/features/quota/schemas.test.ts`** testing:
  - `computeQuotaOverviewStats` across `ok`, `low`, `exhausted`, `not_activated`, `error` connections
  - `filterProviderQuotas` by provider, status, and search query (connection name, plan name, model ID)
  - `formatResetCountdown` for future dates (days/hours/minutes), expired/past dates, and missing/invalid strings
  - `getQuotaStatusVariant` and `getBucketProgressTone`
- [ ] **Step 2: Run `cd web && bun test src/features/quota/schemas.test.ts`** and confirm RED failure.
- [ ] **Step 3: Implement `schemas.ts`, `hooks.ts`, and `QuotaConsole.tsx`** with full `en`/`id` copy, semantic CSS tokens, responsive layout, loading/error/empty states, per-connection & global refresh, and clean TanStack Query invalidation.
- [ ] **Step 4: Run `cd web && bun test`** and verify all frontend tests pass.

---

### Task 5: Full Verification, Web Build, Graph Update & Versioning Audit

- [ ] **Step 1: Run `make check`** (`fmt-check`, `lint`, `go test -race`, `CGO_ENABLED=0`, `bun test`, `typecheck`, `cli` tests).
- [ ] **Step 2: Run `cd web && bun run build`** and verify bundle builds cleanly.
- [ ] **Step 3: Run `graphify update .`** to keep the knowledge graph current.
- [ ] **Step 4: Complete `omeg4mon-verify` and `omeg4mon-versioning` exit tables** with fresh command output.
