# Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver Wave 2C — completing all Wave 2 domains — by adding background job observability and on-demand triggering to `internal/platform/jobs`, implementing the isolated `security` module (security posture snapshot, redacted security audit trail, eventbus auto-audit ingestion, JSON export, and retention pruning), implementing the isolated `notifications` module (persistent notification inbox, unread counters, live SSE stream, eventbus auto-alerts, and `/api/jobs` dashboard control endpoints), and shipping matching Next.js 16 static-export studios embedded into the single binary.

**Architecture:** `internal/platform/jobs.Runner` tracks thread-safe per-job telemetry (`JobStatus`), prevents overlapping executions of the same job, exposes `Snapshot()` and `Trigger(ctx, name)`, and notifies an optional observer callback on completion. `internal/contracts` defines narrow `JobStatus` and `JobInspector` types so `notifications` exposes `/api/jobs` and `/api/jobs/{name}/trigger` without importing `platform/jobs` or other modules. Both `security` and `notifications` subscribe to existing `eventbus.Bus` control-plane mutation topics (`auth.policy_changed`, `keys.changed`, `providers.changed`, `settings.updated`, `combos.changed`, `proxypools.changed`, `promptinjection.config_changed`, `modeleditor.rules_changed`, `tokensaver.config_changed`, plus `security.audit`, `notifications.emit`, and `jobs.failed`) so cross-module observability happens with zero cross-module Go imports.

**Tech Stack:** Go 1.26+ (`net/http`, `chi/v5`, `modernc.org/sqlite` zero-CGO WAL, `eventbus.Bus`, `jobs.Runner`), OpenAPI 3.1 (`openapi-typescript` + `openapi-fetch`), Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`, strict TypeScript), React 19, Tailwind CSS v4, `shadcn/ui`, `react-hook-form` + `@hookform/resolvers` (`zod`), TanStack Query, `golangci-lint` v2, Biome v2.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; all code authored fresh under MIT `Copyright (c) 2026 OMEG4MON`.
- **Module boundary rules:** Feature modules NEVER import another `internal/modules/<other>` package and NEVER touch another module's SQLite tables directly. Cross-module wiring goes strictly through `internal/contracts/` and `eventbus.Bus` in `cmd/omeg4mon-router/main.go`.
- **Zero-CGO SQLite WAL & Low-RAM Budget:** Pure-Go `modernc.org/sqlite` (`db.Writer()` / `db.Reader()`), fixed-width UTC RFC3339Nano timestamps (`"2006-01-02T15:04:05.000000000Z"`), bounded retention (max 5,000 audit events / 90-day TTL; max 1,000 notifications / 30-day TTL).
- **Secret Redaction & Audit Integrity:** `security` redacts API keys, bearer tokens, passwords, and secrets from audit summaries and metadata before persisting. Clearing the audit trail inserts a fresh `security.audit_cleared` record so clearing is itself audited.
- **Streaming discipline:** `GET /api/notifications/stream` mounts on `AdminStreamAPI` (no write timeout), flushes immediately via `http.NewResponseController(w).Flush()`, emits keepalive comments, and unsubscribes cleanly on client disconnect (zero goroutine leaks).
- **Static-export frontend (`web/`):** Next.js `16.4.6` (`output: "export"`, `trailingSlash: true`), `docs/DESIGN.md` strict compliance (`#E56A4A`, `#0a0a0a` / `#FDFAF6`, Geist Sans/Mono, solid opaque surfaces with ZERO `backdrop-blur`/`backdrop-filter`, ZERO UI emojis, WCAG AA contrast).
- **Formatter, Linter & Graphify Discipline:** Every task must pass `golangci-lint fmt --diff ./...`, `golangci-lint run ./...`, Go `-race` tests, and `graphify update .` (cleaning any `graphify-out/20??-??-??/` backup directories). Final verification runs full `make check` and `make build`.

## Review Focus

