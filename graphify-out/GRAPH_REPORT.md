# Graph Report - wave2a-omeg4mon-router  (2026-11-23)

## Corpus Check
- 199 files · ~209,499 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 2471 nodes · 7772 edges · 104 communities (90 shown, 14 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 748 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `c158a8bb`
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
- contracts.ProxyDialerResolver
- compilerOptions
- Server
- Button
- testing.T
- components.json
- KeysConsole.tsx
- repository
- io/fs.FS
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
- ProvidersConsole.tsx
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
- RingBuffer
- Module
- translate_gemini.go
- Module
- repository
- context.Context
- translate_openai.go
- buildApp
- Config
- ModelsRoutingConsole.tsx
- translate_responses.go
- .handleCreateConnection
- ProviderDetailView.tsx
- Module
- stubModule
- CanonicalStreamChunk
- repository
- extractErrorMessage
- ProviderDetailView.tsx
- client.ts
- dependencies
- Module
- Module
- listSettingsWithQueryer
- TestContractInterfacesSatisfied
- teeHandler
- .executeCandidateAttempt
- Module
- resolveStaticFile
- DB
- .AuthenticateAdmin
- ProtocolTranslator
- time.Time
- repository
- File Map
- pricingStore
- time.Duration
- WriteError
- newTestHarness
- CliToolsConsole.tsx
- ProviderAdapter
- net/http.Request
- repository
- App
- newTestHarness
- batcher
- NewStreamTranslator
- net/url.URL
- newTestHarness
- stubModule
- .RecordSuccess
- newTestHarness
- crypto.SessionSigner
- NewClient
- settingDefinition
- Module
- io/fs.FS
- scanner.go
- .postTokenEndpoint

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 76 edges
2. `extractErrorMessage()` - 75 edges
3. `WriteJSON()` - 67 edges
4. `buildApp()` - 53 edges
5. `DB` - 50 edges
6. `Module` - 39 edges
7. `Button` - 37 edges
8. `Module` - 36 edges
9. `ModelsRoutingConsole()` - 36 edges
10. `Module` - 34 edges

## Surprising Connections (you probably didn't know these)
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
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

## Communities (104 total, 14 thin omitted)

### Community 0 - "main.go"
Cohesion: 0.05
Nodes (86): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+78 more)

### Community 1 - "Module"
Cohesion: 0.24
Nodes (5): decodeStrictJSON(), Module, newPoolID(), formatUTC(), parseAndValidateProxyURL()

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
Cohesion: 0.09
Nodes (23): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, scanSummary() (+15 more)

### Community 6 - "package.json"
Cohesion: 0.05
Nodes (42): @biomejs/biome, geist, openapi-typescript, react-dom, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node (+34 more)

