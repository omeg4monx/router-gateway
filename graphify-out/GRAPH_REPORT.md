# Graph Report - omeg4mon-router  (2026-11-27)

## Corpus Check
- 460 files · ~494,070 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5347 nodes · 17772 edges · 224 communities (205 shown, 19 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1733 edges (avg confidence: 0.86)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `4995ea9e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main.go
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
- net/http.Request
- MitmConsole.tsx
- ResolvedTarget
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- Card
- repository
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
- Module
- translate_openai.go
- ModelDescriptor
- Module
- ModelsRoutingConsole.tsx
- translate_responses.go
- oauthLoopbackProxy
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
- repository
- teeHandler
- ProxyPoolsConsole.tsx
- ModelDescriptor
- scripts
- pricingStore
- modelIntelligence.ts
- settingDefinition
- .AuthenticateAdmin
- Module
- newTestHarness
- go_pkg_encoding_json
- arena/hooks.ts
- Module
- NotificationsConsole.tsx
- Load
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- newTestHarness
- paths
- collectDynamicHostTelemetry
- buildModelProbeRequest
- repo
- settingDefinition
- repo
- AppShell
- strategies.go
- batcher
- redact.go
- SettingsConsole.tsx
- Module
- Runner
- contracts.go
- Load
- Module
- setupTokenSaverTest
- WriteError
- setupCatalogTestModule
- DB
- go_pkg_time
- palette.test.ts
- post
- selectProbeHTTPClient
- WriteJSON
- net/http.Header
- ProviderTopology.tsx
- repo
- stubCompleter
- repo
- properties
- GatewayInterceptor
- net/url.URL
- stubModule
- deviceOAuthSession
- Module
- time.Time
- .UpsertPlugin
- properties
- time.Duration
- .identifyProviderModel
- stubModule
- io/fs.FS
- setupMediaTest
- OMEG4MON Router
- stubModule
- UsageRecord
- certManager
- repo
- .executeCandidateAttempt
- Module
- properties
- buildModelProbeRequest
- release-binaries
- steps
- evaluateAntigravityQuota
- NewStreamTranslator
- stubModule
- HeadroomStats
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- parseJSONObjectOrWrap
- properties
- 200
- repository
- properties
- repository
- cli/package.json
- New
- newAuthenticatedDiscoveryRequest
- NewRegistry
- get
- scripts
- post
- schemas
- newOAuthLoopbackProxy
- PreviewPluginsResponse
- Module
- setupTunnelTest
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- stubDep
- Review Focus
- Review Focus
- NewStreamTranslator
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- stubSystem
- version
- redact.go
- stubModule
- parameters
- .list
- ClearImporterHistoryResponse
- source_path
- parseIdentityPayload
- docker-release
- `@omeg4mon/router` CLI Launcher
- Module
- setupModelEditorTest
- context.Context
- listSettingsWithQueryer
- parameters
- Registry
- omeg4mon-router Verification Gate
- palette.test.ts
- single-binary
- keyLimiter
- verify-web
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- .list
- ComboSnapshot
- Feature Suite
- Module
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- security/hooks.ts
- net/http/httptest.ResponseRecorder
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- headroomCache
- time.Duration
- App.test.ts
- jsonMarshalerError
- navigation.ts
- compressor.go
- omeg4mon-router Project Standards
- CanonicalChatResponse
- Module

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 190 edges
2. `WriteJSON()` - 184 edges
3. `WriteError()` - 182 edges
4. `DB` - 93 edges
5. `Button` - 89 edges
6. `buildApp()` - 79 edges
7. `Bus` - 72 edges
8. `Alert` - 69 edges
9. `AlertDescription` - 68 edges
10. `cn()` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3a-teams-vouchers.md → internal/contracts/contracts.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx

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

## Communities (224 total, 19 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.06
Nodes (24): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+16 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.18
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

### Community 6 - "promptinjection/engine.go"
Cohesion: 0.12
Nodes (20): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+12 more)

### Community 7 - "LiveMonitorConsole.tsx"
Cohesion: 0.09
Nodes (41): react-hook-form, KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass() (+33 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "CombosConsole.tsx"
Cohesion: 0.09
Nodes (48): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+40 more)

### Community 10 - "Button"
Cohesion: 0.08
Nodes (62): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, lucide-react, react, react-router, ChatConsole(), STARTER_PROMPTS, useArenaModels() (+54 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (69): testEnv, TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering() (+61 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.13
Nodes (9): catalogFakeProviderStore, fakeProviderStore, stubProviderStore, stubProviderStore, ProviderConnection, UpstreamCredentials, newFakeProviderStore(), cloneUpstreamCredentials() (+1 more)

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

### Community 18 - "MitmConsole.tsx"
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

### Community 28 - "net/http.Request"
Cohesion: 0.06
Nodes (25): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), anthropicErrorType(), extractAPIKey() (+17 more)

### Community 29 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 30 - "ResolvedTarget"
Cohesion: 0.11
Nodes (16): resolveSnapshot, strategyState, weightedCandidate, stubProviderResolver, stubResolver, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare() (+8 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "Card"
Cohesion: 0.08
Nodes (47): ArenaConsole(), useArenaBattles(), useArenaLeaderboard(), useDeleteArenaBattleMutation(), useRunArenaBattleMutation(), useVoteArenaBattleMutation(), useProviderStrategy(), useUpdateProviderStrategyMutation() (+39 more)

### Community 35 - "repository"
Cohesion: 0.17
Nodes (10): boolToInt(), marshalOverrideJSON(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule, ModelOverride (+2 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "New"
Cohesion: 0.10
Nodes (26): loadUI(), loadUI(), doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth(), TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping(), apiTimeoutMiddleware() (+18 more)

### Community 40 - "KeysConsole.tsx"
Cohesion: 0.09
Nodes (52): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt(), KeyExpiryState (+44 more)

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
Nodes (46): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+38 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "Module"
Cohesion: 0.05
Nodes (32): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), buildAwsQuickSigninURL() (+24 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "ModelDescriptor"
Cohesion: 0.18
Nodes (12): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, vite, @vitejs/plugin-react, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin() (+4 more)

### Community 50 - "Module"
Cohesion: 0.05
Nodes (24): Options, Claims, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests (+16 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (26): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+18 more)

### Community 53 - "oauthLoopbackProxy"
Cohesion: 0.09
Nodes (12): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo (+4 more)

### Community 54 - "go_pkg_context"
Cohesion: 0.23
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.23
Nodes (11): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+3 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.10
Nodes (42): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+34 more)

### Community 57 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 58 - "Module"
Cohesion: 0.25
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 59 - "server_test.go"
Cohesion: 0.07
Nodes (55): ModelOption, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse (+47 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (183): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, runModelIdentityProbe(), AddCustomModelModal() (+175 more)

### Community 61 - "Open"
Cohesion: 0.04
Nodes (65): arenaTestEnv, roundTripFunc, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+57 more)

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (75): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), useProvidersCatalog(), useProxyPools(), useSystemTelemetry(), formatActivityOriginLabel(), LazyProviderTopology, OverviewRealtimeTraffic(), OverviewRealtimeTrafficProps (+67 more)

### Community 63 - "repository"
Cohesion: 0.12
Nodes (22): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), oauthEnvPrefix(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice() (+14 more)

### Community 64 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.12
Nodes (33): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation(), useUpsertProxyBindingMutation() (+25 more)

### Community 67 - "ModelDescriptor"
Cohesion: 0.24
Nodes (7): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, isModelAllowed(), isModelDisabled(), stubModelLister

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "settingDefinition"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.10
Nodes (8): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "Module"
Cohesion: 0.13
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 74 - "newTestHarness"
Cohesion: 0.10
Nodes (22): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart() (+14 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (20): providerIDContextKey, firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestParseDiscoveredModelsFormats(), PreviewResult, catalogAdapterDescriptor, catalogOAuthDescriptor (+12 more)

### Community 76 - "arena/hooks.ts"
Cohesion: 0.08
Nodes (37): arenaQueryKeys, computeSessionTelemetry(), filterChatSessions(), formatLatencyMs(), formatTokenCount(), SessionTelemetry, sampleMessages, sampleSessions (+29 more)

### Community 77 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 78 - "NotificationsConsole.tsx"
Cohesion: 0.15
Nodes (24): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+16 more)

### Community 79 - "Load"
Cohesion: 0.22
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.12
Nodes (23): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+15 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.12
Nodes (12): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+4 more)

### Community 83 - "newTestHarness"
Cohesion: 0.31
Nodes (10): Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, decodeJSON(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageStreamForwardsInventoryInvalidations() (+2 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 86 - "buildModelProbeRequest"
Cohesion: 0.13
Nodes (3): fakeHandle, osProcessHandle, ssePipeReader

### Community 87 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.16
Nodes (8): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotificationRow(), EmitPayload, repo

### Community 90 - "AppShell"
Cohesion: 0.28
Nodes (12): Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, PwaInstallButton(), PwaInstallButtonProps, PwaStatusBanner(), PwaStatusBannerProps, BeforeInstallPromptEvent, getPwaCopy(), isStandaloneDisplayMode() (+4 more)

### Community 91 - "strategies.go"
Cohesion: 0.13
Nodes (21): Combo, ComboStep, createComboParams, ProviderStrategyState, rowScanner, updateComboValues, validateComboTargets(), cloneStringSlice() (+13 more)

### Community 92 - "batcher"
Cohesion: 0.10
Nodes (8): newBatcher(), Module, New(), newPricingStore(), newRepository(), batcher, ModelPricing, pricingStore

### Community 93 - "redact.go"
Cohesion: 0.46
Nodes (6): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.07
Nodes (50): @hookform/resolvers, zod, zod@^4.6.7, authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+42 more)

### Community 95 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 96 - "Runner"
Cohesion: 0.15
Nodes (9): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector, tgAPICall (+1 more)

### Community 97 - "contracts.go"
Cohesion: 0.11
Nodes (26): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+18 more)

### Community 98 - "Load"
Cohesion: 0.25
Nodes (9): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable() (+1 more)

### Community 99 - "Module"
Cohesion: 0.11
Nodes (10): decodeStrictJSON(), copyNonHopHeaders(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL() (+2 more)

### Community 100 - "setupTokenSaverTest"
Cohesion: 0.27
Nodes (10): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory(), TestTokenSaverFinalReviewFindings() (+2 more)

### Community 101 - "WriteError"
Cohesion: 0.12
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), validateSQLiteSource(), decodeJSON(), decodeOptionalJSONBody() (+3 more)

### Community 102 - "setupCatalogTestModule"
Cohesion: 0.17
Nodes (9): setupCatalogTestModule(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), BenchmarkResolveStrategies(), TestAllOmniRouteComboStrategiesCRUDAndPersistence(), TestOmniRouteResolverStrategiesExecution() (+1 more)

### Community 103 - "DB"
Cohesion: 0.06
Nodes (46): Options, authHarness, TestAppCloseTeardownOrderIsLIFO(), Options, testHarness, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, New(), newRepo() (+38 more)

### Community 104 - "go_pkg_time"
Cohesion: 0.29
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 105 - "palette.test.ts"
Cohesion: 0.09
Nodes (19): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+11 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "selectProbeHTTPClient"
Cohesion: 0.22
Nodes (9): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar(), SelectHTTPClient(), SharedHTTPClient() (+1 more)

### Community 108 - "WriteJSON"
Cohesion: 0.08
Nodes (10): decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), WriteJSON(), LogEntry, LogListResponse (+2 more)

### Community 109 - "net/http.Header"
Cohesion: 0.11
Nodes (14): antigravityPromptRewriteRule, memoryResponseWriter, extractVendorHeaders(), newMemoryResponseWriter(), parseHeaderInt(), parseResetDuration(), applyOverrideToJSON(), applyPipeline() (+6 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (19): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+11 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "stubCompleter"
Cohesion: 0.33
Nodes (6): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent()

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "GatewayInterceptor"
Cohesion: 0.06
Nodes (49): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway (+41 more)

### Community 116 - "net/url.URL"
Cohesion: 0.13
Nodes (16): stubProxyDialerResolver, stubProxyResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock() (+8 more)

### Community 117 - "stubModule"
Cohesion: 0.24
Nodes (7): buildWhere(), newEventID(), newRepo(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 118 - "deviceOAuthSession"
Cohesion: 0.14
Nodes (7): OAuthConfig, originOfURL(), buildAuthorizeURL(), builtinAdapter, deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 119 - "Module"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 120 - "time.Time"
Cohesion: 0.08
Nodes (18): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock() (+10 more)

### Community 121 - ".UpsertPlugin"
Cohesion: 0.25
Nodes (8): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), PluginSpec, PluginState, repository

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "time.Duration"
Cohesion: 0.07
Nodes (17): stubAccountHealthTracker, attemptOutcome, cachedGatewayResponse, chainStore, failureRecord, inflightCall, responseCoalescer, stubHealthTracker (+9 more)

### Community 124 - ".identifyProviderModel"
Cohesion: 0.09
Nodes (28): ModelIdentityResult, ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+20 more)

### Community 125 - "stubModule"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.03
Nodes (15): stubDep, stubModule, migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile() (+7 more)

### Community 127 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "stubModule"
Cohesion: 0.23
Nodes (3): fakeRunner, osProcessRunner, ProcessHandle

### Community 130 - "UsageRecord"
Cohesion: 0.21
Nodes (6): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, generateID(), recordingUsage

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - ".executeCandidateAttempt"
Cohesion: 0.11
Nodes (18): Task 2: Backend Model Test Endpoint (TDD), Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isAntigravityPromptRejection() (+10 more)

### Community 134 - "Module"
Cohesion: 0.08
Nodes (11): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+3 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "buildModelProbeRequest"
Cohesion: 0.31
Nodes (9): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools(), openCodeEnvFingerprintPrompt() (+1 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 140 - "NewStreamTranslator"
Cohesion: 0.22
Nodes (9): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs() (+1 more)

### Community 141 - "stubModule"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 142 - "HeadroomStats"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 153 - "newAuthenticatedDiscoveryRequest"
Cohesion: 0.40
Nodes (4): isBlockedAntigravityProjectID(), buildDiscoveryCandidateURLs(), newAuthenticatedDiscoveryRequest(), TestAutoDetectModelsDiscoveryCandidates()

### Community 154 - "NewRegistry"
Cohesion: 0.09
Nodes (28): TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification (+20 more)

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

### Community 159 - "newOAuthLoopbackProxy"
Cohesion: 0.50
Nodes (4): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider()

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "Module"
Cohesion: 0.07
Nodes (25): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module (+17 more)

### Community 162 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

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

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "redact.go"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 181 - ".list"
Cohesion: 0.06
Nodes (42): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`) (+34 more)

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

