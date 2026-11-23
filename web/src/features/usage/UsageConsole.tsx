"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Activity,
  AlertCircle,
  BarChart3,
  CheckCircle2,
  Coins,
  Cpu,
  DollarSign,
  Gauge,
  Pencil,
  Radio,
  RefreshCw,
  RotateCcw,
  ScrollText,
  Server,
  ShieldAlert,
  X,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useDeleteModelPricingMutation,
  useModelPricing,
  useUpsertModelPricingMutation,
  useUsageLogs,
  useUsageStream,
  useUsageSummary,
} from "./hooks";
import {
  DEFAULT_LOGS_LIMIT,
  DEFAULT_UPSERT_PRICING_VALUES,
  formatCompactTokens,
  formatDurationMs,
  formatErrorRate,
  formatUSD,
  type ModelPricing,
  toUpsertModelPricingRequest,
  type UpsertPricingFormInput,
  type UpsertPricingFormOutput,
  USAGE_RANGE_DESCRIPTIONS,
  USAGE_RANGE_LABELS,
  USAGE_RANGES,
  type UsageLog,
  type UsageRange,
  upsertPricingSchema,
} from "./schemas";

const LOG_PAGE_SIZES = [25, 50, 100, 200] as const;

const STATUS_FILTER_OPTIONS = [
  { value: "", label: "All statuses" },
  { value: "success", label: "Success (2xx / 3xx)" },
  { value: "error", label: "Errors (4xx / 5xx)" },
  { value: "429", label: "HTTP 429 (Rate limit)" },
  { value: "500", label: "HTTP 500 (Server error)" },
  { value: "502", label: "HTTP 502 (Bad gateway)" },
] as const;

