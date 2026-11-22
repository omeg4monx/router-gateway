"use client";

import { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import { Toggle } from "@/shared/components";
import { nextResetAt, formatDuration } from "@/shared/utils/interval";
import { useNow } from "@/shared/hooks/useNow";
import { cn } from "@/shared/utils/cn";
import {
  formatTokensNumber,
  maskApiKey,
  getExpiryMeta,
  countWhitelistEntries,
} from "../endpointUtils";

function ResetCountdownInline({ resetInterval, lastResetAt }) {
  const now = useNow(true);
  if (!resetInterval || resetInterval === "never") return null;
  const due = nextResetAt(resetInterval, lastResetAt);
  if (!due) return null;
  const remaining = new Date(due).getTime() - now;
  return (
    <span className="text-[11px] font-mono text-text-subtle">
      {remaining <= 0 ? "Resetting…" : `Resets in ${formatDuration(remaining)}`}
    </span>
  );
}

ResetCountdownInline.propTypes = {
  resetInterval: PropTypes.string,
  lastResetAt: PropTypes.string,
};

export default function ApiKeyCard({
  apiKey,
  creatorLabel = "",
  canBulkSelect = false,
  selected = false,
  onToggleSelect,
  isOwnKey = false,
  revealed = false,
  onToggleReveal,
  copied = false,
  onCopy,
  onToggleActive,
  onOpenSnippet,
  onEdit,
  onDuplicate,
  onResetUsage,
  onShare,
  onDelete,
}) {
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef(null);

  useEffect(() => {
    if (!menuOpen) return;
    const onMouseDown = (e) => {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setMenuOpen(false);
      }
    };
    const onKeyDown = (e) => {
      if (e.key === "Escape") {
        setMenuOpen(false);
      }
    };
    document.addEventListener("mousedown", onMouseDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onMouseDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [menuOpen]);

  const limit = Number(apiKey.tokenLimit) || 0;
  const used = Number(apiKey.usedTokens) || 0;
  const isUnlimited = limit <= 0;
  const isExhausted = !isUnlimited && used >= limit;
  const pct = isUnlimited ? 0 : Math.min(100, Math.max(0, Math.round((used / limit) * 100)));
  const remaining = isUnlimited ? 0 : Math.max(0, limit - used);
  const isWarning = !isUnlimited && pct >= 75 && !isExhausted;
  const isDisabled = apiKey.isActive === false;

  const expiry = getExpiryMeta(apiKey.expiresAt);
  const isExpired = expiry.status === "expired";
  const whitelistCount = countWhitelistEntries(apiKey.ipWhitelist);

  const accentBarClass =
    isExhausted || isExpired
      ? "bg-danger"
      : isDisabled
        ? "bg-warning"
        : "bg-primary";

  const statusDotClass =
    isExhausted || isExpired
      ? "bg-danger"
      : isDisabled
        ? "bg-warning"
        : "bg-success";

  const modelsSummary =
    apiKey.allowedModels && apiKey.allowedModels !== "*" ? apiKey.allowedModels : "All";

  const runMenuAction = (fn) => {
    setMenuOpen(false);
    fn?.();
  };

  return (
    <article
      className={cn(
        "group relative rounded-xl border p-4 transition-colors duration-150",
        menuOpen && "z-30",
        isDisabled || isExpired
          ? "border-border-subtle bg-surface-2/60 opacity-85"
          : "border-border bg-surface hover:border-primary/40"
      )}
    >
      <span
        className={cn("pointer-events-none absolute inset-y-3 left-0 w-[3px] rounded-r-full", accentBarClass)}
        aria-hidden="true"
      />

      {/* Row 1: Identity & Safe Action Bar */}
      <div className="flex items-center justify-between gap-3 pl-1">
        <div className="flex items-center gap-2.5 min-w-0 flex-wrap">
          {canBulkSelect && (
            <input
              type="checkbox"
              checked={selected}
              disabled={isOwnKey}
              onChange={() => onToggleSelect?.(apiKey.id)}
              className="size-4 accent-[var(--color-primary)] shrink-0 cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
              aria-label={`Select ${apiKey.name}`}
            />
          )}
          <span className={cn("size-2 rounded-full shrink-0", statusDotClass)} />
          <h4 className="text-sm font-semibold text-text-main truncate max-w-[260px]" title={apiKey.name}>
            {apiKey.name}
          </h4>
          {creatorLabel && (
            <span className="px-2 py-0.5 rounded-md bg-surface-2 border border-border-subtle font-mono text-[10px] text-text-muted shrink-0">
              {creatorLabel}
            </span>
          )}
          {isDisabled && (
            <span className="px-2 py-0.5 rounded-md bg-warning-bg border border-warning-border text-warning-fg text-[10px] font-medium shrink-0">
              Disabled
            </span>
          )}
          {isExpired && (
            <span className="px-2 py-0.5 rounded-md bg-danger-bg border border-danger-border text-danger-fg text-[10px] font-semibold shrink-0">
              Expired
            </span>
          )}
          {isExhausted && !isExpired && (
            <span className="px-2 py-0.5 rounded-md bg-danger-bg border border-danger-border text-danger-fg text-[10px] font-semibold shrink-0">
              Exhausted
            </span>
          )}
        </div>

        <div className="flex items-center gap-1.5 shrink-0">
          <span title={isOwnKey ? "Cannot disable your own active session key" : isDisabled ? "Activate key" : "Disable key"}>
            <Toggle
              size="sm"
              checked={!isDisabled}
              disabled={isOwnKey}
              onChange={(nextActive) => onToggleActive?.(apiKey, nextActive)}
            />
          </span>

          <div className="h-4 w-px bg-border-subtle mx-0.5" aria-hidden="true" />

          <button
            type="button"
            onClick={() => onOpenSnippet?.(apiKey)}
            className="size-8 rounded-lg flex items-center justify-center text-text-muted hover:bg-surface-2 hover:text-primary transition-colors cursor-pointer"
            title="Code snippet"
            aria-label={`Open code snippet for ${apiKey.name}`}
          >
            <span className="material-symbols-outlined text-[17px]">code</span>
          </button>

          <button
            type="button"
            onClick={() => onEdit?.(apiKey)}
            disabled={isOwnKey}
            className="size-8 rounded-lg flex items-center justify-center text-text-muted hover:bg-surface-2 hover:text-primary transition-colors cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
            title={isOwnKey ? "Cannot edit your own active session key" : "Edit key"}
            aria-label={`Edit ${apiKey.name}`}
          >
            <span className="material-symbols-outlined text-[17px]">edit</span>
          </button>

          {/* Overflow menu separates destructive Delete from frequent clicks */}
          <div className="relative" ref={menuRef}>
            <button
              type="button"
              onClick={() => setMenuOpen((prev) => !prev)}
              aria-haspopup="menu"
              aria-expanded={menuOpen}
              aria-label={`More actions for ${apiKey.name}`}
              className={cn(
                "size-8 rounded-lg flex items-center justify-center transition-colors cursor-pointer",
                menuOpen
                  ? "bg-primary/12 text-primary"
                  : "text-text-muted hover:bg-surface-2 hover:text-text-main"
              )}
              title="More actions"
            >
              <span className="material-symbols-outlined text-[18px]">more_vert</span>
            </button>

            {menuOpen && (
              <div
                role="menu"
                className="shell-popover popover-in absolute right-0 top-full mt-1.5 z-40 w-48 rounded-xl p-1.5"
              >
                <button
                  type="button"
                  role="menuitem"
                  disabled={isOwnKey}
                  onClick={() => runMenuAction(() => onDuplicate?.(apiKey))}
                  className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs font-medium text-text-main hover:bg-surface-2 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[16px] text-text-muted">library_add</span>
                  <span>Duplicate Key</span>
                </button>
                <button
                  type="button"
                  role="menuitem"
                  disabled={isOwnKey}
                  onClick={() => runMenuAction(() => onResetUsage?.(apiKey))}
                  className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs font-medium text-text-main hover:bg-surface-2 transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[16px] text-text-muted">restart_alt</span>
                  <span>Reset Token Usage</span>
                </button>
                <button
                  type="button"
                  role="menuitem"
                  onClick={() => runMenuAction(() => onShare?.(apiKey))}
                  className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs font-medium text-text-main hover:bg-surface-2 transition-colors cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[16px] text-text-muted">share</span>
                  <span>Share Usage Link</span>
                </button>
                <div className="my-1 border-t border-border-subtle" role="separator" />
                <button
                  type="button"
                  role="menuitem"
                  disabled={isOwnKey}
                  onClick={() => runMenuAction(() => onDelete?.(apiKey.id))}
                  className="flex w-full items-center gap-2.5 rounded-lg px-2.5 py-2 text-xs font-medium text-danger hover:bg-danger-bg transition-colors disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[16px]">delete</span>
                  <span>Delete Key</span>
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Row 2: Credential Bar */}
      <div className="mt-3 pl-1">
        <div className="flex items-center justify-between gap-2 rounded-lg border border-border-subtle bg-surface-2/80 px-3 py-1.5 font-mono text-xs">
          <span className="text-text-main truncate select-all">
            {revealed ? apiKey.key : maskApiKey(apiKey.key)}
          </span>
          <div className="flex items-center gap-1 shrink-0">
            <button
              type="button"
              onClick={() => onToggleReveal?.(apiKey.id)}
              className="size-7 rounded-md flex items-center justify-center text-text-muted hover:bg-surface-3 hover:text-text-main transition-colors cursor-pointer"
              title={revealed ? "Hide key" : "Reveal key"}
              aria-label={revealed ? `Hide ${apiKey.name}` : `Reveal ${apiKey.name}`}
            >
              <span className="material-symbols-outlined text-[15px]">
                {revealed ? "visibility_off" : "visibility"}
              </span>
            </button>
            <button
              type="button"
              onClick={() => onCopy?.(apiKey.key, apiKey.id)}
              className={cn(
                "size-7 rounded-md flex items-center justify-center transition-colors cursor-pointer",
                copied
                  ? "bg-success-bg text-success"
                  : "text-text-muted hover:bg-surface-3 hover:text-primary"
              )}
              title="Copy API key"
              aria-label={`Copy ${apiKey.name}`}
            >
              <span className="material-symbols-outlined text-[15px]">
                {copied ? "check" : "content_copy"}
              </span>
            </button>
          </div>
        </div>
      </div>

      {/* Row 3: Single Clean Quota Telemetry (zero duplication, zero idle pulse) */}
      <div className="mt-3 pl-1">
        {isUnlimited ? (
          <div className="flex flex-wrap items-center justify-between gap-2 text-xs font-mono">
            <span className="inline-flex items-center gap-1.5 text-text-muted">
              <span className="size-1.5 rounded-full bg-success" />
              <span>Unlimited quota</span>
            </span>
            <span className="text-text-subtle">{formatTokensNumber(used)} tokens used</span>
          </div>
        ) : (
          <div className="space-y-1.5">
            <div className="flex flex-wrap items-center justify-between gap-2 text-[11.5px] font-mono">
              <span className="text-text-muted">
                {pct}% used · {formatTokensNumber(used)} / {formatTokensNumber(limit)}
              </span>
              <div className="flex items-center gap-2.5 flex-wrap">
                {apiKey.resetInterval && apiKey.resetInterval !== "never" && (
                  <ResetCountdownInline
                    resetInterval={apiKey.resetInterval}
                    lastResetAt={apiKey.lastResetAt}
                  />
                )}
                <span
                  className={cn(
                    "font-semibold",
                    isExhausted
                      ? "text-danger"
                      : isWarning
                        ? "text-warning"
                        : "text-text-main"
                  )}
                >
                  {isExhausted ? "0 left" : `${formatTokensNumber(remaining)} left`}
                </span>
              </div>
            </div>
            <div className="h-1.5 w-full rounded-full bg-surface-3 overflow-hidden">
              <div
                className={cn(
                  "h-full rounded-full transition-[width] duration-300 ease-out",
                  isExhausted ? "bg-danger" : isWarning ? "bg-warning" : "bg-primary"
                )}
                style={{ width: `${Math.max(pct, used > 0 ? 2 : 0)}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {/* Row 4: Compact Policy & Lifecycle Chips */}
      <div className="mt-3 pl-1 flex flex-wrap items-center gap-1.5 text-[11px] font-mono">
        <span
          className="px-2 py-0.5 rounded-md bg-info-bg text-info-fg border border-info-border max-w-[260px] truncate"
          title={`Allowed models: ${modelsSummary}`}
        >
          Models: {modelsSummary}
        </span>

        {(apiKey.rpmLimit > 0 || apiKey.tpmLimit > 0) && (
          <span className="px-2 py-0.5 rounded-md bg-surface-2 text-text-muted border border-border-subtle">
            {apiKey.rpmLimit > 0 ? `${apiKey.rpmLimit} RPM` : ""}
            {apiKey.rpmLimit > 0 && apiKey.tpmLimit > 0 ? " · " : ""}
            {apiKey.tpmLimit > 0 ? `${formatTokensNumber(apiKey.tpmLimit)} TPM` : ""}
          </span>
        )}

        {whitelistCount > 0 && (
          <span
            className="px-2 py-0.5 rounded-md bg-success-bg text-success-fg border border-success-border"
            title={apiKey.ipWhitelist}
          >
            IP Guard ({whitelistCount})
          </span>
        )}

        {!isUnlimited && apiKey.resetInterval && apiKey.resetInterval !== "never" && (
          <span className="px-2 py-0.5 rounded-md bg-surface-2 text-text-muted border border-border-subtle">
            Reset: every {apiKey.resetInterval}
          </span>
        )}

        {expiry.status !== "none" && (
          <span
            className={cn(
              "px-2 py-0.5 rounded-md border",
              expiry.status === "expired"
                ? "bg-danger-bg text-danger-fg border-danger-border font-semibold"
                : expiry.status === "soon"
                  ? "bg-warning-bg text-warning-fg border-warning-border font-medium"
                  : "bg-surface-2 text-text-muted border-border-subtle"
            )}
            title={apiKey.expiresAt}
          >
            {expiry.label}
          </span>
        )}
      </div>
    </article>
  );
}

ApiKeyCard.propTypes = {
  apiKey: PropTypes.object.isRequired,
  creatorLabel: PropTypes.string,
  canBulkSelect: PropTypes.bool,
  selected: PropTypes.bool,
  onToggleSelect: PropTypes.func,
  isOwnKey: PropTypes.bool,
  revealed: PropTypes.bool,
  onToggleReveal: PropTypes.func,
  copied: PropTypes.bool,
  onCopy: PropTypes.func,
  onToggleActive: PropTypes.func,
  onOpenSnippet: PropTypes.func,
  onEdit: PropTypes.func,
  onDuplicate: PropTypes.func,
  onResetUsage: PropTypes.func,
  onShare: PropTypes.func,
  onDelete: PropTypes.func,
};
