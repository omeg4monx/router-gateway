import { describe, it, expect, vi, beforeEach } from "vitest";

// proxyFetch captures globalThis.fetch at import as its "direct" egress.
const direct = vi.hoisted(() => {
  const fn = vi.fn();
  globalThis.fetch = fn;
  return fn;
});

vi.mock("../../open-sse/utils/proxyDispatcher.js", () => ({
  createProxyDispatcher: () => ({ proxied: true }),
}));

const { proxyAwareFetch } = await import("../../open-sse/utils/proxyFetch.js");

const proxyOptions = { connectionProxyEnabled: true, connectionProxyUrl: "http://127.0.0.1:9" };

describe("proxyAwareFetch never bypasses a configured proxy", () => {
  beforeEach(() => direct.mockReset());

  const target = "https://api.example.com/v1/x";
  const callsTo = (url) => direct.mock.calls.filter(([u]) => String(u) === url);

  it("throws instead of retrying direct when the proxy connection fails", async () => {
    direct.mockImplementation(async (_url, init) => {
      if (init?.dispatcher?.proxied) throw new Error("ECONNREFUSED");
      return new Response("direct egress");
    });
    await expect(proxyAwareFetch(target, {}, proxyOptions)).rejects.toThrow(/Proxy required but failed/);
    expect(callsTo(target)).toHaveLength(1);
    expect(callsTo(target)[0][1].dispatcher).toEqual({ proxied: true });
  });

  it("returns a FreeTierError rejection instead of replaying it on the direct egress", async () => {
    direct.mockImplementation(async (_url, init) =>
      init?.dispatcher?.proxied ? new Response("FreeTierError: blocked", { status: 403 }) : new Response("direct egress"));
    const res = await proxyAwareFetch(target, {}, proxyOptions);
    expect(res.status).toBe(403);
    expect(callsTo(target)).toHaveLength(1);
  });
});
