"use client";

import {
  AlertTriangle,
  ArrowRight,
  CheckCircle2,
  ChevronDown,
  FilterX,
  Globe,
  KeyRound,
  Layers,
  PauseCircle,
  Play,
  Plus,
  RefreshCw,
  Search,
  Server,
  ShieldCheck,
  Terminal,
  X,
} from "lucide-react";
import { type ReactNode, useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import {
  buildProviderEntries,
  computeFleetTelemetry,
  countCategoryBuckets,
  countStatusBuckets,
  matchesProviderSearch,
  matchesStatusFilter,
  type ProviderCategory,
  type ProviderCategoryTab,
  type ProviderEntry,
  type ProviderStatusFilter,
  sortProviderEntries,
} from "./catalogUtils";
import {
  useDeleteProviderConnectionMutation,
  useDeleteProviderNodeMutation,
  useProviderConnections,
  useProvidersCatalog,
  useTestProviderConnectionMutation,
  useUpdateProviderConnectionMutation,
} from "./hooks";
import { ProviderLogo } from "./ProviderLogo";
import {
  AddApiKeyModal,
  AddCustomNodeModal,
  OAuthWizardModal,
  ProviderConnectionsDrawer,
} from "./ProvidersModals";
import type {
  ProviderConnectionSummary,
  ProviderNode,
  ProviderProtocol,
  TestProviderConnectionResponse,
} from "./schemas";

interface BatchTestSummary {
  label: string;
  total: number;
  passed: number;
  failed: number;
  results: Array<{
    connectionId: string;
    connectionName: string;
    providerId: string;
    ok: boolean;
    statusCode: number;
    latencyMs: number;
    error?: string;
  }>;
}

const KPI_SKELETON_KEYS = ["kpi-1", "kpi-2", "kpi-3", "kpi-4"] as const;
const CARD_SKELETON_KEYS = [
  "card-1",
  "card-2",
  "card-3",
  "card-4",
  "card-5",
  "card-6",
  "card-7",
  "card-8",
] as const;

function ProvidersLoadingSkeleton() {
  return (
    <div className="flex min-w-0 flex-col gap-5 animate-pulse">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 rounded-xl bg-[var(--color-surface-2)]" />
          <div className="space-y-1.5">
            <div className="h-5 w-32 rounded bg-[var(--color-surface-2)]" />
            <div className="h-3 w-64 rounded bg-[var(--color-surface-2)]" />
          </div>
        </div>
        <div className="flex items-center gap-2">
          <div className="h-8 w-24 rounded-lg bg-[var(--color-surface-2)]" />
          <div className="h-8 w-28 rounded-lg bg-[var(--color-surface-2)]" />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-2.5 lg:grid-cols-4">
        {KPI_SKELETON_KEYS.map((key) => (
          <div
            key={key}
            className="h-[72px] rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)]"
          />
        ))}
      </div>

      <div className="h-24 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)]" />

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {CARD_SKELETON_KEYS.map((key) => (
          <div
            key={key}
            className="h-[118px] rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)]"
          />
        ))}
      </div>
    </div>
  );
}

function SectionHeader({
  icon,
  title,
  count,
  action,
}: {
  icon: ReactNode;
  title: string;
  count: number;
  action?: ReactNode;
}) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <div className="flex items-center gap-2">
        <span className="text-[var(--color-primary)]">{icon}</span>
        <h3 className="text-sm font-semibold tracking-tight text-[var(--color-text-main)]">
          {title}
        </h3>
        <span className="rounded-full border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
          {count}
        </span>
      </div>
      {action}
    </div>
  );
}

