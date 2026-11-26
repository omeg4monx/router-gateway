# Graph Report - omeg4mon-router  (2026-11-26)

## Corpus Check
- 449 files · ~485,654 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 10 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5310 nodes · 17212 edges · 246 communities (201 shown, 45 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1565 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b6390e8b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main_test.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- promptinjection/engine.go
- LiveMonitorConsole.tsx
- TypeScript Strict Compiler Config
- CombosConsole.tsx
- Button
- testing.T
- shadcn/ui Component Configuration
- fakeProviderStore
- MediaConsole.tsx
- App.tsx
- Module
- Design System & Wave 0 Plan Docs
- TelegramConsole.tsx
- updates
- translate_anthropic.go
- PluginsStudio.tsx
- with
- jobs
- WriteError
- steps
- web-check
- on
- net/http.ResponseWriter
- Server
- Vault
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- customModelUtils.ts
- media/module_test.go
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
- net/http.Request
- translate_openai.go
- Module
- ModelsRoutingConsole.tsx
- translate_responses.go
- go_pkg_context
- SystemOverview.tsx
- extractErrorMessage
- dependencies
- Module
- server_test.go
- ProvidersModals.tsx
- Open
- UsageConsole.tsx
- repository
- strategyState
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
- card.tsx
- Module
- NotificationsConsole.tsx
- Load
- web/package.json
- NewStreamTranslator
- Module
- newTestHarness
- paths
- collectDynamicHostTelemetry
- fakeHandle
- repo
- settingDefinition
- repo
- Module
- .UpdateCombo
- pricingStore
- promptinjection/engine.go
- SettingsConsole.tsx
- oauthLoopbackProxy
- Runner
- main.go
- Config
- Module
- setupCatalogTestModule
- RingBuffer
- repository
- DB
- certs.go
- ref_node_fs
- post
- ref_node_path
- WriteJSON
- net/http.Header
- ProviderTopology.tsx
- repo
- Module
- repo
- properties
- .Reader
- net/url.URL
- deviceOAuthSession
- Module
- time.Time
- .UpsertPlugin
- properties
- time.Duration
- contracts.UpstreamCredentials
- setupTestMitm
- io/fs.FS
- Module
- OMEG4MON Router
- ProcessHandle
- .Record
- certManager
- repo
- applyUpstreamHeaders
- Module
- properties
- release-binaries
- steps
- evaluateAntigravityQuota
- .ExportCombos
- parseJSONObjectOrWrap
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- responseCoalescer
- properties
- 200
- Module
- properties
- cli/package.json
- setupMediaTest
- setupSecurityTest
- get
- scripts
- post
- schemas
- Review Focus
- PreviewPluginsResponse
- .ProvisionKey
- setupTunnelTest
- telemetryCollector
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- stubDep
- Review Focus
- Review Focus
- io.Reader
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- mobileLayout.test.ts
- version
- Review Focus
- Review Focus
- parameters
- Review Focus
- newHarness
- mapBusEventToAudit
- ClearImporterHistoryResponse
- source_path
- botAPIClient
- docker-release
- `@omeg4mon/router` CLI Launcher
- RouteMounts
- setupModelEditorTest
- context.Context
- listSettingsWithQueryer
- stubDep
- parameters
- schema.go
- omeg4mon-router Verification Gate
- stubDepModule
- single-binary
- keyLimiter
- verify-web
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- readLegacySQLite
- stubUsageRecorder
- plugins/interceptor.go
- net/http.Client
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- newTestHarness
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- repository
- repository
- time.Duration
- App.test.ts
- jsonMarshalerError
- navigation.ts
- compressor.go
- omeg4mon-router Project Standards
- CanonicalChatResponse

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 190 edges
2. `WriteJSON()` - 184 edges
3. `WriteError()` - 182 edges
4. `DB` - 93 edges
5. `Button` - 86 edges
6. `buildApp()` - 79 edges
7. `Bus` - 72 edges
8. `Alert` - 69 edges
9. `AlertDescription` - 68 edges
10. `cn()` - 67 edges

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
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (246 total, 45 thin omitted)

### Community 0 - "main_test.go"
Cohesion: 0.09
Nodes (24): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+16 more)

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
Cohesion: 0.07
Nodes (24): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository() (+16 more)

### Community 6 - "promptinjection/engine.go"
Cohesion: 0.10
Nodes (21): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+13 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.09
Nodes (41): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+33 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "CombosConsole.tsx"
Cohesion: 0.10
Nodes (42): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+34 more)

