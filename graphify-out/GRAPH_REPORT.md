# Graph Report - omeg4mon-router  (2026-11-25)

## Corpus Check
- 457 files · ~494,087 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5667 nodes · 18561 edges · 234 communities (203 shown, 31 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1761 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `9719b39e`
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
- SettingsConsole.tsx
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
- Card
- updates
- ir.go
- usage/schemas.ts
- with
- jobs
- net/http.ResponseWriter
- steps
- web-check
- on
- io/fs.FS
- OverviewRealtimeTraffic.tsx
- TelegramConsole.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- Button
- navigation.ts
- Go Module Root Definition
- CanonicalUsage
- main_test.go
- UsageConsole.tsx
- encoding/json.RawMessage
- New
- models/catalogUtils.ts
- translate_gemini.go
- Module
- repository
- paths
- translate_openai.go
- Runner
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- sync.Mutex
- main.go
- overview.ts
- TranslatorConsole.tsx
- dependencies
- oauth_registry.go
- ModelCatalogConsole.tsx
- ProvidersModals.tsx
- vouchers/hooks.ts
- modelIntelligence.ts
- ProviderAdapter
- UpstreamCredentials
- teeHandler
- Module
- CanonicalChatResponse
- scripts
- pricingStore
- extractErrorMessage
- Module
- .AuthenticateAdmin
- Module
- Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module
- go_pkg_encoding_json
- promptinjection/engine.go
- oauthTokenResponse
- Module
- Load
- web/package.json
- CliToolsConsole.tsx
- Module
- repository
- paths
- collectDynamicHostTelemetry
- collectDynamicHostTelemetry
- TokenSaverConsole.tsx
- settingDefinition
- Load
- buildApp
- net/http.ResponseWriter
- repo
- New
- Module
- deviceOAuthSession
- repo
- listSettingsWithQueryer
- properties
- Module
- Module
- .StatusCode
- WithProviderID
- GatewayInterceptor
- combos/hooks.ts
- properties
- post
- WriteJSON
- Module
- parseJSONObjectOrWrap
- ProviderTopology.tsx
- setupMediaTest
- go_pkg_strings
- repo
- properties
- transform.go
- repository
- Module
- Module
- compressor.go
- CanonicalRequest
- repository
- properties
- application/json
- newTestHarness
- net/url.URL
- plugins/interceptor.go
- properties
- OMEG4MON Router
- applyComboStrategy
- @radix-ui/react-label
- certManager
- repo
- ProviderAdapter
- mapBusEventToAudit
- properties
- release-binaries
- steps
- fakeHandle
- repository
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- 200
- DB
- properties
- cli/package.json
- modeleditor/engine.go
- .executeCandidateAttempt
- schemas
- get
- scripts
- post
- schemas
- ClearCloudSyncHistoryResponse
- PreviewPluginsResponse
- properties
- conflict_policy
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
- DeleteSkillResponse
- version
- Module
- parameters
- .executeCandidateAttempt
- ref_bun_test
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- Module
- docker-release
- time.Time
- `@omeg4mon/router` CLI Launcher
- setupModelEditorTest
- context.Context
- parameters
- repo
- clitools/module_test.go
- omeg4mon-router Verification Gate
- single-binary
- ProcessHandle
- verify-web
- CanonicalStreamChunk
- smoke-test.sh
- ComboSnapshot
- repo
- stubModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- repository
- put
- stubSystemModule
- stubDiscoverySystemModule
- stubSystemModule
- models/hooks.ts
- Module
- Module
- pricingStore
- patch
- Skill
- repository
- stubDep
- stubModule
- net/http.ResponseWriter
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- stubDepModule
- jsonMarshalerError
- vite.config.ts
- stubSystemModule
- stubSystem
- stubModule
- stubDep
- stubModule
- omeg4mon-router Project Standards
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- io.Reader

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 228 edges
2. `WriteError()` - 212 edges
3. `WriteJSON()` - 209 edges
4. `DB` - 103 edges
5. `Button` - 96 edges
6. `buildApp()` - 85 edges
7. `Bus` - 79 edges
8. `Alert` - 73 edges
9. `AlertDescription` - 72 edges
10. `Badge()` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/module_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (234 total, 31 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.08
Nodes (53): ModelOption, configRequest, patchConfigRequest, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse (+45 more)

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
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+19 more)

### Community 6 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (45): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, zod, zod@^4.6.7, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+37 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 10 - "Button"
Cohesion: 0.06
Nodes (83): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Functional & UX Scope, class-variance-authority, lucide-react, react, react-router, AuthGuard(), AuthGuardProps (+75 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (158): testHarness, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+150 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "Module"
Cohesion: 0.11
Nodes (9): decodeStrictJSON(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID() (+1 more)

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

### Community 18 - "Card"
Cohesion: 0.05
Nodes (80): ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+72 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (34): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+26 more)

### Community 21 - "usage/schemas.ts"
Cohesion: 0.10
Nodes (22): FilteredRecentRequestsResult, DEFAULT_LOGS_LIMIT, ListModelPricingResponse, MAX_LOGS_LIMIT, MAX_MODEL_IDENTIFIER_RUNES, ModelPricingSource, ModelUsageBreakdown, parseNonNegativeRate() (+14 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (13): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module (+5 more)

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
Cohesion: 0.10
Nodes (8): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 29 - "OverviewRealtimeTraffic.tsx"
Cohesion: 0.14
Nodes (28): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), formatActivityOriginLabel(), formatUsageBucketLabel(), buildRequestChartModel(), ChartPoint, ChartYTick, filterRecentRequests(), formatChartMetricValue() (+20 more)

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
Cohesion: 0.11
Nodes (27): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, RawModuleEntry, UILocale, ActiveNavContext (+19 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "main_test.go"
Cohesion: 0.07
Nodes (19): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+11 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.11
Nodes (34): Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, useProvidersCatalog(), useProxyPools(), useSystemTelemetry(), formatBytesGB() (+26 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 42 - "New"
Cohesion: 0.08
Nodes (30): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+22 more)

### Community 43 - "models/catalogUtils.ts"
Cohesion: 0.13
Nodes (25): AssignKeysModalProps, computeUpdatedAllowedModels(), countKeyAccess(), DEFAULT_COMBO_FILTERS, DEFAULT_MODEL_FILTERS, keyAllowsIdentifier(), MAX_CACHED_TEST_RESULTS, patternMatchesIdentifier() (+17 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (54): content, description, content, description, content, description, get, patch (+46 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.21
Nodes (17): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+9 more)

### Community 49 - "Runner"
Cohesion: 0.07
Nodes (32): TestAppCloseTeardownOrderIsLIFO(), TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification (+24 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.09
Nodes (27): COMBO_STRATEGIES, CreateComboFormInput, CreateComboFormOutput, createComboSchema, DEFAULT_CREATE_COMBO_VALUES, DEFAULT_DISABLE_MODEL_VALUES, DEFAULT_UPSERT_ALIAS_VALUES, DisableModelFormInput (+19 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (25): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+17 more)

### Community 53 - "sync.Mutex"
Cohesion: 0.22
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 54 - "main.go"
Cohesion: 0.23
Nodes (4): lookupPluginSpec(), ToolDescriptor, PluginSpec, storedSetting

### Community 55 - "overview.ts"
Cohesion: 0.23
Nodes (11): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+3 more)

### Community 56 - "TranslatorConsole.tsx"
Cohesion: 0.09
Nodes (37): react-hook-form, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+29 more)

### Community 57 - "dependencies"
Cohesion: 0.09
Nodes (23): clsx, @fontsource-variable/jetbrains-mono, openapi-fetch, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies (+15 more)

### Community 58 - "oauth_registry.go"
Cohesion: 0.18
Nodes (21): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), decodeJWTPayload(), enrichAntigravity(), enrichGeminiOAuth(), enrichGitHubCopilot() (+13 more)

### Community 59 - "ModelCatalogConsole.tsx"
Cohesion: 0.18
Nodes (23): clearCachedTestResults(), comboTestResultKey(), computeVisibilityCounts(), filterCombos(), filterModels(), formatRelativeTime(), formatTokenCount(), getEnabledCapabilities() (+15 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (181): Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, AssignKeysModal(), runModelIdentityProbe(), useUpdateKeyAllowedModelsMutation(), AddCustomModelModal(), AddCustomModelModalProps (+173 more)

### Community 61 - "vouchers/hooks.ts"
Cohesion: 0.10
Nodes (35): teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation(), useRemoveTeamKeyMutation(), useResetTeamSpendMutation(), useTeamsList(), useUpdateTeamMutation() (+27 more)

### Community 62 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (18): decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository(), scanConnectionRecord() (+10 more)

### Community 64 - "UpstreamCredentials"
Cohesion: 0.05
Nodes (35): catalogFakeProviderStore, fakeProviderStore, resolveSnapshot, weightedCandidate, stubProviderResolver, stubProviderStore, Global Constraints, Provider Setup Modals Redesign & Pre-Save Test Implementation Plan (+27 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "Module"
Cohesion: 0.36
Nodes (3): BinariesStatus, BinaryStatus, Status

### Community 67 - "CanonicalChatResponse"
Cohesion: 0.17
Nodes (18): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString(), RewriteModelField() (+10 more)

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "extractErrorMessage"
Cohesion: 0.04
Nodes (94): CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest, UpdateCloudSyncConfigRequest (+86 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.08
Nodes (10): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, TestContractInterfacesSatisfied(), stubAdminAuth, stubAdminAuth, stubAdminAuth (+2 more)

### Community 73 - "Module"
Cohesion: 0.13
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 74 - "Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module"
Cohesion: 0.16
Nodes (9): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, ProvisionedKey, ProvisionKeyRequest, TeamAssigner, stubKeyProvisioner (+1 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.11
Nodes (24): buildDiscoveryCandidateURLs(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor, catalogOAuthDescriptor (+16 more)

### Community 76 - "promptinjection/engine.go"
Cohesion: 0.07
Nodes (33): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+25 more)

### Community 77 - "oauthTokenResponse"
Cohesion: 0.24
Nodes (12): isAnthropicTokenURL(), defaultExchangePKCECode(), defaultOAuthStrategies(), defaultRefreshFormToken(), exchangeClaudeCode(), newRFC8628Strategy(), refreshAWSOIDCToken(), refreshClaudeToken() (+4 more)

### Community 78 - "Module"
Cohesion: 0.19
Nodes (5): computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 79 - "Load"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "web/package.json"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.16
Nodes (25): CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation(), useRollbackCliToolConfigMutation() (+17 more)

### Community 82 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 83 - "repository"
Cohesion: 0.11
Nodes (20): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues (+12 more)

### Community 84 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.19
Nodes (14): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), TestTelemetryCacheAvoidsRedundantWork() (+6 more)

### Community 86 - "collectDynamicHostTelemetry"
Cohesion: 0.31
Nodes (10): collectDynamicHostTelemetry(), collectStaticHostInfo(), readBoundedFile(), parseCPUModel(), parseHostUptime(), parseLoadAvg(), parseMemInfo(), parseOSPrettyName() (+2 more)

### Community 87 - "TokenSaverConsole.tsx"
Cohesion: 0.21
Nodes (17): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), parsePreviewBodyJson() (+9 more)

### Community 88 - "settingDefinition"
Cohesion: 0.14
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 90 - "buildApp"
Cohesion: 0.04
Nodes (53): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+45 more)

### Community 91 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (17): Module, trackedResponseWriter, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeJSONBody(), isValidMode() (+9 more)

### Community 92 - "repo"
Cohesion: 0.08
Nodes (20): buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), formatUTC(), generateCode() (+12 more)

### Community 93 - "New"
Cohesion: 0.09
Nodes (19): loadUI(), loadUI(), apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement(), TestDevUIHandlerFallback(), TestGatewaySSEStreamsPastAPITimeout() (+11 more)

### Community 94 - "Module"
Cohesion: 0.14
Nodes (7): decodeStrictJSON(), Module, isValidSeverity(), mapBusEventToNotification(), toJobDTO(), JobStatusDTO, Notification

### Community 95 - "deviceOAuthSession"
Cohesion: 0.23
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 96 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 97 - "listSettingsWithQueryer"
Cohesion: 0.10
Nodes (18): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, SettingItem, int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds() (+10 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "Module"
Cohesion: 0.09
Nodes (25): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), isValidAuthType(), isValidProtocol(), oauthEnvPrefix(), validateHTTPURL(), buildAwsQuickSigninURL() (+17 more)

### Community 100 - "Module"
Cohesion: 0.18
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 101 - ".StatusCode"
Cohesion: 0.18
Nodes (4): statusError, statusError, httpStatusError, statusError

### Community 102 - "WithProviderID"
Cohesion: 0.12
Nodes (17): ModelIdentityResult, ModelTestResult, testAllModelsResponse, buildModelProbeRequest(), extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+9 more)

### Community 103 - "GatewayInterceptor"
Cohesion: 0.06
Nodes (54): main(), Options, Global Constraints, Review Focus, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints (+46 more)

### Community 104 - "combos/hooks.ts"
Cohesion: 0.12
Nodes (17): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+9 more)

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "WriteJSON"
Cohesion: 0.09
Nodes (13): validateSQLiteSource(), decodeJSON(), decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), WriteJSON() (+5 more)

