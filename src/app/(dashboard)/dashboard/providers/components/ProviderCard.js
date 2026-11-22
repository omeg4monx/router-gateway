"use client";

import Link from "next/link";
import PropTypes from "prop-types";
import { Badge, Toggle } from "@/shared/components";
import ProviderIcon from "@/shared/components/ProviderIcon";
import { getProviderIconSrc } from "@/shared/utils/providerIcon";
import { getCustomLogo } from "@/shared/utils/providerLogo";
import { cn } from "@/shared/utils/cn";
import { isMoonshotCompatible } from "../utils";

function renderStatusDisplay({
  connected,
  error,
  errorCode,
  total,
  isCustom,
}) {
  const parts = [];
  if (connected > 0) {
    parts.push(
      <Badge key="connected" variant="success" size="sm" dot>
        {connected} Active
      </Badge>,
    );
  }
  if (error > 0) {
    const errText = errorCode
      ? `${error} Error (${errorCode})`
      : `${error} Error`;
    parts.push(
      <Badge key="error" variant="error" size="sm" dot>
        {errText}
      </Badge>,
    );
  }
  if (parts.length > 0) return parts;

  if (total > 0) {
    return (
      <Badge variant="default" size="sm" dot>
        {total} Configured
      </Badge>
    );
  }

  if (isCustom) {
    return (
      <span className="text-[11px] text-primary/90 group-hover:text-primary transition-colors flex items-center gap-1 font-medium">
        <span className="material-symbols-outlined text-[13px]">key</span>
        Add API Key
      </span>
    );
  }

  return (
    <span className="text-[11px] text-text-muted group-hover:text-primary transition-colors flex items-center gap-1 font-medium">
      <span className="material-symbols-outlined text-[13px]">add_circle</span>
      Setup
    </span>
  );
}

