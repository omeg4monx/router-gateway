# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 528 files · ~583,001 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 6049 nodes · 20675 edges · 251 communities (224 shown, 27 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1930 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `daa5cb01`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- SystemOverview.tsx
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- modeleditor/repo.go
- NewRegistry
- TypeScript Strict Compiler Config
- setupArenaTest
- MediaConsole.tsx
- testing.T
- shadcn/ui Component Configuration
- App.tsx
- ProviderDetailView.tsx
- TopProgressBar.tsx
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- SettingsConsole.tsx
- updates
- KeysConsole.tsx
- DB
- with
- jobs
- Module
- steps
- web-check
- on
- vault.go
- TokenSaverConsole.tsx
- ProvidersConsole.tsx
- SystemOverview.tsx
- docker-build
- go-check
- .executeCandidateAttempt
- ModelEditorConsole.tsx
- Go Module Root Definition
- CanonicalUsage
- ProvidersModals.tsx
- QuotaConsole.tsx
- CanonicalStreamChunk
- fingerprint.go
- ir.go
- extractErrorMessage
- main.go
- repository
- PluginsStudio.tsx
- Module
- ModelDescriptor
- setupModelEditorTest
- Module
- SecurityConsole.tsx
- New
- DecodeRequest
- awsQuickRuntime
- Module
- web/package.json
- Module
- migrator.go
- Badge
- oauthLoopbackProxy
- LiveMonitorConsole.tsx
- repository
- newTestHarness
- teeHandler
- proxypools/schemas.ts
- TelegramConsole.tsx
- Module
- pricingStore
- NotificationsConsole.tsx
- repo
- Module
- WriteJSON
- go_pkg_encoding_json
- Label
- go_pkg_context
- telemetry.go
- ComboSnapshotter
- time.Duration
- PWA Support Design — `omeg4mon-router`
- Module
- venice/tool_bridge.go
- paths
- PromptGuardConsole.tsx
- stubCompleter
- net/http.ResponseWriter
- settingDefinition
- translate_anthropic.go
- Review Focus
- MitmConsole.tsx
- Module
- Module
- deviceOAuthSession
- UpstreamCredentials
- net/http.Client
- strategies.go
- responsesStreamEncoder
- net/http.ResponseWriter
- modelIntelligence.ts
- keyLimiter
- net/http.Request
- orderStubModule
- VersionResolver
- evaluateAntigravityQuota
- post
- palette.test.ts
- repo
- New
- ProviderTopology.tsx
- Load
- repository
- repo
- properties
- batcher
- net/url.URL
- Module
- .Reader
- .ListModels
- AddApiKeyModal
- plugins/repo.go
- properties
- fakeHandle
- translate_openai.go
- setupTokenSaverTest
- io/fs.FS
- JobInspector
- OMEG4MON Router
- TunnelConsole.tsx
- ResolvedTarget
- certManager
- repo
- AttestationManager
- applyUpstreamHeaders
- properties
- sync.RWMutex
- release-binaries
- steps
- ResponseMeta
- UpstreamCredentials
- repository
- ref_bun_test
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- Module
- properties
- 200
- newTestHarness
- properties
- dependencies
- cli/package.json
- navigation.ts
- TransformRequest
- net/http.Handler
- get
- scripts
- post
- schemas
- manifest.test.ts
- PreviewPluginsResponse
- translate_gemini.go
- newHarness
- Module
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- .UpdateCombo
- Review Focus
- Review Focus
- routingStateCache
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- stubSystemModule
- version
- mapCodeBuddyQuota
- setupMediaTest
- parameters
- responseCoalescer
- repository
- resolve.js
- ClearImporterHistoryResponse
- source_path
- Runner
- listSettingsWithQueryer
- docker-release
- GatewayInterceptor
- `@omeg4mon/router` CLI Launcher
- stubDep
- stubDiscoverySystemModule
- go_pkg_encoding_json
- context.Context
- setupTunnelTest
- redact.go
- parameters
- quotaCache
- UpstreamCredentials
- omeg4mon-router Verification Gate
- repo
- repository
- single-binary
- omeg4mon-router Project Standards
- verify-web
- stubModule
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- Module
- .ProvisionKey
- Feature Suite
- setupCatalogTestModule
- ProviderEntry
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- newAuthHarness
- Module
- settingDefinition
- ComboSnapshot
- SessionSigner
- net/http.Header
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- scanCombo
- ChatCompleter
- DB
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- readLegacySQLite
- New
- stubDep
- stubModule
- optimization.test.ts
- stubModule
- repository
- jsonMarshalerError
- botAPIClient
- compressor.go
- stubSystemModule
- stubDepModule
- translate_responses.go
- stubSystem
- headroom.go
- newOAuthLoopbackProxy
- Review Focus

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 194 edges
2. `WriteJSON()` - 187 edges
3. `WriteError()` - 184 edges
4. `Button` - 106 edges
5. `cn()` - 97 edges
6. `DB` - 93 edges
7. `buildApp()` - 79 edges
8. `Badge()` - 76 edges
9. `Bus` - 75 edges
10. `Alert` - 75 edges

