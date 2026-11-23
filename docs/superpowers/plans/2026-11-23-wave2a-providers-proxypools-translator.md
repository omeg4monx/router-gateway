# Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 2A by expanding the built-in `providers` adapter catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), adding encrypted HTTP/HTTPS/SOCKS5 outbound proxy pools with latency testing, auto-scanning, and per-provider binding (`proxypools`), and adding an interactive protocol translator workbench plus live console log SSE stream (`translator`) — all wired into the Go binary, OpenAPI 3.1 contract, and Next.js 16 dashboard.

**Architecture:** `providers` gains five dedicated `ProviderAdapter` entries without modifying gateway core. `internal/modules/proxypools` owns its SQLite tables (`proxy_pools`, `proxy_bindings`), encrypts proxy URLs at rest via `crypto.Vault`, and plugs into `httpclient` via `contracts.ProxyDialerResolver` + request-context provider tagging. `internal/modules/translator` consumes `contracts.ProtocolTranslator` (implemented by `gateway.Module`) and `logging.RingBuffer.Subscribe()` to expose interactive request/response translation and live structured log streaming with zero cross-module imports.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `crypto.Vault`, `eventbus.Bus`, `logging.RingBuffer`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339 timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded in-memory caches and bounded non-blocking log subscriber channels.
- **Secret protection at rest and on wire:** Proxy URLs (which may contain `user:pass@host:port` credentials) are encrypted at rest via `*crypto.Vault` and always returned with userinfo credentials redacted (`scheme://***@host:port` or `scheme://host:port` + `has_credentials: bool`).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task must pass `make check` (`golangci-lint` v2 + Biome v2 + Go `-race` + zero-CGO DB + Bun tests + `tsc --noEmit`) and `graphify update .` (cleaning any `graphify-out/20??-??-??/` backup directories).

## Review Focus

1. **Proxy credential redaction & scheme validation in `proxypools`:** Only `http`, `https`, `socks5`, and `socks5h` schemes with non-empty host are accepted; proxy URLs containing username/password are encrypted at rest via `crypto.Vault`, never logged or returned in plaintext in any API response, and still work transparently when returned by `ProxyURLForProvider`. _(Pinned in Task 2)_
2. **Non-blocking `RingBuffer.Subscribe()` under slow SSE clients:** A slow or stalled SSE log consumer MUST NEVER block `slog.Logger` / `RingBuffer.add()` or leak goroutines when the HTTP client disconnects (`r.Context().Done()`). _(Pinned in Task 3)_
3. **Zero-coupling `ProtocolTranslator` fidelity:** Translating requests and responses across `openai-chat`, `anthropic`, `openai-responses`, and `gemini` via `POST /api/translator/translate` executes the exact `gateway` translator pipeline through `contracts.ProtocolTranslator`, rejecting unknown protocols or malformed JSON with `400`. _(Pinned in Task 3)_
4. **Live Gateway -> Proxy Pool E2E routing:** Binding a provider to a local test HTTP proxy pool in `buildApp` routes a real `/v1/chat/completions` upstream request through that proxy, while unbound providers continue dialing directly. _(Pinned in Task 4)_
5. **Static-export UI completeness & zero-blur `DESIGN.md` compliance:** `/dashboard/proxy-pools/` and `/dashboard/translator/` export as static trailing-slash HTML routes, serve HTTP 200 from `bin/omeg4mon-router`, and contain zero `backdrop-blur`/`backdrop-filter` or UI emojis. _(Pinned in Tasks 5 & 6)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Provider expansion | `internal/modules/providers/adapter.go`, `internal/modules/providers/module_test.go` | Register 5 dedicated Wave 2 adapters (`qwen`, `github-copilot`, `kiro`, `cursor`, `antigravity`) in `defaultBuiltInAdapters()` |
| HTTP client proxy hook | `internal/platform/httpclient/client.go`, `internal/platform/httpclient/client_test.go` | Context-scoped provider ID (`WithProviderID` / `ProviderIDFromContext`) and `NewProviderProxyFunc(resolver)` |
| Proxy pools module | `internal/modules/proxypools/{migrations/0001_init.sql,repo.go,scanner.go,module.go,module_test.go}`, `contracts/openapi/modules/proxypools.yaml` | Encrypted HTTP/SOCKS5 proxy pools CRUD, latency tester, concurrent auto-scanner, provider binding, `contracts.ProxyDialerResolver` |
| Log streaming & translator | `internal/platform/logging/{logger.go,logger_test.go}`, `internal/contracts/contracts.go`, `internal/modules/gateway/module.go`, `internal/modules/translator/{module.go,module_test.go}`, `contracts/openapi/modules/translator.yaml` | Non-blocking `RingBuffer.Subscribe()`, `contracts.ProtocolTranslator`, interactive translator workbench + live SSE console log stream |
| App wiring & E2E | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Wire `proxypools` + `translator`, dynamic proxy resolver on `httpClient`, OpenAPI codegen, E2E verification |
| Web Proxy Pools UI | `web/src/features/proxypools/{schemas.ts,schemas.test.ts,hooks.ts,ProxyPoolsConsole.tsx}`, `web/src/app/(dashboard)/dashboard/proxy-pools/page.tsx` | Proxy pool management, latency test, auto-scanner, and provider binding studio |
| Web Translator UI & Nav | `web/src/features/translator/{schemas.ts,schemas.test.ts,hooks.ts,TranslatorConsole.tsx}`, `web/src/app/(dashboard)/dashboard/translator/page.tsx`, `web/src/shared/ui/AppShell.tsx`, `web/src/features/system/SystemOverview.tsx` | Interactive protocol translator workbench, live SSE console log viewer, Wave 2A nav + single-binary smoke test |

