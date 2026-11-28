import { AlertCircle, ChevronDown, ChevronRight, Search, X } from "lucide-react";
import { Fragment, memo, useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  formatCompactTokens,
  formatDurationMs,
  formatUSD,
  getUsageLogStatusVariant,
  type UsageLog,
} from "./schemas";

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

function formatTokenBreakdown(entry: UsageLog): string {
  const base = `${formatCompactTokens(entry.input_tokens)}/${formatCompactTokens(entry.output_tokens)}`;
  return entry.cached_tokens > 0 ? `${base} +${formatCompactTokens(entry.cached_tokens)}` : base;
}

interface UsageLogsPaneProps {
  visibleOnMobile: boolean;
  locale: "en" | "id";
  isId: boolean;
  numberLocale: string;
  logs: UsageLog[];
  logsTotal: number;
  isPending: boolean;
  isFetching: boolean;
  isError: boolean;
  errorMessage: string | undefined;
  searchQuery: string;
  selectedModel: string;
  selectedProvider: string;
  statusFilter: string;
  limit: number;
  offset: number;
  onSearchChange: (value: string) => void;
  onClearModel: () => void;
  onClearProvider: () => void;
  onStatusChange: (value: string) => void;
  onLimitChange: (value: number) => void;
  onClearAllFilters: () => void;
  onPrevPage: () => void;
  onNextPage: () => void;
  onRetry: () => void;
}

