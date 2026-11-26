---
name: omeg4mon-versioning
description: >-
  Use when finishing any task, reviewing diffs before commit, auditing version
  drift, or deciding whether changes in omeg4mon-router require no bump, patch,
  minor, major, or pre-release SemVer increments.
---

# omeg4mon-router Versioning Verificator

Verify SemVer impact at the end of every task so version decisions stay accurate, disciplined, and free of premature bumps or hidden breaking changes.

**Core principle:** Version increments measure **external contract compatibility**, never diff size, refactor effort, or task count. Per-task verification **classifies and reports** impact; it does not spam version bumps or git tags on every task.

## 1. Gather Evidence First

Before declaring any verdict, inspect the repository's actual version state:

```bash
git describe --tags --abbrev=0 2>/dev/null || echo "untagged"
git log "$(git describe --tags --abbrev=0 2>/dev/null)..HEAD" --oneline
git status --short
```

Read the three canonical version manifests and compare them against the latest `vX.Y.Z` tag:

| Artifact | Field | Role |
| --- | --- | --- |
| Git tag `vX.Y.Z` | Annotated tag | Source of truth for `Makefile` (`-X main.version`), `/api/version`, GitHub Releases, and GHCR Docker tags |
| `cli/package.json` | `"version"` | Published npm launcher `@omeg4mon/router` |
| `web/package.json` | `"version"` | Embedded Vite React frontend package |
| `contracts/openapi/openapi.yaml` | `info.version` | Control-plane OpenAPI 3.1 contract (`make codegen` source) |

Never hardcode release versions in Go source (`cmd/omeg4mon-router/main.go` receives `main.version` via `-ldflags`).

## 2. Public Contract Surfaces

A change affects versioning **only** if it alters a shipped public surface:

- **Public surfaces (version-relevant):**
  1. **Gateway wire API:** `/v1/*` and provider-compatible proxy endpoints, headers, SSE framing, error envelopes.
  2. **Control-plane API:** `contracts/openapi/openapi.yaml` & `contracts/openapi/modules/*.yaml` routes, request/response schemas, status codes.
  3. **Operator config & CLI:** `OMEG4MON_*` env vars, CLI flags, `cli/bin/omeg4mon-router.js`, exit codes, `cli/package.json` `engines.node`, Docker port/volume contract.
  4. **Persistent state:** SQLite schema/migrations (`internal/platform/db/`), encrypted vault/config formats, upgrade/rollback compatibility.
  5. **Shipped Web UI (`web/src/`):** User-facing features, workflows, accessibility/theme fixes, runtime behavior.
- **Non-public surfaces (`NONE`):** `README.md`, `docs/`, comments, unit/integration tests (`*_test.go`, `*.test.ts`), CI/workflows, `.agents/`, linters, devDependencies, and internal Go/TS refactors with zero observable runtime or contract change.

## 3. SemVer Decision Matrix

Classify every task using the **highest applicable row**:

| Change Kind | Current `0.y.z` (Pre-1.0) | Current `>= 1.0.0` (Stable) | Conventional Commit |
| --- | --- | --- | --- |
| Non-public only (docs, tests, CI, internal refactor, devDeps) | `NONE` | `NONE` | `docs:`, `test:`, `chore:`, `ci:`, `refactor:` |
| Backward-compatible bugfix, security hardening, memory/leak fix, perf fix, runtime dep fix, UI bug/contrast fix | `PATCH` (`0.Y.Z+1`) | `PATCH` (`X.Y.Z+1`) | `fix:`, `perf:` |
| Backward-compatible new feature: new optional endpoint/field, new provider adapter, new CLI flag, new UI page/capability, additive DB migration | `MINOR` (`0.Y+1.0`) | `MINOR` (`X.Y+1.0`) | `feat:` |
| **Breaking change** on any public surface: removed/renamed route or JSON field, stricter validation, renamed/removed `OMEG4MON_*` env var without fallback, destructive/non-migrating DB change, dropped Node/runtime support | `MINOR` (`0.Y+1.0`) + explicit `BREAKING CHANGE` | `MAJOR` (`X+1.0.0`) | `feat!:` / `refactor!:` + `BREAKING CHANGE:` |
| Explicit GA stability commitment (public API, config, and storage frozen by operator) | `MAJOR` (`1.0.0`) — **operator approval required** | N/A | `chore(release): v1.0.0` |
| High-risk overhaul needing soak testing before GA (major rewrite, risky migration, breaking API preview) | `PRE-RELEASE` (`0.Y+1.0-rc.1` or `1.0.0-rc.1`) | `PRE-RELEASE` (`X+1.0.0-rc.1`) | Tag `v...-rc.N` |

