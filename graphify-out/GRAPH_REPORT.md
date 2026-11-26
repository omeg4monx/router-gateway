# Graph Report - omeg4mon-router  (2026-11-26)

## Corpus Check
- 443 files · ~486,811 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5282 nodes · 17052 edges · 259 communities (212 shown, 47 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1557 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c70d9902`
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
- net/http.ResponseWriter
- MediaConsole.tsx
- App.tsx
- Module
- Design System & Wave 0 Plan Docs
- TelegramConsole.tsx
- updates
- translate_anthropic.go
- combos/module.go
- with
- jobs
- WriteError
- steps
- web-check
- on
- ProviderDetailView.tsx
- TranslatorConsole.tsx
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
- CanonicalStreamChunk
- PromptGuardConsole.tsx
- ir.go
- translate_gemini.go
- buildApp
- repository
- New
- translate_openai.go
- Module
- ModelsRoutingConsole.tsx
- translate_responses.go
- go_pkg_context
- SystemOverview.tsx
- extractErrorMessage
- dependencies
- oauthTokenResponse
- main_test.go
- ProvidersModals.tsx
- .ValidateKey
- UsageConsole.tsx
- repository
- server_test.go
- teeHandler
- ProxyPoolsConsole.tsx
- Module
- scripts
- pricingStore
- modelIntelligence.ts
- TokenSaverConsole.tsx
- .AuthenticateAdmin
- Module
- newTestHarness
- go_pkg_encoding_json
- repository
- Module
- sync.Mutex
- Load
- web/package.json
- newTestHarness
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
- UpstreamCredentials
- repository
- SettingsConsole.tsx
- awsQuickRuntime
- .Record
- main.go
- oauthLoopbackProxy
- net/http.ResponseWriter
- postExchangeCodex
- repository
- WriteJSON
- .Reader
- net/url.URL
- .UpdateCombo
- post
- ref_bun_test
- headroomCache
- oauthProviderStrategy
- ProviderTopology.tsx
- repo
- Module
- repo
- properties
- setupModelEditorTest
- scanner.go
- Module
- modeleditor/engine.go
- time.Duration
- .UpsertPlugin
- properties
- pxpipe/repo.go
- ModelCatalogConsole.tsx
- ModelDescriptor
- stubModule
- NewClient
- OMEG4MON Router
- stubDiscoverySystemModule
- orderStubModule
- certManager
- repo
- applyUpstreamHeaders
- Module
- properties
- release-binaries
- steps
- evaluateAntigravityQuota
- .ExportCombos
- stubSystemModule
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- 200
- io/fs.FS
- properties
- cli/package.json
- newTestHarness
- setupNotificationsTest
- get
- scripts
- post
- schemas
- .executeCandidateAttempt
- PreviewPluginsResponse
- time.Time
- stubSystemModule
- vite.config.ts
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- net/http.ResponseWriter
- Review Focus
- Review Focus
- ComboStep
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- scanCombo
- version
- Module
- NewStaticUIHandler
- parameters
- stubCompleter
- New
- nameMinHeap
- ClearImporterHistoryResponse
- source_path
- Review Focus
- tgTestEnv
- docker-release
- DB
- `@omeg4mon/router` CLI Launcher
- RouteMounts
- setupModelEditorTest
- context.Context
- listSettingsWithQueryer
- stubDep
- parameters
- catalogFakeProviderStore
- OverviewRealtimeTraffic.tsx
- omeg4mon-router Verification Gate
- ModelTestResult
- .resolveIdentifier
- single-binary
- keyLimiter
- verify-web
- deviceOAuthSession
- smoke-test.sh
- .InsertRun
- readLegacySQLite
- responseCoalescer
- plugins/interceptor.go
- NewChromeTLSClient
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- clsx
- setupMediaTest
- repo
- settingDefinition
- analyzeIdentityProbeResponse
- .CompleteChat
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- net/http/httptest.ResponseRecorder
- batcher
- mapBusEventToAudit
- repository
- EventListResponse
- TestContractInterfacesSatisfied
- jsonMarshalerError
- setupTunnelTest
- compressor.go
- pricingListResponse
- stubDepModule
- omeg4mon-router Project Standards
- streamCache.ts
- stubSettings
- .prepareUpstreamBody

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

## Communities (259 total, 47 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.12
Nodes (35): ModelOption, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, Module, ChangedEvent (+27 more)

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
Cohesion: 0.07
Nodes (24): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository() (+16 more)

### Community 6 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.17
Nodes (21): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+13 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "CombosConsole.tsx"
Cohesion: 0.12
Nodes (36): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+28 more)

### Community 10 - "Button"
Cohesion: 0.07
Nodes (79): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Functional & UX Scope, class-variance-authority (+71 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (77): testEnv, TestSentinelErrorsDistinct(), setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies() (+69 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (17): Module, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), anthropicErrorType() (+9 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (37): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+29 more)

### Community 16 - "Module"
Cohesion: 0.14
Nodes (6): Module, defaultRefreshFormToken(), refreshAWSOIDCToken(), refreshClaudeToken(), refreshCall, storedCredentials

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "translate_anthropic.go"
Cohesion: 0.15
Nodes (18): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+10 more)

### Community 21 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "WriteError"
Cohesion: 0.16
Nodes (9): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateDescription(), validateModelIdentifier() (+1 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "ProviderDetailView.tsx"
Cohesion: 0.10
Nodes (41): runModelIdentityProbe(), AddCustomModelModal(), AddCustomModelModalProps, mergeManualModels(), MergeManualModelsResult, normalizeManualModelInput(), removeModelFromList(), stripProviderModelPrefix() (+33 more)

### Community 29 - "TranslatorConsole.tsx"
Cohesion: 0.15
Nodes (16): translatorQueryKeys, useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS, translateFormSchema, TranslateFormValues (+8 more)

### Community 30 - "navigation.ts"
Cohesion: 0.11
Nodes (25): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, RawModuleEntry, UILocale, ActiveNavContext (+17 more)

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
Cohesion: 0.19
Nodes (5): memoryResponseWriter, copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

### Community 38 - "CanonicalUsage"
Cohesion: 0.10
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "New"
Cohesion: 0.09
Nodes (22): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), Server, New() (+14 more)

### Community 40 - "KeysConsole.tsx"
Cohesion: 0.06
Nodes (71): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt(), KeyExpiryState (+63 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.50
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, mapOpenAIFinishReasonToCanonical()

### Community 42 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 43 - "ir.go"
Cohesion: 0.11
Nodes (48): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat (+40 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "buildApp"
Cohesion: 0.08
Nodes (32): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog(), TestWave2BDisabledModulesSkipInterceptors() (+24 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "New"
Cohesion: 0.06
Nodes (51): roundTripFunc, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, stubCombos, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+43 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 50 - "Module"
Cohesion: 0.15
Nodes (6): AdminCredentialStore, clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.16
Nodes (23): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+15 more)

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.23
Nodes (11): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+3 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.04
Nodes (93): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+85 more)

### Community 57 - "dependencies"
Cohesion: 0.09
Nodes (22): @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, dependencies (+14 more)

### Community 58 - "oauthTokenResponse"
Cohesion: 0.23
Nodes (15): postExchangeIFlow(), postExchangeNous(), enrichAntigravity(), enrichGeminiOAuth(), enrichGitHubCopilot(), isHostMatch(), postExchangeAntigravity(), postExchangeGeminiOAuth() (+7 more)

### Community 59 - "main_test.go"
Cohesion: 0.07
Nodes (22): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+14 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (156): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, ADAPTER_TO_REGISTRY_ID, buildProviderEntries() (+148 more)

### Community 61 - ".ValidateKey"
Cohesion: 0.22
Nodes (5): KeyInfo, KeyValidator, stubKeyValidator, stubKeyValidator, stubKeyValidator

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (53): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), @hookform/resolvers, useCombos(), useProviderConnections(), useSystemTelemetry(), formatBytesGB(), formatHostUptimeShort(), LiveClockUptime() (+45 more)

### Community 63 - "repository"
Cohesion: 0.10
Nodes (21): cloneStringMap(), cloneStringSlice(), cloneUpstreamCredentials(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx() (+13 more)

### Community 64 - "server_test.go"
Cohesion: 0.13
Nodes (11): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, testAllModelsRequest, testComboRequest (+3 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (33): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 67 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.07
Nodes (7): AuthVerifier, stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "newTestHarness"
Cohesion: 0.07
Nodes (29): File Map, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota() (+21 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (22): providerIDContextKey, buildDiscoveryCandidateURLs(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor (+14 more)

### Community 76 - "repository"
Cohesion: 0.12
Nodes (20): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+12 more)

### Community 77 - "Module"
Cohesion: 0.19
Nodes (6): decodeJSONBody(), isValidMode(), isValidSensitivity(), newID(), validateRule(), Module

### Community 78 - "sync.Mutex"
Cohesion: 0.17
Nodes (8): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 79 - "Load"
Cohesion: 0.18
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 81 - "newTestHarness"
Cohesion: 0.16
Nodes (14): fakeClock, stubAdminAuth, testHarness, decodeJSON(), newFakeClock(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency() (+6 more)

### Community 82 - "Module"
Cohesion: 0.19
Nodes (10): extractUsageCost(), actorFromContext(), decodeJSONBody(), New(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID() (+2 more)

### Community 83 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 86 - "Module"
Cohesion: 0.05
Nodes (22): streamResponseBody(), decodeStrictJSON(), New(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), newID(), newRepo() (+14 more)

### Community 87 - "arena/hooks.ts"
Cohesion: 0.07
Nodes (41): ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+33 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.14
Nodes (11): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+3 more)

### Community 90 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 91 - "strategyState"
Cohesion: 0.18
Nodes (10): ProviderStrategyState, strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), fastRandFloat64(), fastRandIntN() (+2 more)

### Community 92 - "UpstreamCredentials"
Cohesion: 0.20
Nodes (17): UpstreamCredentials, buildModelProbeRequest(), extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload() (+9 more)

### Community 93 - "repository"
Cohesion: 0.17
Nodes (10): boolToInt(), marshalOverrideJSON(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule, ModelOverride (+2 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (42): zod, zod@^4.6.7, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation() (+34 more)

### Community 95 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 96 - ".Record"
Cohesion: 0.14
Nodes (5): stubUsageRecorder, stubUsageRecorder, UsageRecorder, stubUsageRecorder, recordingUsage

### Community 97 - "main.go"
Cohesion: 0.11
Nodes (11): targetKey, ChatCompletionRequest, ChatMessage, errorDisposition, handledGatewayError, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus() (+3 more)

### Community 98 - "oauthLoopbackProxy"
Cohesion: 0.19
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (34): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol(), validateHTTPURL() (+26 more)

### Community 100 - "postExchangeCodex"
Cohesion: 0.22
Nodes (10): generateTraceID(), resolveAwsQuickTenantBaseURL(), decodeJWTPayload(), extractCodexAccountInfo(), extractCompanionProjectID(), extractEmailFromJWT(), oauthPlatformEnum(), postExchangeCodex() (+2 more)

### Community 101 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 102 - "WriteJSON"
Cohesion: 0.08
Nodes (15): Module, Module, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), providerDisplayName() (+7 more)

### Community 103 - ".Reader"
Cohesion: 0.04
Nodes (55): main(), runMain(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints, Review Focus (+47 more)

### Community 104 - "net/url.URL"
Cohesion: 0.10
Nodes (14): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+6 more)

### Community 105 - ".UpdateCombo"
Cohesion: 0.29
Nodes (9): createComboParams, updateComboValues, cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps(), encodeConfigJSON() (+1 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "ref_bun_test"
Cohesion: 0.06
Nodes (30): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+22 more)

### Community 108 - "headroomCache"
Cohesion: 0.24
Nodes (3): newHeadroomCache(), headroomCache, WarmPrefixSummary

### Community 109 - "oauthProviderStrategy"
Cohesion: 0.24
Nodes (8): isAnthropicTokenURL(), defaultExchangePKCECode(), defaultOAuthStrategies(), exchangeClaudeCode(), newRFC8628Strategy(), oauthProviderStrategy, oauthSessionRecord, testEnv

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (18): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), RouterHubNode() (+10 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "Module"
Cohesion: 0.08
Nodes (12): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), xrouterBuiltInOAuthAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable() (+4 more)

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (8): Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 116 - "scanner.go"
Cohesion: 0.33
Nodes (9): parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL, PoolTestResponse, ScanCandidateResult (+1 more)

### Community 118 - "Module"
Cohesion: 0.05
Nodes (38): Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), setupSecurityTest(), TestSecurityAuditPaginationOffset(), TestSecurityEventBusAutoAuditIngestion() (+30 more)

### Community 119 - "modeleditor/engine.go"
Cohesion: 0.29
Nodes (9): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 120 - "time.Duration"
Cohesion: 0.09
Nodes (12): AccountHealthTracker, stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore() (+4 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "pxpipe/repo.go"
Cohesion: 0.39
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 124 - "ModelCatalogConsole.tsx"
Cohesion: 0.08
Nodes (62): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update, AssignKeysModalProps, classifyTestError(), clearCachedTestResults() (+54 more)

### Community 125 - "ModelDescriptor"
Cohesion: 0.18
Nodes (6): stubModelLister, ModelDescriptor, ModelLister, stubModelLister, stubModelLister, stubModelLister

### Community 127 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.24
Nodes (6): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, storedConfig

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
Cohesion: 0.38
Nodes (3): AliasSnapshot, ComboSnapshot, ComboSnapshotter

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
Cohesion: 0.03
Nodes (10): stubModule, migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), stubSystemModule, stubModule, stubSystem (+2 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "newTestHarness"
Cohesion: 0.14
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 154 - "setupNotificationsTest"
Cohesion: 0.15
Nodes (16): TestAppCloseTeardownOrderIsLIFO(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup() (+8 more)

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
Cohesion: 0.20
Nodes (10): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), buildCompactedResponse(), randomResponseID(), isTransientSingleCandidateStatus() (+2 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "time.Time"
Cohesion: 0.13
Nodes (11): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, newLoginRateLimiter(), newCooldownTracker(), newFakeClock() (+3 more)

### Community 163 - "vite.config.ts"
Cohesion: 0.29
Nodes (6): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (10): trackedResponseWriter, decodeJSONBody(), New(), decodeJSONBody(), decodePathModelParam(), Module, isValidRate(), jsonMethodNotAllowedHandler() (+2 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "ComboStep"
Cohesion: 0.47
Nodes (5): ComboStep, validateComboTargets(), decodeStepsJSON(), normalizeAndValidateSteps(), stepsFromTargets()

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

### Community 180 - "stubCompleter"
Cohesion: 0.33
Nodes (3): stubCompleter, ChatCompleter, ChatCompletionResult

### Community 181 - "New"
Cohesion: 0.07
Nodes (40): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+32 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "Review Focus"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 186 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "DB"
Cohesion: 0.07
Nodes (44): Options, Options, Options, Options, New(), newRepo(), New(), newRepository() (+36 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "RouteMounts"
Cohesion: 0.09
Nodes (9): Module, validateSQLiteSource(), decodeJSON(), actorFromContext(), actorFromContext(), decodeJSONBody(), Module, AdminSubject() (+1 more)

### Community 191 - "setupModelEditorTest"
Cohesion: 0.14
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins() (+11 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (13): authHarness, fakeCredentialStore, fakeProviderStore, stubAdminCredentialStore, stubProviderStore, stubProviderStore, stubResolver, newFakeProviderStore() (+5 more)

### Community 194 - "listSettingsWithQueryer"
Cohesion: 0.24
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 198 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.25
Nodes (14): formatActivityOriginLabel(), LazyProviderTopology, OverviewRealtimeTraffic(), OverviewRealtimeTrafficProps, RecentRequestsRows(), ProviderTopologyProps, FilteredRecentRequestsResult, filterRecentRequests() (+6 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "ModelTestResult"
Cohesion: 0.40
Nodes (3): ModelTestResult, testAllModelsResponse, resolveKeylessProbeCredentials()

### Community 201 - ".resolveIdentifier"
Cohesion: 0.20
Nodes (6): resolveSnapshot, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), isModelAllowed(), isModelDisabled()

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.26
Nodes (5): newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "deviceOAuthSession"
Cohesion: 0.18
Nodes (6): newDeviceOAuthRuntime(), originOfURL(), deviceOAuthResult, deviceOAuthRuntime, deviceOAuthSession, rfc8628DeviceSpec

### Community 207 - ".InsertRun"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 209 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 210 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 211 - "plugins/interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 212 - "NewChromeTLSClient"
Cohesion: 0.40
Nodes (4): bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton()

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "clsx"
Cohesion: 0.33
Nodes (4): clsx, zustand, clsx@^2.1.1, zustand@^5.1.0

### Community 215 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 216 - "repo"
Cohesion: 0.18
Nodes (12): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), newEventID(), normalizeSeverity() (+4 more)

### Community 217 - "settingDefinition"
Cohesion: 0.14
Nodes (11): HistoryItem, SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), historyListResp, historyResponse, settingDefinition (+3 more)

### Community 218 - "analyzeIdentityProbeResponse"
Cohesion: 0.22
Nodes (10): ModelIdentityResult, analyzeIdentityProbeResponse(), extractIDPrefix(), extractVendorHeaders(), inferModelFamily(), modelsLooselyMatch(), normalizeModelForCompare(), parseIdentityPayload() (+2 more)

### Community 219 - ".CompleteChat"
Cohesion: 0.50
Nodes (3): extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.18
Nodes (21): react-hook-form, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES (+13 more)

### Community 223 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.16
Nodes (19): Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+11 more)

### Community 224 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

### Community 225 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 229 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 231 - "TestContractInterfacesSatisfied"
Cohesion: 0.05
Nodes (10): CanonicalRequest, GatewayInterceptor, ProviderResolver, ResolvedTarget, ResponseMeta, StreamEvent, stubInterceptor, stubProviderResolver (+2 more)

### Community 241 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 242 - "compressor.go"
Cohesion: 0.12
Nodes (28): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+20 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 254 - "streamCache.ts"
Cohesion: 0.43
Nodes (4): UsageLogsResponse, matchesLogFilter(), prependUsageLogToResponse(), StreamLogsFilter

### Community 256 - ".prepareUpstreamBody"
Cohesion: 0.22
Nodes (6): ProtocolTranslator, CanFastPath(), IsSupportedProtocol(), RewriteModelField(), TestProtocolConstantsAndFastPath(), TestRewriteModelFieldPreservesProviderSpecificFields()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1026 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **47 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `net/http.ResponseWriter`?**
  _High betweenness centrality (0.195) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11770768913626056 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `net/http.ResponseWriter` to `context.Context`, `go_pkg_encoding_json`, `ProvidersModals.tsx`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.168) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `newTestHarness` to `.Reader`, `New`, `Button`, `Module`, `net/http/httptest.ResponseRecorder`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._