"use client";

import {
  AlertTriangle,
  CheckCircle2,
  Clock,
  Database,
  Layers,
  RefreshCw,
  Server,
  ShieldCheck,
} from "lucide-react";
import Link from "next/link";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { useInitStatus, useSystemHealth, useSystemVersion } from "./hooks";

const PLANNED_WAVE1_MODULES = [
  { name: "auth", scope: "Single-admin session cookies & login policy" },
  { name: "settings", scope: "Typed settings store, history & rollback" },
  { name: "providers", scope: "Upstream LLM adapters, OAuth & credential vault" },
  { name: "keys", scope: "Client API keys, allowed models & quotas" },
  { name: "combos", scope: "Fallback chains, aliases & account cooldowns" },
  { name: "gateway", scope: "OpenAI, Anthropic, Responses & Gemini SSE proxy" },
  { name: "usage", scope: "Async token accounting, pricing & request telemetry" },
];

function formatUptime(totalSeconds: number): string {
  const safe = Math.max(0, Math.floor(totalSeconds));
  const days = Math.floor(safe / 86_400);
  const hours = Math.floor((safe % 86_400) / 3_600);
  const minutes = Math.floor((safe % 3_600) / 60);
  const seconds = safe % 60;

  const pad = (n: number) => String(n).padStart(2, "0");
  if (days > 0) {
    return `${days}d ${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
  }
  return `${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
}

function formatUpdatedAt(timestampMs: number): string {
  if (timestampMs <= 0) {
    return "Not synced";
  }
  return new Date(timestampMs).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

export function SystemOverview() {
  const healthQuery = useSystemHealth();
  const versionQuery = useSystemVersion();
  const initQuery = useInitStatus();

  const isRefreshing = healthQuery.isFetching || versionQuery.isFetching || initQuery.isFetching;

  const handleRefresh = () => {
    void Promise.all([healthQuery.refetch(), versionQuery.refetch(), initQuery.refetch()]);
  };

  const health = healthQuery.data;
  const version = versionQuery.data;
  const initStatus = initQuery.data;

  const isHealthy = health?.status === "ok";
  const isDegraded = health?.status === "degraded";
  const hasError = healthQuery.isError || versionQuery.isError;

  const mountedSet = new Set(health?.modules ?? []);
  const pendingWave1Modules = PLANNED_WAVE1_MODULES.filter((mod) => !mountedSet.has(mod.name));

  const lastSyncTimestamp = Math.max(healthQuery.dataUpdatedAt, versionQuery.dataUpdatedAt);

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            System Telemetry
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Runtime &amp; Control Plane Health
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Live status from <span className="font-mono">/api/health</span> and{" "}
            <span className="font-mono">/api/version</span>.
          </p>
        </div>

        <div className="flex items-center gap-3">
          <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
            Synced: {formatUpdatedAt(lastSyncTimestamp)}
          </span>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
          </Button>
        </div>
      </section>

      {initStatus && !initStatus.initialized ? (
        <Alert
          variant="warning"
          className="flex-col justify-between gap-3 px-4 py-3.5 sm:flex-row sm:items-center"
        >
          <div className="flex items-start gap-2.5">
            <AlertTriangle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
              aria-hidden="true"
            />
            <p className="text-sm text-[var(--text-primary)]">
              Administrator password has not been initialized yet. Complete first-run configuration
              before exposing this instance.
            </p>
          </div>
          <Button asChild size="xs" className="shrink-0">
            <Link href="/init/">Run Setup</Link>
          </Button>
        </Alert>
      ) : null}

      {hasError ? (
        <Alert
          variant="destructive"
          className="flex-col justify-between gap-3 px-4 py-3.5 sm:flex-row sm:items-center"
        >
          <div className="flex items-start gap-2.5">
            <AlertTriangle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <div className="text-sm">
              <AlertTitle>Control plane endpoint unreachable</AlertTitle>
              <AlertDescription className="font-mono">
                {healthQuery.error?.message ??
                  versionQuery.error?.message ??
                  "Failed to reach gateway backend"}
              </AlertDescription>
            </div>
          </div>
          <Button
            type="button"
            variant="destructive"
            size="xs"
            onClick={handleRefresh}
            className="shrink-0"
          >
            Retry
          </Button>
        </Alert>
      ) : null}

      <section
        aria-label="Key system metrics"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Gateway Status
            </span>
            {isHealthy ? (
              <CheckCircle2 className="h-4 w-4 text-[var(--status-ok)]" aria-hidden="true" />
            ) : (
              <AlertTriangle
                className={`h-4 w-4 ${
                  isDegraded ? "text-[var(--status-warn)]" : "text-[var(--text-muted)]"
                }`}
                aria-hidden="true"
              />
            )}
          </div>
          <div className="mt-3 flex items-baseline gap-2">
            <span
              className={`font-mono text-2xl font-semibold tracking-tight ${
                isHealthy
                  ? "text-[var(--status-ok)]"
                  : isDegraded
                    ? "text-[var(--status-warn)]"
                    : healthQuery.isError
                      ? "text-[var(--status-error)]"
                      : "text-[var(--text-secondary)]"
              }`}
            >
              {healthQuery.isPending
                ? "CHECKING"
                : health
                  ? health.status.toUpperCase()
                  : "OFFLINE"}
            </span>
          </div>
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">
            HTTP listener on local control plane
          </p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              SQLite Store
            </span>
            <Database className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <div className="mt-3">
            <span
              className={`font-mono text-2xl font-semibold tracking-tight ${
                health?.db === "ok"
                  ? "text-[var(--text-primary)]"
                  : health?.db === "error"
                    ? "text-[var(--status-error)]"
                    : "text-[var(--text-secondary)]"
              }`}
            >
              {healthQuery.isPending ? "CHECKING" : health ? health.db.toUpperCase() : "UNKNOWN"}
            </span>
          </div>
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">WAL mode embedded state engine</p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Process Uptime
            </span>
            <Clock className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <div className="mt-3">
            <span className="font-mono text-2xl font-semibold tracking-tight tabular-nums text-[var(--text-primary)]">
              {healthQuery.isPending || healthQuery.isError || !health
                ? "--:--:--"
                : formatUptime(health.uptimeSeconds)}
            </span>
          </div>
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">Elapsed since kernel start</p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Loaded Modules
            </span>
            <Layers className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
          </div>
          <div className="mt-3 flex items-baseline gap-2">
            <span className="font-mono text-2xl font-semibold tracking-tight tabular-nums text-[var(--text-primary)]">
              {health ? health.modules.length : 0}
            </span>
            <span className="text-xs text-[var(--text-muted)]">active</span>
          </div>
          <p className="mt-1.5 font-mono text-xs text-[var(--text-muted)]">
            {health && health.modules.length > 0
              ? health.modules.join(", ")
              : "Awaiting health probe"}
          </p>
        </Card>
      </section>

      <EmberRule />

      <section className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Server className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Registered Kernel Modules</CardTitle>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">GET /api/health</span>
          </CardHeader>

          <div className="divide-y divide-[var(--border-subtle)]">
            {health && health.modules.length > 0 ? (
              health.modules.map((mod) => (
                <div key={mod} className="flex items-center justify-between px-5 py-3">
                  <div>
                    <p className="font-mono text-sm font-medium text-[var(--text-primary)]">
                      {mod}
                    </p>
                    <p className="text-xs text-[var(--text-muted)]">
                      Active kernel control plane module
                    </p>
                  </div>
                  <Badge variant="success">Mounted</Badge>
                </div>
              ))
            ) : (
              <div className="px-5 py-4 text-sm text-[var(--text-muted)]">
                {healthQuery.isPending
                  ? "Querying kernel module registry..."
                  : "No active modules reported."}
              </div>
            )}

            {pendingWave1Modules.map((mod) => (
              <div
                key={mod.name}
                className="flex items-center justify-between bg-[var(--bg-canvas)] px-5 py-2.5"
              >
                <div>
                  <p className="font-mono text-xs font-medium text-[var(--text-secondary)]">
                    {mod.name}
                  </p>
                  <p className="text-xs text-[var(--text-muted)]">{mod.scope}</p>
                </div>
                <Badge variant="outline">Wave 1</Badge>
              </div>
            ))}
          </div>
        </Card>

        <Card className="flex flex-col">
          <CardHeader>
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Build &amp; Security Posture</CardTitle>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">GET /api/version</span>
          </CardHeader>

          <dl className="divide-y divide-[var(--border-subtle)] text-sm">
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Version</dt>
              <dd className="col-span-2 font-mono text-[var(--text-primary)]">
                {versionQuery.isPending ? "Loading..." : (version?.version ?? "Unavailable")}
              </dd>
            </div>
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Git Commit</dt>
              <dd className="col-span-2 truncate font-mono text-[var(--text-primary)]">
                {versionQuery.isPending ? "Loading..." : (version?.commit ?? "Unavailable")}
              </dd>
            </div>
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Build Time</dt>
              <dd className="col-span-2 font-mono text-[var(--text-primary)]">
                {versionQuery.isPending ? "Loading..." : (version?.buildTime ?? "Unavailable")}
              </dd>
            </div>
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Go Runtime</dt>
              <dd className="col-span-2 font-mono text-[var(--text-primary)]">
                {versionQuery.isPending ? "Loading..." : (version?.goVersion ?? "Unavailable")}
              </dd>
            </div>
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Initialized</dt>
              <dd className="col-span-2 font-mono text-[var(--text-primary)]">
                {initQuery.isPending
                  ? "Checking..."
                  : initStatus
                    ? initStatus.initialized
                      ? "true"
                      : "false"
                    : "Unavailable"}
              </dd>
            </div>
            <div className="grid grid-cols-3 gap-4 px-5 py-3">
              <dt className="text-[var(--text-muted)]">Login Required</dt>
              <dd className="col-span-2 font-mono text-[var(--text-primary)]">
                {initQuery.isPending
                  ? "Checking..."
                  : initStatus
                    ? initStatus.requireLogin
                      ? "true"
                      : "false"
                    : "Unavailable"}
              </dd>
            </div>
          </dl>
        </Card>
      </section>
    </div>
  );
}
