import type { components } from "@/shared/api/schema";

export type UsageRange = "24h" | "7d" | "30d" | "all";
export type UsageLog = components["schemas"]["UsageLog"];
export type UsageLogsResponse = components["schemas"]["UsageLogsResponse"];
export type UsageSummaryTotals = components["schemas"]["UsageSummaryTotals"];
export type ModelUsageBreakdown = components["schemas"]["ModelUsageBreakdown"];
export type ProviderUsageBreakdown = components["schemas"]["ProviderUsageBreakdown"];
export type UsageSeriesBucket = components["schemas"]["UsageSeriesBucket"];
export type UsageSummaryResponse = components["schemas"]["UsageSummaryResponse"];

export const USAGE_RANGES = ["24h", "7d", "30d", "all"] as const;

export const USAGE_RANGE_LABELS: Record<UsageRange, string> = {
  "24h": "24h",
  "7d": "7d",
  "30d": "30d",
  all: "All",
};

export const DEFAULT_LOGS_LIMIT = 25;

export function isUsageLogError(log: Pick<UsageLog, "status_code" | "error_summary">): boolean {
  return log.status_code < 200 || log.status_code >= 400 || log.error_summary.trim().length > 0;
}

export function getUsageLogStatusVariant(
  log: Pick<UsageLog, "status_code" | "error_summary">,
): "success" | "warning" | "destructive" {
  if (log.status_code >= 400 && log.status_code < 500) {
    return "warning";
  }
  if (log.status_code < 200 || log.status_code >= 500 || log.error_summary.trim().length > 0) {
    return "destructive";
  }
  return "success";
}

export function computeSharePercent(part: number, total: number): number {
  if (!Number.isFinite(part) || !Number.isFinite(total) || part <= 0 || total <= 0) {
    return 0;
  }
  return Math.min(100, Math.round((part / total) * 100));
}

export function computeBarPercent(value: number, max: number): number {
  if (!Number.isFinite(value) || !Number.isFinite(max) || value <= 0 || max <= 0) {
    return 0;
  }
  return Math.min(100, Math.max(4, Math.round((value / max) * 100)));
}

export function computeTotalTokens(counts: {
  input_tokens: number;
  output_tokens: number;
  cached_tokens: number;
}): number {
  const safe = (n: number) => (Number.isFinite(n) && n > 0 ? n : 0);
  return safe(counts.input_tokens) + safe(counts.output_tokens) + safe(counts.cached_tokens);
}

const STANDARD_USD_FORMATTER = new Intl.NumberFormat("en-US", {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

export function formatUSD(amount: number): string {
  if (!Number.isFinite(amount) || amount === 0) {
    return "$0.00";
  }
  const sign = amount < 0 ? "-" : "";
  const abs = Math.abs(amount);
  if (abs < 0.01) {
    const fourDecimals = abs.toFixed(4);
    if (Number(fourDecimals) > 0) {
      return `${sign}$${fourDecimals}`;
    }
    return `${sign}$${abs.toFixed(6)}`;
  }
  return `${sign}$${STANDARD_USD_FORMATTER.format(abs)}`;
}

export function formatCompactTokens(tokens: number): string {
  if (!Number.isFinite(tokens) || tokens <= 0) {
    return "0";
  }
  if (tokens < 1_000) {
    return Math.round(tokens).toString();
  }
  const scales: Array<[number, string]> = [
    [1_000_000_000, "B"],
    [1_000_000, "M"],
    [1_000, "K"],
  ];
  for (const [divisor, suffix] of scales) {
    if (tokens >= divisor) {
      const rounded = Math.round((tokens / divisor) * 10) / 10;
      const formatted = Number.isInteger(rounded) ? rounded.toFixed(0) : rounded.toFixed(1);
      return `${formatted}${suffix}`;
    }
  }
  return Math.round(tokens).toString();
}

export function formatDurationMs(ms: number): string {
  if (!Number.isFinite(ms) || ms <= 0) {
    return "0ms";
  }
  if (ms < 1_000) {
    return `${Math.round(ms)}ms`;
  }
  return `${(ms / 1_000).toFixed(2)}s`;
}

export function formatErrorRate(rate: number): string {
  if (!Number.isFinite(rate) || rate <= 0) {
    return "0.0%";
  }
  const clamped = Math.min(1, rate);
  return `${(clamped * 100).toFixed(1)}%`;
}