---

### Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity)

**Files:**
- Modify: `internal/modules/providers/adapter.go`
- Modify: `internal/modules/providers/module_test.go`

**Interfaces:**
- Consumes: `providers.ProviderAdapter`, `providers.OAuthConfig`, `contracts.UpstreamCredentials`
- Produces: 5 additional built-in adapters registered in `defaultBuiltInAdapters()` (total 15 built-in adapters):
  - `qwen` (`Qwen (DashScope / OAuth)`, category `oauth`, protocol `openai-chat`, auth `oauth`, base `https://portal.qwen.ai/v1`, models `["qwen3-coder-plus", "qwen3-235b-a22b", "qwen-max"]`, OAuth authorize/token URLs).
  - `github-copilot` (`GitHub Copilot (OAuth)`, category `oauth`, protocol `openai-chat`, auth `oauth`, base `https://api.githubcopilot.com`, models `["gpt-4o", "claude-sonnet-4", "o3-mini"]`, OAuth authorize/token URLs).
  - `kiro` (`Kiro (OAuth)`, category `oauth`, protocol `anthropic`, auth `oauth`, base `https://codewhisperer.us-east-1.amazonaws.com`, models `["claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219"]`, OAuth authorize/token URLs).
  - `cursor` (`Cursor`, category `api_key`, protocol `openai-chat`, auth `api_key`, base `https://api2.cursor.sh/v1`, models `["claude-sonnet-4", "gpt-4.1", "gemini-2.5-pro"]`).
  - `antigravity` (`Antigravity (OAuth)`, category `oauth`, protocol `gemini`, auth `oauth`, base `https://cloudcode-pa.googleapis.com`, models `["gemini-2.5-pro", "gemini-2.5-flash", "claude-sonnet-4"]`, OAuth authorize/token URLs with offline consent params).

- [ ] **Step 1: Write failing test in `internal/modules/providers/module_test.go`**
  - Add `TestWave2ProviderAdaptersRegisteredAndTestable`: verify `GET /api/providers/catalog` includes `qwen`, `github-copilot`, `kiro`, `cursor`, and `antigravity` with expected protocol, auth type, non-empty default models, and non-nil `OAuthConfig` for OAuth adapters; verify `BuildTestRequest` produces valid HTTP requests with proper auth headers for each.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race -run TestWave2ProviderAdaptersRegisteredAndTestable ./internal/modules/providers/...`
  - Expected: FAIL (adapters not found in catalog).
- [ ] **Step 3: Add the 5 dedicated adapters to `defaultBuiltInAdapters()` in `internal/modules/providers/adapter.go`**
  - Keep all URLs/configs clean-room and `golangci-lint` (`gosec` G101 annotated for public OAuth URLs) compliant.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/providers/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(providers): add dedicated Wave 2 adapters for Qwen, GitHub Copilot, Kiro, Cursor, and Antigravity`.