### Community 190 - "Module"
Cohesion: 0.17
Nodes (4): newDeviceOAuthRuntime(), Module, deviceOAuthRuntime, refreshCall

### Community 191 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (12): stubAdminCredentialStore, stubSettings, sanitizeBotError(), Alias, AliasListResponse, Config, Event, repository (+4 more)

### Community 194 - "listSettingsWithQueryer"
Cohesion: 0.24
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 198 - "Registry"
Cohesion: 0.20
Nodes (4): Registry, normalizeName(), nameMinHeap, registryState

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.29
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.22
Nodes (8): Approaches Considered, Functional & UX Scope, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification, useSystemHealth(), useSystemVersion()

### Community 208 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 209 - "ComboSnapshot"
Cohesion: 0.23
Nodes (10): legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables(), readLegacyCombos(), readLegacyKeys() (+2 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 211 - "Module"
Cohesion: 0.16
Nodes (7): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages(), Module

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "security/hooks.ts"
Cohesion: 0.09
Nodes (34): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+26 more)

### Community 223 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.09
Nodes (27): fakeCredentialStore, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+19 more)

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "headroomCache"
Cohesion: 0.14
Nodes (8): newHeadroomCache(), New(), boolToInt(), defaultConfig(), newRepository(), Config, headroomCache, repository