### Community 108 - "Module"
Cohesion: 0.15
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 109 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.15
Nodes (18): @xyflow/react, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps, RouterHubNode() (+10 more)

### Community 111 - "setupMediaTest"
Cohesion: 0.18
Nodes (12): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), isLinkUnsupported() (+4 more)

### Community 112 - "go_pkg_strings"
Cohesion: 0.13
Nodes (13): ComboTargetTestResult, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest, testComboRequest, testComboResponse, testModelRequest (+5 more)

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 116 - "repository"
Cohesion: 0.27
Nodes (5): Config, HistoryEvent, repository, storedConfig, boolToInt()

### Community 117 - "Module"
Cohesion: 0.21
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 118 - "Module"
Cohesion: 0.05
Nodes (48): authHarness, Options, TestAppWave1AAuthSettingsKeysEndToEnd(), Claims, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen (+40 more)

### Community 119 - "compressor.go"
Cohesion: 0.49
Nodes (8): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput(), TestTokenSaverFinalReviewFindings()

### Community 120 - "CanonicalRequest"
Cohesion: 0.09
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), injectAnthropicCacheControl(), computePrefixFingerprint()

### Community 121 - "repository"
Cohesion: 0.19
Nodes (10): boolToInt(), cloneMap(), mergeDefaults(), sanitizeStringSlice(), scanSkill(), PluginState, repository, scanner (+2 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.18
Nodes (12): content, content, description, content, content, description, schema, application/json (+4 more)

### Community 124 - "newTestHarness"
Cohesion: 0.25
Nodes (10): testHarness, File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback() (+2 more)

### Community 125 - "net/url.URL"
Cohesion: 0.12
Nodes (17): stubProxyDialerResolver, stubProxyResolver, isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+9 more)

