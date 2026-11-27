import {
  AlertCircle,
  Check,
  Code2,
  Copy,
  Eye,
  EyeOff,
  KeyRound,
  Pencil,
  Plus,
  RotateCw,
  Search,
  Trash2,
  Waypoints,
  X,
} from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useProviderConnections } from "@/features/providers/hooks";
import {
  useStartTunnelMutation,
  useStopTunnelMutation,
  useTunnelStatus,
} from "@/features/tunnel/hooks";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import {
  countKeyStatusBuckets,
  filterKeysByStatusAndSearch,
  type GatewayRouteId,
  getKeyExpiryState,
  isKeyExpiredAt,
  type KeyStatusFilter,
  normalizeOrigin,
  resolveGatewayOrigins,
} from "./endpointUtils";
import { GatewayEndpointBar } from "./GatewayEndpointBar";
import {
  useApiKeys,
  useCreateApiKeyMutation,
  useDeleteApiKeyMutation,
  useRotateApiKeyMutation,
  useUpdateApiKeyMutation,
} from "./hooks";
import { KeyFormModal, KeyReadyModal, KeySnippetModal } from "./KeyModals";
import { type ApiKeySummary, type CreateKeyFormOutput, toCreateApiKeyRequest } from "./schemas";

interface SecretRevealState {
  action: "created" | "rotated";
  keyName: string;
  secret: string;
}

type PendingRowConfirm = { id: string; kind: "rotate" } | { id: string; kind: "delete" } | null;

const SKELETON_CARDS = ["sk-1", "sk-2"] as const;

