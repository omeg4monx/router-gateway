# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 426 files · ~419,503 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5186 nodes · 16776 edges · 242 communities (206 shown, 36 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1706 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1f7cd9f3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- server.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- newTestHarness
- SettingsConsole.tsx
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- io/fs.FS
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
- system/module.go
- steps
- web-check
- on
- Label
- Module
- ProvidersModals.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- JobInspector
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
- Runner
- go_pkg_strings
- VouchersConsole.tsx
- Config
- newTestHarness
- redact.go
- Server
- ProviderAdapter
- log/slog.Logger
- listSettingsWithQueryer
- Bus
- repository
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- time.Duration
- TeamsConsole.tsx
- Module
- go_pkg_encoding_json
- ModelOverride
- properties
- Module
- Load
- repo
- CliToolsConsole.tsx
- web/package.json
- collectMigrations
- paths
- TunnelConsole.tsx
- net/http.Request
- collectDynamicHostTelemetry
- settingDefinition
- Module
- buildApp
- crypto.SessionSigner
- repo
- .executeCandidateAttempt
- Module
- repository
- repo
- Module
- properties
- extractErrorMessage
- CanonicalChatResponse
- Badge
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- .ProvisionKey
- sync.Mutex
- properties
- post
- setupModelEditorTest
- WriteError
- Button
- .StatusCode
- Module
- stubSystemModule
- repo
- properties
- transform.go
- repository
- Module
- Config
- DB
- ModelDescriptor
- repository
- properties
- application/json
- stubSystem
- Module
- Module
- properties
- OMEG4MON Router
- properties
- CanonicalStreamChunk
- certManager
- repo
- Module
- stubDep
- properties
- resolve.js
- release-binaries
- steps
- stubModule
- listAliasesResponse
- stubModule
- stubModule
- evaluateAntigravityQuota
- certs.go
- Module
- properties
- 200
- stubCompleter
- properties
- time.Time
- cli/package.json
- Module
- net/url.URL
- schemas
- get
- cloudsync/hooks.ts
- post
- schemas
- stubSystemModule
- PreviewPluginsResponse
- headroomCache
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
- settingDefinition
- DeleteSkillResponse
- version
- stubDepModule
- stubDep
- parameters
- Config
- batcher
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- compressor.go
- docker-release
- scripts
- `@omeg4mon/router` CLI Launcher
- devDependencies
- setupTokenSaverTest
- context.Context
- Registry
- setupTestMitm
- parameters
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- parseJSONObjectOrWrap
- interceptor.go
- net/http.Header
- Config
- single-binary
- .executeBuiltinFetch
- verify-web
- app/layout.tsx
- smoke-test.sh
- telemetry_windows.go
- repository
- ClearCloudSyncHistoryResponse
- NewStaticUIHandler
- nameMinHeap
- tgTestEnv
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- put
- setupTunnelTest
- stubSystemModule
- patch
- Skill
- media/module.go
- applyPipeline
- NewClient
- repository
- stubModule
- stubModule
- clitools/module.go
- redact.go
- trackedResponseWriter
- HistoryItem
- nullableString
- jsonMarshalerError
- listCombosResponse
- listDisabledModelsResponse
- EventListResponse
- settingsResponse
- pricingListResponse

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 219 edges
2. `WriteError()` - 204 edges
3. `WriteJSON()` - 201 edges
4. `DB` - 103 edges
5. `buildApp()` - 84 edges
6. `Module` - 80 edges
7. `Bus` - 79 edges
8. `Button` - 78 edges
9. `Module` - 63 edges
10. `Alert` - 63 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
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

## Communities (242 total, 36 thin omitted)

