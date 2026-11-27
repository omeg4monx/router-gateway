import type { components } from "@/shared/api/schema";
import type { UILocale } from "@/shared/store/uiStore";

export type ProviderQuotaStatus = components["schemas"]["ProviderQuotaStatus"];
export type ProviderQuotaBucket = components["schemas"]["ProviderQuotaBucket"];
export type ProviderConnectionQuota = components["schemas"]["ProviderConnectionQuota"];
export type ProviderQuotasResponse = components["schemas"]["ProviderQuotasResponse"];

export type QuotaStatusFilter = "all" | ProviderQuotaStatus;
export type QuotaBadgeVariant = "success" | "warning" | "destructive" | "secondary";
export type BucketProgressTone = "ok" | "warn" | "error";

export interface QuotaOverviewStats {
  total: number;
  healthy: number;
  low: number;
  exhausted: number;
  notActivated: number;
  error: number;
}

export interface QuotaFilterOptions {
  query: string;
  status: QuotaStatusFilter;
  providerId: string;
}

export interface ProviderQuotaGroup {
  providerId: string;
  providerName: string;
  fetcherId: string;
  items: ProviderConnectionQuota[];
  stats: QuotaOverviewStats;
  avgRemainingPct: number | null;
}

export function groupProviderQuotas(
  quotas: readonly ProviderConnectionQuota[],
): ProviderQuotaGroup[] {
  const buckets = new Map<
    string,
    {
      providerId: string;
      providerName: string;
      fetcherId: string;
      items: ProviderConnectionQuota[];
    }
  >();

  for (const item of quotas) {
    const existing = buckets.get(item.provider_id);
    if (existing) {
      existing.items.push(item);
    } else {
      buckets.set(item.provider_id, {
        providerId: item.provider_id,
        providerName: item.provider_name,
        fetcherId: item.fetcher_id,
        items: [item],
      });
    }
  }

  return Array.from(buckets.values()).map((entry) => {
    let remainingSum = 0;
    let remainingCount = 0;

    for (const item of entry.items) {
      if (item.buckets.length === 0) {
        continue;
      }
      let minRemaining = 100;
      for (const bucket of item.buckets) {
        const clamped = Math.max(0, Math.min(100, bucket.remaining_pct));
        if (clamped < minRemaining) {
          minRemaining = clamped;
        }
      }
      remainingSum += minRemaining;
      remainingCount += 1;
    }

    const avgRemainingPct =
      remainingCount > 0 ? Math.round((remainingSum / remainingCount) * 10) / 10 : null;

    return {
      providerId: entry.providerId,
      providerName: entry.providerName,
      fetcherId: entry.fetcherId,
      items: entry.items,
      stats: computeQuotaOverviewStats(entry.items),
      avgRemainingPct,
    };
  });
}

export function computeQuotaOverviewStats(
  quotas: readonly ProviderConnectionQuota[],
): QuotaOverviewStats {
  let healthy = 0;
  let low = 0;
  let exhausted = 0;
  let notActivated = 0;
  let error = 0;

  for (const item of quotas) {
    switch (item.status) {
      case "ok":
        healthy += 1;
        break;
      case "low":
        low += 1;
        break;
      case "exhausted":
        exhausted += 1;
        break;
      case "not_activated":
        notActivated += 1;
        break;
      case "error":
        error += 1;
        break;
    }
  }

  return {
    total: quotas.length,
    healthy,
    low,
    exhausted,
    notActivated,
    error,
  };
}

export function patchConnectionQuotaInResponse(
  current: ProviderQuotasResponse,
  updated: ProviderConnectionQuota,
): ProviderQuotasResponse {
  const exists = current.quotas.some((item) => item.connection_id === updated.connection_id);
  return {
    ...current,
    fetched_at: updated.fetched_at,
    quotas: exists
      ? current.quotas.map((item) =>
          item.connection_id === updated.connection_id ? updated : item,
        )
      : [...current.quotas, updated],
  };
}

export function filterProviderQuotas(
  quotas: readonly ProviderConnectionQuota[],
  options: QuotaFilterOptions,
): ProviderConnectionQuota[] {
  const normalizedQuery = options.query.trim().toLowerCase();
  const normalizedProvider = options.providerId.trim().toLowerCase();

  return quotas.filter((item) => {
    if (options.status !== "all" && item.status !== options.status) {
      return false;
    }
    if (
      normalizedProvider &&
      normalizedProvider !== "all" &&
      item.provider_id.toLowerCase() !== normalizedProvider &&
      item.fetcher_id.toLowerCase() !== normalizedProvider
    ) {
      return false;
    }
    if (!normalizedQuery) {
      return true;
    }

    const haystack = [
      item.connection_name,
      item.provider_name,
      item.provider_id,
      item.plan_name ?? "",
      item.error,
      ...item.models,
      ...item.buckets.map((bucket) => bucket.label),
    ]
      .join(" ")
      .toLowerCase();

    return haystack.includes(normalizedQuery);
  });
}

export function getQuotaStatusVariant(status: ProviderQuotaStatus): QuotaBadgeVariant {
  switch (status) {
    case "ok":
      return "success";
    case "low":
      return "warning";
    case "exhausted":
    case "error":
      return "destructive";
    case "not_activated":
      return "secondary";
  }
}

export function getBucketProgressTone(
  remainingPct: number,
  exhausted: boolean,
): BucketProgressTone {
  if (exhausted || remainingPct <= 0) {
    return "error";
  }
  if (remainingPct <= 20) {
    return "warn";
  }
  return "ok";
}

export function formatResetCountdown(
  iso: string | null | undefined,
  locale: UILocale,
  nowMs: number = Date.now(),
): string | null {
  if (!iso) {
    return null;
  }
  const targetMs = Date.parse(iso);
  if (Number.isNaN(targetMs)) {
    return null;
  }
  const diffMs = targetMs - nowMs;
  if (diffMs <= 0) {
    return locale === "id" ? "Segera reset" : "Resetting soon";
  }

  const totalMinutes = Math.max(1, Math.floor(diffMs / 60_000));
  const days = Math.floor(totalMinutes / (60 * 24));
  const hours = Math.floor((totalMinutes % (60 * 24)) / 60);
  const minutes = totalMinutes % 60;

  if (days > 0) {
    return locale === "id" ? `${days}h ${hours}j lagi` : `in ${days}d ${hours}h`;
  }
  if (hours > 0) {
    return locale === "id" ? `${hours}j ${minutes}m lagi` : `in ${hours}h ${minutes}m`;
  }
  return locale === "id" ? `${minutes}m lagi` : `in ${minutes}m`;
}
