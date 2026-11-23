"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Play, Plus, RefreshCw, Sliders, Trash2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useCreateBypassRuleMutation,
  useCreateModelOverrideMutation,
  useDeleteBypassRuleMutation,
  useDeleteModelOverrideMutation,
  useModelBypassRules,
  useModelOverrides,
  usePreviewModelEditorMutation,
  useUpdateBypassRuleMutation,
  useUpdateModelOverrideMutation,
} from "./hooks";
import {
  BYPASS_ACTIONS,
  type CreateBypassRuleFormValues,
  type CreateModelOverrideFormValues,
  createBypassRuleSchema,
  createModelOverrideSchema,
  parseExtraParamsJson,
  parseStripParamsList,
} from "./schemas";

export function ModelEditorConsole() {
  const overridesQuery = useModelOverrides();
  const bypassQuery = useModelBypassRules();

  const createOverride = useCreateModelOverrideMutation();
  const updateOverride = useUpdateModelOverrideMutation();
  const deleteOverride = useDeleteModelOverrideMutation();
  const createBypass = useCreateBypassRuleMutation();
  const updateBypass = useUpdateBypassRuleMutation();
  const deleteBypass = useDeleteBypassRuleMutation();
  const previewMutation = usePreviewModelEditorMutation();

  const [previewModel, setPreviewModel] = useState("coding-fast");
  const [previewBody, setPreviewBody] = useState(
    '{"model":"coding-fast","temperature":0.9,"messages":[{"role":"user","content":"Write a Go HTTP handler"}]}',
  );
  const [previewError, setPreviewError] = useState<string | null>(null);

  const overrideForm = useForm<CreateModelOverrideFormValues>({
    resolver: zodResolver(createModelOverrideSchema),
    defaultValues: {
      name: "",
      model_pattern: "",
      target_model: "",
      system_prompt_prefix: "",
      system_prompt_suffix: "",
      temperature: "",
      max_tokens: "",
      top_p: "",
      extra_params_json: "{}",
      strip_params_csv: "",
      priority: "100",
    },
  });

  const bypassForm = useForm<CreateBypassRuleFormValues>({
    resolver: zodResolver(createBypassRuleSchema),
    defaultValues: {
      name: "",
      model_pattern: "*",
      key_id: "",
      protocol: "",
      header_name: "",
      header_value: "",
      action: "bypass_tokensaver",
      target_model: "",
      priority: "100",
    },
  });

  const handleCreateOverride = overrideForm.handleSubmit((values) => {
    const extra = parseExtraParamsJson(values.extra_params_json ?? "");
    if (!extra.ok) return;
    createOverride.mutate(
      {
        name: values.name,
        model_pattern: values.model_pattern,
        target_model: values.target_model || undefined,
        system_prompt_prefix: values.system_prompt_prefix || undefined,
        system_prompt_suffix: values.system_prompt_suffix || undefined,
        temperature: values.temperature?.trim() ? Number(values.temperature) : undefined,
        max_tokens: values.max_tokens?.trim() ? Number(values.max_tokens) : undefined,
        top_p: values.top_p?.trim() ? Number(values.top_p) : undefined,
        extra_params: extra.value,
        strip_params: parseStripParamsList(values.strip_params_csv ?? ""),
        priority: values.priority?.trim() ? Number(values.priority) : 100,
      },
      { onSuccess: () => overrideForm.reset() },
    );
  });

  const handleCreateBypass = bypassForm.handleSubmit((values) => {
    createBypass.mutate(
      {
        name: values.name,
        model_pattern: values.model_pattern,
        key_id: values.key_id || undefined,
        protocol: values.protocol || undefined,
        header_name: values.header_name || undefined,
        header_value: values.header_value || undefined,
        action: values.action,
        target_model: values.target_model || undefined,
        priority: values.priority?.trim() ? Number(values.priority) : 100,
      },
      { onSuccess: () => bypassForm.reset() },
    );
  });

  const handleRunPreview = () => {
    setPreviewError(null);
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

  const activeError =
    overridesQuery.error?.message ??
    bypassQuery.error?.message ??
    createOverride.error?.message ??
    updateOverride.error?.message ??
    deleteOverride.error?.message ??
    createBypass.error?.message ??
    updateBypass.error?.message ??
    deleteBypass.error?.message ??
    previewMutation.error?.message ??
    previewError;

  const overrideFirstErr =
    overrideForm.formState.errors.name?.message ??
    overrideForm.formState.errors.model_pattern?.message ??
    overrideForm.formState.errors.temperature?.message ??
    overrideForm.formState.errors.max_tokens?.message ??
    overrideForm.formState.errors.top_p?.message ??
    overrideForm.formState.errors.extra_params_json?.message ??
    overrideForm.formState.errors.priority?.message ??
    null;

  const bypassFirstErr =
    bypassForm.formState.errors.name?.message ??
    bypassForm.formState.errors.model_pattern?.message ??
    bypassForm.formState.errors.target_model?.message ??
    bypassForm.formState.errors.priority?.message ??
    null;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="font-mono text-xs uppercase tracking-[0.18em] text-[var(--accent-primary)]">
            Gateway Interceptor · Priority 20
          </p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Model Editor & Bypass
          </h1>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Shape model parameters, wrap system prompts, strip unsupported fields, and define
            conditional interceptor bypass rules.
          </p>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => {
            void Promise.all([overridesQuery.refetch(), bypassQuery.refetch()]);
          }}
        >
          <RefreshCw className="mr-2 h-4 w-4" />
          Refresh
        </Button>
      </div>

      <EmberRule />

      {activeError ? (
        <Alert variant="destructive">
          <Sliders className="h-4 w-4" />
          <AlertTitle>Model Editor Error</AlertTitle>
          <AlertDescription>{activeError}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="p-6">
          <CardHeader className="p-0 pb-4">
            <CardTitle className="text-base">New Model Override</CardTitle>
          </CardHeader>
          <form onSubmit={handleCreateOverride} className="space-y-3">
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="mo-name">Name</Label>
                <Input
                  id="mo-name"
                  placeholder="Reasoning Safe"
                  {...overrideForm.register("name")}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mo-pattern">Model Pattern</Label>
                <Input
                  id="mo-pattern"
                  placeholder="coding-* or o3-*"
                  className="font-mono"
                  {...overrideForm.register("model_pattern")}
                />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-3">
              <div className="space-y-1">
                <Label htmlFor="mo-target">Target Model</Label>
                <Input
                  id="mo-target"
                  placeholder="optional rewrite"
                  className="font-mono"
                  {...overrideForm.register("target_model")}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mo-temp">Temperature</Label>
                <Input id="mo-temp" placeholder="0.2" {...overrideForm.register("temperature")} />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mo-max">Max Tokens</Label>
                <Input id="mo-max" placeholder="4096" {...overrideForm.register("max_tokens")} />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="mo-strip">Strip Params (CSV)</Label>
                <Input
                  id="mo-strip"
                  placeholder="temperature, presence_penalty"
                  className="font-mono"
                  {...overrideForm.register("strip_params_csv")}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mo-extra">Extra Params JSON</Label>
                <Input
                  id="mo-extra"
                  placeholder='{"seed": 42}'
                  className="font-mono"
                  {...overrideForm.register("extra_params_json")}
                />
              </div>
            </div>
            {overrideFirstErr ? (
              <p className="text-xs text-[var(--accent-primary)]">{overrideFirstErr}</p>
            ) : null}
            <Button type="submit" disabled={createOverride.isPending}>
              <Plus className="mr-2 h-4 w-4" />
              Create Override
            </Button>
          </form>

          <div className="mt-5 space-y-2">
            {overridesQuery.isLoading ? (
              <p className="text-sm text-[var(--text-secondary)]">Loading model overrides...</p>
            ) : null}
            {(overridesQuery.data?.overrides ?? []).map((ov) => (
              <div
                key={ov.id}
                className="flex items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2.5 text-xs"
              >
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-medium text-[var(--text-primary)]">{ov.name}</span>
                    <Badge variant="outline" className="font-mono">
                      {ov.model_pattern}
                    </Badge>
                    {ov.target_model ? (
                      <Badge variant="secondary" className="font-mono">
                        -&gt; {ov.target_model}
                      </Badge>
                    ) : null}
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() =>
                      updateOverride.mutate({ id: ov.id, body: { enabled: !ov.enabled } })
                    }
                  >
                    {ov.enabled ? "Disable" : "Enable"}
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => deleteOverride.mutate(ov.id)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </Card>

        <Card className="p-6">
          <CardHeader className="p-0 pb-4">
            <CardTitle className="text-base">New Bypass Rule</CardTitle>
          </CardHeader>
          <form onSubmit={handleCreateBypass} className="space-y-3">
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="mb-name">Name</Label>
                <Input
                  id="mb-name"
                  placeholder="Raw Header Bypass"
                  {...bypassForm.register("name")}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mb-pattern">Model Pattern</Label>
                <Input
                  id="mb-pattern"
                  placeholder="*"
                  className="font-mono"
                  {...bypassForm.register("model_pattern")}
                />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="mb-action">Action</Label>
                <select
                  id="mb-action"
                  className="h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 text-sm text-[var(--text-primary)]"
                  {...bypassForm.register("action")}
                >
                  {BYPASS_ACTIONS.map((a) => (
                    <option key={a} value={a}>
                      {a}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1">
                <Label htmlFor="mb-target">Target Model (rewrite only)</Label>
                <Input
                  id="mb-target"
                  placeholder="gpt-4o-mini"
                  className="font-mono"
                  {...bypassForm.register("target_model")}
                />
              </div>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="mb-hname">Header Name (optional)</Label>
                <Input
                  id="mb-hname"
                  placeholder="X-Bypass-Compression"
                  className="font-mono"
                  {...bypassForm.register("header_name")}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="mb-hval">Header Value (optional)</Label>
                <Input
                  id="mb-hval"
                  placeholder="true"
                  className="font-mono"
                  {...bypassForm.register("header_value")}
                />
              </div>
            </div>
            {bypassFirstErr ? (
              <p className="text-xs text-[var(--accent-primary)]">{bypassFirstErr}</p>
            ) : null}
            <Button type="submit" disabled={createBypass.isPending}>
              <Plus className="mr-2 h-4 w-4" />
              Create Bypass Rule
            </Button>
          </form>

          <div className="mt-5 space-y-2">
            {bypassQuery.isLoading ? (
              <p className="text-sm text-[var(--text-secondary)]">Loading bypass rules...</p>
            ) : null}
            {(bypassQuery.data?.rules ?? []).map((br) => (
              <div
                key={br.id}
                className="flex items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2.5 text-xs"
              >
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-medium text-[var(--text-primary)]">{br.name}</span>
                    <Badge variant="secondary" className="font-mono">
                      {br.action}
                    </Badge>
                    <Badge variant="outline" className="font-mono">
                      {br.model_pattern}
                    </Badge>
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() =>
                      updateBypass.mutate({ id: br.id, body: { enabled: !br.enabled } })
                    }
                  >
                    {br.enabled ? "Disable" : "Enable"}
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => deleteBypass.mutate(br.id)}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card className="p-6">
        <CardHeader className="p-0 pb-4">
          <CardTitle className="text-base">Dry-Run Preview</CardTitle>
        </CardHeader>
        <div className="grid gap-4 lg:grid-cols-2">
          <div className="space-y-3">
            <div className="space-y-1">
              <Label htmlFor="me-prev-model">Model</Label>
              <Input
                id="me-prev-model"
                value={previewModel}
                onChange={(e) => setPreviewModel(e.target.value)}
                className="font-mono"
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="me-prev-body">JSON Request Body</Label>
              <textarea
                id="me-prev-body"
                rows={5}
                value={previewBody}
                onChange={(e) => setPreviewBody(e.target.value)}
                className="w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] p-3 font-mono text-xs text-[var(--text-primary)]"
              />
            </div>
            <Button type="button" variant="outline" onClick={handleRunPreview}>
              <Play className="mr-2 h-4 w-4" />
              Run Preview
            </Button>
          </div>
          <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3">
            <p className="mb-2 font-mono text-xs text-[var(--text-secondary)]">
              Mutated Output {previewMutation.data ? `(model: ${previewMutation.data.model})` : ""}
            </p>
            <pre className="max-h-56 overflow-auto font-mono text-xs text-[var(--text-primary)]">
              {previewMutation.data
                ? JSON.stringify(previewMutation.data, null, 2)
                : "Run preview to inspect transformed request payload and metadata."}
            </pre>
          </div>
        </div>
      </Card>
    </div>
  );
}
