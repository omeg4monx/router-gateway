# Auto-Update from `main` Branch Design

**Date:** 2026-11-28  
**Status:** Draft (Pending Approval)  
**Scope:** CI rolling release, Go system module (`internal/modules/system`), OpenAPI 3.1 contract, and Dashboard AppShell UI (`web/src/features/system`, `web/src/shared/ui/AppShell.tsx`)

---

## 1. Intent & Why

Operators running `omeg4mon-router` want a zero-friction way to stay current with the `main` branch directly from the dashboard: whenever a new commit lands on `main`, an **Update** button appears automatically in the top header; clicking it updates the running router and restarts it seamlessly, whether the instance is deployed as a **standalone zero-CGO binary** or running from a **local git checkout**.

## 2. Goals, Non-Goals, and Anti-Goals

### Goals
- Detect new commits on `origin/main` (GitHub `omeg4monx/router-gateway`, branch `main`) against the running build commit.
- Show an **Update** button in the dashboard header **only when an update is available** (or while an update is in progress).
- Support both deployment modes automatically:
  1. **Standalone Binary Mode:** Download the matching OS/arch binary from the `main-latest` rolling GitHub pre-release, verify SHA256 against `SHA256SUMS.txt`, atomically replace the executable, and self-restart.
  2. **Local Source Mode:** Detect when the running binary lives inside a clean local git checkout with `git`, `go`, and `bun` available; fetch & fast-forward `origin/main`, rebuild UI + Go binary, atomically replace the executable, and self-restart.
- Provide clear confirmation and live stage progress in the UI, followed by automatic page reload once the updated server is back online.

### Non-Goals
- Automatic unattended background updates without operator click.
- Downgrading to arbitrary historical commits or switching between feature branches.
- In-place container image mutation inside Docker (container deployments should detect read-only/containerized state cleanly if binary replacement is not writable).

### Anti-Goals
- Never overwrite uncommitted local changes in a git checkout (`dirty_worktree` must fail safely before touching git state).
- Never replace the running binary without checksum verification (in binary mode) or backup rollback (`.bak`).
- Never leave the server half-stopped or corrupt SQLite WAL on restart (must reuse `App.Close` graceful shutdown before exec/spawn).
- Never expose update trigger on non-loopback binds unless `OMEG4MON_ALLOW_INSECURE_REMOTE=true` (matching `/api/system/shutdown`).

## 3. Approaches Considered

| Approach | Summary | Trade-offs | Decision |
| :--- | :--- | :--- | :--- |
| **A. Hybrid Auto-Detect + Rolling `main-latest` CI** | Add a GitHub Actions workflow publishing `main-latest` pre-release binaries on push to `main`. Backend auto-detects `binary` vs `source` mode, performs verified atomic swap + graceful self-restart. Header button opens a compact confirmation + progress dialog. | Full support for both standalone binary and local source checkouts on every `main` commit; adds one CI workflow on push to `main`. | **Selected** |
| **B. Source-Local for `main`, Binary for `v*` Tags Only** | Local checkout pulls `main`; standalone binary only updates from semver `v*` releases. | Saves CI minutes, but standalone binary users cannot update to `main` commits before a tag is cut. | Rejected |
| **C. Periodic Background Job + Immediate One-Click** | Background job polls `main` and clicking the header button restarts immediately without confirmation. | Faster by one click, but accidental clicks interrupt live gateway traffic. | Rejected |

## 4. Architecture & Contracts

### 4.1 CI Rolling Pre-Release (`.github/workflows/main-rolling-release.yml`)
- Trigger: `push` to `branches: [main]` and `workflow_dispatch`.
- Concurrency: `group: main-rolling-release`, `cancel-in-progress: true`.
- Builds embedded UI (`web/out` → `cmd/omeg4mon-router/ui_dist`) and all 6 zero-CGO binaries (`linux/darwin/windows` × `amd64/arm64`) with `-ldflags` injecting `main.version=main-<shortsha>`, `main.commit=<shortsha>`, `main.buildTime=<utc>`.
- Generates `SHA256SUMS.txt` and `COMMIT.txt` (full 40-char `GITHUB_SHA`).
- Publishes/overwrites GitHub pre-release tag `main-latest` (`prerelease: true`, `make_latest: false`) pointing to `github.sha`.

### 4.2 Backend Updater (`internal/modules/system/updater.go`)
- **Configuration** (`internal/platform/config/config.go`, optional overrides with sane defaults):
  - `OMEG4MON_UPDATE_REPO` (default: `omeg4monx/router-gateway`)
  - `OMEG4MON_UPDATE_BRANCH` (default: `main`)
  - `OMEG4MON_UPDATE_TOKEN` (optional GitHub token for private repos / higher rate limit)
- **Check (`GET /api/system/update?force=1`)**:
  - Protected under `AdminAPI`.
  - Caches remote check for 5 minutes in memory unless `force=1`.
  - Queries GitHub API `GET /repos/{repo}/commits/{branch}` for latest commit SHA, first-line commit message, author/date, and HTML URL.
  - Compares against current running commit (`BuildInfo.Commit`, falling back to `git rev-parse --short HEAD` when `Commit` is `"none"` inside a local repo) using case-insensitive prefix match.
  - Resolves `mode` and `canApply`:
    - Detects if running inside a local git checkout (`detectLocalRepo(exePath)`) with `git`, `go`, and `bun` available.
    - When inside a local git checkout: selects `source` mode (checks clean working tree; if dirty, sets `canApply: false`, `unavailableReason: "dirty_worktree"`).
    - When standalone (not in a git checkout): checks `GET /repos/{repo}/releases/tags/main-latest` for matching `omeg4mon-router-<goos>-<goarch>[.exe]` + `SHA256SUMS.txt` and confirms the release target commit matches `latestCommit`. If CI is still building for the latest commit, returns `updateAvailable: true`, `canApply: false`, `unavailableReason: "release_building"`.
