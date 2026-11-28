# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 488 files · ~544,073 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5643 nodes · 18704 edges · 237 communities (202 shown, 35 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1810 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `001bc310`
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
- Bus
- TypeScript Strict Compiler Config
- ProvidersModals.tsx
- ModelEditorConsole.tsx
- testing.T
- shadcn/ui Component Configuration
- UpstreamCredentials
- MediaConsole.tsx
- App.tsx
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- extractErrorMessage
- updates
- KeysConsole.tsx
- buildApp
- with
- jobs
- net/http/httptest.ResponseRecorder
- steps
- web-check
- on
- GatewayInterceptor
- TokenSaverConsole.tsx
- builtinAdapter
- NotificationsConsole.tsx
- docker-build
- go-check
- Badge
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
- TelegramConsole.tsx
- ModelDescriptor
- DB
- ProvidersConsole.tsx
- translate_responses.go
- providers/hooks.ts
- .Fetch
- Module
- PromptGuardConsole.tsx
- translate_anthropic.go
- repository
- migrator.go
- Button
- MitmConsole.tsx
- ref_bun_test
- repository
- net/http.Header
- teeHandler
- ProxyPoolsConsole.tsx
- TestContractInterfacesSatisfied
- Module
- pricingStore
- modelIntelligence.ts
- translate_openai.go
- repo
- Module
- promptinjection/engine.go
- go_pkg_encoding_json
- arena/schemas.ts
- combos/module.go
- oauthLoopbackProxy
- readLegacySQLite
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- fakeHandle
- paths
- ProviderDetailView.tsx
- modalUtils.ts
- net/http.ResponseWriter
- settingDefinition
- repo
- time.Time
- UsageRecord
- CanonicalChatResponse
- net/http.Client
- SettingsConsole.tsx
- New
- sync.Mutex
- fetch.go
- ref_bun_test
- Module
- awsQuickRuntime
- setupModelEditorTest
- repository
- SecurityConsole.tsx
- settingDefinition
- Module
- post
- .identifyProviderModel
- repo
- New
- ProviderTopology.tsx
- net/http.Request
- repository
- repo
- properties
- evaluateAntigravityQuota
- time.Duration
- newTestHarness
- Module
- ProcessHandle
- Provider Quota Menu & Token Harbor Quota Integration
- plugins/repo.go
- properties
- Module
- Provider Detail: Proportional Model Grid & Per-Model Test
- .ProvisionKey
- io/fs.FS
- repository
- OMEG4MON Router
- batcher
- setupMediaTest
- certManager
- repo
- buildModelProbeRequest
- headroomCache
- properties
- tgTestEnv
- release-binaries
- steps
- .executeCandidateAttempt
- repo
- redact.go
- deviceOAuthSession
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- .StatusCode
- properties
- 200
- collectDynamicHostTelemetry
- properties
- dependencies
- cli/package.json
- setupTestMitm
- vite.config.ts
- .handleOAuthPaste
- get
- scripts
- post
- schemas
- repository
- PreviewPluginsResponse
- NewStreamTranslator
- Module
- stubSystemModule
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- setupTokenSaverTest
- Review Focus
- Review Focus
- stubDep
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- stubModule
- WriteJSON
- parameters
- Module
- New
- responseCoalescer
- ClearImporterHistoryResponse
- source_path
- stubModule
- runtime_test.go
- docker-release
- Module
- `@omeg4mon/router` CLI Launcher
- setupTunnelTest
- .ExportCombos
- go_pkg_encoding_json
- context.Context
- setupNotificationsTest
- stubModule
- parameters
- stubDiscoverySystemModule
- stubSystemModule
- omeg4mon-router Verification Gate
- stubDep
- listSettingsWithQueryer
- single-binary
- keyLimiter
- verify-web
- .list
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- repository
- stubSystemModule
- Feature Suite
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- io.Reader
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- stubDepModule
- stubSystem
- Config
- quotaCache
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- mapBusEventToAudit
- time.Duration
- jsonMarshalerError
- compressor.go

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 189 edges
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
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Approaches Considered` --references--> `quotaFetcher`  [INFERRED]
  docs/superpowers/specs/2026-11-27-provider-quota-design.md → internal/modules/providers/quota.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx
- `Reliability, Memory & Type Safety` --references--> `extractErrorMessage()`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/api/client.ts

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

## Communities (237 total, 35 thin omitted)

### Community 0 - "main_test.go"
Cohesion: 0.08
Nodes (13): collectStaticHostInfo(), readWindowsRegistryString(), isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive(), healthResponse, initCompleteResponse, initRequest (+5 more)

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
Cohesion: 0.07
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+19 more)

### Community 6 - "repository"
Cohesion: 0.09
Nodes (27): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+19 more)

### Community 7 - "Bus"
Cohesion: 0.07
Nodes (24): Options, New(), New(), New(), New(), clientIP(), decodeStrictJSON(), Module (+16 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (59): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), AddApiKeyModal(), AddCustomNodeModal(), DraftTestResultCard(), getCredentialFieldMeta(), MODAL_SIZE_CLASS (+51 more)

### Community 10 - "ModelEditorConsole.tsx"
Cohesion: 0.06
Nodes (71): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation() (+63 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (99): TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), decodeJSON(), newTestHarness() (+91 more)

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
Nodes (43): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety, @tanstack/react-query, @tanstack/react-query@^5.76.1, App() (+35 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.17
Nodes (22): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), decodeJWTPayload(), enrichGeminiOAuth(), enrichGitHubCopilot(), extractCodexAccountInfo() (+14 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "extractErrorMessage"
Cohesion: 0.05
Nodes (75): openapi-fetch, openapi-fetch@^0.14.0, combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCreateComboMutation(), useDeleteComboMutation() (+67 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.06
Nodes (71): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isIpLiteral(), isKeyExpiredAt() (+63 more)

### Community 21 - "buildApp"
Cohesion: 0.06
Nodes (33): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog(), TestWave2BDisabledModulesSkipInterceptors() (+25 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.08
Nodes (36): authHarness, fakeCredentialStore, Options, Claims, File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter (+28 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "GatewayInterceptor"
Cohesion: 0.04
Nodes (65): main(), ChatCompletionResult, File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment (+57 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 30 - "builtinAdapter"
Cohesion: 0.24
Nodes (7): cloneOAuthConfig(), isAnthropicTokenURL(), defaultExchangePKCECode(), exchangeClaudeCode(), builtinAdapter, OAuthConfig, oauthSessionRecord

### Community 31 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "Badge"
Cohesion: 0.05
Nodes (70): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+62 more)

### Community 35 - "Module"
Cohesion: 0.05
Nodes (36): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), oauthEnvPrefix(), validateHTTPURL() (+28 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (18): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+10 more)

### Community 39 - "Open"
Cohesion: 0.04
Nodes (61): Options, roundTripFunc, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest() (+53 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.07
Nodes (59): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit (+51 more)

### Community 43 - "ir.go"
Cohesion: 0.15
Nodes (29): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, StreamTranslator (+21 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.17
Nodes (20): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+12 more)

### Community 45 - "go_pkg_strings"
Cohesion: 0.12
Nodes (19): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+11 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.10
Nodes (46): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, modelCatalogQueryKeys (+38 more)

### Community 48 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 49 - "ModelDescriptor"
Cohesion: 0.26
Nodes (8): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin(), classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "DB"
Cohesion: 0.05
Nodes (42): Options, fakeProviderStore, testHarness, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, chainStore, New(), newRepo() (+34 more)

### Community 51 - "ProvidersConsole.tsx"
Cohesion: 0.08
Nodes (52): AddCustomModelModalProps, ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), buildProviderSectionBuckets(), canAddProviderConnection(), canProbeProviderModels(), computeEntryStats(), computeFleetTelemetry() (+44 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.10
Nodes (32): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+24 more)

### Community 53 - "providers/hooks.ts"
Cohesion: 0.08
Nodes (30): ARENA_MODELS_QUERY_KEY, DiscoveredModel, ExchangeProviderOAuthVariables, FetchConnectionModelsResponse, fetchProviderOAuthStatus(), MODEL_CATALOG_QUERY_KEY, PasteProviderOAuthVariables, providersQueryKeys (+22 more)

### Community 54 - ".Fetch"
Cohesion: 0.11
Nodes (23): TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+15 more)

### Community 55 - "Module"
Cohesion: 0.13
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), New(), toJobDTO(), JobStatusDTO

### Community 56 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 57 - "translate_anthropic.go"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 58 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 59 - "migrator.go"
Cohesion: 0.09
Nodes (44): ModelOption, PolicyChangedEvent, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, buildInstallGuide() (+36 more)

### Community 60 - "Button"
Cohesion: 0.05
Nodes (92): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, class-variance-authority, lucide-react (+84 more)

### Community 61 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 62 - "ref_bun_test"
Cohesion: 0.05
Nodes (81): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, useCombos(), KeyedSystemLog, levelColorClass() (+73 more)

### Community 63 - "repository"
Cohesion: 0.09
Nodes (22): cloneStringMap(), cloneStringSlice(), isCodeBuddyHostFor(), isTokenHarborHost(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice() (+14 more)

### Community 64 - "net/http.Header"
Cohesion: 0.13
Nodes (9): memoryResponseWriter, extractVendorHeaders(), extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration() (+1 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (35): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+27 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.07
Nodes (10): stubAdminAuth, AuthVerifier, stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth, stubAdminAuth (+2 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "translate_openai.go"
Cohesion: 0.19
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+12 more)

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "promptinjection/engine.go"
Cohesion: 0.11
Nodes (21): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+13 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (28): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+20 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (46): ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaModels(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteChatSessionMutation() (+38 more)

### Community 77 - "combos/module.go"
Cohesion: 0.07
Nodes (19): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+11 more)

### Community 78 - "oauthLoopbackProxy"
Cohesion: 0.15
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "web/package.json"
Cohesion: 0.10
Nodes (25): @biomejs/biome, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react (+17 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.17
Nodes (10): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Approaches Considered, Builder Decisions & Deliverables, Constraints, Goals, Non-Goals & Anti-Goals (+2 more)

### Community 82 - "Module"
Cohesion: 0.20
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 83 - "fakeHandle"
Cohesion: 0.11
Nodes (5): NewNDJSONToOpenAISSEReader(), TestNewNDJSONToOpenAISSEReader(), fakeHandle, osProcessHandle, ssePipeReader

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "ProviderDetailView.tsx"
Cohesion: 0.08
Nodes (73): ArenaConsole(), useArenaBattles(), useArenaLeaderboard(), useDeleteArenaBattleMutation(), useRunArenaBattleMutation(), useVoteArenaBattleMutation(), KeyFormModal(), runModelIdentityProbe() (+65 more)

### Community 86 - "modalUtils.ts"
Cohesion: 0.17
Nodes (19): Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useTestProviderNodeDraftMutation(), buildTestFingerprint(), BulkCredentialEntry (+11 more)

### Community 87 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (16): Module, trackedResponseWriter, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), ensureRequestID() (+8 more)

### Community 88 - "settingDefinition"
Cohesion: 0.13
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.17
Nodes (9): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotificationRow(), EmitPayload, Notification (+1 more)

### Community 90 - "time.Time"
Cohesion: 0.09
Nodes (14): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, newLoginRateLimiter() (+6 more)

### Community 91 - "UsageRecord"
Cohesion: 0.11
Nodes (12): stubModelLister, matchesAllowedModel(), stubUsageRecorder, stubModelLister, stubUsageRecorder, stubModelLister, stubUsageRecorder, testEnv (+4 more)

### Community 92 - "CanonicalChatResponse"
Cohesion: 0.12
Nodes (26): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString() (+18 more)

### Community 93 - "net/http.Client"
Cohesion: 0.15
Nodes (14): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar(), AttestationManager, isTransientMintStatus() (+6 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (46): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, zod, zod@^4.6.7, authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse (+38 more)

### Community 95 - "New"
Cohesion: 0.07
Nodes (32): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+24 more)

### Community 96 - "sync.Mutex"
Cohesion: 0.11
Nodes (10): JobInspector, JobStatus, SettingsSnapshotter, stubSettings, Runner, job, JobStatus, stubJobInspector (+2 more)

### Community 97 - "fetch.go"
Cohesion: 0.29
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 98 - "ref_bun_test"
Cohesion: 0.06
Nodes (32): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+24 more)

### Community 99 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 100 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 101 - "setupModelEditorTest"
Cohesion: 0.14
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupModelEditorTest() (+11 more)

### Community 102 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 103 - "SecurityConsole.tsx"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 104 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 105 - "Module"
Cohesion: 0.15
Nodes (5): actorFromContext(), decodeJSONBody(), Module, AdminSubject(), RouteMounts

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.09
Nodes (27): ModelIdentityResult, ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+19 more)

### Community 108 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 109 - "New"
Cohesion: 0.17
Nodes (11): isBlockedAntigravityProjectID(), newDeviceOAuthRuntime(), New(), enrichAntigravity(), TestEnrichAntigravityAutoOnboardsNewAccountAndSendsIDEFingerprint(), TestEnrichAntigravityFiltersBlockedProjectAndSetsHeader(), TestEnrichGitHubCopilotAndKiroProfileARN(), TestOAuthProviderRegistryMetadata() (+3 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.15
Nodes (18): @xyflow/react, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps, RouterHubNode() (+10 more)

### Community 111 - "net/http.Request"
Cohesion: 0.07
Nodes (14): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), extractAPIKey(), decodeJSON() (+6 more)

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

### Community 116 - "time.Duration"
Cohesion: 0.05
Nodes (26): ProxyDialerResolver, stubProxyDialerResolver, failureRecord, stubHealthTracker, trafficShaper, Options, stubProxyResolver, computeBackoff() (+18 more)

### Community 117 - "newTestHarness"
Cohesion: 0.13
Nodes (20): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Options, New(), decodeJSON() (+12 more)

### Community 118 - "Module"
Cohesion: 0.10
Nodes (18): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module (+10 more)

### Community 119 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 120 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.18
Nodes (9): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+1 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.06
Nodes (59): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+51 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "Module"
Cohesion: 0.25
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 124 - "Provider Detail: Proportional Model Grid & Per-Model Test"
Cohesion: 0.18
Nodes (8): Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isModelAllowed(), isModelDisabled()

### Community 125 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 126 - "io/fs.FS"
Cohesion: 0.07
Nodes (11): migrationFile, Module, validateSQLiteSource(), decodeJSON(), collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath() (+3 more)

### Community 127 - "repository"
Cohesion: 0.05
Nodes (37): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, resolveSnapshot, rowScanner (+29 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

### Community 130 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 131 - "certManager"
Cohesion: 0.16
Nodes (6): newCertManager(), randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "tgTestEnv"
Cohesion: 0.29
Nodes (4): arenaTestEnv, stubCompleter, tgAPICall, tgTestEnv

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - ".executeCandidateAttempt"
Cohesion: 0.07
Nodes (32): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, attemptOutcome (+24 more)

### Community 140 - "repo"
Cohesion: 0.53
Nodes (4): buildWhere(), scanAuditEvent(), AuditEvent, ListFilter

### Community 141 - "redact.go"
Cohesion: 0.22
Nodes (10): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), formatUTC(), newEventID(), newRepo() (+2 more)

### Community 142 - "deviceOAuthSession"
Cohesion: 0.14
Nodes (12): codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), originOfURL(), defaultOAuthStrategies(), newRFC8628Strategy(), testEnv, codebuddyAuthSpec (+4 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, statusError, httpStatusError

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

### Community 150 - "dependencies"
Cohesion: 0.10
Nodes (21): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form, tailwind-merge, zustand, dependencies (+13 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "setupTestMitm"
Cohesion: 0.31
Nodes (10): New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly(), TestMitmLiveTLSProxyInterceptAndAliasRewrite() (+2 more)

### Community 153 - "vite.config.ts"
Cohesion: 0.29
Nodes (6): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0

### Community 154 - ".handleOAuthPaste"
Cohesion: 0.29
Nodes (3): parseAwsQuickDiscovery(), parseAwsQuickPastedQuery(), renderAwsQuickResultHTML()

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

### Community 159 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 162 - "Module"
Cohesion: 0.17
Nodes (7): validateModelPattern(), decodeJSONBody(), newID(), TestListVirtualModels(), validateBypassRule(), validateOverride(), Module

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

### Community 178 - "WriteJSON"
Cohesion: 0.14
Nodes (10): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider() (+2 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 181 - "New"
Cohesion: 0.08
Nodes (32): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+24 more)

### Community 182 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 186 - "runtime_test.go"
Cohesion: 0.33
Nodes (11): New(), setupPluginsTest(), TestPluginsBuiltinsUpdateAndInterceptorScoping(), newRepository(), putPluginConfig(), TestPluginsAntiSlopFormatAwareCachePreservingAndIdempotent(), TestPluginsContextSqueezerTokenAwareEvictionAndElision(), TestPluginsImageVisionOnlyInjectsWhenImagesPresent() (+3 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "setupTunnelTest"
Cohesion: 0.25
Nodes (9): New(), doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary() (+1 more)

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.17
Nodes (11): errorDisposition, handledGatewayError, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus(), isClientPayloadErrorBody(), compareSemver(), maybeLiftFreePIClientVersion() (+3 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (9): AdminCredentialStore, stubAdminCredentialStore, extractUsageCost(), sanitizeBotError(), stubModule, orderStubModule, stubCredentialStore, repository (+1 more)

### Community 194 - "setupNotificationsTest"
Cohesion: 0.15
Nodes (16): TestAppCloseTeardownOrderIsLIFO(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup() (+8 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 201 - "listSettingsWithQueryer"
Cohesion: 0.17
Nodes (12): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), HistoryItem, lookupDefinition(), applyResult, historyListResp, historyResponse (+4 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.24
Nodes (5): newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 223 - "Config"
Cohesion: 0.09
Nodes (21): runMain(), Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), setupSecurityTest(), TestSecurityAuditPaginationOffset(), TestSecurityEventBusAutoAuditIngestion() (+13 more)

### Community 224 - "quotaCache"
Cohesion: 0.33
Nodes (6): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 231 - "time.Duration"
Cohesion: 0.05
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID() (+1 more)

### Community 242 - "compressor.go"
Cohesion: 0.12
Nodes (28): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+20 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1051 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `Module`, `providers/hooks.ts`, `modalUtils.ts`?**
  _High betweenness centrality (0.197) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08244680851063829 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `ProvidersModals.tsx`, `go_pkg_encoding_json`, `testing.T`, `UpstreamCredentials`, `New`, `net/http.Request`, `.Fetch`?**
  _High betweenness centrality (0.168) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `net/http/httptest.ResponseRecorder` to `Module`, `GatewayInterceptor`, `DB`, `Button`, `SettingsConsole.tsx`, `New`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._