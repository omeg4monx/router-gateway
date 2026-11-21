# OMEG4MON ROUTER

Local AI routing gateway with a web dashboard. It exposes one OpenAI- and Anthropic-compatible endpoint (`/v1`) and routes each request to the upstream provider you configure, translating request and response formats on the way.

- **Gateway**: Next.js 16 server (`custom-server.js`) on port `20128`.
- **Dashboard**: `http://localhost:20128/dashboard`.
- **API base**: `http://localhost:20128/v1`.
- **CLI launcher**: the `omeg4mon-router` package in `cli/`, which starts the server, can run it from the system tray and writes configs for coding tools.

---

## Features

- **One endpoint, many providers.** About 146 provider definitions (`open-sse/providers/registry/`): OpenAI-compatible APIs through a generic executor, and dedicated executors for non-standard upstreams.
- **Format translation.** OpenAI Chat, Anthropic Messages, OpenAI Responses and Gemini formats. Translators pivot through the OpenAI format, with direct routes for lossy pairs.
- **Fallback.**
  - **Combos** try a list of models in order.
  - **Multi-account rotation** applies per-error cooldowns.
  - **Automatic token refresh** renews OAuth credentials.
  - **Request-scoped errors** (malformed request, unsupported parameter) are returned instead of burning the whole chain.
- **API keys** with per-key permissions, allowed models and usage tracking.
- **Usage and cost dashboard**, request logs and provider limits.
- **CLI tool setup.** Writes configs for Claude Code, Codex, OpenCode, Cline, Kilo, Droid, Crush, Hermes, OpenClaw, Grok Build, OMP and others.
- **Optional components:**
  - Cloudflare quick tunnel or Tailscale for remote access.
  - Proxy pools.
  - Local MITM proxy for IDE traffic.
  - Headroom token compression.
  - SAML/OIDC login.

---

## Requirements

- Node.js 22 (see `.nvmrc`). The CLI accepts Node ≥ 18.
- npm (Bun also works: `dev:bun`, `build:bun`, `start:bun`).
- Optional: Docker, Go 1.26+ for the experimental Go core.

---

## Quick start (from source)

```bash
git clone https://github.com/omeg4monx/router-gateway.git
cd router-gateway
cp .env.example .env        # set JWT_SECRET and INITIAL_PASSWORD
npm install
npm run build
npm run start               # http://localhost:20128/dashboard
```

To develop with hot reload, run `npm run dev`. It serves on `127.0.0.1:20128` only.

> `next dev` does **not** load `custom-server.js`, so the dashboard auth guard is inactive in dev mode. Never expose a dev server to a network.

### First login

- Without `INITIAL_PASSWORD`, the default password is `omeg4mon123`. It is accepted only from a trusted network (loopback, LAN, Docker or Tailscale) unless `ALLOW_REMOTE_DEFAULT_LOGIN=true`.
- To change the password for the first time, enter the current initial password.

---

## Docker

```bash
cp .env.example .env
docker compose up -d        # gateway on :20128, Headroom on 127.0.0.1:8787
```

You can also run `./start.sh`, which rebuilds and runs a single container.

- Data lives in the `omeg4mon-router-data` volume (override with `OMEG4MON_ROUTER_VOLUME`).
- The container listens on all interfaces. Put it behind a reverse proxy or firewall before exposing it.

---

## Linux server install

```bash
sudo ./install.sh --port 20128 --password '<strong-password>'
```

`install.sh` does the following:
- Writes `.env` with generated secrets (mode `600`).
- Builds the app.
- Installs a systemd service named `omeg4mon-router`. The unit runs as root; review it before use.

Run `./install.sh --help` for all flags.

---

## CLI launcher

```bash
npm run cli:pack            # builds and packs the cli/ package
npm i -g ./omeg4mon-router-*.tgz   # pack:cli writes the tarball to the repo root
omeg4mon-router             # starts the server and opens the dashboard
```

| Option | Description |
|---|---|
| `-p, --port <port>` | Port (default `20128`) |
| `-H, --host <host>` | Bind address (default `127.0.0.1`; use `0.0.0.0` for LAN access) |
| `-n, --no-browser` | Don't open the browser |
| `-l, --log` | Show server logs |
| `-t, --tray` | Run in the system tray |
| `--skip-update` | Skip the update check |

