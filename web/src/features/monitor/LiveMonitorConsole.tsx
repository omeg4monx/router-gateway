import {
  ArrowDown,
  Check,
  ChevronDown,
  ChevronRight,
  Copy,
  Download,
  Pause,
  Play,
  Radio,
  Search,
  Terminal,
  Trash2,
  X,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useClearConsoleLogsMutation } from "@/features/translator/hooks";
import type { ConsoleLogEntry } from "@/features/translator/schemas";
import { useUsageLogs } from "@/features/usage/hooks";
import type { UsageLog } from "@/features/usage/schemas";
import { useUIStore } from "@/shared/store/uiStore";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import {
  appendSystemLogBatch,
  computeLiveRps,
  computeMonitorStats,
  expireStaleInFlight,
  filterMonitorEntries,
  formatShortCost,
  formatShortLatency,
  formatShortProtocol,
  formatShortRequestId,
  formatShortTokens,
  formatSystemLogAttrs,
  formatTerminalClock,
  type KeyedSystemLog,
  MAX_MONITOR_LINES,
  type MonitorEntry,
  type MonitorLevel,
  type MonitorStatusFilter,
  mergeMonitorEntries,
  parseSystemLogAttrEntries,
  type SystemLogLevelFilter,
  serializeMonitorEntriesNDJSON,
  type TrafficStartPayload,
  trafficStartToMonitorEntry,
  upsertMonitorEntries,
  usageLogToMonitorEntry,
} from "./monitorFormat";

type MonitorTab = "traffic" | "system";
type StreamState = "connecting" | "live" | "paused" | "reconnecting";

const SYSTEM_LEVELS: readonly SystemLogLevelFilter[] = ["ALL", "INFO", "WARN", "ERROR", "DEBUG"];

function levelSymbol(level: MonitorLevel): string {
  switch (level) {
    case "live":
      return "●";
    case "ok":
      return "✓";
    case "warn":
      return "▲";
    case "error":
      return "✗";
  }
}

function levelColorClass(level: MonitorLevel): string {
  switch (level) {
    case "live":
      return "text-[var(--accent-primary)] animate-pulse";
    case "ok":
      return "text-[var(--status-ok)]";
    case "warn":
      return "text-[var(--status-warn)]";
    case "error":
      return "text-[var(--status-error)]";
  }
}

function systemLevelClass(level: string, message?: string): string {
  if (message?.includes("✗")) {
    return "text-[var(--status-error)] font-medium";
  }
  if (message?.includes("✓")) {
    return "text-[var(--status-ok)] font-medium";
  }
  switch (level.toUpperCase()) {
    case "ERROR":
      return "text-[var(--status-error)] font-medium";
    case "WARN":
      return "text-[var(--status-warn)] font-medium";
    case "DEBUG":
      return "text-[var(--text-muted)]";
    default:
      return "text-[var(--accent-text)]";
  }
}

function attrValueClass(key: string): string {
  const lower = key.toLowerCase();
  if (lower === "error" || lower === "err") {
    return "text-[var(--status-error)]";
  }
  if (lower === "provider" || lower === "model") {
    return "text-[var(--color-highlight)]";
  }
  if (lower === "status" || lower === "latency_ms" || lower === "attempt") {
    return "text-[var(--accent-text)]";
  }
  return "text-[var(--text-secondary)]";
}

function normalizeUsagePayload(raw: string): MonitorEntry | null {
  try {
    const parsed = JSON.parse(raw) as Partial<UsageLog>;
    if (!parsed || typeof parsed.provider_id !== "string") {
      return null;
    }
    const normalized: UsageLog = {
      id: typeof parsed.id === "number" ? parsed.id : 0,
      request_id: typeof parsed.request_id === "string" ? parsed.request_id : "",
      key_id: typeof parsed.key_id === "string" ? parsed.key_id : "",
      model: typeof parsed.model === "string" ? parsed.model : "",
      provider_id: parsed.provider_id,
      connection_id: typeof parsed.connection_id === "string" ? parsed.connection_id : "",
      protocol: typeof parsed.protocol === "string" ? parsed.protocol : "",
      input_tokens: typeof parsed.input_tokens === "number" ? parsed.input_tokens : 0,
      output_tokens: typeof parsed.output_tokens === "number" ? parsed.output_tokens : 0,
      cached_tokens: typeof parsed.cached_tokens === "number" ? parsed.cached_tokens : 0,
      cost_usd: typeof parsed.cost_usd === "number" ? parsed.cost_usd : 0,
      latency_ms: typeof parsed.latency_ms === "number" ? parsed.latency_ms : 0,
      ttft_ms: typeof parsed.ttft_ms === "number" ? parsed.ttft_ms : 0,
      status_code: typeof parsed.status_code === "number" ? parsed.status_code : 200,
      error_summary: typeof parsed.error_summary === "string" ? parsed.error_summary : "",
      stream: Boolean(parsed.stream),
      created_at:
        typeof parsed.created_at === "string" ? parsed.created_at : new Date().toISOString(),
    };
    return usageLogToMonitorEntry(normalized);
  } catch {
    return null;
  }
}

