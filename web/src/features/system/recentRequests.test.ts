import { describe, expect, test } from "bun:test";
import type { UsageLog } from "@/features/usage/schemas";
import { filterRecentRequests, formatRelativeTimeShort, isUsageLogSuccess } from "./recentRequests";

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

describe("features/system/recentRequests", () => {
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
