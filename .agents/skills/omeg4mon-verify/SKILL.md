---
name: omeg4mon-verify
description: >-
  Project-specific verification gate for omeg4mon-router covering Go resource/memory safety,
  OpenAPI contract sync, embedded UI bundle hygiene, React 19/TanStack Query cleanup,
  component hierarchy, naming clarity (files/folders/variables), zero dead code,
  Omegamon Zwart/X-Antibody design-system consistency, and responsive UI/UX checks.
  Activate before completing any non-trivial backend, CLI, contract, or frontend change.
---

# omeg4mon-router Verification Gate

Run this gate before claiming work complete. It complements global `verification-before-completion`, `ponytail-antigravity-core`, and `ui-quality-verification` by checking only what is specific to `omeg4mon-router` and not already caught by linters.

## 1. Scope Detection & Deterministic Gate

1. Inspect changed files (`git status --short` / `git diff --name-only`) to determine active tracks:
   - **Backend/Contract track:** `cmd/`, `internal/`, `cli/`, `contracts/`
   - **Web/UI track:** `web/`
2. Run deterministic checks via `Makefile`:
   - Always run `make check` (`gofmt`, `golangci-lint`, `biome check`, `go test -race`, `CGO_ENABLED=0` db tests, `bun test`, `tsc --noEmit`, `cli` tests).
   - When `golangci-lint` reports issues, run `golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...` so capped `unused` / `misspell` / `gocritic` findings are not hidden.
   - If `contracts/openapi/` changed: run `make codegen` and verify `web/src/shared/api/schema.d.ts` is updated and clean.
   - If `web/` changed materially (new deps, assets, or routes): run `cd web && bun run build` and inspect `web/out` bundle size so `-tags embedui` binary stays lean.

## 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always)

Apply to every changed file across `cmd/`, `internal/`, `cli/`, `contracts/`, and `web/`:

- **Unambiguous file & folder naming:**
  - Every file and directory name reflects a single clear responsibility and follows repo layout (`internal/modules/<domain>/`, `web/src/features/<domain>/`, `web/src/shared/ui/`).
  - Reject vague catch-all or versioned names (`misc/`, `helpers.go`, `utils2.ts`, `temp/`, `new_*`, `old_*`, `v2_*` without an active v1 contract).
- **Unambiguous variable, function & type naming:**
  - Use explicit, domain-accurate names (`strategy`, `connectionID`, `healthStatus`, `errorMessage`) instead of cryptic abbreviations (`strat`, `tmp`, `val`, `res2`, `cfg2`, `d`, `x` outside idiomatic short loop indices or Go receivers).
  - Names must match actual behavior and units (e.g., `timeoutMs` vs `timeoutSec`, `isEnabled` for booleans); no shadowed or contradictory names.
- **Zero dead code & unused artifacts:**
  - Remove unused functions, methods, types, constants, struct fields, exports, React components, hooks, routes, CSS rules, and orphaned files left behind after refactors or migrations.
  - Remove unreachable branches, commented-out code blocks, stale `//nolint` or `biome-ignore` suppressions, and unused dependencies in `go.mod` / `package.json`.
  - For `web/` changes, verify newly added or edited exported symbols/files actually have live call sites (not just passing `biome check`, which does not catch unreferenced cross-file exports).

## 3. Backend & CLI Audit (when Go / CLI / Contracts changed)

Focus on failure modes that `golangci-lint` and `go test -race` do **not** catch:

- **Goroutine & lifecycle safety:** Every spawned goroutine terminates on `ctx.Done()` or `internal/platform/lifecycle` shutdown. No fire-and-forget goroutines holding references indefinitely.
- **Resource & memory leak hygiene:**
  - `sql.Rows` always closed (`defer rows.Close()`) AND `rows.Err()` checked.
  - `time.Ticker` stopped; `time.After` not used inside hot loops.
  - Streaming/SSE/proxy handlers stop reading/writing immediately when request `ctx` cancels.
  - Caches, buffers, and in-memory maps/slices have explicit bounds or eviction.
