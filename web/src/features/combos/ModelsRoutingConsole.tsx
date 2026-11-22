"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertCircle,
  ArrowDown,
  ArrowUp,
  Ban,
  Check,
  CheckCircle2,
  Copy,
  GitFork,
  Globe,
  Layers,
  Plus,
  Power,
  RefreshCw,
  ShieldAlert,
  Tag,
  Terminal,
  Trash2,
  X,
} from "lucide-react";
import { type KeyboardEvent, useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { useProviderConnections } from "@/features/providers/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearModelCooldownMutation,
  useCombos,
  useCreateComboMutation,
  useDeleteComboMutation,
  useDeleteModelAliasMutation,
  useDisabledModels,
  useDisableModelMutation,
  useEnableModelMutation,
  useModelAliases,
  useModelCooldowns,
  useUpdateComboMutation,
  useUpsertModelAliasMutation,
} from "./hooks";
import {
  type Combo,
  type CooldownEntry,
  type CreateComboFormInput,
  type CreateComboFormOutput,
  createComboSchema,
  DEFAULT_CREATE_COMBO_VALUES,
  DEFAULT_DISABLE_MODEL_VALUES,
  DEFAULT_UPSERT_ALIAS_VALUES,
  type DisabledModel,
  type DisableModelFormInput,
  type DisableModelFormOutput,
  disableModelSchema,
  type ModelAlias,
  toCreateComboRequest,
  toDisableModelRequest,
  toUpsertAliasRequest,
  type UpsertAliasFormInput,
  type UpsertAliasFormOutput,
  upsertAliasSchema,
} from "./schemas";

type GatewaySnippetProtocol = "openai-chat" | "openai-responses" | "anthropic" | "gemini";

const DEFAULT_GATEWAY_ORIGIN = "http://127.0.0.1:20128";

const GATEWAY_ENDPOINTS: Array<{
  id: GatewaySnippetProtocol;
  label: string;
  basePath: string;
  operationPath: string;
}> = [
  {
    id: "openai-chat",
    label: "OpenAI Chat",
    basePath: "/v1",
    operationPath: "POST /v1/chat/completions",
  },
  {
    id: "openai-responses",
    label: "OpenAI Responses",
    basePath: "/v1",
    operationPath: "POST /v1/responses",
  },
  {
    id: "anthropic",
    label: "Anthropic Messages",
    basePath: "/v1",
    operationPath: "POST /v1/messages",
  },
  {
    id: "gemini",
    label: "Gemini Native",
    basePath: "/v1beta",
    operationPath: "POST /v1beta/models/{model}:generateContent",
  },
];

