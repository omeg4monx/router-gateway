# Graph Report - omeg4mon-router  (2026-11-24)

## Corpus Check
- 394 files · ~369,607 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 4894 nodes · 15761 edges · 219 communities (195 shown, 24 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1540 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fce3268c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- package.json
- SettingsConsole.tsx
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- KeysConsole.tsx
- client.ts
- io/fs.FS
- catalog.go
- Design System & Wave 0 Plan Docs
- Bus
- updates
- ir.go
- .AuthenticateAdmin
- with
- jobs
- CardTitle
- steps
- web-check
- on
- cross-compile
- Module
- providers/schemas.ts
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- updateConnectionTx
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- net/http.ResponseWriter
- UsageConsole.tsx
- encoding/json.RawMessage
- buildApp
- UpstreamCredentials
- translate_gemini.go
- fakeCredentialStore
- repository
- context.Context
- translate_openai.go
- NewStaticUIHandler
- newTestHarness
- ModelsRoutingConsole.tsx
- translate_responses.go
- WriteError
- Module
- WriteJSON
- Module
- CanonicalStreamChunk
- repository
- Config
- Module
- NewStreamTranslator
- New
- Bus
- Module
- listSettingsWithQueryer
- recordingInterceptor
- newTestHarness
- scripts
- pricingStore
- responsesStreamEncoder
- newTestHarness
- Server
- translate_test.go
- time.Time
- DB
- stubModule
- sync.Mutex
- time.Duration
- Load
- providers/module_test.go
- CliToolsConsole.tsx
- stubModule
- collectMigrations
- dependencies
- Module
- Module
- contracts.go
- settingDefinition
- nameMinHeap
- buildApp
- crypto.SessionSigner
- net/http.Request
- .executeCandidateAttempt
- @radix-ui/react-label
- repository
- repo
- contracts_test.go
- properties
- extractErrorMessage
- properties
- Module
- TunnelConsole.tsx
- io/fs.FS
- .Reader
- properties
- post
- setupModelEditorTest
- Module
- antigravity.go
- .StatusCode
- Module
- repository
- repo
- properties
- ModelDescriptor
- repository
- properties
- Module
- repo
- ProvidersConsole.tsx
- repository
- properties
- JobInspector
- stubSystemModule
- Module
- apply.go
- properties
- OMEG4MON Router
- UpstreamCredentials
- CanonicalStreamChunk
- certManager
- repo
- DB
- Bus
- properties
- resolve.js
- release-binaries
- steps
- providers/handlers.go
- net/http.Client
- net/http.Header
- ProvidersConsole
- transform.go
- providers/hooks.ts
- pricingStore
- properties
- 200
- schema.go
- properties
- setupTokenSaverTest
- cli/package.json
- .executeBuiltinFetch
- Module
- schemas
- get
- paths
- compressor.go
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
- setupMediaTest
- Review Focus
- CloudSyncAliasSnapshot
- schemas
- properties
- RunLegacyImportResponse
- parameters
- DeleteSkillResponse
- version
- providers/schemas.test.ts
- repo
- parameters
- Config
- ClearCloudSyncHistoryResponse
- CloudSyncHistoryListResponse
- ClearImporterHistoryResponse
- source_path
- SkillListResponse
- stubDep
- docker-release
- put
- `@omeg4mon/router` CLI Launcher
- patch
- Skill
- media/module.go
- pxpipe/module.go
- repository
- stubSystemModule
- stubDepModule
- stubSystemModule
- stubSystem
- stubDep
- collectMigrations
- single-binary
- UsageRecord
- verify-web
- scripts
- smoke-test.sh
- modeleditor/engine.go
- Config
- app/layout.tsx
- stubModule
- stubModule
- stubModule
- mapBusEventToAudit
- nameMinHeap
- adapter.go
- zustand
- jsonMarshalerError

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
10. `AlertDescription` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification` --references--> `CliToolsConsole()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → web/src/features/clitools/CliToolsConsole.tsx
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `File Map` --references--> `ModelDescriptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/contracts/contracts.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (219 total, 24 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.10
Nodes (38): ModelOption, configRequest, generateResponse, patchConfigRequest, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest (+30 more)

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

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

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
Cohesion: 0.07
Nodes (76): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, class-variance-authority, lucide-react, next (+68 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (101): stubAdminAuth, testHarness, TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), setupCloudSyncTest() (+93 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.07
Nodes (47): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, DashboardKeysPage(), TeamsPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys() (+39 more)

### Community 14 - "client.ts"
Cohesion: 0.08
Nodes (29): @tanstack/react-query, ImporterPage(), ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse (+21 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.17
Nodes (20): TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation() (+12 more)

### Community 16 - "catalog.go"
Cohesion: 0.21
Nodes (17): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+9 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (34): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+26 more)

### Community 21 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (8): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "CardTitle"
Cohesion: 0.08
Nodes (46): ArenaPage(), ChatPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels() (+38 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "cross-compile"
Cohesion: 0.15
Nodes (22): react-hook-form, VouchersPage(), useClaimVoucherMutation(), useCreateVoucherMutation(), useDeleteVoucherMutation(), useUpdateVoucherMutation(), useVoucherRedemptions(), useVouchersList() (+14 more)

### Community 29 - "Module"
Cohesion: 0.15
Nodes (10): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription() (+2 more)

### Community 30 - "providers/schemas.ts"
Cohesion: 0.08
Nodes (26): CatalogOAuthDescriptor, commaListSchema(), DEFAULT_PRIORITY, DEFAULT_UPDATE_CONNECTION_VALUES, headersTextSchema, isValidHttpUrl(), ListProviderConnectionsResponse, ListProviderNodesResponse (+18 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "updateConnectionTx"
Cohesion: 0.23
Nodes (12): cloneStringMap(), cloneStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), updateConnectionTx(), connectionRecord (+4 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.16
Nodes (9): Module, trackedResponseWriter, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), parseInboundMeta() (+1 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (48): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation() (+40 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.21
Nodes (15): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+7 more)

### Community 42 - "buildApp"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "UpstreamCredentials"
Cohesion: 0.08
Nodes (18): fakeProviderStore, resolveSnapshot, targetKey, stubProviderResolver, stubProviderStore, stubProviderStore, stubResolver, testEnv (+10 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.19
Nodes (19): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+11 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.03
Nodes (15): stubAdminCredentialStore, newRepository(), sanitizeBotError(), stubModule, Alias, AliasListResponse, Config, Event (+7 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "context.Context"
Cohesion: 0.05
Nodes (55): content, description, content, description, content, description, get, patch (+47 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.19
Nodes (19): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+11 more)

### Community 49 - "NewStaticUIHandler"
Cohesion: 0.19
Nodes (13): fakeComboStore, restoreSnapshotRequest, SnapshotBundle, legacyExtract, legacyKeyRecord, stubCombos, AliasSnapshot, ComboSnapshot (+5 more)

### Community 50 - "newTestHarness"
Cohesion: 0.11
Nodes (33): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+25 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (59): Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation() (+51 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.17
Nodes (21): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+13 more)

### Community 53 - "WriteError"
Cohesion: 0.14
Nodes (15): isValidProtocol(), oauthEnvPrefix(), decodeJSONBody(), validateHeaders(), validateName(), validateSlugID(), validateStringSlice(), actorFromContext() (+7 more)

### Community 54 - "Module"
Cohesion: 0.12
Nodes (7): providerIDContextKey, ToolDescriptor, ChangedEvent, PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 55 - "WriteJSON"
Cohesion: 0.10
Nodes (12): Module, decodeJSON(), deriveSessionTitle(), New(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON() (+4 more)

### Community 56 - "Module"
Cohesion: 0.19
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.18
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 58 - "repository"
Cohesion: 0.15
Nodes (9): decodeStringMap(), decodeStringSlice(), newRepository(), scanConnectionRecord(), scanNode(), oauthSessionRecord, ProviderNode, repository (+1 more)

### Community 59 - "Config"
Cohesion: 0.11
Nodes (33): DashboardProxyPoolsPage(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 60 - "Module"
Cohesion: 0.14
Nodes (6): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), New(), formatUTC()

### Community 61 - "NewStreamTranslator"
Cohesion: 0.10
Nodes (10): fakeSettingsStore, stubSettings, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter (+2 more)

### Community 62 - "New"
Cohesion: 0.11
Nodes (24): loadUI(), loadUI(), newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions() (+16 more)

### Community 63 - "Bus"
Cohesion: 0.07
Nodes (35): Options, Global Constraints, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, stubKeyProvisioner, KeyProvisioner (+27 more)

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.10
Nodes (15): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds(), applyResult (+7 more)

### Community 67 - "newTestHarness"
Cohesion: 0.06
Nodes (34): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+26 more)

### Community 68 - "scripts"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "responsesStreamEncoder"
Cohesion: 0.12
Nodes (31): TelegramPage(), telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation() (+23 more)

### Community 71 - "newTestHarness"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - "Server"
Cohesion: 0.06
Nodes (20): testHarness, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback() (+12 more)

### Community 73 - "translate_test.go"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 74 - "time.Time"
Cohesion: 0.06
Nodes (21): CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, chainStore, failureRecord, stubHealthTracker (+13 more)

### Community 75 - "DB"
Cohesion: 0.14
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 76 - "stubModule"
Cohesion: 0.09
Nodes (7): streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ProcessRunner

### Community 77 - "sync.Mutex"
Cohesion: 0.10
Nodes (20): clsx, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, dependencies, class-variance-authority@^0.7.1 (+12 more)

### Community 78 - "time.Duration"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Load"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "providers/module_test.go"
Cohesion: 0.09
Nodes (17): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+9 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "stubModule"
Cohesion: 0.10
Nodes (9): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), cloneOAuthConfig(), defaultBuiltInAdapters(), New(), builtinAdapter, customTestAdapter, OAuthConfig, Options (+1 more)

### Community 83 - "collectMigrations"
Cohesion: 0.15
Nodes (12): Combo, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues, cloneStringSlice() (+4 more)

### Community 84 - "dependencies"
Cohesion: 0.12
Nodes (31): MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders() (+23 more)

### Community 85 - "Module"
Cohesion: 0.11
Nodes (11): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), extractAPIKey(), ensureRequestID(), clientIP() (+3 more)

### Community 87 - "contracts.go"
Cohesion: 0.08
Nodes (43): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, patchComboRequest, upsertAliasRequest (+35 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "nameMinHeap"
Cohesion: 0.12
Nodes (13): stubProxyDialerResolver, Options, stubProxyResolver, matchToolForHost(), newID(), singleJoiningSlash(), parseAndValidateProxyURL(), probeProxy() (+5 more)

### Community 90 - "buildApp"
Cohesion: 0.08
Nodes (45): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+37 more)

### Community 91 - "crypto.SessionSigner"
Cohesion: 0.13
Nodes (25): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+17 more)

### Community 92 - "net/http.Request"
Cohesion: 0.17
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - ".executeCandidateAttempt"
Cohesion: 0.16
Nodes (19): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString() (+11 more)

### Community 94 - "@radix-ui/react-label"
Cohesion: 0.12
Nodes (21): File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator (+13 more)

### Community 95 - "repository"
Cohesion: 0.12
Nodes (23): PluginsPage(), SkillsPage(), CreateSkillRequest, PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse (+15 more)

### Community 96 - "repo"
Cohesion: 0.14
Nodes (7): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo, newRepo()

### Community 98 - "properties"
Cohesion: 0.09
Nodes (26): type, type, type, properties, type, type, type, type (+18 more)

### Community 99 - "extractErrorMessage"
Cohesion: 0.08
Nodes (49): CloudSyncPage(), PxpipePage(), CloudSyncStudio(), CloudSyncActionResponse, CloudSyncConfig, CloudSyncHistoryListResponse, cloudSyncQueryKeys, CloudSyncSnapshotBundle (+41 more)

### Community 100 - "properties"
Cohesion: 0.15
Nodes (13): additionalProperties, properties, required, type, type, $ref, type, type (+5 more)

### Community 101 - "Module"
Cohesion: 0.19
Nodes (5): computeBundleChecksum(), decodeJSON(), Module, maskSecret(), validateEndpointURL()

### Community 102 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 103 - "io/fs.FS"
Cohesion: 0.05
Nodes (8): Module, validateSQLiteSource(), decodeJSON(), actorFromContext(), decodeJSONBody(), Module, AdminSubject(), RouteMounts

### Community 104 - ".Reader"
Cohesion: 0.07
Nodes (38): main(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, File Map (+30 more)

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
Cohesion: 0.11
Nodes (10): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), decodeOptionalJSONBody(), WriteError(), jsonMethodNotAllowedHandler() (+2 more)

### Community 109 - "antigravity.go"
Cohesion: 0.11
Nodes (28): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, attemptOutcome, agGroupMatchesFamily(), agIsSessionBucket() (+20 more)

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.12
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 112 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "ModelDescriptor"
Cohesion: 0.36
Nodes (5): stubModelLister, stubModelLister, stubModelLister, ModelDescriptor, stubModelLister

### Community 116 - "repository"
Cohesion: 0.27
Nodes (5): Config, HistoryEvent, repository, storedConfig, boolToInt()

### Community 117 - "properties"
Cohesion: 0.13
Nodes (15): additionalProperties, properties, required, type, type, type, type, created_at (+7 more)

### Community 118 - "Module"
Cohesion: 0.07
Nodes (31): authHarness, fakeCredentialStore, Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, Module, decodeErrorCode(), doRequest() (+23 more)

### Community 119 - "repo"
Cohesion: 0.14
Nodes (12): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow() (+4 more)

### Community 120 - "ProvidersConsole.tsx"
Cohesion: 0.12
Nodes (17): ActiveOAuthSession, CatalogAdapterDescriptor, CreateConnectionFormInput, CreateConnectionFormOutput, CreateNodeFormInput, CreateNodeFormOutput, DEFAULT_CREATE_CONNECTION_VALUES, DEFAULT_CREATE_NODE_VALUES (+9 more)

### Community 121 - "repository"
Cohesion: 0.16
Nodes (13): boolToInt(), cloneMap(), lookupPluginSpec(), mergeDefaults(), newRepository(), sanitizeStringSlice(), scanSkill(), PluginSpec (+5 more)

### Community 122 - "properties"
Cohesion: 0.16
Nodes (19): type, type, type, type, properties, type, type, type (+11 more)

### Community 123 - "JobInspector"
Cohesion: 0.11
Nodes (25): TestAppCloseTeardownOrderIsLIFO(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat` (+17 more)

### Community 125 - "Module"
Cohesion: 0.16
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 126 - "apply.go"
Cohesion: 0.70
Nodes (4): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic()

### Community 127 - "properties"
Cohesion: 0.12
Nodes (17): additionalProperties, properties, required, type, type, type, type, type (+9 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.17
Nodes (11): 1. Build & Run Single Binary, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Stage), 4. Local Development (Live Reload), Architecture, Configuration, License, Modules (+3 more)

### Community 129 - "UpstreamCredentials"
Cohesion: 0.13
Nodes (8): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), truncateError(), Vault, refreshCall, refreshCooldownEntry

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.16
Nodes (13): CanonicalStreamChunk, geminiStreamDecoder, openAIStreamDecoder, responsesStreamDecoder, streamDecoder, NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamDecoder() (+5 more)

