// Origin / fetch-destination gate for state-changing API requests.
// Complements SameSite=lax cookies: even if a cookie policy ever loosens, a
// cross-site browser POST never reaches a handler. Non-browser clients (no
// Origin / no Sec-Fetch-Site headers) pass untouched.
// Next.js 16 `proxy` convention (formerly src/middleware.js); do not add a
// middleware.js alongside it.
import { NextResponse } from "next/server";

const MUTATING = new Set(["POST", "PUT", "PATCH", "DELETE"]);

// Sec-Fetch-Site values that alone prove the request came from our own page.
// Browsers forbid scripts from setting Sec-Fetch-*, so a foreign page cannot
// forge them. `same-site` is deliberately excluded: a sibling subdomain is
// same-site, so it still has to pass the Origin/host comparison below.
const BROWSER_ATTESTED = new Set(["same-origin", "none"]);

// Env vars that may hold the public origin when the dashboard runs behind a
// tunnel / reverse proxy that rewrites Host.
const ORIGIN_ENV_KEYS = [
  "BASE_URL",
  "PUBLIC_BASE_URL",
  "NEXT_PUBLIC_BASE_URL",
  "APP_URL",
  "SITE_URL",
  "NEXTAUTH_URL",
];

function requestHosts(request) {
  const hosts = new Set();
  for (const header of ["x-forwarded-host", "host"]) {
    const raw = request.headers.get(header);
    if (!raw) continue;
    // Proxy chains append: "public.host, internal.host".
    for (const part of raw.split(",")) {
      const value = part.trim();
      if (!value) continue;
      hosts.add(value);
      hosts.add(value.split(":")[0]);
    }
  }
  for (const key of ORIGIN_ENV_KEYS) {
    const raw = process.env[key];
    if (!raw) continue;
    try {
      const url = new URL(/^[a-z][a-z0-9+.-]*:\/\//i.test(raw) ? raw : `https://${raw}`);
      hosts.add(url.host);
      hosts.add(url.hostname);
    } catch {
      // ignore malformed env values
    }
  }
  return hosts;
}

export function proxy(request) {
  if (!MUTATING.has(request.method)) return NextResponse.next();

  // Auth login/logout and status routes need to be accessible from any client
  const pathname = request.nextUrl.pathname;
  if (pathname.startsWith("/api/auth/")) {
    return NextResponse.next();
  }

  const site = request.headers.get("sec-fetch-site");

  // Host comparison is unreliable behind tunnels/proxies that rewrite Host
  // (Origin = public hostname, Host = 127.0.0.1:<port>), so trust the browser.
  if (site && BROWSER_ATTESTED.has(site)) return NextResponse.next();

  const origin = request.headers.get("origin");
  if (origin) {
    let sameHost = false;
    try {
      const originUrl = new URL(origin);
      const hosts = requestHosts(request);
      sameHost = hosts.has(originUrl.host) || hosts.has(originUrl.hostname);
    } catch {
      sameHost = false;
    }
    if (!sameHost) {
      return NextResponse.json({ error: "Cross-origin request blocked" }, { status: 403 });
    }
  }

  if (site && site !== "same-site") {
    return NextResponse.json({ error: "Cross-site request blocked" }, { status: 403 });
  }

  return NextResponse.next();
}

export const config = {
  // Dashboard API surface only: always same-origin. /v1* is deliberately left
  // out — it serves cross-origin consumers (CORS ACAO:*) with API-key auth.
  matcher: ["/api/:path*"],
};
