import {
  AlertCircle,
  ArrowDown,
  ArrowUp,
  ArrowUpDown,
  Check,
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  Clock,
  Copy,
  Eye,
  EyeOff,
  FlaskConical,
  KeyRound,
  Layers,
  RefreshCw,
  Search,
  Sparkles,
  X,
  XCircle,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router";
import { useCombos } from "@/features/combos/hooks";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { AssignKeysModal } from "./AssignKeysModal";
import {
  clearCachedTestResults,
  comboTestResultKey,
  computeVisibilityCounts,
  countKeyAccess,
  DEFAULT_COMBO_FILTERS,
  DEFAULT_MODEL_FILTERS,
  filterCombos,
  filterModels,
  formatRelativeTime,
  formatTokenCount,
  getEnabledCapabilities,
  loadCachedTestResults,
  modelTestResultKey,
  normalizeComboRows,
  saveCachedTestResults,
  sortCombos,
  sortModels,
} from "./catalogUtils";
import {
  runSingleComboTest,
  runSingleModelTest,
  useCatalogApiKeys,
  useModelCatalog,
  useProviderHealthMatrix,
  useToggleModelVisibilityMutation,
} from "./hooks";
import type {
  BulkTestProgress,
  CatalogComboRow,
  CatalogModelRow,
  CatalogTab,
  CatalogTestResult,
  ComboCatalogFilters,
  ComboSortField,
  ModelCatalogFilters,
  ModelSortField,
  ProviderHealthState,
  SortDirection,
} from "./types";

const PAGE_SIZE = 25;
const BULK_CONFIRM_THRESHOLD = 10;

function HealthDot({ state, isId }: { state: ProviderHealthState; isId: boolean }) {
  const config: Record<ProviderHealthState, { dot: string; labelEn: string; labelId: string }> = {
    healthy: { dot: "bg-[var(--status-ok)]", labelEn: "Healthy", labelId: "Sehat" },
    degraded: { dot: "bg-[var(--status-warn)]", labelEn: "Degraded", labelId: "Terdegradasi" },
    down: { dot: "bg-[var(--status-error)]", labelEn: "Down", labelId: "Mati" },
    untested: { dot: "bg-[var(--text-muted)]", labelEn: "Untested", labelId: "Belum diuji" },
  };
  const item = config[state] ?? config.untested;
  const label = isId ? item.labelId : item.labelEn;
  return (
    <span className="inline-flex items-center gap-1.5 text-[11px] text-[var(--text-secondary)]">
      <span className={`h-2 w-2 rounded-full ${item.dot}`} aria-hidden="true" />
      <span>{label}</span>
    </span>
  );
}

function TestResultCell({
  result,
  isTesting,
  nowMs,
  isId,
}: {
  result?: CatalogTestResult;
  isTesting: boolean;
  nowMs: number;
  isId: boolean;
}) {
  if (isTesting) {
    return (
      <span className="inline-flex items-center gap-1.5 text-xs text-[var(--accent-text)]">
        <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
        <span>{isId ? "Menguji..." : "Testing..."}</span>
      </span>
    );
  }
  if (!result) {
    return (
      <span className="text-xs text-[var(--text-muted)]">{isId ? "Belum diuji" : "Untested"}</span>
    );
  }

  const rel = formatRelativeTime(result.testedAt, nowMs, isId);
  const latency = typeof result.latencyMs === "number" ? `${result.latencyMs}ms` : "";

  if (result.status === "ok") {
    return (
      <div className="space-y-0.5">
        <div className="flex items-center gap-1.5">
          <Badge variant="success" className="gap-1 text-[11px]">
            <CheckCircle2 className="h-3 w-3" aria-hidden="true" />
            <span>OK {latency ? `· ${latency}` : ""}</span>
          </Badge>
        </div>
        <p className="text-[10px] text-[var(--text-muted)]">
          {result.resolvedBy ? `${result.resolvedBy} · ` : ""}
          {rel}
        </p>
      </div>
    );
  }

  if (result.status === "slow") {
    return (
      <div className="space-y-0.5">
        <Badge variant="warning" className="gap-1 text-[11px]">
          <Clock className="h-3 w-3" aria-hidden="true" />
          <span>
            {isId ? "Lambat" : "Slow"} {latency ? `· ${latency}` : ""}
          </span>
        </Badge>
        <p className="text-[10px] text-[var(--text-muted)]">{rel}</p>
      </div>
    );
  }

  const classLabel: Record<string, string> = {
    "rate-limited": isId ? "Rate Limit" : "Rate Limited",
    quota: isId ? "Kuota Habis" : "Quota Exceeded",
    timeout: "Timeout",
    other: isId ? "Gagal" : "Error",
  };
  const badgeText = classLabel[result.errorClass ?? "other"] ?? "Error";

  return (
    <div className="max-w-[220px] space-y-0.5" title={result.error || undefined}>
      <Badge variant="destructive" className="gap-1 text-[11px]">
        <XCircle className="h-3 w-3" aria-hidden="true" />
        <span>
          {badgeText}
          {result.statusCode ? ` (${result.statusCode})` : ""}
        </span>
      </Badge>
      {result.error ? (
        <p className="truncate font-mono text-[10px] text-[var(--status-error)]">{result.error}</p>
      ) : null}
      <p className="text-[10px] text-[var(--text-muted)]">{rel}</p>
    </div>
  );
}

function SortHeaderButton<T extends string>({
  field,
  label,
  activeField,
  direction,
  onSort,
}: {
  field: T;
  label: string;
  activeField: T;
  direction: SortDirection;
  onSort: (f: T) => void;
}) {
  const active = activeField === field;
  return (
    <button
      type="button"
      onClick={() => onSort(field)}
      className="inline-flex items-center gap-1 text-left text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)] transition-colors hover:text-[var(--text-primary)]"
    >
      <span>{label}</span>
      {active ? (
        direction === "asc" ? (
          <ArrowUp className="h-3 w-3 text-[var(--accent-primary)]" aria-hidden="true" />
        ) : (
          <ArrowDown className="h-3 w-3 text-[var(--accent-primary)]" aria-hidden="true" />
        )
      ) : (
        <ArrowUpDown className="h-3 w-3 opacity-40" aria-hidden="true" />
      )}
    </button>
  );
}

async function runWithConcurrency<T>(
  items: T[],
  concurrency: number,
  signal: AbortSignal,
  worker: (item: T) => Promise<void>,
): Promise<void> {
  let idx = 0;
  const active = new Set<Promise<void>>();
  while (idx < items.length && !signal.aborted) {
    while (active.size < concurrency && idx < items.length && !signal.aborted) {
      const current = items[idx++];
      if (current === undefined) break;
      const p = worker(current).finally(() => {
        active.delete(p);
      });
      active.add(p);
    }
    if (active.size > 0) {
      await Promise.race(active);
    }
  }
  await Promise.all(active);
}

export function ModelCatalogConsole() {
  const locale = useUIStore((s) => s.locale);
  const isId = locale === "id";

  const catalogQuery = useModelCatalog();
  const combosQuery = useCombos();
  const healthQuery = useProviderHealthMatrix();
  const keysQuery = useCatalogApiKeys();
  const visibilityMutation = useToggleModelVisibilityMutation();

  const [activeTab, setActiveTab] = useState<CatalogTab>("models");
  const [modelFilters, setModelFilters] = useState<ModelCatalogFilters>(DEFAULT_MODEL_FILTERS);
  const [comboFilters, setComboFilters] = useState<ComboCatalogFilters>(DEFAULT_COMBO_FILTERS);
  const [modelSort, setModelSort] = useState<{ field: ModelSortField; direction: SortDirection }>({
    field: "provider",
    direction: "asc",
  });
  const [comboSort, setComboSort] = useState<{ field: ComboSortField; direction: SortDirection }>({
    field: "name",
    direction: "asc",
  });
  const [modelPage, setModelPage] = useState(1);
  const [comboPage, setComboPage] = useState(1);

  const [selectedModelKeys, setSelectedModelKeys] = useState<Set<string>>(() => new Set());
  const [selectedComboNames, setSelectedComboNames] = useState<Set<string>>(() => new Set());
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const [testResults, setTestResults] = useState<Record<string, CatalogTestResult>>({});
  const [testingKeys, setTestingKeys] = useState<Set<string>>(() => new Set());
  const [bulkRunning, setBulkRunning] = useState(false);
  const [bulkProgress, setBulkProgress] = useState<BulkTestProgress>({
    kind: "models",
    completed: 0,
    total: 0,
    cancelled: false,
  });
  const [bulkConfirmPending, setBulkConfirmPending] = useState<{
    kind: "models" | "combos";
    count: number;
    onConfirm: () => void;
  } | null>(null);
  const bulkAbortRef = useRef<AbortController | null>(null);

  const [assignModalOpen, setAssignModalOpen] = useState(false);
  const [assignTargets, setAssignTargets] = useState<string[]>([]);
  const [visibilityError, setVisibilityError] = useState<string | null>(null);
  const [nowMs, setNowMs] = useState(() => Date.now());

  useEffect(() => {
    setTestResults(loadCachedTestResults());
    const timer = window.setInterval(() => setNowMs(Date.now()), 30_000);
    return () => {
      window.clearInterval(timer);
      bulkAbortRef.current?.abort();
    };
  }, []);

  const models: CatalogModelRow[] = catalogQuery.data ?? [];
  const combos: CatalogComboRow[] = useMemo(
    () => normalizeComboRows(combosQuery.data),
    [combosQuery.data],
  );
  const healthMatrix: Record<string, ProviderHealthState> = healthQuery.data ?? {};
  const apiKeys = keysQuery.data ?? [];

  const visibilityCounts = useMemo(() => computeVisibilityCounts(models), [models]);

  // Derived filter options from live catalog data
  const providerOptions = useMemo(() => {
    const map = new Map<string, string>();
    for (const m of models) {
      if (!map.has(m.providerId)) {
        map.set(m.providerId, m.provider);
      }
    }
    return [...map.entries()].sort((a, b) => a[1].localeCompare(b[1]));
  }, [models]);

  const typeOptions = useMemo(
    () => [...new Set(models.map((m) => m.type).filter(Boolean))].sort(),
    [models],
  );
  const subtypeOptions = useMemo(
    () => [...new Set(models.map((m) => m.subtype).filter((s): s is string => Boolean(s)))].sort(),
    [models],
  );
  const capabilityOptions = useMemo(() => {
    const set = new Set<string>();
    for (const m of models) {
      for (const cap of getEnabledCapabilities(m)) {
        set.add(cap);
      }
    }
    return [...set].sort();
  }, [models]);
  const strategyOptions = useMemo(
    () => [...new Set(combos.map((c) => c.strategy).filter(Boolean))].sort(),
    [combos],
  );

  // Filtered + sorted rows
  const filteredModels = useMemo(() => {
    const filtered = filterModels(models, modelFilters, healthMatrix, testResults);
    return sortModels(filtered, modelSort.field, modelSort.direction, testResults);
  }, [models, modelFilters, healthMatrix, testResults, modelSort]);

  const filteredCombos = useMemo(() => {
    const filtered = filterCombos(combos, comboFilters, testResults);
    return sortCombos(filtered, comboSort.field, comboSort.direction, testResults);
  }, [combos, comboFilters, testResults, comboSort]);

  // Keep pagination within bounds whenever filtered counts change
  const totalModelPages = Math.max(1, Math.ceil(filteredModels.length / PAGE_SIZE));
  const safeModelPage = Math.min(modelPage, totalModelPages);
  const pagedModels = useMemo(() => {
    const start = (safeModelPage - 1) * PAGE_SIZE;
    return filteredModels.slice(start, start + PAGE_SIZE);
  }, [filteredModels, safeModelPage]);

  const totalComboPages = Math.max(1, Math.ceil(filteredCombos.length / PAGE_SIZE));
  const safeComboPage = Math.min(comboPage, totalComboPages);
  const pagedCombos = useMemo(() => {
    const start = (safeComboPage - 1) * PAGE_SIZE;
    return filteredCombos.slice(start, start + PAGE_SIZE);
  }, [filteredCombos, safeComboPage]);

  const filteredVisibleModels = useMemo(
    () => filteredModels.filter((m) => !m.hidden),
    [filteredModels],
  );
  const filteredHiddenModels = useMemo(
    () => filteredModels.filter((m) => m.hidden === true),
    [filteredModels],
  );

  const isRefreshing =
    catalogQuery.isFetching ||
    combosQuery.isFetching ||
    healthQuery.isFetching ||
    keysQuery.isFetching;

  const handleRefreshAll = () => {
    void Promise.all([
      catalogQuery.refetch(),
      combosQuery.refetch(),
      healthQuery.refetch(),
      keysQuery.refetch(),
    ]);
  };

  const handleCopyModelId = async (id: string) => {
    try {
      await navigator.clipboard.writeText(id);
      setCopiedId(id);
      window.setTimeout(() => setCopiedId((cur) => (cur === id ? null : cur)), 1800);
    } catch {
      // Ignore clipboard permission errors
    }
  };

  const persistResult = useCallback((result: CatalogTestResult) => {
    setTestResults(saveCachedTestResults([result]));
  }, []);

  const handleTestSingleModel = async (row: CatalogModelRow) => {
    const key = modelTestResultKey(row.providerId, row.id);
    if (testingKeys.has(key)) return;
    setTestingKeys((prev) => new Set(prev).add(key));
    try {
      const result = await runSingleModelTest(row.providerId, row.id);
      persistResult(result);
    } finally {
      setTestingKeys((prev) => {
        const next = new Set(prev);
        next.delete(key);
        return next;
      });
    }
  };

  const handleTestSingleCombo = async (row: CatalogComboRow) => {
    const key = comboTestResultKey(row.name);
    if (testingKeys.has(key)) return;
    setTestingKeys((prev) => new Set(prev).add(key));
    try {
      const result = await runSingleComboTest(row.name);
      persistResult(result);
    } finally {
      setTestingKeys((prev) => {
        const next = new Set(prev);
        next.delete(key);
        return next;
      });
    }
  };

  const executeBulkModelTests = async (targets: CatalogModelRow[]) => {
    if (targets.length === 0 || bulkRunning) return;
    const ctrl = new AbortController();
    bulkAbortRef.current = ctrl;
    setBulkRunning(true);
    setBulkProgress({ kind: "models", completed: 0, total: targets.length, cancelled: false });

    // Group by provider so concurrency across providers is 2 while keeping per-provider sequential (OmniRoute pattern)
    const byProvider = new Map<string, CatalogModelRow[]>();
    for (const m of targets) {
      const list = byProvider.get(m.providerId) ?? [];
      list.push(m);
      byProvider.set(m.providerId, list);
    }

    let completed = 0;
    try {
      await runWithConcurrency([...byProvider.values()], 2, ctrl.signal, async (group) => {
        for (const row of group) {
          if (ctrl.signal.aborted) return;
          const key = modelTestResultKey(row.providerId, row.id);
          setTestingKeys((prev) => new Set(prev).add(key));
          try {
            const res = await runSingleModelTest(row.providerId, row.id, ctrl.signal);
            if (!ctrl.signal.aborted) {
              persistResult(res);
              completed++;
              setBulkProgress((p) => ({ ...p, completed }));
            }
          } finally {
            setTestingKeys((prev) => {
              const next = new Set(prev);
              next.delete(key);
              return next;
            });
          }
        }
      });
    } finally {
      setBulkRunning(false);
      bulkAbortRef.current = null;
    }
  };

  const executeBulkComboTests = async (targets: CatalogComboRow[]) => {
    if (targets.length === 0 || bulkRunning) return;
    const ctrl = new AbortController();
    bulkAbortRef.current = ctrl;
    setBulkRunning(true);
    setBulkProgress({ kind: "combos", completed: 0, total: targets.length, cancelled: false });

    let completed = 0;
    try {
      await runWithConcurrency(targets, 2, ctrl.signal, async (row) => {
        if (ctrl.signal.aborted) return;
        const key = comboTestResultKey(row.name);
        setTestingKeys((prev) => new Set(prev).add(key));
        try {
          const res = await runSingleComboTest(row.name, ctrl.signal);
          if (!ctrl.signal.aborted) {
            persistResult(res);
            completed++;
            setBulkProgress((p) => ({ ...p, completed }));
          }
        } finally {
          setTestingKeys((prev) => {
            const next = new Set(prev);
            next.delete(key);
            return next;
          });
        }
      });
    } finally {
      setBulkRunning(false);
      bulkAbortRef.current = null;
    }
  };

  const triggerBulkModelTests = (targets: CatalogModelRow[]) => {
    if (targets.length >= BULK_CONFIRM_THRESHOLD) {
      setBulkConfirmPending({
        kind: "models",
        count: targets.length,
        onConfirm: () => {
          setBulkConfirmPending(null);
          void executeBulkModelTests(targets);
        },
      });
      return;
    }
    void executeBulkModelTests(targets);
  };

  const triggerBulkComboTests = (targets: CatalogComboRow[]) => {
    if (targets.length >= BULK_CONFIRM_THRESHOLD) {
      setBulkConfirmPending({
        kind: "combos",
        count: targets.length,
        onConfirm: () => {
          setBulkConfirmPending(null);
          void executeBulkComboTests(targets);
        },
      });
      return;
    }
    void executeBulkComboTests(targets);
  };

  const handleCancelBulkTests = () => {
    bulkAbortRef.current?.abort();
    setBulkProgress((p) => ({ ...p, cancelled: true }));
    setBulkRunning(false);
  };

  const handleClearResults = () => {
    clearCachedTestResults();
    setTestResults({});
    setBulkProgress((p) => ({ ...p, cancelled: false }));
  };

  const handleToggleVisibility = async (rows: CatalogModelRow[], hide: boolean) => {
    if (rows.length === 0 || visibilityMutation.isPending) return;
    setVisibilityError(null);
    try {
      await visibilityMutation.mutateAsync({
        qualifiedIds: rows.map((r) => r.qualifiedId),
        hide,
      });
    } catch (err) {
      setVisibilityError(err instanceof Error ? err.message : "Failed to update visibility");
    }
  };

  const handleSortModels = (field: ModelSortField) => {
    setModelSort((prev) =>
      prev.field === field
        ? { field, direction: prev.direction === "asc" ? "desc" : "asc" }
        : { field, direction: "asc" },
    );
  };

  const handleSortCombos = (field: ComboSortField) => {
    setComboSort((prev) =>
      prev.field === field
        ? { field, direction: prev.direction === "asc" ? "desc" : "asc" }
        : { field, direction: "asc" },
    );
  };

  const allPagedModelsSelected =
    pagedModels.length > 0 && pagedModels.every((m) => selectedModelKeys.has(m.qualifiedId));
  const toggleAllPagedModels = () => {
    setSelectedModelKeys((prev) => {
      const next = new Set(prev);
      if (allPagedModelsSelected) {
        for (const m of pagedModels) next.delete(m.qualifiedId);
      } else {
        for (const m of pagedModels) next.add(m.qualifiedId);
      }
      return next;
    });
  };

  const allPagedCombosSelected =
    pagedCombos.length > 0 && pagedCombos.every((c) => selectedComboNames.has(c.name));
  const toggleAllPagedCombos = () => {
    setSelectedComboNames((prev) => {
      const next = new Set(prev);
      if (allPagedCombosSelected) {
        for (const c of pagedCombos) next.delete(c.name);
      } else {
        for (const c of pagedCombos) next.add(c.name);
      }
      return next;
    });
  };

  const selectedModelRows = useMemo(
    () => models.filter((m) => selectedModelKeys.has(m.qualifiedId)),
    [models, selectedModelKeys],
  );
  const selectedComboRows = useMemo(
    () => combos.filter((c) => selectedComboNames.has(c.name)),
    [combos, selectedComboNames],
  );

  const hasTestResults = Object.keys(testResults).length > 0;
  const hasActiveModelFilters =
    JSON.stringify(modelFilters) !== JSON.stringify(DEFAULT_MODEL_FILTERS);
  const hasActiveComboFilters =
    JSON.stringify(comboFilters) !== JSON.stringify(DEFAULT_COMBO_FILTERS);

  const openAssignModal = (targets: string[]) => {
    setAssignTargets(targets);
    setAssignModalOpen(true);
  };

  const progressPct =
    bulkProgress.total > 0 ? Math.round((bulkProgress.completed / bulkProgress.total) * 100) : 0;

  return (
    <div className="space-y-5">
      {/* Page Header */}
      <section className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            {isId ? "Katalog & Kesehatan Model" : "Model Registry & Health"}
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            {isId ? "Katalog Model" : "Model Catalog"}
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            {isId
              ? "Jelajahi seluruh model dan combo aktif, filter kapabilitas, uji konektivitas langsung, kelola visibilitas, dan atur akses kunci API."
              : "Browse all configured models and combos, filter by capabilities and health, run live probes, manage visibility, and assign API key access."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <div
            data-testid="catalog-visibility-counts"
            className="flex items-center gap-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 py-1.5 text-xs"
          >
            <span className="font-medium text-[var(--text-primary)] tabular-nums">
              {visibilityCounts.visible} {isId ? "terlihat" : "visible"}
            </span>
            <span className="text-[var(--text-muted)]">·</span>
            <span className="text-[var(--text-muted)] tabular-nums">
              {visibilityCounts.hidden} {isId ? "tersembunyi" : "hidden"}
            </span>
            <span className="text-[var(--text-muted)]">·</span>
            <span className="text-[var(--text-secondary)] tabular-nums">
              {visibilityCounts.total} total
            </span>
          </div>

          <Link to="/dashboard/combos/">
            <Button type="button" variant="outline" size="sm">
              <Layers className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
              <span>{isId ? "Kelola Combo" : "Manage Combos"}</span>
            </Button>
          </Link>

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefreshAll}
            disabled={isRefreshing}
          >
            <RefreshCw
              className={`h-3.5 w-3.5 text-[var(--accent-primary)] ${isRefreshing ? "animate-spin" : ""}`}
              aria-hidden="true"
            />
            <span>
              {isRefreshing
                ? isId
                  ? "Menyegarkan..."
                  : "Refreshing..."
                : isId
                  ? "Segarkan"
                  : "Refresh"}
            </span>
          </Button>
        </div>
      </section>

      {/* Underline Tabs (Models | Combos) */}
      <div
        role="tablist"
        aria-label={isId ? "Bagian katalog" : "Catalog sections"}
        className="flex items-center gap-6 border-b border-[var(--border-subtle)]"
      >
        <button
          type="button"
          role="tab"
          id="tab-models"
          aria-selected={activeTab === "models"}
          aria-controls="panel-models"
          onClick={() => setActiveTab("models")}
          className={`relative pb-2.5 text-sm font-medium transition-colors ${
            activeTab === "models"
              ? "text-[var(--text-primary)]"
              : "text-[var(--text-muted)] hover:text-[var(--text-secondary)]"
          }`}
        >
          <span>{isId ? "Model" : "Models"}</span>
          <Badge variant="outline" className="ml-2 text-[11px] tabular-nums">
            {models.length}
          </Badge>
          {activeTab === "models" ? (
            <span className="absolute right-0 bottom-0 left-0 h-0.5 rounded-full bg-[var(--accent-primary)]" />
          ) : null}
        </button>

        <button
          type="button"
          role="tab"
          id="tab-combos"
          aria-selected={activeTab === "combos"}
          aria-controls="panel-combos"
          onClick={() => setActiveTab("combos")}
          className={`relative pb-2.5 text-sm font-medium transition-colors ${
            activeTab === "combos"
              ? "text-[var(--text-primary)]"
              : "text-[var(--text-muted)] hover:text-[var(--text-secondary)]"
          }`}
        >
          <span>{isId ? "Combo" : "Combos"}</span>
          <Badge variant="outline" className="ml-2 text-[11px] tabular-nums">
            {combos.length}
          </Badge>
          {activeTab === "combos" ? (
            <span className="absolute right-0 bottom-0 left-0 h-0.5 rounded-full bg-[var(--accent-primary)]" />
          ) : null}
        </button>
      </div>

      {visibilityError ? (
        <Alert variant="destructive">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
          <div className="flex-1">
            <AlertTitle>
              {isId ? "Gagal memperbarui visibilitas" : "Visibility update failed"}
            </AlertTitle>
            <AlertDescription className="font-mono">{visibilityError}</AlertDescription>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setVisibilityError(null)}
          >
            <X className="h-3.5 w-3.5" />
          </Button>
        </Alert>
      ) : null}

      {/* MODELS TAB PANEL */}
      {activeTab === "models" ? (
        <section
          id="panel-models"
          role="tabpanel"
          aria-labelledby="tab-models"
          className="space-y-4"
        >
          {/* Multi-Facet Filter Bar */}
          <Card className="p-4">
            <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6">
              <div className="relative sm:col-span-2">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  type="text"
                  value={modelFilters.query}
                  onChange={(e) => setModelFilters((f) => ({ ...f, query: e.target.value }))}
                  placeholder={
                    isId ? "Cari model, provider, tipe..." : "Search models, providers, types..."
                  }
                  className="h-9 pl-8 text-xs"
                />
              </div>

              <select
                aria-label={isId ? "Filter provider" : "Filter provider"}
                value={modelFilters.providerId}
                onChange={(e) => setModelFilters((f) => ({ ...f, providerId: e.target.value }))}
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Provider" : "All Providers"}</option>
                {providerOptions.map(([pid, label]) => (
                  <option key={pid} value={pid}>
                    {label}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter tipe" : "Filter type"}
                value={modelFilters.type}
                onChange={(e) => setModelFilters((f) => ({ ...f, type: e.target.value }))}
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Tipe" : "All Types"}</option>
                {typeOptions.map((t) => (
                  <option key={t} value={t}>
                    {t}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter subtipe" : "Filter subtype"}
                value={modelFilters.subtype}
                onChange={(e) => setModelFilters((f) => ({ ...f, subtype: e.target.value }))}
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Subtipe" : "All Subtypes"}</option>
                {subtypeOptions.map((s) => (
                  <option key={s} value={s}>
                    {s}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter kapabilitas" : "Filter capability"}
                value={modelFilters.capability}
                onChange={(e) => setModelFilters((f) => ({ ...f, capability: e.target.value }))}
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Kapabilitas" : "All Capabilities"}</option>
                {capabilityOptions.map((c) => (
                  <option key={c} value={c}>
                    {c}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter harga" : "Filter pricing"}
                value={modelFilters.pricing}
                onChange={(e) =>
                  setModelFilters((f) => ({
                    ...f,
                    pricing: e.target.value as ModelCatalogFilters["pricing"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Harga" : "All Pricing"}</option>
                <option value="free">{isId ? "Gratis saja" : "Free only"}</option>
                <option value="paid">{isId ? "Berbayar" : "Non-free"}</option>
              </select>

              <select
                aria-label={isId ? "Filter kesehatan provider" : "Filter provider health"}
                value={modelFilters.providerHealth}
                onChange={(e) =>
                  setModelFilters((f) => ({
                    ...f,
                    providerHealth: e.target.value as ModelCatalogFilters["providerHealth"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Kesehatan" : "All Health States"}</option>
                <option value="healthy">{isId ? "Sehat" : "Healthy"}</option>
                <option value="degraded">{isId ? "Terdegradasi" : "Degraded"}</option>
                <option value="down">{isId ? "Mati" : "Down"}</option>
                <option value="untested">{isId ? "Belum diuji" : "Untested"}</option>
              </select>

              <select
                aria-label={isId ? "Filter hasil uji" : "Filter test result"}
                value={modelFilters.testResult}
                onChange={(e) =>
                  setModelFilters((f) => ({
                    ...f,
                    testResult: e.target.value as ModelCatalogFilters["testResult"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Hasil Uji" : "All Results"}</option>
                <option value="ok">OK</option>
                <option value="slow">{isId ? "Lambat" : "Slow"}</option>
                <option value="error">{isId ? "Gagal" : "Error"}</option>
                <option value="untested">{isId ? "Belum diuji" : "Untested"}</option>
              </select>

              <select
                aria-label={isId ? "Filter visibilitas" : "Filter visibility"}
                value={modelFilters.visibility}
                onChange={(e) =>
                  setModelFilters((f) => ({
                    ...f,
                    visibility: e.target.value as ModelCatalogFilters["visibility"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Visibilitas" : "All Visibility"}</option>
                <option value="visible">{isId ? "Terlihat saja" : "Visible only"}</option>
                <option value="hidden">{isId ? "Tersembunyi saja" : "Hidden only"}</option>
              </select>

              <Input
                type="number"
                min={0}
                value={modelFilters.minContextLength ?? ""}
                onChange={(e) =>
                  setModelFilters((f) => ({
                    ...f,
                    minContextLength: e.target.value ? Number(e.target.value) : undefined,
                  }))
                }
                placeholder={isId ? "Min konteks" : "Min context"}
                className="h-9 text-xs"
              />

              <div className="flex items-center gap-2">
                <Input
                  type="number"
                  min={0}
                  value={modelFilters.minMaxOutputTokens ?? ""}
                  onChange={(e) =>
                    setModelFilters((f) => ({
                      ...f,
                      minMaxOutputTokens: e.target.value ? Number(e.target.value) : undefined,
                    }))
                  }
                  placeholder={isId ? "Min output" : "Min max output"}
                  className="h-9 text-xs"
                />
                {hasActiveModelFilters ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setModelFilters(DEFAULT_MODEL_FILTERS)}
                    className="shrink-0 text-xs"
                  >
                    {isId ? "Reset" : "Clear"}
                  </Button>
                ) : null}
              </div>
            </div>
          </Card>

          {/* Table Card with OmniRoute Bulk Actions Toolbar */}
          <Card className="overflow-hidden">
            {/* Bulk Actions Toolbar */}
            <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-2.5">
              <div className="flex flex-wrap items-center gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="assign-keys-btn"
                  disabled={selectedModelRows.length === 0 || bulkRunning}
                  onClick={() => openAssignModal(selectedModelRows.map((m) => m.qualifiedId))}
                >
                  <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Tetapkan ke Kunci" : "Assign to Keys"}</span>
                </Button>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="test-selected-btn"
                  disabled={selectedModelRows.length === 0 || bulkRunning}
                  onClick={() => triggerBulkModelTests(selectedModelRows)}
                >
                  <FlaskConical className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {isId
                      ? `Uji Terpilih (${selectedModelRows.length})`
                      : `Test Selected (${selectedModelRows.length})`}
                  </span>
                </Button>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="test-all-filtered-btn"
                  disabled={filteredModels.length === 0 || bulkRunning}
                  onClick={() => triggerBulkModelTests(filteredModels)}
                >
                  <span>
                    {isId
                      ? `Uji Semua Terfilter (${filteredModels.length})`
                      : `Test All Filtered (${filteredModels.length})`}
                  </span>
                </Button>

                {bulkRunning ? (
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    data-testid="cancel-tests-btn"
                    onClick={handleCancelBulkTests}
                  >
                    {isId ? "Batalkan" : "Cancel"}
                  </Button>
                ) : null}

                {hasTestResults && !bulkRunning ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    data-testid="clear-test-results-btn"
                    onClick={handleClearResults}
                  >
                    {isId ? "Bersihkan Hasil" : "Clear Results"}
                  </Button>
                ) : null}

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="hide-filtered-btn"
                  disabled={filteredVisibleModels.length === 0 || visibilityMutation.isPending}
                  onClick={() => void handleToggleVisibility(filteredVisibleModels, true)}
                >
                  <EyeOff className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {isId
                      ? `Sembunyikan (${filteredVisibleModels.length})`
                      : `Hide Filtered (${filteredVisibleModels.length})`}
                  </span>
                </Button>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="show-filtered-btn"
                  disabled={filteredHiddenModels.length === 0 || visibilityMutation.isPending}
                  onClick={() => void handleToggleVisibility(filteredHiddenModels, false)}
                >
                  <Eye className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {isId
                      ? `Tampilkan (${filteredHiddenModels.length})`
                      : `Show Filtered (${filteredHiddenModels.length})`}
                  </span>
                </Button>
              </div>

              {/* Bulk Progress Bar */}
              {bulkRunning && bulkProgress.kind === "models" ? (
                <div className="flex items-center gap-3">
                  <div className="h-2 w-36 overflow-hidden rounded-full bg-[var(--border-subtle)]">
                    <div
                      className="h-full bg-[var(--accent-primary)] transition-all"
                      style={{ width: `${progressPct}%` }}
                    />
                  </div>
                  <span className="text-xs font-medium text-[var(--text-muted)] tabular-nums">
                    {bulkProgress.completed}/{bulkProgress.total}
                  </span>
                </div>
              ) : bulkProgress.cancelled && bulkProgress.kind === "models" ? (
                <p role="status" className="text-xs font-medium text-[var(--text-muted)]">
                  {isId
                    ? `Pengujian dibatalkan (${bulkProgress.completed}/${bulkProgress.total})`
                    : `Tests cancelled (${bulkProgress.completed}/${bulkProgress.total})`}
                </p>
              ) : null}
            </div>

            {/* Models Table */}
            {catalogQuery.isError ? (
              <div className="p-5">
                <Alert variant="destructive">
                  <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
                  <div>
                    <AlertTitle>
                      {isId ? "Gagal memuat katalog model" : "Unable to load model catalog"}
                    </AlertTitle>
                    <AlertDescription className="font-mono">
                      {catalogQuery.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              </div>
            ) : catalogQuery.isPending ? (
              <div className="p-8 text-center text-sm text-[var(--text-muted)]">
                {isId ? "Memuat katalog model..." : "Loading model catalog..."}
              </div>
            ) : filteredModels.length === 0 ? (
              <div className="p-10 text-center">
                <p className="text-sm font-medium text-[var(--text-primary)]">
                  {models.length === 0
                    ? isId
                      ? "Belum ada model dari provider aktif"
                      : "No models from enabled providers yet"
                    : isId
                      ? "Tidak ada model yang cocok dengan filter"
                      : "No models match the current filters"}
                </p>
                <p className="mt-1 text-xs text-[var(--text-muted)]">
                  {models.length === 0
                    ? isId
                      ? "Hubungkan provider di halaman Providers untuk mengisi katalog model."
                      : "Connect a provider on the Providers page to populate the model catalog."
                    : isId
                      ? "Coba ubah atau bersihkan filter di atas."
                      : "Try adjusting or clearing the active filters above."}
                </p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full min-w-[960px] border-collapse text-left">
                  <thead>
                    <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-surface)]">
                      <th className="w-10 px-4 py-2.5">
                        <input
                          type="checkbox"
                          checked={allPagedModelsSelected}
                          onChange={toggleAllPagedModels}
                          aria-label={
                            isId ? "Pilih semua model di halaman ini" : "Select all models on page"
                          }
                          className="h-4 w-4 accent-[var(--accent-primary)]"
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="id"
                          label={isId ? "Model" : "Model"}
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="provider"
                          label="Provider"
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="type"
                          label={isId ? "Tipe" : "Type"}
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="context_length"
                          label={isId ? "Konteks" : "Context"}
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="max_output_tokens"
                          label={isId ? "Output" : "Output"}
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-3 py-2.5 text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                        {isId ? "Kapabilitas" : "Capabilities"}
                      </th>
                      <th className="px-3 py-2.5 text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                        {isId ? "Akses Kunci" : "Key Access"}
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="testResult"
                          label={isId ? "Hasil Uji" : "Test Result"}
                          activeField={modelSort.field}
                          direction={modelSort.direction}
                          onSort={handleSortModels}
                        />
                      </th>
                      <th className="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                        {isId ? "Aksi" : "Actions"}
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[var(--border-subtle)] text-xs">
                    {pagedModels.map((row) => {
                      const testKey = modelTestResultKey(row.providerId, row.id);
                      const result = testResults[testKey];
                      const isRowTesting = testingKeys.has(testKey);
                      const isSelected = selectedModelKeys.has(row.qualifiedId);
                      const caps = getEnabledCapabilities(row);
                      const provHealth = healthMatrix[row.providerId] ?? "untested";
                      const keyAccess = countKeyAccess(apiKeys, row.qualifiedId, row.id);

                      return (
                        <tr
                          key={row.qualifiedId}
                          className={`transition-colors hover:bg-[var(--bg-surface-raised)] ${
                            row.hidden ? "opacity-60" : ""
                          }`}
                        >
                          <td className="px-4 py-3">
                            <input
                              type="checkbox"
                              checked={isSelected}
                              onChange={() => {
                                setSelectedModelKeys((prev) => {
                                  const next = new Set(prev);
                                  if (next.has(row.qualifiedId)) next.delete(row.qualifiedId);
                                  else next.add(row.qualifiedId);
                                  return next;
                                });
                              }}
                              aria-label={`Select ${row.qualifiedId}`}
                              className="h-4 w-4 accent-[var(--accent-primary)]"
                            />
                          </td>

                          <td className="px-3 py-3">
                            <div className="flex items-center gap-1.5">
                              <span className="font-mono font-semibold text-[var(--text-primary)]">
                                {row.id}
                              </span>
                              <button
                                type="button"
                                onClick={() => void handleCopyModelId(row.qualifiedId)}
                                className="rounded p-0.5 text-[var(--text-muted)] hover:bg-[var(--bg-canvas)] hover:text-[var(--text-primary)]"
                                title={isId ? "Salin ID model" : "Copy model ID"}
                              >
                                {copiedId === row.qualifiedId ? (
                                  <Check className="h-3 w-3 text-[var(--status-ok)]" />
                                ) : (
                                  <Copy className="h-3 w-3" />
                                )}
                              </button>
                            </div>
                            <div className="mt-1 flex flex-wrap items-center gap-1">
                              {row.free ? (
                                <Badge variant="success" className="text-[10px]">
                                  {isId ? "Gratis" : "Free"}
                                </Badge>
                              ) : null}
                              {row.custom ? (
                                <Badge variant="outline" className="text-[10px]">
                                  {isId ? "Kustom" : "Custom"}
                                </Badge>
                              ) : null}
                              {row.hidden ? (
                                <Badge variant="warning" className="text-[10px]">
                                  {isId ? "Tersembunyi" : "Hidden"}
                                </Badge>
                              ) : null}
                            </div>
                          </td>

                          <td className="px-3 py-3">
                            <div className="font-medium text-[var(--text-primary)]">
                              {row.provider}
                            </div>
                            <div className="mt-0.5">
                              <HealthDot state={provHealth} isId={isId} />
                            </div>
                          </td>

                          <td className="px-3 py-3">
                            <Badge variant="outline" className="font-mono text-[10px]">
                              {row.type}
                            </Badge>
                            {row.subtype && row.subtype !== row.type ? (
                              <p className="mt-0.5 text-[10px] text-[var(--text-muted)]">
                                {row.subtype}
                              </p>
                            ) : null}
                          </td>

                          <td className="px-3 py-3 font-mono tabular-nums text-[var(--text-secondary)]">
                            {formatTokenCount(row.context_length)}
                          </td>

                          <td className="px-3 py-3 font-mono tabular-nums text-[var(--text-secondary)]">
                            {formatTokenCount(row.max_output_tokens)}
                          </td>

                          <td className="px-3 py-3">
                            {caps.length === 0 ? (
                              <span className="text-[var(--text-muted)]">—</span>
                            ) : (
                              <div className="flex flex-wrap gap-1">
                                {caps.slice(0, 3).map((cap) => (
                                  <Badge key={cap} variant="secondary" className="text-[10px]">
                                    {cap}
                                  </Badge>
                                ))}
                                {caps.length > 3 ? (
                                  <Badge variant="outline" className="text-[10px]">
                                    +{caps.length - 3}
                                  </Badge>
                                ) : null}
                              </div>
                            )}
                          </td>

                          <td className="px-3 py-3">
                            <button
                              type="button"
                              onClick={() => openAssignModal([row.qualifiedId])}
                              className="inline-flex items-center gap-1 rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 py-1 text-[11px] text-[var(--text-secondary)] transition-colors hover:border-[var(--accent-primary)] hover:text-[var(--text-primary)]"
                            >
                              <KeyRound className="h-3 w-3 text-[var(--accent-primary)]" />
                              <span className="tabular-nums">
                                {keyAccess.total === 0
                                  ? isId
                                    ? "0 kunci"
                                    : "0 keys"
                                  : `${keyAccess.allowed}/${keyAccess.total}`}
                              </span>
                            </button>
                          </td>

                          <td className="px-3 py-3">
                            <TestResultCell
                              result={result}
                              isTesting={isRowTesting}
                              nowMs={nowMs}
                              isId={isId}
                            />
                          </td>

                          <td className="px-4 py-3 text-right">
                            <div className="flex items-center justify-end gap-1.5">
                              <Button
                                type="button"
                                variant="outline"
                                size="xs"
                                data-testid={`test-model-${row.id}`}
                                disabled={isRowTesting || bulkRunning}
                                onClick={() => void handleTestSingleModel(row)}
                              >
                                <FlaskConical className="h-3 w-3" aria-hidden="true" />
                                <span>{isId ? "Uji" : "Test"}</span>
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                size="icon-sm"
                                disabled={visibilityMutation.isPending}
                                onClick={() => void handleToggleVisibility([row], !row.hidden)}
                                title={
                                  row.hidden
                                    ? isId
                                      ? "Tampilkan model"
                                      : "Show model"
                                    : isId
                                      ? "Sembunyikan model"
                                      : "Hide model"
                                }
                              >
                                {row.hidden ? (
                                  <Eye className="h-3.5 w-3.5 text-[var(--status-ok)]" />
                                ) : (
                                  <EyeOff className="h-3.5 w-3.5 text-[var(--text-muted)]" />
                                )}
                              </Button>
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {/* Pagination Footer */}
            {filteredModels.length > 0 ? (
              <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-2.5 text-xs text-[var(--text-muted)]">
                <span className="tabular-nums">
                  {isId
                    ? `Menampilkan ${(safeModelPage - 1) * PAGE_SIZE + 1}–${Math.min(safeModelPage * PAGE_SIZE, filteredModels.length)} dari ${filteredModels.length} model`
                    : `Showing ${(safeModelPage - 1) * PAGE_SIZE + 1}–${Math.min(safeModelPage * PAGE_SIZE, filteredModels.length)} of ${filteredModels.length} models`}
                </span>
                <div className="flex items-center gap-1.5">
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    disabled={safeModelPage <= 1}
                    onClick={() => setModelPage((p) => Math.max(1, p - 1))}
                  >
                    <ChevronLeft className="h-3.5 w-3.5" />
                    <span>{isId ? "Sebelumnya" : "Previous"}</span>
                  </Button>
                  <span className="px-2 font-mono tabular-nums">
                    {safeModelPage} / {totalModelPages}
                  </span>
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    disabled={safeModelPage >= totalModelPages}
                    onClick={() => setModelPage((p) => Math.min(totalModelPages, p + 1))}
                  >
                    <span>{isId ? "Berikutnya" : "Next"}</span>
                    <ChevronRight className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
            ) : null}
          </Card>
        </section>
      ) : (
        /* COMBOS TAB PANEL */
        <section
          id="panel-combos"
          role="tabpanel"
          aria-labelledby="tab-combos"
          className="space-y-4"
        >
          {/* Combo Filter Bar */}
          <Card className="p-4">
            <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2 lg:grid-cols-6">
              <div className="relative sm:col-span-2">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  type="text"
                  value={comboFilters.query}
                  onChange={(e) => setComboFilters((f) => ({ ...f, query: e.target.value }))}
                  placeholder={isId ? "Cari combo atau anggota..." : "Search combos or members..."}
                  className="h-9 pl-8 text-xs"
                />
              </div>

              <select
                aria-label={isId ? "Filter strategi" : "Filter strategy"}
                value={comboFilters.strategy}
                onChange={(e) => setComboFilters((f) => ({ ...f, strategy: e.target.value }))}
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Strategi" : "All Strategies"}</option>
                {strategyOptions.map((s) => (
                  <option key={s} value={s}>
                    {s}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter status" : "Filter status"}
                value={comboFilters.status}
                onChange={(e) =>
                  setComboFilters((f) => ({
                    ...f,
                    status: e.target.value as ComboCatalogFilters["status"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Status" : "All Statuses"}</option>
                <option value="active">{isId ? "Aktif" : "Active"}</option>
                <option value="paused">{isId ? "Dijeda" : "Paused"}</option>
              </select>

              <select
                aria-label={isId ? "Filter hasil uji combo" : "Filter combo test result"}
                value={comboFilters.testResult}
                onChange={(e) =>
                  setComboFilters((f) => ({
                    ...f,
                    testResult: e.target.value as ComboCatalogFilters["testResult"],
                  }))
                }
                className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Hasil Uji" : "All Results"}</option>
                <option value="ok">OK</option>
                <option value="slow">{isId ? "Lambat" : "Slow"}</option>
                <option value="error">{isId ? "Gagal" : "Error"}</option>
                <option value="untested">{isId ? "Belum diuji" : "Untested"}</option>
              </select>

              <div className="flex items-center gap-1.5">
                <Input
                  type="number"
                  min={0}
                  value={comboFilters.minMembers ?? ""}
                  onChange={(e) =>
                    setComboFilters((f) => ({
                      ...f,
                      minMembers: e.target.value ? Number(e.target.value) : undefined,
                    }))
                  }
                  placeholder={isId ? "Min" : "Min"}
                  className="h-9 text-xs"
                />
                <Input
                  type="number"
                  min={0}
                  value={comboFilters.maxMembers ?? ""}
                  onChange={(e) =>
                    setComboFilters((f) => ({
                      ...f,
                      maxMembers: e.target.value ? Number(e.target.value) : undefined,
                    }))
                  }
                  placeholder={isId ? "Maks" : "Max"}
                  className="h-9 text-xs"
                />
                {hasActiveComboFilters ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setComboFilters(DEFAULT_COMBO_FILTERS)}
                    className="shrink-0 text-xs"
                  >
                    {isId ? "Reset" : "Clear"}
                  </Button>
                ) : null}
              </div>
            </div>
          </Card>

          <Card className="overflow-hidden">
            {/* Combos Bulk Toolbar */}
            <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-2.5">
              <div className="flex flex-wrap items-center gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="assign-keys-btn"
                  disabled={selectedComboRows.length === 0 || bulkRunning}
                  onClick={() => openAssignModal(selectedComboRows.map((c) => c.name))}
                >
                  <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Tetapkan ke Kunci" : "Assign to Keys"}</span>
                </Button>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="test-selected-btn"
                  disabled={selectedComboRows.length === 0 || bulkRunning}
                  onClick={() => triggerBulkComboTests(selectedComboRows)}
                >
                  <FlaskConical className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {isId
                      ? `Uji Terpilih (${selectedComboRows.length})`
                      : `Test Selected (${selectedComboRows.length})`}
                  </span>
                </Button>

                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  data-testid="test-all-filtered-btn"
                  disabled={filteredCombos.length === 0 || bulkRunning}
                  onClick={() => triggerBulkComboTests(filteredCombos)}
                >
                  <span>
                    {isId
                      ? `Uji Semua Terfilter (${filteredCombos.length})`
                      : `Test All Filtered (${filteredCombos.length})`}
                  </span>
                </Button>

                {bulkRunning ? (
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    data-testid="cancel-tests-btn"
                    onClick={handleCancelBulkTests}
                  >
                    {isId ? "Batalkan" : "Cancel"}
                  </Button>
                ) : null}

                {hasTestResults && !bulkRunning ? (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    data-testid="clear-test-results-btn"
                    onClick={handleClearResults}
                  >
                    {isId ? "Bersihkan Hasil" : "Clear Results"}
                  </Button>
                ) : null}
              </div>

              {bulkRunning && bulkProgress.kind === "combos" ? (
                <div className="flex items-center gap-3">
                  <div className="h-2 w-36 overflow-hidden rounded-full bg-[var(--border-subtle)]">
                    <div
                      className="h-full bg-[var(--accent-primary)] transition-all"
                      style={{ width: `${progressPct}%` }}
                    />
                  </div>
                  <span className="text-xs font-medium text-[var(--text-muted)] tabular-nums">
                    {bulkProgress.completed}/{bulkProgress.total}
                  </span>
                </div>
              ) : null}
            </div>

            {combosQuery.isError ? (
              <div className="p-5">
                <Alert variant="destructive">
                  <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
                  <AlertDescription className="font-mono">
                    {combosQuery.error.message}
                  </AlertDescription>
                </Alert>
              </div>
            ) : combosQuery.isPending ? (
              <div className="p-8 text-center text-sm text-[var(--text-muted)]">
                {isId ? "Memuat combo..." : "Loading combos..."}
              </div>
            ) : filteredCombos.length === 0 ? (
              <div className="p-10 text-center">
                <p className="text-sm font-medium text-[var(--text-primary)]">
                  {combos.length === 0
                    ? isId
                      ? "Belum ada combo dikonfigurasi"
                      : "No combos configured yet"
                    : isId
                      ? "Tidak ada combo yang cocok dengan filter"
                      : "No combos match the current filters"}
                </p>
                <div className="mt-3">
                  <Link to="/dashboard/combos/">
                    <Button type="button" size="sm" variant="outline">
                      <Sparkles className="h-3.5 w-3.5" />
                      <span>{isId ? "Buka Pembuat Combo" : "Open Combo Builder"}</span>
                    </Button>
                  </Link>
                </div>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full min-w-[820px] border-collapse text-left">
                  <thead>
                    <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-surface)]">
                      <th className="w-10 px-4 py-2.5">
                        <input
                          type="checkbox"
                          checked={allPagedCombosSelected}
                          onChange={toggleAllPagedCombos}
                          aria-label={isId ? "Pilih semua combo" : "Select all combos"}
                          className="h-4 w-4 accent-[var(--accent-primary)]"
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="name"
                          label="Combo"
                          activeField={comboSort.field}
                          direction={comboSort.direction}
                          onSort={handleSortCombos}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="strategy"
                          label={isId ? "Strategi" : "Strategy"}
                          activeField={comboSort.field}
                          direction={comboSort.direction}
                          onSort={handleSortCombos}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="status"
                          label="Status"
                          activeField={comboSort.field}
                          direction={comboSort.direction}
                          onSort={handleSortCombos}
                        />
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="memberCount"
                          label={isId ? "Anggota" : "Members"}
                          activeField={comboSort.field}
                          direction={comboSort.direction}
                          onSort={handleSortCombos}
                        />
                      </th>
                      <th className="px-3 py-2.5 text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                        {isId ? "Akses Kunci" : "Key Access"}
                      </th>
                      <th className="px-3 py-2.5">
                        <SortHeaderButton
                          field="testResult"
                          label={isId ? "Hasil Uji" : "Test Result"}
                          activeField={comboSort.field}
                          direction={comboSort.direction}
                          onSort={handleSortCombos}
                        />
                      </th>
                      <th className="px-4 py-2.5 text-right text-[11px] font-semibold uppercase tracking-wider text-[var(--text-muted)]">
                        {isId ? "Aksi" : "Actions"}
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[var(--border-subtle)] text-xs">
                    {pagedCombos.map((row) => {
                      const testKey = comboTestResultKey(row.name);
                      const result = testResults[testKey];
                      const isRowTesting = testingKeys.has(testKey);
                      const isSelected = selectedComboNames.has(row.name);
                      const keyAccess = countKeyAccess(apiKeys, row.name);

                      return (
                        <tr
                          key={row.id}
                          className="transition-colors hover:bg-[var(--bg-surface-raised)]"
                        >
                          <td className="px-4 py-3">
                            <input
                              type="checkbox"
                              checked={isSelected}
                              onChange={() => {
                                setSelectedComboNames((prev) => {
                                  const next = new Set(prev);
                                  if (next.has(row.name)) next.delete(row.name);
                                  else next.add(row.name);
                                  return next;
                                });
                              }}
                              aria-label={`Select combo ${row.name}`}
                              className="h-4 w-4 accent-[var(--accent-primary)]"
                            />
                          </td>
                          <td className="px-3 py-3">
                            <div className="font-mono font-semibold text-[var(--text-primary)]">
                              {row.name}
                            </div>
                            {row.description ? (
                              <p className="mt-0.5 text-[11px] text-[var(--text-muted)]">
                                {row.description}
                              </p>
                            ) : null}
                          </td>
                          <td className="px-3 py-3">
                            <Badge variant="outline" className="font-mono text-[10px]">
                              {row.strategy}
                            </Badge>
                          </td>
                          <td className="px-3 py-3">
                            <Badge variant={row.status === "active" ? "success" : "warning"}>
                              {row.status === "active"
                                ? isId
                                  ? "Aktif"
                                  : "Active"
                                : isId
                                  ? "Dijeda"
                                  : "Paused"}
                            </Badge>
                          </td>
                          <td className="px-3 py-3">
                            <div className="flex flex-wrap gap-1">
                              {row.models.slice(0, 3).map((m, idx) => (
                                <Badge
                                  key={m.model}
                                  variant="secondary"
                                  className="font-mono text-[10px]"
                                >
                                  {idx + 1}. {m.model}
                                </Badge>
                              ))}
                              {row.models.length > 3 ? (
                                <Badge variant="outline" className="text-[10px]">
                                  +{row.models.length - 3} {isId ? "lainnya" : "more"}
                                </Badge>
                              ) : null}
                            </div>
                          </td>
                          <td className="px-3 py-3">
                            <button
                              type="button"
                              onClick={() => openAssignModal([row.name])}
                              className="inline-flex items-center gap-1 rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 py-1 text-[11px] text-[var(--text-secondary)] hover:border-[var(--accent-primary)] hover:text-[var(--text-primary)]"
                            >
                              <KeyRound className="h-3 w-3 text-[var(--accent-primary)]" />
                              <span className="tabular-nums">
                                {keyAccess.allowed}/{keyAccess.total}
                              </span>
                            </button>
                          </td>
                          <td className="px-3 py-3">
                            <TestResultCell
                              result={result}
                              isTesting={isRowTesting}
                              nowMs={nowMs}
                              isId={isId}
                            />
                          </td>
                          <td className="px-4 py-3 text-right">
                            <Button
                              type="button"
                              variant="outline"
                              size="xs"
                              disabled={isRowTesting || bulkRunning}
                              onClick={() => void handleTestSingleCombo(row)}
                            >
                              <FlaskConical className="h-3 w-3" aria-hidden="true" />
                              <span>{isId ? "Uji" : "Test"}</span>
                            </Button>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {filteredCombos.length > 0 ? (
              <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-2.5 text-xs text-[var(--text-muted)]">
                <span className="tabular-nums">
                  {isId
                    ? `Menampilkan ${(safeComboPage - 1) * PAGE_SIZE + 1}–${Math.min(safeComboPage * PAGE_SIZE, filteredCombos.length)} dari ${filteredCombos.length} combo`
                    : `Showing ${(safeComboPage - 1) * PAGE_SIZE + 1}–${Math.min(safeComboPage * PAGE_SIZE, filteredCombos.length)} of ${filteredCombos.length} combos`}
                </span>
                <div className="flex items-center gap-1.5">
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    disabled={safeComboPage <= 1}
                    onClick={() => setComboPage((p) => Math.max(1, p - 1))}
                  >
                    <ChevronLeft className="h-3.5 w-3.5" />
                    <span>{isId ? "Sebelumnya" : "Previous"}</span>
                  </Button>
                  <span className="px-2 font-mono tabular-nums">
                    {safeComboPage} / {totalComboPages}
                  </span>
                  <Button
                    type="button"
                    variant="outline"
                    size="xs"
                    disabled={safeComboPage >= totalComboPages}
                    onClick={() => setComboPage((p) => Math.min(totalComboPages, p + 1))}
                  >
                    <span>{isId ? "Berikutnya" : "Next"}</span>
                    <ChevronRight className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
            ) : null}
          </Card>
        </section>
      )}

      {/* Bulk Test Confirmation Modal */}
      {bulkConfirmPending ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4"
          role="dialog"
          aria-modal="true"
        >
          <div className="w-full max-w-md rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] p-5 shadow-2xl">
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">
              {isId ? "Konfirmasi Pengujian Massal" : "Confirm Bulk Testing"}
            </h3>
            <p className="mt-2 text-xs text-[var(--text-secondary)]">
              {bulkConfirmPending.kind === "models"
                ? isId
                  ? `Anda akan menguji ${bulkConfirmPending.count} model secara langsung ke upstream. Ini dapat menggunakan kuota token minimal.`
                  : `You are about to live-probe ${bulkConfirmPending.count} models against upstream providers. This may consume minimal token quota.`
                : isId
                  ? `Anda akan menguji ${bulkConfirmPending.count} combo secara langsung ke upstream.`
                  : `You are about to live-probe ${bulkConfirmPending.count} combos against upstream providers.`}
            </p>
            <div className="mt-4 flex justify-end gap-2">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => setBulkConfirmPending(null)}
              >
                {isId ? "Batal" : "Cancel"}
              </Button>
              <Button type="button" size="sm" onClick={bulkConfirmPending.onConfirm}>
                {isId ? "Mulai Pengujian" : "Start Testing"}
              </Button>
            </div>
          </div>
        </div>
      ) : null}

      {/* Assign to Keys Modal */}
      <AssignKeysModal
        open={assignModalOpen}
        onClose={() => setAssignModalOpen(false)}
        targets={assignTargets}
        keys={apiKeys}
        keysLoading={keysQuery.isPending}
        keysError={keysQuery.error?.message ?? null}
        onRetryKeys={() => void keysQuery.refetch()}
        isId={isId}
      />
    </div>
  );
}
