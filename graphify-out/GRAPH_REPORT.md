# Graph Report - omeg4mon-router  (2026-11-25)

## Corpus Check
- 467 files · ~502,382 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5775 nodes · 18895 edges · 248 communities (218 shown, 30 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1781 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `415ddb39`
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
- Module
- TypeScript Strict Compiler Config
- NotificationsConsole.tsx
- Button
- testing.T
- shadcn/ui Component Configuration
- Module
- MediaConsole.tsx
- App.tsx
- catalog.go
- Design System & Wave 0 Plan Docs
- Button
- updates
- ir.go
- time.Duration
- with
- jobs
- WriteError
- steps
- web-check
- on
- io/fs.FS
- Module
- TelegramConsole.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- Button
- navigation.ts
- Go Module Root Definition
- CanonicalUsage
- go_pkg_strings
- UsageConsole.tsx
- encoding/json.RawMessage
- New
- ModelCatalogConsole.tsx
- translate_gemini.go
- Module
- repository
- get
- translate_openai.go
- JobInspector
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- Bus
- main.go
- SystemOverview.tsx
- extractErrorMessage
- dependencies
- oauth_registry.go
- Module
- ProvidersModals.tsx
- VouchersConsole.tsx
- modelIntelligence.ts
- ProviderAdapter
- UpstreamCredentials
- teeHandler
- ProxyPoolsConsole.tsx
- OAuthConfig
- scripts
- pricingStore
- Module
- TokenSaverConsole.tsx
- .AuthenticateAdmin
- Module
- patch
- go_pkg_encoding_json
- promptinjection/engine.go
- Module
- Runner
- Load
- web/package.json
- CliToolsConsole.tsx
- net/http.Request
- .UpdateCombo
- paths
- Module
- setupMediaTest
- CombosConsole.tsx
- settingDefinition
- Config
- buildApp
- Module
- repo
- DB
- Module
- ref_node_path
- repo
- listSettingsWithQueryer
- properties
- net/http.ResponseWriter
- modeleditor/engine.go
- newTestHarness
- collectDynamicHostTelemetry
- GatewayInterceptor
- net/url.URL
- properties
- post
- Module
- awsQuickRuntime
- repository
- ProviderTopology.tsx
- New
- settingDefinition
- repo
- properties
- transform.go
- repository
- net/http/httptest.ResponseRecorder
- compressor.go
- CanonicalRequest
- repository
- properties
- application/json
- responsesStreamEncoder
- NewStaticUIHandler
- parameters
- properties
- OMEG4MON Router
- ResolvedTarget
- ref_bun_test
- certManager
- repo
- pxpipe/repo.go
- TranslatorConsole.tsx
- properties
- release-binaries
- steps
- evaluateAntigravityQuota
- repo
- combos/module.go
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- Module
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- 200
- RouteMounts
- properties
- cli/package.json
- repository
- .executeCandidateAttempt
- schemas
- get
- scripts
- post
- schemas
- apply.go
- PreviewPluginsResponse
- properties
- conflict_policy
- omnirouteModelsParity.test.ts
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- oauthLoopbackProxy
- Review Focus
- Review Focus
- CloudSyncAliasSnapshot
- schemas
- properties
- RunLegacyImportResponse
- palette.test.ts
- DeleteSkillResponse
- version
- Module
- modeleditor/engine.go
- parameters
- .executeCandidateAttempt
- tgTestEnv
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- UsageRecord
- docker-release
- time.Time
- `@omeg4mon/router` CLI Launcher
- headroomCache
- setupModelEditorTest
- context.Context
- sync.Mutex
- applyUpstreamHeaders
- parameters
- repo
- .executeCandidateAttempt
- omeg4mon-router Verification Gate
- .list
- client.test.ts
- single-binary
- keyLimiter
- verify-web
- CanonicalStreamChunk
- smoke-test.sh
- HeadroomStats
- mapBusEventToAudit
- ComboSnapshot
- repository
- KeyProvisioner
- stubModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- paths
- put
- stubSystemModule
- nullableString
- repository
- EventListResponse
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- settingsResponse
- setupTunnelTest
- Module
- batcher
- pricingListResponse
- patch
- Skill
- applyUpstreamHeaders
- repository
- stubDep
- stubModule
- WriteJSON
- go_pkg_strings
- jsonMarshalerError
- scanCombo
- stubModule
- stubModule
- omeg4mon-router Project Standards

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 230 edges
2. `WriteError()` - 215 edges
3. `WriteJSON()` - 215 edges
4. `DB` - 103 edges
5. `Button` - 98 edges
6. `buildApp()` - 85 edges
7. `Bus` - 79 edges
8. `Alert` - 75 edges
9. `AlertDescription` - 74 edges
10. `Badge()` - 70 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 3-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 3-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (248 total, 30 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.15
Nodes (27): ModelOption, patchConfigRequest, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, Module, ChangedEvent (+19 more)

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
Cohesion: 0.06
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+19 more)

### Community 6 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 7 - "Module"
Cohesion: 0.06
Nodes (10): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), defaultBuiltInAdapters(), xrouterBuiltInOAuthAdapters(), Module, TestWave2ProviderAdaptersRegisteredAndTestable(), builtinAdapter, customTestAdapter (+2 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 10 - "Button"
Cohesion: 0.05
Nodes (111): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Functional & UX Scope, class-variance-authority, lucide-react, react, react-router, AuthGuard() (+103 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (137): testHarness, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestSentinelErrorsDistinct(), doJSON(), setupArenaTest() (+129 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "Module"
Cohesion: 0.11
Nodes (10): decodeStrictJSON(), copyNonHopHeaders(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL() (+2 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.07
Nodes (41): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CliToolsConsole, CloudSyncStudio, CombosConsole, DashboardLayout() (+33 more)

### Community 16 - "catalog.go"
Cohesion: 0.21
Nodes (17): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+9 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Button"
Cohesion: 0.09
Nodes (43): ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+35 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.13
Nodes (35): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ResponseFormat, SystemBlock, ToolChoice, ToolDefinition (+27 more)

### Community 21 - "time.Duration"
Cohesion: 0.09
Nodes (15): stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), newChainStore(), TestTrafficShaperAdaptiveHeaderPacing() (+7 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "WriteError"
Cohesion: 0.10
Nodes (17): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module (+9 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "io/fs.FS"
Cohesion: 0.03
Nodes (12): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), stubModule, orderStubModule, stubDiscoverySystemModule, stubSystemModule (+4 more)

### Community 29 - "Module"
Cohesion: 0.21
Nodes (8): extractUsageCost(), actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken(), Module

### Community 30 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "Button"
Cohesion: 0.06
Nodes (71): buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isKeyExpiredAt(), KeyExpiryState (+63 more)

### Community 35 - "navigation.ts"
Cohesion: 0.14
Nodes (21): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, RawModuleEntry, UILocale, ActiveNavContext, EXTENSION_GROUPS, findActiveNavContext(), isModuleRegistered() (+13 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (22): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+14 more)

### Community 39 - "go_pkg_strings"
Cohesion: 0.05
Nodes (35): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+27 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (73): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, @hookform/resolvers, useProvidersCatalog(), useSystemTelemetry() (+65 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.16
Nodes (23): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+15 more)

### Community 42 - "New"
Cohesion: 0.06
Nodes (34): File Map, Global Constraints, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, Options (+26 more)

### Community 43 - "ModelCatalogConsole.tsx"
Cohesion: 0.08
Nodes (64): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update, AssignKeysModal(), AssignKeysModalProps, classifyTestError() (+56 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (23): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+15 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "get"
Cohesion: 0.17
Nodes (17): description, get, delete, get, get, operationId, responses, summary (+9 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+12 more)

### Community 49 - "JobInspector"
Cohesion: 0.16
Nodes (17): TestAppCloseTeardownOrderIsLIFO(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, JobInspector, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune() (+9 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.13
Nodes (28): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+20 more)

### Community 53 - "Bus"
Cohesion: 0.06
Nodes (28): Options, New(), clientIP(), decodeStrictJSON(), Module, isValidSeverity(), New(), formatUTC() (+20 more)

### Community 54 - "main.go"
Cohesion: 0.20
Nodes (3): ToolDescriptor, PoolTestResponse, storedSetting

### Community 55 - "SystemOverview.tsx"
Cohesion: 0.23
Nodes (11): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+3 more)

### Community 56 - "extractErrorMessage"
Cohesion: 0.04
Nodes (92): CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest, UpdateCloudSyncConfigRequest (+84 more)

### Community 57 - "dependencies"
Cohesion: 0.08
Nodes (24): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, zustand (+16 more)

### Community 58 - "oauth_registry.go"
Cohesion: 0.12
Nodes (33): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies() (+25 more)

### Community 59 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (185): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, runModelIdentityProbe(), AddCustomModelModal() (+177 more)

### Community 61 - "VouchersConsole.tsx"
Cohesion: 0.07
Nodes (56): zod, zod@^4.6.7, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES (+48 more)

### Community 62 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.10
Nodes (22): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), oauthEnvPrefix(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice() (+14 more)

### Community 64 - "UpstreamCredentials"
Cohesion: 0.06
Nodes (27): catalogFakeProviderStore, fakeProviderStore, ModelIdentityResult, ModelTestResult, resolveSnapshot, testAllModelsResponse, stubProviderStore, stubProviderStore (+19 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (32): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+24 more)

### Community 67 - "OAuthConfig"
Cohesion: 0.19
Nodes (6): OAuthConfig, originOfURL(), buildAuthorizeURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "Module"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 71 - "TokenSaverConsole.tsx"
Cohesion: 0.21
Nodes (17): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), parsePreviewBodyJson() (+9 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.06
Nodes (11): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, readSessionCookie(), Module, stubAdminAuth, stubAdminAuth (+3 more)

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "patch"
Cohesion: 0.13
Nodes (16): content, content, content, description, patch, schema, application/json, operationId (+8 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (21): buildDiscoveryCandidateURLs(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor, catalogOAuthDescriptor (+13 more)

### Community 76 - "promptinjection/engine.go"
Cohesion: 0.06
Nodes (34): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+26 more)

### Community 77 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 78 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 79 - "Load"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.10
Nodes (27): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+19 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.16
Nodes (25): CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation(), useRollbackCliToolConfigMutation() (+17 more)

### Community 82 - "net/http.Request"
Cohesion: 0.08
Nodes (21): Module, clearSessionCookie(), isSecureRequest(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), anthropicErrorType(), extractAPIKey() (+13 more)

### Community 83 - ".UpdateCombo"
Cohesion: 0.23
Nodes (12): ComboStep, createComboParams, updateComboValues, cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps() (+4 more)

### Community 84 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 85 - "Module"
Cohesion: 0.05
Nodes (22): streamResponseBody(), decodeStrictJSON(), New(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), newID(), newRepo() (+14 more)

### Community 86 - "setupMediaTest"
Cohesion: 0.15
Nodes (14): New(), setupCloudSyncTest(), TestCloudSyncPushPullAndSnapshotConflictPolicies(), TestCloudSyncSSRFGuardAndVaultedTokenIsolation(), newRepository(), doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints() (+6 more)

### Community 87 - "CombosConsole.tsx"
Cohesion: 0.11
Nodes (34): ComboBuilderModal(), ComboBuilderModalProps, CombosConsole(), useCombos(), useCreateComboMutation(), useDeleteComboMutation(), useUpdateComboMutation(), Combo (+26 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "Config"
Cohesion: 0.06
Nodes (33): Options, testHarness, TestAppWave1CUsageAndCLIToolsEndToEnd(), File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Task 2: `security` Module (Security Posture & Redacted Audit Trail) (+25 more)

### Community 90 - "buildApp"
Cohesion: 0.05
Nodes (51): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting(), TestWave2AE2ETranslatorAndCatalog() (+43 more)

### Community 91 - "Module"
Cohesion: 0.18
Nodes (7): restoreSnapshotRequest, SnapshotBundle, computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 92 - "repo"
Cohesion: 0.16
Nodes (12): formatUTC(), generateCode(), newID(), newRepo(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput (+4 more)

### Community 93 - "DB"
Cohesion: 0.06
Nodes (28): testEnv, New(), newRepo(), New(), newRepository(), newRepository(), Options, New() (+20 more)

### Community 94 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), New(), toJobDTO(), JobStatusDTO

### Community 95 - "ref_node_path"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 97 - "listSettingsWithQueryer"
Cohesion: 0.17
Nodes (10): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, lookupDefinition(), applyResult, historyListResp, historyResponse, rollbackResult (+2 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (34): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, trackedResponseWriter, buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol() (+26 more)

### Community 100 - "modeleditor/engine.go"
Cohesion: 0.35
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 101 - "newTestHarness"
Cohesion: 0.08
Nodes (23): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart(), TestKeysCRUDAndSecretHygiene() (+15 more)

### Community 102 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 103 - "GatewayInterceptor"
Cohesion: 0.04
Nodes (80): Options, main(), runMain(), Options, File Map, Global Constraints, Review Focus, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment (+72 more)

### Community 104 - "net/url.URL"
Cohesion: 0.09
Nodes (24): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL() (+16 more)

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 108 - "awsQuickRuntime"
Cohesion: 0.14
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 109 - "repository"
Cohesion: 0.29
Nodes (5): boolToInt(), newRepository(), scanBypass(), BypassRule, repository

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (19): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+11 more)

### Community 111 - "New"
Cohesion: 0.08
Nodes (34): loadUI(), loadUI(), Options, AuthVerifier, New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD() (+26 more)

### Community 112 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "transform.go"
Cohesion: 0.19
Nodes (17): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+9 more)

### Community 116 - "repository"
Cohesion: 0.27
Nodes (5): Config, HistoryEvent, repository, storedConfig, boolToInt()

### Community 118 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.10
Nodes (27): Claims, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+19 more)

### Community 119 - "compressor.go"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 120 - "CanonicalRequest"
Cohesion: 0.09
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), computePrefixFingerprint()

### Community 121 - "repository"
Cohesion: 0.16
Nodes (13): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), newRepository(), sanitizeStringSlice(), scanSkill(), PluginSpec (+5 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.18
Nodes (12): content, content, description, content, content, description, schema, application/json (+4 more)

### Community 124 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 125 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 126 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "ResolvedTarget"
Cohesion: 0.14
Nodes (14): strategyState, weightedCandidate, stubProviderResolver, stubResolver, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit() (+6 more)

### Community 130 - "ref_bun_test"
Cohesion: 0.11
Nodes (15): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, consoleSource, detailSource, modalsSource, providersDir, modalsSource (+7 more)

### Community 131 - "certManager"
Cohesion: 0.12
Nodes (9): newCertManager(), randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide (+1 more)

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "pxpipe/repo.go"
Cohesion: 0.39
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 134 - "TranslatorConsole.tsx"
Cohesion: 0.15
Nodes (21): react-hook-form, translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+13 more)

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

### Community 140 - "repo"
Cohesion: 0.18
Nodes (9): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotificationRow(), EmitPayload, Notification (+1 more)

### Community 141 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 142 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "Module"
Cohesion: 0.19
Nodes (3): decodeJSON(), Module, New()

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - "RouteMounts"
Cohesion: 0.15
Nodes (5): actorFromContext(), decodeJSONBody(), Module, AdminSubject(), RouteMounts

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "repository"
Cohesion: 0.24
Nodes (6): marshalOverrideJSON(), nullFloat(), nullInt(), scanOverride(), ModelOverride, scanner

### Community 153 - ".executeCandidateAttempt"
Cohesion: 0.24
Nodes (7): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, isModelAllowed(), isModelDisabled(), stubModelLister

### Community 154 - "schemas"
Cohesion: 0.11
Nodes (19): additionalProperties, properties, required, type, additionalProperties, required, type, additionalProperties (+11 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, check, check:fix, dev, format, format:check, gen:api (+4 more)

### Community 157 - "post"
Cohesion: 0.27
Nodes (11): post, post, requestBody, /api/skills, operationId, requestBody, summary, tags (+3 more)

### Community 158 - "schemas"
Cohesion: 0.15
Nodes (13): components, schemas, additionalProperties, required, type, additionalProperties, required, type (+5 more)

### Community 159 - "apply.go"
Cohesion: 0.70
Nodes (4): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic()

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 161 - "properties"
Cohesion: 0.08
Nodes (26): type, properties, additionalProperties, properties, required, type, type, type (+18 more)

### Community 162 - "conflict_policy"
Cohesion: 0.20
Nodes (10): enum, type, conflict_policy, snapshot, additionalProperties, properties, required, type (+2 more)

### Community 163 - "omnirouteModelsParity.test.ts"
Cohesion: 0.25
Nodes (7): appSource, comboBuilderSource, consoleSource, hooksSource, legacyConsolePath, modalSource, modelsDir

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.20
Nodes (10): $ref, additionalProperties, properties, required, type, items, type, plugins (+2 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "oauthLoopbackProxy"
Cohesion: 0.16
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "CloudSyncAliasSnapshot"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, alias, target_model, CloudSyncAliasSnapshot (+1 more)

### Community 171 - "schemas"
Cohesion: 0.22
Nodes (9): schemas, additionalProperties, required, type, additionalProperties, required, type, ImportedProvisionedKey (+1 more)

### Community 172 - "properties"
Cohesion: 0.22
Nodes (9): type, properties, type, type, id, KeyPrefix, Name, Secret (+1 more)

### Community 173 - "RunLegacyImportResponse"
Cohesion: 0.22
Nodes (9): provisioned_keys, run, type, $ref, additionalProperties, properties, required, type (+1 more)

### Community 174 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 178 - "modeleditor/engine.go"
Cohesion: 0.18
Nodes (6): memoryResponseWriter, extractVendorHeaders(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), syncBuffer

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - ".executeCandidateAttempt"
Cohesion: 0.31
Nodes (5): cachedGatewayResponse, inflightCall, responseCoalescer, newResponseCoalescer(), TestResponseCoalescerDeduplicatesAndCachesPerKey()

### Community 181 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

### Community 182 - "CloudSyncHistoryListResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, events, CloudSyncHistoryListResponse

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "SkillListResponse"
Cohesion: 0.29
Nodes (7): skills, SkillListResponse, additionalProperties, properties, required, type, type

### Community 186 - "UsageRecord"
Cohesion: 0.20
Nodes (7): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newDeviceOAuthRuntime(), recordingUsage, deviceOAuthRuntime

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "time.Time"
Cohesion: 0.11
Nodes (12): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, newLoginRateLimiter(), newFakeClock() (+4 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupModelEditorTest() (+11 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (21): authHarness, fakeCredentialStore, Options, fakeSettingsStore, stubAdminCredentialStore, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, stubSettings, AdminCredentialStore (+13 more)

### Community 194 - "sync.Mutex"
Cohesion: 0.29
Nodes (7): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 195 - "applyUpstreamHeaders"
Cohesion: 0.29
Nodes (6): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 196 - "parameters"
Cohesion: 0.47
Nodes (9): parameters, parameters, in, name, required, schema, parameters, parameters (+1 more)

### Community 197 - "repo"
Cohesion: 0.24
Nodes (7): formatUTC(), newRepo(), newTeamID(), createTeamInput, patchTeamInput, repo, Team

### Community 198 - ".executeCandidateAttempt"
Cohesion: 0.20
Nodes (10): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), buildCompactedResponse(), randomResponseID(), isTransientSingleCandidateStatus() (+2 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 201 - "client.test.ts"
Cohesion: 0.40
Nodes (4): openapi-fetch, openapi-fetch@^0.14.0, createApiClient(), CustomFetch

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.33
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 207 - "HeadroomStats"
Cohesion: 0.50
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 208 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 209 - "ComboSnapshot"
Cohesion: 0.22
Nodes (11): fakeComboStore, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables(), readLegacyCombos() (+3 more)

### Community 210 - "repository"
Cohesion: 0.23
Nodes (3): DisabledModel, ModelAlias, repository

### Community 211 - "KeyProvisioner"
Cohesion: 0.22
Nodes (10): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, TeamAssigner (+2 more)

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "paths"
Cohesion: 0.21
Nodes (14): description, post, post, post, components, paths, /api/cloud-sync/pull, /api/cloud-sync/push (+6 more)

### Community 215 - "put"
Cohesion: 0.40
Nodes (5): put, /api/plugins/{key}, operationId, summary, tags

### Community 218 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 222 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 223 - "Module"
Cohesion: 0.21
Nodes (3): Module, validateSQLiteSource(), decodeJSON()

### Community 224 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

### Community 226 - "patch"
Cohesion: 0.50
Nodes (4): patch, operationId, summary, tags

### Community 227 - "Skill"
Cohesion: 0.50
Nodes (4): Skill, additionalProperties, required, type

### Community 228 - "applyUpstreamHeaders"
Cohesion: 0.16
Nodes (11): Task 2: Backend Model Test Endpoint (TDD), Approaches Considered, Constraints, Goals, Intent & Why, Non-Goals, Provider Detail: Proportional Model Grid & Per-Model Test, isBlockedAntigravityProjectID() (+3 more)

### Community 229 - "repository"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 233 - "WriteJSON"
Cohesion: 0.10
Nodes (12): Module, decodeJSON(), deriveSessionTitle(), New(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON() (+4 more)

### Community 235 - "go_pkg_strings"
Cohesion: 0.11
Nodes (19): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+11 more)

### Community 246 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1123 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `net/http.ResponseWriter`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14500716674629718 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `net/http.ResponseWriter` to `UpstreamCredentials`, `context.Context`, `go_pkg_encoding_json`, `net/http.Request`, `ProvidersModals.tsx`?**
  _High betweenness centrality (0.133) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `GatewayInterceptor` to `context.Context`, `newTestHarness`, `Button`, `New`, `Module`, `net/http/httptest.ResponseRecorder`, `buildApp`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._