## Surprising Connections (you probably didn't know these)
- `Review Focus` --references--> `KeyValidator`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1a-auth-settings-keys.md → internal/contracts/contracts.go
- `Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3a-teams-vouchers.md → internal/contracts/contracts.go
- `Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests` --references--> `GatewayInterceptor`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave3e-plugins-cloudsync-importer-cli.md → internal/contracts/contracts.go
- `Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment` --references--> `Config`  [INFERRED]
  docs/superpowers/plans/2026-11-23-wave1c-usage-clitools.md → internal/platform/config/config.go
- `Design System & UI/UX Identity` --references--> `ErrorBoundary`  [INFERRED]
  .agents/rules/project-standards.md → web/src/shared/ui/ErrorBoundary.tsx

## Import Cycles
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/jetbrains_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/jetbrains_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/jetbrains_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/combo_runtime_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/combo_runtime_test.go`
- 2-file cycle: `internal/modules/gateway/freepi_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/freepi_test.go`
- 2-file cycle: `internal/modules/gateway/freepi_test.go -> internal/modules/gateway/jetbrains_test.go -> internal/modules/gateway/freepi_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (251 total, 27 thin omitted)

### Community 0 - "SystemOverview.tsx"
Cohesion: 0.05
Nodes (62): Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, authKeys, AuthMeResponse (+54 more)

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
Cohesion: 0.06
Nodes (27): stubKeyValidator, stubKeyValidator, KeyInfo, cloneCachedEntry(), cloneKeyInfo(), cloneStringSlice(), newKeyCache(), decodeStringSlice() (+19 more)

### Community 6 - "modeleditor/repo.go"
Cohesion: 0.05
Nodes (48): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+40 more)

### Community 7 - "NewRegistry"
Cohesion: 0.04
Nodes (63): TestAppWave1BGatewayEndToEnd(), roundTripFunc, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestCombosRoutingStateCacheHitsAndInvalidates(), TestKeylessFreeProviderModelProbeAndIdentify() (+55 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (20): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+12 more)

### Community 9 - "setupArenaTest"
Cohesion: 0.12
Nodes (13): arenaTestEnv, roundTripFunc, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation(), newAntigravityQuotaMonitor(), TestAntigravityQuotaMonitorLeaderCancelDoesNotAbortWaiters() (+5 more)

### Community 10 - "MediaConsole.tsx"
Cohesion: 0.13
Nodes (30): mediaQueryKeys, useCreateMediaComboMutation(), useCreateMediaProviderMutation(), useDeleteMediaComboMutation(), useDeleteMediaProviderMutation(), useMediaCombos(), useMediaProviders(), useUpdateMediaComboMutation() (+22 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (103): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+95 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "App.tsx"
Cohesion: 0.08
Nodes (35): @tanstack/react-query, App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio (+27 more)

### Community 14 - "ProviderDetailView.tsx"
Cohesion: 0.05
Nodes (79): StrategyCategory, runModelIdentityProbe(), ARENA_MODELS_QUERY_KEY, BatchToggleConnectionsVariables, DiscoveredModel, ExchangeProviderOAuthVariables, fetchConnectionModels(), FetchConnectionModelsResponse (+71 more)

### Community 15 - "TopProgressBar.tsx"
Cohesion: 0.30
Nodes (10): computeNextTrickleProgress(), createTopProgressController(), getTopProgressBarStyle(), routeProgressController, shouldFinishTopProgress(), TOP_PROGRESS_INITIAL, TOP_PROGRESS_MAX_TRICKLE, TopProgressBar() (+2 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.08
Nodes (48): isBlockedAntigravityProjectID(), generateTraceID(), resolveAwsQuickTenantBaseURL(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), TestCodeBuddyTokenRefresh(), buildJetBrainsAuthURL(), cleanJetBrainsAuthCode() (+40 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "SettingsConsole.tsx"
Cohesion: 0.09
Nodes (45): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, zod, zod@^4.6.7, AuthPolicyResponse, ChangePasswordRequest, ChangePasswordResponse, settingsQueryKeys (+37 more)

### Community 19 - "updates"
Cohesion: 0.11
Nodes (17): patterns, patterns, actions, go-dependencies, web-dependencies, dependency-name, versions, interval (+9 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.08
Nodes (62): buildGatewayUrls(), copyTextToClipboard(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isIpLiteral() (+54 more)

### Community 21 - "DB"
Cohesion: 0.05
Nodes (42): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2EProviderProxyRouting() (+34 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "Module"
Cohesion: 0.15
Nodes (7): CanonicalRequest, decodeJSONBody(), isValidMode(), isValidSensitivity(), newID(), validateRule(), Module

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "vault.go"
Cohesion: 0.15
Nodes (11): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+3 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.11
Nodes (29): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId (+21 more)

### Community 30 - "ProvidersConsole.tsx"
Cohesion: 0.08
Nodes (43): ADAPTER_TO_REGISTRY_ID, buildProviderEntries(), buildProviderSectionBuckets(), canAddProviderConnection(), canProbeProviderModels(), computeEntryStats(), computeFleetTelemetry(), countCategoryBuckets() (+35 more)

### Community 31 - "SystemOverview.tsx"
Cohesion: 0.04
Nodes (135): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), class-variance-authority, lucide-react, react, react-router, AuthGuard(), AuthGuardProps (+127 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - ".executeCandidateAttempt"
Cohesion: 0.08
Nodes (37): attemptOutcome, applyOpenCodeFreeHeaders(), isAntigravityPromptRejection(), moveGeminiSystemInstructionToContents(), extractUpstreamErrorMessage(), isCodeBuddyProvider(), clampFreePIMaxTokens(), resolveFreePIMinMaxTokens() (+29 more)

### Community 35 - "ModelEditorConsole.tsx"
Cohesion: 0.07
Nodes (59): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation() (+51 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.13
Nodes (16): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiStreamEncoder, openAIStreamEncoder, streamEncoder, streamResult, StreamTranslator (+8 more)

### Community 39 - "ProvidersModals.tsx"
Cohesion: 0.06
Nodes (59): getCredentialFieldMeta(), MODAL_SIZE_CLASS, ModalSize, CatalogOAuthDescriptor, commaListSchema(), connectionToEditFormValues(), CreateConnectionFormInput, CreateConnectionFormOutput (+51 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (52): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+44 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.15
Nodes (14): CanonicalStreamChunk, geminiStreamDecoder, openAIStreamDecoder, responsesStreamDecoder, streamDecoder, NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamDecoder() (+6 more)

### Community 42 - "fingerprint.go"
Cohesion: 0.17
Nodes (24): aliasHistoryToolCalls(), appendMissingCanonicalTools(), appendSystemText(), applyEnvFingerprint(), applyFingerprintTools(), canNormalizeInPlace(), canonicalToolDefinitions(), CanonicalTools() (+16 more)

### Community 43 - "ir.go"
Cohesion: 0.16
Nodes (29): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+21 more)

### Community 44 - "extractErrorMessage"
Cohesion: 0.06
Nodes (60): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+52 more)

### Community 45 - "main.go"
Cohesion: 0.10
Nodes (23): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+15 more)

### Community 46 - "repository"
Cohesion: 0.13
Nodes (10): escapeLikePattern(), LogsFilter, LogsResponse, ModelBreakdown, ProviderBreakdown, repository, SeriesBucket, SummaryResponse (+2 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.15
Nodes (31): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+23 more)

### Community 48 - "Module"
Cohesion: 0.10
Nodes (7): stubInterceptor, validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 49 - "ModelDescriptor"
Cohesion: 0.18
Nodes (12): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, vite, @vitejs/plugin-react, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin() (+4 more)

### Community 50 - "setupModelEditorTest"
Cohesion: 0.12
Nodes (20): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, TestGatewayInterceptorMetadataPropagatesToResponseHeaders(), doJSON(), setupModelEditorTest() (+12 more)

### Community 51 - "Module"
Cohesion: 0.15
Nodes (5): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO

### Community 52 - "SecurityConsole.tsx"
Cohesion: 0.18
Nodes (21): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES (+13 more)

### Community 53 - "New"
Cohesion: 0.08
Nodes (26): File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan, ProtocolTranslator, New() (+18 more)

### Community 54 - "DecodeRequest"
Cohesion: 0.17
Nodes (20): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString() (+12 more)

### Community 55 - "awsQuickRuntime"
Cohesion: 0.12
Nodes (9): newAwsQuickRuntime(), parseAwsQuickPastedQuery(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens (+1 more)

### Community 56 - "Module"
Cohesion: 0.16
Nodes (3): BuildInfo, Module, validatePassword()

### Community 57 - "web/package.json"
Cohesion: 0.12
Nodes (23): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+15 more)

### Community 58 - "Module"
Cohesion: 0.10
Nodes (10): decodeStrictJSON(), copyNonHopHeaders(), isLoopbackHost(), matchToolForHost(), newID(), singleJoiningSlash(), validateLoopbackListenAddr(), validateLoopbackURL() (+2 more)

### Community 59 - "migrator.go"
Cohesion: 0.07
Nodes (57): ModelOption, PolicyChangedEvent, ChangedEvent, createComboRequest, deleteDisabledModelRequest, listAliasesResponse, listCombosResponse, listDisabledModelsResponse (+49 more)

### Community 60 - "Badge"
Cohesion: 0.05
Nodes (89): ComboCard, ComboCardProps, buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep(), ProviderAccountGroup (+81 more)

### Community 61 - "oauthLoopbackProxy"
Cohesion: 0.21
Nodes (4): codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 62 - "LiveMonitorConsole.tsx"
Cohesion: 0.07
Nodes (57): attrValueClass(), levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, normalizeUsagePayload(), StreamState, SYSTEM_LEVELS (+49 more)

### Community 63 - "repository"
Cohesion: 0.12
Nodes (17): cloneStringMap(), cloneStringSlice(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice(), insertConnectionTx(), isUniqueConstraintError() (+9 more)

### Community 64 - "newTestHarness"
Cohesion: 0.14
Nodes (19): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+11 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "proxypools/schemas.ts"
Cohesion: 0.11
Nodes (34): proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation(), useUpdateProxyPoolMutation() (+26 more)

### Community 67 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 68 - "Module"
Cohesion: 0.05
Nodes (15): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), allBuiltInAdapters(), cloneOAuthConfig(), defaultBuiltInAdapters(), oauthEnvPrefix(), xrouterBuiltInOAuthAdapters(), newDeviceOAuthRuntime(), Module (+7 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.13
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "WriteJSON"
Cohesion: 0.10
Nodes (11): Module, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID(), decodeStrictJSON(), WriteJSON(), LogEntry (+3 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.10
Nodes (31): providerIDContextKey, extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags(), TestParseDiscoveredModelsFormats() (+23 more)

### Community 76 - "Label"
Cohesion: 0.07
Nodes (53): ArenaConsole(), ChatConsole(), STARTER_PROMPTS, arenaQueryKeys, useArenaBattles(), useArenaLeaderboard(), useArenaModels(), useChatSessionDetail() (+45 more)

### Community 77 - "go_pkg_context"
Cohesion: 0.20
Nodes (3): PoolTestResponse, storedSetting, ValidationError

### Community 78 - "telemetry.go"
Cohesion: 0.13
Nodes (25): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+17 more)

### Community 79 - "ComboSnapshotter"
Cohesion: 0.26
Nodes (11): Global Constraints, Review Focus, Task 1: Cross-Module Snapshot/Import Contracts & `plugins` Module, Task 2: `cloudsync` Module (Vaulted Remote Sync, Snapshots, Conflict Policies & OpenAPI), Task 3: `importer` Module (Read-Only Legacy `data.sqlite` Inspector, Dry-Run & Importer & OpenAPI), Task 4: `cmd/omeg4mon-router` Wiring, OpenAPI Codegen & E2E Tests, Task 5: Thin Multi-Platform npm Binary Launcher (`cli/`), Task 6: Web Studio UI (Plugins, Skills, Cloud Sync, Importer), Navigation, Single-Binary Build & Final Verification (+3 more)

### Community 80 - "time.Duration"
Cohesion: 0.05
Nodes (22): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, stubAccountHealthTracker, chainStore (+14 more)

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.22
Nodes (7): actorFromContext(), decodeJSONBody(), Module, normalizeExpiresAt(), validateName(), validateStringSlice(), generateToken()

### Community 83 - "venice/tool_bridge.go"
Cohesion: 0.19
Nodes (19): ApplyToolBridge(), buildToolCatalog(), collectSchemaProperties(), couldBeToolMarkerPrefix(), extractArgsFromObject(), extractCandidatesFromText(), extractNameFromObject(), ExtractRequestTools() (+11 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 86 - "stubCompleter"
Cohesion: 0.33
Nodes (6): stubCompleter, ChatMessage, ChatCompletionRequest, ChatCompletionResult, extractGatewayError(), stringifyMessageContent()

### Community 87 - "net/http.ResponseWriter"
Cohesion: 0.06
Nodes (36): Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, trackedResponseWriter, providerDisplayName(), splitTargetIdentifier(), isValidAuthType() (+28 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "translate_anthropic.go"
Cohesion: 0.19
Nodes (18): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+10 more)

### Community 90 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 91 - "MitmConsole.tsx"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 92 - "Module"
Cohesion: 0.14
Nodes (4): decodeJSONBody(), decodePathModelParam(), Module, isValidRate()

### Community 93 - "Module"
Cohesion: 0.15
Nodes (4): Module, validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 94 - "deviceOAuthSession"
Cohesion: 0.16
Nodes (7): applyCodeBuddyProbeHeaders(), codebuddySpecFor(), originOfURL(), codebuddyAuthSpec, deviceOAuthResult, deviceOAuthSession, rfc8628DeviceSpec

### Community 95 - "UpstreamCredentials"
Cohesion: 0.11
Nodes (24): UpstreamCredentials, buildModelProbeRequest(), buildIdentityProbeRequest(), applyOpenCodeFreePayload(), isOpenCodeFreeCreds(), marshalProbeJSON(), applyJetBrainsHeaders(), buildProtocolTestRequest() (+16 more)

### Community 96 - "net/http.Client"
Cohesion: 0.07
Nodes (36): ModelIdentityResult, ModelTestResult, testAllModelsResponse, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError(), analyzeIdentityProbeResponse() (+28 more)

### Community 97 - "strategies.go"
Cohesion: 0.18
Nodes (10): ProviderStrategyState, strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), fastRandFloat64(), fastRandIntN() (+2 more)

### Community 98 - "responsesStreamEncoder"
Cohesion: 0.20
Nodes (9): geminiPendingTool, responsesStreamBlockState, responsesStreamEncoder, buildCompactedResponse(), randomResponseID(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), aggregatedAnthropicBlock (+1 more)

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.16
Nodes (8): listCooldownsResponse, actorFromContext(), decodeJSONBody(), decodePathModelParam(), Module, validateDescription(), validateModelIdentifier(), validateStrategy()

### Community 100 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 101 - "keyLimiter"
Cohesion: 0.27
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 102 - "net/http.Request"
Cohesion: 0.05
Nodes (27): stubAdminAuth, stubAuthVerifier, Module, stubAuthVerifier, clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie() (+19 more)

### Community 104 - "VersionResolver"
Cohesion: 0.47
Nodes (3): NewVersionResolver(), TestVersionResolverRefreshesFromUpdater(), VersionResolver

### Community 105 - "evaluateAntigravityQuota"
Cohesion: 0.20
Nodes (13): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+5 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 108 - "repo"
Cohesion: 0.14
Nodes (11): mapBusEventToNotification(), buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), EmitPayload (+3 more)

### Community 109 - "New"
Cohesion: 0.18
Nodes (16): New(), generateSelfSignedCertFiles(), TestAdminAPIAuthEnforcement(), TestAdminAuthEnforcedOverHTTPSOnLoopback(), TestDevUIHandlerFallback(), TestGatewaySSEStreamsPastAPITimeout(), testLogger(), TestPanicRecoveryReturnsJSON500() (+8 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.13
Nodes (21): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, EnergyTopologyEdge(), FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, prefersReducedMotion() (+13 more)

### Community 111 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

### Community 112 - "repository"
Cohesion: 0.50
Nodes (3): boolToInt(), defaultConfig(), Config

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "batcher"
Cohesion: 0.11
Nodes (9): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher, ModelPricing (+1 more)

### Community 116 - "net/url.URL"
Cohesion: 0.09
Nodes (22): stubProxyDialerResolver, stubProxyResolver, ProxyDialerResolver, extractReadableContent(), isBlockedAddr(), resolveAndCheckHostSSRF(), validateFetchURL(), validateRedirectTarget() (+14 more)

### Community 117 - "Module"
Cohesion: 0.16
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 118 - ".Reader"
Cohesion: 0.11
Nodes (19): main(), runMain(), Global Constraints, Review Focus, Wave 1B — Providers, Combos & Core LLM Gateway Implementation Plan, Global Constraints, Review Focus, Wave 1C — Observability (`usage`) & Developer Ergonomics (`clitools`) Implementation Plan (+11 more)

### Community 119 - ".ListModels"
Cohesion: 0.25
Nodes (7): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, ModelDescriptor, compositeModelLister, stubModelLister

### Community 120 - "AddApiKeyModal"
Cohesion: 0.15
Nodes (25): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, useCreateProviderNodeMutation(), useTestProviderNodeDraftMutation() (+17 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.05
Nodes (59): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+51 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "fakeHandle"
Cohesion: 0.07
Nodes (8): prefixReadCloser, streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ssePipeReader

### Community 124 - "translate_openai.go"
Cohesion: 0.22
Nodes (17): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+9 more)

### Community 125 - "setupTokenSaverTest"
Cohesion: 0.24
Nodes (12): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestResponsesAPIFunctionCallOutputRTKFilter(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory() (+4 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.11
Nodes (8): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 127 - "JobInspector"
Cohesion: 0.15
Nodes (18): TestWave2CDisabledModulesDetachCleanly(), TestWave2CSecurityNotificationsAndJobsE2E(), File Map, Task 1: Platform Jobs Observability, Overlap Guard & `contracts.JobInspector`, Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, JobInspector, setupNotificationsTest() (+10 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "TunnelConsole.tsx"
Cohesion: 0.14
Nodes (21): react-hook-form, tunnelQueryKeys, useStartTunnelMutation(), useStopTunnelMutation(), useTunnelEvents(), useTunnelStatus(), useUpdateTunnelConfigMutation(), StartTunnelRequest (+13 more)

### Community 130 - "ResolvedTarget"
Cohesion: 0.11
Nodes (17): Module, resolveSnapshot, stubProviderResolver, stubResolver, ResolvedTarget, applyComboRuntimeMetadata(), buildResolvedTarget(), configIntValue() (+9 more)

### Community 131 - "certManager"
Cohesion: 0.16
Nodes (6): newCertManager(), randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.24
Nodes (6): boolToInt(), defaultStoredConfig(), nowTimestamp(), randomHexID(), Event, storedConfig

### Community 133 - "AttestationManager"
Cohesion: 0.12
Nodes (17): attestationBucketKey(), EvictSession(), AttestationManager, isTransientMintStatus(), NewAttestationManager(), normalizeProxyKey(), ResolveMintURL(), TestAttestationManagerBatchCacheExpiryAndSingleflight() (+9 more)

### Community 134 - "applyUpstreamHeaders"
Cohesion: 0.11
Nodes (17): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+9 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "ResponseMeta"
Cohesion: 0.14
Nodes (4): recordingInterceptor, ResponseBodyInterceptor, ResponseMeta, StreamEvent

### Community 140 - "UpstreamCredentials"
Cohesion: 0.06
Nodes (40): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+32 more)

### Community 141 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 142 - "ref_bun_test"
Cohesion: 0.11
Nodes (14): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, consoleSource, detailSource, modalsSource, providersDir, modalsSource (+6 more)

### Community 143 - "`/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split"
Cohesion: 0.29
Nodes (6): Approaches Considered, Constraints, `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split, Goals, Intent & Why, Non-Goals

### Community 144 - "transform.go"
Cohesion: 0.35
Nodes (9): collapseMessageWhitespace(), collapseRuns(), containsString(), dedupeSystemMessages(), marshalNoEscape(), runTransformPipeline(), stripNulls(), windowHistoryMessages() (+1 more)

### Community 145 - "Module"
Cohesion: 0.21
Nodes (3): actorFromContext(), decodeJSONBody(), Module

### Community 146 - "properties"
Cohesion: 0.13
Nodes (15): type, type, type, type, type, dry_run, import_combos, import_keys (+7 more)

### Community 147 - "200"
Cohesion: 0.19
Nodes (13): content, description, content, description, content, description, schema, application/json (+5 more)

### Community 148 - "newTestHarness"
Cohesion: 0.10
Nodes (22): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, statusError, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart() (+14 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 152 - "navigation.ts"
Cohesion: 0.16
Nodes (19): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, ALL_NAV_ITEMS, countVisibleNavItems(), EXTENSION_GROUPS, findActiveNavContext(), getShellCopy(), isModuleRegistered() (+11 more)

### Community 153 - "TransformRequest"
Cohesion: 0.13
Nodes (18): peekFreePIPinMismatch(), TestPeekFreePIPinMismatch(), aggregateAnthropicSSEToJSON(), AggregateOpenAISSEToJSON(), aggregateResponsesSSEToJSON(), AggregateSSEToJSON(), dropAdaptiveReasoningEffort(), MarshalNoEscapeHTML() (+10 more)

### Community 154 - "net/http.Handler"
Cohesion: 0.10
Nodes (29): loadUI(), loadUI(), New(), setupImporterTest(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD() (+21 more)

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

### Community 159 - "manifest.test.ts"
Cohesion: 0.25
Nodes (6): PNG_MAGIC, PUBLIC_DIR, WEB_ROOT, WebAppManifest, WebAppManifestIcon, THEME_COLOR_BY_MODE

### Community 160 - "PreviewPluginsResponse"
Cohesion: 0.14
Nodes (14): items, type, type, type, items, type, additionalProperties, properties (+6 more)

### Community 161 - "translate_gemini.go"
Cohesion: 0.19
Nodes (19): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+11 more)

### Community 162 - "newHarness"
Cohesion: 0.17
Nodes (16): Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), Options, newHarness(), TestProxyPoolsBindingLifecycleAndCacheInvalidation(), TestProxyPoolsCRUDCredentialRedactionAndValidation(), TestProxyPoolsEndToEndHTTPForwardingRotationAndAuthProbe(), TestProxyPoolsLatencyTestAndAutoScanner(), applyDefaults() (+8 more)

### Community 163 - "Module"
Cohesion: 0.15
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 167 - ".UpdateCombo"
Cohesion: 0.20
Nodes (14): ComboStep, createComboParams, updateComboValues, validateComboTargets(), cloneStringSlice(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig() (+6 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "routingStateCache"
Cohesion: 0.28
Nodes (3): cachedRoutingState, routingStateCache, newRoutingStateCache()

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

### Community 177 - "mapCodeBuddyQuota"
Cohesion: 0.14
Nodes (18): TestCodeBuddyQuotaExcludesExpiredBonusAndZeroCapacityIsExhausted(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota(), parseCodeBuddyNumber() (+10 more)

### Community 178 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "responseCoalescer"
Cohesion: 0.43
Nodes (3): cachedGatewayResponse, inflightCall, responseCoalescer

### Community 181 - "repository"
Cohesion: 0.26
Nodes (3): DisabledModel, ModelAlias, repository

### Community 182 - "resolve.js"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - "Runner"
Cohesion: 0.17
Nodes (7): JobStatus, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 186 - "listSettingsWithQueryer"
Cohesion: 0.18
Nodes (11): listSettingsWithQueryer(), decodeStoredValue(), HistoryItem, lookupDefinition(), applyResult, historyListResp, historyResponse, repository (+3 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "GatewayInterceptor"
Cohesion: 0.28
Nodes (15): File Map, Task 1: Cross-Module Contracts for Providers, Combos & Gateway, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 3: `combos` Module (Fallback Combos, Aliases, Disabled Models, 4-Tier Resolver & Cooldown State Machine) & OpenAPI Fragment, Task 5: `gateway` HTTP Pipeline — Streaming SSE Executor, Auth, Failover, OAuth Refresh, Response Chains & `GatewayInterceptor` Chain, Task 7: Frontend Providers Control Studio (`web/src/features/providers/`, `/dashboard/providers/`, `/dashboard/providers/detail/`), Task 8: Frontend Models, Combos & Endpoint Studio (`web/src/features/combos/`, `/dashboard/models/`, Navigation Promotion & Full Single-Binary Verification), Task 2: `telegram` Module (Vaulted Config, Fail-Closed Webhook, Command Controller, EventBus Alerts & OpenAPI) (+7 more)

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.22
Nodes (18): errorDisposition, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus(), isClientPayloadErrorBody(), applyCodeBuddyHeaders(), cloneShallowMap(), codebuddyDomainFor() (+10 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (13): fakeCredentialStore, stubAdminCredentialStore, trafficShaper, stubSettings, acquireSem(), sleepContext(), extractUsageCost(), TestSystemAdminCredentialStore() (+5 more)

### Community 194 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 195 - "redact.go"
Cohesion: 0.18
Nodes (12): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), newEventID(), normalizeSeverity() (+4 more)

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "quotaCache"
Cohesion: 0.27
Nodes (7): cloneConnectionQuota(), newQuotaCache(), ProviderConnectionQuota, ProviderQuotasResponse, quotaCache, quotaCacheEntry, quotaFlightCall

### Community 198 - "UpstreamCredentials"
Cohesion: 0.09
Nodes (11): catalogFakeProviderStore, fakeProviderStore, staticProviderStore, stubProviderStore, stubProviderStore, ProviderConnection, TestContractInterfacesSatisfied(), newFakeProviderStore() (+3 more)

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 201 - "repository"
Cohesion: 0.15
Nodes (6): newRepository(), Alias, AliasListResponse, Config, Event, repository

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 203 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "Module"
Cohesion: 0.20
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 209 - ".ProvisionKey"
Cohesion: 0.25
Nodes (9): Review Focus, Task 1: `contracts.KeyProvisioner` & `TeamAssigner`, `keys.ProvisionKey`, and `teams` Module, Task 2: `vouchers` Module (Voucher Lifecycle, Atomic Claim & Anti-Abuse), Task 3: Wire Wave 3A in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, stubKeyProvisioner, KeyProvisioner, ProvisionedKey, ProvisionKeyRequest (+1 more)

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & System, Operations & Observability, Routing & Gateway Core

### Community 211 - "setupCatalogTestModule"
Cohesion: 0.14
Nodes (11): setupCatalogTestModule(), TestComboTestHonorsPinnedConnectionAndComboRef(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), TestResolveAttachesComboRuntimeMetadata(), BenchmarkResolveStrategies() (+3 more)

### Community 212 - "ProviderEntry"
Cohesion: 0.22
Nodes (11): ProviderEntry, CustomEndpointPreset, PROVIDER_REGISTRY, PROVIDER_REGISTRY_MAP, ProviderAuthMode, ProviderModelSpec, ProviderRegistrySpec, ProviderSourceCategory (+3 more)

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "newAuthHarness"
Cohesion: 0.09
Nodes (25): Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness(), newFakeStore(), newTestSigner() (+17 more)

### Community 216 - "settingDefinition"
Cohesion: 0.20
Nodes (7): SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), settingDefinition, settingsListResp, validatedUpdate

### Community 217 - "ComboSnapshot"
Cohesion: 0.50
Nodes (3): stubCombos, AliasSnapshot, ComboSnapshot

### Community 218 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

### Community 219 - "net/http.Header"
Cohesion: 0.12
Nodes (9): memoryResponseWriter, extractVendorHeaders(), newMemoryResponseWriter(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), SessionStore, syncBuffer (+1 more)

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 221 - "scanCombo"
Cohesion: 0.36
Nodes (5): Combo, rowScanner, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 222 - "ChatCompleter"
Cohesion: 0.36
Nodes (7): Global Constraints, Review Focus, Task 1: Cross-Module `ChatCompleter` Contract & `gateway.Module.CompleteChat`, Task 3: `arena` Module (Multi-Model Comparison Battles, Leaderboard, Persistent Basic Chat Sessions & OpenAPI), Wave 3D — Telegram Bot Controller, Model Arena & Basic Chat Studio Implementation Plan, ChatCompleter, TestGatewayCompleteChat()

### Community 223 - "DB"
Cohesion: 0.05
Nodes (78): Options, authHarness, Options, Options, testHarness, File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `security` Module (Security Posture & Redacted Audit Trail) (+70 more)

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 229 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 230 - "New"
Cohesion: 0.29
Nodes (5): New(), boolToInt(), newRepository(), Config, repository

### Community 235 - "optimization.test.ts"
Cohesion: 0.25
Nodes (7): consoleSource, customModelModalSource, detailSource, hooksSource, logoSource, modalsSource, providersDir

### Community 238 - "repository"
Cohesion: 0.33
Nodes (4): ImportRun, repository, boolToInt(), newRepository()

### Community 241 - "botAPIClient"
Cohesion: 0.43
Nodes (3): newBotAPIClient(), sanitizeBotError(), botAPIClient

### Community 242 - "compressor.go"
Cohesion: 0.10
Nodes (36): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressRoot(), compressText(), estimateTokens(), hasAnthropicToolResult() (+28 more)

### Community 246 - "translate_responses.go"
Cohesion: 0.14
Nodes (26): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+18 more)

### Community 248 - "headroom.go"
Cohesion: 0.33
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 249 - "newOAuthLoopbackProxy"
Cohesion: 0.50
Nodes (4): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider()

### Community 250 - "Review Focus"
Cohesion: 0.67
Nodes (3): Review Focus, Task 1: `mitm` Module (Root CA Engine, Loopback TLS MITM Proxy, IDE Tool Hosts & Model Aliases), Task 2: `pxpipe` Module (In-Process Transform Pipeline, `GatewayInterceptor`, Health Self-Test & Logs)

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1082 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **27 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `AddApiKeyModal` to `ProviderDetailView.tsx`, `UpstreamCredentials`, `net/http.ResponseWriter`?**
  _High betweenness centrality (0.205) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `SystemOverview.tsx` be split into smaller, more focused modules?**
  _Cohesion score 0.05328005328005328 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `UpstreamCredentials` to `context.Context`, `Module`, `net/http.Request`, `go_pkg_encoding_json`, `oauth_registry.go`, `net/http.ResponseWriter`, `AddApiKeyModal`, `deviceOAuthSession`?**
  _High betweenness centrality (0.173) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `AddApiKeyModal()` connect `AddApiKeyModal` to `ProvidersModals.tsx`, `Label`, `ProviderDetailView.tsx`, `Badge`, `SystemOverview.tsx`?**
  _High betweenness centrality (0.074) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._