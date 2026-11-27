# Provider Quota Menu & Token Harbor Quota Integration

**Date:** 2026-11-27
**Status:** Approved

## Intent & Why

Operators managing multiple upstream provider API keys / OAuth accounts in `omeg4mon-router` need a dedicated **Quota** dashboard view to inspect remaining upstream provider allowances per connection without logging into each provider's portal. Starting with **Token Harbor** (`tokenharbor.ai`), operators can immediately monitor both subscription pass (`plan_used_pct`) and free-tier (`used_pct`) allowances, activation posture, eligible free models, and rolling 7-day reset timestamps across every configured Token Harbor key.

## Goals

- Add a dedicated **Quota** sidebar menu item (`/dashboard/quota/`) in the **Operations** section with full `en`/`id` localization and Dark/Light theme parity.
- Add a built-in `tokenharbor` provider adapter (`https://tokenharbor.ai/v1`) while also automatically detecting existing custom/compatible connections targeting `tokenharbor.ai`.
- Fetch live, read-only Token Harbor quota via `GET <origin>/api/cli/free-tier` with `Authorization: Bearer <API_KEY>` (never exposing upstream keys to the browser).
- Structure the backend quota fetcher registry and response contract generically (`buckets[]`) so additional providers (Antigravity, OpenRouter, Codex, Claude OAuth, etc.) can plug in later without breaking the UI or API contract.
- Support 15-second bounded in-memory caching, manual refresh-all, and single-connection force refresh.

## Non-Goals & Anti-Goals

- **Non-goal:** Historical quota time-series storage in SQLite or background cron polling.
- **Non-goal:** Passive inference response header capture (`x-ratelimit-*`) for generic API-key providers in this iteration.
- **Anti-goal:** Leaking raw upstream API keys in API responses, error messages, or logs.
- **Anti-goal:** Blocking the entire Quota view when a single upstream connection fails or times out.

## Approaches Considered

1. **Approach A (Selected): Extensible `quotaFetcher` registry inside `internal/modules/providers` + `/dashboard/quota/` UI.**
   - Reuses `providers` module's credential vault, `httpClient` (with outbound proxy resolution), and adapter catalog with zero new Go modules or SQLite migrations.
   - Trade-off: No historical quota charts yet, which is unnecessary for live allowance inspection.
2. **Approach B: Standalone `internal/modules/quota/` module with SQLite snapshot table and background polling.**
   - Rejected as overengineered for on-demand allowance viewing (Token Harbor already caches `/api/cli/free-tier` for 15s server-side).
3. **Approach C: Direct browser fetch or Token Harbor-only response schema.**
   - Rejected for security (keys must remain server-side) and maintainability (hardcoding one provider's JSON shape would force UI rewrites for future providers).

## Architecture & Contracts

### Backend (`internal/modules/providers`)

- **Built-in Adapter (`adapter.go`):**
  - `id: "tokenharbor"`, `displayName: "Token Harbor"`, `category: "api_key"`, `protocol: "openai-chat"`, `authType: "api_key"`, `defaultBaseURL: "https://tokenharbor.ai/v1"`, default models `["deepseek-v4.1-flash:free", "mimo-v2.6-flash:free"]`.
- **Quota Registry & Cache (`quota.go`):**
  - `quotaFetcher` interface: `ID() string`, `Matches(conn ConnectionSummary) bool`, `Fetch(ctx, client, conn, creds) (QuotaFetchResult, error)`.
  - `tokenHarborQuotaFetcher`: matches `provider_id == "tokenharbor"` (case-insensitive) OR any connection whose `base_url` host is `tokenharbor.ai` / `*.tokenharbor.ai`. Resolves `<scheme>://<host>/api/cli/free-tier` from `creds.BaseURL` (fallback `https://tokenharbor.ai`).
  - Parses `plan_name`, `on_plan`, `plan_used_pct`, `plan_exhausted`, `used_pct`, `exhausted`, `activated`, `period_started_at`, `period_ends_at`, `reset_at`, `models`.
  - Produces normalized `buckets[]` (`plan` when `on_plan == true`, plus `free_tier`), clamping percentages to `[0, 100]` and computing `remaining_pct = 100 - used_pct`.
  - Derives connection status: `"ok"`, `"low"` (any active bucket `remaining_pct <= 20`), `"exhausted"`, `"not_activated"` (`!on_plan && !activated`), or `"error"`.
  - Bounded in-memory cache (`maxQuotaCacheEntries = 512`, TTL `15s`), invalidated on connection patch/delete.
  - Bounded parallel fan-out (`maxQuotaConcurrency = 6`, per-connection timeout `10s`, response body limit `256 KiB`), honoring `ctx.Done()`.
- **Admin HTTP Routes (`module.go` + OpenAPI):**
  - `GET /api/providers/quotas?provider_id=&refresh=true|false` -> `ProviderQuotasResponse`
  - `POST /api/providers/connections/{id}/quota/refresh` -> `ProviderConnectionQuota`

### Frontend (`web/`)

- **Sidebar (`web/src/shared/ui/navigation.ts`):** Add `Quota` / `Kuota` (`/dashboard/quota/`, `moduleId: "providers"`, icon `Gauge`) to `OPERATIONS_ITEMS`.
- **Router (`web/src/App.tsx`):** Register `"quota"` in `DASHBOARD_ROUTE_PATHS` and lazy route to `QuotaConsole`.
- **Feature (`web/src/features/quota/`):**
  - `schemas.ts` + `schemas.test.ts`: Types from generated `schema.d.ts`, summary stat aggregation, status/provider/search filtering, and reset countdown formatting.
  - `hooks.ts`: TanStack Query hooks using typed `apiClient` (`useProviderQuotas`, `useRefreshAllQuotasMutation`, `useRefreshConnectionQuotaMutation`).
  - `QuotaConsole.tsx`: Summary KPI cards, filter/search controls, Refresh All, per-connection quota cards with bucket progress bars, plan/activation/model badges, per-connection refresh, loading skeletons, inline error alerts, and empty-state CTA to `/dashboard/providers/`.
- **Provider Registry (`web/src/features/providers/providerRegistry.ts`):** Add `tokenharbor` metadata entry (`https://tokenharbor.ai`, API key portal URL).

## Constraints

- Strict OpenAPI contract sync (`make codegen` updates `web/src/shared/api/schema.d.ts`).
- Zero hardcoded colors in React components; use semantic CSS tokens from `web/src/app/globals.css`.
- Preserve `en`/`id` localization and Dark/Light theme contrast.
- Sanitize all upstream error messages so credentials never appear in API responses.
