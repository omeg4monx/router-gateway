# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 482 files · ~533,277 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5561 nodes · 18572 edges · 252 communities (219 shown, 33 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1822 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a6b59b47`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main_test.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- repository
- KeysConsole.tsx
- TypeScript Strict Compiler Config
- Button
- ModelEditorConsole.tsx
- testing.T
- shadcn/ui Component Configuration
- UpstreamCredentials
- MediaConsole.tsx
- App.tsx
- Module
- Design System & Wave 0 Plan Docs
- extractErrorMessage
- updates
- Button
- buildApp
- with
- jobs
- net/http/httptest.ResponseRecorder
- steps
- web-check
- on
- .Reader
- TokenSaverConsole.tsx
- models/hooks.ts
- NotificationsConsole.tsx
- docker-build
- go-check
- combos/schemas.ts
- Module
- Go Module Root Definition
- CanonicalUsage
- providers/catalogUtils.ts
- QuotaConsole.tsx
- CanonicalStreamChunk
- main.go
- ir.go
- translate_gemini.go
- go_pkg_time
- repository
- PluginsStudio.tsx
- RouteMounts
- ModelDescriptor
- newTestHarness
- ComboRegistrationForm.tsx
- translate_responses.go
- oauthLoopbackProxy
- .Fetch
- overview.ts
- PromptGuardConsole.tsx
- Provider Quota Menu & Token Harbor Quota Integration
- repository
- migrator.go
- NewClient
- fakeHandle
- UsageConsole.tsx
- repository
- net/http.Header
- teeHandler
- ProxyPoolsConsole.tsx
- ModelDescriptor
- App.test.ts
- pricingStore
- modelIntelligence.ts
- translate_openai.go
- TelegramConsole.tsx
- Module
- promptinjection/engine.go
- go_pkg_encoding_json
- arena/hooks.ts
- combos/module.go
- awsQuickRuntime
- ComboSnapshot
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- Review Focus
- paths
- Module
- collectDynamicHostTelemetry
- .executeInferenceCandidates
- settingDefinition
- repo
- DB
- UsageRecord
- CanonicalChatResponse
- newTestHarness
- SettingsConsole.tsx
- Server
- Runner
- go_pkg_encoding_json
- ref_node_path
- net/http.Request
- Module
- ConnectionSummary
- AttestationManager
- SecurityConsole.tsx
- Module
- Module
- post
- .identifyProviderModel
- repo
- encoding/json.RawMessage
- ProviderTopology.tsx
- repo
- repository
- repo
- properties
- evaluateAntigravityQuota
- net/url.URL
- setupSecurityTest
- Module
- settingDefinition
- time.Time
- .UpsertPlugin
- properties
- Module
- ResolvedTarget
- .ProvisionKey
- io/fs.FS
- repository
- OMEG4MON Router
- KeyInfo
- stubDep
- certManager
- repo
- ProviderStore
- resolve.js
- properties
- stubCompleter
- release-binaries
- steps
- .executeCandidateAttempt
- setupMediaTest
- .RequireLogin
- ComboSnapshotter
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- stubSystemModule
- properties
- 200
- time.Duration
- properties
- dependencies
- cli/package.json
- setupTestMitm
- botAPIClient
- .UpdateCombo
- get
- scripts
- post
- schemas
- stubDiscoverySystemModule
- PreviewPluginsResponse
- stubModule
- Module
- stubSystemModule
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- setupTokenSaverTest
- Review Focus
- Review Focus
- TranslatorConsole.tsx
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- NewStreamTranslator
- WriteJSON
- parameters
- Module
- New
- stubModule
- ClearImporterHistoryResponse
- source_path
- stubDepModule
- stubSystem
- docker-release
- Module
- `@omeg4mon/router` CLI Launcher
- stubDep
- setupModelEditorTest
- context.Context
- stubSystemModule
- stubModule
- parameters
- setupCatalogTestModule
- applyUpstreamHeaders
- omeg4mon-router Verification Gate
- strategies.go
- listSettingsWithQueryer
- single-binary
- keyLimiter
- verify-web
- GatewayInterceptor
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- headroomCache
- Open
- Feature Suite
- repo
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- .Fetch
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- pricingStore
- SessionSigner
- Load
- quotaCache
- setupTunnelTest
- plugins/interceptor.go
- scanCombo
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- parseJSONObjectOrWrap
- redact.go
- time.Duration
- manifest.test.ts
- Review Focus
- jsonMarshalerError
- repository
- compressor.go
- setupNotificationsTest
- headroom.go
- palette.test.ts
- .list
- NewChromeTLSClient
- Review Focus
- HistoryItem
- omeg4mon-router Project Standards

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 188 edges
2. `WriteJSON()` - 186 edges
3. `WriteError()` - 184 edges
4. `Button` - 94 edges
5. `DB` - 93 edges
6. `cn()` - 83 edges
7. `buildApp()` - 79 edges
8. `Module` - 78 edges
9. `Bus` - 72 edges
10. `Alert` - 71 edges

