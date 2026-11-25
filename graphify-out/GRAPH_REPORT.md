# Graph Report - omeg4mon-router  (2026-11-26)

## Corpus Check
- 441 files · ~483,700 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5253 nodes · 16946 edges · 249 communities (205 shown, 44 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1552 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f79c82be`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- registry.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- MitmConsole.tsx
- LiveMonitorConsole.tsx
- TypeScript Strict Compiler Config
- CombosConsole.tsx
- Button
- testing.T
- shadcn/ui Component Configuration
- ProvidersConsole.tsx
- MediaConsole.tsx
- App.tsx
- ProviderDetailView.tsx
- Design System & Wave 0 Plan Docs
- TelegramConsole.tsx
- updates
- ir.go
- trafficShaper
- with
- jobs
- WriteError
- steps
- web-check
- on
- net/http.Request
- go_pkg_encoding_json
- navigation.ts
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- setupTestMitm
- net/http.Header
- Go Module Root Definition
- CanonicalUsage
- New
- KeysConsole.tsx
- translate_anthropic.go
- New
- CanFastPath
- translate_gemini.go
- buildApp
- repository
- New
- translate_openai.go
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- main.go
- SystemOverview.tsx
- extractErrorMessage
- dependencies
- oauth_registry.go
- media/module.go
- ProvidersModals.tsx
- .ValidateKey
- UsageConsole.tsx
- repository
- UpstreamCredentials
- teeHandler
- ProxyPoolsConsole.tsx
- Module
- scripts
- pricingStore
- Module
- TokenSaverConsole.tsx
- .AuthenticateAdmin
- Module
- newTestHarness
- go_pkg_encoding_json
- promptinjection/engine.go
- Module
- Runner
- Load
- web/package.json
- Module
- Module
- repository
- paths
- collectDynamicHostTelemetry
- Module
- arena/hooks.ts
- settingDefinition
- repo
- .ProvisionKey
- strategyState
- .identifyProviderModel
- repository
- SettingsConsole.tsx
- awsQuickRuntime
- ProcessHandle
- combos/catalog.go
- oauthLoopbackProxy
- net/http.ResponseWriter
- omnirouteModelsParity.test.ts
- repository
- WriteJSON
- .Reader
- net/url.URL
- .UpdateCombo
- post
- ref_bun_test
- time.Duration
- fakeHandle
- ProviderTopology.tsx
- repo
- Module
- repo
- properties
- ref_node_path
- log/slog.Logger
- Config
- modeleditor/engine.go
- TestContractInterfacesSatisfied
- .UpsertPlugin
- properties
- repo
- ModelCatalogConsole.tsx
- ModelDescriptor
- newTestHarness
- newTestHarness
- OMEG4MON Router
- stubCompleter
- Module
- certManager
- repo
- applyUpstreamHeaders
- Module
- properties
- release-binaries
- steps
- evaluateAntigravityQuota
- .ExportCombos
- File Map
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- 200
- io/fs.FS
- properties
- cli/package.json
- strategies.go
- setupNotificationsTest
- get
- scripts
- post
- schemas
- .executeCandidateAttempt
- PreviewPluginsResponse
- time.Time
- .executeBuiltinFetch
- modelIntelligence.ts
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- headroomCache
- Review Focus
- Review Focus
- stubModule
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- scanCombo
- version
- Module
- NewStaticUIHandler
- parameters
- NotificationsConsole.tsx
- providerStrategy.ts
- stubSystem
- ClearImporterHistoryResponse
- source_path
- Review Focus
- telemetryCollector
- docker-release
- DB
- `@omeg4mon/router` CLI Launcher
- parseJSONObjectOrWrap
- setupModelEditorTest
- context.Context
- nameMinHeap
- stubDep
- parameters
- io.Reader
- net/http.ResponseWriter
- omeg4mon-router Verification Gate
- stubSystemModule
- .resolveIdentifier
- single-binary
- keyLimiter
- verify-web
- CanonicalStreamChunk
- smoke-test.sh
- .InsertRun
- readLegacySQLite
- responseCoalescer
- plugins/interceptor.go
- .CompleteChat
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- palette.test.ts
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- net/http/httptest.ResponseRecorder
- batcher
- .Record
- repository
- go_pkg_encoding_json
- jsonMarshalerError
- Open
- compressor.go
- omeg4mon-router Project Standards

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 194 edges
2. `WriteJSON()` - 184 edges
3. `WriteError()` - 182 edges
4. `DB` - 93 edges
5. `Button` - 90 edges
6. `buildApp()` - 79 edges
7. `Bus` - 72 edges
8. `Alert` - 69 edges
9. `Badge()` - 68 edges
10. `AlertDescription` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/platform/lifecycle/registry.go
- `Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch)` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3b-media-tunnel.md → internal/platform/lifecycle/registry.go
- `Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs)` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3c-mitm-pxpipe.md → internal/platform/lifecycle/registry.go
- `Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI)` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3d-telegram-arena-chat.md → internal/platform/lifecycle/registry.go

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 3-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (249 total, 44 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.07
Nodes (40): ModelOption, handledGatewayError, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, mapBusEventToAudit() (+32 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.22
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
Nodes (23): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary() (+15 more)

### Community 6 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.09
Nodes (40): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+32 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "CombosConsole.tsx"
Cohesion: 0.19
Nodes (22): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+14 more)

### Community 10 - "Button"
Cohesion: 0.07
Nodes (67): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Functional & UX Scope, lucide-react (+59 more)

### Community 11 - "testing.T"
Cohesion: 0.05
Nodes (71): testEnv, TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies() (+63 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "ProvidersConsole.tsx"
Cohesion: 0.09
Nodes (42): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets(), FleetTelemetry, getEntrySortRank() (+34 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (37): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+29 more)

### Community 16 - "ProviderDetailView.tsx"
Cohesion: 0.10
Nodes (63): class-variance-authority, ArenaConsole(), ChatConsole(), useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+55 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.13
Nodes (33): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, StreamTranslator (+25 more)

### Community 21 - "trafficShaper"
Cohesion: 0.43
Nodes (3): trafficShaper, acquireSem(), sleepContext()

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "WriteError"
Cohesion: 0.07
Nodes (26): listCooldownsResponse, Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`) (+18 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "net/http.Request"
Cohesion: 0.07
Nodes (17): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), extractAPIKey(), ensureRequestID() (+9 more)

