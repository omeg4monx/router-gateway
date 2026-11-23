"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { AlertTriangle, CheckCircle2, Play, RotateCcw, Square, Workflow } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearPxpipeLogsMutation,
  usePreviewPxpipeMutation,
  usePxpipeHealth,
  usePxpipeLogs,
  usePxpipeStats,
  usePxpipeStatus,
  useResetPxpipeStatsMutation,
  useRestartPxpipeMutation,
  useStartPxpipeMutation,
  useStopPxpipeMutation,
  useUpdatePxpipeConfigMutation,
} from "./hooks";
import {
  type PxpipeConfigFormInput,
  type PxpipeConfigFormOutput,
  type PxpipePreviewFormInput,
  type PxpipePreviewFormOutput,
  pxpipeConfigSchema,
  pxpipePreviewSchema,
} from "./schemas";

const DEFAULT_PREVIEW_PAYLOAD = JSON.stringify(
  {
    model: "gpt-4o",
    user: null,
    messages: [
      { role: "system", content: "You   are   concise." },
      { role: "system", content: "You are concise." },
      { role: "user", content: "Explain    PXPIPE    transforms." },
    ],
  },
  null,
  2,
);

export function PxpipeConsole() {
  const statusQuery = usePxpipeStatus();
  const healthQuery = usePxpipeHealth();
  const statsQuery = usePxpipeStats();
  const logsQuery = usePxpipeLogs();

  const updateConfig = useUpdatePxpipeConfigMutation();
  const startPipe = useStartPxpipeMutation();
  const stopPipe = useStopPxpipeMutation();
  const restartPipe = useRestartPxpipeMutation();
  const resetStats = useResetPxpipeStatsMutation();
  const clearLogs = useClearPxpipeLogsMutation();
  const previewTransform = usePreviewPxpipeMutation();

  const configForm = useForm<PxpipeConfigFormInput, unknown, PxpipeConfigFormOutput>({
    resolver: zodResolver(pxpipeConfigSchema),
    defaultValues: {
      enabled: true,
      auto_start: true,
      strip_null_fields: true,
      collapse_whitespace: true,
      dedupe_system_prompts: true,
      max_history_messages: 0,
    },
  });

  const previewForm = useForm<PxpipePreviewFormInput, unknown, PxpipePreviewFormOutput>({
    resolver: zodResolver(pxpipePreviewSchema),
    defaultValues: {
      payloadJson: DEFAULT_PREVIEW_PAYLOAD,
    },
  });

  const status = statusQuery.data;

  useEffect(() => {
    if (!status) {
      return;
    }
    configForm.reset({
      enabled: status.config.enabled,
      auto_start: status.config.auto_start,
      strip_null_fields: status.config.strip_null_fields,
      collapse_whitespace: status.config.collapse_whitespace,
      dedupe_system_prompts: status.config.dedupe_system_prompts,
      max_history_messages: status.config.max_history_messages,
    });
  }, [status, configForm]);

  const firstError =
    statusQuery.error ??
    healthQuery.error ??
    statsQuery.error ??
    logsQuery.error ??
    updateConfig.error ??
    startPipe.error ??
    stopPipe.error ??
    restartPipe.error ??
    resetStats.error ??
    clearLogs.error ??
    previewTransform.error;

  const isRunning = Boolean(status?.running);
  const stats = statsQuery.data;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            In-Process Transform Pipeline
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            PXPIPE Pipeline Studio
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Zero-hop Go request transform interceptor with live self-test health checks, dry-run
            preview, and savings telemetry.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={isRunning ? "success" : "secondary"}>
            {isRunning ? "RUNNING" : "STOPPED"}
          </Badge>
          {isRunning ? (
            <Button
              type="button"
              variant="destructive"
              size="sm"
              onClick={() => stopPipe.mutate()}
              disabled={stopPipe.isPending}
            >
              <Square className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Stop</span>
            </Button>
          ) : (
            <Button
              type="button"
              size="sm"
              onClick={() => startPipe.mutate()}
              disabled={startPipe.isPending}
            >
              <Play className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Start</span>
            </Button>
          )}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => restartPipe.mutate()}
            disabled={restartPipe.isPending}
          >
            <RotateCcw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>Restart</span>
          </Button>
        </div>
      </section>

      {firstError ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>PXPIPE operation failed</AlertTitle>
            <AlertDescription className="font-mono">{firstError.message}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Card className="p-4">
          <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Requests Seen
          </span>
          <p className="mt-2 font-mono text-2xl font-semibold tabular-nums text-[var(--text-primary)]">
            {stats?.requests_seen ?? 0}
          </p>
        </Card>
        <Card className="p-4">
          <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Transformed
          </span>
          <p className="mt-2 font-mono text-2xl font-semibold tabular-nums text-[var(--text-primary)]">
            {stats?.requests_transformed ?? 0}
          </p>
        </Card>
        <Card className="p-4">
          <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Bytes Saved
          </span>
          <p className="mt-2 font-mono text-2xl font-semibold tabular-nums text-[var(--accent-text)]">
            {stats?.bytes_saved ?? 0}
          </p>
        </Card>
        <Card className="p-4">
          <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Est. Tokens Saved
          </span>
          <p className="mt-2 font-mono text-2xl font-semibold tabular-nums text-[var(--status-ok)]">
            {stats?.estimated_tokens_saved ?? 0}
          </p>
        </Card>
      </section>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Workflow className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Pipeline Stages</CardTitle>
            </div>
            <Button
              type="button"
              variant="secondary"
              size="xs"
              onClick={() => resetStats.mutate()}
              disabled={resetStats.isPending}
            >
              Reset Stats
            </Button>
          </CardHeader>

          <form
            className="space-y-3.5 p-5 text-sm"
            onSubmit={configForm.handleSubmit((values) => updateConfig.mutate(values))}
          >
            <label className="flex cursor-pointer items-center gap-2 text-[var(--text-primary)]">
              <input
                type="checkbox"
                {...configForm.register("enabled")}
                className="h-4 w-4 accent-[#E56A4A]"
              />
              <span>Enable gateway request transforms</span>
            </label>
            <label className="flex cursor-pointer items-center gap-2 text-[var(--text-primary)]">
              <input
                type="checkbox"
                {...configForm.register("auto_start")}
                className="h-4 w-4 accent-[#E56A4A]"
              />
              <span>Auto-start pipeline on boot</span>
            </label>
            <label className="flex cursor-pointer items-center gap-2 text-[var(--text-primary)]">
              <input
                type="checkbox"
                {...configForm.register("strip_null_fields")}
                className="h-4 w-4 accent-[#E56A4A]"
              />
              <span>Strip null JSON fields</span>
            </label>
            <label className="flex cursor-pointer items-center gap-2 text-[var(--text-primary)]">
              <input
                type="checkbox"
                {...configForm.register("collapse_whitespace")}
                className="h-4 w-4 accent-[#E56A4A]"
              />
              <span>Collapse repeated message whitespace</span>
            </label>
            <label className="flex cursor-pointer items-center gap-2 text-[var(--text-primary)]">
              <input
                type="checkbox"
                {...configForm.register("dedupe_system_prompts")}
                className="h-4 w-4 accent-[#E56A4A]"
              />
              <span>Deduplicate identical system prompts</span>
            </label>

            <div className="space-y-1.5 pt-1">
              <Label htmlFor="pxpipe-max-history">Max History Messages (0 = unlimited)</Label>
              <Input
                id="pxpipe-max-history"
                type="number"
                className="font-mono"
                {...configForm.register("max_history_messages")}
              />
              {configForm.formState.errors.max_history_messages ? (
                <p className="text-xs text-[var(--status-error)]">
                  {configForm.formState.errors.max_history_messages.message}
                </p>
              ) : null}
            </div>

            <Button type="submit" size="sm" disabled={updateConfig.isPending}>
              {updateConfig.isPending ? "Saving..." : "Save Pipeline Config"}
            </Button>
          </form>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Live Self-Test Health</CardTitle>
            <Badge variant={healthQuery.data?.healthy ? "success" : "destructive"}>
              {healthQuery.data ? (healthQuery.data.healthy ? "HEALTHY" : "DEGRADED") : "CHECKING"}
            </Badge>
          </CardHeader>

          <div className="divide-y divide-[var(--border-subtle)]">
            {healthQuery.data?.checks.map((check) => (
              <div key={check.id} className="flex items-center justify-between gap-3 px-5 py-3.5">
                <div className="flex items-center gap-2.5">
                  {check.ok ? (
                    <CheckCircle2 className="h-4 w-4 text-[var(--status-ok)]" aria-hidden="true" />
                  ) : (
                    <AlertTriangle
                      className="h-4 w-4 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                  )}
                  <div>
                    <p className="text-sm font-medium text-[var(--text-primary)]">{check.label}</p>
                    <p className="font-mono text-xs text-[var(--text-muted)]">{check.detail}</p>
                  </div>
                </div>
                <Badge variant={check.ok ? "success" : "destructive"}>
                  {check.ok ? "PASS" : "FAIL"}
                </Badge>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <EmberRule />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Dry-Run Transform Preview</CardTitle>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              POST /api/pxpipe/preview
            </span>
          </CardHeader>

          <form
            className="space-y-3 p-5"
            onSubmit={previewForm.handleSubmit((values) => {
              const parsed = JSON.parse(values.payloadJson) as Record<string, unknown>;
              previewTransform.mutate({ payload: parsed });
            })}
          >
            <textarea
              rows={8}
              aria-label="Preview JSON payload"
              {...previewForm.register("payloadJson")}
              className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)]"
            />
            {previewForm.formState.errors.payloadJson ? (
              <p className="text-xs text-[var(--status-error)]">
                {previewForm.formState.errors.payloadJson.message}
              </p>
            ) : null}
            <Button type="submit" size="sm" disabled={previewTransform.isPending}>
              {previewTransform.isPending ? "Running..." : "Run Dry-Run Preview"}
            </Button>

            {previewTransform.data ? (
              <div className="space-y-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3 text-xs">
                <p className="font-mono text-[var(--text-primary)]">
                  {previewTransform.data.original_bytes}B -&gt;{" "}
                  {previewTransform.data.transformed_bytes}B (saved{" "}
                  {previewTransform.data.bytes_saved}B / ~
                  {previewTransform.data.estimated_tokens_saved} tok)
                </p>
                <p className="font-mono text-[var(--text-muted)]">
                  Stages: {previewTransform.data.stages_applied.join(", ") || "none"}
                </p>
              </div>
            ) : null}
          </form>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Transform Logs</CardTitle>
            <Button
              type="button"
              variant="secondary"
              size="xs"
              onClick={() => clearLogs.mutate()}
              disabled={clearLogs.isPending}
            >
              Clear Logs
            </Button>
          </CardHeader>

          <div className="divide-y divide-[var(--border-subtle)]">
            {logsQuery.data && logsQuery.data.logs.length > 0 ? (
              logsQuery.data.logs.map((entry) => (
                <div
                  key={`${entry.timestamp}-${entry.request_id}`}
                  className="flex items-center justify-between gap-3 px-5 py-3 text-xs"
                >
                  <div className="min-w-0 space-y-0.5 font-mono">
                    <p className="truncate text-[var(--text-primary)]">
                      {entry.model || "unknown"} ({entry.request_id || "n/a"})
                    </p>
                    <p className="truncate text-[var(--text-muted)]">
                      {entry.stages_applied.join(", ")}
                    </p>
                  </div>
                  <span className="shrink-0 font-mono text-[var(--accent-text)]">
                    -{entry.bytes_saved}B
                  </span>
                </div>
              ))
            ) : (
              <p className="px-5 py-6 text-sm text-[var(--text-muted)]">
                No gateway requests transformed yet.
              </p>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
}
