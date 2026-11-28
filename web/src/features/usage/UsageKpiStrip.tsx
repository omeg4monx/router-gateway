import { Activity, Cpu, DollarSign, Gauge } from "lucide-react";
import { memo } from "react";
import { cn } from "@/shared/lib/utils";
import { Card } from "@/shared/ui/card";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  computeTotalTokens,
  formatCompactTokens,
  formatDurationMs,
  formatErrorRate,
  formatUSD,
  USAGE_RANGE_LABELS,
  type UsageRange,
  type UsageSummaryTotals,
} from "./schemas";

interface UsageKpiStripProps {
  totals: UsageSummaryTotals | undefined;
  range: UsageRange;
  isPending: boolean;
  isId: boolean;
  numberLocale: string;
}

export const UsageKpiStrip = memo(function UsageKpiStrip({
  totals,
  range,
  isPending,
  isId,
  numberLocale,
}: UsageKpiStripProps) {
  const totalTokens = totals ? computeTotalTokens(totals) : 0;
  const cachedTokens = totals?.cached_tokens ?? 0;
  const avgTtftMs = totals?.avg_ttft_ms ?? 0;

  return (
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
        {isPending ? (
          <Skeleton className="mt-1.5 h-5 w-16" />
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
              {(totals?.error_count ?? 0).toLocaleString(numberLocale)} err
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
        {isPending ? (
          <Skeleton className="mt-1.5 h-5 w-20" />
        ) : (
          <div className="mt-1 flex items-baseline justify-between gap-1.5">
            <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
              {formatCompactTokens(totalTokens)}
            </span>
            <span
              className="truncate font-mono text-[10px] text-[var(--text-muted)] tabular-nums"
              title={
                cachedTokens > 0
                  ? `${formatCompactTokens(totals?.input_tokens ?? 0)} in / ${formatCompactTokens(totals?.output_tokens ?? 0)} out / ${formatCompactTokens(cachedTokens)} cached`
                  : undefined
              }
            >
              {formatCompactTokens(totals?.input_tokens ?? 0)}/
              {formatCompactTokens(totals?.output_tokens ?? 0)}
              {cachedTokens > 0 ? ` +${formatCompactTokens(cachedTokens)}` : ""}
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
        {isPending ? (
          <Skeleton className="mt-1.5 h-5 w-16" />
        ) : (
          <div className="mt-1 flex items-baseline justify-between gap-1.5">
            <span className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg">
              {formatUSD(totals?.cost_usd ?? 0)}
            </span>
            <span className="truncate font-mono text-[10px] text-[var(--text-muted)]">
              {range === "all" && isId ? "Semua" : USAGE_RANGE_LABELS[range]}
            </span>
          </div>
        )}
      </Card>

      <Card className="px-3 py-2.5">
        <div className="flex items-center justify-between gap-1.5">
          <span className="truncate font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
            {isId ? "Latensi / Err" : "Latency / Err"}
          </span>
          <Gauge className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]" aria-hidden="true" />
        </div>
        {isPending ? (
          <Skeleton className="mt-1.5 h-5 w-20" />
        ) : (
          <div className="mt-1 flex items-baseline justify-between gap-1.5">
            <span
              className="font-mono text-base font-semibold tracking-tight text-[var(--text-primary)] tabular-nums sm:text-lg"
              title={avgTtftMs > 0 ? `TTFT: ${formatDurationMs(avgTtftMs)}` : undefined}
            >
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
              {avgTtftMs > 0 ? `TTFT ${formatDurationMs(avgTtftMs)} · ` : ""}
              {formatErrorRate(totals?.error_rate ?? 0)}
            </span>
          </div>
        )}
      </Card>
    </section>
  );
});
