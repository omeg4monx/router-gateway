# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 477 files · ~512,007 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5516 nodes · 18505 edges · 232 communities (208 shown, 24 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1820 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a8daec56`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main.go
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
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- Button
- updates
- TranslatorConsole.tsx
- setupMediaTest
- with
- jobs
- Module
- steps
- web-check
- on
- Module
- MitmConsole.tsx
- CanonicalChatResponse
- NotificationsConsole.tsx
- docker-build
- go-check
- KeysConsole.tsx
- .handleCreateConnection
- Go Module Root Definition
- CanonicalUsage
- newTestHarness
- QuotaConsole.tsx
- CanonicalStreamChunk
- go_pkg_context
- ir.go
- translate_gemini.go
- buildApp
- repository
- Module
- translate_openai.go
- ModelDescriptor
- NewClient
- ModelsRoutingConsole.tsx
- translate_responses.go
- oauthLoopbackProxy
- .Fetch
- SystemOverview.tsx
- PromptGuardConsole.tsx
- dependencies
- repository
- server_test.go
- .executeCandidateAttempt
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
- awsQuickRuntime
- WriteJSON
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- newTestHarness
- paths
- collectDynamicHostTelemetry
- fakeHandle
- CombosConsole.tsx
- settingDefinition
- repo
- .Reader
- Module
- repo
- ResolvedTarget
- client.ts
- New
- Runner
- main.go
- Module
- net/http.Request
- deviceOAuthSession
- ProviderConnection
- Module
- Module
- repository
- ref_bun_test
- post
- .identifyProviderModel
- DB
- OverviewRealtimeTraffic.tsx
- ProviderTopology.tsx
- repo
- repository
- repo
- properties
- stubDep
- net/url.URL
- redact.go
- net/http.ResponseWriter
- settingDefinition
- time.Time
- .UpsertPlugin
- properties
- time.Duration
- Review Focus
- setupTestMitm
- io/fs.FS
- stubModule
- OMEG4MON Router
- go_pkg_time
- ProviderStore
- certManager
- repo
- Module
- Module
- properties
- setupTokenSaverTest
- release-binaries
- steps
- evaluateAntigravityQuota
- repository
- sync.RWMutex
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- stubSystemModule
- properties
- 200
- Module
- properties
- GatewayInterceptor
- cli/package.json
- quotaCache
- Provider Quota Menu & Token Harbor Quota Integration
- net/http.Header
- get
- scripts
- post
- schemas
- setupCatalogTestModule
- PreviewPluginsResponse
- stubCompleter
- modeleditor/engine.go
- analyzeIdentityProbeResponse
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- buildModelProbeRequest
- Review Focus
- Review Focus
- SessionSigner
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- UsageRecorder
- version
- parseJSONObjectOrWrap
- vite.config.ts
- parameters
- orderStubModule
- New
- stubDiscoverySystemModule
- ClearImporterHistoryResponse
- source_path
- stubSystemModule
- AttestationManager
- docker-release
- ChatCompleter
- `@omeg4mon/router` CLI Launcher
- repository
- setupModelEditorTest
- context.Context
- schema.go
- parameters
- New
- WithProviderID
- omeg4mon-router Verification Gate
- headroomCache
- listSettingsWithQueryer
- single-binary
- keyLimiter
- verify-web
- promptinjection/engine.go
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- .ProvisionKey
- Feature Suite
- responseCoalescer
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Load
- batcher
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- net/http/httptest.ResponseRecorder
- AdminCredentialStore
- plugins/interceptor.go
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- time.Duration
- jsonMarshalerError
- compressor.go
- omeg4mon-router Project Standards

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 192 edges
2. `WriteJSON()` - 186 edges
3. `WriteError()` - 184 edges
4. `DB` - 93 edges
5. `Button` - 92 edges
6. `buildApp()` - 79 edges
7. `Module` - 78 edges
8. `cn()` - 75 edges
9. `Bus` - 72 edges
10. `Alert` - 71 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse)` --references--> `KeyProvisioner`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3a-teams-vouchers.md → internal/contracts/contracts.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
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

## Communities (232 total, 24 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.07
Nodes (22): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+14 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.17
Nodes (9): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC(), parseAndValidateProxyURL(), redactRawCandidate(), runConcurrentScan() (+1 more)

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
Nodes (26): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+18 more)

### Community 6 - "repository"
Cohesion: 0.17
Nodes (10): boolToInt(), marshalOverrideJSON(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule, ModelOverride (+2 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.17
Nodes (21): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+13 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "ProvidersModals.tsx"
Cohesion: 0.02
Nodes (199): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, classifyTestError(), comboTestResultKey() (+191 more)

### Community 10 - "Button"
Cohesion: 0.04
Nodes (107): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope (+99 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (88): Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), testEnv, TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine() (+80 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.18
Nodes (6): catalogFakeProviderStore, stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), refreshCall

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (20): mediaQueryKeys, CreateMediaComboFormInput, CreateMediaComboFormOutput, CreateMediaComboRequest, createMediaComboSchema, CreateMediaProviderFormInput, CreateMediaProviderFormOutput, CreateMediaProviderRequest (+12 more)

### Community 15 - "App.tsx"
Cohesion: 0.09
Nodes (35): App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio, InitPage (+27 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.09
Nodes (37): isBlockedAntigravityProjectID(), generateTraceID(), parseAwsQuickPastedQuery(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload() (+29 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Button"
Cohesion: 0.06
Nodes (108): @hookform/resolvers, react-hook-form, ArenaConsole(), useArenaBattles(), useArenaLeaderboard(), useDeleteArenaBattleMutation(), useRunArenaBattleMutation(), useVoteArenaBattleMutation() (+100 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "TranslatorConsole.tsx"
Cohesion: 0.17
Nodes (19): translatorQueryKeys, useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS, translateFormSchema (+11 more)

### Community 21 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

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
Cohesion: 0.10
Nodes (8): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 29 - "MitmConsole.tsx"
Cohesion: 0.06
Nodes (66): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+58 more)

### Community 30 - "CanonicalChatResponse"
Cohesion: 0.13
Nodes (24): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString() (+16 more)

### Community 31 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "KeysConsole.tsx"
Cohesion: 0.06
Nodes (66): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt(), KeyExpiryState (+58 more)

### Community 35 - ".handleCreateConnection"
Cohesion: 0.05
Nodes (33): Global Constraints, Approaches Considered, buildProtocolTestRequest(), cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), isValidAuthType(), isValidProtocol() (+25 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "newTestHarness"
Cohesion: 0.07
Nodes (29): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, newKeyCache(), newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), New(), assertFixedWidthUTC() (+21 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.11
Nodes (41): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+33 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.13
Nodes (16): anthropicStreamDecoder, CanonicalStreamChunk, openAIStreamDecoder, responsesStreamDecoder, streamDecoder, NewStreamTranslator(), mapAnthropicStopReasonToCanonical(), newAnthropicStreamDecoder() (+8 more)

### Community 42 - "go_pkg_context"
Cohesion: 0.18
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 43 - "ir.go"
Cohesion: 0.12
Nodes (39): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, anthropicContentBlock, anthropicImageSource, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource (+31 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "buildApp"
Cohesion: 0.06
Nodes (35): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+27 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "Module"
Cohesion: 0.16
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 48 - "translate_openai.go"
Cohesion: 0.20
Nodes (18): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+10 more)

### Community 49 - "ModelDescriptor"
Cohesion: 0.31
Nodes (6): buildServiceWorkerSource(), CORE_PRECACHE_URLS, classifyFetchRequest(), FetchRouteKind, isNetworkOnlyPathname(), NETWORK_ONLY_PREFIXES

### Community 50 - "NewClient"
Cohesion: 0.08
Nodes (20): arenaTestEnv, Options, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), TestQuotaMatchesCustomTokenHarborHost(), startAuthHTTPProxy() (+12 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (24): anthropicChoice, anthropicMessage, anthropicRequest, anthropicThinking, anthropicTool, responsesContentPart, responsesFormat, responsesIncompleteDetails (+16 more)

### Community 53 - "oauthLoopbackProxy"
Cohesion: 0.21
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 54 - ".Fetch"
Cohesion: 0.13
Nodes (20): TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+12 more)

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.21
Nodes (12): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+4 more)

### Community 56 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 57 - "dependencies"
Cohesion: 0.07
Nodes (27): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+19 more)

### Community 58 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 59 - "server_test.go"
Cohesion: 0.07
Nodes (56): ModelOption, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse (+48 more)

### Community 60 - ".executeCandidateAttempt"
Cohesion: 0.09
Nodes (25): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), extractUpstreamErrorMessage(), parseRetryAfter(), applyCodeBuddyHeaders() (+17 more)

### Community 61 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.10
Nodes (35): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), buildLogsQueryParams(), INVALIDATION_TOPIC_QUERY_KEYS, UsageLogsFilter, usageQueryKeys, useUsageLogs(), computeBarPercent(), computeSharePercent() (+27 more)

### Community 63 - "repository"
Cohesion: 0.09
Nodes (21): isCodeBuddyHostFor(), isTokenHarborHost(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+13 more)

### Community 64 - "repository"
Cohesion: 0.12
Nodes (22): Combo, ComboStep, createComboParams, ProviderStrategyState, rowScanner, updateComboValues, cloneStringSlice(), decodeStringSlice() (+14 more)

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
Cohesion: 0.12
Nodes (22): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ActiveNavContext, ALL_NAV_ITEMS, EXTENSION_GROUPS (+14 more)

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
Cohesion: 0.16
Nodes (17): modelEditorQueryKeys, BYPASS_ACTIONS, CreateBypassRuleFormValues, createBypassRuleSchema, CreateModelBypassRuleRequest, CreateModelOverrideFormValues, CreateModelOverrideRequest, createModelOverrideSchema (+9 more)

### Community 73 - "Module"
Cohesion: 0.13
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "JobInspector"
Cohesion: 0.10
Nodes (26): TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification (+18 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.12
Nodes (27): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+19 more)

### Community 76 - "arena/hooks.ts"
Cohesion: 0.08
Nodes (46): ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaModels(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteChatSessionMutation() (+38 more)

### Community 77 - "Module"
Cohesion: 0.19
Nodes (7): compileRule(), decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), compiledRule, Module

### Community 78 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 79 - "WriteJSON"
Cohesion: 0.09
Nodes (18): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), validateSQLiteSource(), decodeJSON(), decodeJSONBody() (+10 more)

### Community 80 - "web/package.json"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.13
Nodes (13): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, Approaches Considered, Builder Decisions & Deliverables (+5 more)

### Community 82 - "Module"
Cohesion: 0.20
Nodes (7): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken()

### Community 83 - "newTestHarness"
Cohesion: 0.15
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 86 - "fakeHandle"
Cohesion: 0.09
Nodes (7): streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ssePipeReader

### Community 87 - "CombosConsole.tsx"
Cohesion: 0.10
Nodes (40): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+32 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.14
Nodes (11): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+3 more)

### Community 90 - ".Reader"
Cohesion: 0.14
Nodes (15): main(), runMain(), Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints, Global Constraints, Review Focus (+7 more)

### Community 91 - "Module"
Cohesion: 0.18
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 92 - "repo"
Cohesion: 0.27
Nodes (6): buildWhere(), newEventID(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 93 - "ResolvedTarget"
Cohesion: 0.18
Nodes (10): strategyState, weightedCandidate, stubProviderResolver, stubResolver, ResolvedTarget, applyComboStrategy(), applyProviderStrategy(), applyStrategy() (+2 more)

### Community 94 - "client.ts"
Cohesion: 0.04
Nodes (64): @tanstack/react-query, zod, zod@^4.6.7, authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+56 more)

### Community 95 - "New"
Cohesion: 0.11
Nodes (21): loadUI(), loadUI(), apiTimeoutMiddleware(), New(), securityHeadersMiddleware(), generateSelfSignedCertFiles(), TestAdminAPIAuthEnforcement(), TestAdminAuthEnforcedOverHTTPSOnLoopback() (+13 more)

### Community 96 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "main.go"
Cohesion: 0.15
Nodes (16): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+8 more)

### Community 98 - "Module"
Cohesion: 0.18
Nodes (5): decodeStrictJSON(), LogEntry, LogListResponse, Module, Status

### Community 99 - "net/http.Request"
Cohesion: 0.05
Nodes (29): Module, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), Module, clientRateLimitKey() (+21 more)

### Community 100 - "deviceOAuthSession"
Cohesion: 0.17
Nodes (8): codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), originOfURL(), codebuddyAuthSpec, deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 101 - "ProviderConnection"
Cohesion: 0.16
Nodes (9): fakeProviderStore, resolveSnapshot, ProviderConnection, newFakeProviderStore(), buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), IsVeniceProviderID() (+1 more)

### Community 102 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 104 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 105 - "ref_bun_test"
Cohesion: 0.06
Nodes (28): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+20 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.11
Nodes (24): ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), generateID(), extractSSETokenUsage() (+16 more)

### Community 108 - "DB"
Cohesion: 0.06
Nodes (59): Options, authHarness, Options, Options, testHarness, Global Constraints, Wave 3B — Media Gateway (Audio/Images/Videos/Search/Fetch) & Tunnel Manager Implementation Plan, New() (+51 more)

### Community 109 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.29
Nodes (12): formatActivityOriginLabel(), LazyProviderTopology, OverviewRealtimeTraffic(), OverviewRealtimeTrafficProps, RecentRequestsRows(), FilteredRecentRequestsResult, filterRecentRequests(), formatRelativeTimeShort() (+4 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.15
Nodes (18): @xyflow/react, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps, RouterHubNode() (+10 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 116 - "net/url.URL"
Cohesion: 0.10
Nodes (16): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+8 more)

### Community 117 - "redact.go"
Cohesion: 0.33
Nodes (7): antigravityPromptRewriteRule, isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 118 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (14): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), trackedResponseWriter, providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam() (+6 more)

### Community 119 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 120 - "time.Time"
Cohesion: 0.09
Nodes (14): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, newFakeClock() (+6 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "time.Duration"
Cohesion: 0.10
Nodes (13): chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore(), TestTrafficShaperAdaptiveHeaderPacing(), TestTrafficShaperLimitsPerConnectionConcurrencyAndQueues() (+5 more)

### Community 124 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 125 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.03
Nodes (15): migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile() (+7 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "go_pkg_time"
Cohesion: 0.21
Nodes (12): anthropicResponse, anthropicUsage, errorDisposition, handledGatewayError, classifyStatus(), isClientPayloadErrorBody(), compareSemver(), maybeLiftFreePIClientVersion() (+4 more)

### Community 130 - "ProviderStore"
Cohesion: 0.30
Nodes (14): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI) (+6 more)

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

### Community 136 - "setupTokenSaverTest"
Cohesion: 0.27
Nodes (10): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory(), TestTokenSaverFinalReviewFindings() (+2 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 140 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 142 - "repository"
Cohesion: 0.18
Nodes (7): builtinRules(), boolToInt(), defaultConfig(), Config, GuardEvent, repository, Rule

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

### Community 148 - "Module"
Cohesion: 0.19
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "GatewayInterceptor"
Cohesion: 0.18
Nodes (11): Global Constraints, Review Focus, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Global Constraints, Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases) (+3 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "quotaCache"
Cohesion: 0.29
Nodes (7): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotaBucket, ProviderQuotasResponse, quotaCache, quotaCacheEntry

### Community 153 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.18
Nodes (9): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+1 more)

### Community 154 - "net/http.Header"
Cohesion: 0.19
Nodes (5): memoryResponseWriter, copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

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

### Community 159 - "setupCatalogTestModule"
Cohesion: 0.17
Nodes (9): setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies(), TestAllOmniRouteComboStrategiesCRUDAndPersistence(), TestOmniRouteResolverStrategiesExecution() (+1 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "stubCompleter"
Cohesion: 0.29
Nodes (7): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 162 - "modeleditor/engine.go"
Cohesion: 0.35
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 163 - "analyzeIdentityProbeResponse"
Cohesion: 0.22
Nodes (10): ModelIdentityResult, analyzeIdentityProbeResponse(), extractIDPrefix(), extractVendorHeaders(), inferModelFamily(), modelsLooselyMatch(), normalizeModelForCompare(), parseIdentityPayload() (+2 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

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

### Community 175 - "UsageRecorder"
Cohesion: 0.29
Nodes (7): Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan, Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), UsageRecorder

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), decodeOpenAIChatRequest()

### Community 178 - "vite.config.ts"
Cohesion: 0.25
Nodes (7): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, omeg4monPwaPlugin()

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 181 - "New"
Cohesion: 0.06
Nodes (36): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+28 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 186 - "AttestationManager"
Cohesion: 0.12
Nodes (15): bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), DefaultAttestationManager(), getSharedCookieJar(), AttestationManager, isTransientMintStatus() (+7 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "ChatCompleter"
Cohesion: 0.53
Nodes (5): Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan, ChatCompleter

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 191 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (16): fakeCredentialStore, stubAdminCredentialStore, stubSettings, extractUsageCost(), TestSystemAdminCredentialStore(), validatePassword(), sanitizeBotError(), HashPassword() (+8 more)

### Community 194 - "schema.go"
Cohesion: 0.33
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "New"
Cohesion: 0.05
Nodes (55): roundTripFunc, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestKeylessFreeProviderModelProbeAndIdentify(), TestProbeAndIdentifyRecordUsageTelemetry(), newCooldownTracker() (+47 more)

### Community 198 - "WithProviderID"
Cohesion: 0.10
Nodes (17): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+9 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

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

### Community 205 - "promptinjection/engine.go"
Cohesion: 0.33
Nodes (11): actionRank(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity(), redactAndTruncate() (+3 more)

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 209 - ".ProvisionKey"
Cohesion: 0.08
Nodes (27): Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification (+19 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 211 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

### Community 219 - "batcher"
Cohesion: 0.11
Nodes (9): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher, ModelPricing (+1 more)

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 222 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.12
Nodes (12): doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth(), TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping(), doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture() (+4 more)

### Community 223 - "AdminCredentialStore"
Cohesion: 0.09
Nodes (34): TestAppWave1AAuthSettingsKeysEndToEnd(), File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, AdminCredentialStore, AuthVerifier (+26 more)

### Community 226 - "plugins/interceptor.go"
Cohesion: 0.43
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 231 - "time.Duration"
Cohesion: 0.06
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID(), computePrefixFingerprint()

### Community 242 - "compressor.go"
Cohesion: 0.13
Nodes (27): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+19 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1029 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `.handleCreateConnection`?**
  _High betweenness centrality (0.187) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0673903211216644 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `.handleCreateConnection` to `context.Context`, `net/http.Request`, `ProvidersModals.tsx`, `go_pkg_encoding_json`, `testing.T`, `UpstreamCredentials`, `oauth_registry.go`, `.Fetch`?**
  _High betweenness centrality (0.157) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `AdminCredentialStore` to `ProviderStore`, `Module`, `newTestHarness`, `SessionSigner`, `Button`?**
  _High betweenness centrality (0.082) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._