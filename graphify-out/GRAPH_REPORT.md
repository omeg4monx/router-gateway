# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 418 files · ~400,588 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5069 nodes · 16266 edges · 227 communities (201 shown, 26 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1572 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `48f85427`
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
- go_pkg_sync
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
- providers/module_test.go
- steps
- web-check
- on
- Label
- Module
- providers/schemas.ts
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- UsageRecord
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- net/http.ResponseWriter
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
- sync.Mutex
- time.Duration
- VouchersConsole.tsx
- Config
- CardTitle
- redact.go
- Server
- ProviderAdapter
- Module
- listSettingsWithQueryer
- Bus
- repository
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- go_pkg_encoding_json
- TeamsConsole.tsx
- GatewayInterceptor
- go_pkg_time
- repository
- properties
- Module
- Load
- newTestHarness
- CliToolsConsole.tsx
- dependencies
- collectMigrations
- DB
- ProviderAdapter
- WriteJSON
- collectDynamicHostTelemetry
- settingDefinition
- WriteJSON
- buildApp
- crypto.SessionSigner
- repo
- .executeCandidateAttempt
- Module
- repository
- repo
- collectDynamicHostTelemetry
- properties
- extractErrorMessage
- stubDep
- Badge
- stubSystemModule
- Module
- io/fs.FS
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
- Module
- Module
- .ProvisionKey
- time.Duration
- repository
- properties
- JobInspector
- stubSystem
- Module
- Module
- properties
- OMEG4MON Router
- properties
- CanonicalStreamChunk
- certManager
- repo
- RouteMounts
- cloudsync/hooks.ts
- properties
- resolve.js
- release-binaries
- steps
- Module
- combos/module.go
- stubModule
- Module
- evaluateAntigravityQuota
- clitools/module_test.go
- pricingStore
- properties
- put
- resolver.go
- properties
- time.Time
- cli/package.json
- New
- net/url.URL
- schemas
- get
- setupMediaTest
- setupTestMitm
- schemas
- stubCompleter
- PreviewPluginsResponse
- properties
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
- net/http.Header
- DeleteSkillResponse
- version
- stubDepModule
- .executeBuiltinFetch
- parameters
- Config
- ModelDescriptor
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- compressor.go
- docker-release
- scripts
- `@omeg4mon/router` CLI Launcher
- responsesStreamEncoder
- setupTokenSaverTest
- context.Context
- headroomCache
- setupTunnelTest
- NewClient
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- interceptor.go
- redact.go
- app/layout.tsx
- .Reader
- single-binary
- net/http/httptest.ResponseRecorder
- verify-web
- stubDep
- smoke-test.sh
- stubModule
- Config
- ClearCloudSyncHistoryResponse
- stubModule
- ChatCompleter
- stubModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- applyPipeline
- geist
- stubSystemModule
- repository
- repository
- @radix-ui/react-label
- tailwind-merge
- promptinjection/engine.go
- tgTestEnv
- jsonMarshalerError

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 215 edges
2. `WriteError()` - 195 edges
3. `WriteJSON()` - 192 edges
4. `DB` - 103 edges
5. `buildApp()` - 84 edges
6. `Bus` - 79 edges
7. `Button` - 78 edges
8. `Module` - 63 edges
9. `Alert` - 63 edges
10. `AlertDescription` - 62 edges