function ProviderCard({
  entry,
  isId,
  onSelect,
  onQuickTest,
  onToggleProvider,
  isTesting,
  isToggling,
  disableQuickTest,
}: {
  entry: ProviderEntry;
  isId: boolean;
  onSelect: (entry: ProviderEntry) => void;
  onQuickTest: (entry: ProviderEntry) => void;
  onToggleProvider: (entry: ProviderEntry, nextActive: boolean) => void;
  isTesting: boolean;
  isToggling: boolean;
  disableQuickTest: boolean;
}) {
  const { stats } = entry;
  const hasConnections = stats.total > 0;
  const hasError = !stats.allDisabled && stats.error > 0;
  const isActiveHealthy = stats.isHealthy;
  const canQuickTest = hasConnections && !stats.allDisabled;
  const typeLabel =
    entry.category === "oauth"
      ? "OAUTH"
      : entry.category === "free"
        ? entry.sourceCategory === "webCookie"
          ? "WEB SESSION"
          : entry.noAuth
            ? "FREE · NO AUTH"
            : "FREE TIER"
        : entry.category === "custom"
          ? entry.protocol.toUpperCase()
          : `API KEY · ${entry.protocol.toUpperCase()}`;

  return (
    <div
      className={cn(
        "group relative flex flex-col justify-between rounded-xl border bg-[var(--color-surface)] p-4 transition-all",
        hasError
          ? "border-[var(--status-error)]/40 hover:border-[var(--status-error)]/70"
          : isActiveHealthy
            ? "border-[var(--color-border)] hover:border-[var(--color-primary)]/45"
            : "border-[var(--color-border)]/80 hover:border-[var(--color-primary)]/40",
        "hover:bg-[var(--color-surface-2)]/35 hover:shadow-[0_4px_16px_-8px_rgba(11,13,18,0.16)]",
      )}
    >
      {/* Top Row: Monogram + Name + Quick Actions */}
      <div className="flex items-start justify-between gap-2">
        <button
          type="button"
          onClick={() => onSelect(entry)}
          className="flex min-w-0 flex-1 cursor-pointer items-center gap-3 text-left focus:outline-none"
        >
          <div
            className={cn(
              "flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl border p-1.5 transition-colors",
              hasError
                ? "border-[var(--status-error)]/30 bg-[var(--status-error-bg)] text-[var(--status-error)]"
                : isActiveHealthy
                  ? "border-[var(--color-primary)]/25 bg-[var(--color-surface-2)]/70 text-[var(--color-primary)]"
                  : "border-[var(--color-border)] bg-[var(--color-surface-2)]/80 text-[var(--color-text-muted)] group-hover:border-[var(--color-primary)]/35",
            )}
          >
            <ProviderLogo
              providerId={entry.id}
              name={entry.name}
              protocol={entry.protocol}
              isCustom={entry.isCustom}
              size={24}
            />
          </div>
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold text-[var(--color-text-main)] transition-colors group-hover:text-[var(--color-primary)]">
              {entry.name}
            </p>
            <span className="mt-0.5 block truncate font-mono text-[10px] uppercase tracking-wider text-[var(--color-text-muted)]">
              {typeLabel}
            </span>
          </div>
        </button>

        {hasConnections ? (
          <div className="relative z-10 flex shrink-0 items-center gap-1.5">
            {canQuickTest ? (
              <button
                type="button"
                disabled={disableQuickTest}
                onClick={(e) => {
                  e.stopPropagation();
                  onQuickTest(entry);
                }}
                title={isId ? `Tes koneksi ${entry.name}` : `Test ${entry.name} connections`}
                aria-label={`Test ${entry.name}`}
                className={cn(
                  "flex h-7 w-7 cursor-pointer items-center justify-center rounded-lg border transition-colors",
                  isTesting
                    ? "border-[var(--color-primary)]/40 bg-[var(--color-primary)]/15 text-[var(--color-primary)]"
                    : "border-[var(--color-border)] bg-[var(--color-surface-2)]/80 text-[var(--color-text-muted)] hover:border-[var(--color-primary)]/45 hover:text-[var(--color-primary)]",
                  disableQuickTest && !isTesting && "cursor-not-allowed opacity-45",
                )}
              >
                {isTesting ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                ) : (
                  <Play className="h-3.5 w-3.5" aria-hidden="true" />
                )}
              </button>
            ) : null}

            <button
              type="button"
              role="switch"
              aria-checked={!stats.allDisabled}
              disabled={isToggling}
              onClick={(e) => {
                e.stopPropagation();
                onToggleProvider(entry, stats.allDisabled);
              }}
              title={
                stats.allDisabled
                  ? isId
                    ? "Aktifkan provider"
                    : "Enable provider"
                  : isId
                    ? "Nonaktifkan provider"
                    : "Disable provider"
              }
              aria-label={`${stats.allDisabled ? "Enable" : "Disable"} ${entry.name}`}
              className={cn(
                "relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border transition-colors disabled:opacity-50",
                !stats.allDisabled
                  ? "border-[var(--color-primary)] bg-[var(--color-primary)]"
                  : "border-[var(--color-border)] bg-[var(--color-surface-3)]",
              )}
            >
              <span
                className={cn(
                  "inline-block h-3.5 w-3.5 rounded-full bg-[var(--color-on-primary)] shadow-xs transition-transform",
                  !stats.allDisabled ? "translate-x-4" : "translate-x-0.5",
                )}
              />
            </button>
          </div>
        ) : null}
      </div>

      {/* Bottom Row: Status Badges + Hover Arrow */}
      <button
        type="button"
        onClick={() => onSelect(entry)}
        className="mt-3.5 flex w-full cursor-pointer items-center justify-between gap-2 border-t border-[var(--color-border-subtle)] pt-2.5 text-left text-xs focus:outline-none"
      >
        <div className="flex min-w-0 flex-wrap items-center gap-1.5">
          {stats.allDisabled ? (
            <Badge variant="outline">
              <PauseCircle className="mr-1 h-3 w-3" aria-hidden="true" />
              {isId ? "Nonaktif" : "Disabled"}
            </Badge>
          ) : hasConnections ? (
            <>
              {stats.connected > 0 ? (
                <Badge variant="success">
                  ● {stats.connected} {isId ? "Aktif" : "Active"}
                </Badge>
              ) : null}
              {stats.error > 0 ? (
                <Badge variant="destructive">
                  ● {stats.error} {isId ? "Error" : "Error"}
                </Badge>
              ) : null}
            </>
          ) : entry.isCustom ? (
            <span className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--color-primary)]">
              <KeyRound className="h-3 w-3" aria-hidden="true" />
              {isId ? "Tambah API Key" : "Add API Key"}
            </span>
          ) : (
            <span className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--color-text-muted)] transition-colors group-hover:text-[var(--color-primary)]">
              <Plus className="h-3 w-3" aria-hidden="true" />
              {isId ? "Atur Koneksi" : "Setup"}
            </span>
          )}
        </div>

        <ArrowRight
          className="h-3.5 w-3.5 shrink-0 text-[var(--color-text-muted)] opacity-45 transition-all group-hover:translate-x-0.5 group-hover:text-[var(--color-primary)] group-hover:opacity-100"
          aria-hidden="true"
        />
      </button>
    </div>
  );
}

