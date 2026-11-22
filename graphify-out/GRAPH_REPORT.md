# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 137 files · ~132,349 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 1726 nodes · 5083 edges · 66 communities (58 shown, 8 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 532 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `285499b7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- SystemOverview.tsx
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
- SettingsConsole.tsx
- with
- jobs
- resolver.go
- steps
- web-check
- on
- cross-compile
- Module
- Web Package Build Scripts
- Generated OpenAPI TypeScript Schema
- docker-build
- go-check
- JSON Marshaler Error Test Fixture
- Tailwind PostCSS Config
- Go Module Root Definition
- CanonicalUsage
- SettingItem
- SettingsConsole.tsx
- encoding/json.RawMessage
- buildApp
- UpstreamCredentials
- translate_gemini.go
- fakeCredentialStore
- providers/repo.go
- context.Context
- translate_openai.go
- App
- TestContractInterfacesSatisfied
- AppShell.tsx
- translate_responses.go
- New
- Config
- DB
- listSettingsWithQueryer
- CanonicalStreamChunk
- nameMinHeap
- parseJSONObjectOrWrap
- newTestHarness
- NewClient
- RouteMounts
- .postTokenEndpoint
- resolveStaticFile
- .ProxyURLForProvider

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 57 edges
2. `WriteJSON()` - 45 edges
3. `DB` - 35 edges
4. `buildApp()` - 35 edges
5. `Module` - 34 edges
6. `Module` - 30 edges
7. `SettingsConsole()` - 30 edges
8. `Module` - 27 edges
9. `Module` - 26 edges
10. `KeysConsole()` - 25 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `File Map` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (66 total, 8 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.05
Nodes (72): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+64 more)

### Community 1 - "SystemOverview.tsx"
Cohesion: 0.07
Nodes (35): openapi-fetch, @tanstack/react-query, openapi-fetch@^0.14.0, @tanstack/react-query@^5.76.1, authKeys, AuthMeResponse, LoginRequest, LoginResponse (+27 more)

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
Nodes (24): stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary (+16 more)