### Community 126 - "plugins/interceptor.go"
Cohesion: 0.43
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "applyComboStrategy"
Cohesion: 0.21
Nodes (8): CatalogModel, CatalogProviderGroup, modelCatalogResponse, strategyState, inferCatalogModel(), applyComboStrategy(), cryptoRandFloat64(), cryptoRandIntn()

### Community 130 - "@radix-ui/react-label"
Cohesion: 0.33
Nodes (4): @fontsource-variable/plus-jakarta-sans, @radix-ui/react-label, @fontsource-variable/plus-jakarta-sans@^5.3.0, @radix-ui/react-label@^2.1.17

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "ProviderAdapter"
Cohesion: 0.09
Nodes (10): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), defaultBuiltInAdapters(), OAuthConfig, xrouterBuiltInOAuthAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable(), buildAuthorizeURL(), builtinAdapter (+2 more)

### Community 134 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 142 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - "DB"
Cohesion: 0.03
Nodes (75): Options, Options, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI) (+67 more)

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "modeleditor/engine.go"
Cohesion: 0.05
Nodes (35): arenaTestEnv, stubCompleter, ChatMessage, memoryResponseWriter, ChatCompletionRequest, ChatCompletionResult, extractVendorHeaders(), extractGatewayError() (+27 more)

