import { describe, expect, test } from "bun:test";
import type { PluginState } from "./hooks";
import {
  buildPluginDraft,
  buildPluginUpdatePayload,
  buildToggleEnabledPayload,
  formatModelAllowlist,
  isPluginDraftDirty,
  parseModelAllowlist,
  parsePreviewPayload,
  summarizePluginScope,
  toggleModelInAllowlist,
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
  test("parseModelAllowlist trims whitespace, drops empty entries, and deduplicates", () => {
    expect(parseModelAllowlist(" gpt-4o , , claude-* , gpt-4o , ")).toEqual(["gpt-4o", "claude-*"]);
    expect(parseModelAllowlist("   ")).toEqual([]);
  });

  test("toggleModelInAllowlist adds and removes model entries while preserving wildcards", () => {
    expect(toggleModelInAllowlist("claude-*", "gpt-4o")).toBe("claude-*, gpt-4o");
    expect(toggleModelInAllowlist("claude-*, gpt-4o", "gpt-4o")).toBe("claude-*");
  });

  test("formatModelAllowlist joins models with comma-space", () => {
    expect(formatModelAllowlist(["gpt-4o", "claude-*"])).toBe("gpt-4o, claude-*");
    expect(formatModelAllowlist([])).toBe("");
    expect(formatModelAllowlist(undefined)).toBe("");
  });

  test("buildPluginDraft extracts models and plugin-specific configs with defaults", () => {
    const draft = buildPluginDraft(basePlugin);
    expect(draft.modelsText).toBe("gpt-4o, claude-*");
    expect(draft.level).toBe("standard");

    const legacyLight: PluginState = { ...basePlugin, config: { level: "light" } };
    expect(buildPluginDraft(legacyLight).level).toBe("lite");

    const thinkDeeper: PluginState = {
      ...basePlugin,
      key: "thinkDeeper",
      config: { effort: "medium", budget_tokens: 16384 },
    };
    const tdDraft = buildPluginDraft(thinkDeeper);
    expect(tdDraft.effort).toBe("medium");
    expect(tdDraft.budgetTokens).toBe("16384");
  });

  test("isPluginDraftDirty detects allowlist and per-plugin config changes", () => {
    const clean = buildPluginDraft(basePlugin);
    expect(isPluginDraftDirty(basePlugin, clean)).toBe(false);
    expect(isPluginDraftDirty(basePlugin, { ...clean, modelsText: "gpt-4o,  claude-*" })).toBe(
      false,
    );
    expect(isPluginDraftDirty(basePlugin, { ...clean, modelsText: "gpt-4o" })).toBe(true);
    expect(isPluginDraftDirty(basePlugin, { ...clean, level: "ultra" })).toBe(true);

    const squeezer: PluginState = {
      ...basePlugin,
      key: "contextSqueezer",
      config: { context_window: 128000, min_turns: 4 },
    };
    const sqDraft = buildPluginDraft(squeezer);
    expect(isPluginDraftDirty(squeezer, sqDraft)).toBe(false);
    expect(isPluginDraftDirty(squeezer, { ...sqDraft, contextWindow: "64000" })).toBe(true);
  });

  test("buildToggleEnabledPayload preserves persisted config and ignores unsaved drafts", () => {
    expect(buildToggleEnabledPayload(basePlugin, false)).toEqual({
      enabled: false,
      models: ["gpt-4o", "claude-*"],
      config: { level: "standard", preserveFlag: true },
    });
  });

  test("buildPluginUpdatePayload preserves existing config keys and applies typed overrides", () => {
    const draft = { ...buildPluginDraft(basePlugin), modelsText: "gpt-4o", level: "ultra" };
    expect(buildPluginUpdatePayload(basePlugin, draft)).toEqual({
      enabled: true,
      models: ["gpt-4o"],
      config: { level: "ultra", preserveFlag: true },
    });

    const sqPlugin: PluginState = {
      ...basePlugin,
      key: "contextSqueezer",
      config: { context_window: 128000, fill_ratio: 0.75, min_turns: 4 },
    };
    const sqDraft = { ...buildPluginDraft(sqPlugin), contextWindow: "65536", minTurns: "6" };
    expect(buildPluginUpdatePayload(sqPlugin, sqDraft).config).toEqual({
      context_window: 65536,
      fill_ratio: 0.75,
      min_turns: 6,
    });
  });

  test("summarizePluginScope renders concise locale-aware summaries", () => {
    const allModelsPlugin: PluginState = { ...basePlugin, key: "imageVision", models: [] };
    expect(summarizePluginScope(allModelsPlugin, "en")).toBe("All models");
    expect(summarizePluginScope(allModelsPlugin, "id")).toBe("Semua model");
    expect(summarizePluginScope(basePlugin, "en")).toBe("2 models · standard");
    expect(summarizePluginScope(basePlugin, "id")).toBe("2 model · standard");

    const squeezer: PluginState = {
      ...basePlugin,
      key: "contextSqueezer",
      models: [],
      config: { context_window: 128000 },
    };
    expect(summarizePluginScope(squeezer, "en")).toBe("All models · 128k ctx");
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
