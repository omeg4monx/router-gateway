# OMEG4MON ROUTER — Modular Next.js (TypeScript) × Go Rewrite Design

**Date:** 2026-11-22
**Branch:** `migration/ts-go`
**Status:** Approved

---

## 1. Intent & Why

Branch `main` has grown into a large JavaScript monolith (223 Next.js JS API routes, 40 dashboard pages, a JS `open-sse/` engine, and an incomplete 10-file Go prototype). That structure causes three recurring problems:

1. **Runtime overhead:** Node.js consumes ~150–350 MB RAM at idle and faces V8 GC pressure under heavy concurrent SSE streaming.
2. **Structural coupling:** JavaScript routes, storage helpers, provider executors, and UI state are tangled together, making feature upgrades brittle and slow.
3. **Incomplete Go integration & legacy ownership:** Go is only bolted on partially, the frontend stack is still plain JavaScript, and `main` carries upstream `decolua` MIT-licensed code.

This rewrite starts from a clean slate on `migration/ts-go` to deliver a **modern, strictly typed Next.js (TypeScript) dashboard** backed by a **high-performance Go (`net/http` + `chi`) modular backend**, packaged into a **single self-contained binary (`//go:embed`)**. Every feature domain from `main` is rebuilt progressively as an isolated module—upgrading performance, reliability, and code clarity rather than porting JavaScript debt.

---

## 2. Goals, Non-Goals & Anti-Goals

### Goals

- **Single-binary production runtime:** Go serves `/v1/*`, `/v1beta/*`, `/api/*`, and the embedded Next.js static export (`output: "export"`), targeting ~15–30 MB idle RAM and native goroutine SSE concurrency.
- **True modular architecture:** Every product domain lives in an isolated vertical slice (both in Go `internal/modules/<domain>/` and in `web/src/features/<domain>/`) registered via an explicit Go `Module` interface.
- **Full feature coverage via progressive waves:** Rebuild and improve all feature domains present in `main`, starting from a rock-solid platform foundation and core LLM gateway, then adding extended modules cleanly without touching core internals.
- **Contract-first Go ↔ TypeScript boundary:** Modular OpenAPI 3.1 specifications generate TypeScript types and the typed frontend HTTP client so frontend and backend never drift.
- **Clean-room IP ownership:** 100% original code written from behavioral specifications (no verbatim code or third-party asset copying from `main`), allowing a clean owner-held MIT `LICENSE`.
- **Preserve owner brand identity:** Carry forward the owner's `DESIGN.md` direction (`#E56A4A` ember orange, charcoal `#0a0a0a` / warm off-white `#FDFAF6`, Geist Sans & Geist Mono, caduceus logo, ember-line motif, solid opaque surfaces, zero UI emojis, WCAG AA).

### Non-Goals

- **Direct compatibility with legacy `main` `data.sqlite` schema:** The rewrite uses a clean, normalized, per-module SQLite schema. Migrating legacy `main` data is deferred to an optional one-way `importer` module.
- **Production Next.js SSR / Server Actions / Node.js server:** Production runs only the Go binary.
- **Out-of-process RPC/gRPC plugin binaries:** Modules are compiled in-process into the single binary and toggled via configuration.

### Anti-Goals (Failure Conditions)

- Re-creating `main`'s spaghetti in Go (cross-module direct repo imports, god packages, shared mutable globals).
- Placing business logic, credential handling, or upstream LLM calls in `web/`.
- Copying upstream `decolua`-copyrighted source files or assets verbatim from `main`.
- Degrading SSE streaming (buffering full upstream responses, applying global HTTP write timeouts to streaming endpoints, or leaking goroutines when clients disconnect).

---

## 3. Constraints

