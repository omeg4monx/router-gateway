# Wave 2B — Gateway Interceptors (`promptinjection`, `modeleditor` / `modelbypass`, `tokensaver` / `headroom`) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 2B by plugging three isolated vertical-slice modules into `gateway`'s ordered `GatewayInterceptor` chain — `promptinjection` (configurable prompt guard middleware), `modeleditor` (custom model parameter overrides and conditional bypass rules), and `tokensaver` (Go-native token compression and Headroom cache-aware prefix affinity) — complete with SQLite migrations, OpenAPI 3.1 contracts, Next.js 16 dashboard studios, E2E verification, and single-binary embedding.

**Architecture:** `gateway` passes `CanonicalRequest.Metadata` through to `ResponseMeta.Metadata` so interceptors share per-request context cleanly from `OnRequest` to `OnResponseHeaders` with zero cross-module imports. The three modules execute in deterministic priority order: `promptinjection` (`Priority=10`) -> `modeleditor` (`Priority=20`) -> `tokensaver` (`Priority=30`). Each module owns its SQLite schema, bounded in-memory caches, and `/api/*` admin routes, and can be independently disabled via `OMEG4MON_DISABLED_MODULES`.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `eventbus.Bus`, RE2 `regexp`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339 timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded in-memory caches (Headroom LRU capped at `max_cache_entries`, Prompt Guard events bounded/pruned).
- **Safe Pattern Matching:** Custom regex rules in `promptinjection` use Go's linear-time RE2 engine, validated at create/update time with max length 512 bytes (zero ReDoS risk).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task must pass `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and `graphify update .` (cleaning any `graphify-out/20??-??-??/` backup directories). Final verification runs full `make check` and `make build`.

## Review Focus

