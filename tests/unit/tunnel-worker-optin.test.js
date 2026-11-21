import { describe, it, expect, vi, afterEach } from "vitest";

describe("tunnel short-link worker is opt-in", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.resetModules();
  });

  it("has no worker and no short public URL unless TUNNEL_WORKER_URL is set", async () => {
    vi.stubEnv("TUNNEL_WORKER_URL", "");
    const { WORKER_URL, workerPublicUrl } = await import("../../src/lib/tunnel/cloudflare/config.js");
    expect(WORKER_URL).toBe("");
    expect(workerPublicUrl("abc123")).toBe("");
  });

  it("builds r<shortId> on the configured worker's host", async () => {
    vi.stubEnv("TUNNEL_WORKER_URL", "https://links.example.com/");
    const { workerPublicUrl } = await import("../../src/lib/tunnel/cloudflare/config.js");
    expect(workerPublicUrl("abc123")).toBe("https://rabc123.links.example.com");
  });
});
