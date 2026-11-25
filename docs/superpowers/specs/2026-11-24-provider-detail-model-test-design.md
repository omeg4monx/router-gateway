# Provider Detail: Proportional Model Grid & Per-Model Test

**Date:** 2026-11-24
**Status:** Approved

## Intent & Why
On the Provider Detail page (`ProviderStudioDetail`), detected models are currently rendered as variable-width flex-wrap chips where differing ID lengths and optional subtitles cause uneven widths and heights. Operators also have no way to verify whether a specific detected model works on a provider connection without manually opening another screen.

## Goals
1. Render detected models in a proportional, uniform-height responsive grid (`1` column mobile, `2` tablet, `3` desktop) with clean truncation, full ID tooltip, badges (`LIVE`, `FREE`), and aligned actions.
2. Provide a real per-model test action on every model card that runs a lightweight inference ping against the provider's active connection across all supported protocols (`openai-chat`, `openai-responses`, `anthropic`, `gemini`) with OAuth token auto-refresh on `401`/`403`.
3. Display non-blocking inline per-model test states (loading spinner, `OK · <ms>`, `FAIL` + compact error message) without affecting other models or marking the connection as cooled down.

## Non-Goals
- Batch "Test All Models" runner (not requested; can be added later).
- Persisting per-model test history to SQLite.

## Constraints
- Preserve Omegamon Zwart (Dark) and Omegamon X-Antibody (Light) CSS variable tokens and `en`/`id` bilingual UI copy.
- Keep OpenAPI contract (`contracts/openapi/modules/providers.yaml`) synced with `web/src/shared/api/schema.d.ts` via `make codegen`.
- Never leak upstream credentials/tokens in error responses; never trigger connection health cooldowns from an exploratory model probe.

## Approaches Considered
1. **Selected — Direct Gateway Admin Endpoint (`POST /api/providers/connections/{id}/test-model`)**: Mounted on `mounts.AdminAPI` inside `internal/modules/gateway`. Reuses `m.providers.ResolveCredentials`, `prepareUpstreamBody`, `buildUpstreamURL`, `applyUpstreamHeaders`, and `RefreshOAuthToken`. Works for every detected model including live-discovered models not yet saved to `conn.models`, with zero circular module dependencies.
2. **Rejected — Routing via `CompleteChat` / `combos.Resolve`**: Requires models to already be saved/routable in `combos` and triggers connection cooldowns on failure.
3. **Rejected — Duplicating protocol request builders in `providers`**: Duplicates translation logic already maintained in `gateway`.
