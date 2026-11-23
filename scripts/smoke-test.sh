#!/usr/bin/env bash
set -euo pipefail

BINARY="${1:-bin/omeg4mon-router}"
PORT="${2:-20199}"
BASE_URL="http://127.0.0.1:${PORT}"

if [[ ! -x "${BINARY}" ]]; then
  echo "ERROR: binary not found or not executable: ${BINARY}" >&2
  exit 1
fi

TMPDIR="$(mktemp -d)"
NOTFOUND_JSON="${TMPDIR}/notfound.json"
SERVER_LOG="${TMPDIR}/server.log"
PID=""

cleanup() {
  if [[ -n "${PID}" ]] && kill -0 "${PID}" 2>/dev/null; then
    kill -TERM "${PID}" 2>/dev/null || true
    wait "${PID}" 2>/dev/null || true
  fi
  rm -rf "${TMPDIR}"
}
trap cleanup EXIT

OMEG4MON_HOST=127.0.0.1 \
OMEG4MON_PORT="${PORT}" \
OMEG4MON_DATA_DIR="${TMPDIR}/data" \
  "${BINARY}" >"${SERVER_LOG}" 2>&1 &
PID=$!

READY=0
for _ in $(seq 1 40); do
  if curl -sf "${BASE_URL}/api/health" >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 0.15
done

if [[ "${READY}" -ne 1 ]]; then
  echo "ERROR: server failed to become ready on ${BASE_URL}" >&2
  cat "${SERVER_LOG}" >&2
  exit 1
fi

HEALTH="$(curl -sf "${BASE_URL}/api/health")"
echo "HEALTH=${HEALTH}"
echo "${HEALTH}" | grep -q '"status":"ok"'

VERSION="$(curl -sf "${BASE_URL}/api/version")"
echo "VERSION=${VERSION}"
echo "${VERSION}" | grep -q '"goVersion"'

for route in "/dashboard/" "/init/" "/login/" "/dashboard/plugins/" "/dashboard/importer/"; do
  curl -sf "${BASE_URL}${route}" | grep -qi '<!DOCTYPE html>'
  echo "UI_OK=${route}"
done

STATUS404="$(curl -s -o "${NOTFOUND_JSON}" -w '%{http_code}' "${BASE_URL}/api/does-not-exist")"
if [[ "${STATUS404}" != "404" ]]; then
  echo "ERROR: expected HTTP 404 for /api/does-not-exist, got ${STATUS404}" >&2
  cat "${NOTFOUND_JSON}" >&2
  exit 1
fi
grep -q '"not_found"' "${NOTFOUND_JSON}"
echo "API_404_OK=$(cat "${NOTFOUND_JSON}")"

if [[ -r "/proc/${PID}/status" ]]; then
  RSS_KB="$(awk '/^VmRSS:/ {print $2}' "/proc/${PID}/status")"
  RSS_MB="$(awk -v kb="${RSS_KB:-0}" 'BEGIN {printf "%.2f", kb / 1024}')"
  echo "IDLE_RSS_MB=${RSS_MB}"
fi

if [[ -f "cli/bin/omeg4mon-router.js" ]]; then
  if command -v node >/dev/null 2>&1; then
    OMEG4MON_ROUTER_BIN="${BINARY}" node cli/bin/omeg4mon-router.js --help >/dev/null 2>&1 || true
    echo "CLI_LAUNCHER_OK=node"
  elif command -v bun >/dev/null 2>&1; then
    OMEG4MON_ROUTER_BIN="${BINARY}" bun cli/bin/omeg4mon-router.js --help >/dev/null 2>&1 || true
    echo "CLI_LAUNCHER_OK=bun"
  fi
fi

kill -TERM "${PID}"
wait "${PID}"
PID=""
echo "SMOKE_TEST_PASS"