## Surprising Connections (you probably didn't know these)
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/module_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/optimization_e2e_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (252 total, 33 thin omitted)

### Community 0 - "main_test.go"
Cohesion: 0.07
Nodes (27): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+19 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
Nodes (5): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC()

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (46): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+38 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "APIKeySummary"
Cohesion: 0.08
Nodes (21): cloneCachedEntry(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary() (+13 more)

### Community 6 - "repository"
Cohesion: 0.09
Nodes (27): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), hasBlockText(), hasPromptSegment(), injectLatestUserAnchor() (+19 more)

### Community 7 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (74): react-hook-form, buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isIpLiteral() (+66 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "Button"
Cohesion: 0.03
Nodes (164): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, StrategyCategory, runModelIdentityProbe() (+156 more)

### Community 10 - "ModelEditorConsole.tsx"
Cohesion: 0.05
Nodes (71): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation() (+63 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (84): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+76 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.07
Nodes (18): catalogFakeProviderStore, fakeProviderStore, ModelTestResult, resolveSnapshot, testAllModelsResponse, stubProviderStore, stubProviderStore, ProviderConnection (+10 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (37): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+29 more)

### Community 16 - "Module"
Cohesion: 0.07
Nodes (43): isBlockedAntigravityProjectID(), applyCodeBuddyProbeHeaders(), codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), TestCodeBuddyTokenRefresh(), originOfURL(), postExchangeIFlow() (+35 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "extractErrorMessage"
Cohesion: 0.04
Nodes (103): combosQueryKeys, useClearModelCooldownMutation(), useCreateComboMutation(), useDeleteComboMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+95 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "Button"
Cohesion: 0.06
Nodes (92): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, class-variance-authority, lucide-react, react, react-router, AuthGuard() (+84 more)

### Community 21 - "buildApp"
Cohesion: 0.05
Nodes (37): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+29 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.08
Nodes (30): fakeCredentialStore, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+22 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - ".Reader"
Cohesion: 0.11
Nodes (18): main(), Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+10 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 30 - "models/hooks.ts"
Cohesion: 0.23
Nodes (16): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, modelCatalogQueryKeys (+8 more)

### Community 31 - "NotificationsConsole.tsx"
Cohesion: 0.16
Nodes (24): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "combos/schemas.ts"
Cohesion: 0.07
Nodes (35): UpdateComboVariables, UpsertModelAliasVariables, COMBO_STRATEGIES, CooldownEntry, CreateComboFormInput, CreateComboFormOutput, createComboSchema, DEFAULT_CREATE_COMBO_VALUES (+27 more)

### Community 35 - "Module"
Cohesion: 0.05
Nodes (36): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), isValidAuthType() (+28 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "providers/catalogUtils.ts"
Cohesion: 0.08
Nodes (47): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), buildProviderSectionBuckets(), canProbeProviderModels(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets() (+39 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (47): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+39 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.38
Nodes (4): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder, newOpenAIStreamDecoder()

### Community 42 - "main.go"
Cohesion: 0.19
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 43 - "ir.go"
Cohesion: 0.13
Nodes (37): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat (+29 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "go_pkg_time"
Cohesion: 0.09
Nodes (19): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+11 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.19
Nodes (22): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+14 more)

### Community 48 - "RouteMounts"
Cohesion: 0.21
Nodes (3): Module, validateSQLiteSource(), decodeJSON()

### Community 49 - "ModelDescriptor"
Cohesion: 0.18
Nodes (12): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, vite, @vitejs/plugin-react, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin() (+4 more)

### Community 50 - "newTestHarness"
Cohesion: 0.11
Nodes (15): statusError, newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota() (+7 more)

### Community 51 - "ComboRegistrationForm.tsx"
Cohesion: 0.17
Nodes (23): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+15 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (26): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+18 more)

### Community 53 - "oauthLoopbackProxy"
Cohesion: 0.16
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 54 - ".Fetch"
Cohesion: 0.14
Nodes (19): clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber(), parseCodeBuddyTimestamp() (+11 more)

### Community 55 - "overview.ts"
Cohesion: 0.22
Nodes (11): filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY, OverviewCopy (+3 more)

### Community 56 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 57 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.17
Nodes (10): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+2 more)

