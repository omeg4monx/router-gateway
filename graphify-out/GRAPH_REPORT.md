# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 179 files · ~190,482 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 2267 nodes · 7025 edges · 96 communities (79 shown, 17 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 673 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bff6bfd3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- package.json
- Structured Logging & Ring Buffer
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- KeysConsole.tsx
- repository
- catalog.go
- Design System & Wave 0 Plan Docs
- Bus
- updates
- ir.go
- InitForm.tsx
- with
- jobs
- resolver.go
- steps
- web-check
- on
- cross-compile
- Module
- ProvidersConsole.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- JSON Marshaler Error Test Fixture
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- net/http.ResponseWriter
- UsageConsole.tsx
- encoding/json.RawMessage
- buildApp
- UpstreamCredentials
- translate_gemini.go
- fakeCredentialStore
- repository
- context.Context
- translate_openai.go
- NewStaticUIHandler
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- WriteError
- Module
- Module
- Module
- CanonicalStreamChunk
- repository
- Config
- ProviderDetailView.tsx
- NewStreamTranslator
- batcher
- Module
- Module
- stubModule
- recordingInterceptor
- newTestHarness
- scripts
- pricingStore
- responsesStreamEncoder
- newTestHarness
- .RequireLogin
- translate_test.go
- time.Time
- DB
- stubModule
- sync.Mutex
- time.Duration
- Load
- stubModule
- CliToolsConsole.tsx
- stubModule
- collectMigrations
- dependencies
- batcher
- newTestHarness
- newTestHarness
- settingDefinition
- nameMinHeap
- app/layout.tsx
- crypto.SessionSigner
- net/http.Request
- orderStubModule
- @radix-ui/react-label

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 67 edges
2. `extractErrorMessage()` - 62 edges
3. `WriteJSON()` - 58 edges
4. `DB` - 46 edges
5. `buildApp()` - 45 edges
6. `ModelsRoutingConsole()` - 36 edges
7. `Module` - 34 edges
8. `Module` - 34 edges
9. `Module` - 34 edges
10. `ProvidersConsole()` - 34 edges

