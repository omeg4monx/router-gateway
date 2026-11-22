# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 123 files · ~103,549 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 1404 nodes · 4010 edges · 58 communities (52 shown, 6 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 358 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `355a4a7b`
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
- TestContractInterfacesSatisfied
- keyCache
- io/fs.FS
- KeysConsole.tsx
- Design System & Wave 0 Plan Docs
- newAuthHarness
- updates
- Root Layout, Geist & Query Provider
- SettingsConsole.tsx
- with
- jobs
- New
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
- LoginForm.tsx
- SettingItem
- auth/schemas.ts
- dependencies
- Module
- UpstreamCredentials
- providers/repo.go
- fakeCredentialStore
- Module
- context.Context
- DB
- App
- File Map
- newTestHarness
- tailwind-merge
- InitForm.tsx
- buildApp
- newTestHarness
- repository
- RouteMounts

## God Nodes (most connected - your core abstractions)
1. `WriteError()` - 44 edges
2. `WriteJSON()` - 36 edges
3. `buildApp()` - 35 edges
4. `Module` - 30 edges
5. `DB` - 30 edges
6. `SettingsConsole()` - 30 edges
7. `Module` - 27 edges
8. `KeysConsole()` - 25 edges
9. `extractErrorMessage()` - 25 edges
10. `Module` - 24 edges

