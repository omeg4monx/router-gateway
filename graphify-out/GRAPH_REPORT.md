# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 413 files · ~385,327 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5024 nodes · 16103 edges · 234 communities (204 shown, 30 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1572 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `36ed54f2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- server.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- web/package.json
- SettingsConsole.tsx
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- repo
- MediaConsole.tsx
- TokenSaverConsole.tsx
- catalog.go
- Design System & Wave 0 Plan Docs
- Module
- updates
- ir.go
- repository
- with
- jobs
- registry_test.go
- steps
- web-check
- on
- Label
- Module
- ProvidersConsole.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- UsageRecord
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- net/http.Request
- UsageConsole.tsx
- encoding/json.RawMessage
- buildApp
- UpstreamCredentials
- translate_gemini.go
- Module
- repository
- paths
- translate_openai.go
- ComboSnapshot
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- Module
- main.go
- navigation.ts
- Module
- WriteJSON
- VouchersConsole.tsx
- Config
- CardTitle
- repo
- Server
- ProviderAdapter
- io/fs.FS
- schema.go
- Bus
- promptinjection/engine.go
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- GatewayInterceptor
- TeamsConsole.tsx
- Module
- main.go
- repository
- fakeHandle
- Module
- Load
- JobInspector
- CliToolsConsole.tsx
- dependencies
- collectMigrations
- Module
- ProviderAdapter
- net/http.Request
- Module
- settingDefinition
- Module
- buildApp
- crypto.SessionSigner
- repo
- .executeCandidateAttempt
- New
- repository
- repo
- collectDynamicHostTelemetry
- properties
- PxpipeConsole.tsx
- stubDep
- Badge
- stubSystemModule
- Module
- RouteMounts
- properties
- post
- setupModelEditorTest
- Module
- providers/handlers.go
- .StatusCode
- Module
- stubSystemModule
- repo
- properties
- transform.go
- repository
- cloudsync/hooks.ts
- net/http/httptest.ResponseRecorder
- .ProvisionKey
- BypassRule
- repository
- properties
- Runner
- stubSystem
- Module
- stubCompleter
- properties
- OMEG4MON Router
- headroomCache
- CanonicalStreamChunk
- certManager
- repo
- settingDefinition
- antigravity.go
- properties
- resolve.js
- release-binaries
- steps
- Module
- combos/module.go
- stubModule
- scripts
- Module
- clitools/module_test.go
- setupTestMitm
- properties
- 200
- ModelDescriptor
- properties
- time.Time
- cli/package.json
- properties
- net/url.URL
- schemas
- get
- paths
- NewStaticUIHandler
- schemas
- application/json
- PreviewPluginsResponse
- post
- conflict_policy
- TestContractInterfacesSatisfied
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- repo
- Review Focus
- Review Focus
- CloudSyncAliasSnapshot
- schemas
- properties
- RunLegacyImportResponse
- parameters
- DeleteSkillResponse
- version
- stubDepModule
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- parameters
- Config
- DB
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- compressor.go
- docker-release
- put
- `@omeg4mon/router` CLI Launcher
- patch
- Skill
- context.Context
- testEnv
- ProcessHandle
- repo
- batcher
- interceptor.go
- telemetryCollector
- app/layout.tsx
- newTestHarness
- single-binary
- setupTunnelTest
- verify-web
- stubDep
- smoke-test.sh
- stubModule
- Config
- tgTestEnv
- stubModule
- NewStreamTranslator
- stubModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- setupMediaTest
- NewClient
- ClearCloudSyncHistoryResponse
- stubSystemModule
- repository
- repository
- Registry
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- net/http.Header
- geist
- @radix-ui/react-label
- tailwind-merge
- media/module.go
- .list
- io.Reader
- testEnv
- jsonMarshalerError
- clsx

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 215 edges
2. `WriteError()` - 195 edges
3. `WriteJSON()` - 192 edges
4. `DB` - 103 edges
5. `buildApp()` - 84 edges
6. `Bus` - 79 edges
7. `Button` - 73 edges
8. `Module` - 63 edges
9. `Alert` - 59 edges
10. `AlertDescription` - 58 edges