---

### Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding)

**Files:**
- Modify: `internal/platform/httpclient/client.go`
- Modify: `internal/platform/httpclient/client_test.go`
- Create: `internal/modules/proxypools/migrations/0001_init.sql`
- Create: `internal/modules/proxypools/repo.go`
- Create: `internal/modules/proxypools/scanner.go`
- Create: `internal/modules/proxypools/module.go`
- Create: `internal/modules/proxypools/module_test.go`
- Create: `contracts/openapi/modules/proxypools.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*crypto.Vault`, `*eventbus.Bus`, `*slog.Logger`, `contracts.ProxyDialerResolver`
- Produces:
  - `httpclient.WithProviderID(ctx context.Context, providerID string) context.Context`
  - `httpclient.ProviderIDFromContext(ctx context.Context) string`
  - `httpclient.NewProviderProxyFunc(resolver contracts.ProxyDialerResolver) func(*http.Request) (*url.URL, error)` (when `providerID` is in `req.Context()` and `resolver` returns non-nil `*url.URL`, returns that proxy URL; otherwise falls back to `http.ProxyFromEnvironment(req)`).
  - `proxypools.Options{DB *db.DB, Vault *crypto.Vault, Bus *eventbus.Bus, Logger *slog.Logger, ProbeTimeout time.Duration}`
  - `proxypools.New(opts Options) *Module` (`Name() = "proxypools"`, `Dependencies() = []string{"system", "providers"}`)
  - `(*Module)` implements `lifecycle.Module` AND `contracts.ProxyDialerResolver` (`var _ contracts.ProxyDialerResolver = (*Module)(nil)`).
  - SQLite tables (`0001_init.sql`):
    - `proxy_pools(id TEXT PRIMARY KEY, name TEXT NOT NULL, scheme TEXT NOT NULL, url_encrypted TEXT NOT NULL, display_url TEXT NOT NULL, has_credentials INTEGER NOT NULL DEFAULT 0, enabled INTEGER NOT NULL DEFAULT 1, status TEXT NOT NULL DEFAULT 'unknown', last_latency_ms INTEGER NOT NULL DEFAULT 0, last_error TEXT NOT NULL DEFAULT '', last_tested_at TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `proxy_bindings(provider_id TEXT PRIMARY KEY, pool_id TEXT NOT NULL REFERENCES proxy_pools(id) ON DELETE CASCADE, updated_at TEXT NOT NULL)`
  - AdminAPI routes (under `/api`):
    - `GET /proxypools` -> `{ "pools": [...], "bindings": [...] }` (all URLs redacted to `display_url`, never plaintext credentials).
    - `POST /proxypools` -> create proxy pool (`name`, `url`, `enabled?`). Validates `http`/`https`/`socks5`/`socks5h` scheme and non-empty host.
    - `PATCH /proxypools/{id}` -> update `name`, `url` (optional), `enabled`.
    - `DELETE /proxypools/{id}` -> delete pool (cascades bindings).
    - `POST /proxypools/{id}/test` -> dials proxy (TCP handshake + optional `target_url` HTTP probe when provided) within `ProbeTimeout` (default 5s), records `status` (`healthy`/`unhealthy`), `last_latency_ms`, `last_error`, `last_tested_at`, and returns test result.
    - `POST /proxypools/scan` -> accepts `{ "candidates": ["http://...", ...], "auto_add"?: bool, "name_prefix"?: string, "timeout_ms"?: int }` (capped at 50 candidates, concurrency 8), probes each candidate concurrently, optionally saves healthy ones, returns `{ "results": [...] }`.
    - `PUT /proxypools/bindings/{providerId}` -> `{ "pool_id": "..." }` binds provider to pool.
    - `DELETE /proxypools/bindings/{providerId}` -> unbinds provider (returns to direct dial).

- [ ] **Step 1: Write failing tests in `internal/platform/httpclient/client_test.go` and `internal/modules/proxypools/module_test.go`**
  - Pin Review Focus #1:
    - `TestProviderProxyFuncContextRouting` in `httpclient/client_test.go`: verifies provider-tagged requests resolve via `ProxyDialerResolver` and untagged requests fall back cleanly.
    - `TestProxyPoolsCRUDCredentialRedactionAndValidation`: creates HTTP and SOCKS5 pools with `user:pass@` credentials, verifies API responses redact credentials (`display_url` has no password, `has_credentials: true`), verifies invalid schemes (`ftp://`, `file://`, missing host) return `400`, and verifies `ProxyURLForProvider` returns the decrypted full `*url.URL` for bound enabled pools and `nil` for unbound or disabled pools.
    - `TestProxyPoolsLatencyTestAndAutoScanner`: starts local `httptest.NewServer` / TCP listener, runs `POST /api/proxypools/{id}/test` and `POST /api/proxypools/scan` with healthy + unreachable candidates, verifies latency > 0 for healthy, error status for unreachable, and `auto_add` persistence.
    - `TestProxyPoolsBindingLifecycleAndCacheInvalidation`: bind, resolve, update/disable pool, unbind, delete pool cascade.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/platform/httpclient/... ./internal/modules/proxypools/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `httpclient` context proxy hook, `proxypools` module, and `contracts/openapi/modules/proxypools.yaml`**
  - Encrypt `url_encrypted` with `vault.EncryptString`, keep in-memory cache of active bindings + pools protected by `sync.RWMutex`, and invalidate on every write.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/platform/httpclient/... ./internal/modules/proxypools/... && CGO_ENABLED=0 go test -count=1 ./internal/modules/proxypools/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(proxypools): implement encrypted HTTP/SOCKS5 proxy pools, latency tester, scanner, and provider binding`.

---

### Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module

**Files:**
- Modify: `internal/platform/logging/logger.go`
- Modify: `internal/platform/logging/logger_test.go`
- Modify: `internal/contracts/contracts.go`
- Modify: `internal/modules/gateway/module.go`
- Create: `internal/modules/translator/module.go`
- Create: `internal/modules/translator/module_test.go`
- Create: `contracts/openapi/modules/translator.yaml`

**Interfaces:**
- Consumes: `contracts.ProtocolTranslator`, `*logging.RingBuffer`
- Produces:
  - `(*logging.RingBuffer).Subscribe() (<-chan LogEntry, func())` (buffered channel cap 64; `add()` uses non-blocking send so slow subscribers drop excess entries without blocking logging; unsubscribe closes channel safely once).
  - `contracts.ProtocolTranslator`:
    ```go
    type ProtocolTranslator interface {
        SupportedProtocols() []string
        TranslateRequest(fromProto, toProto string, body []byte) ([]byte, error)
        TranslateResponse(fromProto, toProto string, body []byte) ([]byte, error)
    }
    ```
  - `(*gateway.Module)` implements `contracts.ProtocolTranslator` (`var _ contracts.ProtocolTranslator = (*Module)(nil)`) AND attaches `httpclient.WithProviderID(ctx, cand.ProviderID)` on upstream HTTP requests in `executeCandidateAttempt`.
  - `translator.Options{Translator contracts.ProtocolTranslator, LogRing *logging.RingBuffer}`
  - `translator.New(opts Options) *Module` (`Name() = "translator"`, `Dependencies() = []string{"system", "gateway"}`, `Migrations() = nil`).
  - AdminAPI routes (under `/api`):
    - `GET /translator/protocols` -> `{ "protocols": ["openai-chat", "anthropic", "openai-responses", "gemini"], "modes": ["request", "response"] }`.
    - `POST /translator/translate` -> `{ "mode": "request"|"response", "from_protocol": "...", "to_protocol": "...", "payload": <json> }` -> `{ "mode", "from_protocol", "to_protocol", "fast_path": bool, "duration_us": int64, "output": <json> }`.
    - `GET /translator/logs` (`?limit=100&level=`) -> `{ "entries": [...], "total": int }`.
    - `GET /translator/logs/stream` (`?level=&tail=50`) -> SSE stream replaying up to `tail` recent entries then streaming live `event: log\ndata: {...}\n\n` frames via `http.NewResponseController(w).Flush()` with 15s keepalives.

- [ ] **Step 1: Write failing tests in `internal/platform/logging/logger_test.go` and `internal/modules/translator/module_test.go`**
  - Pin Review Focus #2 & #3:
    - `TestRingBufferSubscribeNonBlockingAndUnsubscribe` in `logging/logger_test.go`: fill beyond subscriber buffer capacity without reading, assert `logger.Info` never blocks, verify clean unsubscribe under `-race`.
    - `TestTranslatorTranslateRequestAndResponseAllPairs` in `translator/module_test.go`: test OpenAI Chat ↔ Anthropic, OpenAI Chat ↔ Gemini, Anthropic ↔ OpenAI Responses (both `request` and `response` modes), plus `400` on invalid protocol, invalid mode, or malformed JSON.
    - `TestTranslatorLogsSnapshotAndLiveSSEStream`: verify `GET /api/translator/logs?level=WARN` filtering and `GET /api/translator/logs/stream` receiving both replayed and newly emitted log entries, then terminating cleanly on context cancel.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/platform/logging/... ./internal/modules/translator/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `RingBuffer.Subscribe`, `contracts.ProtocolTranslator`, `gateway` implementation + provider-context tagging, `translator` module, and `contracts/openapi/modules/translator.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/platform/logging/... ./internal/modules/gateway/... ./internal/modules/translator/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(translator): implement protocol translator workbench, live log SSE stream, and gateway provider proxy context`.

---

### Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

**Interfaces:**
- Consumes: `proxypools.New`, `translator.New`, `httpclient.NewProviderProxyFunc`, `contracts.ProxyDialerResolver`
- Produces: Fully wired `buildApp` with dynamic provider proxy resolution (honoring `OMEG4MON_DISABLED_MODULES` — if `proxypools` is disabled, proxy resolver falls back to direct/env dial), updated OpenAPI root spec + generated TS schema, and E2E tests.

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go`**
  - Pin Review Focus #4:
    - `TestWave2AE2EProviderProxyRouting`: spins up an upstream mock LLM server and a local HTTP proxy server (`httptest.NewServer` recording forwarded requests); creates a provider connection + API key + proxy pool pointing at the local proxy + binds the provider; sends `POST /v1/chat/completions` and asserts the request traversed the proxy server and returned `200`.
    - `TestWave2AE2ETranslatorAndLogsStream`: verifies `POST /api/translator/translate` (OpenAI Chat -> Anthropic) and `GET /api/translator/logs` work end-to-end through `buildApp` with session auth.
    - `TestWave2ACatalogIncludesExpandedProviders`: verifies `/api/providers/catalog` exposes the 15 adapters including `qwen`, `github-copilot`, `kiro`, `cursor`, `antigravity`.