### Community 7 - "contracts.ProxyDialerResolver"
Cohesion: 0.19
Nodes (10): File Map, Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, defaultBuiltInAdapters() (+2 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.15
Nodes (14): Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, apiTimeoutMiddleware(), Server, New(), TestDevUIHandlerFallback(), testLogger(), TestPanicRecoveryReturnsJSON500() (+6 more)

### Community 10 - "Button"
Cohesion: 0.10
Nodes (52): class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, InitPage() (+44 more)

### Community 11 - "testing.T"
Cohesion: 0.05
Nodes (72): TestAppCloseTeardownOrderIsLIFO(), stubKeyValidator, testEnv, TestSentinelErrorsDistinct(), setupGatewayTest(), TestGatewayAuthExtractionAndProtocolErrors(), TestGatewayCrossProtocolStreaming(), TestGatewayDirectFastPathIncrementalStreaming() (+64 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 14 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.06
Nodes (5): stubModule, stubModule, RouteMounts, stubSystemModule, stubDep

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.13
Nodes (18): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+10 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (33): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+25 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (43): DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation(), usePatchSettingsMutation() (+35 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.19
Nodes (12): resolveSnapshot, targetKey, stubProviderResolver, stubResolver, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit() (+4 more)

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
Cohesion: 0.17
Nodes (9): actorFromContext(), decodeJSONBody(), decodePathModelParam(), generateID(), Module, validateComboTargets(), validateDescription(), validateModelIdentifier() (+1 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.05
Nodes (79): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation() (+71 more)

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
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, streamResult (+8 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.13
Nodes (11): Module, statusError, trackedResponseWriter, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError() (+3 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (48): @hookform/resolvers, DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation() (+40 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.23
Nodes (14): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+6 more)

### Community 42 - "RingBuffer"
Cohesion: 0.31
Nodes (6): Global Constraints, Review Focus, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, cloneAttrMap(), RingBuffer, LogEntry

### Community 43 - "Module"
Cohesion: 0.15
Nodes (7): ProviderConnection, stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "Module"
Cohesion: 0.08
Nodes (3): Module, stubSystemModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.15
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "context.Context"
Cohesion: 0.07
Nodes (8): AdminCredentialStore, KeyInfo, stubAdminCredentialStore, stubKeyValidator, actorFromContext(), validatePassword(), stubModule, orderStubModule

### Community 48 - "translate_openai.go"
Cohesion: 0.16
Nodes (22): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+14 more)

### Community 49 - "buildApp"
Cohesion: 0.11
Nodes (14): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+6 more)

### Community 50 - "Config"
Cohesion: 0.14
Nodes (17): Task 2: `auth` Module & OpenAPI Fragment, BuildInfo, New(), newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle() (+9 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.18
Nodes (19): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+11 more)

### Community 53 - ".handleCreateConnection"
Cohesion: 0.20
Nodes (13): isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), truncateError(), validateHeaders(), validateName() (+5 more)

### Community 54 - "ProviderDetailView.tsx"
Cohesion: 0.14
Nodes (28): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), DashboardPage(), DashboardProviderDetailPage(), useDeleteProviderConnectionMutation(), useProviderConnection(), useTestProviderConnectionMutation(), useUpdateProviderConnectionMutation(), formatTimestamp() (+20 more)

### Community 55 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 58 - "repository"
Cohesion: 0.11
Nodes (20): cloneStringMap(), cloneStringSlice(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+12 more)

### Community 59 - "extractErrorMessage"
Cohesion: 0.11
Nodes (35): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+27 more)

### Community 60 - "ProviderDetailView.tsx"
Cohesion: 0.15
Nodes (22): DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse, parseJsonObject() (+14 more)

### Community 61 - "client.ts"
Cohesion: 0.15
Nodes (16): openapi-fetch, openapi-fetch@^0.14.0, isHealthResponse(), useSystemHealth(), apiClient, createApiClient(), ErrorDetail, ErrorEnvelope (+8 more)

### Community 62 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-hook-form, tailwind-merge, @tanstack/react-query, zod, zustand (+18 more)

### Community 63 - "Module"
Cohesion: 0.28
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 64 - "Module"
Cohesion: 0.22
Nodes (5): generateInput, testHarness, decodeJSONBody(), Module, writeGenerateError()

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.19
Nodes (7): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "TestContractInterfacesSatisfied"
Cohesion: 0.09
Nodes (13): CanonicalRequest, ModelDescriptor, ResponseMeta, StreamEvent, stubInterceptor, stubModelLister, recordingInterceptor, stubModelLister (+5 more)

### Community 67 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 68 - ".executeCandidateAttempt"
Cohesion: 0.11
Nodes (30): attemptOutcome, CanonicalChatResponse, applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage(), isProtectedRequestHeader(), parseRetryAfter(), CanFastPath() (+22 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "resolveStaticFile"
Cohesion: 0.53
Nodes (4): isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 71 - "DB"
Cohesion: 0.08
Nodes (28): Options, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Event, New(), TestModuleMetadataAndStartValidation(), newRepository(), New(), newRepository() (+20 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (7): stubAdminAuth, stubAuthVerifier, stubAuthVerifier, readSessionCookie(), stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "ProtocolTranslator"
Cohesion: 0.29
Nodes (3): ProtocolTranslator, New(), Options

### Community 74 - "time.Time"
Cohesion: 0.13
Nodes (10): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, newCooldownTracker(), newFakeClock(), newFakeClock() (+2 more)

### Community 75 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 76 - "File Map"
Cohesion: 0.09
Nodes (21): main(), runMain(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment (+13 more)

### Community 78 - "time.Duration"
Cohesion: 0.21
Nodes (5): chainStore, failureRecord, stubHealthTracker, computeBackoff(), newChainStore()

### Community 79 - "WriteError"
Cohesion: 0.15
Nodes (13): setSessionCookie(), decodeJSONBody(), actorFromContext(), decodeJSONBody(), Module, WriteError(), WriteJSON(), AdminSubject() (+5 more)

### Community 80 - "newTestHarness"
Cohesion: 0.32
Nodes (8): Options, New(), decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback()

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "ProviderAdapter"
Cohesion: 0.10
Nodes (7): cloneOAuthConfig(), OAuthConfig, oauthEnvPrefix(), buildAuthorizeURL(), builtinAdapter, customTestAdapter, ProviderAdapter

### Community 83 - "net/http.Request"
Cohesion: 0.18
Nodes (6): clearSessionCookie(), isSecureRequest(), extractAPIKey(), buildProtocolTestRequest(), decodeOptionalJSON(), Module

### Community 84 - "repository"
Cohesion: 0.14
Nodes (10): File Map, Global Constraints, Review Focus, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification (+2 more)

### Community 85 - "App"
Cohesion: 0.14
Nodes (3): Runner, job, App

### Community 86 - "newTestHarness"
Cohesion: 0.16
Nodes (14): fakeProviderStore, stubAdminAuth, testHarness, decodeJSON(), newFakeProviderStore(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency() (+6 more)

### Community 87 - "batcher"
Cohesion: 0.18
Nodes (5): stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), batcher

### Community 88 - "NewStreamTranslator"
Cohesion: 0.18
Nodes (10): geminiPendingTool, geminiStreamEncoder, NewStreamTranslator(), parseJSONObjectOrWrap(), newAnthropicStreamEncoder(), mapCanonicalStopReasonToGemini(), newGeminiStreamEncoder(), newOpenAIStreamDecoder() (+2 more)

### Community 89 - "net/url.URL"
Cohesion: 0.29
Nodes (6): stubProxyDialerResolver, stubProxyResolver, NewProviderProxyFunc(), ProviderIDFromContext(), TestProviderProxyFuncContextRouting(), parsedProxyURL

### Community 90 - "newTestHarness"
Cohesion: 0.22
Nodes (8): assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), stubAdminAuth, testHarness

### Community 93 - "newTestHarness"
Cohesion: 0.27
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 96 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 97 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 100 - "io/fs.FS"
Cohesion: 0.10
Nodes (7): loadUI(), loadUI(), migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), stubDepModule

### Community 102 - "scanner.go"
Cohesion: 0.43
Nodes (7): probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), PoolTestResponse, ScanCandidateResult, ScanResponse

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 419 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **14 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `File Map` to `Server`, `Config`, `net/http.Request`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.183) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.05377126004928696 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `contracts.ProxyDialerResolver` to `RingBuffer`?**
  _High betweenness centrality (0.133) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `useProvidersCatalog()` connect `extractErrorMessage` to `ProvidersConsole.tsx`, `contracts.ProxyDialerResolver`?**
  _High betweenness centrality (0.128) - this node is a cross-community bridge._