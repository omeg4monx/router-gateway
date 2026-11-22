# Infrastructure Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give the Next.js TypeScript and Go rewrite a repeatable source build, two CI gates, and a safe local development recipe. This foundation stage does not change app behavior; later stages are not bound to `main`'s feature set, UI, or implementation.

**Architecture:** Keep the existing Node and Go entry points and run their checks independently in the same PR workflow. Keep npm resolution locked at both package roots; loopback-only dev processes use separate temporary SQLite paths.

**Tech Stack:** Next.js, Node `22.23.4`, npm `11.20.0`, Vitest, Go `1.26.0`, GitHub Actions, SQLite.

**Spec:** `docs/specs/2026-11-22-infra-foundation-design.md`

## Global Constraints

- Work only on `infra/foundation`; `main` remains at `3be429e` for this work. Integrate into `migration/next-go` only with separate authorization. No push, release, or deployment in this plan.
- Keep Node `22.23.4` and Go `1.26.0`; use npm `11.20.0` for the supported local and CI path. Preserve direct dependency ranges, Bun lock, and Go module declarations.
- No production traffic, app source, auth, API contracts, SQLite schema, image build, or release workflow changes. No real database, credentials, or live-provider tests. Use a fresh checkout with no `.env` file for verification.
- Before edits, record baseline results on the existing branch for the current root-scoped Vitest command, Node build, and Go checks. This is a baseline for the infrastructure checks, not a parity gate for future product work. Use the toolchain above and the current install instructions until locked installs exist. Attribute failures to the baseline only when observed; baseline failures block acceptance and must not be hidden by changing or excluding tests. Do not fix unrelated failures without explicit scope approval.

## Review Focus

- Secretless PR: the test workflow cannot write to the repository or reuse checkout credentials (Task 2).
- Stale package lock: npm must reject a manifest and lock mismatch in either package root (Task 1).
- Linux fresh checkout: root Vitest must resolve project aliases and run with `RUN_REAL=0`; build must work without inherited production data (Task 2).
- Go module drift: module verification, tests, and binary build must use Go `1.26.0`, read-only resolution, and no CGO (Task 2).
- Local data isolation: both processes must report healthy without writing the other process's SQLite file or using an existing data directory (Task 3).

---

### Task 1: Track npm dependency resolution

**Files:** `.gitignore`, `tests/.gitignore`, new `package-lock.json`, new `tests/package-lock.json`.

**Interfaces:** Consumes existing `package.json`, `tests/package.json`, `.nvmrc`; produces two tracked npm lockfiles consumed by `npx --yes npm@11.20.0 ci` and `npx --yes npm@11.20.0 --prefix tests ci`.

- [ ] Record the current failure: run both supported `ci` commands in a fresh source copy of `3be429e` without generated lockfiles. Expect both to exit nonzero for a missing lockfile.
- [ ] Replace the root broad npm lock ignore with rules that allow only `/package-lock.json` and `/tests/package-lock.json`; remove the duplicate ignore from `tests/.gitignore`. Check other ignored package locks stay ignored.
- [ ] From each clean package directory, generate `package-lock.json` with `npx --yes npm@11.20.0 install --package-lock-only --ignore-scripts`. Review generated versions, integrity metadata, and sources for unexpected registries or private URLs; do not change either manifest or regenerate `bun.lock`.
- [ ] Run both `ci` commands in a clean checkout or clean copy of the changed files. Expect exit zero from both. For each package root, copy its manifest and lock to a throwaway directory and add `lock-mismatch-fixture: file:./lock-mismatch-fixture` to the manifest dependencies. Create the named local directory with a minimal `package.json` declaring that name and version `1.0.0`; leave the lock untouched. Each `npm ci` must exit nonzero for the missing lock entry. Never modify worktree manifests to test failure.
- [ ] Commit only the scoped files with eval `commit`. Explicitly select package-manager-generated lockfiles if the commit draft reports them as unrecorded.

### Task 2: Gate Node and Go in PR CI

**Review:** per-task — security (workflow permissions, action provenance, secretless builds)

**Files:** `.github/workflows/test.yml`.

**Interfaces:** Consumes the two npm locks from Task 1 and `.nvmrc`, `go.mod`, `go.sum`; produces independent `unit` and `go` jobs. Keep `unit` as the existing check identifier.