export default function ProviderCard({
  entry,
  onToggle,
  onQuickTest,
  isToggling = false,
  isTesting = false,
  disableQuickTest = false,
}) {
  const {
    id: providerId,
    info: provider,
    stats,
    typeLabel,
    isNoAuth,
    isCompatible,
    isAnthropicCompatible,
    category,
  } = entry;

  const { connected, error, total, errorCode, errorTime, allDisabled } = stats;
  const isCustom = category === "custom";
  const hasError = !allDisabled && error > 0;
  const isActiveHealthy = !allDisabled && connected > 0 && !hasError;
  const isMoonshot = isMoonshotCompatible(provider, providerId, isCompatible);

  const getIconPath = () => {
    const customLogo = getCustomLogo(provider);
    if (customLogo) return customLogo;
    if (isMoonshot) return "/providers/custom-moonshot.svg";
    if (isCompatible) return "/providers/custom-openai.svg";
    if (isAnthropicCompatible) return "/providers/custom-anthropic.svg";
    return getProviderIconSrc(provider.id || providerId);
  };

  const canQuickTest = total > 0 && !allDisabled && Boolean(onQuickTest);

  return (
    <Link
      href={`/dashboard/providers/${providerId}`}
      className="group min-w-0 block h-full focus:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 rounded-xl"
    >
      <div
        className={cn(
          "relative flex flex-col justify-between h-full p-4 rounded-xl border bg-surface transition-all duration-150 overflow-hidden shadow-[var(--shadow-soft)]",
          allDisabled
            ? "opacity-65 border-border-subtle bg-surface-2/60 hover:opacity-85"
            : hasError
              ? "border-danger/35 hover:border-danger/60 hover:bg-surface-2/60"
              : "border-border hover:border-primary/45 hover:bg-surface-2/60",
        )}
      >
        {/* Subtle top hairline accent for active or errored providers */}
        <span
          className={cn(
            "absolute top-0 inset-x-0 h-[2px] transition-colors",
            hasError
              ? "bg-danger"
              : isActiveHealthy
                ? "bg-primary"
                : "bg-transparent group-hover:bg-border",
          )}
          aria-hidden="true"
        />

        {/* Top row: Icon + Name + Protocol + Quick Controls */}
        <div className="flex items-start justify-between gap-2.5 min-w-0">
          <div className="flex items-center gap-3 min-w-0">
            <div
              className="size-10 shrink-0 rounded-xl border bg-surface-2 dark:bg-[#1a1a20] flex items-center justify-center p-1.5 transition-transform duration-150 group-hover:scale-[1.03]"
              style={{
                borderColor:
                  isCompatible || isAnthropicCompatible
                    ? "color-mix(in srgb, var(--color-primary) 40%, transparent)"
                    : provider.color
                      ? `${provider.color}45`
                      : undefined,
              }}
            >
              <ProviderIcon
                src={getIconPath()}
                alt={provider.name}
                size={28}
                className="object-contain max-w-[28px] max-h-[28px]"
                fallbackText={
                  provider.textIcon ||
                  (provider.id || providerId).slice(0, 2).toUpperCase()
                }
                fallbackColor={provider.color}
              />
            </div>

            <div className="min-w-0">
              <h3 className="truncate font-semibold text-sm text-text-main group-hover:text-primary transition-colors">
                {provider.name}
              </h3>
              <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted mt-0.5 block truncate">
                {typeLabel}
              </span>
            </div>
          </div>

          {total > 0 && (
            <div className="flex items-center gap-1.5 shrink-0">
              {canQuickTest && (
                <button
                  type="button"
                  disabled={disableQuickTest}
                  onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    onQuickTest(entry);
                  }}
                  title={`Test ${provider.name} connections`}
                  aria-label={`Test ${provider.name} connections`}
                  className={cn(
                    "size-7 rounded-lg border flex items-center justify-center transition-colors cursor-pointer",
                    isTesting
                      ? "bg-primary/15 border-primary/40 text-primary"
                      : "bg-surface-2/80 border-border text-text-muted hover:text-primary hover:border-primary/40 hover:bg-primary/8",
                    disableQuickTest && !isTesting && "opacity-45 cursor-not-allowed",
                  )}
                >
                  <span
                    className={cn(
                      "material-symbols-outlined text-[15px]",
                      isTesting && "animate-spin",
                    )}
                  >
                    {isTesting ? "progress_activity" : "play_arrow"}
                  </span>
                </button>
              )}

              <div
                className="shrink-0"
                title={allDisabled ? "Enable provider" : "Disable provider"}
                onClick={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                }}
              >
                <Toggle
                  size="sm"
                  checked={!allDisabled}
                  disabled={isToggling}
                  onChange={(nextActive) => {
                    if (!isToggling && onToggle) onToggle(entry, nextActive);
                  }}
                />
              </div>
            </div>
          )}
        </div>

        {/* Bottom row: Status + Arrow */}
        <div className="mt-3.5 pt-2.5 border-t border-border-subtle flex items-center justify-between gap-2 text-xs">
          <div className="flex items-center gap-1.5 min-w-0 flex-wrap">
            {allDisabled ? (
              <Badge variant="default" size="sm">
                <span className="flex items-center gap-1">
                  <span className="material-symbols-outlined text-[12px]">
                    pause_circle
                  </span>
                  Disabled
                </span>
              </Badge>
            ) : isNoAuth && total === 0 ? (
              <Badge
                variant={provider.statusBadge?.variant || "success"}
                size="sm"
                dot
              >
                {provider.statusBadge?.label || "Ready"}
              </Badge>
            ) : (
              <>
                {renderStatusDisplay({
                  connected,
                  error,
                  errorCode,
                  total,
                  isCustom,
                })}
                {errorTime && hasError && (
                  <span className="text-text-muted text-[10px] font-mono">
                    {errorTime}
                  </span>
                )}
              </>
            )}
          </div>

          <span className="material-symbols-outlined text-[15px] text-text-muted opacity-40 group-hover:opacity-100 group-hover:text-primary group-hover:translate-x-0.5 transition-all shrink-0">
            arrow_forward
          </span>
        </div>
      </div>
    </Link>
  );
}

ProviderCard.propTypes = {
  entry: PropTypes.shape({
    id: PropTypes.string.isRequired,
    info: PropTypes.shape({
      id: PropTypes.string,
      name: PropTypes.string.isRequired,
      color: PropTypes.string,
      textIcon: PropTypes.string,
      statusBadge: PropTypes.shape({
        variant: PropTypes.string,
        label: PropTypes.string,
      }),
    }).isRequired,
    stats: PropTypes.shape({
      connected: PropTypes.number,
      error: PropTypes.number,
      total: PropTypes.number,
      errorCode: PropTypes.string,
      errorTime: PropTypes.string,
      allDisabled: PropTypes.bool,
    }).isRequired,
    typeLabel: PropTypes.string,
    isNoAuth: PropTypes.bool,
    isCompatible: PropTypes.bool,
    isAnthropicCompatible: PropTypes.bool,
    category: PropTypes.string,
  }).isRequired,
  onToggle: PropTypes.func,
  onQuickTest: PropTypes.func,
  isToggling: PropTypes.bool,
  isTesting: PropTypes.bool,
  disableQuickTest: PropTypes.bool,
};