1. **Prompt Guard fails closed on `block` and never leaks raw secrets in audit snippets:** When `promptinjection` blocks a request, upstream is never dialed, the configured HTTP status (400/403) is returned, and audit event snippets are truncated (<= 200 chars) and redacted. _(Pinned in Tasks 1 & 4)_
2. **Model Editor / Bypass preserves valid JSON across all 4 wire protocols:** Rewriting model names, injecting/stripping parameters, and prefixing/suffixing system prompts works cleanly on OpenAI Chat, Anthropic Messages, OpenAI Responses, and Gemini payloads without corrupting JSON structure or dropping unknown provider fields. _(Pinned in Task 2)_
3. **Token Saver respects `bypass_tokensaver` and never corrupts fenced code blocks:** `whitespace`/`dedupe` compression preserves fenced code blocks (`` ``` ``) verbatim, and `bypass_tokensaver`/`bypass_headroom` flags set by `modeleditor` skip compression/affinity cleanly. _(Pinned in Tasks 3 & 4)_
4. **Headroom LRU cache is strictly bounded and race-free under concurrent traffic:** Concurrent `OnRequest` + `OnResponseHeaders` calls never exceed `max_cache_entries`, evict stale/LRU entries cleanly, and pass `go test -race`. _(Pinned in Task 3)_
5. **Disabled optional modules never intercept live gateway traffic & static-export routes serve 200:** Disabling `promptinjection`, `modeleditor`, or `tokensaver` via `OMEG4MON_DISABLED_MODULES` omits their interceptor from `gateway`, and `/dashboard/prompt-guard/`, `/dashboard/model-editor/`, and `/dashboard/token-saver/` export as static HTML with zero `backdrop-blur` or UI emojis. _(Pinned in Tasks 4, 5 & 6)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Gateway metadata pass-through | `internal/modules/gateway/module.go`, `internal/modules/gateway/module_test.go` | Seed `ResponseMeta.Metadata` from `CanonicalRequest.Metadata` so interceptors share per-request state between `OnRequest` and `OnResponseHeaders` |
| Prompt Injection Guard | `internal/modules/promptinjection/{migrations/0001_init.sql,repo.go,engine.go,module.go,module_test.go}`, `contracts/openapi/modules/promptinjection.yaml` | Prompt guard config, built-in + custom RE2/substring rules, block/sanitize/audit interceptor (`Priority=10`), dry-run tester, bounded event log |
| Model Editor & Bypass | `internal/modules/modeleditor/{migrations/0001_init.sql,repo.go,engine.go,module.go,module_test.go}`, `contracts/openapi/modules/modeleditor.yaml` | Custom model parameter overrides + conditional bypass rules interceptor (`Priority=20`), multi-protocol JSON mutator, dry-run preview |
| Token Saver & Headroom | `internal/modules/tokensaver/{migrations/0001_init.sql,repo.go,compressor.go,headroom.go,module.go,module_test.go}`, `contracts/openapi/modules/tokensaver.yaml` | Go-native token compression + Headroom bounded LRU prefix cache affinity interceptor (`Priority=30`), Anthropic `cache_control` injection, stats & preview |
| App wiring & E2E | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Register all 3 modules & active interceptors, OpenAPI codegen, full live gateway interceptor pipeline E2E tests |
| Web Prompt Guard & Model Editor UI | `web/src/features/promptinjection/{schemas.ts,schemas.test.ts,hooks.ts,PromptGuardConsole.tsx}`, `web/src/app/(dashboard)/dashboard/prompt-guard/page.tsx`, `web/src/features/modeleditor/{schemas.ts,schemas.test.ts,hooks.ts,ModelEditorConsole.tsx}`, `web/src/app/(dashboard)/dashboard/model-editor/page.tsx` | Prompt Guard studio (config, rules CRUD, dry-run tester, audit events) and Model Editor & Bypass studio (overrides, bypass rules, dry-run preview) |
| Web Token Saver UI, Nav & Build | `web/src/features/tokensaver/{schemas.ts,schemas.test.ts,hooks.ts,TokenSaverConsole.tsx}`, `web/src/app/(dashboard)/dashboard/token-saver/page.tsx`, `web/src/shared/ui/AppShell.tsx`, `web/src/features/system/SystemOverview.tsx` | Token Saver & Headroom studio (config, live stats, cache clear, preview), Wave 2B navigation & system module scopes, single-binary verification |

---

### Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware)

**Files:**
- Modify: `internal/modules/gateway/module.go`
- Modify: `internal/modules/gateway/module_test.go`
- Create: `internal/modules/promptinjection/migrations/0001_init.sql`
- Create: `internal/modules/promptinjection/repo.go`
- Create: `internal/modules/promptinjection/engine.go`
- Create: `internal/modules/promptinjection/module.go`
- Create: `internal/modules/promptinjection/module_test.go`
- Create: `contracts/openapi/modules/promptinjection.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.GatewayInterceptor`, `contracts.CanonicalRequest`, `contracts.ResponseMeta`, `contracts.StreamEvent`
- Produces:
  - `gateway` propagates `icReq.Metadata` entries into `ResponseMeta.Metadata` inside `executeCandidateAttempt`.
  - `promptinjection.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger}`
  - `promptinjection.New(opts Options) *Module` (`Name() = "promptinjection"`, `Dependencies() = []string{"system"}`, `Priority() = 10`)
  - `var _ lifecycle.Module = (*Module)(nil)` and `var _ contracts.GatewayInterceptor = (*Module)(nil)`
  - SQLite tables (`0001_init.sql`):
    - `prompt_guard_config(id INTEGER PRIMARY KEY CHECK (id = 1), enabled INTEGER NOT NULL DEFAULT 0, mode TEXT NOT NULL DEFAULT 'audit', sensitivity TEXT NOT NULL DEFAULT 'medium', block_status_code INTEGER NOT NULL DEFAULT 400, updated_at TEXT NOT NULL)`
    - `prompt_guard_rules(id TEXT PRIMARY KEY, name TEXT NOT NULL, pattern TEXT NOT NULL, match_type TEXT NOT NULL, action TEXT NOT NULL, severity TEXT NOT NULL, enabled INTEGER NOT NULL DEFAULT 1, builtin INTEGER NOT NULL DEFAULT 0, description TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `prompt_guard_events(id TEXT PRIMARY KEY, request_id TEXT NOT NULL, model TEXT NOT NULL, key_id TEXT NOT NULL, rule_id TEXT NOT NULL, rule_name TEXT NOT NULL, action TEXT NOT NULL, severity TEXT NOT NULL, snippet TEXT NOT NULL, created_at TEXT NOT NULL)`
  - AdminAPI routes (under `/api`):
    - `GET /prompt-guard/config`, `PATCH /prompt-guard/config`
    - `GET /prompt-guard/rules`, `POST /prompt-guard/rules`, `PATCH /prompt-guard/rules/{id}`, `DELETE /prompt-guard/rules/{id}`
    - `POST /prompt-guard/test` -> dry-run evaluation `{ "text": "...", "protocol"?: "..." }` -> `{ "matched": bool, "action": "pass"|"audit"|"sanitize"|"block", "hits": [...], "sanitized_text": "..." }`
    - `GET /prompt-guard/events?limit=50`, `DELETE /prompt-guard/events`

- [ ] **Step 1: Write failing tests in `internal/modules/gateway/module_test.go` and `internal/modules/promptinjection/module_test.go`**
  - Add `TestGatewayInterceptorMetadataPropagatesToResponseHeaders` in `gateway/module_test.go`.
  - Add `TestPromptInjectionBlockSanitizeAuditAndBuiltins`, `TestPromptInjectionCustomRuleValidationAndRE2`, and `TestPromptInjectionDryRunAndEventLog` in `promptinjection/module_test.go` (pinning Review Focus #1).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race -run TestGatewayInterceptorMetadataPropagatesToResponseHeaders ./internal/modules/gateway/... && go test -race ./internal/modules/promptinjection/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement gateway metadata propagation, `promptinjection` module, and `contracts/openapi/modules/promptinjection.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/gateway/... ./internal/modules/promptinjection/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(promptinjection): implement prompt guard interceptor, rule engine, and OpenAPI contract`.

---

### Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules)

**Files:**
- Create: `internal/modules/modeleditor/migrations/0001_init.sql`
- Create: `internal/modules/modeleditor/repo.go`
- Create: `internal/modules/modeleditor/engine.go`
- Create: `internal/modules/modeleditor/module.go`
- Create: `internal/modules/modeleditor/module_test.go`
- Create: `contracts/openapi/modules/modeleditor.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.GatewayInterceptor`
- Produces:
  - `modeleditor.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger}`
  - `modeleditor.New(opts Options) *Module` (`Name() = "modeleditor"`, `Dependencies() = []string{"system"}`, `Priority() = 20`)
  - `var _ lifecycle.Module = (*Module)(nil)` and `var _ contracts.GatewayInterceptor = (*Module)(nil)`
  - SQLite tables (`0001_init.sql`):
    - `model_overrides(id TEXT PRIMARY KEY, name TEXT NOT NULL, model_pattern TEXT NOT NULL, target_model TEXT NOT NULL DEFAULT '', system_prompt_prefix TEXT NOT NULL DEFAULT '', system_prompt_suffix TEXT NOT NULL DEFAULT '', temperature REAL, max_tokens INTEGER, top_p REAL, extra_params_json TEXT NOT NULL DEFAULT '{}', strip_params_json TEXT NOT NULL DEFAULT '[]', priority INTEGER NOT NULL DEFAULT 100, enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `model_bypass_rules(id TEXT PRIMARY KEY, name TEXT NOT NULL, model_pattern TEXT NOT NULL, key_id TEXT NOT NULL DEFAULT '', protocol TEXT NOT NULL DEFAULT '', header_name TEXT NOT NULL DEFAULT '', header_value TEXT NOT NULL DEFAULT '', action TEXT NOT NULL, target_model TEXT NOT NULL DEFAULT '', priority INTEGER NOT NULL DEFAULT 100, enabled INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - Bypass actions: `rewrite_model`, `bypass_overrides`, `bypass_tokensaver`, `bypass_headroom`.
  - AdminAPI routes (under `/api`):
    - `GET /modeleditor/overrides`, `POST /modeleditor/overrides`, `PATCH /modeleditor/overrides/{id}`, `DELETE /modeleditor/overrides/{id}`
    - `GET /modeleditor/bypass`, `POST /modeleditor/bypass`, `PATCH /modeleditor/bypass/{id}`, `DELETE /modeleditor/bypass/{id}`
    - `POST /modeleditor/preview` -> dry-run `{ "protocol": "...", "model": "...", "body": {...}, "key_id"?: "...", "headers"?: {...} }` -> `{ "model": "...", "body": {...}, "applied_overrides": [...], "applied_bypass_rules": [...], "metadata": {...} }`

- [ ] **Step 1: Write failing tests in `internal/modules/modeleditor/module_test.go`**
  - Pin Review Focus #2:
    - `TestModelEditorOverridesAcrossAllFourProtocols`: verifies model rewrite, system prompt prefix/suffix, parameter injection (`temperature`, `max_tokens`, `extra_params`), and `strip_params` preserve valid JSON across `openai-chat`, `anthropic`, `openai-responses`, and `gemini`.
    - `TestModelBypassRulesPriorityAndMetadataFlags`: verifies `rewrite_model`, `bypass_overrides` (skips override application), `bypass_tokensaver`, and `bypass_headroom` matching on glob pattern, key ID, protocol, and request header.
    - `TestModelEditorCRUDValidationAndDryRunPreview`: CRUD + invalid JSON/action validation + `POST /api/modeleditor/preview`.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/modeleditor/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `modeleditor` module and `contracts/openapi/modules/modeleditor.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/modeleditor/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(modeleditor): implement model parameter overrides, bypass rules interceptor, and OpenAPI contract`.

---

### Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity)

**Files:**
- Create: `internal/modules/tokensaver/migrations/0001_init.sql`
- Create: `internal/modules/tokensaver/repo.go`
- Create: `internal/modules/tokensaver/compressor.go`
- Create: `internal/modules/tokensaver/headroom.go`
- Create: `internal/modules/tokensaver/module.go`
- Create: `internal/modules/tokensaver/module_test.go`
- Create: `contracts/openapi/modules/tokensaver.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.GatewayInterceptor`
- Produces:
  - `tokensaver.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger}`
  - `tokensaver.New(opts Options) *Module` (`Name() = "tokensaver"`, `Dependencies() = []string{"system"}`, `Priority() = 30`)
  - `var _ lifecycle.Module = (*Module)(nil)` and `var _ contracts.GatewayInterceptor = (*Module)(nil)`
  - SQLite table (`0001_init.sql`):
    - `tokensaver_config(id INTEGER PRIMARY KEY CHECK (id = 1), enabled INTEGER NOT NULL DEFAULT 0, strategy TEXT NOT NULL DEFAULT 'whitespace', max_history_messages INTEGER NOT NULL DEFAULT 40, max_tool_output_chars INTEGER NOT NULL DEFAULT 8000, headroom_enabled INTEGER NOT NULL DEFAULT 1, prefix_depth INTEGER NOT NULL DEFAULT 2, max_cache_entries INTEGER NOT NULL DEFAULT 1024, cache_ttl_seconds INTEGER NOT NULL DEFAULT 1800, inject_cache_control INTEGER NOT NULL DEFAULT 1, updated_at TEXT NOT NULL)`
  - AdminAPI routes (under `/api`):
    - `GET /tokensaver/config`, `PATCH /tokensaver/config`
    - `GET /tokensaver/stats` -> compression totals (`requests_processed`, `requests_compressed`, `bytes_saved`, `estimated_tokens_saved`) + Headroom stats (`cache_hits`, `cache_misses`, `hit_rate`, `active_entries`, `max_entries`, `warm_prefixes`)
    - `POST /tokensaver/stats/reset` -> zeroes counters
    - `POST /tokensaver/cache/clear` -> clears Headroom in-memory prefix cache
    - `POST /tokensaver/preview` -> dry-run compression + prefix hash `{ "protocol": "...", "model": "...", "body": {...} }`

- [ ] **Step 1: Write failing tests in `internal/modules/tokensaver/module_test.go`**
  - Pin Review Focus #3 & #4:
    - `TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass`: verifies `whitespace`, `dedupe`, and `aggressive` strategies compress prose/history/tool output while preserving fenced `` ``` `` code blocks verbatim, and `req.Metadata["bypass_tokensaver"] = true` skips compression.
    - `TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl`: verifies prefix fingerprint miss -> `OnResponseHeaders` records provider/connection affinity -> second request gets cache hit + `X-Omeg4mon-Headroom-Cache: hit` + Anthropic `cache_control: {"type":"ephemeral"}` injection.
    - `TestHeadroomLRUBoundedCapacityAndConcurrency`: concurrent `OnRequest`/`OnResponseHeaders` under `-race` never exceeds `max_cache_entries` and evicts LRU/expired entries cleanly.
    - `TestTokenSaverAdminAPIStatsResetCacheClearAndPreview`: config CRUD, stats, reset, cache clear, preview.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/tokensaver/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `tokensaver` module and `contracts/openapi/modules/tokensaver.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/tokensaver/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(tokensaver): implement Go-native token compression, headroom prefix cache affinity, and OpenAPI contract`.

---

### Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Modify: `web/src/shared/api/schema.d.ts` (generated via `make codegen`)

**Interfaces:**
- Consumes: `promptinjection.New`, `modeleditor.New`, `tokensaver.New`, `gateway.Options.Interceptors`
- Produces: Active Wave 2B modules registered in `lifecycle.Registry` and wired as ordered `GatewayInterceptor`s only when not disabled in `cfg.DisabledModules`, plus updated OpenAPI TypeScript types.

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go`**
  - Pin Review Focus #1, #3, #5:
    - `TestWave2BGatewayInterceptorPipelineE2E`: spins up `buildApp` + mock upstream, enables `promptinjection`, `modeleditor`, and `tokensaver`:
      1. Injection payload is blocked with configured status before upstream is hit.
      2. Clean payload with model override + bypass rule rewrites model & strips forbidden param on upstream request.
      3. Repeated request hits Headroom cache (`X-Omeg4mon-Headroom-Cache: miss` then `hit`) and reports `X-Omeg4mon-Tokens-Saved`.
    - `TestWave2BDisabledModulesSkipInterceptors`: sets `DisabledModules: []string{"promptinjection", "modeleditor", "tokensaver"}` and confirms gateway requests pass through untouched and disabled routes return 404.
- [ ] **Step 2: Run E2E tests to verify failure**
  - Run: `go test -race -run TestWave2B ./cmd/omeg4mon-router/...`
  - Expected: FAIL.
- [ ] **Step 3: Wire modules in `cmd/omeg4mon-router/main.go`, register OpenAPI paths/schemas in `contracts/openapi/openapi.yaml`, and run `make codegen`**
- [ ] **Step 4: Run tests, linter, and typecheck**
  - Run: `make check`
  - Expected: PASS, 0 lint/type errors.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(cmd): wire promptinjection, modeleditor, and tokensaver interceptors with OpenAPI codegen and E2E tests`.

---

### Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass

**Files:**
- Create: `web/src/features/promptinjection/schemas.ts`
- Create: `web/src/features/promptinjection/schemas.test.ts`
- Create: `web/src/features/promptinjection/hooks.ts`
- Create: `web/src/features/promptinjection/PromptGuardConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/prompt-guard/page.tsx`
- Create: `web/src/features/modeleditor/schemas.ts`
- Create: `web/src/features/modeleditor/schemas.test.ts`
- Create: `web/src/features/modeleditor/hooks.ts`
- Create: `web/src/features/modeleditor/ModelEditorConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/model-editor/page.tsx`

- [ ] **Step 1: Write failing Zod schema tests in `promptinjection/schemas.test.ts` and `modeleditor/schemas.test.ts`**
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/promptinjection/schemas.test.ts src/features/modeleditor/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement schemas, TanStack Query hooks, Console studios (strict `docs/DESIGN.md`, zero `backdrop-blur`, zero UI emojis), and static App Router pages**
- [ ] **Step 4: Run Bun tests, Biome check, and TypeScript check**
  - Run: `cd web && bun test && bun run typecheck && bun run check`
  - Expected: PASS, 0 errors.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Prompt Guard and Model Editor / Bypass studio consoles`.

---

### Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/tokensaver/schemas.ts`
- Create: `web/src/features/tokensaver/schemas.test.ts`
- Create: `web/src/features/tokensaver/hooks.ts`
- Create: `web/src/features/tokensaver/TokenSaverConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/token-saver/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

- [ ] **Step 1: Write failing Zod schema tests in `web/src/features/tokensaver/schemas.test.ts`**
- [ ] **Step 2: Run Bun test to verify failure**
  - Run: `cd web && bun test src/features/tokensaver/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement Token Saver & Headroom studio, add Wave 2B routes to `AppShell.tsx` (`Prompt Guard`, `Model Editor`, `Token Saver`), and add module descriptions to `SystemOverview.tsx`**
- [ ] **Step 4: Run full `make check`, `make build`, and single-binary smoke test**
  - Verify `bin/omeg4mon-router` starts cleanly and serves HTTP 200 on `/dashboard/prompt-guard/`, `/dashboard/model-editor/`, `/dashboard/token-saver/`, and `/api/health` reporting all 14 active modules.
  - Verify zero `backdrop-blur` / `backdrop-filter` in new Wave 2B TSX files.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Token Saver & Headroom studio, Wave 2B navigation, and verify single binary`.