- **Apply (`POST /api/system/update`)**:
  - Protected under `AdminAPI` + loopback guard (`!cfg.IsLoopback() && !cfg.AllowInsecureRemote` → `403 Forbidden`).
  - Rejects with `409 Conflict` if an update is already running or `canApply` is false.
  - Runs asynchronously with mutex-guarded state transitions:
    `idle` → `checking` → `downloading` (binary) / `building` (source) → `verifying` → `applying` → `restarting` (or `failed` with error message).
  - **Binary Mode:** downloads asset + `SHA256SUMS.txt` using the module HTTP client (bounded size limit), verifies SHA256 digest, writes to `<exeDir>/.omeg4mon-router.update.<pid>.tmp` (`0755`), atomically swaps `<exe>` → `<exe>.bak` and `<tmp>` → `<exe>` (restoring `.bak` on failure).
  - **Source Mode:** runs `git fetch origin main`, verifies clean tree and fast-forwardability, runs `git checkout main && git pull --ff-only origin main` (or `git merge --ff-only origin/main` if already on `main`), builds web (`bun run build` + copy to `cmd/omeg4mon-router/ui_dist`), builds Go binary into temp file with injected `version`/`commit`/`buildTime`, atomically swaps target binary.
  - **Graceful Self-Restart (`cmd/omeg4mon-router/main.go`):**
    - `system.Options` receives a `Restart func(exePath string)` callback wired in `main.go`.
    - `App` records `restartExePath` and triggers `requestShutdown()`.
    - After `app.Run()` cleanly closes HTTP server, jobs, modules, and SQLite WAL, `runMain()` executes the replacement binary:
      - Unix (`linux`, `darwin`): `syscall.Exec(exePath, os.Args, os.Environ())` (preserves PID and supervisor integration).
      - Windows (`windows`): spawns detached child process with `os.Args` and returns exit code `0`.

### 4.3 OpenAPI Contract (`contracts/openapi/modules/system.yaml` & `openapi.yaml`)
- Add paths:
  - `GET /api/system/update` (`getSystemUpdateStatus`)
  - `POST /api/system/update` (`triggerSystemUpdate`)
- Add schemas: `SystemUpdateStatusResponse`, `SystemUpdateRequest`.
- Regenerate TypeScript types via `make codegen` (`cd web && bun run gen:api`).

### 4.4 Frontend UI (`web/src/features/system/` & `web/src/shared/ui/AppShell.tsx`)
- **Hooks (`web/src/features/system/hooks.ts`):**
  - `useSystemUpdateStatus(enabled)`: queries `GET /api/system/update`, caches 5 min while idle, switches to 1.5s polling while `state` is active (`checking`/`downloading`/`building`/`verifying`/`applying`/`restarting`).
  - `useTriggerSystemUpdate()`: mutation calling `POST /api/system/update`.
- **Component (`web/src/features/system/SystemUpdateControl.tsx`):**
  - Renders an **Update** button in `AppShell.tsx` header (adjacent to `PwaInstallButton`) **only when `updateAvailable === true` or an update is actively running**.
  - Matches `PwaInstallButton` sizing (`h-9`), Omegamon Zwart / X-Antibody token styling, and bilingual copy (`en` / `id`).
  - Clicking opens a confirmation & live progress modal (matching the Sign-Out modal pattern in `AppShell.tsx`):
    - Shows current commit → target `main` commit, commit title, and update mode badge (`Binary` / `Source`).
    - If `canApply === false`, explains why (`dirty_worktree`, `release_building`, `not_writable`) and provides a "Check Again" button (`?force=1`).
    - During update, shows live stage status and disables dismiss until finished or failed.
    - When `state === "restarting"`, polls `GET /api/version` until the server responds with the updated commit (or recovers after brief downtime) and automatically calls `window.location.reload()`.

## 5. Failure Modes & Verification Strategy

- **GitHub API unreachable / rate-limited:** `GET /api/system/update` returns `updateAvailable: false` with non-fatal diagnostic metadata rather than breaking the dashboard shell.
- **Checksum mismatch / truncated download / build error:** Temp file is removed, running executable is untouched, state transitions to `failed` with exact error message, and user can retry.
- **Atomic swap failure:** `.bak` is restored immediately before reporting `failed`.
- **Verification:**
  - Go unit tests (`internal/modules/system/updater_test.go`) covering commit comparison, GitHub API mock responses, SHA256 verification pass/fail, atomic swap + rollback, loopback guard, dirty worktree rejection, and concurrency guard (`409`).
  - Web unit/component tests (`web/src/features/system/SystemUpdateControl.test.tsx`) covering hidden when up-to-date, visible when update available, confirmation modal, progress states, and error rendering.
  - Full `make check` (`fmt-check`, `lint`, `go test -race`, `CGO_ENABLED=0`, `bun test`, `typecheck`) and `make build`.

## 6. Left to the Builder
- Exact icon choice from `lucide-react` (`ArrowUpCircle` / `DownloadCloud` / `RefreshCw`) consistent with `AppShell.tsx`.
- Platform-specific restart helper split into `restart_unix.go` and `restart_windows.go` under `cmd/omeg4mon-router/`.
