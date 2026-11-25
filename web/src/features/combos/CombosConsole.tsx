import {
  AlertCircle,
  BookOpen,
  Check,
  CheckCircle2,
  Copy,
  CopyPlus,
  FlaskConical,
  KeyRound,
  Layers,
  Pencil,
  Power,
  RefreshCw,
  Search,
  Sparkles,
  Trash2,
  X,
  XCircle,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { runSingleComboTest, useModelCatalog } from "@/features/models/hooks";
import type { CatalogTestResult } from "@/features/models/types";
import { useProviderConnections, useProvidersCatalog } from "@/features/providers/hooks";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { ComboRegistrationForm } from "./ComboRegistrationForm";
import { formatStepAccountLabel } from "./comboHelpers";
import {
  useCombos,
  useCreateComboMutation,
  useDeleteComboMutation,
  useUpdateComboMutation,
} from "./hooks";
import { StrategyGuideModal } from "./StrategyGuideModal";
import type { Combo, ComboStep, ComboStrategy, CreateComboRequest } from "./schemas";
import { COMBO_STRATEGY_OPTIONS, getStrategyMeta } from "./strategies";

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

  // Inline registration form state
  const [formMode, setFormMode] = useState<"create" | "edit">("create");
  const [editingCombo, setEditingCombo] = useState<Combo | null>(null);
  const [formSubmitError, setFormSubmitError] = useState<string | null>(null);
  const [strategyOverride, setStrategyOverride] = useState<ComboStrategy | null>(null);
  const formContainerRef = useRef<HTMLDivElement | null>(null);

  // Strategy guide helper modal
  const [guideOpen, setGuideOpen] = useState(false);
  const [guideStrategy, setGuideStrategy] = useState<ComboStrategy>("priority");

  const [busyComboId, setBusyComboId] = useState<string | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [copiedName, setCopiedName] = useState<string | null>(null);
  const copyTimerRef = useRef<number | null>(null);

  const [testingNames, setTestingNames] = useState<Set<string>>(() => new Set());
  const [testResults, setTestResults] = useState<Record<string, CatalogTestResult>>({});

  useEffect(() => {
    return () => {
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
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

  const scrollToForm = () => {
    formContainerRef.current?.scrollIntoView({ behavior: "smooth", block: "start" });
  };

  const handleEditCombo = (combo: Combo) => {
    setFormMode("edit");
    setEditingCombo(combo);
    setFormSubmitError(null);
    scrollToForm();
  };

  const handleDuplicateCombo = (combo: Combo) => {
    setFormMode("create");
    setEditingCombo(combo);
    setFormSubmitError(null);
    scrollToForm();
  };

  const handleCancelEdit = () => {
    setFormMode("create");
    setEditingCombo(null);
    setFormSubmitError(null);
  };

  const handleFormSubmit = (
    payload: CreateComboRequest,
    editComboId?: string,
    onDone?: () => void,
  ) => {
    setFormSubmitError(null);
    if (editComboId) {
      updateMutation.mutate(
        { id: editComboId, body: payload },
        {
          onSuccess: () => {
            setFormMode("create");
            setEditingCombo(null);
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
        onDone?.();
      },
      onError: (err) => setFormSubmitError(err.message),
    });
  };

  const handleToggleEnabled = (combo: Combo) => {
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
  };

  const handleDelete = (combo: Combo) => {
    if (busyComboId) return;
    setRowError(null);
    setBusyComboId(combo.id);
    deleteMutation.mutate(combo.id, {
      onSuccess: () => {
        setConfirmDeleteId(null);
        if (editingCombo?.id === combo.id) {
          handleCancelEdit();
        }
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyComboId(null),
    });
  };

  const handleTestCombo = async (combo: Combo) => {
    if (testingNames.has(combo.name)) return;
    setTestingNames((prev) => new Set(prev).add(combo.name));
    try {
      const res = await runSingleComboTest(combo.name);
      setTestResults((prev) => ({ ...prev, [combo.name]: res }));
    } finally {
      setTestingNames((prev) => {
        const next = new Set(prev);
        next.delete(combo.name);
        return next;
      });
    }
  };

  const handleCopyName = async (comboName: string) => {
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
  };

  const openStrategyGuide = useCallback((strategy: ComboStrategy) => {
    setGuideStrategy(strategy);
    setGuideOpen(true);
  }, []);

  const handleStrategyOverrideApplied = useCallback(() => {
    setStrategyOverride(null);
  }, []);

  const isRefreshing =
    combosQuery.isFetching || connectionsQuery.isFetching || modelCatalogQuery.isFetching;

  return (
    <div className="space-y-5">
      {/* Header */}
      <section className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            {isId ? "Routing Combo" : "Combo Routing"}
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            {isId ? "Combo" : "Combos"}
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            {isId
              ? "Gabungkan model lintas provider dalam satu endpoint, pilih akun tertentu atau Dynamic, dan atur strategi routingnya."
              : "Combine models across providers into one endpoint, pin specific accounts or stay Dynamic, and pick your routing strategy."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="outline" className="tabular-nums">
            {activeCount}/{combos.length} {isId ? "Aktif" : "Active"}
          </Badge>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => openStrategyGuide(guideStrategy)}
          >
            <BookOpen className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
            <span>{isId ? "Panduan Strategy" : "Strategy Guide"}</span>
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => {
              void combosQuery.refetch();
              void connectionsQuery.refetch();
              void modelCatalogQuery.refetch();
            }}
            disabled={isRefreshing}
          >
            <RefreshCw
              className={`h-3.5 w-3.5 text-[var(--accent-primary)] ${isRefreshing ? "animate-spin" : ""}`}
              aria-hidden="true"
            />
            <span>{isId ? "Segarkan" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {/* Main 2-Column Workspace: Registration Form (Left) + Registered Combos (Right) */}
      <div className="grid grid-cols-1 items-start gap-5 xl:grid-cols-12">
        <div ref={formContainerRef} className="xl:sticky xl:top-4 xl:col-span-5">
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

        <div className="space-y-4 xl:col-span-7">
          {/* Filter Bar */}
          <Card className="p-3.5">
            <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-4">
              <div className="relative sm:col-span-2">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  type="text"
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder={
                    isId
                      ? "Cari nama combo, strategi, atau model..."
                      : "Search combo name, strategy, or model..."
                  }
                  className="h-9 pl-8 text-xs"
                />
              </div>

              <select
                aria-label={isId ? "Filter strategi" : "Filter strategy"}
                value={strategyFilter}
                onChange={(e) => setStrategyFilter(e.target.value)}
                className="h-9 rounded-md border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
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
                className="h-9 rounded-md border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
              >
                <option value="all">{isId ? "Semua Status" : "All Statuses"}</option>
                <option value="active">{isId ? "Aktif" : "Active"}</option>
                <option value="paused">{isId ? "Dijeda" : "Paused"}</option>
              </select>
            </div>
          </Card>

          {rowError ? (
            <Alert variant="destructive">
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

          {/* Registered Combos List */}
          {combosQuery.isError ? (
            <Alert variant="destructive">
              <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
              <AlertDescription className="font-mono">{combosQuery.error.message}</AlertDescription>
            </Alert>
          ) : combosQuery.isPending ? (
            <Card className="p-10 text-center text-sm text-[var(--text-muted)]">
              {isId ? "Memuat daftar combo..." : "Loading combos..."}
            </Card>
          ) : filteredCombos.length === 0 ? (
            <Card className="p-10 text-center">
              <Layers className="mx-auto h-8 w-8 text-[var(--text-muted)]" aria-hidden="true" />
              <p className="mt-3 text-sm font-medium text-[var(--text-primary)]">
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
                  ? "Isi form registrasi di sebelah kiri untuk membuat combo pertamamu."
                  : "Use the registration form on the left to create your first combo."}
              </p>
            </Card>
          ) : (
            <div className="space-y-3">
              {filteredCombos.map((combo) => {
                const stratMeta = getStrategyMeta(combo.strategy);
                const steps: ComboStep[] =
                  combo.steps && combo.steps.length > 0
                    ? combo.steps
                    : combo.targets.map((t): ComboStep => ({ kind: "model", model: t }));
                const testRes = testResults[combo.name];
                const isTesting = testingNames.has(combo.name);
                const isBusy = busyComboId === combo.id;
                const confirmingDelete = confirmDeleteId === combo.id;
                const isCurrentlyEditing = formMode === "edit" && editingCombo?.id === combo.id;

                return (
                  <Card
                    key={combo.id}
                    className={`overflow-hidden transition-colors ${
                      isCurrentlyEditing ? "border-[var(--accent-primary)]" : ""
                    }`}
                  >
                    <div className="h-[2px] w-full bg-[var(--border-subtle)]">
                      <div
                        className={`h-full transition-all ${
                          combo.enabled
                            ? "w-full bg-[var(--accent-primary)]"
                            : "w-1/3 bg-[var(--status-warn)]"
                        }`}
                      />
                    </div>

                    <div className="space-y-3 p-4">
                      {/* Top Row: Name + Badges + Actions */}
                      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                        <div className="min-w-0 flex-1 space-y-1">
                          <div className="flex flex-wrap items-center gap-1.5">
                            <span className="font-mono text-sm font-semibold text-[var(--text-primary)]">
                              {combo.name}
                            </span>
                            <button
                              type="button"
                              onClick={() => void handleCopyName(combo.name)}
                              className="rounded p-1 text-[var(--text-muted)] hover:bg-[var(--bg-canvas)] hover:text-[var(--text-primary)]"
                              title={isId ? "Salin nama combo" : "Copy combo name"}
                            >
                              {copiedName === combo.name ? (
                                <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" />
                              ) : (
                                <Copy className="h-3.5 w-3.5" />
                              )}
                            </button>
                            <Badge variant={combo.enabled ? "success" : "warning"}>
                              {combo.enabled
                                ? isId
                                  ? "Aktif"
                                  : "Active"
                                : isId
                                  ? "Dijeda"
                                  : "Paused"}
                            </Badge>
                            <button
                              type="button"
                              onClick={() => openStrategyGuide(combo.strategy)}
                              title={
                                isId
                                  ? "Klik untuk lihat penjelasan strategi ini"
                                  : "Click to view guide for this strategy"
                              }
                              className="inline-flex items-center gap-1 rounded-full border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2 py-0.5 font-mono text-[11px] text-[var(--text-primary)] transition-colors hover:border-[var(--accent-primary)] hover:text-[var(--accent-text)]"
                            >
                              <span>{combo.strategy}</span>
                              <span className="font-sans text-[10px] text-[var(--text-muted)]">
                                · {isId ? stratMeta.labelId : stratMeta.labelEn}
                              </span>
                            </button>
                          </div>

                          {combo.description ? (
                            <p className="text-xs text-[var(--text-secondary)]">
                              {combo.description}
                            </p>
                          ) : null}
                        </div>

                        {/* Actions */}
                        <div className="flex shrink-0 flex-wrap items-center gap-1">
                          {confirmingDelete ? (
                            <>
                              <Button
                                type="button"
                                variant="destructive"
                                size="xs"
                                disabled={isBusy}
                                onClick={() => handleDelete(combo)}
                              >
                                {isBusy
                                  ? isId
                                    ? "Menghapus..."
                                    : "Deleting..."
                                  : isId
                                    ? "Konfirmasi Hapus"
                                    : "Confirm Delete"}
                              </Button>
                              <Button
                                type="button"
                                variant="secondary"
                                size="xs"
                                disabled={isBusy}
                                onClick={() => setConfirmDeleteId(null)}
                              >
                                {isId ? "Batal" : "Cancel"}
                              </Button>
                            </>
                          ) : (
                            <>
                              <Button
                                type="button"
                                variant="outline"
                                size="xs"
                                disabled={isTesting}
                                onClick={() => void handleTestCombo(combo)}
                              >
                                <FlaskConical className="h-3 w-3" />
                                <span>
                                  {isTesting
                                    ? isId
                                      ? "Menguji..."
                                      : "Testing..."
                                    : isId
                                      ? "Uji"
                                      : "Test"}
                                </span>
                              </Button>
                              <Button
                                type="button"
                                variant={isCurrentlyEditing ? "default" : "outline"}
                                size="xs"
                                disabled={isBusy}
                                onClick={() => handleEditCombo(combo)}
                              >
                                <Pencil className="h-3 w-3" />
                                <span>Edit</span>
                              </Button>
                              <Button
                                type="button"
                                variant="ghost"
                                size="xs"
                                disabled={isBusy}
                                onClick={() => handleDuplicateCombo(combo)}
                                title={isId ? "Duplikat ke form" : "Duplicate into form"}
                              >
                                <CopyPlus className="h-3.5 w-3.5" />
                              </Button>
                              <Button
                                type="button"
                                variant="secondary"
                                size="xs"
                                disabled={isBusy}
                                onClick={() => handleToggleEnabled(combo)}
                                title={
                                  combo.enabled
                                    ? isId
                                      ? "Jeda combo"
                                      : "Pause combo"
                                    : isId
                                      ? "Aktifkan combo"
                                      : "Enable combo"
                                }
                              >
                                <Power className="h-3 w-3" />
                              </Button>
                              <Button
                                type="button"
                                variant="destructive"
                                size="xs"
                                disabled={isBusy}
                                onClick={() => setConfirmDeleteId(combo.id)}
                                title={isId ? "Hapus combo" : "Delete combo"}
                              >
                                <Trash2 className="h-3 w-3" />
                              </Button>
                            </>
                          )}
                        </div>
                      </div>

                      {/* Steps Chain */}
                      <div className="flex flex-wrap items-center gap-1.5">
                        {steps.map((step, idx) => {
                          const isRef = step.kind === "combo-ref";
                          const ident = isRef ? (step.comboName ?? "") : (step.model ?? "");
                          const isPinned = Boolean(step.connectionId?.trim());
                          const accountLabel = isRef
                            ? null
                            : formatStepAccountLabel(step.connectionId, connectionsById, isId);

                          return (
                            <div
                              key={`${step.kind}-${ident}-${step.connectionId ?? "dynamic"}`}
                              className="inline-flex flex-wrap items-center gap-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-1 text-xs"
                            >
                              <span className="font-mono text-[10px] text-[var(--text-muted)]">
                                #{idx + 1}
                              </span>
                              {isRef ? (
                                <Badge variant="default" className="px-1 py-0 text-[9px]">
                                  ref
                                </Badge>
                              ) : null}
                              <span className="font-mono font-medium text-[var(--text-primary)]">
                                {ident}
                              </span>
                              {accountLabel ? (
                                <span
                                  className={`inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-[10px] ${
                                    isPinned
                                      ? "bg-[var(--bg-surface-raised)] text-[var(--accent-text)]"
                                      : "bg-[var(--bg-surface)] text-[var(--text-muted)]"
                                  }`}
                                >
                                  {isPinned ? (
                                    <KeyRound className="h-2.5 w-2.5" aria-hidden="true" />
                                  ) : (
                                    <Sparkles className="h-2.5 w-2.5" aria-hidden="true" />
                                  )}
                                  <span>{isPinned ? accountLabel : "Dynamic"}</span>
                                </span>
                              ) : null}
                              {step.weight && step.weight > 0 ? (
                                <span className="font-mono text-[10px] font-semibold text-[var(--accent-text)]">
                                  {step.weight}%
                                </span>
                              ) : null}
                            </div>
                          );
                        })}
                      </div>

                      {/* Inline Live Test Result */}
                      {testRes ? (
                        <div className="flex flex-wrap items-center gap-2 border-t border-[var(--border-subtle)] pt-2 text-xs">
                          {testRes.status === "ok" || testRes.status === "slow" ? (
                            <Badge
                              variant={testRes.status === "ok" ? "success" : "warning"}
                              className="gap-1"
                            >
                              <CheckCircle2 className="h-3 w-3" />
                              <span>
                                {testRes.status.toUpperCase()}
                                {typeof testRes.latencyMs === "number"
                                  ? ` · ${testRes.latencyMs}ms`
                                  : ""}
                              </span>
                            </Badge>
                          ) : (
                            <Badge variant="destructive" className="gap-1">
                              <XCircle className="h-3 w-3" />
                              <span>ERROR</span>
                            </Badge>
                          )}
                          {testRes.resolvedBy ? (
                            <span className="font-mono text-[11px] text-[var(--text-secondary)]">
                              resolvedBy: {testRes.resolvedBy}
                            </span>
                          ) : null}
                          {testRes.error ? (
                            <span className="font-mono text-[11px] text-[var(--status-error)]">
                              {testRes.error}
                            </span>
                          ) : null}
                        </div>
                      ) : null}
                    </div>
                  </Card>
                );
              })}
            </div>
          )}
        </div>
      </div>

      <StrategyGuideModal
        open={guideOpen}
        selectedStrategy={guideStrategy}
        onSelectStrategy={(chosen) => setStrategyOverride(chosen)}
        onClose={() => setGuideOpen(false)}
        isId={isId}
      />
    </div>
  );
}
