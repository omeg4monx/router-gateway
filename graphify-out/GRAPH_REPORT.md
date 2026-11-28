# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 490 files · ~541,895 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5665 nodes · 18748 edges · 238 communities (211 shown, 27 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1811 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `83b44c99`
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
- repository
- TypeScript Strict Compiler Config
- ProvidersConsole.tsx
- ModelEditorConsole.tsx
- testing.T
- shadcn/ui Component Configuration
- UpstreamCredentials
- MediaConsole.tsx
- App.tsx
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- TelegramConsole.tsx
- updates
- KeysConsole.tsx
- buildApp
- with
- jobs
- net/http/httptest.ResponseRecorder
- steps
- web-check
- on
- .Reader
- TokenSaverConsole.tsx
- ProviderDetailView.tsx
- LiveMonitorConsole.tsx
- docker-build
- go-check
- combos/schemas.ts
- Module
- Go Module Root Definition
- CanonicalUsage
- Open
- QuotaConsole.tsx
- CanonicalStreamChunk
- main.go
- ir.go
- translate_gemini.go
- go_pkg_strings
- repository
- PluginsStudio.tsx
- Module
- ModelDescriptor
- newTestHarness
- time.Duration
- translate_responses.go
- navigation.ts
- .Fetch
- combos/schemas.ts
- PromptGuardConsole.tsx
- translate_anthropic.go
- CanonicalChatResponse
- migrator.go
- Button
- extractErrorMessage
- UsageConsole.tsx
- repository
- setupTunnelTest
- teeHandler
- KeyInfo
- TestContractInterfacesSatisfied
- WriteJSON
- pricingStore
- net/http.Client
- translate_openai.go
- repo
- Module
- redact.go
- go_pkg_encoding_json
- arena/schemas.ts
- combos/module.go
- oauthLoopbackProxy
- readLegacySQLite
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- Module
- paths
- NotificationsConsole.tsx
- SystemOverview.tsx
- net/http.Request
- settingDefinition
- modeleditor/hooks.ts
- time.Time
- sync.Mutex
- RouteMounts
- Module
- SettingsConsole.tsx
- New
- Runner
- strategyState
- models/hooks.ts
- Module
- awsQuickRuntime
- setupModelEditorTest
- repository
- SecurityConsole.tsx
- Module
- Provider Quota Menu & Token Harbor Quota Integration
- post
- .identifyProviderModel
- headroomCache
- Review Focus
- ProviderTopology.tsx
- ConnectionSummary
- repository
- repo
- properties
- evaluateAntigravityQuota
- net/url.URL
- newTestHarness
- Module
- GatewayInterceptor
- usage/hooks.ts
- plugins/repo.go
- properties
- Module
- Provider Detail: Proportional Model Grid & Per-Model Test
- .ProvisionKey
- io/fs.FS
- .UpdateCombo
- OMEG4MON Router
- buildModelProbeRequest
- OverviewRealtimeTraffic.tsx
- certManager
- repo
- Module
- net/http.Header
- properties
- AttestationManager
- release-binaries
- steps
- .executeCandidateAttempt
- repository
- trafficShaper
- ref_bun_test
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- scanCombo
- properties
- 200
- collectDynamicHostTelemetry
- properties
- TranslatorConsole.tsx
- cli/package.json
- setupTestMitm
- ResolvedTarget
- NewClient
- get
- scripts
- post
- schemas
- ProviderStore
- PreviewPluginsResponse
- WithProviderID
- Module
- stubSystemModule
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- setupTokenSaverTest
- Review Focus
- Review Focus
- ChatCompleter
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- setupCatalogTestModule
- stubSystemModule
- parameters
- .ListModels
- New
- responseCoalescer
- ClearImporterHistoryResponse
- source_path
- setupSecurityTest
- net/http.Handler
- docker-release
- Module
- `@omeg4mon/router` CLI Launcher
- Load
- .ExportCombos
- go_pkg_encoding_json
- context.Context
- pricingStore
- vite.config.ts
- parameters
- SessionSigner
- Review Focus
- omeg4mon-router Verification Gate
- NewStreamTranslator
- redact.go
- single-binary
- keyLimiter
- verify-web
- setupMediaTest
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- omeg4mon-router Project Standards
- stubModule
- Feature Suite
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- stubSystemModule
- stubModule
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- Global Constraints
- schema.go
- DB
- quotaCache
- mapBusEventToAudit
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- time.Duration
- jsonMarshalerError
- compressor.go

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 189 edges
2. `WriteJSON()` - 186 edges
3. `WriteError()` - 184 edges
4. `DB` - 93 edges
5. `Button` - 93 edges
6. `cn()` - 85 edges
7. `buildApp()` - 79 edges
8. `Module` - 78 edges
9. `Bus` - 72 edges
10. `Alert` - 71 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD)` --references--> `ProviderStudioDetail()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-provider-detail-model-test.md → web/src/features/providers/ProviderDetailView.tsx
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Approaches Considered` --references--> `quotaFetcher`  [INFERRED]
  docs/superpowers/specs/2026-11-27-provider-quota-design.md → internal/modules/providers/quota.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go

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

## Communities (238 total, 27 thin omitted)

### Community 0 - "main_test.go"
Cohesion: 0.05
Nodes (36): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+28 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.23
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
Cohesion: 0.12
Nodes (16): cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, scanSummary(), createParams, deleteResult, keySecretResp (+8 more)

### Community 6 - "repository"
Cohesion: 0.05
Nodes (48): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+40 more)

### Community 7 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "ProvidersConsole.tsx"
Cohesion: 0.07
Nodes (57): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), buildProviderSectionBuckets(), canAddProviderConnection(), canProbeProviderModels(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets() (+49 more)

### Community 10 - "ModelEditorConsole.tsx"
Cohesion: 0.11
Nodes (37): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, usePreviewModelEditorMutation(), useRunQuickPromptTestMutation(), EMPTY_FORM_VALUES, FieldLabel(), FieldLabelProps, MobileTab (+29 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (85): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+77 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.12
Nodes (9): catalogFakeProviderStore, ProviderConnection, stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), DefaultCredentials(), refreshCall (+1 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.06
Nodes (46): App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio, InitPage (+38 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.08
Nodes (41): isBlockedAntigravityProjectID(), generateTraceID(), resolveAwsQuickTenantBaseURL(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), originOfURL(), postExchangeIFlow(), postExchangeNous() (+33 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (72): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isIpLiteral(), isKeyExpiredAt() (+64 more)

### Community 21 - "buildApp"
Cohesion: 0.05
Nodes (50): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+42 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.11
Nodes (19): fakeProviderStore, testHarness, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+11 more)

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
Cohesion: 0.12
Nodes (17): main(), runMain(), Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+9 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 30 - "ProviderDetailView.tsx"
Cohesion: 0.03
Nodes (151): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, @hookform/resolvers, runModelIdentityProbe() (+143 more)

### Community 31 - "LiveMonitorConsole.tsx"
Cohesion: 0.17
Nodes (21): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+13 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "combos/schemas.ts"
Cohesion: 0.08
Nodes (52): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+44 more)

### Community 35 - "Module"
Cohesion: 0.05
Nodes (36): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), buildAwsQuickSigninURL() (+28 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.10
Nodes (20): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+12 more)

### Community 39 - "Open"
Cohesion: 0.05
Nodes (55): roundTripFunc, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+47 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (47): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+39 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 42 - "main.go"
Cohesion: 0.17
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 43 - "ir.go"
Cohesion: 0.16
Nodes (32): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+24 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "go_pkg_strings"
Cohesion: 0.11
Nodes (17): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+9 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.16
Nodes (30): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+22 more)

### Community 48 - "Module"
Cohesion: 0.21
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 49 - "ModelDescriptor"
Cohesion: 0.31
Nodes (6): buildServiceWorkerSource(), CORE_PRECACHE_URLS, classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "newTestHarness"
Cohesion: 0.09
Nodes (24): File Map, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, KeyValidator, assertFixedWidthUTC(), newTestHarness(), sha256Hex() (+16 more)

### Community 51 - "time.Duration"
Cohesion: 0.13
Nodes (8): stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore(), computeJitteredRetryDelay(), formatDuration()

### Community 52 - "translate_responses.go"
Cohesion: 0.10
Nodes (34): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+26 more)

### Community 53 - "navigation.ts"
Cohesion: 0.11
Nodes (26): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ParameterGuideEntry, QuickTestPreset, UILocale (+18 more)

### Community 54 - ".Fetch"
Cohesion: 0.12
Nodes (21): TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+13 more)

### Community 55 - "combos/schemas.ts"
Cohesion: 0.09
Nodes (27): COMBO_STRATEGIES, CreateComboFormInput, CreateComboFormOutput, createComboSchema, DEFAULT_CREATE_COMBO_VALUES, DEFAULT_DISABLE_MODEL_VALUES, DEFAULT_UPSERT_ALIAS_VALUES, DisableModelFormInput (+19 more)

### Community 56 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 57 - "translate_anthropic.go"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 58 - "CanonicalChatResponse"
Cohesion: 0.14
Nodes (20): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString(), RewriteModelField() (+12 more)

### Community 59 - "migrator.go"
Cohesion: 0.12
Nodes (31): ModelOption, inspectRequest, runImportRequest, Module, ChangedEvent, createKeyRequest, nullableString, patchKeyRequest (+23 more)

### Community 60 - "Button"
Cohesion: 0.06
Nodes (98): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, lucide-react, react, react-router, ArenaConsole(), useArenaBattles(), useArenaLeaderboard() (+90 more)

### Community 61 - "extractErrorMessage"
Cohesion: 0.05
Nodes (86): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCreateComboMutation(), useDeleteComboMutation(), useDeleteModelAliasMutation(), useDisabledModels() (+78 more)

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.17
Nodes (23): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), computeBarPercent(), computeSharePercent(), computeTotalTokens(), DEFAULT_LOGS_LIMIT, formatCompactTokens(), formatDurationMs(), formatErrorRate() (+15 more)

### Community 63 - "repository"
Cohesion: 0.16
Nodes (8): decodeStringMap(), decodeStringSlice(), newRepository(), scanConnectionRecord(), scanNode(), ProviderNode, repository, rowScanner

### Community 64 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "KeyInfo"
Cohesion: 0.15
Nodes (11): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), newKeyCache(), hashToken(), isValidTokenFormat() (+3 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.07
Nodes (10): stubAdminAuth, AuthVerifier, stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth, stubAdminAuth (+2 more)

### Community 68 - "WriteJSON"
Cohesion: 0.12
Nodes (7): decodeStrictJSON(), Module, WriteJSON(), LogEntry, LogListResponse, Module, Status

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "net/http.Client"
Cohesion: 0.16
Nodes (12): selectProbeHTTPClient(), newBotAPIClient(), sanitizeBotError(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar() (+4 more)

### Community 71 - "translate_openai.go"
Cohesion: 0.21
Nodes (17): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+9 more)

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "redact.go"
Cohesion: 0.18
Nodes (12): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), newEventID(), normalizeSeverity() (+4 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.14
Nodes (27): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+19 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (46): ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaModels(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteChatSessionMutation() (+38 more)

### Community 77 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 78 - "oauthLoopbackProxy"
Cohesion: 0.16
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "web/package.json"
Cohesion: 0.12
Nodes (17): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+9 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.25
Nodes (6): Approaches Considered, Builder Decisions & Deliverables, Constraints, Goals, Non-Goals & Anti-Goals, Intent & Why, PWA Support Design — `omeg4mon-router`

### Community 82 - "Module"
Cohesion: 0.21
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 83 - "Module"
Cohesion: 0.16
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 86 - "SystemOverview.tsx"
Cohesion: 0.05
Nodes (67): Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, authKeys, AuthMeResponse (+59 more)

### Community 87 - "net/http.Request"
Cohesion: 0.05
Nodes (36): Module, Module, trackedResponseWriter, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), clearSessionCookie() (+28 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "modeleditor/hooks.ts"
Cohesion: 0.16
Nodes (16): invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation(), useModelOverrides(), useUpdateModelOverrideMutation() (+8 more)

### Community 90 - "time.Time"
Cohesion: 0.10
Nodes (14): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock() (+6 more)

### Community 91 - "sync.Mutex"
Cohesion: 0.15
Nodes (7): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher

### Community 92 - "RouteMounts"
Cohesion: 0.15
Nodes (4): Module, validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 93 - "Module"
Cohesion: 0.15
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (42): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation() (+34 more)

### Community 95 - "New"
Cohesion: 0.11
Nodes (25): newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify(), VerifyPassword(), apiTimeoutMiddleware() (+17 more)

### Community 96 - "Runner"
Cohesion: 0.14
Nodes (8): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "strategyState"
Cohesion: 0.20
Nodes (8): strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), fastRandFloat64(), fastRandIntN(), newStrategyState()

### Community 98 - "models/hooks.ts"
Cohesion: 0.26
Nodes (14): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, ModelIdentityResult (+6 more)

### Community 99 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 100 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 101 - "setupModelEditorTest"
Cohesion: 0.15
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestConcurrentExtraParamsNoDataRaceOrCacheMutation() (+10 more)

### Community 102 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 103 - "SecurityConsole.tsx"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 105 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.18
Nodes (9): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+1 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.10
Nodes (23): ModelIdentityResult, ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+15 more)

### Community 109 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (20): @xyflow/react, EDGE_TYPES, EnergyTopologyEdge(), FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, prefersReducedMotion(), ProviderTopology() (+12 more)

### Community 111 - "ConnectionSummary"
Cohesion: 0.13
Nodes (17): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), oauthEnvPrefix(), isCodeBuddyHostFor(), isTokenHarborHost(), encodeStringMap(), encodeStringSlice() (+9 more)

### Community 112 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 116 - "net/url.URL"
Cohesion: 0.07
Nodes (22): ProxyDialerResolver, stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+14 more)

### Community 117 - "newTestHarness"
Cohesion: 0.15
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 118 - "Module"
Cohesion: 0.10
Nodes (18): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module (+10 more)

### Community 119 - "GatewayInterceptor"
Cohesion: 0.16
Nodes (14): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI) (+6 more)

### Community 120 - "usage/hooks.ts"
Cohesion: 0.17
Nodes (13): buildLogsQueryParams(), INVALIDATION_TOPIC_QUERY_KEYS, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useUsageLogs(), isUsageLogError(), UsageLogsResponse (+5 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.06
Nodes (59): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+51 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "Module"
Cohesion: 0.05
Nodes (19): streamResponseBody(), New(), newID(), newRepo(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord (+11 more)

### Community 124 - "Provider Detail: Proportional Model Grid & Per-Model Test"
Cohesion: 0.18
Nodes (8): Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isModelAllowed(), isModelDisabled()

### Community 125 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 126 - "io/fs.FS"
Cohesion: 0.03
Nodes (16): stubDep, stubModule, migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile() (+8 more)

### Community 127 - ".UpdateCombo"
Cohesion: 0.23
Nodes (12): ComboStep, createComboParams, updateComboValues, cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps() (+4 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 130 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.24
Nodes (14): formatActivityOriginLabel(), formatClockUptime(), LazyProviderTopology, OverviewRealtimeTraffic, OverviewRealtimeTrafficProps, RecentRequestsRows(), computeRelativeClockIntervalMs(), FilteredRecentRequestsResult (+6 more)

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "Module"
Cohesion: 0.19
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 134 - "net/http.Header"
Cohesion: 0.13
Nodes (9): memoryResponseWriter, extractVendorHeaders(), extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration() (+1 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "AttestationManager"
Cohesion: 0.18
Nodes (11): ApplyHeaders(), DefaultAttestationManager(), AttestationManager, isTransientMintStatus(), NewAttestationManager(), NewNDJSONToOpenAISSEReader(), ResolveMintURL(), TestNewNDJSONToOpenAISSEReader() (+3 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - ".executeCandidateAttempt"
Cohesion: 0.11
Nodes (23): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), extractUpstreamErrorMessage(), parseRetryAfter(), applyCodeBuddyHeaders() (+15 more)

### Community 140 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 141 - "trafficShaper"
Cohesion: 0.24
Nodes (7): trafficShaper, TestTrafficShaperAdaptiveHeaderPacing(), TestTrafficShaperLimitsPerConnectionConcurrencyAndQueues(), TestTrafficShaperShedsLoadWhenGlobalFull(), acquireSem(), newTrafficShaper(), sleepContext()

### Community 142 - "ref_bun_test"
Cohesion: 0.05
Nodes (38): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+30 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "TranslatorConsole.tsx"
Cohesion: 0.05
Nodes (55): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form (+47 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 153 - "ResolvedTarget"
Cohesion: 0.17
Nodes (9): resolveSnapshot, stubProviderResolver, stubResolver, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), prioritizeHeadroomAffinity() (+1 more)

### Community 154 - "NewClient"
Cohesion: 0.15
Nodes (11): arenaTestEnv, stubCompleter, Options, startAuthHTTPProxy(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(), applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout() (+3 more)

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

### Community 159 - "ProviderStore"
Cohesion: 0.32
Nodes (13): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI) (+5 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "WithProviderID"
Cohesion: 0.15
Nodes (12): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, isBlockedAntigravityProjectID() (+4 more)

### Community 162 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

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

### Community 170 - "ChatCompleter"
Cohesion: 0.21
Nodes (10): ChatCompletionRequest, ChatCompletionResult, ChatMessage, Global Constraints, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan (+2 more)

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

### Community 177 - "setupCatalogTestModule"
Cohesion: 0.17
Nodes (9): setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies(), TestAllOmniRouteComboStrategiesCRUDAndPersistence(), TestOmniRouteResolverStrategiesExecution() (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - ".ListModels"
Cohesion: 0.12
Nodes (8): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, testEnv, TestListVirtualModels(), compositeModelLister, stubModelLister

### Community 181 - "New"
Cohesion: 0.07
Nodes (33): File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator (+25 more)

### Community 182 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "setupSecurityTest"
Cohesion: 0.08
Nodes (23): TestAppCloseTeardownOrderIsLIFO(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune() (+15 more)

### Community 186 - "net/http.Handler"
Cohesion: 0.13
Nodes (19): loadUI(), loadUI(), New(), setupPluginsTest(), TestPluginsBuiltinsUpdateAndInterceptorScoping(), newRepository(), putPluginConfig(), TestPluginsAntiSlopFormatAwareCachePreservingAndIdempotent() (+11 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.21
Nodes (8): errorDisposition, handledGatewayError, classifyStatus(), isClientPayloadErrorBody(), compareSemver(), maybeLiftFreePIClientVersion(), parseFreePIRequiredVersion(), TestFreePIParseRequiredVersionAndCompare()

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (21): fakeCredentialStore, AdminCredentialStore, SettingsSnapshotter, stubAdminCredentialStore, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, stubSettings, newFakeStore(), extractUsageCost() (+13 more)

### Community 195 - "vite.config.ts"
Cohesion: 0.29
Nodes (6): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

### Community 198 - "Review Focus"
Cohesion: 0.25
Nodes (8): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, omeg4monPwaPlugin(), ThemeMode

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 201 - "redact.go"
Cohesion: 0.06
Nodes (26): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), listSettingsWithQueryer() (+18 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.29
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "setupMediaTest"
Cohesion: 0.36
Nodes (6): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction()

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "Global Constraints"
Cohesion: 0.33
Nodes (5): Global Constraints, Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), Wave 3C — Local MITM Proxy & PXPIPE Transform Pipeline Implementation Plan

### Community 222 - "schema.go"
Cohesion: 0.33
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 223 - "DB"
Cohesion: 0.05
Nodes (61): Options, authHarness, Options, Options, Global Constraints, Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager) (+53 more)

### Community 224 - "quotaCache"
Cohesion: 0.33
Nodes (6): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry

### Community 225 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 231 - "time.Duration"
Cohesion: 0.05
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID() (+1 more)

### Community 242 - "compressor.go"
Cohesion: 0.12
Nodes (28): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+20 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1060 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **27 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProviderDetailView.tsx` to `Module`?**
  _High betweenness centrality (0.230) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05266106442577031 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `go_pkg_encoding_json`, `testing.T`, `UpstreamCredentials`, `oauth_registry.go`, `.Fetch`, `net/http.Request`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.193) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `ProviderDetailView.tsx` to `ProvidersConsole.tsx`, `combos/schemas.ts`, `Button`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._