- **Clean-room discipline (Owner Decision — 1A):** Builders **may** read `main` strictly as a behavioral reference (extracting endpoint URLs, OAuth/PKCE parameters, upstream wire quirks, default model pricing, and CLI tool config formats), but **0% verbatim copying** of source code, SQL, or third-party assets from `main` is permitted. Every Go package, SQL migration, OpenAPI schema, and React/TypeScript module is authored fresh. `LICENSE` is MIT under `Copyright (c) 2026 OMEG4MON`. Third-party Go/npm dependencies must use permissive licenses compatible with MIT.
- **Single-admin dashboard model (Owner Decision — 2A):** Dashboard authentication remains single-admin. Wave 3 `teams` operates strictly as API-key, budget, and quota grouping—no multi-user dashboard RBAC or per-tenant row isolation is required in Wave 0/1.
- **Zero-CGO Go build:** Persistence uses pure-Go `modernc.org/sqlite` (WAL mode, low-RAM pragma profile suitable for 1 GB RAM / 1-core VPS) so cross-compilation (`linux/darwin/windows`, `amd64/arm64`) requires no C compiler.
- **HTTP stack:** Standard library `net/http` with thin router `github.com/go-chi/chi/v5` (no heavy framework wrappers, full access to `http.ResponseController` / `http.Flusher` / HTTP/2).
- **Frontend stack:** Next.js 16 latest (`^16.4.6`, App Router, Turbopack, `output: "export"`), React 19, strict TypeScript (`strict: true`, `noUncheckedIndexedAccess: true`), Tailwind CSS v4, `shadcn/ui` (Radix UI + CVA + `cn`), `react-hook-form` + `@hookform/resolvers` + `zod`, `zustand`, TanStack Query, Bun (`1.4+`) for workspace/scripts/build.
- **Network defaults:** Default bind `127.0.0.1:20128`, configurable via flags and environment variables (`OMEG4MON_HOST`, `OMEG4MON_PORT`, `OMEG4MON_DATA_DIR`, `OMEG4MON_DISABLED_MODULES`).
- **Local dev prerequisite:** Go toolchain (1.24+, target 1.26) must be available on `PATH` alongside Bun (`1.4+`). _(Observed during recon: `go` is not currently on default shell `PATH` in this workspace and must be installed or added to `PATH` before running Go build/test commands.)_

---

## 4. Approaches Considered

| Approach                                                        | Summary                                                                                                                                                                                                           | Trade-offs                                                                                                                                                                                 | Decision   |
| --------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- |
| **A. Modular Monolith + Module Registry + OpenAPI Contract**    | Vertical-slice modules in `internal/modules/<domain>/` implementing a strict `Module` interface, mirrored by `web/src/features/<domain>/`, glued by `contracts/openapi/` codegen and embedded into one Go binary. | Slightly higher upfront scaffolding discipline (module interface, per-module migrations, OpenAPI pipeline), but scales cleanly across 20+ feature domains without coupling or TS/Go drift. | **Chosen** |
| **B. Layered Monolith + Manual Types**                          | Traditional `handler/` + `service/` + `repo/` folders in Go with hand-written TypeScript fetch types.                                                                                                             | Fastest in week 1, but quickly couples 20+ domains together and causes Go/TS contract drift—risking a Go repeat of `main`'s structural debt.                                               | Rejected   |
| **C. Micro-Kernel + Out-of-Process Plugins (gRPC/`go-plugin`)** | Minimal core binary plus separate plugin binaries for extended features (`mitm`, `telegram`, `vouchers`, `arena`, etc.).                                                                                          | Strong crash isolation, but defeats single-binary operational simplicity and adds heavy IPC/distribution complexity.                                                                       | Rejected   |

---

## 5. System Architecture & Repository Layout

### 5.1 High-Level Runtime Flow

- **Production mode (single binary):**
  - Go binary listens on `HOST:PORT` (default `127.0.0.1:20128`).
  - `/v1/*`, `/v1beta/*`, and compatibility aliases (`/codex/*`, `/responses`, `/models/*`) route to the **Gateway Router** (API key auth, no HTTP write timeout, streaming-optimized `http.ResponseController`, upstream HTTP/2 transport pool).
  - `/api/*` routes to the **Management Router** (CORS/CSRF/session auth, standard request timeouts, JSON envelope, OpenAPI-aligned endpoints).
  - All other GET requests fall through to the **Embedded UI Handler** (`//go:embed` static files from `web/out/` with SPA fallback for client-routed pages).
