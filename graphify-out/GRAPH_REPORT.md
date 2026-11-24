# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 432 files · ~466,006 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5462 nodes · 17897 edges · 238 communities (216 shown, 22 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1732 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6e2dfdbb`
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
- App.tsx
- catalog.go
- Design System & Wave 0 Plan Docs
- Module
- updates
- ir.go
- net/http.Request
- with
- jobs
- main_test.go
- steps
- web-check
- on
- io/fs.FS
- Module
- extractErrorMessage
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- keys/hooks.ts
- RouteMounts
- Go Module Root Definition
- CanonicalUsage
- Module
- UsageConsole.tsx
- encoding/json.RawMessage
- buildApp
- ModelCatalogConsole.tsx
- translate_gemini.go
- Module
- repository
- paths
- translate_openai.go
- Module
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- .handleOAuthExchange
- main.go
- overview.ts
- Runner
- main.go
- oauth_registry.go
- Badge
- ProvidersModals.tsx
- Module
- readLegacySQLite
- ProviderAdapter
- GatewayInterceptor
- File Map
- Module
- promptinjection/engine.go
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- .AuthenticateAdmin
- net/url.URL
- JobInspector
- go_pkg_encoding_json
- repository
- oauthTokenResponse
- Bus
- Load
- devDependencies
- CliToolsConsole.tsx
- orderStubModule
- collectMigrations
- paths
- telemetry.go
- Module
- Module
- settingDefinition
- scripts
- buildApp
- setupMediaTest
- repo
- New
- WriteError
- Button
- repo
- listSettingsWithQueryer
- properties
- net/http.ResponseWriter
- main.tsx
- .StatusCode
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- interceptor.go
- awsQuickRuntime
- properties
- post
- stubSystemModule
- Module
- stubCompleter
- .executeCandidateAttempt
- Module
- go_pkg_encoding_json
- repo
- properties
- transform.go
- repository
- TranslatorConsole.tsx
- net/http/httptest.ResponseRecorder
- OAuthConfig
- context.Context
- repository
- properties
- application/json
- repository
- .executeBuiltinFetch
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
- ref_node_path
- release-binaries
- steps
- settingDefinition
- newTestHarness
- omnirouteModelsParity.test.ts
- @tanstack/react-query
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- clitools/module_test.go
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- patch
- .ProvisionKey
- properties
- time.Time
- cli/package.json
- New
- ModelDescriptor
- schemas
- get
- compressor.go
- paths
- schemas
- ClearCloudSyncHistoryResponse
- PreviewPluginsResponse
- properties
- conflict_policy
- NewClient
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
- BypassRule
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
- runConcurrentScan
- `@omeg4mon/router` CLI Launcher
- headroomCache
- setupModelEditorTest
- Module
- SecurityConsole.tsx
- repo
- parameters
- newTestHarness
- nameMinHeap
- Module
- Module
- setupTunnelTest
- single-binary
- palette.test.ts
- verify-web
- repo
- smoke-test.sh
- .list
- mobileLayout.test.ts
- omnirouteFlow.test.ts
- stubSystemModule
- PluginState
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- net/http.Header
- put
- repository
- setupTestMitm
- repository
- jsonMarshalerError
- modeleditor/engine.go
- mapBusEventToAudit
- applyComboStrategy
- setupTokenSaverTest
- Load
- Config
- repository
- NewStaticUIHandler
- stubModule
- stubDepModule
- stubDep
- stubModule
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
9. `AlertDescription` - 67 edges
10. `Badge()` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification` --references--> `CliToolsConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → web/src/features/clitools/CliToolsConsole.tsx
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/module_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/oauth_flow_integration_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (238 total, 22 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.08
Nodes (53): ModelOption, configRequest, patchConfigRequest, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse (+45 more)

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
Nodes (28): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+20 more)

