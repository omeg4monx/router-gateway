# Graph Report - omeg4mon-router  (2026-11-28)

## Corpus Check
- 512 files · ~561,782 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 11 file(s) not represented in the graph (top: (none) 7, .lock 1, .ico 1)

## Summary
- 5855 nodes · 19573 edges · 248 communities (210 shown, 38 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 1833 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2c1432fe`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- combos/module.go
- extractErrorMessage
- biome.json
- settings
- Rewrite Architecture Specification
- APIKeySummary
- modeleditor/engine.go
- Open
- TypeScript Strict Compiler Config
- combos/schemas.ts
- useUIStore
- testing.T
- shadcn/ui Component Configuration
- TestContractInterfacesSatisfied
- client.ts
- App.tsx
- oauth_registry.go
- Design System & Wave 0 Plan Docs
- client.ts
- updates
- KeysConsole.tsx
- DB
- with
- jobs
- Module
- steps
- web-check
- on
- extractErrorMessage
- TokenSaverConsole.tsx
- ProviderDetailView.tsx
- SystemOverview.tsx
- docker-build
- go-check
- NotificationsConsole.tsx
- ModelEditorConsole.tsx
- Go Module Root Definition
- CanonicalUsage
- TelegramConsole.tsx
- QuotaConsole.tsx
- CanonicalStreamChunk
- main.go
- ir.go
- translate_gemini.go
- go_pkg_strings
- repository
- PluginsStudio.tsx
- Module
- ModelDescriptor
- setupModelEditorTest
- Module
- translate_responses.go
- setupNotificationsTest
- NewClient
- AddCustomModelModal.tsx
- Badge
- web/package.json
- translate_anthropic.go
- migrator.go
- Button
- oauthLoopbackProxy
- LiveMonitorConsole.tsx
- providers/repo.go
- newTestHarness
- teeHandler
- ProxyPoolsConsole.tsx
- TestContractInterfacesSatisfied
- navigation.ts
- pricingStore
- oauthLoopbackProxy
- ResolvedTarget
- repo
- Module
- net/http.Request
- go_pkg_encoding_json
- arena/schemas.ts
- go_pkg_context
- telemetry.go
- readLegacySQLite
- repository
- PWA Support Design — `omeg4mon-router`
- Module
- KeyInfo
- paths
- PromptGuardConsole.tsx
- Module
- Module
- settingDefinition
- modelIntelligence.ts
- redact.go
- Module
- repository
- keyLimiter
- Module
- CanonicalChatResponse
- time.Duration
- .resolveIdentifier
- extractErrorMessage
- net/http.ResponseWriter
- time.Time
- Module
- Module
- evaluateAntigravityQuota
- post
- .Fetch
- settingDefinition
- New
- ProviderTopology.tsx
- Load
- repository
- repo
- properties
- Module
- net/url.URL
- Module
- RouteMounts
- .ListModels
- modeleditor/hooks.ts
- plugins/repo.go
- properties
- fakeHandle
- translate_openai.go
- setupTokenSaverTest
- io/fs.FS
- strategies.go
- OMEG4MON Router
- DB
- scanCombo
- certManager
- repo
- .identifyProviderModel
- Review Focus
- properties
- Module
- release-binaries
- steps
- applyUpstreamHeaders
- Provider Quota Menu & Token Harbor Quota Integration
- TopProgressBar.tsx
- ref_bun_test
- `/dashboard/models` OmniRoute Parity & `/dashboard/combos` Route Split
- transform.go
- fetch.go
- properties
- 200
- newTestHarness
- properties
- dependencies
- cli/package.json
- setupCatalogTestModule
- setupTestMitm
- get
- scripts
- post
- schemas
- manifest.test.ts
- PreviewPluginsResponse
- repository
- net/http.Handler
- ImporterHistoryListResponse
- PluginListResponse
- PreviewPluginsRequest
- Review Focus
- Review Focus
- routingStateCache
- schemas
- properties
- RunLegacyImportResponse
- omeg4mon-router Versioning Verificator
- version
- stubDep
- setupMediaTest
- parameters
- repository
- New
- resolve.js
- ClearImporterHistoryResponse
- source_path
- .Snapshot
- playgroundUtils.test.ts
- docker-release
- palette.test.ts
- `@omeg4mon/router` CLI Launcher
- SessionSigner
- .ExportCombos
- go_pkg_encoding_json
- context.Context
- setupTunnelTest
- repository
- parameters
- ResolvedTarget
- PromptTemplateCombobox.tsx
- omeg4mon-router Verification Gate
- repo
- repo
- single-binary
- stubModule
- verify-web
- headroom.go
- smoke-test.sh
- Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)
- NewStreamTranslator
- stubSystemModule
- Feature Suite
- stubDiscoverySystemModule
- Post-Login Dashboard Layout & Categorized Searchable Sidebar Design
- net/http/httptest.ResponseRecorder
- stubSystem
- net/http.Header
- Provider Setup Modals Redesign & Pre-Save Connectivity Test
- stubModule
- stubSystemModule
- DB
- stubDep
- optimization.test.ts
- AttestationManager
- stubDepModule
- Implementation Plan: `/dashboard/models` OmniRoute Parity
- stubModule
- time.Duration
- stubSystemModule
- omeg4mon-router Project Standards
- jsonMarshalerError
- .StatusCode
- compressor.go
- Status

## God Nodes (most connected - your core abstractions)
1. `extractErrorMessage()` - 194 edges
2. `WriteJSON()` - 187 edges
3. `WriteError()` - 184 edges
4. `Button` - 106 edges
5. `cn()` - 97 edges
6. `DB` - 93 edges
7. `Module` - 79 edges
8. `buildApp()` - 79 edges
9. `Badge()` - 76 edges
10. `Bus` - 75 edges

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
- 2-file cycle: `internal/modules/combos/module_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/module_test.go`
- 2-file cycle: `internal/modules/combos/catalog_test.go -> internal/modules/combos/strategies_test.go -> internal/modules/combos/catalog_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/module_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/antigravity_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/antigravity_test.go`
- 2-file cycle: `internal/modules/gateway/codebuddy_test.go -> internal/modules/gateway/freepi_test.go -> internal/modules/gateway/codebuddy_test.go`
- 2-file cycle: `internal/modules/gateway/module_test.go -> internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/module_test.go`
- 2-file cycle: `internal/modules/gateway/optimization_e2e_test.go -> internal/modules/gateway/venice_test.go -> internal/modules/gateway/optimization_e2e_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/device_oauth_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/aws_quick_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/aws_quick_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/module_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/device_oauth_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/device_oauth_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/oauth_flow_integration_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/provider_hardening_test.go -> internal/modules/providers/module_test.go`
- 2-file cycle: `internal/modules/providers/module_test.go -> internal/modules/providers/xrouter_providers_test.go -> internal/modules/providers/module_test.go`

## Hyperedges (group relationships)
- **OMEG4MON ROUTER Visual Design System** — docs_design_brand_identity, docs_design_color_palette, docs_design_typography, docs_design_solid_surfaces_and_dials [EXTRACTED 1.00]
- **Zero-Stall Gateway Request & Streaming Pipeline** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_gateway_sse_pipeline, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_canonical_ir_translators, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_model_resolution_precedence, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_usage_metadata_retention [EXTRACTED 1.00]
- **Modular Monolith Boundary & Contract Stack** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_modular_monolith_registry, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_internal_contracts, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_eventbus_pubsub, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_sqlite_wal_split_pools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_openapi_contract_codegen [EXTRACTED 1.00]
- **Contract-First Single-Binary Build Pipeline** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_system_module, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_web_shell_and_codegen, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_single_binary_pipeline [EXTRACTED 1.00]
- **Wave 0 Go Platform Foundation Infrastructure** — docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_config, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_logging_and_eventbus, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_crypto, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_platform_db, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpclient_and_jobs, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_contracts_and_registry, docs_superpowers_plans_2026_11_22_wave0_platform_foundation_httpserver [EXTRACTED 1.00]
- **Waves 0-3 Phased Modular Delivery Roadmap** — docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_0_platform_foundation, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_1_core_gateway, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_2_routing_security_devtools, docs_superpowers_specs_2026_11_22_modular_ts_go_rewrite_design_wave_3_extended_ecosystem [EXTRACTED 1.00]

## Communities (248 total, 38 thin omitted)

### Community 0 - "combos/module.go"
Cohesion: 0.08
Nodes (19): changePasswordRequest, changePasswordResponse, loginRequest, loginResponse, logoutResponse, meResponse, policyResponse, updatePolicyRequest (+11 more)

### Community 1 - "extractErrorMessage"
Cohesion: 0.23
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
Cohesion: 0.11
Nodes (17): cloneStringSlice(), decodeStringSlice(), encodeStringSlice(), APIKeySummary, newRepository(), scanSummary(), createParams, deleteResult (+9 more)

### Community 6 - "modeleditor/engine.go"
Cohesion: 0.07
Nodes (39): antigravityPromptRewriteRule, applyOverrideToJSON(), applyPipeline(), bypassRuleMatches(), cloneJSONValue(), compileGlobPattern(), globToRegexp(), hasBlockText() (+31 more)

### Community 7 - "Open"
Cohesion: 0.05
Nodes (59): TestAppWave1BGatewayEndToEnd(), roundTripFunc, Task 6: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Gateway Tests, Task 2: `httpclient` Provider-Context Proxy Hook & `proxypools` Module (HTTP/SOCKS5 Pools, Latency Tester, Auto-Scanner, Provider Binding), `/dashboard` System Overview UI/UX Redesign (Approach A), Intent & Why, Key Proportion & UX Fixes, TestCombosRoutingStateCacheHitsAndInvalidates() (+51 more)

### Community 8 - "TypeScript Strict Compiler Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, allowJs, esModuleInterop, incremental, isolatedModules, jsx, lib (+11 more)

### Community 9 - "combos/schemas.ts"
Cohesion: 0.06
Nodes (44): combosQueryKeys, UpdateComboVariables, UpsertModelAliasVariables, useClearModelCooldownMutation(), useDeleteModelAliasMutation(), useDisabledModels(), useDisableModelMutation(), useEnableModelMutation() (+36 more)

### Community 10 - "useUIStore"
Cohesion: 0.14
Nodes (26): Task 4: Frontend PWA Lifecycle Hook, Install Button, Offline/Update Banner & Full Verification, Functional & UX Scope, useLogoutMutation(), useSystemHealth(), useSystemVersion(), PwaInstallButton(), PwaInstallButtonProps, PwaStatusBanner() (+18 more)

### Community 11 - "testing.T"
Cohesion: 0.03
Nodes (100): TestSentinelErrorsDistinct(), decodeJSON(), newTestHarness(), TestCombosAliasesDisabledCRUDAndAuth(), TestCooldownConcurrency(), TestDisabledModelsAndCooldownStateMachine(), TestListModelsDeduplicationAndFiltering(), TestProviderQualifiedModelsListAndBareResolveSymmetry() (+92 more)

### Community 12 - "shadcn/ui Component Configuration"
Cohesion: 0.11
Nodes (17): aliases, components, hooks, lib, ui, utils, iconLibrary, rsc (+9 more)

### Community 13 - "TestContractInterfacesSatisfied"
Cohesion: 0.07
Nodes (6): fakeCredentialStore, AdminCredentialStore, stubAdminCredentialStore, TestContractInterfacesSatisfied(), stubCredentialStore, repository

### Community 14 - "client.ts"
Cohesion: 0.18
Nodes (21): Task 5: Web Studio UI for Security & Notifications/Jobs, Navigation, System Overview & Single-Binary Verification, securityQueryKeys, useClearSecurityAuditMutation(), useCreateSecurityAuditMutation(), useSecurityAuditEvents(), useSecurityPosture(), AUDIT_PAGE_LIMIT, AUDIT_SEVERITIES (+13 more)

### Community 15 - "App.tsx"
Cohesion: 0.09
Nodes (33): App(), ArenaConsole, ChatConsole, CombosConsole, DashboardLayout(), HomePage, ImporterStudio, InitPage (+25 more)

### Community 16 - "oauth_registry.go"
Cohesion: 0.08
Nodes (42): isBlockedAntigravityProjectID(), resolveAwsQuickTenantBaseURL(), codebuddySpecFor(), newCodeBuddyOAuthStrategy(), refreshCodeBuddyToken(), TestCodeBuddyTokenRefresh(), originOfURL(), postExchangeIFlow() (+34 more)

### Community 17 - "Design System & Wave 0 Plan Docs"
Cohesion: 0.19
Nodes (13): Brand Identity (Caduceus & Ember Line), OMEG4MON ROUTER Design Direction, Cross-Module Contracts & Topological Module Registry, Streaming HTTP/2 Client Pool & Panic-Safe Job Runner, Chi HTTP Server, Streaming Timeout Split & SPA Static Handler, Structured Logging RingBuffer & Non-Blocking EventBus, Platform Config Loader (internal/platform/config), Platform Crypto Vault, Bcrypt & Session Signer (+5 more)

### Community 18 - "client.ts"
Cohesion: 0.05
Nodes (66): Task 7: Frontend API Keys & Settings Control Plane Pages + Full Single-Binary Verification, @hookform/resolvers, @tanstack/react-query, zod, zod@^4.6.7, authKeys, AuthMeResponse, LoginRequest (+58 more)

### Community 19 - "updates"
Cohesion: 0.13
Nodes (14): patterns, patterns, actions, go-dependencies, web-dependencies, interval, updates, directory (+6 more)

### Community 20 - "KeysConsole.tsx"
Cohesion: 0.05
Nodes (79): buildGatewayUrls(), copyTextToClipboard(), countKeyStatusBuckets(), filterKeysByStatusAndSearch(), GatewayRouteId, generateKeySnippet(), getKeyExpiryState(), isIpLiteral() (+71 more)

### Community 21 - "DB"
Cohesion: 0.04
Nodes (48): buildApp(), TestAppBootstrapsMigrationsRoutesAndShutdown(), TestAppCloseTeardownOrderIsLIFO(), TestAppRunAndShutdownEndpointTriggersGracefulExit(), TestAppWave1AAuthSettingsKeysEndToEnd(), TestAppWave1CUsageEndToEnd(), TestBuildAppCleansUpOnPartialFailure(), TestWave2AE2ETranslatorAndCatalog() (+40 more)

### Community 22 - "with"
Cohesion: 0.12
Nodes (16): with, bun-version, cache, cache-from, cache-to, check-latest, context, go-version-file (+8 more)

### Community 23 - "jobs"
Cohesion: 0.15
Nodes (12): concurrency, cancel-in-progress, group, name, needs, runs-on, timeout-minutes, jobs (+4 more)

### Community 24 - "Module"
Cohesion: 0.12
Nodes (5): decodeJSONBody(), isValidMode(), isValidSensitivity(), validateRule(), Module

### Community 25 - "steps"
Cohesion: 0.42
Nodes (9): steps, steps, steps, steps, name, run, uses, working-directory (+1 more)

### Community 26 - "web-check"
Cohesion: 0.29
Nodes (7): run, web-check, working-directory, defaults, name, runs-on, timeout-minutes

### Community 27 - "on"
Cohesion: 0.33
Nodes (6): on, pull_request, push, workflow_dispatch, branches, branches

### Community 28 - "extractErrorMessage"
Cohesion: 0.11
Nodes (35): mitmQueryKeys, useClearMitmEventsMutation(), useCreateMitmAliasMutation(), useDeleteMitmAliasMutation(), useGenerateMitmCAMutation(), useMitmAliases(), useMitmEvents(), useMitmStatus() (+27 more)

### Community 29 - "TokenSaverConsole.tsx"
Cohesion: 0.11
Nodes (29): tokenSaverQueryKeys, useClearHeadroomCacheMutation(), useResetTokenSaverStatsMutation(), useTokenSaverConfig(), useTokenSaverStats(), useUpdateTokenSaverConfigMutation(), HISTORY_LIMIT_PRESETS, HistoryPresetId (+21 more)

### Community 30 - "ProviderDetailView.tsx"
Cohesion: 0.02
Nodes (196): Provider Setup Modals Redesign & Pre-Save Test Implementation Plan, Review Focus, Task 2: OpenAPI Contract, Generated Schema, and React Query Hook, Task 3: Pure Modal Helpers (`modalUtils.ts`) + Unit Tests (TDD), Task 4: Redesign `ModalShell`, `AddCustomNodeModal`, `AddApiKeyModal`, and `OAuthWizardModal`, Task 5: Full Verification & Graphify Sync, runModelIdentityProbe(), ADAPTER_TO_REGISTRY_ID (+188 more)

### Community 31 - "SystemOverview.tsx"
Cohesion: 0.05
Nodes (88): Task 4: Frontend Usage, Telemetry & Pricing Studio (`web/src/features/usage/`, `/dashboard/usage/`), Task 4: Web Studio UI for Teams & Vouchers, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for Media Gateway & Tunnel Manager, Navigation, System Overview & Single-Binary Verification, Task 4: Web Studio UI for MITM Proxy & PXPIPE Pipeline, Navigation, System Overview & Single-Binary Verification, Task 5: Web Studios for Telegram Bot, Model Arena & Basic Chat, Navigation, System Overview & Single-Binary Verification, useSystemTelemetry(), formatBytesGB(), formatHostUptimeShort() (+80 more)

### Community 32 - "docker-build"
Cohesion: 0.40
Nodes (5): name, needs, runs-on, timeout-minutes, docker-build

### Community 33 - "go-check"
Cohesion: 0.50
Nodes (4): name, runs-on, timeout-minutes, go-check

### Community 34 - "NotificationsConsole.tsx"
Cohesion: 0.17
Nodes (23): notificationsQueryKeys, useBackgroundJobs(), useClearNotificationsMutation(), useCreateNotificationMutation(), useDeleteNotificationMutation(), useMarkAllNotificationsReadMutation(), useMarkNotificationReadMutation(), useNotificationsList() (+15 more)

### Community 35 - "ModelEditorConsole.tsx"
Cohesion: 0.14
Nodes (31): useModelOverrides(), usePreviewModelEditorMutation(), useRunQuickPromptTestMutation(), EMPTY_FORM_VALUES, FieldLabel(), FieldLabelProps, MobileTab, ModelEditorConsole() (+23 more)

### Community 38 - "CanonicalUsage"
Cohesion: 0.11
Nodes (19): anthropicStreamEncoder, CanonicalUsage, composedStreamTranslator, geminiPendingTool, geminiStreamEncoder, openAIStreamEncoder, responsesStreamBlockState, responsesStreamEncoder (+11 more)

### Community 39 - "TelegramConsole.tsx"
Cohesion: 0.13
Nodes (30): telegramQueryKeys, useClearTelegramEventsMutation(), useDeleteTelegramWebhookMutation(), useRegisterTelegramWebhookMutation(), useRotateTelegramSecretMutation(), useSendTelegramTestMutation(), useSimulateTelegramCommandMutation(), useTelegramConfig() (+22 more)

### Community 40 - "QuotaConsole.tsx"
Cohesion: 0.09
Nodes (52): Task 4: Quota UI Feature (`web/src/features/quota/`) (TDD), Frontend (`web/`), fetchProviderQuotas(), fetchSingleConnectionQuota(), QuotaLiveStatus, quotaQueryKeys, upsertConnectionQuotaInCache(), useProviderQuotas() (+44 more)

### Community 41 - "CanonicalStreamChunk"
Cohesion: 0.44
Nodes (3): CanonicalStreamChunk, openAIStreamDecoder, streamDecoder

### Community 42 - "main.go"
Cohesion: 0.08
Nodes (13): collectStaticHostInfo(), readWindowsRegistryString(), isLinkUnsupported(), writeKeyAtomically(), writeKeyExclusive(), healthResponse, initCompleteResponse, initRequest (+5 more)

### Community 43 - "ir.go"
Cohesion: 0.14
Nodes (33): Task 4: `gateway` Core — Canonical IR, Direct Fast-Path & Protocol Translators, CacheControl, CanonicalChatRequest, CanonicalMessage, ContentBlock, ImageSource, ResponseFormat, SystemBlock (+25 more)

### Community 44 - "translate_gemini.go"
Cohesion: 0.16
Nodes (21): geminiCandidate, geminiContent, geminiFileData, geminiFunctionCall, geminiFunctionCallingConfig, geminiFunctionDecl, geminiFunctionResponse, geminiGenConfig (+13 more)

### Community 45 - "go_pkg_strings"
Cohesion: 0.08
Nodes (28): CatalogModel, CatalogProviderGroup, ComboTargetTestResult, modelCatalogResponse, ProviderHealthMatrixEntry, providerHealthMatrixResponse, targetKey, testAllModelsRequest (+20 more)

### Community 46 - "repository"
Cohesion: 0.12
Nodes (12): escapeLikePattern(), newRepository(), LogsFilter, LogsResponse, ModelBreakdown, pricingListResponse, ProviderBreakdown, repository (+4 more)

### Community 47 - "PluginsStudio.tsx"
Cohesion: 0.16
Nodes (30): PluginListResponse, pluginsQueryKeys, PluginState, PreviewPluginsRequest, PreviewPluginsResponse, UpdatePluginRequest, usePluginsList(), usePreviewPluginsMutation() (+22 more)

### Community 48 - "Module"
Cohesion: 0.21
Nodes (6): validateModelPattern(), decodeJSONBody(), newID(), validateBypassRule(), validateOverride(), Module

### Community 49 - "ModelDescriptor"
Cohesion: 0.18
Nodes (12): Task 3: Build-Time Service Worker & Vite 8 Precache Plugin, vite, @vitejs/plugin-react, vite@^8.3.8, @vitejs/plugin-react@^6.3.0, buildServiceWorkerSource(), CORE_PRECACHE_URLS, omeg4monPwaPlugin() (+4 more)

### Community 50 - "setupModelEditorTest"
Cohesion: 0.10
Nodes (18): File Map, Task 1: Gateway Metadata Pass-Through & `promptinjection` Module (Prompt Guard Middleware), Task 2: `modeleditor` Module (Custom Model Parameter Overrides & Conditional Bypass Rules), Task 5: Web Studio UI for Prompt Guard & Model Editor / Bypass, Task 6: Web Studio UI for Token Saver & Headroom, Navigation, System Overview & Single-Binary Verification, doJSON(), setupModelEditorTest(), TestConcurrentExtraParamsNoDataRaceOrCacheMutation() (+10 more)

### Community 51 - "Module"
Cohesion: 0.14
Nodes (6): decodeStrictJSON(), Module, isValidSeverity(), toJobDTO(), JobStatusDTO, Notification

### Community 52 - "translate_responses.go"
Cohesion: 0.13
Nodes (28): openAIChatFunctionSpec, openAIChatTool, responsesContentPart, responsesFormat, responsesIncompleteDetails, responsesInputTokensDetails, responsesItem, responsesReason (+20 more)

### Community 53 - "setupNotificationsTest"
Cohesion: 0.47
Nodes (5): Task 3: `notifications` Module (System Notifications Inbox, SSE Stream & `/api/jobs` Control), setupNotificationsTest(), TestNotificationsCRUDReadUnreadAndPrune(), TestNotificationsEventBusAndSSEStreamCleanup(), TestNotificationsJobsListAndTriggerConflict()

### Community 54 - "NewClient"
Cohesion: 0.07
Nodes (24): arenaTestEnv, stubCompleter, roundTripFunc, Options, doJSON(), setupArenaTest(), TestArenaBattlePartialFailureVoteAndLeaderboard(), TestBasicChatSessionsHistoryAndValidation() (+16 more)

### Community 55 - "AddCustomModelModal.tsx"
Cohesion: 0.26
Nodes (14): classifyTestError(), comboTestResultKey(), flattenCatalogResponse(), modelTestResultKey(), normalizeTestStatus(), RawCatalogGroup, SLOW_LATENCY_MS, ModelIdentityResult (+6 more)

### Community 56 - "Badge"
Cohesion: 0.07
Nodes (60): lucide-react, ComboCard, ComboCardProps, buildProviderAccountGroups(), formatAuthTypeBadge(), formatStepAccountLabel(), getModelsForSelection(), isDuplicateModelStep() (+52 more)

### Community 57 - "web/package.json"
Cohesion: 0.12
Nodes (23): @biomejs/biome, openapi-typescript, tailwindcss, @tailwindcss/vite, @types/bun, @types/node, @types/react, @types/react-dom (+15 more)

### Community 58 - "translate_anthropic.go"
Cohesion: 0.21
Nodes (16): anthropicChoice, anthropicContentBlock, anthropicImageSource, anthropicMessage, anthropicRequest, anthropicResponse, anthropicStreamDecoder, anthropicThinking (+8 more)

### Community 59 - "migrator.go"
Cohesion: 0.10
Nodes (30): ModelOption, PolicyChangedEvent, adminSubjectKey, ErrorDetail, ErrorEnvelope, inspectRequest, runImportRequest, Module (+22 more)

### Community 60 - "Button"
Cohesion: 0.07
Nodes (95): Task 6: Frontend Auth Flow (`web/src/features/auth/`, Live `/login/`, `AuthGuard`, Sign-Out), class-variance-authority, react, react-hook-form, react-router, ArenaConsole(), ChatConsole(), STARTER_PROMPTS (+87 more)

### Community 61 - "oauthLoopbackProxy"
Cohesion: 0.12
Nodes (9): newAwsQuickRuntime(), parseAwsQuickPastedQuery(), validateAwsRegion(), awsQuickDiscovery, awsQuickResult, awsQuickRuntime, awsQuickSession, awsQuickTokens (+1 more)

### Community 62 - "LiveMonitorConsole.tsx"
Cohesion: 0.07
Nodes (57): attrValueClass(), levelColorClass(), levelSymbol(), LiveMonitorConsole(), MonitorTab, normalizeUsagePayload(), StreamState, SYSTEM_LEVELS (+49 more)

### Community 63 - "providers/repo.go"
Cohesion: 0.11
Nodes (22): cloneStringSlice(), cloneStringMap(), isCodeBuddyHostFor(), isTokenHarborHost(), decodeStringMap(), decodeStringSlice(), encodeStringMap(), encodeStringSlice() (+14 more)

### Community 64 - "newTestHarness"
Cohesion: 0.15
Nodes (18): File Map, Task 1: `usage` Module (Async Batched Writer, Pricing, Logs, Summary, Live SSE & Retention Pruning) & OpenAPI Fragment, Task 2: `clitools` Module (Stateless Config Generators & Loopback-Only Atomic Local Apply/Rollback) & OpenAPI Fragment, Task 3: Binary Wiring (`cmd/omeg4mon-router/main.go`), Root OpenAPI Aggregation, `make codegen` & Live E2E Tests, Task 5: Frontend CLI Tools Studio (`web/src/features/clitools/`, `/dashboard/cli-tools/`), Navigation Promotion & Final Wave 1 Single-Binary Verification, New(), decodeJSON(), newTestHarness() (+10 more)

### Community 65 - "teeHandler"
Cohesion: 0.17
Nodes (6): addAttrAt(), attrValue(), ensurePath(), parseLevel(), boundAttrs, teeHandler

### Community 66 - "ProxyPoolsConsole.tsx"
Cohesion: 0.11
Nodes (35): useProvidersCatalog(), proxyPoolsQueryKeys, useCreateProxyPoolMutation(), useDeleteProxyBindingMutation(), useDeleteProxyPoolMutation(), useProxyPools(), useScanProxyCandidatesMutation(), useTestProxyPoolMutation() (+27 more)

### Community 67 - "TestContractInterfacesSatisfied"
Cohesion: 0.09
Nodes (9): stubAdminAuth, AuthVerifier, stubAuthVerifier, stubAuthVerifier, stubAdminAuth, stubAdminAuth, stubAdminAuth, stubAdminAuth (+1 more)

### Community 68 - "navigation.ts"
Cohesion: 0.11
Nodes (26): Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan, Tasks, ParameterGuideEntry, QuickTestPreset, TargetModelComboboxProps, UILocale, ActiveNavContext, ALL_NAV_ITEMS (+18 more)

### Community 69 - "pricingStore"
Cohesion: 0.09
Nodes (23): with, build-args, bun-version, cache, cache-from, cache-to, check-latest, context (+15 more)

### Community 70 - "oauthLoopbackProxy"
Cohesion: 0.16
Nodes (8): newCodexOAuthRuntime(), TestCodexLoopbackRuntimeCaptureAndStatus(), newOAuthLoopbackProxy(), TestOAuthLoopbackProxyMultiProvider(), codexOAuthRuntime, loopbackProviderServer, loopbackSessionEntry, oauthLoopbackProxy

### Community 71 - "ResolvedTarget"
Cohesion: 0.20
Nodes (8): strategyState, weightedCandidate, applyComboStrategy(), applyProviderStrategy(), applyStrategy(), fastRandFloat64(), fastRandIntN(), newStrategyState()

### Community 72 - "repo"
Cohesion: 0.15
Nodes (6): Battle, ChatMessage, ChatSession, ContenderResult, LeaderboardEntry, repo

### Community 73 - "Module"
Cohesion: 0.14
Nodes (12): decodeStrictJSON(), extractGatewayAPIKey(), extractMediaSelector(), isModelAllowed(), isRetryableUpstreamStatus(), isValidCapability(), isValidDriver(), normalizeTargets() (+4 more)

### Community 74 - "net/http.Request"
Cohesion: 0.07
Nodes (21): clearSessionCookie(), isSecureRequest(), readSessionCookie(), setSessionCookie(), decodeJSONBody(), clientRateLimitKey(), buildModelProbeRequest(), buildIdentityProbeRequest() (+13 more)

### Community 75 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (27): buildDiscoveryCandidateURLs(), extractDiscoveredModelIsFree(), firstStringField(), mergeConnectionModels(), parseDiscoveredModels(), parseOptionalPrice(), TestAutoDetectModelsDiscoveryCandidates(), TestParseDiscoveredModelsExtractsIsFreeFromPricingAndFlags() (+19 more)

### Community 76 - "arena/schemas.ts"
Cohesion: 0.08
Nodes (26): ArenaBattleFormInput, ArenaBattleFormOutput, arenaBattleFormSchema, ArenaBattleListResponse, ArenaContenderResult, ArenaLeaderboardEntry, ArenaLeaderboardResponse, ArenaModelListResponse (+18 more)

### Community 77 - "go_pkg_context"
Cohesion: 0.19
Nodes (7): mapBusEventToNotification(), mapBusEventToAudit(), payloadToMap(), EmitPayload, PoolTestResponse, AuditInput, storedSetting

### Community 78 - "telemetry.go"
Cohesion: 0.13
Nodes (24): clampPercent(), computeCPUPercent(), collectDynamicHostTelemetry(), collectStaticHostInfo(), estimateDarwinUsedMem(), runSysctl(), collectDynamicHostTelemetry(), collectStaticHostInfo() (+16 more)

### Community 79 - "readLegacySQLite"
Cohesion: 0.39
Nodes (7): legacyExtract, legacyKeyRecord, listTables(), readLegacyCombos(), readLegacyKeys(), readLegacySettings(), readLegacySQLite()

### Community 80 - "repository"
Cohesion: 0.18
Nodes (7): builtinRules(), boolToInt(), defaultConfig(), Config, GuardEvent, repository, Rule

### Community 81 - "PWA Support Design — `omeg4mon-router`"
Cohesion: 0.14
Nodes (12): Global Constraints, PWA Support Implementation Plan, Review Focus, Task 1: Go Static UI Headers for PWA (`sw.js` & `manifest.webmanifest`), Task 2: PWA Icons, Web App Manifest, HTML Meta & Dynamic `theme-color`, Approaches Considered, Builder Decisions & Deliverables, Constraints (+4 more)

### Community 82 - "Module"
Cohesion: 0.16
Nodes (10): extractUsageCost(), actorFromContext(), decodeJSONBody(), normalizeExpiresAt(), validateName(), validateStringSlice(), generateKeyID(), generateToken() (+2 more)

### Community 83 - "KeyInfo"
Cohesion: 0.13
Nodes (12): stubKeyValidator, stubKeyValidator, KeyInfo, extractAPIKey(), cloneCachedEntry(), cloneKeyInfo(), newKeyCache(), hashToken() (+4 more)

### Community 84 - "paths"
Cohesion: 0.22
Nodes (9): get, operationId, responses, summary, tags, paths, /api/plugins, /api/plugins/{key} (+1 more)

### Community 85 - "PromptGuardConsole.tsx"
Cohesion: 0.12
Nodes (31): promptGuardQueryKeys, useClearPromptGuardEventsMutation(), useCreatePromptGuardRuleMutation(), useDeletePromptGuardRuleMutation(), usePromptGuardConfig(), usePromptGuardEvents(), usePromptGuardRules(), useTestPromptGuardMutation() (+23 more)

### Community 86 - "Module"
Cohesion: 0.16
Nodes (5): clientIP(), decodeStrictJSON(), Module, isValidSeverity(), formatUTC()

### Community 87 - "Module"
Cohesion: 0.04
Nodes (45): Task 1: Expand Built-in `providers` Adapter Catalog (Qwen, GitHub Copilot, Kiro, Cursor, Antigravity), Global Constraints, Task 1: Backend Stateless Probe Endpoint (`POST /api/providers/nodes/test`), Approaches Considered, UpstreamCredentials, allBuiltInAdapters(), buildProtocolTestRequest(), cloneOAuthConfig() (+37 more)

### Community 88 - "settingDefinition"
Cohesion: 0.15
Nodes (8): truncateText(), decodeJSON(), formatBusAlert(), generateWebhookSecret(), Module, maskToken(), safeEqual(), ConfigView

### Community 89 - "modelIntelligence.ts"
Cohesion: 0.16
Nodes (14): detailSource, hooksSource, providersDir, CURATED_MODEL_RULES, CuratedRule, estimateHeuristicProfile(), getModelIntelligenceProfile(), MODEL_TIER_BADGE (+6 more)

### Community 90 - "redact.go"
Cohesion: 0.18
Nodes (12): isSensitiveKey(), RedactMetadata(), RedactText(), redactValue(), truncateRunes(), buildWhere(), newEventID(), normalizeSeverity() (+4 more)

### Community 91 - "Module"
Cohesion: 0.11
Nodes (9): stubUsageRecorder, stubUsageRecorder, stubUsageRecorder, UsageRecord, newBatcher(), recordingUsage, batcher, ModelPricing (+1 more)

### Community 92 - "repository"
Cohesion: 0.26
Nodes (3): DisabledModel, ModelAlias, repository

### Community 93 - "keyLimiter"
Cohesion: 0.31
Nodes (3): keyLimiter, keyLimiterShard, keyLimitSnapshot

### Community 95 - "CanonicalChatResponse"
Cohesion: 0.17
Nodes (19): CanonicalChatResponse, CanFastPath(), DecodeRequest(), DecodeResponse(), EncodeRequest(), EncodeResponse(), IsSupportedProtocol(), rawJSONToString() (+11 more)

### Community 96 - "time.Duration"
Cohesion: 0.10
Nodes (14): failureRecord, stubHealthTracker, trafficShaper, computeBackoff(), acquireSem(), sleepContext(), cloneConnectionQuota(), newQuotaCache() (+6 more)

### Community 97 - ".resolveIdentifier"
Cohesion: 0.12
Nodes (11): resolveSnapshot, applyComboRuntimeMetadata(), buildResolvedTarget(), configIntValue(), connectionSupportsBare(), connectionSupportsExplicit(), isNestedExecuteMode(), withMetadataEntry() (+3 more)

### Community 98 - "extractErrorMessage"
Cohesion: 0.06
Nodes (72): ImporterHistoryListResponse, importerQueryKeys, InspectLegacyRequest, InspectLegacyResponse, RunLegacyImportRequest, RunLegacyImportResponse, useClearImporterHistoryMutation(), useImporterHistory() (+64 more)

### Community 99 - "net/http.ResponseWriter"
Cohesion: 0.05
Nodes (33): Module, Module, Task 1: Backend Go Endpoints & Tests (`internal/modules/combos/`), trackedResponseWriter, decodeJSON(), deriveSessionTitle(), nowTimestamp(), randomHexID() (+25 more)

### Community 101 - "time.Time"
Cohesion: 0.06
Nodes (20): loginAttemptWindow, loginRateLimiter, CooldownEntry, cooldownState, cooldownTracker, fakeClock, listCooldownsResponse, stubAccountHealthTracker (+12 more)

### Community 102 - "Module"
Cohesion: 0.09
Nodes (17): cachedGatewayResponse, inflightCall, Module, responseCoalescer, anthropicErrorType(), geminiErrorStatus(), openAIErrorType(), writeGatewayError() (+9 more)

### Community 105 - "evaluateAntigravityQuota"
Cohesion: 0.18
Nodes (14): agQuotaBucket, agQuotaCacheEntry, agQuotaGroup, agQuotaSummary, antigravityQuotaMonitor, agGroupMatchesFamily(), agIsSessionBucket(), agIsWeeklyBucket() (+6 more)

### Community 106 - "post"
Cohesion: 0.13
Nodes (19): content, content, description, post, post, schema, components, application/json (+11 more)

### Community 107 - ".Fetch"
Cohesion: 0.11
Nodes (23): TestCodeBuddyQuotaExcludesExpiredBonusAndZeroCapacityIsExhausted(), TestCodeBuddyQuotaFetcherRefillAndBonusPacks(), clampPercent(), cloneOptionalString(), codebuddyRefillCadence(), computeBucketPercentages(), formatCodeBuddyTimestamp(), mapCodeBuddyQuota() (+15 more)

### Community 108 - "settingDefinition"
Cohesion: 0.15
Nodes (11): HistoryItem, SettingItem, int64Ptr(), isStrictJSONInteger(), TestSettingDefinitionValidateRawIntegerBounds(), historyListResp, historyResponse, settingDefinition (+3 more)

### Community 109 - "New"
Cohesion: 0.15
Nodes (18): apiTimeoutMiddleware(), New(), generateSelfSignedCertFiles(), TestAdminAPIAuthEnforcement(), TestAdminAuthEnforcedOverHTTPSOnLoopback(), TestDevUIHandlerFallback(), TestGatewaySSEStreamsPastAPITimeout(), testLogger() (+10 more)

### Community 110 - "ProviderTopology.tsx"
Cohesion: 0.13
Nodes (21): @xyflow/react, @xyflow/react@^13.0.3, EDGE_TYPES, EnergyTopologyEdge(), FIT_OPTIONS, FLOW_PULSES, NODE_TYPES, prefersReducedMotion() (+13 more)

### Community 111 - "Load"
Cohesion: 0.29
Nodes (8): CoreModules(), Load(), parseDisabledModules(), resolveDefaultDataDir(), TestLoadDefaults(), TestLoadEnvAndFlagOverrides(), TestLoadRejectsInvalidPortAndCoreModuleDisable(), TestLoadTLSConfigAndValidation()

### Community 112 - "repository"
Cohesion: 0.32
Nodes (5): boolToInt(), defaultConfig(), newRepository(), Config, repository

### Community 113 - "repo"
Cohesion: 0.18
Nodes (11): buildUpstreamMediaURL(), isUniqueConstraintErr(), newID(), nowUTC(), scanCombo(), scanProvider(), Combo, providerRecord (+3 more)

### Community 114 - "properties"
Cohesion: 0.11
Nodes (20): type, additionalProperties, properties, required, type, type, type, combos_count (+12 more)

### Community 115 - "Module"
Cohesion: 0.16
Nodes (5): decodeJSONBody(), decodePathModelParam(), Module, isValidRate(), testHarness

### Community 116 - "net/url.URL"
Cohesion: 0.15
Nodes (11): stubProxyDialerResolver, stubProxyResolver, parseAndValidateProxyURL(), probeProxy(), proxyDialAddress(), redactRawCandidate(), runConcurrentScan(), dynamicProxyResolver (+3 more)

### Community 117 - "Module"
Cohesion: 0.15
Nodes (8): decodeStrictJSON(), validateLoopbackLocalURL(), validateProvider(), validateTailscaleMode(), BinariesStatus, BinaryStatus, Module, Status

### Community 118 - "RouteMounts"
Cohesion: 0.16
Nodes (4): Module, validateSQLiteSource(), decodeJSON(), RouteMounts

### Community 119 - ".ListModels"
Cohesion: 0.15
Nodes (7): stubModelLister, matchesAllowedModel(), stubModelLister, stubModelLister, TestListVirtualModels(), compositeModelLister, stubModelLister

### Community 120 - "modeleditor/hooks.ts"
Cohesion: 0.17
Nodes (15): invalidateModelEditorQueries(), modelEditorQueryKeys, QuickPromptTestInput, QuickPromptTestResult, useCreateModelOverrideMutation(), useDeleteModelOverrideMutation(), useUpdateModelOverrideMutation(), CreateModelOverrideFormValues (+7 more)

### Community 121 - "plugins/repo.go"
Cohesion: 0.06
Nodes (59): antiSlopDirective(), normalizeAntiSlopLevel(), buildSqueezeRecap(), estimateConversationTokens(), estimateMessageTokens(), estimateValueTokens(), shortenMessageContent(), squeezeConversation() (+51 more)

### Community 122 - "properties"
Cohesion: 0.18
Nodes (11): type, type, type, type, properties, category, description, key (+3 more)

### Community 123 - "fakeHandle"
Cohesion: 0.09
Nodes (7): streamResponseBody(), fakeHandle, fakeRunner, osProcessHandle, osProcessRunner, ProcessHandle, ssePipeReader

### Community 124 - "translate_openai.go"
Cohesion: 0.21
Nodes (17): openAIChatChoice, openAIChatContentPart, openAIChatFunctionCall, openAIChatMessage, openAIChatOutMessage, openAIChatRequest, openAIChatResponse, openAIChatToolCall (+9 more)

### Community 125 - "setupTokenSaverTest"
Cohesion: 0.23
Nodes (13): Task 3: `tokensaver` Module (Go-Native Token Compression & Headroom Cache-Aware Affinity), doJSON(), setupTokenSaverTest(), TestHeadroomLRUBoundedCapacityAndConcurrency(), TestHeadroomPrefixCacheHitMissAffinityAndAnthropicCacheControl(), TestResponsesAPIFunctionCallOutputRTKFilter(), TestTokenSaverCompressionPreservesCodeBlocksAndRespectsBypass(), TestTokenSaverFeatureTogglesCavemanPonytailAndTrimHistory() (+5 more)

### Community 126 - "io/fs.FS"
Cohesion: 0.11
Nodes (8): migrationFile, collectMigrations(), collectMigrationsFromDir(), parseMigrationVersion(), isMissingAssetPath(), openRegularFile(), resolveStaticFile(), serveFSFile()

### Community 127 - "strategies.go"
Cohesion: 0.22
Nodes (13): ComboStep, ProviderStrategyState, validateComboTargets(), encodeStringSlice(), isUniqueConstraintError(), cloneConfig(), cloneSteps(), decodeStepsJSON() (+5 more)

### Community 128 - "OMEG4MON Router"
Cohesion: 0.18
Nodes (11): 1. Build & Run Single Binary, 23 Vertical-Slice Modules, 2. Multi-Platform CLI Launcher (`@omeg4mon/router`), 3. Docker (Distroless Multi-Arch), Architecture, Configuration, Development & Verification, Gateway Wire API (+3 more)

### Community 129 - "DB"
Cohesion: 0.04
Nodes (78): Options, authHarness, Options, Options, testHarness, Task 2: `providers` Module (Pluggable Adapters, Vaulted Connections, Custom Nodes & Headless-Safe OAuth) & OpenAPI Fragment, Task 2: `security` Module (Security Posture & Redacted Audit Trail), Task 4: Wire Wave 2C in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests (+70 more)

### Community 130 - "scanCombo"
Cohesion: 0.24
Nodes (7): Combo, createComboParams, rowScanner, updateComboValues, decodeStringSlice(), scanCombo(), decodeConfigJSON()

### Community 131 - "certManager"
Cohesion: 0.17
Nodes (5): randomSerial(), writeAtomicFile(), cachedLeaf, CAInfo, certManager

### Community 132 - "repo"
Cohesion: 0.22
Nodes (8): boolToInt(), defaultStoredConfig(), newRepo(), nowTimestamp(), randomHexID(), Event, repo, storedConfig

### Community 133 - ".identifyProviderModel"
Cohesion: 0.05
Nodes (48): ModelIdentityResult, ModelTestResult, testAllModelsResponse, attemptOutcome, extractProbeTokenUsage(), extractUpstreamProbeError(), resolveKeylessProbeCredentials(), sanitizeProbeError() (+40 more)

### Community 134 - "Review Focus"
Cohesion: 0.14
Nodes (12): Global Constraints, Provider Quota Menu & Token Harbor Quota Implementation Plan, Review Focus, Task 1: Built-in Token Harbor Adapter & Backend Quota Engine (TDD), Task 2: OpenAPI Contract & TypeScript Schema Codegen, Task 3: Sidebar Navigation, Routing & Token Harbor Registry Metadata (TDD), Task 5: Full Verification, Web Build, Graph Update & Versioning Audit, components (+4 more)

### Community 135 - "properties"
Cohesion: 0.12
Nodes (17): type, type, type, properties, type, combos_imported, created_at, error (+9 more)

### Community 136 - "Module"
Cohesion: 0.20
Nodes (5): decodeStrictJSON(), matchesLevel(), normalizeLevelFilter(), parseIntParam(), Module

### Community 137 - "release-binaries"
Cohesion: 0.18
Nodes (10): jobs, release-binaries, name, permissions, contents, name, needs, permissions (+2 more)

### Community 138 - "steps"
Cohesion: 0.25
Nodes (11): steps, RELEASE_VERSION, steps, env, id, if, name, run (+3 more)

### Community 139 - "applyUpstreamHeaders"
Cohesion: 0.11
Nodes (17): Global Constraints, Provider Detail Proportional Model Grid & Per-Model Test Implementation Plan, Review Focus, Task 1: OpenAPI Contract & TypeScript Schema Codegen, Task 2: Backend Model Test Endpoint (TDD), Task 3: Frontend Hook, Proportional Model Grid & Per-Model Test UI (TDD), Task 4: Full Verification Gate & Graphify Update, Approaches Considered (+9 more)

### Community 140 - "Provider Quota Menu & Token Harbor Quota Integration"
Cohesion: 0.18
Nodes (9): Approaches Considered, Architecture & Contracts, Backend (`internal/modules/providers`), Constraints, Goals, Intent & Why, Non-Goals & Anti-Goals, Provider Quota Menu & Token Harbor Quota Integration (+1 more)

### Community 141 - "TopProgressBar.tsx"
Cohesion: 0.30
Nodes (10): computeNextTrickleProgress(), createTopProgressController(), getTopProgressBarStyle(), routeProgressController, shouldFinishTopProgress(), TOP_PROGRESS_INITIAL, TOP_PROGRESS_MAX_TRICKLE, TopProgressBar() (+2 more)

### Community 142 - "ref_bun_test"
Cohesion: 0.13
Nodes (14): DASHBOARD_ROUTE_PATHS, walkTsFiles(), WEB_SRC_DIR, consoleSource, detailSource, modalsSource, providersDir, modalsSource (+6 more)

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

### Community 148 - "newTestHarness"
Cohesion: 0.15
Nodes (19): Task 3: `settings` Module (Typed Store, Revision History & Rollback) & OpenAPI Fragment, Task 4: `keys` Module (API Keys CRUD, SHA-256 Cache & `contracts.KeyValidator`) & OpenAPI Fragment, assertFixedWidthUTC(), newTestHarness(), sha256Hex(), TestKeyInterceptorEnforcesRPMAndQuota(), TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart(), TestKeysCRUDAndSecretHygiene() (+11 more)

### Community 149 - "properties"
Cohesion: 0.20
Nodes (10): additionalProperties, type, type, config, enabled, UpdatePluginRequest, additionalProperties, properties (+2 more)

### Community 150 - "dependencies"
Cohesion: 0.08
Nodes (26): clsx, @fontsource-variable/jetbrains-mono, @fontsource-variable/plus-jakarta-sans, openapi-fetch, @radix-ui/react-label, @radix-ui/react-slot, react-dom, tailwind-merge (+18 more)

### Community 151 - "cli/package.json"
Cohesion: 0.14
Nodes (13): bin, omeg4mon-router, description, engines, node, files, keywords, license (+5 more)

### Community 153 - "setupCatalogTestModule"
Cohesion: 0.15
Nodes (11): setupCatalogTestModule(), TestComboTestHonorsPinnedConnectionAndComboRef(), TestModelAndComboLiveProbeEndpoints(), TestModelCatalogAndHealthMatrix(), TestModelProbeOpenAIResponsesAndOAuthRefresh(), TestIdentifyModelEndpoint(), TestResolveAttachesComboRuntimeMetadata(), BenchmarkResolveStrategies() (+3 more)

### Community 154 - "setupTestMitm"
Cohesion: 0.24
Nodes (12): newCertManager(), New(), doMitmJSON(), setupTestMitm(), TestMitmAliasesCRUD(), TestMitmCAGenerateDownloadAndKeyPermissions(), TestMitmConcurrentStartAndChunkedCARotation(), TestMitmConfigValidationLoopbackOnly() (+4 more)

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

### Community 162 - "repository"
Cohesion: 0.19
Nodes (5): boolToInt(), BindingSummary, PoolSummary, repository, storedPool

### Community 163 - "net/http.Handler"
Cohesion: 0.13
Nodes (20): loadUI(), loadUI(), createLegacySQLite(), setupImporterTest(), TestImporterRejectsInvalidFilesAndPreservesReadOnlySource(), setupPluginsTest(), TestPluginsBuiltinsUpdateAndInterceptorScoping(), putPluginConfig() (+12 more)

### Community 164 - "ImporterHistoryListResponse"
Cohesion: 0.20
Nodes (10): additionalProperties, properties, required, type, $ref, runs, items, items (+2 more)

### Community 165 - "PluginListResponse"
Cohesion: 0.22
Nodes (9): $ref, additionalProperties, properties, required, type, items, type, plugins (+1 more)

### Community 166 - "PreviewPluginsRequest"
Cohesion: 0.20
Nodes (10): type, additionalProperties, type, additionalProperties, properties, required, type, model (+2 more)

### Community 168 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan, Review Focus, Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`, Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`), Task 3: Upgraded `LoginForm.tsx`, Task 4: Upgraded `web/src/app/login/page.tsx`, Task 5: Full Verification & Graphify Update

### Community 169 - "Review Focus"
Cohesion: 0.22
Nodes (8): Global Constraints, Review Focus, Task 1: Build Automation, `Dockerfile` OCI Metadata & `scripts/smoke-test.sh`, Task 2: CI & Tag-Triggered Release Pipelines, Task 3: Production `README.md`, Task 4: Full Verification & Graphify Update, Task 5: Merged Worktree Cleanup, PR & Cutover Report, Wave 3F — Release & Cutover Readiness Implementation Plan

### Community 170 - "routingStateCache"
Cohesion: 0.25
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

### Community 178 - "setupMediaTest"
Cohesion: 0.24
Nodes (9): doJSONReq(), setupMediaTest(), TestMediaGatewayAudioVideoSearchEndpoints(), TestMediaGatewayImagesFallbackAndClientErrorNoFallback(), TestMediaProvidersAndCombosCRUD(), TestMediaWebFetchSSRFGuardAndExtraction(), TestLoadOrCreateMasterKey(), LoadOrCreateMasterKey() (+1 more)

### Community 179 - "parameters"
Cohesion: 0.25
Nodes (8): parameters, in, name, required, schema, maximum, minimum, type

### Community 180 - "repository"
Cohesion: 0.40
Nodes (3): ImportRun, repository, boolToInt()

### Community 181 - "New"
Cohesion: 0.07
Nodes (31): TestWave2AE2EProviderProxyRouting(), File Map, Review Focus, Task 3: Live Log Streaming (`logging`), `contracts.ProtocolTranslator` & `translator` Module, Task 4: Wire `proxypools` & `translator` in `cmd/omeg4mon-router`, OpenAPI Codegen & E2E Tests, Task 5: Web UI — Proxy Pools Studio (`/dashboard/proxy-pools/`), Task 6: Web UI — Translator Workbench & Live Console (`/dashboard/translator/`), Navigation Update & Single-Binary Verification, Wave 2A — Provider Expansion, Proxy Pools (`proxypools`) & Protocol Translator Workbench (`translator`) Implementation Plan (+23 more)

### Community 182 - "resolve.js"
Cohesion: 0.20
Nodes (8): child, __dirname, pkgRoot, ARCH_MAP, binaryFileName(), isExecutableFile(), PLATFORM_MAP, resolveBinaryPath()

### Community 183 - "ClearImporterHistoryResponse"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, type, deleted, ClearImporterHistoryResponse

### Community 184 - "source_path"
Cohesion: 0.29
Nodes (7): additionalProperties, properties, required, type, source_path, InspectLegacyRequest, type

### Community 185 - ".Snapshot"
Cohesion: 0.14
Nodes (7): JobInspector, Runner, job, JobStatus, stubJobInspector, jobInspectorAdapter, stubJobInspector

### Community 186 - "playgroundUtils.test.ts"
Cohesion: 0.29
Nodes (9): computeSessionTelemetry(), filterChatSessions(), formatLatencyMs(), formatTokenCount(), SessionTelemetry, sampleMessages, sampleSessions, ChatMessage (+1 more)

### Community 187 - "docker-release"
Cohesion: 0.29
Nodes (7): name, needs, permissions, runs-on, timeout-minutes, docker-release, packages

### Community 188 - "palette.test.ts"
Cohesion: 0.33
Nodes (3): GLOBALS_CSS_PATH, walkSourceFiles(), WEB_SRC_DIR

### Community 189 - "`@omeg4mon/router` CLI Launcher"
Cohesion: 0.50
Nodes (3): Binary Resolution Order, `@omeg4mon/router` CLI Launcher, Usage

### Community 190 - "SessionSigner"
Cohesion: 0.31
Nodes (5): Claims, Global Constraints, Review Focus, Wave 1A — Control Plane Foundation (`auth`, `settings`, `keys`) Implementation Plan, SessionSigner

### Community 191 - ".ExportCombos"
Cohesion: 0.36
Nodes (4): AliasSnapshot, ComboSnapshot, ComboSnapshotter, stubCombos

### Community 192 - "go_pkg_encoding_json"
Cohesion: 0.21
Nodes (14): errorDisposition, handledGatewayError, rewriteAntigravityPromptText(), rewriteAntigravitySystemInstruction(), classifyStatus(), isClientPayloadErrorBody(), cloneShallowMap(), extractCodeBuddyText() (+6 more)

### Community 193 - "context.Context"
Cohesion: 0.03
Nodes (19): catalogFakeProviderStore, fakeProviderStore, staticProviderStore, SettingsSnapshotter, stubProviderStore, stubProviderStore, stubKeyProvisioner, stubSettings (+11 more)

### Community 194 - "setupTunnelTest"
Cohesion: 0.33
Nodes (7): doTunnelReq(), setupTunnelTest(), TestTunnelCloudflareLifecycleAndURLCapture(), TestTunnelConfigValidationAndPersistence(), TestTunnelLongOutputLineNoDeadlockAndTransientOverride(), TestTunnelStopTerminatesRunningProcessOnShutdown(), TestTunnelTailscaleLifecycleAndMissingBinary()

### Community 195 - "repository"
Cohesion: 0.40
Nodes (3): boolToInt(), Config, repository

### Community 196 - "parameters"
Cohesion: 0.20
Nodes (10): put, in, name, required, schema, operationId, parameters, summary (+2 more)

### Community 197 - "ResolvedTarget"
Cohesion: 0.23
Nodes (7): stubProviderResolver, stubResolver, ResolvedTarget, isModelAllowed(), isModelDisabled(), prioritizeHeadroomAffinity(), TestPrioritizeHeadroomAffinityPromotesMatchingCandidateWithoutLosingFallbacks()

### Community 198 - "PromptTemplateCombobox.tsx"
Cohesion: 0.43
Nodes (6): BUILTIN_PROMPT_TEMPLATES, BuiltinPromptTemplate, filterPromptTemplates(), PromptTemplateCombobox(), PromptTemplateComboboxProps, templatePreviewSnippet()

### Community 199 - "omeg4mon-router Verification Gate"
Cohesion: 0.20
Nodes (9): 1. Scope Detection & Deterministic Gate, 2. Cross-Cutting Hygiene: Naming Clarity & Zero Dead Code (always), 3. Backend & CLI Audit (when Go / CLI / Contracts changed), 4.1 Code, Type Safety & Memory Safety, 4.2 Design System, Aesthetics & UX Quality, 4.3 Runtime Visual Verification, 4. Web & UI/UX Audit (when `web/` changed), 5. Exit Report (+1 more)

### Community 200 - "repo"
Cohesion: 0.33
Nodes (5): newID(), nowUTC(), ConfigRecord, Event, repo

### Community 201 - "repo"
Cohesion: 0.09
Nodes (16): buildWhere(), formatUTC(), newNotificationID(), normalizeSeverity(), scanNotification(), scanNotificationRow(), listSettingsWithQueryer(), decodeStoredValue() (+8 more)

### Community 202 - "single-binary"
Cohesion: 0.40
Nodes (5): single-binary, name, needs, runs-on, timeout-minutes

### Community 204 - "verify-web"
Cohesion: 0.50
Nodes (4): verify-web, name, runs-on, timeout-minutes

### Community 205 - "headroom.go"
Cohesion: 0.33
Nodes (4): CompressionStats, HeadroomStats, StatsResponse, WarmPrefixSummary

### Community 207 - "Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)"
Cohesion: 0.33
Nodes (5): Approaches Considered, Intent & Why, Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`), Proportionality & Standardization Contract, Verification