- **Development mode (fast feedback):**
  - Go backend runs via `go run ./cmd/omeg4mon-router` (without `-tags embedui`, serving APIs + helpful dev stub on non-API paths) on `:20128`.
  - Next.js dev server runs via `bun run dev` on `:20127` with dev rewrites proxying `/api/*`, `/v1/*`, `/v1beta/*` to `http://127.0.0.1:20128`.

### 5.2 Repository Layout

```text
omeg4mon-router/
├── cmd/
│   └── omeg4mon-router/
│       ├── main.go                 # Thin bootstrap: config -> platform -> module registry -> server
│       ├── ui_embed.go             # //go:embed (build tag: embedui)
│       └── ui_dev.go               # Dev fallback handler (build tag: !embedui)
├── internal/
│   ├── platform/                   # Cross-cutting infrastructure (no domain business logic)
│   │   ├── config/                 # Env + CLI flag loader, validation
│   │   ├── lifecycle/              # Module interface, dependency topological sorter, Registry
│   │   ├── httpserver/             # net/http + chi server, streaming vs API middleware stacks, SPA handler
│   │   ├── db/                     # modernc.org/sqlite WAL manager (read/write pools) + modular migrator
│   │   ├── crypto/                 # AES-256-GCM secret vault, password hashing, signed session tokens
│   │   ├── eventbus/               # In-process typed pub/sub (SSE streams, telemetry, cross-module hooks)
│   │   ├── httpclient/             # Tuned upstream HTTP/2 client pool, proxy dialer hooks, timeouts
│   │   ├── jobs/                   # Background worker pool & scheduled job runner
│   │   └── logging/                # Structured log/slog logger + ring buffer for live console logs
│   ├── contracts/                  # Narrow Go interfaces for cross-module communication (zero repo coupling)
│   └── modules/                    # Isolated vertical-slice feature modules
│       ├── system/                 # /api/health, /api/version, /api/init, /api/shutdown
│       ├── auth/                   # Admin login/logout, session guard, require-login policy
│       ├── settings/               # Typed settings KV, settings history & rollback, theme/locale
│       ├── providers/              # Provider catalog, connections, provider-nodes, OAuth flows, token refresh
│       ├── keys/                   # Gateway API keys, allowed models, scopes, tags
│       ├── combos/                 # Multi-model fallback chains, aliases, disabled models, account cooldowns
│       ├── gateway/                # /v1/* & /v1beta/* proxy, Canonical IR translators, SSE pipeline
│       ├── usage/                  # Async batch usage/cost writer, pricing, request logs, SSE live stats
│       ├── clitools/               # Config generators for Claude Code, Codex, OpenCode, Cline, etc.
│       └── ...                     # Wave 2 & Wave 3 modules (see Section 7)
├── contracts/
│   └── openapi/
│       ├── openapi.yaml            # Root OpenAPI 3.1 document referencing per-module specs
│       └── modules/                # Per-module OpenAPI fragments (auth.yaml, providers.yaml, etc.)
├── web/                            # Next.js (App Router, strict TS, output: "export", Tailwind v4)
│   ├── package.json
│   ├── next.config.ts
│   ├── tsconfig.json
│   └── src/
│       ├── app/                    # Thin route shells (login, dashboard layout, feature routes)
│       ├── features/               # Feature-modular UI mirroring internal/modules/*
│       └── shared/                 # Design system (DESIGN.md), generated OpenAPI client, query/auth utils
├── docs/
│   ├── DESIGN.md                   # Owner visual direction (brand, palette, typography, rules)
│   └── superpowers/
├── scripts/                        # Build, codegen, and release helper scripts
├── Makefile                        # Unified commands: dev, codegen, test, build, cross-build
├── Dockerfile                      # Multi-stage: Bun build web/out -> Go embedui build -> minimal image
├── go.mod
├── go.sum
└── LICENSE                         # Owner-held MIT license
```

