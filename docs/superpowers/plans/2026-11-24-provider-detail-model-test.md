# Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Render detected models on Provider Detail in a proportional responsive grid and allow testing each model individually via a real lightweight upstream probe.

**Architecture:** Add `POST /api/providers/connections/{id}/test-model` to the OpenAPI spec and implement it in `internal/modules/gateway/test_model.go` (reusing gateway's protocol translation, upstream headers, and OAuth refresh without tripping connection cooldowns). Add a typed TanStack Query mutation hook in `web/src/features/providers/hooks.ts` and update `ProviderStudioDetail` in `web/src/features/providers/ProviderDetailView.tsx` to use a uniform responsive CSS grid with per-model Test button and inline status feedback.

**Tech Stack:** Go (`chi`, `net/http`), OpenAPI 3.1 + `openapi-typescript`, React 19, TanStack Query, Tailwind CSS v4, Bun test.

**Spec:** `docs/superpowers/specs/2026-11-24-provider-detail-model-test-design.md`

## Global Constraints
- All colors/surfaces/borders must use semantic CSS variables from `web/src/app/globals.css`.
- Preserve `en` / `id` bilingual copy in `ProviderDetailView.tsx`.
- API contract in `contracts/openapi/modules/providers.yaml` must stay in sync with `web/src/shared/api/schema.d.ts`.
- Never leak API keys or OAuth tokens in test error responses; never record connection health failures during exploratory model tests.

## Review Focus
1. Live-discovered models not yet saved in `conn.models` must still be testable directly against `connectionID`.
2. Models with `/` in their ID (e.g. OpenRouter `anthropic/claude-3.5-sonnet`) must be sent verbatim without stripping.
3. OAuth connections returning `401`/`403` during a model test must refresh once and retry.
4. Empty or overly long `model` input must return HTTP `400` with `ErrorEnvelope`; unknown connection ID must return HTTP `404`.
5. Providers without an active connection must disable per-model Test buttons with a clear tooltip, and multiple models can be tested independently with isolated loading/result states.

---

### Task 1: OpenAPI Contract & TypeScript Schema Codegen
**Files:**
- Modify: `contracts/openapi/modules/providers.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`

**Interfaces:**
- Produces: `POST /api/providers/connections/{id}/test-model` (`operationId: testProviderConnectionModel`), `TestProviderModelRequest` (`{ model: string }`), `TestProviderModelResponse` (`{ ok: boolean; model: string; status_code?: number; latency_ms: number; error?: string }`).

- [ ] Add `/api/providers/connections/{id}/test-model` path and `TestProviderModelRequest` / `TestProviderModelResponse` schemas to `contracts/openapi/modules/providers.yaml`.
- [ ] Run `make codegen` to update `web/src/shared/api/schema.d.ts`.

### Task 2: Backend Model Test Endpoint (TDD)
**Files:**
- Create: `internal/modules/gateway/test_model.go`
- Create: `internal/modules/gateway/test_model_test.go`
- Modify: `internal/modules/gateway/module.go` (mount route on `mounts.AdminAPI`)

**Interfaces:**
- Consumes: `m.providers.ResolveCredentials`, `m.providers.RefreshOAuthToken`, `m.prepareUpstreamBody`, `buildUpstreamURL`, `applyUpstreamHeaders`.
- Produces: `POST /api/providers/connections/{id}/test-model` HTTP handler.

- [ ] Write failing tests in `internal/modules/gateway/test_model_test.go` covering:
  - OpenAI chat model probe success (`ok: true`, `status_code: 200`, `model` preserved including slash IDs)
  - Anthropic protocol translation & upstream error handling (sanitizes secret from error body)
  - OAuth `401` auto-refresh and retry succeeding
  - Validation (`400` empty/oversized model) and `404` unknown connection
- [ ] Run `go test ./internal/modules/gateway/... -run TestHandleTestConnectionModel` to verify RED.
- [ ] Implement `handleTestConnectionModel` in `internal/modules/gateway/test_model.go` and register route in `internal/modules/gateway/module.go`.
- [ ] Run `go test -race ./internal/modules/gateway/...` to verify GREEN.

### Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD)
**Files:**
- Modify: `web/src/features/providers/schemas.ts`
- Modify: `web/src/features/providers/hooks.ts`
- Modify: `web/src/features/providers/ProviderDetailView.tsx`
- Modify: `web/src/features/providers/loadModels.test.ts`

**Interfaces:**
- Consumes: `POST /api/providers/connections/{id}/test-model` via `apiClient`.
- Produces: `TestProviderModelRequest`, `TestProviderModelResponse`, `useTestProviderModelMutation`.

- [ ] Add failing contract tests in `web/src/features/providers/loadModels.test.ts` verifying:
  - `hooks.ts` exports `useTestProviderModelMutation` calling `/api/providers/connections/{id}/test-model`
  - `ProviderDetailView.tsx` renders Available Models in a responsive grid (`grid-cols-1 sm:grid-cols-2 lg:grid-cols-3`) and includes per-model test state & action
- [ ] Run `cd web && bun test src/features/providers/loadModels.test.ts` to verify RED.
- [ ] Export types in `schemas.ts`, add `useTestProviderModelMutation` in `hooks.ts`, and update `ProviderStudioDetail` in `ProviderDetailView.tsx` with proportional grid cards and per-model test button + inline result/error feedback.
- [ ] Run `cd web && bun test src/features/providers/loadModels.test.ts` to verify GREEN.

### Task 4: Full Verification Gate & Graphify Update
- [ ] Run `make check` and `cd web && bun run build`.
- [ ] Run `graphify update .`.
