# Auto-Update from `main` Branch Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an automatic update button to the dashboard header that appears whenever branch `main` has a newer commit and performs a verified, graceful self-update in both standalone binary mode and local git source-checkout mode.

**Architecture:** A GitHub Actions rolling workflow publishes `main-latest` pre-release binaries + `SHA256SUMS.txt` on every push to `main`. The Go `system` module exposes `GET /api/system/update` and `POST /api/system/update` (admin + loopback guarded) with hybrid auto-detection (`source` in a clean local git repo, `binary` for standalone binaries), atomic swap with `.bak` rollback, and graceful self-restart via `syscall.Exec` (Unix) / spawn+exit (Windows). The React `AppShell` renders a conditional header Update button + confirmation/progress modal with auto-reload.

**Tech Stack:** Go 1.26, Chi v5, OpenAPI 3.1 + openapi-typescript, React 19, TanStack Query, Tailwind CSS, GitHub Actions.

**Spec:** `docs/superpowers/specs/2026-11-28-auto-update-main-design.md`

## Global Constraints
- Default update repo: `omeg4monx/router-gateway`, branch: `main`, rolling release tag: `main-latest`.
- `POST /api/system/update` MUST enforce admin auth + loopback guard (`!cfg.IsLoopback() && !cfg.AllowInsecureRemote` → `403 Forbidden`).
- Local source mode MUST refuse to update when git working tree is dirty (`unavailableReason: "dirty_worktree"`).
- Binary mode MUST verify SHA256 against `SHA256SUMS.txt` and atomically swap `<exe>` with `.bak` rollback on error.
- UI button MUST be hidden when `updateAvailable === false` and `state === "idle"`.
- UI copy MUST support both `en` and `id` locales and semantic theme tokens from `web/src/app/globals.css`.

## Review Focus
- Short commit (7-char) vs full 40-char SHA comparison must match case-insensitively without false positives when already up-to-date.
- GitHub API unreachable / 403 rate-limit must not crash or break `AppShell`; `GET /api/system/update` returns `200` with `updateAvailable: false` and diagnostic `error`.
- Concurrent `POST /api/system/update` calls must return `409 Conflict` without spawning duplicate updates.
- SHA256 mismatch or invalid binary download must abort before touching the live executable and clean up temp files.
- Restart callback must only run after `App.Close` cleanly flushes SQLite WAL and stops the HTTP server.

---

### Task 1: Runtime Config & OpenAPI Contract + TypeScript Codegen

**Files:**
- Modify: `internal/platform/config/config.go`
- Modify: `contracts/openapi/modules/system.yaml`
- Modify: `contracts/openapi/openapi.yaml`
- Modify: `web/src/shared/api/client.ts`
- Test: `internal/platform/config/config_test.go`

**Interfaces:**
- Produces: `Config.UpdateRepo`, `Config.UpdateBranch`, `Config.UpdateToken`
- Produces: OpenAPI paths `/api/system/update` (`GET`, `POST`) and generated TS types `SystemUpdateStatusResponse`, `SystemUpdateRequest`.

- [ ] **Step 1:** Add failing config test for `OMEG4MON_UPDATE_REPO`, `OMEG4MON_UPDATE_BRANCH`, `OMEG4MON_UPDATE_TOKEN` defaults and overrides.
- [ ] **Step 2:** Run `go test ./internal/platform/config/...` and verify it fails.
- [ ] **Step 3:** Add update config fields and env parsing to `internal/platform/config/config.go`.
- [ ] **Step 4:** Run `go test ./internal/platform/config/...` and verify it passes.
- [ ] **Step 5:** Add `/api/system/update` (`GET`, `POST`) and schemas to `contracts/openapi/modules/system.yaml`, reference in `contracts/openapi/openapi.yaml`, export types in `web/src/shared/api/client.ts`, and run `make codegen`.

---

### Task 2: Backend Hybrid Updater & System HTTP Endpoints

**Files:**
- Create: `internal/modules/system/updater.go`
- Create: `internal/modules/system/updater_test.go`
- Modify: `internal/modules/system/module.go`

**Interfaces:**
- Consumes: `Config.UpdateRepo`, `Config.UpdateBranch`, `Config.UpdateToken`, `BuildInfo`
- Produces: `Options.Restart func(exePath string)`, `Options.HTTPClient *http.Client`, `GET /api/system/update`, `POST /api/system/update`.

