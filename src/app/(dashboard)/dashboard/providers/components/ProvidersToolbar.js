"use client";

import { useEffect, useRef } from "react";
import PropTypes from "prop-types";
import { cn } from "@/shared/utils/cn";
import { PROVIDER_STATUS_PILLS, CATEGORY_TABS } from "../utils";

export default function ProvidersToolbar({
  search,
  onSearchChange,
  statusFilter,
  onStatusFilterChange,
  statusCounts,
  categoryTab,
  onCategoryTabChange,
  categoryCounts,
  hasActiveFilters,
  onResetFilters,
}) {
  const searchInputRef = useRef(null);

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key !== "/" || e.metaKey || e.ctrlKey || e.altKey) return;
      const target = e.target;
      if (
        target instanceof HTMLElement &&
        (target.tagName === "INPUT" ||
          target.tagName === "TEXTAREA" ||
          target.tagName === "SELECT" ||
          target.isContentEditable)
      ) {
        return;
      }
      if (document.querySelector('[role="dialog"][aria-modal="true"]')) {
        return;
      }
      e.preventDefault();
      searchInputRef.current?.focus();
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, []);

  const handleSearchKeyDown = (e) => {
    if (e.key === "Escape") {
      if (search) {
        onSearchChange("");
      } else {
        searchInputRef.current?.blur();
      }
    }
  };

  return (
    <div className="flex flex-col gap-3 rounded-2xl border border-border bg-surface p-3.5 shadow-[var(--shadow-soft)]">
      {/* Row 1: Search & Status Pills */}
      <div className="flex flex-col gap-2.5 lg:flex-row lg:items-center lg:justify-between">
        <div className="relative flex-1 max-w-md">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-[17px] pointer-events-none">
            search
          </span>
          <input
            ref={searchInputRef}
            type="text"
            placeholder="Search by name, ID, alias, or protocol..."
            value={search}
            onChange={(e) => onSearchChange(e.target.value)}
            onKeyDown={handleSearchKeyDown}
            aria-label="Search providers by name, ID, alias, or protocol"
            className="w-full h-8.5 pl-9 pr-9 bg-surface-2/80 border border-border rounded-lg text-xs text-text-main placeholder:text-text-muted focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/35 transition-colors"
          />
          {search ? (
            <button
              type="button"
              onClick={() => onSearchChange("")}
              className="absolute right-2 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-main p-1 rounded-md hover:bg-surface-3 transition-colors cursor-pointer"
              aria-label="Clear search"
            >
              <span className="material-symbols-outlined text-[15px] block">
                close
              </span>
            </button>
          ) : (
            <kbd
              aria-hidden="true"
              className="hidden sm:inline-flex items-center justify-center absolute right-2.5 top-1/2 -translate-y-1/2 px-1.5 py-0.5 rounded border border-border-subtle bg-surface text-[10px] font-mono text-text-muted pointer-events-none"
            >
              /
            </kbd>
          )}
        </div>

        {/* Status Filter Pills */}
        <div
          className="flex items-center gap-1.5 overflow-x-auto custom-scrollbar"
          role="group"
          aria-label="Filter providers by status"
        >
          {PROVIDER_STATUS_PILLS.map((pill) => {
            const count = statusCounts?.[pill.value] ?? 0;
            const active = statusFilter === pill.value;
            const isErrorAttention = pill.value === "error" && count > 0;
            return (
              <button
                key={pill.value}
                type="button"
                onClick={() => onStatusFilterChange(pill.value)}
                aria-pressed={active}
                className={cn(
                  "inline-flex items-center gap-1.5 h-7.5 px-2.5 rounded-lg text-xs font-medium whitespace-nowrap border transition-colors shrink-0 cursor-pointer",
                  active
                    ? pill.value === "error"
                      ? "bg-danger text-white border-danger shadow-xs"
                      : "bg-primary text-white border-primary shadow-xs"
                    : isErrorAttention
                      ? "bg-danger/8 border-danger/30 text-danger hover:bg-danger/15"
                      : "bg-surface-2/70 border-border text-text-muted hover:text-text-main hover:bg-surface-3",
                )}
              >
                <span>{pill.label}</span>
                <span
                  className={cn(
                    "text-[10px] font-mono tabular-nums",
                    active ? "opacity-90" : "opacity-75",
                  )}
                >
                  {count}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Row 2: Category Tabs & Reset */}
      <div className="flex items-center justify-between gap-2 pt-2.5 border-t border-border-subtle">
        <div
          className="flex items-center gap-1.5 overflow-x-auto custom-scrollbar text-xs"
          role="tablist"
          aria-label="Provider categories"
        >
          {CATEGORY_TABS.map((tab) => {
            const active = categoryTab === tab.value;
            const count = categoryCounts?.[tab.value] ?? 0;
            return (
              <button
                key={tab.value}
                type="button"
                role="tab"
                aria-selected={active}
                onClick={() => onCategoryTabChange(tab.value)}
                className={cn(
                  "inline-flex items-center gap-1.5 h-7.5 px-3 rounded-lg font-medium whitespace-nowrap border transition-colors shrink-0 cursor-pointer",
                  active
                    ? "bg-primary/12 text-primary border-primary/30 font-semibold"
                    : "bg-transparent border-transparent text-text-muted hover:text-text-main hover:bg-surface-2",
                )}
              >
                <span className="material-symbols-outlined text-[15px]">
                  {tab.icon}
                </span>
                <span>{tab.label}</span>
                <span className="text-[10px] font-mono tabular-nums opacity-80">
                  ({count})
                </span>
              </button>
            );
          })}
        </div>

        {hasActiveFilters && (
          <button
            type="button"
            onClick={onResetFilters}
            className="inline-flex items-center gap-1 h-7.5 px-2.5 rounded-lg text-xs font-medium text-text-muted hover:text-text-main hover:bg-surface-2 transition-colors shrink-0 cursor-pointer"
          >
            <span className="material-symbols-outlined text-[14px]">
              filter_alt_off
            </span>
            <span>Reset</span>
          </button>
        )}
      </div>
    </div>
  );
}

ProvidersToolbar.propTypes = {
  search: PropTypes.string.isRequired,
  onSearchChange: PropTypes.func.isRequired,
  statusFilter: PropTypes.string.isRequired,
  onStatusFilterChange: PropTypes.func.isRequired,
  statusCounts: PropTypes.object,
  categoryTab: PropTypes.string.isRequired,
  onCategoryTabChange: PropTypes.func.isRequired,
  categoryCounts: PropTypes.object,
  hasActiveFilters: PropTypes.bool,
  onResetFilters: PropTypes.func.isRequired,
};
