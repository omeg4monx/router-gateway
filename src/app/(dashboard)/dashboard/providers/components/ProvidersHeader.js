"use client";

import { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import { Button } from "@/shared/components";
import { cn } from "@/shared/utils/cn";
import ModelAvailabilityBadge from "./ModelAvailabilityBadge";

const CUSTOM_PROTOCOL_OPTIONS = [
  {
    variant: "openai",
    label: "OpenAI Compatible",
    hint: "Chat Completions or Responses API",
    icon: "smart_toy",
    iconClass: "text-info",
  },
  {
    variant: "anthropic",
    label: "Anthropic Compatible",
    hint: "Messages API (/v1/messages)",
    icon: "psychology",
    iconClass: "text-warning",
  },
  {
    variant: "moonshot",
    label: "MoonshotAI Compatible",
    hint: "Kimi OpenAI-compatible gateway",
    icon: "dark_mode",
    iconClass: "text-primary",
  },
];

export default function ProvidersHeader({
  telemetry,
  statusFilter,
  onSelectStatusFilter,
  activeAllCount,
  testingState,
  onTestAll,
  onOpenAddModal,
  onCooldownCleared,
}) {
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef(null);

  useEffect(() => {
    if (!menuOpen) return;
    const handleMouseDown = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setMenuOpen(false);
      }
    };
    const handleKeyDown = (e) => {
      if (e.key === "Escape") setMenuOpen(false);
    };
    document.addEventListener("mousedown", handleMouseDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("mousedown", handleMouseDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [menuOpen]);

  const isTestingAll = testingState?.mode === "all";
  const toggleFilter = (target) => {
    onSelectStatusFilter(statusFilter === target ? "all" : target);
  };

  const hasAttention = (telemetry?.attentionProviders || 0) > 0;

  return (
    <div className="flex flex-col gap-3.5">
      {/* Studio Header Band */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3 min-w-0">
          <div className="size-9 shrink-0 rounded-xl bg-primary/10 border border-primary/20 text-primary flex items-center justify-center shadow-xs">
            <span className="material-symbols-outlined text-[20px]">dns</span>
          </div>
          <div className="min-w-0">
            <h1 className="text-lg sm:text-xl font-semibold tracking-tight text-text-main truncate">
              Providers
            </h1>
            <p className="text-xs text-text-muted truncate">
              Monitor upstream health, test connections, and configure AI providers
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2 shrink-0">
          <ModelAvailabilityBadge onCooldownCleared={onCooldownCleared} />

          <Button
            size="sm"
            variant="secondary"
            icon={isTestingAll ? "progress_activity" : "play_arrow"}
            onClick={onTestAll}
            disabled={activeAllCount === 0 || Boolean(testingState)}
            title={
              activeAllCount === 0
                ? "No active connections to test"
                : `Test all ${activeAllCount} active connection(s)`
            }
            className="whitespace-nowrap"
          >
            {isTestingAll ? "Testing..." : "Test All"}
          </Button>

          <div className="relative" ref={menuRef}>
            <Button
              size="sm"
              icon="add"
              onClick={() => setMenuOpen((v) => !v)}
              aria-expanded={menuOpen}
              aria-haspopup="menu"
              className="whitespace-nowrap"
            >
              <span>Add Custom</span>
              <span className="material-symbols-outlined text-[16px] -mr-1">
                {menuOpen ? "expand_less" : "expand_more"}
              </span>
            </Button>

            {menuOpen && (
              <div
                role="menu"
                aria-label="Add custom compatible provider"
                className="absolute left-0 sm:left-auto sm:right-0 top-full mt-1.5 w-60 max-w-[calc(100vw-2rem)] rounded-xl border border-border shell-popover popover-in p-1.5 z-50"
              >
                {CUSTOM_PROTOCOL_OPTIONS.map((opt) => (
                  <button
                    key={opt.variant}
                    type="button"
                    role="menuitem"
                    onClick={() => {
                      setMenuOpen(false);
                      onOpenAddModal(opt.variant);
                    }}
                    className="flex w-full items-start gap-2.5 rounded-lg px-2.5 py-2 text-left hover:bg-surface-2 transition-colors cursor-pointer"
                  >
                    <span
                      className={cn(
                        "material-symbols-outlined text-[18px] mt-0.5 shrink-0",
                        opt.iconClass,
                      )}
                    >
                      {opt.icon}
                    </span>
                    <span className="min-w-0">
                      <span className="block text-xs font-medium text-text-main">
                        {opt.label}
                      </span>
                      <span className="block text-[11px] text-text-muted truncate">
                        {opt.hint}
                      </span>
                    </span>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Clickable Fleet Telemetry Strip */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
        <button
          type="button"
          onClick={() => toggleFilter("active")}
          aria-pressed={statusFilter === "active"}
          className={cn(
            "group flex flex-col justify-between rounded-xl border bg-surface px-3.5 py-2.5 text-left transition-all cursor-pointer shadow-[var(--shadow-soft)]",
            statusFilter === "active"
              ? "border-primary ring-1 ring-primary/30 bg-primary/5"
              : "border-border hover:border-primary/40 hover:bg-surface-2/60",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-text-muted">
              Active Providers
            </span>
            <span className="material-symbols-outlined text-[16px] text-success">
              check_circle
            </span>
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="text-lg font-semibold font-mono tabular-nums text-text-main">
              {telemetry.activeProviders}
            </span>
            <span className="text-[11px] text-text-muted truncate">
              {telemetry.configuredActiveProviders} configured · {telemetry.noAuthProviders} built-in
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => onSelectStatusFilter("all")}
          className={cn(
            "group flex flex-col justify-between rounded-xl border bg-surface px-3.5 py-2.5 text-left transition-all cursor-pointer shadow-[var(--shadow-soft)]",
            "border-border hover:border-primary/40 hover:bg-surface-2/60",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-text-muted">
              Active Connections
            </span>
            <span className="material-symbols-outlined text-[16px] text-primary">
              hub
            </span>
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="text-lg font-semibold font-mono tabular-nums text-text-main">
              {telemetry.enabledConnections}
              <span className="text-xs font-normal text-text-muted">
                /{telemetry.totalConnections}
              </span>
            </span>
            <span className="text-[11px] text-text-muted truncate">
              Enabled credentials
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => toggleFilter("error")}
          aria-pressed={statusFilter === "error"}
          className={cn(
            "group flex flex-col justify-between rounded-xl border bg-surface px-3.5 py-2.5 text-left transition-all cursor-pointer shadow-[var(--shadow-soft)]",
            statusFilter === "error"
              ? "border-danger ring-1 ring-danger/30 bg-danger/5"
              : hasAttention
                ? "border-danger/35 hover:border-danger/60 hover:bg-danger/5"
                : "border-border hover:border-primary/40 hover:bg-surface-2/60",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-text-muted">
              Needs Attention
            </span>
            <span
              className={cn(
                "material-symbols-outlined text-[16px]",
                hasAttention ? "text-danger" : "text-text-muted/60",
              )}
            >
              {hasAttention ? "error" : "verified_user"}
            </span>
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span
              className={cn(
                "text-lg font-semibold font-mono tabular-nums",
                hasAttention ? "text-danger" : "text-text-main",
              )}
            >
              {telemetry.attentionProviders}
            </span>
            <span className="text-[11px] text-text-muted truncate">
              {hasAttention ? "Providers with errors" : "All connections clear"}
            </span>
          </div>
        </button>

        <button
          type="button"
          onClick={() => toggleFilter("inactive")}
          aria-pressed={statusFilter === "inactive"}
          className={cn(
            "group flex flex-col justify-between rounded-xl border bg-surface px-3.5 py-2.5 text-left transition-all cursor-pointer shadow-[var(--shadow-soft)]",
            statusFilter === "inactive"
              ? "border-primary ring-1 ring-primary/30 bg-primary/5"
              : "border-border hover:border-primary/40 hover:bg-surface-2/60",
          )}
        >
          <div className="flex items-center justify-between gap-2">
            <span className="text-[11px] font-medium text-text-muted">
              Disabled
            </span>
            <span className="material-symbols-outlined text-[16px] text-text-muted">
              pause_circle
            </span>
          </div>
          <div className="mt-1 flex items-baseline gap-2">
            <span className="text-lg font-semibold font-mono tabular-nums text-text-main">
              {telemetry.disabledProviders}
            </span>
            <span className="text-[11px] text-text-muted truncate">
              Paused providers
            </span>
          </div>
        </button>
      </div>
    </div>
  );
}

ProvidersHeader.propTypes = {
  telemetry: PropTypes.shape({
    activeProviders: PropTypes.number,
    configuredActiveProviders: PropTypes.number,
    noAuthProviders: PropTypes.number,
    enabledConnections: PropTypes.number,
    totalConnections: PropTypes.number,
    attentionProviders: PropTypes.number,
    disabledProviders: PropTypes.number,
  }).isRequired,
  statusFilter: PropTypes.string.isRequired,
  onSelectStatusFilter: PropTypes.func.isRequired,
  activeAllCount: PropTypes.number.isRequired,
  testingState: PropTypes.shape({
    mode: PropTypes.string,
    providerId: PropTypes.string,
    label: PropTypes.string,
  }),
  onTestAll: PropTypes.func.isRequired,
  onOpenAddModal: PropTypes.func.isRequired,
  onCooldownCleared: PropTypes.func,
};