### Community 131 - "certManager"
Cohesion: 0.13
Nodes (8): randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.24
Nodes (7): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "DB"
Cohesion: 0.07
Nodes (36): Options, Options, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, New(), newRepository(), New(), TestModuleMetadataAndStartValidation() (+28 more)

### Community 134 - "Bus"
Cohesion: 0.17
Nodes (3): decodeJSON(), Module, New()

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "resolve.js"
Cohesion: 0.19
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 137 - "release-binaries"
Cohesion: 0.20
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "providers/handlers.go"
Cohesion: 0.10
Nodes (17): antiSlopDirective(), applyPluginsToBody(), injectSystemDirectives(), intFromAny(), modelMatches(), squeezePayloadMessages(), catalogAdapterDescriptor, catalogOAuthDescriptor (+9 more)

### Community 140 - "net/http.Client"
Cohesion: 0.15
Nodes (16): newBotAPIClient(), New(), doJSON(), setupTelegramTest(), TestTelegramCommandsTestMessageAndEventBusForwarding(), TestTelegramConfigVaultIsolationAndWebhookSecurity(), TestTelegramNetworkErrorRedactsBotToken(), newRepo() (+8 more)

### Community 141 - "net/http.Header"
Cohesion: 0.20
Nodes (9): arenaTestEnv, stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent(), tgAPICall (+1 more)

