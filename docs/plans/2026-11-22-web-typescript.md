# Web TypeScript Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Bun mandatory for web source development and CI, pin Node v26, and ship one strict TypeScript slice without porting the legacy backend.

**Architecture:** Bun owns installs and runs package scripts; Node executes the Next.js app, standalone custom server, and Vitest. Next.js App Router runs typed UI beside legacy JS; pilot changes are local to `/dashboard/endpoint`.

**Tech Stack:** Next.js 16, React 19, TypeScript 5.9 line, Bun `1.4.5`, Node.js `26.13.1` LTS, Vitest, GitHub Actions.

**Spec:** `docs/specs/2026-11-22-web-typescript-design.md`

## Global Constraints

- Implement only on `web/typescript`, based on `migration/next-go` at `1b1182655f`; never edit or merge into `main`. Keep the existing Go CI job unchanged.
- Bun `1.4.5` is mandatory for app/test install and scripts; Node `26.13.1` is mandatory for the Next.js/Vitest runtime. Recheck both official stable lines at implementation start as described in the spec; do not use floating tags, `bun --bun`, or an npm fallback in the migrated source/CI path.
- Keep Docker publishing, global CLI packaging, backend JS, API/session behavior and database schema out of this change. The published migration branch still uses the foundation npm workflow until the Bun cutover is implemented and integrated.
- Use a checkout without production `.env`, disposable `DATA_DIR`, and `RUN_REAL=0`. No live upstream credentials, remote tunnel operations, or release actions during verification.

## Review Focus

- Missing Bun or wrong major Node must not silently use the old workflow. Verify the supported commands and report the actual executing runtime (Task 1).
- Manifest/lock mismatch at either package root must fail before running tests (Task 1).
- Next output must retain the existing standalone assets/custom-server layout and start under Node, not just compile (Task 1).
- A real TypeScript type error must fail CI checks; existing alias resolution must work in fresh and existing checkouts (Task 2).
- Pilot status must remain understandable without color; unreachable health must not be reported online, and no auth or real tunnel action is needed for QA (Task 3).

### Task 1: Cut over source tooling and CI to Bun

**Review:** per-task — security (dependency resolution, CI action pins, and native install scripts)

**Files:** `.nvmrc`, `package.json`, `tests/package.json`, `bun.lock`, new `tests/bun.lock`, remove root and tests `package-lock.json`, `.gitignore`, `.github/workflows/test.yml`, and `README.md`.

**Interfaces:** Preserve app commands `bun run dev`, `bun run build`, `bun run start`; add `bun run test` at repo root. Runtime proof: `process.version` matches `.nvmrc` and `process.versions.bun` is unset within the Node process.

- [ ] Record the current Node v26 baseline: root-scoped Vitest, Next production build, and `/api/health` with an isolated loopback production start. Confirm a clean frozen install cannot yet cover both Bun package roots. Verify official Bun and Node release versions before changes.
- [ ] Pin Bun through root `packageManager: "bun@1.4.5"` and setup-bun CI input, keep `.nvmrc` as the exact Node source, and set CI permissions/read-only checkout as today. Pin setup-bun to a verified official commit SHA; no other CI action or Go-job changes.
- [ ] Regenerate root and test Bun locks with the pinned Bun, preserving manifest dependency ranges and native optional platform packages. Review graph differences against the active npm locks before deleting them; do not inherit stale transitive versions silently. Audit the install scripts Bun blocks; add `trustedDependencies` only for named scripts demonstrably needed by the locked native/build path, not blanket trust. Remove the old `:bun` aliases that force Bun runtime, update the source-only `build:fresh` command, and remove the two npm lock exceptions. Leave the separate CLI npm scripts intact. Use explicit `node` invocations for Next/Vitest entry points where runtime would otherwise be ambiguous, preserve the existing `--webpack` build and ensure `postbuild` staging runs once.
- [ ] Run `bun install --frozen-lockfile` at root and `(cd tests && bun install --frozen-lockfile)` for tests. In disposable copies add a dependency missing from each lock and assert both installs fail. Check in CI that each lock exists before running the frozen install: Bun can otherwise resolve dependencies when no lock exists. Confirm installs leave locks unchanged. From the root, CI runs `RUN_REAL=0 bun run test` and `DATA_DIR="$(mktemp -d)" bun run build`. Test script executes `tests/node_modules/vitest/vitest.mjs` with `--config tests/vitest.config.js` under Node.
- [ ] From the source checkout, run `bun run start` with Node v26, a fresh data directory, and a loopback port; the root custom wrapper loads `.next/standalone/server.js` and locates source guard code. Starting the traced wrapper directly fails closed on the baseline because the guard is not traced, and independent packaging is out of scope. Verify `/api/health`, production login, a staged local static asset, and the SQLite native binding. Update the README with Bun commands, exact Node pin, and Docker/CLI carveout. Commit via eval `commit`.

### Task 2: Enable strict mixed JS/TS checks

**Files:** create `tsconfig.json`, remove `jsconfig.json`, update `package.json`, `bun.lock`, `.github/workflows/test.yml`, `README.md`.