**Pre-1.0 discipline:** At `0.y.z`, breaking changes bump `MINOR` (`0.2.0` → `0.3.0`) and **must** be marked `BREAKING CHANGE`. Never auto-promote `0.y.z` → `1.0.0` because a diff is large or breaking; `1.0.0` is an explicit operator stability decision.

## 4. Per-Task vs Release-Cut Rules

1. **Compute two verdicts every task:**
   - **Task Impact:** Bump required by *this task's diff alone*.
   - **Cumulative Unreleased Impact:** `max(NONE < PATCH < MINOR < MAJOR)` across all commits in `<latest-tag>..HEAD` plus current diff.
2. **Normal task completion (default):**
   - Report the verdict only. **Do NOT** edit `cli/package.json`, `web/package.json`, `contracts/openapi/openapi.yaml`, or run `git tag`.
   - Why: Pushing any `v*` tag triggers `.github/workflows/release.yml` (6 cross-compiled binaries, GitHub Release, multi-arch GHCR Docker push). Per-task tagging creates noisy releases and wastes CI.
3. **Urgent hotfix exception:** If a `PATCH` security/crash fix must ship ahead of unreleased `MINOR` work, recommend cutting a hotfix branch from `<latest-tag>` for `vX.Y.(Z+1)`.
4. **Release cut (only when the operator explicitly asks to cut a release):**
   - Update **all three** manifests (`cli/package.json`, `web/package.json`, `contracts/openapi/openapi.yaml`) to the exact same `X.Y.Z`.
   - If `contracts/openapi/openapi.yaml` changed, run `make codegen` so `web/src/shared/api/schema.d.ts` stays clean.
   - Run `make check` (and `make build && make smoke` before tagging).
   - Commit `chore(release): vX.Y.Z` first, verify clean `git status`, then create annotated tag `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.

## 5. Rationalizations & Red Flags

| Excuse | Reality |
| --- | --- |
| "Big refactor / huge PR, let's jump to `v1.0.0`" | Diff size ≠ stability guarantee. Additive work is `MINOR`; `1.0.0` requires explicit GA approval. |
| "Just internal cleanup, rename endpoint/field/env var is `NONE` or `PATCH`" | Any renamed route, JSON key, or `OMEG4MON_*` env var breaks callers/deployments. Always `BREAKING`. |
| "Docs/tests/lint fix should bump `PATCH`" | No shipped runtime behavior changed. Verdict is `NONE`. |
| "Let's bump & tag after every task to stay tidy" | Accumulate tasks until release cut; use `max` impact since `<latest-tag>`. |
| "OpenAPI `info.version` or `web/package.json` is optional" | Drift causes confusion across CLI, UI, API spec, and Docker labels. Sync all 3 manifests on release cut. |

## 6. Required Exit Report

Always end versioning verification with this exact table (evidence-backed, no guessing):

| Check | Result | Evidence / Notes |
| --- | --- | --- |
| Latest Release Tag | `vX.Y.Z` | Output of `git describe --tags --abbrev=0` |
| Manifest Sync | `SYNCED` / `DRIFT` | `cli/package.json`, `web/package.json`, `contracts/openapi/openapi.yaml` versions |
| Task Impact | `NONE` / `PATCH` / `MINOR` / `MAJOR` (+ `BREAKING` if applicable) | Public surfaces touched by this task |
| Cumulative Unreleased Impact | `NONE` / `PATCH` / `MINOR` / `MAJOR` | Highest impact across `<tag>..HEAD` + current diff |
| Next Target Version | `vX.Y.Z` (or `N/A`) | Derived from Latest Tag + Cumulative Impact |
| Release Action | `Defer to release` / `Cut hotfix` / `Ready to release` | Never tag or bump manifests without explicit release intent |
