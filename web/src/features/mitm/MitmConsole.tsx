import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  Download,
  Play,
  Plus,
  RefreshCw,
  Shield,
  Square,
  Trash2,
} from "lucide-react";
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
  useClearMitmEventsMutation,
  useCreateMitmAliasMutation,
  useDeleteMitmAliasMutation,
  useGenerateMitmCAMutation,
  useMitmAliases,
  useMitmEvents,
  useMitmStatus,
  useStartMitmMutation,
  useStopMitmMutation,
  useUpdateMitmConfigMutation,
} from "./hooks";
import {
  type CreateMitmAliasFormInput,
  type CreateMitmAliasFormOutput,
  createMitmAliasSchema,
  MITM_ALIAS_TOOLS,
  MITM_TOOL_IDS,
  type MitmConfigFormInput,
  type MitmConfigFormOutput,
  mitmConfigSchema,
} from "./schemas";

export function MitmConsole() {
  const statusQuery = useMitmStatus();
  const aliasesQuery = useMitmAliases();
  const eventsQuery = useMitmEvents();

  const updateConfig = useUpdateMitmConfigMutation();
  const startMitm = useStartMitmMutation();
  const stopMitm = useStopMitmMutation();
  const generateCA = useGenerateMitmCAMutation();
  const createAlias = useCreateMitmAliasMutation();
  const deleteAlias = useDeleteMitmAliasMutation();
  const clearEvents = useClearMitmEventsMutation();

  const configForm = useForm<MitmConfigFormInput, unknown, MitmConfigFormOutput>({
    resolver: zodResolver(mitmConfigSchema),
    defaultValues: {
      listen_addr: "127.0.0.1:8443",
      target_gateway_url: "http://127.0.0.1:20128",
      enabled_tools: [...MITM_TOOL_IDS],
      auto_start: false,
    },
  });

  const aliasForm = useForm<CreateMitmAliasFormInput, unknown, CreateMitmAliasFormOutput>({
    resolver: zodResolver(createMitmAliasSchema),
    defaultValues: {
      tool: "copilot",
      source_model: "",
      target_model: "",
      enabled: true,
    },
  });

  const status = statusQuery.data;

  useEffect(() => {
    if (!status || configForm.formState.isDirty) {
      return;
    }
    const validTools = status.enabled_tools.filter((t): t is (typeof MITM_TOOL_IDS)[number] =>
      (MITM_TOOL_IDS as readonly string[]).includes(t),
    );
    configForm.reset({
      listen_addr: status.listen_addr,
      target_gateway_url: status.target_gateway_url,
      enabled_tools: validTools,
      auto_start: status.auto_start,
    });
  }, [status, configForm, configForm.formState.isDirty]);

  const resetMutationErrors = () => {
    updateConfig.reset();
    startMitm.reset();
    stopMitm.reset();
    generateCA.reset();
    createAlias.reset();
    deleteAlias.reset();
    clearEvents.reset();
  };

  const firstError =
    statusQuery.error ??
    aliasesQuery.error ??
    eventsQuery.error ??
    updateConfig.error ??
    startMitm.error ??
    stopMitm.error ??
    generateCA.error ??
    createAlias.error ??
    deleteAlias.error ??
    clearEvents.error;

  const isRunning = status?.status === "running";
  const isErrorState = status?.status === "error";
  const selectedTools = configForm.watch("enabled_tools") ?? [];

  const toggleTool = (toolId: (typeof MITM_TOOL_IDS)[number]) => {
    const next = selectedTools.includes(toolId)
      ? selectedTools.filter((item) => item !== toolId)
      : [...selectedTools, toolId];
    configForm.setValue("enabled_tools", next, { shouldDirty: true, shouldValidate: true });
  };

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Local TLS Interception
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            MITM Proxy &amp; Root CA Studio
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Loopback-only Go TLS MITM listener, ECDSA P-256 Root CA, IDE host routing, and model
            alias rewriting.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <Badge variant={isRunning ? "success" : isErrorState ? "destructive" : "secondary"}>
            {status ? status.status.toUpperCase() : "LOADING"}
          </Badge>
          {isRunning ? (
            <Button
              type="button"
              variant="destructive"
              size="sm"
              onClick={() => {
                resetMutationErrors();
                stopMitm.mutate();
              }}
              disabled={stopMitm.isPending}
            >
              <Square className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{stopMitm.isPending ? "Stopping..." : "Stop Proxy"}</span>
            </Button>
          ) : (
            <Button
              type="button"
              size="sm"
              onClick={() => {
                resetMutationErrors();
                startMitm.mutate();
              }}
              disabled={startMitm.isPending}
            >
              <Play className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{startMitm.isPending ? "Starting..." : "Start Proxy"}</span>
            </Button>
          )}
        </div>
      </section>

      {firstError ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>MITM operation failed</AlertTitle>
            <AlertDescription className="font-mono">{firstError.message}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      {status?.last_error ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>Last MITM runtime error</AlertTitle>
            <AlertDescription className="font-mono">{status.last_error}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Shield className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <CardTitle>Proxy Configuration</CardTitle>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              {status?.bound_addr ? `Bound: ${status.bound_addr}` : "Loopback only"}
            </span>
          </CardHeader>

          <form
            className="space-y-4 p-5"
            onSubmit={configForm.handleSubmit((values) => {
              resetMutationErrors();
              updateConfig.mutate(values, {
                onSuccess: (updated) => {
                  const validTools = updated.enabled_tools.filter(
                    (t): t is (typeof MITM_TOOL_IDS)[number] =>
                      (MITM_TOOL_IDS as readonly string[]).includes(t),
                  );
                  configForm.reset({
                    listen_addr: updated.listen_addr,
                    target_gateway_url: updated.target_gateway_url,
                    enabled_tools: validTools,
                    auto_start: updated.auto_start,
                  });
                },
              });
            })}
          >
            <div className="space-y-1.5">
              <Label htmlFor="mitm-listen-addr">Listen Address (Loopback)</Label>
              <Input
                id="mitm-listen-addr"
                className="font-mono"
                {...configForm.register("listen_addr")}
              />
              {configForm.formState.errors.listen_addr ? (
                <p className="text-xs text-[var(--status-error)]">
                  {configForm.formState.errors.listen_addr.message}
                </p>
              ) : null}
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="mitm-target-url">Target Gateway URL (Loopback)</Label>
              <Input
                id="mitm-target-url"
                className="font-mono"
                {...configForm.register("target_gateway_url")}
              />
              {configForm.formState.errors.target_gateway_url ? (
                <p className="text-xs text-[var(--status-error)]">
                  {configForm.formState.errors.target_gateway_url.message}
                </p>
              ) : null}
            </div>

            <div className="space-y-2">
              <Label>Enabled IDE Interceptors</Label>
              <div className="grid grid-cols-2 gap-2">
                {MITM_TOOL_IDS.map((toolId) => {
                  const checked = selectedTools.includes(toolId);
                  return (
                    <label
                      key={toolId}
                      className="flex cursor-pointer items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2 text-sm text-[var(--text-primary)]"
                    >
                      <input
                        type="checkbox"
                        checked={checked}
                        onChange={() => toggleTool(toolId)}
                        className="h-4 w-4 accent-[var(--accent-primary)]"
                      />
                      <span className="font-mono text-xs uppercase">{toolId}</span>
                    </label>
                  );
                })}
              </div>
            </div>

            <label className="flex cursor-pointer items-center gap-2 text-sm text-[var(--text-secondary)]">
              <input
                type="checkbox"
                {...configForm.register("auto_start")}
                className="h-4 w-4 accent-[var(--accent-primary)]"
              />
              <span>Auto-start MITM proxy on router boot</span>
            </label>

            <Button type="submit" size="sm" disabled={updateConfig.isPending}>
              {updateConfig.isPending ? "Saving..." : "Save Configuration"}
            </Button>
          </form>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Root CA &amp; OS Trust Setup</CardTitle>
            <Badge variant={status?.ca.exists ? "success" : "secondary"}>
              {status?.ca.exists ? "GENERATED" : "MISSING"}
            </Badge>
          </CardHeader>

          <div className="space-y-4 p-5 text-sm">
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                size="sm"
                variant="outline"
                onClick={() => generateCA.mutate({ force: false })}
                disabled={generateCA.isPending}
              >
                <RefreshCw
                  className="h-3.5 w-3.5 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
                <span>{status?.ca.exists ? "Verify CA" : "Generate Root CA"}</span>
              </Button>
              <Button
                type="button"
                size="sm"
                variant="secondary"
                onClick={() => generateCA.mutate({ force: true })}
                disabled={generateCA.isPending}
              >
                <span>Rotate Root CA</span>
              </Button>
              {status?.ca.exists ? (
                <Button asChild size="sm" variant="outline">
                  <a href="/api/mitm/ca/download" download="omeg4mon-mitm-rootCA.crt">
                    <Download className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>Download PEM</span>
                  </a>
                </Button>
              ) : null}
            </div>

            {status?.ca.exists ? (
              <dl className="space-y-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3 font-mono text-xs">
                <div className="flex justify-between gap-3">
                  <dt className="text-[var(--text-muted)]">Fingerprint</dt>
                  <dd className="truncate text-[var(--text-primary)]">{status.ca.fingerprint}</dd>
                </div>
                <div className="flex justify-between gap-3">
                  <dt className="text-[var(--text-muted)]">Valid Until</dt>
                  <dd className="text-[var(--text-primary)]">{status.ca.not_after}</dd>
                </div>
              </dl>
            ) : null}

            {status?.install_guide ? (
              <div className="space-y-2">
                <p className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  Trust Command ({status.install_guide.os})
                </p>
                <pre className="overflow-x-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)]">
                  {status.install_guide.ca_command}
                </pre>
                <p className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  Hosts Entries ({status.install_guide.hosts_path})
                </p>
                <pre className="overflow-x-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 font-mono text-xs text-[var(--text-primary)]">
                  {status.install_guide.hosts_snippet}
                </pre>
              </div>
            ) : null}
          </div>
        </Card>
      </div>

      <EmberRule />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Model Rewrite Aliases</CardTitle>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              {aliasesQuery.data?.aliases.length ?? 0} active
            </span>
          </CardHeader>

          <div className="space-y-4 p-5">
            <form
              className="space-y-2"
              onSubmit={aliasForm.handleSubmit((values) => {
                resetMutationErrors();
                createAlias.mutate(values, {
                  onSuccess: () =>
                    aliasForm.reset({ ...values, source_model: "", target_model: "" }),
                });
              })}
            >
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-4">
                <select
                  aria-label="IDE tool"
                  {...aliasForm.register("tool")}
                  className="h-9 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-2.5 font-mono text-xs text-[var(--text-primary)]"
                >
                  {MITM_ALIAS_TOOLS.map((t) => (
                    <option key={t} value={t}>
                      {t}
                    </option>
                  ))}
                </select>
                <Input
                  placeholder="source model"
                  aria-label="Source model"
                  className="font-mono text-xs"
                  {...aliasForm.register("source_model")}
                />
                <Input
                  placeholder="target model"
                  aria-label="Target model"
                  className="font-mono text-xs"
                  {...aliasForm.register("target_model")}
                />
                <Button type="submit" size="sm" disabled={createAlias.isPending}>
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Add</span>
                </Button>
              </div>
              {aliasForm.formState.errors.source_model ||
              aliasForm.formState.errors.target_model ? (
                <p className="text-xs text-[var(--status-error)]">
                  {aliasForm.formState.errors.source_model?.message ??
                    aliasForm.formState.errors.target_model?.message}
                </p>
              ) : null}
            </form>

            <div className="divide-y divide-[var(--border-subtle)] rounded-md border border-[var(--border-subtle)]">
              {aliasesQuery.data && aliasesQuery.data.aliases.length > 0 ? (
                aliasesQuery.data.aliases.map((alias) => (
                  <div
                    key={alias.id}
                    className="flex items-center justify-between gap-3 px-3.5 py-2.5 text-xs"
                  >
                    <div className="min-w-0 font-mono">
                      <Badge variant="secondary" className="mr-2">
                        {alias.tool}
                      </Badge>
                      <span className="text-[var(--text-secondary)]">{alias.source_model}</span>
                      <span className="mx-1.5 text-[var(--accent-text)]">-&gt;</span>
                      <span className="text-[var(--text-primary)]">{alias.target_model}</span>
                    </div>
                    <Button
                      type="button"
                      variant="secondary"
                      size="icon"
                      aria-label={`Delete alias ${alias.source_model}`}
                      onClick={() => deleteAlias.mutate(alias.id)}
                      disabled={deleteAlias.isPending}
                      className="h-7 w-7 shrink-0"
                    >
                      <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                    </Button>
                  </div>
                ))
              ) : (
                <p className="px-3.5 py-4 text-xs text-[var(--text-muted)]">
                  No MITM model aliases configured yet.
                </p>
              )}
            </div>
          </div>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Intercepted Traffic Log</CardTitle>
            <Button
              type="button"
              variant="secondary"
              size="xs"
              onClick={() => clearEvents.mutate()}
              disabled={clearEvents.isPending}
            >
              Clear
            </Button>
          </CardHeader>

          <div className="divide-y divide-[var(--border-subtle)]">
            {eventsQuery.data && eventsQuery.data.events.length > 0 ? (
              eventsQuery.data.events.map((ev) => (
                <div
                  key={ev.id}
                  className="flex items-center justify-between gap-3 px-5 py-3 text-xs"
                >
                  <div className="min-w-0 space-y-0.5">
                    <p className="truncate font-mono text-[var(--text-primary)]">
                      {ev.method} {ev.host}
                      {ev.path}
                    </p>
                    <p className="font-mono text-[var(--text-muted)]">
                      {ev.tool}
                      {ev.source_model ? ` | ${ev.source_model} -> ${ev.target_model}` : ""}
                    </p>
                  </div>
                  <div className="flex shrink-0 items-center gap-2 font-mono">
                    <Badge variant={ev.status_code < 400 ? "success" : "destructive"}>
                      {ev.status_code}
                    </Badge>
                    <span className="text-[var(--text-muted)]">{ev.latency_ms}ms</span>
                  </div>
                </div>
              ))
            ) : (
              <p className="px-5 py-6 text-sm text-[var(--text-muted)]">
                No intercepted requests recorded yet.
              </p>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
}
