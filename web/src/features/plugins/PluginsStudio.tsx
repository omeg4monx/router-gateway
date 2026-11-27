import {
  Check,
  ChevronDown,
  ChevronUp,
  Play,
  RefreshCw,
  Save,
  Search,
  SlidersHorizontal,
  X,
} from "lucide-react";
import { useMemo, useState } from "react";
import { useModelCatalog } from "@/features/models/hooks";
import { cn } from "@/shared/lib/utils";
import { type UILocale, useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  type PluginState,
  usePluginsList,
  usePreviewPluginsMutation,
  useUpdatePluginMutation,
} from "./hooks";
import {
  buildPluginDraft,
  buildPluginUpdatePayload,
  buildToggleEnabledPayload,
  isPluginDraftDirty,
  type PluginDraft,
  parseModelAllowlist,
  parsePreviewPayload,
  summarizePluginScope,
  toggleModelInAllowlist,
} from "./pluginHelpers";

const SAMPLE_PAYLOAD = JSON.stringify(
  {
    model: "gpt-4o",
    messages: [{ role: "user", content: "Summarize the router architecture in JSON." }],
  },
  null,
  2,
);

const SELECT_CLASS =
  "flex h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 text-xs text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus:outline-none";

const COPY = {
  en: {
    title: "Plugins",
    subtitle: "Toggle and scope request & response gateway plugins.",
    activeCount: (active: number, total: number) => `${active}/${total} active`,
    preview: "Preview",
    hidePreview: "Hide preview",
    retry: "Retry",
    dismiss: "Dismiss",
    empty: "No plugins available.",
    configure: "Configure",
    close: "Close",
    modelsLabel: "Model allowlist & wildcards",
    modelsHint: "Comma-separated (e.g. gpt-4o, claude-*). Leave blank for all models.",
    pickerLabel: "Quick model picker",
    pickerSearch: "Filter catalog models...",
    pickerEmpty: "No matching catalog models.",
    levelLabel: "Strictness level",
    levelLite: "Lite — core filter + obvious patterns",
    levelStandard: "Standard — patterns + swap test",
    levelFull: "Full — craftsmanship + hard gates",
    levelUltra: "Ultra — mandatory PASS/FAIL delivery report",
    effortLabel: "Reasoning effort (OpenAI)",
    effortLow: "Low",
    effortMedium: "Medium",
    effortHigh: "High",
    budgetTokensLabel: "Thinking budget tokens (Claude)",
    maxTokensLabel: "Optional max_tokens cap",
    maxTokensHint: "0 disables cap (reasoning disabled only).",
    forceJsonLabel: "Always enforce JSON mode",
    forceJsonHint: "Off: enforce only when prompt asks for JSON. Always repairs JSON responses.",
    contextWindowLabel: "Context window (tokens)",
    minTurnsLabel: "Minimum kept turns",
    parallelToolCallsLabel: "Allow parallel tool calls",
    parallelToolCallsHint: "Sets parallel_tool_calls default when tools are present.",
    save: "Save",
    cancel: "Cancel",
    saved: (name: string) => `Saved ${name}.`,
    toggled: (name: string, enabled: boolean) => `${name} ${enabled ? "enabled" : "disabled"}.`,
    previewTitle: "Dry-Run Preview",
    previewSubtitle: "Test active plugins against a sample JSON request.",
    runPreview: "Run",
    targetModel: "Target model",
    inputPayload: "Input JSON",
    outputPayload: "Output",
    applied: "Applied",
    none: "none",
    invalidJson: "Payload must be a valid JSON object.",
    outputPlaceholder: "Run preview to inspect transformed payload.",
    toggleAria: (name: string, enabled: boolean) => `${enabled ? "Disable" : "Enable"} ${name}`,
    removeModelAria: (model: string) => `Remove ${model}`,
  },
  id: {
    title: "Plugin",
    subtitle: "Aktifkan dan atur cakupan plugin request & response gateway.",
    activeCount: (active: number, total: number) => `${active}/${total} aktif`,
    preview: "Pratinjau",
    hidePreview: "Tutup pratinjau",
    retry: "Coba lagi",
    dismiss: "Tutup",
    empty: "Belum ada plugin tersedia.",
    configure: "Atur",
    close: "Tutup",
    modelsLabel: "Daftar model & wildcard",
    modelsHint: "Pisahkan dengan koma (mis. gpt-4o, claude-*). Kosongkan untuk semua model.",
    pickerLabel: "Pilih model dari katalog",
    pickerSearch: "Cari model katalog...",
    pickerEmpty: "Tidak ada model katalog yang cocok.",
    levelLabel: "Tingkat ketegasan",
    levelLite: "Lite — filter inti + pola umum",
    levelStandard: "Standar — pola slop + swap test",
    levelFull: "Penuh — standar kualitas + hard gate",
    levelUltra: "Ultra — wajib laporan PASS/FAIL",
    effortLabel: "Reasoning effort (OpenAI)",
    effortLow: "Rendah",
    effortMedium: "Sedang",
    effortHigh: "Tinggi",
    budgetTokensLabel: "Budget token thinking (Claude)",
    maxTokensLabel: "Batas opsional max_tokens",
    maxTokensHint: "0 = tanpa batas (hanya matikan reasoning).",
    forceJsonLabel: "Selalu paksa mode JSON",
    forceJsonHint: "Mati: aktif saat prompt meminta JSON. Perbaikan JSON response selalu aktif.",
    contextWindowLabel: "Batas context window (token)",
    minTurnsLabel: "Minimum turn dipertahankan",
    parallelToolCallsLabel: "Izinkan parallel tool calls",
    parallelToolCallsHint: "Mengisi default parallel_tool_calls saat tools tersedia.",
    save: "Simpan",
    cancel: "Batal",
    saved: (name: string) => `${name} disimpan.`,
    toggled: (name: string, enabled: boolean) =>
      `${name} ${enabled ? "diaktifkan" : "dinonaktifkan"}.`,
    previewTitle: "Pratinjau Dry-Run",
    previewSubtitle: "Uji plugin aktif pada contoh payload JSON.",
    runPreview: "Jalankan",
    targetModel: "Model target",
    inputPayload: "Input JSON",
    outputPayload: "Hasil",
    applied: "Diterapkan",
    none: "tidak ada",
    invalidJson: "Payload harus berupa objek JSON yang valid.",
    outputPlaceholder: "Jalankan pratinjau untuk melihat hasil transformasi.",
    toggleAria: (name: string, enabled: boolean) =>
      `${enabled ? "Nonaktifkan" : "Aktifkan"} ${name}`,
    removeModelAria: (model: string) => `Hapus ${model}`,
  },
} as const satisfies Record<UILocale, Record<string, unknown>>;

