# Docker

Run OMEG4MON ROUTER in a container. CI publishes a multi-platform image (`linux/amd64` and `linux/arm64`) to GHCR at `ghcr.io/omeg4monx/router-gateway`. You can also build the image yourself from this repository.

---

# 👤 For Users

## Quick start (Docker Compose)

```bash
cp .env.example .env        # set JWT_SECRET and INITIAL_PASSWORD
docker compose up -d
```

This starts two services:
- **`omeg4mon-router`** on port `20128`, with data in the `omeg4mon-router-data` volume.
- **Headroom** on `127.0.0.1:8787`.

Open http://localhost:20128/dashboard.

`./start.sh` does the same for a single container. It rebuilds the image and replaces any running `omeg4mon-router` container.

## Quick start (docker run)

```bash
docker run -d \
  -p 20128:20128 \
  -v omeg4mon-router-data:/app/data \
  --env-file .env \
  --name omeg4mon-router \
  ghcr.io/omeg4monx/router-gateway:latest
```

The app listens on port `20128` inside the container (`HOSTNAME=0.0.0.0`).

To keep it reachable only from the host, publish the port on loopback: `-p 127.0.0.1:20128:20128`.

## Manage container

```bash
docker logs -f omeg4mon-router        # view logs
docker stop omeg4mon-router           # stop
docker start omeg4mon-router          # start again
docker rm -f omeg4mon-router          # remove (the data volume is kept)
```

## Data persistence

The image sets `DATA_DIR=/app/data`. Mount a volume or host directory there:

```bash
-v omeg4mon-router-data:/app/data
# or a host path:
-v "$HOME/.omeg4mon-router:/app/data"
```

Outside a container, the default `DATA_DIR` is `~/.omeg4mon-router/` on macOS and Linux, and `%APPDATA%\omeg4mon-router\` on Windows.

Data layout under `$DATA_DIR/`:

```text
$DATA_DIR/
├── db/
│   ├── data.sqlite       # main SQLite database
│   └── backups/          # safety backups taken before schema changes
├── vault.key             # at-rest key for stored credentials (keep with backups)
├── jwt-secret            # session signing key when JWT_SECRET is unset
└── ...                   # certs, logs, runtime configs
```

Back up `vault.key` together with `data.sqlite`. Without it, stored provider credentials cannot be decrypted. Instead of the file, you can supply the key through `OMEG4MON_ROUTER_VAULT_KEY`.

The entrypoint fixes ownership of the data directory, then drops privileges and runs the server as the `node` user.

## Optional env vars

```bash
docker run -d \
  -p 20128:20128 \
  -v omeg4mon-router-data:/app/data \
  -e JWT_SECRET=<long-random-secret> \
  -e INITIAL_PASSWORD=<strong-password> \
  -e OMEG4MON_ROUTER_VAULT_KEY=<long-random-secret> \
  -e DEBUG=true \
  --name omeg4mon-router \
  ghcr.io/omeg4monx/router-gateway:latest
```

`.env.example` documents every variable.

Without `INITIAL_PASSWORD`, the first-login password is `omeg4mon123`. It is accepted only from trusted networks (loopback, LAN, Docker bridge, Tailscale).

## Optional Headroom sidecar

The image does not bundle Python or Headroom. `docker-compose.yml` already runs Headroom as a separate service and points OMEG4MON ROUTER at it:

```yaml
services:
  omeg4mon-router:
    image: ghcr.io/omeg4monx/router-gateway:latest
    ports:
      - "20128:20128"
    volumes:
      - omeg4mon-router-data:/app/data
    environment:
      DATA_DIR: /app/data
      HEADROOM_URL: http://headroom:8787
    depends_on:
      - headroom

  headroom:
    image: ghcr.io/chopratejas/headroom:latest
    ports:
      - "127.0.0.1:8787:8787"

volumes:
  omeg4mon-router-data:
```

To enable it:
1. In the dashboard, open `Endpoint` → `Token Saver` → `Headroom`.
2. Confirm the URL is `http://headroom:8787`.
3. Recheck the status, then enable Headroom.

If Headroom runs on the Docker host instead of as a sidecar:
- On macOS and Windows, use `http://host.docker.internal:8787`.
- On Linux, add `--add-host=host.docker.internal:host-gateway`, or the equivalent compose `extra_hosts` entry.

## Update to latest

```bash
docker pull ghcr.io/omeg4monx/router-gateway:latest
docker rm -f omeg4mon-router
# re-run the quick start command; the data volume is reused
```

To pin a specific version, use a numbered tag (`ghcr.io/omeg4monx/router-gateway:X.Y.Z`). For an immutable deployment, pin the digest instead (`ghcr.io/omeg4monx/router-gateway@sha256:<digest>`).

---

# 🛠 For Developers

## Build image locally

```bash
docker build -t omeg4mon-router .

docker run --rm -p 20128:20128 \
  -v omeg4mon-router-data:/app/data \
  omeg4mon-router
```

The Dockerfile uses the official Alpine and npm registries by default. Regional mirrors can be supplied when needed:

```bash
docker build \
  --build-arg ALPINE_MIRROR=mirrors.aliyun.com \
  --build-arg NPM_REGISTRY=https://registry.npmmirror.com/ \
  -t omeg4mon-router .
```

`docker build -f Dockerfile.gateway -t omeg4mon-router-go .` builds the experimental Go core. That core is not feature-complete and is not used by `docker-compose.yml`.

## Publish (CI)

`.github/workflows/docker-publish.yml` is tag-driven. Pushing a semver git tag `vX.Y.Z` (or a prerelease such as `vX.Y.Z-rc.1`) runs this sequence:

1. Builds `linux/amd64` and `linux/arm64` on native runners.
2. Health-checks each platform image with `/api/health`.
3. Assembles and verifies the multi-platform manifest.
4. Publishes `ghcr.io/omeg4monx/router-gateway:X.Y.Z`. A stable tag also moves `:latest`.

```bash
git tag vX.Y.Z && git push origin vX.Y.Z
```

**Tag rules:**
- The git tag and both `package.json` versions must match exactly.
- SemVer build metadata (`v1.2.3+build.7`) is rejected, because `+` is not valid in a Docker tag.
- Prereleases publish only their numbered tag, unless `promote_latest` is enabled on a manual run.

**Republishing an existing tag:** run the workflow manually and set `release_tag`. `promote_latest` is an explicit opt-in for moving `latest`, for example during a rollback.

**Publishing:**
- GHCR uses the workflow's `GITHUB_TOKEN` with package write permission.
- The workflow publishes to Docker Hub only for the upstream repository it was inherited from, so this repository publishes to GHCR only.
- The optional repository variables `ALPINE_MIRROR` and `NPM_REGISTRY` override the package mirrors used in CI builds.
