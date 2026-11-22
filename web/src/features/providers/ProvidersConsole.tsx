"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Activity,
  AlertCircle,
  Check,
  CheckCircle2,
  Copy,
  ExternalLink,
  Eye,
  EyeOff,
  KeyRound,
  Layers,
  Plus,
  Power,
  RefreshCw,
  Server,
  Settings2,
  ShieldCheck,
  Trash2,
  X,
} from "lucide-react";
import Link from "next/link";
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
  useCreateProviderConnectionMutation,
  useCreateProviderNodeMutation,
  useDeleteProviderConnectionMutation,
  useDeleteProviderNodeMutation,
  useExchangeProviderOAuthMutation,
  useProviderConnections,
  useProvidersCatalog,
  useStartProviderOAuthMutation,
  useTestProviderConnectionMutation,
  useUpdateProviderConnectionMutation,
} from "./hooks";
import {
  type CatalogAdapterDescriptor,
  type CreateConnectionFormInput,
  type CreateConnectionFormOutput,
  type CreateNodeFormInput,
  type CreateNodeFormOutput,
  createConnectionSchema,
  createNodeSchema,
  DEFAULT_CREATE_CONNECTION_VALUES,
  DEFAULT_CREATE_NODE_VALUES,
  DEFAULT_OAUTH_EXCHANGE_VALUES,
  DEFAULT_OAUTH_START_VALUES,
  formatCommaSeparatedList,
  formatHeadersText,
  OAUTH_PROVIDER_IDS,
  type OAuthExchangeFormInput,
  type OAuthExchangeFormOutput,
  type OAuthStartFormInput,
  type OAuthStartFormOutput,
  type OAuthStartResponse,
  oauthExchangeSchema,
  oauthStartSchema,
  PROVIDER_PROTOCOLS,
  type ProviderConnectionSummary,
  type ProviderNode,
  type TestProviderConnectionResponse,
  toCreateConnectionRequest,
  toCreateNodeRequest,
  toOAuthExchangeRequest,
  toOAuthStartRequest,
} from "./schemas";

interface ActiveOAuthSession {
  provider: string;
  session: OAuthStartResponse;
}

const SELECT_CLASS_NAME =
  "flex h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-sm text-[var(--text-primary)] transition-colors focus:border-[#E56A4A] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#E56A4A]";

const TEXTAREA_CLASS_NAME =
  "flex min-h-[82px] w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-xs text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[#E56A4A] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#E56A4A]";

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

function renderTestBadge(status: string, liveResult?: TestProviderConnectionResponse) {
  if (liveResult) {
    return liveResult.ok ? (
      <Badge variant="success" className="tabular-nums">
        OK {liveResult.status_code} · {liveResult.latency_ms}ms
      </Badge>
    ) : (
      <Badge variant="destructive" className="tabular-nums">
        FAIL {liveResult.status_code || "ERR"} · {liveResult.latency_ms}ms
      </Badge>
    );
  }
  if (status === "ok") {
    return <Badge variant="success">Healthy</Badge>;
  }
  if (status === "error") {
    return <Badge variant="destructive">Failed</Badge>;
  }
  return <Badge variant="outline">Untested</Badge>;
}

