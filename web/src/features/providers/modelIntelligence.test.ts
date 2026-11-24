import { describe, expect, test } from "bun:test";
import { getModelIntelligenceProfile, sortProviderModels } from "./modelIntelligence";

describe("modelIntelligence hybrid scoring", () => {
  test("ranks curated models accurately across Flagship, Reasoning, Balanced, and Fast tiers", () => {
    const opus = getModelIntelligenceProfile("claude-opus-4-6");
    const o3 = getModelIntelligenceProfile("o3");
    const sonnet = getModelIntelligenceProfile("claude-sonnet-4-5");
    const mini = getModelIntelligenceProfile("gpt-4o-mini");

    expect(opus.source).toBe("curated");
    expect(opus.tier).toBe("flagship");
    expect(o3.tier).toBe("reasoning");
    expect(sonnet.tier).toBe("balanced");
    expect(mini.tier).toBe("fast");

    expect(opus.score).toBeGreaterThan(sonnet.score);
    expect(o3.score).toBeGreaterThan(mini.score);
    expect(sonnet.score).toBeGreaterThan(mini.score);
  });

  test("normalizes provider prefixes, date suffixes, and :free tags before matching", () => {
    const prefixed = getModelIntelligenceProfile("anthropic/claude-opus-4-6-20260205");
    const freeTagged = getModelIntelligenceProfile("deepseek/deepseek-r1:free");

    expect(prefixed.source).toBe("curated");
    expect(prefixed.tier).toBe("flagship");
    expect(freeTagged.source).toBe("curated");
    expect(freeTagged.tier).toBe("reasoning");
  });

  test("falls back to heuristic scoring for unknown and custom models", () => {
    const bigReasoning = getModelIntelligenceProfile("custom-frontier-405b-thinking");
    const mid70b = getModelIntelligenceProfile("custom-workhorse-70b");
    const tinyFast = getModelIntelligenceProfile("custom-edge-7b-mini");

    expect(bigReasoning.source).toBe("heuristic");
    expect(bigReasoning.tier).toBe("reasoning");
    expect(tinyFast.tier).toBe("fast");
    expect(bigReasoning.score).toBeGreaterThan(mid70b.score);
    expect(mid70b.score).toBeGreaterThan(tinyFast.score);
  });

  test("sortProviderModels sorts by intelligence desc/asc and name asc/desc without mutating input", () => {
    const input = [
      { id: "gpt-4o-mini" },
      { id: "claude-opus-4-6" },
      { id: "gpt-4.1" },
      { id: "o3" },
    ];
    const snapshot = input.map((m) => m.id);

    const smartestFirst = sortProviderModels(input, "intelligence", "desc").map((m) => m.id);
    expect(smartestFirst[0]).toBe("claude-opus-4-6");
    expect(smartestFirst[smartestFirst.length - 1]).toBe("gpt-4o-mini");

    const lightestFirst = sortProviderModels(input, "intelligence", "asc").map((m) => m.id);
    expect(lightestFirst[0]).toBe("gpt-4o-mini");
    expect(lightestFirst[lightestFirst.length - 1]).toBe("claude-opus-4-6");

    expect(sortProviderModels(input, "name", "asc").map((m) => m.id)).toEqual([
      "claude-opus-4-6",
      "gpt-4.1",
      "gpt-4o-mini",
      "o3",
    ]);
    expect(input.map((m) => m.id)).toEqual(snapshot);
  });
});