### Community 208 - "NewStreamTranslator"
Cohesion: 0.25
Nodes (8): NewStreamTranslator(), newAnthropicStreamEncoder(), newGeminiStreamEncoder(), newOpenAIStreamDecoder(), newOpenAIStreamEncoder(), newResponsesStreamEncoder(), TestResponsesStreamEncoderUniqueDefaultIDs(), TestStreamTranslators()

### Community 210 - "Feature Suite"
Cohesion: 0.40
Nodes (5): AI Studios & Extensions, Feature Suite, Infrastructure, Networking & Security, Operations & Observability, Routing & Gateway Core

### Community 213 - "Post-Login Dashboard Layout & Categorized Searchable Sidebar Design"
Cohesion: 0.29
Nodes (6): Approaches Considered, Intent & Why, Navigation Taxonomy (29/29 Routes Preserved), Post-Login Dashboard Layout & Categorized Searchable Sidebar Design, Proportionality & Standardization Contract, Search & Keyboard UX

### Community 214 - "net/http/httptest.ResponseRecorder"
Cohesion: 0.11
Nodes (28): File Map, Task 1: Cross-Module `AdminCredentialStore` Contract & `system` Adapter, Task 2: `auth` Module & OpenAPI Fragment, New(), decodeErrorCode(), doRequest(), findCookie(), newAuthHarness() (+20 more)