- **Error clarity & safety:** Errors use `ErrorEnvelope` at HTTP boundaries, wrap internal causes with `%w`, give actionable context to callers, and never leak credentials/tokens.
- **Zero overengineering:** Reuses `internal/platform/*` utilities; no speculative interfaces, redundant wrappers, or dead code.

## 4. Web & UI/UX Audit (when `web/` changed)

### 4.1 Code, Type Safety & Memory Safety

- **End-to-end typesafe:** API calls use `apiClient` (`src/shared/api/schema.d.ts`) and Zod schemas; zero `any`, zero unsafe `as unknown as` bypasses.
- **Anti memory leak (React 19 + TanStack Query):**
  - Every `useEffect` cleans up `addEventListener`, timers (`setTimeout`/`setInterval`), `AbortController`, and `EventSource`/subscriptions.
  - TanStack Query mutations invalidate exact query keys; no unbounded arrays or leaked closures in Zustand/component state.
- **Reusable component hierarchy:**
  - `pages/` stays thin; domain logic lives in `features/<domain>/` (`hooks.ts`, `schemas.ts`, `*Console.tsx`); generic primitives live in `shared/ui/`.
  - Reuse existing `shared/ui/` primitives (`Button`, `Card`, `Alert`, `Badge`, `Input`, `Label`, `ErrorBoundary`, `EmberRule`) instead of reinventing styles.
- **Lightweight & zero bloat:** No redundant npm packages, no unnecessary re-renders on large lists, no heavy inline assets.

### 4.2 Design System, Aesthetics & UX Quality

- **Design tokens & anti AI-slop:**
  - All colors/surfaces/borders use semantic variables from `web/src/app/globals.css` (no hardcoded hex, no generic neon/gradient/glassmorphism slop).
  - Verify both Dark (_Omegamon Zwart_) and Light (_Omegamon X-Antibody_) themes remain high-contrast and cohesive (`palette.test.ts` passes).
  - Preserve `en`/`id` locale support where copy is added or modified.
- **Layout clarity & responsiveness:**
  - Information hierarchy is direct and scannable — no convoluted nested cards, redundant chrome, or cluttered controls.
  - Proportional layout across mobile (`375px`), tablet (`768px`), and desktop (`1280px+`) with zero unintended horizontal overflow, clipped text, or unreachable actions.
  - Standard CSS/Tailwind features with broad cross-browser compatibility (Chromium, Firefox, Safari/WebKit).
- **Complete interaction states:**
  - **Loading:** Non-blocking skeleton or inline/button pending indicators (avoid full-screen blockers outside auth/init bootstrap).
  - **Error:** Uses `extractErrorMessage(err)` inside `Alert`/inline feedback with clear cause + recovery action (retry/dismiss/reload).
  - **Empty & guidance:** Clear empty states, concise helper text/labels so usage is self-evident.
  - **Accessibility:** Semantic HTML, visible keyboard focus, proper labels, clickable targets.

### 4.3 Runtime Visual Verification

- When UI layout, styling, or interaction changed and a live dev server (`127.0.0.1:20127`) + browser inspection tool are available, verify representative mobile and desktop viewports plus light/dark states.
- If browser/runtime verification cannot be executed, explicitly report `NOT VERIFIED` with the exact reason — never claim visual pass from static code alone.

## 5. Exit Report

Report concisely using evidence only:

| Check                                              | Status                           | Evidence / Notes                          |
| -------------------------------------------------- | -------------------------------- | ----------------------------------------- |
| Deterministic (`make check` / `codegen` / `build`) | PASS / FAIL                      | Exact commands run & exit result          |
| Naming Clarity & Zero Dead Code                    | PASS / FAIL                      | Files/folders/symbols & unused/orphan audit |
| Backend Memory & Robustness                        | PASS / FAIL / N/A                | Specific files & resource paths audited   |
| Frontend Types, Memory & Hierarchy                 | PASS / FAIL / N/A                | Hooks, cleanup, reuse, bundle impact      |
| UI/UX Design, States & Responsive                  | PASS / FAIL / NOT VERIFIED / N/A | Tokens, themes, states, viewports checked |

Fix any confirmed in-scope `FAIL` with minimal diffs and re-run affected checks before finishing.
