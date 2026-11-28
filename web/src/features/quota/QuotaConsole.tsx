import {
  AlertCircle,
  ChevronDown,
  ChevronsDownUp,
  ChevronsUpDown,
  Clock,
  Gauge,
  RefreshCw,
  Search,
  Server,
  X,
} from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router";
import { ProviderLogo } from "@/features/providers/ProviderLogo";
import { extractErrorMessage } from "@/shared/api/client";
import { cn } from "@/shared/lib/utils";
import { type UILocale, useUIStore } from "@/shared/store/uiStore";
import { AnimatedNumber } from "@/shared/ui/AnimatedNumber";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  type QuotaLiveStatus,
  useProviderQuotas,
  useQuotaLiveSync,
  useRefreshAllQuotasMutation,
  useRefreshConnectionQuotaMutation,
} from "./hooks";
import {
  type BucketProgressTone,
  computeQuotaOverviewStats,
  filterProviderQuotas,
  formatResetCountdown,
  getBucketProgressTone,
  getQuotaStatusVariant,
  groupProviderQuotas,
  type ProviderConnectionQuota,
  type ProviderQuotaBucket,
  type ProviderQuotaGroup,
  type ProviderQuotaStatus,
  type QuotaStatusFilter,
} from "./schemas";

const CLOCK_TICK_INTERVAL_MS = 60_000;

interface QuotaCopy {
  title: string;
  accountsUnit: string;
  activeUnit: string;
  lowUnit: string;
  exhaustedUnit: string;
  groupAccountsCount: (count: number) => string;
  liveOn: string;
  liveConnecting: string;
  liveReconnecting: string;
  livePaused: string;
  toggleLiveAria: string;
  expandAll: string;
  collapseAll: string;
  refreshAll: string;
  refreshing: string;
  refreshOne: string;
  searchPlaceholder: string;
  clearSearch: string;
  allProviders: string;
  filterAll: string;
  colAccount: string;
  colRemaining: string;
  colReset: string;
  colStatus: string;
  avgLeftLabel: string;
  disabledTag: string;
  notActivatedHint: string;
  statusLabels: Record<ProviderQuotaStatus, string>;
  errorTitle: string;
  retryLabel: string;
  emptyTitle: string;
  emptyDescription: string;
  emptyAction: string;
  filteredEmptyTitle: string;
  resetFilters: string;
}

