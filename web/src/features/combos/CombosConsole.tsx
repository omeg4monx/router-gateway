"use client";

import {
  AlertCircle,
  Check,
  CheckCircle2,
  Copy,
  CopyPlus,
  FlaskConical,
  Layers,
  Pencil,
  Plus,
  Power,
  RefreshCw,
  Search,
  Trash2,
  X,
  XCircle,
} from "lucide-react";
import { useMemo, useState } from "react";
import { runSingleComboTest } from "@/features/models/hooks";
import type { CatalogTestResult } from "@/features/models/types";
import { useProviderConnections } from "@/features/providers/hooks";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { ComboBuilderModal } from "./ComboBuilderModal";
import {
  useCombos,
  useCreateComboMutation,
  useDeleteComboMutation,
  useUpdateComboMutation,
} from "./hooks";
import type { Combo, ComboStep, CreateComboRequest } from "./schemas";
import { COMBO_STRATEGY_OPTIONS, getStrategyMeta } from "./strategies";

export function CombosConsole() {
  const locale = useUIStore((s) => s.locale);
  const isId = locale === "id";

  const combosQuery = useCombos();
  const connectionsQuery = useProviderConnections();

  const createMutation = useCreateComboMutation();
  const updateMutation = useUpdateComboMutation();
  const deleteMutation = useDeleteComboMutation();

  const [query, setQuery] = useState("");
  const [strategyFilter, setStrategyFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "paused">("all");

  const [modalOpen, setModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<"create" | "edit">("create");
  const [editingCombo, setEditingCombo] = useState<Combo | null>(null);
  const [modalError, setModalError] = useState<string | null>(null);

  const [busyComboId, setBusyComboId] = useState<string | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [copiedName, setCopiedName] = useState<string | null>(null);

  const [testingNames, setTestingNames] = useState<Set<string>>(() => new Set());
  const [testResults, setTestResults] = useState<Record<string, CatalogTestResult>>({});

  const combos = combosQuery.data ?? [];
  const connections = connectionsQuery.data ?? [];

  const availableModels = useMemo(() => {
    const seen = new Set<string>();
    const out: string[] = [];
    for (const conn of connections) {
      if (!conn.enabled) continue;
      for (const raw of conn.models) {
        const trimmed = raw.trim();
        if (!trimmed) continue;
        const qualified = trimmed.includes("/") ? trimmed : `${conn.provider_id}/${trimmed}`;
        if (!seen.has(qualified)) {
          seen.add(qualified);
          out.push(qualified);
        }
      }
    }
    return out.sort();
  }, [connections]);

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

  const openCreateModal = (duplicateFrom?: Combo) => {
    setModalMode("create");
    setEditingCombo(duplicateFrom ?? null);
    setModalError(null);
    setModalOpen(true);
  };

  const openEditModal = (combo: Combo) => {
    setModalMode("edit");
    setEditingCombo(combo);
    setModalError(null);
    setModalOpen(true);
  };

  const handleModalSubmit = (payload: CreateComboRequest, editComboId?: string) => {
    setModalError(null);
    if (editComboId) {
      updateMutation.mutate(
        { id: editComboId, body: payload },
        {
          onSuccess: () => setModalOpen(false),
          onError: (err) => setModalError(err.message),
        },
      );
      return;
    }
    createMutation.mutate(payload, {
      onSuccess: () => setModalOpen(false),
      onError: (err) => setModalError(err.message),
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
      onSuccess: () => setConfirmDeleteId(null),
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
      window.setTimeout(() => setCopiedName((cur) => (cur === comboName ? null : cur)), 1800);
    } catch {
      // Ignore clipboard errors
    }
  };

  const isRefreshing = combosQuery.isFetching || connectionsQuery.isFetching;

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
              ? "Buat dan kelola combo multi-model dengan 22 strategi routing (priority, weighted, round-robin, p2c, cost-optimized, auto, fusion, pipeline, dll.) beserta nested combo-ref."
              : "Create and orchestrate multi-model combos across all 22 routing strategies (priority, weighted, round-robin, p2c, cost-optimized, auto, fusion, pipeline, and more) with nested combo-ref steps."}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant="outline" className="tabular-nums">
            {activeCount}/{combos.length} {isId ? "Aktif" : "Active"}
          </Badge>
          <Badge variant="secondary" className="tabular-nums">
            {COMBO_STRATEGY_OPTIONS.length} {isId ? "Strategi" : "Strategies"}
          </Badge>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => {
              void combosQuery.refetch();
              void connectionsQuery.refetch();
            }}
            disabled={isRefreshing}
          >
            <RefreshCw
              className={`h-3.5 w-3.5 text-[var(--accent-primary)] ${isRefreshing ? "animate-spin" : ""}`}
            />
            <span>{isId ? "Segarkan" : "Refresh"}</span>
          </Button>
          <Button
            type="button"
            size="sm"
            data-testid="create-combo-btn"
            onClick={() => openCreateModal()}
          >
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{isId ? "Buat Combo" : "Create Combo"}</span>
          </Button>
        </div>
      </section>

      {/* Filter Bar */}
      <Card className="p-4">
        <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-4">
          <div className="relative sm:col-span-2">
            <Search className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]" />
            <Input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder={
                isId
                  ? "Cari nama combo, strategi, atau model..."
                  : "Search combo name, strategy, or target model..."
              }
              className="h-9 pl-8 text-xs"
            />
          </div>

          <select
            aria-label={isId ? "Filter strategi" : "Filter strategy"}
            value={strategyFilter}
            onChange={(e) => setStrategyFilter(e.target.value)}
            className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
          >
            <option value="all">{isId ? "Semua Strategi (22)" : "All Strategies (22)"}</option>
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
            className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)]"
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
          <Button type="button" variant="ghost" size="icon-sm" onClick={() => setRowError(null)}>
            <X className="h-3.5 w-3.5" />
          </Button>
        </Alert>
      ) : null}

      {/* Combos List */}
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
        <Card className="p-12 text-center">
          <Layers className="mx-auto h-8 w-8 text-[var(--text-muted)]" aria-hidden="true" />
          <p className="mt-3 text-sm font-medium text-[var(--text-primary)]">
            {combos.length === 0
              ? isId
                ? "Belum ada combo dikonfigurasi"
                : "No combos configured yet"
              : isId
                ? "Tidak ada combo yang cocok dengan filter"
                : "No combos match the current filters"}
          </p>
          <p className="mt-1 text-xs text-[var(--text-muted)]">
            {isId
              ? "Klik Buat Combo untuk membuat rantai failover, load balancing, cost-optimized, auto, fusion, atau pipeline."
              : "Click Create Combo to build priority failover, weighted, round-robin, cost-optimized, auto, fusion, or pipeline chains."}
          </p>
          <div className="mt-4">
            <Button type="button" size="sm" onClick={() => openCreateModal()}>
              <Plus className="h-3.5 w-3.5" />
              <span>{isId ? "Buat Combo Pertama" : "Create First Combo"}</span>
            </Button>
          </div>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-4">
          {filteredCombos.map((combo) => {
            const stratMeta = getStrategyMeta(combo.strategy);
            const steps: ComboStep[] =
              combo.steps && combo.steps.length > 0
                ? combo.steps
                : combo.targets.map((t): ComboStep => ({ kind: "model", model: t }));
            const cfg = (combo.config ?? {}) as Record<string, unknown>;
            const cfgEntries = Object.entries(cfg).filter(
              ([, v]) => v !== undefined && v !== null && v !== "",
            );
            const testRes = testResults[combo.name];
            const isTesting = testingNames.has(combo.name);
            const isBusy = busyComboId === combo.id;
            const confirmingDelete = confirmDeleteId === combo.id;

            return (
              <Card key={combo.id} className="overflow-hidden">
                <div className="h-[2px] w-full bg-[var(--border-subtle)]">
                  <div
                    className={`h-full transition-all ${
                      combo.enabled
                        ? "w-full bg-[var(--accent-primary)]"
                        : "w-1/3 bg-[var(--status-warn)]"
                    }`}
                  />
                </div>

                <div className="flex flex-col gap-4 p-5 lg:flex-row lg:items-start lg:justify-between">
                  <div className="min-w-0 flex-1 space-y-2.5">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-mono text-base font-semibold text-[var(--text-primary)]">
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
                        {combo.enabled ? (isId ? "Aktif" : "Active") : isId ? "Dijeda" : "Paused"}
                      </Badge>
                      <Badge variant="outline" className="font-mono">
                        {combo.strategy}
                      </Badge>
                      <span className="text-xs text-[var(--text-muted)]">
                        {isId ? stratMeta.labelId : stratMeta.labelEn}
                      </span>
                    </div>

                    {combo.description ? (
                      <p className="text-xs text-[var(--text-secondary)]">{combo.description}</p>
                    ) : null}

                    {/* Steps Chain */}
                    <div className="flex flex-wrap items-center gap-1.5">
                      {steps.map((step, idx) => {
                        const ident = step.kind === "combo-ref" ? step.comboName : step.model;
                        return (
                          <div
                            key={`${step.kind}-${ident}`}
                            className="inline-flex items-center gap-1.5 rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 py-1 text-xs"
                          >
                            <span className="font-mono text-[10px] text-[var(--text-muted)]">
                              #{idx + 1}
                            </span>
                            {step.kind === "combo-ref" ? (
                              <Badge variant="default" className="px-1 py-0 text-[9px]">
                                ref
                              </Badge>
                            ) : null}
                            <span className="font-mono font-medium text-[var(--text-primary)]">
                              {ident}
                            </span>
                            {step.weight && step.weight > 0 ? (
                              <span className="font-mono text-[10px] text-[var(--accent-text)]">
                                {step.weight}%
                              </span>
                            ) : null}
                            {step.label ? (
                              <span className="text-[10px] text-[var(--text-muted)]">
                                ({step.label})
                              </span>
                            ) : null}
                          </div>
                        );
                      })}
                    </div>

                    {/* Runtime Config Chips */}
                    {cfgEntries.length > 0 ? (
                      <div className="flex flex-wrap items-center gap-1.5">
                        {cfgEntries.map(([k, v]) => (
                          <Badge key={k} variant="secondary" className="font-mono text-[10px]">
                            {k}={String(v)}
                          </Badge>
                        ))}
                      </div>
                    ) : null}

                    {/* Inline Live Test Result */}
                    {testRes ? (
                      <div className="flex flex-wrap items-center gap-2 pt-1 text-xs">
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

                  {/* Action Buttons */}
                  <div className="flex shrink-0 flex-wrap items-center gap-1.5">
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
                          variant="outline"
                          size="xs"
                          disabled={isBusy}
                          onClick={() => openEditModal(combo)}
                        >
                          <Pencil className="h-3 w-3" />
                          <span>Edit</span>
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="xs"
                          disabled={isBusy}
                          onClick={() => openCreateModal(combo)}
                          title={isId ? "Duplikat combo" : "Duplicate combo"}
                        >
                          <CopyPlus className="h-3.5 w-3.5" />
                        </Button>
                        <Button
                          type="button"
                          variant="secondary"
                          size="xs"
                          disabled={isBusy}
                          onClick={() => handleToggleEnabled(combo)}
                        >
                          <Power className="h-3 w-3" />
                          <span>
                            {combo.enabled
                              ? isId
                                ? "Jeda"
                                : "Pause"
                              : isId
                                ? "Aktifkan"
                                : "Enable"}
                          </span>
                        </Button>
                        <Button
                          type="button"
                          variant="destructive"
                          size="xs"
                          disabled={isBusy}
                          onClick={() => setConfirmDeleteId(combo.id)}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </>
                    )}
                  </div>
                </div>
              </Card>
            );
          })}
        </div>
      )}

      <ComboBuilderModal
        open={modalOpen}
        mode={modalMode}
        initialCombo={editingCombo}
        existingCombos={combos}
        availableModels={availableModels}
        submitting={createMutation.isPending || updateMutation.isPending}
        submitError={modalError}
        onClose={() => setModalOpen(false)}
        onSubmit={handleModalSubmit}
        isId={isId}
      />
    </div>
  );
}
