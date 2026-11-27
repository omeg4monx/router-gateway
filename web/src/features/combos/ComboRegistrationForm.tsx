import {
  AlertCircle,
  ArrowDown,
  ArrowUp,
  ChevronDown,
  ChevronRight,
  HelpCircle,
  KeyRound,
  Layers,
  Plus,
  RotateCcw,
  Scale,
  Sparkles,
  Workflow,
  X,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import type { ProviderConnectionSummary } from "@/features/providers/schemas";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  buildProviderAccountGroups,
  formatAuthTypeBadge,
  formatStepAccountLabel,
  getModelsForSelection,
  isDuplicateModelStep,
  qualifyProviderModel,
} from "./comboHelpers";
import type { Combo, ComboStep, ComboStrategy, CreateComboRequest } from "./schemas";
import { MAX_MODEL_IDENTIFIER_RUNES, MAX_TARGETS_COUNT } from "./schemas";
import {
  AUTO_ROUTER_STRATEGIES,
  COMBO_STRATEGY_OPTIONS,
  COMPRESSION_MODES,
  getStrategyMeta,
  STRATEGY_CATEGORIES,
} from "./strategies";

interface ComboRegistrationFormProps {
  mode: "create" | "edit";
  initialCombo: Combo | null;
  existingCombos: Combo[];
  connections: ProviderConnectionSummary[];
  providerLabels: Record<string, string>;
  providerCatalogModels: Record<string, string[]>;
  modelsLoading: boolean;
  submitting: boolean;
  submitError: string | null;
  strategyOverride: ComboStrategy | null;
  onStrategyOverrideApplied: () => void;
  onOpenStrategyGuide: (strategy: ComboStrategy) => void;
  onCancelEdit: () => void;
  onSubmit: (payload: CreateComboRequest, editComboId?: string, onDone?: () => void) => void;
  isId: boolean;
}

