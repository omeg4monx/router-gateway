"use client";

import {
  ChevronRight,
  ChevronsDownUp,
  ChevronsUpDown,
  LogOut,
  Menu,
  PanelLeftClose,
  PanelLeftOpen,
  Search,
  SearchX,
  X,
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { type ReactNode, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useAuthSession, useLogoutMutation } from "@/features/auth/hooks";
import { useSystemHealth, useSystemVersion } from "@/features/system/hooks";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Button } from "./button";
import { CaduceusLogo } from "./CaduceusLogo";
import { EmberRule } from "./EmberRule";
import { LocaleToggle } from "./LocaleToggle";
import {
  countVisibleNavItems,
  EXTENSION_GROUPS,
  findActiveNavContext,
  getShellCopy,
  isNavItemActive,
  type NavAccordionId,
  type NavItem,
  normalizeRoute,
  resolveVisibleNavigation,
} from "./navigation";
import { ThemeToggle } from "./ThemeToggle";

interface AppShellProps {
  children: ReactNode;
}

function createInitialExpandedAccordions(): Record<NavAccordionId, boolean> {
  return EXTENSION_GROUPS.reduce(
    (acc, group) => {
      acc[group.id] = group.defaultExpanded;
      return acc;
    },
    {} as Record<NavAccordionId, boolean>,
  );
}

function isEditableElement(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false;
  }
  const tag = target.tagName;
  return tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT" || target.isContentEditable;
}