## Surprising Connections (you probably didn't know these)
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `File Map` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/contracts/contracts.go
- `Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/contracts/contracts.go
- `Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests` --references--> `Options`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/modules/gateway/module.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (96 total, 17 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.05
Nodes (80): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+72 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.08
Nodes (34): @tanstack/react-query, @tanstack/react-query@^5.76.1, authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+26 more)

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (47): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+39 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "APIKeySummary"
Cohesion: 0.06
Nodes (32): stubKeyValidator, File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan (+24 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), RingBuffer, parseLevel(), boundAttrs, LogEntry (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.13
Nodes (15): loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement() (+7 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.09
Nodes (54): class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, InitPage() (+46 more)

### Community 11 - "testing.T"
Cohesion: 0.07
Nodes (52): TestAppCloseTeardownOrderIsLIFO(), TestSentinelErrorsDistinct(), setupGatewayTest(), TestGatewayAuthExtractionAndProtocolErrors(), TestGatewayCrossProtocolStreaming(), TestGatewayDirectFastPathIncrementalStreaming(), TestGatewayErrorClassificationAndOAuthSingleRetry(), TestGatewayInterceptorsPriorityOrdering() (+44 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.11
Nodes (23): authHarness, fakeCredentialStore, Options, Event, New(), decodeErrorCode(), doRequest(), findCookie() (+15 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.13
Nodes (33): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+25 more)

### Community 21 - "InitForm.tsx"
Cohesion: 0.13
Nodes (31): @hookform/resolvers, react-hook-form, zod, zod@^4.6.7, DashboardSettingsPage(), buildSettingsPatchRequest(), ChangePasswordFormValues, changePasswordSchema (+23 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.11
Nodes (18): fakeProviderStore, resolveSnapshot, targetKey, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, ModelDescriptor (+10 more)

### Community 25 - "steps"
Cohesion: 0.50
Nodes (8): steps, steps, steps, steps, name, run, uses, steps

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "cross-compile"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, cross-compile

### Community 29 - "Module"
Cohesion: 0.22
Nodes (10): actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription(), validateModelIdentifier() (+2 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.05
Nodes (76): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation() (+68 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 38 - "CanonicalUsage"
Cohesion: 0.12
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+9 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (16): attemptOutcome, statusError, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), applyUpstreamHeaders() (+8 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.20
Nodes (17): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+9 more)

### Community 42 - "buildApp"
Cohesion: 0.07
Nodes (23): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer() (+15 more)

### Community 43 - "UpstreamCredentials"
Cohesion: 0.13
Nodes (7): stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall, storedCredentials

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.07
Nodes (6): Module, validatePassword(), stubSystemModule, stubModule, stubSystemModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.14
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "context.Context"
Cohesion: 0.07
Nodes (3): stubAdminCredentialStore, actorFromContext(), oauthTokenResponse

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 50 - "newTestHarness"
Cohesion: 0.09
Nodes (24): Options, Task 2: `auth` Module & OpenAPI Fragment, File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New() (+16 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.15
Nodes (22): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+14 more)

### Community 53 - "WriteError"
Cohesion: 0.16
Nodes (16): buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), Module, truncateError() (+8 more)

### Community 54 - "Module"
Cohesion: 0.06
Nodes (21): File Map, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification) (+13 more)

### Community 55 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 58 - "repository"
Cohesion: 0.16
Nodes (20): cloneStringMap(), cloneStringSlice(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+12 more)

### Community 59 - "Config"
Cohesion: 0.32
Nodes (7): testHarness, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback()

### Community 60 - "ProviderDetailView.tsx"
Cohesion: 0.12
Nodes (33): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), DashboardPage(), DashboardProviderDetailPage(), useDeleteProviderConnectionMutation(), useProviderConnection(), useTestProviderConnectionMutation(), useUpdateProviderConnectionMutation(), formatTimestamp() (+25 more)

### Community 61 - "NewStreamTranslator"
Cohesion: 0.22
Nodes (9): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder() (+1 more)

### Community 62 - "batcher"
Cohesion: 0.27
Nodes (4): stubUsageRecorder, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, stubUsageRecorder, UsageRecord

### Community 63 - "Module"
Cohesion: 0.28
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 66 - "recordingInterceptor"
Cohesion: 0.12
Nodes (11): GatewayInterceptor, ProxyDialerResolver, stubInterceptor, stubProxyDialerResolver, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent (+3 more)

### Community 67 - "newTestHarness"
Cohesion: 0.11
Nodes (15): fakeClock, stubAdminAuth, testHarness, decodeJSON(), newFakeClock(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency() (+7 more)

### Community 68 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 70 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 71 - "newTestHarness"
Cohesion: 0.18
Nodes (11): decodeJSON(), newFakeClock(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageSummaryLogsAndRetentionPrune() (+3 more)

### Community 72 - ".RequireLogin"
Cohesion: 0.09
Nodes (6): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "translate_test.go"
Cohesion: 0.24
Nodes (15): CanonicalChatResponse, DecodeRequest(), DecodeResponse(), EncodeResponse(), rawJSONToString(), encodeAnthropicResponse(), encodeGeminiResponse(), encodeOpenAIChatResponse() (+7 more)

### Community 74 - "time.Time"
Cohesion: 0.19
Nodes (9): CooldownEntry, cooldownState, cooldownTracker, listCooldownsResponse, Options, newCooldownTracker(), New(), TestModuleMetadataAndStartValidation() (+1 more)

### Community 75 - "DB"
Cohesion: 0.09
Nodes (21): main(), Global Constraints, Global Constraints, Global Constraints, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, chainStore, newChainStore(), New() (+13 more)

### Community 77 - "sync.Mutex"
Cohesion: 0.25
Nodes (3): failureRecord, stubHealthTracker, testEnv

### Community 79 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 83 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 84 - "dependencies"
Cohesion: 0.11
Nodes (19): clsx, openapi-fetch, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies, class-variance-authority@^0.7.1 (+11 more)

### Community 85 - "batcher"
Cohesion: 0.20
Nodes (5): newBatcher(), New(), newPricingStore(), batcher, Options

### Community 86 - "newTestHarness"
Cohesion: 0.18
Nodes (11): assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), New(), TestConcurrentPublishSubscribeUnsubscribe() (+3 more)

### Community 87 - "newTestHarness"
Cohesion: 0.27
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 88 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 90 - "app/layout.tsx"
Cohesion: 0.24
Nodes (7): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 92 - "net/http.Request"
Cohesion: 0.14
Nodes (13): trackedResponseWriter, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), extractAPIKey(), actorFromContext() (+5 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 391 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **17 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `APIKeySummary` to `Module`, `Server`, `newTestHarness`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.313) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.053588516746411484 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `ProviderDetailView.tsx` to `AppShell.tsx`, `APIKeySummary`?**
  _High betweenness centrality (0.220) - this node is a cross-community bridge._
- **Should `extractErrorMessage` be split into smaller, more focused modules?**
  _Cohesion score 0.08076923076923077 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `ProviderDetailView.tsx`, `APIKeySummary`?**
  _High betweenness centrality (0.185) - this node is a cross-community bridge._