### Community 142 - "ProvidersConsole"
Cohesion: 0.13
Nodes (16): DashboardProvidersPage(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation(), useExchangeProviderOAuthMutation(), useProviderConnections(), useProvidersCatalog(), useStartProviderOAuthMutation(), useTestProviderConnectionMutation() (+8 more)

### Community 143 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 144 - "providers/hooks.ts"
Cohesion: 0.14
Nodes (16): ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useDeleteProviderConnectionMutation(), useProviderNodes(), CreateProviderConnectionRequest (+8 more)

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.18
Nodes (15): description, description, content, description, delete, operationId, responses, summary (+7 more)

### Community 149 - "properties"
Cohesion: 0.14
Nodes (15): additionalProperties, type, type, type, properties, config, key, models (+7 more)

### Community 150 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (8): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - ".executeBuiltinFetch"
Cohesion: 0.25
Nodes (9): isBlockedAddr(), extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest (+1 more)

### Community 153 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 154 - "schemas"
Cohesion: 0.25
Nodes (8): additionalProperties, required, type, schemas, CloudSyncConfig, UpdateCloudSyncConfigRequest, additionalProperties, type

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "paths"
Cohesion: 0.26
Nodes (12): get, get, get, get, operationId, summary, tags, paths (+4 more)

### Community 157 - "compressor.go"
Cohesion: 0.56
Nodes (7): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput()

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
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, newID(), injectAnthropicCacheControl(), computePrefixFingerprint()

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

### Community 168 - "setupMediaTest"
Cohesion: 0.19
Nodes (11): New(), doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), newRepo() (+3 more)

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