### Community 10 - "Button"
Cohesion: 0.06
Nodes (85): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Functional & UX Scope, class-variance-authority (+77 more)

### Community 11 - "testing.T"
Cohesion: 0.05
Nodes (72): Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, testEnv, TestSentinelErrorsDistinct(), TestResolverKeylessVeniceFallbackWithoutDBConnections(), TestAntigravityBare429PersistentRefusalDoesNotCooldownAccount(), TestAntigravityBare429RetriesWithSystemMovedToUser(), TestAntigravityHeaderStripsConsumerProjectID(), TestAntigravityProactiveQuotaPrecheckSkipsDrainedAccount() (+64 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "fakeProviderStore"
Cohesion: 0.11
Nodes (5): fakeProviderStore, stubProviderStore, stubProviderStore, newFakeProviderStore(), stubProviderStore

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (36): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+28 more)

### Community 16 - "Module"
Cohesion: 0.12
Nodes (33): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies() (+25 more)

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
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 21 - "PluginsStudio.tsx"
Cohesion: 0.19
Nodes (22): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+14 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "WriteError"
Cohesion: 0.16
Nodes (13): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+5 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (16): Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), Module, trackedResponseWriter, Module, providerDisplayName(), splitTargetIdentifier(), anthropicErrorType(), geminiErrorStatus() (+8 more)

### Community 29 - "Server"
Cohesion: 0.07
Nodes (18): stubAdminAuth, testHarness, decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering() (+10 more)

