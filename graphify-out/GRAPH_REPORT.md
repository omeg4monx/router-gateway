# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 91 files · ~60,643 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 8 file(s) not represented in the graph (top: (none) 6, .lock 1, .css 1)

## Summary
- 998 nodes · 2552 edges · 54 communities (46 shown, 8 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 205 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8e11c08e`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- Next.js 15 App & shadcn/ui Shell
- biome.json
- settings
- Rewrite Architecture Specification
- Options
- Frontend Build & Dev Tooling
- Structured Logging & Ring Buffer
- TypeScript Strict Compiler Config
- DB
- Web Form, Store & UI Dependencies
- testing.T
- shadcn/ui Component Configuration
- TestContractInterfacesSatisfied
- Server
- static_ui.go
- RouteMounts
- Design System & Wave 0 Plan Docs
- newAuthHarness
- updates
- Root Layout, Geist & Query Provider
- App
- with
- jobs
- Module
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
- NewClient
- SettingItem
- File Map
- newTestHarness
- Config
- buildApp
- Module
- fakeCredentialStore
- Runner
- New
- Vault
- nameMinHeap
- New
- newTestHarness
- zod
- UsageRecord

## God Nodes (most connected - your core abstractions)
1. `Module` - 30 edges
2. `buildApp()` - 30 edges
3. `cn()` - 23 edges
4. `WriteError()` - 23 edges
5. `App` - 22 edges
6. `DB` - 20 edges
7. `Server` - 20 edges
8. `WriteJSON()` - 20 edges
9. `with` - 20 edges
10. `New()` - 19 edges

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

## Communities (54 total, 8 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.07
Nodes (39): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent, policyResponse (+31 more)

### Community 1 - "Next.js 15 App & shadcn/ui Shell"
Cohesion: 0.05
Nodes (86): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, @hookform/resolvers, lucide-react, next, openapi-fetch, react, nextConfig (+78 more)

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (47): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+39 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 6 - "Frontend Build & Dev Tooling"
Cohesion: 0.12
Nodes (22): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom (+14 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), RingBuffer, parseLevel(), boundAttrs, LogEntry (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "DB"
Cohesion: 0.20
Nodes (7): Global Constraints, New(), newRepository(), DB, Options, repository, repository

### Community 10 - "Web Form, Store & UI Dependencies"
Cohesion: 0.11
Nodes (19): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form, tailwind-merge, zustand, dependencies (+11 more)

### Community 11 - "testing.T"
Cohesion: 0.15
Nodes (23): TestAppCloseTeardownOrderIsLIFO(), buildDSN(), Open(), TestConcurrentMigrateModuleIdempotent(), TestConcurrentReadersAndWritersUnderWAL(), TestMigrateModuleIdempotentAndAtomicRollback(), TestOpenAppliesWALPragmasAndSplitPools(), NewRunner() (+15 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TestContractInterfacesSatisfied"
Cohesion: 0.14
Nodes (10): GatewayInterceptor, ProviderResolver, stubInterceptor, stubProviderResolver, CanonicalRequest, ResolvedTarget, ResponseMeta, StreamEvent (+2 more)

### Community 14 - "Server"
Cohesion: 0.19
Nodes (6): loadUI(), loadUI(), apiTimeoutMiddleware(), Server, wrapTrackedResponseWriter(), testHarness

### Community 15 - "static_ui.go"
Cohesion: 0.17
Nodes (10): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile() (+2 more)

### Community 16 - "RouteMounts"
Cohesion: 0.10
Nodes (4): RouteMounts, stubModule, orderStubModule, stubSystemModule

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

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - "Module"
Cohesion: 0.24
Nodes (4): authHarness, Event, Module, Bus

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
Cohesion: 0.10
Nodes (19): stubAuthVerifier, stubAuthVerifier, trackedResponseWriter, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody() (+11 more)

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

### Community 38 - "NewClient"
Cohesion: 0.29
Nodes (7): Options, applyDefaults(), NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 39 - "SettingItem"
Cohesion: 0.11
Nodes (14): listSettingsWithQueryer(), decodeStoredValue(), SettingItem, isStrictJSONInteger(), lookupDefinition(), applyResult, patchSettingsRequest, rollbackResult (+6 more)

### Community 40 - "File Map"
Cohesion: 0.21
Nodes (10): KeyValidator, stubKeyValidator, File Map, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification (+2 more)

### Community 41 - "newTestHarness"
Cohesion: 0.15
Nodes (12): indexByKey(), newTestHarness(), TestSettingsDefaultsAndPatch(), TestSettingsHistoryAndRollback(), TestSettingsPublishesEventBus(), TestSettingsValidationRejectsUnknownAndInvalid(), New(), TestConcurrentPublishSubscribeUnsubscribe() (+4 more)

### Community 42 - "Config"
Cohesion: 0.21
Nodes (5): Options, Claims, Config, TestIsLoopbackClassification(), SessionSigner

### Community 43 - "buildApp"
Cohesion: 0.24
Nodes (10): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, BuildInfo, New() (+2 more)

### Community 44 - "Module"
Cohesion: 0.33
Nodes (4): Module, Registry, normalizeName(), registryState

### Community 45 - "fakeCredentialStore"
Cohesion: 0.12
Nodes (3): fakeCredentialStore, stubAdminCredentialStore, AdminCredentialStore

### Community 46 - "Runner"
Cohesion: 0.25
Nodes (3): Runner, job, syncBuffer

### Community 47 - "New"
Cohesion: 0.33
Nodes (9): Options, New(), TestAdminAPIAuthEnforcement(), TestDevUIHandlerFallback(), TestGatewaySSEStreamsPastAPITimeout(), testLogger(), TestPanicRecoveryReturnsJSON500(), TestServerIsolatesJSON404FromSPAFallback() (+1 more)

### Community 48 - "Vault"
Cohesion: 0.33
Nodes (3): TestVaultRoundTripAndTamperRejection(), Vault, NewVault()

### Community 50 - "New"
Cohesion: 0.29
Nodes (6): New(), TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(), TestLevelFilteringAndChildLoggerAttrsAndGroups(), TestRingBufferCapsAndPreservesOrder(), TestRingBufferConcurrentWritesAndReads(), TestRingBufferRecentDeepClonesAttrs()

### Community 51 - "newTestHarness"
Cohesion: 0.29
Nodes (10): Task 2: `auth` Module & OpenAPI Fragment, newTestHarness(), TestSystemAdminCredentialStore(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify() (+2 more)

### Community 53 - "UsageRecord"
Cohesion: 0.60
Nodes (3): stubUsageRecorder, UsageRecorder, UsageRecord

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 253 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `File Map` connect `File Map` to `buildApp`, `newTestHarness`, `Module`, `Next.js 15 App & shadcn/ui Shell`?**
  _High betweenness centrality (0.209) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `Module` (e.g. with `File Map` and `Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter`) actually correct?**
  _`Module` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.06843507751937984 - nodes in this community are weakly interconnected._
- **Why does `Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out)` connect `Next.js 15 App & shadcn/ui Shell` to `File Map`?**
  _High betweenness centrality (0.162) - this node is a cross-community bridge._
- **Are the 13 inferred relationships involving `buildApp()` (e.g. with `.Close()` and `.Start()`) actually correct?**
  _`buildApp()` has 13 INFERRED edges - model-reasoned connections that need verification._
- **Should `Next.js 15 App & shadcn/ui Shell` be split into smaller, more focused modules?**
  _Cohesion score 0.05232646834477498 - nodes in this community are weakly interconnected._