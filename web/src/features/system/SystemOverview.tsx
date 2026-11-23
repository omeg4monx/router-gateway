"use client";

import {
  Activity,
  AlertTriangle,
  ArrowRight,
  Database,
  Flame,
  KeyRound,
  Layers,
  Network,
  RefreshCw,
  Server,
  Sparkles,
} from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";
import { useCombos } from "@/features/combos/hooks";
import { useApiKeys } from "@/features/keys/hooks";
import { useProviderConnections } from "@/features/providers/hooks";
import { useProxyPools } from "@/features/proxypools/hooks";
import { useUsageLogs, useUsageSummary } from "@/features/usage/hooks";
import {
  formatCompactTokens,
  formatDurationMs,
  formatUSD,
  USAGE_RANGES,
  type UsageRange,
} from "@/features/usage/schemas";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { useSystemHealth } from "./hooks";
import { formatClockUptime } from "./overview";

function formatBucketLabel(bucket: string): string {
  const parsed = new Date(bucket);
  if (Number.isNaN(parsed.getTime())) {
    return bucket.slice(0, 10);
  }
  return parsed.toLocaleDateString([], { month: "short", day: "numeric" });
}

export function SystemOverview() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const [range, setRange] = useState<UsageRange>("7d");

  const healthQuery = useSystemHealth();
  const providersQuery = useProviderConnections();
  const keysQuery = useApiKeys();
  const combosQuery = useCombos();
  const proxyPoolsQuery = useProxyPools();
  const usageQuery = useUsageSummary(range);
  const recentLogsQuery = useUsageLogs({ limit: 100 });

  const isRefreshing =
    healthQuery.isFetching ||
    providersQuery.isFetching ||
    keysQuery.isFetching ||
    combosQuery.isFetching ||
    proxyPoolsQuery.isFetching ||
    usageQuery.isFetching ||
    recentLogsQuery.isFetching;

  const handleRefresh = () => {
    void Promise.all([
      healthQuery.refetch(),
      providersQuery.refetch(),
      keysQuery.refetch(),
      combosQuery.refetch(),
      proxyPoolsQuery.refetch(),
      usageQuery.refetch(),
      recentLogsQuery.refetch(),
    ]);
  };

  const health = healthQuery.data;
  const providers = providersQuery.data ?? [];
  const keys = keysQuery.data ?? [];
  const combos = combosQuery.data ?? [];
  const proxyPools = proxyPoolsQuery.data?.pools ?? [];
  const summary = usageQuery.data;
  const totals = summary?.totals;
  const series = summary?.series ?? [];
  const logs = recentLogsQuery.data?.logs ?? [];

  const activeProvidersCount = providers.filter((p) => p.enabled).length;
  const activeKeysCount = keys.filter((k) => !k.disabled).length;
  const activeProxyPoolsCount = proxyPools.filter((p) => p.enabled).length;

  const isHealthy = health?.status === "ok";
  const isDegraded = health?.status === "degraded";
  const hasError = healthQuery.isError;

  const isFreshInstall =
    !providersQuery.isPending &&
    !keysQuery.isPending &&
    providers.length === 0 &&
    keys.length === 0;

  // Map model -> provider_id from recent logs & provider connections
  const modelProviderMap = useMemo(() => {
    const map = new Map<string, string>();
    for (const log of logs) {
      if (log.model && log.provider_id && !map.has(log.model)) {
        map.set(log.model, log.provider_id);
      }
    }
    for (const conn of providers) {
      for (const model of conn.models) {
        if (model && !map.has(model)) {
          map.set(model, conn.provider_id);
        }
      }
    }
    return map;
  }, [logs, providers]);

  const topModels = useMemo(
    () => [...(summary?.by_model ?? [])].sort((a, b) => b.requests - a.requests),
    [summary?.by_model],
  );

  const topProviders = useMemo(
    () => [...(summary?.by_provider ?? [])].sort((a, b) => b.requests - a.requests),
    [summary?.by_provider],
  );

  const totalRequests = totals?.requests ?? 0;
  const totalTokens = totals
    ? totals.input_tokens + totals.output_tokens + totals.cached_tokens
    : 0;
  const cachedTokens = totals?.cached_tokens ?? 0;
  const cacheHitPct =
    totals && totals.input_tokens + totals.cached_tokens > 0
      ? Math.round((totals.cached_tokens / (totals.input_tokens + totals.cached_tokens)) * 100)
      : 0;

  // Spotlight #1: Top Model (with provider)
  const mostUsedModel = useMemo(() => {
    const first = topModels[0];
    if (!first) return null;
    const sharePct = totalRequests > 0 ? Math.round((first.requests / totalRequests) * 100) : null;
    return {
      model: first.model,
      provider: modelProviderMap.get(first.model) ?? null,
      requests: first.requests,
      tokens: first.input_tokens + first.output_tokens + first.cached_tokens,
      sharePct,
    };
  }, [topModels, totalRequests, modelProviderMap]);

  // Spotlight #2: Top Provider
  const mostUsedProvider = useMemo(() => {
    const first = topProviders[0];
    if (!first) return null;
    const sharePct = totalRequests > 0 ? Math.round((first.requests / totalRequests) * 100) : null;
    return {
      name: first.provider_id,
      requests: first.requests,
      tokens: first.input_tokens + first.output_tokens + first.cached_tokens,
      sharePct,
    };
  }, [topProviders, totalRequests]);

  // Spotlight #3: Most Used API Key (from recent logs + key name lookup)
  const mostUsedKey = useMemo(() => {
    if (logs.length === 0) return null;
    const keyNameMap = new Map(keys.map((k) => [k.id, k.name]));
    const counts = new Map<string, { requests: number; tokens: number; provider: string | null }>();
    for (const log of logs) {
      if (!log.key_id) continue;
      const prev = counts.get(log.key_id) ?? {
        requests: 0,
        tokens: 0,
        provider: log.provider_id || null,
      };
      prev.requests += 1;
      prev.tokens += log.input_tokens + log.output_tokens + log.cached_tokens;
      counts.set(log.key_id, prev);
    }
    let bestId: string | null = null;
    let best = { requests: 0, tokens: 0, provider: null as string | null };
    for (const [id, data] of counts) {
      if (data.requests > best.requests) {
        bestId = id;
        best = data;
      }
    }
    if (!bestId) return null;
    const sharePct = logs.length > 0 ? Math.round((best.requests / logs.length) * 100) : null;
    return {
      name: keyNameMap.get(bestId) ?? bestId,
      requests: best.requests,
      tokens: best.tokens,
      provider: best.provider,
      sharePct,
    };
  }, [logs, keys]);

  // Spotlight #4: Provider with Most Connections
  const providerMostConnections = useMemo(() => {
    if (providers.length === 0) return null;
    const grouped = new Map<string, { total: number; active: number }>();
    for (const conn of providers) {
      const prev = grouped.get(conn.provider_id) ?? { total: 0, active: 0 };
      prev.total += 1;
      if (conn.enabled) prev.active += 1;
      grouped.set(conn.provider_id, prev);
    }
    let bestName: string | null = null;
    let best = { total: 0, active: 0 };
    for (const [name, data] of grouped) {
      if (data.total > best.total) {
        bestName = name;
        best = data;
      }
    }
    if (!bestName) return null;
    return {
      name: bestName,
      total: best.total,
      active: best.active,
      sharePct: Math.round((best.total / providers.length) * 100),
    };
  }, [providers]);

  const maxSeriesRequests = useMemo(() => Math.max(1, ...series.map((b) => b.requests)), [series]);
  const maxModelRequests = topModels[0]?.requests ?? 1;
  const maxProviderRequests = topProviders[0]?.requests ?? 1;

  const rangeLabels: Record<UsageRange, string> = {
    "24h": isId ? "24 jam" : "24h",
    "7d": isId ? "7 hari" : "7d",
    "30d": isId ? "30 hari" : "30d",
    all: isId ? "Semua" : "All",
  };

  return (
    <div className="flex min-w-0 flex-col gap-5 pb-10">
      {/* Header */}
      <header className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <div className="flex items-center gap-2.5">
            <h2 className="text-xl font-bold tracking-tight text-[var(--color-text-main)] sm:text-2xl">
              Dashboard
            </h2>
            <Badge
              variant={
                isHealthy
                  ? "success"
                  : isDegraded
                    ? "warning"
                    : hasError
                      ? "destructive"
                      : "secondary"
              }
            >
              {isHealthy
                ? isId
                  ? "Aktif"
                  : "Online"
                : isDegraded
                  ? isId
                    ? "Terdegradasi"
                    : "Degraded"
                  : hasError
                    ? "Offline"
                    : isId
                      ? "Memeriksa"
                      : "Checking"}
            </Badge>
          </div>
          <p className="mt-0.5 text-xs text-[var(--color-text-muted)]">
            {isId
              ? `Ringkasan pemakaian, model terpopuler, dan kesehatan gateway (${rangeLabels[range]})`
              : `Usage overview, most popular models, and gateway health (${rangeLabels[range]})`}
            {health ? ` · Uptime ${formatClockUptime(health.uptimeSeconds)}` : ""}
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <div
            className="inline-flex h-8 items-center rounded-lg border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-0.5"
            role="tablist"
            aria-label={isId ? "Rentang waktu" : "Time range"}
          >
            {USAGE_RANGES.map((item) => {
              const active = range === item;
              return (
                <button
                  key={item}
                  type="button"
                  role="tab"
                  aria-selected={active}
                  onClick={() => setRange(item)}
                  className={cn(
                    "h-6.5 cursor-pointer rounded-md px-2.5 font-mono text-[11px] font-medium transition-colors",
                    active
                      ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                      : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                  )}
                >
                  {rangeLabels[item]}
                </button>
              );
            })}
          </div>

          <Button
            type="button"
            variant="secondary"
            size="xs"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="h-8 border-[var(--color-border)]/80 bg-[var(--color-surface)] text-[var(--color-text-main)] hover:bg-[var(--color-surface-2)]"
          >
            <RefreshCw
              className={cn(
                "h-3.5 w-3.5 text-[var(--color-text-muted)]",
                isRefreshing && "animate-spin text-[var(--color-primary)]",
              )}
              aria-hidden="true"
            />
            <span>
              {isRefreshing ? (isId ? "Memuat" : "Syncing") : isId ? "Segarkan" : "Refresh"}
            </span>
          </Button>
        </div>
      </header>

      {/* Error banner */}
      {hasError ? (
        <div className="flex items-center gap-2.5 rounded-xl border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-4 py-3 text-xs text-[var(--status-error)]">
          <AlertTriangle className="h-4 w-4 shrink-0" aria-hidden="true" />
          <span>
            {isId
              ? "Gagal memuat sebagian data monitoring dari backend gateway."
              : "Failed to load monitoring data from the gateway backend."}
          </span>
        </div>
      ) : null}

      {/* Fresh install 3-step onboarding (disappears once first provider or key exists) */}
      {isFreshInstall ? (
        <section className="rounded-xl border border-[var(--color-primary)]/30 bg-[var(--color-primary)]/6 p-4 sm:p-5">
          <div className="flex items-center gap-2 text-sm font-semibold text-[var(--color-text-main)]">
            <Sparkles className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span>
              {isId ? "Langkah Pertama Menyiapkan Router" : "First Steps to Start Routing"}
            </span>
          </div>
          <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
            {[
              {
                step: "1",
                href: "/dashboard/providers/",
                title: isId ? "Hubungkan Provider" : "Connect a Provider",
                sub: isId
                  ? "OpenAI, Anthropic, Gemini, atau OAuth"
                  : "OpenAI, Anthropic, Gemini, or OAuth",
              },
              {
                step: "2",
                href: "/dashboard/keys/",
                title: isId ? "Buat Kunci API" : "Create an API Key",
                sub: isId ? "Atur batas kuota & cakupan model" : "Set rate limits & model scopes",
              },
              {
                step: "3",
                href: "/dashboard/chat/",
                title: isId ? "Uji di Obrolan Dasar" : "Test in Basic Chat",
                sub: isId
                  ? "Verifikasi routing langsung di studio"
                  : "Verify live routing in the studio",
              },
            ].map((item) => (
              <Link
                key={item.step}
                href={item.href}
                className="group flex items-start gap-3 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface)] p-3.5 transition-colors hover:border-[var(--color-primary)]/40"
              >
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-[var(--color-primary)] font-mono text-xs font-bold text-[var(--color-on-primary)]">
                  {item.step}
                </span>
                <div className="min-w-0">
                  <p className="text-xs font-semibold text-[var(--color-text-main)] group-hover:text-[var(--color-primary)]">
                    {item.title}
                  </p>
                  <p className="mt-0.5 text-[11px] text-[var(--color-text-muted)]">{item.sub}</p>
                </div>
              </Link>
            ))}
          </div>
        </section>
      ) : null}

      {/* Row 1: 4 Overview KPI Cards */}
      <section className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <div className="flex flex-col justify-between gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <span className="text-xs font-medium text-[var(--color-text-muted)]">
            {isId ? "Total Request" : "Total Requests"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {usageQuery.isPending ? "—" : totalRequests.toLocaleString()}
          </span>
          <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
            {isId ? "Latensi rata-rata" : "Avg latency"}:{" "}
            {totals ? formatDurationMs(totals.avg_latency_ms) : "—"}
          </span>
        </div>

        <div className="flex flex-col justify-between gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <span className="text-xs font-medium text-[var(--color-text-muted)]">
            {isId ? "Total Token" : "Total Tokens"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {usageQuery.isPending ? "—" : formatCompactTokens(totalTokens)}
          </span>
          <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
            In {formatCompactTokens(totals?.input_tokens ?? 0)} · Out{" "}
            {formatCompactTokens(totals?.output_tokens ?? 0)}
          </span>
        </div>

        <div className="flex flex-col justify-between gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <span className="text-xs font-medium text-[var(--color-text-muted)]">
            {isId ? "Estimasi Biaya" : "Estimated Spend"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {usageQuery.isPending ? "—" : formatUSD(totals?.cost_usd ?? 0)}
          </span>
          <span className="text-[11px] text-[var(--color-text-muted)]">
            {rangeLabels[range]} · {topModels.length} {isId ? "model aktif" : "active models"}
          </span>
        </div>

        <div className="flex flex-col justify-between gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <span className="text-xs font-medium text-[var(--color-text-muted)]">
            {isId ? "Token Ter-Cache" : "Cached Tokens"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {usageQuery.isPending ? "—" : formatCompactTokens(cachedTokens)}
          </span>
          <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
            {cacheHitPct}% {isId ? "dari input token" : "of input tokens"}
          </span>
        </div>
      </section>

      {/* Row 2: 4 Spotlight Insight Cards (9router MonitoringPageClient parity) */}
      <section className="grid grid-cols-1 gap-3 md:grid-cols-2">
        {/* Spotlight 1: Most Used Model + Provider */}
        <div className="flex flex-col gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex items-center gap-2">
            <Flame className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "Model Paling Sering Dipakai" : "Most Used Model"}
            </span>
          </div>
          {mostUsedModel ? (
            <>
              <div className="flex flex-wrap items-center gap-2">
                <span className="truncate font-mono text-lg font-bold tracking-tight text-[var(--color-text-main)]">
                  {mostUsedModel.model}
                </span>
                {mostUsedModel.provider ? (
                  <Badge variant="default">{mostUsedModel.provider}</Badge>
                ) : null}
              </div>
              <span className="text-xs text-[var(--color-text-muted)]">
                <strong className="font-semibold text-[var(--color-text-main)]">
                  {mostUsedModel.requests.toLocaleString()} request
                </strong>
                {mostUsedModel.sharePct !== null
                  ? ` · ${mostUsedModel.sharePct}% ${isId ? "dari total" : "of total"}`
                  : ""}{" "}
                · {formatCompactTokens(mostUsedModel.tokens)} {isId ? "token" : "tokens"} ·{" "}
                {rangeLabels[range]}
              </span>
            </>
          ) : (
            <span className="text-xs text-[var(--color-text-muted)]">
              {usageQuery.isPending
                ? isId
                  ? "Memuat…"
                  : "Loading…"
                : isId
                  ? `Belum ada pemakaian model (${rangeLabels[range]}).`
                  : `No model usage in the last ${rangeLabels[range]}.`}
            </span>
          )}
        </div>

        {/* Spotlight 2: Most Used Provider */}
        <div className="flex flex-col gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex items-center gap-2">
            <Server className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "Provider Paling Sering Dipakai" : "Most Used Provider"}
            </span>
          </div>
          {mostUsedProvider ? (
            <>
              <span className="truncate font-mono text-lg font-bold tracking-tight text-[var(--color-text-main)]">
                {mostUsedProvider.name}
              </span>
              <span className="text-xs text-[var(--color-text-muted)]">
                <strong className="font-semibold text-[var(--color-text-main)]">
                  {mostUsedProvider.requests.toLocaleString()} request
                </strong>
                {mostUsedProvider.sharePct !== null
                  ? ` · ${mostUsedProvider.sharePct}% ${isId ? "dari total" : "of total"}`
                  : ""}{" "}
                · {formatCompactTokens(mostUsedProvider.tokens)} {isId ? "token" : "tokens"} ·{" "}
                {rangeLabels[range]}
              </span>
            </>
          ) : (
            <span className="text-xs text-[var(--color-text-muted)]">
              {usageQuery.isPending
                ? isId
                  ? "Memuat…"
                  : "Loading…"
                : isId
                  ? `Belum ada pemakaian provider (${rangeLabels[range]}).`
                  : `No provider usage in the last ${rangeLabels[range]}.`}
            </span>
          )}
        </div>

        {/* Spotlight 3: Most Used API Key */}
        <div className="flex flex-col gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex items-center gap-2">
            <KeyRound className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "API Key Paling Sering Dipakai" : "Most Active API Key"}
            </span>
          </div>
          {mostUsedKey ? (
            <>
              <span className="truncate text-lg font-bold tracking-tight text-[var(--color-text-main)]">
                {mostUsedKey.name}
              </span>
              <span className="text-xs text-[var(--color-text-muted)]">
                <strong className="font-semibold text-[var(--color-text-main)]">
                  {mostUsedKey.requests.toLocaleString()} request
                </strong>
                {mostUsedKey.sharePct !== null ? ` · ${mostUsedKey.sharePct}%` : ""} ·{" "}
                {formatCompactTokens(mostUsedKey.tokens)} {isId ? "token" : "tokens"}
                {mostUsedKey.provider ? ` · via ${mostUsedKey.provider}` : ""}
              </span>
            </>
          ) : (
            <span className="text-xs text-[var(--color-text-muted)]">
              {recentLogsQuery.isPending
                ? isId
                  ? "Memuat…"
                  : "Loading…"
                : isId
                  ? "Belum ada pemakaian API key tercatat."
                  : "No API key activity recorded yet."}
            </span>
          )}
        </div>

        {/* Spotlight 4: Provider with Most Connections */}
        <div className="flex flex-col gap-1.5 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex items-center gap-2">
            <Database className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "Provider dengan Koneksi Terbanyak" : "Provider with Most Connections"}
            </span>
          </div>
          {providerMostConnections ? (
            <>
              <span className="truncate font-mono text-lg font-bold tracking-tight text-[var(--color-text-main)]">
                {providerMostConnections.name}
              </span>
              <span className="text-xs text-[var(--color-text-muted)]">
                <strong className="font-semibold text-[var(--color-text-main)]">
                  {providerMostConnections.total} {isId ? "koneksi" : "connections"}
                </strong>{" "}
                · {providerMostConnections.active} {isId ? "aktif" : "active"} ·{" "}
                {providerMostConnections.sharePct}% {isId ? "dari total" : "of total"}
              </span>
            </>
          ) : (
            <span className="text-xs text-[var(--color-text-muted)]">
              {providersQuery.isPending
                ? isId
                  ? "Memuat…"
                  : "Loading…"
                : isId
                  ? "Belum ada koneksi provider."
                  : "No provider connections yet."}
            </span>
          )}
        </div>
      </section>

      {/* Row 3: 4 Inventory MiniCards (9router MonitoringPageClient parity) */}
      <section className="grid grid-cols-2 gap-3 lg:grid-cols-4">
        <Link
          href="/dashboard/providers/"
          className="group flex flex-col gap-1 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 transition-colors hover:border-[var(--color-primary)]/35"
        >
          <Server className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
          <span className="mt-1 text-xs text-[var(--color-text-muted)]">
            {isId ? "Total Provider" : "Total Providers"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {providersQuery.isPending ? "—" : providers.length}
          </span>
          <span className="text-xs text-[var(--color-text-muted)]">
            {activeProvidersCount} {isId ? "aktif" : "active"}
          </span>
        </Link>

        <Link
          href="/dashboard/proxy-pools/"
          className="group flex flex-col gap-1 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 transition-colors hover:border-[var(--color-primary)]/35"
        >
          <Network className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
          <span className="mt-1 text-xs text-[var(--color-text-muted)]">
            {isId ? "Total Proxy Pool" : "Proxy Pools"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {proxyPoolsQuery.isPending ? "—" : proxyPools.length}
          </span>
          <span className="text-xs text-[var(--color-text-muted)]">
            {activeProxyPoolsCount} {isId ? "aktif" : "active"}
          </span>
        </Link>

        <Link
          href="/dashboard/models/"
          className="group flex flex-col gap-1 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 transition-colors hover:border-[var(--color-primary)]/35"
        >
          <Layers className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
          <span className="mt-1 text-xs text-[var(--color-text-muted)]">
            {isId ? "Total Combo" : "Saved Combos"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {combosQuery.isPending ? "—" : combos.length}
          </span>
          <span className="text-xs text-[var(--color-text-muted)]">
            {isId ? "kombinasi tersimpan" : "routing chains"}
          </span>
        </Link>

        <Link
          href="/dashboard/keys/"
          className="group flex flex-col gap-1 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 transition-colors hover:border-[var(--color-primary)]/35"
        >
          <KeyRound className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
          <span className="mt-1 text-xs text-[var(--color-text-muted)]">
            {isId ? "Total API Key" : "Total API Keys"}
          </span>
          <span className="font-mono text-2xl font-bold tracking-tight tabular-nums text-[var(--color-text-main)]">
            {keysQuery.isPending ? "—" : keys.length}
          </span>
          <span className="text-xs text-[var(--color-text-muted)]">
            {activeKeysCount} {isId ? "aktif" : "active"}
          </span>
        </Link>
      </section>

      {/* Row 4: Usage Activity Bar Chart */}
      <section className="overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
        <div className="flex h-11 items-center justify-between border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4">
          <div className="flex items-center gap-2">
            <Activity className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <h3 className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId
                ? `Tren Permintaan (${rangeLabels[range]})`
                : `Request Trend (${rangeLabels[range]})`}
            </h3>
          </div>
          <Link
            href="/dashboard/usage/"
            className="inline-flex items-center gap-1 text-xs font-medium text-[var(--color-primary)] hover:underline"
          >
            <span>{isId ? "Analitik Lengkap" : "Full Analytics"}</span>
            <ArrowRight className="h-3 w-3" aria-hidden="true" />
          </Link>
        </div>

        <div className="p-4">
          {series.length === 0 ? (
            <div className="flex h-36 items-center justify-center text-xs text-[var(--color-text-muted)]">
              {usageQuery.isPending
                ? isId
                  ? "Memuat grafik…"
                  : "Loading chart…"
                : isId
                  ? "Belum ada data seri waktu pada rentang ini."
                  : "No time-series data for this period yet."}
            </div>
          ) : (
            <div className="flex h-36 items-end gap-1.5 pt-2">
              {series.map((bucket) => {
                const heightPct = Math.max(
                  6,
                  Math.round((bucket.requests / maxSeriesRequests) * 100),
                );
                const bucketTokens =
                  bucket.input_tokens + bucket.output_tokens + bucket.cached_tokens;
                return (
                  <div
                    key={bucket.bucket}
                    className="group flex h-full min-w-0 flex-1 flex-col items-center justify-end gap-1.5"
                    title={`${formatBucketLabel(bucket.bucket)}: ${bucket.requests} req · ${formatCompactTokens(bucketTokens)} tokens · ${formatUSD(bucket.cost_usd)}`}
                  >
                    <div className="flex w-full flex-1 items-end">
                      <div
                        className="w-full rounded-t-md bg-[var(--color-primary)]/75 transition-all group-hover:bg-[var(--color-primary)]"
                        style={{ height: `${heightPct}%` }}
                      />
                    </div>
                    <span className="truncate font-mono text-[9.5px] text-[var(--color-text-subtle)]">
                      {formatBucketLabel(bucket.bucket)}
                    </span>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </section>

      {/* Row 5: Top Models & Top Providers Charts */}
      <section className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Top Models Breakdown */}
        <div className="overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex h-11 items-center justify-between border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4">
            <h3 className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "Model Terpopuler" : "Top Models"}
            </h3>
            <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
              {topModels.length}
            </span>
          </div>
          <div className="max-h-72 space-y-2.5 overflow-y-auto p-4">
            {topModels.length === 0 ? (
              <p className="py-6 text-center text-xs text-[var(--color-text-muted)]">
                {isId ? "Belum ada data model." : "No model usage data yet."}
              </p>
            ) : (
              topModels.map((item) => {
                const barPct = Math.max(4, Math.round((item.requests / maxModelRequests) * 100));
                const sharePct =
                  totalRequests > 0 ? Math.round((item.requests / totalRequests) * 100) : 0;
                const providerName = modelProviderMap.get(item.model);
                const tokens = item.input_tokens + item.output_tokens + item.cached_tokens;

                return (
                  <div key={item.model} className="space-y-1.5">
                    <div className="flex items-center justify-between gap-2 text-xs">
                      <div className="flex min-w-0 items-center gap-1.5">
                        <span className="truncate font-mono font-semibold text-[var(--color-text-main)]">
                          {item.model}
                        </span>
                        {providerName ? (
                          <Badge variant="outline" className="shrink-0">
                            {providerName}
                          </Badge>
                        ) : null}
                      </div>
                      <span className="shrink-0 font-mono tabular-nums text-[var(--color-text-muted)]">
                        <strong className="text-[var(--color-text-main)]">
                          {item.requests.toLocaleString()}
                        </strong>{" "}
                        ({sharePct}%)
                      </span>
                    </div>
                    <div className="h-1.5 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                      <div
                        className="h-full rounded-full bg-[var(--color-primary)]"
                        style={{ width: `${barPct}%` }}
                      />
                    </div>
                    <div className="flex items-center justify-between font-mono text-[10.5px] tabular-nums text-[var(--color-text-subtle)]">
                      <span>
                        {formatCompactTokens(tokens)} tok · {formatDurationMs(item.avg_latency_ms)}
                      </span>
                      <span>{formatUSD(item.cost_usd)}</span>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>

        {/* Top Providers Breakdown */}
        <div className="overflow-hidden rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.03)]">
          <div className="flex h-11 items-center justify-between border-b border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 px-4">
            <h3 className="text-xs font-semibold text-[var(--color-text-main)]">
              {isId ? "Distribusi Provider" : "Top Providers"}
            </h3>
            <span className="font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
              {topProviders.length}
            </span>
          </div>
          <div className="max-h-72 space-y-2.5 overflow-y-auto p-4">
            {topProviders.length === 0 ? (
              <p className="py-6 text-center text-xs text-[var(--color-text-muted)]">
                {isId ? "Belum ada data provider." : "No provider usage data yet."}
              </p>
            ) : (
              topProviders.map((prov) => {
                const barPct = Math.max(4, Math.round((prov.requests / maxProviderRequests) * 100));
                const sharePct =
                  totalRequests > 0 ? Math.round((prov.requests / totalRequests) * 100) : 0;
                const tokens = prov.input_tokens + prov.output_tokens + prov.cached_tokens;

                return (
                  <div key={prov.provider_id} className="space-y-1.5">
                    <div className="flex items-center justify-between gap-2 text-xs">
                      <span className="truncate font-mono font-semibold text-[var(--color-text-main)]">
                        {prov.provider_id}
                      </span>
                      <span className="shrink-0 font-mono tabular-nums text-[var(--color-text-muted)]">
                        <strong className="text-[var(--color-text-main)]">
                          {prov.requests.toLocaleString()}
                        </strong>{" "}
                        ({sharePct}%)
                      </span>
                    </div>
                    <div className="h-1.5 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                      <div
                        className="h-full rounded-full bg-[var(--color-primary)]"
                        style={{ width: `${barPct}%` }}
                      />
                    </div>
                    <div className="flex items-center justify-between font-mono text-[10.5px] tabular-nums text-[var(--color-text-subtle)]">
                      <span>
                        {formatCompactTokens(tokens)} tok · {formatDurationMs(prov.avg_latency_ms)}
                      </span>
                      <span>{formatUSD(prov.cost_usd)}</span>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      </section>
    </div>
  );
}
