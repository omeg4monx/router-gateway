"use client";

import { useState, useEffect, useMemo } from "react";
import PropTypes from "prop-types";
import { Button, Toggle, SegmentedControl } from "@/shared/components";
import { useCopyToClipboard } from "@/shared/hooks/useCopyToClipboard";
import { cn } from "@/shared/utils/cn";
import StatusAlert from "./StatusAlert";
import { normalizeOrigin } from "../endpointUtils";
import { generateEndpointSnippet } from "../endpointSnippets";

export default function GatewayConsole({
  isApiKeyUser = false,
  localOrigin = "",
  activeRoute = "local",
  onSelectRoute,
  tunnel,
  tailscale,
  requireApiKey = false,
  tunnelDashboardAccess = false,
  onToggleRequireApiKey,
  onToggleTunnelDashboardAccess,
  onEnableTunnelClick,
  onEnableAuthAndTunnel,
  onDisableTunnelClick,
  onConnectTailscaleClick,
  onDisableTailscaleClick,
  onDismissTunnelStatus,
  onDismissTsStatus,
  keys = [],
  defaultModel = "gpt-4o",
  securitySectionRef,
}) {
  const { copied, copy } = useCopyToClipboard();
  const [showSnippet, setShowSnippet] = useState(false);
  const [snippetFormat, setSnippetFormat] = useState("openai");
  const [snippetLang, setSnippetLang] = useState("curl");
  const [selectedKeyId, setSelectedKeyId] = useState("");
  const [snippetModel, setSnippetModel] = useState(defaultModel);
  const [modelTouched, setModelTouched] = useState(false);

  useEffect(() => {
    if (!modelTouched && defaultModel) {
      setSnippetModel(defaultModel);
    }
  }, [defaultModel, modelTouched]);

  const activeKeys = useMemo(() => keys.filter((k) => k?.isActive !== false), [keys]);
  const chosenKeyObj = useMemo(() => {
    if (selectedKeyId) {
      const found = keys.find((k) => k.id === selectedKeyId);
      if (found) return found;
    }
    return activeKeys[0] || keys[0] || null;
  }, [keys, activeKeys, selectedKeyId]);

  const normalizedLocal = normalizeOrigin(localOrigin) || "http://localhost:20128";
  const cfLive = Boolean(tunnel?.enabled && tunnel?.reachable && (tunnel?.publicUrl || tunnel?.url));
  const tsLive = Boolean(tailscale?.enabled && tailscale?.reachable && tailscale?.url);

  const cfOrigin = normalizeOrigin(tunnel?.publicUrl || tunnel?.url || "");
  const tsOrigin = normalizeOrigin(tailscale?.url || "");

  let effectiveOrigin = normalizedLocal;
  let usingLocalFallback = false;
  let activeRouteLabel = "Local Host";

  if (!isApiKeyUser && activeRoute === "cloudflare") {
    activeRouteLabel = "Cloudflare Tunnel";
    if (cfLive && cfOrigin) {
      effectiveOrigin = cfOrigin;
    } else {
      effectiveOrigin = normalizedLocal;
      usingLocalFallback = true;
    }
  } else if (!isApiKeyUser && activeRoute === "tailscale") {
    activeRouteLabel = "Tailscale Mesh";
    if (tsLive && tsOrigin) {
      effectiveOrigin = tsOrigin;
    } else {
      effectiveOrigin = normalizedLocal;
      usingLocalFallback = true;
    }
  }

  const openaiBaseUrl = `${effectiveOrigin}/v1`;
  const claudeBaseUrl = effectiveOrigin;

  const snippetCode = generateEndpointSnippet({
    format: snippetFormat,
    lang: snippetLang,
    baseUrl: effectiveOrigin,
    apiKey: chosenKeyObj?.key || "",
    model: snippetModel,
  });

  const routeOptions = [
    {
      id: "local",
      label: "Local Host",
      icon: "computer",
      badge: "ACTIVE",
      tone: "success",
      sub: normalizedLocal,
    },
    ...(!isApiKeyUser
      ? [
          {
            id: "cloudflare",
            label: "Cloudflare Tunnel",
            icon: "cloud_sync",
            badge: cfLive ? "PUBLIC" : tunnel?.enabled || tunnel?.loading ? "CONNECTING" : "OFFLINE",
            tone: cfLive ? "success" : tunnel?.enabled || tunnel?.loading ? "warning" : "muted",
            sub: cfLive
              ? cfOrigin
              : tunnel?.loading
                ? tunnel.progress || "Starting tunnel…"
                : "Public HTTPS without open ports",
          },
          {
            id: "tailscale",
            label: "Tailscale Mesh",
            icon: "vpn_lock",
            badge: tsLive ? "MESH ACTIVE" : tailscale?.enabled || tailscale?.loading ? "CONNECTING" : "OFFLINE",
            tone: tsLive ? "success" : tailscale?.enabled || tailscale?.loading ? "warning" : "muted",
            sub: tsLive
              ? tsOrigin
              : tailscale?.loading
                ? tailscale.progress || "Connecting Tailscale…"
                : tailscale?.port
                  ? `Tailnet Funnel → :${tailscale.port}`
                  : "Private tailnet & Funnel access",
          },
        ]
      : []),
  ];

  const badgeClasses = (tone) => {
    if (tone === "success") return "bg-success-bg border-success-border text-success-fg";
    if (tone === "warning") return "bg-warning-bg border-warning-border text-warning-fg";
    return "bg-surface-2 border-border-subtle text-text-muted";
  };

  return (
    <section className="relative rounded-2xl border border-border bg-surface shadow-[var(--shadow-soft)] overflow-hidden">
      <span className="ember-rule w-full" aria-hidden="true" />

      <div className="p-4 sm:p-5 space-y-4">
        {/* Route / Interface Switcher */}
        <div className="space-y-2">
          <div className="flex items-center justify-between gap-2">
            <span className="micro-label">Gateway Route &amp; Exposure</span>
            <span className="text-[11px] font-mono text-text-subtle">
              Active target: <strong className="text-text-main font-medium">{activeRouteLabel}</strong>
            </span>
          </div>

          <div className={cn("grid gap-2", isApiKeyUser ? "grid-cols-1" : "grid-cols-1 sm:grid-cols-3")}>
            {routeOptions.map((opt) => {
              const selected = (isApiKeyUser ? "local" : activeRoute) === opt.id;
              return (
                <button
                  key={opt.id}
                  type="button"
                  onClick={() => onSelectRoute?.(opt.id)}
                  aria-pressed={selected}
                  className={cn(
                    "group flex items-start gap-2.5 rounded-xl border p-3 text-left transition-all duration-150 cursor-pointer",
                    selected
                      ? "border-primary/50 bg-primary/8 shadow-[var(--shadow-soft)]"
                      : "border-border-subtle bg-surface-2/50 hover:border-border hover:bg-surface-2"
                  )}
                >
                  <div
                    className={cn(
                      "size-8 rounded-lg border flex items-center justify-center shrink-0 transition-colors",
                      selected
                        ? "border-primary/30 bg-primary/12 text-primary"
                        : "border-border-subtle bg-surface text-text-muted group-hover:text-primary"
                    )}
                  >
                    <span className="material-symbols-outlined text-[17px]">{opt.icon}</span>
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center justify-between gap-1.5">
                      <span className="text-xs font-semibold text-text-main truncate">{opt.label}</span>
                      <span
                        className={cn(
                          "px-1.5 py-0.5 rounded-md border font-mono text-[9.5px] font-medium shrink-0",
                          badgeClasses(opt.tone)
                        )}
                      >
                        {opt.badge}
                      </span>
                    </div>
                    <p className="mt-0.5 text-[11px] font-mono text-text-muted truncate" title={opt.sub}>
                      {opt.sub}
                    </p>
                  </div>
                </button>
              );
            })}
          </div>
        </div>

        {/* Inline Exposure Control Strip for Remote Routes (Admin Only) */}
        {!isApiKeyUser && activeRoute === "cloudflare" && (
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 rounded-xl border border-border-subtle bg-surface-2/60 px-3.5 py-2.5">
            <div className="min-w-0">
              <p className="text-xs font-medium text-text-main">
                {cfLive
                  ? "Cloudflare Tunnel is live and routing traffic to your gateway."
                  : tunnel?.loading
                    ? tunnel.progress || "Starting Cloudflare Tunnel…"
                    : tunnel?.enabled
                      ? "Tunnel started — verifying reachability…"
                      : "Expose this gateway over HTTPS without port forwarding."}
              </p>
              {tunnel?.authGateWarning && !requireApiKey && (
                <p className="mt-1 text-xs text-warning-fg">
                  API Key Authentication must be active before opening a public Cloudflare tunnel.
                </p>
              )}
            </div>
            <div className="flex flex-wrap items-center gap-2 shrink-0">
              {tunnel?.authGateWarning && !requireApiKey ? (
                <Button size="sm" icon="lock" onClick={onEnableAuthAndTunnel}>
                  {"Enable Auth & Continue"}
                </Button>
              ) : cfLive || tunnel?.enabled ? (
                <Button
                  size="sm"
                  variant="danger"
                  icon="power_settings_new"
                  loading={tunnel?.loading}
                  onClick={onDisableTunnelClick}
                >
                  Disable Tunnel
                </Button>
              ) : (
                <Button
                  size="sm"
                  variant="secondary"
                  icon="cloud_upload"
                  loading={tunnel?.loading}
                  onClick={onEnableTunnelClick}
                >
                  Enable Tunnel
                </Button>
              )}
            </div>
          </div>
        )}

        {!isApiKeyUser && (activeRoute === "tailscale" || Boolean(tailscale?.authUrl)) && (
          <div className="flex flex-col gap-2.5 rounded-xl border border-border-subtle bg-surface-2/60 px-3.5 py-2.5">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2.5">
              <div className="min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <p className="text-xs font-medium text-text-main">
                    {tsLive
                      ? "Tailscale Funnel is active and reachable on your tailnet."
                      : tailscale?.loading || tailscale?.connecting
                        ? tailscale.progress || "Connecting Tailscale…"
                        : tailscale?.enabled
                          ? `Funnel configured → 127.0.0.1:${tailscale.port || 20128}`
                          : "Connect your router to a private Tailscale tailnet with Funnel."}
                  </p>
                  {(tailscale?.port || tailscale?.enabled) && (
                    <span className="px-1.5 py-0.5 rounded-md border border-border bg-surface font-mono text-[10px] text-text-muted">
                      PORT {tailscale.port || 20128}
                    </span>
                  )}
                </div>
              </div>
              <div className="flex items-center gap-2 shrink-0">
                {tsLive || tailscale?.enabled ? (
                  <Button
                    size="sm"
                    variant="danger"
                    icon="power_settings_new"
                    loading={tailscale?.loading}
                    onClick={onDisableTailscaleClick}
                  >
                    Disable Tailscale
                  </Button>
                ) : (
                  <Button
                    size="sm"
                    variant="secondary"
                    icon="vpn_lock"
                    loading={tailscale?.loading || tailscale?.connecting}
                    onClick={onConnectTailscaleClick}
                  >
                    Connect Tailscale
                  </Button>
                )}
              </div>
            </div>

            {tailscale?.authUrl && (
              <div className="flex flex-wrap items-center gap-2 pt-1 border-t border-border-subtle">
                <a
                  href={tailscale.authUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-primary text-on-primary text-xs font-semibold hover:bg-primary-hover transition-colors"
                >
                  <span className="material-symbols-outlined text-[15px]">open_in_new</span>
                  <span>{tailscale.authLabel || "Open Login Page"}</span>
                </a>
                <Button
                  size="sm"
                  variant="secondary"
                  icon={copied === "ts_auth_url" ? "check" : "content_copy"}
                  onClick={() => copy(tailscale.authUrl, "ts_auth_url")}
                >
                  {copied === "ts_auth_url" ? "Copied" : "Copy Auth URL"}
                </Button>
              </div>
            )}
          </div>
        )}

        {/* Inline Dismissible Status Alerts (never swallowed silently) */}
        {!isApiKeyUser && tunnel?.status && (
          <div className="relative">
            <StatusAlert status={tunnel.status} className="pr-9 rounded-xl border border-border-subtle" />
            {onDismissTunnelStatus && (
              <button
                type="button"
                onClick={onDismissTunnelStatus}
                aria-label="Dismiss tunnel status"
                className="absolute right-2 top-1/2 -translate-y-1/2 size-6 rounded-md flex items-center justify-center text-current opacity-70 hover:opacity-100 cursor-pointer"
              >
                <span className="material-symbols-outlined text-[15px]">close</span>
              </button>
            )}
          </div>
        )}

        {!isApiKeyUser && tailscale?.status && (
          <div className="relative">
            <StatusAlert status={tailscale.status} className="pr-9 rounded-xl border border-border-subtle" />
            {onDismissTsStatus && (
              <button
                type="button"
                onClick={onDismissTsStatus}
                aria-label="Dismiss Tailscale status"
                className="absolute right-2 top-1/2 -translate-y-1/2 size-6 rounded-md flex items-center justify-center text-current opacity-70 hover:opacity-100 cursor-pointer"
              >
                <span className="material-symbols-outlined text-[15px]">close</span>
              </button>
            )}
          </div>
        )}

        {/* Synchronized OpenAI & Anthropic Base URL Rows */}
        <div className="space-y-2">
          {usingLocalFallback && (
            <div className="flex items-center justify-between gap-2 text-xs text-text-muted px-1">
              <span className="flex items-center gap-1.5">
                <span className="material-symbols-outlined text-[15px] text-warning">info</span>
                <span>Showing Local URL until {activeRouteLabel} is active.</span>
              </span>
              <button
                type="button"
                onClick={() => onSelectRoute?.("local")}
                className="text-primary hover:underline font-medium shrink-0 cursor-pointer"
              >
                Switch to Local
              </button>
            </div>
          )}

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-2">
            <div className="flex items-center gap-2 rounded-xl border border-border bg-surface-2/70 p-2 transition-colors focus-within:border-primary/60">
              <span className="px-2 py-1 rounded-md bg-surface border border-border-subtle font-mono text-[10.5px] font-semibold text-primary shrink-0">
                OpenAI
              </span>
              <input
                type="text"
                readOnly
                aria-label="OpenAI compatible base URL"
                value={openaiBaseUrl}
                className="flex-1 min-w-0 bg-transparent px-1 font-mono text-xs text-text-main outline-none select-all"
              />
              <Button
                size="sm"
                variant="secondary"
                icon={copied === "openai_url" ? "check" : "content_copy"}
                onClick={() => copy(openaiBaseUrl, "openai_url")}
              >
                {copied === "openai_url" ? "Copied" : "Copy"}
              </Button>
            </div>

            <div className="flex items-center gap-2 rounded-xl border border-border bg-surface-2/70 p-2 transition-colors focus-within:border-primary/60">
              <span className="px-2 py-1 rounded-md bg-surface border border-border-subtle font-mono text-[10.5px] font-semibold text-secondary-fg shrink-0">
                Claude
              </span>
              <input
                type="text"
                readOnly
                aria-label="Anthropic Claude compatible base URL"
                value={claudeBaseUrl}
                className="flex-1 min-w-0 bg-transparent px-1 font-mono text-xs text-text-main outline-none select-all"
              />
              <Button
                size="sm"
                variant="secondary"
                icon={copied === "claude_url" ? "check" : "content_copy"}
                onClick={() => copy(claudeBaseUrl, "claude_url")}
              >
                {copied === "claude_url" ? "Copied" : "Copy"}
              </Button>
            </div>
          </div>
        </div>

        {/* Collapsible Quick Integration Snippet */}
        <div className="pt-3 border-t border-border-subtle">
          <div className="flex items-center justify-between gap-2">
            <button
              type="button"
              onClick={() => setShowSnippet((prev) => !prev)}
              aria-expanded={showSnippet}
              className="inline-flex items-center gap-1.5 text-xs font-medium text-text-muted hover:text-primary transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined text-[16px] text-primary">code</span>
              <span>{showSnippet ? "Hide Quick Integration Snippet" : "Quick Integration Snippet"}</span>
              <span
                className="material-symbols-outlined text-[16px] transition-transform duration-200"
                style={{ transform: showSnippet ? "rotate(180deg)" : "rotate(0deg)" }}
              >
                expand_more
              </span>
            </button>

            {showSnippet && (
              <Button
                size="sm"
                variant="ghost"
                icon={copied === "quick_snippet" ? "check" : "content_copy"}
                onClick={() => copy(snippetCode, "quick_snippet")}
              >
                {copied === "quick_snippet" ? "Copied!" : "Copy Code"}
              </Button>
            )}
          </div>

          <div className="accordion-grid" data-open={showSnippet ? "true" : "false"}>
            <div className="accordion-inner">
              <div className="pt-3 space-y-2.5">
                <div className="flex flex-wrap items-center gap-2">
                  <SegmentedControl
                    size="sm"
                    value={snippetFormat}
                    onChange={setSnippetFormat}
                    options={[
                      { value: "openai", label: "OpenAI" },
                      { value: "anthropic", label: "Anthropic" },
                    ]}
                  />
                  <SegmentedControl
                    size="sm"
                    value={snippetLang}
                    onChange={setSnippetLang}
                    options={[
                      { value: "curl", label: "cURL" },
                      { value: "python", label: "Python" },
                      { value: "node", label: "Node.js" },
                      { value: "go", label: "Go" },
                    ]}
                  />
                  {keys.length > 0 && (
                    <select
                      aria-label="API key for snippet"
                      value={chosenKeyObj?.id || ""}
                      onChange={(e) => setSelectedKeyId(e.target.value)}
                      className="field-control h-8 px-2.5 text-xs max-w-[180px] truncate"
                    >
                      {keys.map((k) => (
                        <option key={k.id} value={k.id}>
                          Key: {k.name}
                        </option>
                      ))}
                    </select>
                  )}
                  <input
                    type="text"
                    aria-label="Model name for snippet"
                    value={snippetModel}
                    onChange={(e) => {
                      setModelTouched(true);
                      setSnippetModel(e.target.value);
                    }}
                    placeholder="Model ID (e.g. gpt-4o)"
                    className="field-control h-8 px-2.5 text-xs font-mono w-40 sm:w-48"
                  />
                </div>

                <pre className="rounded-xl border border-border-subtle bg-bg-alt p-3.5 font-mono text-[11.5px] text-text-main overflow-x-auto custom-scrollbar whitespace-pre-wrap break-all leading-relaxed">
                  {snippetCode}
                </pre>
              </div>
            </div>
          </div>
        </div>

        {/* Admin-Only Security & Access Policy Strip */}
        {!isApiKeyUser && (
          <div
            id="endpoint-security-policy"
            ref={securitySectionRef}
            className="pt-3.5 border-t border-border-subtle flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3.5"
          >
            <Toggle
              checked={requireApiKey}
              onChange={(next) => onToggleRequireApiKey?.(next)}
              label="Require API Key Authentication"
              description={
                requireApiKey
                  ? "Active: unauthenticated requests are rejected with HTTP 401."
                  : "Warning: anyone with the URL can call your gateway without a key."
              }
            />

            {(tunnel?.enabled || tailscale?.enabled) && (
              <div className="pt-2.5 sm:pt-0 sm:border-l sm:border-border-subtle sm:pl-4">
                <Toggle
                  checked={tunnelDashboardAccess}
                  onChange={(next) => onToggleTunnelDashboardAccess?.(next)}
                  label="Allow Dashboard over Tunnel"
                  description="Expose the web GUI through your active tunnel URL."
                />
              </div>
            )}
          </div>
        )}
      </div>
    </section>
  );
}

GatewayConsole.propTypes = {
  isApiKeyUser: PropTypes.bool,
  localOrigin: PropTypes.string,
  activeRoute: PropTypes.oneOf(["local", "cloudflare", "tailscale"]),
  onSelectRoute: PropTypes.func,
  tunnel: PropTypes.object,
  tailscale: PropTypes.object,
  requireApiKey: PropTypes.bool,
  tunnelDashboardAccess: PropTypes.bool,
  onToggleRequireApiKey: PropTypes.func,
  onToggleTunnelDashboardAccess: PropTypes.func,
  onEnableTunnelClick: PropTypes.func,
  onEnableAuthAndTunnel: PropTypes.func,
  onDisableTunnelClick: PropTypes.func,
  onConnectTailscaleClick: PropTypes.func,
  onDisableTailscaleClick: PropTypes.func,
  onDismissTunnelStatus: PropTypes.func,
  onDismissTsStatus: PropTypes.func,
  keys: PropTypes.array,
  defaultModel: PropTypes.string,
  securitySectionRef: PropTypes.object,
};