- [ ] **Step 2: Run E2E tests to verify failure**
  - Run: `go test -race -run TestWave2A ./cmd/omeg4mon-router/...`
  - Expected: FAIL.
- [ ] **Step 3: Wire `proxypools` & `translator` in `cmd/omeg4mon-router/main.go`, update `contracts/openapi/openapi.yaml`, and run `make codegen`**
  - Use a thread-safe `atomic.Pointer` / deferred resolver wrapper passed to `httpclient.NewProviderProxyFunc` so `httpClient` is constructed once and delegates to `proxypoolsMod` whenever active.
- [ ] **Step 4: Run full Go tests, linter & TS typecheck**
  - Run: `make codegen && make check`
  - Expected: PASS, 0 lint/type errors.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(cmd): wire proxypools and translator modules with OpenAPI codegen and E2E tests`.

---

### Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)

**Files:**
- Create: `web/src/features/proxypools/schemas.ts`
- Create: `web/src/features/proxypools/schemas.test.ts`
- Create: `web/src/features/proxypools/hooks.ts`
- Create: `web/src/features/proxypools/ProxyPoolsConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/proxy-pools/page.tsx`

**Interfaces:**
- Consumes: `web/src/shared/api/client.ts` (typed via `schema.d.ts`), `web/src/features/providers/hooks.ts` (`useProvidersCatalog` for provider binding dropdown), `web/src/shared/ui/*`
- Produces: Static-export route `/dashboard/proxy-pools/` with:
  - Proxy Pool table (name, scheme badge, redacted `display_url`, credentials indicator, enabled toggle, latency ms, health badge, `Test` & `Delete` actions).
  - Add / Edit Proxy Pool form (`react-hook-form` + `zod` validating `http/https/socks5/socks5h` URLs).
  - Auto-Scanner panel (multiline candidate URL input, `auto_add` checkbox, live scan results table with latency & status).
  - Provider Binding matrix (bind/unbind any catalog provider to a proxy pool).

- [ ] **Step 1: Write failing Zod schema unit tests in `web/src/features/proxypools/schemas.test.ts`**
  - Verify valid `http://`, `https://`, `socks5://`, `socks5h://` URLs pass; invalid schemes (`ftp://`), empty names, and >50 scanner candidates fail with clear messages.
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/proxypools/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `schemas.ts`, `hooks.ts`, `ProxyPoolsConsole.tsx`, and `page.tsx`**
  - Strict `docs/DESIGN.md` compliance: solid opaque surfaces (`var(--bg-surface)`, `var(--bg-surface-raised)`), `#E56A4A` accent, Geist Sans/Mono, zero `backdrop-blur`/`backdrop-filter`, zero emojis.
- [ ] **Step 4: Run web tests, Biome lint & typecheck**
  - Run: `cd web && bun test && bun run check && bun run typecheck`
  - Expected: PASS, 0 errors.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Proxy Pools studio with latency tester, auto-scanner, and provider bindings`.

---

### Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification

**Files:**
- Create: `web/src/features/translator/schemas.ts`
- Create: `web/src/features/translator/schemas.test.ts`
- Create: `web/src/features/translator/hooks.ts`
- Create: `web/src/features/translator/TranslatorConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/translator/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

**Interfaces:**
- Consumes: `web/src/shared/api/client.ts`, `/api/translator/protocols`, `/api/translator/translate`, `/api/translator/logs`, `/api/translator/logs/stream` (EventSource SSE)
- Produces: Static-export route `/dashboard/translator/` with:
  - Interactive Protocol Translator Workbench (mode: `request`/`response`, source & target protocol selectors, built-in sample payload presets for OpenAI Chat / Anthropic / Responses / Gemini, formatted JSON input/output editors, fast-path badge + microsecond timing, copy output).
  - Live Console Log Stream panel (EventSource connection to `/api/translator/logs/stream` with level filter `ALL`/`DEBUG`/`INFO`/`WARN`/`ERROR`, Pause/Resume, Clear, structured monospace log table).
  - Updated `AppShell.tsx` navigation adding `Proxy Pools` (`/dashboard/proxy-pools/`) and `Translator` (`/dashboard/translator/`), and `WAVE 2A` kernel badge.

- [ ] **Step 1: Write failing Zod schema unit tests in `web/src/features/translator/schemas.test.ts`**
  - Verify valid translate requests pass, invalid JSON payload string fails validation cleanly, and supported protocols/modes are enforced.
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/translator/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `schemas.ts`, `hooks.ts`, `TranslatorConsole.tsx`, `page.tsx`, and update `AppShell.tsx` & `SystemOverview.tsx`**
  - Pin Review Focus #5: zero `backdrop-blur`/`backdrop-filter`, zero UI emojis, clean EventSource cleanup on unmount.
- [ ] **Step 4: Run full verification (`make check` + `make build` + single-binary smoke test)**
  - Run: `make check && make build`
  - Smoke test `bin/omeg4mon-router` on an ephemeral localhost port with a temp data dir: verify HTTP 200 on `/api/health`, `/dashboard/proxy-pools/`, and `/dashboard/translator/`, and verify zero `backdrop-blur` in `web/src/features/{proxypools,translator}/`.
  - Expected: PASS, all checks green.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` (clean any `graphify-out/20??-??-??/` backups) and commit `feat(web): implement Protocol Translator workbench, live log SSE console, and Wave 2A navigation`.
