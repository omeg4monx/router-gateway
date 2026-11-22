# Graph Report - omeg4mon-router  (2026-11-23)

## Corpus Check
- Corpus is ~43,675 words - fits in a single context window. You may not need a graph.

## Summary
- 666 nodes · 1761 edges · 38 communities (29 shown, 9 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 120 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Binary Entrypoint & Embedded UI
- Next.js 15 App & shadcn/ui Shell
- HTTP Server, Config & Auth Contracts
- OpenAPI Client & System Hooks
- Rewrite Architecture Specification
- Embedded FS & SQL Migration Loader
- Frontend Build & Dev Tooling
- Structured Logging & Ring Buffer
- TypeScript Strict Compiler Config
- SQLite WAL Engine & System Repo
- Web Form, Store & UI Dependencies
- Lifecycle & Job Runner Test Suites
- shadcn/ui Component Configuration
- Gateway Interceptor Contracts
- Session Signing & Event Bus
- Secret Vault & App Shutdown
- Key & Provider Resolver Contracts
- Design System & Wave 0 Plan Docs
- Module Registry & Orchestration
- Password Schema & Zod Validation
- Root Layout, Geist & Query Provider
- Module Start/Stop Lifecycle Hooks
- HTTP/2 Client & Periodic Jobs
- App Bootstrap & Crypto Tests
- System Module & Password Hashing
- CLI Flag & Env Config Loading
- Logging Ring Buffer Tests
- HTTP/2 Client Streaming Tests
- SQLite WAL Pool & Migration Tests
- Topological Sort Min-Heap
- Web Package Build Scripts
- Generated OpenAPI TypeScript Schema
- Proxy Dialer Resolver Contract
- Thread-Safe Test Log Buffer
- JSON Marshaler Error Test Fixture
- Tailwind PostCSS Config
- Go Module Root Definition

## God Nodes (most connected - your core abstractions)
1. `buildApp()` - 29 edges
2. `Module` - 23 edges
3. `cn()` - 23 edges
4. `App` - 22 edges
5. `Server` - 19 edges
6. `DB` - 18 edges
7. `NewRegistry()` - 17 edges
8. `SystemOverview()` - 17 edges
9. `compilerOptions` - 17 edges
10. `New()` - 16 edges

## Surprising Connections (you probably didn't know these)
- `App` --references--> `Config`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/config/config.go
- `App` --references--> `SessionSigner`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/crypto/session.go
- `App` --references--> `DB`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/db/db.go
- `App` --references--> `Bus`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/eventbus/bus.go
- `App` --references--> `Server`  [EXTRACTED]
  cmd/omeg4mon-router/main.go → internal/platform/httpserver/server.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (38 total, 9 thin omitted)

### Community 0 - "Binary Entrypoint & Embedded UI"
Cohesion: 0.09
Nodes (16): loadUI(), adminSubjectKey, ErrorDetail, ErrorEnvelope, TestLoadOrCreateMasterKey(), isLinkUnsupported(), LoadOrCreateMasterKey(), readMasterKeyFile() (+8 more)

### Community 1 - "Next.js 15 App & shadcn/ui Shell"
Cohesion: 0.09
Nodes (56): class-variance-authority, lucide-react, next, react, nextConfig, DashboardLayout(), DashboardLayoutProps, DashboardPage() (+48 more)

### Community 2 - "HTTP Server, Config & Auth Contracts"
Cohesion: 0.06
Nodes (27): Config, AuthVerifier, stubAuthVerifier, Options, Server, stubAuthVerifier, trackedResponseWriter, TestIsLoopbackClassification() (+19 more)

### Community 3 - "OpenAPI Client & System Hooks"
Cohesion: 0.13
Nodes (19): openapi-fetch, openapi-fetch@^0.14.0, isHealthResponse(), systemKeys, useInitializeSystem(), useSystemHealth(), apiClient, createApiClient() (+11 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "Embedded FS & SQL Migration Loader"
Cohesion: 0.18
Nodes (12): loadUI(), migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), applyStaticCacheHeaders(), hasUnsafePath(), isMissingAssetPath() (+4 more)

### Community 6 - "Frontend Build & Dev Tooling"
Cohesion: 0.13
Nodes (20): openapi-typescript, tailwindcss, @tailwindcss/postcss, @types/bun, @types/node, @types/react, @types/react-dom, typescript (+12 more)

### Community 7 - "Structured Logging & Ring Buffer"
Cohesion: 0.14
Nodes (9): addAttrAt(), attrValue(), cloneAttrMap(), ensurePath(), parseLevel(), boundAttrs, LogEntry, RingBuffer (+1 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib, module (+11 more)

### Community 9 - "SQLite WAL Engine & System Repo"
Cohesion: 0.17
Nodes (6): DB, New(), newRepository(), BuildInfo, Options, repository

### Community 10 - "Web Form, Store & UI Dependencies"
Cohesion: 0.11
Nodes (19): clsx, @radix-ui/react-label, @radix-ui/react-slot, react-dom, react-hook-form, tailwind-merge, zustand, dependencies (+11 more)

### Community 11 - "Lifecycle & Job Runner Test Suites"
Cohesion: 0.20
Nodes (17): TestAppCloseTeardownOrderIsLIFO(), NewRunner(), TestRunnerIdempotentStartStopAndGuards(), TestRunnerRecoversFromPanicAndStopsCleanly(), TestRunnerStopRespectsContextDeadline(), NewRegistry(), TestRegistryCallbacksCanInspectRegistryWithoutDeadlock(), TestRegistryDeterministicLexicographicalTieBreak() (+9 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "Gateway Interceptor Contracts"
Cohesion: 0.21
Nodes (6): CanonicalRequest, GatewayInterceptor, ResponseMeta, StreamEvent, stubInterceptor, TestContractInterfacesSatisfied()

### Community 14 - "Session Signing & Event Bus"
Cohesion: 0.18
Nodes (8): Claims, SessionSigner, Bus, Event, New(), TestConcurrentPublishSubscribeUnsubscribe(), TestPublishDeliversToMultipleSubscribersOnMatchingTopic(), TestPublishNeverBlocksSlowSubscriber()

### Community 16 - "Key & Provider Resolver Contracts"
Cohesion: 0.18
Nodes (9): KeyInfo, KeyValidator, ProviderResolver, ResolvedTarget, stubKeyValidator, stubProviderResolver, stubUsageRecorder, UsageRecord (+1 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "Module Registry & Orchestration"
Cohesion: 0.23
Nodes (4): normalizeName(), Module, Registry, registryState

### Community 19 - "Password Schema & Zod Validation"
Cohesion: 0.23
Nodes (12): @hookform/resolvers, zod, @hookform/resolvers@^7.0.8, zod@^4.6.7, countPasswordRunes(), initPasswordSchema, InitPasswordValidation, InitPasswordValues (+4 more)

### Community 20 - "Root Layout, Geist & Query Provider"
Cohesion: 0.20
Nodes (9): geist, @tanstack/react-query, geist@^1.4.2, @tanstack/react-query@^5.76.1, metadata, RootLayout(), RootLayoutProps, QueryProvider() (+1 more)

### Community 22 - "HTTP/2 Client & Periodic Jobs"
Cohesion: 0.31
Nodes (4): Options, applyDefaults(), job, Runner

### Community 23 - "App Bootstrap & Crypto Tests"
Cohesion: 0.22
Nodes (9): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestBuildAppCleansUpOnPartialFailure(), TestDevUIHandlerDirectsToNextDevServer(), TestSessionSignVerifyAndExpiry(), TestVaultRoundTripAndTamperRejection(), NewSessionSigner() (+1 more)

### Community 24 - "System Module & Password Hashing"
Cohesion: 0.25
Nodes (9): newTestHarness(), TestSystemHealthAndVersion(), TestSystemInitConcurrentRace(), TestSystemInitLifecycle(), TestSystemShutdownRemoteBindGuard(), TestPasswordHashAndVerify(), HashPassword(), VerifyPassword() (+1 more)

### Community 25 - "CLI Flag & Env Config Loading"
Cohesion: 0.29
Nodes (8): main(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable()

### Community 26 - "Logging Ring Buffer Tests"
Cohesion: 0.29
Nodes (6): New(), TestErrorAttrsSerializeAsMessageInRingBufferAndJSON(), TestLevelFilteringAndChildLoggerAttrsAndGroups(), TestRingBufferCapsAndPreservesOrder(), TestRingBufferConcurrentWritesAndReads(), TestRingBufferRecentDeepClonesAttrs()

### Community 27 - "HTTP/2 Client Streaming Tests"
Cohesion: 0.33
Nodes (5): NewClient(), TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(), TestNewClientAppliesDefaults(), TestNewClientRequestContextCancellationAbortsStream(), TestNewClientUsesCustomProxyFunc()

### Community 28 - "SQLite WAL Pool & Migration Tests"
Cohesion: 0.33
Nodes (6): buildDSN(), Open(), TestConcurrentMigrateModuleIdempotent(), TestConcurrentReadersAndWritersUnderWAL(), TestMigrateModuleIdempotentAndAtomicRollback(), TestOpenAppliesWALPragmasAndSplitPools()

### Community 30 - "Web Package Build Scripts"
Cohesion: 0.33
Nodes (6): scripts, build, dev, gen:api, test, typecheck

### Community 31 - "Generated OpenAPI TypeScript Schema"
Cohesion: 0.33
Nodes (5): components, $defs, operations, paths, webhooks

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Wave 0 Platform Foundation Implementation Plan`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 132 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `buildApp()` connect `App Bootstrap & Crypto Tests` to `Binary Entrypoint & Embedded UI`, `HTTP Server, Config & Auth Contracts`, `SQLite WAL Engine & System Repo`, `Lifecycle & Job Runner Test Suites`, `Session Signing & Event Bus`, `Secret Vault & App Shutdown`, `Module Registry & Orchestration`, `Module Start/Stop Lifecycle Hooks`, `HTTP/2 Client & Periodic Jobs`, `CLI Flag & Env Config Loading`, `Logging Ring Buffer Tests`, `HTTP/2 Client Streaming Tests`, `SQLite WAL Pool & Migration Tests`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **Are the 12 inferred relationships involving `buildApp()` (e.g. with `.Close()` and `.Start()`) actually correct?**
  _`buildApp()` has 12 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/omeg4monx/router-gateway`, `Wave 0 Platform Foundation Implementation Plan`, `Structured Logging RingBuffer & Non-Blocking EventBus` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Binary Entrypoint & Embedded UI` be split into smaller, more focused modules?**
  _Cohesion score 0.09113607990012484 - nodes in this community are weakly interconnected._
- **Why does `App` connect `Secret Vault & App Shutdown` to `Binary Entrypoint & Embedded UI`, `Thread-Safe Test Log Buffer`, `HTTP Server, Config & Auth Contracts`, `Structured Logging & Ring Buffer`, `SQLite WAL Engine & System Repo`, `Session Signing & Event Bus`, `Module Registry & Orchestration`, `HTTP/2 Client & Periodic Jobs`, `App Bootstrap & Crypto Tests`, `HTTP/2 Client Streaming Tests`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `Next.js 15 App & shadcn/ui Shell` be split into smaller, more focused modules?**
  _Cohesion score 0.08724608724608725 - nodes in this community are weakly interconnected._
- **Why does `Module` connect `HTTP Server, Config & Auth Contracts` to `Binary Entrypoint & Embedded UI`, `Embedded FS & SQL Migration Loader`, `SQLite WAL Engine & System Repo`, `Session Signing & Event Bus`, `Secret Vault & App Shutdown`, `Module Registry & Orchestration`, `Module Start/Stop Lifecycle Hooks`, `System Module & Password Hashing`?**
  _High betweenness centrality (0.023) - this node is a cross-community bridge._