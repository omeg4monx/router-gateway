# `@omeg4mon/router` CLI Launcher

Thin, zero-dependency multi-platform npm/Bun launcher for the single-binary `omeg4mon-router` server.

## Binary Resolution Order

1. `OMEG4MON_ROUTER_BIN` environment variable
2. `cli/vendor/omeg4mon-router-<os>-<arch>[.exe]`
3. `../bin/omeg4mon-router[.exe]` (local `make build` output)
4. `../dist/omeg4mon-router-<os>-<arch>[.exe]` (`make cross-build` output)
5. `omeg4mon-router` on `PATH`

## Usage

```bash
npx @omeg4mon/router --host 127.0.0.1 --port 20128
```
