import { describe, expect, test } from "bun:test";
import type { UsageLog, UsageLogsResponse } from "./schemas";
import { prependUsageLogToResponse } from "./streamCache";

function makeLog(overrides: Partial<UsageLog>): UsageLog {
  return {
    id: 1,
    request_id: "req_1",
    key_id: "key_1",
    model: "gpt-4o",
    provider_id: "openai",
    connection_id: "conn_1",
    protocol: "openai-chat",
    input_tokens: 100,
    output_tokens: 50,
    cached_tokens: 0,
    cost_usd: 0.001,
    latency_ms: 420,
    ttft_ms: 120,
    status_code: 200,
    error_summary: "",
    stream: true,
    created_at: "2026-11-24T19:00:00Z",
    ...overrides,
  };
}

describe("features/usage/streamCache", () => {
  test("prepends a new log, increments total, respects limit, and deduplicates by id", () => {
    const initial: UsageLogsResponse = {
      logs: [makeLog({ id: 2 }), makeLog({ id: 1 })],
      total: 2,
      limit: 2,
      offset: 0,
    };

    const incoming = makeLog({ id: 3, model: "claude-sonnet-4" });
    const next = prependUsageLogToResponse(initial, incoming, { limit: 2 });
    expect(next.logs.map((l) => l.id)).toEqual([3, 2]);
    expect(next.total).toBe(3);

    // Duplicate SSE frame for id=3 must be a no-op
    const deduped = prependUsageLogToResponse(next, incoming, { limit: 2 });
    expect(deduped).toBe(next);
  });

  test("ignores non-matching filters and non-zero offset pages", () => {
    const initial: UsageLogsResponse = {
      logs: [makeLog({ id: 1, provider_id: "openai" })],
      total: 1,
      limit: 50,
      offset: 0,
    };

    const otherProvider = makeLog({ id: 2, provider_id: "boz" });
    expect(prependUsageLogToResponse(initial, otherProvider, { provider: "openai" })).toBe(initial);

    const errLog = makeLog({
      id: 3,
      provider_id: "openai",
      status_code: 500,
      error_summary: "boom",
    });
    expect(prependUsageLogToResponse(initial, errLog, { status: "ok" })).toBe(initial);
    expect(prependUsageLogToResponse(initial, errLog, { status: "error" }).logs[0]?.id).toBe(3);

    const paged: UsageLogsResponse = { ...initial, offset: 50 };
    expect(prependUsageLogToResponse(paged, makeLog({ id: 4 }), { offset: 50 })).toBe(paged);
  });

  test("filters incoming SSE logs by exact numeric HTTP status code", () => {
    const initial: UsageLogsResponse = {
      logs: [makeLog({ id: 1, status_code: 429, error_summary: "rate limited" })],
      total: 1,
      limit: 50,
      offset: 0,
    };

    const okLog = makeLog({ id: 2, status_code: 200, error_summary: "" });
    const otherErrLog = makeLog({ id: 3, status_code: 500, error_summary: "boom" });
    const matching429 = makeLog({ id: 4, status_code: 429, error_summary: "too many requests" });

    expect(prependUsageLogToResponse(initial, okLog, { status: "429" })).toBe(initial);
    expect(prependUsageLogToResponse(initial, otherErrLog, { status: "429" })).toBe(initial);

    const next = prependUsageLogToResponse(initial, matching429, { status: "429" });
    expect(next.logs.map((l) => l.id)).toEqual([4, 1]);
    expect(next.total).toBe(2);
  });
});
