"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  Check,
  CheckCircle2,
  Copy,
  Download,
  FileCode2,
  FolderCog,
  HardDriveUpload,
  RefreshCw,
  RotateCcw,
  ShieldAlert,
  Terminal,
} from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useApplyCliToolConfigMutation,
  useCliToolsCatalog,
  useGenerateCliToolConfigMutation,
  useRollbackCliToolConfigMutation,
} from "./hooks";
import {
  type ApplyCLIToolConfigResponse,
  type CLIToolDescriptor,
  type CLIToolID,
  DEFAULT_GENERATE_CLI_CONFIG_VALUES,
  type GenerateCLIToolConfigResponse,
  type GenerateCliConfigFormInput,
  type GenerateCliConfigFormOutput,
  generateCliConfigSchema,
  type RollbackCLIToolConfigResponse,
  toCLIToolConfigRequest,
} from "./schemas";

function formatBadgeLabel(format: CLIToolDescriptor["format"]): string {
  return format.toUpperCase();
}

export function CliToolsConsole() {
  const catalogQuery = useCliToolsCatalog();
  const generateMutation = useGenerateCliToolConfigMutation();
  const applyMutation = useApplyCliToolConfigMutation();
  const rollbackMutation = useRollbackCliToolConfigMutation();

  const [selectedTool, setSelectedTool] = useState<CLIToolID>("claude-code");
  const [preview, setPreview] = useState<GenerateCLIToolConfigResponse | null>(null);
  const [copiedTool, setCopiedTool] = useState<CLIToolID | null>(null);
  const [copyError, setCopyError] = useState<string | null>(null);
  const [confirmApply, setConfirmApply] = useState(false);
  const [confirmRollback, setConfirmRollback] = useState(false);
  const [applyOutcome, setApplyOutcome] = useState<ApplyCLIToolConfigResponse | null>(null);
  const [rollbackOutcome, setRollbackOutcome] = useState<RollbackCLIToolConfigResponse | null>(
    null,
  );
  const [localActionError, setLocalActionError] = useState<string | null>(null);

  const initialPreviewFiredRef = useRef(false);
  const previewSeqRef = useRef(0);

  const {
    register,
    handleSubmit,
    getValues,
    setValue,
    formState: { errors },
  } = useForm<GenerateCliConfigFormInput, unknown, GenerateCliConfigFormOutput>({
    resolver: zodResolver(generateCliConfigSchema),
    defaultValues: DEFAULT_GENERATE_CLI_CONFIG_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const tools = catalogQuery.data?.tools ?? [];
  const localApplyEnabled = catalogQuery.data?.local_apply_enabled ?? false;

  const activeDescriptor =
    tools.find((item) => item.tool === selectedTool) ?? (tools.length > 0 ? tools[0] : null);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    if (!getValues("base_url")) {
      setValue("base_url", window.location.origin, { shouldValidate: false });
    }
  }, [getValues, setValue]);

  useEffect(() => {
    const firstTool = tools[0];
    if (!firstTool) {
      return;
    }
    const currentExists = tools.some((item) => item.tool === selectedTool);
    if (!currentExists) {
      setSelectedTool(firstTool.tool);
    }
  }, [tools, selectedTool]);

  const runGenerate = (tool: CLIToolID, values: GenerateCliConfigFormOutput) => {
    const seq = ++previewSeqRef.current;
    generateMutation.mutate(
      {
        tool,
        body: toCLIToolConfigRequest(values),
      },
      {
        onSuccess: (generated) => {
          if (seq !== previewSeqRef.current) {
            return;
          }
          setPreview(generated);
        },
      },
    );
  };

  useEffect(() => {
    if (initialPreviewFiredRef.current || !activeDescriptor) {
      return;
    }
    const parsed = generateCliConfigSchema.safeParse(getValues());
    if (!parsed.success) {
      return;
    }
    initialPreviewFiredRef.current = true;
    runGenerate(activeDescriptor.tool, parsed.data);
  });

  const isMutating = applyMutation.isPending || rollbackMutation.isPending;

  const handleSelectTool = (tool: CLIToolID) => {
    setSelectedTool(tool);
    setConfirmApply(false);
    setConfirmRollback(false);
    setApplyOutcome(null);
    setRollbackOutcome(null);
    setLocalActionError(null);
    setCopyError(null);

    const parsed = generateCliConfigSchema.safeParse(getValues());
    if (parsed.success) {
      runGenerate(tool, parsed.data);
    } else {
      void handleSubmit((valid) => {
        runGenerate(tool, valid);
      })();
    }
  };

  const onGenerateSubmit = (values: GenerateCliConfigFormOutput) => {
    setCopyError(null);
    runGenerate(selectedTool, values);
  };

  const handleCopyConfig = async () => {
    if (!preview || preview.tool !== selectedTool) {
      return;
    }
    setCopyError(null);
    try {
      await navigator.clipboard.writeText(preview.content);
      setCopiedTool(preview.tool);
      window.setTimeout(() => {
        setCopiedTool((current) => (current === preview.tool ? null : current));
      }, 2200);
    } catch {
      setCopyError("Clipboard access was denied. Copy the generated configuration manually.");
    }
  };

  const handleDownloadConfig = () => {
    if (!preview || preview.tool !== selectedTool) {
      return;
    }
    const mimeType = preview.format === "json" ? "application/json" : "application/toml";
    const blob = new Blob([preview.content], { type: mimeType });
    const objectUrl = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = objectUrl;
    anchor.download = preview.filename;
    anchor.style.display = "none";
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
    window.setTimeout(() => URL.revokeObjectURL(objectUrl), 1000);
  };

  const handleConfirmApply = handleSubmit((values) => {
    if (isMutating || !localApplyEnabled) {
      return;
    }
    setLocalActionError(null);
    setApplyOutcome(null);
    setRollbackOutcome(null);

    applyMutation.mutate(
      {
        tool: selectedTool,
        body: toCLIToolConfigRequest(values),
      },
      {
        onSuccess: (result) => {
          setApplyOutcome(result);
          setConfirmApply(false);
          runGenerate(selectedTool, values);
        },
        onError: (err) => {
          setLocalActionError(err.message);
          setConfirmApply(false);
        },
      },
    );
  });

  const handleConfirmRollback = () => {
    if (isMutating || !localApplyEnabled) {
      return;
    }
    setLocalActionError(null);
    setApplyOutcome(null);
    setRollbackOutcome(null);

    rollbackMutation.mutate(selectedTool, {
      onSuccess: (result) => {
        setRollbackOutcome(result);
        setConfirmRollback(false);
      },
      onError: (err) => {
        setLocalActionError(err.message);
        setConfirmRollback(false);
      },
    });
  };

  const previewStale = preview !== null && preview.tool !== selectedTool;
  const activePreview = preview && !previewStale ? preview : null;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Developer Integrations
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            CLI Tools Studio
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Generate, copy, download, and loopback-apply configuration profiles for supported AI
            coding CLIs.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          {catalogQuery.data ? (
            <Badge variant={localApplyEnabled ? "success" : "warning"}>
              {localApplyEnabled ? "Loopback Apply Enabled" : "Remote-Bind Disabled"}
            </Badge>
          ) : null}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void catalogQuery.refetch()}
            disabled={catalogQuery.isFetching}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{catalogQuery.isFetching ? "Syncing" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {catalogQuery.isError ? (
        <Alert variant="destructive">
          <AlertTriangle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="min-w-0 flex-1">
            <AlertTitle>Unable to load CLI tools catalog</AlertTitle>
            <AlertDescription className="font-mono break-all">
              {catalogQuery.error.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      <section aria-label="Supported CLI coding tools" className="space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Terminal className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            <h3 className="text-sm font-semibold uppercase tracking-wider text-[var(--text-secondary)]">
              1. Select Coding Client
            </h3>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            {tools.length > 0 ? `${tools.length} integrations` : "GET /api/clitools"}
          </span>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {catalogQuery.isPending && tools.length === 0 ? (
            <Card className="p-4 sm:col-span-2 lg:col-span-3">
              <p className="text-sm text-[var(--text-muted)]">
                Loading supported CLI tool catalog...
              </p>
            </Card>
          ) : (
            tools.map((tool) => {
              const isSelected = tool.tool === selectedTool;
              return (
                <button
                  key={tool.tool}
                  type="button"
                  onClick={() => handleSelectTool(tool.tool)}
                  aria-pressed={isSelected}
                  className={`flex h-full flex-col justify-between rounded-lg border p-4 text-left transition-colors ${
                    isSelected
                      ? "border-[#E56A4A] bg-[var(--bg-surface-raised)]"
                      : "border-[var(--border-subtle)] bg-[var(--bg-surface)] hover:border-[var(--border-strong)] hover:bg-[var(--bg-surface-raised)]"
                  }`}
                >
                  <div>
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">
                          {tool.name}
                        </p>
                        <p className="font-mono text-xs text-[var(--text-muted)]">{tool.tool}</p>
                      </div>
                      <Badge variant={isSelected ? "default" : "secondary"}>
                        {formatBadgeLabel(tool.format)}
                      </Badge>
                    </div>
                    <p className="mt-2.5 text-xs leading-relaxed text-[var(--text-secondary)]">
                      {tool.description}
                    </p>
                  </div>

                  <div className="mt-3 border-t border-[var(--border-subtle)] pt-2.5">
                    <p className="truncate font-mono text-[11px] text-[var(--text-muted)]">
                      {tool.target_path}
                    </p>
                  </div>
                </button>
              );
            })
          )}
        </div>
      </section>

      <EmberRule />

      <section className="grid grid-cols-1 gap-6 lg:grid-cols-12">
        <div className="space-y-6 lg:col-span-5">
          <Card>
            <CardHeader>
              <div>
                <div className="flex items-center gap-2">
                  <FolderCog className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                  <CardTitle>2. Profile Parameters</CardTitle>
                </div>
                <CardDescription className="mt-1">
                  Configure gateway URL, optional key, and default model for{" "}
                  <span className="font-mono text-[var(--text-primary)]">
                    {activeDescriptor?.name ?? selectedTool}
                  </span>
                  .
                </CardDescription>
              </div>
            </CardHeader>

            <CardContent className="pt-5">
              <form onSubmit={handleSubmit(onGenerateSubmit)} noValidate className="space-y-4">
                <div className="space-y-1.5">
                  <Label htmlFor="cli-base-url">Gateway Base URL</Label>
                  <Input
                    id="cli-base-url"
                    type="url"
                    placeholder="http://127.0.0.1:20128"
                    className="font-mono text-xs"
                    aria-invalid={errors.base_url ? "true" : undefined}
                    {...register("base_url")}
                  />
                  {errors.base_url ? (
                    <p role="alert" className="text-xs text-[var(--status-error)]">
                      {errors.base_url.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      HTTP or HTTPS root URL of this router (without trailing{" "}
                      <span className="font-mono">/v1</span>).
                    </p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="cli-api-key">Gateway API Key (Optional)</Label>
                  <Input
                    id="cli-api-key"
                    type="password"
                    autoComplete="off"
                    placeholder="Leave blank for YOUR_OMEG4MON_API_KEY placeholder"
                    className="font-mono text-xs"
                    aria-invalid={errors.api_key ? "true" : undefined}
                    {...register("api_key")}
                  />
                  {errors.api_key ? (
                    <p role="alert" className="text-xs text-[var(--status-error)]">
                      {errors.api_key.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      Key is rendered only into your generated snippet or local disk config.
                    </p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="cli-model">Default Model / Combo Alias (Optional)</Label>
                  <Input
                    id="cli-model"
                    type="text"
                    placeholder="e.g. anthropic/claude-sonnet-4 or gpt-4o"
                    className="font-mono text-xs"
                    aria-invalid={errors.model ? "true" : undefined}
                    {...register("model")}
                  />
                  {errors.model ? (
                    <p role="alert" className="text-xs text-[var(--status-error)]">
                      {errors.model.message}
                    </p>
                  ) : null}
                </div>

                <Button
                  type="submit"
                  className="w-full"
                  disabled={generateMutation.isPending || isMutating}
                >
                  <FileCode2 className="h-4 w-4" aria-hidden="true" />
                  <span>
                    {generateMutation.isPending ? "Generating Config..." : "Generate Preview"}
                  </span>
                </Button>
              </form>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <div>
                <div className="flex items-center gap-2">
                  <HardDriveUpload className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                  <CardTitle>Local Apply &amp; Rollback</CardTitle>
                </div>
                <CardDescription className="mt-1 font-mono text-xs">
                  {activeDescriptor?.target_path ?? "~/.config"}
                </CardDescription>
              </div>
              {catalogQuery.isPending ? (
                <Badge variant="secondary">Checking</Badge>
              ) : catalogQuery.data ? (
                <Badge variant={catalogQuery.data.local_apply_enabled ? "success" : "warning"}>
                  {catalogQuery.data.local_apply_enabled ? "Loopback" : "Remote-Bind Disabled"}
                </Badge>
              ) : null}
            </CardHeader>

            <CardContent className="space-y-4 pt-5">
              {catalogQuery.isPending ? (
                <p className="text-xs text-[var(--text-muted)]">
                  Checking local loopback apply availability...
                </p>
              ) : catalogQuery.isError || !catalogQuery.data ? (
                <Alert variant="destructive">
                  <AlertTriangle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertTitle>Local apply status unavailable</AlertTitle>
                    <AlertDescription>
                      Unable to verify loopback bind status while the CLI tools catalog is
                      unreachable.
                    </AlertDescription>
                  </div>
                </Alert>
              ) : !catalogQuery.data.local_apply_enabled ? (
                <Alert variant="warning">
                  <ShieldAlert
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertTitle>Remote-Bind Disabled</AlertTitle>
                    <AlertDescription>
                      Local filesystem writes and backup rollbacks are disabled because the router
                      is bound to a non-loopback address. Use Copy Config or Download Config File to
                      install manually.
                    </AlertDescription>
                  </div>
                </Alert>
              ) : (
                <>
                  <p className="text-xs leading-relaxed text-[var(--text-secondary)]">
                    Writes the rendered profile atomically to{" "}
                    <span className="font-mono text-[var(--text-primary)]">
                      {activeDescriptor?.target_path ?? "target path"}
                    </span>{" "}
                    and snapshots any existing file as a timestamped backup.
                  </p>

                  {localActionError ? (
                    <Alert variant="destructive">
                      <AlertTriangle
                        className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                        aria-hidden="true"
                      />
                      <div className="min-w-0 flex-1">
                        <AlertTitle>Local operation failed</AlertTitle>
                        <AlertDescription className="font-mono break-all">
                          {localActionError}
                        </AlertDescription>
                      </div>
                    </Alert>
                  ) : null}

                  {applyOutcome ? (
                    <Alert className="border-[var(--status-ok-border)] bg-[var(--status-ok-bg)]">
                      <CheckCircle2
                        className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                        aria-hidden="true"
                      />
                      <div className="min-w-0 flex-1">
                        <AlertTitle>Configuration applied locally</AlertTitle>
                        <AlertDescription className="space-y-1 font-mono text-xs">
                          <span className="block break-all">
                            Target: {applyOutcome.target_path}
                          </span>
                          <span className="block break-all">
                            Backup:{" "}
                            {applyOutcome.backup_path || "None (no pre-existing file to snapshot)"}
                          </span>
                        </AlertDescription>
                      </div>
                    </Alert>
                  ) : null}

                  {rollbackOutcome ? (
                    <Alert className="border-[var(--status-ok-border)] bg-[var(--status-ok-bg)]">
                      <CheckCircle2
                        className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                        aria-hidden="true"
                      />
                      <div className="min-w-0 flex-1">
                        <AlertTitle>Latest backup restored</AlertTitle>
                        <AlertDescription className="space-y-1 font-mono text-xs">
                          <span className="block break-all">
                            Target: {rollbackOutcome.target_path}
                          </span>
                          <span className="block break-all">
                            Restored from: {rollbackOutcome.backup_path}
                          </span>
                        </AlertDescription>
                      </div>
                    </Alert>
                  ) : null}

                  <div className="space-y-2.5">
                    {!confirmApply ? (
                      <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        disabled={isMutating}
                        onClick={() => {
                          setConfirmRollback(false);
                          setLocalActionError(null);
                          setConfirmApply(true);
                        }}
                      >
                        <HardDriveUpload className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                        <span>Apply Locally</span>
                      </Button>
                    ) : (
                      <div className="rounded-md border border-[#E56A4A] bg-[var(--bg-surface-raised)] p-3">
                        <p className="text-xs font-medium text-[var(--text-primary)]">
                          Overwrite local config for {activeDescriptor?.name ?? selectedTool}?
                        </p>
                        <p className="mt-1 truncate font-mono text-[11px] text-[var(--text-muted)]">
                          {activeDescriptor?.target_path}
                        </p>
                        <div className="mt-3 flex items-center gap-2">
                          <Button
                            type="button"
                            size="sm"
                            disabled={isMutating}
                            onClick={() => void handleConfirmApply()}
                          >
                            {applyMutation.isPending ? "Applying..." : "Confirm Apply"}
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            size="sm"
                            disabled={isMutating}
                            onClick={() => setConfirmApply(false)}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    )}

                    {!confirmRollback ? (
                      <Button
                        type="button"
                        variant="secondary"
                        className="w-full"
                        disabled={isMutating}
                        onClick={() => {
                          setConfirmApply(false);
                          setLocalActionError(null);
                          setConfirmRollback(true);
                        }}
                      >
                        <RotateCcw className="h-4 w-4" aria-hidden="true" />
                        <span>Rollback Latest Backup</span>
                      </Button>
                    ) : (
                      <div className="rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] p-3">
                        <p className="text-xs font-medium text-[var(--text-primary)]">
                          Restore newest backup for {activeDescriptor?.name ?? selectedTool}?
                        </p>
                        <div className="mt-3 flex items-center gap-2">
                          <Button
                            type="button"
                            variant="destructive"
                            size="sm"
                            disabled={isMutating}
                            onClick={handleConfirmRollback}
                          >
                            {rollbackMutation.isPending ? "Restoring..." : "Confirm Rollback"}
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            size="sm"
                            disabled={isMutating}
                            onClick={() => setConfirmRollback(false)}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    )}
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>

        <div className="lg:col-span-7">
          <Card className="flex h-full flex-col">
            <CardHeader className="flex-wrap gap-3">
              <div>
                <div className="flex items-center gap-2">
                  <FileCode2 className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                  <CardTitle>3. Generated Configuration Preview</CardTitle>
                </div>
                <CardDescription className="mt-1 font-mono text-xs">
                  {activePreview
                    ? `${activePreview.target_path} (${activePreview.filename})`
                    : activeDescriptor
                      ? `${activeDescriptor.target_path} (${activeDescriptor.filename})`
                      : "Select a tool and generate preview"}
                </CardDescription>
              </div>

              <div className="flex flex-wrap items-center gap-2">
                {activePreview ? (
                  <Badge variant="default">{formatBadgeLabel(activePreview.format)}</Badge>
                ) : activeDescriptor ? (
                  <Badge variant="secondary">{formatBadgeLabel(activeDescriptor.format)}</Badge>
                ) : null}
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  disabled={!activePreview}
                  onClick={() => void handleCopyConfig()}
                >
                  {copiedTool === activePreview?.tool && activePreview ? (
                    <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" aria-hidden="true" />
                  ) : (
                    <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                  )}
                  <span>
                    {copiedTool === activePreview?.tool && activePreview ? "Copied" : "Copy Config"}
                  </span>
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={!activePreview}
                  onClick={handleDownloadConfig}
                >
                  <Download className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                  <span>
                    Download {activePreview?.filename ?? activeDescriptor?.filename ?? "File"}
                  </span>
                </Button>
              </div>
            </CardHeader>

            <CardContent className="flex flex-1 flex-col gap-4 pt-5">
              {generateMutation.isError ? (
                <Alert variant="destructive">
                  <AlertTriangle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertTitle>Preview generation failed</AlertTitle>
                    <AlertDescription className="font-mono break-all">
                      {generateMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {copyError ? (
                <Alert variant="warning">
                  <AlertTriangle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertDescription>{copyError}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {previewStale && !generateMutation.isPending ? (
                <Alert variant="warning">
                  <AlertTriangle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 flex-1">
                    <AlertTitle>
                      Preview out of date for {activeDescriptor?.name ?? selectedTool}
                    </AlertTitle>
                    <AlertDescription>
                      Fix any invalid profile parameters and click Generate Preview to render the
                      configuration for {activeDescriptor?.name ?? selectedTool}.
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="flex-1 overflow-x-auto rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                {generateMutation.isPending && !activePreview ? (
                  <p className="font-mono text-xs text-[var(--text-muted)]">
                    Rendering configuration preview for {activeDescriptor?.name ?? selectedTool}...
                  </p>
                ) : activePreview ? (
                  <pre className="font-mono text-xs leading-relaxed text-[var(--text-primary)]">
                    <code>{activePreview.content}</code>
                  </pre>
                ) : previewStale ? (
                  <p className="font-mono text-xs text-[var(--text-muted)]">
                    No current preview available for {activeDescriptor?.name ?? selectedTool}. Enter
                    valid parameters and click Generate Preview.
                  </p>
                ) : (
                  <p className="font-mono text-xs text-[var(--text-muted)]">
                    Enter a valid gateway base URL and click Generate Preview.
                  </p>
                )}
              </div>

              {activePreview ? (
                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-4 py-3">
                  <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                    Setup Instructions
                  </p>
                  <p className="mt-1 text-xs leading-relaxed text-[var(--text-secondary)]">
                    {activePreview.instructions}
                  </p>
                </div>
              ) : null}
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  );
}
