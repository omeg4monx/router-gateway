# Graph Report - omeg4mon-router  (2026-11-27)

## Corpus Check
- 468 files · ~504,810 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5460 nodes · 18204 edges · 218 communities (196 shown, 22 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1785 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7051e00a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- server_test.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- repository
- LiveMonitorConsole.tsx
- TypeScript Strict Compiler Config
- ProvidersModals.tsx
- Button
- testing.T
- shadcn/ui Component Configuration
- UpstreamCredentials
- MediaConsole.tsx
- App.tsx
- Module
- Design System & Wave 0 Plan Docs
- MitmConsole.tsx
- updates
- translate_anthropic.go
- PluginsStudio.tsx
- with
- jobs
- WriteJSON
- steps
- web-check
- on
- Module
- go_pkg_encoding_json
- models/hooks.ts
- OverviewRealtimeTraffic.tsx
- docker-build
- go-check
- KeysConsole.tsx
- Module
- Go Module Root Definition
- CanonicalUsage
- newTestHarness
- QuotaConsole.tsx
- CanonicalStreamChunk
- PromptGuardConsole.tsx
- ir.go
- translate_gemini.go
- buildApp
- repository
- Module
- translate_openai.go
- ModelDescriptor
- Module
- ModelsRoutingConsole.tsx
- translate_responses.go
- oauthLoopbackProxy
- mapTokenHarborQuota
- SystemOverview.tsx
- extractErrorMessage
- dependencies
- repository
- server_test.go
- ProcessHandle
- repo
- UsageConsole.tsx
- repository
- repository
- teeHandler
- ProxyPoolsConsole.tsx
- ModelDescriptor
- navigation.ts
- pricingStore
- modelIntelligence.ts
- settingDefinition
- ModelEditorConsole.tsx
- Module
- JobInspector
- go_pkg_encoding_json
- arena/hooks.ts
- Module
- NotificationsConsole.tsx
- WriteError
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- newTestHarness
- paths
- collectDynamicHostTelemetry
- fakeHandle
- Badge
- settingDefinition
- repo
- vite.config.ts
- stubSystem
- redact.go
- SettingsConsole.tsx
- streamCache.ts
- Runner
- contracts.go
- schema.go
- Module
- nameMinHeap
- Module
- Module
- DB
- io.Reader
- ref_bun_test
- post
- .identifyProviderModel
- WriteJSON
- net/http.Header
- ProviderTopology.tsx
- repo
- stubCompleter
- repo
- properties
- ProviderStore
- net/url.URL
- AuditEvent
- deviceOAuthSession
- Module
- time.Time
- .UpsertPlugin
- properties
- time.Duration
- New
- io/fs.FS
- setupMediaTest
- OMEG4MON Router
- main.go
- UsageRecord
- certManager
- repo
- Module
- properties
- release-binaries
- steps
- evaluateAntigravityQuota
- combos/module.go
- newHarness
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- repository
- properties
- 200
- repository
- properties
- Module
- cli/package.json
- plugins/interceptor.go
- get
- scripts
- post
- schemas
- Card
- PreviewPluginsResponse
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- repository
- Review Focus
- Review Focus
- setupTestMitm
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- stubSystem
- version
- WithProviderID
- parameters
- promptinjection/engine.go
- New
- ClearImporterHistoryResponse
- source_path
- docker-release
- stubModule
- `@omeg4mon/router` CLI Launcher
- setupModelEditorTest
- context.Context
- listSettingsWithQueryer
- parameters
- New
- .executeCandidateAttempt
- omeg4mon-router Verification Gate
- selectProbeHTTPClient
- single-binary
- keyLimiter
- verify-web
- modeleditor/engine.go
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- GatewayInterceptor
- Feature Suite
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- setupTunnelTest
- Config
- responsesStreamEncoder
- Load
- manifest.test.ts
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- headroomCache
- time.Duration
- jsonMarshalerError
- compressor.go
- omeg4mon-router Project Standards
- net/http.Request

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 195 edges
2. `WriteJSON()` - 186 edges
3. `WriteError()` - 184 edges
4. `DB` - 93 edges
5. `Button` - 92 edges
6. `buildApp()` - 79 edges
7. `Bus` - 72 edges
8. `cn()` - 72 edges
9. `Alert` - 71 edges
10. `AlertDescription` - 70 edges

## Surprising Connections (you probably didn't know these)
- `Approaches Considered` --references--> `quotaFetcher`  [INFERRED]
  docs/superpowers/specs/2026-11-27-provider-quota-design.md → internal/modules/providers/quota.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx
- `Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`)` --references--> `TestServerIsolatesJSON404FromSPAFallback()`  [INFERRED]
  docs/superpowers/plans/2026-11-27-pwa-support.md → internal/platform/httpserver/server_test.go

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
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (218 total, 22 thin omitted)

### Community 0 - "server_test.go"
Cohesion: 0.10
Nodes (13): collectStaticHostInfo(), readWindowsRegistryString(), isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive(), healthResponse, initCompleteResponse, initRequest (+5 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.20
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
Cohesion: 0.06
Nodes (28): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+20 more)

### Community 6 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.17
Nodes (21): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+13 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (174): Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, AddCustomModelModal(), AddCustomModelModalProps, ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), canAddProviderConnection() (+166 more)

### Community 10 - "Button"
Cohesion: 0.06
Nodes (97): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, class-variance-authority, lucide-react, react, react-router, ArenaConsole() (+89 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (100): stubAdminAuth, testHarness, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh() (+92 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.07
Nodes (24): catalogFakeProviderStore, fakeProviderStore, ModelTestResult, resolveSnapshot, testAllModelsResponse, stubProviderStore, stubProviderStore, ProviderConnection (+16 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (37): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+29 more)

### Community 16 - "Module"
Cohesion: 0.12
Nodes (33): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies() (+25 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "MitmConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "translate_anthropic.go"
Cohesion: 0.17
Nodes (17): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+9 more)

### Community 21 - "PluginsStudio.tsx"
Cohesion: 0.19
Nodes (22): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+14 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "WriteJSON"
Cohesion: 0.11
Nodes (12): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), trackedResponseWriter, providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam() (+4 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "Module"
Cohesion: 0.07
Nodes (5): stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth

### Community 29 - "go_pkg_encoding_json"
Cohesion: 0.37
Nodes (3): errorDisposition, classifyStatus(), isClientPayloadErrorBody()

### Community 30 - "models/hooks.ts"
Cohesion: 0.21
Nodes (17): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, modelCatalogQueryKeys (+9 more)

### Community 31 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.25
Nodes (14): formatActivityOriginLabel(), LazyProviderTopology, OverviewRealtimeTraffic(), OverviewRealtimeTrafficProps, RecentRequestsRows(), ProviderTopologyProps, FilteredRecentRequestsResult, filterRecentRequests() (+6 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "KeysConsole.tsx"
Cohesion: 0.06
Nodes (66): react-hook-form, buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt() (+58 more)

### Community 35 - "Module"
Cohesion: 0.05
Nodes (41): Global Constraints, Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Task 5: Full Verification & Graphify Sync, Approaches Considered, buildProtocolTestRequest(), cloneOAuthConfig() (+33 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+9 more)

### Community 39 - "newTestHarness"
Cohesion: 0.08
Nodes (26): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart() (+18 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.05
Nodes (68): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit (+60 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.38
Nodes (4): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder, mapOpenAIFinishReasonToCanonical()

### Community 42 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 43 - "ir.go"
Cohesion: 0.10
Nodes (47): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat (+39 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "buildApp"
Cohesion: 0.07
Nodes (39): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+31 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "Module"
Cohesion: 0.23
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "ModelDescriptor"
Cohesion: 0.31
Nodes (6): buildServiceWorkerSource(), CORE_PRECACHE_URLS, classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "Module"
Cohesion: 0.19
Nodes (3): Module, validatePassword(), testHarness

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (26): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+18 more)

### Community 53 - "oauthLoopbackProxy"
Cohesion: 0.08
Nodes (16): newAwsQuickRuntime(), validateAwsRegion(), newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), awsQuickDiscovery, awsQuickResult (+8 more)

### Community 54 - "mapTokenHarborQuota"
Cohesion: 0.13
Nodes (16): clampPercent(), cloneConnectionQuota(), cloneOptionalString(), derefFloat(), deriveTokenHarborStatus(), mapTokenHarborQuota(), newQuotaCache(), normalizeOptionalString() (+8 more)

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.24
Nodes (10): filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY, OverviewCopy (+2 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.04
Nodes (92): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+84 more)

### Community 57 - "dependencies"
Cohesion: 0.07
Nodes (27): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+19 more)

### Community 58 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 59 - "server_test.go"
Cohesion: 0.07
Nodes (53): ModelOption, changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent (+45 more)

### Community 60 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 61 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.07
Nodes (59): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, @hookform/resolvers, useCombos(), useProviderConnections() (+51 more)

### Community 63 - "repository"
Cohesion: 0.10
Nodes (19): isTokenHarborHost(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+11 more)

### Community 64 - "repository"
Cohesion: 0.05
Nodes (41): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, ProviderStrategyState, repository, rowScanner (+33 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (34): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+26 more)

### Community 67 - "ModelDescriptor"
Cohesion: 0.24
Nodes (7): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, isModelAllowed(), isModelDisabled(), stubModelLister

### Community 68 - "navigation.ts"
Cohesion: 0.11
Nodes (27): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, RawModuleEntry, UILocale, ActiveNavContext (+19 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "settingDefinition"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 72 - "ModelEditorConsole.tsx"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 73 - "Module"
Cohesion: 0.12
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 74 - "JobInspector"
Cohesion: 0.15
Nodes (17): File Map, Review Focus, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan, JobInspector, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune() (+9 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.12
Nodes (24): buildDiscoveryCandidateURLs(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor, catalogOAuthDescriptor (+16 more)

### Community 76 - "arena/hooks.ts"
Cohesion: 0.08
Nodes (37): arenaQueryKeys, computeSessionTelemetry(), filterChatSessions(), formatLatencyMs(), formatTokenCount(), SessionTelemetry, sampleMessages, sampleSessions (+29 more)

### Community 77 - "Module"
Cohesion: 0.19
Nodes (7): compileRule(), decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), compiledRule, Module

### Community 78 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 79 - "WriteError"
Cohesion: 0.10
Nodes (12): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), decodeOptionalJSONBody(), actorFromContext(), decodeJSONBody(), Module (+4 more)

### Community 80 - "web/package.json"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.13
Nodes (13): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, Approaches Considered, Builder Decisions & Deliverables (+5 more)

### Community 82 - "Module"
Cohesion: 0.19
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 83 - "newTestHarness"
Cohesion: 0.11
Nodes (21): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newFakeClock() (+13 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 86 - "fakeHandle"
Cohesion: 0.13
Nodes (3): fakeHandle, osProcessHandle, ssePipeReader

### Community 87 - "Badge"
Cohesion: 0.10
Nodes (43): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+35 more)

### Community 88 - "settingDefinition"
Cohesion: 0.08
Nodes (15): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC(), truncateText(), decodeJSON(), formatBusAlert() (+7 more)

### Community 89 - "repo"
Cohesion: 0.13
Nodes (12): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow() (+4 more)

### Community 90 - "vite.config.ts"
Cohesion: 0.25
Nodes (7): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, omeg4monPwaPlugin()

### Community 93 - "redact.go"
Cohesion: 0.33
Nodes (7): antigravityPromptRewriteRule, isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, zod, zod@^4.6.7, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+36 more)

### Community 95 - "streamCache.ts"
Cohesion: 0.43
Nodes (4): UsageLogsResponse, matchesLogFilter(), prependUsageLogToResponse(), StreamLogsFilter

### Community 96 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "contracts.go"
Cohesion: 0.13
Nodes (18): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+10 more)

### Community 98 - "schema.go"
Cohesion: 0.33
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 99 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 101 - "Module"
Cohesion: 0.12
Nodes (7): Module, validateSQLiteSource(), decodeJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), RouteMounts

### Community 102 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 103 - "DB"
Cohesion: 0.09
Nodes (32): Options, Options, Options, New(), newRepo(), New(), newRepo(), decodeJSON() (+24 more)

### Community 105 - "ref_bun_test"
Cohesion: 0.08
Nodes (22): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+14 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.07
Nodes (37): ModelIdentityResult, extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse(), extractIDPrefix(), inferModelFamily(), modelsLooselyMatch() (+29 more)

### Community 108 - "WriteJSON"
Cohesion: 0.11
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON(), LogEntry (+3 more)

### Community 109 - "net/http.Header"
Cohesion: 0.18
Nodes (6): memoryResponseWriter, extractVendorHeaders(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.15
Nodes (17): @xyflow/react, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), RouterHubNode(), buildTopologyLayout() (+9 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "stubCompleter"
Cohesion: 0.26
Nodes (8): arenaTestEnv, stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "ProviderStore"
Cohesion: 0.05
Nodes (60): main(), runMain(), Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints (+52 more)

### Community 116 - "net/url.URL"
Cohesion: 0.08
Nodes (22): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+14 more)

### Community 117 - "AuditEvent"
Cohesion: 0.24
Nodes (7): buildWhere(), newEventID(), newRepo(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 118 - "deviceOAuthSession"
Cohesion: 0.23
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 119 - "Module"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 120 - "time.Time"
Cohesion: 0.08
Nodes (14): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, newLoginRateLimiter() (+6 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "time.Duration"
Cohesion: 0.13
Nodes (9): chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore(), acquireSem(), sleepContext() (+1 more)

### Community 125 - "New"
Cohesion: 0.06
Nodes (34): TestAppCloseTeardownOrderIsLIFO(), loadUI(), loadUI(), testEnv, Options, AuthVerifier, doPxpipeJSON(), setupTestPxpipe() (+26 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.02
Nodes (17): stubDep, migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile() (+9 more)

### Community 127 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 130 - "UsageRecord"
Cohesion: 0.11
Nodes (9): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher, ModelPricing (+1 more)

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "Module"
Cohesion: 0.15
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 140 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 142 - "newHarness"
Cohesion: 0.18
Nodes (11): Options, newMockVeniceServer(), TestVeniceKeylessNonStreamingAndStreamingEndToEnd(), startAuthHTTPProxy(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(), applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout() (+3 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "repository"
Cohesion: 0.29
Nodes (5): ImportRun, repository, New(), boolToInt(), newRepository()

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 154 - "plugins/interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

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

### Community 159 - "Card"
Cohesion: 0.12
Nodes (34): useProviderStrategy(), useUpdateProviderStrategyMutation(), ProviderStudioDetail(), getProviderStrategyMeta(), isProviderAccountStrategy(), ProviderStrategyCard(), ProviderStrategyCardProps, translatorQueryKeys (+26 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

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

### Community 177 - "WithProviderID"
Cohesion: 0.10
Nodes (18): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+10 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "promptinjection/engine.go"
Cohesion: 0.33
Nodes (11): actionRank(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity(), redactAndTruncate() (+3 more)

### Community 181 - "New"
Cohesion: 0.06
Nodes (39): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`) (+31 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (21): File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+13 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (15): stubAdminCredentialStore, TestContractInterfacesSatisfied(), extractUsageCost(), parseAwsQuickPastedQuery(), boolToInt(), newRepository(), Alias, AliasListResponse (+7 more)

### Community 194 - "listSettingsWithQueryer"
Cohesion: 0.22
Nodes (9): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer (+1 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "New"
Cohesion: 0.04
Nodes (61): roundTripFunc, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation() (+53 more)

### Community 198 - ".executeCandidateAttempt"
Cohesion: 0.13
Nodes (14): attemptOutcome, cachedGatewayResponse, inflightCall, responseCoalescer, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction() (+6 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "selectProbeHTTPClient"
Cohesion: 0.16
Nodes (12): selectProbeHTTPClient(), newBotAPIClient(), sanitizeBotError(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar() (+4 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.29
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "modeleditor/engine.go"
Cohesion: 0.35
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 209 - "GatewayInterceptor"
Cohesion: 0.07
Nodes (32): Global Constraints, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module (+24 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.18
Nodes (21): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES (+13 more)

### Community 222 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 223 - "Config"
Cohesion: 0.07
Nodes (42): authHarness, fakeCredentialStore, Options, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 2: `security` Module (Security Posture & Redacted Audit Trail), AdminCredentialStore (+34 more)

### Community 225 - "responsesStreamEncoder"
Cohesion: 0.27
Nodes (6): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder, buildCompactedResponse(), randomResponseID(), newResponsesStreamEncoder()

### Community 226 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

### Community 227 - "manifest.test.ts"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "headroomCache"
Cohesion: 0.24
Nodes (3): newHeadroomCache(), headroomCache, WarmPrefixSummary

### Community 231 - "time.Duration"
Cohesion: 0.07
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), computePrefixFingerprint()

### Community 242 - "compressor.go"
Cohesion: 0.13
Nodes (27): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+19 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 269 - "net/http.Request"
Cohesion: 0.06
Nodes (27): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), anthropicErrorType(), extractAPIKey() (+19 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1030 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `Module` to `ProvidersModals.tsx`?**
  _High betweenness centrality (0.199) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0990990990990991 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `UpstreamCredentials`, `go_pkg_encoding_json`, `net/http.Request`?**
  _High betweenness centrality (0.161) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `Config` to `newTestHarness`, `Button`, `buildApp`, `Module`, `ProviderStore`, `New`, `SettingsConsole.tsx`?**
  _High betweenness centrality (0.096) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._