# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 487 files · ~540,788 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5618 nodes · 18623 edges · 238 communities (206 shown, 32 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 1811 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `26b25bfc`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- main_test.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- repository
- KeysConsole.tsx
- TypeScript Strict Compiler Config
- ProvidersModals.tsx
- ModelEditorConsole.tsx
- testing.T
- shadcn/ui Component Configuration
- UpstreamCredentials
- MediaConsole.tsx
- App.tsx
- Module
- Design System & Wave 0 Plan Docs
- extractErrorMessage
- updates
- Button
- buildApp
- with
- jobs
- net/http/httptest.ResponseRecorder
- steps
- web-check
- on
- GatewayInterceptor
- TokenSaverConsole.tsx
- WriteJSON
- NotificationsConsole.tsx
- docker-build
- go-check
- combos/schemas.ts
- .handleCreateConnection
- Go Module Root Definition
- CanonicalUsage
- providers/hooks.ts
- QuotaConsole.tsx
- CanonicalStreamChunk
- main.go
- ir.go
- translate_gemini.go
- go_pkg_time
- repository
- PluginsStudio.tsx
- RouteMounts
- ModelDescriptor
- newTestHarness
- Badge
- translate_responses.go
- oauthLoopbackProxy
- .Fetch
- overview.ts
- PromptGuardConsole.tsx
- Provider Quota Menu & Token Harbor Quota Integration
- repository
- migrator.go
- ProviderAdapter
- fakeHandle
- UsageConsole.tsx
- repository
- net/http.Header
- teeHandler
- ProxyPoolsConsole.tsx
- TestContractInterfacesSatisfied
- navigation.ts
- pricingStore
- modelIntelligence.ts
- translate_openai.go
- TelegramConsole.tsx
- Module
- promptinjection/engine.go
- go_pkg_encoding_json
- arena/schemas.ts
- combos/module.go
- awsQuickRuntime
- readLegacySQLite
- web/package.json
- PWA Support Design — `omeg4mon-router`
- Module
- Review Focus
- paths
- Module
- collectDynamicHostTelemetry
- net/http.ResponseWriter
- settingDefinition
- repo
- DB
- sync.Mutex
- CanonicalChatResponse
- useUIStore
- SettingsConsole.tsx
- Server
- Runner
- go_pkg_encoding_json
- ref_node_path
- Module
- Module
- ConnectionSummary
- AttestationManager
- SecurityConsole.tsx
- Module
- net/http.Request
- post
- .identifyProviderModel
- repo
- encoding/json.RawMessage
- ProviderTopology.tsx
- repo
- repository
- repo
- properties
- evaluateAntigravityQuota
- net/url.URL
- Runner
- Module
- providers/catalogUtils.ts
- time.Time
- plugins/repo.go
- properties
- Module
- ResolvedTarget
- .ProvisionKey
- io/fs.FS
- repository
- OMEG4MON Router
- Module
- stubDep
- certManager
- repo
- modalUtils.ts
- Module
- properties
- ChatCompleter
- release-binaries
- steps
- codebuddy.go
- AddApiKeyModal
- modeleditor/engine.go
- deviceOAuthSession
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- fetch.go
- properties
- 200
- ProcessHandle
- properties
- dependencies
- cli/package.json
- setupTestMitm
- botAPIClient
- playgroundUtils.test.ts
- get
- scripts
- post
- schemas
- stubDiscoverySystemModule
- PreviewPluginsResponse
- stubModule
- Module
- stubSystemModule
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- setupTokenSaverTest
- Review Focus
- Review Focus
- LiveMonitorConsole.tsx
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- NewStreamTranslator
- Module
- parameters
- Module
- New
- stubModule
- ClearImporterHistoryResponse
- source_path
- BuildChatPayload
- buildModelProbeRequest
- docker-release
- Module
- `@omeg4mon/router` CLI Launcher
- stubDep
- .ExportCombos
- freepi.go
- context.Context
- stubSystemModule
- stubModule
- parameters
- .handleAwsQuickOAuthStart
- refreshCodeBuddyToken
- omeg4mon-router Verification Gate
- nameMinHeap
- redact.go
- single-binary
- keyLimiter
- verify-web
- io.Reader
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- parseIdentityPayload
- mapBusEventToAudit
- Feature Suite
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- Load
- New
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- time.Duration
- manifest.test.ts
- Review Focus
- jsonMarshalerError
- compressor.go
- palette.test.ts
- logger.go
- omeg4mon-router Project Standards

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 188 edges
2. `WriteJSON()` - 186 edges
3. `WriteError()` - 184 edges
4. `Button` - 94 edges
5. `DB` - 93 edges
6. `cn()` - 83 edges
7. `buildApp()` - 79 edges
8. `Module` - 78 edges
9. `Bus` - 72 edges
10. `Alert` - 71 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Approaches Considered` --references--> `quotaFetcher`  [INFERRED]
  docs/superpowers/specs/2026-11-27-provider-quota-design.md → internal/modules/providers/quota.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx
- `4.1 Code, Type Safety & Memory Safety` --references--> `ErrorBoundary`  [INFERRED]
  .agents/skills/omeg4mon-verify/SKILL.md → web/src/shared/ui/ErrorBoundary.tsx

## Import Cycles
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/module_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/optimization_e2e_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/provider_hardening_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (238 total, 32 thin omitted)

### Community 0 - "main_test.go"
Cohesion: 0.08
Nodes (16): collectStaticHostInfo(), readWindowsRegistryString(), TestLoadOrCreateMasterKey(), isLinkUnsupported(), LoadOrCreateMasterKey(), readMasterKeyFile(), writeKeyAtomically(), writeKeyExclusive() (+8 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.21
Nodes (5): decodeOptionalJSON(), decodeStrictJSON(), Module, newPoolID(), formatUTC()

### Community 2 - "biome.json"
Cohesion: 0.05
Nodes (46): source, assist, actions, enabled, noUnusedImports, noUnusedVariables, css, formatter (+38 more)

### Community 3 - "settings"
Cohesion: 0.05
Nodes (37): check-blank, check-type-assertions, generated, presets, rules, formatters, enable, settings (+29 more)

### Community 4 - "Rewrite Architecture Specification"
Cohesion: 0.13
Nodes (19): Single-Admin Auth & Remote-Bind Posture, Canonical IR & Direct Protocol Translators, CLI Tools Generator & Loopback-Only Apply, OMEG4MON Design System (docs/DESIGN.md), In-Process Typed EventBus, Gateway Routing & SSE Streaming Pipeline, Internal Contracts & GatewayInterceptor Chain, Next.js 15 Static Export & Feature-Modular UI (+11 more)

### Community 5 - "APIKeySummary"
Cohesion: 0.07
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+19 more)

### Community 6 - "repository"
Cohesion: 0.16
Nodes (11): boolToInt(), marshalOverrideJSON(), newRepository(), nullFloat(), nullInt(), scanBypass(), scanOverride(), BypassRule (+3 more)

### Community 7 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (70): zod, zod@^4.6.7, buildGatewayUrls(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState() (+62 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (52): formatShortTime(), MODAL_SIZE_CLASS, ModalSize, CatalogOAuthDescriptor, connectionToEditFormValues(), CreateConnectionFormInput, CreateConnectionFormOutput, createConnectionSchema (+44 more)

### Community 10 - "ModelEditorConsole.tsx"
Cohesion: 0.09
Nodes (48): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation(), useModelOverrides() (+40 more)

### Community 11 - "testing.T"
Cohesion: 0.02
Nodes (192): fakeProviderStore, roundTripFunc, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, Options (+184 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "UpstreamCredentials"
Cohesion: 0.08
Nodes (13): catalogFakeProviderStore, ModelTestResult, testAllModelsResponse, ProviderConnection, stubProviderStore, stubProviderStore, UpstreamCredentials, resolveKeylessProbeCredentials() (+5 more)

### Community 14 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 15 - "App.tsx"
Cohesion: 0.08
Nodes (36): App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio, InitPage (+28 more)

### Community 16 - "Module"
Cohesion: 0.10
Nodes (37): isBlockedAntigravityProjectID(), generateTraceID(), resolveAwsQuickTenantBaseURL(), postExchangeIFlow(), postExchangeNous(), isAnthropicTokenURL(), decodeJWTPayload(), defaultExchangePKCECode() (+29 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "extractErrorMessage"
Cohesion: 0.06
Nodes (66): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+58 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "Button"
Cohesion: 0.06
Nodes (110): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, lucide-react, react, react-hook-form, react-router, ArenaConsole(), ChatConsole() (+102 more)

### Community 21 - "buildApp"
Cohesion: 0.08
Nodes (34): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1BGatewayEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting(), TestWave2AE2ETranslatorAndCatalog() (+26 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.05
Nodes (36): Options, Claims, Task 2: `auth` Module & OpenAPI Fragment, Task 2: `tunnel` Module (Cloudflare Quick Tunnel & Tailscale Integration Manager), Module, New(), decodeErrorCode(), doRequest() (+28 more)

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "GatewayInterceptor"
Cohesion: 0.04
Nodes (59): main(), File Map, Global Constraints, Review Focus, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, Global Constraints (+51 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.13
Nodes (24): tokenSaverQueryKeys, useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId, matchHistoryPreset(), TOKEN_SAVER_LEVELS (+16 more)

### Community 30 - "WriteJSON"
Cohesion: 0.12
Nodes (10): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider() (+2 more)

### Community 31 - "NotificationsConsole.tsx"
Cohesion: 0.16
Nodes (24): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation() (+16 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "combos/schemas.ts"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 35 - ".handleCreateConnection"
Cohesion: 0.12
Nodes (22): Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, buildProtocolTestRequest(), isValidAuthType(), isValidProtocol(), validateHTTPURL(), parseAwsQuickPastedQuery() (+14 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.10
Nodes (20): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+12 more)

### Community 39 - "providers/hooks.ts"
Cohesion: 0.07
Nodes (52): buildProviderEntries(), buildProviderSectionBuckets(), computeFleetTelemetry(), countCategoryBuckets(), countStatusBuckets(), getProviderMonogram(), matchesProviderSearch(), matchesStatusFilter() (+44 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (47): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+39 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 42 - "main.go"
Cohesion: 0.21
Nodes (4): PoolTestResponse, ScanCandidateResult, ScanResponse, storedSetting

### Community 43 - "ir.go"
Cohesion: 0.17
Nodes (31): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+23 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "go_pkg_time"
Cohesion: 0.12
Nodes (17): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, ProviderStrategyState, targetKey (+9 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): buildLogsWhere(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.17
Nodes (27): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+19 more)

### Community 48 - "RouteMounts"
Cohesion: 0.24
Nodes (3): validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 49 - "ModelDescriptor"
Cohesion: 0.18
Nodes (12): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, vite, @vitejs/plugin-react, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin() (+4 more)

### Community 50 - "newTestHarness"
Cohesion: 0.09
Nodes (24): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, statusError, New(), assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart() (+16 more)

### Community 51 - "Badge"
Cohesion: 0.07
Nodes (60): buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup, qualifyProviderModel(), stripProviderPrefix() (+52 more)

### Community 52 - "translate_responses.go"
Cohesion: 0.13
Nodes (28): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+20 more)

### Community 53 - "oauthLoopbackProxy"
Cohesion: 0.15
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 54 - ".Fetch"
Cohesion: 0.12
Nodes (21): TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+13 more)

### Community 55 - "overview.ts"
Cohesion: 0.21
Nodes (12): filterModules(), formatClockUptime(), generateEndpointSnippet(), getModuleMetadata(), getOverviewCopy(), MODULE_REGISTRY, ModuleMetadata, OVERVIEW_COPY (+4 more)

### Community 56 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 57 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.15
Nodes (9): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+1 more)

### Community 58 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 59 - "migrator.go"
Cohesion: 0.08
Nodes (51): ModelOption, changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, PolicyChangedEvent (+43 more)

### Community 60 - "ProviderAdapter"
Cohesion: 0.10
Nodes (8): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), defaultBuiltInAdapters(), xrouterBuiltInOAuthAdapters(), builtinAdapter, customTestAdapter, OAuthConfig, ProviderAdapter

### Community 61 - "fakeHandle"
Cohesion: 0.13
Nodes (3): fakeHandle, osProcessHandle, ssePipeReader

### Community 62 - "UsageConsole.tsx"
Cohesion: 0.07
Nodes (61): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, useCombos(), useProviderConnections(), useSystemTelemetry() (+53 more)

### Community 63 - "repository"
Cohesion: 0.16
Nodes (8): decodeStringMap(), decodeStringSlice(), newRepository(), scanConnectionRecord(), scanNode(), ProviderNode, repository, rowScanner

### Community 64 - "net/http.Header"
Cohesion: 0.13
Nodes (9): memoryResponseWriter, extractVendorHeaders(), extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration() (+1 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (35): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+27 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.06
Nodes (17): stubModelLister, matchesAllowedModel(), stubAdminAuth, AuthVerifier, ModelDescriptor, stubAuthVerifier, stubModelLister, stubModelLister (+9 more)

### Community 68 - "navigation.ts"
Cohesion: 0.11
Nodes (23): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, ActiveNavContext, ALL_NAV_ITEMS, EXTENSION_GROUPS (+15 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 71 - "translate_openai.go"
Cohesion: 0.18
Nodes (19): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+11 more)

### Community 72 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 73 - "Module"
Cohesion: 0.13
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "promptinjection/engine.go"
Cohesion: 0.10
Nodes (22): actionRank(), builtinRules(), compileRule(), effectiveRuleAction(), evaluateText(), highestAction(), inspectAndMaybeSanitizeBody(), isPromptTextKey() (+14 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.12
Nodes (24): extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags(), TestParseDiscoveredModelsFormats(), catalogAdapterDescriptor (+16 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (26): ArenaBattleFormInput, ArenaBattleFormOutput, arenaBattleFormSchema, ArenaBattleListResponse, ArenaContenderResult, ArenaLeaderboardEntry, ArenaLeaderboardResponse, ArenaModelListResponse (+18 more)

### Community 77 - "combos/module.go"
Cohesion: 0.12
Nodes (11): ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse, listProviderStrategiesResponse, patchComboRequest (+3 more)

### Community 78 - "awsQuickRuntime"
Cohesion: 0.25
Nodes (6): validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickSession, awsQuickTokens, awsQuickUserInfo

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "web/package.json"
Cohesion: 0.12
Nodes (23): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+15 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.19
Nodes (8): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken()

### Community 83 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 86 - "collectDynamicHostTelemetry"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 87 - "net/http.ResponseWriter"
Cohesion: 0.08
Nodes (23): Module, trackedResponseWriter, isAntigravityPromptRejection(), isBlockedAntigravityProjectID(), moveGeminiSystemInstructionToContents(), anthropicErrorType(), geminiErrorStatus(), openAIErrorType() (+15 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "repo"
Cohesion: 0.14
Nodes (10): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+2 more)

### Community 90 - "DB"
Cohesion: 0.04
Nodes (73): Options, authHarness, Options, testHarness, File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment (+65 more)

### Community 91 - "sync.Mutex"
Cohesion: 0.07
Nodes (14): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), Module, New(), TestUsageStreamMountsOnAdminStreamAPI() (+6 more)

### Community 92 - "CanonicalChatResponse"
Cohesion: 0.15
Nodes (18): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeResponse(), IsSupportedProtocol(), RewriteModelField(), encodeAnthropicResponse() (+10 more)

### Community 93 - "useUIStore"
Cohesion: 0.14
Nodes (26): Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, useLogoutMutation(), useSystemHealth(), useSystemVersion(), PwaInstallButton(), PwaInstallButtonProps, PwaStatusBanner() (+18 more)

### Community 94 - "SettingsConsole.tsx"
Cohesion: 0.05
Nodes (74): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, @tanstack/react-query, authKeys, AuthMeResponse, LoginRequest, LoginResponse, LogoutResponse (+66 more)

### Community 95 - "Server"
Cohesion: 0.08
Nodes (29): TestAppWave1AAuthSettingsKeysEndToEnd(), loadUI(), loadUI(), Task 5: Wire `auth`, `settings`, `keys` into `main.go`, Root OpenAPI & TypeScript Codegen, Options, doPxpipeJSON(), setupTestPxpipe(), TestPxpipeInterceptorTransformsAndRuntimeGating() (+21 more)

### Community 96 - "Runner"
Cohesion: 0.19
Nodes (5): JobInspector, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 97 - "go_pkg_encoding_json"
Cohesion: 0.21
Nodes (8): ChatCompletionRequest, ChatMessage, errorDisposition, handledGatewayError, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus(), isClientPayloadErrorBody()

### Community 98 - "ref_node_path"
Cohesion: 0.09
Nodes (19): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath() (+11 more)

### Community 99 - "Module"
Cohesion: 0.11
Nodes (10): decodeStrictJSON(), copyNonHopHeaders(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL() (+2 more)

### Community 100 - "Module"
Cohesion: 0.23
Nodes (4): clientIP(), decodeStrictJSON(), Module, isValidSeverity()

### Community 101 - "ConnectionSummary"
Cohesion: 0.15
Nodes (17): cloneOAuthConfig(), cloneStringMap(), cloneStringSlice(), oauthEnvPrefix(), isCodeBuddyHostFor(), isTokenHarborHost(), encodeStringMap(), encodeStringSlice() (+9 more)

### Community 102 - "AttestationManager"
Cohesion: 0.13
Nodes (14): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar(), AttestationManager, isTransientMintStatus() (+6 more)

### Community 103 - "SecurityConsole.tsx"
Cohesion: 0.19
Nodes (20): securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES, AuditPaginationState (+12 more)

### Community 104 - "Module"
Cohesion: 0.15
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO, Notification

### Community 105 - "net/http.Request"
Cohesion: 0.07
Nodes (15): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), extractAPIKey(), decodeJSON() (+7 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".identifyProviderModel"
Cohesion: 0.13
Nodes (21): ModelIdentityResult, extractProbeTokenUsage(), extractUpstreamProbeError(), sanitizeProbeError(), analyzeIdentityProbeResponse(), inferModelFamily(), modelsLooselyMatch(), normalizeModelForCompare() (+13 more)

### Community 108 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 109 - "encoding/json.RawMessage"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.14
Nodes (19): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, ProviderTopology(), ProviderTopologyProps (+11 more)

### Community 111 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 112 - "repository"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 116 - "net/url.URL"
Cohesion: 0.16
Nodes (10): ProxyDialerResolver, stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan() (+2 more)

### Community 117 - "Runner"
Cohesion: 0.07
Nodes (32): TestAppCloseTeardownOrderIsLIFO(), TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests (+24 more)

### Community 118 - "Module"
Cohesion: 0.10
Nodes (18): listCooldownsResponse, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), providerDisplayName(), splitTargetIdentifier(), actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module (+10 more)

### Community 119 - "providers/catalogUtils.ts"
Cohesion: 0.13
Nodes (23): ADAPTER_TO_REGISTRY_ID, computeEntryStats(), FleetTelemetry, getEntrySortRank(), ProviderCategory, ProviderCategoryTab, ProviderEntry, ProviderEntryStats (+15 more)

### Community 120 - "time.Time"
Cohesion: 0.04
Nodes (34): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, attemptOutcome (+26 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.06
Nodes (53): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+45 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "Module"
Cohesion: 0.23
Nodes (4): BinariesStatus, BinaryStatus, Module, Status

### Community 124 - "ResolvedTarget"
Cohesion: 0.07
Nodes (25): resolveSnapshot, strategyState, weightedCandidate, stubProviderResolver, Approaches Considered, Constraints, Goals, Intent & Why (+17 more)

### Community 125 - ".ProvisionKey"
Cohesion: 0.29
Nodes (4): KeyProvisioner, ProvisionedKey, ProvisionKeyRequest, stubKeyProvisioner

### Community 126 - "io/fs.FS"
Cohesion: 0.04
Nodes (13): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile() (+5 more)

### Community 127 - "repository"
Cohesion: 0.10
Nodes (20): Combo, ComboStep, createComboParams, DisabledModel, ModelAlias, repository, rowScanner, updateComboValues (+12 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "Module"
Cohesion: 0.17
Nodes (5): truncateError(), buildDiscoveryCandidateURLs(), TestAutoDetectModelsDiscoveryCandidates(), Module, quotaFetcher

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - "modalUtils.ts"
Cohesion: 0.18
Nodes (18): Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useTestProviderNodeDraftMutation(), buildTestFingerprint(), BulkCredentialEntry (+10 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "ChatCompleter"
Cohesion: 0.17
Nodes (10): arenaTestEnv, stubCompleter, ChatCompletionResult, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), ChatCompleter, TestGatewayCompleteChat() (+2 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "codebuddy.go"
Cohesion: 0.25
Nodes (10): applyCodeBuddyHeaders(), cloneShallowMap(), codebuddyDomainFor(), extractCodeBuddyText(), isCodeBuddyProvider(), sanitizeCodeBuddyChannelMarkers(), transformCodeBuddyRequest(), applyFreePIHeaders() (+2 more)

### Community 140 - "AddApiKeyModal"
Cohesion: 0.17
Nodes (16): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, useCreateProviderNodeMutation(), AddApiKeyModal(), AddCustomNodeModal(), DraftTestResultCard(), getCredentialFieldMeta(), ModalShell(), commaListSchema() (+8 more)

### Community 141 - "modeleditor/engine.go"
Cohesion: 0.25
Nodes (14): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), globToRegexp(), hasBlockText(), hasPromptSegment(), injectLatestUserAnchor() (+6 more)

### Community 142 - "deviceOAuthSession"
Cohesion: 0.20
Nodes (5): originOfURL(), generatePKCE(), deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "fetch.go"
Cohesion: 0.28
Nodes (8): extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget(), TestMediaSSRFReservedRangesAndRedirectBlock(), WebFetchRequest, WebFetchResponse

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "ProcessHandle"
Cohesion: 0.25
Nodes (3): fakeRunner, osProcessRunner, ProcessHandle

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

### Community 154 - "playgroundUtils.test.ts"
Cohesion: 0.29
Nodes (9): computeSessionTelemetry(), filterChatSessions(), formatLatencyMs(), formatTokenCount(), SessionTelemetry, sampleMessages, sampleSessions, ChatMessage (+1 more)

### Community 155 - "get"
Cohesion: 0.17
Nodes (13): description, delete, get, operationId, responses, summary, tags, operationId (+5 more)

### Community 156 - "scripts"
Cohesion: 0.17
Nodes (12): scripts, build, check, check:fix, dev, format, format:check, gen:api (+4 more)

### Community 157 - "post"
Cohesion: 0.29
Nodes (8): post, operationId, requestBody, summary, tags, requestBody, content, required

### Community 158 - "schemas"
Cohesion: 0.29
Nodes (6): components, schemas, additionalProperties, required, type, PluginState

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 162 - "Module"
Cohesion: 0.22
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - "setupTokenSaverTest"
Cohesion: 0.14
Nodes (20): File Map, Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestModelBypassRulesPriorityAndMetadataFlags() (+12 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "LiveMonitorConsole.tsx"
Cohesion: 0.09
Nodes (40): KeyedSystemLog, levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, StreamState, systemLevelClass(), appendMonitorEntry() (+32 more)

### Community 171 - "schemas"
Cohesion: 0.22
Nodes (9): schemas, additionalProperties, required, type, additionalProperties, required, type, ImportedProvisionedKey (+1 more)

### Community 172 - "properties"
Cohesion: 0.22
Nodes (9): type, properties, type, type, id, KeyPrefix, Name, Secret (+1 more)

### Community 173 - "RunLegacyImportResponse"
Cohesion: 0.22
Nodes (9): provisioned_keys, run, type, $ref, additionalProperties, properties, required, type (+1 more)

### Community 174 - "omeg4mon-router Versioning Verificator"
Cohesion: 0.25
Nodes (7): 1. Gather Evidence First, 2. Public Contract Surfaces, 3. SemVer Decision Matrix, 4. Per-Task vs Release-Cut Rules, 5. Rationalizations & Red Flags, 6. Required Exit Report, omeg4mon-router Versioning Verificator

### Community 176 - "version"
Cohesion: 0.22
Nodes (9): version, on, push, workflow_dispatch, tags, default, description, required (+1 more)

### Community 177 - "NewStreamTranslator"
Cohesion: 0.20
Nodes (10): NewStreamTranslator(), buildCompactedResponse(), randomResponseID(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder() (+2 more)

### Community 178 - "Module"
Cohesion: 0.16
Nodes (5): decodeStrictJSON(), LogEntry, LogListResponse, Module, Status

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "Module"
Cohesion: 0.21
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 181 - "New"
Cohesion: 0.11
Nodes (23): File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator, New() (+15 more)

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "BuildChatPayload"
Cohesion: 0.18
Nodes (9): BuildChatPayload(), extractContentText(), extractToolArguments(), formatToolsInstruction(), NewNDJSONToOpenAISSEReader(), randomUUIDv4(), StripModelPrefix(), TestBuildChatPayload() (+1 more)

### Community 186 - "buildModelProbeRequest"
Cohesion: 0.27
Nodes (10): buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreeHeaders(), applyOpenCodeFreePayload(), generateOpenCodeSessionID(), isOpenCodeFreeCreds(), marshalProbeJSON(), openCodeCanonicalTools() (+2 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "Module"
Cohesion: 0.21
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "freepi.go"
Cohesion: 0.43
Nodes (6): clampFreePIMaxTokens(), compareSemver(), maybeLiftFreePIClientVersion(), parseFreePIRequiredVersion(), TestFreePIClampMaxTokensFloor(), TestFreePIParseRequiredVersionAndCompare()

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (18): fakeCredentialStore, AdminCredentialStore, SettingsSnapshotter, stubAdminCredentialStore, ImportRun, stubSettings, newFakeStore(), boolToInt() (+10 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - ".handleAwsQuickOAuthStart"
Cohesion: 0.40
Nodes (4): buildAwsQuickSigninURL(), generateAwsQuickNonce(), generateOAuthState(), oauthStartRequest

### Community 198 - "refreshCodeBuddyToken"
Cohesion: 0.40
Nodes (5): applyCodeBuddyProbeHeaders(), codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), codebuddyAuthSpec

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 201 - "redact.go"
Cohesion: 0.07
Nodes (28): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), formatUTC(), newEventID() (+20 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "keyLimiter"
Cohesion: 0.24
Nodes (5): newKeyLimiter(), TestKeyLimiterShardedConcurrencyAndRPM(), keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "parseIdentityPayload"
Cohesion: 1.00
Nodes (3): extractIDPrefix(), parseIdentityPayload(), parseSSEIdentityPayload()

### Community 209 - "mapBusEventToAudit"
Cohesion: 0.67
Nodes (3): mapBusEventToAudit(), payloadToMap(), AuditInput

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 223 - "Load"
Cohesion: 0.25
Nodes (9): runMain(), CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable() (+1 more)

### Community 224 - "New"
Cohesion: 0.11
Nodes (14): newAwsQuickRuntime(), newDeviceOAuthRuntime(), New(), TestEnrichGitHubCopilotAndKiroProfileARN(), TestOAuthProviderRegistryMetadata(), cloneConnectionQuota(), defaultQuotaFetchers(), newQuotaCache() (+6 more)

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 231 - "time.Duration"
Cohesion: 0.05
Nodes (9): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, New(), computePrefixFingerprint() (+1 more)

### Community 233 - "manifest.test.ts"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 235 - "Review Focus"
Cohesion: 0.25
Nodes (7): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update

### Community 242 - "compressor.go"
Cohesion: 0.13
Nodes (27): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressText(), estimateTokens(), injectAnthropicCacheControl(), injectOutputDirectives() (+19 more)

### Community 246 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 250 - "logger.go"
Cohesion: 0.22
Nodes (4): HistoryItem, historyListResp, historyResponse, ValidationError

### Community 253 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1049 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `AddApiKeyModal` to `.handleCreateConnection`, `modalUtils.ts`, `providers/hooks.ts`?**
  _High betweenness centrality (0.215) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `main_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08156028368794327 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `.handleCreateConnection` to `context.Context`, `refreshCodeBuddyToken`, `net/http.Request`, `go_pkg_encoding_json`, `AddApiKeyModal`, `UpstreamCredentials`, `testing.T`, `Module`, `.Fetch`?**
  _High betweenness centrality (0.181) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `AddApiKeyModal` to `modalUtils.ts`, `providers/hooks.ts`, `ProvidersModals.tsx`, `Badge`, `Button`?**
  _High betweenness centrality (0.072) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._