### Community 219 - "net/http.Header"
Cohesion: 0.09
Nodes (12): memoryResponseWriter, extractGatewayError(), newMemoryResponseWriter(), stringifyMessageContent(), copyResponseHeaders(), parseHeaderInt(), parseResetDuration(), copyNonHopHeaders() (+4 more)

### Community 220 - "Provider Setup Modals Redesign & Pre-Save Connectivity Test"
Cohesion: 0.33
Nodes (5): Constraints, Goals, Intent & Why, Non-Goals, Provider Setup Modals Redesign & Pre-Save Connectivity Test

### Community 223 - "DB"
Cohesion: 0.05
Nodes (60): main(), runMain(), ChatCompletionRequest, ChatCompletionResult, ChatMessage, File Map, Global Constraints, Review Focus (+52 more)

### Community 225 - "optimization.test.ts"
Cohesion: 0.25
Nodes (7): consoleSource, customModelModalSource, detailSource, hooksSource, logoSource, modalsSource, providersDir

### Community 226 - "AttestationManager"
Cohesion: 0.15
Nodes (14): selectProbeHTTPClient(), bunBoringSSLClientHelloSpec(), NewChromeTLSClient(), SharedChromeTLSClient(), TestSharedChromeTLSClientHandlesHTTPAndReturnsSingleton(), getSharedCookieJar(), AttestationManager, isTransientMintStatus() (+6 more)