### Community 0 - "server.go"
Cohesion: 0.12
Nodes (45): ModelOption, patchConfigRequest, ChangedEvent, createComboRequest, deleteDisabledModelRequest, patchComboRequest, upsertAliasRequest, upsertDisabledModelRequest (+37 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
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

### Community 6 - "newTestHarness"
Cohesion: 0.19
Nodes (15): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Options, New(), decodeJSON(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain() (+7 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.06
Nodes (61): @hookform/resolvers, openapi-fetch, openapi-fetch@^0.14.0, DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+53 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.16
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.07
Nodes (70): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Functional & UX Scope, next, react, nextConfig, next@^16.4.6, DashboardLayout() (+62 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (123): stubAdminAuth, testHarness, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+115 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "io/fs.FS"
Cohesion: 0.08
Nodes (3): actorFromContext(), decodeJSONBody(), Module

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
Cohesion: 0.18
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (35): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+27 more)

### Community 21 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "system/module.go"
Cohesion: 0.09
Nodes (22): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+14 more)

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
Nodes (47): ArenaPage(), ChatPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels() (+39 more)

### Community 29 - "Module"
Cohesion: 0.15
Nodes (10): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+2 more)

### Community 30 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (148): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, DashboardProvidersPage(), ADAPTER_TO_REGISTRY_ID (+140 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "JobInspector"
Cohesion: 0.12
Nodes (23): File Map, Global Constraints, Review Focus, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan (+15 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.12
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamDecoder, streamEncoder (+9 more)

### Community 39 - "net/http.Request"
Cohesion: 0.12
Nodes (15): Module, clearSessionCookie(), isSecureRequest(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+7 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.07
Nodes (56): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), DashboardPage(), DashboardUsagePage(), useProviderConnections(), useProxyPools(), useSystemTelemetry(), formatBucketLabel(), formatBytesGB() (+48 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.15
Nodes (19): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+11 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "UpstreamCredentials"
Cohesion: 0.07
Nodes (20): fakeProviderStore, resolveSnapshot, stubProviderResolver, stubProviderStore, stubProviderStore, stubResolver, testEnv, ProviderConnection (+12 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "Module"
Cohesion: 0.10
Nodes (7): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (55): content, description, content, description, content, description, get, patch (+47 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.21
Nodes (17): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+9 more)

### Community 49 - "ComboSnapshot"
Cohesion: 0.14
Nodes (20): fakeComboStore, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification (+12 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (59): Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation() (+51 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.14
Nodes (25): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+17 more)

### Community 53 - "Module"
Cohesion: 0.04
Nodes (50): Global Constraints, Approaches Considered, buildProtocolTestRequest(), cloneOAuthConfig(), cloneStringSlice(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol() (+42 more)

### Community 54 - "main.go"
Cohesion: 0.16
Nodes (11): targetKey, migrationFile, lookupPluginSpec(), collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), PluginSpec, PoolTestResponse (+3 more)

### Community 55 - "navigation.ts"
Cohesion: 0.09
Nodes (33): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY (+25 more)

### Community 56 - "Runner"
Cohesion: 0.18
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 57 - "go_pkg_strings"
Cohesion: 0.21
Nodes (4): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), ValidationError

### Community 58 - "VouchersConsole.tsx"
Cohesion: 0.16
Nodes (21): VouchersPage(), useClaimVoucherMutation(), useCreateVoucherMutation(), useDeleteVoucherMutation(), useUpdateVoucherMutation(), useVoucherRedemptions(), useVouchersList(), vouchersQueryKeys (+13 more)

### Community 59 - "Config"
Cohesion: 0.11
Nodes (33): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 60 - "newTestHarness"
Cohesion: 0.43
Nodes (8): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid()

### Community 61 - "redact.go"
Cohesion: 0.16
Nodes (11): mapBusEventToAudit(), payloadToMap(), buildWhere(), formatUTC(), newEventID(), normalizeSeverity(), scanAuditEvent(), AuditEvent (+3 more)

### Community 62 - "Server"
Cohesion: 0.08
Nodes (21): TestAppCloseTeardownOrderIsLIFO(), loadUI(), loadUI(), doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth(), TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping() (+13 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (20): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+12 more)

### Community 64 - "log/slog.Logger"
Cohesion: 0.08
Nodes (8): Module, validateSQLiteSource(), decodeJSON(), clientIP(), decodeStrictJSON(), Module, truncateRunes(), RouteMounts

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.22
Nodes (9): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer (+1 more)

### Community 66 - "Bus"
Cohesion: 0.08
Nodes (9): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+1 more)

### Community 67 - "repository"
Cohesion: 0.12
Nodes (20): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+12 more)

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

### Community 72 - "time.Duration"
Cohesion: 0.13
Nodes (7): stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore(), formatDuration()

### Community 73 - "TeamsConsole.tsx"
Cohesion: 0.12
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 74 - "Module"
Cohesion: 0.06
Nodes (47): main(), runMain(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment (+39 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.09
Nodes (17): providerIDContextKey, wrapSystemPrompt(), wrapText(), ToolDescriptor, PreviewResult, catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse (+9 more)

### Community 76 - "ModelOverride"
Cohesion: 0.24
Nodes (6): marshalOverrideJSON(), nullFloat(), nullInt(), scanOverride(), ModelOverride, scanner

### Community 77 - "properties"
Cohesion: 0.13
Nodes (15): additionalProperties, properties, required, type, type, type, type, created_at (+7 more)

### Community 78 - "Module"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Load"
Cohesion: 0.22
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "repo"
Cohesion: 0.15
Nodes (11): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+3 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.10
Nodes (38): DashboardCliToolsPage(), DashboardProviderDetailPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog() (+30 more)

### Community 82 - "web/package.json"
Cohesion: 0.08
Nodes (31): class-variance-authority, clsx, geist, lucide-react, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form (+23 more)

### Community 83 - "collectMigrations"
Cohesion: 0.15
Nodes (12): Combo, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues, cloneStringSlice() (+4 more)

### Community 84 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 85 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 86 - "net/http.Request"
Cohesion: 0.10
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), readSessionCookie(), decodeStrictJSON(), LogEntry (+3 more)

### Community 87 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (22): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), estimateDarwinUsedMem(), collectDynamicHostTelemetry(), collectStaticHostInfo(), readBoundedFile(), collectDynamicHostTelemetry() (+14 more)

### Community 88 - "settingDefinition"
Cohesion: 0.11
Nodes (13): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), decodeStrictJSON() (+5 more)

### Community 89 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 90 - "buildApp"
Cohesion: 0.07
Nodes (47): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+39 more)

### Community 91 - "crypto.SessionSigner"
Cohesion: 0.15
Nodes (22): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+14 more)