### Community 153 - ".executeCandidateAttempt"
Cohesion: 0.08
Nodes (24): stubModelLister, stubModelLister, Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD) (+16 more)

### Community 154 - "schemas"
Cohesion: 0.15
Nodes (13): additionalProperties, required, type, additionalProperties, required, type, components, schemas (+5 more)

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

### Community 159 - "ClearCloudSyncHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearCloudSyncHistoryResponse

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 161 - "properties"
Cohesion: 0.08
Nodes (26): type, additionalProperties, properties, required, type, properties, type, type (+18 more)

### Community 162 - "conflict_policy"
Cohesion: 0.20
Nodes (10): enum, type, conflict_policy, snapshot, additionalProperties, properties, required, type (+2 more)

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
Cohesion: 0.08
Nodes (16): newAwsQuickRuntime(), validateAwsRegion(), newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), awsQuickDiscovery, awsQuickResult (+8 more)

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

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "Module"
Cohesion: 0.16
Nodes (3): Module, validatePassword(), HashPassword()

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - ".executeCandidateAttempt"
Cohesion: 0.12
Nodes (15): attemptOutcome, isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), NewStreamTranslator(), buildCompactedResponse(), randomResponseID() (+7 more)

### Community 181 - "ref_bun_test"
Cohesion: 0.06
Nodes (29): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+21 more)

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

