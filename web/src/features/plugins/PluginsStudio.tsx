import { ChevronDown, ChevronUp, Play, RefreshCw, Save, SlidersHorizontal, X } from "lucide-react";
import { useState } from "react";
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
  isPluginDraftDirty,
  type PluginDraft,
  parsePreviewPayload,
  summarizePluginScope,
} from "./pluginHelpers";

const SAMPLE_PAYLOAD = JSON.stringify(
  {
    model: "gpt-4o",
    messages: [{ role: "user", content: "Summarize the router architecture in JSON." }],
  },
  null,
  2,
);

const COPY = {
  en: {
    title: "Plugins",
    subtitle: "Toggle and scope gateway request interceptors.",
    activeCount: (active: number, total: number) => `${active}/${total} active`,
    preview: "Preview",
    hidePreview: "Hide preview",
    retry: "Retry",
    dismiss: "Dismiss",
    empty: "No plugins available.",
    configure: "Configure",
    close: "Close",
    modelsLabel: "Model allowlist",
    modelsHint: "Comma-separated. Leave blank for all models.",
    levelLabel: "Strictness level",
    levelLight: "Light",
    levelStandard: "Standard",
    levelFull: "Full",
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
  },
  id: {
    title: "Plugin",
    subtitle: "Aktifkan dan atur cakupan interceptor request gateway.",
    activeCount: (active: number, total: number) => `${active}/${total} aktif`,
    preview: "Pratinjau",
    hidePreview: "Tutup pratinjau",
    retry: "Coba lagi",
    dismiss: "Tutup",
    empty: "Belum ada plugin tersedia.",
    configure: "Atur",
    close: "Tutup",
    modelsLabel: "Daftar model",
    modelsHint: "Pisahkan dengan koma. Kosongkan untuk semua model.",
    levelLabel: "Tingkat ketegasan",
    levelLight: "Ringan",
    levelStandard: "Standar",
    levelFull: "Penuh",
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
  },
} as const satisfies Record<UILocale, Record<string, unknown>>;

export function PluginsStudio() {
  const locale = useUIStore((s) => s.locale);
  const copy = COPY[locale];

  const pluginsQuery = usePluginsList();
  const updateMutation = useUpdatePluginMutation();
  const previewMutation = usePreviewPluginsMutation();

  const [expandedKey, setExpandedKey] = useState<string | null>(null);
  const [drafts, setDrafts] = useState<Record<string, PluginDraft>>({});
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewModel, setPreviewModel] = useState("gpt-4o");
  const [previewPayload, setPreviewPayload] = useState(SAMPLE_PAYLOAD);
  const [previewError, setPreviewError] = useState<string | null>(null);

  const plugins = pluginsQuery.data?.plugins ?? [];
  const activeCount = plugins.filter((plugin) => plugin.enabled).length;

  const getDraft = (plugin: PluginState): PluginDraft =>
    drafts[plugin.key] ?? buildPluginDraft(plugin);

  const handleToggleEnabled = (plugin: PluginState) => {
    const nextEnabled = !plugin.enabled;
    setStatusMessage(null);
    setPendingKey(plugin.key);
    updateMutation.mutate(
      {
        key: plugin.key,
        body: buildPluginUpdatePayload(plugin, getDraft(plugin), nextEnabled),
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
      return;
    }
    setDrafts((prev) => ({
      ...prev,
      [plugin.key]: buildPluginDraft(plugin),
    }));
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
                    <div className="border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-4 sm:px-5">
                      <div className="grid gap-4 sm:grid-cols-2">
                        <div className="space-y-1.5">
                          <Label htmlFor={`plugin-models-${plugin.key}`}>{copy.modelsLabel}</Label>
                          <Input
                            id={`plugin-models-${plugin.key}`}
                            type="text"
                            value={draft.modelsText}
                            onChange={(e) =>
                              setDrafts((prev) => ({
                                ...prev,
                                [plugin.key]: { ...draft, modelsText: e.target.value },
                              }))
                            }
                            placeholder="gpt-4o, claude-*"
                            className="h-9 text-xs"
                          />
                          <p className="text-[11px] text-[var(--text-muted)]">{copy.modelsHint}</p>
                        </div>

                        {plugin.key === "antiSlop" && (
                          <div className="space-y-1.5">
                            <Label htmlFor={`plugin-level-${plugin.key}`}>{copy.levelLabel}</Label>
                            <select
                              id={`plugin-level-${plugin.key}`}
                              value={draft.level}
                              onChange={(e) =>
                                setDrafts((prev) => ({
                                  ...prev,
                                  [plugin.key]: { ...draft, level: e.target.value },
                                }))
                              }
                              className="flex h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 text-xs text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus:outline-none"
                            >
                              <option value="light">{copy.levelLight}</option>
                              <option value="standard">{copy.levelStandard}</option>
                              <option value="full">{copy.levelFull}</option>
                            </select>
                          </div>
                        )}
                      </div>

                      <div className="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-[var(--border-subtle)] pt-3">
                        <span className="font-mono text-[11px] text-[var(--text-muted)]">
                          {plugin.key}
                        </span>
                        <div className="flex items-center gap-2">
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => setExpandedKey(null)}
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