export const UsageLogsPane = memo(function UsageLogsPane({
  visibleOnMobile,
  locale,
  isId,
  numberLocale,
  logs,
  logsTotal,
  isPending,
  isFetching,
  isError,
  errorMessage,
  searchQuery,
  selectedModel,
  selectedProvider,
  statusFilter,
  limit,
  offset,
  onSearchChange,
  onClearModel,
  onClearProvider,
  onStatusChange,
  onLimitChange,
  onClearAllFilters,
  onPrevPage,
  onNextPage,
  onRetry,
}: UsageLogsPaneProps) {
  const [searchDraft, setSearchDraft] = useState(searchQuery);
  const [expandedLogId, setExpandedLogId] = useState<number | null>(null);
  const skipNextDebounceRef = useRef(true);

  useEffect(() => {
    setSearchDraft(searchQuery);
  }, [searchQuery]);

  useEffect(() => {
    if (skipNextDebounceRef.current) {
      skipNextDebounceRef.current = false;
      return;
    }
    const trimmed = searchDraft.trim();
    if (trimmed === searchQuery) {
      return;
    }
    const timer = window.setTimeout(() => {
      onSearchChange(trimmed);
    }, 250);
    return () => window.clearTimeout(timer);
  }, [searchDraft, searchQuery, onSearchChange]);

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

  const hasActiveLogFilters =
    Boolean(searchDraft.trim()) ||
    Boolean(selectedModel) ||
    Boolean(selectedProvider) ||
    Boolean(statusFilter);

  const showingFrom = logsTotal === 0 ? 0 : offset + 1;
  const showingTo = Math.min(offset + logs.length, logsTotal);
  const canPrevPage = offset > 0;
  const canNextPage = offset + limit < logsTotal;

  const toggleExpanded = (id: number) => {
    setExpandedLogId((prev) => (prev === id ? null : id));
  };

  return (
    <Card
      className={cn(
        "min-h-0 flex-col overflow-hidden lg:col-span-7 lg:flex lg:max-h-[calc(100dvh-14.5rem)]",
        visibleOnMobile ? "flex max-h-[calc(100dvh-16rem)] min-h-[260px]" : "hidden",
      )}
    >
      {/* Pinned Compact Filter Bar */}
      <div className="shrink-0 space-y-2 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-2.5">
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-12">
          <div className="relative col-span-2 sm:col-span-6">
            <Search
              className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
              aria-hidden="true"
            />
            <Input
              id="logs-filter-search"
              aria-label={isId ? "Cari log request" : "Search request logs"}
              type="text"
              value={searchDraft}
              onChange={(e) => setSearchDraft(e.target.value)}
              placeholder={
                isId
                  ? "Cari model, provider, request ID, error..."
                  : "Search model, provider, request ID, error..."
              }
              className="h-8 pl-8 text-xs"
            />
          </div>

          <div className="col-span-1 sm:col-span-4">
            <select
              id="logs-filter-status"
              aria-label={isId ? "Filter status HTTP" : "Filter HTTP status"}
              value={statusFilter}
              onChange={(e) => onStatusChange(e.target.value)}
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
              onChange={(e) => onLimitChange(Number(e.target.value))}
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
                onClick={() => {
                  setSearchDraft("");
                  onClearAllFilters();
                }}
                aria-label={isId ? "Hapus filter" : "Clear filters"}
                title={isId ? "Hapus filter" : "Clear filters"}
                className="h-8 w-8 shrink-0"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
              </Button>
            ) : null}
          </div>
        </div>

        {/* Active exact-match chips from Breakdown click */}
        {selectedModel || selectedProvider ? (
          <div className="flex flex-wrap items-center gap-1.5">
            {selectedModel ? (
              <button
                type="button"
                onClick={onClearModel}
                className="inline-flex items-center gap-1 rounded border border-[var(--accent-primary)] bg-[var(--bg-surface-raised)] px-2 py-0.5 font-mono text-[11px] text-[var(--text-primary)] transition-colors hover:bg-[var(--bg-surface)]"
                aria-label={
                  isId
                    ? `Hapus filter model ${selectedModel}`
                    : `Remove model filter ${selectedModel}`
                }
              >
                <span className="text-[var(--text-muted)]">model:</span>
                <span className="max-w-[200px] truncate">{selectedModel}</span>
                <X className="h-3 w-3 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
              </button>
            ) : null}
            {selectedProvider ? (
              <button
                type="button"
                onClick={onClearProvider}
                className="inline-flex items-center gap-1 rounded border border-[var(--accent-primary)] bg-[var(--bg-surface-raised)] px-2 py-0.5 font-mono text-[11px] text-[var(--text-primary)] transition-colors hover:bg-[var(--bg-surface)]"
                aria-label={
                  isId
                    ? `Hapus filter provider ${selectedProvider}`
                    : `Remove provider filter ${selectedProvider}`
                }
              >
                <span className="text-[var(--text-muted)]">provider:</span>
                <span className="max-w-[200px] truncate">{selectedProvider}</span>
                <X className="h-3 w-3 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
              </button>
            ) : null}
          </div>
        ) : null}
      </div>

      {/* Internal Scrollable Body */}
      <div className="min-h-0 flex-1 overflow-y-auto">
        {isError ? (
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
                  {errorMessage}
                </AlertDescription>
              </div>
              <Button type="button" variant="outline" size="xs" onClick={onRetry}>
                {isId ? "Coba Lagi" : "Retry"}
              </Button>
            </Alert>
          </div>
        ) : isPending ? (
          <div className="space-y-2.5 p-3.5">
            {[0, 1, 2, 3, 4].map((idx) => (
              <Skeleton key={idx} className="h-12 w-full rounded-md" />
            ))}
          </div>
        ) : logs.length === 0 ? (
          <div className="flex min-h-[160px] flex-col items-center justify-center p-6 text-center">
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
                onClick={() => {
                  setSearchDraft("");
                  onClearAllFilters();
                }}
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
                const isExpanded = expandedLogId === entry.id;
                return (
                  <article key={entry.id} className="space-y-2 p-3">
                    <button
                      type="button"
                      onClick={() => toggleExpanded(entry.id)}
                      aria-expanded={isExpanded}
                      className="flex w-full items-start justify-between gap-2 text-left"
                    >
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
                    </button>

                    <div className="grid grid-cols-3 gap-2 rounded bg-[var(--bg-canvas)] px-2 py-1.5 font-mono text-[11px] tabular-nums">
                      <div className="min-w-0">
                        <span className="block text-[10px] uppercase text-[var(--text-muted)]">
                          {isId ? "Token" : "Tokens"}
                        </span>
                        <span className="truncate text-[var(--text-primary)]">
                          {formatTokenBreakdown(entry)}
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

                    {isExpanded ? (
                      <dl className="grid grid-cols-2 gap-x-3 gap-y-1.5 rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-2 font-mono text-[11px]">
                        <div>
                          <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                            Request ID
                          </dt>
                          <dd className="break-all text-[var(--text-primary)]">
                            {entry.request_id || "—"}
                          </dd>
                        </div>
                        <div>
                          <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                            TTFT / Protocol
                          </dt>
                          <dd className="text-[var(--text-secondary)]">
                            {entry.ttft_ms > 0 ? formatDurationMs(entry.ttft_ms) : "—"} ·{" "}
                            {entry.protocol || "—"}
                          </dd>
                        </div>
                        <div>
                          <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                            Connection
                          </dt>
                          <dd className="break-all text-[var(--text-secondary)]">
                            {entry.connection_id || "—"}
                          </dd>
                        </div>
                        <div>
                          <dt className="text-[10px] uppercase text-[var(--text-muted)]">Key ID</dt>
                          <dd className="break-all text-[var(--text-secondary)]">
                            {entry.key_id || "—"}
                          </dd>
                        </div>
                      </dl>
                    ) : null}
                  </article>
                );
              })}
            </div>

            {/* Desktop table (md+) with sticky header & expandable telemetry detail */}
            <div className="hidden md:block">
              <table className="w-full border-collapse text-left text-xs">
                <thead className="sticky top-0 z-10 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[10px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <tr>
                    <th scope="col" className="w-7 pl-2.5 pr-0 py-2">
                      <span className="sr-only">{isId ? "Detail" : "Details"}</span>
                    </th>
                    <th scope="col" className="px-2.5 py-2">
                      {isId ? "Waktu" : "Time"}
                    </th>
                    <th scope="col" className="px-3 py-2">
                      Model
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
                  {logs.map((entry) => {
                    const isExpanded = expandedLogId === entry.id;
                    return (
                      <Fragment key={entry.id}>
                        <tr
                          onClick={() => toggleExpanded(entry.id)}
                          className={cn(
                            "cursor-pointer align-top transition-colors hover:bg-[var(--bg-surface-raised)] focus-within:bg-[var(--bg-surface-raised)]",
                            isExpanded && "bg-[var(--bg-surface-raised)]",
                          )}
                        >
                          <td className="pl-2.5 pr-0 py-2.5 text-[var(--text-muted)]">
                            <button
                              type="button"
                              aria-expanded={isExpanded}
                              aria-label={isId ? "Tampilkan detail log" : "Toggle log details"}
                              onClick={(e) => {
                                e.stopPropagation();
                                toggleExpanded(entry.id);
                              }}
                              className="inline-flex items-center justify-center rounded p-0.5 hover:text-[var(--text-primary)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)]"
                            >
                              {isExpanded ? (
                                <ChevronDown className="h-3.5 w-3.5" aria-hidden="true" />
                              ) : (
                                <ChevronRight className="h-3.5 w-3.5" aria-hidden="true" />
                              )}
                            </button>
                          </td>

                          <td className="px-2.5 py-2 font-mono text-[11px] text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
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
                            {formatTokenBreakdown(entry)}
                          </td>

                          <td className="px-3 py-2 font-mono text-[11px] text-[var(--text-primary)] tabular-nums whitespace-nowrap">
                            {formatUSD(entry.cost_usd)}
                          </td>

                          <td className="px-3 py-2 font-mono text-[11px] text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
                            <div>{formatDurationMs(entry.latency_ms)}</div>
                            {entry.ttft_ms > 0 ? (
                              <div className="text-[10px] text-[var(--text-muted)]">
                                TTFT {formatDurationMs(entry.ttft_ms)}
                              </div>
                            ) : null}
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

                        {isExpanded ? (
                          <tr className="bg-[var(--bg-canvas)]">
                            <td colSpan={7} className="px-4 py-2.5">
                              <dl className="grid grid-cols-2 gap-x-4 gap-y-1.5 font-mono text-[11px] xl:grid-cols-4">
                                <div>
                                  <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                                    Request ID
                                  </dt>
                                  <dd className="break-all text-[var(--text-primary)]">
                                    {entry.request_id || "—"}
                                  </dd>
                                </div>
                                <div>
                                  <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                                    Protocol
                                  </dt>
                                  <dd className="text-[var(--text-secondary)]">
                                    {entry.protocol || "—"}
                                  </dd>
                                </div>
                                <div>
                                  <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                                    Connection
                                  </dt>
                                  <dd className="break-all text-[var(--text-secondary)]">
                                    {entry.connection_id || "—"}
                                  </dd>
                                </div>
                                <div>
                                  <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                                    Key ID
                                  </dt>
                                  <dd className="break-all text-[var(--text-secondary)]">
                                    {entry.key_id || "—"}
                                  </dd>
                                </div>
                                <div className="col-span-2 xl:col-span-4">
                                  <dt className="text-[10px] uppercase text-[var(--text-muted)]">
                                    Tokens (In / Out / Cached)
                                  </dt>
                                  <dd className="text-[var(--text-secondary)] tabular-nums">
                                    {entry.input_tokens.toLocaleString(numberLocale)} /{" "}
                                    {entry.output_tokens.toLocaleString(numberLocale)} /{" "}
                                    {entry.cached_tokens.toLocaleString(numberLocale)}
                                  </dd>
                                </div>
                                {entry.error_summary ? (
                                  <div className="col-span-2 xl:col-span-4">
                                    <dt className="text-[10px] uppercase text-[var(--status-error)]">
                                      Error
                                    </dt>
                                    <dd className="break-words text-[var(--status-error)]">
                                      {entry.error_summary}
                                    </dd>
                                  </div>
                                ) : null}
                              </dl>
                            </td>
                          </tr>
                        ) : null}
                      </Fragment>
                    );
                  })}
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
            disabled={!canPrevPage || isFetching}
            onClick={onPrevPage}
          >
            {isId ? "Sebelumnya" : "Prev"}
          </Button>
          <Button
            type="button"
            variant="secondary"
            size="xs"
            disabled={!canNextPage || isFetching}
            onClick={onNextPage}
          >
            {isId ? "Berikutnya" : "Next"}
          </Button>
        </div>
      </div>
    </Card>
  );
});
