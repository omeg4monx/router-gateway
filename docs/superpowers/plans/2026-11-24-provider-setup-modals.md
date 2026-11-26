# Provider Setup Modals Redesign & Pre-Save Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the `/dashboard/providers` setup modals (`ModalShell`, `AddApiKeyModal`, `OAuthWizardModal`, `AddCustomNodeModal`) for proportional UI/UX and add a stateless pre-save connectivity test endpoint (`POST /api/providers/nodes/test`) for custom endpoints and API key drafts.

**Architecture:** Add a stateless admin probe endpoint in the Go `providers` module that reuses `buildProtocolTestRequest` and `sanitizeUpstreamError` without touching SQLite. Expose it through OpenAPI + generated TypeScript types + React Query mutation hook, extract pure modal helpers (`modalUtils.ts`) under TDD, and refactor `ProvidersModals.tsx` with a 3-zone `ModalShell` (sticky header, scrollable body, sticky footer).

**Tech Stack:** Go (`chi`, `net/http`), OpenAPI 3.1 + `openapi-typescript`, Next.js 16 / React 19, Tailwind CSS v4, React Hook Form + Zod v4, TanStack Query, Bun test, Biome.

**Spec:** `docs/superpowers/specs/2026-11-24-provider-setup-modals-design.md`

## Global Constraints

- Reuses existing `buildProtocolTestRequest`, `validateHTTPURL`, `validateHeaders`, `isValidProtocol`, and `sanitizeUpstreamError` in Go; never writes to DB or leaks `api_key` in error output.
- OpenAPI contract (`contracts/openapi/modules/providers.yaml`, `contracts/openapi/openapi.yaml`) and `web/src/shared/api/schema.d.ts` must stay in sync via `bun run gen:api`.
- Bilingual copy (`isId`: Indonesian / English) and existing CSS design tokens (`--color-*`, `--status-*`) must be preserved across all modals.
- Zero new npm or Go dependencies.
- Saving a custom endpoint or connection must never be blocked if the user chooses not to run or ignores a failed pre-save test.

## Review Focus

- Secret redaction in draft probe errors: if an upstream or transport error echoes the draft `api_key`, `sanitizeUpstreamError` must scrub it before returning JSON.
- Local/no-auth custom endpoints (e.g., Ollama on `http://127.0.0.1:11434/v1`): `POST /api/providers/nodes/test` must succeed when `api_key` is omitted or empty.
- Stale test indicator: editing `protocol`, `base_url`, `api_key`, or `headers` after a test must mark the displayed result stale so users are never misled by an outdated green badge.
- Partial failure when saving Custom Endpoint with an optional API Key: if node registration succeeds (`POST /api/providers/nodes`) but initial connection creation fails (`POST /api/providers/connections`), surface a clear error without losing track of the created node.
- OAuth callback URL parsing: handle query parameters (`?code=...&state=...`), hash fragments (`#code=...&state=...`), and non-URL / missing-param inputs gracefully without throwing.

---

### Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`)

**Files:**

- Modify: `internal/modules/providers/module.go:140-158`
- Modify: `internal/modules/providers/handlers.go:640-705`
- Test: `internal/modules/providers/module_test.go`

**Interfaces:**

- Produces: HTTP route `POST /api/providers/nodes/test` (AdminAPI)
  - Request JSON: `{ "provider_id"?: string, "protocol": ProviderProtocol, "base_url": string, "api_key"?: string, "headers"?: Record<string, string> }`
  - Response JSON (200 OK): `testConnectionResponse` `{ "ok": bool, "status_code": int, "latency_ms": int64, "error"?: string }`
  - Validation errors: HTTP 400 `invalid_request`

