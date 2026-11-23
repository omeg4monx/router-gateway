# Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete Wave 1 by delivering async batched usage/cost telemetry, model pricing, metadata-first request logs & live SSE stats (`usage`), plus stateless CLI tool config generators with loopback-only atomic local apply/rollback (`clitools`), wired into the Go binary, OpenAPI 3.1 contract, and Next.js 16 dashboard studios.

**Architecture:** Two isolated vertical-slice Go modules (`internal/modules/usage`, `internal/modules/clitools`) implement `lifecycle.Module`, own their own SQLite migrations/state, communicate with `gateway` strictly via `contracts.UsageRecorder` wired in `cmd/omeg4mon-router/main.go`, and expose OpenAPI 3.1 fragments (`contracts/openapi/modules/{usage,clitools}.yaml`) driving strict TypeScript client codegen for `web/src/features/{usage,clitools}/`.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `jobs.Runner`, `eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), `zustand`, TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. `gateway` emits telemetry to `usage` strictly via `contracts.UsageRecorder`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339 timestamps (`"2006-01-02T15:04:05.000000000Z"`), metadata-first usage persistence (never stalling SSE streaming; bounded async channel + batch writer + scheduled retention pruning).
- **Safe `clitools` Execution:** Always generates copyable/downloadable snippets via API/UI; local filesystem apply & rollback are permitted ONLY when `cfg.IsLoopback()` is true (`403 local_apply_disabled` on non-loopback binds), using atomic temp-file + `fsync` + `os.Rename` with backup files.
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task must pass `make check` (`golangci-lint` v2 + Biome v2 + Go `-race` + zero-CGO DB + Bun tests + `tsc --noEmit`) and `graphify update .` (cleaning `graphify-out/20??-??-??/` backups).

## Review Focus

1. **Non-blocking hot path & graceful flush drain in `usage`:** `Record(ctx, rec)` MUST NEVER block gateway SSE streaming even if SQLite is busy or the channel is full (drop excess with counter/log instead of blocking), and `Stop(ctx)` MUST flush all buffered records cleanly before closing without race conditions or panic on late `Record` calls. _(Pinned in Task 1)_
2. **Accurate cost calculation & retention pruning in `usage`:** When `rec.CostUSD == 0`, cost is computed from `model_pricing` (using exact model match or bare model after `<provider>/` strip); `PruneExpired` deletes only rows older than the retention cutoff without touching newer rows. _(Pinned in Task 1)_
3. **Loopback-only local apply guard & atomic backup/rollback in `clitools`:** On non-loopback bind (`Host: "0.0.0.0"`), `POST /api/clitools/{tool}/apply` and `POST /api/clitools/{tool}/rollback` MUST return `403 local_apply_disabled` without touching the filesystem; on loopback, apply creates a backup before atomically replacing an existing file, and rollback restores exact previous bytes. _(Pinned in Task 2)_
4. **Live Gateway -> Usage E2E wiring:** Sending a real `/v1/chat/completions` request through `buildApp` records the exact model, provider, tokens, latency, and status into `GET /api/usage/logs` and `GET /api/usage/summary`. _(Pinned in Task 3)_
5. **Static-export UI completeness & zero-blur `DESIGN.md` compliance:** `/dashboard/usage/` and `/dashboard/cli-tools/` export as static trailing-slash HTML routes, serve HTTP 200 from `bin/omeg4mon-router`, and have zero `backdrop-blur`/`backdrop-filter` or UI emojis. _(Pinned in Tasks 4 & 5)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Usage module | `internal/modules/usage/{migrations/0001_init.sql,pricing.go,repo.go,batcher.go,module.go,module_test.go}`, `contracts/openapi/modules/usage.yaml` | Async batched SQLite writer, default & custom model pricing, summary/leaderboard/logs/SSE APIs, retention pruning |
| CLI Tools module | `internal/modules/clitools/{catalog.go,apply.go,module.go,module_test.go}`, `contracts/openapi/modules/clitools.yaml` | Config generators for 6 CLI tools + loopback-only atomic local apply & backup/rollback |
| Wiring & E2E | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Register `usage` & `clitools`, wire `Usage: usageMod` into `gateway`, schedule pruning job, OpenAPI codegen & E2E tests |
| Web Usage UI | `web/src/features/usage/{schemas.ts,schemas.test.ts,hooks.ts,UsageConsole.tsx}`, `web/src/app/(dashboard)/dashboard/usage/page.tsx` | Usage KPIs, range filter, model/provider leaderboards, filterable logs table, model pricing editor |
| Web CLI Tools UI | `web/src/features/clitools/{schemas.ts,schemas.test.ts,hooks.ts,CliToolsConsole.tsx}`, `web/src/app/(dashboard)/dashboard/cli-tools/page.tsx`, `web/src/shared/ui/AppShell.tsx`, `web/src/features/system/SystemOverview.tsx` | CLI tool config generator studio, copy/download/loopback apply/rollback, Wave 1 complete nav & overview + full single-binary smoke test |

---

### Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/usage/migrations/0001_init.sql`
- Create: `internal/modules/usage/pricing.go`
- Create: `internal/modules/usage/repo.go`
- Create: `internal/modules/usage/batcher.go`
- Create: `internal/modules/usage/module.go`
- Create: `internal/modules/usage/module_test.go`
- Create: `contracts/openapi/modules/usage.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `usage.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger, BatchSize int, FlushInterval time.Duration, BufferSize int, MetadataRetentionDays int}`
  - `usage.New(opts Options) *Module` (`Name() = "usage"`, `Dependencies() = []string{"system", "settings"}`)
  - `(*Module)` implements `lifecycle.Module` AND `contracts.UsageRecorder` (`var _ contracts.UsageRecorder = (*Module)(nil)`).
  - `(*Module).PruneExpired(ctx context.Context) error` (deletes `usage_logs` older than `MetadataRetentionDays`, default 90).
  - AdminAPI routes (under `/api`):
    - `GET /usage/summary` (`?range=24h|7d|30d|all`, default `24h`) -> totals (`requests`, `input_tokens`, `output_tokens`, `cached_tokens`, `cost_usd`, `avg_latency_ms`, `avg_ttft_ms`, `error_count`, `error_rate`), `by_model[]`, `by_provider[]`, `series[]` (hourly/daily buckets).
    - `GET /usage/logs` (`?model=&provider=&key_id=&status=&limit=50&offset=0`) -> `{ "logs": [...], "total": int64, "limit": int, "offset": int }`.
    - `GET /usage/stream` -> SSE stream subscribing to `usage.recorded` on `eventbus`, flushing `event: usage\ndata: {...}\n\n` immediately via `http.NewResponseController(w).Flush()` with 15s keepalive comments.
    - `GET /usage/pricing`, `PUT /usage/pricing/{model}`, `DELETE /usage/pricing/{model}` -> list/upsert/delete custom model pricing (merged with built-in defaults).

- [ ] **Step 1: Write failing tests in `internal/modules/usage/module_test.go`**
  - Pin Review Focus #1 & #2:
    - `TestUsageAsyncBatcherAndGracefulDrain`: record 75 events concurrently under `-race`, stop module cleanly, verify all 75 persisted in SQLite, and verify `Record` after `Stop` is a safe non-panicking no-op.
    - `TestUsageCostCalculationAndPricingCRUD`: built-in model cost calculation (`openai/gpt-4o`, `claude-sonnet-4`, `gemini-2.5-pro`), custom `PUT /api/usage/pricing/{model}` override, `DELETE` reset.
    - `TestUsageSummaryLogsAndRetentionPrune`: verify `GET /api/usage/summary`, `GET /api/usage/logs` filtering (`model`, `provider`, `status`), and `PruneExpired` deleting only expired rows.
    - `TestUsageLiveSSEStream`: connect to `GET /api/usage/stream`, emit a `Record`, verify SSE frame arrives promptly and closes cleanly on context cancel.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/usage/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `usage` module and `contracts/openapi/modules/usage.yaml`**
  - Follow `db.Reader()` / `db.Writer()` split, fixed-width UTC RFC3339 timestamps, strict JSON decoding, and `golangci-lint` cleanliness.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/usage/... && CGO_ENABLED=0 go test -count=1 ./internal/modules/usage/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(usage): implement async batched usage telemetry, model pricing, live SSE, and OpenAPI contract`.

