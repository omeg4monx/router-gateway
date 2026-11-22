"use client";

import Link from "next/link";
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
import { EmberRule } from "@/shared/ui/EmberRule";
import {
  useInitStatus,
  useSystemHealth,
  useSystemVersion,
} from "./hooks";

const PLANNED_WAVE1_MODULES = [
  { name: "auth", scope: "Session cookies & RBAC" },
  { name: "providers", scope: "Upstream LLM connections & OAuth" },
  { name: "keys", scope: "Client API keys & scopes" },
  { name: "router", scope: "OpenAI & Anthropic gateway translation" },
  { name: "usage", scope: "Token accounting & cost analytics" },
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

  const isRefreshing =
    healthQuery.isFetching || versionQuery.isFetching || initQuery.isFetching;

  const handleRefresh = () => {
    void Promise.all([
      healthQuery.refetch(),
      versionQuery.refetch(),
      initQuery.refetch(),
    ]);
  };

  const health = healthQuery.data;
  const version = versionQuery.data;
  const initStatus = initQuery.data;

  const isHealthy = health?.status === "ok";
  const isDegraded = health?.status === "degraded";
  const hasError = healthQuery.isError || versionQuery.isError;

  const lastSyncTimestamp = Math.max(
    healthQuery.dataUpdatedAt,
    versionQuery.dataUpdatedAt,
  );

  return (
    <div className="space-y-6">
      {/* Top control strip */}
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[#E56A4A]">
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
          <button
            type="button"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="inline-flex h-9 items-center gap-2 rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-3.5 text-xs font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A] disabled:opacity-60"
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
          </button>
        </div>
      </section>

      {/* First-run initialization reminder if uninitialized */}
      {initStatus && !initStatus.initialized ? (
        <div
          role="status"
          className="flex flex-col gap-3 rounded-md border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-4 py-3.5 sm:flex-row sm:items-center sm:justify-between"
        >
          <div className="flex items-start gap-2.5">
            <AlertTriangle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
              aria-hidden="true"
            />
            <p className="text-sm text-[var(--text-primary)]">
              Administrator password has not been initialized yet. Complete first-run
              configuration before exposing this instance.
            </p>
          </div>
          <Link
            href="/init/"
            className="inline-flex shrink-0 items-center justify-center rounded-md bg-[#E56A4A] px-3.5 py-2 text-xs font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)]"
          >
            Run Setup
          </Link>
        </div>
      ) : null}

      {/* API unreachable banner */}
      {hasError ? (
        <div
          role="alert"
          className="flex flex-col gap-3 rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-4 py-3.5 sm:flex-row sm:items-center sm:justify-between"
        >
          <div className="flex items-start gap-2.5">
            <AlertTriangle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <div className="text-sm">
              <p className="font-medium text-[var(--text-primary)]">
                Control plane endpoint unreachable
              </p>
              <p className="mt-0.5 font-mono text-xs text-[var(--text-secondary)]">
                {healthQuery.error?.message ??
                  versionQuery.error?.message ??
                  "Failed to reach gateway backend"}
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={handleRefresh}
            className="inline-flex shrink-0 items-center justify-center rounded-md border border-[var(--status-error-border)] bg-[var(--bg-surface)] px-3 py-1.5 text-xs font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
          >
            Retry
          </button>
        </div>
      ) : null}

      {/* Primary KPI metric cards */}
      <section
        aria-label="Key system metrics"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Gateway Status
            </span>
            {isHealthy ? (
              <CheckCircle2
                className="h-4 w-4 text-[var(--status-ok)]"
                aria-hidden="true"
              />
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
        </div>

        <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-4">
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
              {healthQuery.isPending
                ? "CHECKING"
                : health
                  ? health.db.toUpperCase()
                  : "UNKNOWN"}
            </span>
          </div>
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">
            WAL mode embedded state engine
          </p>
        </div>

        <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Process Uptime
            </span>
            <Clock className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <div className="mt-3">
            <span className="font-mono text-2xl font-semibold tracking-tight tabular-nums text-[var(--text-primary)]">
              {healthQuery.isPending
                ? "--:--:--"
                : health
                  ? formatUptime(health.uptimeSeconds)
                  : "00h 00m 00s"}
            </span>
          </div>
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">
            Elapsed since kernel start
          </p>
        </div>

        <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-4">
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
        </div>
      </section>

      <EmberRule />

      {/* Technical detail panels */}
      <section className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Kernel modules registry */}
        <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)]">
          <div className="flex items-center justify-between border-b border-[var(--border-subtle)] px-5 py-3.5">
            <div className="flex items-center gap-2">
              <Server className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <h3 className="text-sm font-semibold text-[var(--text-primary)]">
                Registered Kernel Modules
              </h3>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              GET /api/health
            </span>
          </div>

          <div className="divide-y divide-[var(--border-subtle)]">
            {health && health.modules.length > 0 ? (
              health.modules.map((mod) => (
                <div
                  key={mod}
                  className="flex items-center justify-between px-5 py-3"
                >
                  <div>
                    <p className="font-mono text-sm font-medium text-[var(--text-primary)]">
                      {mod}
                    </p>
                    <p className="text-xs text-[var(--text-muted)]">
                      Active kernel control plane module
                    </p>
                  </div>
                  <span className="rounded border border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] px-2 py-0.5 font-mono text-[11px] font-medium uppercase text-[var(--status-ok)]">
                    Mounted
                  </span>
                </div>
              ))
            ) : (
              <div className="px-5 py-4 text-sm text-[var(--text-muted)]">
                {healthQuery.isPending
                  ? "Querying kernel module registry..."
                  : "No active modules reported."}
              </div>
            )}

            {PLANNED_WAVE1_MODULES.map((mod) => (
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
                <span className="rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
                  Wave 1
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Build & runtime metadata */}
        <div className="flex flex-col rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)]">
          <div className="flex items-center justify-between border-b border-[var(--border-subtle)] px-5 py-3.5">
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <h3 className="text-sm font-semibold text-[var(--text-primary)]">
                Build &amp; Security Posture
              </h3>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              GET /api/version
            </span>
          </div>

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
        </div>
      </section>
    </div>
  );
}
