import { ArrowDown, Pause, Play, Radio, Search, Terminal, Trash2, X } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useConsoleLogsSnapshot } from "@/features/translator/hooks";
import type { ConsoleLogEntry } from "@/features/translator/schemas";
import { useUsageLogs } from "@/features/usage/hooks";
import type { UsageLog } from "@/features/usage/schemas";
import { useUIStore } from "@/shared/store/uiStore";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import {
  appendMonitorEntry,
  computeMonitorStats,
  filterMonitorEntries,
  formatShortLatency,
  formatShortTokens,
  formatTerminalClock,
  MAX_MONITOR_LINES,
  type MonitorEntry,
  type MonitorLevel,
  type MonitorStatusFilter,
  usageLogToMonitorEntry,
} from "./monitorFormat";

type MonitorTab = "traffic" | "system";
type StreamState = "connecting" | "live" | "paused" | "reconnecting";
type KeyedSystemLog = ConsoleLogEntry & { seq: number };

function levelSymbol(level: MonitorLevel): string {
  switch (level) {
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
    case "ok":
      return "text-[var(--status-ok)]";
    case "warn":
      return "text-[var(--status-warn)]";
    case "error":
      return "text-[var(--status-error)]";
  }
}

function systemLevelClass(level: string): string {
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

export function LiveMonitorConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const [tab, setTab] = useState<MonitorTab>("traffic");
  const [paused, setPaused] = useState(false);
  const [autoScroll, setAutoScroll] = useState(true);
  const [statusFilter, setStatusFilter] = useState<MonitorStatusFilter>("all");
  const [searchQuery, setSearchQuery] = useState("");

  // Traffic stream state
  const [trafficEntries, setTrafficEntries] = useState<MonitorEntry[]>([]);
  const [trafficStreamState, setTrafficStreamState] = useState<StreamState>("connecting");
  const [queuedCount, setQueuedCount] = useState(0);
  const seededTrafficRef = useRef(false);
  const pausedRef = useRef(false);

  // System console stream state
  const [systemEntries, setSystemEntries] = useState<KeyedSystemLog[]>([]);
  const [systemStreamState, setSystemStreamState] = useState<StreamState>("connecting");
  const seededSystemRef = useRef(false);
  const systemSeqRef = useRef(0);

  const viewportRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    pausedRef.current = paused;
    if (!paused) {
      setQueuedCount(0);
    }
  }, [paused]);

  // Initial recent provider/model activity snapshot (oldest -> newest like a terminal)
  const initialTrafficQuery = useUsageLogs({ limit: 50, offset: 0 });
  useEffect(() => {
    if (seededTrafficRef.current || !initialTrafficQuery.data) {
      return;
    }
    seededTrafficRef.current = true;
    const chronological = [...initialTrafficQuery.data.logs].reverse().map(usageLogToMonitorEntry);
    setTrafficEntries(chronological.slice(-MAX_MONITOR_LINES));
  }, [initialTrafficQuery.data]);

  // Initial system console snapshot when switching to System tab
  const initialSystemQuery = useConsoleLogsSnapshot("ALL");
  useEffect(() => {
    if (seededSystemRef.current || !initialSystemQuery.data) {
      return;
    }
    seededSystemRef.current = true;
    const seeded = initialSystemQuery.data.entries.slice(-MAX_MONITOR_LINES).map((item) => {
      systemSeqRef.current += 1;
      return { ...item, seq: systemSeqRef.current };
    });
    setSystemEntries(seeded);
  }, [initialSystemQuery.data]);

  // Live SSE: /api/usage/stream (provider/model traffic)
  useEffect(() => {
    if (typeof EventSource === "undefined") {
      return;
    }
    setTrafficStreamState("connecting");
    const source = new EventSource("/api/usage/stream", { withCredentials: true });

    const onOpen = () => {
      setTrafficStreamState("live");
    };

    const onUsage = (evt: Event) => {
      setTrafficStreamState("live");
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      if (pausedRef.current) {
        setQueuedCount((count) => count + 1);
        return;
      }
      try {
        const parsed = JSON.parse((evt as MessageEvent).data) as Partial<UsageLog>;
        if (!parsed || typeof parsed.provider_id !== "string") {
          return;
        }
        const normalized: UsageLog = {
          id: typeof parsed.id === "number" ? parsed.id : Date.now(),
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
        const entry = usageLogToMonitorEntry(normalized);
        setTrafficEntries((prev) => appendMonitorEntry(prev, entry));
      } catch {
        // Ignore malformed SSE frame.
      }
    };

    const onError = () => {
      setTrafficStreamState("reconnecting");
    };

    source.addEventListener("open", onOpen);
    source.addEventListener("usage", onUsage);
    source.addEventListener("error", onError);

    return () => {
      source.removeEventListener("open", onOpen);
      source.removeEventListener("usage", onUsage);
      source.removeEventListener("error", onError);
      source.close();
    };
  }, []);

  // Live SSE: /api/translator/logs/stream (system console logs, only active on System tab)
  useEffect(() => {
    if (tab !== "system" || typeof EventSource === "undefined") {
      return;
    }
    setSystemStreamState("connecting");
    const source = new EventSource("/api/translator/logs/stream", { withCredentials: true });

    const onOpen = () => setSystemStreamState("live");
    const onLog = (evt: Event) => {
      setSystemStreamState("live");
      if (pausedRef.current) {
        return;
      }
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const entry = JSON.parse((evt as MessageEvent).data) as ConsoleLogEntry;
        if (!entry || typeof entry.message !== "string") {
          return;
        }
        setSystemEntries((prev) => {
          systemSeqRef.current += 1;
          const next = [...prev, { ...entry, seq: systemSeqRef.current }];
          return next.length > MAX_MONITOR_LINES
            ? next.slice(next.length - MAX_MONITOR_LINES)
            : next;
        });
      } catch {
        // Ignore malformed frame.
      }
    };
    const onError = () => setSystemStreamState("reconnecting");

    source.addEventListener("open", onOpen);
    source.addEventListener("log", onLog);
    source.addEventListener("error", onError);

    return () => {
      source.removeEventListener("open", onOpen);
      source.removeEventListener("log", onLog);
      source.removeEventListener("error", onError);
      source.close();
    };
  }, [tab]);

  const visibleTraffic = useMemo(
    () => filterMonitorEntries(trafficEntries, statusFilter, searchQuery),
    [trafficEntries, statusFilter, searchQuery],
  );

  const visibleSystem = useMemo(() => {
    const q = searchQuery.trim().toLowerCase();
    if (!q) {
      return systemEntries;
    }
    return systemEntries.filter((entry) =>
      `${entry.level} ${entry.message} ${JSON.stringify(entry.attrs ?? {})}`
        .toLowerCase()
        .includes(q),
    );
  }, [systemEntries, searchQuery]);

  const stats = useMemo(() => computeMonitorStats(trafficEntries), [trafficEntries]);

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
      setTrafficEntries([]);
      setQueuedCount(0);
    } else {
      setSystemEntries([]);
    }
  };

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
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={
                tab === "traffic"
                  ? isId
                    ? "Filter provider, model, error..."
                    : "Filter provider, model, error..."
                  : isId
                    ? "Filter log sistem..."
                    : "Filter system logs..."
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
            <div className="flex items-center gap-1.5">
              {(
                [
                  { id: "all", label: isId ? "Semua" : "All", count: stats.total },
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
            <span className="font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
              {visibleSystem.length} / {MAX_MONITOR_LINES} {isId ? "baris" : "lines"}
            </span>
          )}
        </div>

        {/* Terminal viewport */}
        <div className="relative">
          <div
            ref={viewportRef}
            onScroll={handleScroll}
            className="h-[calc(100dvh-18.5rem)] min-h-[340px] overflow-y-auto bg-[var(--bg-canvas)] px-3 py-2.5 font-mono text-xs sm:px-4 lg:h-[calc(100dvh-15.5rem)]"
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
                <div className="space-y-1 sm:space-y-0.5">
                  {visibleTraffic.map((entry, idx) => (
                    <div
                      key={entry.id}
                      className="rounded px-2 py-1.5 transition-colors hover:bg-[var(--bg-surface-raised)] sm:flex sm:items-baseline sm:gap-2.5 sm:py-1"
                    >
                      {/* Mobile line 1 / Desktop inline prefix */}
                      <div className="flex items-baseline gap-2 sm:shrink-0">
                        <span className="hidden w-7 select-none text-right text-[11px] tabular-nums text-[var(--text-muted)]/60 sm:inline-block">
                          {idx + 1}
                        </span>
                        <span className="text-[11px] tabular-nums text-[var(--text-muted)]">
                          {formatTerminalClock(entry.timestamp)}
                        </span>
                        <span
                          className={`inline-flex items-center gap-1 font-semibold tabular-nums ${levelColorClass(entry.level)}`}
                        >
                          <span>{levelSymbol(entry.level)}</span>
                          <span>{entry.statusCode}</span>
                        </span>
                        {entry.stream ? (
                          <Badge variant="default" className="px-1 py-0 text-[10px] leading-4">
                            SSE
                          </Badge>
                        ) : null}
                      </div>

                      {/* Provider / Model */}
                      <div className="mt-0.5 flex min-w-0 flex-1 flex-wrap items-baseline gap-x-2 gap-y-0.5 sm:mt-0">
                        <span className="truncate font-semibold text-[var(--accent-text)]">
                          {entry.provider}
                        </span>
                        <span className="text-[var(--text-muted)]">/</span>
                        <span className="truncate text-[var(--text-primary)]">{entry.model}</span>

                        {entry.errorSummary ? (
                          <span className="truncate text-[11px] text-[var(--status-error)]">
                            — {entry.errorSummary}
                          </span>
                        ) : null}
                      </div>

                      {/* Compact telemetry */}
                      <div className="mt-0.5 flex shrink-0 items-center gap-2.5 text-[11px] tabular-nums text-[var(--text-muted)] sm:mt-0">
                        <span>{formatShortLatency(entry.latencyMs)}</span>
                        {entry.ttftMs > 0 ? (
                          <span className="hidden md:inline">
                            TTFT {formatShortLatency(entry.ttftMs)}
                          </span>
                        ) : null}
                        <span>
                          {formatShortTokens(entry.inputTokens)}→
                          {formatShortTokens(entry.outputTokens)} tok
                        </span>
                      </div>
                    </div>
                  ))}
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
                {visibleSystem.map((line, idx) => (
                  <div
                    key={line.seq}
                    className="flex gap-2.5 rounded px-2 py-1 transition-colors hover:bg-[var(--bg-surface-raised)]"
                  >
                    <span className="hidden w-7 shrink-0 select-none text-right text-[11px] tabular-nums text-[var(--text-muted)]/60 sm:inline-block">
                      {idx + 1}
                    </span>
                    <span className="shrink-0 text-[11px] tabular-nums text-[var(--text-muted)]">
                      {formatTerminalClock(line.time)}
                    </span>
                    <span className={`shrink-0 uppercase ${systemLevelClass(line.level)}`}>
                      [{line.level}]
                    </span>
                    <span className="min-w-0 break-words text-[var(--text-primary)]">
                      {line.message}
                    </span>
                  </div>
                ))}
              </div>
            )}
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

        {/* Bottom terminal status bar */}
        <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]/40 px-4 py-2 font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
          <div className="flex items-center gap-3">
            <span>{tab === "traffic" ? "/api/usage/stream" : "/api/translator/logs/stream"}</span>
            <span className="hidden sm:inline">·</span>
            <span className="hidden sm:inline">
              {activeLength} / {MAX_MONITOR_LINES} {isId ? "baris" : "lines"}
            </span>
          </div>

          {tab === "traffic" ? (
            <div className="flex items-center gap-3">
              <span className="text-[var(--status-ok)]">OK {stats.okCount}</span>
              <span className={stats.errorCount > 0 ? "text-[var(--status-error)]" : ""}>
                ERR {stats.errorCount}
              </span>
              <span>avg {formatShortLatency(stats.avgLatencyMs)}</span>
            </div>
          ) : (
            <span>{isId ? "Ring buffer langsung" : "Live ring buffer"}</span>
          )}
        </div>
      </Card>
    </div>
  );
}
