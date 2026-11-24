"use client";

import {
  AlertCircle,
  ArrowDown,
  ArrowUp,
  ChevronDown,
  ChevronRight,
  GitFork,
  Layers,
  Plus,
  Scale,
  Search,
  Sparkles,
  X,
} from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import type { Combo, ComboStep, ComboStrategy, CreateComboRequest } from "./schemas";
import { MAX_MODEL_IDENTIFIER_RUNES, MAX_TARGETS_COUNT } from "./schemas";
import {
  AUTO_ROUTER_STRATEGIES,
  COMBO_STRATEGY_OPTIONS,
  COMPRESSION_MODES,
  getStrategyMeta,
  STRATEGY_CATEGORIES,
} from "./strategies";

interface ComboBuilderModalProps {
  open: boolean;
  mode: "create" | "edit";
  initialCombo?: Combo | null;
  existingCombos: Combo[];
  availableModels: string[];
  submitting: boolean;
  submitError: string | null;
  onClose: () => void;
  onSubmit: (payload: CreateComboRequest, editComboId?: string) => void;
  isId: boolean;
}

export function ComboBuilderModal({
  open,
  mode,
  initialCombo,
  existingCombos,
  availableModels,
  submitting,
  submitError,
  onClose,
  onSubmit,
  isId,
}: ComboBuilderModalProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [strategy, setStrategy] = useState<ComboStrategy>("priority");
  const [enabled, setEnabled] = useState(true);
  const [steps, setSteps] = useState<ComboStep[]>([]);

  // Strategy-specific & runtime config
  const [routerStrategy, setRouterStrategy] = useState("rules");
  const [judgeModel, setJudgeModel] = useState("");
  const [handoffThreshold, setHandoffThreshold] = useState("0.8");
  const [handoffModel, setHandoffModel] = useState("");
  const [maxRetries, setMaxRetries] = useState("");
  const [retryDelayMs, setRetryDelayMs] = useState("");
  const [timeoutMs, setTimeoutMs] = useState("");
  const [nestedComboMode, setNestedComboMode] = useState<"flatten" | "execute">("flatten");
  const [disableSessionStickiness, setDisableSessionStickiness] = useState(false);
  const [compressionMode, setCompressionMode] = useState("");
  const [showAdvanced, setShowAdvanced] = useState(false);

  // Step draft controls
  const [stepKindTab, setStepKindTab] = useState<"model" | "combo-ref">("model");
  const [modelDraft, setModelDraft] = useState("");
  const [modelFilter, setModelFilter] = useState("");
  const [comboRefDraft, setComboRefDraft] = useState("");
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    setFormError(null);
    setModelDraft("");
    setModelFilter("");
    setComboRefDraft("");
    setStepKindTab("model");

    if (initialCombo) {
      setName(mode === "edit" ? initialCombo.name : `${initialCombo.name}-copy`);
      setDescription(initialCombo.description ?? "");
      setStrategy(initialCombo.strategy || "priority");
      setEnabled(initialCombo.enabled);
      const loadedSteps: ComboStep[] =
        initialCombo.steps && initialCombo.steps.length > 0
          ? initialCombo.steps.map((s) => ({ ...s }))
          : (initialCombo.targets ?? []).map((t) => ({ kind: "model", model: t }));
      setSteps(loadedSteps);

      const cfg = (initialCombo.config ?? {}) as Record<string, unknown>;
      setRouterStrategy(typeof cfg.routerStrategy === "string" ? cfg.routerStrategy : "rules");
      setJudgeModel(typeof cfg.judgeModel === "string" ? cfg.judgeModel : "");
      setHandoffThreshold(
        typeof cfg.handoffThreshold === "number" ? String(cfg.handoffThreshold) : "0.8",
      );
      setHandoffModel(typeof cfg.handoffModel === "string" ? cfg.handoffModel : "");
      setMaxRetries(typeof cfg.maxRetries === "number" ? String(cfg.maxRetries) : "");
      setRetryDelayMs(typeof cfg.retryDelayMs === "number" ? String(cfg.retryDelayMs) : "");
      setTimeoutMs(typeof cfg.timeoutMs === "number" ? String(cfg.timeoutMs) : "");
      setNestedComboMode(cfg.nestedComboMode === "execute" ? "execute" : "flatten");
      setDisableSessionStickiness(cfg.disableSessionStickiness === true);
      setCompressionMode(typeof cfg.compressionMode === "string" ? cfg.compressionMode : "");
      setShowAdvanced(
        Boolean(
          cfg.maxRetries ||
            cfg.retryDelayMs ||
            cfg.timeoutMs ||
            cfg.nestedComboMode === "execute" ||
            cfg.disableSessionStickiness ||
            cfg.compressionMode,
        ),
      );
    } else {
      setName("");
      setDescription("");
      setStrategy("priority");
      setEnabled(true);
      setSteps([]);
      setRouterStrategy("rules");
      setJudgeModel("");
      setHandoffThreshold("0.8");
      setHandoffModel("");
      setMaxRetries("");
      setRetryDelayMs("");
      setTimeoutMs("");
      setNestedComboMode("flatten");
      setDisableSessionStickiness(false);
      setCompressionMode("");
      setShowAdvanced(false);
    }
  }, [open, mode, initialCombo]);

  const strategyMeta = useMemo(() => getStrategyMeta(strategy), [strategy]);

  const filteredQuickPickModels = useMemo(() => {
    const q = modelFilter.trim().toLowerCase();
    const list = q ? availableModels.filter((m) => m.toLowerCase().includes(q)) : availableModels;
    return list.slice(0, 18);
  }, [availableModels, modelFilter]);

  const selectableComboRefs = useMemo(() => {
    const trimmedName = name.trim();
    return existingCombos.filter((c) => c.name !== trimmedName);
  }, [existingCombos, name]);

  if (!open) {
    return null;
  }

  const totalWeight = steps.reduce((sum, s) => sum + (s.weight ?? 0), 0);

  const addModelStep = (rawModel: string) => {
    const trimmed = rawModel.trim();
    if (!trimmed) {
      setFormError(isId ? "Masukkan ID model terlebih dahulu." : "Enter a model identifier first.");
      return;
    }
    if (trimmed === "*" || /[\r\n\t]/.test(trimmed)) {
      setFormError(isId ? "ID model tidak valid." : "Invalid model identifier.");
      return;
    }
    if (steps.length >= MAX_TARGETS_COUNT) {
      setFormError(
        isId
          ? `Maksimal ${MAX_TARGETS_COUNT} langkah.`
          : `Maximum ${MAX_TARGETS_COUNT} steps allowed.`,
      );
      return;
    }
    if (steps.some((s) => s.kind === "model" && s.model === trimmed)) {
      setFormError(
        isId ? `Model "${trimmed}" sudah ada di daftar.` : `Model "${trimmed}" is already added.`,
      );
      return;
    }
    setFormError(null);
    setSteps((prev) => [
      ...prev,
      {
        kind: "model",
        model: trimmed,
        weight: strategyMeta.usesWeights ? Math.round(100 / (prev.length + 1)) : 0,
      },
    ]);
    setModelDraft("");
  };

  const addComboRefStep = (refName: string) => {
    const trimmed = refName.trim();
    if (!trimmed) {
      setFormError(isId ? "Pilih combo referensi terlebih dahulu." : "Select a combo reference.");
      return;
    }
    if (trimmed === name.trim()) {
      setFormError(
        isId ? "Combo tidak boleh mereferensikan dirinya sendiri." : "Cannot reference itself.",
      );
      return;
    }
    if (steps.some((s) => s.kind === "combo-ref" && s.comboName === trimmed)) {
      setFormError(
        isId ? `Combo "${trimmed}" sudah ada di daftar.` : `Combo "${trimmed}" is already added.`,
      );
      return;
    }
    setFormError(null);
    setSteps((prev) => [
      ...prev,
      {
        kind: "combo-ref",
        comboName: trimmed,
        weight: strategyMeta.usesWeights ? Math.round(100 / (prev.length + 1)) : 0,
      },
    ]);
    setComboRefDraft("");
  };

  const moveStep = (index: number, delta: -1 | 1) => {
    setSteps((prev) => {
      const nextIdx = index + delta;
      if (nextIdx < 0 || nextIdx >= prev.length) return prev;
      const copy = [...prev];
      const a = copy[index];
      const b = copy[nextIdx];
      if (!a || !b) return prev;
      copy[index] = b;
      copy[nextIdx] = a;
      return copy;
    });
  };

  const updateStep = (index: number, patch: Partial<ComboStep>) => {
    setSteps((prev) => prev.map((s, i) => (i === index ? { ...s, ...patch } : s)));
  };

  const removeStep = (index: number) => {
    setSteps((prev) => prev.filter((_, i) => i !== index));
  };

  const autoBalanceWeights = () => {
    if (steps.length === 0) return;
    const base = Math.floor(100 / steps.length);
    const rem = 100 - base * steps.length;
    setSteps((prev) =>
      prev.map((s, i) => ({
        ...s,
        weight: base + (i < rem ? 1 : 0),
      })),
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const cleanName = name.trim();
    if (
      !cleanName ||
      cleanName === "*" ||
      Array.from(cleanName).length > MAX_MODEL_IDENTIFIER_RUNES
    ) {
      setFormError(isId ? "Nama combo wajib diisi dan valid." : "A valid combo name is required.");
      return;
    }
    if (steps.length === 0) {
      setFormError(
        isId
          ? "Tambahkan minimal satu langkah model atau combo-ref."
          : "Add at least one model or combo-ref step.",
      );
      return;
    }

    const config: Record<string, unknown> = {};
    if (strategyMeta.usesRouterStrategy && routerStrategy) {
      config.routerStrategy = routerStrategy;
    }
    if (strategyMeta.usesJudgeModel && judgeModel.trim()) {
      config.judgeModel = judgeModel.trim();
    }
    if (strategyMeta.usesHandoff) {
      const th = Number(handoffThreshold);
      if (Number.isFinite(th) && th >= 0.5 && th <= 0.94) {
        config.handoffThreshold = th;
      }
      if (handoffModel.trim()) {
        config.handoffModel = handoffModel.trim();
      }
    }
    if (maxRetries.trim() !== "") {
      const n = Number.parseInt(maxRetries, 10);
      if (Number.isFinite(n) && n >= 0) config.maxRetries = n;
    }
    if (retryDelayMs.trim() !== "") {
      const n = Number.parseInt(retryDelayMs, 10);
      if (Number.isFinite(n) && n >= 0) config.retryDelayMs = n;
    }
    if (timeoutMs.trim() !== "") {
      const n = Number.parseInt(timeoutMs, 10);
      if (Number.isFinite(n) && n >= 1000) config.timeoutMs = n;
    }
    if (nestedComboMode === "execute") {
      config.nestedComboMode = "execute";
    }
    if (disableSessionStickiness) {
      config.disableSessionStickiness = true;
    }
    if (compressionMode) {
      config.compressionMode = compressionMode;
    }

    const hasQuotaOnlyRef = steps.some(
      (s) => s.kind === "combo-ref" && s.fallbackOnlyOnQuotaExhaustion,
    );
    if (hasQuotaOnlyRef) {
      config.nestedComboMode = "execute";
    }

    const targets = steps
      .map((s) => (s.kind === "combo-ref" ? (s.comboName ?? "") : (s.model ?? "")).trim())
      .filter(Boolean);

    setFormError(null);
    onSubmit(
      {
        name: cleanName,
        description: description.trim(),
        strategy,
        enabled,
        targets,
        steps,
        config,
      },
      mode === "edit" ? initialCombo?.id : undefined,
    );
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-[2px]"
      role="dialog"
      aria-modal="true"
      aria-labelledby="combo-builder-modal-title"
    >
      <form
        onSubmit={handleSubmit}
        className="flex max-h-[90vh] w-full max-w-3xl flex-col overflow-hidden rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl"
      >
        {/* Sticky Header */}
        <div className="flex items-center justify-between border-b border-[var(--border-subtle)] px-5 py-4">
          <div className="flex items-center gap-2.5">
            <GitFork className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <div>
              <h3
                id="combo-builder-modal-title"
                className="text-sm font-semibold text-[var(--text-primary)]"
              >
                {mode === "edit"
                  ? isId
                    ? `Edit Combo: ${initialCombo?.name ?? ""}`
                    : `Edit Combo: ${initialCombo?.name ?? ""}`
                  : isId
                    ? "Buat Combo Baru"
                    : "Create Combo"}
              </h3>
              <p className="text-xs text-[var(--text-muted)]">
                {isId
                  ? "Pilih salah satu dari 22 strategi routing, susun langkah model atau combo-ref, dan atur konfigurasi runtime."
                  : "Pick from all 22 routing strategies, compose model or combo-ref steps, and tune runtime behavior."}
              </p>
            </div>
          </div>
          <Button type="button" variant="ghost" size="icon-sm" onClick={onClose}>
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        {/* Scrollable Body */}
        <div className="flex-1 space-y-5 overflow-y-auto px-5 py-4">
          {formError || submitError ? (
            <Alert variant="destructive">
              <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
              <AlertDescription className="font-mono text-xs">
                {formError || submitError}
              </AlertDescription>
            </Alert>
          ) : null}

          {/* Basic Info */}
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div>
              <Label htmlFor="omni-combo-name">{isId ? "Nama Combo" : "Combo Name"}</Label>
              <Input
                id="omni-combo-name"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. coding-smart-router"
                className="mt-1.5 font-mono text-xs"
              />
            </div>
            <div>
              <Label htmlFor="omni-combo-desc">
                {isId ? "Deskripsi (Opsional)" : "Description (Optional)"}
              </Label>
              <Input
                id="omni-combo-desc"
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder={
                  isId
                    ? "Failover multi-provider otomatis"
                    : "Multi-provider resilient routing chain"
                }
                className="mt-1.5 text-xs"
              />
            </div>
          </div>

          {/* 22-Strategy Picker */}
          <div className="space-y-2.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <Label htmlFor="combo-strategy">
                {isId ? "Strategi Routing (22 Jenis)" : "Routing Strategy (22 Types)"}
              </Label>
              <label className="flex cursor-pointer items-center gap-2 text-xs text-[var(--text-primary)]">
                <input
                  type="checkbox"
                  checked={enabled}
                  onChange={(e) => setEnabled(e.target.checked)}
                  className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                />
                <span>{isId ? "Aktif segera" : "Active immediately"}</span>
              </label>
            </div>

            <select
              id="combo-strategy"
              data-testid="combo-strategy-select"
              value={strategy}
              onChange={(e) => setStrategy(e.target.value as ComboStrategy)}
              className="h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-3 font-mono text-xs text-[var(--text-primary)]"
            >
              {STRATEGY_CATEGORIES.map((cat) => (
                <optgroup key={cat.id} label={isId ? cat.labelId : cat.labelEn}>
                  {COMBO_STRATEGY_OPTIONS.filter((s) => s.category === cat.id).map((s) => (
                    <option key={s.value} value={s.value}>
                      {s.value} — {isId ? s.labelId : s.labelEn}
                    </option>
                  ))}
                </optgroup>
              ))}
            </select>

            <p className="text-xs text-[var(--text-secondary)]">
              <Sparkles className="mr-1.5 inline h-3.5 w-3.5 text-[var(--accent-primary)]" />
              {isId ? strategyMeta.descId : strategyMeta.descEn}
            </p>

            {/* Dynamic Strategy-Specific Options */}
            {strategyMeta.usesRouterStrategy ? (
              <div className="mt-2 grid grid-cols-1 gap-2 border-t border-[var(--border-subtle)] pt-3 sm:grid-cols-2">
                <div>
                  <Label htmlFor="auto-router-strategy" className="text-xs">
                    {isId ? "Kebijakan Auto Router" : "Auto Router Policy"}
                  </Label>
                  <select
                    id="auto-router-strategy"
                    value={routerStrategy}
                    onChange={(e) => setRouterStrategy(e.target.value)}
                    className="mt-1 h-8 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
                  >
                    {AUTO_ROUTER_STRATEGIES.map((rs) => (
                      <option key={rs} value={rs}>
                        {rs}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            ) : null}

            {strategyMeta.usesJudgeModel ? (
              <div className="mt-2 border-t border-[var(--border-subtle)] pt-3">
                <Label htmlFor="fusion-judge-model" className="text-xs">
                  {isId ? "Model Juri (Fusion Judge Model)" : "Fusion Judge Model"}
                </Label>
                <Input
                  id="fusion-judge-model"
                  type="text"
                  value={judgeModel}
                  onChange={(e) => setJudgeModel(e.target.value)}
                  placeholder="e.g. openai/gpt-4o (defaults to first panel model)"
                  className="mt-1 h-8 font-mono text-xs"
                />
              </div>
            ) : null}

            {strategyMeta.usesHandoff ? (
              <div className="mt-2 grid grid-cols-1 gap-2 border-t border-[var(--border-subtle)] pt-3 sm:grid-cols-2">
                <div>
                  <Label htmlFor="handoff-threshold" className="text-xs">
                    {isId ? "Ambang Handoff (0.50–0.94)" : "Handoff Threshold (0.50–0.94)"}
                  </Label>
                  <Input
                    id="handoff-threshold"
                    type="number"
                    step="0.02"
                    min="0.5"
                    max="0.94"
                    value={handoffThreshold}
                    onChange={(e) => setHandoffThreshold(e.target.value)}
                    className="mt-1 h-8 font-mono text-xs"
                  />
                </div>
                <div>
                  <Label htmlFor="handoff-model" className="text-xs">
                    {isId ? "Model Handoff (Opsional)" : "Handoff Model (Optional)"}
                  </Label>
                  <Input
                    id="handoff-model"
                    type="text"
                    value={handoffModel}
                    onChange={(e) => setHandoffModel(e.target.value)}
                    placeholder="e.g. gemini/gemini-2.5-pro"
                    className="mt-1 h-8 font-mono text-xs"
                  />
                </div>
              </div>
            ) : null}
          </div>

          {/* Step Builder (Model Steps + Combo-Ref Steps) */}
          <div className="space-y-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="flex items-center gap-2">
                <Label>{isId ? "Langkah Routing (Steps)" : "Routing Steps"}</Label>
                <Badge variant="outline" className="text-[11px] tabular-nums">
                  {steps.length}
                </Badge>
                {strategyMeta.usesWeights ? (
                  <Badge
                    variant={totalWeight === 100 ? "success" : "warning"}
                    className="text-[11px] tabular-nums"
                  >
                    {isId ? `Total bobot: ${totalWeight}%` : `Total weight: ${totalWeight}%`}
                  </Badge>
                ) : null}
              </div>

              <div className="flex items-center gap-2">
                {strategyMeta.usesWeights && steps.length > 0 ? (
                  <Button type="button" size="xs" variant="outline" onClick={autoBalanceWeights}>
                    <Scale className="h-3 w-3" />
                    <span>{isId ? "Seimbangkan" : "Auto-Balance"}</span>
                  </Button>
                ) : null}
                <div className="inline-flex rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-0.5 text-[11px]">
                  <button
                    type="button"
                    onClick={() => setStepKindTab("model")}
                    className={`rounded px-2.5 py-0.5 font-medium ${
                      stepKindTab === "model"
                        ? "bg-[var(--accent-primary)] text-white"
                        : "text-[var(--text-secondary)]"
                    }`}
                  >
                    {isId ? "+ Model" : "+ Model"}
                  </button>
                  <button
                    type="button"
                    onClick={() => setStepKindTab("combo-ref")}
                    className={`rounded px-2.5 py-0.5 font-medium ${
                      stepKindTab === "combo-ref"
                        ? "bg-[var(--accent-primary)] text-white"
                        : "text-[var(--text-secondary)]"
                    }`}
                  >
                    {isId ? "+ Combo Ref" : "+ Combo Ref"}
                  </button>
                </div>
              </div>
            </div>

            {stepKindTab === "model" ? (
              <div className="space-y-2">
                <div className="flex gap-2">
                  <Input
                    type="text"
                    value={modelDraft}
                    onChange={(e) => setModelDraft(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") {
                        e.preventDefault();
                        addModelStep(modelDraft);
                      }
                    }}
                    placeholder="provider/model (e.g. openai/gpt-4o)"
                    className="h-9 font-mono text-xs"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => addModelStep(modelDraft)}
                    className="h-9 shrink-0"
                  >
                    <Plus className="h-3.5 w-3.5" />
                    <span>{isId ? "Tambah" : "Add"}</span>
                  </Button>
                </div>

                {availableModels.length > 0 ? (
                  <div className="space-y-1.5">
                    <div className="relative">
                      <Search className="pointer-events-none absolute top-1/2 left-2.5 h-3 w-3 -translate-y-1/2 text-[var(--text-muted)]" />
                      <Input
                        type="text"
                        value={modelFilter}
                        onChange={(e) => setModelFilter(e.target.value)}
                        placeholder={
                          isId
                            ? "Filter cepat model dari provider aktif..."
                            : "Quick-filter models from enabled providers..."
                        }
                        className="h-7 pl-7 text-[11px]"
                      />
                    </div>
                    <div className="flex max-h-24 flex-wrap gap-1 overflow-y-auto">
                      {filteredQuickPickModels.map((m) => {
                        const exists = steps.some((s) => s.kind === "model" && s.model === m);
                        return (
                          <Button
                            key={m}
                            type="button"
                            size="xs"
                            variant={exists ? "secondary" : "outline"}
                            disabled={exists}
                            onClick={() => addModelStep(m)}
                            className="font-mono text-[11px]"
                          >
                            <Plus className="h-2.5 w-2.5" />
                            <span>{m}</span>
                          </Button>
                        );
                      })}
                    </div>
                  </div>
                ) : null}
              </div>
            ) : (
              <div className="flex gap-2">
                <select
                  value={comboRefDraft}
                  onChange={(e) => setComboRefDraft(e.target.value)}
                  className="h-9 flex-1 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
                >
                  <option value="">
                    {selectableComboRefs.length === 0
                      ? isId
                        ? "Belum ada combo lain"
                        : "No other combos available"
                      : isId
                        ? "Pilih combo untuk disarang (combo-ref)..."
                        : "Select a combo to nest (combo-ref)..."}
                  </option>
                  {selectableComboRefs.map((c) => (
                    <option key={c.id} value={c.name}>
                      {c.name} ({c.strategy})
                    </option>
                  ))}
                </select>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={!comboRefDraft}
                  onClick={() => addComboRefStep(comboRefDraft)}
                  className="h-9 shrink-0"
                >
                  <Plus className="h-3.5 w-3.5" />
                  <span>{isId ? "Tambah Ref" : "Add Ref"}</span>
                </Button>
              </div>
            )}

            {/* Configured Steps List */}
            {steps.length === 0 ? (
              <p className="rounded border border-dashed border-[var(--border-subtle)] px-3 py-4 text-center text-xs text-[var(--text-muted)]">
                {isId
                  ? "Belum ada langkah. Tambahkan model atau combo-ref di atas."
                  : "No steps added yet. Add a model or nested combo-ref above."}
              </p>
            ) : (
              <ol className="space-y-2">
                {steps.map((step, idx) => {
                  const ident = step.kind === "combo-ref" ? step.comboName : step.model;
                  return (
                    <li
                      key={`${step.kind}-${ident}`}
                      className="space-y-2 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-2.5"
                    >
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <div className="flex min-w-0 items-center gap-2">
                          <Badge variant="outline" className="tabular-nums text-[10px]">
                            #{idx + 1}
                          </Badge>
                          <Badge
                            variant={step.kind === "combo-ref" ? "default" : "secondary"}
                            className="text-[10px]"
                          >
                            {step.kind === "combo-ref" ? "combo-ref" : "model"}
                          </Badge>
                          <span className="truncate font-mono text-xs font-semibold text-[var(--text-primary)]">
                            {ident}
                          </span>
                        </div>

                        <div className="flex items-center gap-1">
                          {strategyMeta.usesWeights ? (
                            <div className="mr-1 flex items-center gap-1">
                              <Input
                                type="number"
                                min={0}
                                max={100}
                                value={step.weight ?? 0}
                                onChange={(e) =>
                                  updateStep(idx, {
                                    weight: Math.max(0, Math.min(100, Number(e.target.value) || 0)),
                                  })
                                }
                                className="h-7 w-16 px-2 text-right font-mono text-xs"
                                aria-label="Weight percentage"
                              />
                              <span className="text-[11px] text-[var(--text-muted)]">%</span>
                            </div>
                          ) : null}

                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            disabled={idx === 0}
                            onClick={() => moveStep(idx, -1)}
                          >
                            <ArrowUp className="h-3.5 w-3.5" />
                          </Button>
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            disabled={idx === steps.length - 1}
                            onClick={() => moveStep(idx, 1)}
                          >
                            <ArrowDown className="h-3.5 w-3.5" />
                          </Button>
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => removeStep(idx)}
                          >
                            <X className="h-3.5 w-3.5 text-[var(--status-error)]" />
                          </Button>
                        </div>
                      </div>

                      <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                        <Input
                          type="text"
                          value={step.label ?? ""}
                          onChange={(e) => updateStep(idx, { label: e.target.value })}
                          placeholder={isId ? "Label langkah (opsional)" : "Step label (optional)"}
                          className="h-7 text-[11px]"
                        />
                        {step.kind === "combo-ref" ? (
                          <label className="flex cursor-pointer items-center gap-1.5 text-[11px] text-[var(--text-secondary)]">
                            <input
                              type="checkbox"
                              checked={Boolean(step.fallbackOnlyOnQuotaExhaustion)}
                              onChange={(e) =>
                                updateStep(idx, {
                                  fallbackOnlyOnQuotaExhaustion: e.target.checked,
                                })
                              }
                              className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                            />
                            <span>
                              {isId
                                ? "Failover hanya saat kuota habis"
                                : "Fallback only on quota exhaustion"}
                            </span>
                          </label>
                        ) : null}
                      </div>

                      {strategyMeta.usesStepPrompt && step.kind === "model" ? (
                        <textarea
                          rows={2}
                          value={step.prompt ?? ""}
                          onChange={(e) => updateStep(idx, { prompt: e.target.value })}
                          placeholder={
                            isId
                              ? "Instruksi sistem langkah pipeline (opsional)..."
                              : "Pipeline step system instruction (optional)..."
                          }
                          className="w-full rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-2 font-mono text-[11px] text-[var(--text-primary)]"
                        />
                      ) : null}
                    </li>
                  );
                })}
              </ol>
            )}
          </div>

          {/* Collapsible Runtime Config */}
          <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)]">
            <button
              type="button"
              onClick={() => setShowAdvanced((v) => !v)}
              className="flex w-full items-center justify-between px-3.5 py-2.5 text-xs font-medium text-[var(--text-primary)]"
            >
              <span>{isId ? "Konfigurasi Runtime Lanjutan" : "Advanced Runtime Config"}</span>
              {showAdvanced ? (
                <ChevronDown className="h-4 w-4 text-[var(--text-muted)]" />
              ) : (
                <ChevronRight className="h-4 w-4 text-[var(--text-muted)]" />
              )}
            </button>
            {showAdvanced ? (
              <div className="grid grid-cols-1 gap-3 border-t border-[var(--border-subtle)] p-3.5 sm:grid-cols-3">
                <div>
                  <Label htmlFor="cfg-max-retries" className="text-xs">
                    maxRetries
                  </Label>
                  <Input
                    id="cfg-max-retries"
                    type="number"
                    min={0}
                    max={10}
                    value={maxRetries}
                    onChange={(e) => setMaxRetries(e.target.value)}
                    placeholder="2"
                    className="mt-1 h-8 font-mono text-xs"
                  />
                </div>
                <div>
                  <Label htmlFor="cfg-retry-delay" className="text-xs">
                    retryDelayMs
                  </Label>
                  <Input
                    id="cfg-retry-delay"
                    type="number"
                    min={0}
                    max={60000}
                    value={retryDelayMs}
                    onChange={(e) => setRetryDelayMs(e.target.value)}
                    placeholder="500"
                    className="mt-1 h-8 font-mono text-xs"
                  />
                </div>
                <div>
                  <Label htmlFor="cfg-timeout-ms" className="text-xs">
                    timeoutMs
                  </Label>
                  <Input
                    id="cfg-timeout-ms"
                    type="number"
                    min={1000}
                    value={timeoutMs}
                    onChange={(e) => setTimeoutMs(e.target.value)}
                    placeholder="30000"
                    className="mt-1 h-8 font-mono text-xs"
                  />
                </div>
                <div>
                  <Label htmlFor="cfg-nested-mode" className="text-xs">
                    nestedComboMode
                  </Label>
                  <select
                    id="cfg-nested-mode"
                    value={nestedComboMode}
                    onChange={(e) => setNestedComboMode(e.target.value as "flatten" | "execute")}
                    className="mt-1 h-8 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2 font-mono text-xs text-[var(--text-primary)]"
                  >
                    <option value="flatten">flatten</option>
                    <option value="execute">execute</option>
                  </select>
                </div>
                <div>
                  <Label htmlFor="cfg-compression" className="text-xs">
                    compressionMode
                  </Label>
                  <select
                    id="cfg-compression"
                    value={compressionMode}
                    onChange={(e) => setCompressionMode(e.target.value)}
                    className="mt-1 h-8 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2 font-mono text-xs text-[var(--text-primary)]"
                  >
                    {COMPRESSION_MODES.map((m) => (
                      <option key={m || "default"} value={m}>
                        {m || (isId ? "(default global)" : "(global default)")}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="flex items-end pb-1">
                  <label className="flex cursor-pointer items-center gap-2 text-xs text-[var(--text-primary)]">
                    <input
                      type="checkbox"
                      checked={disableSessionStickiness}
                      onChange={(e) => setDisableSessionStickiness(e.target.checked)}
                      className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                    />
                    <span>disableSessionStickiness</span>
                  </label>
                </div>
              </div>
            ) : null}
          </div>
        </div>

        {/* Sticky Footer */}
        <div className="flex items-center justify-end gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3">
          <Button type="button" variant="secondary" size="sm" onClick={onClose}>
            {isId ? "Batal" : "Cancel"}
          </Button>
          <Button
            type="submit"
            size="sm"
            data-testid="combo-builder-submit-btn"
            disabled={submitting}
          >
            <Layers className="h-3.5 w-3.5" aria-hidden="true" />
            <span>
              {submitting
                ? isId
                  ? "Menyimpan..."
                  : "Saving..."
                : mode === "edit"
                  ? isId
                    ? "Simpan Perubahan"
                    : "Save Changes"
                  : isId
                    ? "Buat Combo"
                    : "Create Combo"}
            </span>
          </Button>
        </div>
      </form>
    </div>
  );
}