---

### Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment

**Files:**
- Create: `internal/modules/clitools/catalog.go`
- Create: `internal/modules/clitools/apply.go`
- Create: `internal/modules/clitools/module.go`
- Create: `internal/modules/clitools/module_test.go`
- Create: `contracts/openapi/modules/clitools.yaml`

**Interfaces:**
- Consumes: `config.Config`
- Produces:
  - `clitools.Options{Config config.Config, HomeDir string}` (`HomeDir` defaults to `os.UserHomeDir()` when empty; tests inject `t.TempDir()`).
  - `clitools.New(opts Options) *Module` (`Name() = "clitools"`, `Dependencies() = []string{"system"}`, `Migrations() = nil`).
  - Supported 6 tools: `claude-code` (`~/.claude/settings.json`), `codex-cli` (`~/.codex/config.toml`), `opencode` (`~/.config/opencode/opencode.json`), `cline` (`~/.cline/omeg4mon-router.json`), `roo-code` (`~/.roo-code/omeg4mon-router.json`), `droid` (`~/.factory/config.json`).
  - AdminAPI routes (under `/api`):
    - `GET /clitools` -> `{ "local_apply_enabled": bool, "tools": [...] }`.
    - `POST /clitools/{tool}/generate` -> body `{ "base_url", "api_key"?, "model"? }` -> `{ "tool", "name", "format", "filename", "target_path", "content", "instructions" }`.
    - `POST /clitools/{tool}/apply` -> loopback-only (`cfg.IsLoopback()`; returns `403 local_apply_disabled` on non-loopback binds!), atomically writes config file with backup `<target>.om4.bak.<ts>` if file existed, returns `{ "applied": true, "tool", "target_path", "backup_path" }`.
    - `POST /clitools/{tool}/rollback` -> loopback-only (`403 local_apply_disabled` on non-loopback), restores latest `.om4.bak.*` backup atomically (`404 backup_not_found` if none), returns `{ "restored": true, "tool", "target_path", "backup_path" }`.