- [ ] **Step 1:** Write failing unit tests in `internal/modules/system/updater_test.go` covering:
  - Commit comparison (up-to-date short vs full SHA, newer commit available).
  - Binary mode detection & `release_building` state when `main-latest` target commit lags behind `main` HEAD.
  - Binary mode download + SHA256 verification success + atomic swap + restart callback invocation.
  - Binary mode SHA256 mismatch failure leaves original executable intact and cleans temp file.
  - Source mode detection with clean vs dirty git repo (`dirty_worktree` rejection).
  - Loopback guard (`403 Forbidden` on non-loopback without `AllowInsecureRemote`) and concurrency guard (`409 Conflict`).
  - GitHub API error returns graceful `200` with `updateAvailable: false` and `error` message.
- [ ] **Step 2:** Run `go test -race ./internal/modules/system/...` and verify tests fail.
- [ ] **Step 3:** Implement `internal/modules/system/updater.go` and wire routes + `Options.Restart` / `Options.HTTPClient` in `internal/modules/system/module.go`.
- [ ] **Step 4:** Run `go test -race ./internal/modules/system/...` and verify all tests pass.

---

### Task 3: Graceful Self-Restart Wiring in `cmd/omeg4mon-router`

**Files:**
- Create: `cmd/omeg4mon-router/restart_unix.go`
- Create: `cmd/omeg4mon-router/restart_windows.go`
- Modify: `cmd/omeg4mon-router/main.go`

**Interfaces:**
- Consumes: `system.Options.Restart`
- Produces: `execReplacement(exePath string, args []string, env []string) error`

- [ ] **Step 1:** Add platform-specific `restart_unix.go` (`syscall.Exec`) and `restart_windows.go` (`exec.Command` detached spawn).
- [ ] **Step 2:** Wire `App.requestRestart(exePath string)` in `cmd/omeg4mon-router/main.go`, pass `Restart: app.requestRestart` and `HTTPClient: httpClient` into `system.New`, and execute `execReplacement` in `runMain` after `app.Run` cleanly returns.
- [ ] **Step 3:** Run `go test -race ./cmd/omeg4mon-router/...` and verify compilation and tests pass.

---

### Task 4: CI Rolling Pre-Release Workflow for `main`

**Files:**
- Create: `.github/workflows/main-rolling-release.yml`

- [ ] **Step 1:** Create `.github/workflows/main-rolling-release.yml` triggered on `push: branches: [main]` and `workflow_dispatch`, building embedded web UI + 6 zero-CGO binaries + `SHA256SUMS.txt` + `COMMIT.txt` and updating pre-release tag `main-latest`.

---

### Task 5: Frontend Update Hooks, Header Button & Progress Modal

**Files:**
- Create: `web/src/features/system/updateCopy.ts`
- Create: `web/src/features/system/updateCopy.test.ts`
- Create: `web/src/features/system/SystemUpdateControl.tsx`
- Modify: `web/src/features/system/hooks.ts`
- Modify: `web/src/shared/ui/AppShell.tsx`

**Interfaces:**
- Consumes: `SystemUpdateStatusResponse`, `SystemUpdateRequest` from `@/shared/api/client`.
- Produces: `useSystemUpdateStatus`, `useTriggerSystemUpdate`, `SystemUpdateControl`.

- [ ] **Step 1:** Write failing unit tests in `web/src/features/system/updateCopy.test.ts` testing visibility rules (`shouldShowUpdateButton`), active-state polling helper (`isUpdateInProgress`), and bilingual copy (`en`/`id`) for all states and unavailable reasons.
- [ ] **Step 2:** Run `cd web && bun test src/features/system/updateCopy.test.ts` and verify failure.
- [ ] **Step 3:** Implement `updateCopy.ts`, add `useSystemUpdateStatus` and `useTriggerSystemUpdate` to `web/src/features/system/hooks.ts`, implement `SystemUpdateControl.tsx` (header button + confirmation/progress modal + post-restart `/api/version` polling auto-reload), and mount `<SystemUpdateControl />` in `web/src/shared/ui/AppShell.tsx`.
- [ ] **Step 4:** Run `cd web && bun test && bun run typecheck && bun run check`.

---

### Task 6: Full Verification Gate & Knowledge Graph Update

- [ ] **Step 1:** Run `make check` (`fmt-check`, `lint`, `go test -race`, `CGO_ENABLED=0`, `bun test`, `typecheck`, `cli` tests).
- [ ] **Step 2:** Run `make build` to verify single-binary compilation with embedded UI.
- [ ] **Step 3:** Run `graphify update .` to keep the codebase knowledge graph current.
