import { AlertCircle, BookOpen, Layers, Plus, RefreshCw, Search, X } from "lucide-react";
import { lazy, Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { runSingleComboTest, useModelCatalog } from "@/features/models/hooks";
import type { CatalogTestResult } from "@/features/models/types";
import { useProviderConnections, useProvidersCatalog } from "@/features/providers/hooks";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { ComboCard } from "./ComboCard";
import { ComboRegistrationForm } from "./ComboRegistrationForm";
import {
  useCombos,
  useCreateComboMutation,
  useDeleteComboMutation,
  useUpdateComboMutation,
} from "./hooks";
import type { Combo, ComboStrategy, CreateComboRequest } from "./schemas";
import { COMBO_STRATEGY_OPTIONS } from "./strategies";

const StrategyGuideModal = lazy(() =>
  import("./StrategyGuideModal").then((m) => ({ default: m.StrategyGuideModal })),
);

export function CombosConsole() {
  const locale = useUIStore((s) => s.locale);
  const isId = locale === "id";

  const combosQuery = useCombos();
  const connectionsQuery = useProviderConnections();
  const catalogQuery = useProvidersCatalog();
  const modelCatalogQuery = useModelCatalog();

  const createMutation = useCreateComboMutation();
  const updateMutation = useUpdateComboMutation();
  const deleteMutation = useDeleteComboMutation();

  const [query, setQuery] = useState("");
  const [strategyFilter, setStrategyFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "paused">("all");
  const [mobilePane, setMobilePane] = useState<"list" | "form">("list");

  // Inline registration form state
  const [formMode, setFormMode] = useState<"create" | "edit">("create");
  const [editingCombo, setEditingCombo] = useState<Combo | null>(null);
  const [formSubmitError, setFormSubmitError] = useState<string | null>(null);
  const [strategyOverride, setStrategyOverride] = useState<ComboStrategy | null>(null);

  // Strategy guide helper modal
  const [guideOpen, setGuideOpen] = useState(false);
  const [guideStrategy, setGuideStrategy] = useState<ComboStrategy>("priority");

  const [busyComboId, setBusyComboId] = useState<string | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [copiedName, setCopiedName] = useState<string | null>(null);
  const copyTimerRef = useRef<number | null>(null);

  const [testingIds, setTestingIds] = useState<Set<string>>(() => new Set());
  const [testResults, setTestResults] = useState<Record<string, CatalogTestResult>>({});
  const testControllersRef = useRef<Map<string, AbortController>>(new Map());

  useEffect(() => {
    const controllers = testControllersRef.current;
    return () => {
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
      for (const controller of controllers.values()) {
        controller.abort();
      }
      controllers.clear();
    };
  }, []);

  const combos = combosQuery.data ?? [];
  const connections = connectionsQuery.data ?? [];

  const connectionsById = useMemo(
    () => new Map(connections.map((conn) => [conn.id, conn])),
    [connections],
  );

  const providerLabels = useMemo(() => {
    const labels: Record<string, string> = {};
    for (const item of catalogQuery.data?.adapters ?? []) {
      labels[item.id] = item.display_name;
    }
    return labels;
  }, [catalogQuery.data]);

  const providerCatalogModels = useMemo(() => {
    const grouped: Record<string, string[]> = {};
    const seenByProvider = new Map<string, Set<string>>();

    const add = (providerId: string, modelId: string) => {
      const cleanProvider = providerId.trim();
      const cleanModel = modelId.trim();
      if (!cleanProvider || !cleanModel) return;
      let seen = seenByProvider.get(cleanProvider);
      if (!seen) {
        seen = new Set<string>();
        seenByProvider.set(cleanProvider, seen);
        grouped[cleanProvider] = [];
      }
      if (!seen.has(cleanModel)) {
        seen.add(cleanModel);
        grouped[cleanProvider]?.push(cleanModel);
      }
    };

    for (const row of modelCatalogQuery.data ?? []) {
      if (row.hidden) continue;
      add(row.providerId, row.id);
    }
    for (const adapter of catalogQuery.data?.adapters ?? []) {
      for (const modelId of adapter.default_models ?? []) {
        add(adapter.id, modelId);
      }
    }

    return grouped;
  }, [modelCatalogQuery.data, catalogQuery.data]);

  const filteredCombos = useMemo(() => {
    const q = query.trim().toLowerCase();
    return combos.filter((c) => {
      if (q) {
        const hay = `${c.name} ${c.description} ${c.strategy} ${c.targets.join(" ")}`.toLowerCase();
        if (!hay.includes(q)) return false;
      }
      if (strategyFilter !== "all" && c.strategy !== strategyFilter) {
        return false;
      }
      if (statusFilter === "active" && !c.enabled) return false;
      if (statusFilter === "paused" && c.enabled) return false;
      return true;
    });
  }, [combos, query, strategyFilter, statusFilter]);

  const activeCount = useMemo(() => combos.filter((c) => c.enabled).length, [combos]);

  const clearTestResult = useCallback((comboId: string) => {
    setTestResults((prev) => {
      if (!(comboId in prev)) return prev;
      const next = { ...prev };
      delete next[comboId];
      return next;
    });
  }, []);

  const handleEditCombo = useCallback((combo: Combo) => {
    setFormMode("edit");
    setEditingCombo(combo);
    setFormSubmitError(null);
    setMobilePane("form");
  }, []);

  const handleDuplicateCombo = useCallback((combo: Combo) => {
    setFormMode("create");
    setEditingCombo(combo);
    setFormSubmitError(null);
    setMobilePane("form");
  }, []);

  const handleCancelEdit = useCallback(() => {
    setFormMode("create");
    setEditingCombo(null);
    setFormSubmitError(null);
    setMobilePane("list");
  }, []);

  const handleFormSubmit = useCallback(
    (payload: CreateComboRequest, editComboId?: string, onDone?: () => void) => {
      setFormSubmitError(null);
      if (editComboId) {
        updateMutation.mutate(
          { id: editComboId, body: payload },
          {
            onSuccess: () => {
              clearTestResult(editComboId);
              setFormMode("create");
              setEditingCombo(null);
              setMobilePane("list");
              onDone?.();
            },
            onError: (err) => setFormSubmitError(err.message),
          },
        );
        return;
      }
      createMutation.mutate(payload, {
        onSuccess: () => {
          setEditingCombo(null);
          setMobilePane("list");
          onDone?.();
        },
        onError: (err) => setFormSubmitError(err.message),
      });
    },
    [clearTestResult, createMutation, updateMutation],
  );

  const handleToggleEnabled = useCallback(
    (combo: Combo) => {
      if (busyComboId) return;
      setRowError(null);
      setBusyComboId(combo.id);
      updateMutation.mutate(
        { id: combo.id, body: { enabled: !combo.enabled } },
        {
          onError: (err) => setRowError(err.message),
          onSettled: () => setBusyComboId(null),
        },
      );
    },
    [busyComboId, updateMutation],
  );

  const handleDelete = useCallback(
    (combo: Combo) => {
      if (busyComboId) return;
      setRowError(null);
      setBusyComboId(combo.id);
      deleteMutation.mutate(combo.id, {
        onSuccess: () => {
          setConfirmDeleteId(null);
          clearTestResult(combo.id);
          if (editingCombo?.id === combo.id) {
            handleCancelEdit();
          }
        },
        onError: (err) => setRowError(err.message),
        onSettled: () => setBusyComboId(null),
      });
    },
    [busyComboId, clearTestResult, deleteMutation, editingCombo?.id, handleCancelEdit],
  );

  const handleTestCombo = useCallback(async (combo: Combo) => {
    if (testControllersRef.current.has(combo.id)) return;
    const controller = new AbortController();
    testControllersRef.current.set(combo.id, controller);
    setTestingIds((prev) => new Set(prev).add(combo.id));
    try {
      const res = await runSingleComboTest(combo.name, controller.signal);
      if (!controller.signal.aborted) {
        setTestResults((prev) => ({ ...prev, [combo.id]: res }));
      }
    } finally {
      testControllersRef.current.delete(combo.id);
      setTestingIds((prev) => {
        const next = new Set(prev);
        next.delete(combo.id);
        return next;
      });
    }
  }, []);

  const handleCopyName = useCallback(async (comboName: string) => {
    try {
      await navigator.clipboard.writeText(comboName);
      setCopiedName(comboName);
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
      copyTimerRef.current = window.setTimeout(() => {
        setCopiedName((cur) => (cur === comboName ? null : cur));
      }, 1800);
    } catch {
      // Ignore clipboard permission errors
    }
  }, []);

  const openStrategyGuide = useCallback((strategy: ComboStrategy) => {
    setGuideStrategy(strategy);
    setGuideOpen(true);
  }, []);

  const handleStrategyOverrideApplied = useCallback(() => {
    setStrategyOverride(null);
  }, []);

  const handleRequestDelete = useCallback((comboId: string) => {
    setConfirmDeleteId(comboId);
  }, []);

  const handleCancelDelete = useCallback(() => {
    setConfirmDeleteId(null);
  }, []);

  const isRefreshing =
    combosQuery.isFetching ||
    connectionsQuery.isFetching ||
    catalogQuery.isFetching ||
    modelCatalogQuery.isFetching;

  return (
    <div className="flex flex-col gap-3 lg:h-[calc(100vh-7.25rem)] lg:min-h-[440px] lg:overflow-hidden lg:supports-[height:100dvh]:h-[calc(100dvh-7.25rem)]">
      {/* Compact Single-Row Top Bar */}
      <section className="flex shrink-0 items-center justify-between gap-2">
        {/* Mobile-first segmented switcher (hidden on lg+ where both panes sit side-by-side) */}
        <div
          className="inline-flex rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-0.5 lg:hidden"
          role="tablist"
          aria-label={isId ? "Tampilan Combo" : "Combo workspace view"}
        >
          <button
            type="button"
            role="tab"
            aria-selected={mobilePane === "list"}
            onClick={() => setMobilePane("list")}
            className={`inline-flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
              mobilePane === "list"
                ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
            }`}
          >
            <Layers className="h-3.5 w-3.5" aria-hidden="true" />
            <span>
              {isId ? "Daftar" : "Combos"} ({combos.length})
            </span>
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={mobilePane === "form"}
            onClick={() => setMobilePane("form")}
            className={`inline-flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-medium transition-colors ${
              mobilePane === "form"
                ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
            }`}
          >
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{formMode === "edit" ? "Edit" : isId ? "Baru" : "New"}</span>
          </button>
        </div>

        {/* Desktop compact summary */}
        <div className="hidden min-w-0 items-center gap-2 lg:flex">
          <Badge variant="outline" className="tabular-nums">
            {activeCount}/{combos.length} {isId ? "Aktif" : "Active"}
          </Badge>
          <span className="truncate text-xs text-[var(--text-muted)]">
            {isId
              ? "Gabungkan model lintas provider dalam satu endpoint."
              : "Combine multi-provider models behind one endpoint."}
          </span>
        </div>

        <div className="flex shrink-0 items-center gap-1.5">
          <Badge variant="outline" className="tabular-nums lg:hidden">
            {activeCount}/{combos.length}
          </Badge>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 px-2.5"
            onClick={() => openStrategyGuide(guideStrategy)}
            title={isId ? "Panduan Strategy" : "Strategy Guide"}
          >
            <BookOpen className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
            <span className="hidden sm:inline">{isId ? "Panduan" : "Guide"}</span>
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 px-2.5"
            onClick={() => {
              void combosQuery.refetch();
              void connectionsQuery.refetch();
              void catalogQuery.refetch();
              void modelCatalogQuery.refetch();
            }}
            disabled={isRefreshing}
            aria-label={isId ? "Segarkan" : "Refresh"}
            title={isId ? "Segarkan" : "Refresh"}
          >
            <RefreshCw
              className={`h-3.5 w-3.5 text-[var(--accent-primary)] ${isRefreshing ? "animate-spin" : ""}`}
              aria-hidden="true"
            />
            <span className="hidden sm:inline">{isId ? "Segarkan" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {/* Main 2-Pane Workspace: Form (Left) + Scrollable Combos List (Right) */}
      <div className="grid min-h-0 flex-1 grid-cols-1 gap-4 lg:grid-cols-12 lg:overflow-hidden">
        <div
          className={`min-h-0 lg:col-span-5 lg:block lg:overflow-hidden ${
            mobilePane === "form" ? "block" : "hidden"
          }`}
        >
          <ComboRegistrationForm
            mode={formMode}
            initialCombo={editingCombo}
            existingCombos={combos}
            connections={connections}
            providerLabels={providerLabels}
            providerCatalogModels={providerCatalogModels}
            modelsLoading={modelCatalogQuery.isPending}
            submitting={createMutation.isPending || updateMutation.isPending}
            submitError={formSubmitError}
            strategyOverride={strategyOverride}
            onStrategyOverrideApplied={handleStrategyOverrideApplied}
            onOpenStrategyGuide={openStrategyGuide}
            onCancelEdit={handleCancelEdit}
            onSubmit={handleFormSubmit}
            isId={isId}
          />
        </div>

        <div
          className={`min-h-0 flex-col gap-2.5 lg:col-span-7 lg:flex lg:overflow-hidden ${
            mobilePane === "list" ? "flex" : "hidden"
          }`}
        >
          {/* Compact Filter Bar */}
          <Card className="shrink-0 p-2.5">
            <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
              <div className="relative col-span-2">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  type="text"
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder={
                    isId ? "Cari combo, strategi, model..." : "Search combo, strategy, model..."
                  }
                  className="h-8 pl-8 text-xs"
                />
              </div>

              <select
                aria-label={isId ? "Filter strategi" : "Filter strategy"}
                value={strategyFilter}
                onChange={(e) => setStrategyFilter(e.target.value)}
                className="h-8 rounded-md border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2 font-mono text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Strategi" : "All Strategies"}</option>
                {COMBO_STRATEGY_OPTIONS.map((s) => (
                  <option key={s.value} value={s.value}>
                    {s.value}
                  </option>
                ))}
              </select>

              <select
                aria-label={isId ? "Filter status" : "Filter status"}
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value as "all" | "active" | "paused")}
                className="h-8 rounded-md border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Status" : "All Status"}</option>
                <option value="active">{isId ? "Aktif" : "Active"}</option>
                <option value="paused">{isId ? "Dijeda" : "Paused"}</option>
              </select>
            </div>
          </Card>

          {rowError ? (
            <Alert variant="destructive" className="shrink-0">
              <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
              <div className="flex-1">
                <AlertTitle>{isId ? "Operasi combo gagal" : "Combo operation failed"}</AlertTitle>
                <AlertDescription className="font-mono">{rowError}</AlertDescription>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setRowError(null)}
              >
                <X className="h-3.5 w-3.5" />
              </Button>
            </Alert>
          ) : null}

          {/* Scrollable Registered Combos List */}
          <div className="custom-scrollbar min-h-0 flex-1 pr-0.5 lg:overflow-y-auto">
            {combosQuery.isError ? (
              <Alert variant="destructive">
                <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
                <AlertDescription className="font-mono">
                  {combosQuery.error.message}
                </AlertDescription>
              </Alert>
            ) : combosQuery.isPending ? (
              <Card className="p-8 text-center text-sm text-[var(--text-muted)]">
                {isId ? "Memuat daftar combo..." : "Loading combos..."}
              </Card>
            ) : filteredCombos.length === 0 ? (
              <Card className="p-8 text-center">
                <Layers className="mx-auto h-7 w-7 text-[var(--text-muted)]" aria-hidden="true" />
                <p className="mt-2.5 text-sm font-medium text-[var(--text-primary)]">
                  {combos.length === 0
                    ? isId
                      ? "Belum ada combo terdaftar"
                      : "No combos registered yet"
                    : isId
                      ? "Tidak ada combo yang cocok dengan filter"
                      : "No combos match the current filters"}
                </p>
                <p className="mt-1 text-xs text-[var(--text-muted)]">
                  {isId
                    ? "Buat combo baru untuk menggabungkan beberapa model."
                    : "Create a new combo to route across multiple models."}
                </p>
                {combos.length === 0 ? (
                  <Button
                    type="button"
                    size="sm"
                    className="mt-3 lg:hidden"
                    onClick={() => setMobilePane("form")}
                  >
                    <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Buat Combo" : "New Combo"}</span>
                  </Button>
                ) : null}
              </Card>
            ) : (
              <div className="space-y-2.5 pb-1">
                {filteredCombos.map((combo) => (
                  <ComboCard
                    key={combo.id}
                    combo={combo}
                    connectionsById={connectionsById}
                    testResult={testResults[combo.id]}
                    isTesting={testingIds.has(combo.id)}
                    isBusy={busyComboId === combo.id}
                    isCopied={copiedName === combo.name}
                    isConfirmingDelete={confirmDeleteId === combo.id}
                    isCurrentlyEditing={formMode === "edit" && editingCombo?.id === combo.id}
                    isId={isId}
                    onCopyName={handleCopyName}
                    onOpenStrategyGuide={openStrategyGuide}
                    onTestCombo={handleTestCombo}
                    onEditCombo={handleEditCombo}
                    onDuplicateCombo={handleDuplicateCombo}
                    onToggleEnabled={handleToggleEnabled}
                    onRequestDelete={handleRequestDelete}
                    onCancelDelete={handleCancelDelete}
                    onConfirmDelete={handleDelete}
                  />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {guideOpen ? (
        <Suspense fallback={null}>
          <StrategyGuideModal
            open={guideOpen}
            selectedStrategy={guideStrategy}
            onSelectStrategy={(chosen) => setStrategyOverride(chosen)}
            onClose={() => setGuideOpen(false)}
            isId={isId}
          />
        </Suspense>
      ) : null}
    </div>
  );
}
