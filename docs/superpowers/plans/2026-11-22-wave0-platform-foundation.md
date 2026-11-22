# Wave 0: Platform Foundation & Single-Binary Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the modular Go + Next.js (TypeScript) platform foundation, OpenAPI codegen pipeline, `system` module, `DESIGN.md` dashboard shell, and `//go:embed` single-binary build pipeline on branch `migration/ts-go`.

**Architecture:** Modular Monolith in Go (`net/http` + `chi/v5`) organized into cross-cutting infrastructure (`internal/platform/`), narrow cross-module contracts (`internal/contracts/`), and vertical-slice modules (`internal/modules/<domain>/`) registered via an explicit topological `Module` registry. Next.js (App Router, strict TypeScript, `output: "export"`, Tailwind CSS v4) lives in `web/`, consumes OpenAPI-generated TypeScript types from `contracts/openapi/`, and is embedded into the Go binary via `//go:embed` behind the `embedui` build tag.

**Tech Stack:** Go 1.24+ (target 1.26, pure-Go `modernc.org/sqlite`, `github.com/go-chi/chi/v5`, `golang.org/x/crypto`), Next.js 16 (App Router, `output: "export"`), React 19, TypeScript 5 (`strict`), Tailwind CSS v4, Bun 1.4+, `openapi-typescript`, `openapi-fetch`, `@tanstack/react-query`, `Makefile`, multi-stage `Dockerfile`.

**Spec:** [docs/superpowers/specs/2026-11-22-modular-ts-go-rewrite-design.md](../specs/2026-11-22-modular-ts-go-rewrite-design.md)

## Global Constraints

- **Clean-room ownership (Owner Decision 1A):** 0% verbatim copying of source code, SQL, or third-party assets from `main`. All Go, TypeScript, SQL, OpenAPI, and UI files are authored fresh. `LICENSE` is MIT under `Copyright (c) 2026 OMEG4MON`. Only owner-authored visual rules from `DESIGN.md` are preserved.
- **Zero-CGO Go build:** Persistence MUST use pure-Go `modernc.org/sqlite` (WAL mode, `busy_timeout=10000`, `synchronous=NORMAL`, `foreign_keys=ON`, `temp_store=MEMORY`, `cache_size=-2000`, `mmap_size=0`). `CGO_ENABLED=0` must succeed cleanly.
- **HTTP stack:** Standard library `net/http` + thin router `github.com/go-chi/chi/v5` only. No heavy Go web frameworks (`Gin`, `Echo`, `Fiber`).
- **Single-admin dashboard (Owner Decision 2A):** No multi-user RBAC or tenant row isolation in Wave 0/1.
- **Single-binary static export:** `web/next.config.ts` MUST use `output: "export"` and `images: { unoptimized: true }`. Dynamic entity IDs use static route shells + query parameters (e.g., `/dashboard/providers/detail?id=...`).
- **Visual identity (`docs/DESIGN.md`):** Accent `#E56A4A`, dark base `#0a0a0a`, light base `#FDFAF6`, neutral `#6B7280`, Geist Sans & Geist Mono, caduceus SVG logo, ember-line gradient motif (`from-transparent via-primary to-transparent`), solid opaque surfaces (no glassmorphism), zero UI emojis, WCAG AA contrast.
- **Network defaults:** Default bind `127.0.0.1:20128`; env vars `OMEG4MON_HOST`, `OMEG4MON_PORT`, `OMEG4MON_DATA_DIR`, `OMEG4MON_MASTER_KEY`, `OMEG4MON_ADMIN_PASSWORD`, `OMEG4MON_DISABLED_MODULES`, `OMEG4MON_ALLOW_INSECURE_REMOTE`.
- **Toolchain prerequisite:** Go 1.24+ (`go` on `PATH`) and Bun 1.4+ (`bun` on `PATH`) are required to run verification commands. If `go` is absent on the host `PATH`, install official Go Linux tarball to `$HOME/.local/go` and export `PATH="$HOME/.local/go/bin:$HOME/.bun/bin:$PATH"` before starting Task 1.

## Review Focus

1. **Core module disable attempt & non-loopback host classification:** `OMEG4MON_DISABLED_MODULES` containing a core module (`system`, `auth`, `gateway`, etc.) must be rejected with a clear error, and non-loopback hosts (`0.0.0.0`, `::`, LAN IP) must be detected accurately to enforce remote-bind security (pinned in Task 1 & Task 8).
2. **Slow/stalled `eventbus` subscriber under burst traffic:** A subscriber that stops draining its channel must never block `Publish()` or deadlock the gateway hot path, and calling `unsubscribe()` must cleanly release resources (pinned in Task 2).
3. **Corrupt vault ciphertext, malformed `OMEG4MON_MASTER_KEY`, and expired/tampered session cookies:** Must fail closed with typed errors and never panic; auto-generated `master.key` must be written with `0600` permissions (pinned in Task 3).
4. **Concurrent WAL reads/writes and failed SQL migration rollback:** Parallel readers + writers must complete without `SQLITE_BUSY`, and a broken migration SQL file must roll back atomically without inserting a row into `schema_migrations` (pinned in Task 4).
5. **Streaming `/v1/*` timeout exemption vs. `/api/*` JSON 404 isolation:** Long-lived SSE responses on `/v1/*` must flush chunks past the `/api/*` middleware timeout without truncation, and unknown `/api/*` or `/v1/*` paths must return JSON `404`—never falling through to the SPA `index.html` (pinned in Task 7).

---