export function LiveMonitorConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";
  const clearSystemLogsMutation = useClearConsoleLogsMutation();

  const [tab, setTab] = useState<MonitorTab>("traffic");
  const [paused, setPaused] = useState(false);
  const [autoScroll, setAutoScroll] = useState(true);
  const [statusFilter, setStatusFilter] = useState<MonitorStatusFilter>("all");
  const [systemLevelFilter, setSystemLevelFilter] = useState<SystemLogLevelFilter>("ALL");
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [nowMs, setNowMs] = useState(() => Date.now());

  // Traffic stream state
  const [trafficEntries, setTrafficEntries] = useState<MonitorEntry[]>([]);
  const [trafficStreamState, setTrafficStreamState] = useState<StreamState>("connecting");
  const [queuedTrafficCount, setQueuedTrafficCount] = useState(0);
  const seededTrafficRef = useRef(false);
  const hadTrafficErrorRef = useRef(false);
  const pausedTrafficQueueRef = useRef<MonitorEntry[]>([]);
  const pendingTrafficBatchRef = useRef<MonitorEntry[]>([]);
  const trafficRafRef = useRef<number | null>(null);
  const pausedRef = useRef(false);

  // System console stream state
  const [systemEntries, setSystemEntries] = useState<KeyedSystemLog[]>([]);
  const [systemStreamState, setSystemStreamState] = useState<StreamState>("connecting");
  const [queuedSystemCount, setQueuedSystemCount] = useState(0);
  const pausedSystemQueueRef = useRef<ConsoleLogEntry[]>([]);
  const pendingSystemBatchRef = useRef<ConsoleLogEntry[]>([]);
  const systemRafRef = useRef<number | null>(null);
  const systemSeqRef = useRef(0);
  const copyTimerRef = useRef<number | null>(null);

  const viewportRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  // Terminal keyboard shortcuts: `/` focuses filter, `Space` toggles pause, `Esc` blurs/collapses
  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      const target = event.target as HTMLElement | null;
      const isTyping =
        target !== null &&
        (target.tagName === "INPUT" ||
          target.tagName === "TEXTAREA" ||
          target.tagName === "SELECT" ||
          target.isContentEditable);

      if (event.key === "Escape") {
        if (isTyping && target) {
          target.blur();
        } else if (expandedId !== null) {
          setExpandedId(null);
        }
        return;
      }
      if (isTyping || event.metaKey || event.ctrlKey || event.altKey) {
        return;
      }
      if (event.key === "/") {
        event.preventDefault();
        searchInputRef.current?.focus();
      } else if (event.key === " ") {
        event.preventDefault();
        setPaused((prev) => !prev);
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [expandedId]);

  const nextSystemSeq = useCallback(() => {
    systemSeqRef.current += 1;
    return systemSeqRef.current;
  }, []);

  const scheduleTrafficFlush = useCallback(() => {
    if (trafficRafRef.current !== null) {
      return;
    }
    const flush = () => {
      trafficRafRef.current = null;
      const batch = pendingTrafficBatchRef.current;
      if (batch.length === 0) {
        return;
      }
      pendingTrafficBatchRef.current = [];
      setTrafficEntries((prev) => upsertMonitorEntries(prev, batch));
      setNowMs(Date.now());
    };
    if (typeof requestAnimationFrame === "undefined") {
      flush();
      return;
    }
    trafficRafRef.current = requestAnimationFrame(flush);
  }, []);

  const scheduleSystemFlush = useCallback(() => {
    if (systemRafRef.current !== null) {
      return;
    }
    const flush = () => {
      systemRafRef.current = null;
      const batch = pendingSystemBatchRef.current;
      if (batch.length === 0) {
        return;
      }
      pendingSystemBatchRef.current = [];
      setSystemEntries((prev) => appendSystemLogBatch(prev, batch, nextSystemSeq));
    };
    if (typeof requestAnimationFrame === "undefined") {
      flush();
      return;
    }
    systemRafRef.current = requestAnimationFrame(flush);
  }, [nextSystemSeq]);

  useEffect(() => {
    return () => {
      if (trafficRafRef.current !== null && typeof cancelAnimationFrame !== "undefined") {
        cancelAnimationFrame(trafficRafRef.current);
      }
      if (systemRafRef.current !== null && typeof cancelAnimationFrame !== "undefined") {
        cancelAnimationFrame(systemRafRef.current);
      }
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
    };
  }, []);

  useEffect(() => {
    pausedRef.current = paused;
    if (!paused) {
      const queuedTraffic = pausedTrafficQueueRef.current;
      if (queuedTraffic.length > 0) {
        pausedTrafficQueueRef.current = [];
        setTrafficEntries((prev) => mergeMonitorEntries(prev, queuedTraffic));
      }
      setQueuedTrafficCount(0);

      const queuedSystem = pausedSystemQueueRef.current;
      if (queuedSystem.length > 0) {
        pausedSystemQueueRef.current = [];
        setSystemEntries((prev) => appendSystemLogBatch(prev, queuedSystem, nextSystemSeq));
      }
      setQueuedSystemCount(0);
    }
  }, [paused, nextSystemSeq]);

  // Initial recent provider/model activity snapshot (oldest -> newest like a terminal)
  const initialTrafficQuery = useUsageLogs({ limit: 50, offset: 0 });
  const refetchTrafficSnapshot = initialTrafficQuery.refetch;
  useEffect(() => {
    if (seededTrafficRef.current || !initialTrafficQuery.data) {
      return;
    }
    seededTrafficRef.current = true;
    const chronological = [...initialTrafficQuery.data.logs].reverse().map(usageLogToMonitorEntry);
    setTrafficEntries((prev) => mergeMonitorEntries(prev, chronological));
  }, [initialTrafficQuery.data]);

  // Live SSE: /api/usage/stream (0ms in-flight + 0ms finished + persisted usage)
  useEffect(() => {
    if (typeof EventSource === "undefined") {
      return;
    }
    setTrafficStreamState("connecting");
    const source = new EventSource("/api/usage/stream?live=1", { withCredentials: true });

    const enqueueTraffic = (entry: MonitorEntry) => {
      if (pausedRef.current) {
        pausedTrafficQueueRef.current = upsertMonitorEntries(pausedTrafficQueueRef.current, [
          entry,
        ]);
        setQueuedTrafficCount(pausedTrafficQueueRef.current.length);
        return;
      }
      pendingTrafficBatchRef.current.push(entry);
      scheduleTrafficFlush();
    };

    const onOpen = () => {
      setTrafficStreamState("live");
      if (hadTrafficErrorRef.current) {
        hadTrafficErrorRef.current = false;
        void refetchTrafficSnapshot().then((result) => {
          if (!result.data) {
            return;
          }
          const recovered = [...result.data.logs].reverse().map(usageLogToMonitorEntry);
          setTrafficEntries((prev) => mergeMonitorEntries(prev, recovered));
        });
      }
    };

    const onTrafficStarted = (evt: Event) => {
      setTrafficStreamState("live");
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const payload = JSON.parse((evt as MessageEvent).data) as TrafficStartPayload;
        if (!payload || typeof payload.provider_id !== "string") {
          return;
        }
        enqueueTraffic(trafficStartToMonitorEntry(payload));
      } catch {
        // Ignore malformed SSE frame.
      }
    };

    const onTrafficCompleted = (evt: Event) => {
      setTrafficStreamState("live");
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      const entry = normalizeUsagePayload((evt as MessageEvent).data);
      if (entry) {
        enqueueTraffic(entry);
      }
    };

    const onError = () => {
      hadTrafficErrorRef.current = true;
      setTrafficStreamState("reconnecting");
    };

    source.addEventListener("open", onOpen);
    source.addEventListener("traffic_started", onTrafficStarted);
    source.addEventListener("traffic_finished", onTrafficCompleted);
    source.addEventListener("usage", onTrafficCompleted);
    source.addEventListener("error", onError);

    return () => {
      source.removeEventListener("open", onOpen);
      source.removeEventListener("traffic_started", onTrafficStarted);
      source.removeEventListener("traffic_finished", onTrafficCompleted);
      source.removeEventListener("usage", onTrafficCompleted);
      source.removeEventListener("error", onError);
      source.close();
    };
  }, [refetchTrafficSnapshot, scheduleTrafficFlush]);

  // Live SSE: /api/translator/logs/stream (system console logs, only active on System tab)
  useEffect(() => {
    if (tab !== "system" || typeof EventSource === "undefined") {
      return;
    }
    setSystemStreamState("connecting");
    const source = new EventSource("/api/translator/logs/stream?tail=50", {
      withCredentials: true,
    });

    const onOpen = () => setSystemStreamState("live");
    const onLog = (evt: Event) => {
      setSystemStreamState("live");
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const entry = JSON.parse((evt as MessageEvent).data) as ConsoleLogEntry;
        if (!entry || typeof entry.message !== "string") {
          return;
        }
        if (pausedRef.current) {
          const nextQueue = [...pausedSystemQueueRef.current, entry];
          pausedSystemQueueRef.current =
            nextQueue.length > MAX_MONITOR_LINES
              ? nextQueue.slice(nextQueue.length - MAX_MONITOR_LINES)
              : nextQueue;
          setQueuedSystemCount(pausedSystemQueueRef.current.length);
          return;
        }
        pendingSystemBatchRef.current.push(entry);
        scheduleSystemFlush();
      } catch {
        // Ignore malformed frame.
      }
    };
    const onClear = () => {
      pausedSystemQueueRef.current = [];
      pendingSystemBatchRef.current = [];
      setQueuedSystemCount(0);
      setSystemEntries([]);
    };
    const onError = () => setSystemStreamState("reconnecting");

    source.addEventListener("open", onOpen);
    source.addEventListener("log", onLog);
    source.addEventListener("clear", onClear);
    source.addEventListener("error", onError);

    return () => {
      source.removeEventListener("open", onOpen);
      source.removeEventListener("log", onLog);
      source.removeEventListener("clear", onClear);
      source.removeEventListener("error", onError);
      source.close();
    };
  }, [tab, scheduleSystemFlush]);

  const visibleTraffic = useMemo(
    () => filterMonitorEntries(trafficEntries, statusFilter, searchQuery),
    [trafficEntries, statusFilter, searchQuery],
  );

  const visibleSystem = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    return systemEntries.filter((entry) => {
      if (systemLevelFilter !== "ALL" && entry.level.toUpperCase() !== systemLevelFilter) {
        return false;
      }
      if (!q) {
        return true;
      }
      return `${entry.level} ${entry.message} ${formatSystemLogAttrs(entry.attrs)}`
        .toLowerCase()
        .includes(q);
    });
  }, [systemEntries, systemLevelFilter, searchQuery]);

  const stats = useMemo(() => computeMonitorStats(trafficEntries), [trafficEntries]);
  const liveRps = useMemo(() => computeLiveRps(trafficEntries, nowMs), [trafficEntries, nowMs]);

  // Tick elapsed timer once per second only when at least one in-flight request is active
  // and auto-expire any stale in-flight row (>5m) so zombie LIVE entries never linger.
  useEffect(() => {
    if (stats.liveCount === 0) {
      return;
    }
    const timer = window.setInterval(() => {
      const now = Date.now();
      setNowMs(now);
      setTrafficEntries((prev) => expireStaleInFlight(prev, now));
    }, 1000);
    return () => window.clearInterval(timer);
  }, [stats.liveCount]);

  const scrollToBottom = useCallback(() => {
    const el = viewportRef.current;
    if (!el) return;
    el.scrollTo({ top: el.scrollHeight, behavior: "smooth" });
    setAutoScroll(true);
  }, []);

  // Auto-scroll smoothly when new lines arrive and auto-scroll is active
  const activeLength = tab === "traffic" ? visibleTraffic.length : visibleSystem.length;
  useEffect(() => {
    const el = viewportRef.current;
    if (!autoScroll || !el || activeLength === 0 || !tab) {
      return;
    }
    el.scrollTop = el.scrollHeight;
  }, [autoScroll, activeLength, tab]);

  const handleScroll = () => {
    const el = viewportRef.current;
    if (!el) return;
    const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    if (distanceFromBottom > 36 && autoScroll) {
      setAutoScroll(false);
    } else if (distanceFromBottom <= 16 && !autoScroll) {
      setAutoScroll(true);
    }
  };

  const handleClear = () => {
    if (tab === "traffic") {
      pausedTrafficQueueRef.current = [];
      pendingTrafficBatchRef.current = [];
      setTrafficEntries([]);
      setQueuedTrafficCount(0);
      setExpandedId(null);
    } else {
      pausedSystemQueueRef.current = [];
      pendingSystemBatchRef.current = [];
      setSystemEntries([]);
      setQueuedSystemCount(0);
      clearSystemLogsMutation.mutate();
    }
  };

  const handleExport = () => {
    const payload =
      tab === "traffic"
        ? serializeMonitorEntriesNDJSON(visibleTraffic)
        : visibleSystem
            .map(({ time, level, message, attrs }) =>
              JSON.stringify({ time, level, message, attrs }),
            )
            .join("\n");
    if (!payload) {
      return;
    }
    const blob = new Blob([`${payload}\n`], { type: "application/x-ndjson;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `omeg4mon-${tab}-${new Date().toISOString().replace(/[:.]/g, "-")}.ndjson`;
    link.click();
    URL.revokeObjectURL(url);
  };

  const handleCopyEntry = (entry: MonitorEntry) => {
    if (typeof navigator === "undefined" || !navigator.clipboard) {
      return;
    }
    void navigator.clipboard.writeText(JSON.stringify(entry, null, 2)).then(() => {
      setCopiedId(entry.id);
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
      copyTimerRef.current = window.setTimeout(() => setCopiedId(null), 1500);
    });
  };

  const queuedCount = tab === "traffic" ? queuedTrafficCount : queuedSystemCount;
  const rawStreamState = tab === "traffic" ? trafficStreamState : systemStreamState;
  const effectiveState: StreamState = paused ? "paused" : rawStreamState;

  const statusPill = useMemo(() => {
    switch (effectiveState) {
      case "live":
        return {
          label: isId ? "Aktif" : "Live",
          dotClass: "bg-[var(--status-ok)] animate-pulse",
          textClass: "text-[var(--status-ok)]",
        };
      case "paused":
        return {
          label: isId ? "Dijeda" : "Paused",
          dotClass: "bg-[var(--status-warn)]",
          textClass: "text-[var(--status-warn)]",
        };
      case "reconnecting":
        return {
          label: isId ? "Menghubungkan ulang" : "Reconnecting",
          dotClass: "bg-[var(--status-error)] animate-pulse",
          textClass: "text-[var(--status-error)]",
        };
      default:
        return {
          label: isId ? "Menghubungkan" : "Connecting",
          dotClass: "bg-[var(--text-muted)] animate-pulse",
          textClass: "text-[var(--text-muted)]",
        };
    }
  }, [effectiveState, isId]);

  return (
    <div className="space-y-3">
      <Card className="overflow-hidden">
        {/* Terminal top chrome */}
        <div className="flex flex-col gap-2.5 border-b border-[var(--border-subtle)] px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex flex-wrap items-center gap-2">
            <div className="flex items-center gap-2 pr-1">
              <Terminal className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <span className="text-sm font-semibold text-[var(--text-primary)]">
                {isId ? "Monitor Langsung" : "Live Monitor"}
              </span>
            </div>

            <div
              role="tablist"
              aria-label={isId ? "Sumber log" : "Log source"}
              className="inline-flex rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-0.5"
            >
              <button
                type="button"
                role="tab"
                aria-selected={tab === "traffic"}
                onClick={() => setTab("traffic")}
                className={`cursor-pointer rounded px-2.5 py-1 text-xs font-medium transition-colors ${
                  tab === "traffic"
                    ? "bg-[var(--bg-surface)] text-[var(--text-primary)] shadow-xs"
                    : "text-[var(--text-muted)] hover:text-[var(--text-primary)]"
                }`}
              >
                {isId ? "Koneksi Provider" : "Provider Traffic"}
              </button>
              <button
                type="button"
                role="tab"
                aria-selected={tab === "system"}
                onClick={() => setTab("system")}
                className={`cursor-pointer rounded px-2.5 py-1 text-xs font-medium transition-colors ${
                  tab === "system"
                    ? "bg-[var(--bg-surface)] text-[var(--text-primary)] shadow-xs"
                    : "text-[var(--text-muted)] hover:text-[var(--text-primary)]"
                }`}
              >
                {isId ? "Log Sistem" : "System Logs"}
              </button>
            </div>

            <span
              className={`inline-flex items-center gap-1.5 rounded-full border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-2.5 py-0.5 font-mono text-[11px] font-medium ${statusPill.textClass}`}
            >
              <span className={`h-1.5 w-1.5 rounded-full ${statusPill.dotClass}`} />
              <span>{statusPill.label}</span>
              {paused && queuedCount > 0 ? (
                <span className="tabular-nums text-[var(--text-muted)]">(+{queuedCount})</span>
              ) : null}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              size="xs"
              onClick={() => setPaused((prev) => !prev)}
              className="h-8 flex-1 justify-center sm:flex-none"
            >
              {paused ? (
                <>
                  <Play className="h-3.5 w-3.5 text-[var(--status-ok)]" aria-hidden="true" />
                  <span>{isId ? "Lanjut" : "Resume"}</span>
                </>
              ) : (
                <>
                  <Pause className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
                  <span>{isId ? "Jeda" : "Pause"}</span>
                </>
              )}
            </Button>
            <Button
              type="button"
              variant="outline"
              size="xs"
              onClick={handleExport}
              disabled={activeLength === 0}
              className="h-8 flex-1 justify-center sm:flex-none"
              title={isId ? "Ekspor NDJSON" : "Export NDJSON"}
            >
              <Download className="h-3.5 w-3.5" aria-hidden="true" />
              <span className="hidden sm:inline">NDJSON</span>
            </Button>
            <Button
              type="button"
              variant="outline"
              size="xs"
              onClick={handleClear}
              disabled={activeLength === 0}
              className="h-8 flex-1 justify-center sm:flex-none"
            >
              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{isId ? "Bersihkan" : "Clear"}</span>
            </Button>
          </div>
        </div>

        {/* Slim filter bar */}
        <div className="flex flex-col gap-2 border-b border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]/50 px-4 py-2 sm:flex-row sm:items-center sm:justify-between">
          <div className="relative w-full sm:w-64">
            <Search
              className="pointer-events-none absolute left-2.5 top-2 h-3.5 w-3.5 text-[var(--text-muted)]"
              aria-hidden="true"
            />
            <Input
              ref={searchInputRef}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={
                tab === "traffic"
                  ? isId
                    ? "Filter provider, model, req_id, error..."
                    : "Filter provider, model, req_id, error..."
                  : isId
                    ? "Filter log sistem & atribut..."
                    : "Filter system logs & attrs..."
              }
              className="h-8 w-full pl-8 pr-7 font-mono text-xs"
              aria-label={isId ? "Filter baris terminal" : "Filter terminal lines"}
            />
            {searchQuery ? (
              <button
                type="button"
                onClick={() => setSearchQuery("")}
                aria-label={isId ? "Hapus pencarian" : "Clear search"}
                className="absolute right-2 top-2 cursor-pointer text-[var(--text-muted)] hover:text-[var(--text-primary)]"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
            ) : null}
          </div>

          {tab === "traffic" ? (
            <div className="flex flex-wrap items-center gap-1.5">
              {(
                [
                  { id: "all", label: isId ? "Semua" : "All", count: stats.total },
                  ...(stats.liveCount > 0 || statusFilter === "live"
                    ? [{ id: "live" as const, label: "Live", count: stats.liveCount }]
                    : []),
                  { id: "ok", label: "OK", count: stats.okCount },
                  { id: "error", label: isId ? "Gagal" : "Errors", count: stats.errorCount },
                ] as const
              ).map((pill) => (
                <button
                  key={pill.id}
                  type="button"
                  onClick={() => setStatusFilter(pill.id)}
                  className={`inline-flex cursor-pointer items-center gap-1.5 rounded-md border px-2.5 py-1 font-mono text-[11px] font-medium transition-colors ${
                    statusFilter === pill.id
                      ? "border-[var(--accent-primary)] bg-[var(--accent-primary)]/10 text-[var(--accent-text)]"
                      : "border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-muted)] hover:text-[var(--text-primary)]"
                  }`}
                >
                  <span>{pill.label}</span>
                  <span className="tabular-nums opacity-75">{pill.count}</span>
                </button>
              ))}
            </div>
          ) : (
            <div className="flex flex-wrap items-center gap-1.5">
              {SYSTEM_LEVELS.map((lvl) => (
                <button
                  key={lvl}
                  type="button"
                  onClick={() => setSystemLevelFilter(lvl)}
                  className={`inline-flex cursor-pointer items-center rounded-md border px-2 py-1 font-mono text-[11px] font-medium transition-colors ${
                    systemLevelFilter === lvl
                      ? "border-[var(--accent-primary)] bg-[var(--accent-primary)]/10 text-[var(--accent-text)]"
                      : "border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-muted)] hover:text-[var(--text-primary)]"
                  }`}
                >
                  {lvl}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Authentic dark terminal window pane (Omegamon Zwart tokens even in Light theme) */}
        <div
          data-theme="dark"
          className="relative bg-[var(--bg-canvas)] text-[var(--text-primary)] selection:bg-[var(--accent-primary)]/25"
        >
          {/* Terminal titlebar chrome */}
          <div className="flex items-center justify-between gap-2 border-b border-[var(--border-strong)]/60 bg-[var(--bg-surface)] px-3 py-1.5 font-mono text-[11px]">
            <div className="flex min-w-0 items-center gap-2.5">
              <div className="flex shrink-0 items-center gap-1.5" aria-hidden="true">
                <span className="h-2.5 w-2.5 rounded-full bg-[var(--status-error)]/80" />
                <span className="h-2.5 w-2.5 rounded-full bg-[var(--status-warn)]/80" />
                <span className="h-2.5 w-2.5 rounded-full bg-[var(--status-ok)]/80" />
              </div>
              <span className="truncate text-[var(--text-muted)]">
                <span className="text-[var(--status-ok)]">omeg4mon@router</span>
                <span>:</span>
                <span className="text-[var(--color-highlight)]">~</span>
                <span className="text-[var(--text-secondary)]">
                  {" "}
                  $ tail -f{" "}
                  {tab === "traffic" ? "/api/usage/stream?live=1" : "/api/translator/logs/stream"}
                </span>
              </span>
            </div>
            <div className="hidden shrink-0 items-center gap-3 text-[10px] text-[var(--text-muted)] md:flex">
              <span>[/] filter</span>
              <span>[Space] {paused ? "resume" : "pause"}</span>
              <span>[Esc] close</span>
            </div>
          </div>

          <div
            ref={viewportRef}
            onScroll={handleScroll}
            className="h-[calc(100dvh-20.5rem)] min-h-[340px] overflow-y-auto bg-[var(--bg-canvas)] px-2.5 py-2.5 font-mono text-xs leading-5 sm:px-3.5 lg:h-[calc(100dvh-17.5rem)]"
          >
            {tab === "traffic" ? (
              visibleTraffic.length === 0 ? (
                <div className="flex h-full flex-col items-center justify-center gap-2 text-center text-[var(--text-muted)]">
                  <Radio className="h-5 w-5 opacity-50" aria-hidden="true" />
                  <p className="text-xs">
                    {initialTrafficQuery.isPending
                      ? isId
                        ? "Memuat aktivitas koneksi provider..."
                        : "Loading provider connection activity..."
                      : isId
                        ? "Menunggu aktivitas koneksi provider/model..."
                        : "Waiting for provider/model connection activity..."}
                  </p>
                </div>
              ) : (
                <div className="space-y-0.5">
                  {visibleTraffic.map((entry, idx) => {
                    const isExpanded = expandedId === entry.id;
                    const liveElapsedMs = entry.inFlight
                      ? Math.max(0, nowMs - Date.parse(entry.timestamp))
                      : entry.latencyMs;
                    const shortProto = formatShortProtocol(entry.protocol);
                    const shortReqId = formatShortRequestId(entry.requestId);
                    const shortCost = formatShortCost(entry.costUsd);
                    return (
                      <div key={entry.id} className="rounded">
                        <button
                          type="button"
                          onClick={() => setExpandedId(isExpanded ? null : entry.id)}
                          aria-expanded={isExpanded}
                          className={`w-full cursor-pointer rounded px-2 py-1 text-left transition-colors hover:bg-[var(--bg-surface-raised)]/80 sm:flex sm:items-baseline sm:gap-2.5 sm:py-0.5 ${
                            isExpanded ? "bg-[var(--bg-surface-raised)]" : ""
                          }`}
                        >
                          {/* Mobile line 1 / Desktop inline prefix */}
                          <div className="flex items-baseline gap-2 sm:shrink-0">
                            <span className="hidden w-8 select-none border-r border-[var(--border-strong)]/40 pr-2 text-right text-[11px] tabular-nums text-[var(--text-muted)]/60 sm:inline-block">
                              {idx + 1}
                            </span>
                            {isExpanded ? (
                              <ChevronDown
                                className="h-3 w-3 shrink-0 self-center text-[var(--text-muted)]"
                                aria-hidden="true"
                              />
                            ) : (
                              <ChevronRight
                                className="h-3 w-3 shrink-0 self-center text-[var(--text-muted)]/60"
                                aria-hidden="true"
                              />
                            )}
                            <span className="text-[11px] tabular-nums text-[var(--text-muted)]">
                              [{formatTerminalClock(entry.timestamp)}]
                            </span>
                            <span
                              className={`inline-flex w-12 items-center gap-1 font-semibold tabular-nums ${levelColorClass(entry.level)}`}
                            >
                              <span>{levelSymbol(entry.level)}</span>
                              <span>{entry.inFlight ? "LIVE" : entry.statusCode}</span>
                            </span>
                            {shortProto ? (
                              <span className="rounded bg-[var(--bg-surface-raised)] px-1 text-[10px] uppercase text-[var(--color-highlight)]">
                                {shortProto}
                              </span>
                            ) : null}
                            {entry.stream ? (
                              <Badge variant="default" className="px-1 py-0 text-[10px] leading-4">
                                SSE
                              </Badge>
                            ) : null}
                          </div>

                          {/* Provider / Model / ReqID */}
                          <div className="mt-0.5 flex min-w-0 flex-1 flex-wrap items-baseline gap-x-1.5 gap-y-0.5 sm:mt-0">
                            <span className="truncate font-semibold text-[var(--accent-text)]">
                              {entry.provider}
                            </span>
                            <span className="text-[var(--text-muted)]">/</span>
                            <span className="truncate text-[var(--text-primary)]">
                              {entry.model}
                            </span>
                            {shortReqId ? (
                              <span className="hidden text-[10px] text-[var(--text-muted)]/80 xl:inline">
                                {shortReqId}
                              </span>
                            ) : null}

                            {entry.errorSummary ? (
                              <span className="truncate text-[11px] text-[var(--status-error)]">
                                — {entry.errorSummary}
                              </span>
                            ) : null}
                          </div>

                          {/* Compact telemetry */}
                          <div className="mt-0.5 flex shrink-0 items-center gap-2.5 text-[11px] tabular-nums text-[var(--text-muted)] sm:mt-0">
                            <span className="text-[var(--text-secondary)]">
                              {formatShortLatency(liveElapsedMs)}
                            </span>
                            {entry.ttftMs > 0 ? (
                              <span className="hidden md:inline">
                                ttft={formatShortLatency(entry.ttftMs)}
                              </span>
                            ) : null}
                            {!entry.inFlight ? (
                              <span>
                                {formatShortTokens(entry.inputTokens)}→
                                {formatShortTokens(entry.outputTokens)}
                                {entry.cachedTokens > 0 ? (
                                  <span className="text-[var(--color-highlight)]">
                                    +{formatShortTokens(entry.cachedTokens)}c
                                  </span>
                                ) : null}
                              </span>
                            ) : (
                              <span className="text-[var(--color-highlight)]">
                                {isId ? "streaming..." : "in-flight..."}
                              </span>
                            )}
                            {shortCost ? (
                              <span className="hidden text-[var(--status-warn)] lg:inline">
                                {shortCost}
                              </span>
                            ) : null}
                          </div>
                        </button>

                        {isExpanded ? (
                          <div className="mx-2 mb-1.5 mt-0.5 rounded border border-[var(--border-strong)]/70 bg-[var(--bg-surface)] px-3 py-2 text-[11px] text-[var(--text-muted)]">
                            <div className="flex flex-wrap items-center justify-between gap-2">
                              <div className="flex flex-wrap gap-x-4 gap-y-1">
                                {entry.requestId ? (
                                  <span>
                                    req:{" "}
                                    <span className="text-[var(--text-primary)]">
                                      {entry.requestId}
                                    </span>
                                  </span>
                                ) : null}
                                {entry.connectionId ? (
                                  <span>
                                    conn:{" "}
                                    <span className="text-[var(--text-primary)]">
                                      {entry.connectionId}
                                    </span>
                                  </span>
                                ) : null}
                                {entry.protocol ? (
                                  <span>
                                    proto:{" "}
                                    <span className="text-[var(--color-highlight)]">
                                      {entry.protocol}
                                    </span>
                                  </span>
                                ) : null}
                                {entry.cachedTokens > 0 ? (
                                  <span>
                                    cached:{" "}
                                    <span className="text-[var(--color-highlight)]">
                                      {formatShortTokens(entry.cachedTokens)}
                                    </span>
                                  </span>
                                ) : null}
                                {entry.costUsd > 0 ? (
                                  <span>
                                    cost:{" "}
                                    <span className="text-[var(--status-warn)]">
                                      ${entry.costUsd.toFixed(5)}
                                    </span>
                                  </span>
                                ) : null}
                              </div>
                              <button
                                type="button"
                                onClick={() => handleCopyEntry(entry)}
                                className="inline-flex cursor-pointer items-center gap-1 rounded border border-[var(--border-strong)] px-2 py-0.5 text-[10px] text-[var(--text-primary)] hover:border-[var(--accent-primary)]"
                              >
                                {copiedId === entry.id ? (
                                  <>
                                    <Check
                                      className="h-3 w-3 text-[var(--status-ok)]"
                                      aria-hidden="true"
                                    />
                                    <span>{isId ? "Tersalin" : "Copied"}</span>
                                  </>
                                ) : (
                                  <>
                                    <Copy className="h-3 w-3" aria-hidden="true" />
                                    <span>JSON</span>
                                  </>
                                )}
                              </button>
                            </div>
                            {entry.errorSummary ? (
                              <p className="mt-1.5 break-words text-[var(--status-error)]">
                                {entry.errorSummary}
                              </p>
                            ) : null}
                          </div>
                        ) : null}
                      </div>
                    );
                  })}
                </div>
              )
            ) : visibleSystem.length === 0 ? (
              <div className="flex h-full flex-col items-center justify-center gap-2 text-center text-[var(--text-muted)]">
                <Terminal className="h-5 w-5 opacity-50" aria-hidden="true" />
                <p className="text-xs">
                  {isId ? "Belum ada log sistem." : "No system console logs yet."}
                </p>
              </div>
            ) : (
              <div className="space-y-0.5">
                {visibleSystem.map((line, idx) => {
                  const attrEntries = parseSystemLogAttrEntries(line.attrs);
                  return (
                    <div
                      key={line.seq}
                      className="flex gap-2 rounded px-2 py-0.5 transition-colors hover:bg-[var(--bg-surface-raised)]/80"
                    >
                      <span className="hidden w-8 shrink-0 select-none border-r border-[var(--border-strong)]/40 pr-2 text-right text-[11px] tabular-nums text-[var(--text-muted)]/60 sm:inline-block">
                        {idx + 1}
                      </span>
                      <span className="shrink-0 text-[11px] tabular-nums text-[var(--text-muted)]">
                        [{formatTerminalClock(line.time)}]
                      </span>
                      <span
                        className={`shrink-0 uppercase ${systemLevelClass(line.level, line.message)}`}
                      >
                        [{line.level}]
                      </span>
                      <span className="min-w-0 break-words text-[var(--text-primary)]">
                        <span>{line.message}</span>
                        {attrEntries.length > 0 ? (
                          <span className="ml-2 inline-flex flex-wrap gap-x-2 text-[11px]">
                            {attrEntries.map(({ key, value }) => (
                              <span key={key}>
                                <span className="text-[var(--text-muted)]">{key}=</span>
                                <span className={attrValueClass(key)}>{value}</span>
                              </span>
                            ))}
                          </span>
                        ) : null}
                      </span>
                    </div>
                  );
                })}
              </div>
            )}

            {/* Authentic live terminal prompt + blinking block cursor */}
            <div className="mt-2 flex items-center gap-1.5 px-2 pt-1 text-[11px] text-[var(--text-muted)]">
              <span className="text-[var(--status-ok)]">omeg4mon@router</span>
              <span>:</span>
              <span className="text-[var(--color-highlight)]">~</span>
              <span>$</span>
              {paused ? (
                <span className="text-[var(--status-warn)]">
                  {isId
                    ? "[DIJEDA — tekan Space untuk lanjut]"
                    : "[PAUSED — press Space to resume]"}
                </span>
              ) : (
                <span
                  aria-hidden="true"
                  className="inline-block h-3.5 w-2 animate-pulse bg-[var(--status-ok)]/90 align-middle"
                />
              )}
            </div>
          </div>

          {/* Floating jump-to-latest pill when scrolled up */}
          {!autoScroll && activeLength > 0 ? (
            <button
              type="button"
              onClick={scrollToBottom}
              className="absolute bottom-3 right-4 inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-[var(--border-strong)] bg-[var(--bg-surface)] px-3 py-1 font-mono text-[11px] font-medium text-[var(--text-primary)] shadow-md transition-colors hover:border-[var(--accent-primary)] hover:text-[var(--accent-text)]"
            >
              <ArrowDown className="h-3 w-3" aria-hidden="true" />
              <span>{isId ? "Ke terbaru" : "Latest"}</span>
            </button>
          ) : null}
        </div>

        {/* Tmux-style bottom terminal status bar */}
        <div
          data-theme="dark"
          className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-strong)]/60 bg-[var(--bg-surface)] px-3 py-1.5 font-mono text-[11px] tabular-nums text-[var(--text-muted)]"
        >
          <div className="flex items-center gap-2.5">
            <span
              className={`rounded px-1.5 py-0.5 text-[10px] font-semibold uppercase ${
                effectiveState === "live"
                  ? "bg-[var(--status-ok)]/15 text-[var(--status-ok)]"
                  : effectiveState === "paused"
                    ? "bg-[var(--status-warn)]/15 text-[var(--status-warn)]"
                    : "bg-[var(--status-error)]/15 text-[var(--status-error)]"
              }`}
            >
              {statusPill.label}
            </span>
            <span className="text-[var(--text-secondary)]">
              {tab === "traffic" ? "/api/usage/stream?live=1" : "/api/translator/logs/stream"}
            </span>
            <span className="hidden sm:inline">·</span>
            <span className="hidden sm:inline">
              {activeLength}/{MAX_MONITOR_LINES} {isId ? "baris" : "lines"}
            </span>
          </div>

          {tab === "traffic" ? (
            <div className="flex items-center gap-3">
              {stats.liveCount > 0 ? (
                <span className="text-[var(--color-highlight)]">LIVE {stats.liveCount}</span>
              ) : null}
              <span className="text-[var(--text-secondary)]">{liveRps.toFixed(1)} req/s</span>
              <span className="text-[var(--status-ok)]">OK {stats.okCount}</span>
              <span className={stats.errorCount > 0 ? "text-[var(--status-error)]" : ""}>
                ERR {stats.errorCount}
              </span>
              <span>avg {formatShortLatency(stats.avgLatencyMs)}</span>
            </div>
          ) : (
            <span>{isId ? "Ring buffer SSE 60fps" : "Live ring buffer · SSE 60fps"}</span>
          )}
        </div>
      </Card>
    </div>
  );
}
