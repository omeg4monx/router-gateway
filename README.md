# OMEG4MON Router

Single-binary modular LLM gateway and control plane built with **Go (`net/http` + `chi/v5`, zero-CGO SQLite WAL)** and **Vite + React 19 (TypeScript, static SPA embedded via `//go:embed`)**.

## Architecture

- **Single self-contained binary:** Serves OpenAI/Anthropic/Gemini-compatible gateway routes (`/v1/*`, `/v1beta/*`), admin management API (`/api/*`), and the embedded Vite + React static dashboard (`web/out/`).
- **27 isolated vertical-slice modules:** Each feature domain lives in `internal/modules/<name>/` (with its own SQLite migrations and OpenAPI 3.1 contract in `contracts/openapi/modules/`) and mirrors `web/src/features/<name>/`. Cross-module calls use narrow interfaces in `internal/contracts/` and `internal/platform/eventbus`.
- **Zero-CGO persistence & secret vault:** Pure-Go `modernc.org/sqlite` with split 1-writer / 8-reader WAL pools and AES-256-GCM credential encryption backed by `OMEG4MON_MASTER_KEY` or `<DATA_DIR>/master.key` (`0600`).

## Quick Start

### 1. Build & Run Single Binary

Requires Go `1.26+` and Bun `1.4+`:

```bash
make build
./bin/omeg4mon-router
```

Open `http://127.0.0.1:20128` to complete first-run setup (`/init/`) and open the dashboard.

### 2. Multi-Platform CLI Launcher (`@omeg4mon/router`)

```bash
npx @omeg4mon/router --host 127.0.0.1 --port 20128
```

Binary resolution order: `OMEG4MON_ROUTER_BIN` → `cli/vendor/` → `../bin/` → `../dist/` → `PATH`.

### 3. Docker (Distroless Multi-Stage)

```bash
docker build -t omeg4mon-router .
docker run -d --name omeg4mon-router \
  -p 20128:20128 \
  -v omeg4mon-data:/data \
  -e OMEG4MON_ADMIN_PASSWORD="change-me" \
  omeg4mon-router
```

> Non-loopback binds (`0.0.0.0`) enforce `require_login=true` and disable `/api/shutdown` unless `OMEG4MON_ALLOW_INSECURE_REMOTE=true` is explicitly set.

### 4. Local Development (Live Reload)

```bash
# Terminal 1 — Go API backend (:20128)
go run ./cmd/omeg4mon-router

# Terminal 2 — Vite dev UI (:20127, proxying /api, /v1, /v1beta)
cd web && bun run dev
```

## Configuration

| Environment Variable | CLI Flag | Default | Description |
| --- | --- | --- | --- |
| `OMEG4MON_HOST` | `--host` | `127.0.0.1` | Bind address |
| `OMEG4MON_PORT` | `--port` | `20128` | Listen port (`1–65535`) |
| `OMEG4MON_DATA_DIR` | `--data-dir` | `~/.omeg4mon-router` | SQLite DB and vault directory |
| `OMEG4MON_LOG_LEVEL` | `--log-level` | `info` | Structured `slog` level |
| `OMEG4MON_MASTER_KEY` | — | auto-generated | 32-byte hex AES-256-GCM vault key |
| `OMEG4MON_ADMIN_PASSWORD` | — | — | Optional headless first-run admin password |
| `OMEG4MON_API_TIMEOUT` | — | `30s` | Non-streaming `/api/*` request timeout |
| `OMEG4MON_DISABLED_MODULES` | — | — | Comma-separated optional modules to disable |
| `OMEG4MON_ALLOW_INSECURE_REMOTE` | — | `false` | Override remote-bind security guard |

## Modules

- **Core (always enabled):** `system`, `auth`, `settings`, `providers`, `keys`, `combos`, `gateway`, `usage`
- **Optional (toggleable via `OMEG4MON_DISABLED_MODULES`):** `clitools`, `proxypools`, `translator`, `modeleditor`, `tokensaver`, `promptinjection`, `security`, `notifications`, `teams`, `vouchers`, `media`, `tunnel`, `mitm`, `pxpipe`, `telegram`, `arena`, `plugins`, `cloudsync`, `importer`

## Verification & Release

```bash
make check        # Go fmt/lint/race/zero-CGO + Web biome/tests/typecheck + CLI tests
make build        # OpenAPI codegen + Vite SPA static build + single binary (bin/omeg4mon-router)
make smoke        # Live single-binary health, UI, 404, CLI launcher, and idle RSS check
make cross-build  # 6 zero-CGO binaries (linux/darwin/windows x amd64/arm64) + dist/SHA256SUMS.txt
```

Tagging `v*` triggers `.github/workflows/release.yml`, publishing all 6 verified binaries, `SHA256SUMS.txt`, and multi-arch `linux/amd64,linux/arm64` images to GHCR.

## License

MIT — Copyright (c) 2026 OMEG4MON. See [LICENSE](LICENSE).