1. **Concurrent job overlap guard & panic recovery:** Manually triggering a job via `POST /api/jobs/{name}/trigger` while a scheduled tick is already running must return `409 Conflict` (`ErrJobAlreadyRunning`) instead of overlapping execution; a panicking job must recover cleanly, increment `ErrorCount`, set `LastStatus = "panic"`, and emit a failure notification without crashing the server. _(Pinned in Tasks 1 & 3)_
2. **Audit secret redaction & tamper-evident clear:** `security` must redact `sk-*`, `Bearer ...`, `password`, `secret`, `token`, and `api_key` values in audit summaries and metadata before storage, and `DELETE /api/security/audit` must leave a `security.audit_cleared` marker event rather than an empty unlogged wipe. _(Pinned in Task 2)_
3. **Notification SSE stream cleanup & non-blocking fan-out:** Disconnecting an SSE client from `GET /api/notifications/stream` must exit the handler goroutine cleanly without leaking subscriptions, and slow SSE consumers must never block `POST /api/notifications` or eventbus publishers. _(Pinned in Task 3)_
4. **Disabled optional modules cleanly detach:** Disabling `security` or `notifications` via `OMEG4MON_DISABLED_MODULES` omits their routes, eventbus listeners, and periodic prune jobs cleanly without affecting core gateway traffic or startup/shutdown lifecycle. _(Pinned in Task 4)_
5. **Static-export UI compliance & token discipline:** `/dashboard/security/` and `/dashboard/notifications/` export cleanly to static HTML (`web/out/`), use semantic design tokens (`var(--status-ok)`, `var(--status-warn)`, `var(--status-error)`), contain ZERO `backdrop-blur` or UI emojis, and pass `bun test`, `bun run typecheck`, and `bun run check`. _(Pinned in Task 5)_

---

## File Map

| Layer | Files | Responsibility |
| --- | --- | --- |
| Platform Jobs & Contracts | `internal/platform/jobs/runner.go`, `internal/platform/jobs/runner_test.go`, `internal/contracts/contracts.go` | Thread-safe `JobStatus`, `Snapshot()`, `Trigger(ctx, name)`, overlap guard, observer callback, and narrow `contracts.JobInspector` interface |
| Security Module | `internal/modules/security/{migrations/0001_init.sql,repo.go,redact.go,module.go,module_test.go}`, `contracts/openapi/modules/security.yaml` | Security posture snapshot, redacted audit trail CRUD/filter/export/prune, eventbus auto-audit ingestion |
| Notifications & Jobs Module | `internal/modules/notifications/{migrations/0001_init.sql,repo.go,module.go,module_test.go}`, `contracts/openapi/modules/notifications.yaml` | System notifications inbox, unread counter, mark read/read-all/dismiss/clear, live SSE stream, eventbus alerts, `/api/jobs` status & manual trigger |
| App Wiring, Codegen & E2E | `cmd/omeg4mon-router/main.go`, `cmd/omeg4mon-router/main_test.go`, `contracts/openapi/openapi.yaml`, `web/src/shared/api/schema.d.ts` | Wire `security` & `notifications`, bridge `jobs.Runner` failures to `eventbus`, register prune jobs, OpenAPI codegen, E2E tests |
| Web Security & Notifications UI | `web/src/features/security/{schemas.ts,schemas.test.ts,hooks.ts,SecurityConsole.tsx}`, `web/src/app/(dashboard)/dashboard/security/page.tsx`, `web/src/features/notifications/{schemas.ts,schemas.test.ts,hooks.ts,NotificationsConsole.tsx}`, `web/src/app/(dashboard)/dashboard/notifications/page.tsx`, `web/src/shared/ui/AppShell.tsx`, `web/src/features/system/SystemOverview.tsx` | Security Posture & Audit studio, Notifications & Background Jobs studio, Wave 2C navigation & system module scopes |

---

### Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`

**Files:**
- Modify: `internal/platform/jobs/runner.go`
- Modify: `internal/platform/jobs/runner_test.go`
- Modify: `internal/contracts/contracts.go`

**Interfaces:**
- Consumes: `context.Context`, `slog.Logger`, `time.Duration`
- Produces:
  - `jobs.ErrJobNotFound`, `jobs.ErrJobAlreadyRunning`
  - `jobs.JobStatus{Name string, IntervalMs int64, Running bool, RunCount int64, ErrorCount int64, LastStatus string, LastError string, LastDurationMs int64, LastStartedAt *time.Time, LastFinishedAt *time.Time, NextRunAt *time.Time}`
  - `(*jobs.Runner).SetObserver(fn func(JobStatus))`
  - `(*jobs.Runner).Snapshot() []JobStatus`
  - `(*jobs.Runner).Trigger(ctx context.Context, name string) error`
  - `contracts.ErrJobNotFound`, `contracts.ErrJobAlreadyRunning`, `contracts.JobStatus`, `contracts.JobInspector`

