import { describe, expect, test } from "bun:test";
import type { ProviderConnectionSummary } from "@/features/providers/schemas";
import {
  buildProviderAccountGroups,
  formatStepAccountLabel,
  getModelsForSelection,
  isDuplicateModelStep,
  qualifyProviderModel,
} from "./comboHelpers";
import { COMBO_STRATEGY_OPTIONS, STRATEGY_CATEGORIES } from "./strategies";

function makeConn(
  overrides: Partial<ProviderConnectionSummary> &
    Pick<ProviderConnectionSummary, "id" | "provider_id" | "name">,
): ProviderConnectionSummary {
  return {
    protocol: "openai-chat",
    base_url: "https://api.example.com/v1",
    auth_type: "api_key",
    priority: 1,
    enabled: true,
    models: [],
    headers: {},
    has_credentials: true,
    token_expires_at: null,
    last_tested_at: null,
    last_test_status: "ok",
    last_test_error: "",
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
    ...overrides,
  };
}

describe("comboHelpers", () => {
  const connections: ProviderConnectionSummary[] = [
    makeConn({
      id: "conn-oa-1",
      provider_id: "openai",
      name: "OpenAI Primary",
      auth_type: "api_key",
      models: ["gpt-4o", "openai/gpt-4o-mini"],
    }),
    makeConn({
      id: "conn-oa-2",
      provider_id: "openai",
      name: "OpenAI OAuth Team",
      auth_type: "oauth",
      models: ["gpt-4o", "o3-mini"],
    }),
    makeConn({
      id: "conn-oa-disabled",
      provider_id: "openai",
      name: "OpenAI Disabled",
      enabled: false,
      models: ["gpt-3.5-turbo"],
    }),
    makeConn({
      id: "conn-ant-1",
      provider_id: "anthropic",
      name: "Claude Prod",
      auth_type: "api_key",
      models: ["claude-sonnet-4-20250514"],
    }),
  ];

  test("groups only enabled connections by provider and counts accounts", () => {
    const groups = buildProviderAccountGroups(connections, {
      openai: "OpenAI",
      anthropic: "Anthropic",
    });

    expect(groups).toHaveLength(2);
    const openai = groups.find((g) => g.providerId === "openai");
    const anthropic = groups.find((g) => g.providerId === "anthropic");

    expect(openai?.displayName).toBe("OpenAI");
    expect(openai?.accounts).toHaveLength(2);
    expect(openai?.hasMultipleAccounts).toBe(true);
    expect(anthropic?.accounts).toHaveLength(1);
    expect(anthropic?.hasMultipleAccounts).toBe(false);
  });

  test("resolves union of models for dynamic mode and account-specific models when pinned", () => {
    const groups = buildProviderAccountGroups(connections);
    const openai = groups.find((g) => g.providerId === "openai");
    expect(openai).toBeDefined();
    if (!openai) return;

    const dynamicModels = getModelsForSelection(openai, "");
    expect(dynamicModels).toEqual(["gpt-4o", "gpt-4o-mini", "o3-mini"]);

    const pinnedModels = getModelsForSelection(openai, "conn-oa-2");
    expect(pinnedModels).toEqual(["gpt-4o", "o3-mini"]);

    // Merges catalog models in dynamic mode and falls back to catalog when connection.models is empty
    const mergedDynamic = getModelsForSelection(openai, "", {
      openai: ["openai/gpt-4.1", "gpt-4o"],
    });
    expect(mergedDynamic).toEqual(["gpt-4.1", "gpt-4o", "gpt-4o-mini", "o3-mini"]);

    const emptyConnGroup = buildProviderAccountGroups([
      makeConn({ id: "conn-gem-1", provider_id: "gemini", name: "Gemini OAuth", models: [] }),
    ])[0];
    expect(
      getModelsForSelection(emptyConnGroup, "conn-gem-1", {
        gemini: ["gemini-2.5-pro", "gemini/gemini-2.5-flash"],
      }),
    ).toEqual(["gemini-2.5-flash", "gemini-2.5-pro"]);
  });

  test("qualifies model identifier without duplicating provider prefix", () => {
    expect(qualifyProviderModel("openai", "gpt-4o")).toBe("openai/gpt-4o");
    expect(qualifyProviderModel("openai", "openai/gpt-4o")).toBe("openai/gpt-4o");
    expect(qualifyProviderModel("openrouter", "anthropic/claude-3.5-sonnet")).toBe(
      "anthropic/claude-3.5-sonnet",
    );
  });

  test("allows same model on different accounts while blocking exact duplicates", () => {
    const steps = [{ kind: "model" as const, model: "openai/gpt-4o", connectionId: "conn-oa-1" }];
    expect(isDuplicateModelStep(steps, "openai/gpt-4o", "conn-oa-1")).toBe(true);
    expect(isDuplicateModelStep(steps, "openai/gpt-4o", "conn-oa-2")).toBe(false);
    expect(isDuplicateModelStep(steps, "openai/gpt-4o", "")).toBe(false);
  });

  test("formats step account badge label for dynamic vs pinned account", () => {
    const byId = new Map(connections.map((c) => [c.id, c]));
    expect(formatStepAccountLabel(undefined, byId, true)).toBe("Dynamic (Semua Akun)");
    expect(formatStepAccountLabel("", byId, false)).toBe("Dynamic (All Accounts)");
    expect(formatStepAccountLabel("conn-oa-2", byId, true)).toBe("OpenAI OAuth Team");
    expect(formatStepAccountLabel("missing-id", byId, false)).toBe("missing-id");
  });
});

describe("COMBO_STRATEGY_OPTIONS guide completeness", () => {
  test("all 22 strategies include informative OmniRoute guide fields in EN and ID", () => {
    expect(COMBO_STRATEGY_OPTIONS).toHaveLength(22);
    const validCategories = new Set(STRATEGY_CATEGORIES.map((c) => c.id));

    for (const s of COMBO_STRATEGY_OPTIONS) {
      expect(validCategories.has(s.category)).toBe(true);
      expect(s.howItWorksEn.trim().length).toBeGreaterThan(15);
      expect(s.howItWorksId.trim().length).toBeGreaterThan(15);
      expect(s.bestForEn.trim().length).toBeGreaterThan(10);
      expect(s.bestForId.trim().length).toBeGreaterThan(10);
      expect(s.whenToUseEn.trim().length).toBeGreaterThan(10);
      expect(s.whenToUseId.trim().length).toBeGreaterThan(10);
      expect(s.flowPreview.trim().length).toBeGreaterThan(3);
    }
  });
});
