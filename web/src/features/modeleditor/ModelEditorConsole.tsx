import { zodResolver } from "@hookform/resolvers/zod";
import { CircleHelp, Play, Plus, RefreshCw, Sliders, Trash2, X } from "lucide-react";
import { useCallback, useState } from "react";
import { useForm } from "react-hook-form";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useCreateModelOverrideMutation,
  useDeleteModelOverrideMutation,
  useModelOverrides,
  usePreviewModelEditorMutation,
  useUpdateModelOverrideMutation,
} from "./hooks";
import { ModelEditorHelpModal } from "./ModelEditorHelpModal";
import { MODEL_EDITOR_COPY, type ParameterGuideId } from "./parameterGuide";
import {
  type CreateModelOverrideFormValues,
  createModelOverrideSchema,
  parseExtraParamsJson,
  parseStripParamsList,
} from "./schemas";

type MobileTab = "override" | "rules" | "dryrun";

interface FieldLabelProps {
  htmlFor: string;
  label: string;
  helpTopic?: ParameterGuideId;
  helpAriaPrefix: string;
  onOpenHelp: (topic: ParameterGuideId) => void;
}

function FieldLabel({ htmlFor, label, helpTopic, helpAriaPrefix, onOpenHelp }: FieldLabelProps) {
  return (
    <div className="flex items-center justify-between gap-1">
      <Label htmlFor={htmlFor} className="truncate text-xs">
        {label}
      </Label>
      {helpTopic ? (
        <button
          type="button"
          onClick={() => onOpenHelp(helpTopic)}
          aria-label={`${helpAriaPrefix} ${label}`}
          className="inline-flex h-4 w-4 shrink-0 cursor-pointer items-center justify-center rounded-full text-[var(--text-muted)] transition-colors hover:text-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-[var(--accent-primary)]"
        >
          <CircleHelp className="h-3.5 w-3.5" aria-hidden="true" />
        </button>
      ) : null}
    </div>
  );
}