## Surprising Connections (you probably didn't know these)
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx
- `Key Proportion & UX Fixes` --references--> `Open()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-overview-design.md → internal/platform/db/db.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (227 total, 26 thin omitted)

### Community 0 - "server.go"
Cohesion: 0.09
Nodes (37): ModelOption, configRequest, generateResponse, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest (+29 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
Nodes (6): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC(), testHarness

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
Nodes (26): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+18 more)

### Community 6 - "web/package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): zod, zod@^4.6.7, DashboardSettingsPage(), authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+36 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.16
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.08
Nodes (71): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Functional & UX Scope, class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout() (+63 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (97): testEnv, TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest(), TestCloudSyncPushPullAndSnapshotConflictPolicies() (+89 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "go_pkg_sync"
Cohesion: 0.12
Nodes (11): HealthCheck, HealthReport, patchConfigRequest, previewRequest, PreviewResponse, Stats, CompressionStats, HeadroomStats (+3 more)

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
Cohesion: 0.22
Nodes (7): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken()

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.12
Nodes (38): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ResponseFormat, SystemBlock, ToolChoice, ToolDefinition (+30 more)

### Community 21 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "providers/module_test.go"
Cohesion: 0.11
Nodes (16): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+8 more)

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
Nodes (51): ArenaPage(), ChatPage(), DashboardProviderDetailPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard() (+43 more)

### Community 29 - "Module"
Cohesion: 0.16
Nodes (10): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+2 more)

### Community 30 - "providers/schemas.ts"
Cohesion: 0.04
Nodes (123): DashboardProvidersPage(), ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets(), FleetTelemetry (+115 more)

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
Cohesion: 0.30
Nodes (4): stubUsageRecorder, stubUsageRecorder, UsageRecord, recordingUsage

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (21): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+13 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (14): Module, clearSessionCookie(), isSecureRequest(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+6 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (62): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, DashboardPage(), DashboardUsagePage(), useProviderConnections() (+54 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.14
Nodes (23): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+15 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "UpstreamCredentials"
Cohesion: 0.18
Nodes (6): stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall

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
Nodes (55): content, description, content, description, content, description, get, patch (+47 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "ComboSnapshot"
Cohesion: 0.22
Nodes (11): fakeComboStore, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables(), readLegacyCombos() (+3 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (59): Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation() (+51 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.15
Nodes (23): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+15 more)

### Community 53 - "Module"
Cohesion: 0.10
Nodes (22): buildProtocolTestRequest(), cloneStringSlice(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol(), oauthEnvPrefix(), validateHTTPURL(), decodeJSONBody() (+14 more)

### Community 54 - "main.go"
Cohesion: 0.20
Nodes (7): migrationFile, lookupPluginSpec(), collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), PluginSpec, storedSetting

### Community 55 - "navigation.ts"
Cohesion: 0.08
Nodes (34): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata (+26 more)

### Community 56 - "sync.Mutex"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 57 - "time.Duration"
Cohesion: 0.18
Nodes (3): stubAccountHealthTracker, failureRecord, stubHealthTracker

### Community 58 - "VouchersConsole.tsx"
Cohesion: 0.09
Nodes (38): TeamsPage(), VouchersPage(), teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation(), useRemoveTeamKeyMutation(), useResetTeamSpendMutation() (+30 more)

### Community 59 - "Config"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "CardTitle"
Cohesion: 0.13
Nodes (22): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+14 more)

### Community 61 - "redact.go"
Cohesion: 0.26
Nodes (7): buildWhere(), formatUTC(), newEventID(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 62 - "Server"
Cohesion: 0.08
Nodes (23): loadUI(), loadUI(), Options, trackedResponseWriter, AuthVerifier, apiTimeoutMiddleware(), Server, New() (+15 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (20): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+12 more)

### Community 64 - "Module"
Cohesion: 0.27
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.10
Nodes (18): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, SettingItem, int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds() (+10 more)

### Community 66 - "Bus"
Cohesion: 0.13
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO, Notification

### Community 67 - "repository"
Cohesion: 0.18
Nodes (7): builtinRules(), boolToInt(), defaultConfig(), Config, GuardEvent, repository, Rule

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

### Community 73 - "TeamsConsole.tsx"
Cohesion: 0.11
Nodes (30): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation() (+22 more)

### Community 74 - "GatewayInterceptor"
Cohesion: 0.31
Nodes (14): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI), Options (+6 more)

### Community 75 - "go_pkg_time"
Cohesion: 0.10
Nodes (16): PolicyChangedEvent, errorDisposition, providerIDContextKey, classifyStatus(), isClientPayloadErrorBody(), wrapSystemPrompt(), wrapText(), createComboReq (+8 more)

### Community 76 - "repository"
Cohesion: 0.17
Nodes (10): boolToInt(), marshalOverrideJSON(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule, ModelOverride (+2 more)

### Community 77 - "properties"
Cohesion: 0.13
Nodes (15): additionalProperties, properties, required, type, type, type, type, created_at (+7 more)

### Community 78 - "Module"
Cohesion: 0.22
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Load"
Cohesion: 0.20
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "newTestHarness"
Cohesion: 0.12
Nodes (23): TestAppWave1AAuthSettingsKeysEndToEnd(), Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness() (+15 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.14
Nodes (27): @hookform/resolvers, DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog() (+19 more)

### Community 82 - "dependencies"
Cohesion: 0.11
Nodes (18): clsx, openapi-fetch, @radix-ui/react-slot, react-dom, zustand, dependencies, class-variance-authority@^0.7.1, clsx@^2.1.1 (+10 more)

### Community 83 - "collectMigrations"
Cohesion: 0.15
Nodes (12): Combo, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues, cloneStringSlice() (+4 more)

### Community 84 - "DB"
Cohesion: 0.05
Nodes (68): Options, Options, Options, testHarness, Global Constraints, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Global Constraints, Wave 3B — Media Gateway (Audio/Images/Videos/Search/Fetch) & Tunnel Manager Implementation Plan (+60 more)

### Community 85 - "ProviderAdapter"
Cohesion: 0.10
Nodes (6): cloneOAuthConfig(), buildAuthorizeURL(), builtinAdapter, customTestAdapter, OAuthConfig, ProviderAdapter

### Community 86 - "WriteJSON"
Cohesion: 0.10
Nodes (10): readSessionCookie(), decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), LogEntry, LogListResponse (+2 more)

### Community 87 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (22): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), estimateDarwinUsedMem(), collectDynamicHostTelemetry(), collectStaticHostInfo(), readBoundedFile(), collectDynamicHostTelemetry() (+14 more)

### Community 88 - "settingDefinition"
Cohesion: 0.14
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "WriteJSON"
Cohesion: 0.11
Nodes (8): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeJSON(), Module, WriteJSON()

### Community 90 - "buildApp"
Cohesion: 0.05
Nodes (47): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+39 more)

### Community 91 - "crypto.SessionSigner"
Cohesion: 0.14
Nodes (23): react-hook-form, DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse (+15 more)

### Community 92 - "repo"
Cohesion: 0.08
Nodes (21): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), formatUTC() (+13 more)

### Community 93 - ".executeCandidateAttempt"
Cohesion: 0.12
Nodes (20): attemptOutcome, CanonicalChatResponse, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage() (+12 more)

### Community 94 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 95 - "repository"
Cohesion: 0.12
Nodes (23): PluginsPage(), SkillsPage(), CreateSkillRequest, PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse (+15 more)

### Community 96 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 97 - "collectDynamicHostTelemetry"
Cohesion: 0.05
Nodes (18): streamResponseBody(), New(), newID(), newRepo(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord (+10 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "extractErrorMessage"
Cohesion: 0.06
Nodes (61): @tanstack/react-query, ImporterPage(), PxpipePage(), ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest (+53 more)

### Community 101 - "Badge"
Cohesion: 0.19
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 103 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 104 - "io/fs.FS"
Cohesion: 0.07
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 105 - "properties"
Cohesion: 0.10
Nodes (22): items, type, type, additionalProperties, properties, required, type, items (+14 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+10 more)

### Community 108 - "Module"
Cohesion: 0.14
Nodes (9): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), WriteError(), jsonMethodNotAllowedHandler(), jsonNotFoundHandler() (+1 more)

### Community 109 - "providers/handlers.go"
Cohesion: 0.18
Nodes (11): catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, createConnectionRequest, createNodeRequest, listNodesResponse, oauthExchangeRequest, oauthStartRequest (+3 more)

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.13
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 113 - "repo"
Cohesion: 0.17
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

### Community 118 - "Module"
Cohesion: 0.14
Nodes (21): File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, AdminCredentialStore, New(), decodeErrorCode(), doRequest(), findCookie() (+13 more)

### Community 119 - ".ProvisionKey"
Cohesion: 0.16
Nodes (10): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, ProvisionedKey, ProvisionKeyRequest, TeamAssigner, generateKeyID() (+2 more)

### Community 120 - "time.Duration"
Cohesion: 0.12
Nodes (8): Options, Claims, chainStore, Module, computeBackoff(), newChainStore(), formatDuration(), SessionSigner

### Community 121 - "repository"
Cohesion: 0.19
Nodes (10): boolToInt(), cloneMap(), mergeDefaults(), sanitizeStringSlice(), scanSkill(), PluginState, repository, scanner (+2 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "JobInspector"
Cohesion: 0.11
Nodes (24): TestAppCloseTeardownOrderIsLIFO(), TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification (+16 more)

### Community 125 - "Module"
Cohesion: 0.18
Nodes (7): restoreSnapshotRequest, SnapshotBundle, computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "properties"
Cohesion: 0.15
Nodes (13): additionalProperties, properties, required, type, type, $ref, type, type (+5 more)

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 131 - "certManager"
Cohesion: 0.15
Nodes (7): randomSerial(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "RouteMounts"
Cohesion: 0.11
Nodes (8): validateSQLiteSource(), decodeJSON(), decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), RouteMounts, Module

### Community 134 - "cloudsync/hooks.ts"
Cohesion: 0.20
Nodes (17): CloudSyncPage(), CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest (+9 more)

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
Cohesion: 0.23
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 140 - "combos/module.go"
Cohesion: 0.15
Nodes (9): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, patchComboRequest, upsertAliasRequest (+1 more)

### Community 142 - "Module"
Cohesion: 0.15
Nodes (4): clientIP(), decodeStrictJSON(), Module, isValidSeverity()

### Community 143 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 144 - "clitools/module_test.go"
Cohesion: 0.09
Nodes (10): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), writeAtomicFile(), collectStaticHostInfo(), runSysctl(), collectStaticHostInfo() (+2 more)

### Community 145 - "pricingStore"
Cohesion: 0.14
Nodes (7): newBatcher(), New(), newPricingStore(), newRepository(), batcher, ModelPricing, pricingStore

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "put"
Cohesion: 0.05
Nodes (68): content, description, content, description, content, description, content, description (+60 more)

### Community 148 - "resolver.go"
Cohesion: 0.12
Nodes (15): fakeProviderStore, resolveSnapshot, targetKey, stubProviderResolver, stubResolver, ProviderConnection, ResolvedTarget, newFakeProviderStore() (+7 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, models, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "time.Time"
Cohesion: 0.18
Nodes (8): CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock(), refreshCooldownEntry, cacheEntry

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "New"
Cohesion: 0.07
Nodes (34): File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+26 more)

### Community 153 - "net/url.URL"
Cohesion: 0.13
Nodes (12): stubProxyDialerResolver, stubProxyResolver, copyNonHopHeaders(), matchToolForHost(), newID(), singleJoiningSlash(), parseAndValidateProxyURL(), probeProxy() (+4 more)

### Community 154 - "schemas"
Cohesion: 0.25
Nodes (8): additionalProperties, required, type, schemas, CloudSyncConfig, UpdateCloudSyncConfigRequest, additionalProperties, type

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "setupMediaTest"
Cohesion: 0.18
Nodes (12): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), isLinkUnsupported() (+4 more)

### Community 157 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 158 - "schemas"
Cohesion: 0.14
Nodes (13): components, schemas, additionalProperties, required, type, CreateSkillRequest, Skill, UpdateSkillRequest (+5 more)

### Community 159 - "stubCompleter"
Cohesion: 0.29
Nodes (7): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 161 - "properties"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, key, updated_at, PluginState (+1 more)

### Community 162 - "conflict_policy"
Cohesion: 0.20
Nodes (10): enum, type, conflict_policy, snapshot, additionalProperties, properties, required, type (+2 more)

### Community 163 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID(), injectAnthropicCacheControl() (+1 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, plugins, PluginListResponse

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

### Community 174 - "net/http.Header"
Cohesion: 0.20
Nodes (4): memoryResponseWriter, copyResponseHeaders(), applyStaticCacheHeaders(), syncBuffer

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 178 - ".executeBuiltinFetch"
Cohesion: 0.25
Nodes (9): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Config"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 181 - "ModelDescriptor"
Cohesion: 0.36
Nodes (5): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, stubModelLister

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
Cohesion: 0.20
Nodes (10): $ref, items, skills, SkillListResponse, additionalProperties, properties, required, type (+2 more)

### Community 186 - "compressor.go"
Cohesion: 0.56
Nodes (7): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput()

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 191 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (18): authHarness, fakeCredentialStore, fakeSettingsStore, stubAdminCredentialStore, stubSettings, TestSystemAdminCredentialStore(), newRepository(), sanitizeBotError() (+10 more)

### Community 195 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 196 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 197 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 198 - "interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 199 - "redact.go"
Cohesion: 0.46
Nodes (6): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 200 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): jetbrainsMono, metadata, plusJakartaSans, RootLayout(), RootLayoutProps, QueryProvider()

### Community 201 - ".Reader"
Cohesion: 0.07
Nodes (34): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, File Map (+26 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.20
Nodes (8): testHarness, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback()

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 208 - "Config"
Cohesion: 0.06
Nodes (33): runMain(), Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs) (+25 more)

### Community 209 - "ClearCloudSyncHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearCloudSyncHistoryResponse

### Community 211 - "ChatCompleter"
Cohesion: 0.43
Nodes (6): Global Constraints, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan, ChatCompleter

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 215 - "applyPipeline"
Cohesion: 0.38
Nodes (5): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern()

### Community 219 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 220 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 224 - "promptinjection/engine.go"
Cohesion: 0.30
Nodes (13): actionRank(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity() (+5 more)

### Community 225 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 985 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **26 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `Module` to `TeamsConsole.tsx`, `AppShell.tsx`, `GatewayInterceptor`, `.Reader`, `newTestHarness`, `Server`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09215686274509804 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `GatewayInterceptor` to `.Reader`, `UpstreamCredentials`, `ModelsRoutingConsole.tsx`, `resolver.go`, `ir.go`, `ModelDescriptor`, `ProviderAdapter`, `buildApp`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `ModelsRoutingConsole()` connect `ModelsRoutingConsole.tsx` to `Badge`, `UsageConsole.tsx`, `AppShell.tsx`, `CardTitle`, `VouchersConsole.tsx`, `Label`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._