### Task 1: Repository Skeleton, Clean MIT License, `DESIGN.md`, `go.mod`, and Config Loader

**Files:**
- Create: `LICENSE`
- Create: `.gitignore`
- Create: `docs/DESIGN.md`
- Create: `go.mod`
- Create: `internal/platform/config/config.go`
- Test: `internal/platform/config/config_test.go`

**Interfaces:**
- Consumes: none
- Produces:
  - `config.Config` struct (`Host string`, `Port int`, `DataDir string`, `MasterKeyHex string`, `AdminPassword string`, `DisabledModules map[string]bool`, `AllowInsecureRemote bool`, `LogLevel string`, `APITimeout time.Duration`)
  - `config.CoreModules() []string` (`system`, `auth`, `settings`, `providers`, `keys`, `combos`, `gateway`, `usage`)
  - `config.Load(args []string, getenv func(string) string) (Config, error)`
  - `(c Config) Addr() string`
  - `(c Config) IsLoopback() bool`
  - `(c Config) IsModuleDisabled(name string) bool`

- [ ] **Step 1: Verify Go and Bun toolchains on `PATH`**

Run: `go version && bun --version`
Expected: Go `1.24+` and Bun `1.4+` printed with exit code 0. (If `go` is missing, install official Go binary under `~/.local/go` and ensure `~/.local/go/bin` is on `PATH`.)

- [ ] **Step 2: Create repo metadata (`LICENSE`, `.gitignore`, `docs/DESIGN.md`, `go.mod`)**

- `LICENSE`: standard MIT License text with `Copyright (c) 2026 OMEG4MON`.
- `.gitignore`: ignore `bin/`, `dist/`, `web/node_modules/`, `web/.next/`, `web/out/`, `cmd/omeg4mon-router/ui_dist/`, `*.sqlite*`, `master.key`, `.env`, `.env.local`.
- `docs/DESIGN.md`: owner visual direction (`#E56A4A`, `#0a0a0a`, `#FDFAF6`, `#6B7280`, Geist Sans/Mono, caduceus logo, ember-line motif, solid surfaces, ENERGY 2 / RHYTHM 2 / MOTION 1, no UI emojis, WCAG AA).
- `go.mod`: `module github.com/omeg4monx/router-gateway`, `go 1.24`.

- [ ] **Step 3: Write the failing config tests in `internal/platform/config/config_test.go`**

```go
func TestLoadDefaults(t *testing.T) {
    cfg, err := Load(nil, func(string) string { return "" })
    // assert err == nil, cfg.Host == "127.0.0.1", cfg.Port == 20128, cfg.IsLoopback() == true
}

func TestLoadRejectsInvalidPortAndCoreModuleDisable(t *testing.T) {
    // 1) OMEG4MON_PORT=99999 -> error
    // 2) OMEG4MON_DISABLED_MODULES="mitm,gateway" -> error naming core module "gateway"
    // 3) OMEG4MON_DISABLED_MODULES="mitm, telegram" -> ok, IsModuleDisabled("mitm")==true, IsModuleDisabled("telegram")==true
}

func TestIsLoopbackClassification(t *testing.T) {
    // true: "127.0.0.1", "::1", "localhost"
    // false: "0.0.0.0", "::", "192.168.1.10", "example.com"
}
```

- [ ] **Step 4: Run test to verify it fails**

Run: `go test ./internal/platform/config/...`
Expected: FAIL (`undefined: Load` / no Go files).

- [ ] **Step 5: Implement `internal/platform/config/config.go`**

Use `flag.NewFlagSet("omeg4mon-router", flag.ContinueOnError)` with env defaults (`OMEG4MON_*`) overridden by flags (`--host`, `--port`, `--data-dir`, `--log-level`). Validate port `1..65535`, resolve default `DataDir` to `$HOME/.omeg4mon-router`, parse comma-separated `OMEG4MON_DISABLED_MODULES`, reject any name in `CoreModules()`, and classify loopback via `net.ParseIP(host).IsLoopback()` plus `strings.EqualFold(host, "localhost")`.

- [ ] **Step 6: Run tests to verify pass**

Run: `go test -v -race ./internal/platform/config/...`
Expected: PASS (all 3 tests pass with `-race`).

- [ ] **Step 7: Commit**

```bash
git add LICENSE .gitignore docs/DESIGN.md go.mod internal/platform/config/
git commit -m "chore(foundation): init clean repo metadata and platform config loader"
```

---

### Task 2: Structured Logging (`internal/platform/logging`) & Event Bus (`internal/platform/eventbus`)

**Files:**
- Create: `internal/platform/logging/logger.go`
- Test: `internal/platform/logging/logger_test.go`
- Create: `internal/platform/eventbus/bus.go`
- Test: `internal/platform/eventbus/bus_test.go`

**Interfaces:**
- Consumes: stdlib `log/slog`
- Produces:
  - `logging.LogEntry` (`Time time.Time`, `Level string`, `Message string`, `Attrs map[string]any`)
  - `logging.New(level string, w io.Writer, ringCap int) (*slog.Logger, *RingBuffer)`
  - `(*RingBuffer).Recent(limit int) []LogEntry`
  - `eventbus.Event` (`Topic string`, `Payload any`, `Timestamp time.Time`)
  - `eventbus.New() *Bus`
  - `(*Bus).Publish(topic string, payload any)`
  - `(*Bus).Subscribe(topic string, buf int) (<-chan Event, func())`

- [ ] **Step 1: Write failing tests for logging ring buffer and non-blocking eventbus**