### Community 186 - "Module"
Cohesion: 0.06
Nodes (24): stubUsageRecorder, stubUsageRecorder, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, stubUsageRecorder, testEnv, UsageRecord, newBatcher(), Module (+16 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "time.Time"
Cohesion: 0.04
Nodes (36): CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, agQuotaBucket, agQuotaCacheEntry, agQuotaGroup (+28 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupPromptGuardTest() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (16): fakeCredentialStore, fakeSettingsStore, stubAdminCredentialStore, stubSettings, TestSystemAdminCredentialStore(), sanitizeBotError(), stubModule, Alias (+8 more)

### Community 196 - "parameters"
Cohesion: 0.47
Nodes (9): parameters, parameters, in, name, required, schema, parameters, parameters (+1 more)

### Community 197 - "repo"
Cohesion: 0.27
Nodes (6): formatUTC(), newTeamID(), createTeamInput, patchTeamInput, repo, Team

### Community 198 - "clitools/module_test.go"
Cohesion: 0.24
Nodes (8): generateResponse, applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), applyResponse, catalogResponse, rollbackResponse

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "ProcessHandle"
Cohesion: 0.23
Nodes (3): fakeRunner, osProcessRunner, ProcessHandle

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 209 - "ComboSnapshot"
Cohesion: 0.19
Nodes (13): fakeComboStore, restoreSnapshotRequest, SnapshotBundle, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot (+5 more)

### Community 210 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 215 - "put"
Cohesion: 0.40
Nodes (5): put, /api/plugins/{key}, operationId, summary, tags

### Community 221 - "models/hooks.ts"
Cohesion: 0.31
Nodes (9): classifyTestError(), flattenCatalogResponse(), normalizeTestStatus(), RawCatalogGroup, modelCatalogQueryKeys, ModelIdentityResult, runSingleComboTest(), runSingleModelTest() (+1 more)

### Community 222 - "Module"
Cohesion: 0.15
Nodes (12): New(), doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary() (+4 more)

### Community 223 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 226 - "patch"
Cohesion: 0.50
Nodes (4): patch, operationId, summary, tags

### Community 227 - "Skill"
Cohesion: 0.50
Nodes (4): Skill, additionalProperties, required, type

### Community 229 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 233 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (20): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+12 more)

### Community 235 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 242 - "vite.config.ts"
Cohesion: 0.29
Nodes (6): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 255 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1114 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `UpstreamCredentials` to `ProvidersModals.tsx`?**
  _High betweenness centrality (0.177) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07909604519774012 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `UpstreamCredentials` to `Module`, `context.Context`, `go_pkg_encoding_json`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.149) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `Module` to `SettingsConsole.tsx`, `GatewayInterceptor`, `Button`, `testing.T`, `Module`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._