### Community 29 - "go_pkg_encoding_json"
Cohesion: 0.24
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 30 - "navigation.ts"
Cohesion: 0.14
Nodes (21): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, RawModuleEntry, UILocale, ActiveNavContext, EXTENSION_GROUPS, findActiveNavContext(), isModuleRegistered() (+13 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 35 - "net/http.Header"
Cohesion: 0.18
Nodes (6): memoryResponseWriter, extractVendorHeaders(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamDecoder, streamEncoder (+8 more)

### Community 39 - "New"
Cohesion: 0.09
Nodes (25): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+17 more)

### Community 40 - "KeysConsole.tsx"
Cohesion: 0.06
Nodes (68): react-hook-form, buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt() (+60 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.14
Nodes (22): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+14 more)

### Community 42 - "New"
Cohesion: 0.13
Nodes (19): Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, New(), newHarness(), TestTranslatorLogsSnapshotAndLiveSSEStream(), TestTranslatorLogsStreamSurvivesPastAPITimeout(), TestTranslatorTranslateRequestAndResponseAllPairs() (+11 more)

### Community 43 - "CanFastPath"
Cohesion: 0.16
Nodes (15): ProtocolTranslator, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), IsSupportedProtocol(), RewriteModelField(), decodeOpenAIResponsesRequest() (+7 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "buildApp"
Cohesion: 0.05
Nodes (54): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+46 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "New"
Cohesion: 0.06
Nodes (42): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints() (+34 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.19
Nodes (18): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+10 more)

### Community 54 - "main.go"
Cohesion: 0.23
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.23
Nodes (11): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+3 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.05
Nodes (69): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+61 more)

### Community 57 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 58 - "oauth_registry.go"
Cohesion: 0.12
Nodes (33): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies() (+25 more)

### Community 59 - "media/module.go"
Cohesion: 0.06
Nodes (28): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+20 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.04
Nodes (112): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, ARENA_MODELS_QUERY_KEY, DiscoveredModel (+104 more)

### Community 61 - ".ValidateKey"
Cohesion: 0.22
Nodes (5): KeyInfo, KeyValidator, stubKeyValidator, stubKeyValidator, stubKeyValidator

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (70): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), @hookform/resolvers, useProviderConnections(), useProxyPools(), useSystemTelemetry(), formatActivityOriginLabel(), LazyProviderTopology, OverviewRealtimeTraffic() (+62 more)

### Community 63 - "repository"
Cohesion: 0.10
Nodes (20): cloneStringMap(), cloneUpstreamCredentials(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+12 more)

### Community 64 - "UpstreamCredentials"
Cohesion: 0.10
Nodes (7): catalogFakeProviderStore, ProviderConnection, ProviderStore, stubProviderStore, UpstreamCredentials, stubProviderStore, stubProviderStore

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.12
Nodes (31): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation(), useUpsertProxyBindingMutation() (+23 more)

### Community 67 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "Module"
Cohesion: 0.16
Nodes (3): File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Module

