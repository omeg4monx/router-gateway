import { Activity, ArrowRight, BarChart3, Clock, Network, TrendingUp } from "lucide-react";
import { lazy, Suspense, useEffect, useId, useMemo, useRef, useState } from "react";
import { Link } from "react-router";
import type { UsageStreamStatus } from "@/features/usage/hooks";
import {
  formatCompactTokens,
  formatDurationMs,
  formatUSD,
  type UsageLog,
  type UsageRange,
  type UsageSeriesBucket,
} from "@/features/usage/schemas";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/badge";
import { formatActivityOriginLabel } from "./overview";
import {
  buildRequestChartModel,
  filterRecentRequests,
  formatChartMetricValue,
  formatRelativeTimeShort,
  isUsageLogSuccess,
  type RecentRequestStatusFilter,
  type RequestChartMetric,
  type RequestChartStyle,
} from "./overviewChart";
import type { TopologyProviderInput } from "./topologyLayout";

const LazyProviderTopology = lazy(() => import("./ProviderTopology"));

const ACTIVE_PULSE_MS = 5500;

interface OverviewRealtimeTrafficProps {
  topologyProviders: readonly TopologyProviderInput[];
  series: readonly UsageSeriesBucket[];
  logs: readonly UsageLog[];
  keyNameMap: ReadonlyMap<string, string>;
  range: UsageRange;
  rangeLabel: string;
  isLoadingChart: boolean;
  isLoadingLogs: boolean;
  streamStatus: UsageStreamStatus;
  streamEventCount: number;
  lastStreamEvent: UsageLog | null;
  isId: boolean;
}

// Single shared 1-second relative clock scoped strictly to Recent Requests rows (replaces N per-row timers).
function RecentRequestsRows({
  items,
  freshIds,
  keyNameMap,
  isId,
}: {
  items: readonly UsageLog[];
  freshIds: ReadonlySet<number>;
  keyNameMap: ReadonlyMap<string, string>;
  isId: boolean;
}) {
  const [nowMs, setNowMs] = useState(() => Date.now());

  useEffect(() => {
    if (items.length === 0) return;
    const timer = window.setInterval(() => {
      setNowMs(Date.now());
    }, 1000);
    return () => {
      window.clearInterval(timer);
    };
  }, [items.length]);

  return (
    <tbody className="divide-y divide-[var(--color-border-subtle)]/70">
      {items.map((log) => {
        const ok = isUsageLogSuccess(log);
        const isWarn = !ok && log.status_code >= 400 && log.status_code < 500;
        const isFresh = freshIds.has(log.id);
        const originLabel = formatActivityOriginLabel(log.key_id, keyNameMap, isId);
        const exactTime = (() => {
          const d = new Date(log.created_at);
          return Number.isNaN(d.getTime()) ? log.created_at : d.toLocaleString();
        })();

        return (
          <tr
            key={log.id}
            className={cn(
              "transition-colors hover:bg-[var(--color-surface-2)]/55",
              isFresh && "bg-[var(--color-primary)]/10",
            )}
            title={
              log.error_summary
                ? `HTTP ${log.status_code}: ${log.error_summary}`
                : `${log.model} via ${log.provider_id} · HTTP ${log.status_code} · ${formatDurationMs(log.latency_ms)}`
            }
          >
            <td className="py-2 pl-3 pr-1 align-top">
              <span
                className={cn(
                  "mt-1 block h-2 w-2 rounded-full",
                  ok
                    ? "bg-[var(--status-ok)]"
                    : isWarn
                      ? "bg-[var(--status-warn)]"
                      : "bg-[var(--status-error)]",
                )}
              />
            </td>
            <td className="max-w-[175px] py-1.5 pr-2 align-top">
              <div className="flex items-center gap-1.5">
                <span
                  className="truncate font-mono text-[11.5px] font-medium text-[var(--color-text-main)]"
                  title={log.model}
                >
                  {log.model}
                </span>
                {!ok ? (
                  <Badge
                    variant={isWarn ? "warning" : "destructive"}
                    className="shrink-0 px-1 py-0 font-mono text-[9.5px]"
                  >
                    {log.status_code || "ERR"}
                  </Badge>
                ) : null}
              </div>
              <div className="mt-0.5 flex items-center gap-1.5 truncate text-[10.5px] text-[var(--color-text-muted)]">
                <span className="truncate font-mono">{log.provider_id}</span>
                <span>·</span>
                <span className="truncate">{originLabel}</span>
              </div>
              {log.error_summary ? (
                <p
                  className="mt-0.5 truncate font-mono text-[10px] text-[var(--status-error)]"
                  title={log.error_summary}
                >
                  {log.error_summary}
                </p>
              ) : null}
            </td>
            <td className="py-1.5 pr-2 text-right align-top font-mono text-[11px] tabular-nums whitespace-nowrap">
              <div>
                <span className="text-[var(--color-primary)]">
                  {formatCompactTokens(log.input_tokens)}↑
                </span>{" "}
                <span className="text-[var(--status-ok)]">
                  {formatCompactTokens(log.output_tokens)}↓
                </span>
              </div>
              <div className="mt-0.5 flex items-center justify-end gap-1 text-[10px] text-[var(--color-text-subtle)]">
                {log.cached_tokens > 0 ? (
                  <span
                    className="rounded bg-[var(--color-primary)]/12 px-1 py-0.2 font-semibold text-[var(--color-primary)]"
                    title={`${log.cached_tokens.toLocaleString()} cached tokens`}
                  >
                    CACHE {formatCompactTokens(log.cached_tokens)}
                  </span>
                ) : null}
                <span>{formatDurationMs(log.latency_ms)}</span>
              </div>
            </td>
            <td
              className="py-1.5 pr-3 text-right align-top font-mono text-[10.5px] tabular-nums whitespace-nowrap text-[var(--color-text-muted)]"
              title={exactTime}
            >
              {formatRelativeTimeShort(log.created_at, nowMs, isId)}
            </td>
          </tr>
        );
      })}
    </tbody>
  );
}

