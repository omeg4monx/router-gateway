# omeg4mon-router Project Standards

Concise repository guardrails complementary to global lean-engineering and UI policies. Do not duplicate global rules; enforce what is specific to this codebase.

## Architecture & Boundaries

- **Go backend:** Keep shared infrastructure in `internal/platform/{config,crypto,db,eventbus,httpclient,httpserver,jobs,lifecycle,logging}` and domain features isolated in `internal/modules/<domain>/`.
- **API contracts:** `contracts/openapi/openapi.yaml` (and `contracts/openapi/modules/*.yaml`) is the single source of truth. Any HTTP schema or route change MUST update the OpenAPI spec and regenerate `web/src/shared/api/schema.d.ts` via `make codegen`.
- **Web frontend (`web/src/`):** Follow strict layer boundaries: `pages/` (route entry) → `features/<domain>/` (`hooks.ts`, `schemas.ts`, `*Console.tsx`) → `shared/{api,lib,store,ui}` (reusable primitives). Do not import across unrelated features directly.
- **Embedded UI budget:** `web/out` is embedded into the Go binary (`-tags embedui`). Avoid heavy npm dependencies, unsubsetted fonts, or large static assets.

## Reliability, Memory & Type Safety

- **Go:** Every goroutine must have an explicit exit path (`ctx.Done()` or `internal/platform/lifecycle`). Always close `sql.Rows` and check `rows.Err()`, stop `time.Ticker`/`time.Timer`, drain/close HTTP/SSE streams on client disconnect, and bound in-memory maps/slices/caches.
- **TypeScript/React:** Zero `any` or unsafe `as unknown as` casts. Validate API/form boundaries with generated OpenAPI types + Zod. Clean up all `addEventListener`, timers, `AbortController`, and `EventSource`/subscriptions in `useEffect`.
- **Errors:** Return actionable, user-safe messages (`ErrorEnvelope` on backend; `extractErrorMessage` + recovery action on frontend). Never swallow errors or leak secrets/tokens.

## Design System & UI/UX Identity

- **Tokens only:** Use semantic CSS variables from `web/src/app/globals.css` (`--color-*`, `--bg-*`, `--text-*`, `--border-*`, `--status-*`). Never hardcode raw hex colors or introduce generic AI-template gradients/bloat.
- **Dual theme & locale:** Preserve contrast and visual fidelity in both Dark (_Omegamon Zwart_) and Light (_Omegamon X-Antibody_) themes (verified by `web/src/app/palette.test.ts`), and maintain `en`/`id` copy parity via `useUIStore`.
- **UX discipline:** Reuse `web/src/shared/ui/` primitives (`Button`, `Card`, `Alert`, `Badge`, `Input`, `Label`, `ErrorBoundary`). Every async view must handle loading (non-blocking skeleton/inline/pending state), empty, and error states cleanly, remain proportional across mobile/tablet/desktop without horizontal overflow, and avoid cluttered or convoluted layouts.

## Mandatory Verification

- Run `make check` before claiming any code change complete.
- For non-trivial changes or any UI/UX work, invoke the project skill `omeg4mon-verify` (`.agents/skills/omeg4mon-verify/SKILL.md`).
