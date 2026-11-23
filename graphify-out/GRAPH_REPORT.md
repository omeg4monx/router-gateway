# Graph Report - wave3c-omeg4mon-router  (2026-11-23)

## Corpus Check
- 324 files · ~315,643 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 3711 nodes · 12411 edges · 161 communities (136 shown, 25 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1185 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ebe9f692`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- server.go
- Module
- biome.json
- settings
- Wave 1 — Core LLM Gateway & Essential Control Plane
- APIKeySummary
- package.json
- .executeCandidateAttempt
- compilerOptions
- NotificationsConsole.tsx
- Button
- testing.T
- components.json
- KeysConsole.tsx
- GatewayInterceptor
- Input
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- TeamsConsole.tsx
- steps
- web-check
- on
- cross-compile
- WriteError
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
- TestContractInterfacesSatisfied
- translate_gemini.go
- context.Context
- repository
- CardTitle
- translate_openai.go
- lifecycle.Registry
- PromptGuardConsole.tsx
- extractErrorMessage
- encoding/json.RawMessage
- net/http.ResponseWriter
- go_pkg_context
- WriteJSON
- Module
- CardHeader
- repository
- ProxyPoolsConsole.tsx
- modeleditor/repo.go
- Runner
- Server
- Module
- Module
- listSettingsWithQueryer
- GatewayInterceptor
- repository
- ModelEditorConsole.tsx
- Module
- go_pkg_encoding_json
- MitmConsole.tsx
- Module
- scripts
- time.Time
- Bus
- Module
- dependencies
- Module
- Module
- setupModelEditorTest
- CliToolsConsole.tsx
- OAuthConfig
- repository
- MediaConsole.tsx
- Module
- batcher
- File Map
- DB
- net/url.URL
- buildApp
- Label
- repo
- app/layout.tsx
- RingBuffer
- Module
- providers/module_test.go
- PxpipeConsole.tsx
- stubModule
- client.ts
- .Reader
- openapi-fetch
- settingDefinition
- setupPromptGuardTest
- Module
- setupTokenSaverTest
- .StatusCode
- Module
- repository
- repo
- Module
- .executeBuiltinFetch
- New
- main.go
- newTestHarness
- setupTestMitm
- setupMediaTest
- responsesStreamEncoder
- setupTunnelTest
- setupNotificationsTest
- .Name
- applyConfigFile
- promptinjection/engine.go
- cloneStringSlice
- CanonicalStreamChunk
- certManager
- NewStaticUIHandler
- newTestHarness
- Vault
- tokensaver/module.go
- stubHealthTracker
- repository
- RouteMounts
- providers/handlers.go
- headroomCache
- Open
- Module
- runTransformPipeline
- newTestHarness
- sync.Mutex
- stubModule
- Load
- combos/module.go
- stubModule
- ProcessRunner
- stubModule
- pxpipe/module.go
- crypto.SessionSigner
- stubSystemModule
- stubDepModule
- stubDep
- collectMigrations
- nullableString
- EventListResponse
- @radix-ui/react-label

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 169 edges
2. `WriteError()` - 151 edges
3. `WriteJSON()` - 148 edges
4. `DB` - 85 edges
5. `buildApp()` - 70 edges
6. `Bus` - 66 edges
7. `Button` - 59 edges
8. `Module` - 55 edges
9. `Alert` - 55 edges
10. `AlertDescription` - 54 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (161 total, 25 thin omitted)

### Community 0 - "server.go"
Cohesion: 0.12
Nodes (26): configRequest, adminSubjectKey, ErrorDetail, ErrorEnvelope, applyResponse, rollbackResponse, Module, ChangedEvent (+18 more)

### Community 1 - "Module"
Cohesion: 0.24
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

### Community 5 - "APIKeySummary"
Cohesion: 0.08
Nodes (24): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository() (+16 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - ".executeCandidateAttempt"
Cohesion: 0.09
Nodes (27): ProtocolTranslator, attemptOutcome, CanonicalChatResponse, buildUpstreamURL(), extractUpstreamErrorMessage(), parseRetryAfter(), CanFastPath(), DecodeResponse() (+19 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "NotificationsConsole.tsx"
Cohesion: 0.14
Nodes (26): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+18 more)

### Community 10 - "Button"
Cohesion: 0.08
Nodes (65): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, class-variance-authority, lucide-react, next, react (+57 more)

### Community 11 - "testing.T"
Cohesion: 0.07
Nodes (55): testHarness, TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback() (+47 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (29): @hookform/resolvers, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "GatewayInterceptor"
Cohesion: 0.10
Nodes (18): Global Constraints, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, Global Constraints, Wave 3C — Local MITM Proxy & PXPIPE Transform Pipeline Implementation Plan, CanonicalRequest (+10 more)

### Community 15 - "Input"
Cohesion: 0.15
Nodes (21): @tanstack/react-query, @tanstack/react-query@^5.76.1, TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig() (+13 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.17
Nodes (16): Options, testHarness, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+8 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (35): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+27 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.07
Nodes (49): zod, zod@^4.6.7, DashboardSettingsPage(), authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+41 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "TeamsConsole.tsx"
Cohesion: 0.20
Nodes (17): TeamsPage(), useApiKeys(), teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation(), useRemoveTeamKeyMutation(), useResetTeamSpendMutation() (+9 more)

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

### Community 29 - "WriteError"
Cohesion: 0.19
Nodes (9): actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription(), validateModelIdentifier() (+1 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.05
Nodes (89): DashboardProviderDetailPage(), DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation() (+81 more)

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
Nodes (20): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+12 more)

### Community 39 - "net/http.Request"
Cohesion: 0.08
Nodes (22): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+14 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.21
Nodes (15): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+7 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (21): fakeProviderStore, ModelDescriptor, ModelLister, ProviderConnection, ProviderStore, ResolvedTarget, stubKeyValidator, stubModelLister (+13 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "context.Context"
Cohesion: 0.04
Nodes (11): fakeCredentialStore, AdminCredentialStore, stubAdminCredentialStore, stubUsageRecorder, actorFromContext(), TestSystemAdminCredentialStore(), validatePassword(), stubModule (+3 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "CardTitle"
Cohesion: 0.15
Nodes (22): VouchersPage(), useClaimVoucherMutation(), useCreateVoucherMutation(), useDeleteVoucherMutation(), useUpdateVoucherMutation(), useVoucherRedemptions(), useVouchersList(), vouchersQueryKeys (+14 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.18
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+12 more)

### Community 49 - "lifecycle.Registry"
Cohesion: 0.12
Nodes (9): File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests, Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, normalizeName(), nameMinHeap (+1 more)

### Community 50 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "extractErrorMessage"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "encoding/json.RawMessage"
Cohesion: 0.16
Nodes (22): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+14 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.12
Nodes (16): trackedResponseWriter, isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), sanitizeUpstreamError(), truncateError() (+8 more)

### Community 54 - "go_pkg_context"
Cohesion: 0.21
Nodes (5): mapBusEventToAudit(), payloadToMap(), PreviewResult, AuditInput, storedSetting

### Community 55 - "WriteJSON"
Cohesion: 0.13
Nodes (10): decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), WriteJSON(), LogEntry, LogListResponse (+2 more)

### Community 56 - "Module"
Cohesion: 0.22
Nodes (6): compileRule(), decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 57 - "CardHeader"
Cohesion: 0.19
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 58 - "repository"
Cohesion: 0.10
Nodes (21): cloneStringMap(), cloneUpstreamCredentials(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+13 more)

### Community 59 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "modeleditor/repo.go"
Cohesion: 0.10
Nodes (20): copyResponseHeaders(), copyNonHopHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), wrapSystemPrompt() (+12 more)

### Community 61 - "Runner"
Cohesion: 0.18
Nodes (7): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter

### Community 62 - "Server"
Cohesion: 0.09
Nodes (18): TestAppCloseTeardownOrderIsLIFO(), loadUI(), loadUI(), AuthVerifier, Options, apiTimeoutMiddleware(), Server, New() (+10 more)

### Community 63 - "Module"
Cohesion: 0.23
Nodes (8): actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken(), Module

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.23
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "GatewayInterceptor"
Cohesion: 0.07
Nodes (4): stubInterceptor, recordingInterceptor, ResponseMeta, StreamEvent

### Community 67 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 68 - "ModelEditorConsole.tsx"
Cohesion: 0.13
Nodes (30): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+22 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "go_pkg_encoding_json"
Cohesion: 0.27
Nodes (6): generateResponse, applyResponse, catalogResponse, rollbackResponse, historyListResp, settingsListResp

### Community 71 - "MitmConsole.tsx"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - "Module"
Cohesion: 0.09
Nodes (8): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.14
Nodes (14): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, Options, newCooldownTracker(), New() (+6 more)

### Community 75 - "Bus"
Cohesion: 0.08
Nodes (9): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), clientIP(), decodeStrictJSON(), Module, isValidSeverity() (+1 more)

### Community 76 - "Module"
Cohesion: 0.07
Nodes (14): streamResponseBody(), newID(), newRepo(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord, Event (+6 more)

### Community 77 - "dependencies"
Cohesion: 0.12
Nodes (17): clsx, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies, class-variance-authority@^0.7.1, clsx@^2.1.1 (+9 more)

### Community 78 - "Module"
Cohesion: 0.11
Nodes (16): New(), New(), decodeStrictJSON(), extractUsageSpend(), Module, New(), New(), Bus (+8 more)

### Community 79 - "Module"
Cohesion: 0.19
Nodes (7): validateModelPattern(), decodeJSONBody(), New(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "OAuthConfig"
Cohesion: 0.13
Nodes (4): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), customTestAdapter, ProviderAdapter

### Community 83 - "repository"
Cohesion: 0.08
Nodes (24): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+16 more)

### Community 84 - "MediaConsole.tsx"
Cohesion: 0.12
Nodes (32): react-hook-form, MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos() (+24 more)

### Community 85 - "Module"
Cohesion: 0.18
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 86 - "batcher"
Cohesion: 0.14
Nodes (7): newBatcher(), New(), newPricingStore(), batcher, ModelPricing, Options, pricingStore

### Community 87 - "File Map"
Cohesion: 0.28
Nodes (9): ProviderResolver, File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`) (+1 more)

