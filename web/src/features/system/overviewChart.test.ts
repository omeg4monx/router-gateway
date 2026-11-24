import { describe, expect, test } from "bun:test";
import type { UsageLog, UsageSeriesBucket } from "@/features/usage/schemas";
import {
  buildRequestChartModel,
  filterRecentRequests,
  formatRelativeTimeShort,
  isUsageLogSuccess,
} from "./overviewChart";

const SAMPLE_SERIES: UsageSeriesBucket[] = [
  {
    bucket: "2026-11-24T10:00:00Z",
    requests: 4,
    input_tokens: 1200,
    output_tokens: 400,
    cached_tokens: 200,
    cost_usd: 0.012,
    error_count: 0,
  },
  {
    bucket: "2026-11-24T11:00:00Z",
    requests: 12,
    input_tokens: 5000,
    output_tokens: 2500,
    cached_tokens: 500,
    cost_usd: 0.045,
    error_count: 1,
  },
  {
    bucket: "2026-11-24T12:00:00Z",
    requests: 6,
    input_tokens: 2000,
    output_tokens: 1000,
    cached_tokens: 0,
    cost_usd: 0.02,
    error_count: 0,
  },
];

function makeLog(overrides: Partial<UsageLog>): UsageLog {
  return {
    id: 1,
    request_id: "req_1",
    key_id: "key_1",
    model: "gpt-4o",
    provider_id: "openai",
    connection_id: "conn_1",
    protocol: "openai-chat",
    input_tokens: 120,
    output_tokens: 45,
    cached_tokens: 0,
    cost_usd: 0.001,
    latency_ms: 420,
    ttft_ms: 110,
    status_code: 200,
    error_summary: "",
    stream: true,
    created_at: "2026-11-24T18:50:00Z",
    ...overrides,
  };
}

describe("features/system/overviewChart", () => {
  test("buildRequestChartModel computes SVG geometry, ticks, totals, and peak across metrics", () => {
    const reqModel = buildRequestChartModel(SAMPLE_SERIES, "24h", "requests");
    expect(reqModel.points).toHaveLength(3);
    expect(reqModel.totalValue).toBe(22);
    expect(reqModel.peakPoint?.value).toBe(12);
    expect(reqModel.linePath.startsWith("M ")).toBe(true);
    expect(reqModel.areaPath.endsWith(" Z")).toBe(true);
    expect(reqModel.yTicks.length).toBeGreaterThanOrEqual(2);
    for (const pt of reqModel.points) {
      expect(Number.isFinite(pt.x)).toBe(true);
      expect(Number.isFinite(pt.y)).toBe(true);
    }

    const tokModel = buildRequestChartModel(SAMPLE_SERIES, "24h", "tokens");
    expect(tokModel.totalValue).toBe(1800 + 8000 + 3000);
    expect(tokModel.peakPoint?.value).toBe(8000);

    const costModel = buildRequestChartModel(SAMPLE_SERIES, "24h", "cost");
    expect(costModel.totalValue).toBeCloseTo(0.077, 6);
    expect(costModel.peakPoint?.value).toBeCloseTo(0.045, 6);
  });

  test("buildRequestChartModel handles empty and single-point series without NaN coordinates", () => {
    const empty = buildRequestChartModel([], "7d", "requests");
    expect(empty.points).toEqual([]);
    expect(empty.linePath).toBe("");
    expect(empty.areaPath).toBe("");
    expect(empty.peakPoint).toBeNull();

    const firstSample = SAMPLE_SERIES[0];
    expect(firstSample).toBeDefined();
    if (!firstSample) return;

    const single = buildRequestChartModel([firstSample], "24h", "requests");
    expect(single.points).toHaveLength(1);
    expect(Number.isFinite(single.points[0]?.x)).toBe(true);
    expect(Number.isFinite(single.points[0]?.y)).toBe(true);
    expect(single.linePath.startsWith("M ")).toBe(true);
    expect(single.areaPath.endsWith(" Z")).toBe(true);
  });

  test("formatRelativeTimeShort formats seconds, minutes, hours, and days in EN and ID", () => {
    const now = Date.parse("2026-11-24T19:00:00Z");
    expect(formatRelativeTimeShort("2026-11-24T18:59:58Z", now, false)).toBe("2s ago");
    expect(formatRelativeTimeShort("2026-11-24T18:59:58Z", now, true)).toBe("2d lalu");
    expect(formatRelativeTimeShort("2026-11-24T18:55:00Z", now, false)).toBe("5m ago");
    expect(formatRelativeTimeShort("2026-11-24T18:55:00Z", now, true)).toBe("5m lalu");
    expect(formatRelativeTimeShort("2026-11-24T16:00:00Z", now, false)).toBe("3h ago");
    expect(formatRelativeTimeShort("2026-11-24T16:00:00Z", now, true)).toBe("3j lalu");
    expect(formatRelativeTimeShort("2026-11-22T19:00:00Z", now, false)).toBe("2d ago");
    expect(formatRelativeTimeShort("2026-11-22T19:00:00Z", now, true)).toBe("2h lalu");
    expect(formatRelativeTimeShort("invalid", now, false)).toBe("—");
  });

  test("isUsageLogSuccess and filterRecentRequests separate OK and error logs accurately", () => {
    const okLog = makeLog({ id: 3, status_code: 200, error_summary: "" });
    const rateLimitedLog = makeLog({ id: 2, status_code: 429, error_summary: "rate limit" });
    const abortedLog = makeLog({ id: 1, status_code: 200, error_summary: "stream aborted" });
    const logs = [okLog, rateLimitedLog, abortedLog];

    expect(isUsageLogSuccess(okLog)).toBe(true);
    expect(isUsageLogSuccess(rateLimitedLog)).toBe(false);
    expect(isUsageLogSuccess(abortedLog)).toBe(false);

    const allRes = filterRecentRequests(logs, "all", 2);
    expect(allRes.total).toBe(3);
    expect(allRes.okCount).toBe(1);
    expect(allRes.errorCount).toBe(2);
    expect(allRes.items.map((l) => l.id)).toEqual([3, 2]);

    const okRes = filterRecentRequests(logs, "ok");
    expect(okRes.items.map((l) => l.id)).toEqual([3]);

    const errRes = filterRecentRequests(logs, "error");
    expect(errRes.items.map((l) => l.id)).toEqual([2, 1]);
  });
});
