import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

describe("UsageConsole realtime isolation and range sync", () => {
  const consoleSource = readFileSync(resolve(import.meta.dir, "UsageConsole.tsx"), "utf8");
  const indicatorSource = readFileSync(
    resolve(import.meta.dir, "UsageStreamIndicator.tsx"),
    "utf8",
  );

  test("isolates useUsageStream inside UsageStreamIndicator and mounts it exactly once", () => {
    expect(consoleSource).not.toContain("useUsageStream");
    expect(indicatorSource).toContain("useUsageStream(true)");
    const mountCount = (consoleSource.match(/<UsageStreamIndicator\b/g) ?? []).length;
    expect(mountCount).toBe(1);
  });

  test("synchronizes selected range into logsFilter", () => {
    expect(consoleSource).toMatch(/logsFilter\s*=\s*useMemo\(\s*\(\)\s*=>\s*\(\{\s*range,/);
  });

  test("lazy-loads UsagePricingModal and uses stable refetch callbacks", () => {
    expect(consoleSource).toMatch(/const UsagePricingModal = lazy\(/);
    expect(consoleSource).not.toContain("[summaryQuery, logsQuery]");
    expect(consoleSource).not.toContain("[logsQuery]");
  });
});
