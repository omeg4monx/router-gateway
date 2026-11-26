import { describe, expect, test } from "bun:test";
import { existsSync } from "node:fs";
import { join, resolve } from "node:path";
import {
  buildProviderEntries,
  canAddProviderConnection,
  canProbeProviderModels,
  computeFleetTelemetry,
  countCategoryBuckets,
  countStatusBuckets,
  getProviderMonogram,
  matchesProviderSearch,
  matchesStatusFilter,
  sortProviderEntries,
} from "./catalogUtils";
import { resolveProviderIconSrc } from "./ProviderLogo";
import type { CatalogAdapterDescriptor, ProviderConnectionSummary, ProviderNode } from "./schemas";

const ADAPTERS: CatalogAdapterDescriptor[] = [
  {
    id: "openai",
    display_name: "OpenAI",
    category: "api_key",
    protocol: "openai-chat",
    auth_type: "api_key",
    default_base_url: "https://api.openai.com/v1",
    default_models: ["gpt-4o", "o3-mini"],
    oauth: null,
  },
  {
    id: "claude",
    display_name: "Claude OAuth",
    category: "oauth",
    protocol: "anthropic",
    auth_type: "oauth",
    default_base_url: "https://api.anthropic.com",
    default_models: ["claude-sonnet-4-5"],
    oauth: null,
  },
  {
    id: "gemini",
    display_name: "Google Gemini",
    category: "api_key",
    protocol: "gemini",
    auth_type: "api_key",
    default_base_url: "https://generativelanguage.googleapis.com",
    default_models: ["gemini-2.5-pro"],
    oauth: null,
  },
];

const CUSTOM_NODES: ProviderNode[] = [
  {
    id: "vllm-local",
    name: "Local vLLM",
    protocol: "openai-chat",
    base_url: "http://127.0.0.1:8000/v1",
    models: ["llama-3.3-70b"],
    headers: {},
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
  },
];

const CONNECTIONS: ProviderConnectionSummary[] = [
  {
    id: "conn-1",
    provider_id: "openai",
    name: "OpenAI Prod",
    protocol: "openai-chat",
    auth_type: "api_key",
    has_credentials: true,
    base_url: "https://api.openai.com/v1",
    models: ["gpt-4o"],
    headers: {},
    priority: 100,
    enabled: true,
    token_expires_at: null,
    last_test_status: "ok",
    last_tested_at: "2026-01-02T00:00:00Z",
    last_test_error: "",
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-02T00:00:00Z",
  },
  {
    id: "conn-2",
    provider_id: "claude",
    name: "Claude Team",
    protocol: "anthropic",
    auth_type: "oauth",
    has_credentials: true,
    base_url: "https://api.anthropic.com",
    models: ["claude-sonnet-4-5"],
    headers: {},
    priority: 100,
    enabled: true,
    token_expires_at: null,
    last_test_status: "error",
    last_tested_at: "2026-01-02T01:00:00Z",
    last_test_error: "401 unauthorized",
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-02T01:00:00Z",
  },
  {
    id: "conn-3",
    provider_id: "vllm-local",
    name: "vLLM GPU-0",
    protocol: "openai-chat",
    auth_type: "api_key",
    has_credentials: true,
    base_url: "http://127.0.0.1:8000/v1",
    models: ["llama-3.3-70b"],
    headers: {},
    priority: 100,
    enabled: false,
    token_expires_at: null,
    last_test_status: "untested",
    last_tested_at: null,
    last_test_error: "",
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
  },
];