## Surprising Connections (you probably didn't know these)
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx
- `Key Proportion & UX Fixes` --references--> `Open()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-overview-design.md → internal/platform/db/db.go
- `Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)` --references--> `ModelsRoutingConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → web/src/features/combos/ModelsRoutingConsole.tsx

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (234 total, 30 thin omitted)

### Community 0 - "server.go"
Cohesion: 0.07
Nodes (50): ModelOption, patchConfigRequest, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, mapBusEventToAudit() (+42 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.22
Nodes (5): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC()

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (47): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+39 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "APIKeySummary"
Cohesion: 0.07
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+19 more)

### Community 6 - "web/package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.06
Nodes (61): @tanstack/react-query, DashboardSettingsPage(), authKeys, LoginRequest, LoginResponse, LogoutResponse, AuthPolicyResponse, ChangePasswordRequest (+53 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.16
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.08
Nodes (67): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Functional & UX Scope, class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout() (+59 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (103): testHarness, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest() (+95 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "repo"
Cohesion: 0.16
Nodes (8): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotificationRow(), EmitPayload, repo

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.12
Nodes (31): MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders() (+23 more)

### Community 15 - "TokenSaverConsole.tsx"
Cohesion: 0.19
Nodes (18): TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation() (+10 more)

### Community 16 - "catalog.go"
Cohesion: 0.21
Nodes (17): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+9 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Module"
Cohesion: 0.21
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (32): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+24 more)

### Community 21 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "Label"
Cohesion: 0.08
Nodes (45): ArenaPage(), ChatPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels() (+37 more)

### Community 29 - "Module"
Cohesion: 0.15
Nodes (10): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+2 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.05
Nodes (90): DashboardProviderDetailPage(), DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation() (+82 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "UsageRecord"
Cohesion: 0.29
Nodes (5): stubUsageRecorder, stubUsageRecorder, testEnv, UsageRecord, recordingUsage

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "net/http.Request"
Cohesion: 0.09
Nodes (17): Module, trackedResponseWriter, anthropicErrorType(), extractAPIKey(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError() (+9 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (63): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, DashboardPage(), DashboardUsagePage(), useProxyPools() (+55 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.21
Nodes (15): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+7 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "UpstreamCredentials"
Cohesion: 0.08
Nodes (21): fakeProviderStore, resolveSnapshot, targetKey, stubProviderResolver, stubProviderStore, stubProviderStore, stubResolver, ProviderConnection (+13 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "Module"
Cohesion: 0.09
Nodes (9): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth (+1 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (54): content, description, content, description, content, description, get, patch (+46 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "ComboSnapshot"
Cohesion: 0.09
Nodes (25): fakeComboStore, restoreSnapshotRequest, SnapshotBundle, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests (+17 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (72): ImporterPage(), DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation() (+64 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.16
Nodes (23): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+15 more)

### Community 53 - "Module"
Cohesion: 0.11
Nodes (20): buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol(), oauthEnvPrefix(), validateHTTPURL(), decodeJSONBody(), truncateError() (+12 more)

### Community 54 - "main.go"
Cohesion: 0.19
Nodes (9): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), ChangedEvent, PoolTestResponse, ScanCandidateResult, ScanResponse (+1 more)

### Community 55 - "navigation.ts"
Cohesion: 0.08
Nodes (34): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata (+26 more)

### Community 56 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 57 - "WriteJSON"
Cohesion: 0.14
Nodes (10): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider() (+2 more)

### Community 58 - "VouchersConsole.tsx"
Cohesion: 0.09
Nodes (39): zod, TeamsPage(), VouchersPage(), teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation(), useRemoveTeamKeyMutation() (+31 more)

### Community 59 - "Config"
Cohesion: 0.11
Nodes (33): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 60 - "CardTitle"
Cohesion: 0.14
Nodes (21): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 61 - "repo"
Cohesion: 0.24
Nodes (8): buildWhere(), formatUTC(), newEventID(), normalizeSeverity(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 62 - "Server"
Cohesion: 0.10
Nodes (16): loadUI(), loadUI(), Options, AuthVerifier, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement() (+8 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (21): cloneStringMap(), cloneStringSlice(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+13 more)

### Community 64 - "io/fs.FS"
Cohesion: 0.09
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "schema.go"
Cohesion: 0.23
Nodes (7): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "Bus"
Cohesion: 0.08
Nodes (22): Options, New(), decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), New(), Bus (+14 more)

### Community 67 - "promptinjection/engine.go"
Cohesion: 0.09
Nodes (27): globToRegexp(), actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody() (+19 more)

### Community 68 - "scripts"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "TelegramConsole.tsx"
Cohesion: 0.12
Nodes (31): TelegramPage(), telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation() (+23 more)

### Community 71 - "newTestHarness"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - "GatewayInterceptor"
Cohesion: 0.06
Nodes (52): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway (+44 more)

### Community 73 - "TeamsConsole.tsx"
Cohesion: 0.11
Nodes (31): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation() (+23 more)

### Community 75 - "main.go"
Cohesion: 0.13
Nodes (7): errorDisposition, providerIDContextKey, classifyStatus(), isClientPayloadErrorBody(), ToolDescriptor, PreviewResult, ValidationError

### Community 76 - "repository"
Cohesion: 0.19
Nodes (8): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanOverride(), ModelOverride, repository

### Community 78 - "Module"
Cohesion: 0.22
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Load"
Cohesion: 0.22
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "JobInspector"
Cohesion: 0.12
Nodes (23): TestAppCloseTeardownOrderIsLIFO(), File Map, Review Focus, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan, JobInspector (+15 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "dependencies"
Cohesion: 0.12
Nodes (17): openapi-fetch, @radix-ui/react-slot, react-dom, zustand, dependencies, class-variance-authority@^0.7.1, @hookform/resolvers@^7.0.8, lucide-react@^0.511.0 (+9 more)

### Community 83 - "collectMigrations"
Cohesion: 0.15
Nodes (12): Combo, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues, cloneStringSlice() (+4 more)

### Community 85 - "ProviderAdapter"
Cohesion: 0.10
Nodes (6): cloneOAuthConfig(), buildAuthorizeURL(), builtinAdapter, customTestAdapter, OAuthConfig, ProviderAdapter

### Community 86 - "net/http.Request"
Cohesion: 0.12
Nodes (10): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), decodeStrictJSON(), LogEntry, LogListResponse (+2 more)

### Community 87 - "Module"
Cohesion: 0.15
Nodes (4): clientIP(), decodeStrictJSON(), Module, isValidSeverity()

### Community 88 - "settingDefinition"
Cohesion: 0.14
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "Module"
Cohesion: 0.14
Nodes (5): Module, Options, New(), newPricingStore(), newRepository()

### Community 90 - "buildApp"
Cohesion: 0.08
Nodes (44): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+36 more)

### Community 91 - "crypto.SessionSigner"
Cohesion: 0.15
Nodes (22): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+14 more)

### Community 92 - "repo"
Cohesion: 0.18
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (31): attemptOutcome, CanonicalChatResponse, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage() (+23 more)

### Community 94 - "New"
Cohesion: 0.06
Nodes (36): File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+28 more)

### Community 95 - "repository"
Cohesion: 0.12
Nodes (23): PluginsPage(), SkillsPage(), CreateSkillRequest, PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse (+15 more)

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 97 - "collectDynamicHostTelemetry"
Cohesion: 0.16
Nodes (19): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+11 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "PxpipeConsole.tsx"
Cohesion: 0.11
Nodes (32): react-hook-form, PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats() (+24 more)

### Community 101 - "Badge"
Cohesion: 0.18
Nodes (19): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES (+11 more)

### Community 103 - "Module"
Cohesion: 0.20
Nodes (3): decodeJSON(), Module, New()

### Community 104 - "RouteMounts"
Cohesion: 0.09
Nodes (10): actorFromContext(), decodeJSONBody(), Module, decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), AdminSubject() (+2 more)

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (19): File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+11 more)

### Community 108 - "Module"
Cohesion: 0.12
Nodes (9): validateSQLiteSource(), decodeJSON(), decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), WriteError() (+1 more)

### Community 109 - "providers/handlers.go"
Cohesion: 0.10
Nodes (20): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+12 more)

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.13
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

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

### Community 117 - "cloudsync/hooks.ts"
Cohesion: 0.20
Nodes (17): CloudSyncPage(), CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest (+9 more)

### Community 118 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.07
Nodes (33): authHarness, fakeCredentialStore, Options, Claims, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, AdminCredentialStore (+25 more)

### Community 119 - ".ProvisionKey"
Cohesion: 0.18
Nodes (10): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, TeamAssigner (+2 more)

### Community 120 - "BypassRule"
Cohesion: 0.23
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), matchPattern(), wrapSystemPrompt(), wrapText(), scanBypass(), BypassRule (+1 more)

### Community 121 - "repository"
Cohesion: 0.16
Nodes (13): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), newRepository(), sanitizeStringSlice(), scanSkill(), PluginSpec (+5 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 125 - "Module"
Cohesion: 0.19
Nodes (5): computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 126 - "stubCompleter"
Cohesion: 0.26
Nodes (8): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), TestGatewayCompleteChat()

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload), Architecture, Configuration, License, Modules (+3 more)

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 134 - "antigravity.go"
Cohesion: 0.17
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "resolve.js"
Cohesion: 0.13
Nodes (11): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+3 more)

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "Module"
Cohesion: 0.20
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 140 - "combos/module.go"
Cohesion: 0.15
Nodes (9): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, patchComboRequest, upsertAliasRequest (+1 more)

