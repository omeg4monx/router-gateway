# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 434 files · ~466,431 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 6, .lock 1, .ico 1)

## Summary
- 5477 nodes · 18001 edges · 243 communities (215 shown, 28 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1761 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c7bb82dd`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- registry.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- VouchersConsole.tsx
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
- contracts.go
- steps
- web-check
- on
- io/fs.FS
- Module
- extractErrorMessage
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- ProviderDetailView.tsx
- ProvidersConsole.tsx
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
- AddApiKeyModal
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
- providers/hooks.ts
- ProviderAdapter
- GatewayInterceptor
- New
- Module
- promptinjection/engine.go
- scripts
- pricingStore
- TelegramConsole.tsx
- newTestHarness
- .AuthenticateAdmin
- TunnelConsole.tsx
- JobInspector
- go_pkg_encoding_json
- modeleditor/engine.go
- DB
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
- Module
- repo
- New
- WriteError
- Badge
- repo
- listSettingsWithQueryer
- properties
- Module
- @radix-ui/react-label
- .StatusCode
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- interceptor.go
- Module
- properties
- post
- stubSystemModule
- ComboSnapshot
- stubCompleter
- .executeCandidateAttempt
- Module
- combos/module_test.go
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
- ref_node_path
- release-binaries
- steps
- settingDefinition
- TokenSaverConsole.tsx
- omnirouteModelsParity.test.ts
- @tanstack/react-query
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- main_test.go
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- properties
- patch
- Module
- properties
- time.Time
- cli/package.json
- Module
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
- ref_bun_test
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
- KeyInfo
- `@omeg4mon/router` CLI Launcher
- headroomCache
- setupModelEditorTest
- context.Context
- client.ts
- repo
- parameters
- newTestHarness
- nameMinHeap
- Module
- New
- .testAwsQuickConnection
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
- setupModelEditorTest
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- ProviderLogo.tsx
- put
- auth/handlers.go
- NewStreamTranslator
- stubDep
- stubModule
- stubModule
- stubModule
- setupTestMitm
- stubDiscoverySystemModule
- stubSystemModule
- vite.config.ts
- jsonMarshalerError
- stubSystem
- applyConfigFile
- botAPIClient
- applyComboStrategy
- trackedResponseWriter
- Load
- net/http/httptest.ResponseRecorder
- `/dashboard` System Overview UI/UX Redesign (Approach A)
- repository
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
5. `Button` - 87 edges
6. `buildApp()` - 85 edges
7. `Bus` - 79 edges
8. `Alert` - 70 edges
9. `AlertDescription` - 69 edges
10. `Module` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification` --references--> `CliToolsConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → web/src/features/clitools/CliToolsConsole.tsx
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Intent & Why` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md → web/src/shared/ui/AppShell.tsx
- `Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`` --references--> `LocaleToggle()`  [INFERRED]
  docs/superpowers/plans/2026-11-24-login-ui-ux-redesign.md → web/src/shared/ui/LocaleToggle.tsx
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go

## Import Cycles
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
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
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/module_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 3-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (243 total, 28 thin omitted)

### Community 0 - "registry.go"
Cohesion: 0.09
Nodes (42): ModelOption, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, mapBusEventToAudit(), payloadToMap() (+34 more)

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
Cohesion: 0.09
Nodes (21): cloneCachedEntry(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary() (+13 more)

### Community 6 - "VouchersConsole.tsx"
Cohesion: 0.10
Nodes (36): zod, teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation(), useRemoveTeamKeyMutation(), useResetTeamSpendMutation(), useTeamsList() (+28 more)

### Community 7 - "SettingsConsole.tsx"
Cohesion: 0.10
Nodes (41): AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation(), useRollbackSettingMutation() (+33 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.11
Nodes (18): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+10 more)

### Community 9 - "Server"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.05
Nodes (98): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks (+90 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (101): TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest(), TestCloudSyncPushPullAndSnapshotConflictPolicies(), TestCloudSyncSSRFGuardAndVaultedTokenIsolation() (+93 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (32): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+24 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.05
Nodes (60): App(), ArenaConsole, ChatConsole, CliToolsConsole, CloudSyncStudio, CombosConsole, DashboardLayout(), HomePage (+52 more)

### Community 16 - "catalog.go"
Cohesion: 0.21
Nodes (17): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+9 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Module"
Cohesion: 0.14
Nodes (12): Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, stubKeyProvisioner, ProvisionedKey, ProvisionKeyRequest, actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt() (+4 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.11
Nodes (38): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ResponseFormat, SystemBlock, ToolChoice, ToolDefinition (+30 more)

### Community 21 - "net/http.Request"
Cohesion: 0.08
Nodes (23): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+15 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "contracts.go"
Cohesion: 0.18
Nodes (10): PolicyChangedEvent, patchConfigRequest, isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive(), createComboReq, createProviderReq, updateComboReq (+2 more)

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
Cohesion: 0.07
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 29 - "Module"
Cohesion: 0.16
Nodes (11): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+3 more)

### Community 30 - "extractErrorMessage"
Cohesion: 0.06
Nodes (64): CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle, RestoreCloudSyncSnapshotRequest, UpdateCloudSyncConfigRequest (+56 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "ProviderDetailView.tsx"
Cohesion: 0.09
Nodes (42): keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation(), useUpdateApiKeyMutation(), formatQuotaUsd() (+34 more)

### Community 35 - "ProvidersConsole.tsx"
Cohesion: 0.11
Nodes (36): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets(), FleetTelemetry, getEntrySortRank() (+28 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "Module"
Cohesion: 0.21
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.07
Nodes (56): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), useCombos(), useApiKeys(), useProviderConnections(), useProxyPools(), useSystemTelemetry(), formatBucketLabel(), formatBytesGB() (+48 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.14
Nodes (23): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+15 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "ModelCatalogConsole.tsx"
Cohesion: 0.08
Nodes (61): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update, AssignKeysModalProps, classifyTestError(), clearCachedTestResults() (+53 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "Module"
Cohesion: 0.17
Nodes (3): decodeJSONBody(), Module, stubKeyProvisioner

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (54): content, description, content, description, content, description, get, patch (+46 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "AddApiKeyModal"
Cohesion: 0.12
Nodes (29): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation() (+21 more)

### Community 50 - "newTestHarness"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.09
Nodes (27): COMBO_STRATEGIES, CreateComboFormInput, CreateComboFormOutput, createComboSchema, DEFAULT_CREATE_COMBO_VALUES, DEFAULT_DISABLE_MODEL_VALUES, DEFAULT_UPSERT_ALIAS_VALUES, DisableModelFormInput (+19 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.16
Nodes (22): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+14 more)

### Community 53 - ".handleOAuthExchange"
Cohesion: 0.09
Nodes (30): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), isValidProtocol(), validateHTTPURL() (+22 more)

### Community 54 - "main.go"
Cohesion: 0.16
Nodes (9): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), ToolDescriptor, PoolTestResponse, ScanCandidateResult, ScanResponse (+1 more)

### Community 55 - "overview.ts"
Cohesion: 0.21
Nodes (12): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+4 more)

### Community 56 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 57 - "main.go"
Cohesion: 0.19
Nodes (4): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), ValidationError

### Community 58 - "oauth_registry.go"
Cohesion: 0.18
Nodes (27): postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode(), defaultOAuthStrategies(), defaultRefreshFormToken(), enrichGeminiOAuth() (+19 more)

### Community 59 - "Badge"
Cohesion: 0.09
Nodes (42): ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail(), useChatSessions() (+34 more)

### Community 60 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (49): MODAL_SIZE_CLASS, ModalSize, CatalogOAuthDescriptor, commaListSchema(), connectionToEditFormValues(), CreateConnectionFormInput, CreateConnectionFormOutput, createConnectionSchema (+41 more)

### Community 61 - "Module"
Cohesion: 0.08
Nodes (9): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+1 more)

### Community 62 - "providers/hooks.ts"
Cohesion: 0.09
Nodes (27): DiscoveredModel, ExchangeProviderOAuthVariables, FetchConnectionModelsResponse, fetchProviderOAuthStatus(), PasteProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables (+19 more)

### Community 63 - "ProviderAdapter"
Cohesion: 0.11
Nodes (19): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+11 more)

### Community 64 - "GatewayInterceptor"
Cohesion: 0.05
Nodes (60): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway (+52 more)

### Community 65 - "New"
Cohesion: 0.08
Nodes (36): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+28 more)

### Community 67 - "promptinjection/engine.go"
Cohesion: 0.07
Nodes (34): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+26 more)

### Community 68 - "scripts"
Cohesion: 0.14
Nodes (28): modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides(), usePreviewModelEditorMutation() (+20 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 71 - "newTestHarness"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (6): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): react-hook-form, tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 74 - "JobInspector"
Cohesion: 0.15
Nodes (19): TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, JobInspector, New() (+11 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.07
Nodes (46): providerIDContextKey, TestAwsQuickEntitlementsProbeTestAndAutoRefresh(), TestAwsQuickTwoHopLoopbackOAuthAndRefresh(), TestAwsQuickTwoStepPasteBackFlow(), containsJSONStatus(), TestAllBuiltinOAuthAdaptersHaveClientIDConfiguredAndStartWithoutManualClientID(), TestFreePIDeviceCodeFlow(), TestGitHubCopilotDeviceCodeFlow() (+38 more)

### Community 76 - "modeleditor/engine.go"
Cohesion: 0.07
Nodes (26): memoryResponseWriter, newMemoryResponseWriter(), copyResponseHeaders(), copyNonHopHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp() (+18 more)

### Community 78 - "DB"
Cohesion: 0.05
Nodes (71): Options, Options, Options, Global Constraints, Review Focus, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Global Constraints (+63 more)

### Community 79 - "Load"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "devDependencies"
Cohesion: 0.13
Nodes (21): @biomejs/biome, openapi-typescript, tailwindcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+13 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): @hookform/resolvers, CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 83 - "collectMigrations"
Cohesion: 0.11
Nodes (20): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues (+12 more)

### Community 84 - "paths"
Cohesion: 0.50
Nodes (9): get, get, get, get, operationId, parameters, responses, summary (+1 more)

### Community 85 - "telemetry.go"
Cohesion: 0.05
Nodes (32): streamResponseBody(), clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry() (+24 more)

### Community 86 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

### Community 87 - "Module"
Cohesion: 0.15
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, check, check:fix, dev, format, format:check, gen:api (+4 more)

### Community 90 - "buildApp"
Cohesion: 0.06
Nodes (40): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd() (+32 more)

### Community 91 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 92 - "repo"
Cohesion: 0.17
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - "New"
Cohesion: 0.10
Nodes (22): loadUI(), loadUI(), Options, AuthVerifier, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement() (+14 more)

### Community 94 - "WriteError"
Cohesion: 0.06
Nodes (21): Module, restoreSnapshotRequest, SnapshotBundle, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID() (+13 more)

### Community 95 - "Badge"
Cohesion: 0.15
Nodes (25): class-variance-authority, ComboBuilderModal(), ComboBuilderModalProps, CombosConsole(), useCreateComboMutation(), useDeleteComboMutation(), useUpdateComboMutation(), Combo (+17 more)

### Community 96 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 97 - "listSettingsWithQueryer"
Cohesion: 0.22
Nodes (9): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer (+1 more)

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "Module"
Cohesion: 0.27
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 100 - "@radix-ui/react-label"
Cohesion: 0.33
Nodes (4): @fontsource-variable/plus-jakarta-sans, @radix-ui/react-label, @fontsource-variable/plus-jakarta-sans@^5.3.0, @radix-ui/react-label@^2.1.17

### Community 101 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 102 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 103 - "interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 104 - "Module"
Cohesion: 0.14
Nodes (8): parseAwsQuickPastedQuery(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickSession, awsQuickTokens, awsQuickUserInfo, Module

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 108 - "ComboSnapshot"
Cohesion: 0.20
Nodes (12): fakeComboStore, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot, listTables() (+4 more)

### Community 109 - "stubCompleter"
Cohesion: 0.33
Nodes (6): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent()

### Community 110 - ".executeCandidateAttempt"
Cohesion: 0.15
Nodes (17): attemptOutcome, CanonicalChatResponse, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage() (+9 more)

### Community 111 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 112 - "combos/module_test.go"
Cohesion: 0.10
Nodes (17): ChangedEvent, ComboTargetTestResult, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, patchComboRequest (+9 more)

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
Cohesion: 0.16
Nodes (21): translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS (+13 more)

### Community 118 - "Config"
Cohesion: 0.05
Nodes (47): authHarness, fakeCredentialStore, Options, Claims, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 2: `security` Module (Security Posture & Redacted Audit Trail) (+39 more)

### Community 119 - "OAuthConfig"
Cohesion: 0.08
Nodes (12): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), cloneStringSlice(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), xrouterBuiltInOAuthAdapters() (+4 more)

### Community 120 - "CanonicalRequest"
Cohesion: 0.06
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), injectAnthropicCacheControl(), computePrefixFingerprint()

### Community 121 - "repository"
Cohesion: 0.17
Nodes (12): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), sanitizeStringSlice(), scanSkill(), PluginSpec, PluginState (+4 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "application/json"
Cohesion: 0.13
Nodes (22): content, description, content, description, content, description, content, description (+14 more)

### Community 124 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 125 - "net/url.URL"
Cohesion: 0.09
Nodes (20): stubProxyDialerResolver, stubProxyResolver, isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+12 more)

### Community 126 - "oauthLoopbackProxy"
Cohesion: 0.14
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload & Dev Single Binary), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 130 - "deviceOAuthSession"
Cohesion: 0.24
Nodes (4): originOfURL(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 131 - "certManager"
Cohesion: 0.15
Nodes (7): randomSerial(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

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
Cohesion: 0.19
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 140 - "TokenSaverConsole.tsx"
Cohesion: 0.21
Nodes (17): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), parsePreviewBodyJson() (+9 more)

### Community 141 - "omnirouteModelsParity.test.ts"
Cohesion: 0.25
Nodes (7): appSource, comboBuilderSource, consoleSource, hooksSource, legacyConsolePath, modalSource, modelsDir

### Community 142 - "@tanstack/react-query"
Cohesion: 0.09
Nodes (23): clsx, @fontsource-variable/jetbrains-mono, openapi-fetch, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies (+15 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "main_test.go"
Cohesion: 0.09
Nodes (16): configRequest, generateResponse, applyResponse, catalogResponse, rollbackResponse, applyResponse, catalogResponse, rollbackResponse (+8 more)

### Community 145 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "patch"
Cohesion: 0.22
Nodes (9): delete, patch, operationId, summary, tags, operationId, summary, tags (+1 more)

### Community 148 - "Module"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "time.Time"
Cohesion: 0.14
Nodes (10): CooldownEntry, cooldownState, cooldownTracker, fakeClock, newCooldownTracker(), newFakeClock(), newFakeClock(), refreshCooldownEntry (+2 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 153 - "ModelDescriptor"
Cohesion: 0.22
Nodes (7): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, isModelAllowed(), isModelDisabled(), stubModelLister

### Community 154 - "schemas"
Cohesion: 0.15
Nodes (13): additionalProperties, required, type, additionalProperties, required, type, components, schemas (+5 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "compressor.go"
Cohesion: 0.56
Nodes (7): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput()

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

### Community 174 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 175 - "DeleteSkillResponse"
Cohesion: 0.22
Nodes (9): type, additionalProperties, properties, required, type, type, deleted, id (+1 more)

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "ref_bun_test"
Cohesion: 0.16
Nodes (9): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, detailSource, hooksSource, providersDir, modalsSource, providersDir (+1 more)

### Community 178 - "UpstreamCredentials"
Cohesion: 0.07
Nodes (24): catalogFakeProviderStore, fakeProviderStore, ModelTestResult, resolveSnapshot, testAllModelsResponse, weightedCandidate, stubProviderResolver, stubProviderStore (+16 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "time.Duration"
Cohesion: 0.11
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
Cohesion: 0.15
Nodes (7): stubUsageRecorder, stubUsageRecorder, testEnv, UsageRecord, newBatcher(), recordingUsage, batcher

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "KeyInfo"
Cohesion: 0.26
Nodes (7): stubKeyValidator, stubKeyValidator, KeyInfo, cloneKeyInfo(), hashToken(), isValidTokenFormat(), stubKeyValidator

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - "setupModelEditorTest"
Cohesion: 0.13
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (12): fakeSettingsStore, stubAdminCredentialStore, stubSettings, boolToInt(), Alias, AliasListResponse, Config, Event (+4 more)

### Community 194 - "client.ts"
Cohesion: 0.06
Nodes (44): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, @tanstack/react-query, ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse (+36 more)

### Community 195 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 196 - "parameters"
Cohesion: 0.46
Nodes (8): parameters, in, name, required, schema, parameters, parameters, type

### Community 197 - "newTestHarness"
Cohesion: 0.14
Nodes (19): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation() (+11 more)

### Community 199 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 200 - "New"
Cohesion: 0.13
Nodes (10): newAwsQuickRuntime(), newDeviceOAuthRuntime(), New(), enrichAntigravity(), TestEnrichAntigravityAutoOnboardsNewAccountAndSendsIDEFingerprint(), TestEnrichAntigravityFiltersBlockedProjectAndSetsHeader(), TestEnrichGitHubCopilotAndKiroProfileARN(), TestOAuthProviderRegistryMetadata() (+2 more)

### Community 201 - ".testAwsQuickConnection"
Cohesion: 0.24
Nodes (5): generateTraceID(), resolveAwsQuickTenantBaseURL(), extractCompanionProjectID(), oauthPlatformEnum(), stringClaim()

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

### Community 212 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (8): Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "ProviderLogo.tsx"
Cohesion: 0.24
Nodes (9): getProviderMonogram(), failedIconSrcs, PREFIX_ICON_ALIASES, PROVIDER_ICON_ALIASES, ProviderLogo(), resolveProviderIconSrc(), ProviderCard(), formatShortTime() (+1 more)

### Community 215 - "put"
Cohesion: 0.24
Nodes (10): put, requestBody, /api/plugins/{key}, requestBody, operationId, requestBody, summary, tags (+2 more)

### Community 216 - "auth/handlers.go"
Cohesion: 0.22
Nodes (8): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest

### Community 218 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 223 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 226 - "vite.config.ts"
Cohesion: 0.29
Nodes (6): @tailwindcss/vite, vite, @vitejs/plugin-react, @tailwindcss/vite@^4.3.5, vite@^8.3.8, @vitejs/plugin-react@^6.3.0

### Community 229 - "applyConfigFile"
Cohesion: 0.33
Nodes (5): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), writeAtomicFile()

### Community 231 - "applyComboStrategy"
Cohesion: 0.19
Nodes (8): CatalogModel, CatalogProviderGroup, modelCatalogResponse, strategyState, inferCatalogModel(), applyComboStrategy(), cryptoRandFloat64(), cryptoRandIntn()

### Community 233 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 234 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.08
Nodes (26): testHarness, stubAdminAuth, testHarness, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback() (+18 more)

### Community 235 - "`/dashboard` System Overview UI/UX Redesign (Approach A)"
Cohesion: 0.50
Nodes (3): `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes

### Community 236 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 255 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

## Knowledge Gaps
- **8 isolated node(s):** `Module`, `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan` (+3 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1058 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **28 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `AddApiKeyModal` to `.handleOAuthExchange`, `providers/hooks.ts`?**
  _High betweenness centrality (0.170) - this node is a cross-community bridge._
- **What connects `Module`, `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus` to the rest of the system?**
  _8 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `registry.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08551641071966275 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `.handleOAuthExchange` to `context.Context`, `go_pkg_encoding_json`, `AddApiKeyModal`, `UpstreamCredentials`, `net/http.Request`?**
  _High betweenness centrality (0.137) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `Config` to `GatewayInterceptor`, `newTestHarness`, `AppShell.tsx`, `Module`, `buildApp`, `New`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._