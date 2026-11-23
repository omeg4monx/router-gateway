# Graph Report - wave2a-omeg4mon-router  (2026-11-23)

## Corpus Check
- 186 files · ~199,770 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 2353 nodes · 7390 edges · 106 communities (86 shown, 20 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 729 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2e77854c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main.go
- Module
- biome.json
- settings
- Wave 1 — Core LLM Gateway & Essential Control Plane
- keys/repo.go
- package.json
- Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan
- compilerOptions
- Server
- Button
- testing.T
- components.json
- KeysConsole.tsx
- repository
- catalog.go
- System Module & OpenAPI 3.1 Spec
- Bus
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- resolver.go
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
- net/http.ResponseWriter
- UsageConsole.tsx
- encoding/json.RawMessage
- Server
- UpstreamCredentials
- translate_gemini.go
- Module
- repository
- context.Context
- translate_openai.go
- buildApp
- Config
- ModelsRoutingConsole.tsx
- translate_responses.go
- WriteError
- ProviderAdapter
- Module
- RouteMounts
- CanonicalStreamChunk
- repository
- updateConnectionTx
- ProviderDetailView.tsx
- NewStreamTranslator
- Module
- Module
- Module
- listSettingsWithQueryer
- TestContractInterfacesSatisfied
- newTestHarness
- scripts
- ProvidersConsole.tsx
- responsesStreamEncoder
- DB
- .AuthenticateAdmin
- CanonicalChatResponse
- time.Time
- repository
- stubModule
- ProvidersConsole
- providers/hooks.ts
- teeHandler
- newTestHarness
- CliToolsConsole.tsx
- Module
- Module
- File Map
- Runner
- newTestHarness
- newTestHarness
- settingDefinition
- nameMinHeap
- providers/schemas.test.ts
- NewStaticUIHandler
- stubModule
- orderStubModule
- stubModule
- NewClient
- net/url.URL
- Load
- stubSystemModule
- stubDepModule
- stubModule
- collectMigrations
- trackedResponseWriter
- tokens.go
- .postTokenEndpoint

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 75 edges
2. `WriteJSON()` - 64 edges
3. `extractErrorMessage()` - 62 edges
4. `DB` - 50 edges
5. `buildApp()` - 47 edges
6. `Module` - 37 edges
7. `ModelsRoutingConsole()` - 36 edges
8. `Module` - 35 edges
9. `Module` - 34 edges
10. `ProvidersConsole()` - 34 edges

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

## Communities (106 total, 20 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.05
Nodes (85): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+77 more)

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

### Community 5 - "keys/repo.go"
Cohesion: 0.05
Nodes (36): stubKeyValidator, File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+28 more)

### Community 6 - "package.json"
Cohesion: 0.04
Nodes (56): @biomejs/biome, class-variance-authority, clsx, geist, openapi-typescript, @radix-ui/react-label, @radix-ui/react-slot, react-dom (+48 more)

### Community 7 - "Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan"
Cohesion: 0.20
Nodes (12): ProxyDialerResolver, File Map, Review Focus, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification (+4 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.19
Nodes (14): loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), New(), TestAdminAPIAuthEnforcement(), TestDevUIHandlerFallback() (+6 more)