export function PluginsStudio() {
  const locale = useUIStore((s) => s.locale);
  const copy = COPY[locale];

  const pluginsQuery = usePluginsList();
  const catalogQuery = useModelCatalog();
  const updateMutation = useUpdatePluginMutation();
  const previewMutation = usePreviewPluginsMutation();

  const [expandedKey, setExpandedKey] = useState<string | null>(null);
  const [drafts, setDrafts] = useState<Record<string, PluginDraft>>({});
  const [modelSearch, setModelSearch] = useState("");
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewModel, setPreviewModel] = useState("gpt-4o");
  const [previewPayload, setPreviewPayload] = useState(SAMPLE_PAYLOAD);
  const [previewError, setPreviewError] = useState<string | null>(null);

  const plugins = pluginsQuery.data?.plugins ?? [];
  const activeCount = plugins.filter((plugin) => plugin.enabled).length;

  const catalogModelIds = useMemo(() => {
    const rows = catalogQuery.data ?? [];
    const seen = new Set<string>();
    const out: string[] = [];
    for (const row of rows) {
      const id = row.qualifiedId || row.id;
      if (!id || seen.has(id)) continue;
      seen.add(id);
      out.push(id);
    }
    return out;
  }, [catalogQuery.data]);

  const filteredCatalogModels = useMemo(() => {
    const q = modelSearch.trim().toLowerCase();
    const list = q ? catalogModelIds.filter((id) => id.toLowerCase().includes(q)) : catalogModelIds;
    return list.slice(0, 24);
  }, [catalogModelIds, modelSearch]);

  const getDraft = (plugin: PluginState): PluginDraft =>
    drafts[plugin.key] ?? buildPluginDraft(plugin);

  const updateDraft = (plugin: PluginState, patch: Partial<PluginDraft>) => {
    const current = getDraft(plugin);
    setDrafts((prev) => ({
      ...prev,
      [plugin.key]: { ...current, ...patch },
    }));
  };

  const handleToggleEnabled = (plugin: PluginState) => {
    const nextEnabled = !plugin.enabled;
    setStatusMessage(null);
    setPendingKey(plugin.key);
    updateMutation.mutate(
      {
        key: plugin.key,
        body: buildToggleEnabledPayload(plugin, nextEnabled),
      },
      {
        onSuccess: () => setStatusMessage(copy.toggled(plugin.name, nextEnabled)),
        onSettled: () => setPendingKey(null),
      },
    );
  };

  const handleExpandToggle = (plugin: PluginState) => {
    if (expandedKey === plugin.key) {
      setExpandedKey(null);
      setModelSearch("");
      return;
    }
    setDrafts((prev) => ({
      ...prev,
      [plugin.key]: buildPluginDraft(plugin),
    }));
    setModelSearch("");
    setExpandedKey(plugin.key);
  };

  const handleSaveConfig = (plugin: PluginState) => {
    const draft = getDraft(plugin);
    setStatusMessage(null);
    setPendingKey(plugin.key);
    updateMutation.mutate(
      {
        key: plugin.key,
        body: buildPluginUpdatePayload(plugin, draft),
      },
      {
        onSuccess: () => {
          setStatusMessage(copy.saved(plugin.name));
          setExpandedKey(null);
          setModelSearch("");
        },
        onSettled: () => setPendingKey(null),
      },
    );
  };

  const handlePreview = () => {
    setPreviewError(null);
    const parsed = parsePreviewPayload(previewPayload);
    if (!parsed.ok) {
      setPreviewError(copy.invalidJson);
      return;
    }
    previewMutation.mutate({
      model: previewModel.trim() || "gpt-4o",
      payload: parsed.payload,
    });
  };

  return (
    <div className="space-y-5">
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[var(--border-subtle)] pb-4">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2.5">
            <h1 className="text-xl font-semibold tracking-tight text-[var(--text-primary)]">
              {copy.title}
            </h1>
            {!pluginsQuery.isLoading && !pluginsQuery.error && (
              <Badge variant={activeCount > 0 ? "success" : "secondary"}>
                {copy.activeCount(activeCount, plugins.length)}
              </Badge>
            )}
          </div>
          <p className="mt-0.5 text-sm text-[var(--text-muted)]">{copy.subtitle}</p>
        </div>

        <Button
          type="button"
          variant={previewOpen ? "secondary" : "outline"}
          size="sm"
          onClick={() => setPreviewOpen((open) => !open)}
          aria-expanded={previewOpen}
        >
          <Play className="h-3.5 w-3.5" />
          {previewOpen ? copy.hidePreview : copy.preview}
        </Button>
      </div>

      {statusMessage && (
        <Alert className="items-center justify-between py-2.5">
          <span className="text-xs text-[var(--text-primary)]">{statusMessage}</span>
          <button
            type="button"
            onClick={() => setStatusMessage(null)}
            aria-label={copy.dismiss}
            className="text-[var(--text-muted)] transition-colors hover:text-[var(--text-primary)]"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </Alert>
      )}

      {updateMutation.error && (
        <Alert variant="destructive" className="items-center justify-between py-2.5">
          <AlertDescription className="mt-0 text-[var(--status-error)]">
            {updateMutation.error.message}
          </AlertDescription>
          <button
            type="button"
            onClick={() => updateMutation.reset()}
            aria-label={copy.dismiss}
            className="text-[var(--status-error)] opacity-80 transition-opacity hover:opacity-100"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </Alert>
      )}

      {pluginsQuery.isLoading ? (
        <Card>
          <div className="divide-y divide-[var(--border-subtle)]">
            {[0, 1, 2].map((item) => (
              <div key={item} className="flex items-center justify-between gap-4 px-5 py-4">
                <div className="flex items-center gap-3.5">
                  <div className="h-5 w-9 animate-pulse rounded-full bg-[var(--bg-surface-raised)]" />
                  <div className="space-y-1.5">
                    <div className="h-4 w-36 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
                    <div className="h-3 w-56 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
                  </div>
                </div>
                <div className="h-7 w-16 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
              </div>
            ))}
          </div>
        </Card>
      ) : pluginsQuery.error ? (
        <Alert variant="destructive" className="items-center justify-between">
          <AlertDescription className="mt-0 text-[var(--status-error)]">
            {pluginsQuery.error.message}
          </AlertDescription>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void pluginsQuery.refetch()}
          >
            <RefreshCw className="h-3.5 w-3.5" />
            {copy.retry}
          </Button>
        </Alert>
      ) : plugins.length === 0 ? (
        <Card>
          <CardContent className="py-10 text-center text-sm text-[var(--text-muted)]">
            {copy.empty}
          </CardContent>
        </Card>
      ) : (
        <Card>
          <div className="divide-y divide-[var(--border-subtle)]">
            {plugins.map((plugin) => {
              const isExpanded = expandedKey === plugin.key;
              const isRowPending = updateMutation.isPending && pendingKey === plugin.key;
              const draft = getDraft(plugin);
              const isDirty = isPluginDraftDirty(plugin, draft);
              const selectedModels = parseModelAllowlist(draft.modelsText);

              return (
                <div key={plugin.key}>
                  <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-3.5 sm:px-5">
                    <div className="flex min-w-0 flex-1 items-start gap-3.5">
                      <button
                        type="button"
                        role="switch"
                        aria-checked={plugin.enabled}
                        aria-label={copy.toggleAria(plugin.name, plugin.enabled)}
                        disabled={isRowPending}
                        onClick={() => handleToggleEnabled(plugin)}
                        className={cn(
                          "mt-0.5 inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)] focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50",
                          plugin.enabled
                            ? "border-[var(--accent-primary)] bg-[var(--accent-primary)]"
                            : "border-[var(--border-strong)] bg-[var(--bg-canvas)]",
                        )}
                      >
                        <span
                          className={cn(
                            "pointer-events-none block h-3.5 w-3.5 rounded-full bg-white shadow-sm transition-transform",
                            plugin.enabled ? "translate-x-4" : "translate-x-0.5",
                          )}
                        />
                      </button>

                      <div className="min-w-0 flex-1">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className="text-sm font-medium text-[var(--text-primary)]">
                            {plugin.name}
                          </span>
                          <Badge variant="outline">{plugin.category}</Badge>
                        </div>
                        <p className="mt-0.5 text-xs text-[var(--text-muted)]">
                          {plugin.description}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-3">
                      <span className="hidden font-mono text-xs text-[var(--text-muted)] sm:inline">
                        {summarizePluginScope(plugin, locale)}
                      </span>
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        aria-expanded={isExpanded}
                        onClick={() => handleExpandToggle(plugin)}
                      >
                        <SlidersHorizontal className="h-3.5 w-3.5" />
                        <span className="hidden sm:inline">
                          {isExpanded ? copy.close : copy.configure}
                        </span>
                        {isExpanded ? (
                          <ChevronUp className="h-3.5 w-3.5" />
                        ) : (
                          <ChevronDown className="h-3.5 w-3.5" />
                        )}
                      </Button>
                    </div>
                  </div>

                  {isExpanded && (
                    <div className="space-y-4 border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-4 sm:px-5">
                      <div className="grid gap-4 sm:grid-cols-2">
                        <div className="space-y-2">
                          <Label htmlFor={`plugin-models-${plugin.key}`}>{copy.modelsLabel}</Label>
                          <Input
                            id={`plugin-models-${plugin.key}`}
                            type="text"
                            value={draft.modelsText}
                            onChange={(e) => updateDraft(plugin, { modelsText: e.target.value })}
                            placeholder="gpt-4o, claude-*"
                            className="h-9 text-xs"
                          />
                          <p className="text-[11px] text-[var(--text-muted)]">{copy.modelsHint}</p>

                          {selectedModels.length > 0 && (
                            <div className="flex flex-wrap gap-1.5 pt-1">
                              {selectedModels.map((modelId) => (
                                <Badge
                                  key={modelId}
                                  variant="secondary"
                                  className="gap-1 font-mono text-[11px]"
                                >
                                  <span>{modelId}</span>
                                  <button
                                    type="button"
                                    aria-label={copy.removeModelAria(modelId)}
                                    onClick={() =>
                                      updateDraft(plugin, {
                                        modelsText: toggleModelInAllowlist(
                                          draft.modelsText,
                                          modelId,
                                        ),
                                      })
                                    }
                                    className="text-[var(--text-muted)] transition-colors hover:text-[var(--text-primary)]"
                                  >
                                    <X className="h-3 w-3" />
                                  </button>
                                </Badge>
                              ))}
                            </div>
                          )}

                          {catalogModelIds.length > 0 && (
                            <div className="space-y-1.5 pt-1">
                              <div className="relative">
                                <Search className="pointer-events-none absolute left-2.5 top-2.5 h-3.5 w-3.5 text-[var(--text-muted)]" />
                                <Input
                                  type="text"
                                  value={modelSearch}
                                  onChange={(e) => setModelSearch(e.target.value)}
                                  placeholder={copy.pickerSearch}
                                  aria-label={copy.pickerLabel}
                                  className="h-8 pl-8 text-xs"
                                />
                              </div>
                              {filteredCatalogModels.length === 0 ? (
                                <p className="text-[11px] text-[var(--text-muted)]">
                                  {copy.pickerEmpty}
                                </p>
                              ) : (
                                <div className="flex max-h-28 flex-wrap gap-1.5 overflow-y-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-2">
                                  {filteredCatalogModels.map((modelId) => {
                                    const active = selectedModels.includes(modelId);
                                    return (
                                      <button
                                        key={modelId}
                                        type="button"
                                        onClick={() =>
                                          updateDraft(plugin, {
                                            modelsText: toggleModelInAllowlist(
                                              draft.modelsText,
                                              modelId,
                                            ),
                                          })
                                        }
                                        className={cn(
                                          "inline-flex items-center gap-1 rounded border px-2 py-0.5 font-mono text-[11px] transition-colors",
                                          active
                                            ? "border-[var(--accent-primary)] bg-[var(--accent-subtle)] text-[var(--accent-text)]"
                                            : "border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]",
                                        )}
                                      >
                                        {active && <Check className="h-3 w-3" />}
                                        <span className="max-w-[200px] truncate">{modelId}</span>
                                      </button>
                                    );
                                  })}
                                </div>
                              )}
                            </div>
                          )}
                        </div>

                        {plugin.key === "antiSlop" && (
                          <div className="space-y-1.5">
                            <Label htmlFor={`plugin-level-${plugin.key}`}>{copy.levelLabel}</Label>
                            <select
                              id={`plugin-level-${plugin.key}`}
                              value={draft.level}
                              onChange={(e) => updateDraft(plugin, { level: e.target.value })}
                              className={SELECT_CLASS}
                            >
                              <option value="lite">{copy.levelLite}</option>
                              <option value="standard">{copy.levelStandard}</option>
                              <option value="full">{copy.levelFull}</option>
                              <option value="ultra">{copy.levelUltra}</option>
                            </select>
                          </div>
                        )}

                        {plugin.key === "thinkDeeper" && (
                          <div className="space-y-3">
                            <div className="space-y-1.5">
                              <Label htmlFor={`plugin-effort-${plugin.key}`}>
                                {copy.effortLabel}
                              </Label>
                              <select
                                id={`plugin-effort-${plugin.key}`}
                                value={draft.effort}
                                onChange={(e) => updateDraft(plugin, { effort: e.target.value })}
                                className={SELECT_CLASS}
                              >
                                <option value="low">{copy.effortLow}</option>
                                <option value="medium">{copy.effortMedium}</option>
                                <option value="high">{copy.effortHigh}</option>
                              </select>
                            </div>
                            <div className="space-y-1.5">
                              <Label htmlFor={`plugin-budget-${plugin.key}`}>
                                {copy.budgetTokensLabel}
                              </Label>
                              <Input
                                id={`plugin-budget-${plugin.key}`}
                                type="number"
                                min={1024}
                                step={1024}
                                value={draft.budgetTokens}
                                onChange={(e) =>
                                  updateDraft(plugin, { budgetTokens: e.target.value })
                                }
                                className="h-9 text-xs"
                              />
                            </div>
                          </div>
                        )}

                        {plugin.key === "speedMode" && (
                          <div className="space-y-1.5">
                            <Label htmlFor={`plugin-max-tokens-${plugin.key}`}>
                              {copy.maxTokensLabel}
                            </Label>
                            <Input
                              id={`plugin-max-tokens-${plugin.key}`}
                              type="number"
                              min={0}
                              step={128}
                              value={draft.maxTokens}
                              onChange={(e) => updateDraft(plugin, { maxTokens: e.target.value })}
                              className="h-9 text-xs"
                            />
                            <p className="text-[11px] text-[var(--text-muted)]">
                              {copy.maxTokensHint}
                            </p>
                          </div>
                        )}

                        {plugin.key === "jsonGuard" && (
                          <div className="space-y-1.5">
                            <label className="flex cursor-pointer items-center gap-2.5 pt-1 text-xs font-medium text-[var(--text-primary)]">
                              <input
                                type="checkbox"
                                checked={draft.forceJson}
                                onChange={(e) =>
                                  updateDraft(plugin, { forceJson: e.target.checked })
                                }
                                className="h-4 w-4 rounded border-[var(--border-strong)] accent-[var(--accent-primary)]"
                              />
                              <span>{copy.forceJsonLabel}</span>
                            </label>
                            <p className="text-[11px] text-[var(--text-muted)]">
                              {copy.forceJsonHint}
                            </p>
                          </div>
                        )}

                        {plugin.key === "contextSqueezer" && (
                          <div className="space-y-3">
                            <div className="space-y-1.5">
                              <Label htmlFor={`plugin-ctx-window-${plugin.key}`}>
                                {copy.contextWindowLabel}
                              </Label>
                              <Input
                                id={`plugin-ctx-window-${plugin.key}`}
                                type="number"
                                min={1024}
                                step={4096}
                                value={draft.contextWindow}
                                onChange={(e) =>
                                  updateDraft(plugin, { contextWindow: e.target.value })
                                }
                                className="h-9 text-xs"
                              />
                            </div>
                            <div className="space-y-1.5">
                              <Label htmlFor={`plugin-min-turns-${plugin.key}`}>
                                {copy.minTurnsLabel}
                              </Label>
                              <Input
                                id={`plugin-min-turns-${plugin.key}`}
                                type="number"
                                min={1}
                                step={1}
                                value={draft.minTurns}
                                onChange={(e) => updateDraft(plugin, { minTurns: e.target.value })}
                                className="h-9 text-xs"
                              />
                            </div>
                          </div>
                        )}

                        {plugin.key === "openaiToolBridge" && (
                          <div className="space-y-1.5">
                            <label className="flex cursor-pointer items-center gap-2.5 pt-1 text-xs font-medium text-[var(--text-primary)]">
                              <input
                                type="checkbox"
                                checked={draft.parallelToolCalls}
                                onChange={(e) =>
                                  updateDraft(plugin, { parallelToolCalls: e.target.checked })
                                }
                                className="h-4 w-4 rounded border-[var(--border-strong)] accent-[var(--accent-primary)]"
                              />
                              <span>{copy.parallelToolCallsLabel}</span>
                            </label>
                            <p className="text-[11px] text-[var(--text-muted)]">
                              {copy.parallelToolCallsHint}
                            </p>
                          </div>
                        )}
                      </div>

                      <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] pt-3">
                        <span className="font-mono text-[11px] text-[var(--text-muted)]">
                          {plugin.key}
                        </span>
                        <div className="flex items-center gap-2">
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              setExpandedKey(null);
                              setModelSearch("");
                            }}
                          >
                            {copy.cancel}
                          </Button>
                          <Button
                            type="button"
                            size="sm"
                            disabled={!isDirty || isRowPending}
                            onClick={() => handleSaveConfig(plugin)}
                          >
                            <Save className="h-3.5 w-3.5" />
                            {copy.save}
                          </Button>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </Card>
      )}

      {previewOpen && (
        <Card>
          <CardContent className="space-y-4 p-4 sm:p-5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div>
                <h2 className="text-sm font-semibold text-[var(--text-primary)]">
                  {copy.previewTitle}
                </h2>
                <p className="text-xs text-[var(--text-muted)]">{copy.previewSubtitle}</p>
              </div>
              <Button
                type="button"
                size="sm"
                onClick={handlePreview}
                disabled={previewMutation.isPending}
              >
                <Play className="h-3.5 w-3.5" />
                {copy.runPreview}
              </Button>
            </div>

            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2.5">
                <div className="space-y-1">
                  <Label htmlFor="plugins-preview-model">{copy.targetModel}</Label>
                  <Input
                    id="plugins-preview-model"
                    type="text"
                    value={previewModel}
                    onChange={(e) => setPreviewModel(e.target.value)}
                    className="h-9 text-xs"
                  />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="plugins-preview-payload">{copy.inputPayload}</Label>
                  <textarea
                    id="plugins-preview-payload"
                    rows={7}
                    value={previewPayload}
                    onChange={(e) => setPreviewPayload(e.target.value)}
                    className="w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] p-2.5 font-mono text-xs text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus:outline-none"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-secondary)]">
                    {copy.outputPayload}
                  </span>
                  {previewMutation.data && (
                    <span className="font-mono text-xs text-[var(--accent-text)]">
                      {copy.applied}:{" "}
                      {previewMutation.data.applied_plugins.length > 0
                        ? previewMutation.data.applied_plugins.join(", ")
                        : copy.none}
                    </span>
                  )}
                </div>
                {(previewError || previewMutation.error) && (
                  <Alert variant="destructive" className="py-2">
                    <AlertDescription className="mt-0 text-[var(--status-error)]">
                      {previewError ?? previewMutation.error?.message}
                    </AlertDescription>
                  </Alert>
                )}
                <pre className="h-[214px] overflow-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)]">
                  {previewMutation.data
                    ? JSON.stringify(previewMutation.data.payload, null, 2)
                    : copy.outputPlaceholder}
                </pre>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
