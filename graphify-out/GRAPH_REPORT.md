# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 461 files · ~465,604 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5500 nodes · 17718 edges · 256 communities (220 shown, 36 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1747 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5d0c159b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- registry.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- DB
- SettingsConsole.tsx
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- ProxyPoolsConsole.tsx
- MediaConsole.tsx
- TokenSaverConsole.tsx
- catalog.go
- Design System & Wave 0 Plan Docs
- Module
- updates
- ir.go
- net/http.ResponseWriter
- with
- jobs
- registry_test.go
- steps
- web-check
- on
- io/fs.FS
- Module
- extractErrorMessage
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- KeysConsole.tsx
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- Module
- UsageConsole.tsx
- encoding/json.RawMessage
- buildApp
- models/catalogUtils.ts
- translate_gemini.go
- Module
- repository
- get
- translate_openai.go
- Module
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- .handleOAuthExchange
- main.go
- navigation.ts
- Runner
- go_pkg_encoding_json
- oauth_registry.go
- Badge
- ProvidersModals.tsx
- Module
- ComboSnapshot
- ProviderAdapter
- GatewayInterceptor
- ProvidersConsole.tsx
- net/http.Request
- repository
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- .AuthenticateAdmin
- ProcessHandle
- JobInspector
- go_pkg_encoding_json
- repository
- newTestHarness
- Bus
- Load
- devDependencies
- CliToolsConsole.tsx
- VouchersConsole.tsx
- collectMigrations
- paths
- collectDynamicHostTelemetry
- Module
- Module
- settingDefinition
- PxpipeConsole.tsx
- buildApp
- AddApiKeyModal
- repo
- New
- Module
- Input
- repo
- listSettingsWithQueryer
- properties
- WriteJSON
- TunnelConsole.tsx
- providers/hooks.ts
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- awsQuickRuntime
- properties
- post
- stubSystemModule
- Module
- stubCompleter
- .executeCandidateAttempt
- Module
- contracts.go
- repo
- properties
- transform.go
- repository
- TranslatorConsole.tsx
- Config
- OAuthConfig
- CanonicalRequest
- repository
- properties
- application/json
- repository
- net/url.URL
- oauthLoopbackProxy
- properties
- OMEG4MON Router
- Module
- deviceOAuthSession
- certManager
- repo
- repository
- antigravity.go
- properties
- omnirouteModelsParity.test.ts
- release-binaries
- steps
- settingDefinition
- stubDep
- Module
- @tanstack/react-query
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- apply.go
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- 200
- .ProvisionKey
- properties
- time.Time
- cli/package.json
- New
- ModelDescriptor
- schemas
- get
- compressor.go
- post
- schemas
- patch
- PreviewPluginsResponse
- properties
- conflict_policy
- fakeHandle
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
- CanonicalStreamChunk
- DeleteSkillResponse
- version
- paths
- UpstreamCredentials
- parameters
- time.Duration
- pricingStore
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- UsageRecord
- docker-release
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- `@omeg4mon/router` CLI Launcher
- headroomCache
- setupModelEditorTest
- context.Context
- SecurityConsole.tsx
- repo
- parameters
- Module
- NewStreamTranslator
- Module
- setupTunnelTest
- single-binary
- parameters
- verify-web
- repo
- smoke-test.sh
- awsQuickRuntime
- stubModule
- ModelCatalogConsole.tsx
- stubSystemModule
- io.Reader
- telemetryCollector
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- net/http.Header
- put
- promptinjection/engine.go
- ProviderAdapter
- repository
- patch
- Skill
- security/repo.go
- setupTestMitm
- auth/handlers.go
- repository
- ModelCatalogConsole
- jsonMarshalerError
- models/hooks.ts
- modeleditor/engine.go
- redact.go
- applyComboStrategy
- setupTokenSaverTest
- Load
- newTestHarness
- repository
- NewStaticUIHandler
- stubModule
- stubModule
- media/module.go
- stubModule
- stubDiscoverySystemModule
- stubSystemModule
- stubDepModule
- stubSystem
- stubDep
- stubModule
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- tgTestEnv

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 226 edges
2. `WriteError()` - 211 edges
3. `WriteJSON()` - 208 edges
4. `DB` - 103 edges
5. `buildApp()` - 85 edges
6. `Button` - 85 edges
7. `Bus` - 79 edges
8. `Alert` - 68 edges
9. `Badge()` - 67 edges
10. `AlertDescription` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`)` --references--> `ModelCatalogConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-dashboard-models-omniroute.md → web/src/features/models/ModelCatalogConsole.tsx
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (256 total, 36 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.07
Nodes (48): ModelOption, configRequest, patchConfigRequest, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest (+40 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
Nodes (4): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID()

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
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+19 more)

### Community 6 - "DB"
Cohesion: 0.10
Nodes (31): main(), Global Constraints, Global Constraints, File Map, Global Constraints, Review Focus, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification (+23 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (43): zod, zod@^4.6.7, DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+35 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.16
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.08
Nodes (72): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Functional & UX Scope, class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout() (+64 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (136): TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest(), TestCloudSyncPushPullAndSnapshotConflictPolicies(), TestCloudSyncSSRFGuardAndVaultedTokenIsolation() (+128 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (32): DashboardProxyPoolsPage(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+24 more)

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
Cohesion: 0.11
Nodes (39): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ResponseFormat, SystemBlock, ToolChoice, ToolDefinition (+31 more)

### Community 21 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (12): generateInput, Module, trackedResponseWriter, decodeJSONBody(), Module, writeGenerateError(), anthropicErrorType(), geminiErrorStatus() (+4 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "registry_test.go"
Cohesion: 0.17
Nodes (8): PolicyChangedEvent, generateResponse, applyResponse, catalogResponse, rollbackResponse, isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive()

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
Cohesion: 0.09
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 29 - "Module"
Cohesion: 0.09
Nodes (19): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID() (+11 more)

### Community 30 - "extractErrorMessage"
Cohesion: 0.04
Nodes (87): openapi-fetch, openapi-fetch@^0.14.0, CloudSyncPage(), ImporterPage(), PluginsPage(), SkillsPage(), CloudSyncStudio(), CloudSyncActionResponse (+79 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "KeysConsole.tsx"
Cohesion: 0.07
Nodes (47): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, DashboardKeysPage(), TeamsPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys() (+39 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (18): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+10 more)

### Community 39 - "Module"
Cohesion: 0.21
Nodes (5): clientIP(), decodeStrictJSON(), Module, New(), truncateRunes()

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (60): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, DashboardPage(), DashboardUsagePage(), useProxyPools() (+52 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.16
Nodes (22): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+14 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "models/catalogUtils.ts"
Cohesion: 0.16
Nodes (20): AssignKeysModalProps, comboTestResultKey(), computeUpdatedAllowedModels(), computeVisibilityCounts(), countKeyAccess(), DEFAULT_COMBO_FILTERS, DEFAULT_MODEL_FILTERS, filterCombos() (+12 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "Module"
Cohesion: 0.09
Nodes (3): decodeJSONBody(), New(), Module

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "get"
Cohesion: 0.17
Nodes (17): description, get, delete, get, get, operationId, responses, summary (+9 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.20
Nodes (18): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+10 more)

### Community 49 - "Module"
Cohesion: 0.18
Nodes (7): restoreSnapshotRequest, SnapshotBundle, computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.09
Nodes (27): COMBO_STRATEGIES, CreateComboFormInput, CreateComboFormOutput, createComboSchema, DEFAULT_CREATE_COMBO_VALUES, DEFAULT_DISABLE_MODEL_VALUES, DEFAULT_UPSERT_ALIAS_VALUES, DisableModelFormInput (+19 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.15
Nodes (25): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+17 more)

### Community 53 - ".handleOAuthExchange"
Cohesion: 0.07
Nodes (33): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), cloneStringSlice(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol() (+25 more)

### Community 54 - "main.go"
Cohesion: 0.22
Nodes (3): ToolDescriptor, ChangedEvent, storedSetting

### Community 55 - "navigation.ts"
Cohesion: 0.09
Nodes (32): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata (+24 more)

### Community 56 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 57 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (12): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), HistoryItem, collectStaticHostInfo(), runSysctl(), collectStaticHostInfo(), readWindowsRegistryString() (+4 more)

### Community 58 - "oauth_registry.go"
Cohesion: 0.13
Nodes (34): postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies(), defaultRefreshFormToken(), enrichAntigravity() (+26 more)

### Community 59 - "Badge"
Cohesion: 0.07
Nodes (68): ArenaPage(), ChatPage(), DashboardProviderDetailPage(), DashboardTranslatorPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles() (+60 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (49): MODAL_SIZE_CLASS, ModalSize, CatalogOAuthDescriptor, commaListSchema(), connectionToEditFormValues(), CreateConnectionFormInput, CreateConnectionFormOutput, createConnectionSchema (+41 more)

### Community 61 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), New(), toJobDTO(), JobStatusDTO

### Community 62 - "ComboSnapshot"
Cohesion: 0.22
Nodes (11): fakeComboStore, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables(), readLegacyCombos() (+3 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (19): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+11 more)

### Community 64 - "GatewayInterceptor"
Cohesion: 0.09
Nodes (37): Options, Options, File Map, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain (+29 more)

### Community 65 - "ProvidersConsole.tsx"
Cohesion: 0.09
Nodes (44): DashboardProvidersPage(), ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets(), FleetTelemetry (+36 more)

### Community 66 - "net/http.Request"
Cohesion: 0.06
Nodes (21): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), extractAPIKey(), validateSQLiteSource(), decodeJSON() (+13 more)

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

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (8): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "ProcessHandle"
Cohesion: 0.23
Nodes (4): fakeRunner, osProcessRunner, ProcessHandle, ProcessRunner

### Community 74 - "JobInspector"
Cohesion: 0.13
Nodes (17): File Map, Review Focus, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune() (+9 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.08
Nodes (23): providerIDContextKey, buildDiscoveryCandidateURLs(), firstStringField(), parseDiscoveredModels(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor, catalogOAuthDescriptor (+15 more)

### Community 76 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 77 - "newTestHarness"
Cohesion: 0.21
Nodes (11): Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors() (+3 more)

### Community 78 - "Bus"
Cohesion: 0.07
Nodes (26): New(), New(), New(), clientIP(), decodeStrictJSON(), Module, isValidSeverity(), New() (+18 more)

### Community 79 - "Load"
Cohesion: 0.19
Nodes (7): validateModelPattern(), decodeJSONBody(), New(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "devDependencies"
Cohesion: 0.07
Nodes (37): @biomejs/biome, geist, openapi-typescript, react-dom, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node (+29 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "VouchersConsole.tsx"
Cohesion: 0.16
Nodes (21): VouchersPage(), useClaimVoucherMutation(), useCreateVoucherMutation(), useDeleteVoucherMutation(), useUpdateVoucherMutation(), useVoucherRedemptions(), useVouchersList(), vouchersQueryKeys (+13 more)

### Community 83 - "collectMigrations"
Cohesion: 0.11
Nodes (20): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues (+12 more)

### Community 84 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 85 - "collectDynamicHostTelemetry"
Cohesion: 0.18
Nodes (19): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), estimateDarwinUsedMem(), collectDynamicHostTelemetry(), collectStaticHostInfo(), readBoundedFile(), collectDynamicHostTelemetry() (+11 more)

### Community 86 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 87 - "Module"
Cohesion: 0.23
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "PxpipeConsole.tsx"
Cohesion: 0.12
Nodes (31): PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats(), usePxpipeStatus() (+23 more)

### Community 90 - "buildApp"
Cohesion: 0.05
Nodes (46): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+38 more)

### Community 91 - "AddApiKeyModal"
Cohesion: 0.11
Nodes (31): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation() (+23 more)

### Community 92 - "repo"
Cohesion: 0.16
Nodes (12): formatUTC(), generateCode(), newID(), newRepo(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput (+4 more)

### Community 93 - "New"
Cohesion: 0.07
Nodes (22): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), testHarness, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), Server (+14 more)

### Community 94 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 95 - "Input"
Cohesion: 0.14
Nodes (24): DashboardCombosPage(), ComboBuilderModal(), ComboBuilderModalProps, CombosConsole(), useCombos(), useCreateComboMutation(), useDeleteComboMutation(), useUpdateComboMutation() (+16 more)

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 97 - "listSettingsWithQueryer"
Cohesion: 0.22
Nodes (9): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer (+1 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "WriteJSON"
Cohesion: 0.11
Nodes (9): Module, decodeJSON(), deriveSessionTitle(), decodeStrictJSON(), WriteJSON(), LogEntry, LogListResponse, Module (+1 more)

### Community 100 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 101 - "providers/hooks.ts"
Cohesion: 0.10
Nodes (26): DiscoveredModel, ExchangeProviderOAuthVariables, FetchConnectionModelsResponse, fetchProviderOAuthStatus(), PasteProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables (+18 more)

### Community 102 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 104 - "awsQuickRuntime"
Cohesion: 0.20
Nodes (8): parseAwsQuickDiscovery(), parseAwsQuickPastedQuery(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 108 - "Module"
Cohesion: 0.10
Nodes (22): Options, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI) (+14 more)

### Community 109 - "stubCompleter"
Cohesion: 0.26
Nodes (8): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), TestGatewayCompleteChat()

### Community 110 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (17): attemptOutcome, CanonicalChatResponse, statusError, moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), parseRetryAfter(), EncodeResponse(), rawJSONToString() (+9 more)

### Community 111 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 112 - "contracts.go"
Cohesion: 0.08
Nodes (27): CatalogModel, CatalogProviderGroup, ChangedEvent, ComboTargetTestResult, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse (+19 more)

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
Cohesion: 0.24
Nodes (6): Config, HistoryEvent, repository, storedConfig, boolToInt(), newRepository()

### Community 117 - "TranslatorConsole.tsx"
Cohesion: 0.16
Nodes (20): translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS (+12 more)

### Community 118 - "Config"
Cohesion: 0.05
Nodes (50): authHarness, fakeCredentialStore, Options, Claims, File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment (+42 more)

### Community 119 - "OAuthConfig"
Cohesion: 0.11
Nodes (4): cloneOAuthConfig(), OAuthConfig, builtinAdapter, customTestAdapter

### Community 120 - "CanonicalRequest"
Cohesion: 0.10
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), injectAnthropicCacheControl(), computePrefixFingerprint()

### Community 121 - "repository"
Cohesion: 0.16
Nodes (13): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), newRepository(), sanitizeStringSlice(), scanSkill(), PluginSpec (+5 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.18
Nodes (12): content, content, description, content, content, description, schema, application/json (+4 more)

### Community 124 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 125 - "net/url.URL"
Cohesion: 0.09
Nodes (20): stubProxyDialerResolver, stubProxyResolver, isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+12 more)

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 130 - "deviceOAuthSession"
Cohesion: 0.22
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 134 - "antigravity.go"
Cohesion: 0.24
Nodes (15): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+7 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "omnirouteModelsParity.test.ts"
Cohesion: 0.06
Nodes (33): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+25 more)

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 142 - "@tanstack/react-query"
Cohesion: 0.07
Nodes (27): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-hook-form, tailwind-merge, @tanstack/react-query, zustand, dependencies (+19 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "apply.go"
Cohesion: 0.70
Nodes (4): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic()

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - ".ProvisionKey"
Cohesion: 0.33
Nodes (4): stubKeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "time.Time"
Cohesion: 0.09
Nodes (20): CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), New(), newFakeClock(), newRepository() (+12 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "New"
Cohesion: 0.06
Nodes (37): TestAppCloseTeardownOrderIsLIFO(), TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`) (+29 more)