### Community 71 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.07
Nodes (8): AuthVerifier, stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "newTestHarness"
Cohesion: 0.10
Nodes (18): Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, newKeyCache(), newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), New(), assertFixedWidthUTC(), newTestHarness() (+10 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.09
Nodes (24): providerIDContextKey, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), buildDiscoveryCandidateURLs(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates() (+16 more)

### Community 76 - "promptinjection/engine.go"
Cohesion: 0.06
Nodes (34): antigravityPromptRewriteRule, actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody() (+26 more)

### Community 77 - "Module"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 78 - "Runner"
Cohesion: 0.16
Nodes (8): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 79 - "Load"
Cohesion: 0.22
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.10
Nodes (27): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+19 more)

### Community 81 - "Module"
Cohesion: 0.15
Nodes (6): AdminCredentialStore, clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 82 - "Module"
Cohesion: 0.20
Nodes (9): extractUsageCost(), actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+1 more)

### Community 83 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.16
Nodes (19): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+11 more)

### Community 86 - "Module"
Cohesion: 0.25
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 87 - "arena/hooks.ts"
Cohesion: 0.10
Nodes (28): arenaQueryKeys, useUpdateChatSessionMutation(), ArenaBattle, ArenaBattleFormInput, ArenaBattleFormOutput, arenaBattleFormSchema, ArenaBattleListResponse, ArenaContenderResult (+20 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.06
Nodes (26): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), listSettingsWithQueryer() (+18 more)

### Community 90 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 91 - "strategyState"
Cohesion: 0.22
Nodes (7): strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), fastRandFloat64(), fastRandIntN()

### Community 92 - ".identifyProviderModel"
Cohesion: 0.13
Nodes (16): ModelIdentityResult, ModelTestResult, testAllModelsResponse, buildModelProbeRequest(), extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+8 more)

### Community 93 - "repository"
Cohesion: 0.17
Nodes (10): boolToInt(), marshalOverrideJSON(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule, ModelOverride (+2 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): zod, zod@^4.6.7, authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+36 more)

### Community 95 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 96 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 97 - "combos/catalog.go"
Cohesion: 0.24
Nodes (11): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, testAllModelsRequest, testComboRequest (+3 more)

### Community 98 - "oauthLoopbackProxy"
Cohesion: 0.20
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (33): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol(), validateHTTPURL() (+25 more)

### Community 100 - "omnirouteModelsParity.test.ts"
Cohesion: 0.22
Nodes (8): appSource, comboFormSource, consoleSource, hooksSource, legacyConsolePath, modalSource, modelsDir, strategyGuideSource

### Community 101 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 102 - "WriteJSON"
Cohesion: 0.11
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON(), LogEntry (+3 more)

### Community 103 - ".Reader"
Cohesion: 0.05
Nodes (49): main(), runMain(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan (+41 more)

### Community 104 - "net/url.URL"
Cohesion: 0.09
Nodes (20): stubProxyDialerResolver, File Map, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, stubProxyResolver, ProxyDialerResolver, copyNonHopHeaders(), matchToolForHost() (+12 more)

### Community 105 - ".UpdateCombo"
Cohesion: 0.29
Nodes (9): createComboParams, updateComboValues, cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps(), encodeConfigJSON() (+1 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "ref_bun_test"
Cohesion: 0.11
Nodes (15): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, consoleSource, detailSource, modalsSource, providersDir, modalsSource (+7 more)

### Community 108 - "time.Duration"
Cohesion: 0.10
Nodes (9): AccountHealthTracker, stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore(), computeJitteredRetryDelay() (+1 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (19): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+11 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "Module"
Cohesion: 0.04
Nodes (21): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), cloneStringSlice(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), xrouterBuiltInOAuthAdapters() (+13 more)

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "ref_node_path"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 116 - "log/slog.Logger"
Cohesion: 0.15
Nodes (4): Module, validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 118 - "Config"
Cohesion: 0.06
Nodes (35): Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), setupSecurityTest() (+27 more)

### Community 119 - "modeleditor/engine.go"
Cohesion: 0.35
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 120 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (7): CanonicalRequest, GatewayInterceptor, ResponseMeta, StreamEvent, stubInterceptor, recordingInterceptor, TestContractInterfacesSatisfied()

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 124 - "ModelCatalogConsole.tsx"
Cohesion: 0.10
Nodes (57): AssignKeysModalProps, classifyTestError(), clearCachedTestResults(), comboTestResultKey(), computeUpdatedAllowedModels(), computeVisibilityCounts(), countKeyAccess(), DEFAULT_COMBO_FILTERS (+49 more)

### Community 125 - "ModelDescriptor"
Cohesion: 0.18
Nodes (6): stubModelLister, ModelDescriptor, ModelLister, stubModelLister, stubModelLister, stubModelLister

### Community 126 - "newTestHarness"
Cohesion: 0.08
Nodes (24): fakeClock, fakeProviderStore, stubAdminAuth, testHarness, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), decodeJSON(), newFakeClock(), newFakeProviderStore() (+16 more)

