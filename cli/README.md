# OMEG4MON ROUTER - FREE AI Router & Token Saver

**Never stop coding. Save 20-40% tokens with RTK + auto-fallback to FREE & cheap AI models.**

**Connect All AI Code Tools (Claude Code, Cursor, Antigravity, Copilot, Codex, Gemini, OpenCode, Cline, OpenClaw...) to 40+ AI Providers & 100+ Models.**



[Project documentation](../docs/)

---

## 🤔 Why OMEG4MON ROUTER?

**Stop wasting money, tokens and hitting limits:**

- ❌ Subscription quota expires unused every month
- ❌ Rate limits stop you mid-coding
- ❌ Tool outputs (git diff, grep, ls...) burn tokens fast
- ❌ Expensive APIs ($20-50/month per provider)

**OMEG4MON ROUTER solves this:**

- ✅ **RTK Token Saver** - Auto-compress tool_result, save 20-40% tokens
- ✅ **Maximize subscriptions** - Track quota, use every bit before reset
- ✅ **Auto fallback** - Subscription → Cheap → Free, zero downtime
- ✅ **Multi-account** - Round-robin between accounts per provider
- ✅ **Universal** - Works with any OpenAI/Claude-compatible CLI

---

## ⚡ Quick Start

**Option 1 — CLI package (from the repository root):**

```bash
bun install
bun run cli:pack
# Install the generated .tgz from the repository root using npm install -g <tarball>.
omeg4mon-router
```

**Option 2 — Docker (server/VPS):**

```bash
docker compose up --build -d
```

On upgrades, stop the previous container without removing its data volume.
Set `OMEG4MON_ROUTER_VOLUME` if it used a non-default volume.
For native installations, keep the previous `DATA_DIR` and vault key.
For SAML, keep your existing service-provider entity ID in the dashboard before upgrading.


🎉 Dashboard opens at `http://localhost:20128`

**2. Connect a FREE provider (no signup needed):**

Dashboard → Providers → Connect **Kiro AI** (free Claude unlimited) or **OpenCode Free** (no auth) → Done!

**3. Use in your CLI tool:**

```
Claude Code/Codex/OpenClaw/Cursor/Cline Settings:
  Endpoint: http://localhost:20128/v1
  API Key:  [copy from dashboard]
  Model:    kr/claude-sonnet-4.5
```

That's it! Start coding with FREE AI models.

---

## 🚀 CLI Options

```bash
omeg4mon-router                    # Start with default settings
omeg4mon-router --port 8080        # Custom port
omeg4mon-router --no-browser       # Don't open browser
omeg4mon-router --skip-update      # Skip auto-update check
omeg4mon-router --help             # Show all options
```

**Dashboard**: `http://localhost:20128/dashboard`

---

## 🔌 Connect to a Remote OMEG4MON ROUTER

Already running OMEG4MON ROUTER on another machine (e.g. a team server on your LAN)? Point this machine's CLI tools at it — no local server is started. The server binds `127.0.0.1` by default, so start it with `--host 0.0.0.0` to accept LAN connections:

```bash
omeg4mon-router connect http://<server-host>:20128                       # pick tools interactively
omeg4mon-router connect http://<server-host>:20128 --tools claude,codex  # or choose up front
omeg4mon-router connect --reset --tools claude,codex                     # undo
```

It logs in with the dashboard password (hidden prompt), reuses or creates an API key named `cli-<hostname>`, and writes each tool's config (backing up the original once as `*.bak-omeg4mon-router`).

Supported: `claude`, `codex`, `opencode`, `droid`, `crush`, `kilo`, `cline`, or `all`. Other options: `--model`, `--opus/--sonnet/--haiku/--fable`, `--api-key`, `--key-name`, `--print-env`. See `omeg4mon-router connect --help`.

> ⚠️ Over plain `http://` the password and API key are sent unencrypted — use a trusted LAN/VPN or put HTTPS in front. The API key is stored in each tool's config file.

---

## 🛠️ Supported CLI Tools

Claude-Code • OpenClaw • Codex • OpenCode • Cursor • Antigravity • Cline • Continue • Droid • Roo • Copilot • Kilo Code • Gemini CLI • Qwen Code • iFlow • Crush • Crusher • Aider

Any tool supporting OpenAI/Claude-compatible API works.

---

## 💾 Data Location

- **macOS/Linux**: `~/.omeg4mon-router/db/data.sqlite`
- **Windows**: `%APPDATA%/omeg4mon-router/db/data.sqlite`
- **Docker**: `/app/data/db/data.sqlite` (mount `$HOME/.omeg4mon-router` to persist)

---

## 📚 Documentation

Full docs, advanced setup, video tutorials & development guide:

- **GitHub**: https://github.com/omeg4monx/router-gateway
- **Documentation**: https://github.com/omeg4monx/router-gateway/tree/main/docs

---

## 🙏 Acknowledgments

- **[CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI)** - Original Go implementation

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