## Surprising Connections (you probably didn't know these)
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests` --references--> `buildApp()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1b-providers-combos-gateway.md → cmd/omeg4mon-router/main.go
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment` --references--> `KeyInfo`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `File Map` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/modules/system/module.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (58 total, 6 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.06
Nodes (62): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+54 more)

### Community 1 - "SystemOverview.tsx"
Cohesion: 0.09
Nodes (24): openapi-fetch, @tanstack/react-query, @tanstack/react-query@^5.76.1, authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+16 more)

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
Cohesion: 0.13
Nodes (16): cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, scanSummary(), createParams, deleteResult, keySecretResp (+8 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.11
Nodes (15): loadUI(), loadUI(), Options, trackedResponseWriter, apiTimeoutMiddleware(), Server, New(), TestAdminAPIAuthEnforcement() (+7 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.13
Nodes (32): lucide-react, next, nextConfig, DashboardLayout(), DashboardLayoutProps, InitPage(), LoginPage(), HomePage() (+24 more)

### Community 11 - "testing.T"
Cohesion: 0.08
Nodes (45): TestSentinelErrorsDistinct(), newTestHarness(), TestCatalogAndAdapterRegistry(), TestConnectionQuickTest(), TestConnectionsCRUDAndSecretEncryptionHygiene(), TestCustomProviderNodesCRUD(), TestHeadlessOAuthFlowAndCallbackXSSProtection(), TestOAuthEnvResolutionAndCallbackJSONAndValidation() (+37 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TestContractInterfacesSatisfied"
Cohesion: 0.05
Nodes (30): GatewayInterceptor, ProviderResolver, ProxyDialerResolver, stubAccountHealthTracker, stubInterceptor, stubModelLister, stubProviderResolver, stubProxyDialerResolver (+22 more)

### Community 14 - "keyCache"
Cohesion: 0.20
Nodes (8): stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), hashToken(), isValidTokenFormat(), cachedEntry, keyCache

### Community 15 - "io/fs.FS"
Cohesion: 0.14
Nodes (11): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler() (+3 more)

### Community 16 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (29): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+21 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "newAuthHarness"
Cohesion: 0.22
Nodes (12): New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore(), newTestSigner(), TestAuthLoginLifecycle() (+4 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "Root Layout, Geist & Query Provider"
Cohesion: 0.24
Nodes (7): geist, geist@^1.4.2, metadata, RootLayout(), RootLayoutProps, QueryProvider(), QueryProviderProps

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (44): @hookform/resolvers, zod, DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest (+36 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "New"
Cohesion: 0.17
Nodes (11): TestAppCloseTeardownOrderIsLIFO(), NewRunner(), TestRunnerIdempotentStartStopAndGuards(), TestRunnerRecoversFromPanicAndStopsCleanly(), TestRunnerStopRespectsContextDeadline(), New(), TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(), TestLevelFilteringAndChildLoggerAttrsAndGroups() (+3 more)

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
Nodes (43): stubAuthVerifier, stubAuthVerifier, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), actorFromContext() (+35 more)

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

### Community 38 - "LoginForm.tsx"
Cohesion: 0.14
Nodes (21): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, DashboardPage(), useInitStatus(), useSystemVersion(), formatUpdatedAt(), formatUptime(), MODULE_SCOPES (+13 more)

### Community 39 - "SettingItem"
Cohesion: 0.10
Nodes (15): listSettingsWithQueryer(), decodeStoredValue(), int64Ptr(), isStrictJSONInteger(), lookupDefinition(), TestSettingDefinitionValidateRawIntegerBounds(), applyResult, patchSettingsRequest (+7 more)

### Community 40 - "auth/schemas.ts"
Cohesion: 0.36
Nodes (6): AuthRedirectTarget, AuthRouteState, AuthSessionPosture, LoginFormValues, loginSchema, normalizeAuthPathname()

### Community 41 - "dependencies"
Cohesion: 0.11
Nodes (19): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-dom, zustand, dependencies, class-variance-authority@^0.7.1, clsx@^2.1.1 (+11 more)

### Community 42 - "Module"
Cohesion: 0.15
Nodes (7): Task 2: `auth` Module & OpenAPI Fragment, Module, Registry, normalizeName(), nameMinHeap, registryState, testHarness

### Community 43 - "UpstreamCredentials"
Cohesion: 0.15
Nodes (7): stubProviderStore, ProviderConnection, UpstreamCredentials, cloneUpstreamCredentials(), sanitizeUpstreamError(), refreshCall, storedCredentials

### Community 44 - "providers/repo.go"
Cohesion: 0.13
Nodes (20): cloneStringMap(), cloneStringSlice(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+12 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.06
Nodes (6): Module, TestRegistryStartAllRollsBackOnPartialFailure(), stubSystemModule, orderStubModule, stubSystemModule, stubSystemModule

### Community 46 - "Module"
Cohesion: 0.07
Nodes (13): Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, cloneOAuthConfig(), defaultBuiltInAdapters(), OAuthConfig, oauthEnvPrefix(), Module, New(), buildAuthorizeURL() (+5 more)

### Community 47 - "context.Context"
Cohesion: 0.13
Nodes (9): authHarness, fakeCredentialStore, Options, stubAdminCredentialStore, AdminCredentialStore, actorFromContext(), actorFromContext(), TestSystemAdminCredentialStore() (+1 more)

### Community 48 - "DB"
Cohesion: 0.17
Nodes (10): Event, newKeyCache(), New(), newRepository(), New(), newRepository(), DB, Bus (+2 more)

### Community 49 - "App"
Cohesion: 0.11
Nodes (7): Runner, cloneAttrMap(), RingBuffer, job, syncBuffer, LogEntry, App

### Community 50 - "File Map"
Cohesion: 0.24
Nodes (9): KeyValidator, File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+1 more)

### Community 51 - "newTestHarness"
Cohesion: 0.22
Nodes (9): newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify(), HashPassword(), VerifyPassword() (+1 more)

### Community 53 - "InitForm.tsx"
Cohesion: 0.26
Nodes (15): react, react-hook-form, useLoginMutation(), LoginForm(), LoginFormProps, useInitializeSystem(), InitForm(), Alert (+7 more)

### Community 54 - "buildApp"
Cohesion: 0.13
Nodes (15): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, BuildInfo (+7 more)

### Community 60 - "newTestHarness"
Cohesion: 0.32
Nodes (8): assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), stubAdminAuth, testHarness

### Community 61 - "repository"
Cohesion: 0.19
Nodes (4): Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, repository

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 300 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `File Map` to `LoginForm.tsx`, `Module`, `Module`, `buildApp`?**
  _High betweenness centrality (0.250) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `WriteError()` (e.g. with `.adminAuthMiddleware()` and `.recoverJSON()`) actually correct?**
  _`WriteError()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.05674365554605076 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `LoginForm.tsx` to `File Map`, `AppShell.tsx`?**
  _High betweenness centrality (0.180) - this node is a cross-community bridge._
- **Should `SystemOverview.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.08817204301075268 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `SystemOverview.tsx`, `File Map`, `LoginForm.tsx`?**
  _High betweenness centrality (0.134) - this node is a cross-community bridge._