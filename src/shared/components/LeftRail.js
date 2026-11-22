"use client";

import { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/shared/utils/cn";
import { APP_CONFIG } from "@/shared/constants/config";
import HeaderMenu from "@/shared/components/HeaderMenu";
import LogoMark from "@/shared/components/LogoMark";
import { useNavigation } from "@/shared/hooks/useNavigation";
import {
  MEDIA_ITEMS,
  COMBINED_WEB_ITEM,
  CORE_ROUTING_HREFS,
  SETTINGS_HREF,
  RAIL_PINNED_COOKIE,
} from "@/shared/constants/navigation";

/* Suspended Studio Rail. A floating, rounded dock on the left edge: compact
   (52px, icon-only) or pinned/expanded (212px). Pinned, a group opens inline and
   stays open while one of its pages is active; compact, it opens as a flyout to
   the right of the rail. The active item carries a glowing ember indicator. */

function RailRow({ href, icon, label, active, open, onClick, badge = null }) {
  const inner = (
    <>
      <span
        className={cn(
          "material-symbols-outlined shrink-0 text-[18px] leading-none transition-colors",
          active ? "text-primary font-semibold" : "text-text-main/70"
        )}
      >
        {icon}
      </span>
      {open && (
        <span className={cn("flex-1 truncate whitespace-nowrap text-[12.5px] leading-none", active ? "font-semibold text-primary" : "text-text-main/80")}>
          {label}
        </span>
      )}
      {open && badge}
      {active && (
        <span className="absolute left-0 top-1/2 h-4 w-[2.5px] -translate-y-1/2 rounded-r bg-primary shadow-[0_0_8px_rgba(var(--color-primary-rgb,229,106,74),0.6)]" aria-hidden="true" />
      )}
    </>
  );

  const cls = cn(
    "relative flex h-7.5 w-full items-center gap-2 rounded-[5px] transition-colors duration-150 cursor-pointer",
    open ? "px-2" : "justify-center px-0",
    active ? "bg-primary/15 text-primary font-medium border border-border-subtle" : "text-text-main/80 hover:bg-surface-2 hover:text-text-main"
  );

  if (href) {
    return (
      <Link href={href} onClick={onClick} title={label} className={cls}>
        {inner}
      </Link>
    );
  }
  return (
    <button type="button" onClick={onClick} title={label} className={cls}>
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
  badge: PropTypes.node,
};

function GroupRow({ icon, label, active, open, expanded, onClick }) {
  return (
    <button
      type="button"
      onClick={onClick}
      title={label}
      aria-expanded={expanded}
      className={cn(
        "relative flex h-7.5 w-full items-center gap-2 rounded-[5px] transition-colors duration-150 cursor-pointer",
        open ? "px-2" : "justify-center px-0",
        active || expanded ? "bg-primary/15 text-primary border border-border-subtle font-medium" : "text-text-main/80 hover:bg-surface-2 hover:text-text-main"
      )}
    >
      <span
        className={cn(
          "material-symbols-outlined shrink-0 text-[18px] leading-none transition-colors",
          active || expanded ? "text-primary" : "text-text-main/70"
        )}
      >
        {icon}
      </span>
      {open && (
        <>
          <span className={cn("flex-1 truncate whitespace-nowrap text-left text-[12.5px] leading-none", (active || expanded) && "font-semibold text-primary")}>
            {label}
          </span>
          <span
            className="material-symbols-outlined text-[14px] leading-none text-text-subtle transition-transform duration-200"
            style={{ transform: expanded ? "rotate(90deg)" : "rotate(0deg)" }}
          >
            chevron_right
          </span>
        </>
      )}
      {active && (
        <span className="absolute left-0 top-1/2 h-4 w-[2.5px] -translate-y-1/2 rounded-r bg-primary shadow-[0_0_8px_rgba(var(--color-primary-rgb,229,106,74),0.6)]" aria-hidden="true" />
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
};

/* Fixed, not absolute: the nav list scrolls, and its overflow would clip an
   absolutely positioned child. `top` is the opening group row's viewport offset. */
function Flyout({ top, label, children }) {
  return (
    <div
      role="group"
      aria-label={label}
      className="fade-in custom-scrollbar fixed left-[68px] z-50 w-[240px] overflow-y-auto rounded-[6px] border border-border bg-surface p-1.5 shadow-[var(--shadow-elev)]"
      style={{ top, maxHeight: `calc(100dvh - ${top}px - 8px)` }}
    >
      <div className="px-2.5 pb-1 pt-0.5 font-mono text-[10.5px] uppercase tracking-wider text-text-muted">{label}</div>
      {children}
    </div>
  );
}

Flyout.propTypes = { top: PropTypes.number.isRequired, label: PropTypes.string.isRequired, children: PropTypes.node };

function SectionLabel({ children }) {
  return <div className="px-2 pb-1 font-mono text-[10.5px] uppercase tracking-wider text-text-muted">{children}</div>;
}

SectionLabel.propTypes = { children: PropTypes.node };

function NewBadge() {
  return (
    <span className="rounded-[3px] bg-brand-500/15 px-1.5 py-0.5 text-[9px] font-semibold uppercase tracking-wide text-primary">
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
      className={cn(
        "flex items-center gap-2 rounded-[4px] px-2.5 py-[6px] transition-colors",
        active ? "bg-primary/15 text-primary font-medium" : "text-text-main/80 hover:bg-surface-2 hover:text-text-main"
      )}
    >
      <span className={cn("material-symbols-outlined shrink-0 text-[15px] leading-none transition-colors", active ? "text-primary" : "text-text-main/70")}>
        {icon}
      </span>
      <span className={cn("flex-1 truncate text-[12px] leading-none", active && "font-semibold text-primary")}>{label}</span>
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

  // Only the compact flyout is transient; an inline group stays until toggled.
  useEffect(() => {
    if (open || !expandedGroup) return;
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

  const handleTogglePin = () => {
    const next = !pinned;
    setPinned(next);
    document.cookie = `${RAIL_PINNED_COOKIE}=${next}; path=/; max-age=31536000; samesite=lax`;
  };

  const toggleGroup = (group) => (e) => {
    setFlyoutTop(e.currentTarget.getBoundingClientRect().top);
    setExpandedGroup((prev) => (prev === group ? null : group));
  };

  const closeFlyout = () => {
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
      />
    ));

  const renderGroup = (id, icon, label, active, links) => {
    const expanded = expandedGroup === id;
    return (
      <div>
        <GroupRow icon={icon} label={label} active={active} open={open} expanded={expanded} onClick={toggleGroup(id)} />
        {expanded &&
          (open ? (
            <div className="ml-3 mt-1 space-y-0.5 border-l-2 border-border-subtle pl-2">{links}</div>
          ) : (
            <Flyout top={flyoutTop} label={label}>
              {links}
            </Flyout>
          ))}
      </div>
    );
  };

  return (
    <aside
      ref={railRef}
      className={cn(
        "relative isolate hidden my-2 ml-2 h-[calc(100dvh-16px)] flex-col rounded-[8px] border border-border bg-surface select-none lg:flex z-30 shrink-0 shadow-lg",
        open ? "w-[212px] overflow-hidden" : "w-[52px]"
      )}
      aria-label="Application navigation"
    >
      {/* Brand header */}
      <div className={cn(
        "flex h-12 items-center border-b border-border px-2 shrink-0 bg-surface-2 gap-1.5",
        open ? "justify-between" : "justify-center relative"
      )}>
        {open ? (
          <>
            <Link
              href="/dashboard/endpoint"
              className="flex items-center gap-2 overflow-hidden transition-opacity hover:opacity-90 min-w-0 flex-1"
              title="OMEG4MON ROUTER Dashboard"
            >
              <div
                className="size-7 rounded-[5px] border bg-surface-2 flex items-center justify-center shrink-0 shadow-xs overflow-hidden p-0.5"
                style={{ borderColor: "var(--color-primary, #E56A4A)" }}
              >
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
              className="p-1.5 rounded-[5px] text-text-muted hover:text-text-main hover:bg-surface-3 transition-colors cursor-pointer shrink-0 border border-transparent hover:border-border"
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
            className="size-8 rounded-[6px] text-text-muted hover:text-text-main hover:bg-surface-3 transition-colors cursor-pointer flex items-center justify-center shrink-0 border border-border/60 hover:border-primary/50"
            title="Expand sidebar"
            aria-label="Expand sidebar"
          >
            <span className="material-symbols-outlined text-[18px] leading-none text-primary">
              left_panel_open
            </span>
          </button>
        )}
      </div>

      {/* Main navigation list */}
      <div className="flex-1 overflow-y-auto px-2 py-2.5 custom-scrollbar space-y-3">
        <div className="space-y-0.5">
          {open && <SectionLabel>Routing</SectionLabel>}
          {renderRows(coreRouting)}
        </div>

        <div className="space-y-0.5 pt-1 border-t border-border-subtle">
          {open && <SectionLabel>Operations</SectionLabel>}
          {renderRows(operations)}
          {!isApiKeyUser &&
            renderRows([{ href: "/dashboard/vouchers", icon: "card_giftcard", label: "Vouchers & Market" }])}
        </div>

        {(workshopItems.length > 0 || showSystemGroup) && (
          <div className="space-y-0.5 pt-1 border-t border-border-subtle">
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
                    <div className="mb-1 border-b border-border-subtle pb-1">
                      <div className="px-2.5 py-0.5 text-[10.5px] font-medium uppercase text-text-muted">Media Engines</div>
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

      {settingsItem && <div className="shrink-0 border-t border-border px-2 py-1.5">{renderRows([settingsItem])}</div>}

      {/* Footer controls & user menu */}
      <div className="border-t border-border p-2 bg-surface-2 flex items-center justify-between gap-2 shrink-0">
        <div className={cn("flex items-center gap-2 min-w-0", !open && "w-full justify-center")}>
          {!isApiKeyUser ? (
            <HeaderMenu onLogout={handleLogout} />
          ) : (
            <div className="flex items-center justify-center p-2 rounded-[6px] text-primary" title="API Session Active">
              <span className="material-symbols-outlined text-[20px]">key</span>
            </div>
          )}
          {open && (
            <div className="min-w-0 flex flex-col">
              <span className="text-[11.5px] font-semibold text-text-main truncate">
                {isApiKeyUser ? "API Session" : "Instance Admin"}
              </span>
              <span className="text-[9.5px] text-success flex items-center gap-1 font-mono">
                <span className="size-1.5 rounded-full bg-success shrink-0" />
                Live {displayPort}
              </span>
            </div>
          )}
        </div>
        {open && isApiKeyUser && (
          <button
            onClick={handleLogout}
            className="flex items-center justify-center p-1.5 rounded-[6px] text-primary hover:bg-primary/10 transition-colors"
            title="Log Out"
          >
            <span className="material-symbols-outlined text-[18px]">logout</span>
          </button>
        )}
      </div>
    </aside>
  );
}

LeftRail.propTypes = {
  initialPinned: PropTypes.bool,
};
