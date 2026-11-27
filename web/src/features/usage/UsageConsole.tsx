import {
  Activity,
  AlertCircle,
  BarChart3,
  Cpu,
  DollarSign,
  Gauge,
  RefreshCw,
  ScrollText,
  Server,
  X,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { useUsageLogs, useUsageSummary } from "./hooks";
import {
  computeBarPercent,
  computeSharePercent,
  computeTotalTokens,
  DEFAULT_LOGS_LIMIT,
  formatCompactTokens,
  formatDurationMs,
  formatErrorRate,
  formatUSD,
  getUsageLogStatusVariant,
  USAGE_RANGE_LABELS,
  USAGE_RANGES,
  type UsageLog,
  type UsageRange,
} from "./schemas";

type MobilePane = "summary" | "logs";
type BreakdownTab = "models" | "providers";

const LOG_PAGE_SIZES = [25, 50, 100] as const;

function formatTimestamp(iso: string, locale: "en" | "id"): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString(locale === "id" ? "id-ID" : "en-US", {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function getRangeLabel(range: UsageRange, isId: boolean): string {
  switch (range) {
    case "24h":
      return isId ? "24 jam terakhir" : "Last 24 hours";
    case "7d":
      return isId ? "7 hari terakhir" : "Last 7 days";
    case "30d":
      return isId ? "30 hari terakhir" : "Last 30 days";
    case "all":
      return isId ? "Semua riwayat" : "All history";
  }
}

export function UsageConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";
  const numberLocale = isId ? "id-ID" : "en-US";

  const [range, setRange] = useState<UsageRange>("24h");
  const [mobilePane, setMobilePane] = useState<MobilePane>("summary");
  const [breakdownTab, setBreakdownTab] = useState<BreakdownTab>("models");

  const [modelDraft, setModelDraft] = useState("");
  const [providerDraft, setProviderDraft] = useState("");
  const [debouncedModel, setDebouncedModel] = useState("");
  const [debouncedProvider, setDebouncedProvider] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [limit, setLimit] = useState<number>(DEFAULT_LOGS_LIMIT);
  const [offset, setOffset] = useState(0);
  const hasMountedFilterRef = useRef(false);

  useEffect(() => {
    if (!hasMountedFilterRef.current) {
      hasMountedFilterRef.current = true;
      return;
    }
    const timer = window.setTimeout(() => {
      setDebouncedModel(modelDraft.trim());
      setDebouncedProvider(providerDraft.trim());
      setOffset(0);
    }, 250);
    return () => window.clearTimeout(timer);
  }, [modelDraft, providerDraft]);

  const logsFilter = useMemo(
    () => ({
      model: debouncedModel || undefined,
      provider: debouncedProvider || undefined,
      status: statusFilter || undefined,
      limit,
      offset,
    }),
    [debouncedModel, debouncedProvider, statusFilter, limit, offset],
  );

  const summaryQuery = useUsageSummary(range);
  const logsQuery = useUsageLogs(logsFilter);

  const totals = summaryQuery.data?.totals;
  const byModel = summaryQuery.data?.by_model ?? [];
  const byProvider = summaryQuery.data?.by_provider ?? [];
  const totalTokens = totals ? computeTotalTokens(totals) : 0;

  const maxModelRequests = useMemo(
    () => byModel.reduce((max, item) => Math.max(max, item.requests), 0),
    [byModel],
  );
  const maxProviderRequests = useMemo(
    () => byProvider.reduce((max, item) => Math.max(max, item.requests), 0),
    [byProvider],
  );

  const logs: UsageLog[] = logsQuery.data?.logs ?? [];
  const logsTotal = logsQuery.data?.total ?? 0;
  const hasActiveLogFilters =
    Boolean(modelDraft.trim()) || Boolean(providerDraft.trim()) || Boolean(statusFilter);

  const isRefreshing = summaryQuery.isFetching || logsQuery.isFetching;

  const handleRefresh = () => {
    void summaryQuery.refetch();
    void logsQuery.refetch();
  };

  const handleJumpToModelLogs = (model: string) => {
    setModelDraft(model);
    setProviderDraft("");
    setDebouncedModel(model);
    setDebouncedProvider("");
    setStatusFilter("");
    setOffset(0);
    setMobilePane("logs");
  };

  const handleJumpToProviderLogs = (providerId: string) => {
    setModelDraft("");
    setProviderDraft(providerId);
    setDebouncedModel("");
    setDebouncedProvider(providerId);
    setStatusFilter("");
    setOffset(0);
    setMobilePane("logs");
  };

  const handleClearLogFilters = () => {
    setModelDraft("");
    setProviderDraft("");
    setDebouncedModel("");
    setDebouncedProvider("");
    setStatusFilter("");
    setOffset(0);
  };

  const statusOptions = useMemo(
    () => [
      { value: "", label: isId ? "Semua status" : "All statuses" },
      { value: "success", label: isId ? "Sukses (2xx/3xx)" : "Success (2xx/3xx)" },
      { value: "error", label: isId ? "Error (4xx/5xx)" : "Errors (4xx/5xx)" },
      { value: "429", label: "HTTP 429" },
      { value: "500", label: "HTTP 500" },
      { value: "502", label: "HTTP 502" },
    ],
    [isId],
  );

  const showingFrom = logsTotal === 0 ? 0 : offset + 1;
  const showingTo = Math.min(offset + logs.length, logsTotal);
  const canPrevPage = offset > 0;
  const canNextPage = offset + limit < logsTotal;
  const rangeDescription = getRangeLabel(range, isId);

  return (
    <div className="flex h-[calc(100vh-7.25rem)] min-h-[440px] flex-col gap-3 overflow-hidden supports-[height:100dvh]:h-[calc(100dvh-7.25rem)]">
      {/* Compact Single-Screen Top Bar */}
      <section className="flex shrink-0 flex-wrap items-center justify-between gap-2">
        {/* Mobile/Tablet Pane Switcher (lg+ shows both side-by-side) */}
        <div
          className="inline-flex rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-0.5 lg:hidden"
          role="tablist"
          aria-label={isId ? "Tampilan penggunaan" : "Usage workspace view"}
        >
          <button
            type="button"
            role="tab"
            aria-selected={mobilePane === "summary"}
            onClick={() => setMobilePane("summary")}
            className={cn(
              "inline-flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
              mobilePane === "summary"
                ? "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)]"
                : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
            )}
          >
            <BarChart3 className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{isId ? "Ringkasan" : "Breakdown"}</span>
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={mobilePane === "logs"}
            onClick={() => setMobilePane("logs")}
            className={cn(
              "inline-flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
              mobilePane === "logs"
                ? "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)]"
                : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
            )}
          >
            <ScrollText className="h-3.5 w-3.5" aria-hidden="true" />
            <span>
              {isId ? "Log" : "Logs"} ({logsTotal.toLocaleString(numberLocale)})
            </span>
          </button>
        </div>

        {/* Desktop compact context */}
        <div className="hidden min-w-0 items-center gap-2 lg:flex">
          <Badge variant="outline" className="tabular-nums">
            {rangeDescription}
          </Badge>
          <span className="truncate text-xs text-[var(--text-muted)]">
            {isId
              ? "Klik model atau provider untuk memfilter log request."
              : "Click a model or provider to filter request logs."}
          </span>
        </div>

        {/* Range selector + Refresh */}
        <div className="flex shrink-0 items-center gap-1.5">
          <fieldset
            aria-label={isId ? "Rentang waktu" : "Time range"}
            className="inline-flex rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] p-0.5"
          >
            {USAGE_RANGES.map((item) => {
              const active = range === item;
              return (
                <button
                  key={item}
                  type="button"
                  aria-pressed={active}
                  onClick={() => setRange(item)}
                  className={cn(
                    "rounded px-2 py-1 font-mono text-[11px] transition-colors sm:px-2.5 sm:text-xs",
                    active
                      ? "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)]"
                      : "font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
                  )}
                >
                  {item === "all" && isId ? "Semua" : USAGE_RANGE_LABELS[item]}
                </button>
              );
            })}
          </fieldset>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 px-2.5"
            onClick={handleRefresh}
            disabled={isRefreshing}
            aria-label={isId ? "Muat Ulang" : "Refresh"}
            title={isId ? "Muat Ulang" : "Refresh"}
          >
            <RefreshCw
              className={cn(
                "h-3.5 w-3.5 text-[var(--accent-primary)]",
                isRefreshing && "animate-spin",
              )}
              aria-hidden="true"
            />
            <span className="hidden sm:inline">{isId ? "Muat Ulang" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {summaryQuery.isError ? (
        <Alert variant="destructive" className="shrink-0 py-2.5">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="min-w-0 flex-1">
            <AlertTitle>
              {isId ? "Gagal memuat ringkasan penggunaan" : "Failed to load usage summary"}
            </AlertTitle>
            <AlertDescription className="truncate font-mono">
              {summaryQuery.error.message}
            </AlertDescription>
          </div>
          <Button
            type="button"
            variant="outline"
            size="xs"
            onClick={() => void summaryQuery.refetch()}
          >
            {isId ? "Coba Lagi" : "Retry"}
          </Button>
        </Alert>
      ) : null}

      {/* Compact 4-KPI Strip (2x2 mobile -> 4 cols desktop) */}
      <section
        aria-label={isId ? "Metrik utama penggunaan" : "Usage key metrics"}
        className="grid shrink-0 grid-cols-2 gap-2 lg:grid-cols-4 lg:gap-3"
      >
        <Card className="px-3 py-2.5">
          <div className="flex items-center justify-between gap-1.5">
            <span className="truncate font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Request" : "Requests"}
            </span>
            <Activity
              className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
          </div>
          {summaryQuery.isPending ? (
            <div className="mt-1.5 h-5 w-16 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
          ) : (
            <div className="mt-1 flex items-baseline justify-between gap-1.5">
              <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
                {(totals?.requests ?? 0).toLocaleString(numberLocale)}
              </span>
              <span
                className={cn(
                  "truncate font-mono text-[10px] tabular-nums",
                  (totals?.error_count ?? 0) > 0
                    ? "text-[var(--status-error)]"
                    : "text-[var(--text-muted)]",
                )}
              >
                {(totals?.error_count ?? 0).toLocaleString(numberLocale)} {isId ? "err" : "err"}
              </span>
            </div>
          )}
        </Card>

        <Card className="px-3 py-2.5">
          <div className="flex items-center justify-between gap-1.5">
            <span className="truncate font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Token" : "Tokens"}
            </span>
            <Cpu className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]" aria-hidden="true" />
          </div>
          {summaryQuery.isPending ? (
            <div className="mt-1.5 h-5 w-20 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
          ) : (
            <div className="mt-1 flex items-baseline justify-between gap-1.5">
              <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
                {formatCompactTokens(totalTokens)}
              </span>
              <span className="truncate font-mono text-[10px] text-[var(--text-muted)] tabular-nums">
                {formatCompactTokens(totals?.input_tokens ?? 0)}/
                {formatCompactTokens(totals?.output_tokens ?? 0)}
              </span>
            </div>
          )}
        </Card>

        <Card className="px-3 py-2.5">
          <div className="flex items-center justify-between gap-1.5">
            <span className="truncate font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Est. Biaya" : "Est. Cost"}
            </span>
            <DollarSign
              className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
          </div>
          {summaryQuery.isPending ? (
            <div className="mt-1.5 h-5 w-16 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
          ) : (
            <div className="mt-1 flex items-baseline justify-between gap-1.5">
              <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
                {formatUSD(totals?.cost_usd ?? 0)}
              </span>
              <span className="truncate font-mono text-[10px] text-[var(--text-muted)]">
                {USAGE_RANGE_LABELS[range]}
              </span>
            </div>
          )}
        </Card>

        <Card className="px-3 py-2.5">
          <div className="flex items-center justify-between gap-1.5">
            <span className="truncate font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Latensi / Err" : "Latency / Err"}
            </span>
            <Gauge
              className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
          </div>
          {summaryQuery.isPending ? (
            <div className="mt-1.5 h-5 w-20 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
          ) : (
            <div className="mt-1 flex items-baseline justify-between gap-1.5">
              <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
                {formatDurationMs(totals?.avg_latency_ms ?? 0)}
              </span>
              <span
                className={cn(
                  "truncate font-mono text-[10px] tabular-nums",
                  (totals?.error_rate ?? 0) > 0.05
                    ? "text-[var(--status-error)]"
                    : "text-[var(--text-muted)]",
                )}
              >
                {formatErrorRate(totals?.error_rate ?? 0)}
              </span>
            </div>
          )}
        </Card>
      </section>

      {/* Main 1-Screen Fixed Workspace: Breakdown (Left) + Request Logs (Right) */}
      <div className="grid min-h-0 flex-1 grid-cols-1 gap-3 overflow-hidden lg:grid-cols-12">
        {/* Left Pane: Models & Providers Breakdown */}
        <Card
          className={cn(
            "min-h-0 flex-col overflow-hidden lg:col-span-5 lg:flex",
            mobilePane === "summary" ? "flex h-full" : "hidden",
          )}
        >
          <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-2.5">
            <div
              className="inline-flex rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-0.5"
              role="tablist"
              aria-label={isId ? "Kategori ringkasan" : "Breakdown dimension"}
            >
              <button
                type="button"
                role="tab"
                aria-selected={breakdownTab === "models"}
                onClick={() => setBreakdownTab("models")}
                className={cn(
                  "inline-flex items-center gap-1.5 rounded px-2.5 py-1 text-xs font-medium transition-colors",
                  breakdownTab === "models"
                    ? "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)]"
                    : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
                )}
              >
                <BarChart3 className="h-3.5 w-3.5" aria-hidden="true" />
                <span>
                  {isId ? "Model" : "Models"} ({byModel.length})
                </span>
              </button>
              <button
                type="button"
                role="tab"
                aria-selected={breakdownTab === "providers"}
                onClick={() => setBreakdownTab("providers")}
                className={cn(
                  "inline-flex items-center gap-1.5 rounded px-2.5 py-1 text-xs font-medium transition-colors",
                  breakdownTab === "providers"
                    ? "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)]"
                    : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
                )}
              >
                <Server className="h-3.5 w-3.5" aria-hidden="true" />
                <span>
                  {isId ? "Provider" : "Providers"} ({byProvider.length})
                </span>
              </button>
            </div>

            <span className="hidden font-mono text-[11px] text-[var(--text-muted)] sm:inline">
              {isId ? "Urut biaya" : "By cost"}
            </span>
          </div>

          <div className="min-h-0 flex-1 overflow-y-auto p-2.5 sm:p-3">
            {summaryQuery.isPending ? (
              <div className="space-y-3 p-1">
                {[0, 1, 2, 3].map((idx) => (
                  <div key={idx} className="space-y-1.5">
                    <div className="h-3.5 w-2/3 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
                    <div className="h-1.5 w-full animate-pulse rounded bg-[var(--bg-surface-raised)]" />
                  </div>
                ))}
              </div>
            ) : breakdownTab === "models" ? (
              byModel.length === 0 ? (
                <div className="flex h-full items-center justify-center p-6 text-center">
                  <p className="text-xs text-[var(--text-muted)]">
                    {isId
                      ? `Belum ada data model pada ${rangeDescription.toLowerCase()}.`
                      : `No model telemetry recorded in ${rangeDescription.toLowerCase()}.`}
                  </p>
                </div>
              ) : (
                <div className="space-y-1.5">
                  {byModel.map((item) => {
                    const totalModelTokens = computeTotalTokens(item);
                    const barPct = computeBarPercent(item.requests, maxModelRequests);
                    const sharePct = computeSharePercent(item.requests, totals?.requests ?? 0);
                    const isSelected = debouncedModel === item.model;

                    return (
                      <button
                        key={item.model}
                        type="button"
                        onClick={() => handleJumpToModelLogs(item.model)}
                        className={cn(
                          "block w-full rounded-md border p-2 text-left transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)]",
                          isSelected
                            ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                            : "border-transparent hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)]",
                        )}
                      >
                        <div className="flex items-baseline justify-between gap-2">
                          <span className="min-w-0 truncate font-mono text-xs font-medium text-[var(--text-primary)]">
                            {item.model}
                          </span>
                          <span className="shrink-0 font-mono text-xs font-medium text-[var(--text-primary)] tabular-nums">
                            {formatUSD(item.cost_usd)}
                          </span>
                        </div>
                        <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-[var(--bg-surface-raised)]">
                          <div
                            className="h-full rounded-full bg-[var(--accent-primary)] transition-all"
                            style={{ width: `${barPct}%` }}
                          />
                        </div>
                        <div className="mt-1 flex flex-wrap items-center justify-between gap-x-2 font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                          <span>
                            {item.requests.toLocaleString(numberLocale)} ({sharePct}%) ·{" "}
                            {formatCompactTokens(totalModelTokens)} tok
                          </span>
                          <span>
                            {formatDurationMs(item.avg_latency_ms)}
                            {item.error_count > 0 ? (
                              <span className="ml-1.5 text-[var(--status-error)]">
                                · {item.error_count} err
                              </span>
                            ) : null}
                          </span>
                        </div>
                      </button>
                    );
                  })}
                </div>
              )
            ) : byProvider.length === 0 ? (
              <div className="flex h-full items-center justify-center p-6 text-center">
                <p className="text-xs text-[var(--text-muted)]">
                  {isId
                    ? `Belum ada data provider pada ${rangeDescription.toLowerCase()}.`
                    : `No provider telemetry recorded in ${rangeDescription.toLowerCase()}.`}
                </p>
              </div>
            ) : (
              <div className="space-y-1.5">
                {byProvider.map((item) => {
                  const totalProviderTokens = computeTotalTokens(item);
                  const barPct = computeBarPercent(item.requests, maxProviderRequests);
                  const sharePct = computeSharePercent(item.requests, totals?.requests ?? 0);
                  const isSelected = debouncedProvider === item.provider_id;

                  return (
                    <button
                      key={item.provider_id}
                      type="button"
                      onClick={() => handleJumpToProviderLogs(item.provider_id)}
                      className={cn(
                        "block w-full rounded-md border p-2 text-left transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)]",
                        isSelected
                          ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]"
                          : "border-transparent hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)]",
                      )}
                    >
                      <div className="flex items-baseline justify-between gap-2">
                        <span className="min-w-0 truncate font-mono text-xs font-medium text-[var(--text-primary)]">
                          {item.provider_id}
                        </span>
                        <span className="shrink-0 font-mono text-xs font-medium text-[var(--text-primary)] tabular-nums">
                          {formatUSD(item.cost_usd)}
                        </span>
                      </div>
                      <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-[var(--bg-surface-raised)]">
                        <div
                          className="h-full rounded-full bg-[var(--accent-primary)] transition-all"
                          style={{ width: `${barPct}%` }}
                        />
                      </div>
                      <div className="mt-1 flex flex-wrap items-center justify-between gap-x-2 font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                        <span>
                          {item.requests.toLocaleString(numberLocale)} ({sharePct}%) ·{" "}
                          {formatCompactTokens(totalProviderTokens)} tok
                        </span>
                        <span>
                          {formatDurationMs(item.avg_latency_ms)}
                          {item.error_count > 0 ? (
                            <span className="ml-1.5 text-[var(--status-error)]">
                              · {item.error_count} err
                            </span>
                          ) : null}
                        </span>
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </Card>

        {/* Right Pane: Request Logs (Pinned Filter + Internal Scroll + Pinned Pagination) */}
        <Card
          className={cn(
            "min-h-0 flex-col overflow-hidden lg:col-span-7 lg:flex",
            mobilePane === "logs" ? "flex h-full" : "hidden",
          )}
        >
          {/* Pinned Compact Filter Bar */}
          <div className="shrink-0 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-2.5">
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-12">
              <div className="col-span-1 sm:col-span-4">
                <Input
                  id="logs-filter-model"
                  aria-label={isId ? "Filter model" : "Filter model"}
                  type="text"
                  value={modelDraft}
                  onChange={(e) => setModelDraft(e.target.value)}
                  placeholder={isId ? "Model..." : "Model..."}
                  className="h-8 text-xs"
                />
              </div>

              <div className="col-span-1 sm:col-span-3">
                <Input
                  id="logs-filter-provider"
                  aria-label={isId ? "Filter provider" : "Filter provider"}
                  type="text"
                  value={providerDraft}
                  onChange={(e) => setProviderDraft(e.target.value)}
                  placeholder="Provider..."
                  className="h-8 text-xs"
                />
              </div>

              <div className="col-span-1 sm:col-span-3">
                <select
                  id="logs-filter-status"
                  aria-label={isId ? "Filter status" : "Filter status"}
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value);
                    setOffset(0);
                  }}
                  className="flex h-8 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-2.5 py-1 font-mono text-xs text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus:outline-none"
                >
                  {statusOptions.map((opt) => (
                    <option key={opt.value || "all"} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              </div>

              <div className="col-span-1 flex items-center gap-1.5 sm:col-span-2">
                <select
                  id="logs-page-size"
                  aria-label={isId ? "Jumlah baris" : "Rows per page"}
                  value={limit}
                  onChange={(e) => {
                    setLimit(Number(e.target.value));
                    setOffset(0);
                  }}
                  className="flex h-8 w-full min-w-0 rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-2 py-1 font-mono text-xs text-[var(--text-primary)] tabular-nums transition-colors focus:border-[var(--accent-primary)] focus:outline-none"
                >
                  {LOG_PAGE_SIZES.map((size) => (
                    <option key={size} value={size}>
                      {size}
                    </option>
                  ))}
                </select>

                {hasActiveLogFilters ? (
                  <Button
                    type="button"
                    variant="secondary"
                    size="icon-sm"
                    onClick={handleClearLogFilters}
                    aria-label={isId ? "Hapus filter" : "Clear filters"}
                    title={isId ? "Hapus filter" : "Clear filters"}
                    className="h-8 w-8 shrink-0"
                  >
                    <X className="h-3.5 w-3.5" aria-hidden="true" />
                  </Button>
                ) : null}
              </div>
            </div>
          </div>

          {/* Internal Scrollable Body */}
          <div className="min-h-0 flex-1 overflow-y-auto">
            {logsQuery.isError ? (
              <div className="p-3.5">
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertTitle>
                      {isId ? "Gagal memuat log request" : "Unable to load request logs"}
                    </AlertTitle>
                    <AlertDescription className="font-mono break-words">
                      {logsQuery.error.message}
                    </AlertDescription>
                  </div>
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    onClick={() => void logsQuery.refetch()}
                  >
                    {isId ? "Coba Lagi" : "Retry"}
                  </Button>
                </Alert>
              </div>
            ) : logsQuery.isPending ? (
              <div className="space-y-2.5 p-3.5">
                {[0, 1, 2, 3, 4].map((idx) => (
                  <div
                    key={idx}
                    className="h-12 w-full animate-pulse rounded-md bg-[var(--bg-surface-raised)]"
                  />
                ))}
              </div>
            ) : logs.length === 0 ? (
              <div className="flex h-full flex-col items-center justify-center p-6 text-center">
                <p className="text-sm font-medium text-[var(--text-primary)]">
                  {isId ? "Tidak ada log yang cocok" : "No matching request logs"}
                </p>
                <p className="mt-1 text-xs text-[var(--text-muted)]">
                  {hasActiveLogFilters
                    ? isId
                      ? "Coba ubah atau hapus filter yang aktif."
                      : "Try adjusting or clearing active filters."
                    : isId
                      ? "Kirim request ke gateway untuk mulai mencatat telemetri."
                      : "Send a request through the gateway to record telemetry."}
                </p>
                {hasActiveLogFilters ? (
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    onClick={handleClearLogFilters}
                    className="mt-3"
                  >
                    {isId ? "Hapus Filter" : "Clear Filters"}
                  </Button>
                ) : null}
              </div>
            ) : (
              <>
                {/* Mobile stacked cards (< md) */}
                <div className="divide-y divide-[var(--border-subtle)] md:hidden">
                  {logs.map((entry) => {
                    const totalEntryTokens = computeTotalTokens(entry);
                    return (
                      <article key={entry.id} className="space-y-2 p-3">
                        <div className="flex items-start justify-between gap-2">
                          <div className="min-w-0">
                            <div className="flex flex-wrap items-center gap-1.5">
                              <span className="truncate font-mono text-xs font-medium text-[var(--text-primary)]">
                                {entry.model}
                              </span>
                              {entry.stream ? <Badge variant="default">SSE</Badge> : null}
                            </div>
                            <p className="mt-0.5 truncate font-mono text-[11px] text-[var(--text-muted)]">
                              {entry.provider_id} · {formatTimestamp(entry.created_at, locale)}
                            </p>
                          </div>
                          <Badge
                            variant={getUsageLogStatusVariant(entry)}
                            className="shrink-0 tabular-nums"
                          >
                            {entry.status_code > 0 ? `HTTP ${entry.status_code}` : "ERR"}
                          </Badge>
                        </div>

                        <div className="grid grid-cols-3 gap-2 rounded bg-[var(--bg-canvas)] px-2 py-1.5 font-mono text-[11px] tabular-nums">
                          <div className="min-w-0">
                            <span className="block text-[10px] uppercase text-[var(--text-muted)]">
                              {isId ? "Token" : "Tokens"}
                            </span>
                            <span className="text-[var(--text-primary)]">
                              {formatCompactTokens(totalEntryTokens)}
                            </span>
                          </div>
                          <div className="min-w-0">
                            <span className="block text-[10px] uppercase text-[var(--text-muted)]">
                              {isId ? "Biaya" : "Cost"}
                            </span>
                            <span className="text-[var(--text-primary)]">
                              {formatUSD(entry.cost_usd)}
                            </span>
                          </div>
                          <div className="min-w-0">
                            <span className="block text-[10px] uppercase text-[var(--text-muted)]">
                              {isId ? "Latensi" : "Latency"}
                            </span>
                            <span className="text-[var(--text-secondary)]">
                              {formatDurationMs(entry.latency_ms)}
                            </span>
                          </div>
                        </div>

                        {entry.error_summary ? (
                          <p className="font-mono text-[11px] break-words text-[var(--status-error)]">
                            {entry.error_summary}
                          </p>
                        ) : null}
                      </article>
                    );
                  })}
                </div>

                {/* Desktop table (md+) with sticky header */}
                <div className="hidden md:block">
                  <table className="w-full border-collapse text-left text-xs">
                    <thead className="sticky top-0 z-10 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[10px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                      <tr>
                        <th scope="col" className="px-3.5 py-2">
                          {isId ? "Waktu" : "Time"}
                        </th>
                        <th scope="col" className="px-3 py-2">
                          {isId ? "Model" : "Model"}
                        </th>
                        <th scope="col" className="px-3 py-2">
                          {isId ? "Token" : "Tokens"}
                        </th>
                        <th scope="col" className="px-3 py-2">
                          {isId ? "Biaya" : "Cost"}
                        </th>
                        <th scope="col" className="px-3 py-2">
                          {isId ? "Latensi" : "Latency"}
                        </th>
                        <th scope="col" className="px-3.5 py-2 text-right">
                          Status
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-[var(--border-subtle)]">
                      {logs.map((entry) => (
                        <tr
                          key={entry.id}
                          className="align-top transition-colors hover:bg-[var(--bg-surface-raised)]"
                        >
                          <td className="px-3.5 py-2 font-mono text-[11px] text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
                            {formatTimestamp(entry.created_at, locale)}
                          </td>

                          <td className="max-w-[180px] px-3 py-2">
                            <div className="flex items-center gap-1.5">
                              <span
                                className="truncate font-mono text-xs font-medium text-[var(--text-primary)]"
                                title={entry.model}
                              >
                                {entry.model}
                              </span>
                              {entry.stream ? (
                                <Badge variant="default" className="shrink-0">
                                  SSE
                                </Badge>
                              ) : null}
                            </div>
                            <div className="truncate font-mono text-[10px] text-[var(--text-muted)]">
                              {entry.provider_id}
                            </div>
                          </td>

                          <td className="px-3 py-2 font-mono text-[11px] text-[var(--text-primary)] tabular-nums whitespace-nowrap">
                            {formatCompactTokens(entry.input_tokens)}/
                            {formatCompactTokens(entry.output_tokens)}
                          </td>

                          <td className="px-3 py-2 font-mono text-[11px] text-[var(--text-primary)] tabular-nums whitespace-nowrap">
                            {formatUSD(entry.cost_usd)}
                          </td>

                          <td className="px-3 py-2 font-mono text-[11px] text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
                            {formatDurationMs(entry.latency_ms)}
                          </td>

                          <td className="px-3.5 py-2 text-right">
                            <div className="flex flex-col items-end gap-0.5">
                              <Badge
                                variant={getUsageLogStatusVariant(entry)}
                                className="tabular-nums"
                              >
                                {entry.status_code > 0 ? `HTTP ${entry.status_code}` : "ERR"}
                              </Badge>
                              {entry.error_summary ? (
                                <span
                                  className="max-w-[160px] truncate font-mono text-[10px] text-[var(--status-error)]"
                                  title={entry.error_summary}
                                >
                                  {entry.error_summary}
                                </span>
                              ) : null}
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </>
            )}
          </div>

          {/* Pinned Compact Pagination Footer */}
          <div className="flex shrink-0 items-center justify-between gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-2">
            <span className="truncate font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
              {showingFrom.toLocaleString(numberLocale)}–{showingTo.toLocaleString(numberLocale)} /{" "}
              {logsTotal.toLocaleString(numberLocale)}
            </span>
            <div className="flex shrink-0 items-center gap-1.5">
              <Button
                type="button"
                variant="secondary"
                size="xs"
                disabled={!canPrevPage || logsQuery.isFetching}
                onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
              >
                {isId ? "Prev" : "Prev"}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="xs"
                disabled={!canNextPage || logsQuery.isFetching}
                onClick={() => setOffset((prev) => prev + limit)}
              >
                {isId ? "Next" : "Next"}
              </Button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
