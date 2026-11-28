import { AlertCircle, BarChart3, DollarSign, RefreshCw, ScrollText } from "lucide-react";
import { lazy, Suspense, useCallback, useMemo, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { useUsageLogs, useUsageSummary } from "./hooks";
import {
  DEFAULT_LOGS_LIMIT,
  type ModelUsageBreakdown,
  type ProviderUsageBreakdown,
  USAGE_RANGE_LABELS,
  USAGE_RANGES,
  type UsageLog,
  type UsageRange,
  type UsageSeriesBucket,
} from "./schemas";
import { UsageBreakdownPane } from "./UsageBreakdownPane";
import { UsageKpiStrip } from "./UsageKpiStrip";
import { UsageLogsPane } from "./UsageLogsPane";
import { UsageStreamIndicator } from "./UsageStreamIndicator";

const UsagePricingModal = lazy(() =>
  import("./UsagePricingModal").then((m) => ({ default: m.UsagePricingModal })),
);

const EMPTY_MODELS: ModelUsageBreakdown[] = [];
const EMPTY_PROVIDERS: ProviderUsageBreakdown[] = [];
const EMPTY_SERIES: UsageSeriesBucket[] = [];
const EMPTY_LOGS: UsageLog[] = [];

type MobilePane = "summary" | "logs";

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
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedModel, setSelectedModel] = useState("");
  const [selectedProvider, setSelectedProvider] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [limit, setLimit] = useState<number>(DEFAULT_LOGS_LIMIT);
  const [offset, setOffset] = useState(0);
  const [pricingModalOpen, setPricingModalOpen] = useState(false);

  const logsFilter = useMemo(
    () => ({
      range,
      q: searchQuery || undefined,
      model: selectedModel || undefined,
      provider: selectedProvider || undefined,
      status: statusFilter || undefined,
      limit,
      offset,
    }),
    [range, searchQuery, selectedModel, selectedProvider, statusFilter, limit, offset],
  );

  const summaryQuery = useUsageSummary(range);
  const logsQuery = useUsageLogs(logsFilter);
  const { refetch: refetchSummary } = summaryQuery;
  const { refetch: refetchLogs } = logsQuery;

  const totals = summaryQuery.data?.totals;
  const byModel = summaryQuery.data?.by_model ?? EMPTY_MODELS;
  const byProvider = summaryQuery.data?.by_provider ?? EMPTY_PROVIDERS;
  const series = summaryQuery.data?.series ?? EMPTY_SERIES;

  const logs = logsQuery.data?.logs ?? EMPTY_LOGS;
  const logsTotal = logsQuery.data?.total ?? 0;
  const isRefreshing = summaryQuery.isFetching || logsQuery.isFetching;
  const rangeDescription = getRangeLabel(range, isId);

  const handleRangeChange = useCallback((nextRange: UsageRange) => {
    setRange(nextRange);
    setOffset(0);
  }, []);

  const handleRefresh = useCallback(() => {
    void refetchSummary();
    void refetchLogs();
  }, [refetchSummary, refetchLogs]);

  const handleClosePricingModal = useCallback(() => {
    setPricingModalOpen(false);
  }, []);

  const handleToggleModel = useCallback((model: string) => {
    setSelectedModel((prev) => (prev === model ? "" : model));
    setSelectedProvider("");
    setOffset(0);
    setMobilePane("logs");
  }, []);

  const handleToggleProvider = useCallback((providerId: string) => {
    setSelectedProvider((prev) => (prev === providerId ? "" : providerId));
    setSelectedModel("");
    setOffset(0);
    setMobilePane("logs");
  }, []);

  const handleSearchChange = useCallback((value: string) => {
    setSearchQuery(value);
    setOffset(0);
  }, []);

  const handleClearModel = useCallback(() => {
    setSelectedModel("");
    setOffset(0);
  }, []);

  const handleClearProvider = useCallback(() => {
    setSelectedProvider("");
    setOffset(0);
  }, []);

  const handleStatusChange = useCallback((value: string) => {
    setStatusFilter(value);
    setOffset(0);
  }, []);

  const handleLimitChange = useCallback((value: number) => {
    setLimit(value);
    setOffset(0);
  }, []);

  const handleClearAllFilters = useCallback(() => {
    setSearchQuery("");
    setSelectedModel("");
    setSelectedProvider("");
    setStatusFilter("");
    setOffset(0);
  }, []);

  const handlePrevPage = useCallback(() => {
    setOffset((prev) => Math.max(0, prev - limit));
  }, [limit]);

  const handleNextPage = useCallback(() => {
    setOffset((prev) => prev + limit);
  }, [limit]);

  const handleRetryLogs = useCallback(() => {
    void refetchLogs();
  }, [refetchLogs]);

  return (
    <div className="flex flex-col gap-3">
      {/* Compact Top Bar */}
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
              ? "Klik model atau provider untuk filter log; klik baris log untuk detail."
              : "Click a model or provider to filter logs; click a log row for details."}
          </span>
        </div>

        {/* Live SSE indicator + Range selector + Refresh */}
        <div className="flex shrink-0 items-center gap-1.5">
          <UsageStreamIndicator isId={isId} />

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
                  onClick={() => handleRangeChange(item)}
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
            onClick={() => setPricingModalOpen(true)}
            aria-label={isId ? "Tarif Model" : "Model Pricing"}
            title={isId ? "Tarif Model" : "Model Pricing"}
          >
            <DollarSign className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
            <span className="hidden sm:inline">{isId ? "Tarif" : "Pricing"}</span>
          </Button>

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

      <UsageKpiStrip
        totals={totals}
        range={range}
        isPending={summaryQuery.isPending}
        isId={isId}
        numberLocale={numberLocale}
      />

      {/* Main Responsive Workspace: Breakdown (Left) + Request Logs (Right) */}
      <div className="grid grid-cols-1 items-start gap-3 lg:grid-cols-12">
        <UsageBreakdownPane
          visibleOnMobile={mobilePane === "summary"}
          isPending={summaryQuery.isPending}
          isId={isId}
          numberLocale={numberLocale}
          rangeDescription={rangeDescription}
          totalRequests={totals?.requests ?? 0}
          byModel={byModel}
          byProvider={byProvider}
          series={series}
          selectedModel={selectedModel}
          selectedProvider={selectedProvider}
          onToggleModel={handleToggleModel}
          onToggleProvider={handleToggleProvider}
        />

        <UsageLogsPane
          visibleOnMobile={mobilePane === "logs"}
          locale={locale}
          isId={isId}
          numberLocale={numberLocale}
          logs={logs}
          logsTotal={logsTotal}
          isPending={logsQuery.isPending}
          isFetching={logsQuery.isFetching}
          isError={logsQuery.isError}
          errorMessage={logsQuery.error?.message}
          searchQuery={searchQuery}
          selectedModel={selectedModel}
          selectedProvider={selectedProvider}
          statusFilter={statusFilter}
          limit={limit}
          offset={offset}
          onSearchChange={handleSearchChange}
          onClearModel={handleClearModel}
          onClearProvider={handleClearProvider}
          onStatusChange={handleStatusChange}
          onLimitChange={handleLimitChange}
          onClearAllFilters={handleClearAllFilters}
          onPrevPage={handlePrevPage}
          onNextPage={handleNextPage}
          onRetry={handleRetryLogs}
        />
      </div>

      {pricingModalOpen ? (
        <Suspense fallback={null}>
          <UsagePricingModal isId={isId} onClose={handleClosePricingModal} />
        </Suspense>
      ) : null}
    </div>
  );
}
