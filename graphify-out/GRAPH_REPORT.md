# Graph Report - wave3e-omeg4mon-router  (2026-11-23)

## Corpus Check
- 387 files · ~363,895 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 4768 nodes · 15492 edges · 193 communities (180 shown, 13 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1504 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1af9afa4`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_net_http
- Module
- biome.json
- settings
- Wave 1 — Core LLM Gateway & Essential Control Plane
- KeyInfo
- web/package.json
- .executeCandidateAttempt
- compilerOptions
- NotificationsConsole.tsx
- Button
- testing.T
- components.json
- TeamsConsole.tsx
- KeyProvisioner
- TokenSaverConsole
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- extractErrorMessage
- steps
- web-check
- on
- cross-compile
- Module
- ProvidersConsole.tsx
- schema.d.ts
- docker-build
- go-check
- jsonMarshalerError
- postcss.config.mjs
- github.com/omeg4monx/router-gateway
- CanonicalUsage
- net/http.Request
- UsageConsole.tsx
- translate_anthropic.go
- teeHandler
- UpstreamCredentials
- translate_gemini.go
- context.Context
- repository
- paths
- translate_openai.go
- Module
- PromptGuardConsole.tsx
- ModelsRoutingConsole.tsx
- encoding/json.RawMessage
- .handleCreateConnection
- main.go
- Module
- Module
- SecurityConsole.tsx
- repository
- ProxyPoolsConsole.tsx
- repository
- Runner
- Server
- Module
- Module
- listSettingsWithQueryer
- Module
- repository
- ModelEditorConsole.tsx
- Module
- TelegramConsole.tsx
- MitmConsole.tsx
- .AuthenticateAdmin
- scripts
- time.Time
- Module
- Module
- dependencies
- Module
- Module
- providers/module_test.go
- clitools/schemas.ts
- ProviderAdapter
- repository
- MediaConsole.tsx
- Module
- WriteJSON
- GatewayInterceptor
- Module
- net/url.URL
- buildApp
- TranslatorConsole.tsx
- repo
- app/layout.tsx
- File Map
- plugins/hooks.ts
- repo
- contracts.go
- properties
- PxpipeConsole.tsx
- properties
- Module
- client.ts
- io/fs.FS
- .Reader
- properties
- post
- setupModelEditorTest
- WriteError
- setupTokenSaverTest
- .StatusCode
- Module
- repository
- repo
- properties
- resolver.go
- net/http.Client
- system/module.go
- Config
- setupTestMitm
- setupMediaTest
- repository
- properties
- JobInspector
- stubSystemModule
- Module
- clitools/module_test.go
- properties
- promptinjection/engine.go
- updateConnectionTx
- CanonicalStreamChunk
- certManager
- repo
- newTestHarness
- Bus
- properties
- resolve.js
- repository
- Module
- media/module.go
- OAuthConfig
- stubCompleter
- Module
- transform.go
- newTestHarness
- sync.Mutex
- properties
- 200
- go_pkg_strings
- properties
- Module
- cli/package.json
- .executeBuiltinFetch
- combos/module.go
- schemas
- get
- paths
- Load
- schemas
- application/json
- PreviewPluginsResponse
- post
- conflict_policy
- compressor.go
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- DB
- orderStubModule
- tgTestEnv
- CloudSyncAliasSnapshot
- schemas
- properties
- RunLegacyImportResponse
- parameters
- DeleteSkillResponse
- zustand
- interceptor.go
- repo
- parameters
- Module
- ClearCloudSyncHistoryResponse
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- ProxyDialerResolver
- NewClient
- put
- `@omeg4mon/router` CLI Launcher
- patch
- Skill

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 214 edges
2. `WriteError()` - 195 edges
3. `WriteJSON()` - 191 edges
4. `DB` - 103 edges
5. `buildApp()` - 85 edges
6. `Bus` - 79 edges
7. `Button` - 73 edges
8. `Module` - 63 edges
9. `Alert` - 61 edges
10. `AlertTitle` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification` --references--> `CliToolsConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → web/src/features/clitools/CliToolsConsole.tsx
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

## Communities (193 total, 13 thin omitted)

### Community 0 - "go_pkg_net_http"
Cohesion: 0.10
Nodes (35): ModelOption, patchConfigRequest, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, mapBusEventToNotification() (+27 more)

### Community 1 - "Module"
Cohesion: 0.22
Nodes (5): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC()

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (47): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+39 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Wave 1 — Core LLM Gateway & Essential Control Plane"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "KeyInfo"
Cohesion: 0.07
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice() (+19 more)

### Community 6 - "web/package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (30): attemptOutcome, CanonicalChatResponse, applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage(), isProtectedRequestHeader(), parseRetryAfter(), CanFastPath() (+22 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "NotificationsConsole.tsx"
Cohesion: 0.16
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "Button"
Cohesion: 0.08
Nodes (92): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, class-variance-authority, lucide-react (+84 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (93): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest(), TestCloudSyncPushPullAndSnapshotConflictPolicies() (+85 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TeamsConsole.tsx"
Cohesion: 0.07
Nodes (41): @hookform/resolvers, zod, zod@^4.6.7, TeamsPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys() (+33 more)

### Community 14 - "KeyProvisioner"
Cohesion: 0.23
Nodes (10): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, TeamAssigner (+2 more)

### Community 15 - "TokenSaverConsole"
Cohesion: 0.15
Nodes (18): TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation() (+10 more)

### Community 16 - "catalog.go"
Cohesion: 0.21
Nodes (17): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+9 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.09
Nodes (21): testHarness, decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+13 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.17
Nodes (29): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+21 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.08
Nodes (46): DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation() (+38 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "extractErrorMessage"
Cohesion: 0.04
Nodes (75): CloudSyncPage(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useChatSessionDetail(), useChatSessions(), useCreateChatSessionMutation(), useDeleteArenaBattleMutation() (+67 more)

### Community 25 - "steps"
Cohesion: 0.50
Nodes (8): steps, steps, steps, steps, name, run, uses, steps

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "cross-compile"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, cross-compile

### Community 29 - "Module"
Cohesion: 0.16
Nodes (10): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+2 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.04
Nodes (83): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderConnectionMutation() (+75 more)

### Community 31 - "schema.d.ts"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (18): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+10 more)

### Community 39 - "net/http.Request"
Cohesion: 0.05
Nodes (33): trackedResponseWriter, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey() (+25 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (48): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation() (+40 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.14
Nodes (22): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+14 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "UpstreamCredentials"
Cohesion: 0.14
Nodes (10): fakeProviderStore, stubProviderStore, stubProviderStore, ProviderConnection, UpstreamCredentials, newFakeProviderStore(), cloneUpstreamCredentials(), sanitizeUpstreamError() (+2 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (23): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+15 more)

### Community 45 - "context.Context"
Cohesion: 0.04
Nodes (14): authHarness, fakeCredentialStore, stubAdminCredentialStore, stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent (+6 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "paths"
Cohesion: 0.05
Nodes (54): content, description, content, description, content, description, get, patch (+46 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.19
Nodes (19): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+11 more)

### Community 49 - "Module"
Cohesion: 0.09
Nodes (25): fakeComboStore, restoreSnapshotRequest, SnapshotBundle, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests (+17 more)

### Community 50 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (59): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+51 more)

### Community 52 - "encoding/json.RawMessage"
Cohesion: 0.15
Nodes (25): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+17 more)

### Community 53 - ".handleCreateConnection"
Cohesion: 0.17
Nodes (14): isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), truncateError(), validateHeaders(), validateName() (+6 more)

### Community 54 - "main.go"
Cohesion: 0.19
Nodes (4): lookupPluginSpec(), ToolDescriptor, PluginSpec, storedSetting

### Community 55 - "Module"
Cohesion: 0.16
Nodes (6): decodeStrictJSON(), New(), LogEntry, LogListResponse, Module, Status

### Community 56 - "Module"
Cohesion: 0.19
Nodes (7): compileRule(), decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), compiledRule, Module

### Community 57 - "SecurityConsole.tsx"
Cohesion: 0.20
Nodes (17): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+9 more)

### Community 58 - "repository"
Cohesion: 0.16
Nodes (8): decodeStringMap(), decodeStringSlice(), scanConnectionRecord(), scanNode(), oauthSessionRecord, ProviderNode, repository, rowScanner

### Community 59 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "repository"
Cohesion: 0.05
Nodes (35): memoryResponseWriter, newMemoryResponseWriter(), copyResponseHeaders(), copyNonHopHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp() (+27 more)

### Community 61 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 62 - "Server"
Cohesion: 0.07
Nodes (27): TestAppCloseTeardownOrderIsLIFO(), TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), Server (+19 more)

### Community 63 - "Module"
Cohesion: 0.18
Nodes (7): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken()

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.11
Nodes (15): listSettingsWithQueryer(), decodeStoredValue(), SettingItem, int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds(), applyResult (+7 more)

### Community 67 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 68 - "ModelEditorConsole.tsx"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 70 - "TelegramConsole.tsx"
Cohesion: 0.12
Nodes (31): TelegramPage(), telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation() (+23 more)

### Community 71 - "MitmConsole.tsx"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (9): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth (+1 more)

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.06
Nodes (18): CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker (+10 more)

### Community 75 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), New(), toJobDTO(), JobStatusDTO

### Community 76 - "Module"
Cohesion: 0.07
Nodes (11): streamResponseBody(), BinariesStatus, BinaryStatus, fakeHandle, fakeRunner, Module, osProcessHandle, osProcessRunner (+3 more)

### Community 77 - "dependencies"
Cohesion: 0.10
Nodes (20): clsx, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, dependencies, class-variance-authority@^0.7.1 (+12 more)

### Community 78 - "Module"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Module"
Cohesion: 0.20
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "providers/module_test.go"
Cohesion: 0.10
Nodes (12): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+4 more)

### Community 81 - "clitools/schemas.ts"
Cohesion: 0.12
Nodes (23): CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation(), useRollbackCliToolConfigMutation(), ApplyCLIToolConfigResponse, CLI_TOOL_IDS (+15 more)

### Community 82 - "ProviderAdapter"
Cohesion: 0.11
Nodes (4): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), customTestAdapter, ProviderAdapter

### Community 83 - "repository"
Cohesion: 0.15
Nodes (12): Combo, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues, cloneStringSlice() (+4 more)

### Community 84 - "MediaConsole.tsx"
Cohesion: 0.12
Nodes (31): MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders() (+23 more)

### Community 85 - "Module"
Cohesion: 0.18
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 86 - "WriteJSON"
Cohesion: 0.15
Nodes (7): Module, decodeJSON(), deriveSessionTitle(), New(), nowTimestamp(), randomHexID(), WriteJSON()

### Community 87 - "GatewayInterceptor"
Cohesion: 0.15
Nodes (23): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Review Focus (+15 more)

### Community 88 - "Module"
Cohesion: 0.14
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "net/url.URL"
Cohesion: 0.23
Nodes (8): stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL

### Community 90 - "buildApp"
Cohesion: 0.05
Nodes (48): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting(), TestWave2AE2ETranslatorAndCatalog(), TestWave2BDisabledModulesSkipInterceptors() (+40 more)

### Community 91 - "TranslatorConsole.tsx"
Cohesion: 0.14
Nodes (21): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+13 more)

### Community 92 - "repo"
Cohesion: 0.08
Nodes (21): buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow(), formatUTC() (+13 more)

### Community 93 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider()

### Community 94 - "File Map"
Cohesion: 0.19
Nodes (11): File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator, cloneAttrMap() (+3 more)

### Community 95 - "plugins/hooks.ts"
Cohesion: 0.12
Nodes (23): PluginsPage(), SkillsPage(), CreateSkillRequest, PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse (+15 more)

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 97 - "contracts.go"
Cohesion: 0.24
Nodes (8): PolicyChangedEvent, AuthVerifier, HealthCheck, HealthReport, patchConfigRequest, previewRequest, PreviewResponse, Stats

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "PxpipeConsole.tsx"
Cohesion: 0.12
Nodes (31): PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats(), usePxpipeStatus() (+23 more)

### Community 100 - "properties"
Cohesion: 0.08
Nodes (26): type, additionalProperties, properties, required, type, properties, type, type (+18 more)

### Community 101 - "Module"
Cohesion: 0.19
Nodes (5): computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 102 - "client.ts"
Cohesion: 0.04
Nodes (57): @tanstack/react-query, ImporterPage(), authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, AuthRedirectTarget (+49 more)

### Community 103 - "io/fs.FS"
Cohesion: 0.02
Nodes (14): stubDep, stubModule, stubModule, migrationFile, stubModule, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion() (+6 more)

### Community 104 - ".Reader"
Cohesion: 0.07
Nodes (36): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints (+28 more)

### Community 105 - "properties"
Cohesion: 0.11
Nodes (20): items, type, additionalProperties, properties, required, type, items, type (+12 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (19): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupModelEditorTest() (+11 more)

### Community 108 - "WriteError"
Cohesion: 0.15
Nodes (7): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), WriteError(), Module

### Community 109 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 112 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "resolver.go"
Cohesion: 0.11
Nodes (17): stubModelLister, resolveSnapshot, targetKey, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, testEnv (+9 more)

### Community 116 - "net/http.Client"
Cohesion: 0.09
Nodes (14): Config, HistoryEvent, Options, repository, storedConfig, New(), boolToInt(), newRepository() (+6 more)

### Community 117 - "system/module.go"
Cohesion: 0.15
Nodes (13): File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), AdminCredentialStore, KeyValidator, TestKeysExpirationAndValidationErrors(), healthResponse (+5 more)

### Community 118 - "Config"
Cohesion: 0.06
Nodes (43): Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases) (+35 more)

### Community 119 - "setupTestMitm"
Cohesion: 0.31
Nodes (10): New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly(), TestMitmLiveTLSProxyInterceptAndAliasRewrite() (+2 more)

### Community 120 - "setupMediaTest"
Cohesion: 0.18
Nodes (12): New(), doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), newRepo() (+4 more)

### Community 121 - "repository"
Cohesion: 0.19
Nodes (10): boolToInt(), cloneMap(), mergeDefaults(), sanitizeStringSlice(), scanSkill(), PluginState, repository, scanner (+2 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "JobInspector"
Cohesion: 0.17
Nodes (16): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, JobInspector, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup() (+8 more)

### Community 125 - "Module"
Cohesion: 0.16
Nodes (6): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), New(), formatUTC()

### Community 126 - "clitools/module_test.go"
Cohesion: 0.09
Nodes (12): configRequest, generateResponse, applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), applyResponse, catalogResponse (+4 more)

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "promptinjection/engine.go"
Cohesion: 0.33
Nodes (11): actionRank(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity(), redactAndTruncate() (+3 more)

### Community 129 - "updateConnectionTx"
Cohesion: 0.25
Nodes (11): cloneStringMap(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), updateConnectionTx(), connectionRecord, ConnectionSummary (+3 more)

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.38
Nodes (4): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder, mapOpenAIFinishReasonToCanonical()

### Community 131 - "certManager"
Cohesion: 0.12
Nodes (9): newCertManager(), randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide (+1 more)

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "newTestHarness"
Cohesion: 0.11
Nodes (21): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysValidateAndImmediateCacheInvalidation() (+13 more)

### Community 134 - "Bus"
Cohesion: 0.09
Nodes (21): Options, decodeJSON(), Module, New(), New(), New(), New(), Bus (+13 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "resolve.js"
Cohesion: 0.19
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 137 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 139 - "media/module.go"
Cohesion: 0.08
Nodes (20): providerIDContextKey, validateSlugID(), createComboReq, createProviderReq, updateComboReq, updateProviderReq, catalogAdapterDescriptor, catalogOAuthDescriptor (+12 more)

### Community 140 - "OAuthConfig"
Cohesion: 0.19
Nodes (9): cloneOAuthConfig(), cloneStringSlice(), OAuthConfig, oauthEnvPrefix(), New(), buildAuthorizeURL(), newRepository(), builtinAdapter (+1 more)

### Community 141 - "stubCompleter"
Cohesion: 0.33
Nodes (6): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent()

### Community 143 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 144 - "newTestHarness"
Cohesion: 0.18
Nodes (14): Options, testHarness, TestAppWave1CUsageAndCLIToolsEndToEnd(), File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New() (+6 more)

### Community 145 - "sync.Mutex"
Cohesion: 0.06
Nodes (15): fakeSettingsStore, stubUsageRecorder, stubUsageRecorder, stubSettings, UsageRecord, newBatcher(), Module, Options (+7 more)

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 148 - "go_pkg_strings"
Cohesion: 0.18
Nodes (8): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), HistoryItem, PreviewResult, historyListResp, historyResponse, ValidationError

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "Module"
Cohesion: 0.19
Nodes (4): actorFromContext(), decodeJSONBody(), Module, AdminSubject()

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - ".executeBuiltinFetch"
Cohesion: 0.25
Nodes (9): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+1 more)

### Community 153 - "combos/module.go"
Cohesion: 0.15
Nodes (9): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, patchComboRequest, upsertAliasRequest (+1 more)

### Community 154 - "schemas"
Cohesion: 0.15
Nodes (13): additionalProperties, required, type, additionalProperties, required, type, components, schemas (+5 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 157 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

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

### Community 163 - "compressor.go"
Cohesion: 0.36
Nodes (9): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), leadingBackticks(), truncateToolOutput() (+1 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.20
Nodes (10): $ref, additionalProperties, properties, required, type, items, type, plugins (+2 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "DB"
Cohesion: 0.08
Nodes (22): ImportRun, repository, boolToInt(), newRepository(), newRepository(), boolToInt(), newRepository(), formatUTC() (+14 more)

### Community 169 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

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

### Community 177 - "interceptor.go"
Cohesion: 0.36
Nodes (6): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages()

### Community 178 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 181 - "ClearCloudSyncHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearCloudSyncHistoryResponse

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

### Community 186 - "ProxyDialerResolver"
Cohesion: 0.38
Nodes (5): ProxyDialerResolver, NewProviderProxyFunc(), ProviderIDFromContext(), TestProviderProxyFuncContextRouting(), dynamicProxyResolver

### Community 187 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

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

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 870 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `system/module.go` to `contracts.go`, `newTestHarness`, `.Reader`, `Button`, `Module`, `Config`, `Server`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_net_http` be split into smaller, more focused modules?**
  _Cohesion score 0.09506172839506173 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `extractErrorMessage` to `NotificationsConsole.tsx`, `Button`, `TeamsConsole.tsx`, `TokenSaverConsole`, `SettingsConsole.tsx`, `ProvidersConsole.tsx`, `UsageConsole.tsx`, `PromptGuardConsole.tsx`, `ModelsRoutingConsole.tsx`, `SecurityConsole.tsx`, `ProxyPoolsConsole.tsx`, `ModelEditorConsole.tsx`, `TelegramConsole.tsx`, `MitmConsole.tsx`, `clitools/schemas.ts`, `MediaConsole.tsx`, `TranslatorConsole.tsx`, `plugins/hooks.ts`, `PxpipeConsole.tsx`, `client.ts`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `system/module.go`?**
  _High betweenness centrality (0.067) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._