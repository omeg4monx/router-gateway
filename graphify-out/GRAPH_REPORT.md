# Graph Report - wave2b-omeg4mon-router  (2026-11-23)

## Corpus Check
- 234 files · ~239,252 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: (none) 7, .lock 1, .css 1)

## Summary
- 2822 nodes · 9071 edges · 128 communities (103 shown, 25 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 860 edges (avg confidence: 0.83)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6c8ab9c8`
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
- .executeCandidateAttempt
- compilerOptions
- New
- Button
- testing.T
- components.json
- KeysConsole.tsx
- ResolvedTarget
- TokenSaverConsole.tsx
- catalog.go
- System Module & OpenAPI 3.1 Spec
- net/http/httptest.ResponseRecorder
- updates
- ir.go
- SettingsConsole.tsx
- with
- jobs
- contracts.ProxyDialerResolver
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
- TestContractInterfacesSatisfied
- translate_gemini.go
- context.Context
- repository
- newTestHarness
- translate_openai.go
- buildApp
- PromptGuardConsole.tsx
- ModelsRoutingConsole.tsx
- translate_responses.go
- net/http.ResponseWriter
- client.ts
- RouteMounts
- Module
- dependencies
- providers/repo.go
- extractErrorMessage
- modeleditor/repo.go
- Vault
- Server
- Module
- net/http.Request
- listSettingsWithQueryer
- recordingInterceptor
- promptinjection/engine.go
- ModelEditorConsole.tsx
- Module
- Module
- .Reader
- .AuthenticateAdmin
- scripts
- time.Time
- Module
- newTestHarness
- net/http/httptest.ResponseRecorder
- crypto.SessionSigner
- Module
- setupModelEditorTest
- CliToolsConsole.tsx
- ProviderAdapter
- .UpdateCombo
- DB
- repository
- pricingStore
- batcher
- NewStreamTranslator
- net/url.URL
- buildApp
- TranslatorConsole.tsx
- .ValidateKey
- newTestHarness
- app/layout.tsx
- stubModule
- stubSystemModule
- stubDepModule
- stubModule
- File Map
- time.Duration
- .postTokenEndpoint
- settingDefinition
- setupPromptGuardTest
- CanonicalStreamChunk
- setupTokenSaverTest
- .StatusCode
- NewStaticUIHandler
- repository
- stubModule
- stubModule
- Runner
- New
- Module
- Load
- stubDep
- runConcurrentScan
- headroomCache
- RingBuffer
- collectMigrations
- sync.Mutex
- dynamicProxyResolver
- @radix-ui/react-label

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 102 edges
2. `WriteError()` - 100 edges
3. `WriteJSON()` - 85 edges
4. `DB` - 59 edges
5. `buildApp()` - 53 edges
6. `Button` - 43 edges
7. `Module` - 42 edges
8. `Alert` - 39 edges
9. `Bus` - 38 edges
10. `AlertDescription` - 38 edges

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

## Communities (128 total, 25 thin omitted)

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

### Community 5 - "APIKeySummary"
Cohesion: 0.08
Nodes (25): cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), New(), decodeStringSlice(), encodeStringSlice(), APIKeySummary (+17 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - ".executeCandidateAttempt"
Cohesion: 0.09
Nodes (29): attemptOutcome, CanonicalChatResponse, ProtocolTranslator, applyUpstreamHeaders(), buildUpstreamURL(), extractUpstreamErrorMessage(), isProtectedRequestHeader(), parseRetryAfter() (+21 more)

### Community 8 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "New"
Cohesion: 0.13
Nodes (15): loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, trackedResponseWriter, apiTimeoutMiddleware(), New(), TestAdminAPIAuthEnforcement() (+7 more)

### Community 10 - "Button"
Cohesion: 0.09
Nodes (66): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps (+58 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (63): TestAppCloseTeardownOrderIsLIFO(), testEnv, TestSentinelErrorsDistinct(), setupGatewayTest(), TestGatewayAuthExtractionAndProtocolErrors(), TestGatewayCrossProtocolStreaming(), TestGatewayDirectFastPathIncrementalStreaming(), TestGatewayErrorClassificationAndOAuthSingleRetry() (+55 more)

### Community 12 - "components.json"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "KeysConsole.tsx"
Cohesion: 0.11
Nodes (30): @hookform/resolvers, DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation() (+22 more)

### Community 14 - "ResolvedTarget"
Cohesion: 0.13
Nodes (13): resolveSnapshot, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare() (+5 more)

### Community 15 - "TokenSaverConsole.tsx"
Cohesion: 0.16
Nodes (22): TokenSaverPage(), tokenSaverQueryKeys, useClearHeadroomCacheMutation(), usePreviewTokenSaverMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation() (+14 more)

### Community 16 - "catalog.go"
Cohesion: 0.20
Nodes (18): GenerateResult, ToolDescriptor, toolSpec, apiKeyOrPlaceholder(), generateToolConfig(), hasControlChar(), listToolDescriptors(), lookupToolSpec() (+10 more)

### Community 17 - "System Module & OpenAPI 3.1 Spec"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.20
Nodes (15): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+7 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.14
Nodes (35): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+27 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.06
Nodes (55): @tanstack/react-query, zod, @tanstack/react-query@^5.76.1, zod@^4.6.7, DashboardSettingsPage(), authKeys, AuthMeResponse, LoginRequest (+47 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "contracts.ProxyDialerResolver"
Cohesion: 0.38
Nodes (6): File Map, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification

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
Nodes (81): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderConnectionMutation() (+73 more)

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
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+8 more)

### Community 39 - "net/http.Request"
Cohesion: 0.15
Nodes (11): setSessionCookie(), decodeJSONBody(), anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), ensureRequestID() (+3 more)

### Community 40 - "UsageConsole.tsx"
Cohesion: 0.08
Nodes (47): DashboardUsagePage(), buildLogsQueryParams(), UpsertModelPricingVariables, UsageLogsFilter, usageQueryKeys, UsageStreamStatus, useDeleteModelPricingMutation(), useModelPricing() (+39 more)

### Community 41 - "translate_anthropic.go"
Cohesion: 0.23
Nodes (15): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+7 more)

### Community 42 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 43 - "TestContractInterfacesSatisfied"
Cohesion: 0.11
Nodes (9): fakeProviderStore, stubProviderStore, stubProviderStore, UpstreamCredentials, TestContractInterfacesSatisfied(), newFakeProviderStore(), cloneUpstreamCredentials(), sanitizeUpstreamError() (+1 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "context.Context"
Cohesion: 0.05
Nodes (7): AdminCredentialStore, stubAdminCredentialStore, actorFromContext(), validatePassword(), stubSystemModule, orderStubModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.13
Nodes (11): buildLogsWhere(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket (+3 more)

### Community 47 - "newTestHarness"
Cohesion: 0.15
Nodes (16): File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, newTestHarness() (+8 more)

### Community 48 - "translate_openai.go"
Cohesion: 0.23
Nodes (16): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+8 more)

### Community 49 - "buildApp"
Cohesion: 0.18
Nodes (3): normalizeName(), nameMinHeap, registryState

### Community 50 - "PromptGuardConsole.tsx"
Cohesion: 0.11
Nodes (34): PromptGuardPage(), promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules() (+26 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (60): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+52 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.16
Nodes (23): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+15 more)

### Community 53 - "net/http.ResponseWriter"
Cohesion: 0.15
Nodes (15): isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), truncateError(), validateHeaders(), validateName() (+7 more)

### Community 54 - "client.ts"
Cohesion: 0.16
Nodes (14): isHealthResponse(), useSystemHealth(), apiClient, createApiClient(), ErrorDetail, ErrorEnvelope, HealthResponse, InitCompleteResponse (+6 more)

### Community 55 - "RouteMounts"
Cohesion: 0.09
Nodes (7): Module, actorFromContext(), decodeJSONBody(), Module, RouteMounts, TestRegistryStartAllRollsBackOnPartialFailure(), stubModule

### Community 56 - "Module"
Cohesion: 0.26
Nodes (6): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), WriteError(), Module

### Community 57 - "dependencies"
Cohesion: 0.11
Nodes (19): clsx, openapi-fetch, @radix-ui/react-slot, react-dom, tailwind-merge, zustand, dependencies, class-variance-authority@^0.7.1 (+11 more)

### Community 58 - "providers/repo.go"
Cohesion: 0.11
Nodes (20): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), newRepository() (+12 more)

### Community 59 - "extractErrorMessage"
Cohesion: 0.11
Nodes (34): DashboardProxyPoolsPage(), useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation() (+26 more)

### Community 60 - "modeleditor/repo.go"
Cohesion: 0.09
Nodes (22): copyResponseHeaders(), applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), matchPattern(), validateModelPattern(), wrapSystemPrompt() (+14 more)

### Community 61 - "Vault"
Cohesion: 0.28
Nodes (5): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Global Constraints, Review Focus, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, Vault

### Community 62 - "Server"
Cohesion: 0.12
Nodes (3): Server, App, harness

### Community 63 - "Module"
Cohesion: 0.25
Nodes (9): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+1 more)

### Community 64 - "net/http.Request"
Cohesion: 0.14
Nodes (11): generateInput, clearSessionCookie(), isSecureRequest(), readSessionCookie(), decodeJSONBody(), Module, writeGenerateError(), extractAPIKey() (+3 more)

### Community 65 - "listSettingsWithQueryer"
Cohesion: 0.19
Nodes (7): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, rollbackResult, rowQueryer, UpdatedEvent

### Community 66 - "recordingInterceptor"
Cohesion: 0.07
Nodes (9): stubInterceptor, File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 4: Wire Wave 2B Modules in `cmd/omeg4mon-router`, OpenAPI Codegen & Live E2E Tests, Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, recordingInterceptor (+1 more)

### Community 67 - "promptinjection/engine.go"
Cohesion: 0.07
Nodes (31): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+23 more)

### Community 68 - "ModelEditorConsole.tsx"
Cohesion: 0.13
Nodes (29): ModelEditorPage(), modelEditorQueryKeys, useCreateBypassRuleMutation(), useCreateModelOverrideMutation(), useDeleteBypassRuleMutation(), useDeleteModelOverrideMutation(), useModelBypassRules(), useModelOverrides() (+21 more)

### Community 69 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 70 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), isValidRate(), Module

### Community 71 - ".Reader"
Cohesion: 0.08
Nodes (21): main(), runMain(), Global Constraints, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, File Map (+13 more)

### Community 72 - ".AuthenticateAdmin"
Cohesion: 0.09
Nodes (9): stubAdminAuth, stubAdminAuth, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth (+1 more)

### Community 73 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 74 - "time.Time"
Cohesion: 0.14
Nodes (11): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, stubAccountHealthTracker, newCooldownTracker(), newFakeClock() (+3 more)

### Community 76 - "newTestHarness"
Cohesion: 0.13
Nodes (16): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), defaultBuiltInAdapters(), New(), newTestHarness(), TestCatalogAndAdapterRegistry(), TestConnectionQuickTest(), TestConnectionsCRUDAndSecretEncryptionHygiene(), TestCustomProviderNodesCRUD() (+8 more)

### Community 77 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.18
Nodes (7): testHarness, decodeJSON(), newTestHarness(), TestCLIToolsCatalogAndGenerateAllTools(), TestCLIToolsLoopbackAtomicApplyAndRollback(), TestCLIToolsModuleContract(), TestCLIToolsNonLoopbackRejectsApplyAndRollback()

### Community 79 - "Module"
Cohesion: 0.20
Nodes (6): decodeJSONBody(), New(), newID(), validateBypassRule(), validateOverride(), Module

### Community 80 - "setupModelEditorTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags(), TestModelEditorCRUDValidationAndDryRunPreview(), TestModelEditorFinalReviewFindings(), TestModelEditorOverridesAcrossAllFourProtocols(), testEnv

### Community 81 - "CliToolsConsole.tsx"
Cohesion: 0.15
Nodes (26): DashboardCliToolsPage(), CliToolsConsole(), formatBadgeLabel(), CLIToolConfigMutationVariables, cliToolsQueryKeys, useApplyCliToolConfigMutation(), useCliToolsCatalog(), useGenerateCliToolConfigMutation() (+18 more)

### Community 82 - "ProviderAdapter"
Cohesion: 0.08
Nodes (10): cloneOAuthConfig(), cloneStringSlice(), OAuthConfig, oauthEnvPrefix(), Module, buildAuthorizeURL(), builtinAdapter, customTestAdapter (+2 more)

### Community 83 - ".UpdateCombo"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 84 - "DB"
Cohesion: 0.15
Nodes (23): Options, testHarness, Event, New(), New(), newRepository(), New(), newRepository() (+15 more)

### Community 85 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 86 - "pricingStore"
Cohesion: 0.23
Nodes (5): newBatcher(), New(), newPricingStore(), ModelPricing, pricingStore

### Community 87 - "batcher"
Cohesion: 0.20
Nodes (4): stubUsageRecorder, stubUsageRecorder, UsageRecord, batcher

### Community 88 - "NewStreamTranslator"
Cohesion: 0.15
Nodes (11): streamResult, StreamTranslator, NewStreamTranslator(), captureResponsesFrame(), mergeFastPathSSEUsage(), parseSSEFrame(), newAnthropicStreamEncoder(), newGeminiStreamEncoder() (+3 more)

### Community 89 - "net/url.URL"
Cohesion: 0.27
Nodes (7): stubProxyDialerResolver, stubProxyResolver, probeProxy(), proxyDialAddress(), NewProviderProxyFunc(), ProviderIDFromContext(), TestProviderProxyFuncContextRouting()

### Community 90 - "buildApp"
Cohesion: 0.09
Nodes (22): Options, buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageAndCLIToolsEndToEnd(), TestBuildAppCleansUpOnPartialFailure() (+14 more)

### Community 91 - "TranslatorConsole.tsx"
Cohesion: 0.14
Nodes (24): react-hook-form, DashboardTranslatorPage(), translatorQueryKeys, useConsoleLogsSnapshot(), useTranslatePayloadMutation(), useTranslatorProtocols(), ConsoleLogEntry, ConsoleLogsResponse (+16 more)

### Community 93 - "newTestHarness"
Cohesion: 0.18
Nodes (13): New(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestModuleMetadataAndStartValidation() (+5 more)

### Community 94 - "app/layout.tsx"
Cohesion: 0.24
Nodes (7): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 102 - "File Map"
Cohesion: 0.24
Nodes (8): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

### Community 104 - "time.Duration"
Cohesion: 0.27
Nodes (3): failureRecord, stubHealthTracker, computeBackoff()

### Community 106 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 107 - "setupPromptGuardTest"
Cohesion: 0.31
Nodes (7): Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), doJSON(), setupPromptGuardTest(), TestPromptInjectionBlockSanitizeAuditAndBuiltins(), TestPromptInjectionCustomRuleValidationAndRE2(), TestPromptInjectionFinalReviewFindings(), testEnv

### Community 108 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, mapOpenAIFinishReasonToCanonical()

### Community 109 - "setupTokenSaverTest"
Cohesion: 0.31
Nodes (7): doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFinalReviewFindings(), testEnv

### Community 111 - "NewStaticUIHandler"
Cohesion: 0.39
Nodes (6): hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 112 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 116 - "New"
Cohesion: 0.10
Nodes (22): Options, newChainStore(), New(), New(), newHarness(), TestTranslatorLogsSnapshotAndLiveSSEStream(), TestTranslatorLogsStreamSurvivesPastAPITimeout(), TestTranslatorTranslateRequestAndResponseAllPairs() (+14 more)

### Community 117 - "Module"
Cohesion: 0.22
Nodes (3): decodeJSONBody(), New(), Module

### Community 118 - "Load"
Cohesion: 0.33
Nodes (7): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 120 - "runConcurrentScan"
Cohesion: 0.40
Nodes (6): parseAndValidateProxyURL(), redactRawCandidate(), runConcurrentScan(), parsedProxyURL, ScanCandidateResult, ScanResponse

### Community 121 - "headroomCache"
Cohesion: 0.24
Nodes (3): newHeadroomCache(), headroomCache, WarmPrefixSummary

### Community 122 - "RingBuffer"
Cohesion: 0.60
Nodes (3): cloneAttrMap(), RingBuffer, LogEntry

### Community 123 - "collectMigrations"
Cohesion: 0.50
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 454 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **25 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `newTestHarness` to `New`, `Button`, `Module`, `.Reader`?**
  _High betweenness centrality (0.226) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0544977484664752 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Button` to `newTestHarness`?**
  _High betweenness centrality (0.157) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04521276595744681 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `Button` to `SettingsConsole.tsx`, `newTestHarness`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._