const QUOTA_COPY: Record<UILocale, QuotaCopy> = {
  en: {
    title: "Quota",
    accountsUnit: "accounts",
    activeUnit: "active",
    lowUnit: "low",
    exhaustedUnit: "exhausted",
    groupAccountsCount: (count) => `${count} ${count === 1 ? "account" : "accounts"}`,
    liveOn: "Live",
    liveConnecting: "Connecting",
    liveReconnecting: "Reconnecting",
    livePaused: "Paused",
    toggleLiveAria: "Toggle realtime quota sync",
    expandAll: "Expand",
    collapseAll: "Collapse",
    refreshAll: "Refresh",
    refreshing: "Syncing",
    refreshOne: "Refresh account quota",
    searchPlaceholder: "Search account or plan...",
    clearSearch: "Clear search",
    allProviders: "All providers",
    filterAll: "All",
    colAccount: "Account",
    colRemaining: "Remaining",
    colReset: "Reset",
    colStatus: "Status",
    avgLeftLabel: "avg left",
    disabledTag: "Off",
    notActivatedHint: "Free tier not activated in provider settings",
    statusLabels: {
      ok: "Active",
      low: "Low",
      exhausted: "Exhausted",
      not_activated: "Inactive",
      error: "Error",
    },
    errorTitle: "Failed to load quotas",
    retryLabel: "Retry",
    emptyTitle: "No quota-supported accounts yet",
    emptyDescription:
      "Add a supported provider connection (Token Harbor or CodeBuddy) in Providers to view remaining allowances here.",
    emptyAction: "Open Providers",
    filteredEmptyTitle: "No matching accounts",
    resetFilters: "Reset filters",
  },
  id: {
    title: "Kuota",
    accountsUnit: "akun",
    activeUnit: "aktif",
    lowUnit: "menipis",
    exhaustedUnit: "habis",
    groupAccountsCount: (count) => `${count} akun`,
    liveOn: "Live",
    liveConnecting: "Menghubungkan",
    liveReconnecting: "Menyambung ulang",
    livePaused: "Jeda",
    toggleLiveAria: "Aktifkan atau jeda sinkronisasi kuota realtime",
    expandAll: "Buka",
    collapseAll: "Tutup",
    refreshAll: "Refresh",
    refreshing: "Sinkron",
    refreshOne: "Muat ulang kuota akun ini",
    searchPlaceholder: "Cari nama akun atau paket...",
    clearSearch: "Hapus pencarian",
    allProviders: "Semua provider",
    filterAll: "Semua",
    colAccount: "Akun",
    colRemaining: "Sisa Limit",
    colReset: "Reset",
    colStatus: "Status",
    avgLeftLabel: "rata-rata",
    disabledTag: "Nonaktif",
    notActivatedHint: "Free tier belum diaktifkan di pengaturan provider",
    statusLabels: {
      ok: "Aktif",
      low: "Menipis",
      exhausted: "Habis",
      not_activated: "Belum Aktif",
      error: "Error",
    },
    errorTitle: "Gagal memuat kuota",
    retryLabel: "Coba lagi",
    emptyTitle: "Belum ada akun provider dengan dukungan kuota",
    emptyDescription:
      "Tambahkan koneksi provider yang didukung (Token Harbor atau CodeBuddy) di menu Provider untuk melihat sisa limit di sini.",
    emptyAction: "Buka Provider",
    filteredEmptyTitle: "Tidak ada akun yang cocok",
    resetFilters: "Reset filter",
  },
};

function formatPercent(value: number): string {
  if (!Number.isFinite(value)) {
    return "0";
  }
  const rounded = Math.round(value * 10) / 10;
  return Number.isInteger(rounded) ? String(rounded) : rounded.toFixed(1);
}

function formatPercentWithSymbol(value: number): string {
  return `${formatPercent(value)}%`;
}