- [ ] **Step 1: Write failing tests in `internal/modules/clitools/module_test.go`**
  - Pin Review Focus #3:
    - `TestCLIToolsCatalogAndGenerateAllTools`: generates valid JSON/TOML configs for all 6 tools (`claude-code`, `codex-cli`, `opencode`, `cline`, `roo-code`, `droid`), rejects unknown tool with `404` and non-HTTP `base_url` with `400`.
    - `TestCLIToolsLoopbackAtomicApplyAndRollback`: on `Host: "127.0.0.1"` with `HomeDir: t.TempDir()`, pre-seeds target file, calls `/apply`, verifies new content written + backup created, calls `/rollback`, verifies exact original bytes restored.
    - `TestCLIToolsNonLoopbackRejectsApplyAndRollback`: on `Host: "0.0.0.0"`, both `/apply` and `/rollback` return `403 local_apply_disabled` and leave `HomeDir` completely untouched.
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -race ./internal/modules/clitools/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `clitools` module and `contracts/openapi/modules/clitools.yaml`**
  - Use temp-file + `Sync()` + `Close()` + `os.Rename()` in target parent directory for atomic writes.
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/clitools/...`
  - Expected: PASS, 0 issues.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(clitools): implement CLI tool config generators, loopback-only atomic apply/rollback, and OpenAPI contract`.

---

### Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Modify: `web/src/shared/api/schema.d.ts` (generated via `make codegen`)

**Interfaces:**
- Consumes: `usage.New`, `clitools.New`, `gateway.Options.Usage`, `jobRunner.Every`
- Produces:
  - All 9 Wave 1 modules registered in `cmd/omeg4mon-router/main.go` (`system`, `auth`, `settings`, `keys`, `providers`, `combos`, `usage`, `gateway`, `clitools`), `gatewayMod` wired with `Usage: usageMod`, and periodic `usage-prune` job registered on `jobRunner` (`Every("usage-prune", 1*time.Hour, usageMod.PruneExpired)`).
  - Root OpenAPI aggregating all `usage.yaml` and `clitools.yaml` paths & schemas + freshly generated `web/src/shared/api/schema.d.ts`.

- [ ] **Step 1: Write failing E2E test `TestAppWave1CUsageAndCLIToolsEndToEnd` in `cmd/omeg4mon-router/main_test.go`**
  - Pin Review Focus #4: boot `buildApp`, verify `/api/health` lists all 9 modules, execute a real `/v1/chat/completions` against a mock upstream, poll `GET /api/usage/logs` & `GET /api/usage/summary` until the async batcher flushes and verify the recorded model/tokens/cost/status, and verify `GET /api/clitools` + `POST /api/clitools/claude-code/generate`.
- [ ] **Step 2: Run test to verify failure**
  - Run: `go test -race -run TestAppWave1CUsageAndCLIToolsEndToEnd ./cmd/omeg4mon-router/...`
  - Expected: FAIL.
- [ ] **Step 3: Wire `usage` & `clitools` in `main.go`, aggregate `openapi.yaml`, run `make codegen`**
  - Register modules, wire `Usage: usageMod` into `gateway.Options`, register `usage-prune` job, update `contracts/openapi/openapi.yaml`, run `make codegen`.
- [ ] **Step 4: Run full `make check`**
  - Run: `make check`
  - Expected: PASS (0 lint/format issues, all Go & Web tests green).
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(cmd): wire usage telemetry and clitools modules with OpenAPI codegen and E2E tests`.

---

### Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`)