describe("providers catalogUtils", () => {
  test("merges Go adapters, built-in provider registry, custom nodes, and fleet telemetry accurately", () => {
    const entries = buildProviderEntries(ADAPTERS, CUSTOM_NODES, CONNECTIONS);
    // Includes all 138 built-in registry providers + custom vllm-local
    expect(entries.length).toBeGreaterThanOrEqual(138);

    const openaiEntry = entries.find((e) => e.id === "openai");
    expect(openaiEntry?.isBuiltinAdapter).toBe(true);
    expect(openaiEntry?.apiKeyUrl).toContain("platform.openai.com");

    const deepseekEntry = entries.find((e) => e.id === "deepseek");
    expect(deepseekEntry?.name).toBe("DeepSeek");
    expect(deepseekEntry?.apiKeyUrl).toContain("deepseek.com");
    expect(deepseekEntry?.defaultModels.length).toBeGreaterThan(0);

    const telemetry = computeFleetTelemetry(entries, CONNECTIONS);
    expect(telemetry.activeProviders).toBe(1);
    expect(telemetry.attentionProviders).toBe(1);
    expect(telemetry.disabledProviders).toBe(1);
    expect(telemetry.totalCatalog).toBe(entries.length);
    expect(telemetry.enabledConnections).toBe(2);
    expect(telemetry.totalConnections).toBe(3);
    expect(telemetry.oauthConnections).toBe(1);
    expect(telemetry.apiKeyConnections).toBe(2);

    const statusBuckets = countStatusBuckets(entries);
    expect(statusBuckets.all).toBe(entries.length);
    expect(statusBuckets.active).toBe(1);
    expect(statusBuckets.error).toBe(1);
    expect(statusBuckets.inactive).toBe(1);
    expect(statusBuckets.unconfigured).toBe(entries.length - 3);

    const catBuckets = countCategoryBuckets(entries);
    expect(catBuckets.all).toBe(entries.length);
    expect(catBuckets.oauth).toBeGreaterThan(10);
    expect(catBuckets.apikey).toBeGreaterThan(50);
    expect(catBuckets.free).toBeGreaterThan(10);
    expect(catBuckets.custom).toBe(1);
  });

  test("filters by status, searches across names/models/connections, and sorts errors first", () => {
    const entries = buildProviderEntries(ADAPTERS, CUSTOM_NODES, CONNECTIONS);
    const sorted = sortProviderEntries(entries);
    expect(sorted.slice(0, 3).map((e) => e.id)).toEqual(["claude", "openai", "vllm-local"]);

    expect(entries.filter((e) => matchesStatusFilter(e, "error")).map((e) => e.id)).toEqual([
      "claude",
    ]);
    expect(
      entries
        .filter((e) => matchesProviderSearch(e, "llama-3.3"))
        .some((e) => e.id === "vllm-local"),
    ).toBe(true);
    expect(entries.filter((e) => matchesProviderSearch(e, "OpenAI Prod")).map((e) => e.id)).toEqual(
      ["openai"],
    );
  });

  test("generates clean 2-letter monograms", () => {
    expect(getProviderMonogram("Google Gemini", "gemini")).toBe("GG");
    expect(getProviderMonogram("OpenAI", "openai")).toBe("OP");
    expect(getProviderMonogram("", "vllm-local")).toBe("VL");
  });

  test("resolves existing logo assets in web/public/providers for all built-in Go adapters and custom protocols", () => {
    const publicDir = resolve(import.meta.dir, "../../../public");
    const builtinAdapterIds = [
      "openai",
      "anthropic",
      "openai-responses",
      "gemini",
      "openrouter",
      "openai-compatible",
      "anthropic-compatible",
      "claude",
      "codex",
      "gemini-oauth",
      "qwen",
      "github-copilot",
      "kiro",
      "cursor",
      "antigravity",
      "aws-quick",
      "venice",
      "freepi",
    ];

    for (const id of builtinAdapterIds) {
      const src = resolveProviderIconSrc(id);
      expect(src).not.toBeNull();
      const assetPath = join(publicDir, src ?? "");
      expect(existsSync(assetPath)).toBe(true);
    }

    expect(resolveProviderIconSrc("local-vllm", "openai-chat", true)).toBe(
      "/providers/custom-openai.svg",
    );
    expect(resolveProviderIconSrc("my-claude-proxy", "anthropic", true)).toBe(
      "/providers/custom-anthropic.svg",
    );
    expect(resolveProviderIconSrc("my-moonshot", "openai-chat", true)).toBe(
      "/providers/custom-moonshot.svg",
    );
  });

  test("sorts provider entries alphabetically by asc (A-Z) and desc (Z-A)", () => {
    const entries = buildProviderEntries(ADAPTERS, CUSTOM_NODES, CONNECTIONS).filter((e) =>
      ["openai", "claude", "vllm-local"].includes(e.id),
    );
    const ascNames = sortProviderEntries(entries, "asc").map((e) => e.name);
    const descNames = sortProviderEntries(entries, "desc").map((e) => e.name);
    expect(ascNames).toEqual(["Claude Code", "Local vLLM", "OpenAI"]);
    expect(descNames).toEqual(["OpenAI", "Local vLLM", "Claude Code"]);
  });

  test("enables model probing for ready keyless free providers (venice, opencode, ovh) without connections", () => {
    const entries = buildProviderEntries(ADAPTERS, CUSTOM_NODES, CONNECTIONS);
    const byId = new Map(entries.map((entry) => [entry.id, entry]));

    const opencode = byId.get("opencode");
    const venice = byId.get("venice");
    const ovh = byId.get("ovh");
    const edgeTts = byId.get("edge-tts");
    const openai = byId.get("openai");
    const anthropic = byId.get("anthropic");

    expect(opencode).toBeDefined();
    expect(venice).toBeDefined();
    expect(ovh).toBeDefined();
    expect(edgeTts).toBeDefined();
    expect(openai).toBeDefined();
    expect(anthropic).toBeDefined();

    if (opencode) expect(canProbeProviderModels(opencode)).toBe(true);
    if (venice) expect(canProbeProviderModels(venice)).toBe(true);
    if (ovh) expect(canProbeProviderModels(ovh)).toBe(true);
    // Keyless providers without a baseUrl (e.g. edge-tts) cannot be probed
    if (edgeTts) expect(canProbeProviderModels(edgeTts)).toBe(false);
    // Normal provider with enabled connection can be probed
    if (openai) expect(canProbeProviderModels(openai)).toBe(true);
    // Normal provider without connections cannot be probed
    if (anthropic) expect(canProbeProviderModels(anthropic)).toBe(false);
    // Keyless provider whose connections are all explicitly disabled cannot be probed
    if (venice) {
      expect(
        canProbeProviderModels({
          ...venice,
          stats: { ...venice.stats, total: 1, disabled: 1, connected: 0, allDisabled: true },
        }),
      ).toBe(false);
    }
    // OpenCode Free must point to /inference/openai/v1 and list active free models (not retired 400 models)
    if (opencode) {
      expect(opencode.baseUrl).toBe("https://opencode.ai/inference/openai/v1");
      expect(opencode.defaultModels).toContain("big-pickle");
      expect(opencode.defaultModels).toContain("mimo-v2.6-flash-free");
      expect(opencode.defaultModels).toContain("deepseek-v4.1-flash-free");
      expect(opencode.defaultModels).toContain("iceman-free");
      expect(opencode.defaultModels).not.toContain("gpt-4o-mini-free");
      expect(opencode.defaultModels).not.toContain("space-bunny-free");
    }
  });

  test("hides Add connection action for noAuth providers and lists only keyless Venice models", () => {
    const entries = buildProviderEntries(ADAPTERS, CUSTOM_NODES, CONNECTIONS);
    const byId = new Map(entries.map((entry) => [entry.id, entry]));

    const venice = byId.get("venice");
    const veniceFree = byId.get("venice-free");
    const opencode = byId.get("opencode");
    const openai = byId.get("openai");
    const anthropic = byId.get("anthropic");

    expect(venice).toBeDefined();
    expect(veniceFree).toBeDefined();
    expect(opencode).toBeDefined();
    expect(openai).toBeDefined();
    expect(anthropic).toBeDefined();

    if (venice) {
      expect(canAddProviderConnection(venice)).toBe(false);
      expect(venice.defaultModels).toContain("venice-uncensored-1-2");
      expect(venice.defaultModels).toContain("deepseek-v4-1-flash-uncensored");
      expect(venice.defaultModels).not.toContain("deepseek-v4-pro");
      expect(venice.defaultModels).not.toContain("zai-org-glm-5");
      expect(venice.defaultModels).not.toContain("qwen3-vl-235b-a22b");
      expect(venice.defaultModels).not.toContain("llama-3.3-70b");
    }
    if (veniceFree) {
      expect(veniceFree.noAuth).toBe(true);
      expect(canAddProviderConnection(veniceFree)).toBe(false);
    }
    if (opencode) expect(canAddProviderConnection(opencode)).toBe(false);
    if (openai) expect(canAddProviderConnection(openai)).toBe(true);
    if (anthropic) expect(canAddProviderConnection(anthropic)).toBe(true);
  });
});