- [ ] **Step 1: Write failing Go tests `TestProviderDraftProbe` in `internal/modules/providers/module_test.go`**
  - Subtests:
    1. `openai-chat` with `api_key` and custom `headers` hits `/models`, verifies `Authorization: Bearer ...` and custom header, returns `ok: true, status_code: 200`.
    2. No `api_key` (local Ollama/vLLM) hits `/models` without `Authorization` header, returns `ok: true, status_code: 200`, and creates 0 nodes/connections in DB.
    3. Upstream 401 echoing the submitted `api_key` returns `ok: false, status_code: 401`, and `error` does NOT contain the raw `api_key`.
    4. Invalid `protocol`, invalid `base_url`, or invalid `headers` returns HTTP 400 `invalid_request`.

- [ ] **Step 2: Run Go tests to verify failure**
  - Run: `go test -v ./internal/modules/providers -run TestProviderDraftProbe`
  - Expected: FAIL (404 / route not registered).

- [ ] **Step 3: Implement `testDraftNodeRequest`, `handleTestNodeDraft`, and mount `POST /providers/nodes/test`**
  - Register `mounts.AdminAPI.Post("/providers/nodes/test", m.handleTestNodeDraft)` in `module.go`.
  - Implement `handleTestNodeDraft` in `handlers.go` validating `protocol`, `base_url`, `headers`, building `contracts.UpstreamCredentials`, executing via `m.httpClient.Do(testReq)` (reusing a shared helper with `handleTestConnection`), and returning `testConnectionResponse`.

- [ ] **Step 4: Run Go tests to verify pass**
  - Run: `go test -v ./internal/modules/providers/...`
  - Expected: PASS (all existing + new tests green).

---

### Task 2: OpenAPI Contract, Generated Schema, and React Query Hook

**Files:**

- Modify: `contracts/openapi/modules/providers.yaml`
- Modify: `contracts/openapi/openapi.yaml`
- Regenerate: `web/src/shared/api/schema.d.ts`
- Modify: `web/src/features/providers/schemas.ts`
- Modify: `web/src/features/providers/hooks.ts`

**Interfaces:**

- Consumes: `POST /api/providers/nodes/test` from Task 1
- Produces:
  - `TestProviderNodeRequest` type export in `web/src/features/providers/schemas.ts`
  - `useTestProviderNodeDraftMutation()` hook export in `web/src/features/providers/hooks.ts` returning `UseMutationResult<TestProviderConnectionResponse, Error, TestProviderNodeRequest>`

- [ ] **Step 1: Add `/api/providers/nodes/test` path and `TestProviderNodeRequest` schema to `contracts/openapi/modules/providers.yaml` and `$ref` entries in `contracts/openapi/openapi.yaml`**

- [ ] **Step 2: Regenerate TypeScript OpenAPI types**
  - Run: `cd web && bun run gen:api`
  - Expected: `web/src/shared/api/schema.d.ts` updated with `/api/providers/nodes/test` and `TestProviderNodeRequest`.

- [ ] **Step 3: Export `TestProviderNodeRequest` in `schemas.ts` and implement `useTestProviderNodeDraftMutation` in `hooks.ts`**

- [ ] **Step 4: Verify TypeScript compilation**
  - Run: `cd web && bun run typecheck`
  - Expected: exit code 0.

---

### Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD)

**Files:**

- Create: `web/src/features/providers/modalUtils.ts`
- Create: `web/src/features/providers/modalUtils.test.ts`

**Interfaces:**

- Produces:
  - `CUSTOM_ENDPOINT_PRESETS: ReadonlyArray<CustomEndpointPreset>` (`ollama`, `vllm`, `litellm`, `lmstudio`)
  - `slugifyProviderId(input: string): string`
  - `parseOAuthCallbackUrl(rawUrl: string): ParsedOAuthCallback`
  - `parseBulkCredentials(rawText: string, providerName: string): ParsedBulkCredentials`
  - `buildTestFingerprint(input: DraftTestFingerprintInput): string`

- [ ] **Step 1: Write failing unit tests in `web/src/features/providers/modalUtils.test.ts`**
  - Tests for `slugifyProviderId` (spaces, uppercase, special chars, repeated dashes, length cap `MAX_ID_RUNES`).
  - Tests for `parseOAuthCallbackUrl` (query params, hash fragment params, partial/missing params, invalid URL).
  - Tests for `parseBulkCredentials` (`name|key`, raw `key`, blank lines, invalid `name|` empty-key lines).
  - Tests for `buildTestFingerprint` (identical inputs match; changing any of protocol/baseUrl/apiKey/headers changes fingerprint).