---

## 6. Core Technical Contracts

### 6.1 Go Module Lifecycle Contract (`internal/platform/lifecycle`)

Every feature module implements:

```go
type Module interface {
    Name() string
    Dependencies() []string
    Migrations() fs.FS // nil if module owns no SQL tables
    RegisterRoutes(mounts RouteMounts)
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type RouteMounts struct {
    PublicAPI chi.Router // /api/* routes that bypass dashboard session auth (e.g. /api/health, /api/auth/login)
    AdminAPI  chi.Router // /api/* routes protected by dashboard session auth middleware
    Gateway   chi.Router // /v1/*, /v1beta/*, and compatibility routes (API-key auth, no write timeout)
}
```

**Module boundary rules:**

1. Modules **never** import another module's `internal/modules/<other>` package.
2. Cross-module calls use narrow interfaces in `internal/contracts/` (e.g., `KeyValidator`, `ProviderResolver`, `UsageRecorder`, `ProxyDialerResolver`, `GatewayInterceptor`) wired in `cmd/omeg4mon-router/main.go`, or async events via `internal/platform/eventbus`.
3. **Upfront Gateway Interception Hook Chain _(Architect Call)_:** `internal/contracts/` defines an ordered `GatewayInterceptor` interface (`Priority() int`, `OnRequest(ctx, *CanonicalRequest) error`, `OnResponseHeaders(ctx, *ResponseMeta) error`, `OnStreamChunk(ctx, *StreamEvent) error`) registered at startup and executed by `gateway`. This allows Wave 2 modules (`promptinjection`, `tokensaver`, `headroom`, `modelbypass`) to plug in cleanly without editing `gateway` internals.
4. Every module owns its own SQL migration files (`0001_init.sql`, etc.), tracked in `schema_migrations(module TEXT, version INT, applied_at TEXT, PRIMARY KEY (module, version))`. No module reads or writes another module's tables directly.
5. `Registry` validates dependencies via topological sort on startup, runs migrations in dependency order, mounts routes, starts background workers, and shuts down in reverse dependency order on `SIGINT`/`SIGTERM`.
6. Optional modules (Wave 2 & Wave 3) can be disabled via `OMEG4MON_DISABLED_MODULES`; core modules (`system`, `auth`, `settings`, `providers`, `keys`, `combos`, `gateway`, `usage`) are always active.

### 6.2 SQLite Concurrency, Secret Vault & Usage Retention (`internal/platform/db`, `internal/platform/crypto`)

To keep RAM low and prevent `SQLITE_BUSY` under concurrent SSE traffic:

- Pure-Go `modernc.org/sqlite` with `WAL`, `busy_timeout=10000`, `synchronous=NORMAL`, `foreign_keys=ON`, `temp_store=MEMORY`, `cache_size=-2000` (~2 MB cap), `mmap_size=0`.
- **Split connection pools** against the same SQLite file:
  - **Writer pool:** `MaxOpenConns(1)` — serializes writes cleanly in Go without SQLite lock thrashing.
  - **Reader pool:** `MaxOpenConns(8)` (`mode=ro` or read-only query path) — concurrent WAL reads without blocking on writers.