### Community 6 - "package.json"
Cohesion: 0.05
Nodes (51): @biomejs/biome, class-variance-authority, clsx, geist, openapi-typescript, @radix-ui/react-label, @radix-ui/react-slot, react (+43 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.11
Nodes (17): loadUI(), loadUI(), Options, trackedResponseWriter, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement() (+9 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.12
Nodes (42): lucide-react, next, nextConfig, DashboardPage(), InitPage(), LoginPage(), HomePage(), AuthGuard() (+34 more)

### Community 11 - "testing.T"
Cohesion: 0.07
Nodes (57): TestAppCloseTeardownOrderIsLIFO(), stubAdminAuth, testHarness, TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency() (+49 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "time.Time"
Cohesion: 0.06
Nodes (30): CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, Options, stubAccountHealthTracker, stubUsageRecorder (+22 more)

### Community 14 - "repository"
Cohesion: 0.12
Nodes (15): Combo, createComboParams, DisabledModel, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, ModelAlias, repository (+7 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.20
Nodes (4): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion()

### Community 16 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (30): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+22 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Bus"
Cohesion: 0.09
Nodes (23): authHarness, fakeCredentialStore, Options, Event, New(), decodeErrorCode(), doRequest(), findCookie() (+15 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "ir.go"
Cohesion: 0.10
Nodes (46): CanonicalChatRequest, CanonicalChatResponse, CanonicalMessage, ImageSource, ResponseFormat, StreamTranslator, SystemBlock, ToolChoice (+38 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.14
Nodes (24): @hookform/resolvers, zod, zod@^4.6.7, buildSettingsPatchRequest(), ChangePasswordFormValues, changePasswordSchema, DEFAULT_SETTINGS_FORM_VALUES, isUiLocale() (+16 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "resolver.go"
Cohesion: 0.16
Nodes (13): resolveSnapshot, targetKey, stubModelLister, stubProviderResolver, ModelDescriptor, ResolvedTarget, buildResolvedTarget(), connectionSupportsBare() (+5 more)

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
Cohesion: 0.06
Nodes (49): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), actorFromContext(), decodeJSONBody(), decodePathModelParam() (+41 more)

### Community 30 - "Web Package Build Scripts"
Cohesion: 0.18
Nodes (11): scripts, build, check, check:fix, dev, format, format:check, gen:api (+3 more)

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
Cohesion: 0.15
Nodes (13): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder, streamEncoder, formatSSEFrame() (+5 more)

### Community 39 - "SettingItem"
Cohesion: 0.22
Nodes (6): int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, SettingItem, validatedUpdate

### Community 40 - "SettingsConsole.tsx"
Cohesion: 0.15
Nodes (23): DashboardSettingsPage(), useChangePasswordMutation(), usePatchSettingsMutation(), useRollbackSettingMutation(), useSettings(), useSettingsHistory(), useUpdateAuthPolicyMutation(), formatHistoryTimestamp() (+15 more)

### Community 41 - "encoding/json.RawMessage"
Cohesion: 0.15
Nodes (26): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+18 more)

### Community 42 - "buildApp"
Cohesion: 0.16
Nodes (13): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, BuildInfo (+5 more)

### Community 43 - "UpstreamCredentials"
Cohesion: 0.12
Nodes (9): fakeProviderStore, stubProviderStore, ProviderConnection, UpstreamCredentials, newFakeProviderStore(), cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall (+1 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.05
Nodes (6): Module, stubSystemModule, stubModule, orderStubModule, stubSystemModule, stubSystemModule

### Community 46 - "providers/repo.go"
Cohesion: 0.05
Nodes (33): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), Module, New() (+25 more)

### Community 47 - "context.Context"
Cohesion: 0.08
Nodes (5): stubAdminCredentialStore, stubAuthVerifier, stubAuthVerifier, actorFromContext(), stubAdminAuth

### Community 48 - "translate_openai.go"
Cohesion: 0.19
Nodes (19): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatFunctionSpec, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse (+11 more)

### Community 49 - "App"
Cohesion: 0.11
Nodes (4): Runner, job, syncBuffer, App

### Community 50 - "TestContractInterfacesSatisfied"
Cohesion: 0.15
Nodes (7): GatewayInterceptor, stubInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), applyStaticCacheHeaders()

### Community 51 - "AppShell.tsx"
Cohesion: 0.14
Nodes (16): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), DashboardLayout(), DashboardLayoutProps, useLogoutMutation(), applyThemeToDocument(), ThemeMode, UIState, useUIStore (+8 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.17
Nodes (19): responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason, responsesRequest, responsesResponse (+11 more)

### Community 53 - "New"
Cohesion: 0.20
Nodes (9): cloneAttrMap(), RingBuffer, New(), TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(), TestLevelFilteringAndChildLoggerAttrsAndGroups(), TestRingBufferCapsAndPreservesOrder(), TestRingBufferConcurrentWritesAndReads(), TestRingBufferRecentDeepClonesAttrs() (+1 more)

### Community 54 - "Config"
Cohesion: 0.16
Nodes (11): newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), Config, TestIsLoopbackClassification() (+3 more)

### Community 55 - "DB"
Cohesion: 0.12
Nodes (15): File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+7 more)

### Community 56 - "listSettingsWithQueryer"
Cohesion: 0.21
Nodes (8): listSettingsWithQueryer(), decodeStoredValue(), lookupDefinition(), applyResult, repository, rollbackResult, rowQueryer, UpdatedEvent

### Community 57 - "CanonicalStreamChunk"
Cohesion: 0.36
Nodes (4): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder, newOpenAIStreamDecoder()

### Community 59 - "parseJSONObjectOrWrap"
Cohesion: 0.31
Nodes (5): geminiPendingTool, geminiStreamEncoder, parseJSONObjectOrWrap(), mapCanonicalStopReasonToGemini(), newGeminiStreamEncoder()

### Community 60 - "newTestHarness"
Cohesion: 0.18
Nodes (11): newKeyCache(), New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation() (+3 more)

### Community 61 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 64 - "resolveStaticFile"
Cohesion: 0.53
Nodes (4): isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 312 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `DB` to `Binary Entrypoint & Embedded UI`, `buildApp`, `AppShell.tsx`, `Module`?**
  _High betweenness centrality (0.240) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.05416441482015252 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `AppShell.tsx` to `AppShell.tsx`, `DB`?**
  _High betweenness centrality (0.166) - this node is a cross-community bridge._
- **Should `SystemOverview.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.0664451827242525 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `SettingsConsole.tsx`, `AppShell.tsx`, `DB`?**
  _High betweenness centrality (0.145) - this node is a cross-community bridge._