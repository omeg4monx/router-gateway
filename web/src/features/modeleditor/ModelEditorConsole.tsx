import { zodResolver } from "@hookform/resolvers/zod";
import { AlertTriangle, CircleHelp, Play, Plus, RefreshCw, Sliders, Trash2, X } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useModelCatalog } from "@/features/models/hooks";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { BUILTIN_PROMPT_TEMPLATES } from "./builtinPromptTemplates";
import {
  useCreateModelOverrideMutation,
  useDeleteModelOverrideMutation,
  useModelOverrides,
  useRunQuickPromptTestMutation,
  useUpdateModelOverrideMutation,
} from "./hooks";
import { ModelEditorHelpModal } from "./ModelEditorHelpModal";
import { PromptTemplateCombobox } from "./PromptTemplateCombobox";
import {
  classifyPromptTestOutcome,
  estimatePromptTokens,
  MODEL_EDITOR_COPY,
  type ParameterGuideId,
  type PromptTestVerdict,
  QUICK_TEST_PRESETS,
  type QuickTestPresetId,
} from "./parameterGuide";
import {
  type CreateModelOverrideFormValues,
  createModelOverrideSchema,
  parseExtraParamsJson,
  parseStripParamsList,
} from "./schemas";
import { buildTargetModelOptions, TargetModelCombobox } from "./TargetModelCombobox";

