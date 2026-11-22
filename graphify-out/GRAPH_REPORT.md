# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- 77 files · ~44,486 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 5, .lock 1, .css 1)

## Summary
- 744 nodes · 1845 edges · 38 communities (31 shown, 7 thin omitted)
- Extraction: 94% EXTRACTED · 6% INFERRED · 0% AMBIGUOUS · INFERRED: 115 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5c449926`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- Next.js 15 App & shadcn/ui Shell
- Module
- OpenAPI Client & System Hooks
- Rewrite Architecture Specification
- Embedded FS & SQL Migration Loader
- Frontend Build & Dev Tooling
- Structured Logging & Ring Buffer
- TypeScript Strict Compiler Config
- SQLite WAL Engine & System Repo
- Web Form, Store & UI Dependencies
- testing.T
- shadcn/ui Component Configuration
- contracts.go
- Session Signing & Event Bus
- App
- RouteMounts
- Design System & Wave 0 Plan Docs
- Module Registry & Orchestration
- updates
- Root Layout, Geist & Query Provider
- Module Start/Stop Lifecycle Hooks
- with
- jobs
- .AuthenticateAdmin
- steps
- web-check
- on
- cross-compile
- Topological Sort Min-Heap
- Web Package Build Scripts
- Generated OpenAPI TypeScript Schema
- docker-build
- trackedResponseWriter
- JSON Marshaler Error Test Fixture
- Tailwind PostCSS Config
- Go Module Root Definition

## God Nodes (most connected - your core abstractions)
1. `buildApp()` - 29 edges
2. `Module` - 24 edges
3. `cn()` - 23 edges
4. `App` - 22 edges
5. `Server` - 19 edges
6. `with` - 19 edges
7. `SystemOverview()` - 17 edges
8. `NewRegistry()` - 17 edges
9. `compilerOptions` - 17 edges
10. `New()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `buildApp()` --references--> `BuildInfo`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/modules/system/module.go
- `buildApp()` --calls--> `New()`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/modules/system/module.go
- `App` --references--> `SessionSigner`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/crypto/session.go
- `App` --references--> `DB`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/db/db.go
- `App` --references--> `Bus`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/eventbus/bus.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (38 total, 7 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.09
Nodes (16): loadUI(), adminSubjectKey, ErrorDetail, ErrorEnvelope, TestLoadOrCreateMasterKey(), isLinkUnsupported(), LoadOrCreateMasterKey(), readMasterKeyFile() (+8 more)

### Community 1 - "Next.js 15 App & shadcn/ui Shell"
Cohesion: 0.07
Nodes (64): class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, DashboardPage() (+56 more)

### Community 2 - "Module"
Cohesion: 0.41
Nodes (5): Module, WriteError(), WriteJSON(), jsonMethodNotAllowedHandler(), jsonNotFoundHandler()

### Community 3 - "OpenAPI Client & System Hooks"
Cohesion: 0.14
Nodes (19): openapi-fetch, openapi-fetch@^0.14.0, isHealthResponse(), systemKeys, useInitializeSystem(), useSystemHealth(), apiClient, createApiClient() (+11 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "Embedded FS & SQL Migration Loader"
Cohesion: 0.19
Nodes (10): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), hasUnsafePath(), isMissingAssetPath(), NewStaticUIHandler(), openRegularFile() (+2 more)

### Community 6 - "Frontend Build & Dev Tooling"
Cohesion: 0.13
Nodes (20): openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+12 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), RingBuffer, parseLevel(), boundAttrs, LogEntry (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "SQLite WAL Engine & System Repo"
Cohesion: 0.19
Nodes (6): BuildInfo, New(), newRepository(), DB, Options, repository

### Community 10 - "Web Form, Store & UI Dependencies"
Cohesion: 0.09
Nodes (23): clsx, @hookform/resolvers, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form, tailwind-merge, zod (+15 more)

### Community 11 - "testing.T"
Cohesion: 0.06
Nodes (61): buildApp(), main(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), Options (+53 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "contracts.go"
Cohesion: 0.09
Nodes (18): GatewayInterceptor, KeyValidator, ProviderResolver, ProxyDialerResolver, stubInterceptor, stubKeyValidator, stubProviderResolver, stubProxyDialerResolver (+10 more)

### Community 14 - "Session Signing & Event Bus"
Cohesion: 0.18
Nodes (8): Claims, Event, SessionSigner, Bus, New(), TestConcurrentPublishSubscribeUnsubscribe(), TestPublishDeliversToMultipleSubscribersOnMatchingTopic(), TestPublishNeverBlocksSlowSubscriber()

### Community 15 - "App"
Cohesion: 0.06
Nodes (21): loadUI(), Options, Config, TestIsLoopbackClassification(), Vault, AdminSubject(), apiTimeoutMiddleware(), Server (+13 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Module Registry & Orchestration"
Cohesion: 0.21
Nodes (5): Module, Registry, normalizeName(), registryState, testHarness

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "Root Layout, Geist & Query Provider"
Cohesion: 0.20
Nodes (9): geist, @tanstack/react-query, geist@^1.4.2, @tanstack/react-query@^5.76.1, metadata, RootLayout(), RootLayoutProps, QueryProvider() (+1 more)

### Community 21 - "Module Start/Stop Lifecycle Hooks"
Cohesion: 0.18
Nodes (3): validatePassword(), stubModule, orderStubModule

### Community 22 - "with"
Cohesion: 0.13
Nodes (15): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+7 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, jobs, single-binary, name, permissions, contents (+4 more)

### Community 24 - ".AuthenticateAdmin"
Cohesion: 0.21
Nodes (4): stubAuthVerifier, stubAuthVerifier, AuthVerifier, stubAdminAuth

### Community 25 - "steps"
Cohesion: 0.27
Nodes (11): steps, steps, name, runs-on, steps, timeout-minutes, go-check, steps (+3 more)

### Community 26 - "web-check"
Cohesion: 0.25
Nodes (8): run, web-check, working-directory, defaults, name, runs-on, steps, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "cross-compile"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, cross-compile

### Community 30 - "Web Package Build Scripts"
Cohesion: 0.33
Nodes (6): scripts, build, dev, gen:api, test, typecheck

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 182 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `buildApp()` connect `testing.T` to `Binary Entrypoint & Embedded UI`, `SQLite WAL Engine & System Repo`, `Session Signing & Event Bus`, `App`, `Module Registry & Orchestration`, `Module Start/Stop Lifecycle Hooks`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Are the 12 inferred relationships involving `buildApp()` (e.g. with `.Close()` and `.Start()`) actually correct?**
  _`buildApp()` has 12 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.08982322025800286 - nodes in this community are weakly interconnected._
- **Why does `App` connect `App` to `Binary Entrypoint & Embedded UI`, `Structured Logging & Ring Buffer`, `SQLite WAL Engine & System Repo`, `testing.T`, `Session Signing & Event Bus`, `Module Registry & Orchestration`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Next.js 15 App & shadcn/ui Shell` be split into smaller, more focused modules?**
  _Cohesion score 0.07379979570990806 - nodes in this community are weakly interconnected._
- **Why does `Module` connect `Module` to `Binary Entrypoint & Embedded UI`, `Embedded FS & SQL Migration Loader`, `SQLite WAL Engine & System Repo`, `Session Signing & Event Bus`, `App`, `Module Registry & Orchestration`, `Module Start/Stop Lifecycle Hooks`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._