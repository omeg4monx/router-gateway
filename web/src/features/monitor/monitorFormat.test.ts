import { describe, expect, test } from "bun:test";
import type { UsageLog } from "@/features/usage/schemas";
import {
  appendMonitorEntry,
  computeMonitorStats,
  filterMonitorEntries,
  MAX_MONITOR_LINES,
  usageLogToMonitorEntry,
} from "./monitorFormat";

function makeUsageLog(overrides: Partial<UsageLog> = {}): UsageLog {
  return {
    id: 1,
    request_id: "req_abc",
    key_id: "key_1",
    model: "gpt-4o",
    provider_id: "openai",
    connection_id: "conn_1",
    protocol: "openai-chat",
    input_tokens: 1200,
    output_tokens: 340,
    cached_tokens: 200,
    cost_usd: 0.0042,
    latency_ms: 420,
    ttft_ms: 85,
    status_code: 200,
    error_summary: "",
    stream: true,
    created_at: "2026-11-26T01:55:12.000Z",
    ...overrides,
  };
}

describe("features/monitor/monitorFormat", () => {
  test("maps UsageLog to structured terminal MonitorEntry with status severity", () => {
    const okEntry = usageLogToMonitorEntry(makeUsageLog());
    expect(okEntry.id).toBe("usage-1");
    expect(okEntry.level).toBe("ok");
    expect(okEntry.provider).toBe("openai");
    expect(okEntry.model).toBe("gpt-4o");
    expect(okEntry.statusCode).toBe(200);
    expect(okEntry.stream).toBe(true);

    const rateLimited = usageLogToMonitorEntry(makeUsageLog({ id: 2, status_code: 429 }));
    expect(rateLimited.level).toBe("warn");

    const failed = usageLogToMonitorEntry(
      makeUsageLog({ id: 3, status_code: 502, error_summary: "upstream timeout" }),
    );
    expect(failed.level).toBe("error");
    expect(failed.errorSummary).toBe("upstream timeout");
  });

  test("appends entries with deduplication and caps ring buffer at MAX_MONITOR_LINES", () => {
    const first = usageLogToMonitorEntry(makeUsageLog({ id: 10 }));
    const dup = usageLogToMonitorEntry(makeUsageLog({ id: 10 }));
    expect(appendMonitorEntry([first], dup)).toHaveLength(1);

    let list = [] as ReturnType<typeof usageLogToMonitorEntry>[];
    for (let i = 1; i <= MAX_MONITOR_LINES + 15; i += 1) {
      list = appendMonitorEntry(list, usageLogToMonitorEntry(makeUsageLog({ id: i })));
    }
    expect(list).toHaveLength(MAX_MONITOR_LINES);
    expect(list[0]?.id).toBe("usage-16");
    expect(list[list.length - 1]?.id).toBe(`usage-${MAX_MONITOR_LINES + 15}`);
  });

  test("filters entries by status bucket and search query", () => {
    const entries = [
      usageLogToMonitorEntry(makeUsageLog({ id: 1, provider_id: "openai", model: "gpt-4o" })),
      usageLogToMonitorEntry(
        makeUsageLog({
          id: 2,
          provider_id: "antigravity",
          model: "gemini-2.5-pro",
          status_code: 429,
          error_summary: "quota exhausted",
        }),
      ),
      usageLogToMonitorEntry(
        makeUsageLog({
          id: 3,
          provider_id: "anthropic",
          model: "claude-3-7-sonnet",
          status_code: 502,
          error_summary: "bad gateway",
        }),
      ),
    ];

    expect(filterMonitorEntries(entries, "all", "")).toHaveLength(3);
    expect(filterMonitorEntries(entries, "ok", "")).toHaveLength(1);
    expect(filterMonitorEntries(entries, "error", "")).toHaveLength(2);
    expect(filterMonitorEntries(entries, "all", "antigravity")).toHaveLength(1);
    expect(filterMonitorEntries(entries, "error", "gateway")).toHaveLength(1);
  });

  test("computes summary stats accurately", () => {
    const entries = [
      usageLogToMonitorEntry(makeUsageLog({ id: 1, status_code: 200, latency_ms: 200 })),
      usageLogToMonitorEntry(makeUsageLog({ id: 2, status_code: 200, latency_ms: 400 })),
      usageLogToMonitorEntry(makeUsageLog({ id: 3, status_code: 500, latency_ms: 600 })),
    ];
    expect(computeMonitorStats(entries)).toEqual({
      total: 3,
      okCount: 2,
      errorCount: 1,
      avgLatencyMs: 400,
    });
    expect(computeMonitorStats([])).toEqual({
      total: 0,
      okCount: 0,
      errorCount: 0,
      avgLatencyMs: 0,
    });
  });
});
