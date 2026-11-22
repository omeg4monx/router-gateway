"use client";

/**
 * ModelAvailabilityBadge — compact inline status indicator for Studio Header.
 *
 * Shows healthy state when all models are operational, or an amber issue pill
 * when models are cooling down / unavailable, with a shell-popover for details
 * and one-click cooldown clearing.
 */

import { useState, useEffect, useCallback, useRef } from "react";
import PropTypes from "prop-types";
import { Button } from "@/shared/components";
import { cn } from "@/shared/utils/cn";
import { useNotificationStore } from "@/store/notificationStore";

const STATUS_CONFIG = {
  available: { icon: "check_circle", color: "text-success", label: "Available" },
  cooldown: { icon: "schedule", color: "text-warning", label: "Cooldown" },
  unavailable: { icon: "error", color: "text-danger", label: "Unavailable" },
  unknown: { icon: "help", color: "text-text-muted", label: "Unknown" },
};

export default function ModelAvailabilityBadge({ onCooldownCleared }) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [clearing, setClearing] = useState(null);
  const ref = useRef(null);
  const notify = useNotificationStore();

  const fetchStatus = useCallback(async () => {
    try {
      const res = await fetch("/api/models/availability", { cache: "no-store" });
      if (res.ok) {
        const json = await res.json();
        setData(json);
      }
    } catch {
      // silent retry on next interval
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(() => {
      if (typeof document !== "undefined" && document.hidden) return;
      fetchStatus();
    }, 30000);
    return () => clearInterval(interval);
  }, [fetchStatus]);

  // Close popover on outside click or Escape key
  useEffect(() => {
    if (!expanded) return;
    const handleMouseDown = (e) => {
      if (ref.current && !ref.current.contains(e.target)) setExpanded(false);
    };
    const handleKeyDown = (e) => {
      if (e.key === "Escape") setExpanded(false);
    };
    document.addEventListener("mousedown", handleMouseDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("mousedown", handleMouseDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [expanded]);

  const handleClearCooldown = async (provider, model) => {
    setClearing(`${provider}:${model}`);
    try {
      const res = await fetch("/api/models/availability", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ action: "clearCooldown", provider, model }),
      });
      if (res.ok) {
        notify.success(`Cooldown cleared for ${model}`);
        await fetchStatus();
        if (onCooldownCleared) {
          await onCooldownCleared();
        }
      } else {
        notify.error("Failed to clear cooldown");
      }
    } catch {
      notify.error("Failed to clear cooldown");
    } finally {
      setClearing(null);
    }
  };

  if (loading) return null;

  const models = data?.models || [];
  const unavailableCount =
    data?.unavailableCount ??
    models.filter((m) => m.status !== "available").length;
  const isHealthy = unavailableCount === 0;

  const byProvider = {};
  models.forEach((m) => {
    if (m.status === "available") return;
    const key = m.provider || "unknown";
    if (!byProvider[key]) byProvider[key] = [];
    byProvider[key].push(m);
  });

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setExpanded((v) => !v)}
        aria-expanded={expanded}
        aria-haspopup="dialog"
        className={cn(
          "inline-flex items-center gap-1.5 h-8 px-2.5 rounded-lg text-xs font-medium border transition-colors cursor-pointer whitespace-nowrap",
          isHealthy
            ? "bg-success-bg/70 border-success-border text-success hover:bg-success-bg"
            : "bg-warning-bg/80 border-warning-border text-warning hover:bg-warning-bg",
        )}
      >
        <span className="material-symbols-outlined text-[15px]">
          {isHealthy ? "verified" : "warning"}
        </span>
        <span>
          {isHealthy
            ? "All models healthy"
            : `${unavailableCount} model ${unavailableCount === 1 ? "issue" : "issues"}`}
        </span>
      </button>

      {expanded && (
        <div
          role="dialog"
          aria-label="Model availability status"
          className="absolute top-full right-0 mt-2 w-80 max-w-[calc(100vw-2rem)] rounded-xl border border-border shell-popover popover-in z-50 overflow-hidden"
        >
          <div className="flex items-center justify-between px-3.5 py-2.5 border-b border-border-subtle bg-surface-2/60">
            <div className="flex items-center gap-2">
              <span
                className={cn(
                  "material-symbols-outlined text-[16px]",
                  isHealthy ? "text-success" : "text-warning",
                )}
              >
                {isHealthy ? "verified" : "warning"}
              </span>
              <span className="text-xs font-semibold text-text-main">
                Model Availability
              </span>
            </div>
            <button
              type="button"
              onClick={fetchStatus}
              className="p-1 rounded-lg hover:bg-surface text-text-muted hover:text-text-main transition-colors cursor-pointer"
              title="Refresh model status"
              aria-label="Refresh model status"
            >
              <span className="material-symbols-outlined text-[15px] block">
                refresh
              </span>
            </button>
          </div>

          <div className="px-3.5 py-3 max-h-64 overflow-y-auto">
            {isHealthy ? (
              <p className="text-xs text-text-muted text-center py-2">
                All models are responding normally with no active cooldowns.
              </p>
            ) : (
              <div className="flex flex-col gap-3">
                {Object.entries(byProvider).map(([provider, provModels]) => (
                  <div key={provider}>
                    <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted mb-1.5">
                      {provider}
                    </p>
                    <div className="flex flex-col gap-1.5">
                      {provModels.map((m) => {
                        const status =
                          STATUS_CONFIG[m.status] || STATUS_CONFIG.unknown;
                        const isClearing =
                          clearing === `${m.provider}:${m.model}`;
                        return (
                          <div
                            key={`${m.provider}-${m.model}-${m.connectionId || ""}`}
                            className="flex items-center justify-between gap-2 px-2.5 py-1.5 rounded-lg border border-border-subtle bg-surface-2/60"
                          >
                            <div className="flex items-center gap-1.5 min-w-0">
                              <span
                                className={cn(
                                  "material-symbols-outlined text-[14px] shrink-0",
                                  status.color,
                                )}
                              >
                                {status.icon}
                              </span>
                              <span
                                className="font-mono text-xs text-text-main truncate"
                                title={m.model}
                              >
                                {m.model}
                              </span>
                            </div>
                            {m.status === "cooldown" && (
                              <Button
                                size="sm"
                                variant="ghost"
                                onClick={() =>
                                  handleClearCooldown(m.provider, m.model)
                                }
                                disabled={isClearing}
                                className="text-[11px] h-6 px-2 shrink-0"
                              >
                                {isClearing ? "..." : "Clear"}
                              </Button>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

ModelAvailabilityBadge.propTypes = {
  onCooldownCleared: PropTypes.func,
};
