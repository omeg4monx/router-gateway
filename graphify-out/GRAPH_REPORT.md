# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 517 files · ~566,231 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5893 nodes · 19788 edges · 263 communities (228 shown, 35 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1881 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ec184c35`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- combos/module.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- repository
- Open
- TypeScript Strict Compiler Config
- combos/schemas.ts
- App.tsx
- testing.T
- shadcn/ui Component Configuration
- ProvidersModals.tsx
- client.ts
- TopProgressBar.tsx
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- client.ts
- updates
- KeysConsole.tsx
- DB
- with
- jobs
- Module
- steps
- web-check
- on
- extractErrorMessage
- TokenSaverConsole.tsx
- ProviderDetailView.tsx
- SystemOverview.tsx
- docker-build
- go-check
- NotificationsConsole.tsx
- ModelEditorConsole.tsx
- Go Module Root Definition
- CanonicalUsage
- oauthTokenResponse
- QuotaConsole.tsx
- CanonicalStreamChunk
- main.go
- ir.go
- ProvidersConsole.tsx
- go_pkg_strings
- repository
- PluginsStudio.tsx
- Module
- ModelDescriptor
- setupModelEditorTest
- Module
- translate_responses.go
- New
- newHarness
- New
- Module
- web/package.json
- Module
- migrator.go
- Button
- oauthLoopbackProxy
- LiveMonitorConsole.tsx
- providers/repo.go
- newTestHarness
- teeHandler
- ProxyPoolsConsole.tsx
- TestContractInterfacesSatisfied
- TelegramConsole.tsx
- pricingStore
- oauthLoopbackProxy
- ResolvedTarget
- repo
- Module
- WriteJSON
- go_pkg_encoding_json
- arena/schemas.ts
- go_pkg_context
- telemetry.go
- readLegacySQLite
- repository
- PWA Support Design — `omeg4mon-router`
- Module
- .executeCandidateAttempt
- paths
- PromptGuardConsole.tsx
- AddApiKeyModal
- Module
- settingDefinition
- modelIntelligence.ts
- SystemOverview.tsx
- Module
- repository
- keyLimiter
- CanonicalChatResponse
- time.Duration
- modeleditor/engine.go
- extractErrorMessage
- net/http.ResponseWriter
- time.Time
- net/http.Request
- OverviewRealtimeTraffic.tsx
- evaluateAntigravityQuota
- post
- .Fetch
- Module
- New
- ProviderTopology.tsx
- Config
- repository
- repo
- properties
- Module
- net/url.URL
- Module
- RouteMounts
- .ListModels
- GatewayInterceptor
- plugins/repo.go
- properties
- fakeHandle
- translate_openai.go
- setupTokenSaverTest
- io/fs.FS
- strategies.go
- OMEG4MON Router
- DB
- zod
- certManager
- repo
- .identifyProviderModel
- Review Focus
- properties
- Module
- release-binaries
- steps
- applyUpstreamHeaders
- Provider Quota Menu & Token Harbor Quota Integration
- deviceOAuthSession
- ref_node_path
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- RouteMounts
- properties
- 200
- newTestHarness
- properties
- dependencies
- cli/package.json
- repository
- setupTestMitm
- get
- scripts
- post
- schemas
- manifest.test.ts
- PreviewPluginsResponse
- repository
- runtime_test.go
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- Review Focus
- Review Focus
- routingStateCache
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- promptinjection/engine.go
- setupMediaTest
- parameters
- responseCoalescer
- RingBuffer
- resolve.js
- ClearImporterHistoryResponse
- source_path
- .Snapshot
- scanNode
- docker-release
- palette.test.ts
- `@omeg4mon/router` CLI Launcher
- pricingStore
- .ExportCombos
- go_pkg_encoding_json
- context.Context
- setupTunnelTest
- redact.go
- parameters
- quotaCache
- .resolveIdentifier
- omeg4mon-router Verification Gate
- repo
- repo
- single-binary
- ResolvedTarget
- verify-web
- headroomCache
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- ChatCompleter
- responsesStreamEncoder
- Feature Suite
- stubDiscoverySystemModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- net/http/httptest.ResponseRecorder
- io.Reader
- net/http.Header
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- stubModule
- stubSystemModule
- .Reader
- TestVeniceLiveEndToEnd
- TestJetBrainsQuotaFetcher
- AttestationManager
- selectProbeHTTPClient
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- vite.config.ts
- time.Duration
- fakeRunner
- optimization.test.ts
- newTestHarness
- oauthProviderStrategy
- jsonMarshalerError
- buildModelProbeRequest
- compressor.go
- .handleOAuthPaste
- translate_responses.go
- BuildChatPayload
- Global Constraints
- App.test.ts
- stubSettings
- setupNotificationsTest
- repository
- ModelTestResult
- File Map
- telemetryCollector
- omnirouteFlow.test.ts
- .InsertRun
- ProcessHandle
- parseIdentityPayload

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 194 edges
2. `WriteJSON()` - 187 edges
3. `WriteError()` - 184 edges
4. `Button` - 106 edges
5. `cn()` - 97 edges
6. `DB` - 93 edges
7. `Module` - 87 edges
8. `buildApp()` - 79 edges
9. `Badge()` - 76 edges
10. `Bus` - 75 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Approaches Considered` --references--> `quotaFetcher`  [INFERRED]
  docs/superpowers/specs/2026-11-27-provider-quota-design.md → internal/modules/providers/quota.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx

## Import Cycles
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/optimization_e2e_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (263 total, 35 thin omitted)

### Community 0 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
Nodes (6): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC(), testHarness

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
Cohesion: 0.11
Nodes (18): cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary(), createParams (+10 more)

### Community 6 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 7 - "Open"
Cohesion: 0.05
Nodes (53): roundTripFunc, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation() (+45 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "combos/schemas.ts"
Cohesion: 0.05
Nodes (73): ComboCard, ComboCardProps, buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup (+65 more)

### Community 10 - "App.tsx"
Cohesion: 0.03
Nodes (129): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, class-variance-authority, lucide-react, react (+121 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (98): TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestComboTestHonorsPinnedConnectionAndComboRef(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), TestResolveAttachesComboRuntimeMetadata() (+90 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (57): MODAL_SIZE_CLASS, ModalSize, CatalogOAuthDescriptor, commaListSchema(), connectionToEditFormValues(), CreateConnectionFormInput, CreateConnectionFormOutput, createConnectionSchema (+49 more)

### Community 14 - "client.ts"
Cohesion: 0.17
Nodes (22): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, @tanstack/react-query, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT (+14 more)

### Community 15 - "TopProgressBar.tsx"
Cohesion: 0.16
Nodes (15): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety, computeNextTrickleProgress(), createTopProgressController(), getTopProgressBarStyle() (+7 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.19
Nodes (20): resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), decodeJWTPayload(), enrichAntigravity(), enrichGeminiOAuth(), enrichGitHubCopilot(), extractCodexAccountInfo() (+12 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "client.ts"
Cohesion: 0.09
Nodes (42): @hookform/resolvers, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation() (+34 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (79): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, react-hook-form, buildGatewayUrls(), copyTextToClipboard(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet() (+71 more)

### Community 21 - "DB"
Cohesion: 0.06
Nodes (39): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog(), TestWave2BDisabledModulesSkipInterceptors() (+31 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "extractErrorMessage"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.12
Nodes (30): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId (+22 more)

### Community 30 - "ProviderDetailView.tsx"
Cohesion: 0.05
Nodes (74): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, ModelIdentityResult (+66 more)

### Community 31 - "SystemOverview.tsx"
Cohesion: 0.11
Nodes (36): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), buildLogsQueryParams(), INVALIDATION_TOPIC_QUERY_KEYS, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, computeBarPercent(), computeSharePercent() (+28 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 35 - "ModelEditorConsole.tsx"
Cohesion: 0.07
Nodes (59): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation() (+51 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+8 more)

### Community 39 - "oauthTokenResponse"
Cohesion: 0.17
Nodes (19): refreshCodeBuddyToken(), TestCodeBuddyTokenRefresh(), buildJetBrainsAuthURL(), cleanJetBrainsAuthCode(), exchangeJetBrainsCode(), normalizeJetBrainsRedirectURI(), postExchangeJetBrains(), provisionJetBrainsGrazieJWT() (+11 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.10
Nodes (46): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+38 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.11
Nodes (19): anthropicStreamDecoder, CanonicalStreamChunk, geminiStreamDecoder, openAIStreamDecoder, streamDecoder, NewStreamTranslator(), mergeFastPathSSEUsage(), newAnthropicStreamDecoder() (+11 more)

### Community 42 - "main.go"
Cohesion: 0.06
Nodes (22): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+14 more)

### Community 43 - "ir.go"
Cohesion: 0.13
Nodes (36): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, anthropicContentBlock, anthropicImageSource, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource (+28 more)

### Community 44 - "ProvidersConsole.tsx"
Cohesion: 0.07
Nodes (52): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), buildProviderSectionBuckets(), canAddProviderConnection(), canProbeProviderModels(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets() (+44 more)

### Community 45 - "go_pkg_strings"
Cohesion: 0.09
Nodes (31): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+23 more)

### Community 46 - "repository"
Cohesion: 0.12
Nodes (12): escapeLikePattern(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, pricingListResponse, ProviderBreakdown, repository (+4 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.16
Nodes (30): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+22 more)

### Community 48 - "Module"
Cohesion: 0.25
Nodes (4): decodeJSONBody(), newID(), validateOverride(), Module

### Community 49 - "ModelDescriptor"
Cohesion: 0.31
Nodes (6): buildServiceWorkerSource(), CORE_PRECACHE_URLS, classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "setupModelEditorTest"
Cohesion: 0.15
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestConcurrentExtraParamsNoDataRaceOrCacheMutation() (+10 more)

### Community 51 - "Module"
Cohesion: 0.13
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO, Notification

### Community 52 - "translate_responses.go"
Cohesion: 0.09
Nodes (28): anthropicChoice, anthropicMessage, anthropicRequest, anthropicThinking, anthropicTool, geminiCandidate, geminiContent, geminiFileData (+20 more)

### Community 53 - "New"
Cohesion: 0.09
Nodes (25): TestAppCloseTeardownOrderIsLIFO(), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, New(), newHarness(), TestTranslatorDeleteLogsClearsRingAndBroadcastsSSEClear(), TestTranslatorLogsSnapshotAndLiveSSEStream(), TestTranslatorLogsStreamSurvivesPastAPITimeout() (+17 more)

### Community 54 - "newHarness"
Cohesion: 0.10
Nodes (22): TestWave2AE2EProviderProxyRouting(), File Map, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, newMockVeniceServer(), newHarness() (+14 more)

### Community 55 - "New"
Cohesion: 0.10
Nodes (21): roundTripFunc, newAntigravityQuotaMonitor(), newChainStore(), TestAntigravityQuotaMonitorLeaderCancelDoesNotAbortWaiters(), TestChainStoreThrottlesExpiredPrune(), New(), TestGatewayGlobalConcurrencyLoadShedding(), TestTrafficShaperAdaptiveHeaderPacing() (+13 more)

### Community 56 - "Module"
Cohesion: 0.09
Nodes (13): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), newKeyCache(), extractUsageCost(), hashToken(), isValidTokenFormat() (+5 more)

### Community 57 - "web/package.json"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 58 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 59 - "migrator.go"
Cohesion: 0.07
Nodes (48): ModelOption, handledGatewayError, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, buildInstallGuide() (+40 more)

### Community 60 - "Button"
Cohesion: 0.07
Nodes (101): ArenaConsole(), useArenaBattles(), useArenaLeaderboard(), useDeleteArenaBattleMutation(), useRunArenaBattleMutation(), useVoteArenaBattleMutation(), AuthGuard(), LoginForm() (+93 more)

### Community 61 - "oauthLoopbackProxy"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 62 - "LiveMonitorConsole.tsx"
Cohesion: 0.12
Nodes (42): attrValueClass(), levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, normalizeUsagePayload(), StreamState, SYSTEM_LEVELS (+34 more)

### Community 63 - "providers/repo.go"
Cohesion: 0.16
Nodes (15): cloneStringSlice(), cloneStringMap(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository(), scanConnectionRecord() (+7 more)

### Community 64 - "newTestHarness"
Cohesion: 0.13
Nodes (19): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+11 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.12
Nodes (33): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation(), useUpsertProxyBindingMutation() (+25 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.12
Nodes (7): stubAdminAuth, AuthVerifier, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 68 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "oauthLoopbackProxy"
Cohesion: 0.19
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 71 - "ResolvedTarget"
Cohesion: 0.17
Nodes (11): ProviderStrategyState, strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), decodeConfigJSON(), fastRandFloat64() (+3 more)

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "WriteJSON"
Cohesion: 0.09
Nodes (10): decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), WriteJSON(), LogEntry, LogListResponse (+2 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (27): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+19 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (46): ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaModels(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteChatSessionMutation() (+38 more)

### Community 78 - "telemetry.go"
Cohesion: 0.16
Nodes (21): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+13 more)

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.13
Nodes (13): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, Approaches Considered, Builder Decisions & Deliverables (+5 more)

### Community 82 - "Module"
Cohesion: 0.24
Nodes (7): actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 83 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (21): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), extractUpstreamErrorMessage(), parseRetryAfter(), applyCodeBuddyHeaders() (+13 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 86 - "AddApiKeyModal"
Cohesion: 0.13
Nodes (28): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation() (+20 more)

### Community 87 - "Module"
Cohesion: 0.03
Nodes (63): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), UpstreamCredentials, allBuiltInAdapters() (+55 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 90 - "SystemOverview.tsx"
Cohesion: 0.12
Nodes (24): Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, useCombos(), useApiKeys(), useProviderConnections(), useProxyPools() (+16 more)

### Community 91 - "Module"
Cohesion: 0.13
Nodes (8): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, generateID(), newBatcher(), recordingUsage, batcher

### Community 92 - "repository"
Cohesion: 0.26
Nodes (3): DisabledModel, ModelAlias, repository

### Community 93 - "keyLimiter"
Cohesion: 0.31
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 95 - "CanonicalChatResponse"
Cohesion: 0.12
Nodes (23): CanonicalChatResponse, ProtocolTranslator, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol() (+15 more)

### Community 96 - "time.Duration"
Cohesion: 0.08
Nodes (10): Claims, failureRecord, stubHealthTracker, trafficShaper, Module, computeBackoff(), acquireSem(), sleepContext() (+2 more)

### Community 97 - "modeleditor/engine.go"
Cohesion: 0.20
Nodes (18): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+10 more)

### Community 98 - "extractErrorMessage"
Cohesion: 0.05
Nodes (75): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+67 more)

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (24): Module, Module, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), trackedResponseWriter, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID() (+16 more)

### Community 101 - "time.Time"
Cohesion: 0.10
Nodes (13): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, stubAccountHealthTracker (+5 more)

### Community 102 - "net/http.Request"
Cohesion: 0.07
Nodes (24): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), isBlockedAntigravityProjectID(), anthropicErrorType() (+16 more)

### Community 103 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.22
Nodes (15): formatActivityOriginLabel(), formatClockUptime(), LazyProviderTopology, OverviewRealtimeTraffic, OverviewRealtimeTrafficProps, RecentRequestsRows(), computeRelativeClockIntervalMs(), FilteredRecentRequestsResult (+7 more)

### Community 105 - "evaluateAntigravityQuota"
Cohesion: 0.20
Nodes (13): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+5 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".Fetch"
Cohesion: 0.21
Nodes (18): TestCodeBuddyQuotaExcludesExpiredBonusAndZeroCapacityIsExhausted(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+10 more)

### Community 108 - "Module"
Cohesion: 0.15
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 109 - "New"
Cohesion: 0.08
Nodes (28): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+20 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (21): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, EnergyTopologyEdge(), FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, prefersReducedMotion() (+13 more)

### Community 111 - "Config"
Cohesion: 0.07
Nodes (30): Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), setupSecurityTest(), TestSecurityAuditPaginationOffset(), TestSecurityEventBusAutoAuditIngestion(), TestSecurityPostureAndExportAndPrune(), TestSecurityRedactionAndTamperEvidentClear(), TestSecurityWave2BTopicParityAndExtendedRedaction() (+22 more)

### Community 112 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 116 - "net/url.URL"
Cohesion: 0.10
Nodes (19): stubProxyDialerResolver, Options, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+11 more)

### Community 117 - "Module"
Cohesion: 0.21
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 118 - "RouteMounts"
Cohesion: 0.29
Nodes (3): Module, validateSQLiteSource(), decodeJSON()

### Community 119 - ".ListModels"
Cohesion: 0.12
Nodes (9): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, testEnv, isModelAllowed(), isModelDisabled(), TestListVirtualModels() (+1 more)

### Community 120 - "GatewayInterceptor"
Cohesion: 0.14
Nodes (16): Global Constraints, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Global Constraints, Review Focus (+8 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.06
Nodes (59): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+51 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "fakeHandle"
Cohesion: 0.13
Nodes (3): fakeHandle, osProcessHandle, ssePipeReader

### Community 124 - "translate_openai.go"
Cohesion: 0.19
Nodes (19): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+11 more)

### Community 125 - "setupTokenSaverTest"
Cohesion: 0.23
Nodes (13): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestResponsesAPIFunctionCallOutputRTKFilter(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory() (+5 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.02
Nodes (22): stubDep, stubModule, migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), TestStaticUIHandlerPWAHeaders(), applyStaticCacheHeaders() (+14 more)

### Community 127 - "strategies.go"
Cohesion: 0.16
Nodes (15): Combo, ComboStep, createComboParams, rowScanner, updateComboValues, decodeStringSlice(), encodeStringSlice(), isUniqueConstraintError() (+7 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "DB"
Cohesion: 0.06
Nodes (51): Options, Options, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Global Constraints, Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Wave 3B — Media Gateway (Audio/Images/Videos/Search/Fetch) & Tunnel Manager Implementation Plan, Module (+43 more)

### Community 130 - "zod"
Cohesion: 0.15
Nodes (16): zod, translatorQueryKeys, useTranslatorProtocols(), ConsoleLogEntry, parseJsonObject(), SAMPLE_PAYLOADS, translateFormSchema, TranslateFormValues (+8 more)

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - ".identifyProviderModel"
Cohesion: 0.20
Nodes (13): ModelIdentityResult, extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse(), extractVendorHeaders(), inferModelFamily(), modelsLooselyMatch() (+5 more)

### Community 134 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "Module"
Cohesion: 0.20
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "applyUpstreamHeaders"
Cohesion: 0.12
Nodes (14): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+6 more)

### Community 140 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration

### Community 141 - "deviceOAuthSession"
Cohesion: 0.18
Nodes (6): newDeviceOAuthRuntime(), originOfURL(), deviceOAuthResult, deviceOAuthRuntime, deviceOAuthSession, rfc8628DeviceSpec

### Community 142 - "ref_node_path"
Cohesion: 0.19
Nodes (7): consoleSource, detailSource, modalsSource, providersDir, modalsSource, providersDir, appShellSource

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "RouteMounts"
Cohesion: 0.09
Nodes (7): actorFromContext(), decodeJSONBody(), Module, decodeJSONBody(), AdminSubject(), RouteMounts, Module

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "newTestHarness"
Cohesion: 0.08
Nodes (29): File Map, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), TestValidateKeyQueuesLastUsedAsynchronouslyAndFlushesOnStop(), New() (+21 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.07
Nodes (27): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+19 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 153 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 154 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

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

### Community 159 - "manifest.test.ts"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 162 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 163 - "runtime_test.go"
Cohesion: 0.33
Nodes (11): New(), setupPluginsTest(), TestPluginsBuiltinsUpdateAndInterceptorScoping(), newRepository(), putPluginConfig(), TestPluginsAntiSlopFormatAwareCachePreservingAndIdempotent(), TestPluginsContextSqueezerTokenAwareEvictionAndElision(), TestPluginsImageVisionOnlyInjectsWhenImagesPresent() (+3 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "routingStateCache"
Cohesion: 0.25
Nodes (3): cachedRoutingState, routingStateCache, newRoutingStateCache()

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

### Community 177 - "promptinjection/engine.go"
Cohesion: 0.30
Nodes (13): actionRank(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity() (+5 more)

### Community 178 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "responseCoalescer"
Cohesion: 0.24
Nodes (7): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerLeaderCancelDoesNotAbortWaiters(), TestResponseCoalescerPanicDoesNotDeadlockWaiters(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 181 - "RingBuffer"
Cohesion: 0.27
Nodes (6): Global Constraints, Review Focus, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, cloneAttrMap(), RingBuffer, LogEntry

### Community 182 - "resolve.js"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - ".Snapshot"
Cohesion: 0.07
Nodes (11): arenaTestEnv, stubCompleter, JobInspector, stubKeyProvisioner, Runner, job, JobStatus, stubJobInspector (+3 more)

### Community 186 - "scanNode"
Cohesion: 0.29
Nodes (5): decodeStringMap(), decodeStringSlice(), scanNode(), ProviderNode, rowScanner

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.23
Nodes (13): anthropicResponse, anthropicUsage, errorDisposition, classifyStatus(), isClientPayloadErrorBody(), cloneShallowMap(), extractCodeBuddyText(), sanitizeCodeBuddyChannelMarkers() (+5 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (20): authHarness, fakeCredentialStore, catalogFakeProviderStore, fakeProviderStore, staticProviderStore, AdminCredentialStore, stubAdminCredentialStore, stubProviderStore (+12 more)

### Community 194 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 195 - "redact.go"
Cohesion: 0.06
Nodes (30): buildWhere(), scanNotification(), isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere() (+22 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "quotaCache"
Cohesion: 0.29
Nodes (7): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry, quotaFlightCall

### Community 198 - ".resolveIdentifier"
Cohesion: 0.15
Nodes (11): resolveSnapshot, applyComboRuntimeMetadata(), buildResolvedTarget(), configIntValue(), connectionSupportsBare(), connectionSupportsExplicit(), isNestedExecuteMode(), withMetadataEntry() (+3 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 201 - "repo"
Cohesion: 0.18
Nodes (7): mapBusEventToNotification(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotificationRow(), EmitPayload, repo

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "ResolvedTarget"
Cohesion: 0.33
Nodes (5): stubProviderResolver, stubResolver, ResolvedTarget, prioritizeHeadroomAffinity(), TestPrioritizeHeadroomAffinityPromotesMatchingCandidateWithoutLosingFallbacks()

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "headroomCache"
Cohesion: 0.24
Nodes (3): newHeadroomCache(), headroomCache, WarmPrefixSummary

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "ChatCompleter"
Cohesion: 0.25
Nodes (9): ChatCompletionRequest, ChatCompletionResult, ChatMessage, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan, ChatCompleter (+1 more)

### Community 209 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & System, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.19
Nodes (18): Options, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+10 more)

### Community 219 - "net/http.Header"
Cohesion: 0.09
Nodes (12): memoryResponseWriter, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), copyNonHopHeaders() (+4 more)

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 223 - ".Reader"
Cohesion: 0.11
Nodes (16): main(), runMain(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan (+8 more)

### Community 224 - "TestVeniceLiveEndToEnd"
Cohesion: 0.22
Nodes (8): peekFreePIPinMismatch(), TestPeekFreePIPinMismatch(), ApplyHeaders(), NewNDJSONToOpenAISSEReader(), ParseNDJSONToOpenAI(), TestNewNDJSONToOpenAISSEReader(), TestParseNDJSONToOpenAI(), TestVeniceLiveEndToEnd()

### Community 225 - "TestJetBrainsQuotaFetcher"
Cohesion: 0.22
Nodes (5): TestJetBrainsQuotaFetcher(), makeTestJWT(), TestExtractCodexAccountInfo(), isCodeBuddyHostFor(), isTokenHarborHost()

### Community 226 - "AttestationManager"
Cohesion: 0.31
Nodes (5): AttestationManager, isTransientMintStatus(), ResolveMintURL(), inFlightMint, tokenBucket

### Community 227 - "selectProbeHTTPClient"
Cohesion: 0.25
Nodes (7): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), IsVeniceKeyless(), TestIsVeniceKeyless()

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "vite.config.ts"
Cohesion: 0.25
Nodes (7): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, omeg4monPwaPlugin()

### Community 231 - "time.Duration"
Cohesion: 0.05
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID() (+1 more)

### Community 235 - "optimization.test.ts"
Cohesion: 0.25
Nodes (7): consoleSource, customModelModalSource, detailSource, hooksSource, logoSource, modalsSource, providersDir

### Community 236 - "newTestHarness"
Cohesion: 0.14
Nodes (12): testHarness, decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+4 more)

### Community 238 - "oauthProviderStrategy"
Cohesion: 0.60
Nodes (5): newCodeBuddyOAuthStrategy(), newJetBrainsOAuthStrategy(), defaultOAuthStrategies(), newRFC8628Strategy(), oauthProviderStrategy

### Community 241 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 242 - "compressor.go"
Cohesion: 0.10
Nodes (36): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressRoot(), compressText(), estimateTokens(), hasAnthropicToolResult() (+28 more)

### Community 243 - ".handleOAuthPaste"
Cohesion: 0.33
Nodes (3): parseAwsQuickDiscovery(), parseAwsQuickPastedQuery(), renderAwsQuickResultHTML()

### Community 246 - "translate_responses.go"
Cohesion: 0.18
Nodes (19): responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse, responsesStreamDecoder (+11 more)

### Community 247 - "BuildChatPayload"
Cohesion: 0.29
Nodes (7): BuildChatPayload(), extractContentText(), extractToolArguments(), formatToolsInstruction(), randomUUIDv4(), StripModelPrefix(), TestBuildChatPayload()

### Community 249 - "Global Constraints"
Cohesion: 0.33
Nodes (5): Global Constraints, Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), Wave 3C — Local MITM Proxy & PXPIPE Transform Pipeline Implementation Plan

### Community 250 - "App.test.ts"
Cohesion: 0.50
Nodes (3): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR

### Community 252 - "setupNotificationsTest"
Cohesion: 0.47
Nodes (5): Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict()

### Community 253 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 256 - "ModelTestResult"
Cohesion: 0.40
Nodes (3): ModelTestResult, testAllModelsResponse, resolveKeylessProbeCredentials()

### Community 257 - "File Map"
Cohesion: 0.40
Nodes (4): File Map, Global Constraints, Review Focus, Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan

### Community 258 - "telemetryCollector"
Cohesion: 0.70
Nodes (3): TestTelemetryCacheAvoidsRedundantWork(), staticHostInfo, telemetryCollector

### Community 259 - "omnirouteFlow.test.ts"
Cohesion: 0.40
Nodes (4): consoleSource, detailSource, modalsSource, providersDir

### Community 262 - "parseIdentityPayload"
Cohesion: 1.00
Nodes (3): extractIDPrefix(), parseIdentityPayload(), parseSSEIdentityPayload()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1077 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `AddApiKeyModal` to `Button`, `Module`?**
  _High betweenness centrality (0.200) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `combos/module.go` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `net/http.Request`, `testing.T`, `go_pkg_encoding_json`, `AddApiKeyModal`?**
  _High betweenness centrality (0.173) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `newTestHarness` to `context.Context`, `DB`, `App.tsx`, `New`, `KeysConsole.tsx`, `net/http/httptest.ResponseRecorder`, `.Reader`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._