**Commands:**
- `omeg4mon-router connect http://<host>:20128` configures local coding tools against a remote gateway. That gateway must be started with `--host 0.0.0.0`.
- `omeg4mon-router xai video --help` generates video through the gateway.

---

## Using the API

Create an API key in **Dashboard → Keys**, then call the gateway like any OpenAI or Anthropic endpoint:

```bash
curl http://localhost:20128/v1/chat/completions \
  -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{"model":"<provider>/<model>","messages":[{"role":"user","content":"Hello"}]}'
```

| Endpoint | Format |
|---|---|
| `POST /v1/chat/completions` | OpenAI Chat Completions |
| `POST /v1/messages` | Anthropic Messages |
| `POST /v1/responses` | OpenAI Responses |
| `GET /v1/models` | Model list |
| `/v1/embeddings`, `/v1/images`, `/v1/audio`, `/v1/videos` | Media endpoints (provider-dependent) |
| `/v1beta/models/*` | Gemini-style |

The API key can also be sent as `x-api-key`, `x-goog-api-key` or `?key=`.

Loopback requests are accepted without a key unless **Require API key** is enabled. Remote requests always need a valid key.

---

## Configuration

Every variable is documented in `.env.example`. The most important ones:

| Variable | Purpose |
|---|---|
| `JWT_SECRET` | Signs dashboard sessions (generated into `DATA_DIR/jwt-secret` if unset) |
| `INITIAL_PASSWORD` | Dashboard password until you set one |
| `DATA_DIR` | Data location (default `~/.omeg4mon-router`, `%APPDATA%\omeg4mon-router` on Windows) |
| `PORT` | Listen port (default `20128`) |
| `OMEG4MON_ROUTER_VAULT_KEY` | At-rest key for stored credentials. If unset, a random key is created in `DATA_DIR/vault.key` |
| `TUNNEL_WORKER_URL` | Optional short-link worker for tunnels. Leave unset unless you run the worker yourself |
| `BASE_URL` | Public URL. **Required for SAML** |

### Data and backups

- State is stored in SQLite at `DATA_DIR/db/data.sqlite`.
- The driver is chosen automatically, in this order: `bun:sqlite`, `better-sqlite3`, `node:sqlite`, `sql.js`.
- Back up `DATA_DIR/vault.key` together with the database. Without it, stored provider credentials cannot be decrypted.

---

## Security notes

- The gateway binds to loopback by default. Before exposing it, set strong `JWT_SECRET` and `INITIAL_PASSWORD` values and enable **Require API key**.
- Configured proxies fail closed: a request never falls back to a direct connection when its proxy fails.
- Telegram control needs a configured user ID and a webhook registered from the dashboard. The webhook secret is verified on every update.
- The MITM component installs a local root CA and edits the hosts file. Enable it only on machines you control.
- Using consumer web sessions or IDE identities of third-party services through this gateway may violate their terms of service. You are responsible for how providers are used.

---

## Project layout

```
src/            Next.js app: dashboard, /api routes, auth guard, persistence, MITM
open-sse/       Provider-agnostic engine: translators, executors, fallback, token refresh
cli/            omeg4mon-router CLI launcher (separate package)
cmd/, internal/ Experimental Go core (not feature-complete; see docs/ROADMAP_GO_MIGRATION.md)
tests/          Vitest suite (independent package)
```

- Request flow: `/v1/*` → `src/sse/handlers/chat.js` → `open-sse/handlers/chatCore.js` → executor and translator → upstream.
- See `docs/ARCHITECTURE.md` and `open-sse/AGENTS.md` before changing the engine.

---

## Tests

```bash
npm install
cd tests && npm install
npx vitest run                      # whole suite
npx vitest run unit/<file>.test.js  # one file
```

Some tests are known to fail on a plain checkout: `tests/__baseline__/` holds the regression baselines, and `real/*.real.test.js` calls live providers and needs credentials. Compare against a baseline run instead of expecting an all-green suite.

---

## License

MIT. See [LICENSE](./LICENSE).
