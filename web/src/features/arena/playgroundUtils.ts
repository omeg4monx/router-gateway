import type { ChatMessage, ChatSession } from "./schemas";

export interface SessionTelemetry {
  assistantTurns: number;
  totalTokens: number;
  avgLatencyMs: number;
  lastLatencyMs: number;
}

export function computeSessionTelemetry(messages: readonly ChatMessage[]): SessionTelemetry {
  let assistantTurns = 0;
  let totalTokens = 0;
  let totalLatencyMs = 0;
  let lastLatencyMs = 0;

  for (const message of messages) {
    if (message.role !== "assistant") {
      continue;
    }
    assistantTurns += 1;
    totalTokens += message.total_tokens;
    totalLatencyMs += message.latency_ms;
    lastLatencyMs = message.latency_ms;
  }

  return {
    assistantTurns,
    totalTokens,
    avgLatencyMs: assistantTurns > 0 ? Math.round(totalLatencyMs / assistantTurns) : 0,
    lastLatencyMs,
  };
}

export function filterChatSessions(sessions: readonly ChatSession[], query: string): ChatSession[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) {
    return [...sessions];
  }
  return sessions.filter(
    (session) =>
      session.title.toLowerCase().includes(normalized) ||
      session.model.toLowerCase().includes(normalized),
  );
}

export function formatLatencyMs(ms: number): string {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(2)}s`;
  }
  return `${Math.round(ms)}ms`;
}

export function formatTokenCount(tokens: number): string {
  if (tokens >= 1000) {
    return `${(Math.round(tokens / 100) / 10).toFixed(1)}k`;
  }
  return String(tokens);
}
