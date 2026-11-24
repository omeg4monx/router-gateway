import { describe, expect, test } from "bun:test";
import { buildProviderEntries } from "./catalogUtils";
import { PROVIDER_REGISTRY_MAP } from "./providerRegistry";
import type { CatalogAdapterDescriptor } from "./schemas";

function makeAdapter(
  id: string,
  authType: "oauth" | "api_key",
  defaultModels: string[],
): CatalogAdapterDescriptor {
  return {
    id,
    display_name: id,
    category: authType,
    protocol: "openai-chat",
    auth_type: authType,
    default_base_url: "https://example.com/v1",
    default_models: defaultModels,
    oauth: null,
  };
}

describe("OmniRoute provider model catalog accuracy & sync", () => {
  test("built-in Go adapters use full OmniRoute registry models instead of tiny Go stubs", () => {
    const adapters: CatalogAdapterDescriptor[] = [
      makeAdapter("claude", "oauth", ["claude-sonnet-4-20250514", "claude-3-7-sonnet-20250219"]),
      makeAdapter("codex", "oauth", ["gpt-5-codex", "o4-mini"]),
      makeAdapter("antigravity", "oauth", ["gemini-2.5-pro", "gemini-2.5-flash"]),
      makeAdapter("github-copilot", "oauth", ["gpt-4o", "claude-sonnet-4"]),
      makeAdapter("kiro", "oauth", ["claude-sonnet-4-20250514"]),
      makeAdapter("openai", "api_key", ["gpt-4o", "gpt-4o-mini"]),
      makeAdapter("openai-responses", "api_key", ["gpt-4.1", "gpt-4o"]),
      makeAdapter("freepi", "oauth", ["particle"]),
    ];

    const entries = buildProviderEntries(adapters, [], []);
    const byId = new Map(entries.map((e) => [e.id, e]));

    // Claude OAuth has OmniRoute's 8 current models (including claude-opus-5-5)
    expect(byId.get("claude")?.defaultModels).toContain("claude-opus-5-5");
    expect(byId.get("claude")?.defaultModels.length).toBeGreaterThanOrEqual(8);

    // Codex has all 22 OmniRoute models (not truncated to 12, not 2 stubs)
    expect(byId.get("codex")?.defaultModels).toContain("gpt-6.1-sol");
    expect(byId.get("codex")?.defaultModels.length).toBeGreaterThanOrEqual(22);

    // Antigravity has all 39 OmniRoute models
    expect(byId.get("antigravity")?.defaultModels).toContain("gemini-4-carbon-high");
    expect(byId.get("antigravity")?.defaultModels.length).toBeGreaterThanOrEqual(39);

    // GitHub Copilot maps properly to OmniRoute 'github' registry (28 models)
    expect(byId.get("github-copilot")?.alias).toBe("gh");
    expect(byId.get("github-copilot")?.defaultModels).toContain("gpt-5.4");
    expect(byId.get("github-copilot")?.defaultModels.length).toBeGreaterThanOrEqual(28);

    // Kiro has all 48 OmniRoute models
    expect(byId.get("kiro")?.defaultModels).toContain("claude-opus-5.5");
    expect(byId.get("kiro")?.defaultModels.length).toBeGreaterThanOrEqual(48);

    // OpenAI & OpenAI Responses have all 21 OmniRoute models
    expect(byId.get("openai")?.defaultModels.length).toBeGreaterThanOrEqual(21);
    expect(byId.get("openai-responses")?.defaultModels.length).toBeGreaterThanOrEqual(21);

    // FreePI has all 26 OmniRoute models
    expect(byId.get("freepi")?.defaultModels.length).toBeGreaterThanOrEqual(26);
  });

  test("non-LLM embedding models are excluded from LLM provider catalog (e.g. openrouter)", () => {
    const openrouter = PROVIDER_REGISTRY_MAP.get("openrouter");
    expect(openrouter).toBeDefined();
    expect(openrouter?.defaultModels).not.toContain("openai/text-embedding-3-large");
  });

  test("registry exposes model display names and isFree metadata matching OmniRoute", () => {
    const claude = PROVIDER_REGISTRY_MAP.get("claude");
    const opus55 = claude?.models?.find((m) => m.id === "claude-opus-5-5");
    expect(opus55?.name).toBe("Claude Opus 5.5");

    const freepi = PROVIDER_REGISTRY_MAP.get("freepi");
    const freeModels = freepi?.models?.filter((m) => m.isFree) ?? [];
    expect(freeModels.length).toBeGreaterThan(0);
  });
});