- **Hot-path isolation:** Gateway requests read active provider/key/combo state from an in-memory read-through cache (invalidated via `eventbus` on config mutations) and emit usage/log records into a bounded async channel flushed in batches (e.g. every 250ms or 50 records) so disk I/O never stalls token streaming.
- **Secret Vault Master-Key Lifecycle _(Architect Call)_:** Provider secrets and OAuth tokens are encrypted at rest with AES-256-GCM using a 32-byte master key loaded from `OMEG4MON_MASTER_KEY` if provided, or auto-generated once on first boot and persisted to `<OMEG4MON_DATA_DIR>/master.key` (`0600` permissions). It is **never** derived from the admin login password so unattended VPS/Docker restarts can decrypt credentials automatically.
- **Usage & Request Log Retention _(Architect Call)_:** `usage` persists **metadata by default** (request ID, key ID, model, provider, connection ID, input/output/cached tokens, cost, latency/TTFT, HTTP status, error summary) to protect the 15–30 MB RAM / 1 GB VPS budget. Full request/response payload logging is an **opt-in setting** (off by default) with a strict per-field truncation cap (32 KB), automatic secret/header redaction, and a scheduled retention pruning job (default 7 days for payloads, 90 days for metadata).

### 6.3 Gateway Routing, Translation & SSE Streaming Pipeline (`internal/modules/gateway`)

- **Protocol coverage:** OpenAI Chat Completions (`/v1/chat/completions`), Anthropic Messages (`/v1/messages`, `/v1/messages/count_tokens`), OpenAI Responses (`/v1/responses`, `/v1/responses/compact`), Gemini (`/v1beta/models`, `/v1beta/models/{model}:{action}`), Model Listing (`/v1/models`, `/v1/models/{id}`), plus compatibility aliases (`/codex/*` → `/v1/responses`, `/responses` → `/v1/responses`, `/models/*` → `/v1/models/*`).
- **Model Addressing, Resolution Precedence & Key Scoping _(Architect Call)_:**
  - Clients may pass either `<provider>/<model>` or bare `<model>` names.
  - Deterministic resolution precedence: **(1) Active Combo name → (2) Model Alias → (3) Explicit `<provider>/<model>` → (4) Bare `<model>` matched across enabled connected providers ordered by connection priority**.
  - `keys.allowed_models`: empty/wildcard (`*`) permits all enabled models; when non-empty, authorization checks the **requested identifier** (combo name, alias, or `<provider>/<model>` / bare model). Authorizing a combo authorizes its member targets for that request.
  - `/v1/models` returns the deduplicated union of active combos, aliases, and enabled connected provider models, filtered by the caller's API key `allowed_models`.
- **Translation & Local-Compute Depth _(Architect Call)_:**
  - **Canonical Internal Representation (IR)** preserves messages, system instructions, streaming deltas, parallel tool/function calls, structured JSON outputs, multimodal image blocks, and reasoning/thinking blocks.
  - **Direct fast-path pass-through** when inbound protocol matches upstream protocol (preserving provider-specific fields such as Anthropic `cache_control` verbatim with minimal allocation).
  - **Direct pair translators** where Canonical IR would drop fidelity (e.g., Anthropic thinking/cache-control ↔ OpenAI Responses).
  - **`/v1/responses` (`previous_response_id`):** Passed through natively to OpenAI Responses upstreams; for non-Responses upstreams, Wave 1 keeps a lightweight bounded TTL response-chain table in SQLite (or returns a clean `400` if disabled).
  - **`/v1/messages/count_tokens` & `/v1/responses/compact`:** Proxied natively when the resolved upstream supports them; `/v1/messages/count_tokens` falls back to a lightweight zero-table Go token estimator when an upstream lacks a native count endpoint.
- **Streaming discipline:**
  - Server configured with `ReadHeaderTimeout: 10s`, `IdleTimeout: 120s`, and **`WriteTimeout: 0`** on the HTTP server; non-streaming `/api/*` routes enforce per-route context timeouts via middleware (`chi` `middleware.Timeout`), while `/v1/*` streaming routes respect `r.Context()` cancellation + upstream response-header timeout only.
  - Stream chunks directly using `bufio.Reader` / scanner with immediate `http.NewResponseController(w).Flush()` per SSE event.
