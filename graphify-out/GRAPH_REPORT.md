# Graph Report - wave2c-omeg4mon-router  (2026-11-23)

## Corpus Check
- 256 files · ~257,820 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 3016 nodes · 9797 edges · 139 communities (114 shown, 25 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 930 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e0d53270`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main.go
- Module
- biome.json
- settings
- Wave 1 — Core LLM Gateway & Essential Control Plane
- APIKeySummary
- package.json
- translate_test.go
- compilerOptions
- extractErrorMessage
- Button
- testing.T
- components.json
- KeysConsole.tsx
- resolver.go
- Card
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- RingBuffer
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
- fakeProviderStore
- translate_gemini.go
- context.Context
- repository
- newTestHarness
- translate_openai.go
- lifecycle.Registry
- PromptGuardConsole.tsx
- ModelsRoutingConsole.tsx
- encoding/json.RawMessage
- net/http.ResponseWriter
- client.ts
- RouteMounts
- WriteError
- SecurityConsole.tsx
- repository
- ProxyPoolsConsole.tsx
- repository
- setupNotificationsTest
- Server
- Module
- Module
- listSettingsWithQueryer
- TestContractInterfacesSatisfied
- repository
- ModelEditorConsole.tsx
- Module
- Module
- DB
- newTestHarness
- scripts
- time.Time
- WriteJSON
- newTestHarness
- New
- crypto.SessionSigner
- Module
- Open
- CliToolsConsole.tsx
- ProviderAdapter
- repository
- Module
- repository
- batcher
- testEnv
- NewStreamTranslator
- net/url.URL
- buildApp
- TranslatorConsole.tsx
- repo
- newTestHarness
- ProvidersConsole.tsx
- stubModule
- stubSystemModule
- stubDepModule
- stubModule
- Vault
- time.Duration
- .postTokenEndpoint
- settingDefinition
- setupPromptGuardTest
- CanonicalStreamChunk
- setupTokenSaverTest
- .StatusCode
- NewStaticUIHandler
- ProvidersConsole
- stubModule
- stubModule
- Runner
- NewClient
- Bus
- Config
- stubDep
- providers/hooks.ts
- headroomCache
- .executeCandidateAttempt
- collectMigrations
- syncBuffer
- modeleditor/engine.go
- promptinjection/engine.go
- OAuthConfig
- repo
- compressor.go
- Module
- providers/schemas.test.ts
- redact.go
- .Name
- orderStubModule
- .OnRequest
- nameMinHeap

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 116 edges
2. `WriteError()` - 112 edges
3. `WriteJSON()` - 98 edges
4. `DB` - 67 edges
5. `buildApp()` - 58 edges
6. `Bus` - 47 edges
7. `Button` - 47 edges
8. `Module` - 44 edges
9. `Alert` - 43 edges
10. `AlertDescription` - 42 edges

## Surprising Connections (you probably didn't know these)
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (139 total, 25 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.05
Nodes (89): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+81 more)

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
Nodes (26): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), New(), decodeStringSlice(), encodeStringSlice(), APIKeySummary (+18 more)

### Community 6 - "package.json"
Cohesion: 0.05
Nodes (54): @biomejs/biome, class-variance-authority, clsx, geist, openapi-typescript, @radix-ui/react-label, @radix-ui/react-slot, react-dom (+46 more)

### Community 7 - "translate_test.go"
Cohesion: 0.18
Nodes (14): ProtocolTranslator, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), IsSupportedProtocol(), RewriteModelField(), TestAnthropicToOpenAIChatRoundTrip() (+6 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "extractErrorMessage"
Cohesion: 0.17
Nodes (25): NotificationsPage(), notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+17 more)

### Community 10 - "Button"
Cohesion: 0.08
Nodes (58): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), lucide-react, next, nextConfig, DashboardLayout(), DashboardLayoutProps, DashboardPage(), InitPage() (+50 more)

### Community 11 - "testing.T"
Cohesion: 0.07
Nodes (55): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback(), newChainStore() (+47 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "resolver.go"
Cohesion: 0.16
Nodes (13): resolveSnapshot, targetKey, ModelDescriptor, ModelLister, stubModelLister, stubModelLister, buildResolvedTarget(), connectionSupportsBare() (+5 more)

### Community 15 - "Card"
Cohesion: 0.11
Nodes (39): react, DashboardProviderDetailPage(), TokenSaverPage(), DashboardTranslatorPage(), useProviderConnection(), formatTimestamp(), ProviderDetailView(), toUpdateConnectionRequest() (+31 more)

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
Cohesion: 0.16
Nodes (28): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+20 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.08
Nodes (47): @hookform/resolvers, DashboardSettingsPage(), authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, AuthPolicyResponse (+39 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "RingBuffer"
Cohesion: 0.18
Nodes (12): File Map, Global Constraints, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+4 more)

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
Nodes (15): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+7 more)

### Community 39 - "net/http.Request"
Cohesion: 0.11
Nodes (16): clearSessionCookie(), isSecureRequest(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError() (+8 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.19
Nodes (18): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+10 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "fakeProviderStore"
Cohesion: 0.13
Nodes (7): fakeProviderStore, ProviderConnection, ProviderStore, stubProviderStore, UpstreamCredentials, stubProviderStore, newFakeProviderStore()

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "context.Context"
Cohesion: 0.04
Nodes (13): stubAdminAuth, stubAdminCredentialStore, stubAuthVerifier, stubAuthVerifier, AdminCredentialStore, actorFromContext(), TestSystemAdminCredentialStore(), validatePassword() (+5 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "newTestHarness"
Cohesion: 0.20
Nodes (11): fakeCredentialStore, Task 2: `auth` Module & OpenAPI Fragment, newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify() (+3 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.18
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+12 more)

### Community 49 - "lifecycle.Registry"
Cohesion: 0.16
Nodes (8): File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests, Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, normalizeName(), registryState

### Community 50 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (32): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+24 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "encoding/json.RawMessage"
Cohesion: 0.16
Nodes (23): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+15 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.14
Nodes (13): trackedResponseWriter, isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), validateHeaders(), validateName() (+5 more)

### Community 54 - "client.ts"
Cohesion: 0.12
Nodes (16): openapi-fetch, openapi-fetch@^0.14.0, systemKeys, useInitializeSystem(), apiClient, createApiClient(), ErrorDetail, ErrorEnvelope (+8 more)

### Community 55 - "RouteMounts"
Cohesion: 0.11
Nodes (7): authHarness, readSessionCookie(), Module, actorFromContext(), decodeJSONBody(), Module, RouteMounts

### Community 56 - "WriteError"
Cohesion: 0.17
Nodes (11): setSessionCookie(), decodeJSONBody(), decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), WriteError(), AdminSubject() (+3 more)

### Community 57 - "SecurityConsole.tsx"
Cohesion: 0.20
Nodes (17): SecurityPage(), securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_SEVERITIES, ClearSecurityAuditResponse (+9 more)

### Community 58 - "repository"
Cohesion: 0.10
Nodes (20): cloneStringMap(), cloneUpstreamCredentials(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+12 more)

### Community 59 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (33): DashboardProxyPoolsPage(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+25 more)

### Community 60 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 61 - "setupNotificationsTest"
Cohesion: 0.13
Nodes (19): File Map, Global Constraints, Review Focus, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, Wave 2C — Security Audit Trail, System Notifications & Background Jobs Implementation Plan, JobInspector (+11 more)

### Community 62 - "Server"
Cohesion: 0.08
Nodes (18): TestAppCloseTeardownOrderIsLIFO(), loadUI(), loadUI(), AuthVerifier, Options, apiTimeoutMiddleware(), Server, New() (+10 more)

### Community 63 - "Module"
Cohesion: 0.28
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.23
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "TestContractInterfacesSatisfied"
Cohesion: 0.08
Nodes (9): CanonicalRequest, GatewayInterceptor, ResolvedTarget, ResponseMeta, StreamEvent, stubInterceptor, stubProviderResolver, recordingInterceptor (+1 more)

### Community 67 - "repository"
Cohesion: 0.16
Nodes (8): builtinRules(), boolToInt(), defaultConfig(), newRepository(), Config, GuardEvent, repository, Rule

### Community 68 - "ModelEditorConsole.tsx"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 71 - "DB"
Cohesion: 0.09
Nodes (16): Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Global Constraints, Review Focus, Wave 2B — Gateway Interceptors (`promptinjection`, `modeleditor` / `modelbypass`, `tokensaver` / `headroom`) Implementation Plan, chainStore, newRepository(), New(), newRepository() (+8 more)

### Community 72 - "newTestHarness"
Cohesion: 0.15
Nodes (10): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+2 more)

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.15
Nodes (14): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, Options, newCooldownTracker(), New() (+6 more)

### Community 75 - "WriteJSON"
Cohesion: 0.12
Nodes (7): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), Module, WriteJSON(), JobStatusDTO

### Community 76 - "newTestHarness"
Cohesion: 0.14
Nodes (14): New(), newTestHarness(), TestCatalogAndAdapterRegistry(), TestConnectionQuickTest(), TestConnectionsCRUDAndSecretEncryptionHygiene(), TestCustomProviderNodesCRUD(), TestHeadlessOAuthFlowAndCallbackXSSProtection(), TestOAuthEnvResolutionAndCallbackJSONAndValidation() (+6 more)

### Community 77 - "New"
Cohesion: 0.15
Nodes (12): assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), New(), TestConcurrentPublishSubscribeUnsubscribe() (+4 more)

### Community 79 - "Module"
Cohesion: 0.21
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "Open"
Cohesion: 0.15
Nodes (14): doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), buildDSN(), Open() (+6 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "ProviderAdapter"
Cohesion: 0.10
Nodes (5): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable(), customTestAdapter, ProviderAdapter

### Community 83 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 84 - "Module"
Cohesion: 0.14
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 85 - "repository"
Cohesion: 0.17
Nodes (6): boolToInt(), newRepository(), BindingSummary, PoolSummary, repository, storedPool

### Community 86 - "batcher"
Cohesion: 0.14
Nodes (6): newBatcher(), New(), newPricingStore(), batcher, ModelPricing, pricingStore

### Community 87 - "testEnv"
Cohesion: 0.12
Nodes (9): KeyInfo, stubKeyValidator, stubUsageRecorder, UsageRecord, UsageRecorder, stubKeyValidator, stubResolver, stubUsageRecorder (+1 more)

### Community 88 - "NewStreamTranslator"
Cohesion: 0.18
Nodes (10): geminiPendingTool, geminiStreamEncoder, NewStreamTranslator(), parseJSONObjectOrWrap(), newAnthropicStreamEncoder(), mapCanonicalStopReasonToGemini(), newGeminiStreamEncoder(), newOpenAIStreamDecoder() (+2 more)

### Community 89 - "net/url.URL"
Cohesion: 0.14
Nodes (14): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan() (+6 more)

### Community 90 - "buildApp"
Cohesion: 0.07
Nodes (34): buildApp(), main(), runMain(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd() (+26 more)

### Community 91 - "TranslatorConsole.tsx"
Cohesion: 0.16
Nodes (20): translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject(), SAMPLE_PAYLOADS (+12 more)

### Community 92 - "repo"
Cohesion: 0.15
Nodes (10): buildWhere(), formatUTC(), newNotificationID(), newRepo(), normalizeSeverity(), scanNotification(), scanNotificationRow(), ListFilter (+2 more)

### Community 93 - "newTestHarness"
Cohesion: 0.19
Nodes (12): stubAdminAuth, testHarness, decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering() (+4 more)

### Community 94 - "ProvidersConsole.tsx"
Cohesion: 0.12
Nodes (17): ActiveOAuthSession, CatalogAdapterDescriptor, CreateConnectionFormInput, CreateConnectionFormOutput, CreateNodeFormInput, CreateNodeFormOutput, DEFAULT_CREATE_CONNECTION_VALUES, DEFAULT_CREATE_NODE_VALUES (+9 more)

### Community 102 - "Vault"
Cohesion: 0.14
Nodes (14): ProviderResolver, File Map, Global Constraints, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators (+6 more)

### Community 104 - "time.Duration"
Cohesion: 0.13
Nodes (5): AccountHealthTracker, stubAccountHealthTracker, failureRecord, stubHealthTracker, computeBackoff()

### Community 106 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 107 - "setupPromptGuardTest"
Cohesion: 0.31
Nodes (7): Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins(), TestPromptInjectionCustomRuleValidationAndRE2(), TestPromptInjectionFinalReviewFindings(), testEnv

### Community 108 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 109 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 111 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 112 - "ProvidersConsole"
Cohesion: 0.13
Nodes (16): DashboardProvidersPage(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation(), useExchangeProviderOAuthMutation(), useProviderConnections(), useProvidersCatalog(), useStartProviderOAuthMutation(), useTestProviderConnectionMutation() (+8 more)

### Community 115 - "Runner"
Cohesion: 0.25
Nodes (6): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter

### Community 116 - "NewClient"
Cohesion: 0.14
Nodes (14): File Map, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Options, Options, applyDefaults() (+6 more)

### Community 117 - "Bus"
Cohesion: 0.11
Nodes (15): Event, New(), New(), New(), decodeJSONBody(), New(), Bus, Options (+7 more)

### Community 118 - "Config"
Cohesion: 0.11
Nodes (17): Options, Task 2: `security` Module (Security Posture & Redacted Audit Trail), New(), setupSecurityTest(), TestSecurityEventBusAutoAuditIngestion(), TestSecurityPostureAndExportAndPrune(), TestSecurityRedactionAndTamperEvidentClear(), TestSecurityWave2BTopicParityAndExtendedRedaction() (+9 more)

### Community 120 - "providers/hooks.ts"
Cohesion: 0.14
Nodes (16): ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useDeleteProviderConnectionMutation(), useProviderNodes(), CreateProviderConnectionRequest (+8 more)

### Community 121 - "headroomCache"
Cohesion: 0.24
Nodes (3): newHeadroomCache(), headroomCache, WarmPrefixSummary

### Community 122 - ".executeCandidateAttempt"
Cohesion: 0.19
Nodes (14): attemptOutcome, CanonicalChatResponse, buildUpstreamURL(), extractUpstreamErrorMessage(), parseRetryAfter(), EncodeResponse(), rawJSONToString(), buildCompactedResponse() (+6 more)

### Community 123 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 126 - "modeleditor/engine.go"
Cohesion: 0.22
Nodes (10): copyResponseHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), wrapSystemPrompt(), wrapText() (+2 more)

### Community 128 - "promptinjection/engine.go"
Cohesion: 0.30
Nodes (13): actionRank(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey(), minSeverityForSensitivity() (+5 more)

### Community 129 - "OAuthConfig"
Cohesion: 0.27
Nodes (6): cloneOAuthConfig(), cloneStringSlice(), OAuthConfig, oauthEnvPrefix(), buildAuthorizeURL(), builtinAdapter

### Community 130 - "repo"
Cohesion: 0.24
Nodes (7): buildWhere(), newEventID(), newRepo(), scanAuditEvent(), AuditEvent, ListFilter, repo

### Community 131 - "compressor.go"
Cohesion: 0.36
Nodes (9): compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), leadingBackticks(), truncateToolOutput() (+1 more)

### Community 132 - "Module"
Cohesion: 0.29
Nodes (3): Module, refreshCall, storedCredentials

### Community 133 - "providers/schemas.test.ts"
Cohesion: 0.22
Nodes (8): createConnectionSchema, createNodeSchema, oauthExchangeSchema, oauthStartSchema, toCreateConnectionRequest(), toCreateNodeRequest(), toOAuthExchangeRequest(), updateConnectionSchema

### Community 134 - "redact.go"
Cohesion: 0.46
Nodes (6): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), normalizeSeverity()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 472 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `buildApp` to `main.go`, `Button`, `WriteJSON`, `newTestHarness`?**
  _High betweenness centrality (0.164) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05467232809979922 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `buildApp`?**
  _High betweenness centrality (0.124) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `RingBuffer` to `ProviderAdapter`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._