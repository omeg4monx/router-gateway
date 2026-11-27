import { describe, expect, it } from "bun:test";
import { buildServiceWorkerSource } from "../../../pwaPlugin";
import { classifyFetchRequest } from "./swStrategy";

const ORIGIN = "https://router.local";

describe("Service Worker fetch classification", () => {
  it("ignores cross-origin requests", () => {
    expect(
      classifyFetchRequest(new URL("https://cdn.example.com/app.js"), "GET", "cors", ORIGIN),
    ).toBe("ignore");
  });

  it("keeps non-GET requests strict network-only", () => {
    expect(classifyFetchRequest(new URL(`${ORIGIN}/api/providers`), "POST", "cors", ORIGIN)).toBe(
      "network-only",
    );
    expect(
      classifyFetchRequest(new URL(`${ORIGIN}/dashboard/`), "DELETE", "same-origin", ORIGIN),
    ).toBe("network-only");
  });

  it("keeps API, LLM gateway, and sw.js paths strict network-only", () => {
    const networkOnlyPaths = [
      "/api/auth/session",
      "/api/system/health",
      "/v1/chat/completions",
      "/v1beta/models",
      "/codex/responses",
      "/responses",
      "/models",
      "/sw.js",
    ];
    for (const pathname of networkOnlyPaths) {
      expect(classifyFetchRequest(new URL(`${ORIGIN}${pathname}`), "GET", "cors", ORIGIN)).toBe(
        "network-only",
      );
    }
  });

  it("classifies SPA navigation requests for network-first with index.html fallback", () => {
    expect(classifyFetchRequest(new URL(`${ORIGIN}/dashboard/`), "GET", "navigate", ORIGIN)).toBe(
      "navigate",
    );
    expect(
      classifyFetchRequest(
        new URL(`${ORIGIN}/dashboard/providers/detail`),
        "GET",
        "navigate",
        ORIGIN,
      ),
    ).toBe("navigate");
    expect(classifyFetchRequest(new URL(`${ORIGIN}/login/`), "GET", "navigate", ORIGIN)).toBe(
      "navigate",
    );
  });

  it("classifies static hashed assets and icons for cache-first handling", () => {
    expect(
      classifyFetchRequest(new URL(`${ORIGIN}/assets/index-abc123.js`), "GET", "cors", ORIGIN),
    ).toBe("asset");
    expect(
      classifyFetchRequest(new URL(`${ORIGIN}/icons/pwa-192x192.png`), "GET", "no-cors", ORIGIN),
    ).toBe("asset");
    expect(
      classifyFetchRequest(new URL(`${ORIGIN}/manifest.webmanifest`), "GET", "cors", ORIGIN),
    ).toBe("asset");
  });
});

describe("buildServiceWorkerSource", () => {
  it("emits cache version, precache list, SKIP_WAITING listener, old cache cleanup, and API bypass", () => {
    const src = buildServiceWorkerSource("omeg4mon-router-v1-deadbeef", [
      "/index.html",
      "/manifest.webmanifest",
      "/assets/index-abc.js",
    ]);

    expect(src).toContain('const CACHE_NAME = "omeg4mon-router-v1-deadbeef";');
    expect(src).toContain('"/index.html"');
    expect(src).toContain('"/assets/index-abc.js"');
    expect(src).toContain("SKIP_WAITING");
    expect(src).toContain("omeg4mon-router-");
    expect(src).toContain('"/api"');
    expect(src).toContain('"/v1"');
    expect(src).toContain('"/v1beta"');
  });
});