### Community 58 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 59 - "migrator.go"
Cohesion: 0.11
Nodes (34): ModelOption, inspectRequest, runImportRequest, buildInstallGuide(), mapBusEventToAudit(), payloadToMap(), Module, ChangedEvent (+26 more)

### Community 60 - "NewClient"
Cohesion: 0.12
Nodes (15): Options, startAuthHTTPProxy(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(), applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream() (+7 more)

### Community 61 - "fakeHandle"
Cohesion: 0.08
Nodes (7): streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ssePipeReader

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (63): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass() (+55 more)

### Community 63 - "repository"
Cohesion: 0.16
Nodes (8): decodeStringMap(), decodeStringSlice(), newRepository(), scanConnectionRecord(), scanNode(), ProviderNode, repository, rowScanner

### Community 64 - "net/http.Header"
Cohesion: 0.18
Nodes (6): memoryResponseWriter, extractVendorHeaders(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (35): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+27 more)

### Community 67 - "ModelDescriptor"
Cohesion: 0.20
Nodes (9): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, ModelDescriptor, isModelAllowed(), isModelDisabled(), compositeModelLister (+1 more)

### Community 68 - "App.test.ts"
Cohesion: 0.50
Nodes (3): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 72 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "promptinjection/engine.go"
Cohesion: 0.12
Nodes (20): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+12 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (26): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+18 more)

### Community 76 - "arena/hooks.ts"
Cohesion: 0.07
Nodes (52): ArenaConsole(), ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail() (+44 more)

### Community 77 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 78 - "awsQuickRuntime"
Cohesion: 0.12
Nodes (9): newAwsQuickRuntime(), parseAwsQuickPastedQuery(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens (+1 more)

### Community 79 - "ComboSnapshot"
Cohesion: 0.23
Nodes (10): legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables(), readLegacyCombos(), readLegacyKeys() (+2 more)

### Community 80 - "web/package.json"
Cohesion: 0.12
Nodes (23): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+15 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.19
Nodes (7): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken()

### Community 83 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 86 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 87 - ".executeInferenceCandidates"
Cohesion: 0.13
Nodes (12): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), ensureRequestID(), parseInboundMeta(), TestResponseCoalescerDeduplicatesAndCachesPerKey(), computeJitteredRetryDelay() (+4 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.18
Nodes (8): mapBusEventToNotification(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotificationRow(), EmitPayload, Notification, repo

### Community 90 - "DB"
Cohesion: 0.04
Nodes (75): Options, authHarness, Options, Options, stubAdminAuth, testHarness, File Map, Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests (+67 more)

### Community 91 - "UsageRecord"
Cohesion: 0.15
Nodes (7): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher

### Community 92 - "CanonicalChatResponse"
Cohesion: 0.24
Nodes (13): CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), IsSupportedProtocol(), RewriteModelField(), TestAnthropicToOpenAIChatRoundTrip(), TestDecodeRejectsNonObjectJSONPayloads() (+5 more)

### Community 93 - "newTestHarness"
Cohesion: 0.15
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.10
Nodes (41): AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation(), useRollbackSettingMutation() (+33 more)

### Community 95 - "Server"
Cohesion: 0.08
Nodes (30): loadUI(), loadUI(), doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth(), TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping(), apiTimeoutMiddleware() (+22 more)

### Community 96 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "go_pkg_encoding_json"
Cohesion: 0.21
Nodes (9): errorDisposition, handledGatewayError, classifyStatus(), isClientPayloadErrorBody(), compareSemver(), maybeLiftFreePIClientVersion(), parseFreePIRequiredVersion(), TestFreePIParseRequiredVersionAndCompare() (+1 more)

### Community 98 - "ref_node_path"
Cohesion: 0.14
Nodes (11): consoleSource, detailSource, modalsSource, providersDir, modalsSource, providersDir, consoleSource, detailSource (+3 more)

### Community 99 - "net/http.Request"
Cohesion: 0.06
Nodes (21): stubAuthVerifier, Module, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), clientRateLimitKey(), anthropicErrorType() (+13 more)

### Community 100 - "Module"
Cohesion: 0.16
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 101 - "ConnectionSummary"
Cohesion: 0.25
Nodes (10): encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), updateConnectionTx(), connectionRecord, ConnectionSummary, createConnectionParams (+2 more)