- **Fallback & error classification:**
  - Distinguish **client/request-scoped errors** (400 malformed payload, context length exceeded, invalid params) from **provider/account errors** (401/403 auth, 429 rate limit, 5xx upstream outage, network timeout).
  - Client errors return immediately without burning fallback accounts; provider/account errors trigger automatic OAuth refresh (on 401), per-account cooldown backoff, account rotation, and combo fallback chain traversal.
- **Provider OAuth & CLI Tools Runtime Modes _(Architect Call)_:**
  - **Headless-safe OAuth (`providers`):** Every OAuth flow supports both automatic local callback (via `/api/oauth/callback` or short-lived local callback listener when required by the provider) **and** manual authorization-code / callback-URL paste (or device-code flow) in the dashboard UI so headless VPS/Docker deployments work without a browser on the server.
  - **Safe `clitools` execution:** Always generates copyable/downloadable config snippets via API/UI; optional one-click local config apply (with atomic backup & rollback) is enabled **only** when running locally on loopback and automatically disabled in Docker/remote-bind deployments.

### 6.4 OpenAPI Contract, Security Posture & Static-Export Frontend (`contracts/` + `web/`)

- **Contract pipeline:** Each backend module defines its endpoints/schemas in `contracts/openapi/modules/<name>.yaml`. Running `make codegen` bundles `contracts/openapi/openapi.yaml` and runs `openapi-typescript` to generate `web/src/shared/api/schema.d.ts`, consumed by a thin `openapi-fetch` client + TanStack Query hooks in `web/src/features/<module>/`.
- **Static-export routing rule _(Architect Call)_:**
  - Because `next.config.ts` uses `output: "export"`, entity detail views use static route shells + query params for dynamic entity IDs (e.g., `/dashboard/providers/detail?id=...`) so Next.js `output: "export"` works with zero build-time `generateStaticParams` hacks and direct deep-link refresh always hits a real static HTML file.
- **Auth & Remote-Bind Posture _(Architect Call)_:**
  - Fresh installs default to `require_login=true`, bootstrapped either via first-run `/api/init` onboarding screen or `OMEG4MON_ADMIN_PASSWORD` env var.
  - Dashboard UI authenticates via `POST /api/auth/login` returning an `HttpOnly`, `SameSite=Lax`, signed session cookie (HMAC-SHA256 backed by vault key). Setting `require_login=false` on loopback (`127.0.0.1` / `::1`) bypasses `/api/*` dashboard login **only** (never `/v1/*` API-key auth unless explicitly opted in).
  - When bound to a non-loopback address (e.g. `0.0.0.0` / Docker), the server enforces `require_login=true` and disables `/api/shutdown` unless explicitly overridden via `OMEG4MON_ALLOW_INSECURE_REMOTE=true`.
  - Gateway `/v1/*` and `/v1beta/*` endpoints authenticate via `Authorization: Bearer <key>` or `x-api-key: <key>` (or `key` query param for Gemini compatibility).
- **Visual design (`docs/DESIGN.md`):**
  - Strictly implements owner's direction: `#E56A4A` ember orange accent, solid charcoal dark (`#0a0a0a` family) and warm off-white light (`#FDFAF6`), neutral grays (`#6B7280`), Geist Sans + Geist Mono, caduceus SVG logo, ember-line gradient divider (`from-transparent via-primary to-transparent`), solid opaque surfaces (no glassmorphism), minimal elevation, ENERGY 2 / RHYTHM 2 / MOTION 1, zero UI emojis, WCAG AA contrast.

---

## 7. Phased Modular Delivery Roadmap

Every feature domain from `main` is mapped to an isolated module and delivered in waves. Each wave is independently deployable, testable, and reviewable.

### Wave 0 — Platform Foundation & Build Pipeline