### Community 6 - "DB"
Cohesion: 0.07
Nodes (36): Options, Options, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, New(), newRepository(), New(), TestModuleMetadataAndStartValidation() (+28 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, zod, zod@^4.6.7, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+36 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "Server"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.05
Nodes (72): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, Functional & UX Scope, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, authKeys (+64 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (121): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard() (+113 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "ProxyPoolsConsole.tsx"
Cohesion: 0.12
Nodes (31): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation(), useUpsertProxyBindingMutation() (+23 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.10
Nodes (29): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+21 more)

### Community 15 - "App.tsx"
Cohesion: 0.06
Nodes (58): App(), ArenaConsole, ChatConsole, CliToolsConsole, CloudSyncStudio, CombosConsole, DashboardLayout(), HomePage (+50 more)

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
Cohesion: 0.15
Nodes (29): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+21 more)

### Community 21 - "net/http.Request"
Cohesion: 0.11
Nodes (15): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+7 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "main_test.go"
Cohesion: 0.09
Nodes (23): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+15 more)

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
Nodes (11): stubDep, stubModule, migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), stubModule (+3 more)

### Community 29 - "Module"
Cohesion: 0.16
Nodes (11): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+3 more)

### Community 30 - "extractErrorMessage"
Cohesion: 0.02
Nodes (142): @tanstack/react-query, CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest (+134 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "keys/hooks.ts"
Cohesion: 0.12
Nodes (21): keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation(), useUpdateApiKeyMutation(), ApiKeySecretResponse (+13 more)

### Community 35 - "RouteMounts"
Cohesion: 0.11
Nodes (7): validateSQLiteSource(), decodeJSON(), actorFromContext(), decodeJSONBody(), Module, AdminSubject(), RouteMounts

### Community 38 - "CanonicalUsage"
Cohesion: 0.12
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+9 more)