**Interfaces:** `bun run typecheck` runs Next route type generation then `tsc --noEmit` using Node v26. Aliases remain `@/* -> ./src/*` and `open-sse` / `open-sse/*` unchanged.

- [ ] Read the installed Next.js TypeScript guide. Adopt its bundler-resolution JSX settings, Next TS plugin and generated declaration includes (including the dev directory when generated), add `allowJs: true`, `checkJs: false`, `strict: true`, `noEmit: true`, and migrate `jsconfig.json` aliases into the single `tsconfig.json`. Do not commit generated `.next*` or `next-env.d.ts`. Keep TypeScript at a Next-verified patch in the 5.9 line; update Node 26 and React/DOM 19 ambient types with Bun. No unrelated library upgrade.
- [ ] Add `typecheck` using `next typegen` then the pinned TypeScript compiler and make it a required step in the existing unit job before build. Prove the gate fails by adding a temporary included `src` `.ts` file assigning a string to a `number`, seeing the compiler error, then removing only that fixture. Prove typecheck from a fresh checkout without generated declarations, rerun after a dev start, and run the production build. Verify a JS route and the typed pilot can coexist without `ignoreBuildErrors` or global declaration shims. Commit via eval `commit`.

The production Webpack build also validates legacy JS route exports. It rejected non-route exports in the plugins page and the pricing, auto-scan, usage and count-tokens handlers. Keep private values private, remove unused `GET_DEFAULTS` (not a registered route), and move the existing credential refresh helper verbatim to a sibling JS module with all its consumers updated. This does not migrate backend code or change HTTP routes; keep validation enabled.

### Task 3: Migrate a gateway status slice and record next waves

**Files:** under `src/app/(dashboard)/dashboard/endpoint/`, rename `page.js` to `page.tsx`, `endpointPing.js` to `endpointPing.ts`, and `components/StatusAlert.js` to `components/StatusAlert.tsx`; update only affected importers and tests. Add `tests/unit/endpoint-ping.test.js` if no existing behavioral coverage covers the probe. Update `CHANGELOG.md` with the small initial scope.

**Interfaces:** `clientPingUrl(url: string | null | undefined): Promise<boolean>`; `clientPingAny(...urls: (string | null | undefined)[]): Promise<boolean>`; status props `{ status: { type: 'success' | 'warning' | 'info' | 'error'; message: string }; className?: string }`. Keep `machineId` in the existing server-to-client prop contract. Do not touch the key or tunnel mutation paths.

- [ ] Pin probe behavior in Vitest with a local HTTP responder: HTTP 200 resolves `true`, HTTP 503 resolves `false`, no URL resolves `false`, and `clientPingAny` resolves `true` if at least one endpoint answers successfully. Run against the existing JS before renaming to establish expected behavior.
- [ ] Port the three files with strict types and update only real consumers. Keep the status markup, link safety, theme tokens and current request timeout. Use typed nullable URL input instead of unchecked casts. Do not convert the large `EndpointPageClient.js` or unrelated shared components.
- [ ] Verify typecheck, focused probe tests, root Vitest, and production build. In the browser, use a disposable loopback instance: open `/dashboard/endpoint`, verify online/offline status text with the health path available/unavailable, keyboard reachability and copying of the displayed local URL, and no clipping at mobile and desktop widths. Exercise status links with a disposable local fixture only; do not enable a real tunnel or expose API keys. Record QA and commit via eval `commit`.
- [ ] Inventory remaining JSX-bearing frontend files from tracked paths. Split follow-on slices into dashboard shell/shared UI, provider/config screens, usage/observability, then remaining dashboard and public screens. For each slice record route paths, client contracts, and verification flows before conversion. No claim that this pilot completes the whole frontend migration.

### Follow-on slice inventory

The baseline has 33 grouped dashboard page modules, seven other App Router page modules, and 65 JavaScript files under `src/shared/components`. This stage converts one dashboard page entry and two local files. The remaining candidates start with these contracts; enumerate each slice's full file list before changing it:

| Slice | Current routes or entry points | Client contract and verification |
| --- | --- | --- |
| Dashboard shell and shared UI | `src/app/(dashboard)/layout.js`, `/dashboard` | Navigation, theme, focus, and mobile menu; verify keyboard travel and responsive layout. |
| Provider and configuration UI | `/dashboard/providers`, `/dashboard/connection`, `/dashboard/combos` | Provider lists, form payloads and validation; verify empty, validation-error, and success states with disposable data. |
| Usage and observability | `/dashboard/usage`, `/dashboard/api-key-usage`, `/dashboard/console-log` | Usage filters, chart/empty states and log limits; verify data refresh, failures and narrow layouts. |
| Remaining dashboard and public UI | `/dashboard/endpoint` remainder, `/dashboard/basic-chat`, `/login`, `/share/[token]`, `/claim`, other tracked `page.js` entries | One workflow per conversion; verify streaming/cancellation or auth and expired-link cases where applicable. Keep server-side policy in its current backend until the Go stage. |

## Handoff gate

Run the whole-branch review and all task checks; merge only into `migration/next-go` after explicit approval and reported PR checks. No route parity with `main` is mandated beyond the deliberately unchanged pilot contracts. Rollback is reverting the web tooling and TS commits on the migration branch; no data migration or production rollback is required.