- Repo structure, clean MIT `LICENSE`, `docs/DESIGN.md`, `Makefile`, `Dockerfile`.
- `internal/platform/*` (config, `chi` server, SQLite split-pool WAL + modular migrator, AES-256-GCM vault, event bus, HTTP/2 client pool, `slog` logger, module registry).
- `contracts/openapi/` + `make codegen` pipeline.
- `web/` Next.js (strict TS, Tailwind v4, `output: "export"`, Geist fonts, `DESIGN.md` primitives, dashboard shell, login screen).
- `internal/modules/system` (`/api/health`, `/api/version`, `/api/init`, `/api/shutdown`).
- Verification: `make test` and `make build` produce a single working binary serving the embedded dashboard shell and `/api/health`.

### Wave 1 — Core LLM Gateway & Essential Control Plane _(Target of Initial Implementation Plan: Waves 0–1)_

- `auth` — single-admin local password setup/login, signed session cookie, `require_login` toggle, non-loopback safety guard.
- `settings` — typed settings store, history & rollback, theme/locale preferences.
- `providers` — pluggable `ProviderAdapter` registry shipping in Wave 1 with:
  - Generic **OpenAI-compatible**, **Anthropic-compatible**, **OpenAI Responses**, and **Gemini** API-key executors (covering OpenAI, Anthropic, OpenRouter, custom compatible provider-nodes, etc.).
  - **3 core OAuth adapters:** **Claude**, **Codex**, and **Gemini** (including encrypted credential vault, automatic OAuth token refresh, headless-safe callback/code-paste wizard, and connection quick-test).
  - Remaining dedicated provider adapters (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity, etc.) plug into the exact same `ProviderAdapter` interface in subsequent increments without modifying gateway core.
- `keys` — API keys CRUD, per-key `allowed_models`, rate/quota limits, tags.
- `combos` — ordered multi-model fallback combos, model aliases, disabled models, multi-account rotation & cooldown state machine.
- `gateway` — `/v1/chat/completions`, `/v1/messages`, `/v1/messages/count_tokens`, `/v1/responses`, `/v1/models`, `/v1beta/models/*`, compatibility aliases, Canonical IR + direct fast-path translators, streaming SSE pipeline, error classification & fallback execution, and ordered `GatewayInterceptor` hook chain.
- `usage` — async batched SQLite writer, model pricing table, metadata-first request logs (+ opt-in redacted payload capture), latency telemetry, real-time SSE stats stream, charts/leaderboard API.
- `clitools` — stateless config generators for Claude Code, Codex CLI, OpenCode, Cline, Roo, Droid, etc., plus optional loopback-only atomic local apply with backup/rollback.
- Matching `web/src/features/*` pages: Login / First-Run Init, Dashboard Overview, Endpoint, Providers Control Studio, Combos, API Keys, Usage & Request Logs, CLI Tools, Settings & Pricing.

### Wave 2 — Routing Enhancements, Provider Expansion, Security & Developer Tools _(Subsequent Plan)_

- Additional dedicated provider adapters plugged into `providers` (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity, etc.).
- `proxypools` — HTTP/SOCKS5 proxy pools, latency tester, auto-scanner, provider binding.
- `translator` — interactive protocol translator workbench + live console log SSE stream.
- `modeleditor` & `modelbypass` — custom model definitions/overrides and bypass rules (plugged via `GatewayInterceptor`).
- `tokensaver` & `headroom` — Go-native token compression, cache-aware affinity, and in-memory semantic/KV cache stats (plugged via `GatewayInterceptor`).
- `promptinjection` — configurable prompt guard/injection middleware (plugged via `GatewayInterceptor`).
- `security` & `notifications` — security audit trail, system notifications, background dashboard job status (`jobs`).

### Wave 3 — Extended Ecosystem, Key Grouping & Packaging _(Subsequent Plan)_