export function ModelEditorConsole() {
  const locale = useUIStore((state) => state.locale);
  const copy = MODEL_EDITOR_COPY[locale];

  const overridesQuery = useModelOverrides();
  const createOverride = useCreateModelOverrideMutation();
  const updateOverride = useUpdateModelOverrideMutation();
  const deleteOverride = useDeleteModelOverrideMutation();
  const previewMutation = usePreviewModelEditorMutation();

  const [mobileTab, setMobileTab] = useState<MobileTab>("override");
  const [helpOpen, setHelpOpen] = useState(false);
  const [helpTopic, setHelpTopic] = useState<ParameterGuideId>("temperature");
  const [dismissedError, setDismissedError] = useState<string | null>(null);

  const [previewModel, setPreviewModel] = useState("coding-fast");
  const [previewBody, setPreviewBody] = useState(
    '{"model":"coding-fast","temperature":0.9,"messages":[{"role":"user","content":"Write a Go HTTP handler"}]}',
  );
  const [previewError, setPreviewError] = useState<string | null>(null);

  const openHelp = useCallback((topic: ParameterGuideId) => {
    setHelpTopic(topic);
    setHelpOpen(true);
  }, []);

  const closeHelp = useCallback(() => {
    setHelpOpen(false);
  }, []);

  const overrideForm = useForm<CreateModelOverrideFormValues>({
    resolver: zodResolver(createModelOverrideSchema),
    defaultValues: {
      name: "",
      model_pattern: "",
      target_model: "",
      temperature: "",
      max_tokens: "",
      extra_params_json: "{}",
      strip_params_csv: "",
    },
  });

  const handleCreateOverride = overrideForm.handleSubmit((values) => {
    const extra = parseExtraParamsJson(values.extra_params_json ?? "");
    if (!extra.ok) return;
    setDismissedError(null);
    createOverride.mutate(
      {
        name: values.name,
        model_pattern: values.model_pattern,
        target_model: values.target_model || undefined,
        temperature: values.temperature?.trim() ? Number(values.temperature) : undefined,
        max_tokens: values.max_tokens?.trim() ? Number(values.max_tokens) : undefined,
        extra_params: extra.value,
        strip_params: parseStripParamsList(values.strip_params_csv ?? ""),
        priority: 100,
      },
      {
        onSuccess: () => {
          overrideForm.reset();
        },
      },
    );
  });

  const handleRunPreview = () => {
    setPreviewError(null);
    setDismissedError(null);
    const parsed = parseExtraParamsJson(previewBody);
    if (!parsed.ok) {
      setPreviewError(parsed.error);
      return;
    }
    previewMutation.mutate({
      protocol: "openai-chat",
      model: previewModel.trim() || "gpt-4o",
      body: parsed.value,
    });
  };

  const rawError =
    overridesQuery.error?.message ??
    createOverride.error?.message ??
    updateOverride.error?.message ??
    deleteOverride.error?.message ??
    previewMutation.error?.message ??
    previewError;

  const activeError = rawError && rawError !== dismissedError ? rawError : null;

  const overrideFirstErr =
    overrideForm.formState.errors.name?.message ??
    overrideForm.formState.errors.model_pattern?.message ??
    overrideForm.formState.errors.temperature?.message ??
    overrideForm.formState.errors.max_tokens?.message ??
    overrideForm.formState.errors.extra_params_json?.message ??
    null;

  const overrides = overridesQuery.data?.overrides ?? [];

  return (
    <div className="flex h-[calc(100dvh-7.25rem)] flex-col gap-3 overflow-hidden">
      {/* Compact top toolbar */}
      <div className="flex shrink-0 items-center justify-between gap-2">
        <p className="min-w-0 truncate text-xs text-[var(--text-secondary)]">{copy.subtitle}</p>
        <div className="flex shrink-0 items-center gap-1.5">
          <Button type="button" variant="outline" size="xs" onClick={() => openHelp("temperature")}>
            <CircleHelp className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{copy.guideButton}</span>
          </Button>
          <Button
            type="button"
            variant="outline"
            size="xs"
            onClick={() => {
              setDismissedError(null);
              void overridesQuery.refetch();
            }}
            disabled={overridesQuery.isFetching}
            aria-label={copy.refreshButton}
          >
            <RefreshCw
              className={cn("h-3.5 w-3.5", overridesQuery.isFetching && "animate-spin")}
              aria-hidden="true"
            />
            <span className="hidden sm:inline">{copy.refreshButton}</span>
          </Button>
        </div>
      </div>

      {/* Dismissible compact error banner */}
      {activeError ? (
        <Alert variant="destructive" className="shrink-0 py-2">
          <Sliders className="h-4 w-4" aria-hidden="true" />
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <AlertTitle className="text-xs">{copy.errorTitle}</AlertTitle>
              <AlertDescription className="truncate text-xs">{activeError}</AlertDescription>
            </div>
            <button
              type="button"
              onClick={() => setDismissedError(activeError)}
              aria-label={copy.dismissLabel}
              className="shrink-0 cursor-pointer rounded p-0.5 text-[var(--text-muted)] hover:text-[var(--text-primary)]"
            >
              <X className="h-3.5 w-3.5" aria-hidden="true" />
            </button>
          </div>
        </Alert>
      ) : null}

      {/* Mobile-first segmented tabs (hidden on lg+) */}
      <div
        role="tablist"
        aria-label="Model Editor sections"
        className="grid shrink-0 grid-cols-3 gap-1 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-1 lg:hidden"
      >
        {(
          [
            { id: "override", label: copy.tabOverride },
            { id: "rules", label: `${copy.tabRules} (${overrides.length})` },
            { id: "dryrun", label: copy.tabDryRun },
          ] as const
        ).map((tab) => {
          const active = mobileTab === tab.id;
          return (
            <button
              key={tab.id}
              type="button"
              role="tab"
              aria-selected={active}
              onClick={() => setMobileTab(tab.id)}
              className={cn(
                "h-8 cursor-pointer truncate rounded px-2 text-xs font-medium transition-colors",
                active
                  ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)] shadow-xs"
                  : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]",
              )}
            >
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Main 1-screen split workspace */}
      <div className="grid min-h-0 flex-1 grid-cols-1 gap-3 overflow-hidden lg:grid-cols-2 lg:gap-4">
        {/* Left column: New Model Override + Active Overrides */}
        <Card
          className={cn(
            "min-h-0 flex-col overflow-hidden",
            mobileTab === "dryrun" ? "hidden lg:flex" : "flex",
          )}
        >
          {/* New Model Override form */}
          <div
            className={cn(
              "shrink-0 flex-col p-3.5 sm:p-4",
              mobileTab === "rules" ? "hidden lg:flex" : "flex",
            )}
          >
            <div className="mb-2.5 flex items-center justify-between gap-2">
              <h2 className="text-sm font-semibold tracking-tight text-[var(--text-primary)]">
                {copy.newOverrideTitle}
              </h2>
              <Badge variant="outline" className="font-mono text-[10px]">
                POST /overrides
              </Badge>
            </div>

            <form onSubmit={handleCreateOverride} className="space-y-2.5">
              <div className="grid grid-cols-2 gap-2.5">
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-name"
                    label={copy.nameLabel}
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-name"
                    placeholder={copy.namePlaceholder}
                    className="h-9 text-xs lg:h-8"
                    {...overrideForm.register("name")}
                  />
                </div>
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-pattern"
                    label={copy.patternLabel}
                    helpTopic="model_pattern"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-pattern"
                    placeholder={copy.patternPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("model_pattern")}
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-2">
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-target"
                    label={copy.targetLabel}
                    helpTopic="target_model"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-target"
                    placeholder={copy.targetPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("target_model")}
                  />
                </div>
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-temp"
                    label={copy.tempLabel}
                    helpTopic="temperature"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-temp"
                    inputMode="decimal"
                    placeholder={copy.tempPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("temperature")}
                  />
                </div>
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-max"
                    label={copy.maxTokensLabel}
                    helpTopic="max_tokens"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-max"
                    inputMode="numeric"
                    placeholder={copy.maxTokensPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("max_tokens")}
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-strip"
                    label={copy.stripLabel}
                    helpTopic="strip_params"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-strip"
                    placeholder={copy.stripPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("strip_params_csv")}
                  />
                </div>
                <div className="space-y-1">
                  <FieldLabel
                    htmlFor="mo-extra"
                    label={copy.extraLabel}
                    helpTopic="extra_params"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <Input
                    id="mo-extra"
                    placeholder={copy.extraPlaceholder}
                    className="h-9 font-mono text-xs lg:h-8"
                    {...overrideForm.register("extra_params_json")}
                  />
                </div>
              </div>

              {overrideFirstErr ? (
                <p role="alert" className="text-xs text-[var(--status-error)]">
                  {overrideFirstErr}
                </p>
              ) : null}

              <Button
                type="submit"
                size="sm"
                className="w-full sm:w-auto"
                disabled={createOverride.isPending}
              >
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                {createOverride.isPending ? copy.creatingButton : copy.createButton}
              </Button>
            </form>
          </div>

          {/* Active Overrides list (internal scroll only) */}
          <div
            className={cn(
              "min-h-0 flex-1 flex-col overflow-hidden lg:border-t lg:border-[var(--border-subtle)]",
              mobileTab === "rules" ? "flex" : "hidden lg:flex",
            )}
          >
            <div className="flex shrink-0 items-center justify-between px-3.5 py-2 sm:px-4">
              <span className="text-xs font-semibold text-[var(--text-secondary)]">
                {copy.activeOverridesTitle}
              </span>
              <Badge variant="secondary" className="font-mono text-[10px]">
                {overrides.length}
              </Badge>
            </div>

            <div className="min-h-0 flex-1 space-y-2 overflow-y-auto px-3.5 pb-3 sm:px-4">
              {overridesQuery.isLoading ? (
                <p className="py-4 text-xs text-[var(--text-muted)]">{copy.loadingOverrides}</p>
              ) : overrides.length === 0 ? (
                <p className="rounded-md border border-dashed border-[var(--border-subtle)] px-3 py-4 text-xs text-[var(--text-muted)]">
                  {copy.emptyOverrides}
                </p>
              ) : (
                overrides.map((override) => (
                  <div
                    key={override.id}
                    className="flex items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2 text-xs"
                  >
                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center gap-1.5">
                        <span className="truncate font-medium text-[var(--text-primary)]">
                          {override.name}
                        </span>
                        <Badge variant="outline" className="font-mono text-[10px]">
                          {override.model_pattern}
                        </Badge>
                        {override.target_model ? (
                          <Badge variant="secondary" className="font-mono text-[10px]">
                            -&gt; {override.target_model}
                          </Badge>
                        ) : null}
                        {override.temperature !== undefined ? (
                          <span className="font-mono text-[10px] text-[var(--text-muted)]">
                            T={override.temperature}
                          </span>
                        ) : null}
                        {override.max_tokens !== undefined ? (
                          <span className="font-mono text-[10px] text-[var(--text-muted)]">
                            max={override.max_tokens}
                          </span>
                        ) : null}
                      </div>
                    </div>
                    <div className="flex shrink-0 items-center gap-1.5">
                      <Button
                        type="button"
                        size="xs"
                        variant={override.enabled ? "outline" : "secondary"}
                        disabled={updateOverride.isPending}
                        onClick={() =>
                          updateOverride.mutate({
                            id: override.id,
                            body: { enabled: !override.enabled },
                          })
                        }
                      >
                        {override.enabled ? copy.disableLabel : copy.enableLabel}
                      </Button>
                      <Button
                        type="button"
                        size="icon-sm"
                        variant="outline"
                        disabled={deleteOverride.isPending}
                        aria-label={`${copy.deleteAriaLabel} ${override.name}`}
                        onClick={() => deleteOverride.mutate(override.id)}
                      >
                        <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                      </Button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </Card>

        {/* Right column: Dry-Run Preview */}
        <Card
          className={cn(
            "min-h-0 flex-col overflow-hidden",
            mobileTab === "dryrun" ? "flex" : "hidden lg:flex",
          )}
        >
          <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] px-3.5 py-2.5 sm:px-4">
            <div className="flex items-center gap-1.5">
              <h2 className="text-sm font-semibold tracking-tight text-[var(--text-primary)]">
                {copy.dryRunTitle}
              </h2>
              <button
                type="button"
                onClick={() => openHelp("dry_run")}
                aria-label={`${copy.helpAriaPrefix} ${copy.dryRunTitle}`}
                className="inline-flex h-4 w-4 cursor-pointer items-center justify-center rounded-full text-[var(--text-muted)] transition-colors hover:text-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-[var(--accent-primary)]"
              >
                <CircleHelp className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
            </div>
            <Badge variant="outline" className="font-mono text-[10px]">
              openai-chat
            </Badge>
          </div>

          <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-hidden p-3.5 sm:p-4">
            <div className="flex shrink-0 items-end gap-2">
              <div className="min-w-0 flex-1 space-y-1">
                <Label htmlFor="me-prev-model" className="text-xs">
                  {copy.previewModelLabel}
                </Label>
                <Input
                  id="me-prev-model"
                  value={previewModel}
                  onChange={(e) => setPreviewModel(e.target.value)}
                  className="h-9 font-mono text-xs lg:h-8"
                />
              </div>
              <Button
                type="button"
                size="sm"
                className="h-9 shrink-0 lg:h-8"
                disabled={previewMutation.isPending}
                onClick={handleRunPreview}
              >
                <Play className="h-3.5 w-3.5" aria-hidden="true" />
                {previewMutation.isPending ? copy.runningPreviewButton : copy.runPreviewButton}
              </Button>
            </div>

            <div className="flex shrink-0 flex-col gap-1">
              <Label htmlFor="me-prev-body" className="text-xs">
                {copy.previewBodyLabel}
              </Label>
              <textarea
                id="me-prev-body"
                rows={4}
                value={previewBody}
                onChange={(e) => setPreviewBody(e.target.value)}
                className="w-full resize-none rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] p-2.5 font-mono text-xs text-[var(--text-primary)] focus:border-[var(--accent-primary)] focus:outline-none"
              />
            </div>

            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
              <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] px-3 py-1.5">
                <span className="font-mono text-[11px] font-medium text-[var(--text-secondary)]">
                  {copy.mutatedOutputLabel}
                </span>
                {previewMutation.data?.model ? (
                  <Badge variant="secondary" className="font-mono text-[10px]">
                    {previewMutation.data.model}
                  </Badge>
                ) : null}
              </div>
              <pre className="min-h-0 flex-1 overflow-auto p-3 font-mono text-xs leading-relaxed text-[var(--text-primary)]">
                {previewMutation.data
                  ? JSON.stringify(previewMutation.data, null, 2)
                  : copy.previewPlaceholder}
              </pre>
            </div>
          </div>
        </Card>
      </div>

      <ModelEditorHelpModal
        open={helpOpen}
        activeTopic={helpTopic}
        locale={locale}
        onSelectTopic={setHelpTopic}
        onClose={closeHelp}
      />
    </div>
  );
}