### Community 39 - "Module"
Cohesion: 0.21
Nodes (5): clientIP(), decodeStrictJSON(), Module, New(), truncateRunes()

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.06
Nodes (57): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), ComboBuilderModalProps, CombosConsole(), useCombos(), useCreateComboMutation(), useDeleteComboMutation(), useUpdateComboMutation(), Combo (+49 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.18
Nodes (20): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicThinking, anthropicTool (+12 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "ModelCatalogConsole.tsx"
Cohesion: 0.08
Nodes (63): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update, AssignKeysModal(), AssignKeysModalProps, classifyTestError() (+55 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.13
Nodes (25): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+17 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (54): content, description, content, description, content, description, get, patch (+46 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+13 more)

### Community 49 - "Module"
Cohesion: 0.16
Nodes (7): restoreSnapshotRequest, SnapshotBundle, computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.16
Nodes (23): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+15 more)

### Community 53 - ".handleOAuthExchange"
Cohesion: 0.06
Nodes (38): Global Constraints, Review Focus, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Task 5: Full Verification & Graphify Sync, Approaches Considered, buildProtocolTestRequest(), cloneOAuthConfig(), cloneStringMap() (+30 more)

### Community 54 - "main.go"
Cohesion: 0.26
Nodes (4): lookupPluginSpec(), ToolDescriptor, PluginSpec, storedSetting

### Community 55 - "overview.ts"
Cohesion: 0.24
Nodes (10): filterModules(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY, OverviewCopy (+2 more)

### Community 56 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 57 - "main.go"
Cohesion: 0.10
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 58 - "oauth_registry.go"
Cohesion: 0.18
Nodes (21): generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), decodeJWTPayload(), enrichAntigravity(), enrichGeminiOAuth(), enrichGitHubCopilot() (+13 more)

### Community 59 - "Badge"
Cohesion: 0.07
Nodes (38): arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteArenaBattleMutation(), useDeleteChatSessionMutation() (+30 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.03
Nodes (154): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), computeEntryStats(), computeFleetTelemetry() (+146 more)

### Community 61 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), New(), toJobDTO(), JobStatusDTO

### Community 62 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (18): decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository(), scanConnectionRecord() (+10 more)

### Community 64 - "GatewayInterceptor"
Cohesion: 0.06
Nodes (55): main(), TestAppWave1BGatewayEndToEnd(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints, Review Focus (+47 more)

### Community 65 - "File Map"
Cohesion: 0.18
Nodes (10): File Map, Review Focus, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator, WithProviderID(), cloneAttrMap() (+2 more)

### Community 66 - "Module"
Cohesion: 0.19
Nodes (3): decodeJSON(), Module, New()

### Community 67 - "promptinjection/engine.go"
Cohesion: 0.06
Nodes (34): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+26 more)

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "TelegramConsole.tsx"
Cohesion: 0.10
Nodes (29): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+21 more)

### Community 71 - "newTestHarness"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.07
Nodes (8): stubAdminAuth, stubAuthVerifier, Options, stubAuthVerifier, AuthVerifier, Module, stubAdminAuth, stubAdminAuth

### Community 73 - "net/url.URL"
Cohesion: 0.20
Nodes (9): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, probeProxy(), proxyDialAddress(), NewProviderProxyFunc(), ProviderIDFromContext(), TestProviderProxyFuncContextRouting() (+1 more)

### Community 74 - "JobInspector"
Cohesion: 0.16
Nodes (15): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict() (+7 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (20): providerIDContextKey, firstStringField(), parseDiscoveredModels(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, ChangedEvent (+12 more)

### Community 76 - "repository"
Cohesion: 0.19
Nodes (8): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanOverride(), ModelOverride, repository

### Community 77 - "oauthTokenResponse"
Cohesion: 0.22
Nodes (12): isAnthropicTokenURL(), defaultExchangePKCECode(), defaultOAuthStrategies(), defaultRefreshFormToken(), exchangeClaudeCode(), newRFC8628Strategy(), refreshAWSOIDCToken(), refreshClaudeToken() (+4 more)

### Community 78 - "Bus"
Cohesion: 0.06
Nodes (28): Options, Global Constraints, Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, TeamAssigner, clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+20 more)

### Community 79 - "Load"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "devDependencies"
Cohesion: 0.10
Nodes (27): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+19 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.12
Nodes (23): CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation(), useRollbackCliToolConfigMutation(), ApplyCLIToolConfigResponse, CLI_TOOL_IDS (+15 more)

### Community 82 - "orderStubModule"
Cohesion: 0.18
Nodes (3): sanitizeBotError(), orderStubModule, botAPIClient

### Community 83 - "collectMigrations"
Cohesion: 0.11
Nodes (20): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues (+12 more)

### Community 84 - "paths"
Cohesion: 0.50
Nodes (9): get, get, get, get, operationId, parameters, responses, summary (+1 more)

### Community 85 - "telemetry.go"
Cohesion: 0.06
Nodes (32): streamResponseBody(), clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry() (+24 more)

### Community 86 - "Module"
Cohesion: 0.11
Nodes (10): decodeStrictJSON(), copyNonHopHeaders(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL() (+2 more)

### Community 87 - "Module"
Cohesion: 0.15
Nodes (9): decodeStrictJSON(), New(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module (+1 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, check, check:fix, dev, format, format:check, gen:api (+4 more)

### Community 90 - "buildApp"
Cohesion: 0.06
Nodes (44): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting(), TestWave2AE2ETranslatorAndCatalog(), TestWave2BDisabledModulesSkipInterceptors() (+36 more)

### Community 91 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 92 - "repo"
Cohesion: 0.16
Nodes (12): formatUTC(), generateCode(), newID(), newRepo(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput (+4 more)

### Community 93 - "New"
Cohesion: 0.08
Nodes (23): loadUI(), loadUI(), testEnv, New(), doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth() (+15 more)

### Community 94 - "WriteError"
Cohesion: 0.07
Nodes (18): Module, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), decodeJSON(), deriveSessionTitle(), New(), nowTimestamp(), randomHexID(), providerDisplayName() (+10 more)

### Community 95 - "Button"
Cohesion: 0.13
Nodes (80): class-variance-authority, @hookform/resolvers, lucide-react, react, react-hook-form, react-router, ArenaConsole(), ChatConsole() (+72 more)

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 97 - "listSettingsWithQueryer"
Cohesion: 0.24
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (13): generateInput, trackedResponseWriter, decodeJSONBody(), writeGenerateError(), decodeStrictJSON(), decodeStrictJSON(), matchesLevel(), normalizeLevelFilter() (+5 more)

### Community 100 - "main.tsx"
Cohesion: 0.22
Nodes (7): @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, react-dom, @fontsource-variable/jetbrains-mono@^5.3.0, @fontsource-variable/plus-jakarta-sans@^5.3.0, react-dom@^19.1.0, rootEl

### Community 101 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 102 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 103 - "interceptor.go"
Cohesion: 0.43
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 104 - "awsQuickRuntime"
Cohesion: 0.15
Nodes (8): newAwsQuickRuntime(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 108 - "Module"
Cohesion: 0.13
Nodes (16): fakeComboStore, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification (+8 more)

### Community 109 - "stubCompleter"
Cohesion: 0.29
Nodes (7): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent()

### Community 110 - ".executeCandidateAttempt"
Cohesion: 0.14
Nodes (18): attemptOutcome, CanonicalChatResponse, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage() (+10 more)

### Community 111 - "Module"
Cohesion: 0.12
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 112 - "go_pkg_encoding_json"
Cohesion: 0.15
Nodes (11): ComboTargetTestResult, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest, testComboRequest, testComboResponse, testModelRequest (+3 more)

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

### Community 117 - "TranslatorConsole.tsx"
Cohesion: 0.13
Nodes (18): translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS (+10 more)

### Community 118 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.08
Nodes (31): Options, stubAdminAuth, testHarness, Claims, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest() (+23 more)

### Community 119 - "OAuthConfig"
Cohesion: 0.10
Nodes (7): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), defaultBuiltInAdapters(), xrouterBuiltInOAuthAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable(), customTestAdapter, ProviderAdapter

### Community 120 - "context.Context"
Cohesion: 0.05
Nodes (9): fakeSettingsStore, stubInterceptor, recordingInterceptor, stubSettings, CanonicalRequest, ResponseMeta, StreamEvent, newID() (+1 more)

### Community 121 - "repository"
Cohesion: 0.18
Nodes (11): boolToInt(), cloneMap(), mergeDefaults(), newRepository(), sanitizeStringSlice(), scanSkill(), PluginState, repository (+3 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.13
Nodes (22): content, description, content, description, content, description, content, description (+14 more)

### Community 124 - "repository"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 125 - ".executeBuiltinFetch"
Cohesion: 0.25
Nodes (9): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+1 more)

### Community 126 - "oauthLoopbackProxy"
Cohesion: 0.14
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

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
Cohesion: 0.16
Nodes (7): OAuthConfig, originOfURL(), buildAuthorizeURL(), builtinAdapter, deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

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
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "ref_node_path"
Cohesion: 0.12
Nodes (14): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+6 more)

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 140 - "newTestHarness"
Cohesion: 0.29
Nodes (8): newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify(), VerifyPassword(), testHarness

### Community 141 - "omnirouteModelsParity.test.ts"
Cohesion: 0.25
Nodes (7): appSource, comboBuilderSource, consoleSource, hooksSource, legacyConsolePath, modalSource, modelsDir

### Community 142 - "@tanstack/react-query"
Cohesion: 0.10
Nodes (20): clsx, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, tailwind-merge, zustand, dependencies, class-variance-authority@^0.7.1 (+12 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "clitools/module_test.go"
Cohesion: 0.24
Nodes (8): generateResponse, applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), applyResponse, catalogResponse, rollbackResponse

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "patch"
Cohesion: 0.22
Nodes (9): delete, patch, operationId, summary, tags, operationId, summary, tags (+1 more)

### Community 148 - ".ProvisionKey"
Cohesion: 0.20
Nodes (10): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest (+2 more)

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "time.Time"
Cohesion: 0.14
Nodes (10): CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock(), newFakeClock(), refreshCooldownEntry (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "New"
Cohesion: 0.15
Nodes (15): TestAppCloseTeardownOrderIsLIFO(), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, New(), newHarness(), TestTranslatorLogsSnapshotAndLiveSSEStream(), TestTranslatorLogsStreamSurvivesPastAPITimeout(), TestTranslatorTranslateRequestAndResponseAllPairs(), New() (+7 more)

### Community 153 - "ModelDescriptor"
Cohesion: 0.36
Nodes (5): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, stubModelLister

### Community 154 - "schemas"
Cohesion: 0.15
Nodes (13): additionalProperties, required, type, additionalProperties, required, type, components, schemas (+5 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "compressor.go"
Cohesion: 0.47
Nodes (8): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), leadingBackticks(), truncateToolOutput()

### Community 157 - "paths"
Cohesion: 0.21
Nodes (11): post, post, components, paths, /api/plugins, /api/plugins/preview, /api/skills, /api/skills/{id}/export (+3 more)

### Community 158 - "schemas"
Cohesion: 0.17
Nodes (12): schemas, additionalProperties, required, type, CreateSkillRequest, Skill, UpdateSkillRequest, additionalProperties (+4 more)

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

### Community 163 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

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
Cohesion: 0.16
Nodes (13): anthropicStreamDecoder, CanonicalStreamChunk, geminiStreamDecoder, openAIStreamDecoder, responsesStreamDecoder, streamDecoder, NewStreamTranslator(), newAnthropicStreamDecoder() (+5 more)

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "BypassRule"
Cohesion: 0.53
Nodes (3): scanBypass(), BypassRule, scanner

### Community 178 - "UpstreamCredentials"
Cohesion: 0.07
Nodes (25): catalogFakeProviderStore, fakeProviderStore, ModelTestResult, resolveSnapshot, testAllModelsResponse, weightedCandidate, stubProviderResolver, stubProviderStore (+17 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "time.Duration"
Cohesion: 0.12
Nodes (8): stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore(), formatDuration(), newHeadroomCache()

### Community 181 - "pricingStore"
Cohesion: 0.31
Nodes (3): newPricingStore(), ModelPricing, pricingStore

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
Cohesion: 0.13
Nodes (7): stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher, stubTeamAssigner

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "runConcurrentScan"
Cohesion: 0.40
Nodes (6): parseAndValidateProxyURL(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL, ScanCandidateResult, ScanResponse

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupModelEditorTest() (+11 more)

### Community 193 - "Module"
Cohesion: 0.06
Nodes (12): authHarness, fakeCredentialStore, stubAdminCredentialStore, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, AdminCredentialStore, Module, TestSystemAdminCredentialStore() (+4 more)

### Community 194 - "SecurityConsole.tsx"
Cohesion: 0.18
Nodes (14): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse, CreateSecurityAuditFormValues (+6 more)

### Community 195 - "repo"
Cohesion: 0.29
Nodes (6): newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 196 - "parameters"
Cohesion: 0.46
Nodes (8): parameters, in, name, required, schema, parameters, parameters, type

### Community 197 - "newTestHarness"
Cohesion: 0.14
Nodes (19): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation() (+11 more)

### Community 200 - "Module"
Cohesion: 0.15
Nodes (3): Module, deviceOAuthRuntime, refreshCall

### Community 201 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "repo"
Cohesion: 0.18
Nodes (9): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotificationRow(), EmitPayload, Notification (+1 more)

### Community 207 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 208 - "mobileLayout.test.ts"
Cohesion: 0.40
Nodes (4): consoleSource, detailSource, modalsSource, providersDir

### Community 209 - "omnirouteFlow.test.ts"
Cohesion: 0.40
Nodes (4): consoleSource, detailSource, modalsSource, providersDir

### Community 211 - "PluginState"
Cohesion: 0.50
Nodes (4): additionalProperties, required, type, PluginState

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 215 - "put"
Cohesion: 0.24
Nodes (10): put, requestBody, /api/plugins/{key}, requestBody, operationId, requestBody, summary, tags (+2 more)

### Community 218 - "repository"
Cohesion: 0.16
Nodes (5): Alias, AliasListResponse, Config, Event, repository

### Community 223 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 225 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 229 - "modeleditor/engine.go"
Cohesion: 0.28
Nodes (10): copyResponseHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt() (+2 more)

### Community 230 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 231 - "applyComboStrategy"
Cohesion: 0.19
Nodes (8): CatalogModel, CatalogProviderGroup, modelCatalogResponse, strategyState, inferCatalogModel(), applyComboStrategy(), cryptoRandFloat64(), cryptoRandIntn()

### Community 232 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 233 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 234 - "Config"
Cohesion: 0.07
Nodes (27): Options, testHarness, TestAppWave1CUsageAndCLIToolsEndToEnd(), File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Task 2: `security` Module (Security Posture & Redacted Audit Trail) (+19 more)

### Community 236 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 237 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 255 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1056 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **22 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProvidersModals.tsx` to `.handleOAuthExchange`?**
  _High betweenness centrality (0.138) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07780320366132723 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `.handleOAuthExchange` to `go_pkg_encoding_json`, `UpstreamCredentials`, `net/http.Request`, `context.Context`, `ProvidersModals.tsx`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `Module` to `GatewayInterceptor`, `newTestHarness`, `SettingsConsole.tsx`, `.AuthenticateAdmin`, `AppShell.tsx`, `net/http/httptest.ResponseRecorder`, `buildApp`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._