### Community 228 - "Implementation Plan: `/dashboard/models` OmniRoute Parity"
Cohesion: 0.33
Nodes (5): Implementation Plan: `/dashboard/models` OmniRoute Parity, Task 2: OpenAPI Contract & TypeScript Codegen, Task 3: Route Split & Navigation Update, Task 4: Frontend OmniRoute Model Catalog (`web/src/features/models/`), Task 5: Full Verification & Graphify Update

### Community 231 - "time.Duration"
Cohesion: 0.09
Nodes (7): stubInterceptor, recordingInterceptor, CanonicalRequest, ResponseBodyInterceptor, ResponseMeta, StreamEvent, newID()

### Community 235 - "omeg4mon-router Project Standards"
Cohesion: 0.33
Nodes (5): Architecture & Boundaries, Design System & UI/UX Identity, Mandatory Verification, omeg4mon-router Project Standards, Reliability, Memory & Type Safety

### Community 241 - ".StatusCode"
Cohesion: 0.25
Nodes (3): statusError, statusError, httpStatusError

### Community 242 - "compressor.go"
Cohesion: 0.10
Nodes (36): buildOutputDirective(), compressContentBlocks(), compressMessages(), compressPayload(), compressRoot(), compressText(), estimateTokens(), hasAnthropicToolResult() (+28 more)