export function AppShell({ children }: AppShellProps) {
  const pathname = usePathname();
  const router = useRouter();
  const locale = useUIStore((state) => state.locale);
  const copy = getShellCopy(locale);

  const authSessionQuery = useAuthSession();
  const logoutMutation = useLogoutMutation();
  const healthQuery = useSystemHealth();
  const versionQuery = useSystemVersion();

  const mobileNavOpen = useUIStore((state) => state.mobileNavOpen);
  const setMobileNavOpen = useUIStore((state) => state.setMobileNavOpen);
  const closeMobileNav = useUIStore((state) => state.closeMobileNav);
  const sidebarCollapsed = useUIStore((state) => state.sidebarCollapsed);
  const toggleSidebarCollapsed = useUIStore((state) => state.toggleSidebarCollapsed);

  const [searchQuery, setSearchQuery] = useState("");
  const [expandedAccordions, setExpandedAccordions] = useState<Record<NavAccordionId, boolean>>(
    createInitialExpandedAccordions,
  );
  const desktopSearchRef = useRef<HTMLInputElement>(null);
  const mobileSearchRef = useRef<HTMLInputElement>(null);

  const isAuthenticated = Boolean(authSessionQuery.data?.authenticated);
  const normalizedPathname = normalizeRoute(pathname);
  const activeContext = useMemo(
    () => findActiveNavContext(normalizedPathname),
    [normalizedPathname],
  );

  useEffect(() => {
    const accId = activeContext?.accordionId;
    if (!accId) {
      return;
    }
    setExpandedAccordions((prev) => (prev[accId] ? prev : { ...prev, [accId]: true }));
  }, [activeContext]);

  const focusSidebarSearch = useCallback(() => {
    if (typeof window !== "undefined" && window.innerWidth < 1024) {
      setMobileNavOpen(true);
      window.setTimeout(() => {
        mobileSearchRef.current?.focus();
        mobileSearchRef.current?.select();
      }, 40);
      return;
    }

    if (sidebarCollapsed) {
      toggleSidebarCollapsed();
      window.setTimeout(() => {
        desktopSearchRef.current?.focus();
        desktopSearchRef.current?.select();
      }, 60);
      return;
    }

    desktopSearchRef.current?.focus();
    desktopSearchRef.current?.select();
  }, [setMobileNavOpen, sidebarCollapsed, toggleSidebarCollapsed]);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        focusSidebarSearch();
        return;
      }

      if (event.key === "/" && !event.metaKey && !event.ctrlKey && !event.altKey) {
        if (!isEditableElement(event.target)) {
          event.preventDefault();
          focusSidebarSearch();
        }
        return;
      }

      if (event.key === "Escape") {
        if (searchQuery.trim().length > 0) {
          setSearchQuery("");
        } else if (mobileNavOpen) {
          closeMobileNav();
        }
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [closeMobileNav, focusSidebarSearch, mobileNavOpen, searchQuery]);

  const handleSignOut = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        closeMobileNav();
        router.replace("/login/");
      },
    });
  };

  const toggleAccordion = (id: NavAccordionId) => {
    setExpandedAccordions((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  const isSearching = searchQuery.trim().length > 0;
  const visibleNav = useMemo(
    () =>
      resolveVisibleNavigation({
        query: searchQuery,
        registeredModules: healthQuery.data?.modules ?? null,
      }),
    [searchQuery, healthQuery.data?.modules],
  );
  const totalMatches = useMemo(() => countVisibleNavItems(visibleNav), [visibleNav]);

  const expandedCount = useMemo(
    () => visibleNav.extensions.filter((group) => expandedAccordions[group.id]).length,
    [expandedAccordions, visibleNav.extensions],
  );
  const allAccordionsExpanded =
    visibleNav.extensions.length > 0 && expandedCount === visibleNav.extensions.length;

  const handleToggleAllAccordions = () => {
    const nextExpanded = !allAccordionsExpanded;
    setExpandedAccordions(
      EXTENSION_GROUPS.reduce(
        (acc, group) => {
          acc[group.id] = nextExpanded;
          return acc;
        },
        {} as Record<NavAccordionId, boolean>,
      ),
    );
  };

  const activeCategoryLabel = activeContext
    ? activeContext.sectionLabel[locale]
    : copy.controlPlane;
  const activeRouteLabel = activeContext ? activeContext.item.labels[locale] : copy.fallbackTitle;

  const gatewayStatus = healthQuery.isPending
    ? { label: copy.gatewayChecking, dotClass: "bg-[var(--text-muted)]", pulse: false }
    : healthQuery.isError || !healthQuery.data
      ? { label: copy.gatewayOffline, dotClass: "bg-[var(--status-error)]", pulse: false }
      : healthQuery.data.status === "ok"
        ? { label: copy.gatewayOnline, dotClass: "bg-[var(--status-ok)]", pulse: true }
        : { label: copy.gatewayDegraded, dotClass: "bg-[var(--status-warn)]", pulse: true };

  const versionTag = versionQuery.data?.version
    ? `v${versionQuery.data.version.replace(/^v/i, "")}`
    : "—";

  const renderRailRow = (item: NavItem, compact: boolean) => {
    const Icon = item.icon;
    const isActive = isNavItemActive(normalizedPathname, item.href);
    const label = item.labels[locale];

    return (
      <Link
        key={item.href}
        href={item.href}
        onClick={closeMobileNav}
        title={label}
        aria-label={label}
        aria-current={isActive ? "page" : undefined}
        className={cn(
          "relative flex h-8.5 w-full items-center gap-2.5 rounded-[10px] border transition-[background-color,border-color,color,transform] duration-150 active:scale-[0.98]",
          compact ? "justify-center px-0" : "px-2.5",
          isActive
            ? "border-[var(--color-primary)]/20 bg-[var(--color-primary)]/12 font-medium text-[var(--color-primary)]"
            : "border-transparent text-[var(--color-text-main)]/80 hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)]",
        )}
      >
        <Icon
          className={cn(
            "h-4 w-4 shrink-0",
            isActive ? "text-[var(--color-primary)]" : "text-[var(--color-text-main)]/70",
          )}
          aria-hidden="true"
        />
        {!compact ? (
          <span
            className={cn(
              "flex-1 truncate whitespace-nowrap text-[12.5px] leading-none",
              isActive
                ? "font-semibold text-[var(--color-primary)]"
                : "text-[var(--color-text-main)]/85",
            )}
          >
            {label}
          </span>
        ) : null}
        {isActive ? (
          <span
            aria-hidden="true"
            className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-[var(--color-primary)] shadow-[0_0_10px_color-mix(in_srgb,var(--color-primary)_65%,transparent)]"
          />
        ) : null}
      </Link>
    );
  };

  const renderSidebar = (compact: boolean, surface: "desktop" | "mobile") => {
    const searchInputRef = surface === "mobile" ? mobileSearchRef : desktopSearchRef;

    return (
      <div className="flex h-full min-h-0 flex-col text-[var(--color-text-main)]">
        {/* Brand header (main LeftRail parity) */}
        <div
          className={cn(
            "flex h-13 shrink-0 items-center gap-2 border-b border-[var(--color-border-subtle)]/80 px-2.5",
            compact ? "justify-center" : "justify-between",
          )}
        >
          {!compact ? (
            <>
              <Link
                href="/dashboard/"
                className="flex min-w-0 flex-1 items-center gap-2.5 overflow-hidden transition-opacity hover:opacity-90"
                onClick={closeMobileNav}
                aria-label="OMEG4MON ROUTER Dashboard"
              >
                <div className="flex size-8 shrink-0 items-center justify-center rounded-xl border border-[var(--color-primary)]/25 bg-[var(--color-primary)]/10 p-1 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.25)]">
                  <CaduceusLogo size={20} decorative />
                </div>
                <div className="flex min-w-0 flex-1 flex-col">
                  <span className="truncate text-[11.5px] font-bold uppercase tracking-wider text-[var(--color-text-main)]">
                    OMEG4MON
                  </span>
                  <span className="truncate font-mono text-[9.5px] text-[var(--color-text-muted)]">
                    ROUTER · {versionTag}
                  </span>
                </div>
              </Link>

              {surface === "desktop" ? (
                <button
                  type="button"
                  onClick={toggleSidebarCollapsed}
                  className="flex size-7.5 shrink-0 cursor-pointer items-center justify-center rounded-[9px] border border-transparent text-[var(--color-primary)] transition-colors hover:border-[var(--color-border-subtle)] hover:bg-[var(--color-surface-2)]/80 active:scale-[0.96]"
                  title={copy.collapseSidebarLabel}
                  aria-label={copy.collapseSidebarLabel}
                >
                  <PanelLeftClose className="h-4 w-4" aria-hidden="true" />
                </button>
              ) : (
                <button
                  type="button"
                  onClick={closeMobileNav}
                  className="flex size-7.5 shrink-0 cursor-pointer items-center justify-center rounded-[9px] border border-transparent text-[var(--color-text-muted)] transition-colors hover:border-[var(--color-border-subtle)] hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)] active:scale-[0.96]"
                  aria-label={copy.closeMenuLabel}
                >
                  <X className="h-4 w-4" aria-hidden="true" />
                </button>
              )}
            </>
          ) : (
            <button
              type="button"
              onClick={toggleSidebarCollapsed}
              className="flex size-9 shrink-0 cursor-pointer items-center justify-center rounded-[10px] border border-[var(--color-border-subtle)]/80 text-[var(--color-primary)] transition-colors hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-surface-2)]/80 active:scale-[0.96]"
              title={copy.expandSidebarLabel}
              aria-label={copy.expandSidebarLabel}
            >
              <PanelLeftOpen className="h-4 w-4" aria-hidden="true" />
            </button>
          )}
        </div>

        {/* Inline search bar */}
        <div
          className={cn(
            "shrink-0 border-b border-[var(--color-border-subtle)]/80 py-2",
            compact ? "px-2" : "px-2.5",
          )}
        >
          {compact ? (
            <button
              type="button"
              onClick={focusSidebarSearch}
              title={`${copy.quickSearchLabel} (Ctrl+K)`}
              aria-label={copy.quickSearchLabel}
              className="mx-auto flex size-9 cursor-pointer items-center justify-center rounded-[10px] border border-transparent text-[var(--color-text-main)]/75 transition-colors hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)] active:scale-[0.98]"
            >
              <Search className="h-4 w-4" aria-hidden="true" />
            </button>
          ) : (
            <div className="space-y-1.5">
              <div className="relative">
                <Search
                  className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-subtle)]"
                  aria-hidden="true"
                />
                <input
                  ref={searchInputRef}
                  type="search"
                  value={searchQuery}
                  onChange={(event) => setSearchQuery(event.target.value)}
                  placeholder={copy.searchPlaceholder}
                  aria-label={copy.searchPlaceholder}
                  className="h-8.5 w-full rounded-[10px] border border-[var(--color-field-border)] bg-[var(--color-field)] pl-8 pr-8 text-[12px] text-[var(--color-text-main)] placeholder:text-[var(--color-text-subtle)] transition-colors focus:border-[var(--color-primary)] focus:outline-none"
                />
                {isSearching ? (
                  <button
                    type="button"
                    onClick={() => {
                      setSearchQuery("");
                      searchInputRef.current?.focus();
                    }}
                    aria-label={copy.clearSearchLabel}
                    className="absolute right-2 top-1/2 flex h-5 w-5 -translate-y-1/2 items-center justify-center rounded text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
                  >
                    <X className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                ) : (
                  <kbd className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 rounded border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--color-text-subtle)]">
                    /
                  </kbd>
                )}
              </div>

              {isSearching ? (
                <p className="px-1 font-mono text-[10px] uppercase tracking-wider text-[var(--color-text-subtle)] tabular-nums">
                  {copy.searchResultsLabel(totalMatches)}
                </p>
              ) : null}
            </div>
          )}
        </div>

        {/* Main navigation list (exact main LeftRail section order: Routing -> Operations -> Extensions) */}
        <nav
          className="custom-scrollbar min-h-0 flex-1 space-y-3 overflow-y-auto overflow-x-hidden px-2 py-2.5"
          aria-label="Control plane navigation"
        >
          {totalMatches === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-[var(--color-border)] bg-[var(--color-surface-2)]/40 px-3 py-8 text-center">
              <SearchX className="h-5 w-5 text-[var(--color-text-subtle)]" aria-hidden="true" />
              <p className="mt-2 text-xs font-semibold text-[var(--color-text-main)]">
                {copy.emptySearchTitle}
              </p>
              <p className="mt-1 text-[11px] leading-relaxed text-[var(--color-text-muted)]">
                {copy.emptySearchSubtitle}
              </p>
              <Button
                type="button"
                variant="secondary"
                size="xs"
                onClick={() => {
                  setSearchQuery("");
                  searchInputRef.current?.focus();
                }}
                className="mt-3"
              >
                {copy.resetSearchLabel}
              </Button>
            </div>
          ) : (
            <>
              {/* Section 1: Routing */}
              {visibleNav.routing.length > 0 ? (
                <div className="space-y-0.5">
                  {!compact ? (
                    <div className="px-2.5 pb-1 font-mono text-[10px] font-semibold uppercase tracking-wider text-[var(--color-text-subtle)]">
                      {copy.sectionRouting}
                    </div>
                  ) : null}
                  {visibleNav.routing.map((item) => renderRailRow(item, compact))}
                </div>
              ) : null}

              {/* Section 2: Operations */}
              {visibleNav.operations.length > 0 ? (
                <div
                  className={cn(
                    "space-y-0.5",
                    visibleNav.routing.length > 0 &&
                      "border-t border-[var(--color-border-subtle)]/80 pt-2",
                  )}
                >
                  {!compact ? (
                    <div className="px-2.5 pb-1 font-mono text-[10px] font-semibold uppercase tracking-wider text-[var(--color-text-subtle)]">
                      {copy.sectionOperations}
                    </div>
                  ) : null}
                  {visibleNav.operations.map((item) => renderRailRow(item, compact))}
                </div>
              ) : null}

              {/* Section 3: Extensions (Features & System smooth accordions + Collapse/Expand All) */}
              {visibleNav.extensions.length > 0 ? (
                <div
                  className={cn(
                    "space-y-1",
                    (visibleNav.routing.length > 0 || visibleNav.operations.length > 0) &&
                      "border-t border-[var(--color-border-subtle)]/80 pt-2",
                  )}
                >
                  {!compact ? (
                    <div className="flex items-center justify-between px-2.5 pb-1">
                      <span className="font-mono text-[10px] font-semibold uppercase tracking-wider text-[var(--color-text-subtle)]">
                        {copy.sectionExtensions}
                      </span>
                      {!isSearching ? (
                        <button
                          type="button"
                          onClick={handleToggleAllAccordions}
                          className="inline-flex cursor-pointer items-center gap-1 rounded px-1 py-0.5 text-[10px] font-medium text-[var(--color-text-subtle)] transition-colors hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)]"
                        >
                          {allAccordionsExpanded ? (
                            <>
                              <ChevronsDownUp className="h-3 w-3 shrink-0" aria-hidden="true" />
                              <span>{copy.collapseAllLabel}</span>
                            </>
                          ) : (
                            <>
                              <ChevronsUpDown className="h-3 w-3 shrink-0" aria-hidden="true" />
                              <span>{copy.expandAllLabel}</span>
                            </>
                          )}
                        </button>
                      ) : null}
                    </div>
                  ) : null}

                  {visibleNav.extensions.map((group) => {
                    const GroupIcon = group.icon;
                    const isExpanded = compact || isSearching || expandedAccordions[group.id];
                    const groupHasActive = group.items.some((item) =>
                      isNavItemActive(normalizedPathname, item.href),
                    );

                    return (
                      <div key={group.id}>
                        {!compact ? (
                          <button
                            type="button"
                            onClick={() => toggleAccordion(group.id)}
                            aria-expanded={isExpanded}
                            className={cn(
                              "relative flex h-8.5 w-full cursor-pointer items-center gap-2.5 rounded-[10px] border px-2.5 transition-[background-color,border-color,color,transform] duration-150 active:scale-[0.98]",
                              groupHasActive
                                ? "border-[var(--color-primary)]/20 bg-[var(--color-primary)]/12 font-medium text-[var(--color-primary)]"
                                : isExpanded
                                  ? "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/80 font-medium text-[var(--color-text-main)]"
                                  : "border-transparent text-[var(--color-text-main)]/80 hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)]",
                            )}
                          >
                            <GroupIcon
                              className={cn(
                                "h-4 w-4 shrink-0",
                                groupHasActive || isExpanded
                                  ? "text-[var(--color-primary)]"
                                  : "text-[var(--color-text-main)]/70",
                              )}
                              aria-hidden="true"
                            />
                            <span
                              className={cn(
                                "flex-1 truncate whitespace-nowrap text-left text-[12.5px] leading-none",
                                groupHasActive
                                  ? "font-semibold text-[var(--color-primary)]"
                                  : isExpanded && "font-semibold text-[var(--color-text-main)]",
                              )}
                            >
                              {group.labels[locale]}
                            </span>
                            <ChevronRight
                              className="h-3.5 w-3.5 shrink-0 text-[var(--color-text-subtle)] transition-transform duration-200 ease-[var(--ease-spring)]"
                              style={{ transform: isExpanded ? "rotate(90deg)" : "rotate(0deg)" }}
                              aria-hidden="true"
                            />
                            {groupHasActive ? (
                              <span
                                aria-hidden="true"
                                className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-[var(--color-primary)] shadow-[0_0_10px_color-mix(in_srgb,var(--color-primary)_65%,transparent)]"
                              />
                            ) : null}
                          </button>
                        ) : null}

                        <div
                          className={cn(
                            "grid transition-[grid-template-rows,opacity] duration-200 ease-[var(--ease-spring)]",
                            isExpanded
                              ? "grid-rows-[1fr] opacity-100"
                              : "grid-rows-[0fr] opacity-0",
                          )}
                        >
                          <div className="overflow-hidden">
                            <div
                              className={cn(
                                "space-y-0.5",
                                !compact &&
                                  "ml-3.5 mt-1 border-l border-[var(--color-border-subtle)] pl-2 pb-0.5",
                              )}
                            >
                              {group.items.map((item) => {
                                const Icon = item.icon;
                                const isActive = isNavItemActive(normalizedPathname, item.href);
                                const itemLabel = item.labels[locale];

                                return (
                                  <Link
                                    key={item.href}
                                    href={item.href}
                                    onClick={closeMobileNav}
                                    tabIndex={isExpanded ? undefined : -1}
                                    title={itemLabel}
                                    aria-current={isActive ? "page" : undefined}
                                    className={cn(
                                      "relative flex items-center gap-2.5 border transition-[background-color,border-color,color,transform] duration-150 active:scale-[0.98]",
                                      compact
                                        ? "h-8.5 justify-center rounded-[10px] px-0"
                                        : "rounded-[8px] px-2.5 py-2",
                                      isActive
                                        ? "border-[var(--color-primary)]/20 bg-[var(--color-primary)]/12 font-medium text-[var(--color-primary)]"
                                        : "border-transparent text-[var(--color-text-main)]/80 hover:bg-[var(--color-surface-2)]/80 hover:text-[var(--color-text-main)]",
                                    )}
                                  >
                                    <Icon
                                      className={cn(
                                        "shrink-0",
                                        compact ? "h-4 w-4" : "h-3.5 w-3.5",
                                        isActive
                                          ? "text-[var(--color-primary)]"
                                          : "text-[var(--color-text-main)]/70",
                                      )}
                                      aria-hidden="true"
                                    />
                                    {!compact ? (
                                      <span
                                        className={cn(
                                          "flex-1 truncate text-[12px] leading-none",
                                          isActive && "font-semibold text-[var(--color-primary)]",
                                        )}
                                      >
                                        {itemLabel}
                                      </span>
                                    ) : null}
                                    {isActive && compact ? (
                                      <span
                                        aria-hidden="true"
                                        className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-[var(--color-primary)] shadow-[0_0_10px_color-mix(in_srgb,var(--color-primary)_65%,transparent)]"
                                      />
                                    ) : null}
                                  </Link>
                                );
                              })}
                            </div>
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              ) : null}
            </>
          )}
        </nav>

        {/* Docked Settings row (exact main LeftRail placement) */}
        {visibleNav.settings ? (
          <div className="shrink-0 border-t border-[var(--color-border-subtle)]/80 px-2 py-1.5">
            {renderRailRow(visibleNav.settings, compact)}
          </div>
        ) : null}

        {/* Footer status bar (main LeftRail parity) */}
        <div
          className={cn(
            "flex shrink-0 gap-1.5 rounded-b-2xl border-t border-[var(--color-border-subtle)]/80 bg-[var(--color-surface-2)]/45 p-2.5",
            compact ? "flex-col items-center justify-center" : "items-center justify-between",
          )}
        >
          {compact ? (
            <div
              className="flex items-center justify-center"
              title={`Gateway ${gatewayStatus.label} · ${versionTag}`}
            >
              <span className="relative flex h-2.5 w-2.5">
                {gatewayStatus.pulse ? (
                  <span
                    className={cn(
                      "absolute inline-flex h-full w-full animate-ping rounded-full opacity-55",
                      gatewayStatus.dotClass,
                    )}
                  />
                ) : null}
                <span
                  className={cn(
                    "relative inline-flex h-2.5 w-2.5 rounded-full",
                    gatewayStatus.dotClass,
                  )}
                />
              </span>
            </div>
          ) : (
            <>
              <div className="flex min-w-0 flex-col">
                <span className="truncate text-[11.5px] font-semibold text-[var(--color-text-main)]">
                  Instance Admin
                </span>
                <span className="flex items-center gap-1 whitespace-nowrap font-mono text-[10px] text-[var(--status-ok)]">
                  <span
                    className={cn("size-1.5 shrink-0 rounded-full", gatewayStatus.dotClass)}
                    aria-hidden="true"
                  />
                  <span>{gatewayStatus.label}</span>
                </span>
              </div>
              <span className="rounded-md border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] tabular-nums text-[var(--color-text-muted)]">
                {versionTag}
              </span>
            </>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="flex min-h-dvh bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      {/* Floating rounded desktop studio rail (pinned to viewport while main content scrolls) */}
      <aside
        className={cn(
          "shell-rail sticky top-2.5 isolate z-[45] my-2.5 ml-2.5 hidden h-[calc(100dvh-20px)] shrink-0 self-start select-none flex-col overflow-hidden rounded-2xl transition-[width] duration-250 ease-[var(--ease-spring)] lg:flex",
          sidebarCollapsed ? "w-[60px]" : "w-[244px]",
        )}
        aria-label="Application navigation"
      >
        <div className="flex h-full min-h-0 flex-col">
          {renderSidebar(sidebarCollapsed, "desktop")}
        </div>
      </aside>

      {/* Mobile navigation drawer with translucent scrim */}
      {mobileNavOpen ? (
        <div
          className="fixed inset-0 z-50 flex p-2.5 lg:hidden"
          role="dialog"
          aria-modal="true"
          aria-label="Navigation"
        >
          <div className="shell-rail relative z-10 w-[264px] max-w-[85vw] overflow-hidden rounded-2xl">
            {renderSidebar(false, "mobile")}
          </div>
          <button
            type="button"
            className="fixed inset-0 bg-black/45 backdrop-blur-[2px]"
            aria-label={copy.dismissMenuLabel}
            onClick={closeMobileNav}
          />
        </div>
      ) : null}

      {/* Main content area */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-20 bg-[var(--bg-canvas)]/90 backdrop-blur-md">
          <div className="flex h-16 items-center justify-between gap-3 px-4 sm:px-6 lg:px-8">
            <div className="flex min-w-0 items-center gap-3">
              <Button
                type="button"
                variant="secondary"
                size="icon"
                onClick={() => setMobileNavOpen(true)}
                className="h-9 w-9 shrink-0 lg:hidden"
                aria-label={copy.openMenuLabel}
              >
                <Menu className="h-4 w-4" aria-hidden="true" />
              </Button>

              <div className="min-w-0">
                <div className="flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-[0.14em] text-[var(--text-muted)]">
                  <span className="hidden sm:inline">{copy.controlPlane}</span>
                  <ChevronRight className="hidden h-3 w-3 shrink-0 sm:inline" aria-hidden="true" />
                  <span className="truncate text-[var(--accent-text)]">{activeCategoryLabel}</span>
                </div>
                <h1 className="truncate text-base font-semibold tracking-tight text-[var(--text-primary)]">
                  {activeRouteLabel}
                </h1>
              </div>
            </div>

            <div className="flex shrink-0 items-center gap-2">
              <button
                type="button"
                onClick={focusSidebarSearch}
                className="hidden h-9 cursor-pointer items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-3 text-xs font-medium text-[var(--text-muted)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)] md:inline-flex"
                aria-label={`${copy.quickSearchLabel} (Ctrl+K)`}
              >
                <Search className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{copy.quickSearchLabel}</span>
                <kbd className="rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--text-muted)]">
                  Ctrl K
                </kbd>
              </button>

              <LocaleToggle />
              <ThemeToggle />

              {isAuthenticated ? (
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={handleSignOut}
                  disabled={logoutMutation.isPending}
                  title={logoutMutation.error?.message}
                  className="h-9"
                >
                  <LogOut className="h-3.5 w-3.5" aria-hidden="true" />
                  <span className="hidden sm:inline">
                    {logoutMutation.isPending ? copy.signingOutLabel : copy.signOutLabel}
                  </span>
                </Button>
              ) : null}
            </div>
          </div>
          <EmberRule />
        </header>

        <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 lg:px-8">
          {children}
        </main>
      </div>
    </div>
  );
}