### Community 10 - "Button"
Cohesion: 0.10
Nodes (52): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, DashboardPage() (+44 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (58): TestAppCloseTeardownOrderIsLIFO(), TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback() (+50 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "repository"
Cohesion: 0.15
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.16
Nodes (17): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+9 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.13
Nodes (35): CacheControl, CanonicalChatRequest, CanonicalMessage, ResponseFormat, SystemBlock, ToolChoice, ToolDefinition, CanFastPath() (+27 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): zod, zod@^4.6.7, DashboardSettingsPage(), authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+36 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.13
Nodes (15): resolveSnapshot, targetKey, stubModelLister, stubProviderResolver, stubModelLister, ModelDescriptor, ProviderConnection, ResolvedTarget (+7 more)

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
Cohesion: 0.23
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
Cohesion: 0.12
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+9 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.09
Nodes (19): attemptOutcome, statusError, clearSessionCookie(), isSecureRequest(), readSessionCookie(), anthropicErrorType(), extractAPIKey(), geminiErrorStatus() (+11 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (48): @hookform/resolvers, DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation() (+40 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.16
Nodes (24): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+16 more)

### Community 42 - "Server"
Cohesion: 0.14
Nodes (3): testHarness, Server, App

### Community 43 - "UpstreamCredentials"
Cohesion: 0.20
Nodes (7): stubProviderStore, stubProviderStore, UpstreamCredentials, TestContractInterfacesSatisfied(), cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "Module"
Cohesion: 0.08
Nodes (3): Module, stubSystemModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.15
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "context.Context"
Cohesion: 0.06
Nodes (9): stubAdminCredentialStore, actorFromContext(), boolToInt(), validatePassword(), stubModule, BindingSummary, PoolSummary, repository (+1 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.17
Nodes (21): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+13 more)

### Community 49 - "buildApp"
Cohesion: 0.18
Nodes (11): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer() (+3 more)

### Community 50 - "Config"
Cohesion: 0.12
Nodes (19): Options, Task 2: `auth` Module & OpenAPI Fragment, New(), BuildInfo, New(), newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion() (+11 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.05
Nodes (82): openapi-fetch, @tanstack/react-query, DashboardModelsPage(), AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, useLoginMutation() (+74 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.21
Nodes (16): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+8 more)

### Community 53 - "WriteError"
Cohesion: 0.11
Nodes (24): setSessionCookie(), decodeJSONBody(), buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody() (+16 more)

### Community 54 - "ProviderAdapter"
Cohesion: 0.10
Nodes (7): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), OAuthConfig, buildAuthorizeURL(), builtinAdapter, customTestAdapter, ProviderAdapter

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 58 - "repository"
Cohesion: 0.15
Nodes (9): decodeStringMap(), decodeStringSlice(), scanConnectionRecord(), scanNode(), listNodesResponse, oauthSessionRecord, ProviderNode, repository (+1 more)

### Community 59 - "updateConnectionTx"
Cohesion: 0.19
Nodes (14): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), oauthEnvPrefix(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+6 more)

### Community 60 - "ProviderDetailView.tsx"
Cohesion: 0.21
Nodes (18): DashboardProviderDetailPage(), useProviderConnection(), formatTimestamp(), ProviderDetailView(), toUpdateConnectionRequest(), cn(), Badge(), BadgeProps (+10 more)

### Community 61 - "NewStreamTranslator"
Cohesion: 0.20
Nodes (10): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder() (+2 more)

### Community 62 - "Module"
Cohesion: 0.17
Nodes (7): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Module, New(), newRepository(), Vault, Options, storedCredentials

### Community 63 - "Module"
Cohesion: 0.25
Nodes (9): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+1 more)

### Community 64 - "Module"
Cohesion: 0.30
Nodes (4): generateInput, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.19
Nodes (7): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "TestContractInterfacesSatisfied"
Cohesion: 0.14
Nodes (9): GatewayInterceptor, stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, copyResponseHeaders(), New() (+1 more)

### Community 67 - "newTestHarness"
Cohesion: 0.16
Nodes (14): fakeProviderStore, stubAdminAuth, testHarness, decodeJSON(), newFakeProviderStore(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency() (+6 more)

### Community 68 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 69 - "ProvidersConsole.tsx"
Cohesion: 0.12
Nodes (17): ActiveOAuthSession, CatalogAdapterDescriptor, CreateConnectionFormInput, CreateConnectionFormOutput, CreateNodeFormInput, CreateNodeFormOutput, DEFAULT_CREATE_CONNECTION_VALUES, DEFAULT_CREATE_NODE_VALUES (+9 more)

### Community 70 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 71 - "DB"
Cohesion: 0.09
Nodes (24): Options, Event, stubResolver, testEnv, newCooldownTracker(), New(), TestModuleMetadataAndStartValidation(), newRepository() (+16 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (6): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "CanonicalChatResponse"
Cohesion: 0.52
Nodes (7): CanonicalChatResponse, EncodeResponse(), rawJSONToString(), encodeAnthropicResponse(), encodeGeminiResponse(), encodeOpenAIChatResponse(), encodeOpenAIResponsesResponse()

### Community 74 - "time.Time"
Cohesion: 0.05
Nodes (18): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, stubAccountHealthTracker, stubUsageRecorder, Claims (+10 more)

### Community 75 - "repository"
Cohesion: 0.09
Nodes (18): main(), runMain(), Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, File Map, Global Constraints, Review Focus (+10 more)

### Community 77 - "ProvidersConsole"
Cohesion: 0.13
Nodes (16): DashboardProvidersPage(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation(), useExchangeProviderOAuthMutation(), useProviderConnections(), useProvidersCatalog(), useStartProviderOAuthMutation(), useTestProviderConnectionMutation() (+8 more)

### Community 78 - "providers/hooks.ts"
Cohesion: 0.14
Nodes (16): ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useDeleteProviderConnectionMutation(), useProviderNodes(), CreateProviderConnectionRequest (+8 more)

### Community 79 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 80 - "newTestHarness"
Cohesion: 0.24
Nodes (9): decodeJSON(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageSummaryLogsAndRetentionPrune(), waitForLogCount() (+1 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "Module"
Cohesion: 0.46
Nodes (3): actorFromContext(), decodeJSONBody(), Module

### Community 84 - "File Map"
Cohesion: 0.28
Nodes (8): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

### Community 86 - "newTestHarness"
Cohesion: 0.14
Nodes (14): newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation() (+6 more)

### Community 87 - "newTestHarness"
Cohesion: 0.27
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 88 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 90 - "providers/schemas.test.ts"
Cohesion: 0.22
Nodes (8): createConnectionSchema, createNodeSchema, oauthExchangeSchema, oauthStartSchema, toCreateConnectionRequest(), toCreateNodeRequest(), toOAuthExchangeRequest(), updateConnectionSchema

### Community 91 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 96 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 97 - "net/url.URL"
Cohesion: 0.18
Nodes (12): stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), NewProviderProxyFunc() (+4 more)

### Community 98 - "Load"
Cohesion: 0.33
Nodes (7): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 102 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 104 - "tokens.go"
Cohesion: 0.67
Nodes (3): estimateBlocks(), EstimateTokens(), TestEstimateTokens()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 401 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **20 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `keys/repo.go` to `Module`, `Server`, `Config`, `Button`?**
  _High betweenness centrality (0.198) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05422049956933678 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `keys/repo.go`?**
  _High betweenness centrality (0.139) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `Button` to `ModelsRoutingConsole.tsx`, `ProviderDetailView.tsx`, `keys/repo.go`, `package.json`?**
  _High betweenness centrality (0.116) - this node is a cross-community bridge._