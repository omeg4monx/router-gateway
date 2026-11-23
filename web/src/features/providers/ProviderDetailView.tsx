"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Activity,
  AlertCircle,
  ArrowLeft,
  CheckCircle2,
  Eye,
  EyeOff,
  KeyRound,
  RefreshCw,
  Save,
  Server,
  Trash2,
} from "lucide-react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useDeleteProviderConnectionMutation,
  useProviderConnection,
  useTestProviderConnectionMutation,
  useUpdateProviderConnectionMutation,
} from "./hooks";
import {
  connectionToEditFormValues,
  DEFAULT_UPDATE_CONNECTION_VALUES,
  type TestProviderConnectionResponse,
  toUpdateConnectionRequest,
  type UpdateConnectionFormInput,
  type UpdateConnectionFormOutput,
  updateConnectionSchema,
} from "./schemas";

const TEXTAREA_CLASS_NAME =
  "flex min-h-[92px] w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-xs text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]";

function formatTimestamp(iso: string | null): string {
  if (!iso) {
    return "Never";
  }
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString([], {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

export function ProviderDetailView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const connectionId = (searchParams.get("id") ?? "").trim();

  const connectionQuery = useProviderConnection(connectionId);
  const updateMutation = useUpdateProviderConnectionMutation();
  const deleteMutation = useDeleteProviderConnectionMutation();
  const testMutation = useTestProviderConnectionMutation();

  const [saveFeedback, setSaveFeedback] = useState<string | null>(null);
  const [liveTestResult, setLiveTestResult] = useState<TestProviderConnectionResponse | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAccessToken, setShowAccessToken] = useState(false);
  const [showRefreshToken, setShowRefreshToken] = useState(false);

  const editForm = useForm<UpdateConnectionFormInput, unknown, UpdateConnectionFormOutput>({
    resolver: zodResolver(updateConnectionSchema),
    defaultValues: DEFAULT_UPDATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const {
    reset,
    formState: { errors, isDirty },
  } = editForm;

  useEffect(() => {
    if (connectionQuery.data && !isDirty) {
      reset(connectionToEditFormValues(connectionQuery.data));
    }
  }, [connectionQuery.data, isDirty, reset]);

  if (!connectionId) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Provider Connection Not Specified</CardTitle>
          <Button asChild variant="outline" size="xs">
            <Link href="/dashboard/providers/">
              <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Back to Providers</span>
            </Link>
          </Button>
        </CardHeader>
        <CardContent>
          <Alert variant="warning">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
              aria-hidden="true"
            />
            <div>
              <AlertTitle>Missing connection query parameter</AlertTitle>
              <AlertDescription>
                Open a connection from the Providers Control Studio to inspect or rotate its
                configuration.
              </AlertDescription>
            </div>
          </Alert>
        </CardContent>
      </Card>
    );
  }

  const connection = connectionQuery.data;

  const onSaveSubmit = (values: UpdateConnectionFormOutput) => {
    if (!connection) {
      return;
    }
    setSaveFeedback(null);
    updateMutation.reset();

    updateMutation.mutate(
      {
        id: connection.id,
        body: toUpdateConnectionRequest(values, connection.auth_type),
      },
      {
        onSuccess: (updated) => {
          reset(connectionToEditFormValues(updated));
          setSaveFeedback(`Saved changes to "${updated.name}".`);
        },
      },
    );
  };

  const handleRunTest = () => {
    if (!connection || testMutation.isPending) {
      return;
    }
    testMutation.reset();
    testMutation.mutate(connection.id, {
      onSuccess: (result) => {
        setLiveTestResult(result);
      },
    });
  };

  const handleConfirmDelete = () => {
    if (!connection || deleteMutation.isPending) {
      return;
    }
    deleteMutation.mutate(connection.id, {
      onSuccess: () => {
        router.push("/dashboard/providers/");
      },
    });
  };

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <div className="flex items-center gap-2">
            <Button asChild variant="ghost" size="xs" className="-ml-2">
              <Link href="/dashboard/providers/">
                <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Providers Studio</span>
              </Link>
            </Button>
            <span className="font-mono text-xs text-[var(--text-muted)]">/ {connectionId}</span>
          </div>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            {connection ? connection.name : "Provider Connection Detail"}
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Configure routing priority, allowed model list, upstream headers, and rotate vaulted
            credentials.
          </p>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => void connectionQuery.refetch()}
          disabled={connectionQuery.isFetching}
        >
          <RefreshCw className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
          <span>{connectionQuery.isFetching ? "Refreshing" : "Refresh"}</span>
        </Button>
      </section>

      {connectionQuery.isError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div>
            <AlertTitle>Unable to load provider connection</AlertTitle>
            <AlertDescription className="font-mono">
              {connectionQuery.error.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : connectionQuery.isPending || !connection ? (
        <Card>
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading connection details...
          </CardContent>
        </Card>
      ) : (
        <>
          <Card>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Server className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
                <CardTitle>Connection Telemetry &amp; Connectivity</CardTitle>
              </div>
              <Button
                type="button"
                size="xs"
                onClick={handleRunTest}
                disabled={testMutation.isPending}
              >
                <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{testMutation.isPending ? "Testing..." : "Run Connectivity Test"}</span>
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Status &amp; Protocol
                  </p>
                  <div className="mt-2 flex flex-wrap gap-1.5">
                    <Badge variant={connection.enabled ? "success" : "warning"}>
                      {connection.enabled ? "Enabled" : "Disabled"}
                    </Badge>
                    <Badge variant="outline">{connection.protocol}</Badge>
                  </div>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Adapter &amp; Auth
                  </p>
                  <div className="mt-2 flex flex-wrap gap-1.5">
                    <Badge variant="secondary">{connection.provider_id}</Badge>
                    <Badge variant={connection.auth_type === "oauth" ? "default" : "secondary"}>
                      {connection.auth_type === "oauth" ? "OAuth" : "API Key"}
                    </Badge>
                    <Badge variant={connection.has_credentials ? "success" : "destructive"}>
                      {connection.has_credentials ? "Vaulted" : "Missing"}
                    </Badge>
                  </div>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Last Test
                  </p>
                  <div className="mt-2 flex items-center gap-2">
                    <Badge
                      variant={
                        connection.last_test_status === "ok"
                          ? "success"
                          : connection.last_test_status === "error"
                            ? "destructive"
                            : "outline"
                      }
                    >
                      {connection.last_test_status || "Untested"}
                    </Badge>
                  </div>
                  <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
                    {formatTimestamp(connection.last_tested_at)}
                  </p>
                </div>

                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                  <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Timestamps
                  </p>
                  <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
                    Updated: {formatTimestamp(connection.updated_at)}
                  </p>
                  {connection.token_expires_at ? (
                    <p className="mt-0.5 font-mono text-xs text-[var(--text-secondary)]">
                      Token Exp: {formatTimestamp(connection.token_expires_at)}
                    </p>
                  ) : null}
                </div>
              </div>

              {testMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Connectivity test request failed</AlertTitle>
                    <AlertDescription className="font-mono">
                      {testMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {liveTestResult ? (
                <Alert variant={liveTestResult.ok ? "default" : "destructive"}>
                  {liveTestResult.ok ? (
                    <CheckCircle2
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                      aria-hidden="true"
                    />
                  ) : (
                    <AlertCircle
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                  )}
                  <div>
                    <AlertTitle>
                      {liveTestResult.ok
                        ? `Live Test Succeeded (HTTP ${liveTestResult.status_code} · ${liveTestResult.latency_ms}ms)`
                        : `Live Test Failed (${liveTestResult.status_code ? `HTTP ${liveTestResult.status_code}` : "Network Error"} · ${liveTestResult.latency_ms}ms)`}
                    </AlertTitle>
                    {liveTestResult.error ? (
                      <AlertDescription className="font-mono">
                        {liveTestResult.error}
                      </AlertDescription>
                    ) : null}
                  </div>
                </Alert>
              ) : connection.last_test_error ? (
                <Alert variant="warning">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Last recorded test error</AlertTitle>
                    <AlertDescription className="font-mono">
                      {connection.last_test_error}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}
            </CardContent>
          </Card>

          <EmberRule />

          <Card>
            <CardHeader>
              <div>
                <div className="flex items-center gap-2">
                  <KeyRound className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
                  <CardTitle>Configuration &amp; Credential Rotation</CardTitle>
                </div>
                <CardDescription className="mt-0.5">
                  Leave credential fields empty to keep existing encrypted secrets untouched.
                </CardDescription>
              </div>
              <span className="font-mono text-xs text-[var(--text-muted)]">
                PATCH /api/providers/connections/{connection.id}
              </span>
            </CardHeader>
            <CardContent>
              <form onSubmit={editForm.handleSubmit(onSaveSubmit)} noValidate className="space-y-5">
                {updateMutation.error ? (
                  <Alert variant="destructive">
                    <AlertCircle
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                    <div>
                      <AlertTitle>Failed to update connection</AlertTitle>
                      <AlertDescription className="font-mono">
                        {updateMutation.error.message}
                      </AlertDescription>
                    </div>
                  </Alert>
                ) : null}

                {saveFeedback ? (
                  <Alert>
                    <CheckCircle2
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                      aria-hidden="true"
                    />
                    <div>
                      <AlertTitle>Connection Updated</AlertTitle>
                      <AlertDescription>{saveFeedback}</AlertDescription>
                    </div>
                  </Alert>
                ) : null}

                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div>
                    <Label htmlFor="detail-name">Connection Name</Label>
                    <Input
                      id="detail-name"
                      type="text"
                      aria-invalid={Boolean(errors.name)}
                      className="mt-1.5"
                      {...editForm.register("name")}
                    />
                    {errors.name ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.name.message}
                      </p>
                    ) : null}
                  </div>

                  <div>
                    <Label htmlFor="detail-priority">Routing Priority (lower runs first)</Label>
                    <Input
                      id="detail-priority"
                      type="number"
                      min={0}
                      max={100000}
                      step={1}
                      aria-invalid={Boolean(errors.priority)}
                      className="mt-1.5 tabular-nums"
                      {...editForm.register("priority")}
                    />
                    {errors.priority ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.priority.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-base-url">Upstream Base URL</Label>
                    <Input
                      id="detail-base-url"
                      type="url"
                      aria-invalid={Boolean(errors.base_url)}
                      className="mt-1.5"
                      {...editForm.register("base_url")}
                    />
                    {errors.base_url ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.base_url.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-models">
                      Allowed / Advertised Models (comma-separated)
                    </Label>
                    <Input
                      id="detail-models"
                      type="text"
                      aria-invalid={Boolean(errors.models)}
                      className="mt-1.5"
                      {...editForm.register("models")}
                    />
                    {errors.models ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.models.message}
                      </p>
                    ) : null}
                  </div>

                  <div className="md:col-span-2">
                    <Label htmlFor="detail-headers">
                      Custom Upstream Headers (one Key: Value per line)
                    </Label>
                    <textarea
                      id="detail-headers"
                      aria-invalid={Boolean(errors.headers)}
                      className={`${TEXTAREA_CLASS_NAME} mt-1.5`}
                      {...editForm.register("headers")}
                    />
                    {errors.headers ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.headers.message}
                      </p>
                    ) : null}
                  </div>
                </div>

                {connection.auth_type === "api_key" ? (
                  <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                    <Label htmlFor="detail-api-key">Rotate API Key (optional)</Label>
                    <div className="relative mt-1.5">
                      <Input
                        id="detail-api-key"
                        type={showApiKey ? "text" : "password"}
                        autoComplete="off"
                        placeholder="Leave blank to keep existing vaulted API key"
                        className="pr-10"
                        {...editForm.register("api_key")}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon-sm"
                        onClick={() => setShowApiKey((prev) => !prev)}
                        aria-label={showApiKey ? "Hide API key" : "Show API key"}
                        className="absolute top-1.5 right-1.5"
                      >
                        {showApiKey ? (
                          <EyeOff className="h-4 w-4" aria-hidden="true" />
                        ) : (
                          <Eye className="h-4 w-4" aria-hidden="true" />
                        )}
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span className="text-sm font-medium text-[var(--text-primary)]">
                        OAuth Token Rotation (optional)
                      </span>
                      <Link
                        href="/dashboard/providers/"
                        className="font-mono text-xs text-[var(--accent-text)] hover:underline"
                      >
                        Re-authorize via OAuth Wizard
                      </Link>
                    </div>
                    <div className="grid grid-cols-1 gap-3.5 md:grid-cols-2">
                      <div>
                        <Label htmlFor="detail-access-token">New Access Token</Label>
                        <div className="relative mt-1.5">
                          <Input
                            id="detail-access-token"
                            type={showAccessToken ? "text" : "password"}
                            autoComplete="off"
                            placeholder="Leave blank to keep existing"
                            className="pr-10"
                            {...editForm.register("access_token")}
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setShowAccessToken((prev) => !prev)}
                            aria-label={showAccessToken ? "Hide access token" : "Show access token"}
                            className="absolute top-1.5 right-1.5"
                          >
                            {showAccessToken ? (
                              <EyeOff className="h-4 w-4" aria-hidden="true" />
                            ) : (
                              <Eye className="h-4 w-4" aria-hidden="true" />
                            )}
                          </Button>
                        </div>
                      </div>

                      <div>
                        <Label htmlFor="detail-refresh-token">New Refresh Token</Label>
                        <div className="relative mt-1.5">
                          <Input
                            id="detail-refresh-token"
                            type={showRefreshToken ? "text" : "password"}
                            autoComplete="off"
                            placeholder="Leave blank to keep existing"
                            className="pr-10"
                            {...editForm.register("refresh_token")}
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => setShowRefreshToken((prev) => !prev)}
                            aria-label={
                              showRefreshToken ? "Hide refresh token" : "Show refresh token"
                            }
                            className="absolute top-1.5 right-1.5"
                          >
                            {showRefreshToken ? (
                              <EyeOff className="h-4 w-4" aria-hidden="true" />
                            ) : (
                              <Eye className="h-4 w-4" aria-hidden="true" />
                            )}
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                <label
                  htmlFor="detail-enabled"
                  className="flex cursor-pointer items-center justify-between rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-3"
                >
                  <span className="text-sm text-[var(--text-primary)]">
                    Connection enabled for gateway routing
                  </span>
                  <input
                    id="detail-enabled"
                    type="checkbox"
                    className="h-4 w-4 accent-[var(--accent-primary)]"
                    {...editForm.register("enabled")}
                  />
                </label>

                <div className="flex items-center justify-end gap-3 pt-1">
                  <Button
                    type="button"
                    variant="secondary"
                    size="sm"
                    onClick={() => {
                      updateMutation.reset();
                      setSaveFeedback(null);
                      reset(connectionToEditFormValues(connection));
                    }}
                  >
                    Discard Changes
                  </Button>
                  <Button type="submit" size="sm" disabled={updateMutation.isPending}>
                    <Save className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{updateMutation.isPending ? "Saving..." : "Save Changes"}</span>
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          <Card className="border-[var(--status-error-border)]">
            <CardHeader>
              <div>
                <CardTitle>Delete Connection</CardTitle>
                <CardDescription className="mt-0.5">
                  Permanently remove this connection and purge its vaulted credentials.
                </CardDescription>
              </div>
              {confirmDelete ? (
                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant="destructive"
                    size="xs"
                    onClick={handleConfirmDelete}
                    disabled={deleteMutation.isPending}
                  >
                    {deleteMutation.isPending ? "Deleting..." : "Confirm Delete"}
                  </Button>
                  <Button
                    type="button"
                    variant="secondary"
                    size="xs"
                    onClick={() => setConfirmDelete(false)}
                    disabled={deleteMutation.isPending}
                  >
                    Cancel
                  </Button>
                </div>
              ) : (
                <Button
                  type="button"
                  variant="destructive"
                  size="xs"
                  onClick={() => setConfirmDelete(true)}
                >
                  <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Delete Connection</span>
                </Button>
              )}
            </CardHeader>
            {deleteMutation.error ? (
              <CardContent>
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to delete connection</AlertTitle>
                    <AlertDescription className="font-mono">
                      {deleteMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              </CardContent>
            ) : null}
          </Card>
        </>
      )}
    </div>
  );
}