### Community 30 - "Vault"
Cohesion: 0.13
Nodes (13): File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`) (+5 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "customModelUtils.ts"
Cohesion: 0.48
Nodes (5): mergeManualModels(), MergeManualModelsResult, normalizeManualModelInput(), removeModelFromList(), stripProviderModelPrefix()

### Community 35 - "media/module_test.go"
Cohesion: 0.28
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "New"
Cohesion: 0.11
Nodes (25): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+17 more)

### Community 40 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (67): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt(), KeyExpiryState (+59 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.38
Nodes (4): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder, newOpenAIStreamDecoder()

### Community 42 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 43 - "ir.go"
Cohesion: 0.13
Nodes (37): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat (+29 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "buildApp"
Cohesion: 0.06
Nodes (40): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting(), TestWave2AE2ETranslatorAndCatalog() (+32 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "net/http.Request"
Cohesion: 0.05
Nodes (49): Global Constraints, Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Task 5: Full Verification & Graphify Sync, Approaches Considered, clearSessionCookie(), isSecureRequest() (+41 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (26): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+18 more)

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.24
Nodes (10): filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY, OverviewCopy (+2 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.05
Nodes (78): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+70 more)

### Community 57 - "dependencies"
Cohesion: 0.07
Nodes (27): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+19 more)

### Community 58 - "Module"
Cohesion: 0.25
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 59 - "server_test.go"
Cohesion: 0.08
Nodes (47): ModelOption, PolicyChangedEvent, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse (+39 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (172): Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, runModelIdentityProbe(), AddCustomModelModal(), ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), canAddProviderConnection() (+164 more)

### Community 61 - "Open"
Cohesion: 0.05
Nodes (53): roundTripFunc, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation() (+45 more)

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.05
Nodes (75): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), useCombos(), useSystemTelemetry(), formatActivityOriginLabel(), formatClockUptime(), LazyProviderTopology, OverviewRealtimeTraffic(), OverviewRealtimeTrafficProps (+67 more)

### Community 63 - "repository"
Cohesion: 0.10
Nodes (21): cloneStringMap(), cloneStringSlice(), cloneUpstreamCredentials(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx() (+13 more)

### Community 64 - "strategyState"
Cohesion: 0.11
Nodes (15): resolveSnapshot, strategyState, weightedCandidate, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit(), applyComboStrategy() (+7 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (34): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+26 more)

### Community 67 - "Module"
Cohesion: 0.10
Nodes (10): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode() (+2 more)

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
Cohesion: 0.07
Nodes (39): isHealthResponse(), systemKeys, tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId (+31 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.07
Nodes (8): AuthVerifier, stubAuthVerifier, stubAuthVerifier, Module, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "Module"
Cohesion: 0.12
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 74 - "newTestHarness"
Cohesion: 0.08
Nodes (28): File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex() (+20 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.11
Nodes (18): providerIDContextKey, buildDiscoveryCandidateURLs(), TestAutoDetectModelsDiscoveryCandidates(), catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, createConnectionRequest, createNodeRequest (+10 more)

### Community 76 - "card.tsx"
Cohesion: 0.06
Nodes (65): ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+57 more)

### Community 77 - "Module"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 78 - "NotificationsConsole.tsx"
Cohesion: 0.16
Nodes (24): zod, notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 79 - "Load"
Cohesion: 0.20
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.10
Nodes (27): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+19 more)

### Community 81 - "NewStreamTranslator"
Cohesion: 0.22
Nodes (9): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs() (+1 more)

### Community 82 - "Module"
Cohesion: 0.19
Nodes (10): extractUsageCost(), actorFromContext(), decodeJSONBody(), New(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID() (+2 more)

### Community 83 - "newTestHarness"
Cohesion: 0.16
Nodes (17): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+9 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.16
Nodes (19): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+11 more)

### Community 86 - "fakeHandle"
Cohesion: 0.13
Nodes (3): fakeHandle, osProcessHandle, ssePipeReader

### Community 87 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.09
Nodes (17): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), buildWhere() (+9 more)

### Community 90 - "Module"
Cohesion: 0.15
Nodes (6): AdminCredentialStore, clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 91 - ".UpdateCombo"
Cohesion: 0.15
Nodes (17): Combo, ComboStep, createComboParams, rowScanner, updateComboValues, cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+9 more)

### Community 93 - "promptinjection/engine.go"
Cohesion: 0.08
Nodes (26): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt() (+18 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.10
Nodes (40): AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation(), useRollbackSettingMutation() (+32 more)

### Community 95 - "oauthLoopbackProxy"
Cohesion: 0.09
Nodes (12): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo (+4 more)

### Community 96 - "Runner"
Cohesion: 0.16
Nodes (8): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "main.go"
Cohesion: 0.09
Nodes (26): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+18 more)

### Community 98 - "Config"
Cohesion: 0.09
Nodes (23): Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), setupSecurityTest() (+15 more)

### Community 99 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 100 - "setupCatalogTestModule"
Cohesion: 0.12
Nodes (19): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies(), TestAllOmniRouteComboStrategiesCRUDAndPersistence() (+11 more)

### Community 101 - "RingBuffer"
Cohesion: 0.21
Nodes (10): File Map, Review Focus, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProxyDialerResolver, cloneAttrMap(), RingBuffer (+2 more)

### Community 102 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 103 - "DB"
Cohesion: 0.05
Nodes (49): Options, Options, Module, Options, repository, New(), newRepo(), validateSQLiteSource() (+41 more)

### Community 104 - "certs.go"
Cohesion: 0.14
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 105 - "ref_node_fs"
Cohesion: 0.12
Nodes (10): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR, modalsSource, providersDir, consoleSource, detailSource, modalsSource (+2 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "ref_node_path"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 108 - "WriteJSON"
Cohesion: 0.09
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON(), LogEntry (+3 more)

### Community 109 - "net/http.Header"
Cohesion: 0.13
Nodes (9): memoryResponseWriter, extractVendorHeaders(), extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration() (+1 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (19): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+11 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "Module"
Cohesion: 0.08
Nodes (11): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), defaultBuiltInAdapters(), oauthEnvPrefix(), xrouterBuiltInOAuthAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable(), builtinAdapter (+3 more)

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - ".Reader"
Cohesion: 0.09
Nodes (23): main(), runMain(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+15 more)

### Community 116 - "net/url.URL"
Cohesion: 0.13
Nodes (12): stubProxyDialerResolver, stubProxyResolver, copyNonHopHeaders(), matchToolForHost(), newID(), singleJoiningSlash(), parseAndValidateProxyURL(), probeProxy() (+4 more)

### Community 118 - "deviceOAuthSession"
Cohesion: 0.23
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 119 - "Module"
Cohesion: 0.10
Nodes (8): firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestParseDiscoveredModelsFormats(), Module, discoveredModel, fetchConnectionModelsResponse, refreshCall

### Community 120 - "time.Time"
Cohesion: 0.08
Nodes (16): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock() (+8 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "time.Duration"
Cohesion: 0.08
Nodes (15): stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore(), TestTrafficShaperAdaptiveHeaderPacing() (+7 more)

### Community 124 - "contracts.UpstreamCredentials"
Cohesion: 0.07
Nodes (41): catalogFakeProviderStore, ModelIdentityResult, ModelTestResult, testAllModelsResponse, buildModelProbeRequest(), extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials() (+33 more)

### Community 125 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.02
Nodes (20): stubModule, migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), applyStaticCacheHeaders(), hasUnsafePath() (+12 more)

### Community 127 - "Module"
Cohesion: 0.15
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 130 - ".Record"
Cohesion: 0.18
Nodes (3): stubUsageRecorder, stubUsageRecorder, recordingUsage

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "applyUpstreamHeaders"
Cohesion: 0.07
Nodes (24): Task 2: Backend Model Test Endpoint (TDD), Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, attemptOutcome (+16 more)

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
Cohesion: 0.17
Nodes (15): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+7 more)

### Community 140 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 141 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 142 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "responseCoalescer"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "setupMediaTest"
Cohesion: 0.36
Nodes (6): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction()

### Community 154 - "setupSecurityTest"
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

### Community 159 - "Review Focus"
Cohesion: 0.29
Nodes (7): Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 162 - "setupTunnelTest"
Cohesion: 0.25
Nodes (9): New(), doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary() (+1 more)

### Community 163 - "telemetryCollector"
Cohesion: 0.48
Nodes (5): collectDynamicHostTelemetry(), TestTelemetryCacheAvoidsRedundantWork(), hostTelemetrySnapshot, staticHostInfo, telemetryCollector

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

### Community 175 - "mobileLayout.test.ts"
Cohesion: 0.40
Nodes (4): consoleSource, detailSource, modalsSource, providersDir

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "Review Focus"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 178 - "Review Focus"
Cohesion: 0.50
Nodes (4): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Review Focus"
Cohesion: 0.50
Nodes (4): Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI), Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI)

### Community 181 - "newHarness"
Cohesion: 0.11
Nodes (20): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Options, newMockVeniceServer(), newHarness(), startAuthHTTPProxy(), TestProxyPoolsBindingLifecycleAndCacheInvalidation(), TestProxyPoolsCRUDCredentialRedactionAndValidation(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe() (+12 more)

### Community 182 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "botAPIClient"
Cohesion: 0.43
Nodes (3): newBotAPIClient(), sanitizeBotError(), botAPIClient

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "RouteMounts"
Cohesion: 0.21
Nodes (3): actorFromContext(), decodeJSONBody(), Module

### Community 191 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (19): arenaTestEnv, stubCompleter, authHarness, fakeCredentialStore, ChatCompletionResult, SettingsSnapshotter, stubAdminCredentialStore, ImportRun (+11 more)

### Community 194 - "listSettingsWithQueryer"
Cohesion: 0.11
Nodes (15): listSettingsWithQueryer(), decodeStoredValue(), SettingItem, int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds(), applyResult (+7 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "schema.go"
Cohesion: 0.33
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.26
Nodes (5): newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 209 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 211 - "plugins/interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 212 - "net/http.Client"
Cohesion: 0.15
Nodes (14): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar(), isTransientMintStatus(), ResolveMintURL() (+6 more)

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 223 - "newTestHarness"
Cohesion: 0.08
Nodes (34): Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+26 more)

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 230 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 231 - "time.Duration"
Cohesion: 0.04
Nodes (15): stubModelLister, ResponseMeta, StreamEvent, stubInterceptor, stubKeyValidator, stubModelLister, stubProviderResolver, Options (+7 more)

### Community 235 - "App.test.ts"
Cohesion: 0.40
Nodes (4): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ALL_NAV_ITEMS

### Community 241 - "navigation.ts"
Cohesion: 0.14
Nodes (21): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, RawModuleEntry, UILocale, ActiveNavContext, EXTENSION_GROUPS, findActiveNavContext(), isModuleRegistered() (+13 more)

### Community 242 - "compressor.go"
Cohesion: 0.12
Nodes (28): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+20 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 256 - "CanonicalChatResponse"
Cohesion: 0.24
Nodes (13): CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), IsSupportedProtocol(), RewriteModelField(), TestAnthropicToOpenAIChatRoundTrip(), TestDecodeRejectsNonObjectJSONPayloads() (+5 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1024 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `net/http.Request` to `ProvidersModals.tsx`?**
  _High betweenness centrality (0.197) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08509615384615385 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `net/http.Request` to `Module`, `context.Context`, `go_pkg_encoding_json`, `contracts.UpstreamCredentials`?**
  _High betweenness centrality (0.165) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `ProvidersModals.tsx` to `Button`, `net/http.Request`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._