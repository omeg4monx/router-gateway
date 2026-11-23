# Graph Report - wave2b-omeg4mon-router  (2026-11-23)

## Corpus Check
- 212 files · ~223,774 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 2630 nodes · 8427 edges · 116 communities (93 shown, 23 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 834 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3edba52f`
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
- contracts.ProxyDialerResolver
- compilerOptions
- Server
- Button
- testing.T
- components.json
- KeysConsole.tsx
- resolver.go
- stubModule
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- engine.go
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
- RingBuffer
- TestContractInterfacesSatisfied
- translate_gemini.go
- Module
- repository
- context.Context
- translate_openai.go
- buildApp
- Config
- ModelsRoutingConsole.tsx
- translate_responses.go
- net/http.ResponseWriter
- ProviderDetailView.tsx
- WriteJSON
- Module
- CanonicalStreamChunk
- providers/repo.go
- extractErrorMessage
- modeleditor/repo.go
- client.ts
- newTestHarness
- Module
- net/http.ResponseWriter
- listSettingsWithQueryer
- TestContractInterfacesSatisfied
- teeHandler
- translate_test.go
- Module
- Module
- DB
- .AuthenticateAdmin
- .prepareUpstreamBody
- time.Time
- repository
- File Map
- pricingStore
- time.Duration
- Module
- New
- CliToolsConsole.tsx
- ProviderAdapter
- Module
- newTestHarness
- App
- newTestHarness
- batcher
- responsesStreamEncoder
- net/url.URL
- Load
- scripts
- .RecordSuccess
- newTestHarness
- crypto.SessionSigner
- NewClient
- settingDefinition
- NewStaticUIHandler
- File Map
- stubModule
- syncBuffer
- OAuthConfig
- stubDepModule
- stubModule
- stubSystemModule
- stubDep
- .StatusCode
- stubModule
- log/slog.Logger
- nameMinHeap
- collectMigrations
- .postTokenEndpoint

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 96 edges
2. `WriteJSON()` - 81 edges
3. `extractErrorMessage()` - 75 edges
4. `DB` - 56 edges
5. `buildApp()` - 54 edges
6. `Module` - 41 edges
7. `Module` - 37 edges
8. `Button` - 37 edges
9. `Bus` - 36 edges
10. `ModelsRoutingConsole()` - 36 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`)` --references--> `useProvidersCatalog()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → web/src/features/providers/hooks.ts
- `Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity)` --references--> `UpstreamCredentials`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave2a-providers-proxypools-translator.md → internal/contracts/contracts.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (116 total, 23 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.06
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

### Community 5 - "APIKeySummary"
Cohesion: 0.08
Nodes (26): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), New(), decodeStringSlice(), encodeStringSlice(), APIKeySummary (+18 more)

### Community 6 - "package.json"
Cohesion: 0.05
Nodes (53): @biomejs/biome, clsx, geist, openapi-typescript, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form (+45 more)