### Community 231 - "time.Duration"
Cohesion: 0.06
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID(), computePrefixFingerprint()

### Community 235 - "App.test.ts"
Cohesion: 0.40
Nodes (4): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ALL_NAV_ITEMS

### Community 241 - "navigation.ts"
Cohesion: 0.13
Nodes (26): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, RawModuleEntry, UILocale, AppShellProps, createInitialExpandedAccordions(), isEditableElement(), ActiveNavContext (+18 more)

### Community 242 - "compressor.go"
Cohesion: 0.13
Nodes (27): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+19 more)

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 256 - "CanonicalChatResponse"
Cohesion: 0.24
Nodes (13): CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), IsSupportedProtocol(), RewriteModelField(), TestAnthropicToOpenAIChatRoundTrip(), TestDecodeRejectsNonObjectJSONPayloads() (+5 more)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1012 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `Module`?**
  _High betweenness centrality (0.170) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06347469220246238 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `go_pkg_encoding_json`, `UpstreamCredentials`, `net/http.Request`, `newAuthenticatedDiscoveryRequest`, `ProvidersModals.tsx`?**
  _High betweenness centrality (0.143) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `Module` to `newTestHarness`, `Button`, `buildApp`, `GatewayInterceptor`, `net/http/httptest.ResponseRecorder`?**
  _High betweenness centrality (0.080) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._