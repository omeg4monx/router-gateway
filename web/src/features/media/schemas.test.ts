import { describe, expect, test } from "bun:test";
import { createMediaComboSchema, createMediaProviderSchema, parseTargetsList } from "./schemas";

describe("media schemas", () => {
  test("validates valid media provider creation", () => {
    const parsed = createMediaProviderSchema.safeParse({
      name: "openai-images",
      capability: "image",
      driver: "openai_compatible",
      base_url: "https://api.openai.com",
      api_key: "sk-test",
      default_model: "gpt-image-1",
      priority: "10",
      enabled: true,
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data.priority).toBe(10);
    }
  });

  test("rejects missing base_url for remote drivers and allows empty for builtin_fetch", () => {
    const bad = createMediaProviderSchema.safeParse({
      name: "remote-no-url",
      capability: "audio",
      driver: "openai_compatible",
      base_url: "",
      priority: "100",
      enabled: true,
    });
    expect(bad.success).toBe(false);

    const okBuiltin = createMediaProviderSchema.safeParse({
      name: "local-fetch",
      capability: "fetch",
      driver: "builtin_fetch",
      base_url: "",
      priority: "100",
      enabled: true,
    });
    expect(okBuiltin.success).toBe(true);
  });

  test("parses comma/newline targets list and validates media combo", () => {
    expect(parseTargetsList("p1, p2\np1 ,  p3")).toEqual(["p1", "p2", "p3"]);

    const okCombo = createMediaComboSchema.safeParse({
      name: "img-fallback",
      capability: "image",
      targets_text: "mp_1, mp_2",
      enabled: true,
    });
    expect(okCombo.success).toBe(true);

    const badCombo = createMediaComboSchema.safeParse({
      name: "empty-combo",
      capability: "image",
      targets_text: " , ",
      enabled: true,
    });
    expect(badCombo.success).toBe(false);
  });
});