### Community 127 - "newTestHarness"
Cohesion: 0.31
Nodes (10): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, New(), indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC() (+2 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "stubCompleter"
Cohesion: 0.25
Nodes (4): arenaTestEnv, stubCompleter, ChatCompleter, ChatCompletionResult

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "applyUpstreamHeaders"
Cohesion: 0.17
Nodes (11): Task 2: Backend Model Test Endpoint (TDD), Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isBlockedAntigravityProjectID() (+3 more)

### Community 134 - "Module"
Cohesion: 0.15
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 140 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 141 - "File Map"
Cohesion: 0.25
Nodes (7): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

### Community 142 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "io/fs.FS"
Cohesion: 0.02
Nodes (13): migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), stubSystemModule, stubModule, orderStubModule (+5 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "strategies.go"
Cohesion: 0.31
Nodes (7): ComboStep, ProviderStrategyState, validateComboTargets(), decodeStepsJSON(), normalizeAndValidateSteps(), stepsFromTargets(), validateProviderStrategy()

### Community 154 - "setupNotificationsTest"
Cohesion: 0.16
Nodes (15): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict() (+7 more)

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

### Community 159 - ".executeCandidateAttempt"
Cohesion: 0.16
Nodes (17): attemptOutcome, CanonicalChatResponse, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), EncodeResponse(), rawJSONToString() (+9 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "time.Time"
Cohesion: 0.15
Nodes (10): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, newCooldownTracker(), newFakeClock(), refreshCooldownEntry (+2 more)

### Community 162 - ".executeBuiltinFetch"
Cohesion: 0.29
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 163 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

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

### Community 175 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 178 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 181 - "providerStrategy.ts"
Cohesion: 0.14
Nodes (20): ComboStrategy, AUTO_ROUTER_STRATEGIES, COMBO_STRATEGY_OPTIONS, ComboStrategyMeta, COMPRESSION_MODES, STRATEGY_CATEGORIES, StrategyCategory, StrategyGuideModalProps (+12 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "Review Focus"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 186 - "telemetryCollector"
Cohesion: 0.48
Nodes (5): collectDynamicHostTelemetry(), TestTelemetryCacheAvoidsRedundantWork(), hostTelemetrySnapshot, staticHostInfo, telemetryCollector

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "DB"
Cohesion: 0.05
Nodes (60): Options, authHarness, Options, Global Constraints, Wave 3B — Media Gateway (Audio/Images/Videos/Search/Fetch) & Tunnel Manager Implementation Plan, Options, Options, New() (+52 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "parseJSONObjectOrWrap"
Cohesion: 0.31
Nodes (5): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini(), decodeOpenAIChatRequest()

### Community 191 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (20): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins() (+12 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (17): fakeCredentialStore, ProviderResolver, ResolvedTarget, SettingsSnapshotter, stubAdminCredentialStore, stubProviderResolver, stubResolver, stubSettings (+9 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 198 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (9): Module, trackedResponseWriter, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), parseInboundMeta() (+1 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 201 - ".resolveIdentifier"
Cohesion: 0.20
Nodes (6): resolveSnapshot, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), isModelAllowed(), isModelDisabled()

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.33
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 207 - ".InsertRun"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 209 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 210 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 211 - "plugins/interceptor.go"
Cohesion: 0.43
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 212 - ".CompleteChat"
Cohesion: 0.50
Nodes (3): extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 216 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 223 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.15
Nodes (20): Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+12 more)

### Community 224 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

### Community 225 - ".Record"
Cohesion: 0.14
Nodes (6): stubUsageRecorder, stubUsageRecorder, UsageRecord, UsageRecorder, stubUsageRecorder, recordingUsage

### Community 229 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 235 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (19): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+11 more)

### Community 241 - "Open"
Cohesion: 0.08
Nodes (26): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestProbeAndIdentifyRecordUsageTelemetry(), New(), TestModuleMetadataAndStartValidation(), newRepository(), newStrategyState() (+18 more)

### Community 242 - "compressor.go"
Cohesion: 0.12
Nodes (28): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+20 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1023 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **44 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `ProviderDetailView.tsx`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.207) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0735042735042735 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `net/http.ResponseWriter` to `context.Context`, `go_pkg_encoding_json`, `ProvidersModals.tsx`, `net/http.Request`?**
  _High betweenness centrality (0.172) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `ProvidersModals.tsx` to `ProviderDetailView.tsx`, `Button`, `ProvidersConsole.tsx`?**
  _High betweenness centrality (0.074) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._