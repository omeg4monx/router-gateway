"use client";

import { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/shared/utils/cn";
import { APP_CONFIG } from "@/shared/constants/config";
import HeaderMenu from "@/shared/components/HeaderMenu";
import NotificationBell from "@/shared/components/NotificationBell";
import LogoMark from "@/shared/components/LogoMark";
import { useNavigation } from "@/shared/hooks/useNavigation";
import {
  MEDIA_ITEMS,
  COMBINED_WEB_ITEM,
  CORE_ROUTING_HREFS,
  SETTINGS_HREF,
  RAIL_PINNED_COOKIE,
} from "@/shared/constants/navigation";

/* Tactile Omegamon Studio Rail. A floating, rounded dock on the left edge:
   compact (60px, icon-only) or pinned/expanded (228px). Pinned, groups expand
   inline via .accordion-grid and stay open while a child route is active;
   compact, groups open as a floating flyout to the right of the rail. */

function RailRow({
  href,
  icon,
  label,
  active,
  open,
  onClick,
  onHover,
  onLeave,
  badge = null,
}) {
  const inner = (
    <>
      <span
        className={cn(
          "material-symbols-outlined rail-icon shrink-0 text-[18px] leading-none",
          active ? "text-primary fill-1" : "text-text-main/70"
        )}
      >
        {icon}
      </span>
      {open && (
        <span
          className={cn(
            "flex-1 truncate whitespace-nowrap text-[12.5px] leading-none",
            active ? "font-semibold text-primary" : "text-text-main/80"
          )}
        >
          {label}
        </span>
      )}
      {open && badge}
      {active && (
        <span
          className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-primary shadow-[0_0_10px_color-mix(in_srgb,var(--color-primary)_65%,transparent)]"
          aria-hidden="true"
        />
      )}
    </>
  );

  const cls = cn(
    "relative flex h-8.5 w-full items-center gap-2.5 rounded-[10px] border transition-[background-color,border-color,color,transform] duration-150 cursor-pointer active:scale-[0.98]",
    open ? "px-2.5" : "justify-center px-0",
    active
      ? "bg-primary/12 border-primary/20 text-primary font-medium"
      : "border-transparent text-text-main/80 hover:bg-surface-2/80 hover:text-text-main"
  );

  const sharedEvents = {
    onMouseEnter: !open ? (e) => onHover?.(label, e.currentTarget, false) : undefined,
    onMouseLeave: !open ? onLeave : undefined,
    onFocus: !open ? (e) => onHover?.(label, e.currentTarget, false) : undefined,
    onBlur: !open ? onLeave : undefined,
    onClick: (e) => {
      onLeave?.();
      onClick?.(e);
    },
    title: open ? label : undefined,
    "aria-label": label,
  };

  if (href) {
    return (
      <Link
        href={href}
        aria-current={active ? "page" : undefined}
        className={cls}
        {...sharedEvents}
      >
        {inner}
      </Link>
    );
  }
  return (
    <button type="button" className={cls} {...sharedEvents}>
      {inner}
    </button>
  );
}

RailRow.propTypes = {
  href: PropTypes.string,
  icon: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  active: PropTypes.bool,
  open: PropTypes.bool,
  onClick: PropTypes.func,
  onHover: PropTypes.func,
  onLeave: PropTypes.func,
  badge: PropTypes.node,
};

function GroupRow({ icon, label, active, open, expanded, onClick, onHover, onLeave }) {
  return (
    <button
      type="button"
      onClick={(e) => {
        onLeave?.();
        onClick?.(e);
      }}
      onMouseEnter={!open ? (e) => onHover?.(label, e.currentTarget, true) : undefined}
      onMouseLeave={!open ? onLeave : undefined}
      onFocus={!open ? (e) => onHover?.(label, e.currentTarget, true) : undefined}
      onBlur={!open ? onLeave : undefined}
      title={open ? label : undefined}
      aria-label={label}
      aria-expanded={expanded}
      className={cn(
        "relative flex h-8.5 w-full items-center gap-2.5 rounded-[10px] border transition-[background-color,border-color,color,transform] duration-150 cursor-pointer active:scale-[0.98]",
        open ? "px-2.5" : "justify-center px-0",
        active
          ? "bg-primary/12 border-primary/20 text-primary font-medium"
          : expanded
            ? "bg-surface-2/80 border-border-subtle text-text-main font-medium"
            : "border-transparent text-text-main/80 hover:bg-surface-2/80 hover:text-text-main"
      )}
    >
      <span
        className={cn(
          "material-symbols-outlined rail-icon shrink-0 text-[18px] leading-none",
          active ? "text-primary fill-1" : expanded ? "text-primary" : "text-text-main/70"
        )}
      >
        {icon}
      </span>
      {open && (
        <>
          <span
            className={cn(
              "flex-1 truncate whitespace-nowrap text-left text-[12.5px] leading-none",
              active ? "font-semibold text-primary" : expanded && "font-semibold text-text-main"
            )}
          >
            {label}
          </span>
          <span
            className="material-symbols-outlined text-[15px] leading-none text-text-subtle transition-transform duration-200 ease-[var(--ease-spring)] motion-reduce:transition-none"
            style={{ transform: expanded ? "rotate(90deg)" : "rotate(0deg)" }}
          >
            chevron_right
          </span>
        </>
      )}
      {active && (
        <span
          className="absolute left-0 top-1/2 h-4 w-[3px] -translate-y-1/2 rounded-r-full bg-primary shadow-[0_0_10px_color-mix(in_srgb,var(--color-primary)_65%,transparent)]"
          aria-hidden="true"
        />
      )}
    </button>
  );
}

GroupRow.propTypes = {
  icon: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  active: PropTypes.bool,
  open: PropTypes.bool,
  expanded: PropTypes.bool,
  onClick: PropTypes.func,
  onHover: PropTypes.func,
  onLeave: PropTypes.func,
};

/* Fixed, not absolute: the nav list scrolls, and its overflow would clip an
   absolutely positioned child. `top` is the opening group row's viewport offset. */
function Flyout({ top, label, children }) {
  return (
    <div
      role="group"
      aria-label={label}
      className="flyout-in custom-scrollbar fixed left-[78px] z-50 w-[240px] overflow-y-auto rounded-xl shell-popover p-1.5"
      style={{ top, maxHeight: `calc(100dvh - ${top}px - 12px)` }}
    >
      <div className="px-2.5 pb-1 pt-0.5 font-mono text-[10px] font-semibold uppercase tracking-wider text-text-muted">
        {label}
      </div>
      <div className="space-y-0.5">{children}</div>
    </div>
  );
}

Flyout.propTypes = { top: PropTypes.number.isRequired, label: PropTypes.string.isRequired, children: PropTypes.node };

function SectionLabel({ children }) {
  return (
    <div className="px-2.5 pb-1 font-mono text-[10px] font-semibold uppercase tracking-wider text-text-subtle">
      {children}
    </div>
  );
}

SectionLabel.propTypes = { children: PropTypes.node };

function NewBadge() {
  return (
    <span className="rounded-[4px] bg-primary/15 px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wide text-primary">
      New
    </span>
  );
}

export function FlyoutLink({ href, icon, label, active, onClick, badge = null }) {
  return (
    <Link
      href={href}
      onClick={onClick}
      title={label}
      aria-current={active ? "page" : undefined}
      className={cn(
        "relative flex items-center gap-2.5 rounded-[8px] border px-2.5 py-2 transition-[background-color,border-color,color,transform] duration-150 active:scale-[0.98]",
        active
          ? "bg-primary/12 border-primary/20 text-primary font-medium"
          : "border-transparent text-text-main/80 hover:bg-surface-2/80 hover:text-text-main"
      )}
    >
      <span
        className={cn(
          "material-symbols-outlined rail-icon shrink-0 text-[16px] leading-none",
          active ? "text-primary fill-1" : "text-text-main/70"
        )}
      >
        {icon}
      </span>
      <span className={cn("flex-1 truncate text-[12px] leading-none", active && "font-semibold text-primary")}>
        {label}
      </span>
      {badge}
    </Link>
  );
}

FlyoutLink.propTypes = {
  href: PropTypes.string.isRequired,
  icon: PropTypes.string,
  label: PropTypes.string.isRequired,
  active: PropTypes.bool,
  onClick: PropTypes.func,
  badge: PropTypes.node,
};

export default function LeftRail({ initialPinned = true }) {
  const pathname = usePathname();
  const [pinned, setPinned] = useState(initialPinned);
  const [expandedGroup, setExpandedGroup] = useState(null);
  const [flyoutTop, setFlyoutTop] = useState(0);
  const [tooltip, setTooltip] = useState(null);
  const [displayPort, setDisplayPort] = useState("8080");
  const railRef = useRef(null);

  const {
    navItems,
    workshopItems,
    debugItems,
    systemItems,
    canOpenMedia,
    isApiKeyUser,
    hasMediaNew,
    isActive,
  } = useNavigation();

  const open = pinned;
  const coreRouting = navItems.filter((i) => CORE_ROUTING_HREFS.includes(i.href));
  const operations = navItems.filter((i) => !CORE_ROUTING_HREFS.includes(i.href));
  const settingsItem = systemItems.find((i) => i.href === SETTINGS_HREF);
  const groupSystemItems = systemItems.filter((i) => i.href !== SETTINGS_HREF);
  const mediaItems = canOpenMedia ? MEDIA_ITEMS : [];
  const showSystemGroup = mediaItems.length > 0 || debugItems.length > 0 || groupSystemItems.length > 0;
  const workshopActive = workshopItems.some((i) => isActive(i.href));
  const systemActive = [...mediaItems, ...debugItems, ...groupSystemItems].some((i) => isActive(i.href));
  const activeGroup = workshopActive ? "workshop" : systemActive ? "system" : null;

  useEffect(() => {
    // Derive live port from window location or health endpoint
    if (window.location.port) {
      setDisplayPort(window.location.port);
      return;
    }
    // Reverse proxy without port (80/443), query server health probe
    fetch("/api/health", { cache: "no-store" })
      .then((r) => r.json())
      .then((d) => {
        if (d?.port) setDisplayPort(String(d.port));
      })
      .catch(() => {});
  }, []);

  // Pinned: keep the group holding the current page open so its active row stays
  // visible. Compact: navigating (or collapsing the rail) closes the flyout.
  useEffect(() => {
    setExpandedGroup(open ? activeGroup : null);
  }, [pathname, open, activeGroup]);

  // Dismiss any compact floating tooltip whenever the rail expands/collapses or route changes.
  useEffect(() => {
    setTooltip(null);
  }, [open, pathname]);

  // Dismiss compact tooltip on Escape (WCAG 1.4.13).
  useEffect(() => {
    if (!tooltip) return undefined;
    const onKeyDown = (e) => {
      if (e.key === "Escape") setTooltip(null);
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [tooltip]);

  // Only the compact flyout is transient; an inline group stays until toggled.
  useEffect(() => {
    if (open || !expandedGroup) return undefined;
    const close = () => setExpandedGroup(null);
    const onDocClick = (e) => {
      if (railRef.current && !railRef.current.contains(e.target)) close();
    };
    const onKeyDown = (e) => e.key === "Escape" && close();
    document.addEventListener("click", onDocClick);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("click", onDocClick);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open, expandedGroup]);

  const showTooltip = (label, el, isGroup = false) => {
    if (open || !el) return;
    const rect = el.getBoundingClientRect();
    const top = Math.max(16, Math.min(rect.top + rect.height / 2, window.innerHeight - 24));
    setTooltip({ label, top, isGroup });
  };

  const clearTooltip = () => setTooltip(null);

  const handleTogglePin = () => {
    setTooltip(null);
    const next = !pinned;
    setPinned(next);
    document.cookie = `${RAIL_PINNED_COOKIE}=${next}; path=/; max-age=31536000; samesite=lax`;
  };

  const toggleGroup = (group) => (e) => {
    setTooltip(null);
    setFlyoutTop(e.currentTarget.getBoundingClientRect().top);
    setExpandedGroup((prev) => (prev === group ? null : group));
  };

  const closeFlyout = () => {
    setTooltip(null);
    if (!open) setExpandedGroup(null);
  };

  const handleLogout = async () => {
    try {
      const res = await fetch("/api/auth/logout", { method: "POST" });
      if (res.ok) window.location.assign("/login");
    } catch (err) {
      console.error("Failed to logout:", err);
    }
  };

  const renderLinks = (items) =>
    items.map((item) => (
      <FlyoutLink
        key={item.href}
        href={item.href}
        icon={item.icon}
        label={item.label}
        active={isActive(item.href)}
        onClick={closeFlyout}
        badge={item.href === COMBINED_WEB_ITEM.href && hasMediaNew ? <NewBadge /> : null}
      />
    ));

  const renderRows = (items) =>
    items.map((item) => (
      <RailRow
        key={item.href}
        href={item.href}
        icon={item.icon}
        label={item.label}
        active={isActive(item.href)}
        open={open}
        onClick={closeFlyout}
        onHover={showTooltip}
        onLeave={clearTooltip}
      />
    ));

  const renderGroup = (id, icon, label, active, links) => {
    const expanded = expandedGroup === id;
    return (
      <div>
        <GroupRow
          icon={icon}
          label={label}
          active={active}
          open={open}
          expanded={expanded}
          onClick={toggleGroup(id)}
          onHover={showTooltip}
          onLeave={clearTooltip}
        />
        {open ? (
          <div
            className="accordion-grid"
            data-open={expanded ? "true" : "false"}
            aria-hidden={!expanded}
            inert={!expanded || undefined}
          >
            <div className="accordion-inner">
              <div className="ml-3.5 mt-1 space-y-0.5 border-l border-border-subtle pl-2 pb-0.5">
                {links}
              </div>
            </div>
          </div>
        ) : (
          expanded && (
            <Flyout top={flyoutTop} label={label}>
              {links}
            </Flyout>
          )
        )}
      </div>
    );
  };

  return (
    <aside
      ref={railRef}
      className={cn(
        "relative isolate hidden my-2.5 ml-2.5 h-[calc(100dvh-20px)] flex-col rounded-2xl shell-rail overflow-visible select-none lg:flex z-[45] shrink-0 transition-[width] duration-250 ease-[var(--ease-spring)] motion-reduce:transition-none",
        open ? "w-[228px]" : "w-[60px]"
      )}
      aria-label="Application navigation"
    >
      {/* Brand header */}
      <div
        className={cn(
          "flex h-13 items-center border-b border-border-subtle/80 px-2.5 shrink-0 gap-2",
          open ? "justify-between" : "justify-center"
        )}
      >
        {open ? (
          <>
            <Link
              href="/dashboard/endpoint"
              className="flex items-center gap-2.5 overflow-hidden transition-opacity hover:opacity-90 min-w-0 flex-1"
              title="OMEG4MON ROUTER Dashboard"
            >
              <div className="size-8 rounded-xl border border-primary/25 bg-primary/10 flex items-center justify-center shrink-0 p-1 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.25)]">
                <LogoMark className="size-full" />
              </div>
              <div className="min-w-0 flex flex-col flex-1">
                <span className="text-[11.5px] font-bold tracking-wider uppercase text-text-main truncate">
                  {APP_CONFIG.name}
                </span>
                <span className="text-[9.5px] font-mono text-text-muted truncate">
                  v{APP_CONFIG.version}
                </span>
              </div>
            </Link>

            <button
              type="button"
              onClick={handleTogglePin}
              className="size-7.5 rounded-[9px] text-text-muted hover:text-text-main hover:bg-surface-2/80 transition-colors cursor-pointer flex items-center justify-center shrink-0 border border-transparent hover:border-border-subtle active:scale-[0.96]"
              title="Collapse sidebar"
              aria-label="Collapse sidebar"
            >
              <span className="material-symbols-outlined text-[17px] leading-none text-primary">
                left_panel_close
              </span>
            </button>
          </>
        ) : (
          <button
            type="button"
            onClick={handleTogglePin}
            onMouseEnter={(e) => showTooltip("Expand sidebar", e.currentTarget, false)}
            onMouseLeave={clearTooltip}
            onFocus={(e) => showTooltip("Expand sidebar", e.currentTarget, false)}
            onBlur={clearTooltip}
            className="size-9 rounded-[10px] text-text-muted hover:text-text-main hover:bg-surface-2/80 transition-colors cursor-pointer flex items-center justify-center shrink-0 border border-border-subtle/80 hover:border-primary/40 active:scale-[0.96]"
            aria-label="Expand sidebar"
          >
            <span className="material-symbols-outlined text-[18px] leading-none text-primary">
              left_panel_open
            </span>
          </button>
        )}
      </div>

      {/* Main navigation list */}
      <div
        onScroll={clearTooltip}
        className="flex-1 overflow-y-auto overflow-x-hidden px-2 py-2.5 custom-scrollbar space-y-3"
      >
        <div className="space-y-0.5">
          {open && <SectionLabel>Routing</SectionLabel>}
          {renderRows(coreRouting)}
        </div>

        <div className="space-y-0.5 pt-2 border-t border-border-subtle/80">
          {open && <SectionLabel>Operations</SectionLabel>}
          {renderRows(operations)}
          {!isApiKeyUser &&
            renderRows([{ href: "/dashboard/vouchers", icon: "card_giftcard", label: "Vouchers & Market" }])}
        </div>

        {(workshopItems.length > 0 || showSystemGroup) && (
          <div className="space-y-0.5 pt-2 border-t border-border-subtle/80">
            {open && <SectionLabel>Extensions</SectionLabel>}
            {workshopItems.length > 0 &&
              renderGroup("workshop", "auto_awesome", "Features", workshopActive, renderLinks(workshopItems))}
            {showSystemGroup &&
              renderGroup(
                "system",
                "tune",
                "System",
                systemActive,
                <>
                  {mediaItems.length > 0 && (
                    <div className="mb-1 border-b border-border-subtle/80 pb-1">
                      <div className="px-2.5 py-0.5 font-mono text-[10px] font-medium uppercase tracking-wider text-text-subtle">
                        Media Engines
                      </div>
                      {renderLinks(mediaItems)}
                    </div>
                  )}
                  {renderLinks(debugItems)}
                  {renderLinks(groupSystemItems)}
                </>
              )}
          </div>
        )}
      </div>

      {settingsItem && (
        <div className="shrink-0 border-t border-border-subtle/80 px-2 py-1.5">
          {renderRows([settingsItem])}
        </div>
      )}

      {/* Footer controls & user menu */}
      <div
        className={cn(
          "border-t border-border-subtle/80 p-2 bg-surface-2/45 rounded-b-2xl flex shrink-0 gap-1.5",
          open ? "items-center justify-between" : "flex-col-reverse items-center"
        )}
      >
        <div className={cn("flex items-center gap-2 min-w-0", !open && "justify-center")}>
          {!isApiKeyUser ? (
            <HeaderMenu onLogout={handleLogout} />
          ) : (
            <div
              className="flex items-center justify-center size-8 rounded-[10px] text-primary"
              title="API Session Active"
            >
              <span className="material-symbols-outlined text-[20px]">key</span>
            </div>
          )}
          {open && (
            <div className="min-w-0 flex flex-col">
              <span className="text-[11.5px] font-semibold text-text-main truncate">
                {isApiKeyUser ? "API Session" : "Instance Admin"}
              </span>
              <span className="text-[9.5px] text-success flex items-center gap-1 font-mono whitespace-nowrap">
                <span className="size-1.5 rounded-full bg-success shrink-0" />
                Live {displayPort}
              </span>
            </div>
          )}
        </div>

        {(!isApiKeyUser || open) && (
          <div className={cn("flex items-center gap-1 shrink-0", !open && "flex-col")}>
            <NotificationBell placement="right-bottom" />
            {open && isApiKeyUser && (
              <button
                type="button"
                onClick={handleLogout}
                className="flex items-center justify-center size-8 rounded-[10px] text-primary hover:bg-primary/10 transition-colors cursor-pointer active:scale-[0.97]"
                title="Log Out"
                aria-label="Log Out"
              >
                <span className="material-symbols-outlined text-[18px]">logout</span>
              </button>
            )}
          </div>
        )}
      </div>

      {!open && !expandedGroup && tooltip && (
        <div
          role="tooltip"
          className="fixed left-[78px] -translate-y-1/2 z-[60] pointer-events-none shell-popover popover-in rounded-lg px-2.5 py-1.5 text-[12px] font-medium text-text-main whitespace-nowrap flex items-center gap-1.5"
          style={{ top: tooltip.top }}
        >
          <span>{tooltip.label}</span>
          {tooltip.isGroup && (
            <span className="material-symbols-outlined text-[13px] leading-none text-text-subtle">
              chevron_right
            </span>
          )}
        </div>
      )}
    </aside>
  );
}

LeftRail.propTypes = {
  initialPinned: PropTypes.bool,
};
