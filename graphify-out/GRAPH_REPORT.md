# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 517 files · ~566,099 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5892 nodes · 19782 edges · 256 communities (217 shown, 39 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1877 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `84b89c97`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- combos/module.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- modeleditor/engine.go
- Open
- TypeScript Strict Compiler Config
- combos/schemas.ts
- useUIStore
- testing.T
- shadcn/ui Component Configuration
- TestContractInterfacesSatisfied
- client.ts
- App.tsx
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
- ProviderAdapter
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
- setupModelEditorTest
- Module
- translate_responses.go
- setupNotificationsTest
- NewClient
- Open
- ref_shared_api_schema
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
- navigation.ts
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
- media/schemas.ts
- WriteError
- settingDefinition
- modelIntelligence.ts
- .handleFetchConnectionModels
- Module
- repository
- keyLimiter
- Module
- CanonicalChatResponse
- time.Duration
- UpstreamCredentials
- extractErrorMessage
- net/http.Request
- time.Time
- .executeInferenceCandidates
- Module
- evaluateAntigravityQuota
- post
- .Fetch
- settingDefinition
- New
- ProviderTopology.tsx
- Load
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
- scanCombo
- certManager
- repo
- .identifyProviderModel
- Review Focus
- properties
- Module
- release-binaries
- steps
- applyUpstreamHeaders
- Module
- deviceOAuthSession
- ref_bun_test
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- Module
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
- net/http.Handler
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
- stubDep
- setupMediaTest
- parameters
- ProviderStore
- RingBuffer
- resolve.js
- ClearImporterHistoryResponse
- source_path
- .Snapshot
- New
- docker-release
- palette.test.ts
- `@omeg4mon/router` CLI Launcher
- SessionSigner
- .ExportCombos
- go_pkg_encoding_json
- context.Context
- setupTunnelTest
- listSettingsWithQueryer
- parameters
- quotaCache
- .resolveIdentifier
- omeg4mon-router Verification Gate
- repo
- repo
- single-binary
- stubModule
- verify-web
- headroomCache
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- ChatCompleter
- stubSystemModule
- Feature Suite
- stubDiscoverySystemModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- net/http/httptest.ResponseRecorder
- stubSystem
- net/http.Header
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- stubModule
- stubSystemModule
- .Reader
- stubDep
- .streamUpstreamSSE
- AttestationManager
- stubDepModule
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- stubModule
- time.Duration
- stubSystemModule
- omeg4mon-router Project Standards
- newTestHarness
- oauthProviderStrategy
- jsonMarshalerError
- buildModelProbeRequest
- compressor.go
- trackedResponseWriter
- decodeOpenAIResponsesResponse
- stubModule
- Global Constraints
- App.test.ts
- mobileLayout.test.ts
- mapBusEventToAudit
- HistoryItem

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
- `Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD)` --references--> `defaultBuiltInAdapters()`  [INFERRED]
  docs/superpowers/plans/2026-11-27-provider-quota.md → internal/modules/providers/adapter.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx

## Import Cycles
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/freepi_test.go -> internal/modules/gateway/jetbrains_test.go -> internal/modules/gateway/freepi_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/optimization_e2e_test.go`
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

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (256 total, 39 thin omitted)

### Community 0 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

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
Cohesion: 0.06
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+19 more)

### Community 6 - "modeleditor/engine.go"
Cohesion: 0.06
Nodes (39): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+31 more)

### Community 7 - "Open"
Cohesion: 0.06
Nodes (45): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), roundTripFunc, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), TestAntigravityQuotaMonitorLeaderCancelDoesNotAbortWaiters(), TestFetchConnectionModelsEndpoint() (+37 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "combos/schemas.ts"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 10 - "useUIStore"
Cohesion: 0.07
Nodes (59): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope (+51 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (96): TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestComboTestHonorsPinnedConnectionAndComboRef(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), TestResolveAttachesComboRuntimeMetadata() (+88 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TestContractInterfacesSatisfied"
Cohesion: 0.10
Nodes (7): catalogFakeProviderStore, fakeProviderStore, stubProviderStore, stubProviderStore, TestContractInterfacesSatisfied(), newFakeProviderStore(), stubProviderStore

### Community 14 - "client.ts"
Cohesion: 0.17
Nodes (16): securityQueryKeys, useClearSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState, ClearSecurityAuditResponse (+8 more)

### Community 15 - "App.tsx"
Cohesion: 0.06
Nodes (53): App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio, InitPage (+45 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.12
Nodes (34): cleanJetBrainsAuthCode(), exchangeJetBrainsCode(), postExchangeJetBrains(), refreshJetBrainsToken(), TestJetBrainsPostExchangeAndRefresh(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL() (+26 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "client.ts"
Cohesion: 0.09
Nodes (43): zod, zod@^4.6.7, authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+35 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (81): @hookform/resolvers, buildGatewayUrls(), copyTextToClipboard(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState() (+73 more)

### Community 21 - "DB"
Cohesion: 0.04
Nodes (45): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+37 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "Module"
Cohesion: 0.21
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
Cohesion: 0.02
Nodes (201): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, @tanstack/react-query, runModelIdentityProbe() (+193 more)

### Community 31 - "SystemOverview.tsx"
Cohesion: 0.06
Nodes (66): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), useSystemTelemetry(), formatBytesGB(), formatHostUptimeShort(), HostTelemetryCard, HostTelemetryCardProps, formatActivityOriginLabel(), formatClockUptime() (+58 more)

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
Cohesion: 0.08
Nodes (56): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation() (+48 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (14): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamDecoder (+6 more)

### Community 39 - "ProviderAdapter"
Cohesion: 0.08
Nodes (14): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), cloneStringSlice(), defaultBuiltInAdapters(), oauthEnvPrefix(), xrouterBuiltInOAuthAdapters(), buildJetBrainsAuthURL() (+6 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (52): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+44 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.13
Nodes (16): anthropicStreamDecoder, CanonicalStreamChunk, geminiStreamDecoder, openAIStreamDecoder, responsesStreamDecoder, NewStreamTranslator(), newAnthropicStreamDecoder(), newAnthropicStreamEncoder() (+8 more)

### Community 42 - "main.go"
Cohesion: 0.07
Nodes (21): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+13 more)

### Community 43 - "ir.go"
Cohesion: 0.13
Nodes (37): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, anthropicContentBlock, anthropicImageSource, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource (+29 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.13
Nodes (15): geminiCandidate, geminiContent, geminiFunctionCallingConfig, geminiFunctionDecl, geminiGenConfig, geminiRequest, geminiResponse, geminiThinkingCfg (+7 more)

### Community 45 - "go_pkg_strings"
Cohesion: 0.10
Nodes (22): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+14 more)

### Community 46 - "repository"
Cohesion: 0.12
Nodes (12): escapeLikePattern(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, pricingListResponse, ProviderBreakdown, repository (+4 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.10
Nodes (45): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, modelCatalogQueryKeys (+37 more)

### Community 48 - "Module"
Cohesion: 0.22
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 49 - "ModelDescriptor"
Cohesion: 0.26
Nodes (8): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin(), classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "setupModelEditorTest"
Cohesion: 0.14
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, New(), doJSON(), setupModelEditorTest() (+11 more)

### Community 51 - "Module"
Cohesion: 0.07
Nodes (11): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+3 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.11
Nodes (25): anthropicChoice, anthropicMessage, anthropicRequest, anthropicThinking, anthropicTool, geminiFileData, geminiFunctionCall, geminiFunctionResponse (+17 more)

### Community 53 - "setupNotificationsTest"
Cohesion: 0.12
Nodes (19): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup() (+11 more)

### Community 54 - "NewClient"
Cohesion: 0.12
Nodes (17): Options, newMockVeniceServer(), startAuthHTTPProxy(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(), applyDefaults(), NewClient(), NewProviderProxyFunc(), ProviderIDFromContext() (+9 more)

### Community 55 - "Open"
Cohesion: 0.07
Nodes (31): TestAppWave1BGatewayEndToEnd(), roundTripFunc, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestCombosRoutingStateCacheHitsAndInvalidates(), TestKeylessFreeProviderModelProbeAndIdentify() (+23 more)

### Community 56 - "ref_shared_api_schema"
Cohesion: 0.11
Nodes (24): StrategyCategory, ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation() (+16 more)

### Community 57 - "web/package.json"
Cohesion: 0.10
Nodes (27): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+19 more)

### Community 58 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 59 - "migrator.go"
Cohesion: 0.09
Nodes (41): ModelOption, PolicyChangedEvent, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, buildInstallGuide() (+33 more)

### Community 60 - "Button"
Cohesion: 0.06
Nodes (122): class-variance-authority, lucide-react, react, ArenaConsole(), useArenaBattles(), useArenaLeaderboard(), useDeleteArenaBattleMutation(), useRunArenaBattleMutation() (+114 more)

### Community 61 - "oauthLoopbackProxy"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 62 - "LiveMonitorConsole.tsx"
Cohesion: 0.07
Nodes (57): attrValueClass(), levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, normalizeUsagePayload(), StreamState, SYSTEM_LEVELS (+49 more)

### Community 63 - "providers/repo.go"
Cohesion: 0.11
Nodes (21): cloneStringSlice(), cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+13 more)

### Community 64 - "newTestHarness"
Cohesion: 0.15
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (35): react-hook-form, proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+27 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.09
Nodes (9): stubAdminAuth, AuthVerifier, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth (+1 more)

### Community 68 - "navigation.ts"
Cohesion: 0.14
Nodes (20): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, ActiveNavContext, countVisibleNavItems(), EXTENSION_GROUPS, findActiveNavContext(), getShellCopy(), isModuleRegistered() (+12 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "oauthLoopbackProxy"
Cohesion: 0.21
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 71 - "ResolvedTarget"
Cohesion: 0.14
Nodes (13): strategyState, weightedCandidate, stubProviderResolver, stubResolver, ResolvedTarget, applyComboStrategy(), applyProviderStrategy(), applyStrategy() (+5 more)

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "WriteJSON"
Cohesion: 0.10
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON(), LogEntry (+3 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.17
Nodes (16): catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, ChangedEvent, connectionQuotaDTO, createConnectionRequest, createNodeRequest, listNodesResponse (+8 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (46): ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaModels(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteChatSessionMutation() (+38 more)

### Community 78 - "telemetry.go"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "repository"
Cohesion: 0.11
Nodes (21): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+13 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.16
Nodes (10): extractUsageCost(), actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+2 more)

### Community 83 - ".executeCandidateAttempt"
Cohesion: 0.09
Nodes (22): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), applyCodeBuddyHeaders(), codebuddyDomainFor(), isCodeBuddyProvider() (+14 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 86 - "media/schemas.ts"
Cohesion: 0.13
Nodes (20): mediaQueryKeys, CreateMediaComboFormInput, CreateMediaComboFormOutput, CreateMediaComboRequest, createMediaComboSchema, CreateMediaProviderFormInput, CreateMediaProviderFormOutput, CreateMediaProviderRequest (+12 more)

### Community 87 - "WriteError"
Cohesion: 0.10
Nodes (25): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, isValidAuthType(), isValidProtocol(), validateHTTPURL(), buildAwsQuickSigninURL(), generateAwsQuickNonce() (+17 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 90 - ".handleFetchConnectionModels"
Cohesion: 0.12
Nodes (15): generateTraceID(), resolveAwsQuickTenantBaseURL(), sanitizeUpstreamError(), truncateError(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels() (+7 more)

### Community 91 - "Module"
Cohesion: 0.10
Nodes (10): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, generateID(), newBatcher(), recordingUsage, batcher (+2 more)

### Community 92 - "repository"
Cohesion: 0.26
Nodes (3): DisabledModel, ModelAlias, repository

### Community 93 - "keyLimiter"
Cohesion: 0.27
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 95 - "CanonicalChatResponse"
Cohesion: 0.14
Nodes (21): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), RewriteModelField() (+13 more)

### Community 96 - "time.Duration"
Cohesion: 0.13
Nodes (9): failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), acquireSem(), sleepContext(), formatDuration(), bunBoringSSLClientHelloSpec() (+1 more)

### Community 97 - "UpstreamCredentials"
Cohesion: 0.09
Nodes (16): staticProviderStore, UpstreamCredentials, applyJetBrainsHeaders(), buildProtocolTestRequest(), cloneUpstreamCredentials(), isBlockedAntigravityProjectID(), cloneProviderConnections(), TestCodeBuddyAndFreePIProtocolTestRequests() (+8 more)

### Community 98 - "extractErrorMessage"
Cohesion: 0.05
Nodes (79): openapi-fetch, openapi-fetch@^0.14.0, useProviderNodes(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs() (+71 more)

### Community 99 - "net/http.Request"
Cohesion: 0.05
Nodes (35): Module, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), Module, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody() (+27 more)

### Community 101 - "time.Time"
Cohesion: 0.09
Nodes (15): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, stubAccountHealthTracker (+7 more)

### Community 102 - ".executeInferenceCandidates"
Cohesion: 0.16
Nodes (9): cachedGatewayResponse, inflightCall, responseCoalescer, ensureRequestID(), parseComboMetadataInt(), parseInboundMeta(), computeJitteredRetryDelay(), isTransientSingleCandidateStatus() (+1 more)

### Community 105 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".Fetch"
Cohesion: 0.10
Nodes (26): TestCodeBuddyQuotaExcludesExpiredBonusAndZeroCapacityIsExhausted(), TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), provisionJetBrainsGrazieJWT(), TestJetBrainsQuotaFetcher(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages() (+18 more)

### Community 108 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 109 - "New"
Cohesion: 0.18
Nodes (16): New(), generateSelfSignedCertFiles(), TestAdminAPIAuthEnforcement(), TestAdminAuthEnforcedOverHTTPSOnLoopback(), TestDevUIHandlerFallback(), TestGatewaySSEStreamsPastAPITimeout(), testLogger(), TestPanicRecoveryReturnsJSON500() (+8 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (21): @xyflow/react, OverviewRealtimeTrafficProps, EDGE_TYPES, EnergyTopologyEdge(), FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, prefersReducedMotion() (+13 more)

### Community 111 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

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
Cohesion: 0.18
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 116 - "net/url.URL"
Cohesion: 0.07
Nodes (23): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+15 more)

### Community 117 - "Module"
Cohesion: 0.15
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 118 - "RouteMounts"
Cohesion: 0.16
Nodes (4): Module, validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 119 - ".ListModels"
Cohesion: 0.13
Nodes (9): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, isModelAllowed(), isModelDisabled(), TestListVirtualModels(), compositeModelLister (+1 more)

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
Cohesion: 0.07
Nodes (8): prefixReadCloser, streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ssePipeReader

### Community 124 - "translate_openai.go"
Cohesion: 0.18
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+12 more)

### Community 125 - "setupTokenSaverTest"
Cohesion: 0.23
Nodes (13): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestResponsesAPIFunctionCallOutputRTKFilter(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory() (+5 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.11
Nodes (8): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 127 - "strategies.go"
Cohesion: 0.23
Nodes (11): ComboStep, createComboParams, updateComboValues, encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps(), decodeStepsJSON() (+3 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "DB"
Cohesion: 0.04
Nodes (73): Options, authHarness, Options, Options, testHarness, Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Global Constraints, Review Focus (+65 more)

### Community 130 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - ".identifyProviderModel"
Cohesion: 0.09
Nodes (28): ModelIdentityResult, ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+20 more)

### Community 134 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "Module"
Cohesion: 0.23
Nodes (4): matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "applyUpstreamHeaders"
Cohesion: 0.11
Nodes (17): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+9 more)

### Community 140 - "Module"
Cohesion: 0.07
Nodes (19): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+11 more)

### Community 141 - "deviceOAuthSession"
Cohesion: 0.22
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 142 - "ref_bun_test"
Cohesion: 0.11
Nodes (14): modalsSource, providersDir, consoleSource, detailSource, modalsSource, providersDir, consoleSource, customModelModalSource (+6 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "Module"
Cohesion: 0.19
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "newTestHarness"
Cohesion: 0.09
Nodes (24): File Map, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, KeyValidator, assertFixedWidthUTC(), newTestHarness(), sha256Hex() (+16 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.08
Nodes (25): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, zustand (+17 more)

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

### Community 163 - "net/http.Handler"
Cohesion: 0.13
Nodes (20): loadUI(), loadUI(), createLegacySQLite(), setupImporterTest(), TestImporterRejectsInvalidFilesAndPreservesReadOnlySource(), setupPluginsTest(), TestPluginsBuiltinsUpdateAndInterceptorScoping(), putPluginConfig() (+12 more)

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

### Community 178 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "ProviderStore"
Cohesion: 0.32
Nodes (13): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI) (+5 more)

### Community 181 - "RingBuffer"
Cohesion: 0.19
Nodes (11): File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator, cloneAttrMap() (+3 more)

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
Cohesion: 0.14
Nodes (7): JobInspector, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 186 - "New"
Cohesion: 0.14
Nodes (14): newProviderRuntimeCache(), setupProvidersCacheTest(), TestProviderModuleCacheServesAndInvalidatesConcurrently(), TestProviderRuntimeCacheVersionedAndCloned(), newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), TestJetBrainsAdapterAndOAuthRegistry(), New() (+6 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.14
Nodes (17): anthropicResponse, anthropicUsage, errorDisposition, handledGatewayError, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus(), isClientPayloadErrorBody() (+9 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (17): arenaTestEnv, stubCompleter, fakeCredentialStore, AdminCredentialStore, SettingsSnapshotter, stubAdminCredentialStore, ImportRun, stubKeyProvisioner (+9 more)

### Community 194 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 195 - "listSettingsWithQueryer"
Cohesion: 0.24
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "quotaCache"
Cohesion: 0.29
Nodes (7): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry, quotaFlightCall

### Community 198 - ".resolveIdentifier"
Cohesion: 0.18
Nodes (9): resolveSnapshot, applyComboRuntimeMetadata(), buildResolvedTarget(), configIntValue(), connectionSupportsBare(), connectionSupportsExplicit(), isNestedExecuteMode(), withMetadataEntry() (+1 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 201 - "repo"
Cohesion: 0.14
Nodes (11): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+3 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

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
Cohesion: 0.23
Nodes (10): ChatCompletionRequest, ChatCompletionResult, ChatMessage, Global Constraints, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan (+2 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & System, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.09
Nodes (27): Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+19 more)

### Community 219 - "net/http.Header"
Cohesion: 0.13
Nodes (9): memoryResponseWriter, extractVendorHeaders(), extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration() (+1 more)

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 223 - ".Reader"
Cohesion: 0.12
Nodes (16): main(), runMain(), Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+8 more)

### Community 225 - ".streamUpstreamSSE"
Cohesion: 0.29
Nodes (6): openAIStreamEncoder, streamResult, formatSSEFrame(), captureResponsesFrame(), parseSSEFrame(), mapCanonicalStopReasonToOpenAI()

### Community 226 - "AttestationManager"
Cohesion: 0.27
Nodes (6): DefaultAttestationManager(), AttestationManager, isTransientMintStatus(), ResolveMintURL(), inFlightMint, tokenBucket

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 231 - "time.Duration"
Cohesion: 0.10
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, newID()

### Community 235 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 236 - "newTestHarness"
Cohesion: 0.22
Nodes (11): decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry(), TestResolverAllowedModelsAuthorization() (+3 more)

### Community 238 - "oauthProviderStrategy"
Cohesion: 0.24
Nodes (10): applyCodeBuddyProbeHeaders(), codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), TestCodeBuddyTokenRefresh(), newJetBrainsOAuthStrategy(), defaultOAuthStrategies(), newRFC8628Strategy() (+2 more)

### Community 241 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 242 - "compressor.go"
Cohesion: 0.10
Nodes (36): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressRoot(), compressText(), estimateTokens(), hasAnthropicToolResult() (+28 more)

### Community 243 - "trackedResponseWriter"
Cohesion: 0.22
Nodes (3): trackedResponseWriter, apiTimeoutMiddleware(), wrapTrackedResponseWriter()

### Community 246 - "decodeOpenAIResponsesResponse"
Cohesion: 0.25
Nodes (8): responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesResponse, responsesSummaryText, responsesUsage, decodeOpenAIResponsesResponse(), mapResponsesStatusToCanonical()

### Community 249 - "Global Constraints"
Cohesion: 0.33
Nodes (5): Global Constraints, Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), Wave 3C — Local MITM Proxy & PXPIPE Transform Pipeline Implementation Plan

### Community 250 - "App.test.ts"
Cohesion: 0.40
Nodes (4): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ALL_NAV_ITEMS

### Community 251 - "mobileLayout.test.ts"
Cohesion: 0.40
Nodes (4): consoleSource, detailSource, modalsSource, providersDir

### Community 252 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 253 - "HistoryItem"
Cohesion: 0.67
Nodes (3): HistoryItem, historyListResp, historyResponse

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1078 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **39 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProviderDetailView.tsx` to `UpstreamCredentials`, `.handleFetchConnectionModels`, `WriteError`?**
  _High betweenness centrality (0.222) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `combos/module.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `UpstreamCredentials` to `context.Context`, `net/http.Request`, `go_pkg_encoding_json`, `oauthProviderStrategy`, `WriteError`, `.handleFetchConnectionModels`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.184) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `ProviderDetailView.tsx` to `Button`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._