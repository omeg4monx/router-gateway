# Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the core LLM routing kernel of `omeg4mon-router` — pluggable `ProviderAdapter` registry (generic API-key executors + 3 headless-safe OAuth providers: Claude, Codex, Gemini with encrypted vault), deterministic multi-model fallback combos/aliases/cooldowns (`combos`), and multi-protocol streaming SSE `/v1/*` + `/v1beta/*` gateway (`gateway`), plus matching Next.js 16 static-export dashboard studios.

**Architecture:** Three isolated vertical-slice Go modules (`internal/modules/providers`, `internal/modules/combos`, `internal/modules/gateway`) implement `lifecycle.Module`, own their own SQLite migrations, communicate exclusively via narrow interfaces in `internal/contracts/` wired in `cmd/omeg4mon-router/main.go`, and invalidate hot-path caches via `eventbus`. Modular OpenAPI 3.1 fragments (`contracts/openapi/modules/{providers,combos}.yaml`) drive strict TypeScript client generation for `web/src/features/{providers,combos}/`.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `http.ResponseController`, `bufio.Reader`, `modernc.org/sqlite` zero-CGO WAL, AES-256-GCM `crypto.Vault`, PKCE S256), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), `zustand`, TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying of source code, SQL, or third-party assets from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`. Never hardcode third-party OAuth client secrets into source code (configurable via env/connection config).
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module calls use narrow interfaces in `internal/contracts/` wired in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL:** Pure-Go `modernc.org/sqlite` using `db.Writer()` (max 1 conn) for writes and `db.Reader()` (max 8 conns) for reads; fixed-width UTC RFC3339 timestamps (`"2006-01-02T15:04:05.000000000Z"`) for deterministic chronological `ORDER BY` queries.
- **Streaming discipline:** Zero full-body buffering on SSE responses; flush every SSE event immediately via `http.NewResponseController(w).Flush()`; respect `r.Context()` cancellation and close upstream response bodies cleanly (zero goroutine/connection leaks).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`); dynamic detail pages use static route shells + query params (`/dashboard/providers/detail/?id=...`); strict adherence to `docs/DESIGN.md` (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task must pass `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, and `cd web && bun run check` (Biome v2) — verified via `make check` before commit — plus `graphify update .` (cleaning `graphify-out/20??-??-??/` backups).

## Review Focus

1. **Secret encryption & leakage prevention in `providers`:** API keys and OAuth `access_token`/`refresh_token` MUST be sealed with AES-256-GCM (`*crypto.Vault`) before SQLite storage, NEVER returned in plaintext by any HTTP endpoint or emitted on `eventbus`, and OAuth token refresh on concurrent `401`s MUST coalesce per connection ID so refresh tokens are not raced/invalidated. _(Pinned in Task 2)_
2. **Deterministic 4-tier resolution precedence & `allowed_models` authorization in `combos`:** Resolution MUST follow `(1) Active Combo → (2) Model Alias → (3) Explicit <provider>/<model> → (4) Bare <model> by connection priority`; non-empty `allowed_models` authorizes the **requested identifier** (authorizing a combo permits its member targets, while rejecting unauthorized direct targets with `ErrModelForbidden`); disabled models and accounts in cooldown are skipped cleanly. _(Pinned in Task 3)_
3. **Cross-protocol Canonical IR fidelity & direct pass-through preservation:** When inbound protocol matches upstream protocol, fast-path pass-through preserves provider-specific fields (e.g. Anthropic `cache_control`, `thinking`) verbatim; when translating across OpenAI Chat / Anthropic Messages / OpenAI Responses / Gemini, parallel tool calls, system prompts, reasoning/thinking deltas, and SSE `[DONE]` / `message_stop` framing translate losslessly. _(Pinned in Task 4)_
4. **Error classification, single-retry OAuth refresh & zero-leak SSE failover in `gateway`:** Client errors (`400`, `422`, context length) return immediately without burning fallback accounts; `401` on an OAuth connection triggers at most ONE `RefreshOAuthToken` retry before failing over; `429`/`5xx`/network errors record cooldown (`Retry-After` respected) and advance to the next candidate; client disconnect cancels upstream request immediately without panics or leaked bodies. _(Pinned in Task 5)_
5. **Headless-safe OAuth wizard & static-export detail routing in `web/`:** Dashboard OAuth flow works without any local browser popup on the server (Start -> copy/open Authorize URL -> paste callback URL or authorization code -> Exchange), and `/dashboard/providers/detail/?id=...` + `/dashboard/models/` export as static trailing-slash HTML files that survive direct page refresh. _(Pinned in Tasks 7 & 8)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Contracts | `internal/contracts/contracts.go`, `internal/contracts/contracts_test.go` | Add `ProviderConnection`, `UpstreamCredentials`, `ProviderStore`, `AccountHealthTracker`, `ModelDescriptor`, `ModelLister`, sentinel errors |
| Providers module | `internal/modules/providers/{migrations/0001_init.sql,adapter.go,oauth.go,repo.go,module.go,handlers.go,module_test.go}`, `contracts/openapi/modules/providers.yaml` | Pluggable `ProviderAdapter` registry (4 generic API-key + 3 OAuth: Claude, Codex, Gemini), encrypted vault CRUD, custom provider nodes, quick-test, headless-safe OAuth + singleflight refresh |
| Combos module | `internal/modules/combos/{migrations/0001_init.sql,cooldown.go,resolver.go,repo.go,module.go,module_test.go}`, `contracts/openapi/modules/combos.yaml` | Combos, aliases, disabled models CRUD, 4-tier `ProviderResolver`, `ModelLister`, `AccountHealthTracker` backoff state machine |
| Gateway translators | `internal/modules/gateway/{ir.go,translate_openai.go,translate_anthropic.go,translate_responses.go,translate_gemini.go,tokens.go,translate_test.go}` | Canonical IR, direct fast-path detection, bidirectional protocol translators, direct Anthropic ↔ Responses pair translator, local `count_tokens` estimator |
| Gateway pipeline | `internal/modules/gateway/{migrations/0001_init.sql,auth.go,classify.go,stream.go,chains.go,module.go,module_test.go}` | `/v1/*`, `/v1beta/*`, compatibility aliases, API-key auth, `GatewayInterceptor` chain, streaming SSE executor, OAuth refresh retry & failover, bounded `response_chains` |
| Wiring & E2E | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Register `providers`, `combos`, `gateway`, aggregate OpenAPI, run `make codegen`, live mock-upstream E2E tests |
| Web Providers UI | `web/src/features/providers/{schemas.ts,schemas.test.ts,hooks.ts,ProvidersConsole.tsx,ProviderDetailView.tsx}`, `web/src/app/(dashboard)/dashboard/providers/{page.tsx,detail/page.tsx}` | Providers Control Studio, API-key & custom node forms, headless-safe OAuth wizard, quick-test, detail query-param page |
| Web Models/Combos UI | `web/src/features/combos/{schemas.ts,schemas.test.ts,hooks.ts,ModelsRoutingConsole.tsx}`, `web/src/app/(dashboard)/dashboard/models/page.tsx`, `web/src/shared/ui/AppShell.tsx`, `web/src/features/system/SystemOverview.tsx` | Combos builder, aliases, disabled models, cooldown monitor, endpoint snippet panel, nav promotion & full single-binary verification |

---

### Task 1: Cross-Module Contracts for Providers, Combos & Gateway

**Files:**
- Modify: `internal/contracts/contracts.go`
- Modify: `internal/contracts/contracts_test.go`

**Interfaces:**
- Produces in `internal/contracts/contracts.go`:
  - Sentinel errors: `ErrModelNotFound`, `ErrModelForbidden`, `ErrAllAccountsCoolingDown`, `ErrConnectionNotFound`, `ErrOAuthRefreshUnsupported`.
  - `type ProviderConnection struct { ID, ProviderID, Name, Protocol, BaseURL, AuthType string; Priority int; Enabled bool; Models []string; ExtraHeaders map[string]string }`
  - `type UpstreamCredentials struct { ConnectionID, ProviderID, Protocol, BaseURL, AuthType, APIKey, AccessToken string; ExtraHeaders map[string]string }`
  - `type ProviderStore interface { ListEnabledConnections(ctx context.Context) ([]ProviderConnection, error); ResolveCredentials(ctx context.Context, connectionID string) (UpstreamCredentials, error); RefreshOAuthToken(ctx context.Context, connectionID string) (UpstreamCredentials, error) }`
  - `type AccountHealthTracker interface { IsInCooldown(connectionID string, now time.Time) bool; RecordFailure(connectionID string, statusCode int, retryAfter time.Duration); RecordSuccess(connectionID string) }`
  - `type ModelDescriptor struct { ID, OwnedBy, Kind string }` (`Kind`: `"combo"` | `"alias"` | `"model"`)
  - `type ModelLister interface { ListModels(ctx context.Context, allowed []string) ([]ModelDescriptor, error) }`

- [ ] **Step 1: Write failing contract assertion tests in `internal/contracts/contracts_test.go`**
  - Add compile-time assertions for `ProviderStore`, `AccountHealthTracker`, `ModelLister` and verify sentinel error distinctness (`!errors.Is` across all sentinels).
- [ ] **Step 2: Run test to verify failure**
  - Run: `go test -race ./internal/contracts/...`
  - Expected: FAIL (new interfaces/errors undefined).
- [ ] **Step 3: Implement types, interfaces, and sentinel errors in `internal/contracts/contracts.go`**
  - Keep package stdlib-only (`context`, `errors`, `net/http`, `net/url`, `time`).
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/contracts/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(contracts): add ProviderStore, AccountHealthTracker, and ModelLister interfaces`.

---

### Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/providers/migrations/0001_init.sql`
- Create: `internal/modules/providers/adapter.go`
- Create: `internal/modules/providers/oauth.go`
- Create: `internal/modules/providers/repo.go`
- Create: `internal/modules/providers/module.go`
- Create: `internal/modules/providers/handlers.go`
- Create: `internal/modules/providers/module_test.go`
- Create: `contracts/openapi/modules/providers.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*crypto.Vault`, `*http.Client`, `*eventbus.Bus`
- Produces:
  - `providers.Options{DB *db.DB, Vault *crypto.Vault, HTTPClient *http.Client, Bus *eventbus.Bus, OAuthOverrides map[string]OAuthConfig}`
  - `providers.New(opts Options) *Module` (`Name() = "providers"`, `Dependencies() = []string{"system"}`)
  - `(*Module)` implements `lifecycle.Module` AND `contracts.ProviderStore`.
  - Pluggable `ProviderAdapter` interface (`ID()`, `DisplayName()`, `Category()`, `Protocol()`, `AuthType()`, `DefaultBaseURL()`, `DefaultModels()`, `OAuthConfig() *OAuthConfig`, `BuildTestRequest(ctx, creds) (*http.Request, error)`).
  - Built-in adapters registered at construction (extensible via `RegisterAdapter(ProviderAdapter)`):
    - Generic API-key adapters: `openai` (`openai-chat`, `https://api.openai.com/v1`), `anthropic` (`anthropic`, `https://api.anthropic.com`), `openai-responses` (`openai-responses`, `https://api.openai.com/v1`), `gemini` (`gemini`, `https://generativelanguage.googleapis.com`), `openrouter` (`openai-chat`, `https://openrouter.ai/api/v1`), `openai-compatible` (`openai-chat`, custom base URL), `anthropic-compatible` (`anthropic`, custom base URL).
    - 3 core OAuth adapters (PKCE S256, headless-safe): `claude` (`anthropic`), `codex` (`openai-responses`), `gemini-oauth` (`gemini`). Client ID/secret read from env (`OMEG4MON_OAUTH_<PROVIDER>_CLIENT_ID` / `_SECRET`), per-request/connection override, or `Options.OAuthOverrides` for testing — never hardcoded secrets.
  - SQLite schema (`0001_init.sql`):
    - `provider_nodes (id TEXT PRIMARY KEY, name TEXT NOT NULL, protocol TEXT NOT NULL, base_url TEXT NOT NULL, models_json TEXT NOT NULL DEFAULT '[]', headers_json TEXT NOT NULL DEFAULT '{}', created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `provider_connections (id TEXT PRIMARY KEY, provider_id TEXT NOT NULL, name TEXT NOT NULL, protocol TEXT NOT NULL, base_url TEXT NOT NULL, auth_type TEXT NOT NULL, credentials_enc TEXT NOT NULL, priority INTEGER NOT NULL DEFAULT 100, enabled INTEGER NOT NULL DEFAULT 1, models_json TEXT NOT NULL DEFAULT '[]', headers_json TEXT NOT NULL DEFAULT '{}', token_expires_at TEXT, last_tested_at TEXT, last_test_status TEXT NOT NULL DEFAULT '', last_test_error TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `oauth_sessions (state TEXT PRIMARY KEY, provider_id TEXT NOT NULL, connection_name TEXT NOT NULL, code_verifier_enc TEXT NOT NULL, redirect_uri TEXT NOT NULL, client_id TEXT NOT NULL, client_secret_enc TEXT NOT NULL DEFAULT '', base_url TEXT NOT NULL DEFAULT '', models_json TEXT NOT NULL DEFAULT '[]', expires_at TEXT NOT NULL, created_at TEXT NOT NULL)`
  - Routes:
    - AdminAPI:
      - `GET /providers/catalog` — list registered built-in adapters + custom provider nodes.
      - `GET /providers/nodes`, `POST /providers/nodes`, `DELETE /providers/nodes/{id}` — custom compatible provider-nodes CRUD.
      - `GET /providers/connections`, `POST /providers/connections`, `GET /providers/connections/{id}`, `PATCH /providers/connections/{id}`, `DELETE /providers/connections/{id}` — connection CRUD (encrypts credentials via `Vault.Encrypt`, NEVER returns plaintext secrets, publishes `providers.changed`).
      - `POST /providers/connections/{id}/test` — runs upstream connectivity quick-test via adapter + `HTTPClient`, updates `last_tested_at`/`last_test_status`/`last_test_error`, returns `{ "ok": bool, "status_code": int, "latency_ms": int64, "error": string }`.
      - `POST /providers/oauth/{provider}/start` — generates PKCE S256 verifier/challenge + CSPRNG `state`, stores encrypted verifier in `oauth_sessions` (10m TTL), returns `{ "state", "authorize_url", "redirect_uri", "expires_at" }`.
      - `POST /providers/oauth/{provider}/exchange` — accepts `{ "state"?, "code"?, "callback_url"? }` (parses `code` + `state` cleanly from pasted callback URL OR direct `code` + `state`), exchanges authorization code at provider token endpoint, seals access+refresh tokens in `Vault`, creates/updates connection, deletes consumed `oauth_sessions` row, publishes `providers.changed`, returns `201` connection summary.
    - PublicAPI:
      - `GET /oauth/callback` — receives browser redirect `?code=...&state=...`, renders a minimal safe plain-HTML/JSON handoff page displaying the verified status or copyable callback parameters (no unescaped reflected XSS!).
  - Single-flight per-connection `RefreshOAuthToken(ctx, connectionID)`: coalesces concurrent refresh calls for the same `connectionID` using an in-flight map + mutex/channel so only one token-endpoint POST fires, updates vaulted tokens atomically, and returns fresh `UpstreamCredentials`.

- [ ] **Step 1: Write failing tests in `internal/modules/providers/module_test.go`**
  - Cover catalog, custom provider nodes, connection CRUD, secret encryption hygiene _(Review Focus #1)_, quick-test against `httptest.Server`, headless-safe OAuth start + callback-URL paste exchange against mock OAuth token server, and concurrent `RefreshOAuthToken` singleflight coalescing (asserting mock token endpoint receives exactly 1 refresh request across 8 concurrent callers!).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/providers/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `providers` module and `contracts/openapi/modules/providers.yaml`**
  - Implement migration, adapters, PKCE OAuth & singleflight refresh, repo, module, handlers, and OpenAPI fragment.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/providers/... && CGO_ENABLED=0 go test -count=1 ./internal/modules/providers/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(providers): implement pluggable ProviderAdapter registry, vaulted connections, and headless-safe OAuth`.

---

### Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/combos/migrations/0001_init.sql`
- Create: `internal/modules/combos/cooldown.go`
- Create: `internal/modules/combos/resolver.go`
- Create: `internal/modules/combos/repo.go`
- Create: `internal/modules/combos/module.go`
- Create: `internal/modules/combos/module_test.go`
- Create: `contracts/openapi/modules/combos.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `contracts.ProviderStore`, `*eventbus.Bus`
- Produces:
  - `combos.Options{DB *db.DB, Providers contracts.ProviderStore, Bus *eventbus.Bus, Now func() time.Time}`
  - `combos.New(opts Options) *Module` (`Name() = "combos"`, `Dependencies() = []string{"system", "providers"}`)
  - `(*Module)` implements `lifecycle.Module`, `contracts.ProviderResolver`, `contracts.AccountHealthTracker`, AND `contracts.ModelLister`.
  - SQLite schema (`0001_init.sql`):
    - `combos (id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE, description TEXT NOT NULL DEFAULT '', targets_json TEXT NOT NULL, strategy TEXT NOT NULL DEFAULT 'ordered', enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `model_aliases (alias TEXT PRIMARY KEY, target_model TEXT NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `disabled_models (model TEXT PRIMARY KEY, reason TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
  - Resolution algorithm (`resolver.go`, Review Focus #2):
    - Check `allowed` against the **requested identifier** (`model`): if `len(allowed) > 0` and `"*"` is not present and `model` does not match any allowed entry (exact match or `<provider>/*` prefix wildcard), return `contracts.ErrModelForbidden`.
    - Tier 1: Active Combo (`combos.name == model && enabled`) -> expand ordered `targets` (each target resolved via Alias / `<provider>/<model>` / bare `<model>`, preventing alias/combo cycles via visited set).
    - Tier 2: Model Alias (`model_aliases.alias == model`) -> resolve `target_model` (protecting against cycles).
    - Tier 3: Explicit `<provider>/<model>` (contains `/`) -> match enabled connections where `conn.ProviderID == provider` and model is supported (or connection allows wildcard `*` / empty list) and not in `disabled_models`.
    - Tier 4: Bare `<model>` -> match across all enabled connections supporting `<model>` ordered by `(conn.Priority ASC, conn.ID ASC)`.
    - Filter out connections currently in cooldown (`cooldown.IsInCooldown(conn.ID, now)`); if candidates existed but all are cooling down, return `contracts.ErrAllAccountsCoolingDown`; if zero candidates matched, return `contracts.ErrModelNotFound`.
  - Cooldown state machine (`cooldown.go`):
    - Thread-safe map of `connectionID -> cooldownEntry{Failures, Until, LastStatus, UpdatedAt}`.
    - `RecordFailure(id, statusCode, retryAfter)`: if `retryAfter > 0`, use `max(retryAfter, base)`; otherwise exponential backoff (`15s -> 30s -> 60s -> 120s -> 300s` capped at 5m, with `401`/`403` starting at `60s`).
    - `RecordSuccess(id)`: clears failure count and cooldown immediately.
    - `ListCooldowns(now)` & `ClearCooldown(id)` for admin observability/override.
  - AdminAPI routes:
    - `GET /combos`, `POST /combos`, `GET /combos/{id}`, `PATCH /combos/{id}`, `DELETE /combos/{id}`
    - `GET /aliases`, `PUT /aliases/{alias}`, `DELETE /aliases/{alias}`
    - `GET /models/disabled`, `PUT /models/disabled`, `DELETE /models/disabled` (body or query `?model=...` since model names can contain `/`!)
    - `GET /models/cooldowns`, `DELETE /models/cooldowns/{connectionId}`

- [ ] **Step 1: Write failing tests in `internal/modules/combos/module_test.go`**
  - Cover CRUD for combos, aliases, disabled models, cooldown inspect/clear, `ListModels` deduplication + `allowed_models` filtering, and Review Focus #2:
    - Exact 4-tier precedence (Combo beats Alias beats `<provider>/<model>` beats Bare model by priority).
    - `allowed_models`: authorizing combo name authorizes member targets; unauthorized model returns `contracts.ErrModelForbidden`; empty/`*` allows all.
    - Cooldown backoff progression, `Retry-After` override, `ErrAllAccountsCoolingDown`, and `RecordSuccess` reset (using injectable `Now` clock — zero sleeps!).
    - Cycle protection on self-referential aliases.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/combos/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `combos` module and `contracts/openapi/modules/combos.yaml`**
  - Implement migration, cooldown state machine, 4-tier resolver, repo, module, and OpenAPI fragment.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/combos/... && CGO_ENABLED=0 go test -count=1 ./internal/modules/combos/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(combos): implement fallback combos, aliases, 4-tier model resolver, and cooldown state machine`.

---

### Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators

**Files:**
- Create: `internal/modules/gateway/ir.go`
- Create: `internal/modules/gateway/translate_openai.go`
- Create: `internal/modules/gateway/translate_anthropic.go`
- Create: `internal/modules/gateway/translate_responses.go`
- Create: `internal/modules/gateway/translate_gemini.go`
- Create: `internal/modules/gateway/tokens.go`
- Create: `internal/modules/gateway/translate_test.go`

**Interfaces:**
- Produces inside `internal/modules/gateway`:
  - Protocols: `ProtocolOpenAIChat = "openai-chat"`, `ProtocolOpenAIResponses = "openai-responses"`, `ProtocolAnthropic = "anthropic"`, `ProtocolGemini = "gemini"`.
  - `CanonicalChatRequest`, `CanonicalMessage`, `ContentBlock` (`text`, `image`, `thinking`, `tool_use`, `tool_result`), `ToolDefinition`, `ResponseFormat`, `CanonicalChatResponse`, `CanonicalUsage`, `CanonicalStreamChunk`.
  - `CanFastPath(inProto, outProto string) bool` + `RewriteModelField(body []byte, proto, upstreamModel string) ([]byte, error)` (preserves all unknown provider-specific fields such as Anthropic `cache_control`, `thinking`, `betas` verbatim!).
  - `DecodeRequest(proto string, body []byte, modelOverride string) (*CanonicalChatRequest, error)`
  - `EncodeRequest(proto string, req *CanonicalChatRequest) ([]byte, error)`
  - `DecodeResponse(proto string, body []byte) (*CanonicalChatResponse, error)`
  - `EncodeResponse(proto string, resp *CanonicalChatResponse) ([]byte, error)`
  - `StreamTranslator` interface / stateful pair translators: converts upstream SSE lines/events (`outProto`) into inbound SSE frames (`inProto`), including direct Anthropic thinking/cache-control ↔ OpenAI Responses pair mapping.
  - `EstimateTokens(req *CanonicalChatRequest) int64` (zero-table deterministic heuristic for `/v1/messages/count_tokens` fallback).

- [ ] **Step 1: Write failing translator tests in `internal/modules/gateway/translate_test.go`**
  - Pin Review Focus #3:
    - Fast-path `RewriteModelField` preserves Anthropic `cache_control` and `thinking` blocks byte-for-byte while updating `model`.
    - Cross-protocol round-trips: OpenAI Chat -> Anthropic Messages (system + parallel tool_calls + image block + streaming deltas), Anthropic Messages -> OpenAI Chat (tool_use + tool_result + thinking block), OpenAI Responses <-> Anthropic Messages direct pair (instructions, function_call, reasoning), Gemini `generateContent` <-> OpenAI Chat.
    - `EstimateTokens` returns deterministic positive token counts for system, text, and tool payloads.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/gateway/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement Canonical IR, fast-path model rewriter, bidirectional translators, and token estimator**
  - Keep each protocol translator cleanly isolated in its own file.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/gateway/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(gateway): implement Canonical IR, direct fast-path, and multi-protocol translators`.

---

### Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain

**Files:**
- Create: `internal/modules/gateway/migrations/0001_init.sql`
- Create: `internal/modules/gateway/auth.go`
- Create: `internal/modules/gateway/classify.go`
- Create: `internal/modules/gateway/chains.go`
- Create: `internal/modules/gateway/stream.go`
- Create: `internal/modules/gateway/module.go`
- Create: `internal/modules/gateway/module_test.go`

**Interfaces:**
- Consumes: `*db.DB`, `*http.Client`, `contracts.KeyValidator`, `contracts.ProviderResolver`, `contracts.ModelLister`, `contracts.ProviderStore`, `contracts.AccountHealthTracker`, `[]contracts.GatewayInterceptor`, optional `contracts.UsageRecorder` (nil-safe), `*slog.Logger`
- Produces:
  - `gateway.Options{DB *db.DB, HTTPClient *http.Client, Keys contracts.KeyValidator, Resolver contracts.ProviderResolver, Models contracts.ModelLister, Providers contracts.ProviderStore, Health contracts.AccountHealthTracker, Interceptors []contracts.GatewayInterceptor, Usage contracts.UsageRecorder, Logger *slog.Logger, ChainTTL time.Duration}`
  - `gateway.New(opts Options) *Module` (`Name() = "gateway"`, `Dependencies() = []string{"system", "keys", "providers", "combos"}`)
  - SQLite migration (`0001_init.sql`): `response_chains (response_id TEXT PRIMARY KEY, key_id TEXT NOT NULL, model TEXT NOT NULL, messages_json TEXT NOT NULL, expires_at TEXT NOT NULL, created_at TEXT NOT NULL)` + index `idx_response_chains_expires ON response_chains(expires_at)`.
  - Mounted on `mounts.Gateway` (zero write timeout!):
    - `POST /v1/chat/completions`
    - `POST /v1/messages`, `POST /v1/messages/count_tokens`
    - `POST /v1/responses`, `POST /v1/responses/compact`
    - `GET /v1/models`, `GET /v1/models/{id}`
    - `GET /v1beta/models`, `POST /v1beta/models/{modelAction}` (`:generateContent`, `:streamGenerateContent`)
    - Compatibility aliases: `POST /codex/responses`, `POST /responses`, `GET /models`, `GET /models/{id}`.
  - Auth extraction (`auth.go`): `Authorization: Bearer <key>`, `x-api-key: <key>`, or `?key=<key>` (Gemini); returns protocol-appropriate `401` JSON error on invalid/disabled/expired key.
  - Execution pipeline (`stream.go`, `classify.go`, Review Focus #4):
    - Runs `GatewayInterceptor` chain sorted by `Priority() ASC` (`OnRequest`, `OnResponseHeaders`, `OnStreamChunk`).
    - Resolves candidates via `Resolver.Resolve(ctx, model, keyInfo.AllowedModels)`.
    - For each candidate:
      - Loads `UpstreamCredentials` via `Providers.ResolveCredentials`.
      - Dispatches upstream HTTP request using `r.Context()` and `HTTPClient`.
      - On `401` with `AuthType == "oauth"`: calls `Providers.RefreshOAuthToken` ONCE and retries the same candidate before classifying as failure.
      - Classifies non-2xx errors: client errors (`400`, `404`, `422`) return immediately without touching `Health.RecordFailure` or trying fallback; account/upstream errors (`401`/`403`/`429`/`5xx`/transport error) call `Health.RecordFailure` (parsing `Retry-After`) and continue to next candidate!
      - On `2xx`: calls `Health.RecordSuccess`, runs `OnResponseHeaders`, streams SSE line-by-line via `bufio.Reader` + `http.NewResponseController(w).Flush()` (or forwards JSON body for non-streaming), emits `UsageRecord` if `Usage != nil`.

- [ ] **Step 1: Write failing tests in `internal/modules/gateway/module_test.go`**
  - Pin Review Focus #4 & all protocol routes:
    - Auth extraction (`Bearer`, `x-api-key`, `?key=` query param) and `401`/`403` errors.
    - Direct fast-path SSE streaming (verifying chunks flush incrementally before upstream closes!) and cross-protocol streaming (OpenAI Chat -> mock Anthropic SSE upstream, Anthropic Messages -> mock OpenAI SSE upstream, Gemini `:streamGenerateContent` -> mock OpenAI upstream).
    - Error classification: upstream `400` returns immediately without calling fallback #2; upstream `429` (with `Retry-After`) or `500` records cooldown and fails over seamlessly to candidate #2; upstream `401` on OAuth triggers single `RefreshOAuthToken` retry and succeeds.
    - `GatewayInterceptor` priority ordering (`Priority() 10` runs before `Priority() 20`) on `OnRequest`, `OnResponseHeaders`, `OnStreamChunk`.
    - `/v1/models` filtered union, `/v1/messages/count_tokens` heuristic fallback, `/v1/responses` `previous_response_id` SQLite chain persistence, and compatibility aliases (`/responses`, `/codex/responses`, `/models`).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/gateway/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement gateway HTTP pipeline, auth, classifier, chains, streaming executor, and module**
  - Ensure all response bodies are unconditionally closed (`bodyclose` linter clean!) and `golangci-lint` passes with 0 issues.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/gateway/... && CGO_ENABLED=0 go test -count=1 ./internal/modules/gateway/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(gateway): implement multi-protocol SSE gateway, OAuth refresh failover, and interceptor chain`.

---

### Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Modify: `web/src/shared/api/schema.d.ts` (generated via `make codegen`)

**Interfaces:**
- Consumes: `providers.New`, `combos.New`, `gateway.New`, `keysMod` (`contracts.KeyValidator`)
- Produces:
  - Full Wave 1B registration (`system` -> `auth`, `settings`, `keys`, `providers` -> `combos` -> `gateway`) in `cmd/omeg4mon-router/main.go`.
  - Root OpenAPI document referencing all `providers.yaml` and `combos.yaml` paths & schemas, plus freshly generated `web/src/shared/api/schema.d.ts`.

- [ ] **Step 1: Write failing end-to-end integration test in `cmd/omeg4mon-router/main_test.go`**
  - Add `TestAppWave1BGatewayEndToEnd`:
    - Boot real `buildApp` with temp `DataDir` and `AdminPassword`.
    - Verify `GET /api/health` includes `"providers"`, `"combos"`, `"gateway"` alongside Wave 1A modules.
    - Authenticate admin -> create a provider connection pointing at a local mock OpenAI SSE `httptest.Server` + a fallback mock Anthropic `httptest.Server` -> create a combo spanning both -> create a gateway API key scoped to that combo.
    - Send `POST /v1/chat/completions` (`stream: true`) with `Authorization: Bearer sk-om4-...` -> verify failover from primary `429` to secondary mock Anthropic upstream translated back into OpenAI SSE chunks + `[DONE]`!
    - Verify `GET /v1/models` with the scoped API key returns only the authorized combo.
- [ ] **Step 2: Run test to verify failure**
  - Run: `go test -race -run TestAppWave1BGatewayEndToEnd ./cmd/omeg4mon-router/...`
  - Expected: FAIL.
- [ ] **Step 3: Wire modules in `cmd/omeg4mon-router/main.go`, update `contracts/openapi/openapi.yaml`, run `make codegen`**
  - Wire `providersMod`, `combosMod`, `gatewayMod` cleanly with zero cross-module imports.
  - Aggregate all `providers` & `combos` tags/paths/schemas in `contracts/openapi/openapi.yaml` and run `make codegen`.
- [ ] **Step 4: Run full `make check`**
  - Run: `make check`
  - Expected: PASS (0 lint/format issues, all Go & Web tests green).
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(cmd): wire providers, combos, and gateway modules with OpenAPI codegen and E2E tests`.

---

### Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`)

**Files:**
- Create: `web/src/features/providers/schemas.ts`
- Create: `web/src/features/providers/schemas.test.ts`
- Create: `web/src/features/providers/hooks.ts`
- Create: `web/src/features/providers/ProvidersConsole.tsx`
- Create: `web/src/features/providers/ProviderDetailView.tsx`
- Create: `web/src/app/(dashboard)/dashboard/providers/page.tsx`
- Create: `web/src/app/(dashboard)/dashboard/providers/detail/page.tsx`

**Interfaces:**
- Consumes: typed `apiClient` (`/api/providers/*`), `react-hook-form` + `zodResolver`, `shadcn/ui`, `useSearchParams` (wrapped in `<Suspense>` for Next.js 16 static export on `/dashboard/providers/detail/`!).
- Produces:
  - `/dashboard/providers/` — Catalog overview, Connected Accounts table (status badge, protocol, priority, last test latency/result, quick-test button, link to `/dashboard/providers/detail/?id=<id>`, enable/disable, delete), API-Key Connection creation form, Custom Compatible Provider Node form, and **Headless-Safe OAuth Wizard** (select `claude` / `codex` / `gemini-oauth` -> Start -> displays Authorize URL + copy/open button -> paste callback URL or auth code -> Exchange) _(Review Focus #5)_.
  - `/dashboard/providers/detail/?id=<id>` — Static-export-safe query-param detail view allowing priority, model list, base URL, headers, credential rotation, and live connectivity test.

- [ ] **Step 1: Write failing tests in `web/src/features/providers/schemas.test.ts`**
  - Validate `createConnectionSchema` (URL check, non-empty API key, priority bounds, CSV model parser), `createNodeSchema`, and `oauthExchangeSchema` (requires either `callback_url` or both `code` & `state`, rejects empty submission).
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/providers/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `web/src/features/providers/*` and static-export route pages**
  - Wrap `useSearchParams()` inside `<Suspense>` in `detail/page.tsx` so `NEXT_OUTPUT_EXPORT=true bun run build` succeeds cleanly. Strictly follow `docs/DESIGN.md` (no `backdrop-blur`, no emojis, Geist fonts, `#E56A4A` accent).
- [ ] **Step 4: Run Web checks & static export build**
  - Run: `cd web && bun run check && bun test && bun run typecheck && NEXT_OUTPUT_EXPORT=true bun run build`
  - Expected: PASS, `web/out/dashboard/providers/index.html` and `web/out/dashboard/providers/detail/index.html` generated.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Providers Control Studio, headless-safe OAuth wizard, and query-param detail page`.

---

### Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

**Files:**
- Create: `web/src/features/combos/schemas.ts`
- Create: `web/src/features/combos/schemas.test.ts`
- Create: `web/src/features/combos/hooks.ts`
- Create: `web/src/features/combos/ModelsRoutingConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/models/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx` (promote `Providers` -> `/dashboard/providers/` and `Models & Routing` -> `/dashboard/models/` into `ACTIVE_NAV_ITEMS`)
- Modify: `web/src/features/system/SystemOverview.tsx` (reflect active Wave 1A + 1B modules: `system`, `auth`, `settings`, `keys`, `providers`, `combos`, `gateway`; remaining Wave 1C: `usage`, `clitools`)

**Interfaces:**
- Consumes: typed `apiClient` (`/api/combos*`, `/api/aliases*`, `/api/models/disabled`, `/api/models/cooldowns`, `/api/providers/connections`), `react-hook-form` + `zodResolver`, `shadcn/ui`.
- Produces:
  - `/dashboard/models/` — Ordered Fallback Combos builder (add/reorder/remove targets, toggle enabled, delete), Model Aliases editor, Disabled Models panel, Account Cooldowns live table with one-click Clear Cooldown, and Gateway Endpoint Reference card (OpenAI `/v1`, Anthropic `/v1`, Gemini `/v1beta` base URLs + copyable curl snippet).

- [ ] **Step 1: Write failing tests in `web/src/features/combos/schemas.test.ts`**
  - Validate `createComboSchema` (non-empty name, at least 1 unique trimmed target), `upsertAliasSchema` (rejects self-referential `alias === target_model`), and `disableModelSchema`.
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/combos/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `ModelsRoutingConsole`, `/dashboard/models/page.tsx`, and update `AppShell.tsx` & `SystemOverview.tsx`**
  - Strictly follow `docs/DESIGN.md` (solid opaque surfaces, zero `backdrop-blur`, zero emojis, WCAG AA contrast).
- [ ] **Step 4: Run full verification suite (`make check` + `make build` + live single-binary smoke test)**
  - Run: `make check && make build`
  - Verify `web/out/dashboard/providers/index.html`, `web/out/dashboard/providers/detail/index.html`, and `web/out/dashboard/models/index.html` exist.
  - Smoke-test `bin/omeg4mon-router` on temp port/dir: `/api/health` includes all 7 active modules (`auth`, `combos`, `gateway`, `keys`, `providers`, `settings`, `system`), and HTTP 200 HTML on `/dashboard/providers/`, `/dashboard/providers/detail/`, `/dashboard/models/`.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Models & Routing studio (combos, aliases, cooldowns) and Wave 1B navigation`.
