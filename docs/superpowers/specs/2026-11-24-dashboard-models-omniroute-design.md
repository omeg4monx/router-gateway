# `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split

**Date:** 2026-11-24
**Status:** Approved (Option 1)
**Scope:** Backend Go endpoints (`internal/modules/combos/`, `contracts/openapi/modules/combos.yaml`), frontend route split (`/dashboard/combos` + `/dashboard/models`), and full OmniRoute Model Catalog UI (`web/src/features/models/`).

## Intent & Why

Currently `/dashboard/models` renders `ModelsRoutingConsole` (Gateway Quick-Reference, Combo Builder, Model Aliases, Disabled Models, Account Cooldowns), whereas OmniRoute v3.8.54's `/dashboard/models` is a dedicated **Model Catalog** with multi-facet filtering, provider health matrix, per-model/bulk live testing, visibility management, and API-key access assignment. Moving the legacy routing/combo builder to `/dashboard/combos` preserves all existing functionality while allowing `/dashboard/models` to achieve 100% parity with OmniRoute's Model Catalog.

## Goals

- Move `ModelsRoutingConsole` to `/dashboard/combos` and add a separate `Combos` item alongside `Models` in `ROUTING_ITEMS` (`web/src/shared/ui/navigation.ts`).
- Implement OmniRoute's backend endpoints in Go (`internal/modules/combos/`):
  - `GET /api/models/catalog` (with `?includeHidden=1`) returning provider-grouped models with rich metadata (type, subtype, capabilities, context length, max output tokens, free/custom flags, hidden status synced with `disabled_models`).
  - `POST /api/models/test` & `POST /api/models/test-all` probing upstream provider connections for specific models with latency, HTTP status, rate-limit/quota/timeout classification.
  - `POST /api/combos/test` probing a combo's ordered target chain and returning per-target results + `resolvedBy`.
  - `GET /api/providers/health-matrix` aggregating connection enabled/cooldown/test state into `healthy`, `degraded`, `down`, or `untested` per provider.
- Rebuild `/dashboard/models` (`web/src/features/models/ModelCatalogConsole.tsx`) matching OmniRoute's full Model Catalog UX:
  - Header with visibility counts (`visible / hidden / total`) and Refresh button.
  - Underline tabs (`Models` and `Combos`).
  - Multi-facet filter bars for both tabs.
  - Bulk toolbar (`Assign to Keys`, `Test Selected`, `Test All Filtered` with confirm dialog + cancel + progress bar, `Clear Results`, `Hide Filtered`, `Show Filtered`).
  - Persistent test result cache in `localStorage` (`omniroute.catalogTestResults.v1`).
  - `AssignKeysModal` supporting search, add/remove mode, and restricted-mode switch confirmation.
  - Bilingual EN/ID copy and existing CSS token theming.

## Non-Goals

- Redesigning `/dashboard/combos` beyond relocating `ModelsRoutingConsole` cleanly.
- Introducing new third-party UI or state libraries.

## Approaches Considered

1. **Option 1 — Split `/dashboard/combos` & Rebuild `/dashboard/models` as OmniRoute Model Catalog (Selected):** Preserves Combo Builder, Aliases, Cooldowns, and Gateway Quick-Reference at `/dashboard/combos` while giving `/dashboard/models` 100% OmniRoute parity.
2. **Option 2 — Replace `/dashboard/models` In-Place:** Drops Combo Builder UI until a future pass. Rejected to avoid regressing combo/alias creation.

## Constraints

- OpenAPI spec (`contracts/openapi/modules/combos.yaml` + `openapi.yaml`) must stay in sync with generated `web/src/shared/api/schema.d.ts`.
- Never leak provider secrets/API keys in test probe errors.
- Static Next.js export (`NEXT_OUTPUT_EXPORT=true`) must succeed cleanly.