function formatTimestamp(iso: string | null): string {
  if (!iso) {
    return "Unknown";
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

function formatRemainingCountdown(untilIso: string, nowMs: number): string {
  const targetMs = Date.parse(untilIso);
  if (Number.isNaN(targetMs)) {
    return "Active";
  }
  const diffSeconds = Math.max(0, Math.ceil((targetMs - nowMs) / 1000));
  if (diffSeconds === 0) {
    return "Expiring";
  }
  const minutes = Math.floor(diffSeconds / 60);
  const seconds = diffSeconds % 60;
  if (minutes > 0) {
    return `${minutes}m ${String(seconds).padStart(2, "0")}s remaining`;
  }
  return `${seconds}s remaining`;
}

function buildCurlSnippet(
  protocol: GatewaySnippetProtocol,
  origin: string,
  sampleModel: string,
): string {
  const model = sampleModel || "gpt-4o";
  switch (protocol) {
    case "openai-chat":
      return `curl -sS ${origin}/v1/chat/completions \\
  -H "Authorization: Bearer sk-om4-YOUR_GATEWAY_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${model}","messages":[{"role":"user","content":"ping"}]}'`;
    case "openai-responses":
      return `curl -sS ${origin}/v1/responses \\
  -H "Authorization: Bearer sk-om4-YOUR_GATEWAY_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${model}","input":"ping"}'`;
    case "anthropic":
      return `curl -sS ${origin}/v1/messages \\
  -H "x-api-key: sk-om4-YOUR_GATEWAY_KEY" \\
  -H "anthropic-version: 2023-06-01" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${model}","max_tokens":256,"messages":[{"role":"user","content":"ping"}]}'`;
    case "gemini":
      return `curl -sS "${origin}/v1beta/models/${encodeURIComponent(model)}:generateContent" \\
  -H "x-goog-api-key: sk-om4-YOUR_GATEWAY_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"contents":[{"role":"user","parts":[{"text":"ping"}]}]}'`;
  }
}

export function ModelsRoutingConsole() {
  const combosQuery = useCombos();
  const aliasesQuery = useModelAliases();
  const disabledModelsQuery = useDisabledModels();
  const cooldownsQuery = useModelCooldowns();
  const connectionsQuery = useProviderConnections();

  const createComboMutation = useCreateComboMutation();
  const updateComboMutation = useUpdateComboMutation();
  const deleteComboMutation = useDeleteComboMutation();
  const upsertAliasMutation = useUpsertModelAliasMutation();
  const deleteAliasMutation = useDeleteModelAliasMutation();
  const disableModelMutation = useDisableModelMutation();
  const enableModelMutation = useEnableModelMutation();
  const clearCooldownMutation = useClearModelCooldownMutation();

  const [gatewayOrigin, setGatewayOrigin] = useState(DEFAULT_GATEWAY_ORIGIN);
  const [selectedProtocol, setSelectedProtocol] = useState<GatewaySnippetProtocol>("openai-chat");
  const [copiedKey, setCopiedKey] = useState<string | null>(null);
  const [copyError, setCopyError] = useState<string | null>(null);

  const [targetDraft, setTargetDraft] = useState("");
  const [targetDraftError, setTargetDraftError] = useState<string | null>(null);
  const [comboFeedback, setComboFeedback] = useState<string | null>(null);
  const [comboRowError, setComboRowError] = useState<string | null>(null);
  const [busyComboId, setBusyComboId] = useState<string | null>(null);
  const [confirmDeleteComboId, setConfirmDeleteComboId] = useState<string | null>(null);

  const [aliasFeedback, setAliasFeedback] = useState<string | null>(null);
  const [aliasRowError, setAliasRowError] = useState<string | null>(null);
  const [busyAlias, setBusyAlias] = useState<string | null>(null);
  const [confirmDeleteAlias, setConfirmDeleteAlias] = useState<string | null>(null);

  const [disabledFeedback, setDisabledFeedback] = useState<string | null>(null);
  const [disabledRowError, setDisabledRowError] = useState<string | null>(null);
  const [busyDisabledModel, setBusyDisabledModel] = useState<string | null>(null);

  const [cooldownFeedback, setCooldownFeedback] = useState<string | null>(null);
  const [cooldownRowError, setCooldownRowError] = useState<string | null>(null);
  const [busyCooldownId, setBusyCooldownId] = useState<string | null>(null);
  const [nowMs, setNowMs] = useState(() => Date.now());

  useEffect(() => {
    if (typeof window !== "undefined" && window.location.origin) {
      setGatewayOrigin(window.location.origin);
    }
  }, []);

  const cooldownCount = cooldownsQuery.data?.length ?? 0;
  useEffect(() => {
    if (cooldownCount === 0) {
      return;
    }
    const timer = window.setInterval(() => {
      setNowMs(Date.now());
    }, 1000);
    return () => window.clearInterval(timer);
  }, [cooldownCount]);

  const createComboForm = useForm<CreateComboFormInput, unknown, CreateComboFormOutput>({
    resolver: zodResolver(createComboSchema),
    defaultValues: DEFAULT_CREATE_COMBO_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const upsertAliasForm = useForm<UpsertAliasFormInput, unknown, UpsertAliasFormOutput>({
    resolver: zodResolver(upsertAliasSchema),
    defaultValues: DEFAULT_UPSERT_ALIAS_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const disableModelForm = useForm<DisableModelFormInput, unknown, DisableModelFormOutput>({
    resolver: zodResolver(disableModelSchema),
    defaultValues: DEFAULT_DISABLE_MODEL_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const combos: Combo[] = combosQuery.data ?? [];
  const aliases: ModelAlias[] = aliasesQuery.data ?? [];
  const disabledModels: DisabledModel[] = disabledModelsQuery.data ?? [];
  const cooldowns: CooldownEntry[] = cooldownsQuery.data ?? [];
  const connections = connectionsQuery.data ?? [];

  const connectionLookup = useMemo(() => {
    const map = new Map<string, { name: string; providerId: string }>();
    for (const conn of connections) {
      map.set(conn.id, { name: conn.name, providerId: conn.provider_id });
    }
    return map;
  }, [connections]);

  const suggestedProviderModels = useMemo(() => {
    const seen = new Set<string>();
    const out: string[] = [];
    for (const conn of connections) {
      if (!conn.enabled) {
        continue;
      }
      for (const rawModel of conn.models) {
        const model = rawModel.trim();
        if (!model) {
          continue;
        }
        const qualified = model.includes("/") ? model : `${conn.provider_id}/${model}`;
        if (!seen.has(qualified)) {
          seen.add(qualified);
          out.push(qualified);
        }
      }
    }
    return out;
  }, [connections]);

  const routingTargetSuggestions = useMemo(() => {
    const seen = new Set<string>();
    const out: string[] = [];
    for (const combo of combos) {
      if (combo.enabled && !seen.has(combo.name)) {
        seen.add(combo.name);
        out.push(combo.name);
      }
    }
    for (const model of suggestedProviderModels) {
      if (!seen.has(model)) {
        seen.add(model);
        out.push(model);
      }
    }
    return out;
  }, [combos, suggestedProviderModels]);

  const watchedTargets = createComboForm.watch("targets");
  const sampleModel =
    combos.find((c) => c.enabled)?.name ??
    aliases[0]?.alias ??
    suggestedProviderModels[0] ??
    "gpt-4o";
  const curlSnippet = buildCurlSnippet(selectedProtocol, gatewayOrigin, sampleModel);

  const isRefreshing =
    combosQuery.isFetching ||
    aliasesQuery.isFetching ||
    disabledModelsQuery.isFetching ||
    cooldownsQuery.isFetching;

  const handleRefreshAll = () => {
    void Promise.all([
      combosQuery.refetch(),
      aliasesQuery.refetch(),
      disabledModelsQuery.refetch(),
      cooldownsQuery.refetch(),
      connectionsQuery.refetch(),
    ]);
  };

  const handleCopyText = async (key: string, text: string) => {
    setCopyError(null);
    try {
      await navigator.clipboard.writeText(text);
      setCopiedKey(key);
      window.setTimeout(() => {
        setCopiedKey((current) => (current === key ? null : current));
      }, 2200);
    } catch {
      setCopyError("Clipboard access was denied. Copy the endpoint or command manually.");
    }
  };

  const addComboTarget = (rawTarget: string) => {
    const trimmed = rawTarget.trim();
    if (!trimmed) {
      setTargetDraftError("Enter a model or provider/model identifier first.");
      return;
    }
    if (trimmed === "*") {
      setTargetDraftError("Target model must not be wildcard '*'.");
      return;
    }
    if (/[\r\n\t]/.test(trimmed)) {
      setTargetDraftError("Target model must not contain control whitespace.");
      return;
    }
    const current = createComboForm.getValues("targets");
    if (current.includes(trimmed)) {
      setTargetDraftError(`"${trimmed}" is already in the fallback order.`);
      return;
    }
    setTargetDraftError(null);
    createComboForm.setValue("targets", [...current, trimmed], {
      shouldDirty: true,
      shouldValidate: true,
    });
    setTargetDraft("");
  };

  const handleTargetDraftKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      event.preventDefault();
      addComboTarget(targetDraft);
    }
  };

  const moveComboTarget = (index: number, direction: -1 | 1) => {
    const current = [...createComboForm.getValues("targets")];
    const nextIndex = index + direction;
    if (nextIndex < 0 || nextIndex >= current.length) {
      return;
    }
    const item = current[index];
    const swapItem = current[nextIndex];
    if (item === undefined || swapItem === undefined) {
      return;
    }
    current[index] = swapItem;
    current[nextIndex] = item;
    createComboForm.setValue("targets", current, {
      shouldDirty: true,
      shouldValidate: true,
    });
  };

  const removeComboTarget = (index: number) => {
    const current = createComboForm.getValues("targets");
    createComboForm.setValue(
      "targets",
      current.filter((_, idx) => idx !== index),
      { shouldDirty: true, shouldValidate: true },
    );
  };

  const onCreateComboSubmit = (values: CreateComboFormOutput) => {
    setComboFeedback(null);
    setTargetDraftError(null);
    createComboMutation.reset();

    createComboMutation.mutate(toCreateComboRequest(values), {
      onSuccess: (created) => {
        setComboFeedback(
          `Created fallback combo "${created.name}" (${created.targets.length} targets).`,
        );
        createComboForm.reset(DEFAULT_CREATE_COMBO_VALUES);
        setTargetDraft("");
      },
    });
  };

  const handleToggleCombo = (combo: Combo) => {
    if (busyComboId !== null) {
      return;
    }
    setComboRowError(null);
    setConfirmDeleteComboId(null);
    setBusyComboId(combo.id);

    updateComboMutation.mutate(
      { id: combo.id, body: { enabled: !combo.enabled } },
      {
        onError: (err) => setComboRowError(err.message),
        onSettled: () => setBusyComboId(null),
      },
    );
  };

  const handleConfirmDeleteCombo = (combo: Combo) => {
    if (busyComboId !== null) {
      return;
    }
    setComboRowError(null);
    setBusyComboId(combo.id);

    deleteComboMutation.mutate(combo.id, {
      onSuccess: () => {
        setConfirmDeleteComboId(null);
      },
      onError: (err) => setComboRowError(err.message),
      onSettled: () => setBusyComboId(null),
    });
  };

  const onUpsertAliasSubmit = (values: UpsertAliasFormOutput) => {
    setAliasFeedback(null);
    upsertAliasMutation.reset();

    upsertAliasMutation.mutate(
      {
        alias: values.alias,
        body: toUpsertAliasRequest(values),
      },
      {
        onSuccess: (saved) => {
          setAliasFeedback(`Saved alias "${saved.alias}" -> "${saved.target_model}".`);
          upsertAliasForm.reset(DEFAULT_UPSERT_ALIAS_VALUES);
        },
      },
    );
  };

  const handleConfirmDeleteAlias = (item: ModelAlias) => {
    if (busyAlias !== null) {
      return;
    }
    setAliasRowError(null);
    setBusyAlias(item.alias);

    deleteAliasMutation.mutate(item.alias, {
      onSuccess: () => {
        setConfirmDeleteAlias(null);
      },
      onError: (err) => setAliasRowError(err.message),
      onSettled: () => setBusyAlias(null),
    });
  };

  const onDisableModelSubmit = (values: DisableModelFormOutput) => {
    setDisabledFeedback(null);
    disableModelMutation.reset();

    disableModelMutation.mutate(toDisableModelRequest(values), {
      onSuccess: (saved) => {
        setDisabledFeedback(`Suppressed model "${saved.model}" from routing and listing.`);
        disableModelForm.reset(DEFAULT_DISABLE_MODEL_VALUES);
      },
    });
  };

  const handleEnableModel = (item: DisabledModel) => {
    if (busyDisabledModel !== null) {
      return;
    }
    setDisabledRowError(null);
    setBusyDisabledModel(item.model);

    enableModelMutation.mutate(item.model, {
      onError: (err) => setDisabledRowError(err.message),
      onSettled: () => setBusyDisabledModel(null),
    });
  };

  const handleClearCooldown = (entry: CooldownEntry) => {
    if (busyCooldownId !== null) {
      return;
    }
    setCooldownRowError(null);
    setCooldownFeedback(null);
    setBusyCooldownId(entry.connection_id);

    clearCooldownMutation.mutate(entry.connection_id, {
      onSuccess: () => {
        setCooldownFeedback(`Cleared cooldown window for connection ${entry.connection_id}.`);
      },
      onError: (err) => setCooldownRowError(err.message),
      onSettled: () => setBusyCooldownId(null),
    });
  };

  const anyComboBusy = busyComboId !== null;
  const anyAliasBusy = busyAlias !== null;
  const anyDisabledBusy = busyDisabledModel !== null;
  const anyCooldownBusy = busyCooldownId !== null;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Routing &amp; Fallback Studio
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Models &amp; Routing
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Configure ordered fallback combos, logical model aliases, model suppression rules, live
            cooldown recovery, and gateway endpoint routes.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <Badge variant="outline" className="tabular-nums">
            {combos.length} Combos
          </Badge>
          <Badge variant="outline" className="tabular-nums">
            {aliases.length} Aliases
          </Badge>
          <Badge variant="secondary" className="tabular-nums">
            {disabledModels.length} Disabled
          </Badge>
          <Badge variant={cooldowns.length > 0 ? "warning" : "success"} className="tabular-nums">
            {cooldowns.length} Cooldowns
          </Badge>
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

      {/* Gateway Endpoint Quick-Reference */}
      <Card>
        <CardHeader>
          <div>
            <div className="flex items-center gap-2">
              <Globe className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Gateway Endpoint Quick-Reference</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Connect OpenAI, Anthropic, or Gemini clients directly to this single-binary listener.
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">{gatewayOrigin}</span>
        </CardHeader>
        <CardContent className="space-y-4">
          {copyError ? (
            <Alert variant="warning">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                aria-hidden="true"
              />
              <AlertDescription>{copyError}</AlertDescription>
            </Alert>
          ) : null}

          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
            {GATEWAY_ENDPOINTS.map((endpoint) => {
              const fullBaseUrl = `${gatewayOrigin}${endpoint.basePath}`;
              const copyId = `endpoint-${endpoint.id}`;
              const isCopied = copiedKey === copyId;

              return (
                <div
                  key={endpoint.id}
                  className="flex flex-col justify-between rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5"
                >
                  <div>
                    <div className="flex items-center justify-between gap-2">
                      <span className="text-xs font-semibold text-[var(--text-primary)]">
                        {endpoint.label}
                      </span>
                      <Badge variant="outline">{endpoint.basePath}</Badge>
                    </div>
                    <p className="mt-2 truncate font-mono text-xs text-[var(--text-primary)]">
                      {fullBaseUrl}
                    </p>
                    <p className="mt-1 truncate font-mono text-[11px] text-[var(--text-muted)]">
                      {endpoint.operationPath}
                    </p>
                  </div>
                  <div className="mt-3 flex justify-end">
                    <Button
                      type="button"
                      variant="secondary"
                      size="xs"
                      onClick={() => void handleCopyText(copyId, fullBaseUrl)}
                    >
                      {isCopied ? (
                        <Check className="h-3 w-3 text-[var(--status-ok)]" aria-hidden="true" />
                      ) : (
                        <Copy className="h-3 w-3" aria-hidden="true" />
                      )}
                      <span>{isCopied ? "Copied" : "Copy Base URL"}</span>
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>

          <div className="rounded-md border border-[var(--border-strong)] bg-[var(--bg-canvas)] p-4">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div className="flex flex-wrap items-center gap-1.5">
                <Terminal className="mr-1 h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                {GATEWAY_ENDPOINTS.map((endpoint) => (
                  <Button
                    key={endpoint.id}
                    type="button"
                    size="xs"
                    variant={selectedProtocol === endpoint.id ? "default" : "secondary"}
                    onClick={() => setSelectedProtocol(endpoint.id)}
                  >
                    {endpoint.label}
                  </Button>
                ))}
              </div>
              <Button
                type="button"
                variant="outline"
                size="xs"
                onClick={() => void handleCopyText("curl-snippet", curlSnippet)}
              >
                {copiedKey === "curl-snippet" ? (
                  <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" aria-hidden="true" />
                ) : (
                  <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                )}
                <span>{copiedKey === "curl-snippet" ? "Copied cURL" : "Copy cURL"}</span>
              </Button>
            </div>
            <pre className="mt-3 overflow-x-auto rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-3 font-mono text-xs leading-relaxed text-[var(--text-primary)]">
              <code>{curlSnippet}</code>
            </pre>
          </div>
        </CardContent>
      </Card>

      <EmberRule />

      {/* Ordered Fallback Combos Section */}
      <section className="grid grid-cols-1 gap-6 xl:grid-cols-5">
        <Card className="xl:col-span-2">
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <GitFork className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Ordered Fallback Combo Builder</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Chain upstream targets in deterministic failover order.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">POST /api/combos</span>
          </CardHeader>
          <CardContent>
            <form
              onSubmit={createComboForm.handleSubmit(onCreateComboSubmit)}
              noValidate
              className="space-y-4"
            >
              {createComboMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to create combo</AlertTitle>
                    <AlertDescription className="font-mono">
                      {createComboMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {comboFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Combo Created</AlertTitle>
                    <AlertDescription>{comboFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div>
                <Label htmlFor="combo-name">Combo Name</Label>
                <Input
                  id="combo-name"
                  type="text"
                  placeholder="e.g. coding-fallback"
                  aria-invalid={Boolean(createComboForm.formState.errors.name)}
                  className="mt-1.5 font-mono"
                  {...createComboForm.register("name")}
                />
                {createComboForm.formState.errors.name ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createComboForm.formState.errors.name.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="combo-description">Description (Optional)</Label>
                <Input
                  id="combo-description"
                  type="text"
                  placeholder="Primary reasoning model with automatic cross-provider failover"
                  aria-invalid={Boolean(createComboForm.formState.errors.description)}
                  className="mt-1.5"
                  {...createComboForm.register("description")}
                />
                {createComboForm.formState.errors.description ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {createComboForm.formState.errors.description.message}
                  </p>
                ) : null}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label htmlFor="combo-strategy">Strategy</Label>
                  <Input
                    id="combo-strategy"
                    type="text"
                    readOnly
                    className="mt-1.5 font-mono"
                    {...createComboForm.register("strategy")}
                  />
                </div>
                <div className="flex items-end pb-2">
                  <label className="flex cursor-pointer items-center gap-2 text-sm text-[var(--text-primary)]">
                    <input
                      type="checkbox"
                      className="h-4 w-4 accent-[#E56A4A]"
                      {...createComboForm.register("enabled")}
                    />
                    <span>Enabled immediately</span>
                  </label>
                </div>
              </div>

              <div className="space-y-2.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3.5">
                <div className="flex items-center justify-between">
                  <Label htmlFor="combo-target-draft">Ordered Fallback Targets</Label>
                  <span className="font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
                    {watchedTargets.length} selected
                  </span>
                </div>

                <div className="flex gap-2">
                  <Input
                    id="combo-target-draft"
                    type="text"
                    value={targetDraft}
                    onChange={(event) => {
                      setTargetDraft(event.target.value);
                      if (targetDraftError) {
                        setTargetDraftError(null);
                      }
                    }}
                    onKeyDown={handleTargetDraftKeyDown}
                    placeholder="provider/model or model ID"
                    className="font-mono"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => addComboTarget(targetDraft)}
                    className="h-10 shrink-0"
                  >
                    <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>Add</span>
                  </Button>
                </div>

                {targetDraftError ? (
                  <p className="text-xs text-[var(--status-error)]">{targetDraftError}</p>
                ) : null}

                {suggestedProviderModels.length > 0 ? (
                  <div className="space-y-1.5 pt-1">
                    <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                      Quick-pick from enabled providers
                    </p>
                    <div className="flex flex-wrap gap-1.5">
                      {suggestedProviderModels.slice(0, 12).map((modelId) => {
                        const alreadyAdded = watchedTargets.includes(modelId);
                        return (
                          <Button
                            key={modelId}
                            type="button"
                            size="xs"
                            variant={alreadyAdded ? "secondary" : "outline"}
                            disabled={alreadyAdded}
                            onClick={() => addComboTarget(modelId)}
                            className="font-mono text-[11px]"
                          >
                            <Plus className="h-3 w-3" aria-hidden="true" />
                            <span>{modelId}</span>
                          </Button>
                        );
                      })}
                    </div>
                  </div>
                ) : null}

                {watchedTargets.length === 0 ? (
                  <p className="rounded border border-dashed border-[var(--border-subtle)] px-3 py-3 text-center text-xs text-[var(--text-muted)]">
                    No targets queued yet. Add at least one upstream model in failover priority
                    order.
                  </p>
                ) : (
                  <ol className="space-y-1.5 pt-1">
                    {watchedTargets.map((target, index) => (
                      <li
                        key={target}
                        className="flex items-center justify-between gap-2 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-3 py-2"
                      >
                        <div className="flex min-w-0 items-center gap-2">
                          <Badge variant="outline" className="tabular-nums">
                            #{index + 1}
                          </Badge>
                          <span className="truncate font-mono text-xs text-[var(--text-primary)]">
                            {target}
                          </span>
                        </div>
                        <div className="flex shrink-0 items-center gap-1">
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => moveComboTarget(index, -1)}
                            disabled={index === 0}
                            aria-label={`Move ${target} up`}
                          >
                            <ArrowUp className="h-3.5 w-3.5" aria-hidden="true" />
                          </Button>
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => moveComboTarget(index, 1)}
                            disabled={index === watchedTargets.length - 1}
                            aria-label={`Move ${target} down`}
                          >
                            <ArrowDown className="h-3.5 w-3.5" aria-hidden="true" />
                          </Button>
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => removeComboTarget(index)}
                            aria-label={`Remove ${target}`}
                          >
                            <X
                              className="h-3.5 w-3.5 text-[var(--status-error)]"
                              aria-hidden="true"
                            />
                          </Button>
                        </div>
                      </li>
                    ))}
                  </ol>
                )}

                {createComboForm.formState.errors.targets?.message ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {createComboForm.formState.errors.targets.message}
                  </p>
                ) : null}
              </div>

              <Button type="submit" disabled={createComboMutation.isPending} className="w-full">
                <Layers className="h-4 w-4" aria-hidden="true" />
                <span>{createComboMutation.isPending ? "Creating Combo..." : "Create Combo"}</span>
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card className="xl:col-span-3">
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Layers className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Configured Fallback Combos</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Deterministic ordered fallback chains resolved prior to provider selection.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">GET /api/combos</span>
          </CardHeader>

          {comboRowError ? (
            <div className="px-5 pt-4">
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Combo operation failed</AlertTitle>
                  <AlertDescription className="font-mono">{comboRowError}</AlertDescription>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setComboRowError(null)}
                  aria-label="Dismiss combo error"
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </Alert>
            </div>
          ) : null}

          {combosQuery.isError ? (
            <CardContent>
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Unable to load combos</AlertTitle>
                  <AlertDescription className="font-mono">
                    {combosQuery.error.message}
                  </AlertDescription>
                </div>
              </Alert>
            </CardContent>
          ) : combosQuery.isPending ? (
            <CardContent className="text-sm text-[var(--text-muted)]">
              Loading fallback combos...
            </CardContent>
          ) : combos.length === 0 ? (
            <CardContent className="py-10 text-center">
              <p className="text-sm font-medium text-[var(--text-primary)]">
                No fallback combos configured yet
              </p>
              <p className="mt-1 text-xs text-[var(--text-muted)]">
                Build an ordered fallback chain on the left to route resilient multi-provider model
                names.
              </p>
            </CardContent>
          ) : (
            <div className="divide-y divide-[var(--border-subtle)]">
              {combos.map((combo) => {
                const isRowBusy = busyComboId === combo.id;
                const confirmingDelete = confirmDeleteComboId === combo.id;

                return (
                  <div
                    key={combo.id}
                    className="flex flex-col gap-3 px-5 py-4 transition-colors hover:bg-[var(--bg-surface-raised)] sm:flex-row sm:items-start sm:justify-between"
                  >
                    <div className="min-w-0 flex-1 space-y-2">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="font-mono text-sm font-semibold text-[var(--text-primary)]">
                          {combo.name}
                        </span>
                        <Badge variant={combo.enabled ? "success" : "warning"}>
                          {combo.enabled ? "Enabled" : "Disabled"}
                        </Badge>
                        <Badge variant="outline">{combo.strategy}</Badge>
                      </div>

                      {combo.description ? (
                        <p className="text-xs text-[var(--text-secondary)]">{combo.description}</p>
                      ) : null}

                      <div className="flex flex-wrap items-center gap-1.5">
                        {combo.targets.map((target, idx) => (
                          <Badge key={target} variant="secondary" className="tabular-nums">
                            {idx + 1}. {target}
                          </Badge>
                        ))}
                      </div>

                      <p className="font-mono text-[11px] text-[var(--text-muted)]">
                        {combo.id} · Updated {formatTimestamp(combo.updated_at)}
                      </p>
                    </div>

                    <div className="flex shrink-0 items-center gap-1.5 sm:pt-0.5">
                      {confirmingDelete ? (
                        <>
                          <Button
                            type="button"
                            variant="destructive"
                            size="xs"
                            onClick={() => handleConfirmDeleteCombo(combo)}
                            disabled={anyComboBusy}
                          >
                            {isRowBusy ? "Deleting..." : "Confirm Delete"}
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            size="xs"
                            onClick={() => setConfirmDeleteComboId(null)}
                            disabled={anyComboBusy}
                          >
                            Cancel
                          </Button>
                        </>
                      ) : (
                        <>
                          <Button
                            type="button"
                            variant="secondary"
                            size="xs"
                            onClick={() => handleToggleCombo(combo)}
                            disabled={anyComboBusy}
                          >
                            <Power className="h-3.5 w-3.5" aria-hidden="true" />
                            <span>{combo.enabled ? "Disable" : "Enable"}</span>
                          </Button>
                          <Button
                            type="button"
                            variant="destructive"
                            size="xs"
                            onClick={() => setConfirmDeleteComboId(combo.id)}
                            disabled={anyComboBusy}
                            aria-label={`Delete combo ${combo.name}`}
                          >
                            <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                            <span>Delete</span>
                          </Button>
                        </>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </Card>
      </section>

      <EmberRule />

      {/* Model Aliases & Disabled Models */}
      <section className="grid grid-cols-1 gap-6 xl:grid-cols-2">
        {/* Model Aliases Editor */}
        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Tag className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Model Aliases Editor</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Map logical client names to target models or fallback combos.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">PUT /api/aliases</span>
          </CardHeader>
          <CardContent className="space-y-5">
            <form
              onSubmit={upsertAliasForm.handleSubmit(onUpsertAliasSubmit)}
              noValidate
              className="space-y-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4"
            >
              {upsertAliasMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to save alias</AlertTitle>
                    <AlertDescription className="font-mono">
                      {upsertAliasMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {aliasFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Alias Saved</AlertTitle>
                    <AlertDescription>{aliasFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div>
                  <Label htmlFor="alias-name">Alias Name</Label>
                  <Input
                    id="alias-name"
                    type="text"
                    placeholder="e.g. smart-default"
                    aria-invalid={Boolean(upsertAliasForm.formState.errors.alias)}
                    className="mt-1.5 font-mono"
                    {...upsertAliasForm.register("alias")}
                  />
                  {upsertAliasForm.formState.errors.alias ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {upsertAliasForm.formState.errors.alias.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="alias-target">Target Model / Combo</Label>
                  <Input
                    id="alias-target"
                    type="text"
                    placeholder="e.g. coding-fallback or openai/gpt-4o"
                    aria-invalid={Boolean(upsertAliasForm.formState.errors.target_model)}
                    className="mt-1.5 font-mono"
                    {...upsertAliasForm.register("target_model")}
                  />
                  {upsertAliasForm.formState.errors.target_model ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {upsertAliasForm.formState.errors.target_model.message}
                    </p>
                  ) : null}
                </div>
              </div>

              {routingTargetSuggestions.length > 0 ? (
                <div className="flex flex-wrap items-center gap-1.5">
                  <span className="text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    Target quick-pick:
                  </span>
                  {routingTargetSuggestions.slice(0, 8).map((suggestion) => (
                    <Button
                      key={suggestion}
                      type="button"
                      size="xs"
                      variant="outline"
                      onClick={() =>
                        upsertAliasForm.setValue("target_model", suggestion, {
                          shouldDirty: true,
                          shouldValidate: true,
                        })
                      }
                      className="font-mono text-[11px]"
                    >
                      {suggestion}
                    </Button>
                  ))}
                </div>
              ) : null}

              <div className="flex justify-end">
                <Button type="submit" size="sm" disabled={upsertAliasMutation.isPending}>
                  <Tag className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{upsertAliasMutation.isPending ? "Saving..." : "Save Alias"}</span>
                </Button>
              </div>
            </form>

            {aliasRowError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Alias action failed</AlertTitle>
                  <AlertDescription className="font-mono">{aliasRowError}</AlertDescription>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setAliasRowError(null)}
                  aria-label="Dismiss alias error"
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </Alert>
            ) : null}

            {aliasesQuery.isError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <AlertDescription className="font-mono">
                  {aliasesQuery.error.message}
                </AlertDescription>
              </Alert>
            ) : aliasesQuery.isPending ? (
              <p className="text-sm text-[var(--text-muted)]">Loading model aliases...</p>
            ) : aliases.length === 0 ? (
              <p className="py-4 text-center text-xs text-[var(--text-muted)]">
                No model aliases configured.
              </p>
            ) : (
              <div className="divide-y divide-[var(--border-subtle)] rounded-md border border-[var(--border-subtle)]">
                {aliases.map((item) => {
                  const isRowBusy = busyAlias === item.alias;
                  const confirming = confirmDeleteAlias === item.alias;

                  return (
                    <div
                      key={item.alias}
                      className="flex flex-col gap-2 px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
                    >
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-2 font-mono text-xs">
                          <span className="font-semibold text-[var(--text-primary)]">
                            {item.alias}
                          </span>
                          <span className="text-[var(--text-muted)]">-&gt;</span>
                          <Badge variant="default">{item.target_model}</Badge>
                        </div>
                        <p className="mt-1 font-mono text-[11px] text-[var(--text-muted)]">
                          Updated {formatTimestamp(item.updated_at)}
                        </p>
                      </div>

                      <div className="flex shrink-0 items-center gap-1.5">
                        {confirming ? (
                          <>
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => handleConfirmDeleteAlias(item)}
                              disabled={anyAliasBusy}
                            >
                              {isRowBusy ? "Deleting..." : "Confirm"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              onClick={() => setConfirmDeleteAlias(null)}
                              disabled={anyAliasBusy}
                            >
                              Cancel
                            </Button>
                          </>
                        ) : (
                          <>
                            <Button
                              type="button"
                              variant="outline"
                              size="xs"
                              onClick={() => {
                                upsertAliasForm.setValue("alias", item.alias, {
                                  shouldDirty: true,
                                });
                                upsertAliasForm.setValue("target_model", item.target_model, {
                                  shouldDirty: true,
                                });
                              }}
                              disabled={anyAliasBusy}
                            >
                              Edit
                            </Button>
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              onClick={() => setConfirmDeleteAlias(item.alias)}
                              disabled={anyAliasBusy}
                              aria-label={`Delete alias ${item.alias}`}
                            >
                              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>Delete</span>
                            </Button>
                          </>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Disabled Models Panel */}
        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Ban className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Disabled Models Panel</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Suppress bare or provider-qualified model IDs from gateway routing and /v1/models.
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">/api/models/disabled</span>
          </CardHeader>
          <CardContent className="space-y-5">
            <form
              onSubmit={disableModelForm.handleSubmit(onDisableModelSubmit)}
              noValidate
              className="space-y-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4"
            >
              {disableModelMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to disable model</AlertTitle>
                    <AlertDescription className="font-mono">
                      {disableModelMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {disabledFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Model Suppressed</AlertTitle>
                    <AlertDescription>{disabledFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div>
                  <Label htmlFor="disable-model-id">Model Identifier</Label>
                  <Input
                    id="disable-model-id"
                    type="text"
                    placeholder="e.g. openai/gpt-4-turbo-preview"
                    aria-invalid={Boolean(disableModelForm.formState.errors.model)}
                    className="mt-1.5 font-mono"
                    {...disableModelForm.register("model")}
                  />
                  {disableModelForm.formState.errors.model ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {disableModelForm.formState.errors.model.message}
                    </p>
                  ) : null}
                </div>

                <div>
                  <Label htmlFor="disable-model-reason">Reason (Optional)</Label>
                  <Input
                    id="disable-model-reason"
                    type="text"
                    placeholder="e.g. Deprecated upstream or cost guard"
                    aria-invalid={Boolean(disableModelForm.formState.errors.reason)}
                    className="mt-1.5"
                    {...disableModelForm.register("reason")}
                  />
                  {disableModelForm.formState.errors.reason ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {disableModelForm.formState.errors.reason.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div className="flex justify-end">
                <Button type="submit" size="sm" disabled={disableModelMutation.isPending}>
                  <Ban className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{disableModelMutation.isPending ? "Disabling..." : "Disable Model"}</span>
                </Button>
              </div>
            </form>

            {disabledRowError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Model enable failed</AlertTitle>
                  <AlertDescription className="font-mono">{disabledRowError}</AlertDescription>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setDisabledRowError(null)}
                  aria-label="Dismiss disabled model error"
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </Alert>
            ) : null}

            {disabledModelsQuery.isError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <AlertDescription className="font-mono">
                  {disabledModelsQuery.error.message}
                </AlertDescription>
              </Alert>
            ) : disabledModelsQuery.isPending ? (
              <p className="text-sm text-[var(--text-muted)]">Loading disabled models...</p>
            ) : disabledModels.length === 0 ? (
              <p className="py-4 text-center text-xs text-[var(--text-muted)]">
                No models are currently suppressed.
              </p>
            ) : (
              <div className="divide-y divide-[var(--border-subtle)] rounded-md border border-[var(--border-subtle)]">
                {disabledModels.map((item) => {
                  const isRowBusy = busyDisabledModel === item.model;

                  return (
                    <div
                      key={item.model}
                      className="flex flex-col gap-2 px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
                    >
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className="font-mono text-xs font-semibold text-[var(--text-primary)]">
                            {item.model}
                          </span>
                          <Badge variant="warning">Disabled</Badge>
                        </div>
                        <p className="mt-1 text-xs text-[var(--text-secondary)]">
                          {item.reason || "No suppression reason recorded"}
                        </p>
                        <p className="mt-0.5 font-mono text-[11px] text-[var(--text-muted)]">
                          Since {formatTimestamp(item.created_at)}
                        </p>
                      </div>

                      <Button
                        type="button"
                        variant="outline"
                        size="xs"
                        onClick={() => handleEnableModel(item)}
                        disabled={anyDisabledBusy}
                        className="shrink-0"
                      >
                        <Power className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>{isRowBusy ? "Enabling..." : "Re-enable"}</span>
                      </Button>
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>
      </section>

      <EmberRule />

      {/* Account Cooldowns Monitor */}
      <Card>
        <CardHeader>
          <div>
            <div className="flex items-center gap-2">
              <ShieldAlert className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Account Cooldowns Monitor</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Live upstream connection backoff windows (auto-refreshes every 10s). Release a
              connection manually to resume routing immediately.
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            GET /api/models/cooldowns
          </span>
        </CardHeader>

        {cooldownRowError ? (
          <div className="px-5 pt-4">
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div className="flex-1">
                <AlertTitle>Failed to clear cooldown</AlertTitle>
                <AlertDescription className="font-mono">{cooldownRowError}</AlertDescription>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setCooldownRowError(null)}
                aria-label="Dismiss cooldown error"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
              </Button>
            </Alert>
          </div>
        ) : null}

        {cooldownFeedback ? (
          <div className="px-5 pt-4">
            <Alert>
              <CheckCircle2
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                aria-hidden="true"
              />
              <AlertDescription>{cooldownFeedback}</AlertDescription>
            </Alert>
          </div>
        ) : null}

        {cooldownsQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load cooldown windows</AlertTitle>
                <AlertDescription className="font-mono">
                  {cooldownsQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : cooldownsQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Checking upstream cooldown state...
          </CardContent>
        ) : cooldowns.length === 0 ? (
          <CardContent className="py-8 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              All upstream connections healthy
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              No provider connections are currently in failure backoff cooldown.
            </p>
          </CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <th scope="col" className="px-5 py-3">
                    Connection ID
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Failure Count
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Last HTTP Status
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Cooldown Expiry
                  </th>
                  <th scope="col" className="px-5 py-3 text-right">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {cooldowns.map((entry) => {
                  const isRowBusy = busyCooldownId === entry.connection_id;
                  const meta = connectionLookup.get(entry.connection_id);

                  return (
                    <tr
                      key={entry.connection_id}
                      className="transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3.5">
                        {meta ? (
                          <div className="font-medium text-[var(--text-primary)]">
                            {meta.name}{" "}
                            <span className="font-mono text-xs text-[var(--text-muted)]">
                              ({meta.providerId})
                            </span>
                          </div>
                        ) : null}
                        <p className="font-mono text-xs text-[var(--text-secondary)]">
                          {entry.connection_id}
                        </p>
                      </td>
                      <td className="px-4 py-3.5 font-mono text-xs tabular-nums text-[var(--text-primary)]">
                        {entry.failures}
                      </td>
                      <td className="px-4 py-3.5">
                        <Badge variant="destructive" className="tabular-nums">
                          {entry.last_status > 0 ? `HTTP ${entry.last_status}` : "NET ERR"}
                        </Badge>
                      </td>
                      <td className="px-4 py-3.5">
                        <p className="font-mono text-xs font-medium tabular-nums text-[var(--status-warn)]">
                          {formatRemainingCountdown(entry.until, nowMs)}
                        </p>
                        <p className="mt-0.5 font-mono text-[11px] text-[var(--text-muted)]">
                          Until {formatTimestamp(entry.until)}
                        </p>
                      </td>
                      <td className="px-5 py-3.5 text-right">
                        <Button
                          type="button"
                          variant="outline"
                          size="xs"
                          onClick={() => handleClearCooldown(entry)}
                          disabled={anyCooldownBusy}
                        >
                          <Power className="h-3.5 w-3.5" aria-hidden="true" />
                          <span>{isRowBusy ? "Clearing..." : "Clear Cooldown"}</span>
                        </Button>
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
