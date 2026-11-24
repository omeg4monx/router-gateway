import {
  formatCompactTokens,
  formatUSD,
  type UsageLog,
  type UsageRange,
  type UsageSeriesBucket,
} from "@/features/usage/schemas";
import { formatUsageBucketLabel } from "./overview";

export type RequestChartMetric = "requests" | "tokens" | "cost";
export type RequestChartStyle = "area" | "bar";
export type RecentRequestStatusFilter = "all" | "ok" | "error";

export interface ChartPoint {
  index: number;
  bucket: string;
  label: string;
  value: number;
  x: number;
  y: number;
  barHeight: number;
  requests: number;
  inputTokens: number;
  outputTokens: number;
  cachedTokens: number;
  totalTokens: number;
  costUsd: number;
  errorCount: number;
}

export interface ChartYTick {
  value: number;
  y: number;
  label: string;
}

export interface RequestChartModel {
  width: number;
  height: number;
  padLeft: number;
  padRight: number;
  padTop: number;
  padBottom: number;
  plotWidth: number;
  plotHeight: number;
  baselineY: number;
  maxValue: number;
  totalValue: number;
  peakPoint: ChartPoint | null;
  points: ChartPoint[];
  yTicks: ChartYTick[];
  linePath: string;
  areaPath: string;
}

const CHART_WIDTH = 640;
const CHART_HEIGHT = 216;
const PAD_LEFT = 48;
const PAD_RIGHT = 14;
const PAD_TOP = 14;
const PAD_BOTTOM = 28;

export function getBucketMetricValue(
  bucket: UsageSeriesBucket,
  metric: RequestChartMetric,
): number {
  if (metric === "requests") {
    return Math.max(0, bucket.requests);
  }
  if (metric === "tokens") {
    return Math.max(0, bucket.input_tokens + bucket.output_tokens + bucket.cached_tokens);
  }
  return Math.max(0, bucket.cost_usd);
}

export function formatChartMetricValue(value: number, metric: RequestChartMetric): string {
  if (!Number.isFinite(value) || value <= 0) {
    return metric === "cost" ? "$0.00" : "0";
  }
  if (metric === "requests") {
    return Math.round(value).toLocaleString();
  }
  if (metric === "tokens") {
    return formatCompactTokens(value);
  }
  return formatUSD(value);
}

