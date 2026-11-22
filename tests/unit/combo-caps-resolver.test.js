import { describe, expect, it } from "vitest";

import { aggregateComboCapabilities, getCapabilitiesForModel } from "../../open-sse/providers/capabilities.js";

// A combo's context window is the smallest of its members (a request must fit whichever
// member serves it); maxOutput = max. Resolving those members needs the synced model catalog, which
// is server-only (it reads a file), so the browser bundle falls back to the
// generic patterns. The dashboard computed its badges there and under-reported:
// /v1/models and pi-settings (both server-side) said 1M while the badge said 200k.
//
// resolveCaps lets a caller hand in the server's answer. It must only override
// what it carries - the local tables still own tools/pdf/audio/video/thinking*.
const GLM53_FED = { vision: true, search: false, reasoning: true, contextWindow: 500_000, maxOutput: 131_072 };
const DEEPSEEK_LOCAL_CTX = getCapabilitiesForModel(null, "deepseek-v4.1-flash").contextWindow;

describe("aggregateComboCapabilities: resolveCaps override", () => {
  const models = ["glm-cn/glm-5.3", "deepseek-v4.1-flash"];

  it("falls back to the local tables without a resolver", () => {
    const caps = aggregateComboCapabilities(models);
    const local = Math.min(
      getCapabilitiesForModel("glm-cn", "glm-5.3").contextWindow,
      getCapabilitiesForModel(null, "deepseek-v4.1-flash").contextWindow
    );
    expect(caps.contextWindow).toBe(local);
  });

  it("uses the fed limits when a resolver supplies them", () => {
    const resolver = (fullId) => (fullId === "glm-cn/glm-5.3" ? GLM53_FED : null);
    const caps = aggregateComboCapabilities(models, null, resolver);
    expect(caps.contextWindow).toBe(Math.min(500_000, DEEPSEEK_LOCAL_CTX));
    expect(caps.contextWindow).not.toBe(aggregateComboCapabilities(models).contextWindow);
  });

  it("keeps the fields the override does not carry", () => {
    const plain = aggregateComboCapabilities(models);
    const fed = aggregateComboCapabilities(models, null, (id) => (id === "glm-cn/glm-5.3" ? GLM53_FED : null));
    // The override carries no tools/pdf/thinking fields, so those must be unchanged.
    for (const field of ["tools", "pdf", "audioInput", "videoInput", "imageOutput", "audioOutput", "thinkingFormat"]) {
      expect(fed[field]).toEqual(plain[field]);
    }
  });

  it("takes the smallest window and the largest output across members", () => {
    // Feeding a smaller value for one member pulls the window down to it.
    const smaller = aggregateComboCapabilities(models, null, (id) => (id === "glm-cn/glm-5.3" ? { ...GLM53_FED, contextWindow: 64_000 } : null));
    expect(smaller.contextWindow).toBe(64_000);
    const caps = aggregateComboCapabilities(models, null, (id) => (id === "glm-cn/glm-5.3" ? GLM53_FED : null));
    expect(caps.maxOutput).toBe(384_000); // max across members, from deepseek
  });

  it("passes the resolver into nested combos", () => {
    const lookup = {
      zap: ["deepseek-v4.1-flash", "glm-cn/glm-5.3-flash"],
      "deepseek-v4.1-flash": ["cmc/deepseek/deepseek-v4.1-flash", "ocg/deepseek-v4.1-flash"],
    };
    const seen = [];
    const resolver = (fullId) => { seen.push(fullId); return fullId === "glm-cn/glm-5.3-flash" ? { contextWindow: 1_000_000 } : null; };
    aggregateComboCapabilities(lookup.zap, lookup, resolver);
    // The nested combo's own members were resolved with the same resolver.
    expect(seen).toContain("cmc/deepseek/deepseek-v4.1-flash");
    expect(seen).toContain("ocg/deepseek-v4.1-flash");
  });

  it("leaves the plain two-argument call unchanged", () => {
    const caps = aggregateComboCapabilities(["kimi/kimi-k3"], null);
    expect(caps).toEqual(aggregateComboCapabilities(["kimi/kimi-k3"]));
    expect(caps.contextWindow).toBe(getCapabilitiesForModel("kimi", "kimi-k3").contextWindow);
  });

  it("claude-haiku-5-5 combo: 1M only through the provider-alias resolver", () => {
    // The seat is 1M in the opencode-go provider table, but a bare `ocg/` alias
    // resolves to 200k without the resolver: the unknown provider falls to the
    // global *claude*haiku* pattern. Callers that know the alias → provider id
    // map pass it as resolveCaps.
    const seats = ["ocg/claude-haiku-5-5"];
    const seatCaps = (fullId) => {
      const slash = fullId.indexOf("/");
      if (slash <= 0) return null;
      const providerId = { ocg: "opencode-go", cmc: "commandcode" }[fullId.slice(0, slash)] || fullId.slice(0, slash);
      return getCapabilitiesForModel(providerId, fullId.slice(slash + 1));
    };
    expect(aggregateComboCapabilities(seats).contextWindow).toBe(200_000);
    expect(aggregateComboCapabilities(seats, null, seatCaps).contextWindow).toBe(1_000_000);
  });
});
