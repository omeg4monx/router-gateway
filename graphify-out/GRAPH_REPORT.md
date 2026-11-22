# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 151 files · ~157,467 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 1934 nodes · 5908 edges · 82 communities (70 shown, 12 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 605 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f74db31a`
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
- time.Time
- repository
- io/fs.FS
- KeysConsole.tsx
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
- Module
- translate_gemini.go
- fakeCredentialStore
- repository
- context.Context
- translate_openai.go
- App
- Config
- dependencies
- translate_responses.go
- WriteError
- Module
- net/http.Request
- listSettingsWithQueryer
- CanonicalStreamChunk
- nameMinHeap
- tokens.go
- DB
- auth/hooks.ts
- RouteMounts
- Module
- NewStreamTranslator
- Module
- TestContractInterfacesSatisfied
- New
- resolver.go
- Module
- newTestHarness
- NewClient
- .RequireLogin
- parseJSONObjectOrWrap
- sync.Mutex
- time.Duration
- .RecordSuccess
- UsageRecord
- stubModule
- Load
- crypto.SessionSigner
- trackedResponseWriter

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 57 edges
2. `WriteJSON()` - 50 edges
3. `DB` - 40 edges
4. `buildApp()` - 39 edges
5. `extractErrorMessage()` - 38 edges
6. `ProvidersConsole()` - 34 edges
7. `Module` - 34 edges
8. `Module` - 34 edges
9. `Module` - 30 edges
10. `SettingsConsole()` - 30 edges

## Surprising Connections (you probably didn't know these)
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain` --references--> `UsageRecord`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → internal/contracts/contracts.go
- `Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment` --references--> `KeyInfo`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter` --references--> `repository`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/modules/system/repo.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (82 total, 12 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.06
Nodes (66): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+58 more)

### Community 1 - "settings/hooks.ts"
Cohesion: 0.09
Nodes (28): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), openapi-fetch, openapi-fetch@^0.14.0, DashboardPage(), isHealthResponse(), systemKeys, useInitStatus(), useSystemHealth() (+20 more)

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
Cohesion: 0.08
Nodes (25): stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice(), encodeStringSlice() (+17 more)

### Community 6 - "package.json"
Cohesion: 0.07
Nodes (35): @biomejs/biome, openapi-typescript, react-dom, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react (+27 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), RingBuffer, parseLevel(), boundAttrs, LogEntry (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.15
Nodes (14): loadUI(), loadUI(), Options, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement(), TestDevUIHandlerFallback() (+6 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.11
Nodes (47): lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, InitPage(), LoginPage() (+39 more)

### Community 11 - "testing.T"
Cohesion: 0.07
Nodes (57): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+49 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "time.Time"
Cohesion: 0.16
Nodes (8): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, chainStore, newCooldownTracker(), newFakeClock()

### Community 14 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.13
Nodes (10): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile() (+2 more)

### Community 16 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (30): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+22 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.16
Nodes (17): authHarness, fakeCredentialStore, Options, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+9 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.12
Nodes (44): CacheControl, CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+36 more)

### Community 21 - "InitForm.tsx"
Cohesion: 0.09
Nodes (45): @hookform/resolvers, zod, @hookform/resolvers@^7.0.8, zod@^4.6.7, DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse (+37 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.11
Nodes (13): ProxyDialerResolver, stubModelLister, stubProviderResolver, stubProxyDialerResolver, stubKeyValidator, stubModelLister, stubResolver, testEnv (+5 more)

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
Cohesion: 0.12
Nodes (17): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamDecoder, streamEncoder (+9 more)

### Community 39 - "net/http.ResponseWriter"
Cohesion: 0.13
Nodes (16): attemptOutcome, statusError, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError(), writeResolutionError(), applyUpstreamHeaders() (+8 more)

### Community 40 - "SettingsConsole.tsx"
Cohesion: 0.18
Nodes (19): DashboardProviderDetailPage(), useDeleteProviderConnectionMutation(), useProviderConnection(), useTestProviderConnectionMutation(), useUpdateProviderConnectionMutation(), formatTimestamp(), ProviderDetailView(), renderTestBadge() (+11 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.20
Nodes (17): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+9 more)

### Community 42 - "buildApp"
Cohesion: 0.16
Nodes (12): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1BGatewayEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), TestVaultRoundTripAndTamperRejection() (+4 more)

### Community 43 - "Module"
Cohesion: 0.14
Nodes (7): stubProviderStore, stubProviderStore, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall, storedCredentials

### Community 44 - "translate_gemini.go"
Cohesion: 0.15
Nodes (22): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+14 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.05
Nodes (6): Module, validatePassword(), stubSystemModule, stubModule, stubSystemModule, stubSystemModule

### Community 46 - "repository"
Cohesion: 0.11
Nodes (19): cloneStringMap(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError(), scanConnectionRecord() (+11 more)

### Community 47 - "context.Context"
Cohesion: 0.10
Nodes (4): stubAdminCredentialStore, actorFromContext(), oauthTokenResponse, repository

### Community 48 - "translate_openai.go"
Cohesion: 0.18
Nodes (20): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+12 more)

### Community 49 - "App"
Cohesion: 0.12
Nodes (8): TestAppCloseTeardownOrderIsLIFO(), Runner, NewRunner(), TestRunnerIdempotentStartStopAndGuards(), TestRunnerRecoversFromPanicAndStopsCleanly(), TestRunnerStopRespectsContextDeadline(), job, App

### Community 50 - "Config"
Cohesion: 0.13
Nodes (18): Task 2: `auth` Module & OpenAPI Fragment, Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, BuildInfo, New(), newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace() (+10 more)

### Community 51 - "dependencies"
Cohesion: 0.08
Nodes (24): class-variance-authority, clsx, geist, @radix-ui/react-label, @radix-ui/react-slot, react-hook-form, tailwind-merge, zustand (+16 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.18
Nodes (19): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+11 more)

### Community 53 - "WriteError"
Cohesion: 0.17
Nodes (17): buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), decodeJSONBody(), decodeOptionalJSONBody(), Module, truncateError() (+9 more)

