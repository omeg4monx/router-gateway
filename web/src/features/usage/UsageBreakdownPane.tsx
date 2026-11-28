import { BarChart3, Server } from "lucide-react";
import { memo, useMemo, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Card } from "@/shared/ui/card";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  computeBarPercent,
  computeSharePercent,
  computeTotalTokens,
  formatCompactTokens,
  formatDurationMs,
  formatUSD,
  type ModelUsageBreakdown,
  type ProviderUsageBreakdown,
  type UsageSeriesBucket,
} from "./schemas";

type BreakdownTab = "models" | "providers";

interface UsageBreakdownPaneProps {
  visibleOnMobile: boolean;
  isPending: boolean;
  isId: boolean;
  numberLocale: string;
  rangeDescription: string;
  totalRequests: number;
  byModel: ModelUsageBreakdown[];
  byProvider: ProviderUsageBreakdown[];
  series: UsageSeriesBucket[];
  selectedModel: string;
  selectedProvider: string;
  onToggleModel: (model: string) => void;
  onToggleProvider: (providerId: string) => void;
}

export const UsageBreakdownPane = memo(function UsageBreakdownPane({
  visibleOnMobile,
  isPending,
  isId,
  numberLocale,
  rangeDescription,
  totalRequests,
  byModel,
  byProvider,
  series,
  selectedModel,
  selectedProvider,
  onToggleModel,
  onToggleProvider,
}: UsageBreakdownPaneProps) {
  const [breakdownTab, setBreakdownTab] = useState<BreakdownTab>("models");

  const maxModelRequests = useMemo(
    () => byModel.reduce((max, item) => Math.max(max, item.requests), 0),
    [byModel],
  );
  const maxProviderRequests = useMemo(
    () => byProvider.reduce((max, item) => Math.max(max, item.requests), 0),
    [byProvider],
  );
  const maxSeriesRequests = useMemo(
    () => series.reduce((max, item) => Math.max(max, item.requests), 0),
    [series],
  );

  return (
    <Card
      className={cn(
        "min-h-0 flex-col overflow-hidden lg:col-span-5 lg:flex lg:max-h-[calc(100dvh-14.5rem)]",
        visibleOnMobile ? "flex max-h-[calc(100dvh-16rem)] min-h-[260px]" : "hidden",
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

      {/* Compact Activity Sparkline using summary.series */}
      {!isPending && series.length > 0 ? (
        <section
          className="shrink-0 border-b border-[var(--border-subtle)] bg-[var(--bg-surface)] px-3.5 py-2"
          aria-label={isId ? "Tren aktivitas" : "Activity trend"}
        >
          <div className="mb-1 flex items-center justify-between font-mono text-[10px] text-[var(--text-muted)]">
            <span className="uppercase tracking-wider">{isId ? "Tren" : "Trend"}</span>
            <span className="tabular-nums">
              {series.length} {isId ? "bucket" : series.length === 1 ? "bucket" : "buckets"}
            </span>
          </div>
          <div className="flex h-7 items-end gap-0.5">
            {series.map((bucket) => {
              const heightPct = computeBarPercent(bucket.requests, maxSeriesRequests);
              const hasErrors = bucket.error_count > 0;
              const tooltip = `${bucket.bucket} · ${bucket.requests.toLocaleString(numberLocale)} req · ${bucket.error_count} err · ${formatUSD(bucket.cost_usd)}`;
              return (
                <div
                  key={bucket.bucket}
                  title={tooltip}
                  className="flex h-full min-w-[3px] flex-1 items-end"
                >
                  <div
                    className={cn(
                      "w-full rounded-t-[2px] transition-all",
                      hasErrors ? "bg-[var(--status-error)]/80" : "bg-[var(--accent-primary)]/85",
                    )}
                    style={{ height: `${heightPct}%` }}
                  />
                </div>
              );
            })}
          </div>
        </section>
      ) : null}

      <div className="min-h-0 flex-1 overflow-y-auto p-2.5 sm:p-3">
        {isPending ? (
          <div className="space-y-3 p-1" role="status" aria-busy="true">
            {[0, 1, 2, 3].map((idx) => (
              <div key={idx} className="space-y-1.5">
                <Skeleton className="h-3.5 w-2/3" />
                <Skeleton className="h-1.5 w-full" />
              </div>
            ))}
          </div>
        ) : breakdownTab === "models" ? (
          byModel.length === 0 ? (
            <div className="flex min-h-[140px] items-center justify-center p-6 text-center">
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
                const sharePct = computeSharePercent(item.requests, totalRequests);
                const isSelected = selectedModel === item.model;

                return (
                  <button
                    key={item.model}
                    type="button"
                    aria-pressed={isSelected}
                    onClick={() => onToggleModel(item.model)}
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
          <div className="flex min-h-[140px] items-center justify-center p-6 text-center">
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
              const sharePct = computeSharePercent(item.requests, totalRequests);
              const isSelected = selectedProvider === item.provider_id;

              return (
                <button
                  key={item.provider_id}
                  type="button"
                  aria-pressed={isSelected}
                  onClick={() => onToggleProvider(item.provider_id)}
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
  );
});
