# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 113 files · ~81,314 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 1191 nodes · 3289 edges · 55 communities (49 shown, 6 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 243 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b3f947fd`
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
- login/page.tsx
- io/fs.FS
- KeysConsole.tsx
- Design System & Wave 0 Plan Docs
- newAuthHarness
- updates
- Root Layout, Geist & Query Provider
- SettingsConsole.tsx
- with
- jobs
- newTestHarness
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
- auth/hooks.ts
- dependencies
- repository
- .AuthenticateAdmin
- buildApp
- fakeCredentialStore
- NewClient
- Bus
- newTestHarness
- App
- nameMinHeap
- newTestHarness
- .Resolve
- UsageRecord
- openapi-fetch

## God Nodes (most connected - your core abstractions)
1. `buildApp()` - 34 edges
2. `SettingsConsole()` - 30 edges
3. `Module` - 30 edges
4. `WriteError()` - 30 edges
5. `KeysConsole()` - 25 edges
6. `DB` - 25 edges
7. `extractErrorMessage()` - 25 edges
8. `WriteJSON()` - 25 edges
9. `cn()` - 23 edges
10. `Button` - 23 edges

## Surprising Connections (you probably didn't know these)
- `Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification` --references--> `AppShell()`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → web/src/shared/ui/AppShell.tsx
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `File Map` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/modules/system/module.go
- `Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter` --references--> `Module`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/modules/system/module.go
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

## Communities (55 total, 6 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.07
Nodes (38): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+30 more)

### Community 1 - "SystemOverview.tsx"
Cohesion: 0.11
Nodes (24): DashboardPage(), isHealthResponse(), useInitStatus(), useSystemHealth(), useSystemVersion(), formatUpdatedAt(), formatUptime(), MODULE_SCOPES (+16 more)

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
Nodes (36): KeyValidator, stubKeyValidator, File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+28 more)

### Community 6 - "package.json"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), RingBuffer, parseLevel(), boundAttrs, LogEntry (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "Server"
Cohesion: 0.11
Nodes (17): loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, New(), newRepository(), Server, New() (+9 more)

### Community 10 - "AppShell.tsx"
Cohesion: 0.15
Nodes (17): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), DashboardLayout(), DashboardLayoutProps, useLogoutMutation(), applyThemeToDocument(), ThemeMode, UIState, useUIStore (+9 more)