### Community 54 - "Module"
Cohesion: 0.08
Nodes (14): cloneOAuthConfig(), cloneStringSlice(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), Module, New(), TestOAuthEnvResolutionAndCallbackJSONAndValidation() (+6 more)

### Community 55 - "net/http.Request"
Cohesion: 0.12
Nodes (10): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), extractAPIKey(), ensureRequestID(), AdminSubject() (+2 more)

### Community 56 - "listSettingsWithQueryer"
Cohesion: 0.10
Nodes (15): listSettingsWithQueryer(), newRepository(), decodeStoredValue(), int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds(), applyResult (+7 more)

### Community 59 - "tokens.go"
Cohesion: 0.67
Nodes (3): estimateBlocks(), EstimateTokens(), TestEstimateTokens()

### Community 60 - "DB"
Cohesion: 0.08
Nodes (25): main(), runMain(), Options, stubAdminAuth, testHarness, File Map, Global Constraints, Review Focus (+17 more)

### Community 61 - "auth/hooks.ts"
Cohesion: 0.12
Nodes (13): @tanstack/react-query, @tanstack/react-query@^5.76.1, metadata, RootLayout(), RootLayoutProps, authKeys, AuthMeResponse, LoginRequest (+5 more)

### Community 62 - "RouteMounts"
Cohesion: 0.11
Nodes (3): stubModule, RouteMounts, orderStubModule

### Community 63 - "Module"
Cohesion: 0.28
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 64 - "NewStreamTranslator"
Cohesion: 0.29
Nodes (7): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestStreamTranslators()

### Community 66 - "TestContractInterfacesSatisfied"
Cohesion: 0.07
Nodes (21): GatewayInterceptor, stubInterceptor, File Map, Review Focus, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain (+13 more)

### Community 67 - "New"
Cohesion: 0.29
Nodes (6): New(), TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(), TestLevelFilteringAndChildLoggerAttrsAndGroups(), TestRingBufferCapsAndPreservesOrder(), TestRingBufferConcurrentWritesAndReads(), TestRingBufferRecentDeepClonesAttrs()

### Community 68 - "resolver.go"
Cohesion: 0.22
Nodes (10): fakeProviderStore, resolveSnapshot, targetKey, ProviderConnection, newFakeProviderStore(), buildResolvedTarget(), connectionSupportsBare(), connectionSupportsExplicit() (+2 more)

### Community 69 - "Module"
Cohesion: 0.39
Nodes (3): actorFromContext(), decodeJSONBody(), Module

### Community 70 - "newTestHarness"
Cohesion: 0.19
Nodes (11): New(), indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid() (+3 more)

### Community 71 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 72 - ".RequireLogin"
Cohesion: 0.17
Nodes (5): stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth

### Community 73 - "parseJSONObjectOrWrap"
Cohesion: 0.36
Nodes (4): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini()

### Community 74 - "sync.Mutex"
Cohesion: 0.25
Nodes (3): failureRecord, stubHealthTracker, syncBuffer

### Community 77 - "UsageRecord"
Cohesion: 0.48
Nodes (3): stubUsageRecorder, stubUsageRecorder, UsageRecord

### Community 79 - "Load"
Cohesion: 0.33
Nodes (7): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 339 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **12 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `DB` to `Binary Entrypoint & Embedded UI`, `settings/hooks.ts`, `Config`, `Module`?**
  _High betweenness centrality (0.275) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.059753231492361926 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `settings/hooks.ts` to `AppShell.tsx`, `DB`?**
  _High betweenness centrality (0.196) - this node is a cross-community bridge._
- **Should `settings/hooks.ts` be split into smaller, more focused modules?**
  _Cohesion score 0.09243697478991597 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `settings/hooks.ts`, `SettingsConsole.tsx`, `dependencies`, `DB`, `auth/hooks.ts`?**
  _High betweenness centrality (0.155) - this node is a cross-community bridge._