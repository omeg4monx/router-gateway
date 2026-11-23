# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 157 files · ~166,379 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 2014 nodes · 6142 edges · 80 communities (67 shown, 13 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 622 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7689a491`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- settings/hooks.ts
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- package.json
- Structured Logging & Ring Buffer
- TypeScript Strict Compiler Config
- Server
- AppShell.tsx
- testing.T
- shadcn/ui Component Configuration
- stubModule
- repository
- io/fs.FS
- .postTokenEndpoint
- Design System & Wave 0 Plan Docs
- Bus
- updates
- ir.go
- InitForm.tsx
- with
- jobs
- resolver.go
- steps
- web-check
- on
- cross-compile
- Module
- ProvidersConsole.tsx
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- JSON Marshaler Error Test Fixture
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- net/http.ResponseWriter
- SettingsConsole.tsx
- encoding/json.RawMessage
- buildApp
- UpstreamCredentials
- translate_gemini.go
- fakeCredentialStore
- context.Context
- translate_openai.go
- App
- Config
- ModelsRoutingConsole.tsx
- translate_responses.go
- WriteError
- Module
- net/http.Request
- settingDefinition
- CanonicalStreamChunk
- repository
- DB
- Bus
- stubModule
- Module
- NewStreamTranslator
- WriteJSON
- recordingInterceptor
- GatewayInterceptor
- scripts
- Options
- .RequireLogin
- time.Time
- listSettingsWithQueryer
- sync.Mutex
- time.Duration
- app/layout.tsx
- .RecordSuccess
- orderStubModule
- nameMinHeap
- AppShell.tsx
- crypto.SessionSigner
- net/http.ResponseWriter

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 57 edges
2. `extractErrorMessage()` - 51 edges
3. `WriteJSON()` - 50 edges
4. `buildApp()` - 41 edges
5. `DB` - 40 edges
6. `ModelsRoutingConsole()` - 36 edges
7. `Module` - 34 edges
8. `Module` - 34 edges
9. `ProvidersConsole()` - 34 edges
10. `Module` - 30 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → cmd/omeg4mon-router/main.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan` --references--> `Bus`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/eventbus/bus.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (80 total, 13 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.05
Nodes (78): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+70 more)

### Community 1 - "settings/hooks.ts"
Cohesion: 0.07
Nodes (36): @tanstack/react-query, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables (+28 more)

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
Cohesion: 0.06
Nodes (31): stubKeyValidator, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, stubKeyValidator, KeyInfo (+23 more)

### Community 6 - "package.json"
Cohesion: 0.06
Nodes (49): @biomejs/biome, class-variance-authority, clsx, geist, openapi-fetch, openapi-typescript, @radix-ui/react-label, @radix-ui/react-slot (+41 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.19
Nodes (5): addAttrAt(), attrValue(), ensurePath(), boundAttrs, teeHandler

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.15
Nodes (14): loadUI(), loadUI(), Options, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement(), TestDevUIHandlerFallback() (+6 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.12
Nodes (45): lucide-react, next, react, nextConfig, DashboardPage(), InitPage(), LoginPage(), HomePage() (+37 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (68): TestAppCloseTeardownOrderIsLIFO(), testEnv, TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine() (+60 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 14 - "repository"
Cohesion: 0.15
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.13
Nodes (10): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile() (+2 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.18
Nodes (17): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+9 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.10
Nodes (46): CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock, ToolChoice (+38 more)

### Community 21 - "InitForm.tsx"
Cohesion: 0.09
Nodes (44): @hookform/resolvers, zod, DashboardSettingsPage(), authKeys, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+36 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.10
Nodes (18): fakeProviderStore, resolveSnapshot, targetKey, stubModelLister, stubProviderResolver, stubModelLister, stubResolver, ModelDescriptor (+10 more)

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
Nodes (82): DashboardProvidersPage(), ExchangeProviderOAuthVariables, providersQueryKeys, StartProviderOAuthVariables, UpdateProviderConnectionVariables, useCreateProviderConnectionMutation(), useCreateProviderNodeMutation(), useDeleteProviderNodeMutation() (+74 more)

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 38 - "CanonicalUsage"
Cohesion: 0.10
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.10
Nodes (21): attemptOutcome, errorDisposition, statusError, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError() (+13 more)

### Community 40 - "SettingsConsole.tsx"
Cohesion: 0.13
Nodes (32): DashboardKeysPage(), DashboardProviderDetailPage(), useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation(), useUpdateApiKeyMutation(), formatQuotaUsd() (+24 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.13
Nodes (26): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+18 more)

### Community 42 - "buildApp"
Cohesion: 0.16
Nodes (12): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), TestVaultRoundTripAndTamperRejection() (+4 more)

### Community 43 - "UpstreamCredentials"
Cohesion: 0.21
Nodes (6): stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall

### Community 44 - "translate_gemini.go"
Cohesion: 0.17
Nodes (20): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+12 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.05
Nodes (7): Module, RouteMounts, TestRegistryStartAllRollsBackOnPartialFailure(), stubSystemModule, stubModule, stubSystemModule, stubSystemModule

### Community 47 - "context.Context"
Cohesion: 0.09
Nodes (5): ProxyDialerResolver, stubAdminCredentialStore, stubProxyDialerResolver, actorFromContext(), TestSystemAdminCredentialStore()

### Community 48 - "translate_openai.go"
Cohesion: 0.20
Nodes (18): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatTool (+10 more)

### Community 49 - "App"
Cohesion: 0.12
Nodes (4): Runner, job, syncBuffer, App

### Community 50 - "Config"
Cohesion: 0.11
Nodes (18): Task 2: `auth` Module & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, BuildInfo, New(), newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle() (+10 more)

### Community 51 - "ModelsRoutingConsole.tsx"
Cohesion: 0.07
Nodes (58): DashboardModelsPage(), combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useCombos(), useCreateComboMutation(), useDeleteComboMutation() (+50 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.17
Nodes (21): openAIChatFunctionSpec, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest (+13 more)

### Community 53 - "WriteError"
Cohesion: 0.16
Nodes (16): buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), Module, truncateError() (+8 more)

### Community 54 - "Module"
Cohesion: 0.06
Nodes (17): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, cloneOAuthConfig(), cloneStringSlice(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), Module, New() (+9 more)

### Community 56 - "settingDefinition"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 58 - "repository"
Cohesion: 0.11
Nodes (20): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), scanConnectionRecord() (+12 more)

### Community 60 - "DB"
Cohesion: 0.14
Nodes (10): main(), runMain(), Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan (+2 more)

### Community 61 - "Bus"
Cohesion: 0.09
Nodes (21): Options, stubAdminAuth, testHarness, Event, New(), TestModuleMetadataAndStartValidation(), newRepository(), New() (+13 more)

### Community 63 - "Module"
Cohesion: 0.28
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 64 - "NewStreamTranslator"
Cohesion: 0.29
Nodes (7): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs()

### Community 65 - "WriteJSON"
Cohesion: 0.22
Nodes (5): actorFromContext(), decodeJSONBody(), Module, Module, WriteJSON()

### Community 66 - "recordingInterceptor"
Cohesion: 0.15
Nodes (8): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), copyResponseHeaders(), applyStaticCacheHeaders()

### Community 67 - "GatewayInterceptor"
Cohesion: 0.29
Nodes (8): GatewayInterceptor, File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification)

### Community 68 - "scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

### Community 71 - "Options"
Cohesion: 0.10
Nodes (18): File Map, Review Focus, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+10 more)

### Community 72 - ".RequireLogin"
Cohesion: 0.13
Nodes (4): stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth

### Community 74 - "time.Time"
Cohesion: 0.17
Nodes (8): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, chainStore, newCooldownTracker(), newFakeClock()

### Community 75 - "listSettingsWithQueryer"
Cohesion: 0.19
Nodes (7): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, rollbackResult, rowQueryer, UpdatedEvent

### Community 77 - "sync.Mutex"
Cohesion: 0.39
Nodes (3): stubUsageRecorder, stubUsageRecorder, UsageRecord

### Community 78 - "time.Duration"
Cohesion: 0.27
Nodes (3): failureRecord, stubHealthTracker, computeBackoff()

### Community 82 - "app/layout.tsx"
Cohesion: 0.32
Nodes (5): metadata, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 90 - "AppShell.tsx"
Cohesion: 0.14
Nodes (17): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), DashboardLayout(), DashboardLayoutProps, useLogoutMutation(), applyThemeToDocument(), ThemeMode, UIState, useUIStore (+9 more)

### Community 92 - "net/http.ResponseWriter"
Cohesion: 0.16
Nodes (8): trackedResponseWriter, clearSessionCookie(), isSecureRequest(), setSessionCookie(), decodeJSONBody(), AdminSubject(), jsonMethodNotAllowedHandler(), jsonNotFoundHandler()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 355 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **13 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `APIKeySummary` to `WriteJSON`, `Config`, `DB`, `AppShell.tsx`?**
  _High betweenness centrality (0.297) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.05258616313892696 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `AppShell.tsx` to `AppShell.tsx`, `APIKeySummary`?**
  _High betweenness centrality (0.206) - this node is a cross-community bridge._
- **Should `settings/hooks.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.06767676767676768 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `SettingsConsole.tsx`, `AppShell.tsx`, `APIKeySummary`?**
  _High betweenness centrality (0.172) - this node is a cross-community bridge._