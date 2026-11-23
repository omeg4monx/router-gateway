# Graph Report - wave3d-omeg4mon-router  (2026-11-23)

## Corpus Check
- 350 files · ~339,769 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 3983 nodes · 13488 edges · 177 communities (152 shown, 25 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1268 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `71d2da11`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- contracts.go
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
- .ProvisionKey
- TokenSaverConsole.tsx
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- Card
- steps
- web-check
- on
- cross-compile
- Module
- providers/schemas.ts
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
- pricingStore
- CardTitle
- translate_openai.go
- lifecycle.Registry
- PromptGuardConsole.tsx
- ModelsRoutingConsole.tsx
- encoding/json.RawMessage
- net/http.ResponseWriter
- go_pkg_context
- Module
- Module
- Label
- repository
- ProxyPoolsConsole.tsx
- repository
- sync.Mutex
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
- setupModelEditorTest
- CliToolsConsole.tsx
- ProviderAdapter
- repository
- MediaConsole.tsx
- Module
- WriteJSON
- Vault
- Module
- net/url.URL
- buildApp
- TranslatorConsole.tsx
- repo
- app/layout.tsx
- RingBuffer
- repo
- go_pkg_strings
- extractErrorMessage
- ProvidersConsole.tsx
- TunnelConsole.tsx
- io/fs.FS
- .Reader
- ProvidersConsole
- settingDefinition
- setupPromptGuardTest
- Module
- setupTokenSaverTest
- .StatusCode
- Module
- repository
- repo
- Module
- .ListModels
- Options
- system/module.go
- Config
- setupTestMitm
- setupMediaTest
- providers/hooks.ts
- providers/handlers.go
- setupNotificationsTest
- Module
- apply.go
- promptinjection/engine.go
- providers/repo.go
- CanonicalStreamChunk
- certManager
- repo
- newTestHarness
- DB
- modeleditor/engine.go
- stubHealthTracker
- repository
- RouteMounts
- main.go
- headroomCache
- stubCompleter
- Module
- transform.go
- time.Duration
- .Record
- stubModule
- redact.go
- go_pkg_encoding_json
- stubModule
- repo
- stubModule
- .executeBuiltinFetch
- crypto.SessionSigner
- stubSystemModule
- newTestHarness
- providers/schemas.test.ts
- Load
- .ValidateKey
- stubDep
- .buildProxyHandler
- .list
- NewStreamTranslator
- compressor.go
- headroom.go
- runConcurrentScan
- Config
- repository
- orderStubModule
- tgTestEnv
- net/http.Header
- mapBusEventToAudit
- .postTokenEndpoint
- nullableString
- pricingListResponse
- zustand

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 192 edges
2. `WriteError()` - 174 edges
3. `WriteJSON()` - 171 edges
4. `DB` - 91 edges
5. `buildApp()` - 74 edges
6. `Bus` - 72 edges
7. `Button` - 65 edges
8. `Alert` - 61 edges
9. `AlertDescription` - 60 edges
10. `AlertTitle` - 60 edges

## Surprising Connections (you probably didn't know these)
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/platform/lifecycle/registry.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (177 total, 25 thin omitted)

### Community 0 - "contracts.go"
Cohesion: 0.15
Nodes (23): ModelOption, ChangedEvent, createComboRequest, deleteDisabledModelRequest, patchComboRequest, upsertAliasRequest, upsertDisabledModelRequest, adminSubjectKey (+15 more)

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
Nodes (23): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary() (+15 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - ".executeCandidateAttempt"
Cohesion: 0.10
Nodes (28): attemptOutcome, CanonicalChatResponse, ProtocolTranslator, buildUpstreamURL(), extractUpstreamErrorMessage(), parseRetryAfter(), CanFastPath(), DecodeRequest() (+20 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "NotificationsConsole.tsx"
Cohesion: 0.15
Nodes (24): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 10 - "Button"
Cohesion: 0.08
Nodes (63): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, class-variance-authority, lucide-react, next (+55 more)

### Community 11 - "testing.T"
Cohesion: 0.04
Nodes (103): TestAppCloseTeardownOrderIsLIFO(), TestSentinelErrorsDistinct(), doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), decodeJSON(), newTestHarness() (+95 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (29): @hookform/resolvers, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - ".ProvisionKey"
Cohesion: 0.25
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 15 - "TokenSaverConsole.tsx"
Cohesion: 0.20
Nodes (18): TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation() (+10 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.15
Nodes (17): Options, testHarness, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore() (+9 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.15
Nodes (32): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+24 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.05
Nodes (63): @tanstack/react-query, @tanstack/react-query@^5.76.1, DashboardSettingsPage(), authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+55 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "Card"
Cohesion: 0.07
Nodes (59): ArenaPage(), ChatPage(), DashboardProviderDetailPage(), ArenaConsole(), ChatConsole(), arenaQueryKeys, useArenaBattles(), useArenaLeaderboard() (+51 more)

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
Cohesion: 0.19
Nodes (9): actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription(), validateModelIdentifier() (+1 more)

### Community 30 - "providers/schemas.ts"
Cohesion: 0.08
Nodes (26): CatalogOAuthDescriptor, commaListSchema(), DEFAULT_PRIORITY, DEFAULT_UPDATE_CONNECTION_VALUES, headersTextSchema, isValidHttpUrl(), ListProviderConnectionsResponse, ListProviderNodesResponse (+18 more)

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
Cohesion: 0.13
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+8 more)

### Community 39 - "net/http.Request"
Cohesion: 0.10
Nodes (17): clearSessionCookie(), isSecureRequest(), readSessionCookie(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus(), openAIErrorType(), writeGatewayError() (+9 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.24
Nodes (13): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+5 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "TestContractInterfacesSatisfied"
Cohesion: 0.12
Nodes (8): fakeProviderStore, stubProviderStore, stubProviderStore, ProviderConnection, UpstreamCredentials, TestContractInterfacesSatisfied(), newFakeProviderStore(), stubProviderStore

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "context.Context"
Cohesion: 0.05
Nodes (12): fakeCredentialStore, AdminCredentialStore, stubAdminCredentialStore, actorFromContext(), validatePassword(), formatUTC(), newTeamID(), stubCredentialStore (+4 more)

### Community 46 - "pricingStore"
Cohesion: 0.10
Nodes (13): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ModelPricing, pricingStore, ProviderBreakdown (+5 more)

### Community 47 - "CardTitle"
Cohesion: 0.09
Nodes (40): zod, TeamsPage(), VouchersPage(), useApiKeys(), teamsQueryKeys, useAssignTeamKeyMutation(), useCreateTeamMutation(), useDeleteTeamMutation() (+32 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "lifecycle.Registry"
Cohesion: 0.18
Nodes (3): normalizeName(), nameMinHeap, registryState

### Community 50 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "encoding/json.RawMessage"
Cohesion: 0.17
Nodes (22): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+14 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.07
Nodes (30): trackedResponseWriter, setSessionCookie(), decodeJSONBody(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody() (+22 more)

### Community 54 - "go_pkg_context"
Cohesion: 0.23
Nodes (4): targetKey, ToolDescriptor, ChangedEvent, storedSetting

### Community 55 - "Module"
Cohesion: 0.16
Nodes (5): decodeStrictJSON(), LogEntry, LogListResponse, Module, Status

### Community 56 - "Module"
Cohesion: 0.23
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 57 - "Label"
Cohesion: 0.18
Nodes (18): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+10 more)

### Community 58 - "repository"
Cohesion: 0.14
Nodes (9): decodeStringMap(), decodeStringSlice(), newRepository(), scanConnectionRecord(), scanNode(), oauthSessionRecord, ProviderNode, repository (+1 more)

### Community 59 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (33): DashboardProxyPoolsPage(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 60 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 61 - "sync.Mutex"
Cohesion: 0.09
Nodes (10): TeamAssigner, JobStatus, Module, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter (+2 more)

### Community 62 - "Server"
Cohesion: 0.07
Nodes (26): loadUI(), loadUI(), AuthVerifier, testEnv, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+18 more)

### Community 63 - "Module"
Cohesion: 0.23
Nodes (8): actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken(), Module

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.23
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "Module"
Cohesion: 0.06
Nodes (6): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, Module

### Community 67 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 68 - "ModelEditorConsole.tsx"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "TelegramConsole.tsx"
Cohesion: 0.12
Nodes (31): TelegramPage(), telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation() (+23 more)

### Community 71 - "MitmConsole.tsx"
Cohesion: 0.10
Nodes (36): MitmPage(), mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents() (+28 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (7): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.18
Nodes (10): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, newCooldownTracker(), newFakeClock(), newFakeClock() (+2 more)

### Community 75 - "Module"
Cohesion: 0.13
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), mapBusEventToNotification(), toJobDTO(), JobStatusDTO

### Community 76 - "Module"
Cohesion: 0.06
Nodes (15): streamResponseBody(), newID(), nowUTC(), BinariesStatus, BinaryStatus, ConfigRecord, Event, fakeHandle (+7 more)

### Community 77 - "dependencies"
Cohesion: 0.10
Nodes (20): clsx, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, dependencies, class-variance-authority@^0.7.1 (+12 more)

### Community 78 - "Module"
Cohesion: 0.22
Nodes (3): decodeStrictJSON(), extractUsageSpend(), Module

### Community 79 - "Module"
Cohesion: 0.22
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "ProviderAdapter"
Cohesion: 0.09
Nodes (10): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), cloneOAuthConfig(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), TestWave2ProviderAdaptersRegisteredAndTestable(), buildAuthorizeURL(), builtinAdapter (+2 more)

### Community 83 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 84 - "MediaConsole.tsx"
Cohesion: 0.12
Nodes (31): MediaPage(), mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders() (+23 more)

### Community 85 - "Module"
Cohesion: 0.18
Nodes (4): clientIP(), decodeStrictJSON(), Module, truncateRunes()

### Community 86 - "WriteJSON"
Cohesion: 0.14
Nodes (10): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider() (+2 more)

### Community 87 - "Vault"
Cohesion: 0.07
Nodes (27): File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain (+19 more)

### Community 88 - "Module"
Cohesion: 0.14
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "net/url.URL"
Cohesion: 0.19
Nodes (9): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, probeProxy(), proxyDialAddress(), NewProviderProxyFunc(), ProviderIDFromContext(), TestProviderProxyFuncContextRouting() (+1 more)

### Community 90 - "buildApp"
Cohesion: 0.07
Nodes (34): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+26 more)

### Community 91 - "TranslatorConsole.tsx"
Cohesion: 0.14
Nodes (23): react-hook-form, DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse (+15 more)

### Community 92 - "repo"
Cohesion: 0.17
Nodes (11): formatUTC(), generateCode(), newID(), normalizeCode(), scanVoucher(), scanVoucherRow(), createVoucherInput, patchVoucherInput (+3 more)

### Community 93 - "app/layout.tsx"
Cohesion: 0.24
Nodes (7): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 94 - "RingBuffer"
Cohesion: 0.18
Nodes (12): File Map, Global Constraints, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+4 more)

### Community 96 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 97 - "go_pkg_strings"
Cohesion: 0.12
Nodes (8): PolicyChangedEvent, generateResponse, applyResponse, catalogResponse, rollbackResponse, isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive()

### Community 99 - "extractErrorMessage"
Cohesion: 0.13
Nodes (32): PxpipePage(), pxpipeQueryKeys, useClearPxpipeLogsMutation(), usePreviewPxpipeMutation(), usePxpipeHealth(), usePxpipeLogs(), usePxpipeStats(), usePxpipeStatus() (+24 more)

### Community 100 - "ProvidersConsole.tsx"
Cohesion: 0.12
Nodes (17): ActiveOAuthSession, CatalogAdapterDescriptor, CreateConnectionFormInput, CreateConnectionFormOutput, CreateNodeFormInput, CreateNodeFormOutput, DEFAULT_CREATE_CONNECTION_VALUES, DEFAULT_CREATE_NODE_VALUES (+9 more)

### Community 102 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): TunnelPage(), tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 103 - "io/fs.FS"
Cohesion: 0.04
Nodes (16): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler() (+8 more)

### Community 104 - ".Reader"
Cohesion: 0.05
Nodes (38): main(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+30 more)

### Community 105 - "ProvidersConsole"
Cohesion: 0.13
Nodes (16): DashboardProvidersPage(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation(), useExchangeProviderOAuthMutation(), useProviderConnections(), useProvidersCatalog(), useStartProviderOAuthMutation(), useTestProviderConnectionMutation() (+8 more)

### Community 106 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 107 - "setupPromptGuardTest"
Cohesion: 0.14
Nodes (16): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests, Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders() (+8 more)

### Community 108 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), isLoopbackHost(), validateLoopbackListenAddr(), validateLoopbackURL(), validToolID(), Module

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

### Community 115 - ".ListModels"
Cohesion: 0.10
Nodes (17): stubModelLister, resolveSnapshot, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, ModelDescriptor, ProviderResolver (+9 more)

### Community 116 - "Options"
Cohesion: 0.10
Nodes (25): Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI), Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan, Options, Options, ChatCompleter (+17 more)

### Community 117 - "system/module.go"
Cohesion: 0.14
Nodes (10): createAliasRequest, EventListResponse, generateCARequest, updateConfigRequest, healthResponse, initCompleteResponse, initRequest, initStatusResponse (+2 more)

### Community 118 - "Config"
Cohesion: 0.10
Nodes (22): Options, Task 2: `auth` Module & OpenAPI Fragment, Task 2: `security` Module (Security Posture & Redacted Audit Trail), New(), setupSecurityTest(), TestSecurityEventBusAutoAuditIngestion(), TestSecurityPostureAndExportAndPrune(), TestSecurityRedactionAndTamperEvidentClear() (+14 more)

### Community 119 - "setupTestMitm"
Cohesion: 0.31
Nodes (10): New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly(), TestMitmLiveTLSProxyInterceptAndAliasRewrite() (+2 more)

### Community 120 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 121 - "providers/hooks.ts"
Cohesion: 0.14
Nodes (16): ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useDeleteProviderConnectionMutation(), useProviderNodes(), CreateProviderConnectionRequest (+8 more)

### Community 122 - "providers/handlers.go"
Cohesion: 0.19
Nodes (11): catalogAdapterDescriptor, catalogOAuthDescriptor, catalogResponse, createConnectionRequest, createNodeRequest, listNodesResponse, oauthExchangeRequest, oauthStartRequest (+3 more)

### Community 123 - "setupNotificationsTest"
Cohesion: 0.16
Nodes (15): File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict() (+7 more)

### Community 124 - "Module"
Cohesion: 0.07
Nodes (5): authHarness, Module, TestRegistryStartAllRollsBackOnPartialFailure(), stubSystemModule, stubSystemModule

### Community 126 - "apply.go"
Cohesion: 0.43
Nodes (5): applyConfigFile(), nextBackupPath(), rollbackConfigFile(), writeFileAtomic(), writeAtomicFile()

### Community 128 - "promptinjection/engine.go"
Cohesion: 0.20
Nodes (14): actionRank(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity() (+6 more)

### Community 129 - "providers/repo.go"
Cohesion: 0.25
Nodes (13): cloneStringMap(), cloneStringSlice(), cloneUpstreamCredentials(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), updateConnectionTx() (+5 more)

### Community 131 - "certManager"
Cohesion: 0.14
Nodes (8): newCertManager(), randomSerial(), buildInstallGuide(), cachedLeaf, CAInfo, certManager, InstallGuide, Status

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "newTestHarness"
Cohesion: 0.27
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 134 - "DB"
Cohesion: 0.05
Nodes (48): Options, Options, testHarness, New(), newRepo(), New(), TestModuleMetadataAndStartValidation(), newRepository() (+40 more)

### Community 135 - "modeleditor/engine.go"
Cohesion: 0.27
Nodes (9): applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt(), wrapText() (+1 more)

### Community 136 - "stubHealthTracker"
Cohesion: 0.17
Nodes (4): stubAccountHealthTracker, failureRecord, stubHealthTracker, AccountHealthTracker

### Community 137 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 139 - "main.go"
Cohesion: 0.07
Nodes (9): errorDisposition, providerIDContextKey, classifyStatus(), isClientPayloadErrorBody(), createComboReq, createProviderReq, updateComboReq, updateProviderReq (+1 more)

### Community 141 - "stubCompleter"
Cohesion: 0.33
Nodes (6): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent()

### Community 142 - "Module"
Cohesion: 0.29
Nodes (3): Module, refreshCall, storedCredentials

### Community 143 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - ".Record"
Cohesion: 0.19
Nodes (5): stubUsageRecorder, stubUsageRecorder, UsageRecord, UsageRecorder, recordingUsage

### Community 147 - "redact.go"
Cohesion: 0.18
Nodes (12): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), newEventID(), normalizeSeverity() (+4 more)

### Community 148 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (19): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+11 more)

### Community 150 - "repo"
Cohesion: 0.22
Nodes (7): formatUTC(), newNotificationID(), normalizeSeverity(), scanNotificationRow(), EmitPayload, Notification, repo

### Community 152 - ".executeBuiltinFetch"
Cohesion: 0.29
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 155 - "newTestHarness"
Cohesion: 0.25
Nodes (9): newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation() (+1 more)

### Community 156 - "providers/schemas.test.ts"
Cohesion: 0.22
Nodes (8): createConnectionSchema, createNodeSchema, oauthExchangeSchema, oauthStartSchema, toCreateConnectionRequest(), toCreateNodeRequest(), toOAuthExchangeRequest(), updateConnectionSchema

### Community 157 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 158 - ".ValidateKey"
Cohesion: 0.36
Nodes (4): stubKeyValidator, stubKeyValidator, KeyInfo, stubKeyValidator

### Community 160 - ".buildProxyHandler"
Cohesion: 0.29
Nodes (3): matchToolForHost(), newID(), singleJoiningSlash()

### Community 161 - ".list"
Cohesion: 0.50
Nodes (3): buildWhere(), scanNotification(), ListFilter

### Community 162 - "NewStreamTranslator"
Cohesion: 0.16
Nodes (11): streamResult, StreamTranslator, NewStreamTranslator(), captureResponsesFrame(), mergeFastPathSSEUsage(), parseSSEFrame(), newAnthropicStreamEncoder(), newGeminiStreamEncoder() (+3 more)

### Community 163 - "compressor.go"
Cohesion: 0.36
Nodes (9): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), leadingBackticks(), truncateToolOutput() (+1 more)

### Community 164 - "headroom.go"
Cohesion: 0.33
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 165 - "runConcurrentScan"
Cohesion: 0.40
Nodes (6): parseAndValidateProxyURL(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL, ScanCandidateResult, ScanResponse

### Community 166 - "Config"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 167 - "repository"
Cohesion: 0.33
Nodes (4): boolToInt(), newRepository(), Config, repository

### Community 169 - "tgTestEnv"
Cohesion: 0.50
Nodes (3): arenaTestEnv, tgAPICall, tgTestEnv

### Community 170 - "net/http.Header"
Cohesion: 0.18
Nodes (5): memoryResponseWriter, newMemoryResponseWriter(), copyResponseHeaders(), copyNonHopHeaders(), syncBuffer

### Community 171 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 590 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `.Reader` to `buildApp`, `Button`, `Config`, `RouteMounts`?**
  _High betweenness centrality (0.110) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `extractErrorMessage()` connect `extractErrorMessage` to `NotificationsConsole.tsx`, `Button`, `KeysConsole.tsx`, `TokenSaverConsole.tsx`, `SettingsConsole.tsx`, `Card`, `UsageConsole.tsx`, `CardTitle`, `PromptGuardConsole.tsx`, `ModelsRoutingConsole.tsx`, `Label`, `ProxyPoolsConsole.tsx`, `ModelEditorConsole.tsx`, `TelegramConsole.tsx`, `MitmConsole.tsx`, `CliToolsConsole.tsx`, `MediaConsole.tsx`, `TranslatorConsole.tsx`, `TunnelConsole.tsx`, `ProvidersConsole`, `providers/hooks.ts`?**
  _High betweenness centrality (0.108) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._
- **Why does `SystemOverview()` connect `Button` to `Card`, `SettingsConsole.tsx`, `CardTitle`?**
  _High betweenness centrality (0.106) - this node is a cross-community bridge._
- **Should `Wave 1 — Core LLM Gateway & Essential Control Plane` be split into smaller, more focused modules?**
  _Cohesion score 0.1341991341991342 - nodes in this community are weakly interconnected._