export function ProvidersConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const catalogQuery = useProvidersCatalog();
  const connectionsQuery = useProviderConnections();

  const updateConnMutation = useUpdateProviderConnectionMutation();
  const deleteConnMutation = useDeleteProviderConnectionMutation();
  const testConnMutation = useTestProviderConnectionMutation();
  const deleteNodeMutation = useDeleteProviderNodeMutation();

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<ProviderStatusFilter>("all");
  const [categoryTab, setCategoryTab] = useState<ProviderCategoryTab>("all");
  const [showAllApiKey, setShowAllApiKey] = useState(false);

  const [addMenuOpen, setAddMenuOpen] = useState(false);
  const addMenuRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const [apiKeyModal, setApiKeyModal] = useState<{
    open: boolean;
    providerId: string | null;
  }>({
    open: false,
    providerId: null,
  });
  const [oauthModal, setOauthModal] = useState<{
    open: boolean;
    providerId: string | null;
  }>({
    open: false,
    providerId: null,
  });
  const [customNodeModal, setCustomNodeModal] = useState<{
    open: boolean;
    protocol: ProviderProtocol;
  }>({ open: false, protocol: "openai-chat" });
  const [selectedProviderId, setSelectedProviderId] = useState<string | null>(null);

  const [feedbackBanner, setFeedbackBanner] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [busyConnectionId, setBusyConnectionId] = useState<string | null>(null);
  const [togglingProviderId, setTogglingProviderId] = useState<string | null>(null);
  const [testingScope, setTestingScope] = useState<string | null>(null);
  const [liveTestResults, setLiveTestResults] = useState<
    Record<string, TestProviderConnectionResponse>
  >({});
  const [batchSummary, setBatchSummary] = useState<BatchTestSummary | null>(null);

  // Close Add dropdown on outside click or Escape
  useEffect(() => {
    if (!addMenuOpen) return;
    const onMouseDown = (e: MouseEvent) => {
      if (addMenuRef.current && !addMenuRef.current.contains(e.target as Node)) {
        setAddMenuOpen(false);
      }
    };
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setAddMenuOpen(false);
    };
    document.addEventListener("mousedown", onMouseDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onMouseDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [addMenuOpen]);

  // "/" shortcut to focus provider search (matching main)
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key !== "/" || e.metaKey || e.ctrlKey || e.altKey) return;
      const target = e.target;
      if (
        target instanceof HTMLElement &&
        (target.tagName === "INPUT" ||
          target.tagName === "TEXTAREA" ||
          target.tagName === "SELECT" ||
          target.isContentEditable)
      ) {
        return;
      }
      if (document.querySelector('[role="dialog"][aria-modal="true"]')) {
        return;
      }
      e.preventDefault();
      searchInputRef.current?.focus();
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, []);

  const adapters = catalogQuery.data?.adapters ?? [];
  const customNodes = catalogQuery.data?.custom_nodes ?? [];
  const connections = connectionsQuery.data ?? [];

  const allEntries = useMemo(
    () => buildProviderEntries(adapters, customNodes, connections),
    [adapters, customNodes, connections],
  );

  const telemetry = useMemo(
    () => computeFleetTelemetry(allEntries, connections),
    [allEntries, connections],
  );

  const statusCounts = useMemo(() => countStatusBuckets(allEntries), [allEntries]);
  const categoryCounts = useMemo(() => countCategoryBuckets(allEntries), [allEntries]);

  const filteredEntries = useMemo(() => {
    const matched = allEntries.filter((entry) => {
      if (categoryTab !== "all" && entry.category !== categoryTab) return false;
      if (!matchesStatusFilter(entry, statusFilter)) return false;
      if (!matchesProviderSearch(entry, search)) return false;
      return true;
    });
    return sortProviderEntries(matched);
  }, [allEntries, categoryTab, statusFilter, search]);

  const selectedEntry = useMemo(
    () => allEntries.find((e) => e.id === selectedProviderId) ?? null,
    [allEntries, selectedProviderId],
  );

  const hasActiveFilters =
    search.trim().length > 0 || statusFilter !== "all" || categoryTab !== "all";
  const isRefreshing = catalogQuery.isFetching || connectionsQuery.isFetching;

  const resetFilters = () => {
    setSearch("");
    setStatusFilter("all");
    setCategoryTab("all");
  };

  const handleRefresh = () => {
    void Promise.all([catalogQuery.refetch(), connectionsQuery.refetch()]);
  };

  const runBatchTests = async (
    scopeKey: string,
    label: string,
    targetConnections: ProviderConnectionSummary[],
  ) => {
    const enabledTargets = targetConnections.filter((c) => c.enabled);
    if (enabledTargets.length === 0 || testingScope !== null) {
      return;
    }
    setActionError(null);
    setTestingScope(scopeKey);
    setBatchSummary(null);

    const results: BatchTestSummary["results"] = [];
    for (const conn of enabledTargets) {
      try {
        const res = await testConnMutation.mutateAsync(conn.id);
        setLiveTestResults((prev) => ({ ...prev, [conn.id]: res }));
        results.push({
          connectionId: conn.id,
          connectionName: conn.name,
          providerId: conn.provider_id,
          ok: res.ok,
          statusCode: res.status_code,
          latencyMs: res.latency_ms,
          error: res.error || undefined,
        });
      } catch (err) {
        results.push({
          connectionId: conn.id,
          connectionName: conn.name,
          providerId: conn.provider_id,
          ok: false,
          statusCode: 0,
          latencyMs: 0,
          error: err instanceof Error ? err.message : "Test failed",
        });
      }
    }

    const passed = results.filter((r) => r.ok).length;
    setBatchSummary({
      label,
      total: results.length,
      passed,
      failed: results.length - passed,
      results,
    });
    setTestingScope(null);
  };

  const handleToggleProvider = async (entry: ProviderEntry, nextActive: boolean) => {
    if (togglingProviderId !== null || entry.connections.length === 0) return;
    setActionError(null);
    setTogglingProviderId(entry.id);
    try {
      await Promise.all(
        entry.connections.map((conn) =>
          updateConnMutation.mutateAsync({
            id: conn.id,
            body: { enabled: nextActive },
          }),
        ),
      );
      setFeedbackBanner(
        nextActive
          ? isId
            ? `${entry.name} diaktifkan.`
            : `Enabled ${entry.name}.`
          : isId
            ? `${entry.name} dinonaktifkan.`
            : `Disabled ${entry.name}.`,
      );
    } catch (err) {
      setActionError(err instanceof Error ? err.message : "Failed to update provider state.");
    } finally {
      setTogglingProviderId(null);
    }
  };

  const handleTestSingleConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    testConnMutation.mutate(conn.id, {
      onSuccess: (res) => {
        setLiveTestResults((prev) => ({ ...prev, [conn.id]: res }));
      },
      onError: (err) => setActionError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const handleToggleSingleConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    updateConnMutation.mutate(
      { id: conn.id, body: { enabled: !conn.enabled } },
      {
        onError: (err) => setActionError(err.message),
        onSettled: () => setBusyConnectionId(null),
      },
    );
  };

  const handleDeleteSingleConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) return;
    setActionError(null);
    setBusyConnectionId(conn.id);
    deleteConnMutation.mutate(conn.id, {
      onSuccess: () => {
        setFeedbackBanner(
          isId ? `Koneksi "${conn.name}" dihapus.` : `Deleted connection "${conn.name}".`,
        );
      },
      onError: (err) => setActionError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const handleDeleteCustomNode = (node: ProviderNode) => {
    setActionError(null);
    deleteNodeMutation.mutate(node.id, {
      onSuccess: () => {
        setFeedbackBanner(
          isId
            ? `Custom endpoint "${node.name}" dihapus.`
            : `Removed custom endpoint "${node.name}".`,
        );
      },
      onError: (err) => setActionError(err.message),
    });
  };

  const handleCardSelect = (entry: ProviderEntry) => {
    if (entry.connections.length > 0 || entry.isCustom) {
      setSelectedProviderId(entry.id);
      return;
    }
    if (entry.isNativeOAuth) {
      setOauthModal({ open: true, providerId: entry.id });
    } else if (entry.authType === "oauth") {
      setOauthModal({ open: true, providerId: entry.id });
    } else {
      setApiKeyModal({ open: true, providerId: entry.id });
    }
  };

  const handleAddConnectionForEntry = (entry: ProviderEntry) => {
    setSelectedProviderId(null);
    if (entry.isNativeOAuth) {
      setOauthModal({ open: true, providerId: entry.id });
    } else {
      setApiKeyModal({ open: true, providerId: entry.id });
    }
  };

  if (catalogQuery.isPending && connectionsQuery.isPending) {
    return <ProvidersLoadingSkeleton />;
  }

  const toggleStatusKpi = (target: ProviderStatusFilter) => {
    setStatusFilter((prev) => (prev === target ? "all" : target));
  };

  const statusPills: Array<{ value: ProviderStatusFilter; label: string }> = [
    { value: "all", label: isId ? "Semua" : "All" },
    { value: "active", label: isId ? "Aktif" : "Active" },
    { value: "error", label: "Error" },
    { value: "inactive", label: isId ? "Nonaktif" : "Disabled" },
    { value: "unconfigured", label: isId ? "Belum Diatur" : "Unconfigured" },
  ];

  const categoryTabs: Array<{
    value: ProviderCategoryTab;
    label: string;
    icon: ReactNode;
  }> = [
    {
      value: "all",
      label: isId ? "Semua" : "All",
      icon: <Layers className="h-3.5 w-3.5" />,
    },
    {
      value: "oauth",
      label: isId ? "OAuth" : "OAuth",
      icon: <ShieldCheck className="h-3.5 w-3.5" />,
    },
    {
      value: "apikey",
      label: isId ? "API Key" : "API Key",
      icon: <KeyRound className="h-3.5 w-3.5" />,
    },
    {
      value: "free",
      label: isId ? "Gratis & Web" : "Free & Web",
      icon: <Globe className="h-3.5 w-3.5" />,
    },
    {
      value: "custom",
      label: isId ? "Custom Endpoint" : "Custom Endpoints",
      icon: <Terminal className="h-3.5 w-3.5" />,
    },
  ];

  const allSections: Array<{
    key: ProviderCategory;
    title: string;
    icon: ReactNode;
    entries: ProviderEntry[];
  }> = [
    {
      key: "custom",
      title: isId ? "Custom Endpoints" : "Custom Endpoints",
      icon: <Terminal className="h-4 w-4" aria-hidden="true" />,
      entries: filteredEntries.filter((e) => e.category === "custom"),
    },
    {
      key: "oauth",
      title: isId ? "Provider OAuth" : "OAuth Providers",
      icon: <ShieldCheck className="h-4 w-4" aria-hidden="true" />,
      entries: filteredEntries.filter((e) => e.category === "oauth"),
    },
    {
      key: "apikey",
      title: isId ? "Provider API Key" : "API Key Providers",
      icon: <KeyRound className="h-4 w-4" aria-hidden="true" />,
      entries: filteredEntries.filter((e) => e.category === "apikey"),
    },
    {
      key: "free",
      title: isId ? "Free Tier & Web Session" : "Free Tier & Web Session Providers",
      icon: <Globe className="h-4 w-4" aria-hidden="true" />,
      entries: filteredEntries.filter((e) => e.category === "free"),
    },
  ];
  const sectionsToRender = allSections.filter((sec) => sec.entries.length > 0);

  return (
    <div className="flex min-w-0 flex-col gap-5">
      {/* Studio Header Band (main ProvidersHeader parity) */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-[var(--color-primary)]/20 bg-[var(--color-primary)]/10 text-[var(--color-primary)]">
            <Server className="h-4 w-4" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h2 className="truncate text-lg font-semibold tracking-tight text-[var(--color-text-main)] sm:text-xl">
              {isId ? "Provider AI" : "Providers"}
            </h2>
            <p className="truncate text-xs text-[var(--color-text-muted)]">
              {isId
                ? "Pantau kesehatan upstream, uji koneksi, dan kelola kredensial provider AI."
                : "Monitor upstream health, test connections, and configure AI providers."}
            </p>
          </div>
        </div>

        <div className="flex shrink-0 flex-wrap items-center gap-2">
          <Button
            type="button"
            size="sm"
            variant="secondary"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            <RefreshCw className={cn("h-3.5 w-3.5", isRefreshing && "animate-spin")} />
            <span>{isId ? "Segarkan" : "Refresh"}</span>
          </Button>

          <Button
            type="button"
            size="sm"
            variant="secondary"
            disabled={telemetry.enabledConnections === 0 || testingScope !== null}
            onClick={() =>
              runBatchTests("all", isId ? "Semua Provider" : "All Providers", connections)
            }
          >
            {testingScope === "all" ? (
              <RefreshCw className="h-3.5 w-3.5 animate-spin" />
            ) : (
              <Play className="h-3.5 w-3.5" />
            )}
            <span>
              {testingScope === "all"
                ? isId
                  ? "Menguji..."
                  : "Testing..."
                : isId
                  ? "Tes Semua"
                  : "Test All"}
            </span>
          </Button>

          {/* Add Provider Dropdown */}
          <div className="relative" ref={addMenuRef}>
            <Button
              type="button"
              size="sm"
              aria-expanded={addMenuOpen}
              aria-haspopup="menu"
              onClick={() => setAddMenuOpen((v) => !v)}
            >
              <Plus className="h-3.5 w-3.5" />
              <span>{isId ? "Tambah Provider" : "Add Provider"}</span>
              <ChevronDown
                className={cn("h-3.5 w-3.5 transition-transform", addMenuOpen && "rotate-180")}
              />
            </Button>

            {addMenuOpen ? (
              <div
                role="menu"
                aria-label={isId ? "Menu tambah provider" : "Add provider menu"}
                className="absolute right-0 top-full z-40 mt-1.5 w-64 max-w-[calc(100vw-2rem)] rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-1.5 shadow-xl"
              >
                <button
                  type="button"
                  role="menuitem"
                  onClick={() => {
                    setAddMenuOpen(false);
                    setApiKeyModal({ open: true, providerId: null });
                  }}
                  className="flex w-full cursor-pointer items-start gap-2.5 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <KeyRound className="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-primary)]" />
                  <span className="min-w-0">
                    <span className="block text-xs font-semibold text-[var(--color-text-main)]">
                      {isId ? "Koneksi API Key" : "API Key Connection"}
                    </span>
                    <span className="block truncate text-[11px] text-[var(--color-text-muted)]">
                      OpenAI, Anthropic, Gemini, DeepSeek, Groq...
                    </span>
                  </span>
                </button>

                <button
                  type="button"
                  role="menuitem"
                  onClick={() => {
                    setAddMenuOpen(false);
                    setOauthModal({ open: true, providerId: null });
                  }}
                  className="flex w-full cursor-pointer items-start gap-2.5 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]" />
                  <span className="min-w-0">
                    <span className="block text-xs font-semibold text-[var(--color-text-main)]">
                      {isId ? "Akun OAuth (PKCE)" : "OAuth Account (PKCE)"}
                    </span>
                    <span className="block truncate text-[11px] text-[var(--color-text-muted)]">
                      Claude, Codex, Gemini OAuth
                    </span>
                  </span>
                </button>

                <div className="my-1 border-t border-[var(--color-border-subtle)]" />

                <button
                  type="button"
                  role="menuitem"
                  onClick={() => {
                    setAddMenuOpen(false);
                    setCustomNodeModal({ open: true, protocol: "openai-chat" });
                  }}
                  className="flex w-full cursor-pointer items-start gap-2.5 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <Terminal className="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-text-muted)]" />
                  <span className="min-w-0">
                    <span className="block text-xs font-semibold text-[var(--color-text-main)]">
                      {isId ? "OpenAI Compatible" : "OpenAI Compatible"}
                    </span>
                    <span className="block truncate text-[11px] text-[var(--color-text-muted)]">
                      vLLM, Ollama, LiteLLM, OpenRouter...
                    </span>
                  </span>
                </button>

                <button
                  type="button"
                  role="menuitem"
                  onClick={() => {
                    setAddMenuOpen(false);
                    setCustomNodeModal({ open: true, protocol: "anthropic" });
                  }}
                  className="flex w-full cursor-pointer items-start gap-2.5 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-[var(--color-surface-2)]"
                >
                  <Terminal className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]" />
                  <span className="min-w-0">
                    <span className="block text-xs font-semibold text-[var(--color-text-main)]">
                      {isId ? "Anthropic Compatible" : "Anthropic Compatible"}
                    </span>
                    <span className="block truncate text-[11px] text-[var(--color-text-muted)]">
                      Messages API (/v1/messages)
                    </span>
                  </span>
                </button>
              </div>
            ) : null}
          </div>
        </div>
      </div>

      {/* 4 Clickable Fleet Telemetry KPI Cards */}
      <div className="grid grid-cols-2 gap-2.5 lg:grid-cols-4">
        <button
          type="button"
          onClick={() => toggleStatusKpi("active")}
          aria-pressed={statusFilter === "active"}
          className={cn(
            "flex cursor-pointer flex-col justify-between rounded-xl border bg-[var(--color-surface)] px-3.5 py-3 text-left transition-all",
            statusFilter === "active"
              ? "border-[var(--color-primary)] bg-[var(--color-primary)]/5 ring-1 ring-[var(--color-primary)]/25"
              : "border-[var(--color-border)] hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-surface-2)]/50",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-[var(--color-text-muted)]">
              {isId ? "Provider Aktif" : "Active Providers"}
            </span>
            <CheckCircle2 className="h-4 w-4 text-[var(--status-ok)]" aria-hidden="true" />
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="font-mono text-lg font-semibold tabular-nums text-[var(--color-text-main)]">
              {telemetry.activeProviders}
            </span>
            <span className="truncate text-[11px] text-[var(--color-text-muted)]">
              / {telemetry.totalCatalog} {isId ? "katalog" : "in catalog"}
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => setStatusFilter("all")}
          className="flex cursor-pointer flex-col justify-between rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3.5 py-3 text-left transition-all hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-surface-2)]/50"
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-[var(--color-text-muted)]">
              {isId ? "Koneksi Aktif" : "Active Connections"}
            </span>
            <Server className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="font-mono text-lg font-semibold tabular-nums text-[var(--color-text-main)]">
              {telemetry.enabledConnections}
              <span className="text-xs font-normal text-[var(--color-text-muted)]">
                /{telemetry.totalConnections}
              </span>
            </span>
            <span className="truncate text-[11px] text-[var(--color-text-muted)]">
              {telemetry.oauthConnections} OAuth · {telemetry.apiKeyConnections} API Key
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => toggleStatusKpi("error")}
          aria-pressed={statusFilter === "error"}
          className={cn(
            "flex cursor-pointer flex-col justify-between rounded-xl border bg-[var(--color-surface)] px-3.5 py-3 text-left transition-all",
            statusFilter === "error"
              ? "border-[var(--status-error)] bg-[var(--status-error-bg)] ring-1 ring-[var(--status-error)]/30"
              : telemetry.attentionProviders > 0
                ? "border-[var(--status-error)]/40 hover:border-[var(--status-error)]/70"
                : "border-[var(--color-border)] hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-surface-2)]/50",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-[var(--color-text-muted)]">
              {isId ? "Perlu Perhatian" : "Needs Attention"}
            </span>
            <AlertTriangle
              className={cn(
                "h-4 w-4",
                telemetry.attentionProviders > 0
                  ? "text-[var(--status-error)]"
                  : "text-[var(--color-text-muted)]/60",
              )}
              aria-hidden="true"
            />
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span
              className={cn(
                "font-mono text-lg font-semibold tabular-nums",
                telemetry.attentionProviders > 0
                  ? "text-[var(--status-error)]"
                  : "text-[var(--color-text-main)]",
              )}
            >
              {telemetry.attentionProviders}
            </span>
            <span className="truncate text-[11px] text-[var(--color-text-muted)]">
              {telemetry.attentionProviders > 0
                ? isId
                  ? "Provider bermasalah"
                  : "Providers with errors"
                : isId
                  ? "Semua koneksi sehat"
                  : "All connections clear"}
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => toggleStatusKpi("inactive")}
          aria-pressed={statusFilter === "inactive"}
          className={cn(
            "flex cursor-pointer flex-col justify-between rounded-xl border bg-[var(--color-surface)] px-3.5 py-3 text-left transition-all",
            statusFilter === "inactive"
              ? "border-[var(--color-primary)] bg-[var(--color-primary)]/5 ring-1 ring-[var(--color-primary)]/25"
              : "border-[var(--color-border)] hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-surface-2)]/50",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-[var(--color-text-muted)]">
              {isId ? "Nonaktif" : "Disabled"}
            </span>
            <PauseCircle className="h-4 w-4 text-[var(--color-text-muted)]" aria-hidden="true" />
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="font-mono text-lg font-semibold tabular-nums text-[var(--color-text-main)]">
              {telemetry.disabledProviders}
            </span>
            <span className="truncate text-[11px] text-[var(--color-text-muted)]">
              {isId ? "Provider dijeda" : "Paused providers"}
            </span>
          </div>
        </button>
      </div>

      {/* Feedback / Error / Batch Test Banners */}
      {catalogQuery.isError || connectionsQuery.isError ? (
        <Alert variant="destructive">
          <AlertTitle>
            {isId ? "Gagal memuat data provider" : "Failed to load providers"}
          </AlertTitle>
          <AlertDescription>
            {catalogQuery.error?.message ?? connectionsQuery.error?.message}
          </AlertDescription>
        </Alert>
      ) : null}

      {actionError ? (
        <Alert variant="destructive">
          <div className="flex items-center justify-between gap-2">
            <AlertDescription>{actionError}</AlertDescription>
            <button
              type="button"
              onClick={() => setActionError(null)}
              aria-label="Dismiss error"
              className="cursor-pointer rounded p-1 hover:bg-black/10"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        </Alert>
      ) : null}

      {feedbackBanner ? (
        <Alert>
          <div className="flex items-center justify-between gap-2">
            <AlertDescription>{feedbackBanner}</AlertDescription>
            <button
              type="button"
              onClick={() => setFeedbackBanner(null)}
              aria-label="Dismiss notice"
              className="cursor-pointer rounded p-1 hover:bg-black/10"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        </Alert>
      ) : null}

      {batchSummary ? (
        <div className="rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-4">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-xs font-semibold text-[var(--color-text-main)]">
                {isId ? "Hasil Tes:" : "Test Results:"} {batchSummary.label}
              </span>
              <Badge variant="success">{batchSummary.passed} Passed</Badge>
              {batchSummary.failed > 0 ? (
                <Badge variant="destructive">{batchSummary.failed} Failed</Badge>
              ) : null}
            </div>
            <button
              type="button"
              onClick={() => setBatchSummary(null)}
              aria-label="Close test results"
              className="cursor-pointer rounded p-1 text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
          <div className="mt-2.5 flex flex-wrap gap-2">
            {batchSummary.results.map((r) => (
              <span
                key={r.connectionId}
                className={cn(
                  "inline-flex items-center gap-1.5 rounded-lg border px-2 py-1 font-mono text-[11px] tabular-nums",
                  r.ok
                    ? "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]"
                    : "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]",
                )}
              >
                <span>{r.connectionName}</span>
                <span>{r.ok ? `${r.statusCode} · ${r.latencyMs}ms` : r.error || "ERR"}</span>
              </span>
            ))}
          </div>
        </div>
      ) : null}

      {/* Search & Filter Toolbar Card (main ProvidersToolbar parity) */}
      <div className="flex flex-col gap-3 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-3.5">
        {/* Row 1: Search Input + Status Filter Pills */}
        <div className="flex flex-col gap-2.5 lg:flex-row lg:items-center lg:justify-between">
          <div className="relative max-w-md flex-1">
            <Search
              className="pointer-events-none absolute top-1/2 left-3 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-muted)]"
              aria-hidden="true"
            />
            <input
              ref={searchInputRef}
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Escape") {
                  if (search) setSearch("");
                  else searchInputRef.current?.blur();
                }
              }}
              placeholder={
                isId
                  ? "Cari nama, ID, model, atau protokol..."
                  : "Search by name, ID, model, or protocol..."
              }
              aria-label={isId ? "Cari provider" : "Search providers"}
              className="h-9 w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)]/75 pr-9 pl-9 text-xs text-[var(--color-text-main)] placeholder:text-[var(--color-text-muted)] transition-colors focus:border-[var(--color-primary)] focus:outline-none"
            />
            {search ? (
              <button
                type="button"
                onClick={() => setSearch("")}
                aria-label="Clear search"
                className="absolute top-1/2 right-2 -translate-y-1/2 cursor-pointer rounded p-1 text-[var(--color-text-muted)] hover:bg-[var(--color-surface-3)] hover:text-[var(--color-text-main)]"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            ) : (
              <kbd className="pointer-events-none absolute top-1/2 right-2.5 hidden -translate-y-1/2 rounded border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--color-text-muted)] sm:inline-flex">
                /
              </kbd>
            )}
          </div>

          <fieldset
            className="custom-scrollbar m-0 flex items-center gap-1.5 overflow-x-auto border-0 p-0"
            aria-label={isId ? "Filter status provider" : "Filter providers by status"}
          >
            {statusPills.map((pill) => {
              const count = statusCounts[pill.value];
              const active = statusFilter === pill.value;
              const isErrorAttention = pill.value === "error" && count > 0;
              return (
                <button
                  key={pill.value}
                  type="button"
                  aria-pressed={active}
                  onClick={() => setStatusFilter(pill.value)}
                  className={cn(
                    "inline-flex h-8 shrink-0 cursor-pointer items-center gap-1.5 whitespace-nowrap rounded-lg border px-2.5 text-xs font-medium transition-colors",
                    active
                      ? pill.value === "error"
                        ? "border-[var(--status-error)] bg-[var(--status-error)] text-[var(--color-on-primary)]"
                        : "border-[var(--color-primary)] bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                      : isErrorAttention
                        ? "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]"
                        : "border-[var(--color-border)] bg-[var(--color-surface-2)]/70 text-[var(--color-text-muted)] hover:bg-[var(--color-surface-3)] hover:text-[var(--color-text-main)]",
                  )}
                >
                  <span>{pill.label}</span>
                  <span className="font-mono text-[10px] tabular-nums opacity-85">{count}</span>
                </button>
              );
            })}
          </fieldset>
        </div>

        {/* Row 2: Category Tabs & Reset */}
        <div className="flex items-center justify-between gap-2 border-t border-[var(--color-border-subtle)] pt-2.5">
          <div
            className="custom-scrollbar flex items-center gap-1.5 overflow-x-auto text-xs"
            role="tablist"
            aria-label={isId ? "Kategori provider" : "Provider categories"}
          >
            {categoryTabs.map((tab) => {
              const active = categoryTab === tab.value;
              const count = categoryCounts[tab.value];
              return (
                <button
                  key={tab.value}
                  type="button"
                  role="tab"
                  aria-selected={active}
                  onClick={() => setCategoryTab(tab.value)}
                  className={cn(
                    "inline-flex h-8 shrink-0 cursor-pointer items-center gap-1.5 whitespace-nowrap rounded-lg border px-3 font-medium transition-colors",
                    active
                      ? "border-[var(--color-primary)]/30 bg-[var(--color-primary)]/12 font-semibold text-[var(--color-primary)]"
                      : "border-transparent bg-transparent text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]",
                  )}
                >
                  {tab.icon}
                  <span>{tab.label}</span>
                  <span className="font-mono text-[10px] tabular-nums opacity-80">({count})</span>
                </button>
              );
            })}
          </div>

          {hasActiveFilters ? (
            <button
              type="button"
              onClick={resetFilters}
              className="inline-flex h-8 shrink-0 cursor-pointer items-center gap-1 rounded-lg px-2.5 text-xs font-medium text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
            >
              <FilterX className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Reset</span>
            </button>
          ) : null}
        </div>
      </div>

      {/* Categorized Provider Card Grids */}
      {filteredEntries.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-[var(--color-border)] bg-[var(--color-surface)] px-6 py-12 text-center">
          <Server className="h-8 w-8 text-[var(--color-text-muted)]/60" aria-hidden="true" />
          <p className="mt-3 text-sm font-semibold text-[var(--color-text-main)]">
            {isId ? "Tidak ada provider yang cocok" : "No matching providers"}
          </p>
          <p className="mt-1 text-xs text-[var(--color-text-muted)]">
            {isId
              ? "Coba ubah kata kunci pencarian atau reset filter aktif."
              : "Try adjusting your search query or clearing active filters."}
          </p>
          {hasActiveFilters ? (
            <Button
              type="button"
              size="sm"
              variant="secondary"
              className="mt-4"
              onClick={resetFilters}
            >
              <FilterX className="h-3.5 w-3.5" />
              <span>{isId ? "Reset Filter" : "Reset Filters"}</span>
            </Button>
          ) : null}
        </div>
      ) : (
        <div className="space-y-6">
          {sectionsToRender.map((section) => {
            const sectionConns = section.entries.flatMap((e) => e.connections);
            const activeInSection = sectionConns.filter((c) => c.enabled).length;
            const isTestingSection = testingScope === `cat:${section.key}`;

            return (
              <section key={section.key} className="space-y-3">
                <SectionHeader
                  icon={section.icon}
                  title={section.title}
                  count={section.entries.length}
                  action={
                    <button
                      type="button"
                      disabled={activeInSection === 0 || testingScope !== null}
                      onClick={() =>
                        runBatchTests(`cat:${section.key}`, section.title, sectionConns)
                      }
                      className={cn(
                        "inline-flex cursor-pointer items-center gap-1.5 rounded-lg border px-2.5 py-1 text-xs font-medium transition-colors",
                        isTestingSection
                          ? "border-[var(--color-primary)]/40 bg-[var(--color-primary)]/15 text-[var(--color-primary)]"
                          : "border-[var(--color-border)] bg-[var(--color-surface-2)]/80 text-[var(--color-text-muted)] hover:border-[var(--color-primary)]/45 hover:text-[var(--color-text-main)]",
                        (activeInSection === 0 || (testingScope !== null && !isTestingSection)) &&
                          "cursor-not-allowed opacity-45",
                      )}
                    >
                      {isTestingSection ? (
                        <RefreshCw className="h-3 w-3 animate-spin" />
                      ) : (
                        <Play className="h-3 w-3" />
                      )}
                      <span>
                        {isTestingSection
                          ? isId
                            ? "Menguji..."
                            : "Testing..."
                          : isId
                            ? "Tes Grup"
                            : "Test All"}
                      </span>
                    </button>
                  }
                />

                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {(section.key === "apikey" && !showAllApiKey && !hasActiveFilters
                    ? section.entries.slice(0, 20)
                    : section.entries
                  ).map((entry) => (
                    <ProviderCard
                      key={entry.id}
                      entry={entry}
                      isId={isId}
                      onSelect={handleCardSelect}
                      onQuickTest={(e) => runBatchTests(`provider:${e.id}`, e.name, e.connections)}
                      onToggleProvider={handleToggleProvider}
                      isTesting={testingScope === `provider:${entry.id}`}
                      isToggling={togglingProviderId === entry.id}
                      disableQuickTest={testingScope !== null}
                    />
                  ))}
                </div>
                {section.key === "apikey" && !hasActiveFilters && section.entries.length > 20 ? (
                  <div className="flex justify-center pt-1">
                    <Button
                      type="button"
                      size="sm"
                      variant="secondary"
                      onClick={() => setShowAllApiKey((v) => !v)}
                    >
                      {showAllApiKey
                        ? isId
                          ? "Tampilkan Lebih Sedikit"
                          : "Show Fewer API Key Providers"
                        : isId
                          ? `Tampilkan Semua (${section.entries.length} provider)`
                          : `Show All ${section.entries.length} API Key Providers`}
                    </Button>
                  </div>
                ) : null}
              </section>
            );
          })}
        </div>
      )}

      {/* Modals & Slide-Over Drawer */}
      <AddApiKeyModal
        open={apiKeyModal.open}
        onClose={() => setApiKeyModal({ open: false, providerId: null })}
        initialProviderId={apiKeyModal.providerId}
        allEntries={allEntries}
        isId={isId}
        onCreated={(msg) => setFeedbackBanner(msg)}
      />

      <OAuthWizardModal
        open={oauthModal.open}
        onClose={() => setOauthModal({ open: false, providerId: null })}
        initialProviderId={oauthModal.providerId}
        allEntries={allEntries}
        isId={isId}
        onConnected={(msg) => setFeedbackBanner(msg)}
        onSwitchToApiKey={(providerId) => setApiKeyModal({ open: true, providerId })}
      />

      <AddCustomNodeModal
        open={customNodeModal.open}
        onClose={() => setCustomNodeModal((prev) => ({ ...prev, open: false }))}
        initialProtocol={customNodeModal.protocol}
        isId={isId}
        onCreated={(node, meta) =>
          setFeedbackBanner(
            meta?.connectionCreated
              ? isId
                ? `Custom endpoint "${node.name}" beserta koneksi awalnya berhasil disimpan.`
                : `Registered custom endpoint "${node.name}" with its initial connection.`
              : isId
                ? `Custom endpoint "${node.name}" berhasil didaftarkan.`
                : `Registered custom endpoint "${node.name}".`,
          )
        }
      />

      <ProviderConnectionsDrawer
        entry={selectedEntry}
        onClose={() => setSelectedProviderId(null)}
        isId={isId}
        busyConnectionId={busyConnectionId}
        liveTestResults={liveTestResults}
        onTestConnection={handleTestSingleConnection}
        onToggleConnection={handleToggleSingleConnection}
        onDeleteConnection={handleDeleteSingleConnection}
        onAddConnection={handleAddConnectionForEntry}
        onDeleteCustomNode={handleDeleteCustomNode}
        isDeletingNode={deleteNodeMutation.isPending}
      />
    </div>
  );
}
