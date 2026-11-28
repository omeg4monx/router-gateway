import {
  Check,
  ChevronDown,
  Cloud,
  Code2,
  Copy,
  ExternalLink,
  Globe,
  Info,
  Laptop,
  Network,
  Power,
} from "lucide-react";
import { memo, useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router";
import {
  useStartTunnelMutation,
  useStopTunnelMutation,
  useTunnelStatus,
} from "@/features/tunnel/hooks";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/button";
import {
  buildGatewayUrls,
  copyTextToClipboard,
  type GatewayRouteId,
  generateKeySnippet,
  type KeySnippetFormat,
  type KeySnippetLang,
  normalizeOrigin,
} from "./endpointUtils";

export interface ResolvedActiveRoute {
  origin: string;
  label: string;
}

interface GatewayEndpointBarProps {
  isId: boolean;
  localOrigin: string;
  domainOrigin?: string | null;
  isDomainHost?: boolean;
  activeRoute: GatewayRouteId;
  onSelectRoute: (route: GatewayRouteId) => void;
  onResolvedRouteChange?: (target: ResolvedActiveRoute) => void;
  onTunnelError?: (message: string) => void;
  defaultModel: string;
  sampleApiKey?: string;
}

const SNIPPET_LANGS: ReadonlyArray<{ id: KeySnippetLang; label: string }> = [
  { id: "curl", label: "cURL" },
  { id: "python", label: "Python" },
  { id: "node", label: "Node.js" },
  { id: "go", label: "Go" },
];

export const GatewayEndpointBar = memo(function GatewayEndpointBar({
  isId,
  localOrigin,
  domainOrigin,
  isDomainHost = false,
  activeRoute,
  onSelectRoute,
  onResolvedRouteChange,
  onTunnelError,
  defaultModel,
  sampleApiKey,
}: GatewayEndpointBarProps) {
  // Encapsulate 5s tunnel status polling inside this leaf section so parent KeysConsole
  // and API key cards do not re-render on every poll tick.
  const tunnelQuery = useTunnelStatus();
  const startTunnelMutation = useStartTunnelMutation();
  const stopTunnelMutation = useStopTunnelMutation();
  const tunnel = tunnelQuery.data;
  const tunnelPending = startTunnelMutation.isPending || stopTunnelMutation.isPending;

  const [copiedId, setCopiedId] = useState<string | null>(null);
  const [showSnippet, setShowSnippet] = useState(false);
  const [snippetFormat, setSnippetFormat] = useState<KeySnippetFormat>("openai");
  const [snippetLang, setSnippetLang] = useState<KeySnippetLang>("curl");
  const copyTimerRef = useRef<number | null>(null);

  useEffect(() => {
    return () => {
      if (copyTimerRef.current !== null) {
        window.clearTimeout(copyTimerRef.current);
      }
    };
  }, []);

  const copyValue = async (text: string, id: string) => {
    const ok = await copyTextToClipboard(text);
    if (!ok) {
      return;
    }
    setCopiedId(id);
    if (copyTimerRef.current !== null) {
      window.clearTimeout(copyTimerRef.current);
    }
    copyTimerRef.current = window.setTimeout(() => {
      copyTimerRef.current = null;
      setCopiedId(null);
    }, 2000);
  };

  const handleStartTunnel = (provider: "cloudflare" | "tailscale") => {
    startTunnelMutation.mutate(
      { provider },
      {
        onError: (err) => onTunnelError?.(err.message),
      },
    );
  };

  const handleStopTunnel = () => {
    stopTunnelMutation.mutate(undefined, {
      onError: (err) => onTunnelError?.(err.message),
    });
  };

  const cfLive =
    tunnel?.provider === "cloudflare" &&
    tunnel.status === "running" &&
    Boolean(normalizeOrigin(tunnel.public_url));
  const cfStarting = tunnel?.provider === "cloudflare" && tunnel.status === "starting";
  const tsLive =
    tunnel?.provider === "tailscale" &&
    tunnel.status === "running" &&
    Boolean(normalizeOrigin(tunnel.public_url));
  const tsStarting = tunnel?.provider === "tailscale" && tunnel.status === "starting";

  const cfOrigin = cfLive ? normalizeOrigin(tunnel?.public_url) : "";
  const tsOrigin = tsLive ? normalizeOrigin(tunnel?.public_url) : "";
  const cleanDomain = normalizeOrigin(domainOrigin);
  const cleanLocal = normalizeOrigin(localOrigin);
  const fallbackRoute: "domain" | "local" = cleanDomain ? "domain" : "local";

  const effectiveOrigin = useMemo(() => {
    if (activeRoute === "domain" && cleanDomain) return cleanDomain;
    if (activeRoute === "cloudflare" && cfLive) return cfOrigin;
    if (activeRoute === "tailscale" && tsLive) return tsOrigin;
    if (activeRoute === "local") return cleanLocal;
    return cleanDomain || cleanLocal;
  }, [activeRoute, cleanDomain, cfLive, cfOrigin, tsLive, tsOrigin, cleanLocal]);

  const usingPrimaryFallback =
    (activeRoute === "cloudflare" && !cfLive) || (activeRoute === "tailscale" && !tsLive);

  const activeRouteLabel =
    activeRoute === "domain" && cleanDomain
      ? isId
        ? "Host Server"
        : "Server Host"
      : activeRoute === "cloudflare"
        ? "Cloudflare Tunnel"
        : activeRoute === "tailscale"
          ? "Tailscale Mesh"
          : isId
            ? "Host Lokal"
            : "Local Host";

  const effectiveRouteLabel = usingPrimaryFallback
    ? cleanDomain
      ? isId
        ? "Host Server"
        : "Server Host"
      : isId
        ? "Host Lokal"
        : "Local Host"
    : activeRouteLabel;

  useEffect(() => {
    onResolvedRouteChange?.({
      origin: effectiveOrigin,
      label: effectiveRouteLabel,
    });
  }, [effectiveOrigin, effectiveRouteLabel, onResolvedRouteChange]);

  const { openaiBaseUrl, claudeBaseUrl } = useMemo(
    () => buildGatewayUrls(effectiveOrigin),
    [effectiveOrigin],
  );

  const snippetCode = useMemo(
    () =>
      generateKeySnippet({
        format: snippetFormat,
        lang: snippetLang,
        baseUrl: effectiveOrigin,
        apiKey: sampleApiKey,
        model: defaultModel,
      }),
    [snippetFormat, snippetLang, effectiveOrigin, sampleApiKey, defaultModel],
  );

  const routeOptions: ReadonlyArray<{
    id: GatewayRouteId;
    label: string;
    icon: typeof Laptop;
    badge: string;
    tone: "success" | "warning" | "muted";
    sub: string;
  }> = [
    ...(cleanDomain
      ? [
          {
            id: "domain" as const,
            label: isId ? "Host Server" : "Server Host",
            icon: Globe,
            badge: isDomainHost ? "DOMAIN" : "HOST",
            tone: "success" as const,
            sub: cleanDomain,
          },
        ]
      : []),
    {
      id: "local",
      label: isId ? "Host Lokal" : "Local Host",
      icon: Laptop,
      badge: "ONLINE",
      tone: "success",
      sub: cleanLocal || "http://127.0.0.1:20128",
    },
    {
      id: "cloudflare",
      label: "Cloudflare Tunnel",
      icon: Cloud,
      badge: cfLive ? "PUBLIC" : cfStarting ? (isId ? "MEMULAI" : "STARTING") : "OFFLINE",
      tone: cfLive ? "success" : cfStarting ? "warning" : "muted",
      sub: cfLive
        ? cfOrigin
        : isId
          ? "HTTPS publik tanpa buka port"
          : "Public HTTPS without open ports",
    },
    {
      id: "tailscale",
      label: "Tailscale Mesh",
      icon: Network,
      badge: tsLive ? "ACTIVE" : tsStarting ? (isId ? "MEMULAI" : "STARTING") : "OFFLINE",
      tone: tsLive ? "success" : tsStarting ? "warning" : "muted",
      sub: tsLive
        ? tsOrigin
        : isId
          ? "Akses privat tailnet & Funnel"
          : "Private tailnet & Funnel access",
    },
  ];

  const badgeToneClass = (tone: "success" | "warning" | "muted") => {
    if (tone === "success") {
      return "border-[var(--color-success-border)] bg-[var(--color-success-bg)] text-[var(--color-success-fg)]";
    }
    if (tone === "warning") {
      return "border-[var(--color-warning-border)] bg-[var(--color-warning-bg)] text-[var(--color-warning-fg)]";
    }
    return "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] text-[var(--color-text-muted)]";
  };

  return (
    <section className="overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] shadow-[0_1px_2px_rgba(11,13,18,0.04)]">
      <div className="space-y-4 p-4 sm:p-5">
        {/* Route Switcher */}
        <div className="space-y-2">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <span className="font-mono text-[10.5px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">
              {isId ? "Rute & Endpoint Gateway" : "Gateway Route & Exposure"}
            </span>
            <span className="font-mono text-[11px] text-[var(--color-text-subtle)]">
              {isId ? "Target aktif" : "Active target"}:{" "}
              <strong className="font-medium text-[var(--color-text-main)]">
                {activeRouteLabel}
              </strong>
            </span>
          </div>

          <div
            className={cn(
              "grid grid-cols-1 gap-2",
              cleanDomain ? "sm:grid-cols-2 lg:grid-cols-4" : "sm:grid-cols-3",
            )}
          >
            {routeOptions.map((opt) => {
              const Icon = opt.icon;
              const selected = activeRoute === opt.id;
              return (
                <button
                  key={opt.id}
                  type="button"
                  onClick={() => onSelectRoute(opt.id)}
                  aria-pressed={selected}
                  className={cn(
                    "group flex cursor-pointer items-start gap-2.5 rounded-xl border p-3 text-left transition-all duration-150",
                    selected
                      ? "border-[var(--color-primary)]/50 bg-[var(--color-primary)]/8"
                      : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 hover:border-[var(--color-border)] hover:bg-[var(--color-surface-2)]",
                  )}
                >
                  <div
                    className={cn(
                      "flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border transition-colors",
                      selected
                        ? "border-[var(--color-primary)]/30 bg-[var(--color-primary)]/12 text-[var(--color-primary)]"
                        : "border-[var(--color-border-subtle)] bg-[var(--color-surface)] text-[var(--color-text-muted)] group-hover:text-[var(--color-primary)]",
                    )}
                  >
                    <Icon className="h-4 w-4" aria-hidden="true" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between gap-1.5">
                      <span className="truncate text-xs font-semibold text-[var(--color-text-main)]">
                        {opt.label}
                      </span>
                      <span
                        className={cn(
                          "shrink-0 rounded-md border px-1.5 py-0.5 font-mono text-[9.5px] font-medium",
                          badgeToneClass(opt.tone),
                        )}
                      >
                        {opt.badge}
                      </span>
                    </div>
                    <p
                      className="mt-0.5 truncate font-mono text-[11px] text-[var(--color-text-muted)]"
                      title={opt.sub}
                    >
                      {opt.sub}
                    </p>
                  </div>
                </button>
              );
            })}
          </div>
        </div>

        {/* Compact inline tunnel control strip when Cloudflare/Tailscale selected */}
        {activeRoute === "cloudflare" || activeRoute === "tailscale" ? (
          <div className="flex flex-col justify-between gap-2.5 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/60 px-3.5 py-2.5 sm:flex-row sm:items-center">
            <p className="min-w-0 text-xs text-[var(--color-text-main)]">
              {activeRoute === "cloudflare"
                ? cfLive
                  ? isId
                    ? "Cloudflare Tunnel aktif dan meneruskan trafik ke gateway."
                    : "Cloudflare Tunnel is live and routing traffic to your gateway."
                  : isId
                    ? "Aktifkan Cloudflare Tunnel untuk ekspos gateway via HTTPS."
                    : "Start Cloudflare Tunnel to expose this gateway over HTTPS."
                : tsLive
                  ? isId
                    ? "Tailscale aktif dan dapat diakses melalui tailnet Anda."
                    : "Tailscale is active and reachable on your tailnet."
                  : isId
                    ? "Aktifkan Tailscale untuk akses privat di tailnet Anda."
                    : "Start Tailscale for private access across your tailnet."}
            </p>
            <div className="flex shrink-0 flex-wrap items-center gap-2">
              {(activeRoute === "cloudflare" && cfLive) ||
              (activeRoute === "tailscale" && tsLive) ? (
                <Button
                  type="button"
                  size="xs"
                  variant="destructive"
                  disabled={tunnelPending}
                  onClick={handleStopTunnel}
                >
                  <Power className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Matikan" : "Stop"}</span>
                </Button>
              ) : (
                <Button
                  type="button"
                  size="xs"
                  variant="secondary"
                  disabled={tunnelPending}
                  onClick={() => handleStartTunnel(activeRoute)}
                >
                  <Power className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {tunnelPending
                      ? isId
                        ? "Memulai..."
                        : "Starting..."
                      : isId
                        ? "Aktifkan"
                        : "Start"}
                  </span>
                </Button>
              )}
              <Link
                to="/dashboard/tunnel/"
                className="inline-flex h-7 items-center gap-1 rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-2.5 text-xs font-medium text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text-main)]"
              >
                <span>{isId ? "Pengaturan Tunnel" : "Tunnel Settings"}</span>
                <ExternalLink className="h-3 w-3" aria-hidden="true" />
              </Link>
            </div>
          </div>
        ) : null}

        {/* Synchronized OpenAI & Claude Base URLs */}
        <div className="space-y-2">
          {usingPrimaryFallback ? (
            <div className="flex flex-wrap items-center justify-between gap-2 px-1 text-xs text-[var(--color-text-muted)]">
              <span className="inline-flex items-center gap-1.5">
                <Info className="h-3.5 w-3.5 text-[var(--color-warning)]" aria-hidden="true" />
                <span>
                  {isId
                    ? `Menampilkan URL ${cleanDomain ? "Host Server" : "Lokal"} hingga ${activeRouteLabel} aktif.`
                    : `Showing ${cleanDomain ? "Server Host" : "Local"} URL until ${activeRouteLabel} is active.`}
                </span>
              </span>
              <button
                type="button"
                onClick={() => onSelectRoute(fallbackRoute)}
                className="shrink-0 cursor-pointer font-medium text-[var(--color-primary)] hover:underline"
              >
                {cleanDomain
                  ? isId
                    ? "Pindah ke Host Server"
                    : "Switch to Server Host"
                  : isId
                    ? "Pindah ke Lokal"
                    : "Switch to Local"}
              </button>
            </div>
          ) : null}

          <div className="grid grid-cols-1 gap-2 lg:grid-cols-2">
            <div className="flex items-center gap-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/70 p-2 transition-colors focus-within:border-[var(--color-primary)]/60">
              <span className="shrink-0 rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-2 py-1 font-mono text-[10.5px] font-semibold text-[var(--color-primary)]">
                OpenAI
              </span>
              <input
                type="text"
                readOnly
                aria-label={isId ? "Base URL kompatibel OpenAI" : "OpenAI compatible base URL"}
                value={openaiBaseUrl}
                className="min-w-0 flex-1 bg-transparent px-1 font-mono text-xs text-[var(--color-text-main)] outline-none select-all"
              />
              <Button
                type="button"
                size="xs"
                variant="secondary"
                onClick={() => void copyValue(openaiBaseUrl, "openai_url")}
              >
                {copiedId === "openai_url" ? (
                  <>
                    <Check className="h-3.5 w-3.5 text-[var(--color-success)]" aria-hidden="true" />
                    <span>{isId ? "Tersalin" : "Copied"}</span>
                  </>
                ) : (
                  <>
                    <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Salin" : "Copy"}</span>
                  </>
                )}
              </Button>
            </div>

            <div className="flex items-center gap-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/70 p-2 transition-colors focus-within:border-[var(--color-primary)]/60">
              <span className="shrink-0 rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-2 py-1 font-mono text-[10.5px] font-semibold text-[var(--color-secondary-fg)]">
                Claude
              </span>
              <input
                type="text"
                readOnly
                aria-label={
                  isId
                    ? "Base URL kompatibel Anthropic Claude"
                    : "Anthropic Claude compatible base URL"
                }
                value={claudeBaseUrl}
                className="min-w-0 flex-1 bg-transparent px-1 font-mono text-xs text-[var(--color-text-main)] outline-none select-all"
              />
              <Button
                type="button"
                size="xs"
                variant="secondary"
                onClick={() => void copyValue(claudeBaseUrl, "claude_url")}
              >
                {copiedId === "claude_url" ? (
                  <>
                    <Check className="h-3.5 w-3.5 text-[var(--color-success)]" aria-hidden="true" />
                    <span>{isId ? "Tersalin" : "Copied"}</span>
                  </>
                ) : (
                  <>
                    <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Salin" : "Copy"}</span>
                  </>
                )}
              </Button>
            </div>
          </div>
        </div>

        {/* Collapsible Quick Integration Snippet */}
        <div className="border-t border-[var(--color-border-subtle)] pt-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <button
              type="button"
              onClick={() => setShowSnippet((prev) => !prev)}
              aria-expanded={showSnippet}
              className="inline-flex cursor-pointer items-center gap-1.5 text-xs font-medium text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-primary)]"
            >
              <Code2 className="h-3.5 w-3.5 text-[var(--color-primary)]" aria-hidden="true" />
              <span>
                {showSnippet
                  ? isId
                    ? "Sembunyikan Cuplikan Integrasi"
                    : "Hide Quick Integration Snippet"
                  : isId
                    ? "Cuplikan Integrasi Cepat"
                    : "Quick Integration Snippet"}
              </span>
              <ChevronDown
                className={cn(
                  "h-3.5 w-3.5 transition-transform duration-200",
                  showSnippet && "rotate-180",
                )}
                aria-hidden="true"
              />
            </button>

            {showSnippet ? (
              <Button
                type="button"
                size="xs"
                variant="ghost"
                onClick={() => void copyValue(snippetCode, "quick_snippet")}
              >
                {copiedId === "quick_snippet" ? (
                  <>
                    <Check className="h-3.5 w-3.5 text-[var(--color-success)]" aria-hidden="true" />
                    <span>{isId ? "Tersalin" : "Copied"}</span>
                  </>
                ) : (
                  <>
                    <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Salin Kode" : "Copy Code"}</span>
                  </>
                )}
              </Button>
            ) : null}
          </div>

          {showSnippet ? (
            <div className="mt-3 space-y-2.5">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="inline-flex rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] p-0.5">
                  {(
                    [
                      { id: "openai", label: "OpenAI (/v1)" },
                      { id: "anthropic", label: "Anthropic (/v1/messages)" },
                    ] as const
                  ).map((item) => (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => setSnippetFormat(item.id)}
                      className={cn(
                        "cursor-pointer rounded-md px-2.5 py-1 text-[11px] font-medium transition-colors",
                        snippetFormat === item.id
                          ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                          : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                      )}
                    >
                      {item.label}
                    </button>
                  ))}
                </div>

                <div className="inline-flex rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] p-0.5">
                  {SNIPPET_LANGS.map((lang) => (
                    <button
                      key={lang.id}
                      type="button"
                      onClick={() => setSnippetLang(lang.id)}
                      className={cn(
                        "cursor-pointer rounded-md px-2 py-1 font-mono text-[11px] font-medium transition-colors",
                        snippetLang === lang.id
                          ? "bg-[var(--color-surface)] text-[var(--color-text-main)] shadow-xs"
                          : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                      )}
                    >
                      {lang.label}
                    </button>
                  ))}
                </div>
              </div>

              <pre className="custom-scrollbar max-h-64 overflow-x-auto whitespace-pre-wrap break-all rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-bg-alt)] p-3.5 font-mono text-xs leading-relaxed text-[var(--color-text-main)]">
                {snippetCode}
              </pre>
            </div>
          ) : null}
        </div>
      </div>
    </section>
  );
});