### Community 92 - "repo"
Cohesion: 0.17
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - ".executeCandidateAttempt"
Cohesion: 0.15
Nodes (14): attemptOutcome, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage(), isProtectedRequestHeader() (+6 more)

### Community 94 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 95 - "repository"
Cohesion: 0.12
Nodes (23): PluginsPage(), SkillsPage(), CreateSkillRequest, PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse (+15 more)

### Community 96 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 97 - "Module"
Cohesion: 0.06
Nodes (16): streamResponseBody(), newID(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord, Event, fakeHandle (+8 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "extractErrorMessage"
Cohesion: 0.08
Nodes (48): PxpipePage(), TeamsPage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats() (+40 more)

### Community 100 - "CanonicalChatResponse"
Cohesion: 0.14
Nodes (19): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString(), RewriteModelField() (+11 more)

### Community 101 - "Badge"
Cohesion: 0.19
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 102 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 103 - ".ProvisionKey"
Cohesion: 0.16
Nodes (11): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest (+3 more)

### Community 104 - "sync.Mutex"
Cohesion: 0.16
Nodes (6): fakeSettingsStore, stubUsageRecorder, stubUsageRecorder, stubSettings, UsageRecord, recordingUsage

### Community 105 - "properties"
Cohesion: 0.10
Nodes (22): items, type, type, additionalProperties, properties, required, type, items (+14 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (8): Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 108 - "WriteError"
Cohesion: 0.08
Nodes (15): decodeStrictJSON(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID() (+7 more)

### Community 109 - "Button"
Cohesion: 0.17
Nodes (15): ImporterPage(), ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation() (+7 more)

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.14
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

### Community 118 - "Config"
Cohesion: 0.07
Nodes (37): authHarness, fakeCredentialStore, Options, Claims, File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter (+29 more)

### Community 119 - "DB"
Cohesion: 0.05
Nodes (53): Options, Options, Options, Options, repository, New(), newRepo(), New() (+45 more)

### Community 120 - "ModelDescriptor"
Cohesion: 0.36
Nodes (5): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, stubModelLister

### Community 121 - "repository"
Cohesion: 0.19
Nodes (10): boolToInt(), cloneMap(), mergeDefaults(), sanitizeStringSlice(), scanSkill(), PluginState, repository, scanner (+2 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.18
Nodes (12): content, content, description, content, content, description, schema, application/json (+4 more)

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

### Community 131 - "certManager"
Cohesion: 0.15
Nodes (7): randomSerial(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

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

### Community 143 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 144 - "certs.go"
Cohesion: 0.18
Nodes (5): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), writeAtomicFile()

### Community 145 - "Module"
Cohesion: 0.15
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - "stubCompleter"
Cohesion: 0.29
Nodes (7): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "time.Time"
Cohesion: 0.13
Nodes (12): CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock(), newDeviceOAuthRuntime(), newFakeClock() (+4 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "Module"
Cohesion: 0.07
Nodes (35): File Map, Review Focus, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+27 more)

### Community 153 - "net/url.URL"
Cohesion: 0.23
Nodes (8): stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL

### Community 154 - "schemas"
Cohesion: 0.25
Nodes (8): additionalProperties, required, type, schemas, CloudSyncConfig, UpdateCloudSyncConfigRequest, additionalProperties, type

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "cloudsync/hooks.ts"
Cohesion: 0.20
Nodes (17): CloudSyncPage(), CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest (+9 more)

### Community 157 - "post"
Cohesion: 0.27
Nodes (11): post, post, requestBody, /api/skills, operationId, requestBody, summary, tags (+3 more)

### Community 158 - "schemas"
Cohesion: 0.15
Nodes (13): components, schemas, additionalProperties, required, type, additionalProperties, required, type (+5 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 162 - "conflict_policy"
Cohesion: 0.20
Nodes (10): enum, type, conflict_policy, snapshot, additionalProperties, properties, required, type (+2 more)

### Community 163 - "TestContractInterfacesSatisfied"
Cohesion: 0.07
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), newID(), injectAnthropicCacheControl() (+1 more)

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

### Community 174 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Config"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 181 - "batcher"
Cohesion: 0.16
Nodes (4): newBatcher(), batcher, ModelPricing, pricingStore

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

### Community 190 - "devDependencies"
Cohesion: 0.11
Nodes (19): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+11 more)

### Community 191 - "setupTokenSaverTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (14): stubAdminCredentialStore, ImportRun, boolToInt(), newRepository(), sanitizeBotError(), Alias, AliasListResponse, Config (+6 more)

### Community 194 - "Registry"
Cohesion: 0.16
Nodes (7): testHarness, Registry, normalizeName(), TestRegistryStartAllRollsBackOnPartialFailure(), registryState, stubAdminAuth, testHarness

### Community 195 - "setupTestMitm"
Cohesion: 0.22
Nodes (13): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+5 more)

### Community 196 - "parameters"
Cohesion: 0.47
Nodes (9): parameters, parameters, in, name, required, schema, parameters, parameters (+1 more)

### Community 197 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 198 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 199 - "interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 200 - "net/http.Header"
Cohesion: 0.20
Nodes (4): memoryResponseWriter, copyResponseHeaders(), copyNonHopHeaders(), syncBuffer

### Community 201 - "Config"
Cohesion: 0.09
Nodes (22): Options, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Options, setupSecurityTest(), TestSecurityEventBusAutoAuditIngestion() (+14 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - ".executeBuiltinFetch"
Cohesion: 0.25
Nodes (9): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+1 more)

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): jetbrainsMono, metadata, plusJakartaSans, RootLayout(), RootLayoutProps, QueryProvider()

### Community 207 - "telemetry_windows.go"
Cohesion: 0.20
Nodes (5): collectStaticHostInfo(), runSysctl(), collectStaticHostInfo(), readWindowsRegistryString(), memoryStatusEx

### Community 208 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), scanBypass(), BypassRule, repository

### Community 209 - "ClearCloudSyncHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearCloudSyncHistoryResponse

### Community 210 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 212 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 215 - "put"
Cohesion: 0.40
Nodes (5): put, /api/plugins/{key}, operationId, summary, tags

### Community 216 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 219 - "patch"
Cohesion: 0.50
Nodes (4): patch, operationId, summary, tags

### Community 220 - "Skill"
Cohesion: 0.50
Nodes (4): Skill, additionalProperties, required, type

### Community 223 - "media/module.go"
Cohesion: 0.29
Nodes (4): createComboReq, createProviderReq, updateComboReq, updateProviderReq

### Community 224 - "applyPipeline"
Cohesion: 0.38
Nodes (5): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern()

### Community 225 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 226 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 229 - "clitools/module.go"
Cohesion: 0.33
Nodes (4): configRequest, applyResponse, catalogResponse, rollbackResponse

### Community 230 - "redact.go"
Cohesion: 0.67
Nodes (5): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes()

### Community 232 - "HistoryItem"
Cohesion: 0.67
Nodes (3): HistoryItem, historyListResp, historyResponse

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 999 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **36 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `Module`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11662174303683738 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `go_pkg_encoding_json`, `UpstreamCredentials`, `net/http.Request`, `ProvidersModals.tsx`?**
  _High betweenness centrality (0.163) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `extractErrorMessage` to `SettingsConsole.tsx`, `Server`, `AppShell.tsx`, `MediaConsole.tsx`, `TokenSaverConsole.tsx`, `cloudsync/hooks.ts`, `Label`, `ProvidersModals.tsx`, `UsageConsole.tsx`, `newTestHarness`, `ModelsRoutingConsole.tsx`, `VouchersConsole.tsx`, `Config`, `scripts`, `TelegramConsole.tsx`, `newTestHarness`, `TeamsConsole.tsx`, `CliToolsConsole.tsx`, `TunnelConsole.tsx`, `crypto.SessionSigner`, `repository`, `Badge`, `Button`?**
  _High betweenness centrality (0.066) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._