In `logger_test.go`:
```go
func TestRingBufferCapsAndPreservesOrder(t *testing.T) {
    // ringCap=3, write 5 Info logs -> Recent(10) returns last 3 messages in chronological order
}
```
In `bus_test.go`:
```go
func TestPublishNeverBlocksSlowSubscriber(t *testing.T) {
    // Subscribe("usage.recorded", 1), publish 100 events without reading -> completes <100ms
    // Call unsubscribe() twice (idempotent) -> channel closed, subsequent Publish succeeds
}

func TestPublishDeliversToMultipleSubscribersOnMatchingTopic(t *testing.T) {
    // Two subscribers on "config.changed" receive payload; subscriber on "other" receives nothing
}
```

- [ ] **Step 2: Run tests to verify failure**

Run: `go test ./internal/platform/logging/... ./internal/platform/eventbus/...`
Expected: FAIL (missing packages/symbols).

- [ ] **Step 3: Implement `logging/logger.go` and `eventbus/bus.go`**

- `logging`: wrap `slog.NewJSONHandler` with a custom `slog.Handler` tee that records entries into a mutex-protected circular slice (`RingBuffer`).
- `eventbus`: `sync.RWMutex` map of `topic -> map[uint64]chan Event`; `Publish` uses non-blocking `select { case ch <- evt: default: }` so slow subscribers drop excess events instead of stalling publishers; `unsubscribe` uses `sync.Once` to delete and close the channel safely under lock.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/platform/logging/... ./internal/platform/eventbus/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/platform/logging/ internal/platform/eventbus/
git commit -m "feat(platform): add structured slog ring buffer and non-blocking event bus"
```

---

### Task 3: Secret Vault, Master Key Lifecycle & Session Crypto (`internal/platform/crypto`)

**Files:**
- Create: `internal/platform/crypto/vault.go`
- Create: `internal/platform/crypto/password.go`
- Create: `internal/platform/crypto/session.go`
- Test: `internal/platform/crypto/crypto_test.go`

**Interfaces:**
- Consumes: `golang.org/x/crypto/bcrypt` (add to `go.mod` via `go get`)
- Produces:
  - `crypto.LoadOrCreateMasterKey(dataDir, envHex string) ([]byte, error)`
  - `crypto.NewVault(key []byte) (*Vault, error)`
  - `(*Vault).Encrypt(plaintext []byte) (string, error)`
  - `(*Vault).Decrypt(encoded string) ([]byte, error)`
  - `crypto.HashPassword(password string) (string, error)`
  - `crypto.VerifyPassword(hash, password string) bool`
  - `crypto.Claims` (`Subject string`, `ExpiresAt time.Time`)
  - `crypto.NewSessionSigner(key []byte) (*SessionSigner, error)`
  - `(*SessionSigner).Sign(subject string, ttl time.Duration) (string, error)`
  - `(*SessionSigner).Verify(token string) (Claims, error)`

- [ ] **Step 1: Add `golang.org/x/crypto` and write failing crypto tests in `crypto_test.go`**

```go
func TestLoadOrCreateMasterKey(t *testing.T) {
    // 1) Creates <tmp>/master.key with 0600 permissions, reloads identical 32 bytes on second call
    // 2) Valid 64-char hex env override takes precedence without creating file
    // 3) Malformed hex or wrong length returns error
}

func TestVaultRoundTripAndTamperRejection(t *testing.T) {
    // Encrypt same plaintext twice -> distinct ciphertexts (random 12-byte nonce)
    // Decrypt succeeds; bit-flipping ciphertext or using wrong key fails with error
}

func TestPasswordHashAndVerify(t *testing.T) {
    // Rejects empty password; HashPassword -> VerifyPassword true for match, false for mismatch
}

