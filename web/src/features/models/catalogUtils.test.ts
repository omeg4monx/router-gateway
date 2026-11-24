import { describe, expect, test } from "bun:test";
import {
  classifyTestError,
  comboTestResultKey,
  computeUpdatedAllowedModels,
  computeVisibilityCounts,
  countKeyAccess,
  DEFAULT_COMBO_FILTERS,
  DEFAULT_MODEL_FILTERS,
  filterCombos,
  filterModels,
  flattenCatalogResponse,
  formatTokenCount,
  modelTestResultKey,
  normalizeComboRows,
  normalizeTestStatus,
  sortModels,
} from "./catalogUtils";
import type { ApiKeyAccessItem, CatalogTestResult, ProviderHealthState } from "./types";

describe("models/catalogUtils (OmniRoute Model Catalog parity)", () => {
  test("flattens provider-grouped catalog and computes visibility counts", () => {
    const rows = flattenCatalogResponse({
      openai: {
        provider: "openai",
        models: [
          {
            id: "gpt-4o",
            name: "gpt-4o",
            type: "chat",
            subtype: "general",
            context_length: 128000,
            max_output_tokens: 16384,
            capabilities: { vision: true, streaming: true },
          },
          {
            id: "o3-mini",
            name: "o3-mini",
            type: "chat",
            subtype: "reasoning",
            context_length: 200000,
            max_output_tokens: 100000,
            capabilities: { reasoning: true },
            hidden: true,
          },
        ],
      },
      ollama: {
        provider: "ollama",
        models: [{ id: "llama3.3", type: "chat", free: true, custom: true }],
      },
    });

    expect(rows).toHaveLength(3);
    expect(rows[0]?.qualifiedId).toBe("openai/gpt-4o");
    expect(computeVisibilityCounts(rows)).toEqual({ visible: 2, hidden: 1, total: 3 });
  });

  test("filters models by query, capability, pricing, visibility, provider health, test result, and token thresholds", () => {
    const rows = flattenCatalogResponse({
      openai: {
        provider: "openai",
        models: [
          {
            id: "gpt-4o",
            type: "chat",
            subtype: "general",
            context_length: 128000,
            max_output_tokens: 16384,
            capabilities: { vision: true },
          },
          {
            id: "o3-mini",
            type: "chat",
            subtype: "reasoning",
            context_length: 200000,
            max_output_tokens: 100000,
            capabilities: { reasoning: true },
            hidden: true,
          },
        ],
      },
      ollama: {
        provider: "ollama",
        models: [
          {
            id: "qwen2.5-coder",
            type: "chat",
            subtype: "code",
            free: true,
            context_length: 32000,
            max_output_tokens: 8192,
            capabilities: { code: true },
          },
        ],
      },
    });

    const health: Record<string, ProviderHealthState> = {
      openai: "degraded",
      ollama: "healthy",
    };
    const results: Record<string, CatalogTestResult> = {
      [modelTestResultKey("openai", "gpt-4o")]: {
        id: modelTestResultKey("openai", "gpt-4o"),
        targetType: "model",
        providerId: "openai",
        modelId: "gpt-4o",
        status: "ok",
        latencyMs: 320,
        testedAt: Date.now(),
      },
    };

    expect(
      filterModels(
        rows,
        { ...DEFAULT_MODEL_FILTERS, capability: "reasoning" },
        health,
        results,
      ).map((r) => r.id),
    ).toEqual(["o3-mini"]);

    expect(
      filterModels(rows, { ...DEFAULT_MODEL_FILTERS, pricing: "free" }, health, results).map(
        (r) => r.id,
      ),
    ).toEqual(["qwen2.5-coder"]);

    expect(
      filterModels(rows, { ...DEFAULT_MODEL_FILTERS, visibility: "hidden" }, health, results).map(
        (r) => r.id,
      ),
    ).toEqual(["o3-mini"]);

    expect(
      filterModels(
        rows,
        { ...DEFAULT_MODEL_FILTERS, providerHealth: "healthy" },
        health,
        results,
      ).map((r) => r.id),
    ).toEqual(["qwen2.5-coder"]);

    expect(
      filterModels(rows, { ...DEFAULT_MODEL_FILTERS, testResult: "ok" }, health, results).map(
        (r) => r.id,
      ),
    ).toEqual(["gpt-4o"]);

    expect(
      filterModels(
        rows,
        { ...DEFAULT_MODEL_FILTERS, minContextLength: 150000, minMaxOutputTokens: 50000 },
        health,
        results,
      ).map((r) => r.id),
    ).toEqual(["o3-mini"]);

    const sorted = sortModels(rows, "context_length", "desc", results);
    expect(sorted.map((r) => r.id)).toEqual(["o3-mini", "gpt-4o", "qwen2.5-coder"]);
  });

  test("classifies test errors and normalizes slow threshold (>=5000ms)", () => {
    expect(classifyTestError("rate limit exceeded", 429)).toBe("rate-limited");
    expect(classifyTestError("insufficient balance on account", 402)).toBe("quota");
    expect(classifyTestError("request timed out", 504)).toBe("timeout");
    expect(classifyTestError("internal server error", 500)).toBe("other");

    expect(normalizeTestStatus("ok", 1200)).toBe("ok");
    expect(normalizeTestStatus("ok", 5200)).toBe("slow");
    expect(normalizeTestStatus("error", 100)).toBe("error");
  });

  test("filters combos and computes key access & restricted-mode updates", () => {
    const combos = normalizeComboRows([
      {
        id: "c1",
        name: "coding-primary",
        description: "Primary coding chain",
        targets: ["openai/gpt-4o", "anthropic/claude-sonnet-4-5"],
        steps: [
          { kind: "model", model: "openai/gpt-4o" },
          { kind: "model", model: "anthropic/claude-sonnet-4-5" },
        ],
        strategy: "priority",
        config: {},
        enabled: true,
        created_at: "2026-11-24T00:00:00Z",
        updated_at: "2026-11-24T00:00:00Z",
      },
      {
        id: "c2",
        name: "paused-chain",
        description: "",
        targets: ["ollama/llama3.3"],
        steps: [{ kind: "model", model: "ollama/llama3.3" }],
        strategy: "cost-optimized",
        config: {},
        enabled: false,
        created_at: "2026-11-24T00:00:00Z",
        updated_at: "2026-11-24T00:00:00Z",
      },
    ]);

    const comboResults: Record<string, CatalogTestResult> = {
      [comboTestResultKey("coding-primary")]: {
        id: comboTestResultKey("coding-primary"),
        targetType: "combo",
        comboName: "coding-primary",
        status: "ok",
        latencyMs: 410,
        testedAt: Date.now(),
      },
    };

    expect(
      filterCombos(
        combos,
        { ...DEFAULT_COMBO_FILTERS, status: "active", minMembers: 2 },
        comboResults,
      ).map((c) => c.name),
    ).toEqual(["coding-primary"]);

    const keys: ApiKeyAccessItem[] = [
      {
        id: "k1",
        name: "Wildcard Key",
        prefix: "sk-om4-1",
        disabled: false,
        allowedModels: ["*"],
        modelAccessMode: "all",
      },
      {
        id: "k2",
        name: "Restricted Key",
        prefix: "sk-om4-2",
        disabled: false,
        allowedModels: ["openai/*"],
        modelAccessMode: "restricted",
      },
    ];

    expect(countKeyAccess(keys, "openai/gpt-4o", "gpt-4o")).toEqual({ allowed: 2, total: 2 });
    expect(countKeyAccess(keys, "anthropic/claude-sonnet-4-5")).toEqual({ allowed: 1, total: 2 });

    const [wildcardKey, restrictedKey] = keys;
    if (!wildcardKey || !restrictedKey) {
      throw new Error("Expected fixture keys");
    }
    expect(computeUpdatedAllowedModels(wildcardKey, ["openai/gpt-4o"], "add")).toEqual([
      "openai/gpt-4o",
    ]);
    expect(
      computeUpdatedAllowedModels(restrictedKey, ["anthropic/claude-sonnet-4-5"], "add"),
    ).toEqual(["openai/*", "anthropic/claude-sonnet-4-5"]);
    expect(computeUpdatedAllowedModels(restrictedKey, ["openai/*"], "remove")).toEqual([]);
    expect(formatTokenCount(128000)).toBe("128K");
    expect(formatTokenCount(1048576)).toBe("1.0M");
  });
});
