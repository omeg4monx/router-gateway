import fs from "node:fs";
import path from "node:path";
import { describe, expect, it } from "vitest";
import {
  STATUS_FILTER_OPTIONS,
  PROVIDER_STATUS_PILLS,
  CATEGORY_TABS,
  getConnectionStatus,
  matchesStatusFilter,
  getConnectionErrorTag,
  getProviderStats,
  resolveProviderAuthTypes,
  getProviderTypeLabel,
  matchesProviderSearch,
  buildProviderEntries,
  computeFleetTelemetry,
  countStatusBuckets,
  countCategoryBuckets,
  partitionConfiguredAndCatalog,
  sortFilteredEntries,
  countActiveConnectionsByMode,
} from "@/app/(dashboard)/dashboard/providers/utils.js";

const NOW = new Date("2026-11-22T12:00:00.000Z").getTime();

describe("providers-ux Task 1: pure utilities", () => {
  it("preserves STATUS_FILTER_OPTIONS and exposes extended PROVIDER_STATUS_PILLS and CATEGORY_TABS", () => {
    expect(STATUS_FILTER_OPTIONS.map((o) => o.value)).toEqual([
      "all",
      "active",
      "inactive",
      "none",
    ]);
    expect(PROVIDER_STATUS_PILLS.map((o) => o.value)).toEqual([
      "all",
      "active",
      "error",
      "inactive",
      "none",
    ]);
    expect(CATEGORY_TABS.map((t) => t.value)).toEqual([
      "all",
      "custom",
      "oauth",
      "apikey",
      "free",
    ]);
  });

  it("matchesStatusFilter supports error filter without breaking existing statuses", () => {
    expect(
      matchesStatusFilter("error", { total: 1, error: 1, allDisabled: false }, false),
    ).toBe(true);
    expect(
      matchesStatusFilter("error", { total: 1, error: 1, allDisabled: true }, false),
    ).toBe(false);
    expect(
      matchesStatusFilter("error", { total: 0, error: 1, allDisabled: false }, true),
    ).toBe(false);
    expect(
      matchesStatusFilter("error", { total: 1, error: 0, allDisabled: false }, false),
    ).toBe(false);
    expect(getConnectionStatus({ total: 1, allDisabled: false })).toBe("active");
  });

  it("getConnectionErrorTag classifies explicit types, HTTP codes, and fallback messages", () => {
    expect(getConnectionErrorTag(null)).toBeNull();
    expect(getConnectionErrorTag({ lastErrorType: "runtime_error" })).toBe("RUNTIME");
    expect(getConnectionErrorTag({ lastErrorType: "upstream_auth_error" })).toBe("AUTH");
    expect(getConnectionErrorTag({ lastErrorType: "token_expired" })).toBe("AUTH");
    expect(getConnectionErrorTag({ lastErrorType: "upstream_rate_limited" })).toBe("429");
    expect(getConnectionErrorTag({ lastErrorType: "upstream_unavailable" })).toBe("5XX");
    expect(getConnectionErrorTag({ lastErrorType: "network_error" })).toBe("NET");
    expect(getConnectionErrorTag({ errorCode: 404 })).toBe("404");
    expect(getConnectionErrorTag({ lastError: "invalid api key provided" })).toBe("AUTH");
  });

  it("getProviderStats handles expired vs active modelLock cooldowns and enabled counts", () => {
    const expiredLock = new Date(NOW - 60_000).toISOString();
    const activeLock = new Date(NOW + 60_000).toISOString();
    const connections = [
      {
        id: "c1",
        provider: "openai",
        authType: "apikey",
        isActive: true,
        testStatus: "unavailable",
        modelLock_gpt4: expiredLock,
      },
      {
        id: "c2",
        provider: "openai",
        authType: "apikey",
        isActive: true,
        testStatus: "unavailable",
        modelLock_gpt4: activeLock,
        lastErrorType: "upstream_rate_limited",
        lastErrorAt: "2026-11-22T11:55:00.000Z",
      },
      {
        id: "c3",
        provider: "openai",
        authType: "apikey",
        isActive: false,
        testStatus: "active",
      },
    ];

    const stats = getProviderStats(connections, "openai", "apikey", NOW);
    expect(stats.total).toBe(3);
    expect(stats.enabled).toBe(2);
    expect(stats.connected).toBe(1);
    expect(stats.error).toBe(1);
    expect(stats.errorCode).toBe("429");
    expect(stats.allDisabled).toBe(false);
    expect(stats.latestErrorAtMs).toBe(new Date("2026-11-22T11:55:00.000Z").getTime());
  });

  it("matchesProviderSearch searches name, id, alias, prefix, baseUrl, category, and typeLabel safely", () => {
    const meta = {
      name: "OpenAI Compatible (Prod)",
      id: "openai-compatible-prod",
      alias: "oc-prod",
      prefix: "oc-prod",
      baseUrl: "https://gateway.example.com/v1",
      typeLabel: "OpenAI Responses",
      category: "custom",
    };
    expect(matchesProviderSearch("", meta)).toBe(true);
    expect(matchesProviderSearch("   ", meta)).toBe(true);
    expect(matchesProviderSearch("prod", meta)).toBe(true);
    expect(matchesProviderSearch("openai-compatible-prod", meta)).toBe(true);
    expect(matchesProviderSearch("oc-prod", meta)).toBe(true);
    expect(matchesProviderSearch("gateway.example", meta)).toBe(true);
    expect(matchesProviderSearch("responses", meta)).toBe(true);
    expect(matchesProviderSearch("custom", meta)).toBe(true);
    expect(matchesProviderSearch("anthropic", meta)).toBe(false);
    expect(matchesProviderSearch("foo", {})).toBe(false);
  });

  it("resolveProviderAuthTypes and getProviderTypeLabel preserve dual-auth and compatible branding", () => {
    expect(resolveProviderAuthTypes({}, "kiro")).toEqual(["oauth", "apikey", "api_key"]);
    expect(resolveProviderAuthTypes({ authModes: ["oauth"] }, "claude")).toBe("oauth");
    expect(resolveProviderAuthTypes({ authModes: ["oauth", "apikey"] }, "xai")).toEqual([
      "oauth",
      "apikey",
      "api_key",
    ]);
    expect(
      getProviderTypeLabel(
        { name: "My Kimi", brand: "moonshot", apiType: "chat" },
        "openai-compatible-kimi",
        "compatible",
      ),
    ).toBe("MoonshotAI Compatible");
    expect(
      getProviderTypeLabel(
        { name: "Responses Proxy", apiType: "responses" },
        "openai-compatible-resp",
        "compatible",
      ),
    ).toBe("OpenAI Responses");
    expect(
      getProviderTypeLabel({ name: "Claude" }, "claude", "oauth"),
    ).toBe("OAuth 2.0");
  });

  it("partitionConfiguredAndCatalog places configured and 0-conn custom nodes in Zona 1 with zero Zona 2 overlap", () => {
    const entries = [
      {
        id: "claude",
        category: "oauth",
        info: { name: "Claude", priority: 1 },
        isNoAuth: false,
        stats: { total: 1, connected: 1, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
      {
        id: "codex",
        category: "oauth",
        info: { name: "Codex", priority: 2 },
        isNoAuth: false,
        stats: {
          total: 1,
          connected: 0,
          error: 1,
          allDisabled: false,
          latestErrorAtMs: NOW - 10_000,
        },
      },
      {
        id: "openai-compatible-new",
        category: "custom",
        info: { name: "Fresh Custom" },
        isNoAuth: false,
        stats: { total: 0, connected: 0, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
      {
        id: "openai",
        category: "apikey",
        info: { name: "OpenAI" },
        isNoAuth: false,
        stats: { total: 1, connected: 0, error: 0, allDisabled: true, latestErrorAtMs: 0 },
      },
      {
        id: "deepseek",
        category: "apikey",
        info: { name: "DeepSeek" },
        isNoAuth: false,
        stats: { total: 0, connected: 0, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
      {
        id: "gemini-cli",
        category: "oauth",
        info: { name: "Gemini CLI", priority: 3 },
        isNoAuth: false,
        stats: { total: 0, connected: 0, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
    ];

    const { configured, catalog } = partitionConfiguredAndCatalog(entries);
    expect(configured.map((e) => e.id)).toEqual([
      "codex", // errored first
      "claude", // active second
      "openai-compatible-new", // 0-conn custom node third
      "openai", // disabled last
    ]);

    const catalogIds = [
      ...catalog.oauth.map((e) => e.id),
      ...catalog.free.map((e) => e.id),
      ...catalog.apikey.map((e) => e.id),
    ];
    expect(catalogIds).toEqual(["gemini-cli", "deepseek"]);
    for (const id of configured.map((e) => e.id)) {
      expect(catalogIds).not.toContain(id);
    }
  });

  it("computeFleetTelemetry and countStatusBuckets maintain 1-to-1 parity with matchesStatusFilter", () => {
    const entries = [
      {
        id: "claude",
        category: "oauth",
        isNoAuth: false,
        stats: { total: 2, enabled: 2, connected: 2, error: 0, allDisabled: false },
      },
      {
        id: "codex",
        category: "oauth",
        isNoAuth: false,
        stats: { total: 1, enabled: 1, connected: 0, error: 1, allDisabled: false },
      },
      {
        id: "openai",
        category: "apikey",
        isNoAuth: false,
        stats: { total: 1, enabled: 0, connected: 0, error: 0, allDisabled: true },
      },
      {
        id: "opencode",
        category: "free",
        isNoAuth: true,
        stats: { total: 0, enabled: 0, connected: 0, error: 0, allDisabled: false },
      },
      {
        id: "deepseek",
        category: "apikey",
        isNoAuth: false,
        stats: { total: 0, enabled: 0, connected: 0, error: 0, allDisabled: false },
      },
    ];
    const connections = [
      { id: "1", isActive: true },
      { id: "2", isActive: true },
      { id: "3", isActive: true },
      { id: "4", isActive: false },
    ];

    const telemetry = computeFleetTelemetry(entries, connections);
    const buckets = countStatusBuckets(entries);
    const catBuckets = countCategoryBuckets(entries);

    expect(telemetry.activeProviders).toBe(buckets.active);
    expect(telemetry.attentionProviders).toBe(buckets.error);
    expect(telemetry.disabledProviders).toBe(buckets.inactive);
    expect(telemetry.configuredActiveProviders).toBe(2);
    expect(telemetry.noAuthProviders).toBe(1);
    expect(telemetry.enabledConnections).toBe(3);
    expect(telemetry.totalConnections).toBe(4);
    expect(buckets).toEqual({ all: 5, active: 3, error: 1, inactive: 1, none: 1 });
    expect(catBuckets).toEqual({ all: 5, custom: 0, oauth: 2, apikey: 2, free: 1 });
  });

  it("countActiveConnectionsByMode aligns with /api/providers/test-batch grouping", () => {
    const connections = [
      { id: "1", provider: "claude", authType: "oauth", isActive: true },
      { id: "2", provider: "kiro", authType: "oauth", isActive: true }, // in FREE_PROVIDERS or OAUTH
      { id: "3", provider: "openrouter", authType: "apikey", isActive: true }, // FREE_TIER_PROVIDERS (apikey)
      { id: "4", provider: "openai-compatible-1", authType: "apikey", isActive: true },
      { id: "5", provider: "openai", authType: "apikey", isActive: false },
    ];

    expect(countActiveConnectionsByMode(connections, "all")).toBe(4);
    expect(countActiveConnectionsByMode(connections, "compatible")).toBe(1);
    expect(countActiveConnectionsByMode(connections, "oauth")).toBeGreaterThanOrEqual(1);
    // FREE_TIER_PROVIDERS apikey connections (openrouter) must NOT count as mode "free"
    const onlyFreeTier = [
      { id: "3", provider: "openrouter", authType: "apikey", isActive: true },
    ];
    expect(countActiveConnectionsByMode(onlyFreeTier, "free")).toBe(0);
  });

  it("buildProviderEntries and sortFilteredEntries produce normalized, sorted provider entries", () => {
    const nodes = [
      {
        id: "openai-compatible-test",
        type: "openai-compatible",
        name: "My Gateway",
        prefix: "gw",
        baseUrl: "https://gw.test/v1",
        apiType: "chat",
      },
    ];
    const connections = [
      {
        id: "conn-1",
        provider: "openai-compatible-test",
        authType: "apikey",
        isActive: true,
        testStatus: "active",
      },
    ];
    const entries = buildProviderEntries(nodes, connections, NOW);
    const customEntry = entries.find((e) => e.id === "openai-compatible-test");
    expect(customEntry).toBeDefined();
    expect(customEntry.category).toBe("custom");
    expect(customEntry.stats.connected).toBe(1);
    expect(customEntry.prefix).toBe("gw");

    const sorted = sortFilteredEntries([
      {
        id: "idle",
        info: { name: "Alpha" },
        isNoAuth: false,
        stats: { total: 0, connected: 0, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
      {
        id: "err",
        info: { name: "Zeta" },
        isNoAuth: false,
        stats: { total: 1, connected: 0, error: 1, allDisabled: false, latestErrorAtMs: NOW },
      },
      {
        id: "ok",
        info: { name: "Beta" },
        isNoAuth: false,
        stats: { total: 1, connected: 1, error: 0, allDisabled: false, latestErrorAtMs: 0 },
      },
    ]);
    expect(sorted.map((e) => e.id)).toEqual(["err", "ok", "idle"]);
  });
});

const readSrc = (rel) =>
  fs.readFileSync(path.join(process.cwd(), rel), "utf8");

describe("providers-ux Task 2: AddCompatibleModal & ModelAvailabilityBadge", () => {
  it("AddCompatibleModal uses sticky Modal footer, toasts on creation, and does not wipe baseUrl on apiType change", () => {
    const src = readSrc(
      "src/app/(dashboard)/dashboard/providers/components/AddCompatibleModal.js",
    );
    expect(src).toMatch(/<Modal[\s\S]*footer=/);
    expect(src).toContain("useNotificationStore");
    expect(src).not.toContain("config.hasApiType ? formData.apiType : isOpen");
    expect(src).toMatch(/useEffect\([\s\S]*\[isOpen,\s*variant\]/);
  });

  it("ModelAvailabilityBadge renders an active trigger button, uses shell-popover, supports Escape, guards document.hidden, and invokes onCooldownCleared", () => {
    const src = readSrc(
      "src/app/(dashboard)/dashboard/providers/components/ModelAvailabilityBadge.js",
    );
    // Strip block comments to ensure trigger button is genuinely uncommented
    const withoutBlockComments = src.replace(/\/\*[\s\S]*?\*\//g, "");
    expect(withoutBlockComments).toMatch(/<button[\s\S]*setExpanded/);
    expect(withoutBlockComments).toContain("shell-popover");
    expect(withoutBlockComments).toContain("Escape");
    expect(withoutBlockComments).toContain("document.hidden");
    expect(withoutBlockComments).toContain("onCooldownCleared");
  });
});