function formatDateTime(iso: string | null | undefined, locale: UILocale): string | undefined {
  if (!iso) {
    return undefined;
  }
  const parsed = new Date(iso);
  if (Number.isNaN(parsed.getTime())) {
    return undefined;
  }
  return parsed.toLocaleString(locale === "id" ? "id-ID" : "en-US", {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function barFillClass(tone: BucketProgressTone): string {
  switch (tone) {
    case "ok":
      return "bg-[var(--status-ok)]";
    case "warn":
      return "bg-[var(--status-warn)]";
    case "error":
      return "bg-[var(--status-error)]";
  }
}

function barTextClass(tone: BucketProgressTone): string {
  switch (tone) {
    case "ok":
      return "text-[var(--status-ok)]";
    case "warn":
      return "text-[var(--status-warn)]";
    case "error":
      return "text-[var(--status-error)]";
  }
}

function liveIndicatorMeta(liveEnabled: boolean, status: QuotaLiveStatus, copy: QuotaCopy) {
  if (!liveEnabled) {
    return {
      label: copy.livePaused,
      dotClass: "bg-[var(--text-muted)]",
      pulse: false,
      active: false,
    };
  }
  switch (status) {
    case "live":
      return {
        label: copy.liveOn,
        dotClass: "bg-[var(--status-ok)]",
        pulse: true,
        active: true,
      };
    case "reconnecting":
      return {
        label: copy.liveReconnecting,
        dotClass: "bg-[var(--status-warn)]",
        pulse: true,
        active: true,
      };
    default:
      return {
        label: copy.liveConnecting,
        dotClass: "bg-[var(--accent-primary)]",
        pulse: true,
        active: true,
      };
  }
}

function CompactBucketBar({
  bucket,
  showLabel,
}: {
  bucket: ProviderQuotaBucket;
  showLabel: boolean;
}) {
  const tone = getBucketProgressTone(bucket.remaining_pct, bucket.exhausted);
  const remaining = Math.max(0, Math.min(100, bucket.remaining_pct));

  return (
    <div className="flex min-w-0 items-center gap-2.5">
      {showLabel && (
        <span className="w-16 shrink-0 truncate text-[11px] text-[var(--text-muted)]">
          {bucket.label}
        </span>
      )}
      <div
        role="progressbar"
        aria-label={`${bucket.label}: ${formatPercent(remaining)}%`}
        aria-valuenow={Math.round(remaining)}
        aria-valuemin={0}
        aria-valuemax={100}
        className="h-2 flex-1 overflow-hidden rounded-full bg-[var(--bg-surface-raised)] sm:h-1.5"
      >
        <div
          className={cn(
            "h-full rounded-full transition-[width,background-color] duration-500 ease-out motion-reduce:transition-none",
            barFillClass(tone),
          )}
          style={{ width: `${remaining}%` }}
        />
      </div>
      <span
        className={cn(
          "w-12 shrink-0 text-right font-mono text-xs font-semibold tabular-nums transition-colors duration-300",
          barTextClass(tone),
        )}
      >
        <AnimatedNumber value={remaining} format={formatPercentWithSymbol} durationMs={480} />
      </span>
    </div>
  );
}

function QuotaRow({
  item,
  locale,
  nowMs,
  copy,
  isRefreshing,
  onRefresh,
}: {
  item: ProviderConnectionQuota;
  locale: UILocale;
  nowMs: number;
  copy: QuotaCopy;
  isRefreshing: boolean;
  onRefresh: (connectionId: string) => void;
}) {
  const primaryResetIso = item.buckets[0]?.reset_at ?? item.buckets[0]?.period_ends_at ?? null;
  const resetCountdown = formatResetCountdown(primaryResetIso, locale, nowMs);
  const resetFull = formatDateTime(primaryResetIso, locale);
  const showBucketLabels = item.buckets.length > 1;

  return (
    <div className="grid grid-cols-1 gap-2 px-3.5 py-3 transition-colors duration-200 hover:bg-[var(--bg-canvas)]/60 sm:px-4 md:grid-cols-[minmax(160px,1.2fr)_minmax(180px,1.8fr)_110px_96px_36px] md:items-center md:gap-4 md:py-2.5">
      {/* Row top on mobile / Col 1 on md+ */}
      <div className="flex min-w-0 items-center justify-between gap-2 md:justify-start">
        <div className="flex min-w-0 items-center gap-1.5">
          <span className="truncate text-sm font-medium text-[var(--text-primary)]">
            {item.connection_name}
          </span>
          {item.on_plan && item.plan_name && (
            <Badge variant="default" className="shrink-0">
              {item.plan_name}
            </Badge>
          )}
          {!item.enabled && (
            <Badge variant="outline" className="shrink-0">
              {copy.disabledTag}
            </Badge>
          )}
        </div>

        {/* Mobile-only status + touch-friendly refresh */}
        <div className="flex shrink-0 items-center gap-1.5 md:hidden">
          <Badge
            variant={getQuotaStatusVariant(item.status)}
            className="transition-colors duration-300"
          >
            {copy.statusLabels[item.status]}
          </Badge>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onRefresh(item.connection_id)}
            disabled={isRefreshing}
            aria-label={copy.refreshOne}
            className="h-8 w-8 p-0"
          >
            <RefreshCw
              className={cn("h-3.5 w-3.5", isRefreshing && "animate-spin")}
              aria-hidden="true"
            />
          </Button>
        </div>
      </div>

      {/* Progress bar(s) or inline error */}
      <div className="min-w-0 space-y-1.5">
        {item.status === "error" && item.error ? (
          <p
            className="break-words text-xs text-[var(--status-error)] md:truncate"
            title={item.error}
          >
            {item.error}
          </p>
        ) : (
          <>
            {item.buckets.map((bucket) => (
              <CompactBucketBar key={bucket.id} bucket={bucket} showLabel={showBucketLabels} />
            ))}
            {item.fetcher_id === "tokenharbor" && !item.activated && (
              <p className="text-[11px] text-[var(--status-warn)]">{copy.notActivatedHint}</p>
            )}
          </>
        )}
      </div>

      {/* Reset countdown */}
      <div
        className="flex items-center justify-between gap-2 text-[11px] text-[var(--text-muted)] md:block md:text-xs"
        title={resetFull}
      >
        <span className="inline-flex items-center gap-1 md:hidden">
          <Clock className="h-3 w-3 shrink-0" aria-hidden="true" />
          <span>{copy.colReset}</span>
        </span>
        <span className="font-mono tabular-nums">{resetCountdown ?? "—"}</span>
      </div>

      {/* Desktop status */}
      <div className="hidden md:flex md:justify-start">
        <Badge
          variant={getQuotaStatusVariant(item.status)}
          className="transition-colors duration-300"
        >
          {copy.statusLabels[item.status]}
        </Badge>
      </div>

      {/* Desktop refresh */}
      <div className="hidden md:flex md:justify-end">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={() => onRefresh(item.connection_id)}
          disabled={isRefreshing}
          aria-label={copy.refreshOne}
          title={copy.refreshOne}
          className="h-7 w-7 p-0"
        >
          <RefreshCw
            className={cn("h-3.5 w-3.5", isRefreshing && "animate-spin")}
            aria-hidden="true"
          />
        </Button>
      </div>
    </div>
  );
}

function QuotaProviderAccordion({
  group,
  expanded,
  onToggle,
  locale,
  nowMs,
  copy,
  refreshingConnectionId,
  syncingIds,
  onRefreshOne,
}: {
  group: ProviderQuotaGroup;
  expanded: boolean;
  onToggle: () => void;
  locale: UILocale;
  nowMs: number;
  copy: QuotaCopy;
  refreshingConnectionId: string | null;
  syncingIds: ReadonlySet<string>;
  onRefreshOne: (connectionId: string) => void;
}) {
  const panelId = `quota-group-panel-${group.providerId}`;
  const avgTone =
    group.avgRemainingPct !== null
      ? getBucketProgressTone(group.avgRemainingPct, group.avgRemainingPct <= 0)
      : null;

  return (
    <Card className="overflow-hidden">
      <button
        type="button"
        onClick={onToggle}
        aria-expanded={expanded}
        aria-controls={panelId}
        className="flex w-full items-center justify-between gap-2.5 px-3.5 py-3 text-left transition-colors hover:bg-[var(--bg-canvas)]/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-[var(--accent-primary)] sm:gap-3 sm:px-4"
      >
        <div className="flex min-w-0 flex-1 items-center gap-2.5 sm:gap-3">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)]">
            <ProviderLogo
              providerId={group.providerId}
              name={group.providerName}
              isCustom={group.providerId !== group.fetcherId}
              size={20}
            />
          </div>

          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-1.5 sm:gap-2">
              <span className="truncate text-sm font-semibold text-[var(--text-primary)]">
                {group.providerName}
              </span>
              <Badge variant="secondary">{copy.groupAccountsCount(group.stats.total)}</Badge>
              {group.stats.low > 0 && (
                <Badge variant="warning" className="transition-colors duration-300">
                  {group.stats.low} {copy.statusLabels.low}
                </Badge>
              )}
              {group.stats.exhausted > 0 && (
                <Badge variant="destructive" className="transition-colors duration-300">
                  {group.stats.exhausted} {copy.statusLabels.exhausted}
                </Badge>
              )}
              {group.stats.error > 0 && (
                <Badge variant="destructive" className="transition-colors duration-300">
                  {group.stats.error} {copy.statusLabels.error}
                </Badge>
              )}
            </div>
          </div>
        </div>

        <div className="flex shrink-0 items-center gap-2 sm:gap-3">
          {group.avgRemainingPct !== null && avgTone && (
            <div className="flex items-center gap-2">
              <div className="hidden h-1.5 w-20 overflow-hidden rounded-full bg-[var(--bg-surface-raised)] sm:block">
                <div
                  className={cn(
                    "h-full rounded-full transition-[width,background-color] duration-500 ease-out motion-reduce:transition-none",
                    barFillClass(avgTone),
                  )}
                  style={{ width: `${group.avgRemainingPct}%` }}
                />
              </div>
              <span
                className={cn(
                  "font-mono text-xs font-semibold tabular-nums transition-colors duration-300",
                  barTextClass(avgTone),
                )}
              >
                <AnimatedNumber
                  value={group.avgRemainingPct}
                  format={formatPercentWithSymbol}
                  durationMs={480}
                />
              </span>
              <span className="hidden text-[11px] text-[var(--text-muted)] md:inline">
                {copy.avgLeftLabel}
              </span>
            </div>
          )}

          <ChevronDown
            className={cn(
              "h-4 w-4 text-[var(--text-muted)] transition-transform duration-200 ease-out motion-reduce:transition-none",
              expanded && "rotate-180",
            )}
            aria-hidden="true"
          />
        </div>
      </button>

      <div
        id={panelId}
        className={cn(
          "grid transition-[grid-template-rows,opacity] duration-200 ease-out motion-reduce:transition-none",
          expanded ? "grid-rows-[1fr] opacity-100" : "grid-rows-[0fr] opacity-0",
        )}
      >
        <div className="overflow-hidden">
          <div className="border-t border-[var(--border-subtle)]">
            <div className="hidden border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)]/70 px-4 py-2 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)] md:grid md:grid-cols-[minmax(160px,1.2fr)_minmax(180px,1.8fr)_110px_96px_36px] md:gap-4">
              <span>{copy.colAccount}</span>
              <span>{copy.colRemaining}</span>
              <span>{copy.colReset}</span>
              <span>{copy.colStatus}</span>
              <span />
            </div>
            <div className="divide-y divide-[var(--border-subtle)]">
              {group.items.map((item) => (
                <QuotaRow
                  key={item.connection_id}
                  item={item}
                  locale={locale}
                  nowMs={nowMs}
                  copy={copy}
                  isRefreshing={
                    refreshingConnectionId === item.connection_id ||
                    syncingIds.has(item.connection_id)
                  }
                  onRefresh={onRefreshOne}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}

export function QuotaConsole() {
  const locale = useUIStore((state) => state.locale);
  const copy = QUOTA_COPY[locale];

  const [liveEnabled, setLiveEnabled] = useState(true);
  const [nowMs, setNowMs] = useState(() => Date.now());
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<QuotaStatusFilter>("all");
  const [providerFilter, setProviderFilter] = useState("all");
  const [collapsedGroups, setCollapsedGroups] = useState<Record<string, boolean>>({});
  const [refreshingConnectionId, setRefreshingConnectionId] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  useEffect(() => {
    const timer = window.setInterval(() => {
      setNowMs(Date.now());
    }, CLOCK_TICK_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, []);

  const quotasQuery = useProviderQuotas({ liveEnabled });
  const refreshAllMutation = useRefreshAllQuotasMutation();
  const refreshOneMutation = useRefreshConnectionQuotaMutation();

  const quotas = useMemo(() => quotasQuery.data?.quotas ?? [], [quotasQuery.data]);
  const trackedConnectionIds = useMemo(() => quotas.map((item) => item.connection_id), [quotas]);
  const { liveStatus, syncingIds } = useQuotaLiveSync(liveEnabled, trackedConnectionIds);

  const stats = useMemo(() => computeQuotaOverviewStats(quotas), [quotas]);

  const distinctProviders = useMemo(() => {
    const map = new Map<string, string>();
    for (const item of quotas) {
      if (!map.has(item.provider_id)) {
        map.set(item.provider_id, item.provider_name);
      }
    }
    return Array.from(map.entries());
  }, [quotas]);

  const filteredQuotas = useMemo(
    () =>
      filterProviderQuotas(quotas, {
        query: searchQuery,
        status: statusFilter,
        providerId: providerFilter,
      }),
    [quotas, searchQuery, statusFilter, providerFilter],
  );

  const groupedQuotas = useMemo(() => groupProviderQuotas(filteredQuotas), [filteredQuotas]);

  const isSearching = searchQuery.trim().length > 0;
  const isGroupExpanded = (providerId: string) => isSearching || !collapsedGroups[providerId];
  const allGroupsExpanded =
    groupedQuotas.length > 0 && groupedQuotas.every((group) => isGroupExpanded(group.providerId));

  const toggleGroup = (providerId: string) => {
    setCollapsedGroups((prev) => ({
      ...prev,
      [providerId]: !prev[providerId],
    }));
  };

  const handleToggleAllGroups = () => {
    if (allGroupsExpanded) {
      const next: Record<string, boolean> = {};
      for (const group of groupedQuotas) {
        next[group.providerId] = true;
      }
      setCollapsedGroups(next);
    } else {
      setCollapsedGroups({});
    }
  };

  const handleRefreshAll = async () => {
    setActionError(null);
    try {
      await refreshAllMutation.mutateAsync(undefined);
    } catch (err) {
      setActionError(extractErrorMessage(err));
    }
  };

  const handleRefreshOne = async (connectionId: string) => {
    setActionError(null);
    setRefreshingConnectionId(connectionId);
    try {
      await refreshOneMutation.mutateAsync(connectionId);
    } catch (err) {
      setActionError(extractErrorMessage(err));
    } finally {
      setRefreshingConnectionId((current) => (current === connectionId ? null : current));
    }
  };

  const resetFilters = () => {
    setSearchQuery("");
    setStatusFilter("all");
    setProviderFilter("all");
  };

  const isRefreshingAll = refreshAllMutation.isPending || quotasQuery.isFetching;
  const liveMeta = liveIndicatorMeta(liveEnabled, liveStatus, copy);

  const statusTabs = useMemo(() => {
    const candidates: ReadonlyArray<{ id: QuotaStatusFilter; label: string; count: number }> = [
      { id: "all", label: copy.filterAll, count: stats.total },
      { id: "ok", label: copy.statusLabels.ok, count: stats.healthy },
      { id: "low", label: copy.statusLabels.low, count: stats.low },
      { id: "exhausted", label: copy.statusLabels.exhausted, count: stats.exhausted },
      { id: "not_activated", label: copy.statusLabels.not_activated, count: stats.notActivated },
      { id: "error", label: copy.statusLabels.error, count: stats.error },
    ];
    return candidates.filter((tab) => tab.id === "all" || tab.count > 0 || statusFilter === tab.id);
  }, [copy, stats, statusFilter]);

  return (
    <div className="min-w-0 space-y-3.5 sm:space-y-4">
      {/* Mobile-first header */}
      <header className="flex flex-col gap-2.5 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0 space-y-0.5 sm:flex sm:flex-wrap sm:items-baseline sm:gap-x-3 sm:gap-y-1 sm:space-y-0">
          <h1 className="flex items-center gap-2 text-lg font-semibold tracking-tight text-[var(--text-primary)]">
            <Gauge className="h-4 w-4 shrink-0 text-[var(--accent-text)]" aria-hidden="true" />
            <span>{copy.title}</span>
          </h1>
          {stats.total > 0 && (
            <p className="flex flex-wrap items-center gap-x-1.5 font-mono text-xs text-[var(--text-muted)]">
              <span>
                <AnimatedNumber value={stats.total} /> {copy.accountsUnit}
              </span>
              <span aria-hidden="true">·</span>
              <span className="text-[var(--status-ok)]">
                <AnimatedNumber value={stats.healthy} /> {copy.activeUnit}
              </span>
              {stats.low > 0 && (
                <>
                  <span aria-hidden="true">·</span>
                  <span className="text-[var(--status-warn)]">
                    <AnimatedNumber value={stats.low} /> {copy.lowUnit}
                  </span>
                </>
              )}
              {stats.exhausted > 0 && (
                <>
                  <span aria-hidden="true">·</span>
                  <span className="text-[var(--status-error)]">
                    <AnimatedNumber value={stats.exhausted} /> {copy.exhaustedUnit}
                  </span>
                </>
              )}
            </p>
          )}
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => setLiveEnabled((prev) => !prev)}
            aria-pressed={liveEnabled}
            aria-label={copy.toggleLiveAria}
            className={cn(
              "inline-flex h-9 items-center gap-1.5 rounded-md border px-2.5 text-xs font-medium transition-colors sm:h-8",
              liveMeta.active
                ? "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]"
                : "border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-muted)] hover:text-[var(--text-primary)]",
            )}
          >
            <span className="relative flex h-2 w-2 shrink-0" aria-hidden="true">
              {liveMeta.pulse && (
                <span
                  className={cn(
                    "absolute inline-flex h-full w-full animate-ping rounded-full opacity-70 motion-reduce:animate-none",
                    liveMeta.dotClass,
                  )}
                />
              )}
              <span
                className={cn("relative inline-flex h-2 w-2 rounded-full", liveMeta.dotClass)}
              />
            </span>
            <span>{liveMeta.label}</span>
          </button>

          {groupedQuotas.length > 0 && !isSearching && (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={handleToggleAllGroups}
              className="h-9 px-2.5 text-xs text-[var(--text-muted)] sm:h-8"
            >
              {allGroupsExpanded ? (
                <>
                  <ChevronsDownUp className="mr-1.5 h-3.5 w-3.5" aria-hidden="true" />
                  {copy.collapseAll}
                </>
              ) : (
                <>
                  <ChevronsUpDown className="mr-1.5 h-3.5 w-3.5" aria-hidden="true" />
                  {copy.expandAll}
                </>
              )}
            </Button>
          )}

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefreshAll}
            disabled={isRefreshingAll}
            className="h-9 sm:h-8"
          >
            <RefreshCw
              className={cn("mr-1.5 h-3.5 w-3.5", isRefreshingAll && "animate-spin")}
              aria-hidden="true"
            />
            {isRefreshingAll ? copy.refreshing : copy.refreshAll}
          </Button>
        </div>
      </header>

      {quotasQuery.isError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" aria-hidden="true" />
          <AlertTitle>{copy.errorTitle}</AlertTitle>
          <AlertDescription className="flex flex-wrap items-center justify-between gap-2">
            <span>{extractErrorMessage(quotasQuery.error)}</span>
            <Button type="button" variant="outline" size="sm" onClick={() => quotasQuery.refetch()}>
              {copy.retryLabel}
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {actionError && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" aria-hidden="true" />
          <AlertTitle>{copy.errorTitle}</AlertTitle>
          <AlertDescription className="flex flex-wrap items-center justify-between gap-2">
            <span>{actionError}</span>
            <Button type="button" variant="ghost" size="sm" onClick={() => setActionError(null)}>
              <X className="h-3.5 w-3.5" aria-hidden="true" />
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Mobile-first search + horizontally scrollable status tabs */}
      {quotas.length > 0 && (
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div className="relative w-full sm:max-w-xs">
            <Search
              className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
              aria-hidden="true"
            />
            <Input
              value={searchQuery}
              onChange={(event) => setSearchQuery(event.target.value)}
              placeholder={copy.searchPlaceholder}
              aria-label={copy.searchPlaceholder}
              className="h-9 pl-8 pr-7 text-xs sm:h-8"
            />
            {searchQuery && (
              <button
                type="button"
                onClick={() => setSearchQuery("")}
                aria-label={copy.clearSearch}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-[var(--text-muted)] hover:text-[var(--text-primary)]"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
            )}
          </div>

          <div className="flex min-w-0 items-center gap-1.5 overflow-x-auto pb-0.5 sm:flex-wrap sm:overflow-visible sm:pb-0">
            {distinctProviders.length > 1 && (
              <select
                value={providerFilter}
                onChange={(event) => setProviderFilter(event.target.value)}
                aria-label={copy.allProviders}
                className="h-9 shrink-0 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 text-xs text-[var(--text-primary)] focus:outline-none focus:ring-1 focus:ring-[var(--accent-primary)] sm:h-8"
              >
                <option value="all">{copy.allProviders}</option>
                {distinctProviders.map(([id, name]) => (
                  <option key={id} value={id}>
                    {name}
                  </option>
                ))}
              </select>
            )}

            {statusTabs.map((tab) => {
              const active = statusFilter === tab.id;
              return (
                <button
                  key={tab.id}
                  type="button"
                  onClick={() => setStatusFilter(tab.id)}
                  className={cn(
                    "inline-flex h-9 shrink-0 items-center gap-1.5 rounded-md border px-2.5 text-xs font-medium transition-colors sm:h-8",
                    active
                      ? "border-[var(--accent-primary)] bg-[var(--bg-surface-raised)] text-[var(--accent-text)]"
                      : "border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-muted)] hover:text-[var(--text-primary)]",
                  )}
                >
                  <span>{tab.label}</span>
                  <span className="font-mono text-[11px] opacity-75">{tab.count}</span>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {quotasQuery.isPending ? (
        <Card>
          <div className="divide-y divide-[var(--border-subtle)]">
            {[0, 1, 2, 3].map((idx) => (
              <div
                key={idx}
                className="flex items-center justify-between gap-4 px-3.5 py-3.5 sm:px-4"
              >
                <Skeleton className="h-4 w-28 sm:w-36" />
                <Skeleton className="h-2 flex-1 max-w-md" />
                <Skeleton className="h-4 w-12 sm:w-16" />
              </div>
            ))}
          </div>
        </Card>
      ) : quotas.length === 0 && !quotasQuery.isError ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center gap-3 px-4 py-10 text-center">
            <Gauge className="h-5 w-5 text-[var(--text-muted)]" aria-hidden="true" />
            <div className="space-y-1">
              <h2 className="text-sm font-semibold text-[var(--text-primary)]">
                {copy.emptyTitle}
              </h2>
              <p className="max-w-sm text-xs text-[var(--text-muted)]">{copy.emptyDescription}</p>
            </div>
            <Link
              to="/dashboard/providers/"
              className="inline-flex h-9 items-center gap-1.5 rounded-md border border-[var(--accent-primary)] bg-[var(--bg-surface-raised)] px-3.5 text-xs font-medium text-[var(--accent-text)] transition-colors hover:opacity-90 sm:h-8"
            >
              <Server className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{copy.emptyAction}</span>
            </Link>
          </CardContent>
        </Card>
      ) : groupedQuotas.length === 0 ? (
        <Card>
          <CardContent className="flex flex-wrap items-center justify-between gap-3 px-4 py-6">
            <span className="text-sm text-[var(--text-muted)]">{copy.filteredEmptyTitle}</span>
            <Button type="button" variant="outline" size="sm" onClick={resetFilters}>
              {copy.resetFilters}
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {groupedQuotas.map((group) => (
            <QuotaProviderAccordion
              key={group.providerId}
              group={group}
              expanded={isGroupExpanded(group.providerId)}
              onToggle={() => toggleGroup(group.providerId)}
              locale={locale}
              nowMs={nowMs}
              copy={copy}
              refreshingConnectionId={refreshingConnectionId}
              syncingIds={syncingIds}
              onRefreshOne={handleRefreshOne}
            />
          ))}
        </div>
      )}
    </div>
  );
}