const METRIC_ACCENT_VAR: Record<RequestChartMetric, string> = {
  requests: "var(--color-primary)",
  tokens: "var(--status-ok)",
  cost: "var(--status-warn)",
};

export function OverviewRealtimeTraffic({
  topologyProviders,
  series,
  logs,
  keyNameMap,
  range,
  rangeLabel,
  isLoadingChart,
  isLoadingLogs,
  streamStatus,
  streamEventCount,
  lastStreamEvent,
  isId,
}: OverviewRealtimeTrafficProps) {
  const gradientBaseId = useId();
  const [metric, setMetric] = useState<RequestChartMetric>("requests");
  const [chartStyle, setChartStyle] = useState<RequestChartStyle>("area");
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const [statusFilter, setStatusFilter] = useState<RecentRequestStatusFilter>("all");
  const [freshIds, setFreshIds] = useState<Set<number>>(() => new Set());
  const [pulsingProviders, setPulsingProviders] = useState<Set<string>>(() => new Set());

  const prevMaxLogIdRef = useRef<number | null>(null);
  const clearFreshTimerRef = useRef<number | null>(null);
  const providerPulseTimersRef = useRef<Map<string, number>>(new Map());

  // Pulse the React Flow edge whenever a live SSE usage event arrives for a provider.
  useEffect(() => {
    if (!lastStreamEvent?.provider_id) return;
    const providerKey = lastStreamEvent.provider_id.trim().toLowerCase();
    if (!providerKey) return;

    setPulsingProviders((prev) => {
      if (prev.has(providerKey)) return prev;
      const next = new Set(prev);
      next.add(providerKey);
      return next;
    });

    const existingTimer = providerPulseTimersRef.current.get(providerKey);
    if (existingTimer !== undefined) {
      window.clearTimeout(existingTimer);
    }
    const timerId = window.setTimeout(() => {
      providerPulseTimersRef.current.delete(providerKey);
      setPulsingProviders((prev) => {
        if (!prev.has(providerKey)) return prev;
        const next = new Set(prev);
        next.delete(providerKey);
        return next;
      });
    }, ACTIVE_PULSE_MS);
    providerPulseTimersRef.current.set(providerKey, timerId);
  }, [lastStreamEvent]);

  // Highlight newly arrived rows in Recent Requests & pulse provider on first live detection.
  useEffect(() => {
    const currentMaxId = logs.reduce((max, item) => Math.max(item.id, max), 0);
    if (prevMaxLogIdRef.current === null) {
      prevMaxLogIdRef.current = currentMaxId;
      return;
    }
    if (currentMaxId > prevMaxLogIdRef.current) {
      const prevMax = prevMaxLogIdRef.current;
      const incoming = new Set<number>();
      for (const item of logs) {
        if (item.id > prevMax) {
          incoming.add(item.id);
        }
      }
      prevMaxLogIdRef.current = currentMaxId;
      if (incoming.size > 0) {
        setFreshIds(incoming);
        if (clearFreshTimerRef.current !== null) {
          window.clearTimeout(clearFreshTimerRef.current);
        }
        clearFreshTimerRef.current = window.setTimeout(() => {
          setFreshIds(new Set());
          clearFreshTimerRef.current = null;
        }, 4000);
      }
    }
  }, [logs]);

  useEffect(() => {
    return () => {
      if (clearFreshTimerRef.current !== null) {
        window.clearTimeout(clearFreshTimerRef.current);
      }
      for (const timerId of providerPulseTimersRef.current.values()) {
        window.clearTimeout(timerId);
      }
      providerPulseTimersRef.current.clear();
    };
  }, []);

  const latestLog = logs[0] ?? null;
  const lastProvider = lastStreamEvent?.provider_id || latestLog?.provider_id || "";
  const errorProvider = useMemo(() => {
    if (lastStreamEvent && !isUsageLogSuccess(lastStreamEvent)) {
      return lastStreamEvent.provider_id;
    }
    if (latestLog && !isUsageLogSuccess(latestLog)) {
      return latestLog.provider_id;
    }
    return "";
  }, [lastStreamEvent, latestLog]);

  const chart = useMemo(
    () => buildRequestChartModel(series, range, metric),
    [series, range, metric],
  );

  const recentFiltered = useMemo(
    () => filterRecentRequests(logs, statusFilter, 40),
    [logs, statusFilter],
  );

  const successRatePct =
    recentFiltered.total > 0
      ? Math.round((recentFiltered.okCount / recentFiltered.total) * 100)
      : null;

  const hoveredPoint = hoveredIndex !== null ? (chart.points[hoveredIndex] ?? null) : null;
  const activePoint = hoveredPoint ?? chart.peakPoint;
  const accentColor = METRIC_ACCENT_VAR[metric];
  const gradientId = `${gradientBaseId}-${metric}`;
  const xLabelStep = chart.points.length <= 8 ? 1 : Math.ceil(chart.points.length / 7);

  const metricOptions: Array<{ value: RequestChartMetric; label: string }> = [
    { value: "requests", label: isId ? "Request" : "Requests" },
    { value: "tokens", label: isId ? "Token" : "Tokens" },
    { value: "cost", label: isId ? "Biaya" : "Cost" },
  ];

  const isStreamLive = streamStatus === "live";

  return (
    <div className="flex min-w-0 flex-col gap-4">
      {/* Row A: React Flow Provider Topology (left) + Realtime Recent Requests (right) — balanced proportional height */}
      <section className="grid grid-cols-1 items-stretch gap-4 lg:grid-cols-12">
        {/* Left: Interactive React Flow Provider Topology (7 cols) */}
        <div className="flex min-w-0 flex-col overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)] lg:col-span-7">
          <div className="flex min-h-[44px] flex-wrap items-center justify-between gap-2 border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4 py-2">
            <div className="flex min-w-0 items-center gap-2">
              <Network
                className="h-4 w-4 shrink-0 text-[var(--color-primary)]"
                aria-hidden="true"
              />
              <h3 className="truncate text-xs font-semibold text-[var(--color-text-main)]">
                {isId ? "Topologi Request Realtime" : "Live Request Topology"}
              </h3>
              <span
                className={cn(
                  "inline-flex items-center gap-1.5 rounded-full border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-2 py-0.5 font-mono text-[10px] font-semibold tracking-wider",
                  isStreamLive
                    ? "text-[var(--status-ok)]"
                    : streamStatus === "error"
                      ? "text-[var(--status-error)]"
                      : "text-[var(--color-text-muted)]",
                )}
                title={
                  isStreamLive
                    ? isId
                      ? "Tersambung realtime via SSE (/api/usage/stream)"
                      : "Connected realtime via SSE (/api/usage/stream)"
                    : isId
                      ? "Menghubungkan stream realtime..."
                      : "Connecting realtime stream..."
                }
              >
                <span className="relative flex h-1.5 w-1.5" aria-hidden="true">
                  {isStreamLive ? (
                    <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[var(--status-ok)] opacity-75" />
                  ) : null}
                  <span
                    className={cn(
                      "relative inline-flex h-1.5 w-1.5 rounded-full",
                      isStreamLive
                        ? "bg-[var(--status-ok)]"
                        : streamStatus === "error"
                          ? "bg-[var(--status-error)]"
                          : "bg-[var(--status-warn)]",
                    )}
                  />
                </span>
                {isStreamLive ? "LIVE" : streamStatus === "error" ? "RECONNECT" : "SYNC"}
              </span>
            </div>

            <div className="flex flex-wrap items-center gap-3 text-[10.5px] text-[var(--color-text-muted)]">
              <span className="inline-flex items-center gap-1">
                <span className="h-2 w-2 rounded-full bg-[var(--color-highlight)]" />
                <span>{isId ? "Aktif" : "Active"}</span>
              </span>
              <span className="inline-flex items-center gap-1">
                <span className="h-2 w-2 rounded-full bg-[var(--status-warn)]" />
                <span>{isId ? "Terakhir" : "Last"}</span>
              </span>
              <span className="inline-flex items-center gap-1">
                <span className="h-2 w-2 rounded-full bg-[var(--status-error)]" />
                <span>Error</span>
              </span>
            </div>
          </div>

          <div className="relative flex-1 bg-[var(--color-bg-alt)]/30">
            <Suspense
              fallback={
                <div className="flex h-[360px] items-center justify-center text-xs text-[var(--color-text-muted)] sm:h-[410px]">
                  {isId ? "Memuat topologi flow..." : "Loading flow topology..."}
                </div>
              }
            >
              <LazyProviderTopology
                providers={topologyProviders}
                activeProviders={pulsingProviders}
                lastProvider={lastProvider}
                errorProvider={errorProvider}
                isId={isId}
              />
            </Suspense>
          </div>

          <div className="flex min-h-[38px] flex-wrap items-center justify-between gap-2 border-t border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/30 px-4 py-2 text-[11px] text-[var(--color-text-muted)]">
            <div className="flex flex-wrap items-center gap-x-3 gap-y-1 font-mono tabular-nums">
              <span>
                {topologyProviders.length} {isId ? "provider" : "providers"}
              </span>
              {pulsingProviders.size > 0 ? (
                <span className="font-semibold text-[var(--color-highlight)]">
                  ● {pulsingProviders.size} {isId ? "aliran aktif" : "active flow"}
                </span>
              ) : lastProvider ? (
                <span>
                  {isId ? "Terakhir" : "Last"}:{" "}
                  <strong className="text-[var(--color-text-main)]">{lastProvider}</strong>
                </span>
              ) : null}
              {streamEventCount > 0 ? (
                <span className="text-[var(--status-ok)]">
                  +{streamEventCount} {isId ? "live" : "live"}
                </span>
              ) : null}
            </div>
            <Link
              to="/dashboard/providers/"
              className="inline-flex items-center gap-1 font-medium text-[var(--color-primary)] hover:underline"
            >
              <span>{isId ? "Kelola Provider" : "Manage Providers"}</span>
              <ArrowRight className="h-3 w-3" aria-hidden="true" />
            </Link>
          </div>
        </div>

        {/* Right: Realtime Recent Requests Feed (OmniRoute parity, 5 cols) */}
        <div className="flex min-w-0 flex-col overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)] lg:col-span-5">
          <div className="flex min-h-[44px] flex-wrap items-center justify-between gap-2 border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4 py-2">
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 shrink-0 text-[var(--color-primary)]" aria-hidden="true" />
              <h3 className="text-xs font-semibold text-[var(--color-text-main)]">
                {isId ? "Request Terbaru" : "Recent Requests"}
              </h3>
              <span className="rounded-full bg-[var(--color-surface)] px-2 py-0.5 font-mono text-[10.5px] tabular-nums text-[var(--color-text-muted)]">
                {recentFiltered.total}
              </span>
            </div>

            <div className="flex items-center gap-1">
              {(
                [
                  { key: "all", label: isId ? "Semua" : "All", count: recentFiltered.total },
                  { key: "ok", label: "OK", count: recentFiltered.okCount },
                  { key: "error", label: "Err", count: recentFiltered.errorCount },
                ] as const
              ).map((pill) => {
                const active = statusFilter === pill.key;
                return (
                  <button
                    key={pill.key}
                    type="button"
                    onClick={() => setStatusFilter(pill.key)}
                    className={cn(
                      "cursor-pointer rounded-md px-2 py-0.5 font-mono text-[10.5px] font-medium transition-colors",
                      active
                        ? pill.key === "error"
                          ? "bg-[var(--status-error)] text-[var(--color-on-primary)]"
                          : "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                        : "text-[var(--color-text-muted)] hover:bg-[var(--color-surface)] hover:text-[var(--color-text-main)]",
                    )}
                  >
                    {pill.label}
                    {pill.count > 0 ? ` ${pill.count}` : ""}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="h-[360px] flex-1 overflow-y-auto sm:h-[410px]">
            {isLoadingLogs && logs.length === 0 ? (
              <div className="space-y-2 p-4">
                {["sk-1", "sk-2", "sk-3", "sk-4", "sk-5"].map((key) => (
                  <div
                    key={key}
                    className="h-9 animate-pulse rounded-lg bg-[var(--color-surface-2)]/70"
                  />
                ))}
              </div>
            ) : recentFiltered.items.length === 0 ? (
              <div className="flex h-full flex-col items-center justify-center gap-2.5 px-4 text-center">
                <p className="text-xs text-[var(--color-text-muted)]">
                  {recentFiltered.total === 0
                    ? isId
                      ? "Belum ada request yang tercatat."
                      : "No requests recorded yet."
                    : isId
                      ? "Tidak ada request yang cocok dengan filter status ini."
                      : "No requests match this status filter."}
                </p>
                {recentFiltered.total === 0 ? (
                  <div className="flex flex-wrap items-center justify-center gap-2">
                    <Link
                      to="/dashboard/keys/"
                      className="rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)]/60 px-3 py-1.5 text-xs font-medium text-[var(--color-primary)] transition-colors hover:bg-[var(--color-surface-2)]"
                    >
                      {isId ? "Buat API Key →" : "Get an API key →"}
                    </Link>
                    <Link
                      to="/dashboard/chat/"
                      className="rounded-lg border border-[var(--color-border)] px-3 py-1.5 text-xs font-medium text-[var(--color-text-main)] transition-colors hover:bg-[var(--color-surface-2)]"
                    >
                      {isId ? "Uji di Chat →" : "Test in Chat →"}
                    </Link>
                  </div>
                ) : null}
              </div>
            ) : (
              <table className="w-full border-collapse text-xs">
                <thead className="sticky top-0 z-10 border-b border-[var(--color-border-subtle)] bg-[var(--color-surface)]">
                  <tr>
                    <th className="w-4 py-2 pl-3 pr-1 text-left font-semibold text-[var(--color-text-muted)]">
                      <span className="sr-only">Status</span>
                    </th>
                    <th className="py-2 pr-2 text-left text-[11px] font-semibold text-[var(--color-text-muted)]">
                      {isId ? "Model / Asal" : "Model / Origin"}
                    </th>
                    <th className="py-2 pr-2 text-right text-[11px] font-semibold whitespace-nowrap text-[var(--color-text-muted)]">
                      In / Out
                    </th>
                    <th className="py-2 pr-3 text-right text-[11px] font-semibold whitespace-nowrap text-[var(--color-text-muted)]">
                      {isId ? "Waktu" : "When"}
                    </th>
                  </tr>
                </thead>
                <RecentRequestsRows
                  items={recentFiltered.items}
                  freshIds={freshIds}
                  keyNameMap={keyNameMap}
                  isId={isId}
                />
              </table>
            )}
          </div>

          <div className="flex min-h-[38px] items-center justify-between border-t border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/30 px-4 py-2 text-[11px] text-[var(--color-text-muted)]">
            <span>
              {isStreamLive
                ? isId
                  ? "Pembaruan realtime aktif"
                  : "Realtime stream active"
                : isId
                  ? "Menunggu koneksi stream…"
                  : "Waiting for stream…"}
            </span>
            <Link
              to="/dashboard/usage/"
              className="inline-flex items-center gap-1 font-medium text-[var(--color-primary)] hover:underline"
            >
              <span>{isId ? "Semua Log" : "All Logs"}</span>
              <ArrowRight className="h-3 w-3" aria-hidden="true" />
            </Link>
          </div>
        </div>
      </section>

      {/* Row B: Time-Series Activity Chart (Requests / Tokens / Cost) */}
      <section className="flex min-w-0 flex-col overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
        <div className="flex flex-wrap items-center justify-between gap-2 border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4 py-2.5">
          <div className="flex min-w-0 items-center gap-2">
            <Activity className="h-4 w-4 shrink-0 text-[var(--color-primary)]" aria-hidden="true" />
            <h3 className="truncate text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? `Tren Aktivitas (${rangeLabel})` : `Activity Trend (${rangeLabel})`}
            </h3>
          </div>

          <div className="flex flex-wrap items-center gap-1.5">
            <div
              className="inline-flex items-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] p-0.5"
              role="tablist"
              aria-label={isId ? "Metrik grafik" : "Chart metric"}
            >
              {metricOptions.map((opt) => {
                const active = metric === opt.value;
                return (
                  <button
                    key={opt.value}
                    type="button"
                    role="tab"
                    aria-selected={active}
                    onClick={() => {
                      setMetric(opt.value);
                      setHoveredIndex(null);
                    }}
                    className={cn(
                      "cursor-pointer rounded-md px-2 py-1 text-[11px] font-medium transition-colors",
                      active
                        ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                        : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                    )}
                  >
                    {opt.label}
                  </button>
                );
              })}
            </div>

            <div className="inline-flex items-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] p-0.5">
              <button
                type="button"
                onClick={() => setChartStyle("area")}
                aria-label={isId ? "Tampilan grafik area" : "Area chart view"}
                title={isId ? "Grafik Area" : "Area Chart"}
                className={cn(
                  "flex h-6 w-6 cursor-pointer items-center justify-center rounded-md transition-colors",
                  chartStyle === "area"
                    ? "bg-[var(--color-surface-2)] text-[var(--color-primary)]"
                    : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                )}
              >
                <TrendingUp className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
              <button
                type="button"
                onClick={() => setChartStyle("bar")}
                aria-label={isId ? "Tampilan grafik batang" : "Bar chart view"}
                title={isId ? "Grafik Batang" : "Bar Chart"}
                className={cn(
                  "flex h-6 w-6 cursor-pointer items-center justify-center rounded-md transition-colors",
                  chartStyle === "bar"
                    ? "bg-[var(--color-surface-2)] text-[var(--color-primary)]"
                    : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                )}
              >
                <BarChart3 className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-3 p-4">
          <div className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-3 py-2 text-xs">
            {activePoint ? (
              <>
                <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                  <span className="font-mono text-[11px] font-semibold text-[var(--color-text-main)]">
                    {hoveredPoint
                      ? activePoint.label
                      : `${isId ? "Puncak" : "Peak"} · ${activePoint.label}`}
                  </span>
                  <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
                    <strong className="text-[var(--color-text-main)]">
                      {activePoint.requests.toLocaleString()}
                    </strong>{" "}
                    req
                  </span>
                  <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
                    <span className="text-[var(--color-primary)]">
                      {formatCompactTokens(activePoint.inputTokens)}↑
                    </span>{" "}
                    <span className="text-[var(--status-ok)]">
                      {formatCompactTokens(activePoint.outputTokens)}↓
                    </span>
                    {activePoint.cachedTokens > 0
                      ? ` · ${formatCompactTokens(activePoint.cachedTokens)} cache`
                      : ""}
                  </span>
                  <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
                    {formatUSD(activePoint.costUsd)}
                  </span>
                  {activePoint.errorCount > 0 ? (
                    <span className="font-mono text-[11px] font-medium tabular-nums text-[var(--status-error)]">
                      {activePoint.errorCount} err
                    </span>
                  ) : null}
                </div>
                <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
                  Total:{" "}
                  <strong className="text-[var(--color-text-main)]">
                    {formatChartMetricValue(chart.totalValue, metric)}
                  </strong>
                </span>
              </>
            ) : (
              <span className="text-[11px] text-[var(--color-text-muted)]">
                {isLoadingChart
                  ? isId
                    ? "Memuat data seri waktu..."
                    : "Loading time-series data..."
                  : isId
                    ? "Arahkan kursor ke grafik untuk melihat detail per interval."
                    : "Hover any interval on the chart to inspect request metrics."}
              </span>
            )}
          </div>

          {chart.points.length === 0 ? (
            <div className="flex h-44 flex-col items-center justify-center gap-1.5 rounded-lg border border-dashed border-[var(--color-border-subtle)] text-center text-xs text-[var(--color-text-muted)]">
              <Activity className="h-5 w-5 opacity-45" aria-hidden="true" />
              <span>
                {isLoadingChart
                  ? isId
                    ? "Memuat grafik..."
                    : "Loading chart..."
                  : isId
                    ? "Belum ada data seri waktu pada rentang ini."
                    : "No time-series data recorded for this period yet."}
              </span>
            </div>
          ) : (
            <div className="relative w-full overflow-hidden">
              <svg
                viewBox={`0 0 ${chart.width} ${chart.height}`}
                className="h-44 w-full select-none"
                role="img"
                aria-label={
                  isId
                    ? `Grafik ${metric} selama ${rangeLabel}`
                    : `${metric} activity chart for ${rangeLabel}`
                }
              >
                <defs>
                  <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor={accentColor} stopOpacity="0.28" />
                    <stop offset="100%" stopColor={accentColor} stopOpacity="0.02" />
                  </linearGradient>
                </defs>

                {chart.yTicks.map((tick) => (
                  <g key={tick.y}>
                    <line
                      x1={chart.padLeft}
                      y1={tick.y}
                      x2={chart.width - chart.padRight}
                      y2={tick.y}
                      stroke="var(--color-border-subtle)"
                      strokeDasharray={tick.value === 0 ? undefined : "3 3"}
                      strokeWidth="1"
                    />
                    <text
                      x={chart.padLeft - 6}
                      y={tick.y + 3.5}
                      textAnchor="end"
                      fill="var(--color-text-subtle)"
                      fontSize="10"
                      fontFamily="var(--font-mono, monospace)"
                    >
                      {tick.label}
                    </text>
                  </g>
                ))}

                {chartStyle === "area" ? (
                  <>
                    <path d={chart.areaPath} fill={`url(#${gradientId})`} />
                    <path
                      d={chart.linePath}
                      fill="none"
                      stroke={accentColor}
                      strokeWidth="2.25"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </>
                ) : (
                  chart.points.map((pt) => {
                    const slotWidth = chart.plotWidth / Math.max(1, chart.points.length);
                    const barW = Math.max(6, Math.min(28, slotWidth * 0.62));
                    const barX = pt.x - barW / 2;
                    const barY = chart.baselineY - pt.barHeight;
                    const isHovered = hoveredIndex === pt.index;
                    return (
                      <rect
                        key={pt.bucket}
                        x={barX}
                        y={barY}
                        width={barW}
                        height={pt.barHeight}
                        rx="3"
                        fill={accentColor}
                        fillOpacity={isHovered ? 1 : 0.72}
                      />
                    );
                  })
                )}

                {hoveredPoint ? (
                  <line
                    x1={hoveredPoint.x}
                    y1={chart.padTop}
                    x2={hoveredPoint.x}
                    y2={chart.baselineY}
                    stroke={accentColor}
                    strokeOpacity="0.45"
                    strokeDasharray="3 3"
                    strokeWidth="1"
                  />
                ) : null}

                {chart.points.map((pt) => {
                  const isHovered = hoveredIndex === pt.index;
                  const showXLabel =
                    pt.index === 0 ||
                    pt.index === chart.points.length - 1 ||
                    pt.index % xLabelStep === 0;

                  return (
                    <g key={pt.bucket}>
                      {chartStyle === "area" ? (
                        <circle
                          cx={pt.x}
                          cy={pt.y}
                          r={isHovered ? 4.5 : chart.points.length <= 24 ? 2.75 : 0}
                          fill={isHovered ? accentColor : "var(--color-surface)"}
                          stroke={accentColor}
                          strokeWidth="2"
                        />
                      ) : null}
                      {showXLabel ? (
                        <text
                          x={pt.x}
                          y={chart.height - 8}
                          textAnchor={
                            pt.index === 0
                              ? "start"
                              : pt.index === chart.points.length - 1
                                ? "end"
                                : "middle"
                          }
                          fill="var(--color-text-subtle)"
                          fontSize="10"
                          fontFamily="var(--font-mono, monospace)"
                        >
                          {pt.label}
                        </text>
                      ) : null}
                    </g>
                  );
                })}
              </svg>

              <div
                className="absolute inset-x-0 top-0 bottom-6 flex"
                style={{
                  paddingLeft: `${(chart.padLeft / chart.width) * 100}%`,
                  paddingRight: `${(chart.padRight / chart.width) * 100}%`,
                }}
              >
                {chart.points.map((pt) => (
                  <button
                    key={pt.bucket}
                    type="button"
                    aria-label={`${pt.label}: ${pt.requests} requests, ${formatCompactTokens(pt.totalTokens)} tokens, ${formatUSD(pt.costUsd)}`}
                    title={`${pt.label}: ${pt.requests} req · ${formatCompactTokens(pt.totalTokens)} tok · ${formatUSD(pt.costUsd)}`}
                    onMouseEnter={() => setHoveredIndex(pt.index)}
                    onMouseLeave={() => setHoveredIndex(null)}
                    onFocus={() => setHoveredIndex(pt.index)}
                    onBlur={() => setHoveredIndex(null)}
                    className="h-full flex-1 cursor-crosshair bg-transparent focus-visible:outline-1 focus-visible:outline-[var(--color-primary)]"
                  />
                ))}
              </div>
            </div>
          )}

          <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--color-border-subtle)] pt-2.5 text-[11px] text-[var(--color-text-muted)]">
            <div className="flex flex-wrap items-center gap-x-4 gap-y-1 font-mono tabular-nums">
              <span>
                {isId ? "Sukses (terbaru)" : "Recent success"}:{" "}
                <strong
                  className={cn(
                    "font-semibold",
                    successRatePct === null
                      ? "text-[var(--color-text-main)]"
                      : successRatePct >= 95
                        ? "text-[var(--status-ok)]"
                        : successRatePct >= 80
                          ? "text-[var(--status-warn)]"
                          : "text-[var(--status-error)]",
                  )}
                >
                  {successRatePct !== null ? `${successRatePct}%` : "—"}
                </strong>
              </span>
              <span>
                {isId ? "Interval" : "Buckets"}:{" "}
                <strong className="font-semibold text-[var(--color-text-main)]">
                  {chart.points.length}
                </strong>
              </span>
            </div>
            <Link
              to="/dashboard/usage/"
              className="inline-flex items-center gap-1 font-medium text-[var(--color-primary)] hover:underline"
            >
              <span>{isId ? "Analitik Lengkap" : "Full Analytics"}</span>
              <ArrowRight className="h-3 w-3" aria-hidden="true" />
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