**Files:**
- Create: `web/src/features/usage/schemas.ts`
- Create: `web/src/features/usage/schemas.test.ts`
- Create: `web/src/features/usage/hooks.ts`
- Create: `web/src/features/usage/UsageConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/usage/page.tsx`

**Interfaces:**
- Consumes: typed `apiClient` (`/api/usage/summary`, `/api/usage/logs`, `/api/usage/pricing*`), `react-hook-form` + `zodResolver`, `shadcn/ui`.
- Produces:
  - `/dashboard/usage/` — Range switcher (`24h`, `7d`, `30d`, `all`), KPI metric cards (Requests, Total Tokens, Est. Cost USD, Avg Latency / TTFT, Error Rate), Top Models & Top Providers visual breakdown bars, filterable Request Logs table (`model`, `provider`, `status`, pagination), and Model Pricing Editor form (upsert & reset custom per-million-token pricing).

- [ ] **Step 1: Write failing tests in `web/src/features/usage/schemas.test.ts`**
  - Validate `upsertPricingSchema` (non-empty model, non-negative finite USD rates, rejects negative/NaN) and formatting helpers (`formatUSD`, `formatCompactTokens`, `formatDurationMs`).
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/usage/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `web/src/features/usage/*` and `/dashboard/usage/page.tsx`**
  - Strictly follow `docs/DESIGN.md` (solid opaque surfaces, zero `backdrop-blur`, zero UI emojis, `tabular-nums`, Geist fonts, `#E56A4A` accent).
- [ ] **Step 4: Run Web checks & static export build**
  - Run: `cd web && bun run check:fix && bun run check && bun test && bun run typecheck && NEXT_OUTPUT_EXPORT=true bun run build`
  - Expected: PASS, `web/out/dashboard/usage/index.html` generated.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement Usage, Telemetry & Model Pricing studio`.

---

### Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification

**Files:**
- Create: `web/src/features/clitools/schemas.ts`
- Create: `web/src/features/clitools/schemas.test.ts`
- Create: `web/src/features/clitools/hooks.ts`
- Create: `web/src/features/clitools/CliToolsConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/cli-tools/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx` (promote `Usage & Telemetry` -> `/dashboard/usage/` and `CLI Tools` -> `/dashboard/cli-tools/` into `ACTIVE_NAV_ITEMS`, update badge to `W1` / `WAVE 1 COMPLETE`)
- Modify: `web/src/features/system/SystemOverview.tsx` (reflect all 9 Wave 1 modules active!)

**Interfaces:**
- Consumes: typed `apiClient` (`/api/clitools`, `/api/clitools/{tool}/generate`, `/api/clitools/{tool}/apply`, `/api/clitools/{tool}/rollback`), `react-hook-form` + `zodResolver`, `shadcn/ui`.
- Produces:
  - `/dashboard/cli-tools/` — 6 tool selector cards, generator form (`base_url`, optional `api_key`, optional `model`), live generated config preview with Copy & Download buttons, and Loopback-Only Apply + Rollback actions (with clear `Remote-Bind Disabled` status when `local_apply_enabled === false`).

- [ ] **Step 1: Write failing tests in `web/src/features/clitools/schemas.test.ts`**
  - Validate `generateCliConfigSchema` (valid `http://`/`https://` `base_url`, optional `api_key`/`model` trimming, rejects invalid URL).
- [ ] **Step 2: Run Bun tests to verify failure**
  - Run: `cd web && bun test src/features/clitools/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement `CliToolsConsole`, `/dashboard/cli-tools/page.tsx`, and update `AppShell.tsx` & `SystemOverview.tsx`**
  - Strictly follow `docs/DESIGN.md` (zero `backdrop-blur`, zero UI emojis, WCAG AA contrast).
- [ ] **Step 4: Run FULL Wave 1 verification (`make fmt && make check && make build` + live single-binary smoke test)**
  - Verify `web/out/dashboard/usage/index.html` and `web/out/dashboard/cli-tools/index.html` exist.
  - Smoke-test `bin/omeg4mon-router` on temp port/dir: `GET /api/health` includes all 9 Wave 1 modules (`auth`, `clitools`, `combos`, `gateway`, `keys`, `providers`, `settings`, `system`, `usage`), and HTTP 200 HTML across ALL 8 dashboard routes (`/dashboard/`, `/dashboard/providers/`, `/dashboard/providers/detail/`, `/dashboard/models/`, `/dashboard/keys/`, `/dashboard/usage/`, `/dashboard/cli-tools/`, `/dashboard/settings/`).
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` and commit `feat(web): implement CLI Tools studio, complete Wave 1 navigation, and verify single binary`.
