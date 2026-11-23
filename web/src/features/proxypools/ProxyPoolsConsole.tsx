"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Activity, Lock, Network, Pencil, Plus, Radar, Trash2 } from "lucide-react";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useProvidersCatalog } from "@/features/providers/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useCreateProxyPoolMutation,
  useDeleteProxyBindingMutation,
  useDeleteProxyPoolMutation,
  useProxyPools,
  useScanProxyCandidatesMutation,
  useTestProxyPoolMutation,
  useUpdateProxyPoolMutation,
  useUpsertProxyBindingMutation,
} from "./hooks";
import {
  type CreateProxyPoolFormInput,
  type CreateProxyPoolFormValues,
  createProxyPoolSchema,
  parseCandidateLines,
  type ScanCandidatesFormInput,
  type ScanCandidatesFormValues,
  type ScanProxyCandidateResult,
  scanCandidatesSchema,
} from "./schemas";

function statusBadgeVariant(status: string): "success" | "destructive" | "secondary" {
  if (status === "healthy") {
    return "success";
  }
  if (status === "unhealthy") {
    return "destructive";
  }
  return "secondary";
}

export function ProxyPoolsConsole() {
  const poolsQuery = useProxyPools();
  const catalogQuery = useProvidersCatalog();

  const createPool = useCreateProxyPoolMutation();
  const updatePool = useUpdateProxyPoolMutation();
  const deletePool = useDeleteProxyPoolMutation();
  const testPool = useTestProxyPoolMutation();
  const scanCandidates = useScanProxyCandidatesMutation();
  const upsertBinding = useUpsertProxyBindingMutation();
  const deleteBinding = useDeleteProxyBindingMutation();

  const [notice, setNotice] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [scanResults, setScanResults] = useState<ScanProxyCandidateResult[]>([]);
  const [testingPoolId, setTestingPoolId] = useState<string | null>(null);
  const [editingPoolId, setEditingPoolId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [editUrl, setEditUrl] = useState("");

  const createForm = useForm<CreateProxyPoolFormInput, unknown, CreateProxyPoolFormValues>({
    resolver: zodResolver(createProxyPoolSchema),
    defaultValues: { name: "", url: "", enabled: true },
  });

  const scanForm = useForm<ScanCandidatesFormInput, unknown, ScanCandidatesFormValues>({
    resolver: zodResolver(scanCandidatesSchema),
    defaultValues: { candidatesText: "", autoAdd: false, namePrefix: "Scanned", timeoutMs: 3000 },
  });

  const pools = poolsQuery.data?.pools ?? [];
  const bindings = poolsQuery.data?.bindings ?? [];
  const bindableProviders = useMemo(() => {
    const list: Array<{ id: string; label: string }> = [];
    for (const adapter of catalogQuery.data?.adapters ?? []) {
      list.push({ id: adapter.id, label: adapter.display_name });
    }
    for (const node of catalogQuery.data?.custom_nodes ?? []) {
      list.push({ id: node.id, label: `${node.name} (Custom)` });
    }
    return list;
  }, [catalogQuery.data]);

  const bindingByProvider = useMemo(() => {
    const map = new Map<string, string>();
    for (const b of bindings) {
      map.set(b.provider_id, b.pool_id);
    }
    return map;
  }, [bindings]);

  const onCreateSubmit = createForm.handleSubmit((values) => {
    setErrorMsg(null);
    setNotice(null);
    createPool.mutate(
      { name: values.name, url: values.url, enabled: values.enabled },
      {
        onSuccess: (created) => {
          createForm.reset({ name: "", url: "", enabled: true });
          setNotice(`Created proxy pool "${created.name}" (${created.display_url}).`);
        },
        onError: (err) => setErrorMsg(err.message),
      },
    );
  });

  const onScanSubmit = scanForm.handleSubmit((values) => {
    setErrorMsg(null);
    setNotice(null);
    const candidates = parseCandidateLines(values.candidatesText);
    scanCandidates.mutate(
      {
        candidates,
        auto_add: values.autoAdd,
        name_prefix: values.namePrefix || undefined,
        timeout_ms: values.timeoutMs,
      },
      {
        onSuccess: (res) => {
          setScanResults(res.results);
          const healthyCount = res.results.filter((r) => r.healthy).length;
          setNotice(`Scanned ${res.results.length} candidates (${healthyCount} healthy).`);
        },
        onError: (err) => setErrorMsg(err.message),
      },
    );
  });

  const handleTestPool = (id: string) => {
    setErrorMsg(null);
    setNotice(null);
    setTestingPoolId(id);
    testPool.mutate(
      { id },
      {
        onSuccess: (res) => {
          setTestingPoolId(null);
          setNotice(
            res.healthy
              ? `Proxy probe healthy (${res.latency_ms} ms).`
              : `Proxy probe failed: ${res.error ?? "unreachable"}`,
          );
        },
        onError: (err) => {
          setTestingPoolId(null);
          setErrorMsg(err.message);
        },
      },
    );
  };

  const handleBindingChange = (providerId: string, poolId: string) => {
    setErrorMsg(null);
    setNotice(null);
    if (!poolId) {
      deleteBinding.mutate(providerId, {
        onSuccess: () => setNotice(`Unbound ${providerId} (now dialing direct).`),
        onError: (err) => setErrorMsg(err.message),
      });
      return;
    }
    upsertBinding.mutate(
      { providerId, poolId },
      {
        onSuccess: () => setNotice(`Bound ${providerId} to selected proxy pool.`),
        onError: (err) => setErrorMsg(err.message),
      },
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold tracking-tight text-[var(--text-primary)]">
            Proxy Pools &amp; Provider Routing
          </h2>
          <p className="text-sm text-[var(--text-secondary)]">
            Encrypted HTTP/HTTPS/SOCKS5 outbound proxies with latency probes, batch auto-scanner,
            and per-provider bindings.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="secondary">{pools.length} Pools</Badge>
          <Badge variant="default">{bindings.length} Bound Providers</Badge>
        </div>
      </div>

      {notice ? (
        <Alert variant="default">
          <AlertTitle>Updated</AlertTitle>
          <AlertDescription>{notice}</AlertDescription>
        </Alert>
      ) : null}

      {errorMsg || poolsQuery.error ? (
        <Alert variant="destructive">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{errorMsg ?? poolsQuery.error?.message}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Plus className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              Add Proxy Pool
            </CardTitle>
            <CardDescription>
              Credentials in URLs are encrypted with AES-256-GCM at rest and redacted in responses.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={onCreateSubmit} className="space-y-4">
              <div className="space-y-1.5">
                <Label htmlFor="proxy-pool-name">Pool Name</Label>
                <Input
                  id="proxy-pool-name"
                  placeholder="SG Edge SOCKS5"
                  {...createForm.register("name")}
                />
                {createForm.formState.errors.name ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {createForm.formState.errors.name.message}
                  </p>
                ) : null}
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="proxy-pool-url">Proxy URL</Label>
                <Input
                  id="proxy-pool-url"
                  placeholder="socks5h://user:pass@proxy.internal:1080"
                  className="font-mono text-xs"
                  {...createForm.register("url")}
                />
                {createForm.formState.errors.url ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {createForm.formState.errors.url.message}
                  </p>
                ) : null}
              </div>

              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 text-sm text-[var(--text-secondary)]">
                  <input type="checkbox" {...createForm.register("enabled")} />
                  Enabled immediately
                </label>
                <Button type="submit" disabled={createPool.isPending}>
                  {createPool.isPending ? "Saving..." : "Create Pool"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Radar className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              Batch Auto-Scanner
            </CardTitle>
            <CardDescription>
              Probe up to 50 candidate proxy URLs concurrently (one per line) and optionally save
              healthy entries.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={onScanSubmit} className="space-y-3">
              <div className="space-y-1.5">
                <Label htmlFor="scan-candidates">Candidate URLs</Label>
                <textarea
                  id="scan-candidates"
                  rows={4}
                  placeholder={"http://127.0.0.1:8080\nsocks5h://127.0.0.1:1080"}
                  className="w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-2 font-mono text-xs text-[var(--text-primary)] focus:border-[var(--accent-primary)] focus:outline-none"
                  {...scanForm.register("candidatesText")}
                />
                {scanForm.formState.errors.candidatesText ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {scanForm.formState.errors.candidatesText.message}
                  </p>
                ) : null}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <Label htmlFor="scan-prefix">Name Prefix</Label>
                  <Input id="scan-prefix" {...scanForm.register("namePrefix")} />
                </div>
                <div className="space-y-1">
                  <Label htmlFor="scan-timeout">Timeout (ms)</Label>
                  <Input id="scan-timeout" type="number" {...scanForm.register("timeoutMs")} />
                </div>
              </div>

              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 text-sm text-[var(--text-secondary)]">
                  <input type="checkbox" {...scanForm.register("autoAdd")} />
                  Auto-add healthy pools
                </label>
                <Button type="submit" variant="secondary" disabled={scanCandidates.isPending}>
                  {scanCandidates.isPending ? "Scanning..." : "Run Scan"}
                </Button>
              </div>
            </form>

            {scanResults.length > 0 ? (
              <div className="mt-4 space-y-1.5 border-t border-[var(--border-subtle)] pt-3">
                <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                  Scan Results
                </p>
                {scanResults.map((res) => (
                  <div
                    key={`${res.display_url}-${res.scheme ?? "raw"}-${res.added_pool_id ?? res.error ?? "ok"}`}
                    className="flex items-center justify-between rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 py-1.5 text-xs"
                  >
                    <span className="truncate font-mono text-[var(--text-primary)]">
                      {res.display_url}
                    </span>
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-[var(--text-muted)]">{res.latency_ms}ms</span>
                      <Badge variant={res.healthy ? "success" : "destructive"}>
                        {res.healthy ? "healthy" : "failed"}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            ) : null}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Network className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            Configured Proxy Pools
          </CardTitle>
        </CardHeader>
        <CardContent>
          {pools.length === 0 ? (
            <p className="py-6 text-center text-sm text-[var(--text-muted)]">
              No proxy pools configured yet. Add one above or run the batch scanner.
            </p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead className="border-b border-[var(--border-subtle)] font-mono text-xs uppercase text-[var(--text-muted)]">
                  <tr>
                    <th className="py-2 pr-3">Name</th>
                    <th className="py-2 pr-3">Endpoint</th>
                    <th className="py-2 pr-3">Status</th>
                    <th className="py-2 pr-3">Latency</th>
                    <th className="py-2 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[var(--border-subtle)]">
                  {pools.map((pool) => {
                    const isEditing = editingPoolId === pool.id;
                    return (
                      <tr key={pool.id}>
                        <td className="py-3 pr-3 font-medium text-[var(--text-primary)]">
                          {isEditing ? (
                            <Input
                              aria-label="Edit pool name"
                              value={editName}
                              onChange={(e) => setEditName(e.target.value)}
                              className="h-8 text-xs"
                            />
                          ) : (
                            <div className="flex items-center gap-2">
                              <span>{pool.name}</span>
                              <Badge variant="default">{pool.scheme}</Badge>
                              {pool.has_credentials ? (
                                <Lock
                                  className="h-3.5 w-3.5 text-[var(--text-muted)]"
                                  aria-label="Credentials encrypted in vault"
                                />
                              ) : null}
                            </div>
                          )}
                        </td>
                        <td className="py-3 pr-3 font-mono text-xs text-[var(--text-secondary)]">
                          {isEditing ? (
                            <Input
                              aria-label="Replace proxy URL (optional)"
                              placeholder="Leave blank to keep existing URL"
                              value={editUrl}
                              onChange={(e) => setEditUrl(e.target.value)}
                              className="h-8 font-mono text-xs"
                            />
                          ) : (
                            pool.display_url
                          )}
                        </td>
                        <td className="py-3 pr-3">
                          <Badge variant={statusBadgeVariant(pool.status)}>{pool.status}</Badge>
                        </td>
                        <td className="py-3 pr-3 font-mono text-xs text-[var(--text-secondary)]">
                          {pool.last_tested_at ? `${pool.last_latency_ms} ms` : "—"}
                        </td>
                        <td className="py-3 text-right">
                          {isEditing ? (
                            <div className="inline-flex items-center gap-2">
                              <Button
                                type="button"
                                size="sm"
                                disabled={updatePool.isPending || !editName.trim()}
                                onClick={() => {
                                  setErrorMsg(null);
                                  setNotice(null);
                                  updatePool.mutate(
                                    {
                                      id: pool.id,
                                      body: {
                                        name: editName.trim(),
                                        ...(editUrl.trim() ? { url: editUrl.trim() } : {}),
                                      },
                                    },
                                    {
                                      onSuccess: (updated) => {
                                        setEditingPoolId(null);
                                        setEditUrl("");
                                        setNotice(`Updated proxy pool "${updated.name}".`);
                                      },
                                      onError: (err) => setErrorMsg(err.message),
                                    },
                                  );
                                }}
                              >
                                Save
                              </Button>
                              <Button
                                type="button"
                                size="sm"
                                variant="secondary"
                                onClick={() => {
                                  setEditingPoolId(null);
                                  setEditUrl("");
                                }}
                              >
                                Cancel
                              </Button>
                            </div>
                          ) : (
                            <div className="inline-flex items-center gap-2">
                              <Button
                                type="button"
                                size="sm"
                                variant="secondary"
                                onClick={() => {
                                  setEditingPoolId(pool.id);
                                  setEditName(pool.name);
                                  setEditUrl("");
                                }}
                                aria-label={`Edit ${pool.name}`}
                              >
                                <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
                              </Button>
                              <Button
                                type="button"
                                size="sm"
                                variant="secondary"
                                onClick={() => {
                                  setErrorMsg(null);
                                  setNotice(null);
                                  updatePool.mutate(
                                    {
                                      id: pool.id,
                                      body: { enabled: !pool.enabled },
                                    },
                                    {
                                      onSuccess: (updated) =>
                                        setNotice(
                                          `${updated.name} ${updated.enabled ? "enabled" : "disabled"}.`,
                                        ),
                                      onError: (err) => setErrorMsg(err.message),
                                    },
                                  );
                                }}
                              >
                                {pool.enabled ? "Disable" : "Enable"}
                              </Button>
                              <Button
                                type="button"
                                size="sm"
                                variant="secondary"
                                disabled={testingPoolId === pool.id}
                                onClick={() => handleTestPool(pool.id)}
                              >
                                <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                                <span>{testingPoolId === pool.id ? "Testing..." : "Test"}</span>
                              </Button>
                              <Button
                                type="button"
                                size="sm"
                                variant="destructive"
                                onClick={() => {
                                  setErrorMsg(null);
                                  setNotice(null);
                                  deletePool.mutate(pool.id, {
                                    onSuccess: () =>
                                      setNotice(`Deleted proxy pool "${pool.name}".`),
                                    onError: (err) => setErrorMsg(err.message),
                                  });
                                }}
                                aria-label={`Delete ${pool.name}`}
                              >
                                <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
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
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Provider Proxy Bindings</CardTitle>
          <CardDescription>
            Bind any upstream provider to an active proxy pool. Unbound providers dial directly.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {bindableProviders.map((provider) => {
              const boundPoolId = bindingByProvider.get(provider.id) ?? "";
              return (
                <div
                  key={provider.id}
                  className="flex flex-col gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3"
                >
                  <div className="flex items-center justify-between gap-2">
                    <span className="text-sm font-medium text-[var(--text-primary)]">
                      {provider.label}
                    </span>
                    <span className="font-mono text-[11px] text-[var(--text-muted)]">
                      {provider.id}
                    </span>
                  </div>
                  <select
                    aria-label={`Proxy pool for ${provider.label}`}
                    value={boundPoolId}
                    onChange={(e) => handleBindingChange(provider.id, e.target.value)}
                    className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 py-1.5 text-xs text-[var(--text-primary)] focus:border-[var(--accent-primary)] focus:outline-none"
                  >
                    <option value="">Direct (No Proxy)</option>
                    {pools.map((pool) => (
                      <option key={pool.id} value={pool.id}>
                        {pool.name} ({pool.scheme})
                      </option>
                    ))}
                  </select>
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
