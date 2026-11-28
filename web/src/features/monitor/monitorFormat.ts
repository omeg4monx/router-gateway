import type { ConsoleLogEntry } from "@/features/translator/schemas";
import type { UsageLog } from "@/features/usage/schemas";

export const MAX_MONITOR_LINES = 250;

export type KeyedSystemLog = ConsoleLogEntry & { seq: number; fingerprint: string };
export type MonitorLevel = "live" | "ok" | "warn" | "error";
export type MonitorStatusFilter = "all" | "live" | "ok" | "error";
export type SystemLogLevelFilter = "ALL" | "DEBUG" | "INFO" | "WARN" | "ERROR";

export interface TrafficStartPayload {
  request_id?: string;
  provider_id?: string;
  connection_id?: string;
  model?: string;
  protocol?: string;
  stream?: boolean;
  started_at?: string;
}

export interface MonitorEntry {
  id: string;
  requestId: string;
  connectionId: string;
  timestamp: string;
  level: MonitorLevel;
  inFlight: boolean;
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
  costUsd: number;
  errorSummary: string;
}

export interface MonitorStats {
  total: number;
  liveCount: number;
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

export function formatShortProtocol(protocol: string): string {
  const normalized = protocol.trim().toLowerCase();
  if (!normalized) {
    return "";
  }
  switch (normalized) {
    case "openai-chat":
      return "oai";
    case "openai-responses":
      return "resp";
    case "anthropic":
      return "anth";
    case "gemini":
      return "gem";
    default:
      return normalized.slice(0, 6);
  }
}

export function formatShortRequestId(requestId: string): string {
  const trimmed = requestId.trim();
  if (!trimmed) {
    return "";
  }
  const suffix = trimmed.length > 6 ? trimmed.slice(-6) : trimmed;
  return `#${suffix}`;
}

export function formatShortCost(costUsd: number): string {
  if (!Number.isFinite(costUsd) || costUsd <= 0) {
    return "";
  }
  if (costUsd < 0.0001) {
    return "<$0.0001";
  }
  return `$${costUsd.toFixed(4)}`;
}

export function trafficStartToMonitorEntry(payload: TrafficStartPayload): MonitorEntry {
  const requestId = (payload.request_id ?? "").trim();
  const timestamp = payload.started_at || new Date().toISOString();
  return {
    id: requestId ? `req-${requestId}` : `live-${Date.now()}`,
    requestId,
    connectionId: payload.connection_id || "",
    timestamp,
    level: "live",
    inFlight: true,
    statusCode: 0,
    provider: payload.provider_id || "unknown",
    model: payload.model || "unknown",
    protocol: payload.protocol || "",
    stream: Boolean(payload.stream),
    latencyMs: 0,
    ttftMs: 0,
    inputTokens: 0,
    outputTokens: 0,
    cachedTokens: 0,
    costUsd: 0,
    errorSummary: "",
  };
}

export function usageLogToMonitorEntry(log: UsageLog): MonitorEntry {
  const errorSummary = (log.error_summary ?? "").trim();
  const requestId = (log.request_id ?? "").trim();
  const id = log.id > 0 ? `usage-${log.id}` : requestId ? `req-${requestId}` : `usage-${log.id}`;
  return {
    id,
    requestId,
    connectionId: log.connection_id || "",
    timestamp: log.created_at,
    level: classifyStatusLevel(log.status_code, errorSummary),
    inFlight: false,
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
    costUsd: log.cost_usd ?? 0,
    errorSummary,
  };
}

function findMatchingIndex(list: readonly MonitorEntry[], entry: MonitorEntry): number {
  for (let i = list.length - 1; i >= 0; i -= 1) {
    const item = list[i];
    if (!item) continue;
    if (item.id === entry.id) {
      return i;
    }
    if (entry.requestId !== "" && item.requestId === entry.requestId) {
      // Collapse only when upgrading/merging an unpersisted live event (in-flight or traffic_finished).
      const itemIsLive = item.inFlight || item.id.startsWith("req-");
      const entryIsLive = entry.inFlight || entry.id.startsWith("req-");
      if (itemIsLive || entryIsLive) {
        return i;
      }
    }
  }
  return -1;
}

export function upsertMonitorEntries(
  list: readonly MonitorEntry[],
  incoming: readonly MonitorEntry[],
): MonitorEntry[] {
  if (incoming.length === 0) {
    return [...list];
  }
  const next = [...list];
  for (const entry of incoming) {
    const idx = findMatchingIndex(next, entry);
    if (idx >= 0) {
      const existing = next[idx];
      // Never downgrade a finished entry back to in-flight if events arrive out of order.
      if (existing && !existing.inFlight && entry.inFlight) {
        continue;
      }
      next[idx] = entry;
    } else {
      next.push(entry);
    }
  }
  if (next.length > MAX_MONITOR_LINES) {
    return next.slice(next.length - MAX_MONITOR_LINES);
  }
  return next;
}

export function appendMonitorEntry(
  list: readonly MonitorEntry[],
  entry: MonitorEntry,
): MonitorEntry[] {
  return upsertMonitorEntries(list, [entry]);
}

export function filterMonitorEntries(
  entries: readonly MonitorEntry[],
  statusFilter: MonitorStatusFilter,
  query: string,
): MonitorEntry[] {
  const q = query.trim().toLowerCase();
  return entries.filter((entry) => {
    if (statusFilter === "live" && !entry.inFlight) {
      return false;
    }
    if (statusFilter === "ok" && entry.level !== "ok") {
      return false;
    }
    if (statusFilter === "error" && (entry.level === "ok" || entry.inFlight)) {
      return false;
    }
    if (!q) {
      return true;
    }
    const haystack =
      `${entry.provider} ${entry.model} ${entry.protocol} ${entry.requestId} ${entry.connectionId} ${entry.statusCode} ${entry.errorSummary}`.toLowerCase();
    return haystack.includes(q);
  });
}

export function computeMonitorStats(entries: readonly MonitorEntry[]): MonitorStats {
  if (entries.length === 0) {
    return { total: 0, liveCount: 0, okCount: 0, errorCount: 0, avgLatencyMs: 0 };
  }
  let liveCount = 0;
  let okCount = 0;
  let errorCount = 0;
  let completedCount = 0;
  let latencySum = 0;
  for (const entry of entries) {
    if (entry.inFlight) {
      liveCount += 1;
      continue;
    }
    completedCount += 1;
    if (entry.level === "ok") {
      okCount += 1;
    } else {
      errorCount += 1;
    }
    latencySum += entry.latencyMs;
  }
  return {
    total: entries.length,
    liveCount,
    okCount,
    errorCount,
    avgLatencyMs: completedCount > 0 ? Math.round(latencySum / completedCount) : 0,
  };
}

export function computeLiveRps(
  entries: readonly MonitorEntry[],
  nowMs: number = Date.now(),
  windowMs: number = 10_000,
): number {
  if (entries.length === 0 || windowMs <= 0) {
    return 0;
  }
  const cutoff = nowMs - windowMs;
  let count = 0;
  for (let i = entries.length - 1; i >= 0; i -= 1) {
    const ts = Date.parse(entries[i]?.timestamp ?? "");
    if (Number.isFinite(ts) && ts >= cutoff && ts <= nowMs + 1000) {
      count += 1;
    }
  }
  return count / (windowMs / 1000);
}

export function serializeMonitorEntriesNDJSON(entries: readonly MonitorEntry[]): string {
  return entries.map((entry) => JSON.stringify(entry)).join("\n");
}

export function expireStaleInFlight(
  entries: readonly MonitorEntry[],
  nowMs: number = Date.now(),
  maxAgeMs: number = 300_000,
): MonitorEntry[] {
  let changed = false;
  const next = entries.map((entry) => {
    if (!entry.inFlight) {
      return entry;
    }
    const startedMs = Date.parse(entry.timestamp);
    if (!Number.isFinite(startedMs) || nowMs - startedMs < maxAgeMs) {
      return entry;
    }
    changed = true;
    return {
      ...entry,
      inFlight: false,
      level: "warn" as const,
      statusCode: 499,
      latencyMs: Math.max(0, nowMs - startedMs),
      errorSummary: entry.errorSummary || "stream ended / timed out",
    };
  });
  return changed ? next : [...entries];
}

function stableSerializeAttrs(attrs?: Record<string, unknown>): string {
  if (!attrs) {
    return "";
  }
  const keys = Object.keys(attrs).sort();
  if (keys.length === 0) {
    return "";
  }
  return keys.map((key) => `${key}:${JSON.stringify(attrs[key])}`).join(",");
}

export function systemLogFingerprint(entry: ConsoleLogEntry): string {
  return `${entry.time}|${entry.level}|${entry.message}|${stableSerializeAttrs(entry.attrs)}`;
}

export function appendSystemLogBatch(
  list: readonly KeyedSystemLog[],
  incoming: readonly ConsoleLogEntry[],
  nextSeq: () => number,
): KeyedSystemLog[] {
  if (incoming.length === 0) {
    return [...list];
  }
  const seen = new Set<string>();
  for (const item of list) {
    seen.add(item.fingerprint);
  }
  const next = [...list];
  for (const entry of incoming) {
    const fingerprint = systemLogFingerprint(entry);
    if (seen.has(fingerprint)) {
      continue;
    }
    seen.add(fingerprint);
    next.push({ ...entry, seq: nextSeq(), fingerprint });
  }
  if (next.length > MAX_MONITOR_LINES) {
    return next.slice(next.length - MAX_MONITOR_LINES);
  }
  return next;
}

export function appendSystemLogEntry(
  list: readonly KeyedSystemLog[],
  entry: ConsoleLogEntry,
  nextSeq: () => number,
): KeyedSystemLog[] {
  return appendSystemLogBatch(list, [entry], nextSeq);
}

export function parseSystemLogAttrEntries(
  attrs?: Record<string, unknown>,
): Array<{ key: string; value: string }> {
  if (!attrs) {
    return [];
  }
  const entries: Array<{ key: string; value: string }> = [];
  for (const [key, value] of Object.entries(attrs)) {
    if (value === undefined || value === null || value === "") {
      continue;
    }
    const rendered =
      typeof value === "string" || typeof value === "number" || typeof value === "boolean"
        ? String(value)
        : JSON.stringify(value);
    entries.push({ key, value: rendered });
  }
  return entries;
}

export function formatSystemLogAttrs(attrs?: Record<string, unknown>): string {
  return parseSystemLogAttrEntries(attrs)
    .map(({ key, value }) => `${key}=${value}`)
    .join(" ");
}

export function mergeMonitorEntries(
  existing: readonly MonitorEntry[],
  incoming: readonly MonitorEntry[],
): MonitorEntry[] {
  const combined = upsertMonitorEntries(existing, incoming).sort((a, b) => {
    const timeDiff = Date.parse(a.timestamp) - Date.parse(b.timestamp);
    if (Number.isFinite(timeDiff) && timeDiff !== 0) {
      return timeDiff;
    }
    return a.id.localeCompare(b.id);
  });
  if (combined.length > MAX_MONITOR_LINES) {
    return combined.slice(combined.length - MAX_MONITOR_LINES);
  }
  return combined;
}