type MobileTab = "override" | "rules" | "test";

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
  const modelCatalogQuery = useModelCatalog();
  const createOverride = useCreateModelOverrideMutation();
  const updateOverride = useUpdateModelOverrideMutation();
  const deleteOverride = useDeleteModelOverrideMutation();
  const quickTestMutation = useRunQuickPromptTestMutation();

  const targetModelOptions = useMemo(
    () => buildTargetModelOptions(modelCatalogQuery.data ?? []),
    [modelCatalogQuery.data],
  );

  const [mobileTab, setMobileTab] = useState<MobileTab>("override");
  const [helpOpen, setHelpOpen] = useState(false);
  const [helpTopic, setHelpTopic] = useState<ParameterGuideId>("temperature");
  const [dismissedError, setDismissedError] = useState<string | null>(null);

  const [activePresetId, setActivePresetId] = useState<QuickTestPresetId>(QUICK_TEST_PRESETS[0].id);
  const [selectedTemplateId, setSelectedTemplateId] = useState<string>("custom");
  const [customTestModel, setCustomTestModel] = useState("");
  const [testPrompt, setTestPrompt] = useState(QUICK_TEST_PRESETS[0].prompt[locale]);
  const [testValidationError, setTestValidationError] = useState<string | null>(null);

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
      target_model: "",
      system_prompt: "",
      system_prompt_position: "suffix",
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
    const trimmedPrompt = values.system_prompt?.trim() || undefined;
    const isDual = values.system_prompt_position === "dual" && Boolean(trimmedPrompt);
    const extraParams = isDual ? { ...extra.value, _prompt_mode: "dual" } : extra.value;
    createOverride.mutate(
      {
        name: values.name,
        model_pattern: values.name,
        target_model: values.target_model,
        system_prompt_prefix:
          values.system_prompt_position === "prefix" ? trimmedPrompt : undefined,
        system_prompt_suffix:
          values.system_prompt_position === "suffix" || isDual ? trimmedPrompt : undefined,
        temperature: values.temperature?.trim() ? Number(values.temperature) : undefined,
        max_tokens: values.max_tokens?.trim() ? Number(values.max_tokens) : undefined,
        extra_params: extraParams,
        strip_params: parseStripParamsList(values.strip_params_csv ?? ""),
        priority: 100,
      },
      {
        onSuccess: () => {
          overrideForm.reset();
          setSelectedTemplateId("custom");
        },
      },
    );
  });

  const targetModelValue = overrideForm.watch("target_model") ?? "";
  const selectedTargetModel = useMemo(
    () => targetModelOptions.find((opt) => opt.id === targetModelValue),
    [targetModelOptions, targetModelValue],
  );
  const systemPromptValue = overrideForm.watch("system_prompt") ?? "";
  const promptPosition = overrideForm.watch("system_prompt_position");
  const promptTokenEstimate = estimatePromptTokens(systemPromptValue);
  const effectiveTestModel = customTestModel || targetModelValue;

  const handleSelectPreset = (presetId: QuickTestPresetId) => {
    setActivePresetId(presetId);
    const preset = QUICK_TEST_PRESETS.find((item) => item.id === presetId) ?? QUICK_TEST_PRESETS[0];
    setTestPrompt(preset.prompt[locale]);
    setTestValidationError(null);
  };

  const handleRunQuickTest = () => {
    setTestValidationError(null);
    setDismissedError(null);
    const model = effectiveTestModel.trim();
    if (!model) {
      setTestValidationError(copy.quickTestSelectModelError);
      return;
    }
    const prompt = testPrompt.trim();
    if (!prompt) return;

    const systemPrompt = overrideForm.getValues("system_prompt")?.trim() || undefined;
    const tempRaw = overrideForm.getValues("temperature")?.trim();
    const maxRaw = overrideForm.getValues("max_tokens")?.trim();
    const temperature = tempRaw && Number.isFinite(Number(tempRaw)) ? Number(tempRaw) : undefined;
    const max_tokens = maxRaw && Number.isFinite(Number(maxRaw)) ? Number(maxRaw) : undefined;

    quickTestMutation.mutate({
      model,
      prompt,
      system_prompt: systemPrompt,
      temperature,
      max_tokens,
    });
  };

  const rawError =
    overridesQuery.error?.message ??
    modelCatalogQuery.error?.message ??
    createOverride.error?.message ??
    updateOverride.error?.message ??
    deleteOverride.error?.message ??
    quickTestMutation.error?.message ??
    testValidationError;

  const activeError = rawError && rawError !== dismissedError ? rawError : null;

  const overrideFirstErr =
    overrideForm.formState.errors.name?.message ??
    overrideForm.formState.errors.target_model?.message ??
    overrideForm.formState.errors.temperature?.message ??
    overrideForm.formState.errors.max_tokens?.message ??
    overrideForm.formState.errors.extra_params_json?.message ??
    null;

  const overrides = overridesQuery.data?.overrides ?? [];

  const testOutcome =
    quickTestMutation.isError || quickTestMutation.data
      ? classifyPromptTestOutcome(
          {
            content: quickTestMutation.data?.content,
            error: quickTestMutation.error?.message ?? null,
            completionTokens: quickTestMutation.data?.completionTokens,
          },
          locale,
        )
      : null;

  const verdictBadgeClass: Record<PromptTestVerdict, string> = {
    answered: "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]",
    refusal:
      "border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] text-[var(--status-warn)]",
    empty: "border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-muted)]",
    error:
      "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]",
  };

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
              void Promise.all([overridesQuery.refetch(), modelCatalogQuery.refetch()]);
            }}
            disabled={overridesQuery.isFetching || modelCatalogQuery.isFetching}
            aria-label={copy.refreshButton}
          >
            <RefreshCw
              className={cn(
                "h-3.5 w-3.5",
                (overridesQuery.isFetching || modelCatalogQuery.isFetching) && "animate-spin",
              )}
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
            { id: "test", label: copy.tabQuickTest },
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
            mobileTab === "test" ? "hidden lg:flex" : "flex",
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
              <div className="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
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
                    htmlFor="mo-target"
                    label={copy.targetLabel}
                    helpTopic="target_model"
                    helpAriaPrefix={copy.helpAriaPrefix}
                    onOpenHelp={openHelp}
                  />
                  <TargetModelCombobox
                    id="mo-target"
                    value={targetModelValue}
                    options={targetModelOptions}
                    isLoading={modelCatalogQuery.isLoading}
                    locale={locale}
                    onChange={(nextValue) =>
                      overrideForm.setValue("target_model", nextValue, {
                        shouldDirty: true,
                        shouldValidate: true,
                      })
                    }
                  />
                </div>
              </div>

              <div className="space-y-1">
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-1.5">
                    <Label htmlFor="mo-sys-prompt" className="text-xs">
                      {copy.systemPromptLabel}
                    </Label>
                    <button
                      type="button"
                      onClick={() => openHelp("system_prompt")}
                      aria-label={`${copy.helpAriaPrefix} ${copy.systemPromptLabel}`}
                      className="inline-flex h-4 w-4 shrink-0 cursor-pointer items-center justify-center rounded-full text-[var(--text-muted)] transition-colors hover:text-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-[var(--accent-primary)]"
                    >
                      <CircleHelp className="h-3.5 w-3.5" aria-hidden="true" />
                    </button>
                    {promptTokenEstimate > 0 ? (
                      <span className="font-mono text-[10px] text-[var(--text-muted)]">
                        ~{promptTokenEstimate} tok
                      </span>
                    ) : null}
                  </div>
                  <div className="inline-flex rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-0.5 text-[10px]">
                    {(["suffix", "dual", "prefix"] as const).map((pos) => (
                      <button
                        key={pos}
                        type="button"
                        onClick={() => overrideForm.setValue("system_prompt_position", pos)}
                        className={cn(
                          "cursor-pointer rounded px-1.5 py-0.5 font-mono font-medium transition-colors",
                          promptPosition === pos
                            ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                            : "text-[var(--text-muted)] hover:text-[var(--text-primary)]",
                        )}
                      >
                        {pos === "suffix"
                          ? copy.promptSuffixLabel
                          : pos === "dual"
                            ? copy.promptDualLabel
                            : copy.promptPrefixLabel}
                      </button>
                    ))}
                  </div>
                </div>
                <div className="flex items-center justify-between gap-2 rounded border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-2 py-1 text-[10.5px] text-[var(--status-warn)]">
                  <div className="flex min-w-0 items-center gap-1.5">
                    <AlertTriangle className="h-3 w-3 shrink-0" aria-hidden="true" />
                    <span className="truncate">{copy.promptDisclaimerWarning}</span>
                  </div>
                  <button
                    type="button"
                    onClick={() => openHelp("disclaimer")}
                    className="shrink-0 cursor-pointer font-semibold underline underline-offset-2 hover:opacity-80"
                  >
                    {copy.promptDisclaimerAction}
                  </button>
                </div>
                <PromptTemplateCombobox
                  selectedId={selectedTemplateId}
                  templates={BUILTIN_PROMPT_TEMPLATES}
                  locale={locale}
                  onSelectCustom={() => {
                    setSelectedTemplateId("custom");
                    overrideForm.setValue("system_prompt", "", { shouldDirty: true });
                  }}
                  onSelectTemplate={(tpl) => {
                    setSelectedTemplateId(tpl.id);
                    overrideForm.setValue("system_prompt", tpl.content, { shouldDirty: true });
                    overrideForm.setValue("system_prompt_position", tpl.position, {
                      shouldDirty: true,
                    });
                  }}
                />
                <textarea
                  id="mo-sys-prompt"
                  rows={2}
                  placeholder={copy.systemPromptPlaceholder}
                  className="w-full resize-none rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-1.5 text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:border-[var(--accent-primary)] focus:outline-none"
                  {...overrideForm.register("system_prompt")}
                />
              </div>

              <div className="grid grid-cols-2 gap-2.5">
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
                  <div className="flex items-center justify-between gap-1">
                    <div className="flex min-w-0 items-center gap-1.5">
                      <Label htmlFor="mo-max" className="truncate text-xs">
                        {copy.maxTokensLabel}
                      </Label>
                      {selectedTargetModel?.maxOutputTokens ? (
                        <button
                          type="button"
                          onClick={() =>
                            overrideForm.setValue(
                              "max_tokens",
                              String(selectedTargetModel.maxOutputTokens),
                              { shouldDirty: true, shouldValidate: true },
                            )
                          }
                          className="shrink-0 cursor-pointer font-mono text-[10px] text-[var(--accent-text)] hover:underline"
                          title={
                            locale === "id"
                              ? "Klik untuk isi batas maks model"
                              : "Click to fill model max limit"
                          }
                        >
                          max {selectedTargetModel.maxOutputTokens}
                        </button>
                      ) : null}
                    </div>
                    <button
                      type="button"
                      onClick={() => openHelp("max_tokens")}
                      aria-label={`${copy.helpAriaPrefix} ${copy.maxTokensLabel}`}
                      className="inline-flex h-4 w-4 shrink-0 cursor-pointer items-center justify-center rounded-full text-[var(--text-muted)] transition-colors hover:text-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-[var(--accent-primary)]"
                    >
                      <CircleHelp className="h-3.5 w-3.5" aria-hidden="true" />
                    </button>
                  </div>
                  <Input
                    id="mo-max"
                    inputMode="numeric"
                    placeholder={
                      selectedTargetModel?.maxOutputTokens
                        ? String(Math.min(4096, selectedTargetModel.maxOutputTokens))
                        : copy.maxTokensPlaceholder
                    }
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
                        {override.extra_params?._prompt_mode === "dual" ? (
                          <Badge variant="secondary" className="font-mono text-[10px]">
                            dual
                          </Badge>
                        ) : (
                          <>
                            {override.system_prompt_suffix ? (
                              <Badge variant="secondary" className="font-mono text-[10px]">
                                append
                              </Badge>
                            ) : null}
                            {override.system_prompt_prefix ? (
                              <Badge variant="secondary" className="font-mono text-[10px]">
                                prepend
                              </Badge>
                            ) : null}
                          </>
                        )}
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

        {/* Right column: Quick Prompt Test (4 presets) */}
        <Card
          className={cn(
            "min-h-0 flex-col overflow-hidden",
            mobileTab === "test" ? "flex" : "hidden lg:flex",
          )}
        >
          <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] px-3.5 py-2.5 sm:px-4">
            <h2 className="text-sm font-semibold tracking-tight text-[var(--text-primary)]">
              {copy.quickTestTitle}
            </h2>
            <Badge variant="outline" className="font-mono text-[10px]">
              4 Presets
            </Badge>
          </div>

          <div className="flex min-h-0 flex-1 flex-col gap-2.5 overflow-hidden p-3.5 sm:p-4">
            {/* Test model selector + Run button */}
            <div className="flex shrink-0 items-end gap-2">
              <div className="min-w-0 flex-1 space-y-1">
                <Label htmlFor="me-test-model" className="text-xs">
                  {copy.quickTestModelLabel}
                </Label>
                <TargetModelCombobox
                  id="me-test-model"
                  value={effectiveTestModel}
                  options={targetModelOptions}
                  isLoading={modelCatalogQuery.isLoading}
                  locale={locale}
                  onChange={(nextValue) => {
                    setCustomTestModel(nextValue);
                    setTestValidationError(null);
                  }}
                />
              </div>
              <Button
                type="button"
                size="sm"
                className="h-9 shrink-0"
                disabled={quickTestMutation.isPending || !testPrompt.trim()}
                onClick={handleRunQuickTest}
              >
                <Play className="h-3.5 w-3.5" aria-hidden="true" />
                {quickTestMutation.isPending
                  ? copy.quickTestRunningButton
                  : copy.quickTestRunButton}
              </Button>
            </div>

            {/* 4 preset prompt buttons */}
            <div className="grid shrink-0 grid-cols-2 gap-1.5">
              {QUICK_TEST_PRESETS.map((preset) => {
                const isSelected = preset.id === activePresetId;
                return (
                  <button
                    key={preset.id}
                    type="button"
                    onClick={() => handleSelectPreset(preset.id)}
                    className={cn(
                      "cursor-pointer truncate rounded-md border px-2.5 py-1.5 text-left text-xs font-medium transition-colors",
                      isSelected
                        ? "border-[var(--accent-primary)] bg-[var(--accent-primary)]/12 text-[var(--accent-text)]"
                        : "border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]",
                    )}
                  >
                    {preset.label[locale]}
                  </button>
                );
              })}
            </div>

            {/* Compact editable prompt box */}
            <div className="flex shrink-0 flex-col gap-1">
              <Label htmlFor="me-test-prompt" className="text-xs">
                {copy.quickTestPromptLabel}
              </Label>
              <textarea
                id="me-test-prompt"
                rows={3}
                value={testPrompt}
                onChange={(event) => setTestPrompt(event.target.value)}
                className="w-full resize-none rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] p-2.5 text-xs text-[var(--text-primary)] focus:border-[var(--accent-primary)] focus:outline-none"
              />
            </div>

            {/* Scrollable model response output */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
              <div className="flex shrink-0 flex-wrap items-center justify-between gap-1.5 border-b border-[var(--border-subtle)] px-3 py-1.5">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-[11px] font-medium text-[var(--text-secondary)]">
                    {copy.quickTestOutputLabel}
                  </span>
                  {testOutcome && !quickTestMutation.isPending ? (
                    <span
                      className={cn(
                        "inline-flex items-center rounded border px-1.5 py-0.5 font-mono text-[10px] font-semibold",
                        verdictBadgeClass[testOutcome.verdict],
                      )}
                    >
                      {testOutcome.badgeLabel}
                    </span>
                  ) : null}
                </div>
                {quickTestMutation.data ? (
                  <div className="flex items-center gap-1.5">
                    <Badge variant="secondary" className="font-mono text-[10px]">
                      {quickTestMutation.data.model}
                    </Badge>
                    <Badge variant="outline" className="font-mono text-[10px]">
                      {quickTestMutation.data.latencyMs}ms
                    </Badge>
                  </div>
                ) : null}
              </div>

              {testOutcome && !quickTestMutation.isPending ? (
                <div className="shrink-0 border-b border-[var(--border-subtle)] bg-[var(--bg-surface)]/70 px-3 py-1.5 text-[11px] text-[var(--text-secondary)]">
                  {testOutcome.summary}
                </div>
              ) : null}

              <div className="min-h-0 flex-1 overflow-y-auto whitespace-pre-wrap p-3 text-xs leading-relaxed text-[var(--text-primary)]">
                {quickTestMutation.isPending
                  ? copy.quickTestRunningButton
                  : quickTestMutation.isError
                    ? quickTestMutation.error.message
                    : quickTestMutation.data?.content || copy.quickTestPlaceholder}
              </div>
            </div>

            {/* Compact technical disclaimer strip */}
            <div className="flex shrink-0 items-center justify-between gap-2 rounded-md border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-2.5 py-1.5 text-[11px] text-[var(--status-warn)]">
              <div className="flex min-w-0 items-center gap-1.5">
                <AlertTriangle className="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                <span className="truncate">{copy.disclaimerBannerText}</span>
              </div>
              <button
                type="button"
                onClick={() => openHelp("disclaimer")}
                className="shrink-0 cursor-pointer font-semibold underline underline-offset-2 hover:opacity-80"
              >
                {copy.disclaimerBannerAction}
              </button>
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