export function ProvidersConsole() {
  const catalogQuery = useProvidersCatalog();
  const connectionsQuery = useProviderConnections();

  const createConnMutation = useCreateProviderConnectionMutation();
  const updateConnMutation = useUpdateProviderConnectionMutation();
  const deleteConnMutation = useDeleteProviderConnectionMutation();
  const testConnMutation = useTestProviderConnectionMutation();
  const createNodeMutation = useCreateProviderNodeMutation();
  const deleteNodeMutation = useDeleteProviderNodeMutation();
  const startOAuthMutation = useStartProviderOAuthMutation();
  const exchangeOAuthMutation = useExchangeProviderOAuthMutation();

  const [showApiKey, setShowApiKey] = useState(false);
  const [connFeedback, setConnFeedback] = useState<string | null>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [busyConnectionId, setBusyConnectionId] = useState<string | null>(null);
  const [confirmDeleteConnId, setConfirmDeleteConnId] = useState<string | null>(null);
  const [liveTestResults, setLiveTestResults] = useState<
    Record<string, TestProviderConnectionResponse>
  >({});

  const [nodeFeedback, setNodeFeedback] = useState<string | null>(null);
  const [nodeRowError, setNodeRowError] = useState<string | null>(null);
  const [busyNodeId, setBusyNodeId] = useState<string | null>(null);
  const [confirmDeleteNodeId, setConfirmDeleteNodeId] = useState<string | null>(null);

  const [activeOAuth, setActiveOAuth] = useState<ActiveOAuthSession | null>(null);
  const [copiedAuthUrl, setCopiedAuthUrl] = useState(false);
  const [copyAuthUrlError, setCopyAuthUrlError] = useState<string | null>(null);
  const [completedOAuthConn, setCompletedOAuthConn] = useState<ProviderConnectionSummary | null>(
    null,
  );

  const createConnForm = useForm<CreateConnectionFormInput, unknown, CreateConnectionFormOutput>({
    resolver: zodResolver(createConnectionSchema),
    defaultValues: DEFAULT_CREATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const createNodeForm = useForm<CreateNodeFormInput, unknown, CreateNodeFormOutput>({
    resolver: zodResolver(createNodeSchema),
    defaultValues: DEFAULT_CREATE_NODE_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const oauthStartForm = useForm<OAuthStartFormInput, unknown, OAuthStartFormOutput>({
    resolver: zodResolver(oauthStartSchema),
    defaultValues: DEFAULT_OAUTH_START_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const oauthExchangeForm = useForm<OAuthExchangeFormInput, unknown, OAuthExchangeFormOutput>({
    resolver: zodResolver(oauthExchangeSchema),
    defaultValues: DEFAULT_OAUTH_EXCHANGE_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const adapters: CatalogAdapterDescriptor[] = catalogQuery.data?.adapters ?? [];
  const customNodes: ProviderNode[] = catalogQuery.data?.custom_nodes ?? [];
  const connections: ProviderConnectionSummary[] = connectionsQuery.data ?? [];

  const apiKeyAdapters = adapters.filter((item) => item.auth_type === "api_key");
  const oauthAdapters = adapters.filter((item) => item.auth_type === "oauth");

  const enabledCount = connections.filter((conn) => conn.enabled).length;
  const oauthCount = connections.filter((conn) => conn.auth_type === "oauth").length;
  const isRefreshing = catalogQuery.isFetching || connectionsQuery.isFetching;

  const selectedOAuthProvider = oauthStartForm.watch("provider");
  const selectedOAuthDescriptor = oauthAdapters.find((item) => item.id === selectedOAuthProvider);

  const handleRefreshAll = () => {
    void Promise.all([catalogQuery.refetch(), connectionsQuery.refetch()]);
  };

  const applyProviderDefaults = (providerId: string) => {
    const adapter = apiKeyAdapters.find((item) => item.id === providerId);
    if (adapter) {
      createConnForm.setValue("base_url", adapter.default_base_url, { shouldDirty: false });
      createConnForm.setValue("models", formatCommaSeparatedList(adapter.default_models), {
        shouldDirty: false,
      });
      createConnForm.setValue("headers", "", { shouldDirty: false });
      return;
    }
    const node = customNodes.find((item) => item.id === providerId);
    if (node) {
      createConnForm.setValue("base_url", node.base_url, { shouldDirty: false });
      createConnForm.setValue("models", formatCommaSeparatedList(node.models), {
        shouldDirty: false,
      });
      createConnForm.setValue("headers", formatHeadersText(node.headers), { shouldDirty: false });
    }
  };

  const onCreateConnSubmit = (values: CreateConnectionFormOutput) => {
    setConnFeedback(null);
    createConnMutation.reset();

    createConnMutation.mutate(toCreateConnectionRequest(values), {
      onSuccess: (created) => {
        setConnFeedback(`Created connection "${created.name}" (${created.id}).`);
        createConnForm.reset({
          ...DEFAULT_CREATE_CONNECTION_VALUES,
          provider_id: created.provider_id,
          base_url: created.base_url,
          models: formatCommaSeparatedList(created.models),
        });
      },
    });
  };

  const handleToggleConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) {
      return;
    }
    setRowError(null);
    setConfirmDeleteConnId(null);
    setBusyConnectionId(conn.id);

    updateConnMutation.mutate(
      { id: conn.id, body: { enabled: !conn.enabled } },
      {
        onError: (err) => setRowError(err.message),
        onSettled: () => setBusyConnectionId(null),
      },
    );
  };

  const handleTestConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) {
      return;
    }
    setRowError(null);
    setConfirmDeleteConnId(null);
    setBusyConnectionId(conn.id);

    testConnMutation.mutate(conn.id, {
      onSuccess: (result) => {
        setLiveTestResults((prev) => ({ ...prev, [conn.id]: result }));
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const handleConfirmDeleteConnection = (conn: ProviderConnectionSummary) => {
    if (busyConnectionId !== null) {
      return;
    }
    setRowError(null);
    setBusyConnectionId(conn.id);

    deleteConnMutation.mutate(conn.id, {
      onSuccess: () => {
        setConfirmDeleteConnId(null);
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyConnectionId(null),
    });
  };

  const onOAuthStartSubmit = (values: OAuthStartFormOutput) => {
    setCompletedOAuthConn(null);
    setCopiedAuthUrl(false);
    setCopyAuthUrlError(null);
    startOAuthMutation.reset();
    exchangeOAuthMutation.reset();

    startOAuthMutation.mutate(
      {
        provider: values.provider,
        body: toOAuthStartRequest(values),
      },
      {
        onSuccess: (session) => {
          setActiveOAuth({ provider: values.provider, session });
          oauthExchangeForm.setValue("state", session.state, { shouldDirty: false });
          if (values.connection_name && !oauthExchangeForm.getValues("connection_name").trim()) {
            oauthExchangeForm.setValue("connection_name", values.connection_name, {
              shouldDirty: false,
            });
          }
        },
      },
    );
  };

  const handleCopyAuthorizeUrl = async () => {
    if (!activeOAuth) {
      return;
    }
    setCopyAuthUrlError(null);
    try {
      await navigator.clipboard.writeText(activeOAuth.session.authorize_url);
      setCopiedAuthUrl(true);
      window.setTimeout(() => setCopiedAuthUrl(false), 2500);
    } catch {
      setCopyAuthUrlError("Clipboard access was denied. Copy the authorization URL manually.");
    }
  };

  const onOAuthExchangeSubmit = (values: OAuthExchangeFormOutput) => {
    const provider = activeOAuth?.provider ?? selectedOAuthProvider;
    setCompletedOAuthConn(null);
    exchangeOAuthMutation.reset();

    exchangeOAuthMutation.mutate(
      {
        provider,
        body: toOAuthExchangeRequest(values),
      },
      {
        onSuccess: (conn) => {
          setCompletedOAuthConn(conn);
          setActiveOAuth(null);
          setCopiedAuthUrl(false);
          setCopyAuthUrlError(null);
          oauthExchangeForm.reset(DEFAULT_OAUTH_EXCHANGE_VALUES);
          oauthStartForm.reset({
            ...DEFAULT_OAUTH_START_VALUES,
            provider,
          });
        },
      },
    );
  };

  const onCreateNodeSubmit = (values: CreateNodeFormOutput) => {
    setNodeFeedback(null);
    createNodeMutation.reset();

    createNodeMutation.mutate(toCreateNodeRequest(values), {
      onSuccess: (created) => {
        setNodeFeedback(`Registered custom provider node "${created.name}" (${created.id}).`);
        createNodeForm.reset(DEFAULT_CREATE_NODE_VALUES);
      },
    });
  };

  const handleConfirmDeleteNode = (node: ProviderNode) => {
    if (busyNodeId !== null) {
      return;
    }
    setNodeRowError(null);
    setBusyNodeId(node.id);

    deleteNodeMutation.mutate(node.id, {
      onSuccess: () => {
        setConfirmDeleteNodeId(null);
      },
      onError: (err) => setNodeRowError(err.message),
      onSettled: () => setBusyNodeId(null),
    });
  };

  const anyConnBusy = busyConnectionId !== null;
  const anyNodeBusy = busyNodeId !== null;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Upstream Credentials &amp; Adapters
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Providers Control Studio
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Manage vaulted API-key accounts, headless-safe OAuth PKCE sessions, and custom
            compatible provider nodes.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <Badge variant="outline" className="tabular-nums">
            {connections.length} Connections
          </Badge>
          <Badge variant="success" className="tabular-nums">
            {enabledCount} Enabled
          </Badge>
          {oauthCount > 0 ? (
            <Badge variant="default" className="tabular-nums">
              {oauthCount} OAuth
            </Badge>
          ) : null}
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefreshAll}
            disabled={isRefreshing}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {rowError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="flex-1">
            <AlertTitle>Connection action failed</AlertTitle>
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
            <Server className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            <CardTitle>Connected Upstream Accounts</CardTitle>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            GET /api/providers/connections
          </span>
        </CardHeader>

        {connectionsQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load provider connections</AlertTitle>
                <AlertDescription className="font-mono">
                  {connectionsQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : connectionsQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading provider connections...
          </CardContent>
        ) : connections.length === 0 ? (
          <CardContent className="py-10 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              No upstream provider connections configured yet
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              Create an API-key connection or complete an OAuth PKCE flow below to route upstream
              requests.
            </p>
          </CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <th scope="col" className="px-5 py-3">
                    Connection &amp; Provider
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Protocol &amp; Endpoint
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Priority &amp; Models
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Last Health Test
                  </th>
                  <th scope="col" className="px-5 py-3 text-right">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {connections.map((conn) => {
                  const isRowBusy = busyConnectionId === conn.id;
                  const confirmingDelete = confirmDeleteConnId === conn.id;
                  const liveResult = liveTestResults[conn.id];
                  const visibleModels = conn.models.slice(0, 3);
                  const remainingModels = conn.models.length - visibleModels.length;

                  return (
                    <tr
                      key={conn.id}
                      className="align-top transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3.5">
                        <div className="flex flex-wrap items-center gap-2">
                          <Link
                            href={`/dashboard/providers/detail/?id=${encodeURIComponent(conn.id)}`}
                            className="font-medium text-[var(--text-primary)] hover:text-[var(--accent-text)]"
                          >
                            {conn.name}
                          </Link>
                          <Badge variant={conn.enabled ? "success" : "warning"}>
                            {conn.enabled ? "Enabled" : "Disabled"}
                          </Badge>
                        </div>
                        <div className="mt-1.5 flex flex-wrap items-center gap-2 font-mono text-xs text-[var(--text-muted)]">
                          <span>{conn.provider_id}</span>
                          <span>·</span>
                          <span>{conn.id}</span>
                        </div>
                      </td>

                      <td className="px-4 py-3.5">
                        <div className="flex flex-wrap items-center gap-1.5">
                          <Badge variant="outline">{conn.protocol}</Badge>
                          <Badge variant={conn.auth_type === "oauth" ? "default" : "secondary"}>
                            {conn.auth_type === "oauth" ? "OAuth" : "API Key"}
                          </Badge>
                        </div>
                        <p className="mt-1.5 max-w-xs truncate font-mono text-xs text-[var(--text-secondary)]">
                          {conn.base_url}
                        </p>
                      </td>

                      <td className="px-4 py-3.5">
                        <Badge variant="outline" className="tabular-nums">
                          Priority {conn.priority}
                        </Badge>
                        <div className="mt-1.5 flex max-w-xs flex-wrap gap-1.5">
                          {conn.models.length === 0 ? (
                            <span className="font-mono text-xs text-[var(--text-muted)]">
                              No model filter
                            </span>
                          ) : (
                            <>
                              {visibleModels.map((model) => (
                                <Badge key={model} variant="secondary">
                                  {model}
                                </Badge>
                              ))}
                              {remainingModels > 0 ? (
                                <Badge variant="outline" className="tabular-nums">
                                  +{remainingModels} more
                                </Badge>
                              ) : null}
                            </>
                          )}
                        </div>
                      </td>

                      <td className="px-4 py-3.5">
                        <div className="flex flex-wrap items-center gap-2">
                          {renderTestBadge(conn.last_test_status, liveResult)}
                        </div>
                        <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
                          {formatTimestamp(conn.last_tested_at)}
                        </p>
                        {(liveResult && !liveResult.ok && liveResult.error) ||
                        (!liveResult && conn.last_test_error) ? (
                          <p className="mt-1 max-w-xs truncate font-mono text-xs text-[var(--status-error)]">
                            {liveResult?.error || conn.last_test_error}
                          </p>
                        ) : null}
                      </td>

                      <td className="px-5 py-3.5 text-right">
                        {confirmingDelete ? (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => handleConfirmDeleteConnection(conn)}
                              disabled={anyConnBusy}
                            >
                              {isRowBusy ? "Deleting..." : "Confirm Delete"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => setConfirmDeleteConnId(null)}
                              disabled={anyConnBusy}
                            >
                              Cancel
                            </Button>
                          </div>
                        ) : (
                          <div className="flex flex-wrap items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="outline"
                              size="xs"
                              onClick={() => handleTestConnection(conn)}
                              disabled={anyConnBusy}
                            >
                              <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>{isRowBusy ? "Testing..." : "Quick Test"}</span>
                            </Button>
                            <Button asChild variant="outline" size="xs">
                              <Link
                                href={`/dashboard/providers/detail/?id=${encodeURIComponent(conn.id)}`}
                              >
                                <Settings2 className="h-3.5 w-3.5" aria-hidden="true" />
                                <span>Detail</span>
                              </Link>
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => handleToggleConnection(conn)}
                              disabled={anyConnBusy}
                            >
                              <Power className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>{conn.enabled ? "Disable" : "Enable"}</span>
                            </Button>
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => setConfirmDeleteConnId(conn.id)}
                              disabled={anyConnBusy}
                              aria-label={`Delete connection ${conn.name}`}
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

      <EmberRule />

      <section className="grid grid-cols-1 gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <KeyRound className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Create API-Key Connection</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Encrypt an upstream API key inside the AES-256-GCM credential vault.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              POST /api/providers/connections
            </span>
          </CardHeader>
          <CardContent>
            <form
              onSubmit={createConnForm.handleSubmit(onCreateConnSubmit)}
              noValidate
              className="space-y-4"
            >
              {createConnMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to create connection</AlertTitle>
                    <AlertDescription className="font-mono">
                      {createConnMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {connFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Connection Vaulted</AlertTitle>
                    <AlertDescription>{connFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <Label htmlFor="conn-provider">Provider Adapter / Node</Label>
                  <select
                    id="conn-provider"
                    className={`${SELECT_CLASS_NAME} mt-1.5`}
                    {...createConnForm.register("provider_id", {
                      onChange: (event: React.ChangeEvent<HTMLSelectElement>) => {
                        applyProviderDefaults(event.target.value);
                      },
                    })}
                  >
                    {apiKeyAdapters.length > 0 ? (
                      <optgroup label="Built-in API-Key Adapters">
                        {apiKeyAdapters.map((adapter) => (
                          <option key={adapter.id} value={adapter.id}>
                            {adapter.display_name} ({adapter.protocol})
                          </option>
                        ))}
                      </optgroup>
                    ) : (
                      <option value="openai">OpenAI (openai-chat)</option>
                    )}
                    {customNodes.length > 0 ? (
                      <optgroup label="Custom Compatible Nodes">
                        {customNodes.map((node) => (
                          <option key={node.id} value={node.id}>
                            {node.name} ({node.id})
                          </option>
                        ))}
                      </optgroup>
                    ) : null}
                  </select>
                  {createConnForm.formState.errors.provider_id ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createConnForm.formState.errors.provider_id.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="conn-name">Connection Name</Label>
                  <Input
                    id="conn-name"
                    type="text"
                    placeholder="e.g. primary-openai"
                    aria-invalid={Boolean(createConnForm.formState.errors.name)}
                    className="mt-1.5"
                    {...createConnForm.register("name")}
                  />
                  {createConnForm.formState.errors.name ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createConnForm.formState.errors.name.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="conn-base-url">Base URL</Label>
                  <Input
                    id="conn-base-url"
                    type="url"
                    placeholder="https://api.openai.com/v1"
                    aria-invalid={Boolean(createConnForm.formState.errors.base_url)}
                    className="mt-1.5"
                    {...createConnForm.register("base_url")}
                  />
                  {createConnForm.formState.errors.base_url ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createConnForm.formState.errors.base_url.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="conn-priority">Routing Priority (lower runs first)</Label>
                  <Input
                    id="conn-priority"
                    type="number"
                    min={0}
                    max={100000}
                    step={1}
                    aria-invalid={Boolean(createConnForm.formState.errors.priority)}
                    className="mt-1.5 tabular-nums"
                    {...createConnForm.register("priority")}
                  />
                  {createConnForm.formState.errors.priority ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createConnForm.formState.errors.priority.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div>
                <Label htmlFor="conn-api-key">Upstream API Key</Label>
                <div className="relative mt-1.5">
                  <Input
                    id="conn-api-key"
                    type={showApiKey ? "text" : "password"}
                    autoComplete="off"
                    placeholder="sk-..."
                    aria-invalid={Boolean(createConnForm.formState.errors.api_key)}
                    className="pr-10"
                    {...createConnForm.register("api_key")}
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
                {createConnForm.formState.errors.api_key ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createConnForm.formState.errors.api_key.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="conn-models">Allowed / Advertised Models (comma-separated)</Label>
                <Input
                  id="conn-models"
                  type="text"
                  placeholder="gpt-4o, gpt-4o-mini, o3-mini"
                  aria-invalid={Boolean(createConnForm.formState.errors.models)}
                  className="mt-1.5"
                  {...createConnForm.register("models")}
                />
                {createConnForm.formState.errors.models ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createConnForm.formState.errors.models.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="conn-headers">
                  Custom Headers (optional, one Key: Value per line)
                </Label>
                <textarea
                  id="conn-headers"
                  placeholder={"HTTP-Referer: https://router.local\nX-Title: omeg4mon-router"}
                  aria-invalid={Boolean(createConnForm.formState.errors.headers)}
                  className={`${TEXTAREA_CLASS_NAME} mt-1.5`}
                  {...createConnForm.register("headers")}
                />
                {createConnForm.formState.errors.headers ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createConnForm.formState.errors.headers.message}
                  </p>
                ) : null}
              </div>

              <label
                htmlFor="conn-enabled"
                className="flex cursor-pointer items-center justify-between rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-2.5"
              >
                <span className="text-sm text-[var(--text-primary)]">
                  Enable connection for gateway routing immediately
                </span>
                <input
                  id="conn-enabled"
                  type="checkbox"
                  className="h-4 w-4 accent-[#E56A4A]"
                  {...createConnForm.register("enabled")}
                />
              </label>

              <div className="flex items-center justify-end gap-3 pt-1">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    createConnMutation.reset();
                    setConnFeedback(null);
                    createConnForm.reset(DEFAULT_CREATE_CONNECTION_VALUES);
                  }}
                >
                  Reset
                </Button>
                <Button type="submit" size="sm" disabled={createConnMutation.isPending}>
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {createConnMutation.isPending ? "Saving..." : "Create API-Key Connection"}
                  </span>
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Headless-Safe OAuth PKCE Wizard</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Start PKCE S256 authorization, open the URL in any browser, and exchange either the
                callback URL or authorization code.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">OAuth 2.0 PKCE</span>
          </CardHeader>
          <CardContent className="space-y-5">
            <form
              onSubmit={oauthStartForm.handleSubmit(onOAuthStartSubmit)}
              noValidate
              className="space-y-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4"
            >
              <div className="flex flex-wrap items-center justify-between gap-2">
                <span className="font-mono text-xs uppercase tracking-wider text-[var(--accent-text)]">
                  Step 1 · Initialize Authorization Session
                </span>
                {selectedOAuthDescriptor?.oauth ? (
                  <Badge
                    variant={
                      selectedOAuthDescriptor.oauth.client_id_configured ? "success" : "warning"
                    }
                  >
                    {selectedOAuthDescriptor.oauth.client_id_configured
                      ? "Server Client ID Ready"
                      : "Client ID Required"}
                  </Badge>
                ) : null}
              </div>

              {startOAuthMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>OAuth start failed</AlertTitle>
                    <AlertDescription className="font-mono">
                      {startOAuthMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
                <div>
                  <Label htmlFor="oauth-provider">OAuth Provider</Label>
                  <select
                    id="oauth-provider"
                    className={`${SELECT_CLASS_NAME} mt-1.5`}
                    {...oauthStartForm.register("provider")}
                  >
                    {oauthAdapters.length > 0
                      ? oauthAdapters.map((adapter) => (
                          <option key={adapter.id} value={adapter.id}>
                            {adapter.display_name} ({adapter.id})
                          </option>
                        ))
                      : OAUTH_PROVIDER_IDS.map((providerId) => (
                          <option key={providerId} value={providerId}>
                            {providerId}
                          </option>
                        ))}
                  </select>
                </div>

                <div>
                  <Label htmlFor="oauth-start-name">Connection Name (optional)</Label>
                  <Input
                    id="oauth-start-name"
                    type="text"
                    placeholder={selectedOAuthDescriptor?.display_name || "e.g. Claude Max"}
                    aria-invalid={Boolean(oauthStartForm.formState.errors.connection_name)}
                    className="mt-1.5"
                    {...oauthStartForm.register("connection_name")}
                  />
                </div>

                <div>
                  <Label htmlFor="oauth-client-id">Custom Client ID (optional)</Label>
                  <Input
                    id="oauth-client-id"
                    type="text"
                    placeholder="Overrides server env if provided"
                    className="mt-1.5"
                    {...oauthStartForm.register("client_id")}
                  />
                </div>

                <div>
                  <Label htmlFor="oauth-client-secret">Custom Client Secret (optional)</Label>
                  <Input
                    id="oauth-client-secret"
                    type="password"
                    autoComplete="off"
                    placeholder="Optional confidential client secret"
                    className="mt-1.5"
                    {...oauthStartForm.register("client_secret")}
                  />
                </div>

                <div>
                  <Label htmlFor="oauth-redirect-uri">Redirect URI (optional)</Label>
                  <Input
                    id="oauth-redirect-uri"
                    type="url"
                    placeholder={
                      selectedOAuthDescriptor?.oauth?.default_redirect_uri ||
                      "http://localhost:20128/api/oauth/callback"
                    }
                    aria-invalid={Boolean(oauthStartForm.formState.errors.redirect_uri)}
                    className="mt-1.5"
                    {...oauthStartForm.register("redirect_uri")}
                  />
                  {oauthStartForm.formState.errors.redirect_uri ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {oauthStartForm.formState.errors.redirect_uri.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="oauth-base-url">Upstream Base URL (optional)</Label>
                  <Input
                    id="oauth-base-url"
                    type="url"
                    placeholder={selectedOAuthDescriptor?.default_base_url || "Adapter default"}
                    aria-invalid={Boolean(oauthStartForm.formState.errors.base_url)}
                    className="mt-1.5"
                    {...oauthStartForm.register("base_url")}
                  />
                  {oauthStartForm.formState.errors.base_url ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {oauthStartForm.formState.errors.base_url.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div>
                <Label htmlFor="oauth-models">Model Override List (optional comma-separated)</Label>
                <Input
                  id="oauth-models"
                  type="text"
                  placeholder={
                    selectedOAuthDescriptor?.default_models.join(", ") || "Uses adapter defaults"
                  }
                  aria-invalid={Boolean(oauthStartForm.formState.errors.models)}
                  className="mt-1.5"
                  {...oauthStartForm.register("models")}
                />
              </div>

              <div className="flex justify-end">
                <Button type="submit" size="sm" disabled={startOAuthMutation.isPending}>
                  <ShieldCheck className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {startOAuthMutation.isPending
                      ? "Generating PKCE URL..."
                      : "Start OAuth Session"}
                  </span>
                </Button>
              </div>
            </form>

            {activeOAuth ? (
              <div className="space-y-3 rounded-md border border-[#E56A4A] bg-[var(--bg-surface-raised)] p-4">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="font-mono text-xs font-medium uppercase tracking-wider text-[var(--accent-text)]">
                    Authorize URL Ready ({activeOAuth.provider})
                  </span>
                  <Badge variant="outline" className="tabular-nums">
                    Expires {formatTimestamp(activeOAuth.session.expires_at)}
                  </Badge>
                </div>

                <code className="block max-h-28 overflow-x-auto break-all rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-xs text-[var(--text-primary)] select-all">
                  {activeOAuth.session.authorize_url}
                </code>

                <div className="flex flex-wrap items-center gap-2">
                  <Button type="button" size="xs" onClick={() => void handleCopyAuthorizeUrl()}>
                    {copiedAuthUrl ? (
                      <>
                        <Check className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>Copied URL</span>
                      </>
                    ) : (
                      <>
                        <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>Copy URL</span>
                      </>
                    )}
                  </Button>
                  <Button asChild variant="outline" size="xs">
                    <a
                      href={activeOAuth.session.authorize_url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>Open in New Tab</span>
                    </a>
                  </Button>
                </div>

                {copyAuthUrlError ? (
                  <p className="font-mono text-xs text-[var(--status-warn)]">{copyAuthUrlError}</p>
                ) : null}

                <div className="grid grid-cols-1 gap-1.5 font-mono text-xs text-[var(--text-muted)] sm:grid-cols-2">
                  <div className="truncate">State: {activeOAuth.session.state}</div>
                  <div className="truncate">Redirect: {activeOAuth.session.redirect_uri}</div>
                </div>
              </div>
            ) : null}

            <form
              onSubmit={oauthExchangeForm.handleSubmit(onOAuthExchangeSubmit)}
              noValidate
              className="space-y-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4"
            >
              <div className="flex flex-wrap items-center justify-between gap-2">
                <span className="font-mono text-xs uppercase tracking-wider text-[var(--accent-text)]">
                  Step 2 · Exchange Callback URL or Code
                </span>
                <span className="font-mono text-xs text-[var(--text-muted)]">
                  POST /api/providers/oauth/{activeOAuth?.provider ?? selectedOAuthProvider}
                  /exchange
                </span>
              </div>

              {exchangeOAuthMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>OAuth token exchange failed</AlertTitle>
                    <AlertDescription className="font-mono">
                      {exchangeOAuthMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {completedOAuthConn ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div className="flex flex-1 flex-wrap items-center justify-between gap-2">
                    <div>
                      <AlertTitle>OAuth Connection Vaulted</AlertTitle>
                      <AlertDescription>
                        Saved "{completedOAuthConn.name}" ({completedOAuthConn.id}).
                      </AlertDescription>
                    </div>
                    <Button asChild variant="outline" size="xs">
                      <Link
                        href={`/dashboard/providers/detail/?id=${encodeURIComponent(completedOAuthConn.id)}`}
                      >
                        Open Detail
                      </Link>
                    </Button>
                  </div>
                </Alert>
              ) : null}

              <div>
                <Label htmlFor="oauth-callback-url">Full Callback URL (recommended headless)</Label>
                <Input
                  id="oauth-callback-url"
                  type="url"
                  placeholder="http://localhost:20128/api/oauth/callback?code=...&state=..."
                  aria-invalid={Boolean(oauthExchangeForm.formState.errors.callback_url)}
                  className="mt-1.5"
                  {...oauthExchangeForm.register("callback_url")}
                />
                {oauthExchangeForm.formState.errors.callback_url ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {oauthExchangeForm.formState.errors.callback_url.message}
                  </p>
                ) : null}
              </div>

              <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
                <div>
                  <Label htmlFor="oauth-code">Or Manual Authorization Code</Label>
                  <Input
                    id="oauth-code"
                    type="text"
                    placeholder="Authorization code"
                    aria-invalid={Boolean(oauthExchangeForm.formState.errors.code)}
                    className="mt-1.5"
                    {...oauthExchangeForm.register("code")}
                  />
                  {oauthExchangeForm.formState.errors.code ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {oauthExchangeForm.formState.errors.code.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="oauth-state">State (auto-filled from Step 1)</Label>
                  <Input
                    id="oauth-state"
                    type="text"
                    placeholder="OAuth session state"
                    aria-invalid={Boolean(oauthExchangeForm.formState.errors.state)}
                    className="mt-1.5"
                    {...oauthExchangeForm.register("state")}
                  />
                  {oauthExchangeForm.formState.errors.state ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {oauthExchangeForm.formState.errors.state.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="oauth-exchange-name">Connection Name Override</Label>
                  <Input
                    id="oauth-exchange-name"
                    type="text"
                    placeholder="Optional override"
                    className="mt-1.5"
                    {...oauthExchangeForm.register("connection_name")}
                  />
                </div>

                <div>
                  <Label htmlFor="oauth-exchange-priority">Routing Priority</Label>
                  <Input
                    id="oauth-exchange-priority"
                    type="number"
                    min={0}
                    max={100000}
                    step={1}
                    aria-invalid={Boolean(oauthExchangeForm.formState.errors.priority)}
                    className="mt-1.5 tabular-nums"
                    {...oauthExchangeForm.register("priority")}
                  />
                  {oauthExchangeForm.formState.errors.priority ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {oauthExchangeForm.formState.errors.priority.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div className="flex flex-wrap items-center justify-between gap-3">
                <label
                  htmlFor="oauth-exchange-enabled"
                  className="flex cursor-pointer items-center gap-2 text-xs text-[var(--text-secondary)]"
                >
                  <input
                    id="oauth-exchange-enabled"
                    type="checkbox"
                    className="h-4 w-4 accent-[#E56A4A]"
                    {...oauthExchangeForm.register("enabled")}
                  />
                  <span>Enable connection after exchange</span>
                </label>

                <Button type="submit" size="sm" disabled={exchangeOAuthMutation.isPending}>
                  <Check className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {exchangeOAuthMutation.isPending
                      ? "Exchanging..."
                      : "Exchange & Save Connection"}
                  </span>
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </section>

      <EmberRule />

      <section className="grid grid-cols-1 gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Layers className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Custom Compatible Provider Nodes</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Register reusable OpenAI, Anthropic, or Gemini compatible base endpoints (e.g. vLLM,
                LiteLLM, Ollama).
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              POST /api/providers/nodes
            </span>
          </CardHeader>
          <CardContent className="space-y-5">
            <form
              onSubmit={createNodeForm.handleSubmit(onCreateNodeSubmit)}
              noValidate
              className="space-y-4"
            >
              {createNodeMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to create provider node</AlertTitle>
                    <AlertDescription className="font-mono">
                      {createNodeMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {nodeFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Custom Node Registered</AlertTitle>
                    <AlertDescription>{nodeFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <Label htmlFor="node-id">Slug ID (optional)</Label>
                  <Input
                    id="node-id"
                    type="text"
                    placeholder="e.g. local-vllm (auto if blank)"
                    aria-invalid={Boolean(createNodeForm.formState.errors.id)}
                    className="mt-1.5"
                    {...createNodeForm.register("id")}
                  />
                  {createNodeForm.formState.errors.id ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createNodeForm.formState.errors.id.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="node-name">Node Name</Label>
                  <Input
                    id="node-name"
                    type="text"
                    placeholder="e.g. Local vLLM Cluster"
                    aria-invalid={Boolean(createNodeForm.formState.errors.name)}
                    className="mt-1.5"
                    {...createNodeForm.register("name")}
                  />
                  {createNodeForm.formState.errors.name ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createNodeForm.formState.errors.name.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="node-protocol">Wire Protocol</Label>
                  <select
                    id="node-protocol"
                    className={`${SELECT_CLASS_NAME} mt-1.5`}
                    {...createNodeForm.register("protocol")}
                  >
                    {PROVIDER_PROTOCOLS.map((protocol) => (
                      <option key={protocol} value={protocol}>
                        {protocol}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <Label htmlFor="node-base-url">Base URL</Label>
                  <Input
                    id="node-base-url"
                    type="url"
                    placeholder="http://127.0.0.1:8000/v1"
                    aria-invalid={Boolean(createNodeForm.formState.errors.base_url)}
                    className="mt-1.5"
                    {...createNodeForm.register("base_url")}
                  />
                  {createNodeForm.formState.errors.base_url ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {createNodeForm.formState.errors.base_url.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div>
                <Label htmlFor="node-models">Default Models (comma-separated)</Label>
                <Input
                  id="node-models"
                  type="text"
                  placeholder="qwen-2.5-72b, llama-3.3-70b"
                  aria-invalid={Boolean(createNodeForm.formState.errors.models)}
                  className="mt-1.5"
                  {...createNodeForm.register("models")}
                />
                {createNodeForm.formState.errors.models ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createNodeForm.formState.errors.models.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="node-headers">
                  Default Headers (optional, Key: Value per line)
                </Label>
                <textarea
                  id="node-headers"
                  placeholder="X-Cluster: gpu-01"
                  aria-invalid={Boolean(createNodeForm.formState.errors.headers)}
                  className={`${TEXTAREA_CLASS_NAME} mt-1.5`}
                  {...createNodeForm.register("headers")}
                />
                {createNodeForm.formState.errors.headers ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createNodeForm.formState.errors.headers.message}
                  </p>
                ) : null}
              </div>

              <div className="flex justify-end gap-3">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    createNodeMutation.reset();
                    setNodeFeedback(null);
                    createNodeForm.reset(DEFAULT_CREATE_NODE_VALUES);
                  }}
                >
                  Reset
                </Button>
                <Button type="submit" size="sm" disabled={createNodeMutation.isPending}>
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{createNodeMutation.isPending ? "Registering..." : "Register Node"}</span>
                </Button>
              </div>
            </form>

            {nodeRowError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Node operation failed</AlertTitle>
                  <AlertDescription className="font-mono">{nodeRowError}</AlertDescription>
                </div>
              </Alert>
            ) : null}

            <div className="space-y-2.5 border-t border-[var(--border-subtle)] pt-4">
              <div className="flex items-center justify-between">
                <span className="font-mono text-xs uppercase tracking-wider text-[var(--text-muted)]">
                  Registered Custom Nodes ({customNodes.length})
                </span>
              </div>
              {customNodes.length === 0 ? (
                <p className="text-xs text-[var(--text-muted)]">
                  No custom compatible nodes registered yet.
                </p>
              ) : (
                <div className="space-y-2">
                  {customNodes.map((node) => {
                    const confirmingDelete = confirmDeleteNodeId === node.id;
                    const isRowBusy = busyNodeId === node.id;
                    return (
                      <div
                        key={node.id}
                        className="flex flex-col justify-between gap-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 sm:flex-row sm:items-center"
                      >
                        <div className="min-w-0">
                          <div className="flex flex-wrap items-center gap-2">
                            <span className="text-sm font-medium text-[var(--text-primary)]">
                              {node.name}
                            </span>
                            <Badge variant="outline">{node.protocol}</Badge>
                            <code className="font-mono text-xs text-[var(--text-muted)]">
                              {node.id}
                            </code>
                          </div>
                          <p className="mt-1 truncate font-mono text-xs text-[var(--text-secondary)]">
                            {node.base_url}
                          </p>
                        </div>
                        {confirmingDelete ? (
                          <div className="flex shrink-0 items-center gap-1.5">
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => handleConfirmDeleteNode(node)}
                              disabled={anyNodeBusy}
                            >
                              {isRowBusy ? "Deleting..." : "Confirm Delete"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => setConfirmDeleteNodeId(null)}
                              disabled={anyNodeBusy}
                            >
                              Cancel
                            </Button>
                          </div>
                        ) : (
                          <Button
                            type="button"
                            variant="destructive"
                            size="xs"
                            className="shrink-0"
                            onClick={() => setConfirmDeleteNodeId(node.id)}
                            disabled={anyNodeBusy}
                          >
                            <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                            <span>Delete</span>
                          </Button>
                        )}
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Server className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Built-in Provider Adapter Catalog</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Native upstream adapters bundled with the gateway and their default endpoints.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              GET /api/providers/catalog
            </span>
          </CardHeader>
          <CardContent>
            {catalogQuery.isError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Unable to load catalog</AlertTitle>
                  <AlertDescription className="font-mono">
                    {catalogQuery.error.message}
                  </AlertDescription>
                </div>
              </Alert>
            ) : catalogQuery.isPending ? (
              <p className="text-sm text-[var(--text-muted)]">Loading adapter catalog...</p>
            ) : (
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                {adapters.map((adapter) => (
                  <div
                    key={adapter.id}
                    className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5"
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="text-sm font-medium text-[var(--text-primary)]">
                        {adapter.display_name}
                      </span>
                      <Badge variant={adapter.auth_type === "oauth" ? "default" : "secondary"}>
                        {adapter.auth_type === "oauth" ? "OAuth" : "API Key"}
                      </Badge>
                    </div>
                    <div className="mt-1.5 flex flex-wrap items-center gap-1.5">
                      <Badge variant="outline">{adapter.id}</Badge>
                      <Badge variant="outline">{adapter.protocol}</Badge>
                    </div>
                    <p className="mt-2 truncate font-mono text-xs text-[var(--text-muted)]">
                      {adapter.default_base_url || "Custom base_url required"}
                    </p>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </section>
    </div>
  );
}
