"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertCircle,
  Check,
  Copy,
  KeyRound,
  Plus,
  Power,
  RefreshCw,
  RotateCw,
  ShieldAlert,
  Trash2,
  X,
} from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useApiKeys,
  useCreateApiKeyMutation,
  useDeleteApiKeyMutation,
  useRotateApiKeyMutation,
  useUpdateApiKeyMutation,
} from "./hooks";
import {
  type ApiKeySummary,
  type CreateKeyFormInput,
  type CreateKeyFormOutput,
  createKeySchema,
  DEFAULT_CREATE_KEY_VALUES,
  toCreateApiKeyRequest,
} from "./schemas";

interface RevealedSecretState {
  action: "created" | "rotated";
  keyName: string;
  keyPrefix: string;
  secret: string;
}

type PendingRowConfirm = { id: string; kind: "rotate" } | { id: string; kind: "delete" } | null;

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
  });
}

function isKeyExpired(key: ApiKeySummary): boolean {
  if (!key.expires_at) {
    return false;
  }
  const expiresMs = Date.parse(key.expires_at);
  return !Number.isNaN(expiresMs) && expiresMs <= Date.now();
}

function formatQuotaUsd(quotaUsd: number): string {
  if (quotaUsd <= 0) {
    return "Unlimited";
  }
  return `$${quotaUsd.toFixed(2)}`;
}