### Community 243 - "Status"
Cohesion: 0.50
Nodes (4): buildInstallGuide(), InstallGuide, Status, ToolDescriptor

## Knowledge Gaps
- **7 isolated node(s):** `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer`, `Wave 0 Platform Foundation Implementation Plan`, `Canonical IR & Direct Protocol Translators` (+2 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 1077 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **38 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Provider Setup Modals Redesign & Pre-Save Test Implementation Plan` connect `ProviderDetailView.tsx` to `Module`?**
  _High betweenness centrality (0.187) - this node is a cross-community bridge._
- **What connects `github.com/omeg4monx/router-gateway`, `Structured Logging RingBuffer & Non-Blocking EventBus`, `Platform Crypto Vault, Bcrypt & Session Signer` to the rest of the system?**
  _7 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `combos/module.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07692307692307693 - nodes in this community are weakly interconnected._
- **Why does `buildProtocolTestRequest()` connect `Module` to `context.Context`, `net/http.Request`, `go_pkg_encoding_json`, `.Fetch`, `testing.T`, `oauth_registry.go`, `ProviderDetailView.tsx`?**
  _High betweenness centrality (0.152) - this node is a cross-community bridge._
- **Should `biome.json` be split into smaller, more focused modules?**
  _Cohesion score 0.04625346901017576 - nodes in this community are weakly interconnected._
- **Why does `File Map` connect `net/http/httptest.ResponseRecorder` to `SessionSigner`, `client.ts`, `newTestHarness`, `DB`, `Button`, `Module`, `DB`?**
  _High betweenness centrality (0.068) - this node is a cross-community bridge._
- **Should `settings` be split into smaller, more focused modules?**
  _Cohesion score 0.05263157894736842 - nodes in this community are weakly interconnected._