### Community 142 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 143 - "Module"
Cohesion: 0.23
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 144 - "clitools/module_test.go"
Cohesion: 0.05
Nodes (26): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+18 more)

### Community 145 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - "ModelDescriptor"
Cohesion: 0.36
Nodes (5): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, stubModelLister

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "time.Time"
Cohesion: 0.11
Nodes (11): CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, newCooldownTracker(), newFakeClock(), newFakeClock() (+3 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "properties"
Cohesion: 0.08
Nodes (26): type, additionalProperties, properties, required, type, properties, type, type (+18 more)

### Community 153 - "net/url.URL"
Cohesion: 0.13
Nodes (12): stubProxyDialerResolver, stubProxyResolver, copyNonHopHeaders(), matchToolForHost(), newID(), singleJoiningSlash(), parseAndValidateProxyURL(), probeProxy() (+4 more)

### Community 154 - "schemas"
Cohesion: 0.15
Nodes (13): additionalProperties, required, type, additionalProperties, required, type, components, schemas (+5 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 157 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 158 - "schemas"
Cohesion: 0.15
Nodes (13): components, schemas, additionalProperties, required, type, additionalProperties, required, type (+5 more)

### Community 159 - "application/json"
Cohesion: 0.18
Nodes (12): content, content, description, content, content, description, schema, application/json (+4 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 161 - "post"
Cohesion: 0.27
Nodes (11): post, post, requestBody, /api/skills, operationId, requestBody, summary, tags (+3 more)

### Community 162 - "conflict_policy"
Cohesion: 0.20
Nodes (10): enum, type, conflict_policy, snapshot, additionalProperties, properties, required, type (+2 more)

### Community 163 - "TestContractInterfacesSatisfied"
Cohesion: 0.10
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID()

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.20
Nodes (10): $ref, additionalProperties, properties, required, type, items, type, plugins (+2 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "repo"
Cohesion: 0.27
Nodes (6): formatUTC(), newTeamID(), createTeamInput, patchTeamInput, repo, Team

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

### Community 174 - "parameters"
Cohesion: 0.47
Nodes (9): parameters, parameters, in, name, required, schema, parameters, parameters (+1 more)

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 178 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Config"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 181 - "DB"
Cohesion: 0.07
Nodes (36): Options, Options, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, New(), newRepository(), New(), TestModuleMetadataAndStartValidation() (+28 more)

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

### Community 186 - "compressor.go"
Cohesion: 0.36
Nodes (9): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), leadingBackticks(), truncateToolOutput() (+1 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "put"
Cohesion: 0.40
Nodes (5): put, /api/plugins/{key}, operationId, summary, tags

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "patch"
Cohesion: 0.50
Nodes (4): patch, operationId, summary, tags

### Community 191 - "Skill"
Cohesion: 0.50
Nodes (4): Skill, additionalProperties, required, type

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (15): fakeSettingsStore, stubAdminCredentialStore, stubSettings, newRepository(), sanitizeBotError(), stubModule, Alias, AliasListResponse (+7 more)

### Community 194 - "testEnv"
Cohesion: 0.19
Nodes (6): chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore(), formatDuration()

### Community 195 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 196 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 197 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

### Community 198 - "interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 199 - "telemetryCollector"
Cohesion: 0.33
Nodes (5): collectDynamicHostTelemetry(), TestTelemetryCacheAvoidsRedundantWork(), hostTelemetrySnapshot, staticHostInfo, telemetryCollector

### Community 200 - "app/layout.tsx"
Cohesion: 0.24
Nodes (7): jetbrainsMono, metadata, plusJakartaSans, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 201 - "newTestHarness"
Cohesion: 0.13
Nodes (20): testHarness, File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools() (+12 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 208 - "Config"
Cohesion: 0.09
Nodes (22): Options, runMain(), Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), setupSecurityTest(), TestSecurityEventBusAutoAuditIngestion() (+14 more)

### Community 209 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

### Community 211 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "setupMediaTest"
Cohesion: 0.19
Nodes (11): New(), doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), newRepo() (+3 more)

### Community 215 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 216 - "ClearCloudSyncHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearCloudSyncHistoryResponse

### Community 219 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 220 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 221 - "Registry"
Cohesion: 0.13
Nodes (6): Registry, normalizeName(), TestRegistryStartAllRollsBackOnPartialFailure(), nameMinHeap, registryState, testHarness

### Community 223 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 224 - "net/http.Header"
Cohesion: 0.22
Nodes (3): memoryResponseWriter, copyResponseHeaders(), syncBuffer

### Community 228 - "media/module.go"
Cohesion: 0.29
Nodes (4): createComboReq, createProviderReq, updateComboReq, updateProviderReq

### Community 229 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 973 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **30 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `net/http/httptest.ResponseRecorder` to `GatewayInterceptor`, `TeamsConsole.tsx`, `AppShell.tsx`, `DB`, `buildApp`, `Server`?**
  _High betweenness centrality (0.091) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0727522306108442 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `ModelsRoutingConsole.tsx` to `SettingsConsole.tsx`, `Server`, `AppShell.tsx`, `MediaConsole.tsx`, `TokenSaverConsole.tsx`, `Label`, `ProvidersConsole.tsx`, `UsageConsole.tsx`, `newTestHarness`, `VouchersConsole.tsx`, `Config`, `CardTitle`, `scripts`, `TelegramConsole.tsx`, `newTestHarness`, `TeamsConsole.tsx`, `CliToolsConsole.tsx`, `crypto.SessionSigner`, `repository`, `PxpipeConsole.tsx`, `Badge`, `cloudsync/hooks.ts`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `AppShell.tsx` to `net/http/httptest.ResponseRecorder`?**
  _High betweenness centrality (0.064) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._