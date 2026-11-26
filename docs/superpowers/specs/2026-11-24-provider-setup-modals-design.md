# Provider Setup Modals Redesign & Pre-Save Connectivity Test

**Date:** 2026-11-24
**Status:** Approved (Approach A)
**Scope:** `/dashboard/providers` setup modals (`ModalShell`, `AddApiKeyModal`, `OAuthWizardModal`, `AddCustomNodeModal`) + Go stateless probe endpoint (`POST /api/providers/nodes/test`)

## Intent & Why

The provider setup modals on `/dashboard/providers` currently share a narrow `max-w-lg` shell without a sticky footer, use an unsearchable 80+ item native `<select>`, stack both OAuth wizard steps vertically at once, and provide no way to verify a custom compatible endpoint before saving it. Redesigning the modals with proportional layout, sticky action footers, progressive disclosure, searchable provider selection, and a stateless pre-save connectivity test makes onboarding new providers faster and less error-prone.

## Goals

- Make all three setup modals (`AddApiKeyModal`, `OAuthWizardModal`, `AddCustomNodeModal`) proportional (`max-w-xl`), with sticky headers and sticky footers plus body scroll lock.
- Add a stateless backend probe endpoint (`POST /api/providers/nodes/test`) so users can test a custom endpoint (and single API-key credentials) before saving, without creating orphan DB records or hitting browser CORS limits.
- Streamline `AddCustomNodeModal` with 1-click presets (`Ollama`, `vLLM`, `LiteLLM`, `LM Studio`), auto-slugified `Provider ID`, optional API key (used for both test and optional initial connection creation on save), and collapsible advanced defaults.
- Upgrade `AddApiKeyModal` with a searchable provider picker (or compact locked header card with "Change" toggle when opened from a provider card), live bulk-import line counting/validation, and pre-save connectivity testing.
- Upgrade `OAuthWizardModal` with a 2-step visual stepper, progressive disclosure for Step 2, and automatic `code`/`state` extraction from pasted callback URLs.

## Non-Goals

- Redesigning `ProviderConnectionsDrawer` or `ProviderDetailView`.
- Blocking save if a user skips or fails the pre-save connectivity test (offline/air-gapped setup must remain possible).
- Changing existing database schema or encrypted credential storage format.

## Approaches Considered

1. **Approach A — Backend Stateless Probe (`POST /api/providers/nodes/test`) + Full 3-Modal Redesign (Selected):** Reuses Go's existing `buildProtocolTestRequest`, `validateHTTPURL`, `validateHeaders`, and `sanitizeUpstreamError`. Avoids CORS issues, works with localhost/LAN upstreams, leaves zero orphan DB rows, and powers pre-save testing for both Custom Endpoint and API Key modals.
2. **Approach B — Frontend Create-then-Test:** Creates temporary node/connection rows, tests via `/connections/{id}/test`, and deletes on cancel. Rejected because it leaves orphan records on tab close/failure and triggers unnecessary event-bus/cache churn.

## Constraints

- Go backend must validate `protocol`, `base_url` (http/https), and `headers` using existing validators and never echo raw `api_key` in error responses.
- OpenAPI spec (`contracts/openapi/modules/providers.yaml` and `contracts/openapi/openapi.yaml`) must stay in sync with generated TypeScript types (`web/src/shared/api/schema.d.ts`).
- Must preserve full bilingual support (`isId` for Indonesian and English) and existing CSS variable design tokens (`--color-*`, `--status-*`).
- Zero new third-party dependencies.
