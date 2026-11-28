import { describe, expect, test } from "bun:test";
import type { UsageLog } from "@/features/usage/schemas";
import {
  appendMonitorEntry,
  appendSystemLogBatch,
  appendSystemLogEntry,
  computeLiveRps,
  computeMonitorStats,
  expireStaleInFlight,
  filterMonitorEntries,
  formatShortCost,
  formatShortProtocol,
  formatShortRequestId,
  formatSystemLogAttrs,
  MAX_MONITOR_LINES,
  mergeMonitorEntries,
  parseSystemLogAttrEntries,
  serializeMonitorEntriesNDJSON,
  trafficStartToMonitorEntry,
  upsertMonitorEntries,
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
      liveCount: 0,
      okCount: 2,
      errorCount: 1,
      avgLatencyMs: 400,
    });
    expect(computeMonitorStats([])).toEqual({
      total: 0,
      liveCount: 0,
      okCount: 0,
      errorCount: 0,
      avgLatencyMs: 0,
    });
  });

  test("deduplicates system log entries on reconnect replay and caps at MAX_MONITOR_LINES", () => {
    let seq = 0;
    const nextSeq = () => {
      seq += 1;
      return seq;
    };
    const entry = {
      time: "2026-11-28T07:00:00.000Z",
      level: "INFO",
      message: "▶ gateway request",
      attrs: { provider: "openai", model: "gpt-4o" },
    };

    const first = appendSystemLogEntry([], entry, nextSeq);
    expect(first).toHaveLength(1);
    const deduped = appendSystemLogEntry(
      first,
      { ...entry, attrs: { model: "gpt-4o", provider: "openai" } },
      nextSeq,
    );
    expect(deduped).toHaveLength(1);

    let list = first;
    for (let i = 0; i < MAX_MONITOR_LINES + 5; i += 1) {
      list = appendSystemLogEntry(
        list,
        {
          time: `2026-11-28T07:00:${String(i).padStart(2, "0")}Z`,
          level: "INFO",
          message: `msg-${i}`,
        },
        nextSeq,
      );
    }
    expect(list).toHaveLength(MAX_MONITOR_LINES);
    expect(list[list.length - 1]?.message).toBe(`msg-${MAX_MONITOR_LINES + 4}`);
  });

  test("formats structured slog attrs cleanly for inline terminal display", () => {
    expect(formatSystemLogAttrs(undefined)).toBe("");
    expect(formatSystemLogAttrs({})).toBe("");
    expect(
      formatSystemLogAttrs({
        provider: "openai",
        status: 200,
        stream: true,
        error: "",
      }),
    ).toBe("provider=openai status=200 stream=true");
  });

  test("merges monitor entries chronologically without duplicates and preserves cap", () => {
    const older = usageLogToMonitorEntry(
      makeUsageLog({ id: 1, created_at: "2026-11-28T07:00:01Z" }),
    );
    const middle = usageLogToMonitorEntry(
      makeUsageLog({ id: 2, created_at: "2026-11-28T07:00:02Z" }),
    );
    const newer = usageLogToMonitorEntry(
      makeUsageLog({ id: 3, created_at: "2026-11-28T07:00:03Z" }),
    );

    const merged = mergeMonitorEntries([older, newer], [middle, newer]);
    expect(merged.map((e) => e.id)).toEqual(["usage-1", "usage-2", "usage-3"]);
  });

  test("tracks in-flight traffic and upgrades it in-place on completion without duplicates", () => {
    const started = trafficStartToMonitorEntry({
      request_id: "req_live_99",
      provider_id: "anthropic",
      connection_id: "conn_ant",
      model: "claude-sonnet-4",
      protocol: "anthropic",
      stream: true,
      started_at: "2026-11-28T07:10:00.000Z",
    });
    expect(started.inFlight).toBe(true);
    expect(started.level).toBe("live");
    expect(started.id).toBe("req-req_live_99");

    const finished = usageLogToMonitorEntry(
      makeUsageLog({
        id: 0,
        request_id: "req_live_99",
        provider_id: "anthropic",
        model: "claude-sonnet-4",
        status_code: 200,
        latency_ms: 850,
        created_at: "2026-11-28T07:10:00.850Z",
      }),
    );
    const afterFinish = upsertMonitorEntries([started], [finished]);
    expect(afterFinish).toHaveLength(1);
    expect(afterFinish[0]?.inFlight).toBe(false);
    expect(afterFinish[0]?.statusCode).toBe(200);
    expect(afterFinish[0]?.latencyMs).toBe(850);

    // Persisted SQLite usage event with same request_id must not create a duplicate row.
    const persisted = usageLogToMonitorEntry(
      makeUsageLog({
        id: 42,
        request_id: "req_live_99",
        provider_id: "anthropic",
        model: "claude-sonnet-4",
        status_code: 200,
        latency_ms: 850,
        created_at: "2026-11-28T07:10:00.850Z",
      }),
    );
    const afterPersist = upsertMonitorEntries(afterFinish, [persisted]);
    expect(afterPersist).toHaveLength(1);
  });

  test("batches system log entries with O(1) fingerprint dedup and computes live RPS & NDJSON", () => {
    let seq = 0;
    const nextSeq = () => {
      seq += 1;
      return seq;
    };
    const e1 = { time: "2026-11-28T07:10:00Z", level: "INFO", message: "one" };
    const e2 = { time: "2026-11-28T07:10:01Z", level: "WARN", message: "two" };

    const batched = appendSystemLogBatch([], [e1, e1, e2], nextSeq);
    expect(batched).toHaveLength(2);
    expect(batched[0]?.fingerprint).toBeTruthy();

    const nowMs = Date.parse("2026-11-28T07:10:10.000Z");
    const recent = [
      usageLogToMonitorEntry(makeUsageLog({ id: 1, created_at: "2026-11-28T07:10:05.000Z" })),
      usageLogToMonitorEntry(makeUsageLog({ id: 2, created_at: "2026-11-28T07:10:08.000Z" })),
      usageLogToMonitorEntry(makeUsageLog({ id: 3, created_at: "2026-11-28T07:09:00.000Z" })),
    ];
    expect(computeLiveRps(recent, nowMs, 10_000)).toBeCloseTo(0.2, 2);

    const ndjson = serializeMonitorEntriesNDJSON(recent.slice(0, 2));
    const lines = ndjson.trim().split("\n");
    expect(lines).toHaveLength(2);
    expect(JSON.parse(lines[0] ?? "{}").provider).toBe("openai");
  });

  test("formats terminal protocol tags, short request IDs, cost, and structured attr pairs", () => {
    expect(formatShortProtocol("openai-chat")).toBe("oai");
    expect(formatShortProtocol("openai-responses")).toBe("resp");
    expect(formatShortProtocol("anthropic")).toBe("anth");
    expect(formatShortProtocol("gemini")).toBe("gem");
    expect(formatShortProtocol("")).toBe("");

    expect(formatShortRequestId("req_01928374abcdef")).toBe("#abcdef");
    expect(formatShortRequestId("abc")).toBe("#abc");
    expect(formatShortRequestId("")).toBe("");

    expect(formatShortCost(0)).toBe("");
    expect(formatShortCost(0.00001)).toBe("<$0.0001");
    expect(formatShortCost(0.0042)).toBe("$0.0042");

    expect(parseSystemLogAttrEntries({ provider: "openai", status: 200, empty: "" })).toEqual([
      { key: "provider", value: "openai" },
      { key: "status", value: "200" },
    ]);
  });

  test("expires stale in-flight entries after maxAgeMs so zombie LIVE rows never linger", () => {
    const stale = trafficStartToMonitorEntry({
      request_id: "req_zombie",
      provider_id: "openai",
      model: "gpt-4o",
      started_at: "2026-11-28T07:00:00.000Z",
    });
    const fresh = trafficStartToMonitorEntry({
      request_id: "req_active",
      provider_id: "openai",
      model: "gpt-4o",
      started_at: "2026-11-28T07:04:30.000Z",
    });

    const nowMs = Date.parse("2026-11-28T07:05:10.000Z");
    const result = expireStaleInFlight([stale, fresh], nowMs, 300_000);
    expect(result[0]?.inFlight).toBe(false);
    expect(result[0]?.statusCode).toBe(499);
    expect(result[0]?.level).toBe("warn");
    expect(result[1]?.inFlight).toBe(true);
  });
});
