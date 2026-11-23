"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Play, RefreshCw, RotateCcw, Trash2, Zap } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearHeadroomCacheMutation,
  usePreviewTokenSaverMutation,
  useResetTokenSaverStatsMutation,
  useTokenSaverConfig,
  useTokenSaverStats,
  useUpdateTokenSaverConfigMutation,
} from "./hooks";
import {
  parsePreviewBodyJson,
  TOKEN_SAVER_STRATEGIES,
  type TokenSaverConfigFormValues,
  tokenSaverConfigSchema,
} from "./schemas";

export function TokenSaverConsole() {
  const configQuery = useTokenSaverConfig();
  const statsQuery = useTokenSaverStats();

  const updateConfig = useUpdateTokenSaverConfigMutation();
  const resetStats = useResetTokenSaverStatsMutation();
  const clearCache = useClearHeadroomCacheMutation();
  const previewMutation = usePreviewTokenSaverMutation();

  const [feedback, setFeedback] = useState<string | null>(null);
  const [previewErr, setPreviewErr] = useState<string | null>(null);
  const [previewJson, setPreviewJson] = useState(
    '{"model":"gpt-4o","messages":[{"role":"system","content":"You   are   a   concise   assistant.\\n\\n\\n\\nFollow instructions."},{"role":"user","content":"Explain   Go   channels."}]}',
  );

  const form = useForm<TokenSaverConfigFormValues>({
    resolver: zodResolver(tokenSaverConfigSchema),
    defaultValues: {
      enabled: false,
      strategy: "whitespace",
      max_history_messages: 40,
      max_tool_output_chars: 8000,
      headroom_enabled: true,
      prefix_depth: 2,
      max_cache_entries: 1024,
      cache_ttl_seconds: 1800,
      inject_cache_control: true,
    },
  });

  useEffect(() => {
    if (configQuery.data) {
      form.reset({
        enabled: configQuery.data.enabled,
        strategy: configQuery.data.strategy,
        max_history_messages: configQuery.data.max_history_messages,
        max_tool_output_chars: configQuery.data.max_tool_output_chars,
        headroom_enabled: configQuery.data.headroom_enabled,
        prefix_depth: configQuery.data.prefix_depth,
        max_cache_entries: configQuery.data.max_cache_entries,
        cache_ttl_seconds: configQuery.data.cache_ttl_seconds,
        inject_cache_control: configQuery.data.inject_cache_control,
      });
    }
  }, [configQuery.data, form]);

  const handleSave = form.handleSubmit((values) => {
    setFeedback(null);
    updateConfig.mutate(values, {
      onSuccess: () => setFeedback("Token Saver & Headroom settings saved."),
    });
  });

  const handlePreview = () => {
    setPreviewErr(null);
    const parsed = parsePreviewBodyJson(previewJson);
    if (!parsed.ok) {
      setPreviewErr(parsed.error);
      return;
    }
    previewMutation.mutate({
      protocol: "openai-chat",
      model: "gpt-4o",
      strategy: form.getValues("strategy"),
      body: parsed.value,
    });
  };

  const stats = statsQuery.data;
  const hitRatePct = stats ? `${(stats.headroom.hit_rate * 100).toFixed(1)}%` : "0.0%";
  const activeError =
    configQuery.error?.message ??
    statsQuery.error?.message ??
    updateConfig.error?.message ??
    previewMutation.error?.message ??
    previewErr;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <p className="font-mono text-xs uppercase tracking-[0.18em] text-[var(--accent-primary)]">
            Gateway Interceptor · Priority 30
          </p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Token Saver & Headroom
          </h1>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Go-native token compression, Anthropic ephemeral cache injection, and bounded LRU prefix
            affinity tracking.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => {
              void Promise.all([configQuery.refetch(), statsQuery.refetch()]);
            }}
          >
            <RefreshCw className="mr-2 h-4 w-4" />
            Refresh
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => clearCache.mutate()}
            disabled={clearCache.isPending}
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Clear Cache
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => resetStats.mutate()}
            disabled={resetStats.isPending}
          >
            <RotateCcw className="mr-2 h-4 w-4" />
            Reset Counters
          </Button>
        </div>
      </div>

      <EmberRule />

      {activeError ? (
        <Alert variant="destructive">
          <Zap className="h-4 w-4" />
          <AlertTitle>Token Saver Error</AlertTitle>
          <AlertDescription>{activeError}</AlertDescription>
        </Alert>
      ) : null}

      {feedback ? (
        <Alert>
          <AlertTitle>Updated</AlertTitle>
          <AlertDescription>{feedback}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card className="p-4">
          <p className="text-xs text-[var(--text-secondary)]">Estimated Tokens Saved</p>
          <p className="mt-1 font-mono text-2xl font-semibold text-[var(--text-primary)]">
            {stats?.compression.estimated_tokens_saved ?? 0}
          </p>
          <p className="mt-1 font-mono text-xs text-[var(--text-secondary)]">
            {stats?.compression.bytes_saved ?? 0} bytes across{" "}
            {stats?.compression.requests_compressed ?? 0} reqs
          </p>
        </Card>
        <Card className="p-4">
          <p className="text-xs text-[var(--text-secondary)]">Headroom Hit Rate</p>
          <p className="mt-1 font-mono text-2xl font-semibold text-[var(--accent-primary)]">
            {hitRatePct}
          </p>
          <p className="mt-1 font-mono text-xs text-[var(--text-secondary)]">
            {stats?.headroom.cache_hits ?? 0} hits / {stats?.headroom.cache_misses ?? 0} misses
          </p>
        </Card>
        <Card className="p-4">
          <p className="text-xs text-[var(--text-secondary)]">Active LRU Prefixes</p>
          <p className="mt-1 font-mono text-2xl font-semibold text-[var(--text-primary)]">
            {stats?.headroom.active_entries ?? 0} / {stats?.headroom.max_entries ?? 1024}
          </p>
          <p className="mt-1 text-xs text-[var(--text-secondary)]">Bounded low-RAM cache</p>
        </Card>
        <Card className="p-4">
          <p className="text-xs text-[var(--text-secondary)]">Requests Processed</p>
          <p className="mt-1 font-mono text-2xl font-semibold text-[var(--text-primary)]">
            {stats?.compression.requests_processed ?? 0}
          </p>
          <p className="mt-1 text-xs text-[var(--text-secondary)]">Active interceptor counter</p>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="p-6">
          <CardHeader className="p-0 pb-4">
            <CardTitle className="text-base">Compression & Headroom Settings</CardTitle>
          </CardHeader>
          <form onSubmit={handleSave} className="space-y-4">
            <div className="flex flex-wrap gap-6">
              <label className="flex items-center gap-2 text-sm text-[var(--text-primary)]">
                <input
                  type="checkbox"
                  className="h-4 w-4 accent-[var(--accent-primary)]"
                  {...form.register("enabled")}
                />
                Enable Token Saver
              </label>
              <label className="flex items-center gap-2 text-sm text-[var(--text-primary)]">
                <input
                  type="checkbox"
                  className="h-4 w-4 accent-[var(--accent-primary)]"
                  {...form.register("headroom_enabled")}
                />
                Enable Headroom Cache
              </label>
              <label className="flex items-center gap-2 text-sm text-[var(--text-primary)]">
                <input
                  type="checkbox"
                  className="h-4 w-4 accent-[var(--accent-primary)]"
                  {...form.register("inject_cache_control")}
                />
                Inject Anthropic cache_control
              </label>
            </div>

            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="ts-strategy">Compression Strategy</Label>
                <select
                  id="ts-strategy"
                  className="h-10 w-full rounded-md border border-[var(--border-default)] bg-[var(--bg-muted)] px-3 text-sm text-[var(--text-primary)]"
                  {...form.register("strategy")}
                >
                  {TOKEN_SAVER_STRATEGIES.map((s) => (
                    <option key={s} value={s}>
                      {s}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1">
                <Label htmlFor="ts-depth">Prefix Depth (messages)</Label>
                <Input
                  id="ts-depth"
                  type="number"
                  {...form.register("prefix_depth", { valueAsNumber: true })}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="ts-history">Max History Messages</Label>
                <Input
                  id="ts-history"
                  type="number"
                  {...form.register("max_history_messages", { valueAsNumber: true })}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="ts-tool">Max Tool Output Chars</Label>
                <Input
                  id="ts-tool"
                  type="number"
                  {...form.register("max_tool_output_chars", { valueAsNumber: true })}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="ts-entries">Max LRU Cache Entries</Label>
                <Input
                  id="ts-entries"
                  type="number"
                  {...form.register("max_cache_entries", { valueAsNumber: true })}
                />
              </div>
              <div className="space-y-1">
                <Label htmlFor="ts-ttl">Cache TTL (seconds)</Label>
                <Input
                  id="ts-ttl"
                  type="number"
                  {...form.register("cache_ttl_seconds", { valueAsNumber: true })}
                />
              </div>
            </div>

            <Button type="submit" disabled={updateConfig.isPending}>
              {updateConfig.isPending ? "Saving..." : "Save Configuration"}
            </Button>
          </form>
        </Card>

        <Card className="p-6">
          <CardHeader className="p-0 pb-4">
            <CardTitle className="text-base">Dry-Run Compression Preview</CardTitle>
          </CardHeader>
          <div className="space-y-3">
            <Label htmlFor="ts-preview">Sample JSON Payload</Label>
            <textarea
              id="ts-preview"
              rows={5}
              value={previewJson}
              onChange={(e) => setPreviewJson(e.target.value)}
              className="w-full rounded-md border border-[var(--border-default)] bg-[var(--bg-muted)] p-3 font-mono text-xs text-[var(--text-primary)]"
            />
            <Button
              type="button"
              variant="outline"
              onClick={handlePreview}
              disabled={previewMutation.isPending}
            >
              <Play className="mr-2 h-4 w-4" />
              Test Compression
            </Button>
            {previewMutation.data ? (
              <div className="space-y-2 rounded-md border border-[var(--border-default)] bg-[var(--bg-muted)] p-3 font-mono text-xs">
                <div className="flex flex-wrap gap-2">
                  <Badge variant="secondary">
                    Saved {previewMutation.data.bytes_saved} B (~
                    {previewMutation.data.estimated_tokens_saved} tok)
                  </Badge>
                  <Badge variant="outline">Prefix: {previewMutation.data.cache_key}</Badge>
                </div>
                <pre className="max-h-40 overflow-auto text-[var(--text-primary)]">
                  {JSON.stringify(previewMutation.data.body, null, 2)}
                </pre>
              </div>
            ) : null}
          </div>
        </Card>
      </div>

      <Card className="p-6">
        <CardHeader className="p-0 pb-4">
          <CardTitle className="text-base">Warm Headroom Prefixes</CardTitle>
        </CardHeader>
        {(stats?.headroom.warm_prefixes ?? []).length === 0 ? (
          <p className="text-sm text-[var(--text-secondary)]">
            No warm prefix affinity entries yet. Send repeated prompts to populate Headroom cache.
          </p>
        ) : (
          <div className="space-y-2">
            {(stats?.headroom.warm_prefixes ?? []).map((p) => (
              <div
                key={p.key}
                className="flex flex-wrap items-center justify-between gap-2 rounded-md border border-[var(--border-default)] bg-[var(--bg-muted)] px-4 py-2.5 font-mono text-xs"
              >
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant="outline">{p.key}</Badge>
                  <span className="text-[var(--text-primary)]">{p.model}</span>
                  <span className="text-[var(--text-secondary)]">
                    {p.provider_id || "unbound"} / {p.connection_id || "-"}
                  </span>
                </div>
                <Badge variant="secondary">{p.hits} hit(s)</Badge>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