### Community 88 - "DB"
Cohesion: 0.07
Nodes (25): Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), KeyValidator, New(), newRepo(), boolToInt(), newRepository(), New(), newRepository() (+17 more)

### Community 89 - "net/url.URL"
Cohesion: 0.10
Nodes (16): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, matchToolForHost(), newID(), singleJoiningSlash(), parseAndValidateProxyURL(), probeProxy() (+8 more)

### Community 90 - "buildApp"
Cohesion: 0.05
Nodes (46): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure() (+38 more)

### Community 91 - "Label"
Cohesion: 0.13
Nodes (25): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+17 more)

### Community 92 - "repo"
Cohesion: 0.08
Nodes (23): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow() (+15 more)

### Community 93 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider()

### Community 94 - "RingBuffer"
Cohesion: 0.20
Nodes (11): File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+3 more)

### Community 97 - "providers/module_test.go"
Cohesion: 0.10
Nodes (18): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+10 more)

### Community 99 - "PxpipeConsole.tsx"
Cohesion: 0.13
Nodes (32): PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats(), usePxpipeStatus() (+24 more)

### Community 102 - "client.ts"
Cohesion: 0.06
Nodes (35): TunnelPage(), isHealthResponse(), useSystemHealth(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus() (+27 more)

### Community 104 - ".Reader"
Cohesion: 0.06
Nodes (32): main(), runMain(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment (+24 more)

### Community 106 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 107 - "setupPromptGuardTest"
Cohesion: 0.27
Nodes (8): Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins(), TestPromptInjectionCustomRuleValidationAndRE2(), TestPromptInjectionFinalReviewFindings(), testEnv

### Community 108 - "Module"
Cohesion: 0.15
Nodes (7): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), WriteError(), Module

### Community 109 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 112 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 115 - ".executeBuiltinFetch"
Cohesion: 0.29
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 116 - "New"
Cohesion: 0.09
Nodes (24): Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Options, Options, newChainStore(), New(), New(), newHarness(), TestTranslatorLogsSnapshotAndLiveSSEStream() (+16 more)

### Community 117 - "main.go"
Cohesion: 0.06
Nodes (9): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), ToolDescriptor, PoolTestResponse, ScanCandidateResult, ScanResponse, HistoryItem (+1 more)

