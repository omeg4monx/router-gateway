# Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 3D by adding a narrow `contracts.ChatCompleter` gateway contract and two isolated vertical-slice modules — `telegram` (vault-encrypted Telegram Bot configuration, fail-closed `X-Telegram-Bot-Api-Secret-Token` + admin user ID webhook handler, in-process command controller & dashboard simulator, eventbus alert forwarder, and delivery audit log) and `arena` (concurrent 2–4 model side-by-side comparison battles, winner/tie voting, live model leaderboard, and persistent multi-session Basic Chat workbench) — complete with OpenAPI 3.1 contracts, Next.js 16 static-export studios, E2E tests, and single-binary embedding.

**Architecture:** `internal/contracts` defines `ChatCompleter` (`CompleteChat`), implemented by `gateway.Module` via an internal admin-scoped invocation of its existing inference pipeline (reusing interceptors, resolution, OAuth refresh, fallback, translation, and usage accounting). `telegram` owns `telegram_config` and `telegram_events`, stores `bot_token` and `webhook_secret` encrypted in `*crypto.Vault`, exposes `POST /api/telegram/webhook` on `PublicAPI` with constant-time secret + configured-user-ID fail-closed validation, subscribes to `notifications.emit` / `jobs.failed` / `auth.policy_changed` on `*eventbus.Bus`, and uses narrow contracts (`ProviderStore`, `ModelLister`, `JobInspector`, `ActiveModules`) for command responses. `arena` owns `arena_battles`, `arena_contenders`, `chat_sessions`, and `chat_messages`, executes concurrent model battles and multi-turn chat turns via `contracts.ChatCompleter`, and lists available models via `contracts.ModelLister` with zero cross-module Go imports.

**Tech Stack:** Go 1.26+ (`net/http`, `crypto/subtle`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `*crypto.Vault`, `*eventbus.Bus`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded request/response body limits, bounded goroutines and per-call timeouts.
- **Security Posture:** `telegram` encrypts `bot_token` and `webhook_secret` at rest via `*crypto.Vault`, never returns raw secrets in HTTP responses, fails closed when `user_id` or `webhook_secret` is unset, validates `X-Telegram-Bot-Api-Secret-Token` using `crypto/subtle.ConstantTimeCompare`, and fixes the production Telegram API base URL to `https://api.telegram.org` (overridable only via `Options.APIBaseURL` in tests).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task passes `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and final verification runs `make check`, `make build`, and `graphify update .`.

## Review Focus

1. **Telegram fail-closed webhook security & secret isolation:** Missing/mismatched `X-Telegram-Bot-Api-Secret-Token` returns `401`; empty `user_id` or mismatched sender `from.id` returns `403`; `GET /api/telegram/config` and `PATCH /api/telegram/config` never return plaintext `bot_token` or `webhook_secret`. _(Pinned in Task 2 & Task 4)_
2. **Telegram eventbus delivery & clean lifecycle teardown:** Publishing `jobs.failed` or `notifications.emit` when enabled dispatches a single bounded HTTP `sendMessage` call and logs the event; `Stop()` drains the background subscriber cleanly without goroutine leaks or panics. _(Pinned in Task 2)_
3. **Arena concurrent execution isolation & partial failure resilience:** Running a battle where 1 of 2+ contenders fails (e.g., unknown model or upstream 500) still completes and persists the battle with the succeeding contender's response and the failing contender's error summary, without aborting sibling contenders. _(Pinned in Task 3 & Task 4)_
4. **Basic Chat session history ordering & model validation:** Sending a message to a nonexistent session returns `404`; sending an empty prompt returns `400`; multi-turn turns pass prior `user`/`assistant` messages in chronological order to `ChatCompleter` and persist both turns atomically. _(Pinned in Task 3 & Task 4)_
5. **Disabled optional modules cleanly detach & static-export UI compliance:** Disabling `telegram` or `arena` via `OMEG4MON_DISABLED_MODULES` omits their routes cleanly without affecting core gateway traffic; `/dashboard/telegram/`, `/dashboard/arena/`, and `/dashboard/chat/` export cleanly to `web/out/` with zero `backdrop-blur` or UI emojis. _(Pinned in Task 4 & Task 5)_

---

### Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`

**Files:**
- Modify: `internal/contracts/contracts.go`
- Create: `internal/modules/gateway/complete_chat.go`
- Modify: `internal/modules/gateway/module.go`
- Modify: `internal/modules/gateway/module_test.go`

**Interfaces:**
- Produces in `internal/contracts`:
  - `ChatMessage{Role string, Content string}`
  - `ChatCompletionRequest{Model string, Messages []ChatMessage, Temperature *float64, MaxTokens int}`
  - `ChatCompletionResult{Model string, Content string, PromptTokens int64, CompletionTokens int64, TotalTokens int64, LatencyMs int64}`
  - `ChatCompleter` interface: `CompleteChat(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResult, error)`
- Produces in `internal/modules/gateway`:
  - `(*Module).CompleteChat(ctx context.Context, req contracts.ChatCompletionRequest) (contracts.ChatCompletionResult, error)` satisfying `contracts.ChatCompleter`

- [ ] **Step 1: Write failing tests in `internal/modules/gateway/module_test.go`** (`TestGatewayCompleteChat`, verifying non-streaming execution through resolver/provider/interceptors, token extraction, and error propagation on missing model)
- [ ] **Step 2: Run tests to verify failure** (`go test -run TestGatewayCompleteChat ./internal/modules/gateway/...`)
- [ ] **Step 3: Add `ChatCompleter` types to `internal/contracts/contracts.go`, extract `handleInferenceWithKey` in `module.go`, and implement `CompleteChat` in `complete_chat.go`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI)

