"use client";

import PropTypes from "prop-types";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/shared/utils/cn";
import ProviderIcon from "@/shared/components/ProviderIcon";
import { getPageInfo } from "@/shared/utils/pageInfo";
import { translate } from "@/i18n/runtime";
import NotificationBell from "@/shared/components/NotificationBell";

/* Full-bleed heading band. Every page opens with its own clean instrument band —
   crumb path / eyebrow, large title, description, and page-specific action slots or sub-tabs.
   Pages that draw their own hero pass `bare` to suppress it. */

export default function PageHeading({
  title,
  description,
  icon,
  breadcrumbs,
  actions,
  tabs,
  bare = false,
  showBell = true,
  className = "",
}) {
  const pathname = usePathname();
  const info = getPageInfo(pathname);

  const resolvedTitle = title ?? info.title;
  const resolvedDescription = description ?? info.description;
  const resolvedIcon = icon ?? info.icon;
  const resolvedCrumbs = breadcrumbs ?? info.breadcrumbs;

  if (bare) return null;

  const hasCrumbs = Array.isArray(resolvedCrumbs) && resolvedCrumbs.length > 0;
  if (!resolvedTitle && !hasCrumbs) return null;

  return (
    <div className={cn("relative z-20 border-b border-border/80 bg-surface/65 backdrop-blur-md", className)}>
      <div className="flex flex-wrap items-end justify-between gap-3 px-4 pb-3 pt-3.5 lg:px-7">
        <div className="min-w-0">
          {/* Eyebrow: ember tick + crumb path, or the section label. */}
          <div className="mb-1 flex min-w-0 items-center gap-2">
            <span className="ember-rule-left w-5 shrink-0" aria-hidden="true" />
            {hasCrumbs ? (
              <div className="flex min-w-0 items-center gap-1.5">
                {resolvedCrumbs.map((crumb, index) => (
                  <div key={`${crumb.label}-${crumb.href || "current"}`} className="flex min-w-0 items-center gap-1.5">
                    {index > 0 && (
                      <span className="material-symbols-outlined shrink-0 text-[13px] leading-none text-text-subtle">
                        chevron_right
                      </span>
                    )}
                    {crumb.href ? (
                      <Link
                        href={crumb.href}
                        className="truncate rounded px-1 -mx-1 py-0.5 text-[11px] uppercase tracking-[0.06em] text-text-muted transition-colors hover:bg-surface-2/70 hover:text-primary"
                      >
                        {translate(crumb.label)}
                      </Link>
                    ) : (
                      <span className="flex min-w-0 items-center gap-1.5">
                        {crumb.image && (
                          <ProviderIcon
                            src={crumb.image}
                            alt={crumb.label}
                            size={16}
                            className="max-h-[16px] max-w-[16px] shrink-0 rounded object-contain"
                            fallbackText={crumb.label.slice(0, 2).toUpperCase()}
                          />
                        )}
                        <span className="micro-label truncate">{translate(crumb.label)}</span>
                      </span>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <span className="micro-label truncate">{translate(resolvedTitle)}</span>
            )}
          </div>

          <div className="flex min-w-0 items-center gap-3 py-0.5">
            {resolvedIcon && (
              <span className="size-9 rounded-xl bg-primary/10 border border-primary/20 text-primary flex items-center justify-center shadow-xs shrink-0">
                <span className="material-symbols-outlined text-[20px] leading-none">
                  {resolvedIcon}
                </span>
              </span>
            )}
            <h1 className="truncate text-[20px] font-semibold leading-normal tracking-tight text-text-main lg:text-[22px] py-0.5">
              {translate(resolvedTitle)}
            </h1>
          </div>

          {resolvedDescription ? (
            <p className="mt-1 max-w-3xl text-[12px] leading-snug text-text-muted">
              {translate(resolvedDescription)}
            </p>
          ) : null}
        </div>

        {/* NotificationBell in PageHeading actions or docked */}
        {(actions || showBell) ? (
          <div className="flex shrink-0 items-center gap-2">
            {actions ? <div className="flex shrink-0 flex-wrap items-center gap-2">{actions}</div> : null}
            {showBell ? <div className="hidden lg:block shrink-0"><NotificationBell /></div> : null}
          </div>
        ) : null}
      </div>

      {/* Optional Page-Specific Sub-tabs (Level-2 navigation) */}
      {tabs && (
        <div className="flex items-center gap-1 overflow-x-auto px-4 lg:px-7 pb-2 pt-1 custom-scrollbar border-t border-border-subtle">
          {tabs}
        </div>
      )}

      <span className="ember-rule pointer-events-none absolute inset-x-0 bottom-0 opacity-45" aria-hidden="true" />
    </div>
  );
}

PageHeading.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
  icon: PropTypes.string,
  breadcrumbs: PropTypes.array,
  actions: PropTypes.node,
  tabs: PropTypes.node,
  bare: PropTypes.bool,
  showBell: PropTypes.bool,
  className: PropTypes.string,
};