### Community 118 - "newTestHarness"
Cohesion: 0.22
Nodes (10): Task 2: `auth` Module & OpenAPI Fragment, newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify(), HashPassword() (+2 more)

### Community 119 - "setupTestMitm"
Cohesion: 0.31
Nodes (10): New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly(), TestMitmLiveTLSProxyInterceptAndAliasRewrite() (+2 more)

### Community 120 - "setupMediaTest"
Cohesion: 0.12
Nodes (16): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), newHarness(), TestProxyPoolsBindingLifecycleAndCacheInvalidation() (+8 more)

### Community 121 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 122 - "setupTunnelTest"
Cohesion: 0.25
Nodes (9): New(), doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary() (+1 more)

### Community 123 - "setupNotificationsTest"
Cohesion: 0.16
Nodes (15): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict() (+7 more)

### Community 124 - ".Name"
Cohesion: 0.08
Nodes (5): authHarness, Module, TestRegistryStartAllRollsBackOnPartialFailure(), stubSystemModule, stubSystemModule

### Community 126 - "applyConfigFile"
Cohesion: 0.33
Nodes (5): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), writeAtomicFile()

### Community 128 - "promptinjection/engine.go"
Cohesion: 0.07
Nodes (36): actionRank(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity(), redactAndTruncate() (+28 more)

### Community 129 - "cloneStringSlice"
Cohesion: 0.24
Nodes (6): cloneOAuthConfig(), cloneStringSlice(), OAuthConfig, oauthEnvPrefix(), buildAuthorizeURL(), builtinAdapter

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 131 - "certManager"
Cohesion: 0.14
Nodes (8): newCertManager(), randomSerial(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 133 - "newTestHarness"
Cohesion: 0.28
Nodes (8): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), testHarness

### Community 134 - "Vault"
Cohesion: 0.29
Nodes (5): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, New(), newRepository(), Vault, Options

### Community 135 - "tokensaver/module.go"
Cohesion: 0.32
Nodes (5): CompressionStats, HeadroomStats, PreviewResponse, StatsResponse, WarmPrefixSummary

### Community 136 - "stubHealthTracker"
Cohesion: 0.16
Nodes (5): AccountHealthTracker, stubAccountHealthTracker, failureRecord, stubHealthTracker, computeBackoff()

### Community 137 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 138 - "RouteMounts"
Cohesion: 0.15
Nodes (5): actorFromContext(), decodeJSONBody(), Module, Module, RouteMounts

### Community 139 - "providers/handlers.go"
Cohesion: 0.19
Nodes (11): catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, createConnectionRequest, createNodeRequest, listNodesResponse, oauthExchangeRequest, oauthStartRequest (+3 more)

### Community 141 - "Open"
Cohesion: 0.08
Nodes (28): New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), doPxpipeJSON() (+20 more)

