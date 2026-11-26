# Implementation Plan: `/dashboard/models` OmniRoute Parity

## Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`)
1. Add `catalog.go` in `internal/modules/combos/`:
   - `handleModelCatalog` (`GET /api/models/catalog?includeHidden=1`)
   - `handleTestModel` (`POST /api/models/test`)
   - `handleTestAllModels` (`POST /api/models/test-all`)
   - `handleTestCombo` (`POST /api/combos/test`)
   - `handleProviderHealthMatrix` (`GET /api/providers/health-matrix`)
   - Curated metadata inference (type, subtype, capabilities, context_length, max_output_tokens, free flag) for known model families + custom/provider connections + hidden state from `disabled_models`.
   - Real upstream probe via `m.providers.ResolveCredentials` + `http.Client` (with injectable `HTTPClient` in `combos.Options` for deterministic unit tests and timeout/error sanitization).
2. Wire routes in `internal/modules/combos/module.go` and pass `HTTPClient: httpClient` in `cmd/omeg4mon-router/main.go`.
3. Add comprehensive TDD tests in `internal/modules/combos/catalog_test.go`.

## Task 2: OpenAPI Contract & TypeScript Codegen
1. Update `contracts/openapi/modules/combos.yaml` and `contracts/openapi/openapi.yaml` (if needed) with the 5 new endpoints and schemas.
2. Run `make codegen` (`cd web && bun run gen:api`) to update `web/src/shared/api/schema.d.ts`.

## Task 3: Route Split & Navigation Update
1. Create `web/src/app/(dashboard)/dashboard/combos/page.tsx` rendering `ModelsRoutingConsole`.
2. Update `web/src/shared/ui/navigation.ts` to have separate `Combos` (`/dashboard/combos/`) and `Models` (`/dashboard/models/`) items in `ROUTING_ITEMS`.
3. Update `web/src/shared/ui/navigation.test.ts` and `web/src/features/system/overview.ts` if needed.

## Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`)
1. Write `types.ts` and `catalogUtils.ts` + TDD unit tests in `catalogUtils.test.ts`.
2. Write `hooks.ts` for catalog, combos, health matrix, keys, model/combo testing, visibility toggling, and key model access updates.
3. Write `AssignKeysModal.tsx` (search keys, add/remove mode, restricted-mode confirmation, per-key results).
4. Write `ModelCatalogConsole.tsx` (header + visibility counts + refresh, Models/Combos tabs, multi-facet filter bars, bulk toolbar + confirm modal + progress bar + cancel, sortable/paginated tables, localStorage test result cache, EN/ID localization).
5. Point `web/src/app/(dashboard)/dashboard/models/page.tsx` to `ModelCatalogConsole`.
6. Write parity tests in `web/src/features/models/omnirouteModelsParity.test.ts`.

## Task 5: Full Verification & Graphify Update
1. Run `go test -count=1 -race ./internal/modules/combos/... ./cmd/omeg4mon-router/...`
2. Run `cd web && bun test && bun run typecheck && bun run check`
3. Run `cd web && NEXT_OUTPUT_EXPORT=true bun run build`
4. Run `graphify update .`
