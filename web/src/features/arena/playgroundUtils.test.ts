import { describe, expect, test } from "bun:test";
import {
  computeSessionTelemetry,
  filterChatSessions,
  formatLatencyMs,
  formatTokenCount,
} from "./playgroundUtils";
import type { ChatMessage, ChatSession } from "./schemas";

const sampleSessions: ChatSession[] = [
  {
    id: "s1",
    title: "Gateway Failover Debug",
    model: "openai/gpt-4o",
    system_prompt: "",
    max_tokens: 0,
    message_count: 4,
    created_at: "2026-11-26T10:00:00Z",
    updated_at: "2026-11-26T10:05:00Z",
  },
  {
    id: "s2",
    title: "Prompt Guard Tuning",
    model: "anthropic/claude-sonnet",
    system_prompt: "Be strict.",
    max_tokens: 1024,
    message_count: 2,
    created_at: "2026-11-26T11:00:00Z",
    updated_at: "2026-11-26T11:02:00Z",
  },
];

const sampleMessages: ChatMessage[] = [
  {
    id: "m1",
    session_id: "s1",
    role: "user",
    content: "Hello",
    model: "openai/gpt-4o",
    latency_ms: 0,
    prompt_tokens: 0,
    completion_tokens: 0,
    total_tokens: 0,
    created_at: "2026-11-26T10:00:00Z",
  },
  {
    id: "m2",
    session_id: "s1",
    role: "assistant",
    content: "Hi there!",
    model: "openai/gpt-4o",
    latency_ms: 420,
    prompt_tokens: 12,
    completion_tokens: 28,
    total_tokens: 40,
    created_at: "2026-11-26T10:00:01Z",
  },
  {
    id: "m3",
    session_id: "s1",
    role: "assistant",
    content: "Follow up.",
    model: "openai/gpt-4o",
    latency_ms: 1580,
    prompt_tokens: 50,
    completion_tokens: 110,
    total_tokens: 160,
    created_at: "2026-11-26T10:01:01Z",
  },
];

describe("playgroundUtils", () => {
  test("computes aggregate assistant telemetry accurately", () => {
    const stats = computeSessionTelemetry(sampleMessages);
    expect(stats.assistantTurns).toBe(2);
    expect(stats.totalTokens).toBe(200);
    expect(stats.avgLatencyMs).toBe(1000);
    expect(stats.lastLatencyMs).toBe(1580);

    const empty = computeSessionTelemetry([]);
    expect(empty.assistantTurns).toBe(0);
    expect(empty.totalTokens).toBe(0);
    expect(empty.avgLatencyMs).toBe(0);
    expect(empty.lastLatencyMs).toBe(0);
  });

  test("filters sessions by title and model case-insensitively", () => {
    expect(filterChatSessions(sampleSessions, "")).toHaveLength(2);
    expect(filterChatSessions(sampleSessions, "claude").map((s) => s.id)).toEqual(["s2"]);
    expect(filterChatSessions(sampleSessions, "FAILOVER").map((s) => s.id)).toEqual(["s1"]);
    expect(filterChatSessions(sampleSessions, "nonexistent")).toHaveLength(0);
  });

  test("formats latency and token counts concisely", () => {
    expect(formatLatencyMs(420)).toBe("420ms");
    expect(formatLatencyMs(1580)).toBe("1.58s");
    expect(formatTokenCount(240)).toBe("240");
    expect(formatTokenCount(1450)).toBe("1.5k");
  });
});