### Community 142 - "Module"
Cohesion: 0.16
Nodes (5): Module, New(), Options, refreshCall, storedCredentials

### Community 143 - "runTransformPipeline"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 144 - "newTestHarness"
Cohesion: 0.24
Nodes (9): decodeJSON(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageSummaryLogsAndRetentionPrune(), waitForLogCount() (+1 more)

### Community 147 - "Load"
Cohesion: 0.33
Nodes (7): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 148 - "combos/module.go"
Cohesion: 0.29
Nodes (6): ChangedEvent, createComboRequest, deleteDisabledModelRequest, patchComboRequest, upsertAliasRequest, upsertDisabledModelRequest

### Community 150 - "ProcessRunner"
Cohesion: 0.28
Nodes (3): fakeRunner, osProcessRunner, ProcessRunner

### Community 152 - "pxpipe/module.go"
Cohesion: 0.33
Nodes (6): HealthCheck, HealthReport, patchConfigRequest, previewRequest, PreviewResponse, Stats

### Community 157 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 559 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `.Reader` to `buildApp`, `Button`, `newTestHarness`, `RouteMounts`?**
  _High betweenness centrality (0.100) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11882716049382716 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `PxpipeConsole.tsx` to `NotificationsConsole.tsx`, `Button`, `KeysConsole.tsx`, `Input`, `SettingsConsole.tsx`, `TeamsConsole.tsx`, `ProvidersConsole.tsx`, `UsageConsole.tsx`, `CardTitle`, `PromptGuardConsole.tsx`, `extractErrorMessage`, `CardHeader`, `ProxyPoolsConsole.tsx`, `ModelEditorConsole.tsx`, `MitmConsole.tsx`, `CliToolsConsole.tsx`, `MediaConsole.tsx`, `Label`, `client.ts`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `.Reader`?**
  _High betweenness centrality (0.081) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._