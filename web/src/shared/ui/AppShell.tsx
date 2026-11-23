"use client";

import {
  ChevronDown,
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
  filterNavGroups,
  findActiveNavContext,
  getShellCopy,
  isNavItemActive,
  NAV_GROUPS,
  type NavGroupId,
  normalizeRoute,
} from "./navigation";
import { ThemeToggle } from "./ThemeToggle";

interface AppShellProps {
  children: ReactNode;
}

function createInitialExpandedGroups(): Record<NavGroupId, boolean> {
  return NAV_GROUPS.reduce(
    (acc, group) => {
      acc[group.id] = group.defaultExpanded;
      return acc;
    },
    {} as Record<NavGroupId, boolean>,
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
  const [expandedGroups, setExpandedGroups] = useState<Record<NavGroupId, boolean>>(
    createInitialExpandedGroups,
  );
  const desktopSearchRef = useRef<HTMLInputElement>(null);
  const mobileSearchRef = useRef<HTMLInputElement>(null);

  const isAuthenticated = Boolean(authSessionQuery.data?.authenticated);
  const normalizedPathname = normalizeRoute(pathname);
  const activeContext = useMemo(
    () => findActiveNavContext(normalizedPathname),
    [normalizedPathname],
  );

  // Keep the active route's category expanded automatically on navigation.
  useEffect(() => {
    if (!activeContext) {
      return;
    }
    setExpandedGroups((prev) =>
      prev[activeContext.group.id] ? prev : { ...prev, [activeContext.group.id]: true },
    );
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

  const toggleGroup = (groupId: NavGroupId) => {
    setExpandedGroups((prev) => ({ ...prev, [groupId]: !prev[groupId] }));
  };

  const expandedCount = useMemo(
    () => NAV_GROUPS.filter((group) => expandedGroups[group.id]).length,
    [expandedGroups],
  );
  const allGroupsExpanded = expandedCount === NAV_GROUPS.length;

  const handleToggleAllGroups = () => {
    const nextExpanded = !allGroupsExpanded;
    setExpandedGroups(
      NAV_GROUPS.reduce(
        (acc, group) => {
          acc[group.id] = nextExpanded;
          return acc;
        },
        {} as Record<NavGroupId, boolean>,
      ),
    );
  };

  const isSearching = searchQuery.trim().length > 0;
  const visibleGroups = useMemo(() => filterNavGroups(searchQuery, locale), [searchQuery, locale]);
  const totalMatches = useMemo(
    () => visibleGroups.reduce((sum, group) => sum + group.items.length, 0),
    [visibleGroups],
  );

  const activeCategoryLabel = activeContext
    ? activeContext.group.labels[locale]
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

  const renderSidebar = (compact: boolean, surface: "desktop" | "mobile") => {
    const searchInputRef = surface === "mobile" ? mobileSearchRef : desktopSearchRef;

    return (
      <div className="flex h-full flex-col bg-[var(--color-sidebar)] text-[var(--text-primary)] backdrop-blur-sm">
        {/* Brand header (h-16 aligned with topbar) */}
        <div
          className={cn(
            "flex h-16 shrink-0 items-center bg-[var(--bg-surface)]/65",
            compact ? "justify-center px-3" : "justify-between px-4",
          )}
        >
          <Link
            href="/dashboard/"
            className="flex items-center gap-3 rounded-md focus-visible:outline-none"
            onClick={closeMobileNav}
            aria-label="OMEG4MON ROUTER Dashboard"
          >
            <CaduceusLogo size={26} decorative />
            {!compact ? (
              <div className="flex flex-col leading-tight">
                <span className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                  OMEG4MON
                </span>
                <span className="font-mono text-[10px] tracking-[0.22em] text-[var(--text-muted)]">
                  ROUTER
                </span>
              </div>
            ) : null}
          </Link>
          {!compact ? (
            <Button
              type="button"
              variant="secondary"
              size="icon"
              onClick={closeMobileNav}
              className="h-8 w-8 lg:hidden"
              aria-label={copy.closeMenuLabel}
            >
              <X className="h-4 w-4" aria-hidden="true" />
            </Button>
          ) : null}
        </div>

        <EmberRule />

        {/* Sidebar search control */}
        <div
          className={cn(
            "shrink-0 border-b border-[var(--border-strong)]/70 py-3",
            compact ? "px-3" : "px-3.5",
          )}
        >
          {compact ? (
            <button
              type="button"
              onClick={focusSidebarSearch}
              title={`${copy.quickSearchLabel} (Ctrl+K)`}
              aria-label={copy.quickSearchLabel}
              className="mx-auto flex h-9 w-9 items-center justify-center rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-secondary)] shadow-2xs transition-colors hover:border-[var(--accent-primary)]/50 hover:text-[var(--text-primary)]"
            >
              <Search className="h-4 w-4" aria-hidden="true" />
            </button>
          ) : (
            <div className="space-y-1.5">
              <div className="relative">
                <Search
                  className="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
                <input
                  ref={searchInputRef}
                  type="search"
                  value={searchQuery}
                  onChange={(event) => setSearchQuery(event.target.value)}
                  placeholder={copy.searchPlaceholder}
                  aria-label={copy.searchPlaceholder}
                  className="h-9 w-full rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] pl-8 pr-9 text-xs text-[var(--text-primary)] shadow-2xs placeholder:text-[var(--text-muted)] transition-colors focus:border-[var(--accent-primary)] focus:outline-none"
                />
                {isSearching ? (
                  <button
                    type="button"
                    onClick={() => {
                      setSearchQuery("");
                      searchInputRef.current?.focus();
                    }}
                    aria-label={copy.clearSearchLabel}
                    className="absolute right-2 top-1/2 flex h-5 w-5 -translate-y-1/2 items-center justify-center rounded text-[var(--text-muted)] transition-colors hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                  >
                    <X className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                ) : (
                  <kbd className="pointer-events-none absolute right-2 top-1/2 -translate-y-1/2 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--text-muted)]">
                    /
                  </kbd>
                )}
              </div>

              <div className="flex items-center justify-between gap-2 px-0.5 pt-0.5">
                <span className="font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)] tabular-nums">
                  {isSearching
                    ? copy.searchResultsLabel(totalMatches)
                    : copy.categoriesLabel(expandedCount, NAV_GROUPS.length)}
                </span>
                {!isSearching ? (
                  <button
                    type="button"
                    onClick={handleToggleAllGroups}
                    className="inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 text-[11px] font-medium text-[var(--text-secondary)] transition-colors hover:bg-[var(--bg-surface)] hover:text-[var(--text-primary)]"
                  >
                    {allGroupsExpanded ? (
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
            </div>
          )}
        </div>

        {/* Categorized navigation list */}
        <nav
          className={cn(
            "flex-1 overflow-y-auto py-3",
            compact ? "space-y-3 px-2.5" : "space-y-2.5 px-3",
          )}
          aria-label="Control plane navigation"
        >
          {visibleGroups.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-[var(--border-strong)] bg-[var(--bg-canvas)]/60 px-4 py-8 text-center">
              <SearchX className="h-5 w-5 text-[var(--text-muted)]" aria-hidden="true" />
              <p className="mt-2.5 text-xs font-semibold text-[var(--text-primary)]">
                {copy.emptySearchTitle}
              </p>
              <p className="mt-1 text-[11px] leading-relaxed text-[var(--text-muted)]">
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
            visibleGroups.map((group, groupIndex) => {
              const GroupIcon = group.icon;
              const isExpanded = compact || isSearching || expandedGroups[group.id];
              const groupHasActive = group.items.some((item) =>
                isNavItemActive(normalizedPathname, item.href),
              );

              return (
                <div
                  key={group.id}
                  className={cn(
                    compact && groupIndex > 0 && "border-t border-[var(--border-subtle)] pt-3",
                  )}
                >
                  {!compact ? (
                    <button
                      type="button"
                      onClick={() => toggleGroup(group.id)}
                      aria-expanded={isExpanded}
                      className={cn(
                        "flex h-9 w-full items-center justify-between gap-2 rounded-lg border px-2.5 text-left transition-all",
                        groupHasActive
                          ? "border-[var(--accent-primary)]/40 bg-[color-mix(in_srgb,var(--accent-primary)_8%,var(--bg-surface))] text-[var(--text-primary)] shadow-2xs"
                          : isExpanded
                            ? "border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-primary)] shadow-2xs hover:border-[var(--accent-primary)]/35"
                            : "border-[var(--border-strong)]/75 bg-[var(--bg-surface)]/80 text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:bg-[var(--bg-surface)] hover:text-[var(--text-primary)]",
                      )}
                    >
                      <span className="flex min-w-0 items-center gap-2">
                        <GroupIcon
                          className={cn(
                            "h-3.5 w-3.5 shrink-0",
                            groupHasActive
                              ? "text-[var(--accent-primary)]"
                              : "text-[var(--text-muted)]",
                          )}
                          aria-hidden="true"
                        />
                        <span className="truncate text-xs font-semibold tracking-tight">
                          {group.labels[locale]}
                        </span>
                      </span>
                      <span className="flex shrink-0 items-center gap-1.5">
                        <span
                          className={cn(
                            "rounded-md px-1.5 py-0.5 font-mono text-[10px] font-medium tabular-nums",
                            groupHasActive
                              ? "bg-[color-mix(in_srgb,var(--accent-primary)_14%,transparent)] text-[var(--accent-text)]"
                              : "bg-[var(--bg-surface-raised)] text-[var(--text-secondary)]",
                          )}
                        >
                          {group.items.length}
                        </span>
                        <ChevronDown
                          className={cn(
                            "h-3.5 w-3.5 shrink-0 transition-transform duration-200 ease-out",
                            groupHasActive
                              ? "text-[var(--accent-primary)]"
                              : "text-[var(--text-secondary)]",
                            !isExpanded && "-rotate-90",
                          )}
                          aria-hidden="true"
                        />
                      </span>
                    </button>
                  ) : null}

                  <div
                    className={cn(
                      "grid transition-[grid-template-rows,opacity] duration-200 ease-out",
                      isExpanded ? "grid-rows-[1fr] opacity-100" : "grid-rows-[0fr] opacity-0",
                    )}
                  >
                    <div className="overflow-hidden">
                      <div
                        className={cn(
                          "space-y-1",
                          !compact &&
                            "ml-3 mt-1.5 border-l-2 border-[var(--border-strong)]/80 pl-2.5 pb-0.5",
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
                              title={compact ? `${group.labels[locale]} · ${itemLabel}` : undefined}
                              aria-current={isActive ? "page" : undefined}
                              className={cn(
                                "group relative flex h-8.5 items-center rounded-lg border text-[13px] transition-all",
                                compact ? "h-9 justify-center px-0" : "gap-2 px-2",
                                isActive
                                  ? "border-[var(--accent-primary)]/45 bg-[var(--bg-surface)] font-semibold text-[var(--text-primary)] shadow-xs"
                                  : "border-transparent font-medium text-[var(--text-secondary)] hover:border-[var(--border-strong)]/70 hover:bg-[var(--bg-surface)]/85 hover:text-[var(--text-primary)]",
                              )}
                            >
                              {isActive && !compact ? (
                                <span
                                  aria-hidden="true"
                                  className="absolute -left-[12px] top-1.5 bottom-1.5 w-0.5 rounded-full bg-[var(--accent-primary)]"
                                />
                              ) : null}
                              <span
                                className={cn(
                                  "flex h-6 w-6 shrink-0 items-center justify-center rounded-md transition-colors",
                                  isActive
                                    ? "bg-[color-mix(in_srgb,var(--accent-primary)_14%,transparent)] text-[var(--accent-primary)]"
                                    : "text-[var(--text-muted)] group-hover:text-[var(--text-primary)]",
                                )}
                              >
                                <Icon className="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                              </span>
                              {!compact ? <span className="truncate">{itemLabel}</span> : null}
                            </Link>
                          );
                        })}
                      </div>
                    </div>
                  </div>
                </div>
              );
            })
          )}
        </nav>

        {/* Live gateway status footer */}
        <div
          className={cn(
            "shrink-0 border-t border-[var(--border-strong)]/80 bg-[var(--bg-surface)]/75 py-3",
            compact ? "px-2" : "px-4",
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
            <div className="flex items-center justify-between gap-2 text-xs">
              <span className="inline-flex items-center gap-2 text-[var(--text-secondary)]">
                <span className="relative flex h-2 w-2" aria-hidden="true">
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
                      "relative inline-flex h-2 w-2 rounded-full",
                      gatewayStatus.dotClass,
                    )}
                  />
                </span>
                <span className="font-medium text-[var(--text-primary)]">
                  {gatewayStatus.label}
                </span>
              </span>
              <span className="rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-1.5 py-0.5 font-mono text-[10px] tabular-nums text-[var(--text-muted)]">
                {versionTag}
              </span>
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="flex min-h-screen bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden shrink-0 border-r border-[var(--border-strong)] shadow-[1px_0_0_0_var(--border-subtle)] transition-[width] duration-200 lg:block",
          sidebarCollapsed ? "w-[72px]" : "w-[272px]",
        )}
      >
        <div className="sticky top-0 h-screen">{renderSidebar(sidebarCollapsed, "desktop")}</div>
      </aside>

      {/* Mobile navigation drawer with translucent scrim */}
      {mobileNavOpen ? (
        <div
          className="fixed inset-0 z-40 flex lg:hidden"
          role="dialog"
          aria-modal="true"
          aria-label="Navigation"
        >
          <div className="w-[280px] max-w-[85vw] border-r border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
            {renderSidebar(false, "mobile")}
          </div>
          <button
            type="button"
            className="flex-1 bg-black/45 backdrop-blur-[2px]"
            aria-label={copy.dismissMenuLabel}
            onClick={closeMobileNav}
          />
        </div>
      ) : null}

      {/* Main content area */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-20 bg-[var(--bg-surface)]/90 backdrop-blur-md">
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
              <Button
                type="button"
                variant="secondary"
                size="icon"
                onClick={toggleSidebarCollapsed}
                className="hidden h-9 w-9 shrink-0 lg:inline-flex"
                aria-label={sidebarCollapsed ? copy.expandSidebarLabel : copy.collapseSidebarLabel}
              >
                {sidebarCollapsed ? (
                  <PanelLeftOpen className="h-4 w-4" aria-hidden="true" />
                ) : (
                  <PanelLeftClose className="h-4 w-4" aria-hidden="true" />
                )}
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
                className="hidden h-9 items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 text-xs font-medium text-[var(--text-muted)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)] md:inline-flex"
                aria-label={`${copy.quickSearchLabel} (Ctrl+K)`}
              >
                <Search className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{copy.quickSearchLabel}</span>
                <kbd className="rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--text-muted)]">
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
