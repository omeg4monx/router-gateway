import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import {
  filterModules,
  formatActivityOriginLabel,
  formatClockUptime,
  formatUsageBucketLabel,
  generateEndpointSnippet,
  getModuleMetadata,
  getOverviewCopy,
} from "./overview";

describe("features/system/overview", () => {
  test("provides localized overview copy for en and id", () => {
    const en = getOverviewCopy("en");
    const id = getOverviewCopy("id");

    expect(en.title).not.toBe(id.title);
    expect(en.quickActionsHeading).not.toBe(id.quickActionsHeading);
    expect(en.modulesTitle).not.toBe(id.modulesTitle);
    expect(en.postureTitle).not.toBe(id.postureTitle);
  });

  test("maps known kernel modules to descriptions and dashboard links, with safe fallback", () => {
    const providers = getModuleMetadata("providers", "en");
    expect(providers.href).toBe("/dashboard/providers/");
    expect(providers.description.length).toBeGreaterThan(10);

    const plugins = getModuleMetadata("plugins", "id");
    expect(plugins.href).toBe("/dashboard/plugins/");

    const unknown = getModuleMetadata("custom_experimental", "en");
    expect(unknown.href).toBeNull();
    expect(unknown.description.length).toBeGreaterThan(5);
  });

  test("filters modules across module IDs and bilingual descriptions", () => {
    const sample = ["providers", "keys", "telegram", "security", "tokensaver"];
    expect(filterModules(sample, "", "en")).toEqual(sample);
    expect(filterModules(sample, "bot", "en")).toEqual(["telegram"]);
    expect(filterModules(sample, "keamanan", "id")).toEqual(["security"]);
    expect(filterModules(sample, "nonexistent", "en")).toEqual([]);
  });

  test("formats clock uptime cleanly across zero, hours, and days", () => {
    expect(formatClockUptime(-10)).toBe("00h 00m 00s");
    expect(formatClockUptime(3665)).toBe("01h 01m 05s");
    expect(formatClockUptime(90061)).toBe("1d 01h 01m 01s");
  });

  test("generates OpenAI and Anthropic snippets for curl, python, and node", () => {
    const oaiCurl = generateEndpointSnippet({
      format: "openai",
      lang: "curl",
      baseUrl: "http://127.0.0.1:20128/",
      apiKey: "sk-test-123",
      model: "gpt-4o",
    });
    expect(oaiCurl).toContain("http://127.0.0.1:20128/v1/chat/completions");
    expect(oaiCurl).toContain("Bearer sk-test-123");

    const claudePy = generateEndpointSnippet({
      format: "anthropic",
      lang: "python",
      baseUrl: "http://127.0.0.1:20128",
      model: "claude-sonnet-4-5",
    });
    expect(claudePy).toContain("Anthropic(");
    expect(claudePy).toContain("claude-sonnet-4-5");
  });

  test("formats 24h hourly buckets with time and multi-day buckets with date", () => {
    const hourly = formatUsageBucketLabel("2026-11-24T17:00:00Z", "24h");
    expect(hourly).toMatch(/\d{1,2}[:.]\d{2}/);

    const daily = formatUsageBucketLabel("2026-11-24", "7d");
    expect(daily).not.toMatch(/^\d{1,2}[:.]\d{2}$/);
    expect(daily.length).toBeGreaterThan(2);
  });

  test("labels internal admin/probe usage origins and falls back to API key names", () => {
    const keyMap = new Map([["key_abc", "Production CLI"]]);
    expect(formatActivityOriginLabel("key_abc", keyMap, false)).toBe("Production CLI");
    expect(formatActivityOriginLabel("admin-internal", keyMap, false)).toBe("Admin Studio");
    expect(formatActivityOriginLabel("admin-internal", keyMap, true)).toBe("Studio Admin");
    expect(formatActivityOriginLabel("admin-probe", keyMap, false)).toBe("Model Probe");
    expect(formatActivityOriginLabel("admin-probe", keyMap, true)).toBe("Uji Model");
  });

  test("SystemOverview wires live SSE usage stream, /dashboard/combos/ link, and full query error aggregation", () => {
    const overviewSource = readFileSync(resolve(import.meta.dirname, "SystemOverview.tsx"), "utf8");
    const realtimeSource = readFileSync(
      resolve(import.meta.dirname, "OverviewRealtimeTraffic.tsx"),
      "utf8",
    );
    expect(overviewSource).toContain("useUsageStream(true)");
    expect(overviewSource).toContain('to="/dashboard/combos/"');
    expect(overviewSource).toContain("usageQuery.isError");
    expect(overviewSource).toContain("providersQuery.isError");
    expect(overviewSource).toContain("OverviewRealtimeTraffic");
    expect(overviewSource).toContain("groupTopologyProviders");
    expect(realtimeSource).toContain("./ProviderTopology");
  });
});