function formatTimestamp(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString([], {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function statusBadgeVariant(statusCode: number): "success" | "warning" | "destructive" {
  if (statusCode >= 200 && statusCode < 400) {
    return "success";
  }
  if (statusCode >= 400 && statusCode < 500) {
    return "warning";
  }
  return "destructive";
}

function streamBadgeLabel(status: "idle" | "connecting" | "live" | "error", eventCount: number) {
  switch (status) {
    case "connecting":
      return "Connecting SSE";
    case "live":
      return eventCount > 0 ? `Live (${eventCount})` : "Live SSE";
    case "error":
      return "SSE Reconnecting";
    default:
      return "Live Stream Off";
  }
}

export function UsageConsole() {
  const [range, setRange] = useState<UsageRange>("24h");
  const [liveStreamEnabled, setLiveStreamEnabled] = useState(false);

  const [modelDraft, setModelDraft] = useState("");
  const [providerDraft, setProviderDraft] = useState("");
  const [debouncedModel, setDebouncedModel] = useState("");
  const [debouncedProvider, setDebouncedProvider] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [limit, setLimit] = useState<number>(DEFAULT_LOGS_LIMIT);
  const [offset, setOffset] = useState(0);

  const [pricingSearch, setPricingSearch] = useState("");
  const [busyPricingModel, setBusyPricingModel] = useState<string | null>(null);
  const [confirmResetModel, setConfirmResetModel] = useState<string | null>(null);
  const [pricingFeedback, setPricingFeedback] = useState<string | null>(null);
  const [pricingRowError, setPricingRowError] = useState<string | null>(null);

  const pricingModelInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedModel(modelDraft.trim());
      setDebouncedProvider(providerDraft.trim());
      setOffset(0);
    }, 250);
    return () => window.clearTimeout(timer);
  }, [modelDraft, providerDraft]);

  const logsFilter = useMemo(
    () => ({
      model: debouncedModel || undefined,
      provider: debouncedProvider || undefined,
      status: statusFilter || undefined,
      limit,
      offset,
    }),
    [debouncedModel, debouncedProvider, statusFilter, limit, offset],
  );

  const summaryQuery = useUsageSummary(range);
  const logsQuery = useUsageLogs(logsFilter);
  const pricingQuery = useModelPricing();
  const upsertPricingMutation = useUpsertModelPricingMutation();
  const deletePricingMutation = useDeleteModelPricingMutation();
  const streamState = useUsageStream(liveStreamEnabled);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<UpsertPricingFormInput, unknown, UpsertPricingFormOutput>({
    resolver: zodResolver(upsertPricingSchema),
    defaultValues: DEFAULT_UPSERT_PRICING_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const modelField = register("model");

  const totals = summaryQuery.data?.totals;
  const byModel = summaryQuery.data?.by_model ?? [];
  const byProvider = summaryQuery.data?.by_provider ?? [];
  const totalTokens =
    (totals?.input_tokens ?? 0) + (totals?.output_tokens ?? 0) + (totals?.cached_tokens ?? 0);

  const maxModelRequests = useMemo(
    () => byModel.reduce((max, item) => Math.max(max, item.requests), 0),
    [byModel],
  );
  const maxProviderRequests = useMemo(
    () => byProvider.reduce((max, item) => Math.max(max, item.requests), 0),
    [byProvider],
  );

  const logs: UsageLog[] = logsQuery.data?.logs ?? [];
  const logsTotal = logsQuery.data?.total ?? 0;
  const hasActiveLogFilters =
    Boolean(modelDraft.trim()) || Boolean(providerDraft.trim()) || Boolean(statusFilter);

  const pricingCatalog: ModelPricing[] = pricingQuery.data ?? [];
  const filteredPricing = useMemo(() => {
    const query = pricingSearch.trim().toLowerCase();
    if (!query) {
      return pricingCatalog;
    }
    return pricingCatalog.filter((entry) => entry.model.toLowerCase().includes(query));
  }, [pricingCatalog, pricingSearch]);

  const customPricingCount = useMemo(
    () => pricingCatalog.filter((entry) => entry.source === "custom").length,
    [pricingCatalog],
  );

  const handleRefreshAll = () => {
    void summaryQuery.refetch();
    void logsQuery.refetch();
    void pricingQuery.refetch();
  };

  const isRefreshing = summaryQuery.isFetching || logsQuery.isFetching || pricingQuery.isFetching;

  const onSavePricingSubmit = (values: UpsertPricingFormOutput) => {
    if (busyPricingModel !== null) {
      return;
    }
    setPricingFeedback(null);
    setPricingRowError(null);
    setConfirmResetModel(null);
    setBusyPricingModel(values.model);

    upsertPricingMutation.mutate(
      {
        model: values.model,
        body: toUpsertModelPricingRequest(values),
      },
      {
        onSuccess: (saved) => {
          setPricingFeedback(
            `Saved custom pricing for ${saved.model} (${formatUSD(saved.input_per_million)} in / ${formatUSD(saved.output_per_million)} out / ${formatUSD(saved.cached_per_million)} cached per 1M).`,
          );
          reset(DEFAULT_UPSERT_PRICING_VALUES);
        },
        onError: (err) => {
          setPricingRowError(err.message);
        },
        onSettled: () => {
          setBusyPricingModel(null);
        },
      },
    );
  };

  const handlePopulatePricingForm = (entry: ModelPricing) => {
    setPricingFeedback(null);
    setPricingRowError(null);
    upsertPricingMutation.reset();
    setValue("model", entry.model, { shouldValidate: false });
    setValue("input_per_million", String(entry.input_per_million), { shouldValidate: false });
    setValue("output_per_million", String(entry.output_per_million), { shouldValidate: false });
    setValue("cached_per_million", String(entry.cached_per_million), { shouldValidate: false });
    pricingModelInputRef.current?.focus();
    pricingModelInputRef.current?.scrollIntoView({ behavior: "smooth", block: "center" });
  };

  const handleConfirmResetPricing = (model: string) => {
    if (busyPricingModel !== null) {
      return;
    }
    setPricingFeedback(null);
    setPricingRowError(null);
    setBusyPricingModel(model);

    deletePricingMutation.mutate(model, {
      onSuccess: () => {
        setConfirmResetModel(null);
        setPricingFeedback(
          `Removed custom pricing override for ${model}. Built-in defaults (if any) are now active.`,
        );
      },
      onError: (err) => {
        setPricingRowError(err.message);
      },
      onSettled: () => {
        setBusyPricingModel(null);
      },
    });
  };

  const showingFrom = logsTotal === 0 ? 0 : offset + 1;
  const showingTo = Math.min(offset + logs.length, logsTotal);
  const canPrevPage = offset > 0;
  const canNextPage = offset + limit < logsTotal;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Telemetry &amp; Cost Attribution
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Usage &amp; Pricing
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Monitor gateway request throughput, token volume, latency, live logs, and per-million
            model pricing rates.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <fieldset
            aria-label="Summary time range"
            className="inline-flex rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] p-0.5"
          >
            {USAGE_RANGES.map((item) => {
              const active = range === item;
              return (
                <button
                  key={item}
                  type="button"
                  aria-pressed={active}
                  onClick={() => setRange(item)}
                  className={
                    active
                      ? "rounded px-3 py-1.5 font-mono text-xs font-semibold bg-[#E56A4A] text-[#0a0a0a] transition-colors"
                      : "rounded px-3 py-1.5 font-mono text-xs font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors"
                  }
                >
                  {USAGE_RANGE_LABELS[item]}
                </button>
              );
            })}
          </fieldset>

          <Button
            type="button"
            variant={liveStreamEnabled ? "default" : "outline"}
            size="sm"
            aria-pressed={liveStreamEnabled}
            onClick={() => setLiveStreamEnabled((prev) => !prev)}
          >
            <Radio className="h-3.5 w-3.5" aria-hidden="true" />
            <span className="tabular-nums">
              {streamBadgeLabel(streamState.status, streamState.eventCount)}
            </span>
          </Button>

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

      {summaryQuery.isError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div>
            <AlertTitle>Failed to load usage summary</AlertTitle>
            <AlertDescription className="font-mono">{summaryQuery.error.message}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      <section
        aria-label="Usage key performance indicators"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5"
      >
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                Requests
              </span>
              <Activity className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            </div>
            <p className="mt-2 font-mono text-2xl font-semibold tracking-tight text-[var(--text-primary)] tabular-nums">
              {summaryQuery.isPending ? "--" : (totals?.requests ?? 0).toLocaleString("en-US")}
            </p>
            <p className="mt-1 font-mono text-xs text-[var(--text-muted)] tabular-nums">
              {summaryQuery.isPending
                ? "Loading window..."
                : `${totals?.error_count ?? 0} failed · ${USAGE_RANGE_DESCRIPTIONS[range]}`}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                Total Tokens
              </span>
              <Cpu className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            </div>
            <p className="mt-2 font-mono text-2xl font-semibold tracking-tight text-[var(--text-primary)] tabular-nums">
              {summaryQuery.isPending ? "--" : formatCompactTokens(totalTokens)}
            </p>
            <p className="mt-1 font-mono text-xs text-[var(--text-muted)] tabular-nums">
              {summaryQuery.isPending
                ? "In / Out / Cached"
                : `In ${formatCompactTokens(totals?.input_tokens ?? 0)} · Out ${formatCompactTokens(totals?.output_tokens ?? 0)} · Cache ${formatCompactTokens(totals?.cached_tokens ?? 0)}`}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                Est. Cost USD
              </span>
              <DollarSign className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            </div>
            <p className="mt-2 font-mono text-2xl font-semibold tracking-tight text-[var(--text-primary)] tabular-nums">
              {summaryQuery.isPending ? "--" : formatUSD(totals?.cost_usd ?? 0)}
            </p>
            <p className="mt-1 font-mono text-xs text-[var(--text-muted)]">
              Attributed via active rates
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                Avg Latency / TTFT
              </span>
              <Gauge className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            </div>
            <p className="mt-2 font-mono text-2xl font-semibold tracking-tight text-[var(--text-primary)] tabular-nums">
              {summaryQuery.isPending
                ? "--"
                : `${formatDurationMs(totals?.avg_latency_ms ?? 0)} / ${formatDurationMs(totals?.avg_ttft_ms ?? 0)}`}
            </p>
            <p className="mt-1 font-mono text-xs text-[var(--text-muted)]">
              Total duration / First token
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                Error Rate
              </span>
              <ShieldAlert className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            </div>
            <p
              className={`mt-2 font-mono text-2xl font-semibold tracking-tight tabular-nums ${
                (totals?.error_rate ?? 0) > 0.05
                  ? "text-[var(--status-error)]"
                  : "text-[var(--text-primary)]"
              }`}
            >
              {summaryQuery.isPending ? "--" : formatErrorRate(totals?.error_rate ?? 0)}
            </p>
            <p className="mt-1 font-mono text-xs text-[var(--text-muted)] tabular-nums">
              {(totals?.error_count ?? 0).toLocaleString("en-US")} non-2xx/3xx responses
            </p>
          </CardContent>
        </Card>
      </section>

      <section className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <BarChart3 className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Top Models</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Request share, token volume, and cost breakdown by model
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)] tabular-nums">
              {byModel.length} models
            </span>
          </CardHeader>
          <CardContent>
            {summaryQuery.isPending ? (
              <p className="text-sm text-[var(--text-muted)]">Loading model breakdown...</p>
            ) : byModel.length === 0 ? (
              <p className="py-6 text-center text-sm text-[var(--text-muted)]">
                No model telemetry recorded in {USAGE_RANGE_DESCRIPTIONS[range].toLowerCase()}.
              </p>
            ) : (
              <div className="space-y-4">
                {byModel.map((item) => {
                  const totalModelTokens =
                    item.input_tokens + item.output_tokens + item.cached_tokens;
                  const barPct =
                    maxModelRequests > 0
                      ? Math.max(4, Math.round((item.requests / maxModelRequests) * 100))
                      : 0;
                  const sharePct =
                    totals && totals.requests > 0
                      ? Math.round((item.requests / totals.requests) * 100)
                      : 0;

                  return (
                    <div key={item.model} className="space-y-1.5">
                      <div className="flex flex-wrap items-baseline justify-between gap-2">
                        <span className="font-mono text-xs font-medium text-[var(--text-primary)]">
                          {item.model}
                        </span>
                        <div className="flex flex-wrap items-center gap-3 font-mono text-xs text-[var(--text-secondary)] tabular-nums">
                          <span>
                            {item.requests.toLocaleString("en-US")} req ({sharePct}%)
                          </span>
                          <span>{formatCompactTokens(totalModelTokens)} tok</span>
                          <span className="text-[var(--text-primary)]">
                            {formatUSD(item.cost_usd)}
                          </span>
                        </div>
                      </div>
                      <div className="h-2 w-full overflow-hidden rounded-sm border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
                        <div
                          className="h-full bg-[#E56A4A] transition-all"
                          style={{ width: `${barPct}%` }}
                        />
                      </div>
                      <div className="flex items-center justify-between font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                        <span>Avg latency {formatDurationMs(item.avg_latency_ms)}</span>
                        {item.error_count > 0 ? (
                          <span className="text-[var(--status-error)]">
                            {item.error_count} errors
                          </span>
                        ) : (
                          <span>0 errors</span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div>
              <div className="flex items-center gap-2">
                <Server className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                <CardTitle>Top Providers</CardTitle>
              </div>
              <CardDescription className="mt-0.5">
                Upstream provider distribution, cost, and average latency
              </CardDescription>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)] tabular-nums">
              {byProvider.length} providers
            </span>
          </CardHeader>
          <CardContent>
            {summaryQuery.isPending ? (
              <p className="text-sm text-[var(--text-muted)]">Loading provider breakdown...</p>
            ) : byProvider.length === 0 ? (
              <p className="py-6 text-center text-sm text-[var(--text-muted)]">
                No provider telemetry recorded in {USAGE_RANGE_DESCRIPTIONS[range].toLowerCase()}.
              </p>
            ) : (
              <div className="space-y-4">
                {byProvider.map((item) => {
                  const totalProviderTokens =
                    item.input_tokens + item.output_tokens + item.cached_tokens;
                  const barPct =
                    maxProviderRequests > 0
                      ? Math.max(4, Math.round((item.requests / maxProviderRequests) * 100))
                      : 0;
                  const sharePct =
                    totals && totals.requests > 0
                      ? Math.round((item.requests / totals.requests) * 100)
                      : 0;

                  return (
                    <div key={item.provider_id} className="space-y-1.5">
                      <div className="flex flex-wrap items-baseline justify-between gap-2">
                        <span className="font-mono text-xs font-medium text-[var(--text-primary)]">
                          {item.provider_id}
                        </span>
                        <div className="flex flex-wrap items-center gap-3 font-mono text-xs text-[var(--text-secondary)] tabular-nums">
                          <span>
                            {item.requests.toLocaleString("en-US")} req ({sharePct}%)
                          </span>
                          <span>{formatCompactTokens(totalProviderTokens)} tok</span>
                          <span className="text-[var(--text-primary)]">
                            {formatUSD(item.cost_usd)}
                          </span>
                        </div>
                      </div>
                      <div className="h-2 w-full overflow-hidden rounded-sm border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
                        <div
                          className="h-full bg-[#E56A4A] transition-all"
                          style={{ width: `${barPct}%` }}
                        />
                      </div>
                      <div className="flex items-center justify-between font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                        <span>Avg latency {formatDurationMs(item.avg_latency_ms)}</span>
                        {item.error_count > 0 ? (
                          <span className="text-[var(--status-error)]">
                            {item.error_count} errors
                          </span>
                        ) : (
                          <span>0 errors</span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </CardContent>
        </Card>
      </section>

      <EmberRule />

      <Card>
        <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center">
          <div>
            <div className="flex items-center gap-2">
              <ScrollText className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Request Logs</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Metadata-first gateway request history ordered newest-first
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">GET /api/usage/logs</span>
        </CardHeader>

        <div className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3.5">
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <Label htmlFor="logs-filter-model">Model</Label>
              <Input
                id="logs-filter-model"
                type="text"
                value={modelDraft}
                onChange={(e) => setModelDraft(e.target.value)}
                placeholder="e.g. openai/gpt-4o"
                className="mt-1.5 h-9 text-xs"
              />
            </div>

            <div>
              <Label htmlFor="logs-filter-provider">Provider</Label>
              <Input
                id="logs-filter-provider"
                type="text"
                value={providerDraft}
                onChange={(e) => setProviderDraft(e.target.value)}
                placeholder="e.g. openai"
                className="mt-1.5 h-9 text-xs"
              />
            </div>

            <div>
              <Label htmlFor="logs-filter-status">Status</Label>
              <select
                id="logs-filter-status"
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value);
                  setOffset(0);
                }}
                className="mt-1.5 flex h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-1.5 font-mono text-xs text-[var(--text-primary)] transition-colors focus:border-[#E56A4A]"
              >
                {STATUS_FILTER_OPTIONS.map((opt) => (
                  <option key={opt.value || "all"} value={opt.value}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex items-end gap-2">
              <div className="flex-1">
                <Label htmlFor="logs-page-size">Page Size</Label>
                <select
                  id="logs-page-size"
                  value={limit}
                  onChange={(e) => {
                    setLimit(Number(e.target.value));
                    setOffset(0);
                  }}
                  className="mt-1.5 flex h-9 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-1.5 font-mono text-xs text-[var(--text-primary)] tabular-nums transition-colors focus:border-[#E56A4A]"
                >
                  {LOG_PAGE_SIZES.map((size) => (
                    <option key={size} value={size}>
                      {size} rows
                    </option>
                  ))}
                </select>
              </div>
              {hasActiveLogFilters ? (
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    setModelDraft("");
                    setProviderDraft("");
                    setDebouncedModel("");
                    setDebouncedProvider("");
                    setStatusFilter("");
                    setOffset(0);
                  }}
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Clear</span>
                </Button>
              ) : null}
            </div>
          </div>
        </div>

        {logsQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load request logs</AlertTitle>
                <AlertDescription className="font-mono">{logsQuery.error.message}</AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : logsQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading request logs...
          </CardContent>
        ) : logs.length === 0 ? (
          <CardContent className="py-10 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              No matching request logs found
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              {hasActiveLogFilters
                ? "Try clearing or loosening your model, provider, or status filters."
                : "Send a request through /v1/chat/completions, /v1/responses, or /v1/messages to record telemetry."}
            </p>
          </CardContent>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full border-collapse text-left text-sm">
                <thead>
                  <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                    <th scope="col" className="px-5 py-3">
                      Timestamp
                    </th>
                    <th scope="col" className="px-4 py-3">
                      Model &amp; Provider
                    </th>
                    <th scope="col" className="px-4 py-3">
                      Tokens (In / Out / Cache)
                    </th>
                    <th scope="col" className="px-4 py-3">
                      Cost
                    </th>
                    <th scope="col" className="px-4 py-3">
                      Latency / TTFT
                    </th>
                    <th scope="col" className="px-5 py-3 text-right">
                      Status
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[var(--border-subtle)]">
                  {logs.map((entry) => (
                    <tr
                      key={entry.id}
                      className="align-top transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3 font-mono text-xs text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
                        <div>{formatTimestamp(entry.created_at)}</div>
                        <div className="mt-0.5 text-[11px] text-[var(--text-muted)]">
                          {entry.request_id}
                        </div>
                      </td>

                      <td className="px-4 py-3">
                        <div className="flex flex-wrap items-center gap-1.5">
                          <span className="font-mono text-xs font-medium text-[var(--text-primary)]">
                            {entry.model}
                          </span>
                          {entry.stream ? <Badge variant="default">SSE</Badge> : null}
                        </div>
                        <div className="mt-1 flex flex-wrap items-center gap-2 font-mono text-[11px] text-[var(--text-muted)]">
                          <span>{entry.provider_id}</span>
                          <span>·</span>
                          <span>{entry.protocol}</span>
                          {entry.key_id ? (
                            <>
                              <span>·</span>
                              <span>key:{entry.key_id}</span>
                            </>
                          ) : null}
                        </div>
                      </td>

                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-primary)] tabular-nums whitespace-nowrap">
                        {entry.input_tokens.toLocaleString("en-US")} /{" "}
                        {entry.output_tokens.toLocaleString("en-US")} /{" "}
                        <span className="text-[var(--text-muted)]">
                          {entry.cached_tokens.toLocaleString("en-US")}
                        </span>
                      </td>

                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-primary)] tabular-nums whitespace-nowrap">
                        {formatUSD(entry.cost_usd)}
                      </td>

                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-secondary)] tabular-nums whitespace-nowrap">
                        {formatDurationMs(entry.latency_ms)} /{" "}
                        <span className="text-[var(--text-muted)]">
                          {entry.ttft_ms > 0 ? formatDurationMs(entry.ttft_ms) : "--"}
                        </span>
                      </td>

                      <td className="px-5 py-3 text-right">
                        <div className="flex flex-col items-end gap-1">
                          <Badge
                            variant={statusBadgeVariant(entry.status_code)}
                            className="tabular-nums"
                          >
                            HTTP {entry.status_code}
                          </Badge>
                          {entry.error_summary ? (
                            <span
                              className="max-w-xs truncate font-mono text-[11px] text-[var(--status-error)]"
                              title={entry.error_summary}
                            >
                              {entry.error_summary}
                            </span>
                          ) : null}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="flex flex-col gap-3 border-t border-[var(--border-subtle)] px-5 py-3 sm:flex-row sm:items-center sm:justify-between">
              <span className="font-mono text-xs text-[var(--text-muted)] tabular-nums">
                Showing {showingFrom.toLocaleString("en-US")}–{showingTo.toLocaleString("en-US")} of{" "}
                {logsTotal.toLocaleString("en-US")} logs
              </span>
              <div className="flex items-center gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="xs"
                  disabled={!canPrevPage || logsQuery.isFetching}
                  onClick={() => setOffset((prev) => Math.max(0, prev - limit))}
                >
                  Previous
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  size="xs"
                  disabled={!canNextPage || logsQuery.isFetching}
                  onClick={() => setOffset((prev) => prev + limit)}
                >
                  Next
                </Button>
              </div>
            </div>
          </>
        )}
      </Card>

      <EmberRule />

      <Card>
        <CardHeader>
          <div>
            <div className="flex items-center gap-2">
              <Coins className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Model Pricing Editor</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Override USD per-million-token rates for any model identifier, or edit an existing
              catalog entry below.
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            PUT /api/usage/pricing/{"{model}"}
          </span>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSavePricingSubmit)} noValidate className="space-y-4">
            {pricingFeedback ? (
              <Alert>
                <CheckCircle2
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Pricing updated</AlertTitle>
                  <AlertDescription className="font-mono">{pricingFeedback}</AlertDescription>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setPricingFeedback(null)}
                  aria-label="Dismiss pricing notice"
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </Alert>
            ) : null}

            {pricingRowError ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <AlertTitle>Pricing operation failed</AlertTitle>
                  <AlertDescription className="font-mono">{pricingRowError}</AlertDescription>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setPricingRowError(null)}
                  aria-label="Dismiss pricing error"
                >
                  <X className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </Alert>
            ) : null}

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div>
                <Label htmlFor="pricing-model">Model Identifier</Label>
                <Input
                  id="pricing-model"
                  type="text"
                  placeholder="e.g. openai/gpt-4o"
                  aria-invalid={Boolean(errors.model)}
                  className="mt-1.5"
                  {...modelField}
                  ref={(node) => {
                    modelField.ref(node);
                    pricingModelInputRef.current = node;
                  }}
                />
                {errors.model ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">{errors.model.message}</p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="pricing-input-rate">Input USD / 1M Tokens</Label>
                <Input
                  id="pricing-input-rate"
                  type="number"
                  min={0}
                  step="any"
                  placeholder="2.50"
                  aria-invalid={Boolean(errors.input_per_million)}
                  className="mt-1.5 tabular-nums"
                  {...register("input_per_million")}
                />
                {errors.input_per_million ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.input_per_million.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="pricing-output-rate">Output USD / 1M Tokens</Label>
                <Input
                  id="pricing-output-rate"
                  type="number"
                  min={0}
                  step="any"
                  placeholder="10.00"
                  aria-invalid={Boolean(errors.output_per_million)}
                  className="mt-1.5 tabular-nums"
                  {...register("output_per_million")}
                />
                {errors.output_per_million ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.output_per_million.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="pricing-cached-rate">Cached USD / 1M Tokens</Label>
                <Input
                  id="pricing-cached-rate"
                  type="number"
                  min={0}
                  step="any"
                  placeholder="0"
                  aria-invalid={Boolean(errors.cached_per_million)}
                  className="mt-1.5 tabular-nums"
                  {...register("cached_per_million")}
                />
                {errors.cached_per_million ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {errors.cached_per_million.message}
                  </p>
                ) : null}
              </div>
            </div>

            <div className="flex items-center justify-end gap-3 pt-1">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                disabled={busyPricingModel !== null}
                onClick={() => {
                  upsertPricingMutation.reset();
                  setPricingRowError(null);
                  reset(DEFAULT_UPSERT_PRICING_VALUES);
                }}
              >
                Clear Form
              </Button>
              <Button type="submit" size="sm" disabled={busyPricingModel !== null}>
                <DollarSign className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{upsertPricingMutation.isPending ? "Saving..." : "Save Custom Pricing"}</span>
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-col items-start gap-3 sm:flex-row sm:items-center">
          <div>
            <div className="flex items-center gap-2">
              <Coins className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              <CardTitle>Active Model Pricing Catalog</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Built-in default rates merged with persisted custom overrides
            </CardDescription>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Badge variant="outline" className="tabular-nums">
              {pricingCatalog.length} Models
            </Badge>
            <Badge variant="default" className="tabular-nums">
              {customPricingCount} Custom
            </Badge>
          </div>
        </CardHeader>

        <div className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3">
          <div className="max-w-sm">
            <Label htmlFor="pricing-catalog-search">Filter Pricing Catalog</Label>
            <Input
              id="pricing-catalog-search"
              type="text"
              value={pricingSearch}
              onChange={(e) => setPricingSearch(e.target.value)}
              placeholder="Search model name..."
              className="mt-1.5 h-9 text-xs"
            />
          </div>
        </div>

        {pricingQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load model pricing catalog</AlertTitle>
                <AlertDescription className="font-mono">
                  {pricingQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : pricingQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading model pricing catalog...
          </CardContent>
        ) : filteredPricing.length === 0 ? (
          <CardContent className="py-10 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              No model pricing entries match your filter
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              Clear the catalog filter or add a custom rate using the editor above.
            </p>
          </CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <th scope="col" className="px-5 py-3">
                    Model
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Source
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Input / 1M
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Output / 1M
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Cached / 1M
                  </th>
                  <th scope="col" className="px-5 py-3 text-right">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {filteredPricing.map((entry) => {
                  const isRowBusy = busyPricingModel === entry.model;
                  const anyRowBusy = busyPricingModel !== null;
                  const confirmingReset = confirmResetModel === entry.model;

                  return (
                    <tr
                      key={entry.model}
                      className="align-middle transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3 font-mono text-xs font-medium text-[var(--text-primary)]">
                        {entry.model}
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant={entry.source === "custom" ? "default" : "secondary"}>
                          {entry.source}
                        </Badge>
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-primary)] tabular-nums">
                        {formatUSD(entry.input_per_million)}
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-primary)] tabular-nums">
                        {formatUSD(entry.output_per_million)}
                      </td>
                      <td className="px-4 py-3 font-mono text-xs text-[var(--text-secondary)] tabular-nums">
                        {formatUSD(entry.cached_per_million)}
                      </td>
                      <td className="px-5 py-3 text-right">
                        {confirmingReset ? (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="destructive"
                              size="xs"
                              disabled={anyRowBusy}
                              onClick={() => handleConfirmResetPricing(entry.model)}
                            >
                              {isRowBusy ? "Resetting..." : "Confirm Reset"}
                            </Button>
                            <Button
                              type="button"
                              variant="secondary"
                              size="xs"
                              disabled={anyRowBusy}
                              onClick={() => setConfirmResetModel(null)}
                            >
                              Cancel
                            </Button>
                          </div>
                        ) : (
                          <div className="flex items-center justify-end gap-1.5">
                            <Button
                              type="button"
                              variant="outline"
                              size="xs"
                              disabled={anyRowBusy}
                              onClick={() => handlePopulatePricingForm(entry)}
                            >
                              <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
                              <span>Edit</span>
                            </Button>
                            {entry.source === "custom" ? (
                              <Button
                                type="button"
                                variant="destructive"
                                size="xs"
                                disabled={anyRowBusy}
                                onClick={() => setConfirmResetModel(entry.model)}
                              >
                                <RotateCcw className="h-3.5 w-3.5" aria-hidden="true" />
                                <span>Reset</span>
                              </Button>
                            ) : null}
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