- [ ] **Step 2: Run tests to verify failure**
  - Run: `cd web && bun test src/features/providers/modalUtils.test.ts`
  - Expected: FAIL (module not found).

- [ ] **Step 3: Implement `web/src/features/providers/modalUtils.ts`**

- [ ] **Step 4: Run tests to verify pass**
  - Run: `cd web && bun test src/features/providers/modalUtils.test.ts`
  - Expected: PASS.

---

### Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`

**Files:**

- Modify: `web/src/features/providers/ProvidersModals.tsx`
- Modify: `web/src/features/providers/ProvidersConsole.tsx` (update `onCreated` message if custom node also created an initial connection)

**Interfaces:**

- Consumes:
  - `useTestProviderNodeDraftMutation` from `hooks.ts`
  - `CUSTOM_ENDPOINT_PRESETS`, `slugifyProviderId`, `parseOAuthCallbackUrl`, `parseBulkCredentials`, `buildTestFingerprint` from `modalUtils.ts`
- Produces:
  - Preserves exported component signatures `AddApiKeyModal`, `OAuthWizardModal`, `AddCustomNodeModal`, `ProviderConnectionsDrawer` (with backward-compatible optional enhancements).

- [ ] **Step 1: Upgrade `ModalShell`**
  - Add `size?: "md" | "lg" | "xl"` (`max-w-lg`, `max-w-xl`, `max-w-2xl`), `footer?: ReactNode` sticky footer slot, `badge?: ReactNode`, and body scroll lock on open.
  - Use proportional spacing (`px-6 py-4` header/footer, `px-6 py-5` body).

- [ ] **Step 2: Redesign `AddCustomNodeModal` with Pre-Save Test**
  - Add 1-click preset bar (`Ollama`, `vLLM`, `LiteLLM`, `LM Studio`), auto-slug `Provider ID` from `Display Name` until manually edited, optional `API Key` field with Show/Hide toggle, collapsible _Default Models & Custom Headers_ section.
  - Add **"Tes Koneksi"** button in sticky footer calling `useTestProviderNodeDraftMutation`, inline result card (`OK 200 · Nms` or error), and stale indicator via `buildTestFingerprint`.
  - When `API Key` is provided on submit, register custom node then create its initial connection (with clear partial-error messaging if step 2 fails).

- [ ] **Step 3: Redesign `AddApiKeyModal`**
  - Add compact active-provider card + searchable provider picker (with `ProviderLogo`, name, id, protocol badge) replacing the 80+ item native `<select>`.
  - Keep provider visible in both Single and Bulk modes; add live valid/invalid count indicator in Bulk mode via `parseBulkCredentials`.
  - Add pre-save **"Tes Koneksi"** button in Single mode sticky footer + inline result card with stale detection.
  - Enhance collapsible _Advanced Settings_ header with active defaults summary (`P100 · N models`).

- [ ] **Step 4: Redesign `OAuthWizardModal`**
  - Add 2-step visual progress stepper (`1. Buat Sesi & Login` → `2. Verifikasi Callback`).
  - Progressive disclosure for Step 2 (prominent after Step 1 succeeds, with _"Sudah punya callback URL?"_ shortcut).
  - Auto-parse `callback_url` via `parseOAuthCallbackUrl` to populate `code` & `state` with green detected badge, moving manual `code`/`state` fields into a collapsible toggle.

---

### Task 5: Full Verification & Graphify Sync

- [ ] **Step 1: Run Go tests, Bun tests, TypeScript check, Biome check, and Next.js production build**
  - Run: `go test ./internal/modules/providers/... && cd web && bun test && bun run typecheck && bun run check && bun run build`
  - Expected: All commands exit 0 with zero errors.

- [ ] **Step 2: Update knowledge graph**
  - Run: `graphify update .`
  - Expected: Graph updated cleanly.
