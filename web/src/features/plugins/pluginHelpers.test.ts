import { describe, expect, test } from "bun:test";
import type { PluginState } from "./hooks";
import {
  buildPluginDraft,
  buildPluginUpdatePayload,
  formatModelAllowlist,
  isPluginDraftDirty,
  parseModelAllowlist,
  parsePreviewPayload,
  summarizePluginScope,
} from "./pluginHelpers";

const basePlugin: PluginState = {
  key: "antiSlop",
  name: "Anti-Slop",
  description: "Reduces filler phrases in model responses.",
  category: "quality",
  enabled: true,
  models: ["gpt-4o", "claude-*"],
  config: { level: "standard", preserveFlag: true },
  updated_at: "2026-11-26T00:00:00Z",
};

describe("pluginHelpers", () => {
  test("parseModelAllowlist trims whitespace and drops empty entries", () => {
    expect(parseModelAllowlist(" gpt-4o , , claude-* , ")).toEqual(["gpt-4o", "claude-*"]);
    expect(parseModelAllowlist("   ")).toEqual([]);
  });

  test("formatModelAllowlist joins models with comma-space", () => {
    expect(formatModelAllowlist(["gpt-4o", "claude-*"])).toBe("gpt-4o, claude-*");
    expect(formatModelAllowlist([])).toBe("");
    expect(formatModelAllowlist(undefined)).toBe("");
  });

  test("buildPluginDraft extracts models and antiSlop level with fallback", () => {
    expect(buildPluginDraft(basePlugin)).toEqual({
      modelsText: "gpt-4o, claude-*",
      level: "standard",
    });

    const withoutLevel: PluginState = { ...basePlugin, config: {} };
    expect(buildPluginDraft(withoutLevel)).toEqual({
      modelsText: "gpt-4o, claude-*",
      level: "full",
    });
  });

  test("isPluginDraftDirty detects allowlist or antiSlop level changes only", () => {
    const clean = buildPluginDraft(basePlugin);
    expect(isPluginDraftDirty(basePlugin, clean)).toBe(false);
    expect(isPluginDraftDirty(basePlugin, { ...clean, modelsText: "gpt-4o,  claude-*" })).toBe(
      false,
    );
    expect(isPluginDraftDirty(basePlugin, { ...clean, modelsText: "gpt-4o" })).toBe(true);
    expect(isPluginDraftDirty(basePlugin, { ...clean, level: "full" })).toBe(true);

    const otherPlugin: PluginState = { ...basePlugin, key: "jsonMode", config: {} };
    expect(
      isPluginDraftDirty(otherPlugin, { modelsText: "gpt-4o, claude-*", level: "light" }),
    ).toBe(false);
  });

  test("buildPluginUpdatePayload preserves existing config keys and applies overrides", () => {
    const draft = { modelsText: "gpt-4o", level: "light" };
    expect(buildPluginUpdatePayload(basePlugin, draft)).toEqual({
      enabled: true,
      models: ["gpt-4o"],
      config: { level: "light", preserveFlag: true },
    });

    expect(buildPluginUpdatePayload(basePlugin, draft, false)).toEqual({
      enabled: false,
      models: ["gpt-4o"],
      config: { level: "light", preserveFlag: true },
    });
  });

  test("summarizePluginScope renders concise locale-aware summaries", () => {
    const allModelsPlugin: PluginState = { ...basePlugin, key: "jsonMode", models: [] };
    expect(summarizePluginScope(allModelsPlugin, "en")).toBe("All models");
    expect(summarizePluginScope(allModelsPlugin, "id")).toBe("Semua model");
    expect(summarizePluginScope(basePlugin, "en")).toBe("2 models · standard");
    expect(summarizePluginScope(basePlugin, "id")).toBe("2 model · standard");
  });

  test("parsePreviewPayload validates JSON object payloads", () => {
    expect(parsePreviewPayload('{"model":"gpt-4o"}')).toEqual({
      ok: true,
      payload: { model: "gpt-4o" },
    });
    expect(parsePreviewPayload("[1, 2, 3]")).toEqual({ ok: false });
    expect(parsePreviewPayload("{broken")).toEqual({ ok: false });
  });
});