- [ ] Add `migration/next-go` to push triggers, keep pushes to `main`, `pull_request`, `workflow_dispatch`, the existing concurrency policy, and `ubuntu-latest`. Set workflow `permissions: contents: read` and `persist-credentials: false` on each checkout. Keep checkout v5 and setup-node v6; use setup-go v6. Resolve each to a full SHA from its official repository and annotate the version. Set Go module caching to use `go.sum`; leave npm caching out of this change.
- [ ] For `unit`, use `.nvmrc` as the Node version source. Run both Task 1 install commands. From the repo root, run `RUN_REAL=0 tests/node_modules/.bin/vitest run --config tests/vitest.config.js`; then run `DATA_DIR="$(mktemp -d)" npm run build`, with no production `.env` or secrets. Keep the existing job timeout.
- [ ] Add a `go` job with a 20-minute timeout. Set `CGO_ENABLED: '0'` and `GOTOOLCHAIN: local` in job env, use `go-version-file: go.mod`, and run `go mod verify`, `go test -mod=readonly ./...`, `go build -mod=readonly -o "$RUNNER_TEMP/omeg4mon-router" ./cmd/omeg4mon-router`, and `git diff --exit-code -- go.mod go.sum` as separate steps. Any failure must fail the job.
- [ ] Verify the workflow with an available Actions linter or YAML parser plus explicit inspection of permissions, persisted checkout credentials, event triggers, and pinned action references. Run both job command sequences on a clean Linux checkout with the pinned runtimes and no app credentials. Expect passing unit/build checks, verified Go modules, a passing Go suite, and an executable binary. After authorized publication, require the two check results on the integration PR; a local run is not proof that GitHub ran them.
- [ ] Commit the workflow only with eval `commit`.

### Task 3: Document and smoke-test isolated development

**Files:** `README.md` (Requirements, Quick start, Tests), `CHANGELOG.md` (unreleased development note).

**Interfaces:** Consumes both Task 1 install commands and Task 2 test commands; uses existing `npm run dev`, Go entry point `./cmd/omeg4mon-router`, and `GET /api/health` on both runtimes.

- [ ] Update the README to pin the supported npm path, use the locked installs and root-scoped test command, and distinguish live opt-in tests from default checks. Keep existing Node production, Bun, CLI, and Docker instructions compatible; only change the source quick-start install line and adjacent requirements where necessary.
- [ ] Document the dual-runtime recipe from a checkout without production `.env`: Node terminal creates `node_data=$(mktemp -d)` and starts `DATA_DIR="$node_data" INITIAL_PASSWORD=dev-node-disposable npm run dev`; Go terminal creates a different `go_data=$(mktemp -d)`, builds with `CGO_ENABLED=0 go build -mod=readonly -o "$go_data/router" ./cmd/omeg4mon-router`, and starts `HOST=127.0.0.1 PORT=8080 DATA_DIR="$go_data" INITIAL_PASSWORD=dev-go-disposable "$go_data/router"`. Note the Next dev auth-guard limitation, no routing between the processes, Ctrl-C for both, and optional cleanup of only those disposable directories once stopped.
- [ ] Execute the recipe exactly. Before startup, assert the directories differ and neither contains a database. With both up, `curl --fail --silent --show-error http://127.0.0.1:20128/api/health` must show `ok: true`; the same check at port `8080` must show `ok: true` and `engine: "golang"`. Assert the Go DB exists only under `$go_data/db/data.sqlite` for that process; if Node creates a DB, confirm it is under `$node_data/db/data.sqlite`. Stop Go and confirm Node health still succeeds; stop Node and remove only the owned disposable directories. Never stop an existing service to free a port.
- [ ] Record the infrastructure-only scope in `CHANGELOG.md`; do not assign a release version. Commit the docs with eval `commit`.

## Integration gate

Complete the whole-branch review and require Task 2 command results plus the dev smoke result. Verify only approved infra files changed and the baseline `main` ref has not moved due to this work. Report any failed or unrun checks. Roll back before integration by leaving the branch unmerged; after integration, revert the task commits on `migration/next-go` in reverse order. No database, deployment, or release rollback is part of this stage.

Subsequent stages use `main` for comparison and define their own acceptance criteria before implementation. The gateway stage records controlled measurements for gateway overhead, time to first token, p95 latency under concurrent streaming, and RSS; the dashboard stage verifies task flows and useful operational feedback. Numerical performance thresholds and UI designs are not deliverables of `infra/foundation`.
