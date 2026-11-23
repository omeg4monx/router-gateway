# Graph Report - wave3c-omeg4mon-router  (2026-11-23)

## Corpus Check
- 324 files · ~315,664 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 3712 nodes · 12411 edges · 173 communities (144 shown, 29 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1185 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8ccb971a`
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
- CanonicalChatResponse
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
- Vault
- repo
- net/url.URL
- buildApp
- Label
- repo
- app/layout.tsx
- RingBuffer
- Module
- go_pkg_strings
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
- resolver.go
- New
- mitm/module.go
- newTestHarness
- setupTestMitm
- setupMediaTest
- newTestHarness
- setupTunnelTest
- New
- .Name
- applyConfigFile
- promptinjection/engine.go
- updateConnectionTx
- CanonicalStreamChunk
- certManager
- NewStaticUIHandler
- newTestHarness
- DB
- modeleditor/engine.go
- stubHealthTracker
- repository
- RouteMounts
- media/module.go
- headroomCache
- Open
- Module
- runTransformPipeline
- newTestHarness
- contracts.UsageRecord
- stubModule
- repo
- combos/module.go
- stubModule
- repo
- stubModule
- newTestHarness
- crypto.SessionSigner
- stubSystemModule
- stubDepModule
- stubDep
- collectMigrations
- pricingStore
- NewRunner
- @radix-ui/react-label
- Notification
- NewStreamTranslator
- compressor.go
- .OnStreamChunk
- redact.go
- repository
- repository
- orderStubModule
- ProtocolTranslator
- syncBuffer
- mapBusEventToAudit

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

## Communities (173 total, 29 thin omitted)

### Community 0 - "server.go"
Cohesion: 0.12
Nodes (28): Task 1: `media` Module (Providers, Combos, Gateway Media Endpoints & SSRF-Safe Web Fetch), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs), adminSubjectKey, ErrorDetail, ErrorEnvelope, GatewayInterceptor, KeyValidator, Module (+20 more)

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

### Community 7 - "CanonicalChatResponse"
Cohesion: 0.60
Nodes (6): CanonicalChatResponse, EncodeResponse(), rawJSONToString(), encodeGeminiResponse(), encodeOpenAIChatResponse(), encodeOpenAIResponsesResponse()

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
Nodes (50): TestSentinelErrorsDistinct(), setupGatewayTest(), TestGatewayAuthExtractionAndProtocolErrors(), TestGatewayCrossProtocolStreaming(), TestGatewayDirectFastPathIncrementalStreaming(), TestGatewayErrorClassificationAndOAuthSingleRetry(), TestGatewayInterceptorsPriorityOrdering(), TestGatewayLargeSSELineAndAPIKey401Failover() (+42 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (29): @hookform/resolvers, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "GatewayInterceptor"
Cohesion: 0.17
Nodes (11): Global Constraints, Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Wave 3A — Teams (API-Key Grouping & Budgets) & Vouchers Implementation Plan, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest (+3 more)

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
Cohesion: 0.11
Nodes (19): authHarness, fakeCredentialStore, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+11 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.11
Nodes (43): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+35 more)

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
Nodes (18): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+10 more)

### Community 39 - "net/http.Request"
Cohesion: 0.10
Nodes (18): attemptOutcome, readSessionCookie(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError() (+10 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.20
Nodes (17): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+9 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (16): ProviderConnection, ProviderStore, ResolvedTarget, stubKeyValidator, stubProviderResolver, stubProviderStore, stubUsageRecorder, UpstreamCredentials (+8 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "context.Context"
Cohesion: 0.05
Nodes (5): AdminCredentialStore, stubAdminCredentialStore, actorFromContext(), oauthTokenResponse, stubCredentialStore

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "CardTitle"
Cohesion: 0.15
Nodes (22): VouchersPage(), useClaimVoucherMutation(), useCreateVoucherMutation(), useDeleteVoucherMutation(), useUpdateVoucherMutation(), useVoucherRedemptions(), useVouchersList(), vouchersQueryKeys (+14 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

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
Cohesion: 0.17
Nodes (22): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+14 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (20): clearSessionCookie(), isSecureRequest(), setSessionCookie(), decodeJSONBody(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody() (+12 more)

### Community 55 - "WriteJSON"
Cohesion: 0.13
Nodes (10): decodeStrictJSON(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), WriteJSON(), LogEntry, LogListResponse (+2 more)

### Community 56 - "Module"
Cohesion: 0.18
Nodes (8): compileRule(), decodeJSONBody(), isValidMode(), isValidSensitivity(), New(), validateRule(), compiledRule, Module

### Community 57 - "CardHeader"
Cohesion: 0.19
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 58 - "repository"
Cohesion: 0.15
Nodes (8): decodeStringMap(), decodeStringSlice(), scanConnectionRecord(), scanNode(), oauthSessionRecord, ProviderNode, repository, rowScanner

### Community 59 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "modeleditor/repo.go"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 61 - "Runner"
Cohesion: 0.18
Nodes (7): JobInspector, JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter

### Community 62 - "Server"
Cohesion: 0.09
Nodes (20): loadUI(), loadUI(), AuthVerifier, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, trackedResponseWriter, New(), newRepository() (+12 more)

### Community 63 - "Module"
Cohesion: 0.21
Nodes (9): actorFromContext(), decodeJSONBody(), New(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+1 more)

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.23
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "GatewayInterceptor"
Cohesion: 0.06
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, newID(), injectAnthropicCacheControl(), computePrefixFingerprint()

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
Cohesion: 0.24
Nodes (9): configRequest, generateResponse, applyResponse, rollbackResponse, applyResponse, catalogResponse, rollbackResponse, historyListResp (+1 more)

### Community 71 - "MitmConsole.tsx"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - "Module"
Cohesion: 0.11
Nodes (5): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.19
Nodes (7): CooldownEntry, cooldownState, cooldownTracker, listCooldownsResponse, chainStore, newCooldownTracker(), cacheEntry

### Community 75 - "Bus"
Cohesion: 0.15
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 76 - "Module"
Cohesion: 0.05
Nodes (18): streamResponseBody(), New(), newID(), newRepo(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord (+10 more)

### Community 77 - "dependencies"
Cohesion: 0.12
Nodes (17): clsx, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies, class-variance-authority@^0.7.1, clsx@^2.1.1 (+9 more)

### Community 78 - "Module"
Cohesion: 0.24
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Module"
Cohesion: 0.20
Nodes (4): decodeJSONBody(), newID(), validateOverride(), Module

### Community 80 - "setupModelEditorTest"
Cohesion: 0.27
Nodes (8): New(), doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "OAuthConfig"
Cohesion: 0.11
Nodes (5): cloneOAuthConfig(), OAuthConfig, buildAuthorizeURL(), builtinAdapter, customTestAdapter

### Community 83 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 84 - "MediaConsole.tsx"
Cohesion: 0.12
Nodes (32): react-hook-form, MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos() (+24 more)

### Community 85 - "Module"
Cohesion: 0.23
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 86 - "batcher"
Cohesion: 0.21
Nodes (4): newBatcher(), New(), newPricingStore(), batcher

### Community 87 - "Vault"
Cohesion: 0.12
Nodes (13): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification) (+5 more)

### Community 88 - "repo"
Cohesion: 0.27
Nodes (6): formatUTC(), newTeamID(), createTeamInput, patchTeamInput, repo, Team

### Community 89 - "net/url.URL"
Cohesion: 0.09
Nodes (23): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+15 more)

### Community 90 - "buildApp"
Cohesion: 0.13
Nodes (26): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+18 more)

### Community 91 - "Label"
Cohesion: 0.13
Nodes (25): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+17 more)

### Community 92 - "repo"
Cohesion: 0.18
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - "app/layout.tsx"
Cohesion: 0.29
Nodes (6): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider()

### Community 94 - "RingBuffer"
Cohesion: 0.18
Nodes (12): File Map, Global Constraints, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+4 more)

### Community 96 - "Module"
Cohesion: 0.20
Nodes (3): decodeJSONBody(), New(), Module

### Community 97 - "go_pkg_strings"
Cohesion: 0.07
Nodes (18): PolicyChangedEvent, CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable() (+10 more)

### Community 99 - "PxpipeConsole.tsx"
Cohesion: 0.13
Nodes (32): PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats(), usePxpipeStatus() (+24 more)

### Community 102 - "client.ts"
Cohesion: 0.06
Nodes (35): TunnelPage(), isHealthResponse(), useSystemHealth(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus() (+27 more)

### Community 104 - ".Reader"
Cohesion: 0.06
Nodes (33): main(), runMain(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment (+25 more)

### Community 106 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 107 - "setupPromptGuardTest"
Cohesion: 0.27
Nodes (8): Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins(), TestPromptInjectionCustomRuleValidationAndRE2(), TestPromptInjectionFinalReviewFindings(), testEnv

### Community 108 - "Module"
Cohesion: 0.11
Nodes (12): decodeStrictJSON(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID() (+4 more)

### Community 109 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 110 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, httpStatusError, statusError

### Community 111 - "Module"
Cohesion: 0.12
Nodes (13): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+5 more)

### Community 112 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "Module"
Cohesion: 0.19
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 115 - "resolver.go"
Cohesion: 0.16
Nodes (13): resolveSnapshot, targetKey, ModelDescriptor, ModelLister, stubModelLister, stubModelLister, buildResolvedTarget(), connectionSupportsBare() (+5 more)

### Community 116 - "New"
Cohesion: 0.10
Nodes (22): Options, newChainStore(), New(), New(), newHarness(), TestTranslatorLogsSnapshotAndLiveSSEStream(), TestTranslatorLogsStreamSurvivesPastAPITimeout(), TestTranslatorTranslateRequestAndResponseAllPairs() (+14 more)

### Community 117 - "mitm/module.go"
Cohesion: 0.14
Nodes (9): errorDisposition, classifyStatus(), isClientPayloadErrorBody(), createAliasRequest, EventListResponse, generateCARequest, updateConfigRequest, HistoryItem (+1 more)

### Community 118 - "newTestHarness"
Cohesion: 0.29
Nodes (10): Task 2: `auth` Module & OpenAPI Fragment, newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify() (+2 more)

### Community 119 - "setupTestMitm"
Cohesion: 0.31
Nodes (10): New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly(), TestMitmLiveTLSProxyInterceptAndAliasRewrite() (+2 more)

### Community 120 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 121 - "newTestHarness"
Cohesion: 0.13
Nodes (16): fakeClock, fakeProviderStore, stubAdminAuth, testHarness, decodeJSON(), newFakeClock(), newFakeProviderStore(), newTestHarness() (+8 more)

### Community 122 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 123 - "New"
Cohesion: 0.13
Nodes (17): File Map, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict() (+9 more)

### Community 124 - ".Name"
Cohesion: 0.06
Nodes (6): Module, validatePassword(), TestRegistryStartAllRollsBackOnPartialFailure(), stubSystemModule, stubModule, stubSystemModule

### Community 126 - "applyConfigFile"
Cohesion: 0.70
Nodes (4): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic()

### Community 128 - "promptinjection/engine.go"
Cohesion: 0.33
Nodes (11): actionRank(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity(), redactAndTruncate() (+3 more)

### Community 129 - "updateConnectionTx"
Cohesion: 0.23
Nodes (13): cloneStringMap(), cloneStringSlice(), cloneUpstreamCredentials(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), updateConnectionTx() (+5 more)

### Community 130 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 131 - "certManager"
Cohesion: 0.12
Nodes (9): newCertManager(), randomSerial(), writeAtomicFile(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide (+1 more)

### Community 132 - "NewStaticUIHandler"
Cohesion: 0.33
Nodes (7): applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 133 - "newTestHarness"
Cohesion: 0.21
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 134 - "DB"
Cohesion: 0.05
Nodes (42): Options, Options, Options, New(), newRepo(), New(), newRepository(), New() (+34 more)

### Community 135 - "modeleditor/engine.go"
Cohesion: 0.19
Nodes (12): copyResponseHeaders(), copyNonHopHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern() (+4 more)

### Community 136 - "stubHealthTracker"
Cohesion: 0.15
Nodes (6): AccountHealthTracker, stubAccountHealthTracker, failureRecord, stubHealthTracker, computeBackoff(), parseRetryAfter()

### Community 137 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 138 - "RouteMounts"
Cohesion: 0.23
Nodes (4): actorFromContext(), decodeJSONBody(), Module, RouteMounts

### Community 139 - "media/module.go"
Cohesion: 0.10
Nodes (17): providerIDContextKey, createComboReq, createProviderReq, updateComboReq, updateProviderReq, catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse (+9 more)

### Community 141 - "Open"
Cohesion: 0.11
Nodes (18): New(), TestModuleMetadataAndStartValidation(), newRepository(), doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating(), TestPxpipeLifecycleStatusAndHealth(), TestPxpipePreservesCodeFencesAvoidsOrphanToolsAndEscaping() (+10 more)

### Community 142 - "Module"
Cohesion: 0.14
Nodes (9): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), oauthEnvPrefix(), Module, New(), newRepository(), ProviderAdapter, refreshCall (+1 more)

### Community 143 - "runTransformPipeline"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 144 - "newTestHarness"
Cohesion: 0.18
Nodes (11): decodeJSON(), newFakeClock(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageSummaryLogsAndRetentionPrune() (+3 more)

### Community 147 - "repo"
Cohesion: 0.23
Nodes (8): buildWhere(), formatUTC(), newEventID(), newRepo(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 148 - "combos/module.go"
Cohesion: 0.09
Nodes (19): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+11 more)

### Community 150 - "repo"
Cohesion: 0.20
Nodes (7): mapBusEventToNotification(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), EmitPayload, repo

### Community 152 - "newTestHarness"
Cohesion: 0.24
Nodes (9): Options, testHarness, New(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract() (+1 more)

### Community 157 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 159 - "NewRunner"
Cohesion: 0.28
Nodes (9): TestAppCloseTeardownOrderIsLIFO(), Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, NewRunner(), TestRunnerIdempotentStartStopAndGuards(), TestRunnerOverlapGuardReturnsAlreadyRunning(), TestRunnerPanicRecoveryAndObserverCallback(), TestRunnerRecoversFromPanicAndStopsCleanly(), TestRunnerSnapshotAndManualTrigger() (+1 more)

### Community 161 - "Notification"
Cohesion: 0.31
Nodes (5): buildWhere(), scanNotification(), scanNotificationRow(), ListFilter, Notification

### Community 162 - "NewStreamTranslator"
Cohesion: 0.22
Nodes (9): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder() (+1 more)

### Community 163 - "compressor.go"
Cohesion: 0.56
Nodes (7): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), leadingBackticks(), truncateToolOutput()

### Community 165 - "redact.go"
Cohesion: 0.46
Nodes (6): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

### Community 166 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 167 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 171 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 560 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **29 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `.Reader` to `Server`, `Button`, `DB`, `newTestHarness`?**
  _High betweenness centrality (0.098) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `server.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11716975008114248 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `PxpipeConsole.tsx` to `NotificationsConsole.tsx`, `Button`, `KeysConsole.tsx`, `Input`, `SettingsConsole.tsx`, `TeamsConsole.tsx`, `ProvidersConsole.tsx`, `UsageConsole.tsx`, `CardTitle`, `PromptGuardConsole.tsx`, `extractErrorMessage`, `CardHeader`, `ProxyPoolsConsole.tsx`, `ModelEditorConsole.tsx`, `MitmConsole.tsx`, `CliToolsConsole.tsx`, `MediaConsole.tsx`, `Label`, `client.ts`?**
  _High betweenness centrality (0.095) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `Review Focus` connect `GatewayInterceptor` to `buildApp`, `Button`, `GatewayInterceptor`?**
  _High betweenness centrality (0.078) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._