### Community 7 - "contracts.ProxyDialerResolver"
Cohesion: 0.21
Nodes (9): File Map, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, defaultBuiltInAdapters(), TestWave2ProviderAdaptersRegisteredAndTestable() (+1 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.08
Nodes (24): Options, testHarness, loadUI(), loadUI(), New(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools() (+16 more)

### Community 10 - "Button"
Cohesion: 0.09
Nodes (60): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps (+52 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (58): TestAppCloseTeardownOrderIsLIFO(), TestSentinelErrorsDistinct(), TestModuleMetadataAndStartValidation(), setupGatewayTest(), TestGatewayAuthExtractionAndProtocolErrors(), TestGatewayCrossProtocolStreaming(), TestGatewayDirectFastPathIncrementalStreaming(), TestGatewayErrorClassificationAndOAuthSingleRetry() (+50 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "resolver.go"
Cohesion: 0.07
Nodes (27): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+19 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.15
Nodes (17): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+9 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (36): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+28 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.08
Nodes (48): @hookform/resolvers, DashboardSettingsPage(), authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, AuthPolicyResponse (+40 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "engine.go"
Cohesion: 0.18
Nodes (7): builtinRules(), boolToInt(), defaultConfig(), Config, GuardEvent, repository, Rule

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

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.05
Nodes (77): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation() (+69 more)

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
Cohesion: 0.12
Nodes (15): attemptOutcome, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), applyUpstreamHeaders(), buildUpstreamURL() (+7 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.22
Nodes (14): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+6 more)

### Community 42 - "RingBuffer"
Cohesion: 0.31
Nodes (6): Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, cloneAttrMap(), RingBuffer, LogEntry

### Community 43 - "TestContractInterfacesSatisfied"
Cohesion: 0.08
Nodes (13): fakeProviderStore, stubKeyValidator, stubModelLister, stubProviderStore, stubKeyValidator, stubModelLister, stubProviderStore, testEnv (+5 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "Module"
Cohesion: 0.08
Nodes (3): Module, stubSystemModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.14
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "context.Context"
Cohesion: 0.06
Nodes (6): AdminCredentialStore, stubAdminCredentialStore, actorFromContext(), validatePassword(), stubModule, orderStubModule

### Community 48 - "translate_openai.go"
Cohesion: 0.20
Nodes (18): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+10 more)

### Community 49 - "buildApp"
Cohesion: 0.10
Nodes (20): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+12 more)

### Community 50 - "Config"
Cohesion: 0.13
Nodes (18): Task 2: `auth` Module & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, BuildInfo, New(), newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion() (+10 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (60): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+52 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.17
Nodes (21): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+13 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.11
Nodes (23): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), extractAPIKey(), buildProtocolTestRequest(), isValidAuthType() (+15 more)

### Community 54 - "ProviderDetailView.tsx"
Cohesion: 0.09
Nodes (44): DashboardProviderDetailPage(), DashboardTranslatorPage(), useDeleteProviderConnectionMutation(), useProviderConnection(), useTestProviderConnectionMutation(), useUpdateProviderConnectionMutation(), formatTimestamp(), ProviderDetailView() (+36 more)

### Community 55 - "WriteJSON"
Cohesion: 0.17
Nodes (6): actorFromContext(), decodeJSONBody(), Module, Module, WriteJSON(), RouteMounts

### Community 56 - "Module"
Cohesion: 0.22
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 58 - "providers/repo.go"
Cohesion: 0.11
Nodes (19): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+11 more)

### Community 59 - "extractErrorMessage"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "modeleditor/repo.go"
Cohesion: 0.11
Nodes (22): copyResponseHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), matchPattern(), wrapSystemPrompt(), wrapText(), New() (+14 more)

### Community 61 - "client.ts"
Cohesion: 0.13
Nodes (15): openapi-fetch, openapi-fetch@^0.14.0, isHealthResponse(), apiClient, createApiClient(), ErrorDetail, ErrorEnvelope, HealthResponse (+7 more)

### Community 62 - "newTestHarness"
Cohesion: 0.24
Nodes (9): decodeJSON(), newTestHarness(), TestUsageAsyncBatcherAndGracefulDrain(), TestUsageBufferFullNeverBlocks(), TestUsageCostCalculationAndPricingCRUD(), TestUsageLiveSSEStream(), TestUsageSummaryLogsAndRetentionPrune(), waitForLogCount() (+1 more)

### Community 63 - "Module"
Cohesion: 0.18
Nodes (12): Event, actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID() (+4 more)

### Community 64 - "net/http.ResponseWriter"
Cohesion: 0.21
Nodes (5): generateInput, trackedResponseWriter, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.21
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (24): stubInterceptor, File Map, Review Focus, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests, Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification (+16 more)

### Community 67 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 68 - "translate_test.go"
Cohesion: 0.13
Nodes (22): CanonicalChatResponse, DecodeResponse(), EncodeResponse(), buildCompactedResponse(), randomResponseID(), estimateBlocks(), EstimateTokens(), encodeAnthropicResponse() (+14 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 71 - "DB"
Cohesion: 0.10
Nodes (22): main(), Global Constraints, Global Constraints, Global Constraints, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Global Constraints, Global Constraints, chainStore (+14 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (6): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - ".prepareUpstreamBody"
Cohesion: 0.27
Nodes (4): ProtocolTranslator, CanFastPath(), IsSupportedProtocol(), RewriteModelField()

### Community 74 - "time.Time"
Cohesion: 0.17
Nodes (12): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, Options, newCooldownTracker(), New() (+4 more)

### Community 75 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 76 - "File Map"
Cohesion: 0.28
Nodes (7): File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan

### Community 78 - "time.Duration"
Cohesion: 0.24
Nodes (4): failureRecord, stubHealthTracker, computeBackoff(), parseRetryAfter()

### Community 79 - "Module"
Cohesion: 0.26
Nodes (5): decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "New"
Cohesion: 0.15
Nodes (16): Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorOverridesAcrossAllFourProtocols(), doJSON(), setupPromptGuardTest() (+8 more)

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 83 - "Module"
Cohesion: 0.14
Nodes (8): Review Focus, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Module, Vault, Options, refreshCall, storedCredentials

### Community 84 - "newTestHarness"
Cohesion: 0.15
Nodes (14): New(), newTestHarness(), TestCatalogAndAdapterRegistry(), TestConnectionQuickTest(), TestConnectionsCRUDAndSecretEncryptionHygiene(), TestCustomProviderNodesCRUD(), TestHeadlessOAuthFlowAndCallbackXSSProtection(), TestOAuthEnvResolutionAndCallbackJSONAndValidation() (+6 more)

### Community 85 - "App"
Cohesion: 0.18
Nodes (3): Runner, job, App

### Community 86 - "newTestHarness"
Cohesion: 0.19
Nodes (12): stubAdminAuth, testHarness, decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering() (+4 more)

### Community 87 - "batcher"
Cohesion: 0.18
Nodes (4): stubUsageRecorder, stubUsageRecorder, UsageRecord, batcher

### Community 88 - "responsesStreamEncoder"
Cohesion: 0.42
Nodes (3): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder

### Community 89 - "net/url.URL"
Cohesion: 0.20
Nodes (11): stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), NewProviderProxyFunc() (+3 more)

### Community 90 - "Load"
Cohesion: 0.29
Nodes (8): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 91 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 93 - "newTestHarness"
Cohesion: 0.27
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 96 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 97 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 99 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 102 - "File Map"
Cohesion: 0.25
Nodes (7): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

### Community 105 - "OAuthConfig"
Cohesion: 0.27
Nodes (6): cloneOAuthConfig(), cloneStringSlice(), OAuthConfig, oauthEnvPrefix(), buildAuthorizeURL(), builtinAdapter

### Community 112 - "log/slog.Logger"
Cohesion: 0.40
Nodes (4): newBatcher(), New(), newPricingStore(), Options

### Community 114 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 430 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **23 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `File Map` to `Config`, `Button`, `WriteJSON`?**
  _High betweenness centrality (0.181) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.055135750901841654 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `File Map`?**
  _High betweenness centrality (0.121) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `Button` to `File Map`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._