### Community 102 - "AttestationManager"
Cohesion: 0.31
Nodes (5): AttestationManager, isTransientMintStatus(), ResolveMintURL(), inFlightMint, tokenBucket

### Community 103 - "SecurityConsole.tsx"
Cohesion: 0.11
Nodes (29): @hookform/resolvers, zod, zod@^4.6.7, AuthRedirectTarget, AuthRouteState, AuthSessionPosture, LoginFormValues, loginSchema (+21 more)

### Community 104 - "Module"
Cohesion: 0.14
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 105 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.07
Nodes (41): ModelIdentityResult, buildModelProbeRequest(), extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse(), buildIdentityProbeRequest(), extractIDPrefix() (+33 more)

### Community 108 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 109 - "encoding/json.RawMessage"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.13
Nodes (20): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+12 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "repository"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "evaluateAntigravityQuota"
Cohesion: 0.17
Nodes (15): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+7 more)

### Community 116 - "net/url.URL"
Cohesion: 0.09
Nodes (20): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+12 more)

### Community 117 - "setupSecurityTest"
Cohesion: 0.12
Nodes (23): TestAppWave1AAuthSettingsKeysEndToEnd(), File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 2: `security` Module (Security Posture & Redacted Audit Trail), AdminCredentialStore (+15 more)

### Community 118 - "Module"
Cohesion: 0.06
Nodes (27): Module, listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), trackedResponseWriter, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID() (+19 more)

### Community 119 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 120 - "time.Time"
Cohesion: 0.08
Nodes (15): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, newCooldownTracker() (+7 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "Module"
Cohesion: 0.15
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 124 - "ResolvedTarget"
Cohesion: 0.18
Nodes (12): strategyState, weightedCandidate, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), applyComboStrategy(), applyProviderStrategy() (+4 more)

### Community 125 - ".ProvisionKey"
Cohesion: 0.30
Nodes (7): Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, generateKeyID()

### Community 126 - "io/fs.FS"
Cohesion: 0.09
Nodes (7): migrationFile, actorFromContext(), decodeJSONBody(), Module, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 127 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "KeyInfo"
Cohesion: 0.19
Nodes (9): stubKeyValidator, stubKeyValidator, stubResolver, testEnv, KeyInfo, cloneKeyInfo(), hashToken(), isValidTokenFormat() (+1 more)

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "ProviderStore"
Cohesion: 0.30
Nodes (14): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI) (+6 more)

### Community 134 - "resolve.js"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "stubCompleter"
Cohesion: 0.15
Nodes (15): arenaTestEnv, stubCompleter, ChatMessage, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan (+7 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (22): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), extractUpstreamErrorMessage(), parseRetryAfter(), applyCodeBuddyHeaders() (+14 more)