export function buildRequestChartModel(
  series: readonly UsageSeriesBucket[],
  range: UsageRange,
  metric: RequestChartMetric,
): RequestChartModel {
  const plotWidth = CHART_WIDTH - PAD_LEFT - PAD_RIGHT;
  const plotHeight = CHART_HEIGHT - PAD_TOP - PAD_BOTTOM;
  const baselineY = PAD_TOP + plotHeight;

  if (series.length === 0) {
    return {
      width: CHART_WIDTH,
      height: CHART_HEIGHT,
      padLeft: PAD_LEFT,
      padRight: PAD_RIGHT,
      padTop: PAD_TOP,
      padBottom: PAD_BOTTOM,
      plotWidth,
      plotHeight,
      baselineY,
      maxValue: 0,
      totalValue: 0,
      peakPoint: null,
      points: [],
      yTicks: [],
      linePath: "",
      areaPath: "",
    };
  }

  let rawMax = 0;
  let totalValue = 0;
  for (const b of series) {
    const val = getBucketMetricValue(b, metric);
    totalValue += val;
    if (val > rawMax) {
      rawMax = val;
    }
  }

  const maxValue = rawMax > 0 ? rawMax : 1;
  const count = series.length;
  const points: ChartPoint[] = series.map((bucket, index) => {
    const value = getBucketMetricValue(bucket, metric);
    const x = count === 1 ? PAD_LEFT + plotWidth / 2 : PAD_LEFT + (index / (count - 1)) * plotWidth;
    const ratio = Math.min(1, Math.max(0, value / maxValue));
    const y = baselineY - ratio * plotHeight;
    const barHeight = value > 0 ? Math.max(3, ratio * plotHeight) : 0;
    const totalTokens = bucket.input_tokens + bucket.output_tokens + bucket.cached_tokens;

    return {
      index,
      bucket: bucket.bucket,
      label: formatUsageBucketLabel(bucket.bucket, range),
      value,
      x: Number(x.toFixed(2)),
      y: Number(y.toFixed(2)),
      barHeight: Number(barHeight.toFixed(2)),
      requests: bucket.requests,
      inputTokens: bucket.input_tokens,
      outputTokens: bucket.output_tokens,
      cachedTokens: bucket.cached_tokens,
      totalTokens,
      costUsd: bucket.cost_usd,
      errorCount: Math.max(0, bucket.error_count ?? 0),
    };
  });

  let peakPoint: ChartPoint | null = null;
  for (const pt of points) {
    if (!peakPoint || pt.value > peakPoint.value) {
      peakPoint = pt;
    }
  }

  const tickRatios = [0, 0.5, 1];
  const yTicks: ChartYTick[] = tickRatios.map((ratio) => {
    const val = rawMax > 0 ? rawMax * ratio : ratio;
    const y = Number((baselineY - ratio * plotHeight).toFixed(2));
    return {
      value: val,
      y,
      label: formatChartMetricValue(val, metric),
    };
  });

  const first = points[0];
  const last = points[points.length - 1];
  let linePath = "";
  let areaPath = "";

  if (first && last) {
    if (points.length === 1) {
      const span = Math.min(44, plotWidth / 4);
      const x1 = Number(Math.max(PAD_LEFT, first.x - span).toFixed(2));
      const x2 = Number(Math.min(PAD_LEFT + plotWidth, first.x + span).toFixed(2));
      linePath = `M ${x1} ${first.y} L ${x2} ${first.y}`;
      areaPath = `M ${x1} ${baselineY} L ${x1} ${first.y} L ${x2} ${first.y} L ${x2} ${baselineY} Z`;
    } else {
      linePath = points.map((pt, idx) => `${idx === 0 ? "M" : "L"} ${pt.x} ${pt.y}`).join(" ");
      areaPath = `${linePath} L ${last.x} ${baselineY} L ${first.x} ${baselineY} Z`;
    }
  }

  return {
    width: CHART_WIDTH,
    height: CHART_HEIGHT,
    padLeft: PAD_LEFT,
    padRight: PAD_RIGHT,
    padTop: PAD_TOP,
    padBottom: PAD_BOTTOM,
    plotWidth,
    plotHeight,
    baselineY,
    maxValue: rawMax,
    totalValue,
    peakPoint,
    points,
    yTicks,
    linePath,
    areaPath,
  };
}

export function formatRelativeTimeShort(iso: string, nowMs: number, isId: boolean): string {
  const ts = Date.parse(iso);
  if (Number.isNaN(ts)) {
    return "—";
  }
  const diffSec = Math.max(0, Math.floor((nowMs - ts) / 1000));
  if (diffSec < 60) {
    return isId ? `${diffSec}d lalu` : `${diffSec}s ago`;
  }
  if (diffSec < 3600) {
    const mins = Math.floor(diffSec / 60);
    return isId ? `${mins}m lalu` : `${mins}m ago`;
  }
  if (diffSec < 86400) {
    const hours = Math.floor(diffSec / 3600);
    return isId ? `${hours}j lalu` : `${hours}h ago`;
  }
  const days = Math.floor(diffSec / 86400);
  return isId ? `${days}h lalu` : `${days}d ago`;
}

export function isUsageLogSuccess(log: Pick<UsageLog, "status_code" | "error_summary">): boolean {
  return log.status_code >= 200 && log.status_code < 400 && !log.error_summary.trim();
}

export interface FilteredRecentRequestsResult {
  items: UsageLog[];
  total: number;
  okCount: number;
  errorCount: number;
}

export function filterRecentRequests(
  logs: readonly UsageLog[],
  filter: RecentRequestStatusFilter,
  limit = 40,
): FilteredRecentRequestsResult {
  let okCount = 0;
  let errorCount = 0;
  const matched: UsageLog[] = [];

  for (const log of logs) {
    const ok = isUsageLogSuccess(log);
    if (ok) {
      okCount += 1;
    } else {
      errorCount += 1;
    }
    if (filter === "all" || (filter === "ok" && ok) || (filter === "error" && !ok)) {
      matched.push(log);
    }
  }

  return {
    items: matched.slice(0, Math.max(1, limit)),
    total: logs.length,
    okCount,
    errorCount,
  };
}