- `teams` — API-key, budget, and quota grouping under the single-admin dashboard (no multi-user RBAC).
- `vouchers` — voucher creation, claim flow, and anti-abuse enforcement.
- `media` — audio (`/v1/audio/*`), images (`/v1/images/*`), videos (`/v1/videos/*`), web search/fetch (`/v1/search`, `/v1/web/fetch`) and media provider combos.
- `tunnel` — Cloudflare quick tunnel and Tailscale integration manager.
- `mitm` — Go-native local MITM proxy, root CA cert generation/install, IDE traffic interceptors.
- `pxpipe`, `telegram` (bot notifications/commands), `arena` & `basicchat` (comparison & chat UI), `plugins` & `skills`, `cloudsync` (optional remote sync), `importer` (optional one-way legacy `data.sqlite` importer), and `cli/` (thin multi-platform npm binary launcher).

---

## 8. Confirmed Owner Decisions & Architect Calls

### Owner Decisions

1. **Stack & Runtime:** Next.js (strict TypeScript, `output: "export"`) + Go (`net/http` + `chi/v5`), packaged as a **Single Binary (`//go:embed`)**.
2. **Architecture:** **Opsi A — Modular Monolith + Module Registry + OpenAPI Contract**, rebuilt progressively across all feature domains.
3. **Clean-Room Rule (1A):** Reading `main` is permitted strictly as a behavioral specification reference; **0% verbatim copying** of code or upstream assets.
4. **License:** MIT License — `Copyright (c) 2026 OMEG4MON`.
5. **Teams Model (2A):** Single-admin dashboard; Wave 3 `teams` operates as API-key, budget, and quota groups.
6. **Initial Plan Horizon:** **Waves 0–1** (Platform Foundation + Core Gateway & Control Plane with generic executors + 3 core OAuth providers: Claude, Codex, Gemini); remaining providers and Waves 2–3 follow modularly.

### Architect Calls (Check While Reviewing)

1. **Dynamic routes under static export:** Static route shells + query params (e.g., `/dashboard/providers/detail?id=<id>`) for refresh-safe `output: "export"` builds.
2. **Database:** Clean per-module SQLite schema (`modernc.org/sqlite`, WAL mode, split 1-writer / 8-reader pools); legacy `data.sqlite` migration deferred to optional Wave 3 `importer`.
3. **Secret vault master key:** Loaded from `OMEG4MON_MASTER_KEY` or auto-generated to `<DATA_DIR>/master.key` (`0600`), never derived from admin password, enabling unattended VPS/Docker restarts.
4. **Auth & remote-bind posture:** `require_login=true` by default (first-run `/api/init` or `OMEG4MON_ADMIN_PASSWORD`); `require_login=false` allowed on loopback only for `/api/*`; non-loopback binds enforce login and disable `/api/shutdown` unless overridden.
5. **Model resolution & key scoping:** Precedence `Combo → Alias → <provider>/<model> → Bare model (by priority)`; `allowed_models` authorizes the requested identifier; `/v1/models` returns the per-key filtered union.
6. **Gateway translation & hooks:** Canonical IR + direct pass-through + direct Anthropic ↔ Responses pair translator; lightweight heuristic fallback for `count_tokens`; ordered `GatewayInterceptor` hook chain wired in Wave 1 for Wave 2 plug-ins.
7. **Headless-safe OAuth & safe `clitools`:** OAuth supports both auto-callback and manual code/URL paste + device code; `clitools` always generates copyable snippets, with atomic local file apply enabled only on loopback.
8. **Usage log retention:** Metadata persisted by default; full request/response bodies opt-in (off by default, 32 KB cap, secret redaction, auto-prune).
9. **Frontend tooling:** Bun (`1.4+`) + `openapi-typescript` / `openapi-fetch` + TanStack Query, preserving `docs/DESIGN.md` verbatim.

---

## 9. Left to the Builder

- Exact Go struct/package internal organization inside each module directory (`internal/modules/<name>/`).
- Detailed SQL column definitions and secondary indexes per module (within the per-module ownership and migration rules in Section 6.1).
- Exact React component hierarchy, state hooks, and Tailwind v4 utility composition (within `docs/DESIGN.md` constraints).
- Detailed Wave-by-wave implementation task breakdowns (produced via `writing-plans` after this design is approved).