function formatCompactLastUsed(iso: string | null, isId: boolean): string {
  if (!iso) {
    return isId ? "Belum pernah dipakai" : "Never used";
  }
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  const formatted = date.toLocaleString([], {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
  return isId ? `Dipakai ${formatted}` : `Last used ${formatted}`;
}

export function KeysConsole() {
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";

  const keysQuery = useApiKeys();
  const providersQuery = useProviderConnections();
  const tunnelQuery = useTunnelStatus();
  const startTunnelMutation = useStartTunnelMutation();
  const stopTunnelMutation = useStopTunnelMutation();

  const createMutation = useCreateApiKeyMutation();
  const updateMutation = useUpdateApiKeyMutation();
  const rotateMutation = useRotateApiKeyMutation();
  const deleteMutation = useDeleteApiKeyMutation();

  const [resolvedOrigins] = useState(() =>
    resolveGatewayOrigins(typeof window !== "undefined" ? window.location?.origin : ""),
  );
  const [activeRoute, setActiveRoute] = useState<GatewayRouteId>(resolvedOrigins.defaultRoute);
  const [keySearch, setKeySearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<KeyStatusFilter>("all");

  // Live gateway health probe (RTT latency chip, main parity)
  const [healthState, setHealthState] = useState<"checking" | "online" | "offline">("checking");
  const [healthLatencyMs, setHealthLatencyMs] = useState<number | null>(null);
  const mountedRef = useRef(true);

  // Unified Create / Edit modal
  const [formModal, setFormModal] = useState<{
    open: boolean;
    mode: "create" | "edit";
    initialKey: ApiKeySummary | null;
  }>({ open: false, mode: "create", initialKey: null });
  const [formError, setFormError] = useState<string | null>(null);

  // Secondary modals + ephemeral secrets created/rotated during this session
  const [readySecret, setReadySecret] = useState<SecretRevealState | null>(null);
  const [sessionSecrets, setSessionSecrets] = useState<Record<string, string>>({});
  const [visibleKeyIds, setVisibleKeyIds] = useState<Set<string>>(() => new Set());
  const [snippetKey, setSnippetKey] = useState<ApiKeySummary | null>(null);

  // Row feedback & copy state
  const [pendingConfirm, setPendingConfirm] = useState<PendingRowConfirm>(null);
  const [rowError, setRowError] = useState<string | null>(null);
  const [busyKeyId, setBusyKeyId] = useState<string | null>(null);
  const [copiedKeyId, setCopiedKeyId] = useState<string | null>(null);
  const copyTimerRef = useRef<number | null>(null);

  const probeGatewayHealth = useCallback(async () => {
    setHealthState("checking");
    const t0 = typeof performance !== "undefined" ? performance.now() : Date.now();
    try {
      const res = await fetch("/api/health", { cache: "no-store" });
      const elapsed = Math.max(
        1,
        Math.round((typeof performance !== "undefined" ? performance.now() : Date.now()) - t0),
      );
      if (!mountedRef.current) return;
      if (res.ok) {
        setHealthState("online");
        setHealthLatencyMs(elapsed);
      } else {
        setHealthState("offline");
        setHealthLatencyMs(null);
      }
    } catch {
      if (!mountedRef.current) return;
      setHealthState("offline");
      setHealthLatencyMs(null);
    }
  }, []);

  useEffect(() => {
    mountedRef.current = true;
    void probeGatewayHealth();
    return () => {
      mountedRef.current = false;
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
    };
  }, [probeGatewayHealth]);

  const keys = keysQuery.data ?? [];
  const providers = providersQuery.data ?? [];
  const tunnel = tunnelQuery.data;

  const defaultModel = useMemo(() => {
    for (const conn of providers) {
      if (conn.enabled && conn.models.length > 0 && conn.models[0]) {
        return conn.models[0];
      }
    }
    return "gpt-4o";
  }, [providers]);

  const activeEffectiveOrigin = useMemo(() => {
    if (activeRoute === "domain" && resolvedOrigins.domainOrigin) {
      return resolvedOrigins.domainOrigin;
    }
    if (
      activeRoute === "cloudflare" &&
      tunnel?.provider === "cloudflare" &&
      tunnel.status === "running" &&
      normalizeOrigin(tunnel.public_url)
    ) {
      return normalizeOrigin(tunnel.public_url);
    }
    if (
      activeRoute === "tailscale" &&
      tunnel?.provider === "tailscale" &&
      tunnel.status === "running" &&
      normalizeOrigin(tunnel.public_url)
    ) {
      return normalizeOrigin(tunnel.public_url);
    }
    if (activeRoute === "local") {
      return normalizeOrigin(resolvedOrigins.localOrigin);
    }
    return resolvedOrigins.domainOrigin || normalizeOrigin(resolvedOrigins.localOrigin);
  }, [activeRoute, tunnel, resolvedOrigins]);

  const activeRouteLabel =
    activeRoute === "domain" && resolvedOrigins.domainOrigin
      ? isId
        ? "Host Server"
        : "Server Host"
      : activeRoute === "cloudflare" &&
          tunnel?.provider === "cloudflare" &&
          tunnel.status === "running"
        ? "Cloudflare Tunnel"
        : activeRoute === "tailscale" &&
            tunnel?.provider === "tailscale" &&
            tunnel.status === "running"
          ? "Tailscale Mesh"
          : activeRoute === "local" || !resolvedOrigins.domainOrigin
            ? isId
              ? "Host Lokal"
              : "Local Host"
            : isId
              ? "Host Server"
              : "Server Host";

  const statusBuckets = useMemo(() => countKeyStatusBuckets(keys), [keys]);
  const filteredKeys = useMemo(
    () => filterKeysByStatusAndSearch(keys, statusFilter, keySearch),
    [keys, statusFilter, keySearch],
  );

  const sampleApiKey = useMemo(() => {
    for (const k of keys) {
      if (!k.disabled && sessionSecrets[k.id]) {
        return sessionSecrets[k.id];
      }
    }
    return undefined;
  }, [keys, sessionSecrets]);

  const openCreateModal = () => {
    setFormError(null);
    createMutation.reset();
    setFormModal({ open: true, mode: "create", initialKey: null });
  };

  const openEditModal = (key: ApiKeySummary) => {
    setFormError(null);
    updateMutation.reset();
    setFormModal({ open: true, mode: "edit", initialKey: key });
  };

  const handleFormSubmit = (values: CreateKeyFormOutput, editingKeyId?: string) => {
    setFormError(null);
    setRowError(null);
    if (formModal.mode === "edit" && editingKeyId) {
      updateMutation.mutate(
        {
          id: editingKeyId,
          body: {
            name: values.name,
            allowed_models: values.allowed_models,
            rate_limit_rpm: values.rate_limit_rpm,
            quota_usd: values.quota_usd,
            expires_at: values.expires_at,
          },
        },
        {
          onSuccess: () => setFormModal({ open: false, mode: "create", initialKey: null }),
          onError: (err) => setFormError(err.message),
        },
      );
      return;
    }

    createMutation.mutate(toCreateApiKeyRequest(values), {
      onSuccess: (response) => {
        setFormModal({ open: false, mode: "create", initialKey: null });
        setSessionSecrets((prev) => ({ ...prev, [response.key.id]: response.secret }));
        setReadySecret({
          action: "created",
          keyName: response.key.name,
          secret: response.secret,
        });
      },
      onError: (err) => setFormError(err.message),
    });
  };

  const handleToggleDisabled = (key: ApiKeySummary) => {
    if (busyKeyId !== null) return;
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
    if (busyKeyId !== null) return;
    setRowError(null);
    setBusyKeyId(key.id);
    rotateMutation.mutate(key.id, {
      onSuccess: (response) => {
        setPendingConfirm(null);
        setSessionSecrets((prev) => ({ ...prev, [response.key.id]: response.secret }));
        setReadySecret({
          action: "rotated",
          keyName: response.key.name,
          secret: response.secret,
        });
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyKeyId(null),
    });
  };

  const handleConfirmDelete = (key: ApiKeySummary) => {
    if (busyKeyId !== null) return;
    setRowError(null);
    setBusyKeyId(key.id);
    deleteMutation.mutate(key.id, {
      onSuccess: () => {
        setPendingConfirm(null);
        setSessionSecrets((prev) => {
          const next = { ...prev };
          delete next[key.id];
          return next;
        });
      },
      onError: (err) => setRowError(err.message),
      onSettled: () => setBusyKeyId(null),
    });
  };

  const handleCopyKeyValue = async (keyId: string, value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      setCopiedKeyId(keyId);
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
      copyTimerRef.current = window.setTimeout(() => {
        copyTimerRef.current = null;
        setCopiedKeyId(null);
      }, 2000);
    } catch {
      // ignore
    }
  };

  const toggleRevealSecret = (keyId: string) => {
    setVisibleKeyIds((prev) => {
      const next = new Set(prev);
      if (next.has(keyId)) next.delete(keyId);
      else next.add(keyId);
      return next;
    });
  };

  const handleStartTunnel = (provider: "cloudflare" | "tailscale") => {
    setRowError(null);
    startTunnelMutation.mutate(
      { provider },
      {
        onError: (err) => setRowError(err.message),
      },
    );
  };

  const handleStopTunnel = () => {
    setRowError(null);
    stopTunnelMutation.mutate(undefined, {
      onError: (err) => setRowError(err.message),
    });
  };

  const filterTabs: ReadonlyArray<{ id: KeyStatusFilter; label: string }> = [
    { id: "all", label: isId ? "Semua" : "All" },
    { id: "active", label: isId ? "Aktif" : "Active" },
    { id: "disabled", label: isId ? "Nonaktif" : "Disabled" },
    { id: "expired", label: isId ? "Kedaluwarsa" : "Expired" },
  ];

  return (
    <div className="flex min-w-0 flex-col gap-5 pb-12 lg:pb-8">
      {/* Studio Page Header (main Endpoint & Keys parity) */}
      <header className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-[var(--color-primary)]/20 bg-[var(--color-primary)]/10 text-[var(--color-primary)]">
            <Waypoints className="h-5 w-5" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h1 className="truncate text-lg font-semibold tracking-tight text-[var(--color-text-main)] sm:text-xl">
              {isId ? "Endpoint & Kunci API" : "Endpoint & Keys"}
            </h1>
            <p className="truncate text-xs text-[var(--color-text-muted)]">
              {isId
                ? "Hubungkan klien, salin endpoint gateway, dan kelola kunci API."
                : "Connect clients, copy gateway endpoints, and govern API keys."}
            </p>
          </div>
        </div>

        <div className="flex shrink-0 flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => void probeGatewayHealth()}
            title={
              isId ? "Klik untuk cek ulang kesehatan gateway" : "Click to re-check gateway health"
            }
            className={cn(
              "inline-flex h-8 cursor-pointer items-center gap-1.5 rounded-lg border px-2.5 font-mono text-[11px] transition-colors",
              healthState === "online"
                ? "border-[var(--color-success-border)] bg-[var(--color-success-bg)] text-[var(--color-success-fg)] hover:brightness-95"
                : healthState === "offline"
                  ? "border-[var(--color-danger-border)] bg-[var(--color-danger-bg)] text-[var(--color-danger-fg)]"
                  : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] text-[var(--color-text-muted)]",
            )}
          >
            <span
              className={cn(
                "h-1.5 w-1.5 shrink-0 rounded-full",
                healthState === "online"
                  ? "bg-[var(--color-success)]"
                  : healthState === "offline"
                    ? "bg-[var(--color-danger)]"
                    : "bg-[var(--color-text-subtle)]",
              )}
            />
            <span>
              {healthState === "online"
                ? `Online · ${healthLatencyMs}ms`
                : healthState === "offline"
                  ? isId
                    ? "Tidak Terjangkau"
                    : "Unreachable"
                  : isId
                    ? "Memeriksa…"
                    : "Checking…"}
            </span>
          </button>

          <Button type="button" size="sm" onClick={openCreateModal} className="h-8">
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{isId ? "Buat Key" : "Create Key"}</span>
          </Button>
        </div>
      </header>

      {/* Zona 1: Unified Gateway Route, Base URLs & Quick Integration Snippet */}
      <GatewayEndpointBar
        isId={isId}
        localOrigin={resolvedOrigins.localOrigin}
        domainOrigin={resolvedOrigins.domainOrigin}
        isDomainHost={resolvedOrigins.isDomainHost}
        activeRoute={activeRoute}
        onSelectRoute={setActiveRoute}
        tunnel={tunnel}
        tunnelPending={startTunnelMutation.isPending || stopTunnelMutation.isPending}
        onStartTunnel={handleStartTunnel}
        onStopTunnel={handleStopTunnel}
        defaultModel={defaultModel}
        sampleApiKey={sampleApiKey}
      />

      {/* Error alert for key/tunnel actions */}
      {rowError ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="flex-1">
            <AlertTitle>{isId ? "Operasi gagal" : "Operation failed"}</AlertTitle>
            <AlertDescription className="font-mono">{rowError}</AlertDescription>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setRowError(null)}
            aria-label={isId ? "Tutup error" : "Dismiss error"}
          >
            <X className="h-3.5 w-3.5" aria-hidden="true" />
          </Button>
        </Alert>
      ) : null}

      {/* Zona 2: API Keys Vault */}
      <section className="space-y-3.5">
        <div className="flex flex-col gap-3">
          <div className="flex flex-col justify-between gap-2.5 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <KeyRound className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
              <h2 className="text-sm font-semibold tracking-tight text-[var(--color-text-main)] sm:text-base">
                {isId ? "Kunci API" : "API Keys"}
              </h2>
              <span className="rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] tabular-nums text-[var(--color-text-muted)]">
                {filteredKeys.length === keys.length
                  ? keys.length
                  : `${filteredKeys.length} / ${keys.length}`}
              </span>
            </div>

            {keys.length > 0 ? (
              <div className="relative w-full sm:w-64">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-subtle)]"
                  aria-hidden="true"
                />
                <input
                  type="text"
                  value={keySearch}
                  onChange={(e) => setKeySearch(e.target.value)}
                  placeholder={
                    isId ? "Cari nama atau prefix kunci…" : "Search by key name or prefix…"
                  }
                  aria-label={isId ? "Cari kunci API" : "Search API keys"}
                  className="h-9 w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] pr-8 pl-8 text-xs text-[var(--color-text-main)] placeholder:text-[var(--color-text-subtle)] focus:border-[var(--color-primary)] focus:outline-none"
                />
                {keySearch ? (
                  <button
                    type="button"
                    onClick={() => setKeySearch("")}
                    aria-label={isId ? "Hapus pencarian" : "Clear search"}
                    className="absolute top-1/2 right-2 flex h-5 w-5 -translate-y-1/2 cursor-pointer items-center justify-center rounded text-[var(--color-text-subtle)] hover:text-[var(--color-text-main)]"
                  >
                    <X className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                ) : null}
              </div>
            ) : null}
          </div>

          {keys.length > 0 ? (
            <div
              className="flex flex-wrap items-center gap-1.5"
              role="tablist"
              aria-label={isId ? "Filter status kunci API" : "Filter API keys by status"}
            >
              {filterTabs.map((tab) => {
                const count = statusBuckets[tab.id];
                const active = statusFilter === tab.id;
                return (
                  <button
                    key={tab.id}
                    type="button"
                    role="tab"
                    aria-selected={active}
                    onClick={() => setStatusFilter(tab.id)}
                    className={cn(
                      "inline-flex h-7 cursor-pointer items-center gap-1.5 rounded-lg border px-2.5 text-xs font-medium transition-colors",
                      active
                        ? "border-[var(--color-primary)]/30 bg-[var(--color-primary)]/12 text-[var(--color-primary)]"
                        : "border-[var(--color-border-subtle)] bg-[var(--color-surface)] text-[var(--color-text-muted)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]",
                    )}
                  >
                    <span>{tab.label}</span>
                    <span
                      className={cn(
                        "rounded px-1 font-mono text-[10.5px] tabular-nums",
                        active
                          ? "bg-[var(--color-primary)]/15 text-[var(--color-primary)]"
                          : "text-[var(--color-text-subtle)]",
                      )}
                    >
                      {count}
                    </span>
                  </button>
                );
              })}
            </div>
          ) : null}
        </div>

        {/* Keys list / states */}
        {keysQuery.isError ? (
          <Alert variant="destructive">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <div className="flex-1">
              <AlertTitle>{isId ? "Gagal memuat kunci API" : "Unable to load API keys"}</AlertTitle>
              <AlertDescription className="font-mono">{keysQuery.error.message}</AlertDescription>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="xs"
              onClick={() => void keysQuery.refetch()}
              disabled={keysQuery.isFetching}
            >
              {isId ? "Coba Lagi" : "Retry"}
            </Button>
          </Alert>
        ) : keysQuery.isPending ? (
          <div className="space-y-2.5">
            {SKELETON_CARDS.map((id) => (
              <div
                key={id}
                className="h-28 animate-pulse rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface)] p-4"
              />
            ))}
          </div>
        ) : filteredKeys.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-[var(--color-border)] bg-[var(--color-surface)] p-8 text-center">
            <KeyRound
              className="mx-auto mb-2.5 h-7 w-7 text-[var(--color-text-subtle)] opacity-65"
              aria-hidden="true"
            />
            <p className="text-sm font-medium text-[var(--color-text-main)]">
              {keys.length === 0
                ? isId
                  ? "Belum ada kunci API"
                  : "No API keys yet"
                : isId
                  ? "Tidak ada kunci yang cocok dengan filter"
                  : "No keys match your current filter"}
            </p>
            <p className="mt-1 mb-4 text-xs text-[var(--color-text-muted)]">
              {keys.length === 0
                ? isId
                  ? "Buat kunci API pertama untuk mengautentikasi klien IDE dan SDK."
                  : "Create your first API key to authenticate IDE and SDK clients."
                : isId
                  ? "Coba ubah kata kunci pencarian atau kembali ke tab Semua."
                  : "Try adjusting your search query or switching back to All keys."}
            </p>
            {keys.length === 0 ? (
              <Button type="button" size="sm" onClick={openCreateModal}>
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Buat Key" : "Create Key"}</span>
              </Button>
            ) : (
              <Button
                type="button"
                size="sm"
                variant="secondary"
                onClick={() => {
                  setKeySearch("");
                  setStatusFilter("all");
                }}
              >
                {isId ? "Reset Filter" : "Reset Filters"}
              </Button>
            )}
          </div>
        ) : (
          <div className="space-y-2.5">
            {filteredKeys.map((key) => {
              const expired = isKeyExpiredAt(key.expires_at);
              const expiry = getKeyExpiryState(key.expires_at, Date.now(), isId);
              const isDisabled = key.disabled;
              const isRowBusy = busyKeyId === key.id;
              const anyRowBusy = busyKeyId !== null;
              const confirmingRotate =
                pendingConfirm?.id === key.id && pendingConfirm.kind === "rotate";
              const confirmingDelete =
                pendingConfirm?.id === key.id && pendingConfirm.kind === "delete";
              const knownSecret = sessionSecrets[key.id];
              const revealed = Boolean(knownSecret && visibleKeyIds.has(key.id));
              const modelsSummary =
                key.allowed_models.length > 0
                  ? key.allowed_models.join(", ")
                  : isId
                    ? "Semua"
                    : "All";

              const accentBarClass = expired
                ? "bg-[var(--color-danger)]"
                : isDisabled
                  ? "bg-[var(--color-warning)]"
                  : "bg-[var(--color-primary)]";
              const statusDotClass = expired
                ? "bg-[var(--color-danger)]"
                : isDisabled
                  ? "bg-[var(--color-warning)]"
                  : "bg-[var(--color-success)]";

              return (
                <article
                  key={key.id}
                  className={cn(
                    "group relative rounded-xl border p-4 transition-colors duration-150",
                    isDisabled || expired
                      ? "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/60 opacity-85"
                      : "border-[var(--color-border)] bg-[var(--color-surface)] hover:border-[var(--color-primary)]/40",
                  )}
                >
                  <span
                    className={cn(
                      "pointer-events-none absolute inset-y-3 left-0 w-[3px] rounded-r-full",
                      accentBarClass,
                    )}
                    aria-hidden="true"
                  />

                  {/* Row 1: Identity & Action Bar (main ApiKeyCard parity) */}
                  <div className="flex flex-wrap items-center justify-between gap-2.5 pl-1.5">
                    <div className="flex min-w-0 flex-wrap items-center gap-2">
                      <span className={cn("h-2 w-2 shrink-0 rounded-full", statusDotClass)} />
                      <h3
                        className="max-w-[260px] truncate text-sm font-semibold text-[var(--color-text-main)]"
                        title={key.name}
                      >
                        {key.name}
                      </h3>
                      {isDisabled ? (
                        <span className="shrink-0 rounded-md border border-[var(--color-warning-border)] bg-[var(--color-warning-bg)] px-2 py-0.5 text-[10px] font-medium text-[var(--color-warning-fg)]">
                          {isId ? "Nonaktif" : "Disabled"}
                        </span>
                      ) : null}
                      {expired && !isDisabled ? (
                        <span className="shrink-0 rounded-md border border-[var(--color-danger-border)] bg-[var(--color-danger-bg)] px-2 py-0.5 text-[10px] font-semibold text-[var(--color-danger-fg)]">
                          {isId ? "Kedaluwarsa" : "Expired"}
                        </span>
                      ) : null}
                    </div>

                    <div className="flex shrink-0 items-center gap-1">
                      {confirmingRotate ? (
                        <>
                          <Button
                            type="button"
                            size="xs"
                            onClick={() => handleConfirmRotate(key)}
                            disabled={anyRowBusy}
                          >
                            {isRowBusy
                              ? isId
                                ? "Merotasi..."
                                : "Rotating..."
                              : isId
                                ? "Konfirmasi Rotasi"
                                : "Confirm Rotate"}
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            size="xs"
                            onClick={() => setPendingConfirm(null)}
                            disabled={anyRowBusy}
                          >
                            {isId ? "Batal" : "Cancel"}
                          </Button>
                        </>
                      ) : confirmingDelete ? (
                        <>
                          <Button
                            type="button"
                            variant="destructive"
                            size="xs"
                            onClick={() => handleConfirmDelete(key)}
                            disabled={anyRowBusy}
                          >
                            {isRowBusy
                              ? isId
                                ? "Menghapus..."
                                : "Deleting..."
                              : isId
                                ? "Hapus"
                                : "Delete"}
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            size="xs"
                            onClick={() => setPendingConfirm(null)}
                            disabled={anyRowBusy}
                          >
                            {isId ? "Batal" : "Cancel"}
                          </Button>
                        </>
                      ) : (
                        <>
                          {/* Active Toggle Switch */}
                          <button
                            type="button"
                            role="switch"
                            aria-checked={!isDisabled}
                            aria-label={
                              isDisabled
                                ? isId
                                  ? `Aktifkan ${key.name}`
                                  : `Enable ${key.name}`
                                : isId
                                  ? `Nonaktifkan ${key.name}`
                                  : `Disable ${key.name}`
                            }
                            title={
                              isDisabled
                                ? isId
                                  ? "Aktifkan kunci"
                                  : "Enable key"
                                : isId
                                  ? "Nonaktifkan kunci"
                                  : "Disable key"
                            }
                            disabled={anyRowBusy}
                            onClick={() => handleToggleDisabled(key)}
                            className={cn(
                              "relative inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border transition-colors disabled:cursor-not-allowed disabled:opacity-50",
                              !isDisabled
                                ? "border-[var(--color-primary)] bg-[var(--color-primary)]"
                                : "border-[var(--color-border)] bg-[var(--color-surface-3)]",
                            )}
                          >
                            <span
                              className={cn(
                                "inline-block h-3.5 w-3.5 rounded-full bg-[var(--color-surface)] shadow-xs transition-transform",
                                !isDisabled ? "translate-x-4" : "translate-x-0.5",
                              )}
                            />
                          </button>

                          <div
                            className="mx-1 h-4 w-px bg-[var(--color-border-subtle)]"
                            aria-hidden="true"
                          />

                          <button
                            type="button"
                            onClick={() => setSnippetKey(key)}
                            className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)]"
                            title={isId ? "Cuplikan kode" : "Code snippet"}
                            aria-label={
                              isId
                                ? `Buka cuplikan kode ${key.name}`
                                : `Open code snippet for ${key.name}`
                            }
                          >
                            <Code2 className="h-4 w-4" aria-hidden="true" />
                          </button>

                          <button
                            type="button"
                            onClick={() => openEditModal(key)}
                            disabled={anyRowBusy}
                            className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)] disabled:opacity-40"
                            title={isId ? "Edit kunci" : "Edit key"}
                            aria-label={isId ? `Edit ${key.name}` : `Edit ${key.name}`}
                          >
                            <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
                          </button>

                          <button
                            type="button"
                            onClick={() => setPendingConfirm({ id: key.id, kind: "rotate" })}
                            disabled={anyRowBusy}
                            className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-primary)] disabled:opacity-40"
                            title={isId ? "Rotasi secret" : "Rotate secret"}
                            aria-label={isId ? `Rotasi ${key.name}` : `Rotate ${key.name}`}
                          >
                            <RotateCw className="h-3.5 w-3.5" aria-hidden="true" />
                          </button>

                          <button
                            type="button"
                            onClick={() => setPendingConfirm({ id: key.id, kind: "delete" })}
                            disabled={anyRowBusy}
                            className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-danger-bg)] hover:text-[var(--color-danger-fg)] disabled:opacity-40"
                            title={isId ? "Hapus kunci" : "Delete key"}
                            aria-label={isId ? `Hapus ${key.name}` : `Delete ${key.name}`}
                          >
                            <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                          </button>
                        </>
                      )}
                    </div>
                  </div>

                  {/* Row 2: Masked Key Prefix Bar */}
                  <div className="mt-2.5 flex items-center justify-between gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/65 px-3 py-1.5 pl-3">
                    <code className="min-w-0 truncate font-mono text-xs text-[var(--color-text-main)]">
                      {revealed ? knownSecret : `${key.key_prefix}••••••••••••`}
                    </code>
                    <div className="flex shrink-0 items-center gap-1">
                      {knownSecret ? (
                        <>
                          <button
                            type="button"
                            onClick={() => toggleRevealSecret(key.id)}
                            className="flex h-7 w-7 cursor-pointer items-center justify-center rounded-md text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-3)] hover:text-[var(--color-text-main)]"
                            title={
                              revealed
                                ? isId
                                  ? "Sembunyikan"
                                  : "Hide"
                                : isId
                                  ? "Tampilkan"
                                  : "Reveal"
                            }
                            aria-label={
                              revealed
                                ? isId
                                  ? `Sembunyikan secret ${key.name}`
                                  : `Hide secret for ${key.name}`
                                : isId
                                  ? `Tampilkan secret ${key.name}`
                                  : `Reveal secret for ${key.name}`
                            }
                          >
                            {revealed ? (
                              <EyeOff className="h-3.5 w-3.5" aria-hidden="true" />
                            ) : (
                              <Eye className="h-3.5 w-3.5" aria-hidden="true" />
                            )}
                          </button>
                          <button
                            type="button"
                            onClick={() => void handleCopyKeyValue(key.id, knownSecret)}
                            className={cn(
                              "flex h-7 w-7 cursor-pointer items-center justify-center rounded-md transition-colors",
                              copiedKeyId === key.id
                                ? "bg-[var(--color-success-bg)] text-[var(--color-success)]"
                                : "text-[var(--color-text-muted)] hover:bg-[var(--color-surface-3)] hover:text-[var(--color-primary)]",
                            )}
                            title={isId ? "Salin secret" : "Copy API key"}
                            aria-label={isId ? `Salin ${key.name}` : `Copy ${key.name}`}
                          >
                            {copiedKeyId === key.id ? (
                              <Check className="h-3.5 w-3.5" aria-hidden="true" />
                            ) : (
                              <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                            )}
                          </button>
                        </>
                      ) : (
                        <span className="font-mono text-[10.5px] text-[var(--color-text-subtle)]">
                          {formatCompactLastUsed(key.last_used_at, isId)}
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Row 3: Compact Policy & Lifecycle Chips */}
                  <div className="mt-2.5 flex flex-wrap items-center gap-1.5 pl-1.5 font-mono text-[11px]">
                    <span
                      className="max-w-[260px] truncate rounded-md border border-[var(--color-info-border)] bg-[var(--color-info-bg)] px-2 py-0.5 text-[var(--color-info-fg)]"
                      title={`${isId ? "Model diizinkan" : "Allowed models"}: ${modelsSummary}`}
                    >
                      {isId ? "Model" : "Models"}: {modelsSummary}
                    </span>

                    {key.rate_limit_rpm > 0 ? (
                      <span className="rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 tabular-nums text-[var(--color-text-muted)]">
                        {key.rate_limit_rpm} RPM
                      </span>
                    ) : null}

                    {key.quota_usd > 0 ? (
                      <span className="rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 tabular-nums text-[var(--color-text-muted)]">
                        {isId ? "Kuota" : "Quota"}: ${key.quota_usd.toFixed(2)}
                      </span>
                    ) : null}

                    {knownSecret ? (
                      <span className="rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 text-[var(--color-text-subtle)]">
                        {formatCompactLastUsed(key.last_used_at, isId)}
                      </span>
                    ) : null}

                    {expiry.status !== "none" ? (
                      <span
                        className={cn(
                          "rounded-md border px-2 py-0.5",
                          expiry.status === "expired"
                            ? "border-[var(--color-danger-border)] bg-[var(--color-danger-bg)] font-semibold text-[var(--color-danger-fg)]"
                            : expiry.status === "soon"
                              ? "border-[var(--color-warning-border)] bg-[var(--color-warning-bg)] font-medium text-[var(--color-warning-fg)]"
                              : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] text-[var(--color-text-muted)]",
                        )}
                        title={key.expires_at ?? undefined}
                      >
                        {expiry.label}
                      </span>
                    ) : null}
                  </div>
                </article>
              );
            })}
          </div>
        )}
      </section>

      {/* Unified Create / Edit Key Modal */}
      <KeyFormModal
        isOpen={formModal.open}
        isId={isId}
        mode={formModal.mode}
        initialKey={formModal.initialKey}
        isPending={createMutation.isPending || updateMutation.isPending}
        errorMessage={formError}
        onClose={() => setFormModal({ open: false, mode: "create", initialKey: null })}
        onSubmit={handleFormSubmit}
      />

      {/* One-Time Key Ready Modal (Copy Key + Copy Ready cURL) */}
      <KeyReadyModal
        isId={isId}
        secretState={readySecret}
        activeOrigin={activeEffectiveOrigin}
        defaultModel={defaultModel}
        onClose={() => setReadySecret(null)}
      />

      {/* Code Snippet Modal (OpenAI & Anthropic x cURL/Python/Node/Go) */}
      <KeySnippetModal
        isId={isId}
        apiKey={snippetKey}
        revealedSecret={snippetKey ? sessionSecrets[snippetKey.id] : undefined}
        activeOrigin={activeEffectiveOrigin}
        activeRouteLabel={activeRouteLabel}
        defaultModel={defaultModel}
        onClose={() => setSnippetKey(null)}
      />
    </div>
  );
}