### Community 11 - "testing.T"
Cohesion: 0.12
Nodes (29): TestAppCloseTeardownOrderIsLIFO(), buildDSN(), Open(), TestConcurrentMigrateModuleIdempotent(), TestConcurrentReadersAndWritersUnderWAL(), TestMigrateModuleIdempotentAndAtomicRollback(), TestOpenAppliesWALPragmasAndSplitPools(), NewRunner() (+21 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TestContractInterfacesSatisfied"
Cohesion: 0.19
Nodes (7): GatewayInterceptor, stubInterceptor, CanonicalRequest, ResponseMeta, StreamEvent, TestContractInterfacesSatisfied(), applyStaticCacheHeaders()

### Community 14 - "login/page.tsx"
Cohesion: 0.24
Nodes (23): lucide-react, next, nextConfig, InitPage(), LoginPage(), HomePage(), AuthGuard(), AuthGuardProps (+15 more)

### Community 15 - "io/fs.FS"
Cohesion: 0.15
Nodes (10): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile() (+2 more)

### Community 16 - "KeysConsole.tsx"
Cohesion: 0.12
Nodes (30): DashboardKeysPage(), keysQueryKeys, ListApiKeysResponse, UpdateApiKeyVariables, useApiKeys(), useCreateApiKeyMutation(), useDeleteApiKeyMutation(), useRotateApiKeyMutation() (+22 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "newAuthHarness"
Cohesion: 0.22
Nodes (14): New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore(), newTestSigner(), TestAuthLoginLifecycle() (+6 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "Root Layout, Geist & Query Provider"
Cohesion: 0.20
Nodes (9): geist, @tanstack/react-query, geist@^1.4.2, @tanstack/react-query@^5.76.1, metadata, RootLayout(), RootLayoutProps, QueryProvider() (+1 more)

### Community 21 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (43): @hookform/resolvers, DashboardSettingsPage(), AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys, UpdateAuthPolicyRequest, useChangePasswordMutation() (+35 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "newTestHarness"
Cohesion: 0.15
Nodes (12): assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeysCRUDAndSecretHygiene(), TestKeysExpirationAndValidationErrors(), TestKeysValidateAndImmediateCacheInvalidation(), New(), TestConcurrentPublishSubscribeUnsubscribe() (+4 more)

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
Cohesion: 0.07
Nodes (26): trackedResponseWriter, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), Module, actorFromContext() (+18 more)

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
Cohesion: 0.15
Nodes (18): class-variance-authority, react, react-hook-form, LoginFormProps, cn(), AlertProps, alertVariants, Badge() (+10 more)

### Community 39 - "SettingItem"
Cohesion: 0.10
Nodes (17): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, SettingItem, isStrictJSONInteger(), lookupDefinition(), applyResult, historyListResp (+9 more)

### Community 40 - "auth/hooks.ts"
Cohesion: 0.15
Nodes (12): authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse, useLoginMutation(), AuthRedirectTarget, AuthRouteState (+4 more)

### Community 41 - "dependencies"
Cohesion: 0.10
Nodes (21): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge, zod, zustand, dependencies (+13 more)

### Community 42 - "repository"
Cohesion: 0.16
Nodes (9): main(), runMain(), Global Constraints, New(), newRepository(), DB, Options, repository (+1 more)

### Community 43 - ".AuthenticateAdmin"
Cohesion: 0.22
Nodes (3): stubAuthVerifier, stubAuthVerifier, stubAdminAuth

### Community 44 - "buildApp"
Cohesion: 0.18
Nodes (12): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), TestVaultRoundTripAndTamperRejection(), NewVault() (+4 more)

### Community 45 - "fakeCredentialStore"
Cohesion: 0.08
Nodes (5): stubAdminCredentialStore, stubSystemModule, stubModule, orderStubModule, stubSystemModule

### Community 46 - "NewClient"
Cohesion: 0.33
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 47 - "Bus"
Cohesion: 0.11
Nodes (9): authHarness, fakeCredentialStore, Options, Claims, Event, AdminCredentialStore, SessionSigner, Bus (+1 more)

### Community 48 - "newTestHarness"
Cohesion: 0.21
Nodes (9): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsTimestampsFixedWidthUTC(), TestSettingsValidationRejectsUnknownAndInvalid(), stubAdminAuth (+1 more)

### Community 49 - "App"
Cohesion: 0.10
Nodes (5): Vault, Runner, job, syncBuffer, App

### Community 51 - "newTestHarness"
Cohesion: 0.15
Nodes (14): Task 2: `auth` Module & OpenAPI Fragment, BuildInfo, newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard() (+6 more)

### Community 52 - ".Resolve"
Cohesion: 0.22
Nodes (5): ProviderResolver, ProxyDialerResolver, stubProviderResolver, stubProxyDialerResolver, ResolvedTarget

### Community 53 - "UsageRecord"
Cohesion: 0.60
Nodes (3): stubUsageRecorder, UsageRecorder, UsageRecord

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 282 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `APIKeySummary` to `Server`, `AppShell.tsx`, `newTestHarness`, `Module`?**
  _High betweenness centrality (0.261) - this node is a cross-community bridge._
- **Are the 14 inferred relationships involving `buildApp()` (e.g. with `.Close()` and `.Start()`) actually correct?**
  _`buildApp()` has 14 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.0690631808278867 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `AppShell.tsx` to `SystemOverview.tsx`, `APIKeySummary`?**
  _High betweenness centrality (0.187) - this node is a cross-community bridge._
- **Should `SystemOverview.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.10752688172043011 - nodes in this community are weakly interconnected._
- **Why does `AppShell()` connect `AppShell.tsx` to `LoginForm.tsx`, `APIKeySummary`, `login/page.tsx`?**
  _High betweenness centrality (0.143) - this node is a cross-community bridge._