### Community 140 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 141 - ".RequireLogin"
Cohesion: 0.15
Nodes (4): stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 142 - "ComboSnapshotter"
Cohesion: 0.24
Nodes (12): Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification (+4 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "time.Duration"
Cohesion: 0.11
Nodes (12): chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore(), TestTrafficShaperAdaptiveHeaderPacing(), TestTrafficShaperLimitsPerConnectionConcurrencyAndQueues() (+4 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 153 - "botAPIClient"
Cohesion: 0.43
Nodes (3): newBotAPIClient(), sanitizeBotError(), botAPIClient

### Community 154 - ".UpdateCombo"
Cohesion: 0.29
Nodes (9): createComboParams, updateComboValues, cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps(), encodeConfigJSON() (+1 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, check, check:fix, dev, format, format:check, gen:api (+4 more)

### Community 157 - "post"
Cohesion: 0.29
Nodes (8): post, operationId, requestBody, summary, tags, requestBody, content, required

### Community 158 - "schemas"
Cohesion: 0.29
Nodes (6): components, schemas, additionalProperties, required, type, PluginState

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 162 - "Module"
Cohesion: 0.24
Nodes (4): decodeJSONBody(), newID(), validateBypassRule(), Module

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "setupTokenSaverTest"
Cohesion: 0.27
Nodes (10): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory(), TestTokenSaverFinalReviewFindings() (+2 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "TranslatorConsole.tsx"
Cohesion: 0.17
Nodes (19): translatorQueryKeys, useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS, translateFormSchema (+11 more)

### Community 171 - "schemas"
Cohesion: 0.22
Nodes (9): schemas, additionalProperties, required, type, additionalProperties, required, type, ImportedProvisionedKey (+1 more)

### Community 172 - "properties"
Cohesion: 0.22
Nodes (9): type, properties, type, type, id, KeyPrefix, Name, Secret (+1 more)

### Community 173 - "RunLegacyImportResponse"
Cohesion: 0.22
Nodes (9): provisioned_keys, run, type, $ref, additionalProperties, properties, required, type (+1 more)

### Community 174 - "omeg4mon-router Versioning Verificator"
Cohesion: 0.25
Nodes (7): 1. Gather Evidence First, 2. Public Contract Surfaces, 3. SemVer Decision Matrix, 4. Per-Task vs Release-Cut Rules, 5. Rationalizations & Red Flags, 6. Required Exit Report, omeg4mon-router Versioning Verificator

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "NewStreamTranslator"
Cohesion: 0.22
Nodes (9): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs() (+1 more)

### Community 178 - "WriteJSON"
Cohesion: 0.08
Nodes (11): decodeJSON(), Module, decodeStrictJSON(), Module, validatePassword(), WriteJSON(), LogEntry, LogListResponse (+3 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 181 - "New"
Cohesion: 0.06
Nodes (39): TestWave2AE2EProviderProxyRouting(), File Map, Global Constraints, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`) (+31 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.15
Nodes (17): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+9 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (14): stubAdminCredentialStore, stubSettings, extractUsageCost(), boolToInt(), stubModule, Alias, AliasListResponse, Config (+6 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "setupCatalogTestModule"
Cohesion: 0.17
Nodes (9): setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies(), TestAllOmniRouteComboStrategiesCRUDAndPersistence(), TestOmniRouteResolverStrategiesExecution() (+1 more)

### Community 198 - "applyUpstreamHeaders"
Cohesion: 0.17
Nodes (11): Task 2: Backend Model Test Endpoint (TDD), Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isBlockedAntigravityProjectID() (+3 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "strategies.go"
Cohesion: 0.24
Nodes (8): ComboStep, ProviderStrategyState, decodeStepsJSON(), fastRandFloat64(), fastRandIntN(), newStrategyState(), stepsFromTargets(), validateProviderStrategy()

### Community 201 - "listSettingsWithQueryer"
Cohesion: 0.24
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.29
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "GatewayInterceptor"
Cohesion: 0.20
Nodes (10): Global Constraints, Review Focus, Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs) (+2 more)

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 209 - "Open"
Cohesion: 0.05
Nodes (59): roundTripFunc, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+51 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 211 - "repo"
Cohesion: 0.27
Nodes (6): buildWhere(), newEventID(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 219 - ".Fetch"
Cohesion: 0.20
Nodes (5): TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), isCodeBuddyHostFor(), resolveCodeBuddyQuotaURL(), isTokenHarborHost(), codebuddyQuotaFetcher

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 222 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

### Community 223 - "Load"
Cohesion: 0.25
Nodes (9): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable() (+1 more)

### Community 224 - "quotaCache"
Cohesion: 0.33
Nodes (6): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry

### Community 225 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 226 - "plugins/interceptor.go"
Cohesion: 0.43
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 227 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 230 - "redact.go"
Cohesion: 0.46
Nodes (6): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 231 - "time.Duration"
Cohesion: 0.06
Nodes (9): stubInterceptor, stubProviderResolver, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID() (+1 more)

### Community 233 - "manifest.test.ts"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 235 - "Review Focus"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 241 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 242 - "compressor.go"
Cohesion: 0.13
Nodes (27): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+19 more)

### Community 243 - "setupNotificationsTest"
Cohesion: 0.47
Nodes (5): Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict()

### Community 244 - "headroom.go"
Cohesion: 0.33
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 246 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 247 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 248 - "NewChromeTLSClient"
Cohesion: 0.40
Nodes (4): bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton()

### Community 249 - "Review Focus"
Cohesion: 0.50
Nodes (4): Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification

### Community 250 - "HistoryItem"
Cohesion: 0.67
Nodes (3): HistoryItem, historyListResp, historyResponse

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1044 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `Button` to `Module`?**
  _High betweenness centrality (0.215) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06918918918918919 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `net/http.Request`, `Button`, `go_pkg_encoding_json`, `UpstreamCredentials`, `Module`, `.Fetch`?**
  _High betweenness centrality (0.177) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `setupSecurityTest` to `ProviderStore`, `WriteJSON`, `Button`, `net/http/httptest.ResponseRecorder`, `DB`, `.Reader`, `SessionSigner`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._