### Community 177 - "providers/schemas.test.ts"
Cohesion: 0.22
Nodes (8): createConnectionSchema, createNodeSchema, oauthExchangeSchema, oauthStartSchema, toCreateConnectionRequest(), toCreateNodeRequest(), toOAuthExchangeRequest(), updateConnectionSchema

### Community 178 - "repo"
Cohesion: 0.29
Nodes (6): newID(), newRepo(), nowUTC(), ConfigRecord, Event, repo

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Config"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

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

### Community 193 - "media/module.go"
Cohesion: 0.29
Nodes (4): createComboReq, createProviderReq, updateComboReq, updateProviderReq

### Community 194 - "pxpipe/module.go"
Cohesion: 0.33
Nodes (6): HealthCheck, HealthReport, patchConfigRequest, previewRequest, PreviewResponse, Stats

### Community 195 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 201 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "UsageRecord"
Cohesion: 0.16
Nodes (6): stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 207 - "modeleditor/engine.go"
Cohesion: 0.06
Nodes (32): memoryResponseWriter, newMemoryResponseWriter(), copyResponseHeaders(), copyNonHopHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp() (+24 more)

### Community 208 - "Config"
Cohesion: 0.09
Nodes (22): Options, runMain(), Task 2: `security` Module (Security Posture & Redacted Audit Trail), Review Focus, Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Options (+14 more)

### Community 209 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider()

### Community 213 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 215 - "adapter.go"
Cohesion: 0.32
Nodes (5): buildProtocolTestRequest(), isBlockedAntigravityProjectID(), isValidAuthType(), validateHTTPURL(), oauthTokenResponse

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 933 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **24 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `contracts.go` to `DB`, `.Reader`, `AppShell.tsx`, `KeysConsole.tsx`, `Module`, `Module`, `buildApp`?**
  _High betweenness centrality (0.113) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.09758771929824561 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `extractErrorMessage` to `SettingsConsole.tsx`, `Server`, `AppShell.tsx`, `KeysConsole.tsx`, `client.ts`, `ProvidersConsole`, `providers/hooks.ts`, `io/fs.FS`, `CardTitle`, `cross-compile`, `UsageConsole.tsx`, `newTestHarness`, `ModelsRoutingConsole.tsx`, `CanonicalStreamChunk`, `Config`, `scripts`, `responsesStreamEncoder`, `newTestHarness`, `CliToolsConsole.tsx`, `dependencies`, `crypto.SessionSigner`, `repository`, `TunnelConsole.tsx`?**
  _High betweenness centrality (0.113) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `AppShell.tsx` to `contracts.go`?**
  _High betweenness centrality (0.088) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._