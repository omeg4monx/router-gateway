import { describe, expect, test } from "bun:test";
import type { PluginState } from "./hooks";
import {
  buildPluginCatalogModelOptions,
  buildPluginDraft,
  buildPluginUpdatePayload,
  buildToggleEnabledPayload,
  filterPluginCatalogModelOptions,
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

  test("buildPluginCatalogModelOptions merges active catalog with built-in tokenharbor and codebuddy defaults", () => {
    const options = buildPluginCatalogModelOptions(
      [
        {
          providerId: "openai",
          provider: "OpenAI",
          id: "gpt-4o",
          qualifiedId: "openai/gpt-4o",
          name: "gpt-4o",
          type: "chat",
        },
      ],
      {
        adapters: [
          {
            id: "tokenharbor",
            display_name: "Token Harbor",
            category: "api_key",
            protocol: "openai-chat",
            auth_type: "api_key",
            default_base_url: "https://tokenharbor.ai/v1",
            default_models: ["deepseek-v4.1-flash:free", "mimo-v2.6-flash:free"],
            oauth: null,
          },
          {
            id: "codebuddy-intl",
            display_name: "CodeBuddy / WorkBuddy AI",
            category: "oauth",
            protocol: "openai-chat",
            auth_type: "oauth",
            default_base_url: "https://www.codebuddy.ai/v2",
            default_models: ["glm-5.2", "deepseek-v4.1-flash"],
            oauth: null,
          },
        ],
        custom_nodes: [],
      },
    );

    expect(options.map((o) => o.id)).toEqual([
      "openai/gpt-4o",
      "tokenharbor/deepseek-v4.1-flash:free",
      "tokenharbor/mimo-v2.6-flash:free",
      "codebuddy-intl/glm-5.2",
      "codebuddy-intl/deepseek-v4.1-flash",
    ]);

    expect(filterPluginCatalogModelOptions(options, "token harbor").map((o) => o.id)).toEqual([
      "tokenharbor/deepseek-v4.1-flash:free",
      "tokenharbor/mimo-v2.6-flash:free",
    ]);
    expect(filterPluginCatalogModelOptions(options, "codebuddy").map((o) => o.id)).toEqual([
      "codebuddy-intl/glm-5.2",
      "codebuddy-intl/deepseek-v4.1-flash",
    ]);
  });
});