**Files:**
- Create: `internal/modules/telegram/migrations/0001_init.sql`
- Create: `internal/modules/telegram/repo.go`
- Create: `internal/modules/telegram/client.go`
- Create: `internal/modules/telegram/commands.go`
- Create: `internal/modules/telegram/module.go`
- Create: `internal/modules/telegram/module_test.go`
- Create: `contracts/openapi/modules/telegram.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*crypto.Vault`, `*http.Client`, `*eventbus.Bus`, `*slog.Logger`, `contracts.ProviderStore`, `contracts.ModelLister`, `contracts.JobInspector`, `ActiveModules func() []string`, `APIBaseURL string` (optional test override)
- Produces:
  - `telegram.New(opts Options) *Module` (`Name() = "telegram"`, `Dependencies() = []string{"system"}`)
  - Implements `lifecycle.Module`
  - SQLite tables (`0001_init.sql`):
    - `telegram_config(id INTEGER PRIMARY KEY CHECK (id = 1), enabled INTEGER NOT NULL DEFAULT 0, bot_token_enc TEXT NOT NULL DEFAULT '', bot_token_mask TEXT NOT NULL DEFAULT '', chat_id TEXT NOT NULL DEFAULT '', user_id TEXT NOT NULL DEFAULT '', webhook_url TEXT NOT NULL DEFAULT '', webhook_secret_enc TEXT NOT NULL DEFAULT '', notify_on_job_failure INTEGER NOT NULL DEFAULT 1, notify_on_auth_change INTEGER NOT NULL DEFAULT 1, notify_on_emit INTEGER NOT NULL DEFAULT 0, updated_at TEXT NOT NULL)`
    - `telegram_events(id TEXT PRIMARY KEY, direction TEXT NOT NULL, kind TEXT NOT NULL, command TEXT NOT NULL DEFAULT '', status TEXT NOT NULL, summary TEXT NOT NULL DEFAULT '', error TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
  - PublicAPI route (`/api`): `POST /telegram/webhook`
  - AdminAPI routes (`/api`):
    - `GET /telegram/config`, `PATCH /telegram/config`
    - `POST /telegram/secret/rotate`
    - `POST /telegram/webhook/register`, `POST /telegram/webhook/delete`
    - `POST /telegram/test`, `POST /telegram/command`
    - `GET /telegram/events`, `DELETE /telegram/events`

- [ ] **Step 1: Write failing tests in `internal/modules/telegram/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `telegram` migrations, repo, Bot API client, command dispatcher, module handlers, eventbus listener, and `contracts/openapi/modules/telegram.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI)

**Files:**
- Create: `internal/modules/arena/migrations/0001_init.sql`
- Create: `internal/modules/arena/repo.go`
- Create: `internal/modules/arena/module.go`
- Create: `internal/modules/arena/module_test.go`
- Create: `contracts/openapi/modules/arena.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `contracts.ChatCompleter`, `contracts.ModelLister`, `*eventbus.Bus`, `*slog.Logger`
- Produces:
  - `arena.New(opts Options) *Module` (`Name() = "arena"`, `Dependencies() = []string{"system", "gateway"}`)
  - Implements `lifecycle.Module`
  - SQLite tables (`0001_init.sql`):
    - `arena_battles(id TEXT PRIMARY KEY, prompt TEXT NOT NULL, system_prompt TEXT NOT NULL DEFAULT '', temperature REAL, max_tokens INTEGER NOT NULL DEFAULT 0, winner_model TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
    - `arena_contenders(id TEXT PRIMARY KEY, battle_id TEXT NOT NULL REFERENCES arena_battles(id) ON DELETE CASCADE, slot_index INTEGER NOT NULL, model TEXT NOT NULL, content TEXT NOT NULL DEFAULT '', status TEXT NOT NULL, error TEXT NOT NULL DEFAULT '', latency_ms INTEGER NOT NULL DEFAULT 0, prompt_tokens INTEGER NOT NULL DEFAULT 0, completion_tokens INTEGER NOT NULL DEFAULT 0, total_tokens INTEGER NOT NULL DEFAULT 0)`
    - `chat_sessions(id TEXT PRIMARY KEY, title TEXT NOT NULL, model TEXT NOT NULL, system_prompt TEXT NOT NULL DEFAULT '', temperature REAL, max_tokens INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL, updated_at TEXT NOT NULL)`
    - `chat_messages(id TEXT PRIMARY KEY, session_id TEXT NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE, role TEXT NOT NULL, content TEXT NOT NULL, model TEXT NOT NULL DEFAULT '', latency_ms INTEGER NOT NULL DEFAULT 0, prompt_tokens INTEGER NOT NULL DEFAULT 0, completion_tokens INTEGER NOT NULL DEFAULT 0, total_tokens INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL)`
  - AdminAPI routes (`/api`):
    - `GET /arena/models`
    - `GET /arena/battles`, `POST /arena/battles/run`, `GET /arena/battles/{id}`, `POST /arena/battles/{id}/vote`, `DELETE /arena/battles/{id}`
    - `GET /arena/leaderboard`
    - `GET /chat/sessions`, `POST /chat/sessions`, `GET /chat/sessions/{id}`, `PATCH /chat/sessions/{id}`, `DELETE /chat/sessions/{id}`, `POST /chat/sessions/{id}/messages`

- [ ] **Step 1: Write failing tests in `internal/modules/arena/module_test.go`**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement `arena` migrations, repo, battle runner, chat session handlers, and `contracts/openapi/modules/arena.yaml`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 4: Wire Wave 3D in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go` (`TestWave3DTelegramArenaAndChatE2E`, `TestWave3DDisabledModulesDetachCleanly`)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Wire `telegram` and `arena` in `main.go`, update `openapi.yaml`, run `make codegen`**
- [ ] **Step 4: Run `golangci-lint` & `-race` tests**
- [ ] **Step 5: Commit**

---

### Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/telegram/{schemas.ts,schemas.test.ts,hooks.ts,TelegramConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/telegram/page.tsx`
- Create: `web/src/features/arena/{schemas.ts,schemas.test.ts,hooks.ts,ArenaConsole.tsx,ChatConsole.tsx}`
- Create: `web/src/app/(dashboard)/dashboard/arena/page.tsx`
- Create: `web/src/app/(dashboard)/dashboard/chat/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

- [ ] **Step 1: Write failing Zod schema tests (`web/src/features/telegram/schemas.test.ts`, `web/src/features/arena/schemas.test.ts`)**
- [ ] **Step 2: Run tests to verify failure**
- [ ] **Step 3: Implement schemas, hooks, `TelegramConsole`, `ArenaConsole`, `ChatConsole`, static-export pages, `AppShell`, and `SystemOverview`**
- [ ] **Step 4: Run `make check`, `make build`, and `graphify update .`**
- [ ] **Step 5: Commit**