### Community 153 - "ModelDescriptor"
Cohesion: 0.08
Nodes (19): stubModelLister, resolveSnapshot, weightedCandidate, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, ModelDescriptor (+11 more)

### Community 154 - "schemas"
Cohesion: 0.11
Nodes (19): additionalProperties, properties, required, type, additionalProperties, required, type, additionalProperties (+11 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "compressor.go"
Cohesion: 0.56
Nodes (7): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput()

### Community 157 - "post"
Cohesion: 0.27
Nodes (11): post, post, requestBody, /api/skills, operationId, requestBody, summary, tags (+3 more)

### Community 158 - "schemas"
Cohesion: 0.15
Nodes (13): components, schemas, additionalProperties, required, type, additionalProperties, required, type (+5 more)

### Community 159 - "patch"
Cohesion: 0.13
Nodes (16): content, content, content, description, patch, schema, application/json, operationId (+8 more)

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.17
Nodes (12): items, type, type, type, items, additionalProperties, properties, required (+4 more)

### Community 161 - "properties"
Cohesion: 0.08
Nodes (26): type, properties, additionalProperties, properties, required, type, type, type (+18 more)

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

### Community 167 - "repo"
Cohesion: 0.24
Nodes (7): formatUTC(), newRepo(), newTeamID(), createTeamInput, patchTeamInput, repo, Team

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

### Community 174 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "paths"
Cohesion: 0.21
Nodes (14): description, post, post, post, components, paths, /api/cloud-sync/pull, /api/cloud-sync/push (+6 more)

### Community 178 - "UpstreamCredentials"
Cohesion: 0.08
Nodes (17): catalogFakeProviderStore, fakeProviderStore, ModelTestResult, testAllModelsResponse, stubProviderStore, stubProviderStore, ProviderConnection, UpstreamCredentials (+9 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "time.Duration"
Cohesion: 0.13
Nodes (6): stubAccountHealthTracker, failureRecord, stubHealthTracker, testEnv, computeBackoff(), formatDuration()

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
Cohesion: 0.16
Nodes (6): stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.10
Nodes (21): File Map, Review Focus, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, Wave 2B — Gateway Interceptors (`promptinjection`, `modeleditor` / `modelbypass`, `tokensaver` / `headroom`) Implementation Plan, TestGatewayInterceptorMetadataPropagatesToResponseHeaders() (+13 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (11): fakeSettingsStore, stubAdminCredentialStore, stubSettings, newRepository(), newBotAPIClient(), sanitizeBotError(), orderStubModule, stubCredentialStore (+3 more)

### Community 194 - "SecurityConsole.tsx"
Cohesion: 0.20
Nodes (17): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+9 more)

### Community 195 - "repo"
Cohesion: 0.25
Nodes (7): New(), newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 196 - "parameters"
Cohesion: 0.47
Nodes (9): parameters, parameters, in, name, required, schema, parameters, parameters (+1 more)

### Community 197 - "Module"
Cohesion: 0.14
Nodes (12): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Module, New(), indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus() (+4 more)

### Community 199 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 201 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "repo"
Cohesion: 0.14
Nodes (12): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow() (+4 more)

### Community 207 - "awsQuickRuntime"
Cohesion: 0.25
Nodes (3): awsQuickRuntime, loopbackProviderServer, loopbackSessionEntry

### Community 209 - "ModelCatalogConsole.tsx"
Cohesion: 0.23
Nodes (14): formatRelativeTime(), TestResultCell(), BulkTestProgress, CatalogComboMember, CatalogComboRow, CatalogModelRow, CatalogTab, CatalogTestResult (+6 more)

### Community 212 - "telemetryCollector"
Cohesion: 0.70
Nodes (3): TestTelemetryCacheAvoidsRedundantWork(), staticHostInfo, telemetryCollector

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "net/http.Header"
Cohesion: 0.18
Nodes (5): memoryResponseWriter, copyResponseHeaders(), copyNonHopHeaders(), applyStaticCacheHeaders(), syncBuffer

### Community 215 - "put"
Cohesion: 0.40
Nodes (5): put, /api/plugins/{key}, operationId, summary, tags

### Community 216 - "promptinjection/engine.go"
Cohesion: 0.30
Nodes (13): actionRank(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity() (+5 more)

### Community 217 - "ProviderAdapter"
Cohesion: 0.33
Nodes (5): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), defaultBuiltInAdapters(), xrouterBuiltInOAuthAdapters(), ProviderAdapter

### Community 218 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 219 - "patch"
Cohesion: 0.50
Nodes (4): patch, operationId, summary, tags

### Community 220 - "Skill"
Cohesion: 0.50
Nodes (4): Skill, additionalProperties, required, type

### Community 222 - "security/repo.go"
Cohesion: 0.29
Nodes (8): buildWhere(), formatUTC(), newEventID(), newRepo(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 223 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 224 - "auth/handlers.go"
Cohesion: 0.22
Nodes (8): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest

### Community 225 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 226 - "ModelCatalogConsole"
Cohesion: 0.22
Nodes (12): DashboardModelsPage(), clearCachedTestResults(), filterModels(), getEnabledCapabilities(), loadCachedTestResults(), modelTestResultKey(), saveCachedTestResults(), sortModels() (+4 more)

### Community 228 - "models/hooks.ts"
Cohesion: 0.24
Nodes (11): classifyTestError(), flattenCatalogResponse(), normalizeTestStatus(), RawCatalogGroup, modelCatalogQueryKeys, runSingleComboTest(), runSingleModelTest(), useCatalogApiKeys() (+3 more)

### Community 229 - "modeleditor/engine.go"
Cohesion: 0.33
Nodes (8): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), wrapSystemPrompt(), wrapText(), PreviewResult

### Community 230 - "redact.go"
Cohesion: 0.29
Nodes (9): mapBusEventToAudit(), payloadToMap(), isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity() (+1 more)

### Community 231 - "applyComboStrategy"
Cohesion: 0.29
Nodes (4): strategyState, applyComboStrategy(), cryptoRandFloat64(), cryptoRandIntn()

### Community 232 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 233 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 234 - "newTestHarness"
Cohesion: 0.33
Nodes (8): testHarness, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback()

### Community 236 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 237 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 242 - "media/module.go"
Cohesion: 0.29
Nodes (4): createComboReq, createProviderReq, updateComboReq, updateProviderReq

### Community 254 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 255 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1074 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **36 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `AddApiKeyModal` to `providers/hooks.ts`, `.handleOAuthExchange`?**
  _High betweenness centrality (0.174) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07410004675081813 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `.handleOAuthExchange` to `context.Context`, `net/http.Request`, `go_pkg_encoding_json`, `UpstreamCredentials`, `AddApiKeyModal`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `extractErrorMessage` to `SettingsConsole.tsx`, `Server`, `AppShell.tsx`, `ProxyPoolsConsole.tsx`, `MediaConsole.tsx`, `TokenSaverConsole.tsx`, `KeysConsole.tsx`, `UsageConsole.tsx`, `newTestHarness`, `Badge`, `SecurityConsole.tsx`, `scripts`, `TelegramConsole.tsx`, `newTestHarness`, `CliToolsConsole.tsx`, `VouchersConsole.tsx`, `PxpipeConsole.tsx`, `AddApiKeyModal`, `Input`, `models/hooks.ts`, `providers/hooks.ts`, `TunnelConsole.tsx`, `TranslatorConsole.tsx`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._