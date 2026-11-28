import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { formatActivityOriginLabel, formatClockUptime } from "./overview";

describe("features/system/overview", () => {
  test("formats clock uptime cleanly across zero, hours, and days", () => {
    expect(formatClockUptime(-10)).toBe("00h 00m 00s");
    expect(formatClockUptime(3665)).toBe("01h 01m 05s");
    expect(formatClockUptime(90061)).toBe("1d 01h 01m 01s");
  });

  test("labels internal admin/probe usage origins and falls back to API key names", () => {
    const keyMap = new Map([["key_abc", "Production CLI"]]);
    expect(formatActivityOriginLabel("key_abc", keyMap, false)).toBe("Production CLI");
    expect(formatActivityOriginLabel("admin-internal", keyMap, false)).toBe("Admin Studio");
    expect(formatActivityOriginLabel("admin-internal", keyMap, true)).toBe("Studio Admin");
    expect(formatActivityOriginLabel("admin-probe", keyMap, false)).toBe("Model Probe");
    expect(formatActivityOriginLabel("admin-probe", keyMap, true)).toBe("Uji Model");
  });

  test("overview.ts stays free of unused legacy registry/copy/snippet dead code", () => {
    const helperSource = readFileSync(resolve(import.meta.dirname, "overview.ts"), "utf8");
    expect(helperSource).not.toContain("MODULE_REGISTRY");
    expect(helperSource).not.toContain("OVERVIEW_COPY");
    expect(helperSource).not.toContain("getOverviewCopy");
    expect(helperSource).not.toContain("generateEndpointSnippet");
    expect(helperSource).not.toContain("filterModules");
    expect(helperSource).not.toContain("getModuleMetadata");
  });

  test("SystemOverview wires /dashboard/combos/ link, catalog refresh, and full query error aggregation while isolating SSE in OverviewRealtimeTraffic", () => {
    const overviewSource = readFileSync(resolve(import.meta.dirname, "SystemOverview.tsx"), "utf8");
    const realtimeSource = readFileSync(
      resolve(import.meta.dirname, "OverviewRealtimeTraffic.tsx"),
      "utf8",
    );
    expect(overviewSource).not.toContain("useUsageStream");
    expect(realtimeSource).toContain("useUsageStream(true)");
    expect(overviewSource).toContain('to="/dashboard/combos/"');
    expect(overviewSource).toContain("usageQuery.isError");
    expect(overviewSource).toContain("providersQuery.isError");
    expect(overviewSource).toContain("catalogQuery.refetch()");
    expect(overviewSource).toContain("OverviewRealtimeTraffic");
    expect(overviewSource).toContain("groupTopologyProviders");
    expect(realtimeSource).toContain("./ProviderTopology");
    expect(realtimeSource).not.toContain("Activity Trend");
    expect(realtimeSource).not.toContain("Tren Aktivitas");
  });

  test("useUsageSummary keeps previous data across range switches and useProvidersCatalog caches static catalog", () => {
    const usageHooksSource = readFileSync(
      resolve(import.meta.dirname, "../usage/hooks.ts"),
      "utf8",
    );
    const providersHooksSource = readFileSync(
      resolve(import.meta.dirname, "../providers/hooks.ts"),
      "utf8",
    );
    const summaryIdx = usageHooksSource.indexOf("export function useUsageSummary");
    const logsIdx = usageHooksSource.indexOf("export function useUsageLogs");
    expect(usageHooksSource.slice(summaryIdx, logsIdx)).toContain(
      "placeholderData: keepPreviousData",
    );

    const catalogIdx = providersHooksSource.indexOf("export function useProvidersCatalog");
    const nodesIdx = providersHooksSource.indexOf("export function useProviderNodes");
    expect(providersHooksSource.slice(catalogIdx, nodesIdx)).toContain("staleTime:");
  });

  test("HostTelemetryCard isolates 2s telemetry polling away from SystemOverview root renders", () => {
    const overviewSource = readFileSync(resolve(import.meta.dirname, "SystemOverview.tsx"), "utf8");
    const hostCardSource = readFileSync(
      resolve(import.meta.dirname, "HostTelemetryCard.tsx"),
      "utf8",
    );
    expect(overviewSource).not.toContain("useSystemTelemetry");
    expect(hostCardSource).toContain("useSystemTelemetry");
    expect(hostCardSource).toContain("Host & OS Specifications");
    expect(hostCardSource).toContain("motion-reduce:animate-none");

    const hostIdx = overviewSource.indexOf("<HostTelemetryCard");
    const realtimeIdx = overviewSource.indexOf("<OverviewRealtimeTraffic");
    const spotlightIdx = overviewSource.indexOf("Most Used Model");

    expect(hostIdx).toBeGreaterThan(-1);
    expect(realtimeIdx).toBeGreaterThan(hostIdx);
    expect(spotlightIdx).toBeGreaterThan(realtimeIdx);
  });

  test("OverviewRealtimeTraffic is memoized, bounds column heights, and exposes accessible filter/motion semantics", () => {
    const realtimeSource = readFileSync(
      resolve(import.meta.dirname, "OverviewRealtimeTraffic.tsx"),
      "utf8",
    );
    expect(realtimeSource).toContain("memo(function OverviewRealtimeTraffic");
    expect(realtimeSource).toContain("max-h-[360px]");
    expect(realtimeSource).toContain("sm:max-h-[410px]");
    expect(realtimeSource).not.toContain("h-[360px] flex-1 overflow-y-auto");
    expect(realtimeSource).toContain("aria-pressed={active}");
    expect(realtimeSource).toContain("motion-reduce:animate-none");
  });

  test("ProviderTopology uses defined text tokens, debounces ResizeObserver fitView, and respects reduced motion", () => {
    const topologySource = readFileSync(
      resolve(import.meta.dirname, "ProviderTopology.tsx"),
      "utf8",
    );
    expect(topologySource).not.toContain("--color-text-secondary");
    expect(topologySource).toContain("scheduleFitView()");
    expect(topologySource).toContain("prefers-reduced-motion: reduce");
    expect(topologySource).not.toMatch(/new ResizeObserver\(\(\) => \{\s*void rfInstanceRef/);
  });

  test("SystemOverview clarifies Spotlight #3 recent-log sampling and preserves 3 tinted KPI cards with AnimatedNumber", () => {
    const overviewSource = readFileSync(resolve(import.meta.dirname, "SystemOverview.tsx"), "utf8");
    const globalsCss = readFileSync(resolve(import.meta.dirname, "../../app/globals.css"), "utf8");

    expect(overviewSource).not.toContain("Estimated Spend");
    expect(overviewSource).not.toContain("Estimasi Biaya");
    expect(overviewSource).toContain("overview-kpi-requests");
    expect(overviewSource).toContain("overview-kpi-tokens");
    expect(overviewSource).toContain("overview-kpi-cached");
    expect(overviewSource).toContain("<AnimatedNumber");
    expect(overviewSource).toContain("100 req");

    const kpiStart = overviewSource.indexOf("Row 1:");
    const kpiEnd = overviewSource.indexOf("Row 2:");
    expect(kpiStart).toBeGreaterThan(-1);
    expect(kpiEnd).toBeGreaterThan(kpiStart);
    const kpiSection = overviewSource.slice(kpiStart, kpiEnd);
    expect(kpiSection).not.toContain("formatCompactTokens");

    expect(globalsCss).toContain(".overview-kpi-requests");
    expect(globalsCss).toContain(".overview-kpi-tokens");
    expect(globalsCss).toContain(".overview-kpi-cached");
    expect(globalsCss).toContain("@keyframes metric-value-tick");
  });
});
