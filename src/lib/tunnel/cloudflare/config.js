// Cloudflare quick tunnel: DNS propagates fast, short timeouts OK
export const HEALTH_CHECK = {
  intervalMs: 2000,
  timeoutMs: 60000,
  fetchTimeoutMs: 5000,
  dnsTimeoutMs: 2000,
};

// Optional short-link worker. Its /api/tunnel/register is unauthenticated, so anyone
// who sees a public r<shortId> URL can re-point it; opt in only for a worker you run.
export const WORKER_URL = (process.env.TUNNEL_WORKER_URL || "").replace(/\/+$/, "");

// Public short URL for a shortId on the configured worker, or "" when none is set.
export function workerPublicUrl(shortId) {
  if (!WORKER_URL || !shortId) return "";
  return `https://r${shortId}.${new URL(WORKER_URL).host}`;
}
