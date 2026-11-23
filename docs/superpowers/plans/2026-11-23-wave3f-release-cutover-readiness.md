# Wave 3F — Release & Cutover Readiness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Finalize release engineering, cross-platform packaging, automated CI/Release pipelines, root documentation, live single-binary smoke/RAM verification, merged worktree cleanup, and cutover readiness from `migration/ts-go` to `main`.

**Architecture:** Extends the existing `Makefile`, `Dockerfile`, and `.github/workflows/ci.yml` to cover all 6 zero-CGO platform targets (`linux/darwin/windows` × `amd64/arm64`) aligned with `cli/lib/resolve.js`, adds a reusable `scripts/smoke-test.sh` harness shared by local `make smoke` and GitHub Actions, introduces a tag-triggered `.github/workflows/release.yml` (6 binaries + `SHA256SUMS.txt` + multi-arch GHCR Docker publish), and adds a concise production `README.md`.

**Tech Stack:** Go `1.26+` (`CGO_ENABLED=0`, `-tags embedui`), Bun `1.4+`, Next.js `16.4.6` (`output: "export"`), GitHub Actions, Docker Buildx (`distroless/static-debian12:nonroot`), Bash.

**Spec:** [2026-11-22-modular-ts-go-rewrite-design.md](file:///home/omeg4mon/workspace/projects/omeg4mon-router/docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership:** 0% verbatim copying from `main`; MIT `Copyright (c) 2026 OMEG4MON`.
- **Zero-CGO cross-compilation:** All 6 release binaries (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`, `windows/arm64`) compile with `CGO_ENABLED=0 -tags embedui -trimpath`.
- **Low-RAM target:** Single binary boots cleanly and reports idle RSS within the spec envelope.
- **Safety on branch cutover:** Clean up only already-merged, zero-dirty wave worktrees; do not merge `migration/ts-go` into `main` without explicit operator approval.

## Review Focus

1. **CLI ↔ Release binary naming parity:** Every filename produced by `make cross-build` and `.github/workflows/release.yml` matches `cli/lib/resolve.js` (`omeg4mon-router-<os>-<arch>[.exe]`).
2. **Smoke test isolation & cleanup:** `scripts/smoke-test.sh` uses a temporary `OMEG4MON_DATA_DIR`, traps `EXIT` to kill the child process and remove temp state, and verifies graceful `SIGTERM` exit.
3. **CI / Release workflow permissions:** `ci.yml` keeps `contents: read`; `release.yml` scopes `contents: write` and `packages: write` strictly to release jobs and skips publishing on dry-run `workflow_dispatch`.
4. **Documentation accuracy:** Every env var, CLI flag, and module name in `README.md` matches `internal/platform/config/config.go` and `internal/modules/*`.
5. **Worktree cleanup safety:** Unmerged branches/worktrees (`go/gateway`, `infra/foundation`, `migration/next-go`, `web/typescript`) remain untouched.

---

### Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`

**Files:**
- Modify: `Makefile`
- Modify: `Dockerfile`
- Create: `scripts/smoke-test.sh`

- [ ] **Step 1:** Add `cd cli && bun test` to `Makefile` `test`, expand `cross-build` to 6 targets + `dist/SHA256SUMS.txt`, and add `smoke` target.
- [ ] **Step 2:** Add OCI `LABEL` metadata to `Dockerfile`.
- [ ] **Step 3:** Create executable `scripts/smoke-test.sh` verifying `/api/health`, `/api/version`, embedded UI HTML routes, 404 JSON envelope, CLI launcher `--help`, idle RSS, and `SIGTERM` graceful shutdown.

### Task 2: CI & Tag-Triggered Release Pipelines

**Files:**
- Modify: `.github/workflows/ci.yml`
- Create: `.github/workflows/release.yml`

- [ ] **Step 1:** Update `.github/workflows/ci.yml` to run `cli` Bun tests, invoke `scripts/smoke-test.sh`, and cross-compile all 6 targets + `SHA256SUMS.txt`.
- [ ] **Step 2:** Create `.github/workflows/release.yml` for `v*` tag releases (6 binaries, `SHA256SUMS.txt`, GitHub Release, multi-arch GHCR Docker image).

### Task 3: Production `README.md`

**Files:**
- Create: `README.md`

- [ ] **Step 1:** Author concise root `README.md` covering architecture, quick start, Docker, CLI launcher, env vars/flags, 27-module catalog, and development/release commands.

### Task 4: Full Verification & Graphify Update

- [ ] **Step 1:** Run `make check`, `make build`, `make smoke`, `make cross-build`, verify `sha256sum -c dist/SHA256SUMS.txt`, and run `graphify update .`.
- [ ] **Step 2:** Commit all changes cleanly on `wave3f/release-cutover-readiness`.

### Task 5: Merged Worktree Cleanup, PR & Cutover Report

- [ ] **Step 1:** Remove the 8 clean, already-merged wave worktrees (`wave2a`–`wave3e`) and run `git worktree prune`.
- [ ] **Step 2:** Push `wave3f/release-cutover-readiness` and create PR into `migration/ts-go`.
- [ ] **Step 3:** Report full verification evidence and cutover checklist for `migration/ts-go` → `main`.