- [ ] **Step 1: Write failing tests in `internal/platform/jobs/runner_test.go`**
  - Add `TestRunnerSnapshotAndManualTrigger` (verifies `Snapshot()` sorted output, `RunCount`, `LastStatus == "ok"`, `LastDurationMs >= 0`, `LastStartedAt`/`LastFinishedAt`, and `ErrJobNotFound`).
  - Add `TestRunnerOverlapGuardReturnsAlreadyRunning` (blocks a job on a channel, invokes `Trigger` concurrently, asserts `errors.Is(err, ErrJobAlreadyRunning)` — pinning Review Focus #1).
  - Add `TestRunnerPanicRecoveryAndObserverCallback` (panicking job recovers, increments `ErrorCount`, sets `LastStatus == "panic"`, invokes observer with panic status).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -count=1 -race ./internal/platform/jobs/...`
  - Expected: FAIL (missing `Snapshot`, `Trigger`, `SetObserver`, `ErrJobAlreadyRunning`).
- [ ] **Step 3: Implement job telemetry, overlap guard, observer, and `contracts.JobInspector`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/platform/jobs/... ./internal/contracts/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Commit**
  - Commit `feat(jobs): add job status snapshots, overlap-guarded manual trigger, and JobInspector contract`.

---

### Task 2: `security` Module (Security Posture & Redacted Audit Trail)

**Files:**
- Create: `internal/modules/security/migrations/0001_init.sql`
- Create: `internal/modules/security/redact.go`
- Create: `internal/modules/security/repo.go`
- Create: `internal/modules/security/module.go`
- Create: `internal/modules/security/module_test.go`
- Create: `contracts/openapi/modules/security.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `config.Config`, `contracts.AdminCredentialStore` (for live `RequireLogin` check), `ActiveModules func() []string`
- Produces:
  - `security.TopicAudit = "security.audit"`
  - `security.AuditInput{Category, Action, Severity, Actor, ResourceID, Summary, IPAddress string, Metadata map[string]any}`
  - `security.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger, Config config.Config, Store contracts.AdminCredentialStore, ActiveModules func() []string}`
  - `security.New(opts Options) *Module` (`Name() = "security"`, `Dependencies() = []string{"system"}`)
  - `(*Module).PruneExpired(ctx context.Context) error`
  - SQLite table (`0001_init.sql`):
    - `security_audit_events(id TEXT PRIMARY KEY, category TEXT NOT NULL, action TEXT NOT NULL, severity TEXT NOT NULL DEFAULT 'info', actor TEXT NOT NULL DEFAULT 'system', resource_id TEXT NOT NULL DEFAULT '', summary TEXT NOT NULL, metadata_json TEXT NOT NULL DEFAULT '{}', ip_address TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
  - AdminAPI routes (under `/api`):
    - `GET /security/posture`
    - `GET /security/audit` (`limit`, `category`, `severity`, `q`)
    - `POST /security/audit`
    - `DELETE /security/audit`
    - `GET /security/audit/export`

- [ ] **Step 1: Write failing tests in `internal/modules/security/module_test.go`**
  - Add `TestSecurityRedactionAndTamperEvidentClear` (pins Review Focus #2: verifies `sk-...`, `Bearer ...`, and secret metadata keys are redacted before persistence, and `DELETE /api/security/audit` leaves a `security.audit_cleared` record).
  - Add `TestSecurityEventBusAutoAuditIngestion` (publishes `auth.policy_changed`, `keys.changed`, and `security.audit` on `eventbus.Bus`, verifies persisted audit records and filter/search queries).
  - Add `TestSecurityPostureAndExportAndPrune` (verifies remote-bind warnings, master key source detection, JSON download export, and retention pruning).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -count=1 -race ./internal/modules/security/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `security` module and `contracts/openapi/modules/security.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/security/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Commit**
  - Commit `feat(security): implement security posture report, redacted audit trail, eventbus auto-audit, and OpenAPI contract`.

---

### Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control)

**Files:**
- Create: `internal/modules/notifications/migrations/0001_init.sql`
- Create: `internal/modules/notifications/repo.go`
- Create: `internal/modules/notifications/module.go`
- Create: `internal/modules/notifications/module_test.go`
- Create: `contracts/openapi/modules/notifications.yaml`

**Interfaces:**
- Consumes: `*db.DB`, `*eventbus.Bus`, `*slog.Logger`, `contracts.JobInspector`
- Produces:
  - `notifications.TopicEmit = "notifications.emit"`
  - `notifications.TopicJobFailed = "jobs.failed"`
  - `notifications.EmitPayload{Title, Message, Severity, Category, Source string, Metadata map[string]any}`
  - `notifications.Options{DB *db.DB, Bus *eventbus.Bus, Logger *slog.Logger, Jobs contracts.JobInspector}`
  - `notifications.New(opts Options) *Module` (`Name() = "notifications"`, `Dependencies() = []string{"system"}`)
  - `(*Module).PruneExpired(ctx context.Context) error`
  - SQLite table (`0001_init.sql`):
    - `system_notifications(id TEXT PRIMARY KEY, title TEXT NOT NULL, message TEXT NOT NULL, severity TEXT NOT NULL DEFAULT 'info', category TEXT NOT NULL DEFAULT 'system', source TEXT NOT NULL DEFAULT 'system', metadata_json TEXT NOT NULL DEFAULT '{}', read_at TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)`
  - AdminAPI routes (under `/api`):
    - `GET /notifications` (`unread_only`, `severity`, `limit`) -> `{ items, total, unread_count }`
    - `POST /notifications`
    - `POST /notifications/{id}/read`
    - `POST /notifications/read-all`
    - `DELETE /notifications/{id}`
    - `DELETE /notifications`
    - `GET /jobs` -> `{ jobs: [...] }`
    - `POST /jobs/{name}/trigger` -> `{ job: {...} }` (404 on not found, 409 on already running)
  - AdminStreamAPI routes (under `/api`):
    - `GET /notifications/stream` (live SSE stream of new notifications + keepalive)

- [ ] **Step 1: Write failing tests in `internal/modules/notifications/module_test.go`**
  - Add `TestNotificationsCRUDReadUnreadAndPrune` (create, list with `unread_only` & `severity`, unread count, mark single read, mark all read, delete single, clear all, prune).
  - Add `TestNotificationsEventBusAndSSEStreamCleanup` (pins Review Focus #3: subscribes to `/api/notifications/stream`, emits via `POST /api/notifications` and `eventbus`, verifies SSE delivery and clean goroutine teardown on context cancel).
  - Add `TestNotificationsJobsListAndTriggerConflict` (pins Review Focus #1: verifies `GET /api/jobs`, `POST /api/jobs/{name}/trigger` 200/404/409).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -count=1 -race ./internal/modules/notifications/...`
  - Expected: FAIL.
- [ ] **Step 3: Implement `notifications` module and `contracts/openapi/modules/notifications.yaml`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./internal/modules/notifications/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Commit**
  - Commit `feat(notifications): implement system notifications inbox, SSE stream, background jobs API, and OpenAPI contract`.

---

### Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

**Files:**
- Modify: `cmd/omeg4mon-router/main.go`
- Modify: `cmd/omeg4mon-router/main_test.go`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

**Interfaces:**
- Consumes: `security.New`, `notifications.New`, `jobs.Runner`, `contracts.JobInspector`
- Produces:
  - Active `security` and `notifications` modules registered in `lifecycle.Registry`.
  - `jobs.Runner` adapter implementing `contracts.JobInspector` and observer emitting `jobs.failed` on job error/panic.
  - Periodic jobs registered conditionally for active modules: `usage-prune`, `security-audit-prune`, `notifications-prune`.
  - Updated `contracts/openapi/openapi.yaml` and generated `web/src/shared/api/schema.d.ts`.

- [ ] **Step 1: Write failing E2E tests in `cmd/omeg4mon-router/main_test.go`**
  - Add `TestWave2CSecurityNotificationsAndJobsE2E`: verifies `/api/security/posture`, live auto-audit when creating an API key, `/api/notifications` + `/api/jobs` listing (`usage-prune`, `security-audit-prune`, `notifications-prune`), and manual `POST /api/jobs/security-audit-prune/trigger`.
  - Add `TestWave2CDisabledModulesDetachCleanly` (pins Review Focus #4: boots with `DisabledModules: {"security": true, "notifications": true}`, asserts 404 on `/api/security/posture`, `/api/notifications`, `/api/jobs`, and healthy `/api/health`).
- [ ] **Step 2: Run tests to verify failure**
  - Run: `go test -count=1 -race -run TestWave2C ./cmd/omeg4mon-router/...`
  - Expected: FAIL.
- [ ] **Step 3: Wire `cmd/omeg4mon-router/main.go`, update `contracts/openapi/openapi.yaml`, and run `make codegen`**
- [ ] **Step 4: Run tests & linter**
  - Run: `golangci-lint fmt --diff ./... && golangci-lint run ./... && go test -count=1 -race ./cmd/omeg4mon-router/...`
  - Expected: PASS, 0 lint issues.
- [ ] **Step 5: Commit**
  - Commit `feat(cmd): wire security, notifications, and background jobs observability with OpenAPI codegen and E2E tests`.

---

### Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification

**Files:**
- Create: `web/src/features/security/schemas.ts`
- Create: `web/src/features/security/schemas.test.ts`
- Create: `web/src/features/security/hooks.ts`
- Create: `web/src/features/security/SecurityConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/security/page.tsx`
- Create: `web/src/features/notifications/schemas.ts`
- Create: `web/src/features/notifications/schemas.test.ts`
- Create: `web/src/features/notifications/hooks.ts`
- Create: `web/src/features/notifications/NotificationsConsole.tsx`
- Create: `web/src/app/(dashboard)/dashboard/notifications/page.tsx`
- Modify: `web/src/shared/ui/AppShell.tsx`
- Modify: `web/src/features/system/SystemOverview.tsx`

**Interfaces:**
- Consumes: Generated OpenAPI types in `web/src/shared/api/schema.d.ts`, `apiClient` in `web/src/shared/api/client.ts`, `docs/DESIGN.md` primitives (`Card`, `Button`, `Badge`, `Alert`, `EmberRule`).
- Produces:
  - `/dashboard/security/` — Security Posture cards + warnings, filterable/searchable Audit Trail table, manual audit event form, JSON export download, and tamper-evident Clear Trail action.
  - `/dashboard/notifications/` — System Notifications inbox with live SSE badge/subscription, unread filter, severity filter, mark read / mark all read / dismiss / clear / create alert form, and Background Jobs panel with status badges, run/error counters, duration, and one-click "Run Now" trigger.
  - `AppShell.tsx` adds "Security" (`/dashboard/security/`) and "Notifications & Jobs" (`/dashboard/notifications/`), and updates kernel badge to `WAVE 2C` / `W2C`.
  - `SystemOverview.tsx` adds `security` and `notifications` to `MODULE_SCOPES`.

- [ ] **Step 1: Write failing Zod schema tests (`web/src/features/security/schemas.test.ts` and `web/src/features/notifications/schemas.test.ts`)**
- [ ] **Step 2: Run tests to verify failure**
  - Run: `cd web && bun test src/features/security/schemas.test.ts src/features/notifications/schemas.test.ts`
  - Expected: FAIL.
- [ ] **Step 3: Implement schemas, TanStack Query hooks, `SecurityConsole`, `NotificationsConsole`, pages, `AppShell.tsx`, and `SystemOverview.tsx`**
- [ ] **Step 4: Run web tests, Biome check, typecheck, and full `make check` + `make build`**
  - Verify static export contains `web/out/dashboard/security/index.html` and `web/out/dashboard/notifications/index.html`, zero `backdrop-blur`, zero UI emojis, and single-binary `bin/omeg4mon-router` builds cleanly.
- [ ] **Step 5: Update Graphify & commit**
  - Run `graphify update .` (clean any dated backup dirs) and commit `feat(web): implement Security and Notifications/Jobs studios, Wave 2C navigation, and verify single binary`.
