import type { UsageLog } from "@/features/usage/schemas";

export const MAX_MONITOR_LINES = 250;

export type MonitorLevel = "ok" | "warn" | "error";
export type MonitorStatusFilter = "all" | "ok" | "error";

export interface MonitorEntry {
  id: string;
  timestamp: string;
  level: MonitorLevel;
  statusCode: number;
  provider: string;
  model: string;
  protocol: string;
  stream: boolean;
  latencyMs: number;
  ttftMs: number;
  inputTokens: number;
  outputTokens: number;
  cachedTokens: number;
  errorSummary: string;
}

export interface MonitorStats {
  total: number;
  okCount: number;
  errorCount: number;
  avgLatencyMs: number;
}

export function classifyStatusLevel(statusCode: number, errorSummary: string): MonitorLevel {
  if (statusCode >= 500 || (statusCode === 0 && errorSummary.trim() !== "")) {
    return "error";
  }
  if (statusCode >= 400 || errorSummary.trim() !== "") {
    return statusCode === 429 ? "warn" : "error";
  }
  return "ok";
}

export function formatTerminalClock(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return "--:--:--";
  }
  return date.toLocaleTimeString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

export function formatShortTokens(count: number): string {
  if (!Number.isFinite(count) || count <= 0) {
    return "0";
  }
  if (count >= 1000) {
    return `${(count / 1000).toFixed(1)}k`;
  }
  return String(Math.round(count));
}

export function formatShortLatency(ms: number): string {
  if (!Number.isFinite(ms) || ms <= 0) {
    return "0ms";
  }
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(2)}s`;
  }
  return `${Math.round(ms)}ms`;
}

export function usageLogToMonitorEntry(log: UsageLog): MonitorEntry {
  const errorSummary = (log.error_summary ?? "").trim();
  return {
    id: `usage-${log.id}`,
    timestamp: log.created_at,
    level: classifyStatusLevel(log.status_code, errorSummary),
    statusCode: log.status_code,
    provider: log.provider_id || "unknown",
    model: log.model || "unknown",
    protocol: log.protocol || "",
    stream: Boolean(log.stream),
    latencyMs: log.latency_ms ?? 0,
    ttftMs: log.ttft_ms ?? 0,
    inputTokens: log.input_tokens ?? 0,
    outputTokens: log.output_tokens ?? 0,
    cachedTokens: log.cached_tokens ?? 0,
    errorSummary,
  };
}

export function appendMonitorEntry(
  list: readonly MonitorEntry[],
  entry: MonitorEntry,
): MonitorEntry[] {
  if (list.some((item) => item.id === entry.id)) {
    return [...list];
  }
  const next = [...list, entry];
  if (next.length > MAX_MONITOR_LINES) {
    return next.slice(next.length - MAX_MONITOR_LINES);
  }
  return next;
}

export function filterMonitorEntries(
  entries: readonly MonitorEntry[],
  statusFilter: MonitorStatusFilter,
  query: string,
): MonitorEntry[] {
  const q = query.trim().toLowerCase();
  return entries.filter((entry) => {
    if (statusFilter === "ok" && entry.level !== "ok") {
      return false;
    }
    if (statusFilter === "error" && entry.level === "ok") {
      return false;
    }
    if (!q) {
      return true;
    }
    const haystack =
      `${entry.provider} ${entry.model} ${entry.protocol} ${entry.statusCode} ${entry.errorSummary}`.toLowerCase();
    return haystack.includes(q);
  });
}

export function computeMonitorStats(entries: readonly MonitorEntry[]): MonitorStats {
  if (entries.length === 0) {
    return { total: 0, okCount: 0, errorCount: 0, avgLatencyMs: 0 };
  }
  let okCount = 0;
  let errorCount = 0;
  let latencySum = 0;
  for (const entry of entries) {
    if (entry.level === "ok") {
      okCount += 1;
    } else {
      errorCount += 1;
    }
    latencySum += entry.latencyMs;
  }
  return {
    total: entries.length,
    okCount,
    errorCount,
    avgLatencyMs: Math.round(latencySum / entries.length),
  };
}
