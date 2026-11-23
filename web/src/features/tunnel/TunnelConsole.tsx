"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Play, Radio, RefreshCw, Save, Square } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useStartTunnelMutation,
  useStopTunnelMutation,
  useTunnelEvents,
  useTunnelStatus,
  useUpdateTunnelConfigMutation,
} from "./hooks";
import {
  TAILSCALE_MODES,
  TUNNEL_PROVIDERS,
  type TunnelConfigFormInput,
  type TunnelConfigFormOutput,
  tunnelConfigSchema,
} from "./schemas";

export function TunnelConsole() {
  const statusQuery = useTunnelStatus();
  const eventsQuery = useTunnelEvents();
  const updateConfig = useUpdateTunnelConfigMutation();
  const startTunnel = useStartTunnelMutation();
  const stopTunnel = useStopTunnelMutation();

  const form = useForm<TunnelConfigFormInput, unknown, TunnelConfigFormOutput>({
    resolver: zodResolver(tunnelConfigSchema),
    defaultValues: {
      provider: "cloudflare",
      local_url: "http://127.0.0.1:20128",
      tailscale_mode: "funnel",
      auto_start: false,
    },
  });

  const status = statusQuery.data;
  const events = eventsQuery.data?.events ?? [];

  useEffect(() => {
    if (!status || form.formState.isDirty) {
      return;
    }
    form.reset({
      provider: status.provider,
      local_url: status.local_url,
      tailscale_mode: status.tailscale_mode,
      auto_start: status.auto_start,
    });
  }, [status, form]);

  const activeError =
    statusQuery.error ??
    eventsQuery.error ??
    updateConfig.error ??
    startTunnel.error ??
    stopTunnel.error;

  const isRunning = status?.status === "running" || status?.status === "starting";

  const clearMutationErrors = () => {
    updateConfig.reset();
    startTunnel.reset();
    stopTunnel.reset();
  };

  const onSaveConfig = form.handleSubmit((values) => {
    clearMutationErrors();
    updateConfig.mutate(values, {
      onSuccess: (updated) => {
        form.reset({
          provider: updated.provider,
          local_url: updated.local_url,
          tailscale_mode: updated.tailscale_mode,
          auto_start: updated.auto_start,
        });
      },
    });
  });

  const handleStart = form.handleSubmit((values) => {
    clearMutationErrors();
    startTunnel.mutate({
      provider: values.provider,
      local_url: values.local_url,
      tailscale_mode: values.tailscale_mode,
    });
  });

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Wave 3B — Remote Exposure
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Cloudflare Quick Tunnel &amp; Tailscale
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Expose this loopback router safely via <span className="font-mono">cloudflared</span> or{" "}
            <span className="font-mono">tailscale serve/funnel</span> with lifecycle auditing.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="secondary"
            size="sm"
            onClick={() => {
              void Promise.all([statusQuery.refetch(), eventsQuery.refetch()]);
            }}
            disabled={statusQuery.isFetching || eventsQuery.isFetching}
          >
            <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Refresh</span>
          </Button>
          {isRunning ? (
            <Button
              type="button"
              variant="destructive"
              size="sm"
              onClick={() => {
                clearMutationErrors();
                stopTunnel.mutate();
              }}
              disabled={stopTunnel.isPending}
            >
              <Square className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{stopTunnel.isPending ? "Stopping..." : "Stop Tunnel"}</span>
            </Button>
          ) : (
            <Button type="button" size="sm" onClick={handleStart} disabled={startTunnel.isPending}>
              <Play className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{startTunnel.isPending ? "Starting..." : "Start Tunnel"}</span>
            </Button>
          )}
        </div>
      </section>

      <EmberRule />

      {activeError ? (
        <Alert variant="destructive">
          <AlertTitle>Tunnel Error</AlertTitle>
          <AlertDescription>{activeError.message}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Live State</CardDescription>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Radio className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <span className="uppercase">{status?.status ?? "loading"}</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="font-mono text-xs text-[var(--text-muted)]">
            Provider: {status?.provider ?? "-"} {status?.pid ? `| PID ${status.pid}` : ""}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Public URL</CardDescription>
            <CardTitle className="truncate font-mono text-sm">
              {status?.public_url ? (
                <a
                  href={status.public_url}
                  target="_blank"
                  rel="noreferrer"
                  className="text-[var(--accent-primary)] underline-offset-4 hover:underline"
                >
                  {status.public_url}
                </a>
              ) : (
                <span className="text-[var(--text-muted)]">Not active</span>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent className="font-mono text-xs text-[var(--text-muted)]">
            Target: {status?.local_url ?? "-"}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Host CLI Binaries</CardDescription>
            <CardTitle className="flex flex-wrap gap-2 text-sm">
              <Badge variant={status?.binaries.cloudflared.available ? "default" : "outline"}>
                cloudflared: {status?.binaries.cloudflared.available ? "ready" : "missing"}
              </Badge>
              <Badge variant={status?.binaries.tailscale.available ? "default" : "outline"}>
                tailscale: {status?.binaries.tailscale.available ? "ready" : "missing"}
              </Badge>
            </CardTitle>
          </CardHeader>
          <CardContent className="truncate font-mono text-xs text-[var(--text-muted)]">
            {status?.last_error ? `Last error: ${status.last_error}` : "No runtime errors"}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Tunnel Configuration</CardTitle>
          <CardDescription>
            Local target URL is strictly restricted to loopback hosts (127.0.0.1, localhost, ::1).
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSaveConfig} className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="space-y-1.5">
                <Label htmlFor="tunnel-provider">Provider</Label>
                <select
                  id="tunnel-provider"
                  className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
                  {...form.register("provider")}
                >
                  {TUNNEL_PROVIDERS.map((p) => (
                    <option key={p} value={p}>
                      {p}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="tunnel-ts-mode">Tailscale Mode</Label>
                <select
                  id="tunnel-ts-mode"
                  className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
                  {...form.register("tailscale_mode")}
                >
                  {TAILSCALE_MODES.map((m) => (
                    <option key={m} value={m}>
                      {m}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="tunnel-local-url">Loopback Local URL</Label>
                <Input
                  id="tunnel-local-url"
                  placeholder="http://127.0.0.1:20128"
                  {...form.register("local_url")}
                />
                {form.formState.errors.local_url ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.local_url.message}
                  </p>
                ) : null}
              </div>
            </div>

            <div className="flex items-center gap-2">
              <input
                id="tunnel-auto-start"
                type="checkbox"
                className="h-4 w-4 rounded border-[var(--border-subtle)]"
                {...form.register("auto_start")}
              />
              <Label htmlFor="tunnel-auto-start">Auto-start tunnel when router boots</Label>
            </div>

            <Button type="submit" size="sm" disabled={updateConfig.isPending}>
              <Save className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{updateConfig.isPending ? "Saving..." : "Save Configuration"}</span>
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Lifecycle Events ({events.length})</CardTitle>
        </CardHeader>
        <CardContent>
          {events.length === 0 ? (
            <p className="text-sm text-[var(--text-muted)]">No tunnel events recorded yet.</p>
          ) : (
            <div className="space-y-2">
              {events.map((ev) => (
                <div
                  key={ev.id}
                  className="flex flex-col gap-1 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-2.5 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge variant="outline" className="font-mono uppercase">
                      {ev.provider}
                    </Badge>
                    <Badge variant="secondary" className="font-mono">
                      {ev.event_type}
                    </Badge>
                    <span className="text-sm text-[var(--text-primary)]">{ev.message}</span>
                    {ev.public_url ? (
                      <span className="font-mono text-xs text-[var(--accent-primary)]">
                        {ev.public_url}
                      </span>
                    ) : null}
                  </div>
                  <span className="font-mono text-xs text-[var(--text-muted)]">
                    {ev.created_at}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