func TestSessionSignVerifyAndExpiry(t *testing.T) {
    // Sign("admin", time.Hour) -> Verify returns Subject="admin"
    // Tampered signature -> error; negative TTL -> expired error
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./internal/platform/crypto/...`
Expected: FAIL.

- [ ] **Step 3: Implement `vault.go`, `password.go`, `session.go`**

- `LoadOrCreateMasterKey`: validate 32-byte key from `hex.DecodeString(envHex)` if non-empty; otherwise `os.MkdirAll(dataDir, 0700)`, read `<dataDir>/master.key`, or generate 32 bytes via `crypto/rand` and write atomically via temp file + `os.Chmod(0600)` + `os.Rename`.
- `Vault`: AES-256-GCM; format `v1:<base64url(nonce || ciphertext)>`.
- `Password`: `bcrypt.GenerateFromPassword` with cost `bcrypt.DefaultCost` (reject empty password).
- `SessionSigner`: derive subkey via HMAC-SHA256("session-key", masterKey); encode `base64url(json(Claims)).base64url(hmac)`; verify with `hmac.Equal` and `time.Now().UTC().Before(claims.ExpiresAt)`.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/platform/crypto/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum internal/platform/crypto/
git commit -m "feat(platform): add AES-256-GCM vault, master key lifecycle, and session crypto"
```

---

### Task 4: Zero-CGO Split-Pool SQLite WAL Engine & Modular Migrator (`internal/platform/db`)

**Files:**
- Create: `internal/platform/db/db.go`
- Create: `internal/platform/db/migrator.go`
- Test: `internal/platform/db/db_test.go`

**Interfaces:**
- Consumes: `modernc.org/sqlite` (add via `go get modernc.org/sqlite`)
- Produces:
  - `db.Open(path string) (*DB, error)`
  - `(*DB).Writer() *sql.DB` (max 1 open conn)
  - `(*DB).Reader() *sql.DB` (max 8 open conns, read-only query path)
  - `(*DB).Ping(ctx context.Context) error`
  - `(*DB).MigrateModule(ctx context.Context, module string, fsys fs.FS) error`
  - `(*DB).AppliedVersions(ctx context.Context, module string) ([]int, error)`
  - `(*DB).Close() error`

- [ ] **Step 1: Add `modernc.org/sqlite` and write failing DB & migrator tests in `db_test.go`**

```go
func TestOpenAppliesWALPragmasAndSplitPools(t *testing.T) {
    // Open(<tmp>/data.sqlite) -> PRAGMA journal_mode == "wal", PRAGMA foreign_keys == 1
    // Writer().Stats().MaxOpenConnections == 1; Reader().Stats().MaxOpenConnections == 8
}

func TestMigrateModuleIdempotentAndAtomicRollback(t *testing.T) {
    // fstest.MapFS with 0001_init.sql and 0002_add_col.sql -> applied once, re-run is no-op
    // fstest.MapFS with 0003_bad.sql containing valid DDL + syntax error -> returns error,
    // version 3 NOT in schema_migrations, and partial DDL rolled back
}

func TestConcurrentReadersAndWritersUnderWAL(t *testing.T) {
    // 8 reader goroutines + 4 writer goroutines performing 50 ops each -> zero SQLITE_BUSY errors
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./internal/platform/db/...`
Expected: FAIL.

- [ ] **Step 3: Implement `db.go` and `migrator.go`**

- Apply DSN pragmas (`busy_timeout(10000)`, `journal_mode(WAL)`, `synchronous(NORMAL)`, `foreign_keys(ON)`, `temp_store(MEMORY)`, `cache_size(-2000)`, `mmap_size(0)`). Ensure parent dir exists (`0750`). Bootstrap `CREATE TABLE IF NOT EXISTS schema_migrations (module TEXT NOT NULL, version INTEGER NOT NULL, applied_at TEXT NOT NULL, PRIMARY KEY (module, version))`.
- `MigrateModule`: skip if `fsys == nil`; parse leading integer version from `NNNN_*.sql`, sort ascending, execute each unapplied SQL file inside `Writer().BeginTx`, insert version row, commit (or rollback on error).

- [ ] **Step 4: Run tests with zero-CGO verification**

Run: `CGO_ENABLED=0 go test -v ./internal/platform/db/...`
Expected: PASS (all 3 tests pass with `CGO_ENABLED=0`; also run `go test -race ./internal/platform/db/...` when a C compiler is present).

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum internal/platform/db/
git commit -m "feat(platform): add zero-CGO split-pool SQLite WAL engine and modular migrator"
```

---

### Task 5: Upstream HTTP/2 Client Pool (`internal/platform/httpclient`) & Job Runner (`internal/platform/jobs`)

**Files:**
- Create: `internal/platform/httpclient/client.go`
- Test: `internal/platform/httpclient/client_test.go`
- Create: `internal/platform/jobs/runner.go`
- Test: `internal/platform/jobs/runner_test.go`

**Interfaces:**
- Consumes: `log/slog`
- Produces:
  - `httpclient.Options` (`DialTimeout time.Duration`, `ResponseHeaderTimeout time.Duration`, `IdleConnTimeout time.Duration`, `MaxIdleConns int`, `MaxIdleConnsPerHost int`, `ProxyFunc func(*http.Request) (*url.URL, error)`)
  - `httpclient.NewClient(opts Options) *http.Client`
  - `jobs.NewRunner(logger *slog.Logger) *Runner`
  - `(*Runner).Every(name string, interval time.Duration, fn func(context.Context) error)`
  - `(*Runner).Start(ctx context.Context)`
  - `(*Runner).Stop(ctx context.Context) error`

- [ ] **Step 1: Write failing tests for HTTP client streaming/header timeout and job runner panic recovery**

```go
// client_test.go
func TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(t *testing.T) {
    // 1) httptest server delays headers 150ms vs ResponseHeaderTimeout=50ms -> fails with timeout
    // 2) httptest server sends headers immediately then streams 2 chunks 80ms apart -> succeeds (Timeout==0)
}

// runner_test.go
func TestRunnerRecoversFromPanicAndStopsCleanly(t *testing.T) {
    // Periodic job panics on run 1, succeeds on run 2 -> runner stays alive, Stop() returns nil cleanly
}
```

- [ ] **Step 2: Run tests to verify failure**

Run: `go test ./internal/platform/httpclient/... ./internal/platform/jobs/...`
Expected: FAIL.

- [ ] **Step 3: Implement `httpclient/client.go` and `jobs/runner.go`**

- `httpclient`: configure `http.Transport{ForceAttemptHTTP2: true, ResponseHeaderTimeout: opts.ResponseHeaderTimeout, ...}` with `http.Client{Timeout: 0}` so SSE streams are bounded only by request context cancellation and header timeout.
- `jobs`: goroutine per registered periodic task wrapped with `defer recover()` logging stack/error via `slog`, synchronized by `sync.WaitGroup` on `Stop(ctx)`.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/platform/httpclient/... ./internal/platform/jobs/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/platform/httpclient/ internal/platform/jobs/
git commit -m "feat(platform): add streaming HTTP/2 client pool and panic-safe job runner"
```

---

### Task 6: Cross-Module Contracts (`internal/contracts`) & Module Registry (`internal/platform/lifecycle`)

**Files:**
- Create: `internal/contracts/contracts.go`
- Create: `internal/platform/lifecycle/registry.go`
- Test: `internal/platform/lifecycle/registry_test.go`

**Interfaces:**
- Consumes: `github.com/go-chi/chi/v5`, `internal/platform/db`
- Produces:
  - `contracts.AuthVerifier` (`AuthenticateAdmin(r *http.Request) (string, error)`, `RequireLogin(ctx context.Context) bool`)
  - `contracts.KeyInfo`, `contracts.KeyValidator` (`ValidateKey(ctx context.Context, rawKey string) (KeyInfo, error)`)
  - `contracts.ResolvedTarget`, `contracts.ProviderResolver` (`Resolve(ctx context.Context, model string, allowed []string) ([]ResolvedTarget, error)`)
  - `contracts.CanonicalRequest`, `contracts.ResponseMeta`, `contracts.StreamEvent`, `contracts.GatewayInterceptor` (`Name() string`, `Priority() int`, `OnRequest(ctx, *CanonicalRequest) error`, `OnResponseHeaders(ctx, *ResponseMeta) error`, `OnStreamChunk(ctx, *StreamEvent) error`)
  - `contracts.UsageRecord`, `contracts.UsageRecorder` (`Record(ctx context.Context, rec UsageRecord)`)
  - `contracts.ProxyDialerResolver` (`ProxyURLForProvider(ctx context.Context, providerID string) (*url.URL, error)`)
  - `lifecycle.RouteMounts` (`PublicAPI chi.Router`, `AdminAPI chi.Router`, `Gateway chi.Router`)
  - `lifecycle.Module` (`Name() string`, `Dependencies() []string`, `Migrations() fs.FS`, `RegisterRoutes(RouteMounts)`, `Start(context.Context) error`, `Stop(context.Context) error`)
  - `lifecycle.NewRegistry(disabled map[string]bool) *Registry`
  - `(*Registry).Register(m Module) error`
  - `(*Registry).Ordered() ([]Module, error)`
  - `(*Registry).ActiveNames() []string`
  - `(*Registry).MigrateAll(ctx context.Context, database *db.DB) error`
  - `(*Registry).MountAll(mounts RouteMounts) error`
  - `(*Registry).StartAll(ctx context.Context) error`
  - `(*Registry).StopAll(ctx context.Context) error`

- [ ] **Step 1: Add `github.com/go-chi/chi/v5` and write failing registry tests in `registry_test.go`**

```go
func TestRegistryTopologicalOrderAndReverseStop(t *testing.T) {
    // Register C (deps: B), B (deps: A), A (deps: nil) -> StartAll order [A, B, C], StopAll order [C, B, A]
}

func TestRegistryRejectsCyclesAndMissingDeps(t *testing.T) {
    // 1) A -> B -> A returns cycle error
    // 2) B depends on missing X returns missing-dependency error
    // 3) Duplicate module name returns error
}

func TestRegistrySkipsDisabledOptionalModule(t *testing.T) {
    // Disabled "mitm" is skipped from ActiveNames/StartAll; if active module depends on disabled module -> error
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./internal/contracts/... ./internal/platform/lifecycle/...`
Expected: FAIL.

- [ ] **Step 3: Implement `internal/contracts/contracts.go` and `internal/platform/lifecycle/registry.go`**

Use deterministic Kahn topological sort (lexicographical tie-break by module name so startup order is repeatable). If `StartAll` fails at module `i`, automatically invoke `Stop` on already-started modules `i-1..0` before returning the error.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/contracts/... ./internal/platform/lifecycle/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum internal/contracts/ internal/platform/lifecycle/
git commit -m "feat(platform): add cross-module contracts and topological module registry"
```

---

### Task 7: `net/http` + `chi` Server, Split Middleware Stacks & SPA Static UI Handler (`internal/platform/httpserver`)

**Files:**
- Create: `internal/platform/httpserver/response.go`
- Create: `internal/platform/httpserver/static_ui.go`
- Create: `internal/platform/httpserver/server.go`
- Test: `internal/platform/httpserver/server_test.go`

**Interfaces:**
- Consumes: `github.com/go-chi/chi/v5`, `internal/contracts.AuthVerifier`, `internal/platform/config.Config`, `internal/platform/lifecycle.RouteMounts`
- Produces:
  - `httpserver.ErrorEnvelope` (`Error ErrorDetail` with `Code`, `Message`, `RequestID`)
  - `httpserver.WriteJSON(w http.ResponseWriter, status int, v any)`
  - `httpserver.WriteError(w http.ResponseWriter, r *http.Request, status int, code, msg string)`
  - `httpserver.NewStaticUIHandler(fsys fs.FS) http.Handler`
  - `httpserver.Options` (`Config config.Config`, `Logger *slog.Logger`, `Auth contracts.AuthVerifier`, `UI fs.FS`, `DevUIHandler http.Handler`)
  - `httpserver.New(opts Options) *Server`
  - `(*Server).Mounts() lifecycle.RouteMounts`
  - `(*Server).Handler() http.Handler`
  - `(*Server).Start(ctx context.Context) error`
  - `(*Server).Shutdown(ctx context.Context) error`

- [ ] **Step 1: Write failing server tests in `server_test.go`**

```go
func TestServerIsolatesJSON404FromSPAFallback(t *testing.T) {
    // Mount static FS containing index.html and dashboard/index.html
    // GET /api/unknown -> 404 application/json (NOT index.html)
    // GET /v1/unknown -> 404 application/json (NOT index.html)
    // GET /dashboard -> 200 text/html serving dashboard/index.html with Cache-Control: no-store
    // GET /_next/static/app.js -> 200 with Cache-Control containing immutable
}

func TestGatewaySSEStreamsPastAPITimeout(t *testing.T) {
    // Config.APITimeout = 60ms
    // 1) AdminAPI slow handler (150ms) -> ctx cancelled (504 / timeout)
    // 2) Gateway SSE handler flushing 2 chunks 90ms apart -> 200 OK receiving both chunks cleanly
}

func TestAdminAPIAuthEnforcement(t *testing.T) {
    // With AuthVerifier requiring login: unauthenticated AdminAPI -> 401 JSON; PublicAPI -> 200
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./internal/platform/httpserver/...`
Expected: FAIL.

- [ ] **Step 3: Implement `response.go`, `static_ui.go`, and `server.go`**

- Configure `http.Server{ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 120 * time.Second, WriteTimeout: 0}`.
- Attach `middleware.RequestID`, `middleware.RealIP`, panic recovery returning JSON `500` on API/Gateway routes.
- Mount `PublicAPI` and `AdminAPI` under `/api` (both with `middleware.Timeout(cfg.APITimeout)`; `AdminAPI` also checks `AuthVerifier` when non-nil). Register explicit JSON `NotFound` for `/api/*`, `/v1/*`, `/v1beta/*`.
- Mount `Gateway` under `/v1`, `/v1beta`, and compatibility aliases (`/codex/*`, `/responses`, `/models/*`) with **no** `middleware.Timeout`.
- `NewStaticUIHandler`: serve exact path, `<path>/index.html`, or `<path>.html` from embedded `fs.FS`, falling back to root `index.html` for non-asset client paths; set `Cache-Control: public, max-age=31536000, immutable` for `/_next/static/*` and `no-store` for `.html`.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/platform/httpserver/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/platform/httpserver/
git commit -m "feat(platform): add chi HTTP server, streaming timeout split, and SPA static handler"
```

---

### Task 8: `system` Module (`/api/health`, `/api/version`, `/api/init`, `/api/shutdown`) & OpenAPI Contract

**Files:**
- Create: `internal/modules/system/migrations/0001_init.sql`
- Create: `internal/modules/system/repo.go`
- Create: `internal/modules/system/module.go`
- Test: `internal/modules/system/module_test.go`
- Create: `contracts/openapi/openapi.yaml`
- Create: `contracts/openapi/modules/system.yaml`

**Interfaces:**
- Consumes: `internal/platform/db.DB`, `internal/platform/config.Config`, `internal/platform/crypto`, `internal/platform/httpserver`, `internal/platform/lifecycle`
- Produces:
  - `system.BuildInfo` (`Version string`, `Commit string`, `BuildTime string`)
  - `system.Options` (`DB *db.DB`, `Config config.Config`, `Build BuildInfo`, `ActiveModules func() []string`, `Shutdown func()`)
  - `system.New(opts Options) *Module` (implements `lifecycle.Module`, name `"system"`)
  - Public endpoints: `GET /api/health`, `GET /api/version`, `GET /api/init`, `POST /api/init`
  - Admin endpoint: `POST /api/shutdown`
  - OpenAPI 3.1 contract in `contracts/openapi/openapi.yaml` + `contracts/openapi/modules/system.yaml`

- [ ] **Step 1: Write failing tests in `internal/modules/system/module_test.go`**

```go
func TestSystemHealthAndVersion(t *testing.T) {
    // GET /api/health -> 200 {status:"ok", db:"ok", modules:["system"], uptimeSeconds:>=0}
    // GET /api/version -> 200 {version:"0.1.0-dev", commit:"abc", goVersion:...}
}

func TestSystemInitLifecycle(t *testing.T) {
    // 1) GET /api/init -> {initialized:false, requireLogin:true}
    // 2) POST /api/init with short password (<8 chars) -> 400
    // 3) POST /api/init with valid password -> 200 {initialized:true}
    // 4) Second POST /api/init -> 409 Conflict
    // 5) When Config.AdminPassword env is set -> Start() auto-bootstraps initialized=true
}

func TestSystemShutdownRemoteBindGuard(t *testing.T) {
    // 1) Loopback config -> POST /api/shutdown returns 202 and invokes Shutdown callback
    // 2) Non-loopback (0.0.0.0) with AllowInsecureRemote=false -> 403 Forbidden, callback NOT called
    // 3) Non-loopback with AllowInsecureRemote=true -> 202 and invokes callback
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./internal/modules/system/...`
Expected: FAIL.

- [ ] **Step 3: Implement `migrations/0001_init.sql`, `repo.go`, `module.go`, and OpenAPI specs**

- `0001_init.sql`: `CREATE TABLE system_state (key TEXT PRIMARY KEY, value TEXT NOT NULL, updated_at TEXT NOT NULL);`
- `module.go`: embed `migrations/*.sql`, store bcrypt hash under `admin_password_hash` and `require_login` (`"true"` default), enforce non-loopback shutdown guard, and keep `/api/health` fast using `db.Ping`.
- Write `contracts/openapi/openapi.yaml` and `contracts/openapi/modules/system.yaml` with exact request/response schemas matching the Go JSON tags.

- [ ] **Step 4: Run tests with `-race`**

Run: `go test -v -race ./internal/modules/system/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/modules/system/ contracts/openapi/
git commit -m "feat(system): add system module endpoints and modular OpenAPI 3.1 contract"
```

---

### Task 9: Next.js TypeScript Static-Export Frontend (`web/`), OpenAPI Codegen & `DESIGN.md` Shell

**Files:**
- Create: `web/package.json`
- Create: `web/tsconfig.json`
- Create: `web/next.config.ts`
- Create: `web/postcss.config.mjs`
- Create: `web/src/app/globals.css`
- Create: `web/src/app/layout.tsx`
- Create: `web/src/app/page.tsx`
- Create: `web/src/app/init/page.tsx`
- Create: `web/src/app/login/page.tsx`
- Create: `web/src/app/(dashboard)/dashboard/layout.tsx`
- Create: `web/src/app/(dashboard)/dashboard/page.tsx`
- Create: `web/src/shared/api/schema.d.ts` *(generated by codegen)*
- Create: `web/src/shared/api/client.ts`
- Create: `web/src/shared/ui/CaduceusLogo.tsx`
- Create: `web/src/shared/ui/EmberRule.tsx`
- Create: `web/src/shared/ui/AppShell.tsx`
- Create: `web/src/features/system/hooks.ts`
- Create: `web/src/features/system/SystemOverview.tsx`
- Create: `web/src/features/system/InitForm.tsx`
- Test: `web/src/shared/api/client.test.ts`

**Interfaces:**
- Consumes: `contracts/openapi/openapi.yaml` via `openapi-typescript`
- Produces:
  - `web/src/shared/api/client.ts`: `apiClient` (typed `openapi-fetch` instance with `credentials: "include"`), `extractErrorMessage(err: unknown): string`
  - `web/src/shared/ui/`: `CaduceusLogo`, `EmberRule`, `AppShell` implementing `docs/DESIGN.md`
  - `web/out/` static export directory ready for Go `//go:embed`

- [ ] **Step 1: Initialize `web/` config files and install dependencies with Bun**

- `web/package.json`: scripts `dev` (`next dev --port 20127 --hostname 127.0.0.1`), `build` (`next build`), `typecheck` (`tsc --noEmit`), `gen:api` (`openapi-typescript ../contracts/openapi/openapi.yaml -o src/shared/api/schema.d.ts`), `test` (`bun test`).
- Dependencies: `next`, `react`, `react-dom`, `geist`, `@tanstack/react-query`, `openapi-fetch`, `lucide-react`.
- DevDependencies: `typescript`, `@types/react`, `@types/react-dom`, `@types/node`, `tailwindcss`, `@tailwindcss/postcss`, `openapi-typescript`.
- `web/next.config.ts`: `output: "export"`, `trailingSlash: true`, `images: { unoptimized: true }`, plus dev-only API proxy note (in dev mode Next.js `rewrites` only apply when `output` is not `"export"`, so use `const isStaticExport = process.env.NEXT_OUTPUT_EXPORT === "true"` with `output: isStaticExport ? "export" : undefined` and dev rewrites `/api/:path*`, `/v1/:path*`, `/v1beta/:path*` → `http://127.0.0.1:20128` when `isStaticExport` is false; `build` script sets `NEXT_OUTPUT_EXPORT=true`).
- `web/tsconfig.json`: `"strict": true`, `"noUncheckedIndexedAccess": true`.
- Run: `cd web && bun install && bun run gen:api`
- Expected: `web/src/shared/api/schema.d.ts` generated containing `/api/health`, `/api/version`, `/api/init`, `/api/shutdown`.

- [ ] **Step 2: Write failing frontend tests in `web/src/shared/api/client.test.ts`**

```ts
import { describe, expect, test } from "bun:test";
import { createApiClient, extractErrorMessage } from "./client";

describe("apiClient", () => {
  test("extractErrorMessage reads structured Go ErrorEnvelope", () => {
    expect(extractErrorMessage({ error: { code: "conflict", message: "Already initialized" } }))
      .toBe("Already initialized");
  });

  test("createApiClient parses /api/health response using custom fetch", async () => {
    // mock fetch returning { status: "ok", db: "ok", modules: ["system"], uptimeSeconds: 5 }
    // assert data.status === "ok" && data.modules.includes("system")
  });
});
```

- [ ] **Step 3: Run test to verify failure**

Run: `cd web && bun test`
Expected: FAIL (`Cannot find module './client'`).

- [ ] **Step 4: Implement typed API client, `DESIGN.md` UI primitives, and route shells**

- `globals.css`: define CSS variables for dark `#0a0a0a` / light `#FDFAF6`, primary `#E56A4A`, neutral `#6B7280`, Geist Sans & Mono variables, solid opaque surfaces (zero `backdrop-blur`), no emojis anywhere.
- `CaduceusLogo.tsx`: clean custom geometric SVG caduceus mark in ember orange.
- `EmberRule.tsx`: hairline gradient divider `bg-gradient-to-r from-transparent via-[#E56A4A] to-transparent h-px`.
- `AppShell.tsx`: opaque charcoal/warm-paper sidebar + header + main Monitor layout with nav slots prepared for Wave 1 modules.
- Pages:
  - `/` checks `/api/init` and routes cleanly to `/init` (if uninitialized) or `/dashboard`.
  - `/init` renders `InitForm` (Configure archetype) posting to `POST /api/init`.
  - `/login` renders split Configure archetype login shell (wired fully in Wave 1 `auth`).
  - `/dashboard` renders `SystemOverview` consuming live `GET /api/health` and `GET /api/version` via TanStack Query.

- [ ] **Step 5: Verify tests, strict typecheck, and static export build**

Run: `cd web && bun run gen:api && bun test && bun run typecheck && NEXT_OUTPUT_EXPORT=true bun run build && test -f out/index.html && test -f out/dashboard/index.html`
Expected: All tests pass, zero TS errors, `web/out/index.html` and `web/out/dashboard/index.html` exist.

- [ ] **Step 6: Commit**

```bash
git add web/
git commit -m "feat(web): scaffold Next.js TypeScript static-export shell and OpenAPI client"
```

---

### Task 10: Single-Binary Wiring (`cmd/omeg4mon-router`), `Makefile`, `Dockerfile` & End-to-End Verification

**Files:**
- Create: `cmd/omeg4mon-router/main.go`
- Create: `cmd/omeg4mon-router/ui_embed.go`
- Create: `cmd/omeg4mon-router/ui_dev.go`
- Create: `cmd/omeg4mon-router/main_test.go`
- Create: `Makefile`
- Create: `Dockerfile`

**Interfaces:**
- Consumes: all `internal/platform/*`, `internal/contracts`, `internal/modules/system`, and `web/out/` (staged into `cmd/omeg4mon-router/ui_dist/` during `make build`)
- Produces:
  - `cmd/omeg4mon-router/main.go`: `buildApp(ctx context.Context, cfg config.Config, build system.BuildInfo) (*App, error)` (testable entrypoint helper) + `main()` handling `SIGINT`/`SIGTERM`
  - Single static binary `./bin/omeg4mon-router` compiled with `CGO_ENABLED=0 go build -tags embedui`

- [ ] **Step 1: Write failing end-to-end app wiring test in `cmd/omeg4mon-router/main_test.go`**

```go
func TestAppBootstrapsMigrationsRoutesAndShutdown(t *testing.T) {
    // cfg with DataDir=t.TempDir(), Host="127.0.0.1", Port=0
    // app, err := buildApp(ctx, cfg, system.BuildInfo{Version: "test"})
    // 1) httptest GET /api/health -> 200 JSON with modules containing "system"
    // 2) httptest GET /api/init -> 200 {initialized: false}
    // 3) httptest GET /api/nonexistent -> 404 JSON
    // 4) app.Close(ctx) -> nil
}
```

- [ ] **Step 2: Run test to verify failure**

Run: `go test ./cmd/omeg4mon-router/...`
Expected: FAIL.

- [ ] **Step 3: Implement `main.go`, `ui_embed.go`, `ui_dev.go`, `Makefile`, and `Dockerfile`**

- `ui_embed.go` (`//go:build embedui`): `//go:embed all:ui_dist` returning sub-FS rooted at `ui_dist`.
- `ui_dev.go` (`//go:build !embedui`): returns `nil` FS and a lightweight dev handler telling developers to visit `http://127.0.0.1:20127` (or build with `make build`).
- `Makefile` targets:
  - `codegen`: runs `cd web && bun run gen:api`
  - `test`: runs `go test -race ./... && cd web && bun test && bun run typecheck`
  - `build-web`: runs `cd web && NEXT_OUTPUT_EXPORT=true bun run build`, syncs `web/out/` -> `cmd/omeg4mon-router/ui_dist/`
  - `build`: `codegen` + `build-web` + `CGO_ENABLED=0 go build -tags embedui -trimpath -ldflags="-s -w" -o bin/omeg4mon-router ./cmd/omeg4mon-router`
  - `dev`: prints/runs Go backend + Next.js dev server instructions cleanly
  - `cross-build`: builds `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64` with `CGO_ENABLED=0 -tags embedui`.
- `Dockerfile`: multi-stage (`oven/bun` builds `web/out` -> `golang:1.24-alpine` builds static `embedui` binary -> `gcr.io/distroless/static-debian12:nonroot`, exposes `20128`, sets `OMEG4MON_HOST=0.0.0.0`, `OMEG4MON_DATA_DIR=/data`).

- [ ] **Step 4: Run full verification and live single-binary smoke test**

Run:
```bash
make test && make build && (
  TMPDIR=$(mktemp -d)
  OMEG4MON_HOST=127.0.0.1 OMEG4MON_PORT=20199 OMEG4MON_DATA_DIR="$TMPDIR" ./bin/omeg4mon-router &
  PID=$!
  trap "kill $PID 2>/dev/null || true; rm -rf $TMPDIR" EXIT
  sleep 1
  curl -sf http://127.0.0.1:20199/api/health | grep -q '"status":"ok"'
  curl -sf http://127.0.0.1:20199/dashboard/ | grep -qi '<!DOCTYPE html>'
  STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:20199/api/does-not-exist)
  test "$STATUS" = "404"
)
```
Expected: All Go + Web tests pass, `bin/omeg4mon-router` builds with `CGO_ENABLED=0`, `/api/health` returns 200 JSON, `/dashboard/` serves embedded Next.js HTML, and `/api/does-not-exist` returns 404 JSON.

- [ ] **Step 5: Commit**

```bash
git add cmd/omeg4mon-router/ Makefile Dockerfile
git commit -m "feat(cmd): wire single-binary embedui build, Makefile, and multi-stage Dockerfile"
```

---

## Next Plan After Wave 0

Once Wave 0 is green and merged/verified on `migration/ts-go`, create **`docs/superpowers/plans/2026-11-22-wave1-core-gateway.md`** covering the 8 Wave 1 domain modules (`auth`, `settings`, `providers`, `keys`, `combos`, `gateway`, `usage`, `clitools`) and their `web/src/features/*` dashboard screens directly on top of Wave 0's verified platform contracts.