export function KeysConsole() {
  const keysQuery = useApiKeys();
  const createMutation = useCreateApiKeyMutation();
  const updateMutation = useUpdateApiKeyMutation();
  const rotateMutation = useRotateApiKeyMutation();
  const deleteMutation = useDeleteApiKeyMutation();

  const [revealedSecret, setRevealedSecret] = useState<RevealedSecretState | null>(null);
  const [copiedSecret, setCopiedSecret] = useState(false);
  const [copyError, setCopyError] = useState<string | null>(null);
  const [pendingConfirm, setPendingConfirm] = useState<PendingRowConfirm>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [busyKeyId, setBusyKeyId] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateKeyFormInput, unknown, CreateKeyFormOutput>({
    resolver: zodResolver(createKeySchema),
    defaultValues: DEFAULT_CREATE_KEY_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const keys = keysQuery.data ?? [];
  const activeCount = keys.filter((item) => !item.disabled && !isKeyExpired(item)).length;
  const disabledCount = keys.filter((item) => item.disabled).length;

  const handleCopySecret = async () => {
    if (!revealedSecret) {
      return;
    }
    setCopyError(null);
    try {
      await navigator.clipboard.writeText(revealedSecret.secret);
      setCopiedSecret(true);
      window.setTimeout(() => setCopiedSecret(false), 2500);
    } catch {
      setCopyError("Clipboard access was denied. Select and copy the secret manually.");
    }
  };

  const handleDismissSecret = () => {
    setRevealedSecret(null);
    setCopiedSecret(false);
    setCopyError(null);
  };

  const onCreateSubmit = (values: CreateKeyFormOutput) => {
    setRowError(null);
    createMutation.mutate(toCreateApiKeyRequest(values), {
      onSuccess: (response) => {
        setRevealedSecret({
          action: "created",
          keyName: response.key.name,
          keyPrefix: response.key.key_prefix,
          secret: response.secret,
        });
        setCopiedSecret(false);
        setCopyError(null);
        reset(DEFAULT_CREATE_KEY_VALUES);
      },
    });
  };

  const handleToggleDisabled = (key: ApiKeySummary) => {
    if (busyKeyId !== null) {
      return;
    }
    setRowError(null);
    setPendingConfirm(null);
    setBusyKeyId(key.id);
    updateMutation.mutate(
      { id: key.id, body: { disabled: !key.disabled } },
      {
        onError: (err) => setRowError(err.message),
        onSettled: () => setBusyKeyId(null),
      },
    );
  };

  const handleConfirmRotate = (key: ApiKeySummary) => {
    if (busyKeyId !== null) {
      return;
    }
    setRowError(null);
    setBusyKeyId(key.id);
    rotateMutation.mutate(key.id, {
      onSuccess: (response) => {
        setRevealedSecret({
          action: "rotated",
          keyName: response.key.name,
          keyPrefix: response.key.key_prefix,
          secret: response.secret,
        });
        setCopiedSecret(false);
        setCopyError(null);
        setPendingConfirm(null);
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyKeyId(null),
    });
  };

  const handleConfirmDelete = (key: ApiKeySummary) => {
    if (busyKeyId !== null) {
      return;
    }
    setRowError(null);
    setBusyKeyId(key.id);
    deleteMutation.mutate(key.id, {
      onSuccess: () => {
        setPendingConfirm(null);
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyKeyId(null),
    });
  };

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Gateway Credentials
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            API Keys
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Issue, constrain, rotate, and revoke client gateway tokens stored as SHA-256 hashes.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <Badge variant="outline" className="tabular-nums">
            {keys.length} Total
          </Badge>
          <Badge variant="success" className="tabular-nums">
            {activeCount} Active
          </Badge>
          {disabledCount > 0 ? (
            <Badge variant="warning" className="tabular-nums">
              {disabledCount} Disabled
            </Badge>
          ) : null}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void keysQuery.refetch()}
            disabled={keysQuery.isFetching}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
            <span>{keysQuery.isFetching ? "Refreshing" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {revealedSecret ? (
        <Card className="border-[var(--accent-primary)] bg-[var(--bg-surface-raised)]">
          <CardHeader className="border-[var(--border-strong)]">
            <div className="flex items-center gap-2">
              <ShieldAlert className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <CardTitle>
                {revealedSecret.action === "created"
                  ? "One-Time API Key Secret Created"
                  : "One-Time Rotated API Key Secret"}
              </CardTitle>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="xs"
              onClick={handleDismissSecret}
              aria-label="Dismiss one-time secret panel"
            >
              <X className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Dismiss</span>
            </Button>
          </CardHeader>
          <CardContent className="space-y-3">
            <p className="text-sm text-[var(--text-secondary)]">
              Copy this token for{" "}
              <span className="font-medium text-[var(--text-primary)]">
                {revealedSecret.keyName}
              </span>{" "}
              now. Once dismissed, the raw secret cannot be recovered from the server.
            </p>
            <div className="flex flex-col gap-2.5 sm:flex-row sm:items-center">
              <code className="block flex-1 overflow-x-auto rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3.5 py-2.5 font-mono text-sm text-[var(--text-primary)] select-all">
                {revealedSecret.secret}
              </code>
              <Button type="button" onClick={() => void handleCopySecret()} className="shrink-0">
                {copiedSecret ? (
                  <>
                    <Check className="h-4 w-4" aria-hidden="true" />
                    <span>Copied</span>
                  </>
                ) : (
                  <>
                    <Copy className="h-4 w-4" aria-hidden="true" />
                    <span>Copy Secret</span>
                  </>
                )}
              </Button>
            </div>
            {copyError ? (
              <p className="font-mono text-xs text-[var(--status-warn)]">{copyError}</p>
            ) : null}
          </CardContent>
        </Card>
      ) : null}

      <Card>
        <CardHeader>
          <div>
            <div className="flex items-center gap-2">
              <Plus className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <CardTitle>Create Gateway API Key</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Leave allowed models blank to permit all configured models. Use 0 for unlimited RPM or
              USD quota.
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">POST /api/keys</span>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onCreateSubmit)} noValidate className="space-y-4">
            {createMutation.error ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Failed to create API key</AlertTitle>
                  <AlertDescription className="font-mono">
                    {createMutation.error.message}
                  </AlertDescription>
                </div>
              </Alert>
            ) : null}

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
              <div>
                <Label htmlFor="key-name">Key Name</Label>
                <Input
                  id="key-name"
                  type="text"
                  placeholder="e.g. production-backend"
                  aria-invalid={Boolean(errors.name)}
                  className="mt-1.5"
                  {...register("name")}
                />
                {errors.name ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">{errors.name.message}</p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="key-allowed-models">Allowed Models (comma-separated)</Label>
                <Input
                  id="key-allowed-models"
                  type="text"
                  placeholder="gpt-4o, claude-sonnet-4"
                  aria-invalid={Boolean(errors.allowed_models)}
                  className="mt-1.5"
                  {...register("allowed_models")}
                />
                {errors.allowed_models ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.allowed_models.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="key-tags">Tags (comma-separated)</Label>
                <Input
                  id="key-tags"
                  type="text"
                  placeholder="prod, cli, team-a"
                  aria-invalid={Boolean(errors.tags)}
                  className="mt-1.5"
                  {...register("tags")}
                />
                {errors.tags ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">{errors.tags.message}</p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="key-rpm">Rate Limit (RPM)</Label>
                <Input
                  id="key-rpm"
                  type="number"
                  min={0}
                  step={1}
                  placeholder="0"
                  aria-invalid={Boolean(errors.rate_limit_rpm)}
                  className="mt-1.5 tabular-nums"
                  {...register("rate_limit_rpm")}
                />
                {errors.rate_limit_rpm ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.rate_limit_rpm.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="key-quota">Quota (USD)</Label>
                <Input
                  id="key-quota"
                  type="number"
                  min={0}
                  step="0.01"
                  placeholder="0"
                  aria-invalid={Boolean(errors.quota_usd)}
                  className="mt-1.5 tabular-nums"
                  {...register("quota_usd")}
                />
                {errors.quota_usd ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.quota_usd.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="key-expires">Expires At (optional)</Label>
                <Input
                  id="key-expires"
                  type="datetime-local"
                  aria-invalid={Boolean(errors.expires_at)}
                  className="mt-1.5"
                  {...register("expires_at")}
                />
                {errors.expires_at ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.expires_at.message}
                  </p>
                ) : null}
              </div>
            </div>

            <div className="flex items-center justify-end gap-3 pt-1">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => {
                  createMutation.reset();
                  reset(DEFAULT_CREATE_KEY_VALUES);
                }}
              >
                Reset
              </Button>
              <Button type="submit" size="sm" disabled={createMutation.isPending}>
                <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{createMutation.isPending ? "Creating..." : "Create API Key"}</span>
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <EmberRule />

      {rowError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="flex-1">
            <AlertTitle>Key operation failed</AlertTitle>
            <AlertDescription className="font-mono">{rowError}</AlertDescription>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setRowError(null)}
            aria-label="Dismiss error"
          >
            <X className="h-3.5 w-3.5" aria-hidden="true" />
          </Button>
        </Alert>
      ) : null}

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <KeyRound className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <CardTitle>Issued API Keys</CardTitle>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">GET /api/keys</span>
        </CardHeader>

        {keysQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load API keys</AlertTitle>
                <AlertDescription className="font-mono">{keysQuery.error.message}</AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : keysQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading issued API keys...
          </CardContent>
        ) : keys.length === 0 ? (
          <CardContent className="py-10 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              No gateway API keys issued yet
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              Create your first key above to authenticate OpenAI or Anthropic compatible client
              requests.
            </p>
          </CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <th scope="col" className="px-5 py-3">
                    Key &amp; Status
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Prefix
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Allowed Models
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Limits
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Activity
                  </th>
                  <th scope="col" className="px-5 py-3 text-right">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {keys.map((key) => {
                  const expired = isKeyExpired(key);
                  const isRowBusy = busyKeyId === key.id;
                  const anyRowBusy = busyKeyId !== null;
                  const confirmingRotate =
                    pendingConfirm?.id === key.id && pendingConfirm.kind === "rotate";
                  const confirmingDelete =
                    pendingConfirm?.id === key.id && pendingConfirm.kind === "delete";

                  return (
                    <tr
                      key={key.id}
                      className="align-top transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3.5">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className="font-medium text-[var(--text-primary)]">{key.name}</span>
                          {key.disabled ? (
                            <Badge variant="warning">Disabled</Badge>
                          ) : expired ? (
                            <Badge variant="destructive">Expired</Badge>
                          ) : (
                            <Badge variant="success">Active</Badge>
                          )}
                        </div>
                        {key.tags.length > 0 ? (
                          <div className="mt-1.5 flex flex-wrap gap-1.5">
                            {key.tags.map((tag) => (
                              <Badge key={tag} variant="secondary">
                                {tag}
                              </Badge>
                            ))}
                          </div>
                        ) : null}
                      </td>

                      <td className="px-4 py-3.5">
                        <code className="rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 py-0.5 font-mono text-xs text-[var(--text-primary)]">
                          {key.key_prefix}...
                        </code>
                      </td>

                      <td className="px-4 py-3.5">
                        {key.allowed_models.length === 0 ? (
                          <Badge variant="outline">All Models</Badge>
                        ) : (
                          <div className="flex max-w-xs flex-wrap gap-1.5">
                            {key.allowed_models.map((model) => (
                              <Badge key={model} variant="default">
                                {model}
                              </Badge>
                            ))}
                          </div>
                        )}
                      </td>

                      <td className="px-4 py-3.5">
                        <div className="flex flex-wrap gap-1.5">
                          <Badge variant="outline" className="tabular-nums">
                            {key.rate_limit_rpm > 0 ? `${key.rate_limit_rpm} RPM` : "No RPM Cap"}
                          </Badge>
                          <Badge variant="outline" className="tabular-nums">
                            Quota: {formatQuotaUsd(key.quota_usd)}
                          </Badge>
                        </div>
                      </td>

                      <td className="px-4 py-3.5 font-mono text-xs text-[var(--text-muted)]">
                        <div>Created: {formatTimestamp(key.created_at)}</div>
                        <div className="mt-0.5">Last used: {formatTimestamp(key.last_used_at)}</div>
                        {key.expires_at ? (
                          <div className="mt-0.5 text-[var(--text-secondary)]">
                            Expires: {formatTimestamp(key.expires_at)}
                          </div>
                        ) : null}
                      </td>

                      <td className="px-5 py-3.5 text-right">
                        {confirmingRotate ? (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              size="xs"
                              onClick={() => handleConfirmRotate(key)}
                              disabled={anyRowBusy}
                            >
                              {isRowBusy ? "Rotating..." : "Confirm Rotate"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => setPendingConfirm(null)}
                              disabled={anyRowBusy}
                            >
                              Cancel
                            </Button>
                          </div>
                        ) : confirmingDelete ? (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => handleConfirmDelete(key)}
                              disabled={anyRowBusy}
                            >
                              {isRowBusy ? "Deleting..." : "Confirm Delete"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => setPendingConfirm(null)}
                              disabled={anyRowBusy}
                            >
                              Cancel
                            </Button>
                          </div>
                        ) : (
                          <div className="flex flex-wrap items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => handleToggleDisabled(key)}
                              disabled={anyRowBusy}
                            >
                              <Power className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>{key.disabled ? "Enable" : "Disable"}</span>
                            </Button>
                            <Button
                              type="button"
                              variant="outline"
                              size="xs"
                              onClick={() =>
                                setPendingConfirm({
                                  id: key.id,
                                  kind: "rotate",
                                })
                              }
                              disabled={anyRowBusy}
                            >
                              <RotateCw className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>Rotate</span>
                            </Button>
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() =>
                                setPendingConfirm({
                                  id: key.id,
                                  kind: "delete",
                                })
                              }
                              disabled={anyRowBusy}
                              aria-label={`Delete API key ${key.name}`}
                            >
                              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>Delete</span>
                            </Button>
                          </div>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
}