export function ComboRegistrationForm({
  mode,
  initialCombo,
  existingCombos,
  connections,
  providerLabels,
  providerCatalogModels,
  modelsLoading,
  submitting,
  submitError,
  strategyOverride,
  onStrategyOverrideApplied,
  onOpenStrategyGuide,
  onCancelEdit,
  onSubmit,
  isId,
}: ComboRegistrationFormProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [strategy, setStrategy] = useState<ComboStrategy>("priority");
  const [enabled, setEnabled] = useState(true);
  const [steps, setSteps] = useState<ComboStep[]>([]);

  // Provider -> Account -> Model picker state
  const [selectedProviderId, setSelectedProviderId] = useState("");
  const [selectedConnectionId, setSelectedConnectionId] = useState("");
  const [selectedModel, setSelectedModel] = useState("");
  const [customModelMode, setCustomModelMode] = useState(false);
  const [customModelInput, setCustomModelInput] = useState("");

  // Strategy-specific options
  const [routerStrategy, setRouterStrategy] = useState("rules");
  const [judgeModel, setJudgeModel] = useState("");
  const [handoffThreshold, setHandoffThreshold] = useState("0.8");
  const [handoffModel, setHandoffModel] = useState("");

  // Advanced / nested combo-ref options (collapsed by default)
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [comboRefDraft, setComboRefDraft] = useState("");
  const [maxRetries, setMaxRetries] = useState("");
  const [retryDelayMs, setRetryDelayMs] = useState("");
  const [timeoutMs, setTimeoutMs] = useState("");
  const [nestedComboMode, setNestedComboMode] = useState<"flatten" | "execute">("flatten");
  const [disableSessionStickiness, setDisableSessionStickiness] = useState(false);
  const [compressionMode, setCompressionMode] = useState("");

  const [formError, setFormError] = useState<string | null>(null);

  const providerGroups = useMemo(
    () => buildProviderAccountGroups(connections, providerLabels),
    [connections, providerLabels],
  );

  const connectionsById = useMemo(
    () => new Map(connections.map((conn) => [conn.id, conn])),
    [connections],
  );

  const activeProviderGroup = useMemo(
    () => providerGroups.find((g) => g.providerId === selectedProviderId),
    [providerGroups, selectedProviderId],
  );

  const availableModelsForProvider = useMemo(
    () => getModelsForSelection(activeProviderGroup, selectedConnectionId, providerCatalogModels),
    [activeProviderGroup, selectedConnectionId, providerCatalogModels],
  );

  const strategyMeta = useMemo(() => getStrategyMeta(strategy), [strategy]);

  const selectableComboRefs = useMemo(() => {
    const trimmedName = name.trim();
    return existingCombos.filter((c) => c.name !== trimmedName);
  }, [existingCombos, name]);

  // Keep default provider selected once connections load
  useEffect(() => {
    if (providerGroups.length === 0) {
      if (selectedProviderId !== "") setSelectedProviderId("");
      return;
    }
    const stillExists = providerGroups.some((g) => g.providerId === selectedProviderId);
    if (!stillExists) {
      setSelectedProviderId(providerGroups[0]?.providerId ?? "");
      setSelectedConnectionId("");
      setSelectedModel("");
      setCustomModelMode(false);
    }
  }, [providerGroups, selectedProviderId]);

  // Keep selected model in sync with available models when provider/account changes
  useEffect(() => {
    if (availableModelsForProvider.length === 0) {
      setSelectedModel("");
      return;
    }
    if (!availableModelsForProvider.includes(selectedModel)) {
      setSelectedModel(availableModelsForProvider[0] ?? "");
    }
  }, [availableModelsForProvider, selectedModel]);

  // Apply strategy chosen from StrategyGuideModal
  useEffect(() => {
    if (!strategyOverride) return;
    setStrategy(strategyOverride);
    onStrategyOverrideApplied();
  }, [strategyOverride, onStrategyOverrideApplied]);

  const resetFormState = useCallback(() => {
    setName("");
    setDescription("");
    setStrategy("priority");
    setEnabled(true);
    setSteps([]);
    setRouterStrategy("rules");
    setJudgeModel("");
    setHandoffThreshold("0.8");
    setHandoffModel("");
    setComboRefDraft("");
    setMaxRetries("");
    setRetryDelayMs("");
    setTimeoutMs("");
    setNestedComboMode("flatten");
    setDisableSessionStickiness(false);
    setCompressionMode("");
    setShowAdvanced(false);
    setCustomModelInput("");
    setFormError(null);
  }, []);

  // Populate form when switching between create / edit / duplicate
  useEffect(() => {
    setFormError(null);
    if (!initialCombo) {
      resetFormState();
      return;
    }

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

    const hasAdvanced = Boolean(
      cfg.maxRetries ||
        cfg.retryDelayMs ||
        cfg.timeoutMs ||
        cfg.nestedComboMode === "execute" ||
        cfg.disableSessionStickiness ||
        cfg.compressionMode ||
        loadedSteps.some((s) => s.kind === "combo-ref"),
    );
    setShowAdvanced(hasAdvanced);
  }, [initialCombo, mode, resetFormState]);

  const totalWeight = steps.reduce((sum, s) => sum + (s.weight ?? 0), 0);

  const handleProviderChange = (nextProviderId: string) => {
    setSelectedProviderId(nextProviderId);
    setSelectedConnectionId("");
    setCustomModelMode(false);
    setCustomModelInput("");
    setFormError(null);
  };

  const handleAddModel = () => {
    const useCustom = customModelMode || providerGroups.length === 0;
    const rawModel = useCustom ? customModelInput.trim() : selectedModel.trim();

    if (!rawModel) {
      setFormError(
        isId ? "Pilih atau masukkan nama model terlebih dahulu." : "Select or enter a model first.",
      );
      return;
    }
    if (rawModel === "*" || /[\r\n\t]/.test(rawModel)) {
      setFormError(isId ? "ID model tidak valid." : "Invalid model identifier.");
      return;
    }
    if (steps.length >= MAX_TARGETS_COUNT) {
      setFormError(
        isId
          ? `Maksimal ${MAX_TARGETS_COUNT} model dalam 1 combo.`
          : `Maximum ${MAX_TARGETS_COUNT} steps per combo.`,
      );
      return;
    }

    const qualified = selectedProviderId
      ? qualifyProviderModel(selectedProviderId, rawModel)
      : rawModel;
    const pinnedConnectionId = activeProviderGroup?.hasMultipleAccounts
      ? selectedConnectionId.trim()
      : "";

    if (isDuplicateModelStep(steps, qualified, pinnedConnectionId)) {
      setFormError(
        isId
          ? `Model "${qualified}" dengan akun yang sama sudah ada di daftar.`
          : `Model "${qualified}" with the same account is already in the list.`,
      );
      return;
    }

    setFormError(null);
    setSteps((prev) => [
      ...prev,
      {
        kind: "model",
        model: qualified,
        ...(pinnedConnectionId ? { connectionId: pinnedConnectionId } : {}),
        weight: strategyMeta.usesWeights ? Math.round(100 / (prev.length + 1)) : 0,
      },
    ]);
    if (useCustom) {
      setCustomModelInput("");
    }
  };

  const handleAddComboRef = () => {
    const trimmed = comboRefDraft.trim();
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
      const current = copy[index];
      const target = copy[nextIdx];
      if (!current || !target) return prev;
      copy[index] = target;
      copy[nextIdx] = current;
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

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
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
        isId ? "Tambahkan minimal 1 model ke combo ini." : "Add at least one model to this combo.",
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
      const parsed = Number.parseInt(maxRetries, 10);
      if (Number.isFinite(parsed) && parsed >= 0) config.maxRetries = parsed;
    }
    if (retryDelayMs.trim() !== "") {
      const parsed = Number.parseInt(retryDelayMs, 10);
      if (Number.isFinite(parsed) && parsed >= 0) config.retryDelayMs = parsed;
    }
    if (timeoutMs.trim() !== "") {
      const parsed = Number.parseInt(timeoutMs, 10);
      if (Number.isFinite(parsed) && parsed >= 1000) config.timeoutMs = parsed;
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
      () => {
        if (mode === "create") {
          resetFormState();
        }
      },
    );
  };

  return (
    <Card className="flex h-full min-h-0 flex-col overflow-hidden border-[var(--border-strong)]">
      {/* Top Accent Bar */}
      <div className="h-[2px] w-full shrink-0 bg-[var(--accent-primary)]" />

      <form
        onSubmit={handleSubmit}
        className="flex min-h-0 flex-1 flex-col overflow-hidden"
        data-testid="combo-registration-form"
      >
        {/* Pinned Form Header */}
        <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] px-4 py-2.5">
          <div className="flex min-w-0 items-center gap-2">
            <h3 className="truncate text-sm font-semibold text-[var(--text-primary)]">
              {mode === "edit"
                ? isId
                  ? "Edit Combo"
                  : "Edit Combo"
                : isId
                  ? "Combo Baru"
                  : "New Combo"}
            </h3>
            {mode === "edit" && initialCombo ? (
              <Badge variant="default" className="truncate font-mono text-[10px]">
                {initialCombo.name}
              </Badge>
            ) : null}
          </div>

          {mode === "edit" ? (
            <Button type="button" variant="ghost" size="xs" onClick={onCancelEdit}>
              <X className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{isId ? "Batal" : "Cancel"}</span>
            </Button>
          ) : null}
        </div>

        {/* Scrollable Form Body */}
        <div className="custom-scrollbar min-h-0 flex-1 space-y-3 overflow-y-auto px-4 py-3">
          {formError || submitError ? (
            <Alert variant="destructive">
              <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]" />
              <AlertDescription className="font-mono text-xs">
                {formError || submitError}
              </AlertDescription>
            </Alert>
          ) : null}

          {/* 1. Combo Name & Status */}
          <div className="space-y-2">
            <div className="flex items-center justify-between gap-2">
              <Label htmlFor="combo-name-input">{isId ? "1. Nama Combo" : "1. Combo Name"}</Label>
              <label className="flex cursor-pointer items-center gap-1.5 text-xs text-[var(--text-secondary)]">
                <input
                  type="checkbox"
                  checked={enabled}
                  onChange={(e) => setEnabled(e.target.checked)}
                  className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                />
                <span>{isId ? "Aktif" : "Active"}</span>
              </label>
            </div>

            <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
              <Input
                id="combo-name-input"
                data-testid="combo-name-input"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder={isId ? "cth. coding-utama" : "e.g. coding-primary"}
                className="h-8 font-mono text-xs"
              />

              <Input
                id="combo-desc-input"
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder={isId ? "Catatan singkat (opsional)" : "Short note (optional)"}
                className="h-8 text-xs"
              />
            </div>
          </div>

          {/* 2. Provider -> Account -> Model Picker */}
          <div className="space-y-3 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
            <div className="flex items-center justify-between gap-2">
              <Label>{isId ? "2. Tambah Model ke Combo" : "2. Add Models to Combo"}</Label>
              <span className="text-[11px] text-[var(--text-muted)]">
                {isId ? "Boleh campur beda provider" : "Multi-provider supported"}
              </span>
            </div>

            {providerGroups.length === 0 ? (
              <div className="rounded-md border border-dashed border-[var(--border-strong)] bg-[var(--bg-surface)] p-3 text-xs text-[var(--text-muted)]">
                {isId
                  ? "Belum ada provider aktif. Tambahkan API key atau OAuth di menu Providers, atau ketik ID model manual di bawah."
                  : "No enabled providers found. Add an API key or OAuth account in Providers, or enter a custom model ID below."}
              </div>
            ) : (
              <div className="space-y-2.5">
                {/* Provider Select */}
                <div>
                  <Label
                    htmlFor="combo-provider-select"
                    className="text-[11px] text-[var(--text-secondary)]"
                  >
                    {isId ? "Provider" : "Provider"}
                  </Label>
                  <select
                    id="combo-provider-select"
                    data-testid="combo-provider-select"
                    value={selectedProviderId}
                    onChange={(e) => handleProviderChange(e.target.value)}
                    className="mt-1 h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2.5 text-xs text-[var(--text-primary)]"
                  >
                    {providerGroups.map((group) => (
                      <option key={group.providerId} value={group.providerId}>
                        {group.displayName} ({group.accounts.length}{" "}
                        {isId ? "akun" : group.accounts.length === 1 ? "account" : "accounts"})
                      </option>
                    ))}
                  </select>
                </div>

                {/* Conditional Account / API Key Picker */}
                {activeProviderGroup ? (
                  activeProviderGroup.hasMultipleAccounts ? (
                    <div>
                      <Label
                        htmlFor="combo-account-select"
                        className="text-[11px] text-[var(--text-secondary)]"
                      >
                        {isId ? "Akun / API Key (OAuth)" : "Account / API Key (OAuth)"}
                      </Label>
                      <select
                        id="combo-account-select"
                        data-testid="combo-account-select"
                        value={selectedConnectionId}
                        onChange={(e) => setSelectedConnectionId(e.target.value)}
                        className="mt-1 h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2.5 text-xs text-[var(--text-primary)]"
                      >
                        <option value="">
                          {isId
                            ? `Dynamic — Semua Akun (${activeProviderGroup.accounts.length} akun)`
                            : `Dynamic — All Accounts (${activeProviderGroup.accounts.length} accounts)`}
                        </option>
                        {activeProviderGroup.accounts.map((acc) => (
                          <option key={acc.id} value={acc.id}>
                            {acc.name} · {formatAuthTypeBadge(acc.auth_type, isId)}
                          </option>
                        ))}
                      </select>
                      <p className="mt-1 text-[11px] text-[var(--text-muted)]">
                        {selectedConnectionId
                          ? isId
                            ? "Model akan dikunci khusus ke akun yang dipilih."
                            : "Model will be pinned specifically to this account."
                          : isId
                            ? "Mode Dynamic otomatis memakai semua akun aktif di provider ini."
                            : "Dynamic mode automatically rotates across all enabled accounts for this provider."}
                      </p>
                    </div>
                  ) : (
                    <div className="flex items-center gap-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 py-1.5 text-[11px] text-[var(--text-secondary)]">
                      <KeyRound
                        className="h-3 w-3 shrink-0 text-[var(--accent-primary)]"
                        aria-hidden="true"
                      />
                      <span className="truncate">
                        {isId ? "Akun aktif:" : "Active account:"}{" "}
                        <strong className="font-medium text-[var(--text-primary)]">
                          {activeProviderGroup.accounts[0]?.name}
                        </strong>{" "}
                        (
                        {formatAuthTypeBadge(
                          activeProviderGroup.accounts[0]?.auth_type ?? "",
                          isId,
                        )}{" "}
                        · Dynamic)
                      </span>
                    </div>
                  )
                ) : null}

                {/* Model Dropdown (with optional manual toggle) */}
                <div>
                  <div className="flex items-center justify-between gap-2">
                    <Label
                      htmlFor="combo-model-select"
                      className="text-[11px] text-[var(--text-secondary)]"
                    >
                      {isId ? "Model" : "Model"}{" "}
                      {availableModelsForProvider.length > 0 ? (
                        <span className="font-mono text-[10px] text-[var(--text-muted)]">
                          ({availableModelsForProvider.length})
                        </span>
                      ) : null}
                    </Label>
                    <button
                      type="button"
                      onClick={() => setCustomModelMode((prev) => !prev)}
                      className="text-[11px] font-medium text-[var(--accent-text)] hover:underline"
                    >
                      {customModelMode
                        ? isId
                          ? "Pilih dari dropdown"
                          : "Pick from dropdown"
                        : isId
                          ? "Ketik manual"
                          : "Custom model"}
                    </button>
                  </div>

                  <div className="mt-1 flex gap-2">
                    {!customModelMode ? (
                      <select
                        id="combo-model-select"
                        data-testid="combo-model-select"
                        value={selectedModel}
                        disabled={availableModelsForProvider.length === 0}
                        onChange={(e) => setSelectedModel(e.target.value)}
                        className="h-9 min-w-0 flex-1 rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2.5 font-mono text-xs text-[var(--text-primary)] disabled:opacity-60"
                      >
                        {availableModelsForProvider.length === 0 ? (
                          <option value="">
                            {modelsLoading
                              ? isId
                                ? "Memuat daftar model..."
                                : "Loading models..."
                              : isId
                                ? "Belum ada model terload pada provider ini"
                                : "No loaded models for this provider"}
                          </option>
                        ) : (
                          availableModelsForProvider.map((modelId) => (
                            <option key={modelId} value={modelId}>
                              {modelId}
                            </option>
                          ))
                        )}
                      </select>
                    ) : (
                      <Input
                        id="combo-model-select"
                        data-testid="combo-model-custom-input"
                        type="text"
                        value={customModelInput}
                        onChange={(e) => setCustomModelInput(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter") {
                            e.preventDefault();
                            handleAddModel();
                          }
                        }}
                        placeholder={
                          selectedProviderId
                            ? `cth. gpt-4o atau ${selectedProviderId}/model-id`
                            : "provider/model-id"
                        }
                        className="h-9 min-w-0 flex-1 font-mono text-xs"
                      />
                    )}

                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      data-testid="combo-add-model-btn"
                      disabled={!customModelMode && availableModelsForProvider.length === 0}
                      onClick={handleAddModel}
                      className="h-9 shrink-0"
                    >
                      <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>{isId ? "Add Model" : "Add Model"}</span>
                    </Button>
                  </div>
                </div>
              </div>
            )}

            {/* Fallback custom add when no providers are enabled */}
            {providerGroups.length === 0 ? (
              <div className="flex gap-2">
                <Input
                  type="text"
                  value={customModelInput}
                  onChange={(e) => setCustomModelInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") {
                      e.preventDefault();
                      handleAddModel();
                    }
                  }}
                  placeholder="provider/model-id (e.g. openai/gpt-4o)"
                  className="h-9 flex-1 font-mono text-xs"
                />
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  data-testid="combo-add-model-btn"
                  onClick={handleAddModel}
                  className="h-9 shrink-0"
                >
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Add Model" : "Add Model"}</span>
                </Button>
              </div>
            ) : null}

            {/* Added Models Queue */}
            <div className="space-y-2 border-t border-[var(--border-subtle)] pt-3">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-medium text-[var(--text-primary)]">
                    {isId ? "Daftar Model Terpilih" : "Selected Models"}
                  </span>
                  <Badge variant="outline" className="text-[10px] tabular-nums">
                    {steps.length}
                  </Badge>
                  {strategyMeta.usesWeights && steps.length > 0 ? (
                    <Badge
                      variant={totalWeight === 100 ? "success" : "warning"}
                      className="text-[10px] tabular-nums"
                    >
                      {isId ? `Bobot: ${totalWeight}%` : `Weight: ${totalWeight}%`}
                    </Badge>
                  ) : null}
                </div>

                {strategyMeta.usesWeights && steps.length > 0 ? (
                  <Button type="button" size="xs" variant="outline" onClick={autoBalanceWeights}>
                    <Scale className="h-3 w-3" aria-hidden="true" />
                    <span>{isId ? "Seimbangkan" : "Auto-Balance"}</span>
                  </Button>
                ) : null}
              </div>

              {steps.length === 0 ? (
                <div className="rounded-md border border-dashed border-[var(--border-strong)] bg-[var(--bg-surface)]/50 px-3 py-5 text-center">
                  <Layers className="mx-auto h-5 w-5 text-[var(--text-muted)]" aria-hidden="true" />
                  <p className="mt-1.5 text-xs text-[var(--text-muted)]">
                    {isId
                      ? "Belum ada model ditambahkan. Pilih provider & model lalu klik Add Model."
                      : "No models added yet. Pick a provider & model above and click Add Model."}
                  </p>
                </div>
              ) : (
                <ol className="space-y-1.5">
                  {steps.map((step, idx) => {
                    const isRef = step.kind === "combo-ref";
                    const ident = isRef ? (step.comboName ?? "") : (step.model ?? "");
                    const accountLabel = isRef
                      ? null
                      : formatStepAccountLabel(step.connectionId, connectionsById, isId);
                    const isPinned = Boolean(step.connectionId?.trim());

                    return (
                      <li
                        key={`${step.kind}-${ident}-${step.connectionId ?? "dynamic"}`}
                        className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-2.5"
                      >
                        <div className="flex items-center justify-between gap-2">
                          <div className="min-w-0 flex-1 space-y-1">
                            <div className="flex flex-wrap items-center gap-1.5">
                              <span className="font-mono text-[10px] font-semibold text-[var(--text-muted)]">
                                #{idx + 1}
                              </span>
                              {isRef ? (
                                <Badge variant="default" className="px-1.5 py-0 text-[9px]">
                                  combo-ref
                                </Badge>
                              ) : null}
                              <span className="truncate font-mono text-xs font-semibold text-[var(--text-primary)]">
                                {ident}
                              </span>
                            </div>

                            {accountLabel ? (
                              <div className="flex items-center gap-1 text-[11px] text-[var(--text-secondary)]">
                                {isPinned ? (
                                  <KeyRound
                                    className="h-3 w-3 shrink-0 text-[var(--accent-primary)]"
                                    aria-hidden="true"
                                  />
                                ) : (
                                  <Sparkles
                                    className="h-3 w-3 shrink-0 text-[var(--status-ok)]"
                                    aria-hidden="true"
                                  />
                                )}
                                <span className="truncate">{accountLabel}</span>
                              </div>
                            ) : null}
                          </div>

                          <div className="flex shrink-0 items-center gap-1">
                            {strategyMeta.usesWeights ? (
                              <div className="mr-1 flex items-center gap-1">
                                <Input
                                  type="number"
                                  min={0}
                                  max={100}
                                  value={step.weight ?? 0}
                                  onChange={(e) =>
                                    updateStep(idx, {
                                      weight: Math.max(
                                        0,
                                        Math.min(100, Number(e.target.value) || 0),
                                      ),
                                    })
                                  }
                                  className="h-7 w-14 px-1.5 text-right font-mono text-xs"
                                  aria-label={isId ? "Bobot persen" : "Weight percentage"}
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
                              aria-label={isId ? "Naikkan urutan" : "Move up"}
                            >
                              <ArrowUp className="h-3.5 w-3.5" />
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon-sm"
                              disabled={idx === steps.length - 1}
                              onClick={() => moveStep(idx, 1)}
                              aria-label={isId ? "Turunkan urutan" : "Move down"}
                            >
                              <ArrowDown className="h-3.5 w-3.5" />
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon-sm"
                              onClick={() => removeStep(idx)}
                              aria-label={isId ? "Hapus langkah" : "Remove step"}
                            >
                              <X className="h-3.5 w-3.5 text-[var(--status-error)]" />
                            </Button>
                          </div>
                        </div>

                        {isRef ? (
                          <label className="mt-2 flex cursor-pointer items-center gap-1.5 border-t border-[var(--border-subtle)] pt-2 text-[11px] text-[var(--text-secondary)]">
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
                      </li>
                    );
                  })}
                </ol>
              )}
            </div>
          </div>

          {/* 3. Strategy Selection + Helper Modal Trigger */}
          <div className="space-y-2.5 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <Label htmlFor="combo-strategy">
                {isId ? "3. Jenis Strategy" : "3. Routing Strategy"}
              </Label>
              <Button
                type="button"
                size="xs"
                variant="outline"
                data-testid="strategy-helper-btn"
                onClick={() => onOpenStrategyGuide(strategy)}
                className="gap-1.5 border-[var(--accent-primary)]/40 text-[var(--accent-text)] hover:border-[var(--accent-primary)]"
              >
                <HelpCircle className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Panduan Strategy" : "Strategy Guide"}</span>
              </Button>
            </div>

            <select
              id="combo-strategy"
              data-testid="combo-strategy-select"
              value={strategy}
              onChange={(e) => setStrategy(e.target.value as ComboStrategy)}
              className="h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
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

            {/* Compact Strategy Quick Summary */}
            <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-2.5 text-xs">
              <div className="flex items-center gap-1.5 font-mono text-[11px] text-[var(--accent-text)]">
                <Workflow className="h-3 w-3 shrink-0" aria-hidden="true" />
                <span className="truncate">{strategyMeta.flowPreview}</span>
              </div>
              <p className="mt-1 leading-relaxed text-[var(--text-secondary)]">
                {isId ? strategyMeta.descId : strategyMeta.descEn}
              </p>
              <button
                type="button"
                onClick={() => onOpenStrategyGuide(strategy)}
                className="mt-1.5 inline-flex items-center gap-1 text-[11px] font-medium text-[var(--accent-text)] hover:underline"
              >
                <span>
                  {isId
                    ? "Lihat cara kerja, kecocokan & kapan dipakai →"
                    : "See how it works, best fit & when to use →"}
                </span>
              </button>
            </div>

            {/* Strategy-specific inputs */}
            {strategyMeta.usesRouterStrategy ? (
              <div className="border-t border-[var(--border-subtle)] pt-2.5">
                <Label htmlFor="auto-router-strategy" className="text-[11px]">
                  {isId ? "Kebijakan Auto Router" : "Auto Router Policy"}
                </Label>
                <select
                  id="auto-router-strategy"
                  value={routerStrategy}
                  onChange={(e) => setRouterStrategy(e.target.value)}
                  className="mt-1 h-8 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
                >
                  {AUTO_ROUTER_STRATEGIES.map((rs) => (
                    <option key={rs} value={rs}>
                      {rs}
                    </option>
                  ))}
                </select>
              </div>
            ) : null}

            {strategyMeta.usesJudgeModel ? (
              <div className="border-t border-[var(--border-subtle)] pt-2.5">
                <Label htmlFor="fusion-judge-model" className="text-[11px]">
                  {isId ? "Model Juri (Fusion Judge)" : "Fusion Judge Model"}
                </Label>
                <Input
                  id="fusion-judge-model"
                  type="text"
                  value={judgeModel}
                  onChange={(e) => setJudgeModel(e.target.value)}
                  placeholder="e.g. openai/gpt-4o (default: model pertama)"
                  className="mt-1 h-8 font-mono text-xs"
                />
              </div>
            ) : null}

            {strategyMeta.usesHandoff ? (
              <div className="grid grid-cols-1 gap-2 border-t border-[var(--border-subtle)] pt-2.5 sm:grid-cols-2">
                <div>
                  <Label htmlFor="handoff-threshold" className="text-[11px]">
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
                  <Label htmlFor="handoff-model" className="text-[11px]">
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

          {/* 4. Collapsible Advanced Options (combo-ref + runtime config) */}
          <div className="rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/50">
            <button
              type="button"
              onClick={() => setShowAdvanced((prev) => !prev)}
              className="flex w-full items-center justify-between px-3.5 py-2.5 text-left text-xs font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
            >
              <span>
                {isId
                  ? "Opsi Lanjutan (Nested Combo-Ref & Runtime)"
                  : "Advanced Options (Nested Combo-Ref & Runtime)"}
              </span>
              {showAdvanced ? (
                <ChevronDown className="h-3.5 w-3.5" aria-hidden="true" />
              ) : (
                <ChevronRight className="h-3.5 w-3.5" aria-hidden="true" />
              )}
            </button>

            {showAdvanced ? (
              <div className="space-y-3 border-t border-[var(--border-subtle)] p-3.5">
                {/* Nested Combo-Ref */}
                <div>
                  <Label className="text-[11px] text-[var(--text-secondary)]">
                    {isId ? "Sisipkan Combo Lain (combo-ref)" : "Nest Another Combo (combo-ref)"}
                  </Label>
                  <div className="mt-1 flex gap-2">
                    <select
                      value={comboRefDraft}
                      onChange={(e) => setComboRefDraft(e.target.value)}
                      className="h-8 flex-1 rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2 font-mono text-xs text-[var(--text-primary)]"
                    >
                      <option value="">
                        {selectableComboRefs.length === 0
                          ? isId
                            ? "Belum ada combo lain"
                            : "No other combos"
                          : isId
                            ? "Pilih combo..."
                            : "Select combo..."}
                      </option>
                      {selectableComboRefs.map((c) => (
                        <option key={c.id} value={c.name}>
                          {c.name} ({c.strategy})
                        </option>
                      ))}
                    </select>
                    <Button
                      type="button"
                      size="xs"
                      variant="outline"
                      disabled={!comboRefDraft}
                      onClick={handleAddComboRef}
                    >
                      <Plus className="h-3 w-3" />
                      <span>Ref</span>
                    </Button>
                  </div>
                </div>

                <div className="grid grid-cols-3 gap-2">
                  <div>
                    <Label htmlFor="adv-max-retries" className="text-[11px]">
                      maxRetries
                    </Label>
                    <Input
                      id="adv-max-retries"
                      type="number"
                      min={0}
                      value={maxRetries}
                      onChange={(e) => setMaxRetries(e.target.value)}
                      placeholder="default"
                      className="mt-1 h-8 font-mono text-xs"
                    />
                  </div>
                  <div>
                    <Label htmlFor="adv-retry-delay" className="text-[11px]">
                      retryDelayMs
                    </Label>
                    <Input
                      id="adv-retry-delay"
                      type="number"
                      min={0}
                      value={retryDelayMs}
                      onChange={(e) => setRetryDelayMs(e.target.value)}
                      placeholder="default"
                      className="mt-1 h-8 font-mono text-xs"
                    />
                  </div>
                  <div>
                    <Label htmlFor="adv-timeout" className="text-[11px]">
                      timeoutMs
                    </Label>
                    <Input
                      id="adv-timeout"
                      type="number"
                      min={1000}
                      value={timeoutMs}
                      onChange={(e) => setTimeoutMs(e.target.value)}
                      placeholder="default"
                      className="mt-1 h-8 font-mono text-xs"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  <div>
                    <Label htmlFor="adv-nested-mode" className="text-[11px]">
                      nestedComboMode
                    </Label>
                    <select
                      id="adv-nested-mode"
                      value={nestedComboMode}
                      onChange={(e) =>
                        setNestedComboMode(e.target.value === "execute" ? "execute" : "flatten")
                      }
                      className="mt-1 h-8 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2 font-mono text-xs text-[var(--text-primary)]"
                    >
                      <option value="flatten">flatten</option>
                      <option value="execute">execute</option>
                    </select>
                  </div>
                  <div>
                    <Label htmlFor="adv-compression" className="text-[11px]">
                      compressionMode
                    </Label>
                    <select
                      id="adv-compression"
                      value={compressionMode}
                      onChange={(e) => setCompressionMode(e.target.value)}
                      className="mt-1 h-8 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] px-2 font-mono text-xs text-[var(--text-primary)]"
                    >
                      {COMPRESSION_MODES.map((m) => (
                        <option key={m || "default"} value={m}>
                          {m || (isId ? "default" : "default")}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <label className="flex cursor-pointer items-center gap-2 text-xs text-[var(--text-secondary)]">
                  <input
                    type="checkbox"
                    checked={disableSessionStickiness}
                    onChange={(e) => setDisableSessionStickiness(e.target.checked)}
                    className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                  />
                  <span>
                    {isId ? "Nonaktifkan session stickiness" : "Disable session stickiness"}
                  </span>
                </label>
              </div>
            ) : null}
          </div>
        </div>

        {/* Pinned Submit Footer */}
        <div className="flex shrink-0 items-center justify-end gap-2 border-t border-[var(--border-subtle)] bg-[var(--bg-surface)] px-4 py-2.5">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            disabled={submitting}
            onClick={() => {
              if (mode === "edit") {
                onCancelEdit();
              } else {
                resetFormState();
              }
            }}
          >
            <RotateCcw className="h-3.5 w-3.5" aria-hidden="true" />
            <span>
              {mode === "edit" ? (isId ? "Batal Edit" : "Cancel Edit") : isId ? "Reset" : "Reset"}
            </span>
          </Button>
          <Button type="submit" size="sm" data-testid="submit-combo-btn" disabled={submitting}>
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
                    ? "Daftarkan Combo"
                    : "Register Combo"}
            </span>
          </Button>
        </div>
      </form>
    </Card>
  );
}
