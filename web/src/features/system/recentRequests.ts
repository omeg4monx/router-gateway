import type { UsageLog } from "@/features/usage/schemas";

export type RecentRequestStatusFilter = "all" | "ok" | "error";

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

export function computeRelativeClockIntervalMs(
  items: readonly Pick<UsageLog, "created_at">[],
  nowMs: number,
): number | null {
  if (items.length === 0) {
    return null;
  }
  for (const item of items) {
    const ts = Date.parse(item.created_at);
    if (!Number.isNaN(ts) && nowMs - ts < 60_000) {
      return 1_000;
    }
  }
  return 30_000;
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
