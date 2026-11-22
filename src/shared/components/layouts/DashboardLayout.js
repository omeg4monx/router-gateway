"use client";

import { useEffect } from "react";
import { usePathname } from "next/navigation";
import { useNotificationStore } from "@/store/notificationStore";
import LeftRail from "../LeftRail";
import MobileNav from "../MobileNav";
import PageHeading from "../PageHeading";
import StatusBar from "../StatusBar";
import WelcomeModal from "../WelcomeModal";
import UpdateBanner from "../UpdateBanner";
import TaskDock from "../TaskDock";

/* Toast variants pull from the semantic token triplet, so the same toast is
   legible in Omegamon light and Zwart dark without a dark: override per variant. */
function getToastStyle(type) {
  if (type === "success") {
    return {
      accent: "bg-success",
      iconWrap: "text-success",
      border: "border-success/35",
      icon: "check_circle",
    };
  }
  if (type === "error") {
    return {
      accent: "bg-danger",
      iconWrap: "text-danger",
      border: "border-danger/35",
      icon: "error",
    };
  }
  if (type === "warning") {
    return {
      accent: "bg-warning",
      iconWrap: "text-warning",
      border: "border-warning/35",
      icon: "warning",
    };
  }
  return {
    accent: "bg-primary",
    iconWrap: "text-primary",
    border: "border-border",
    icon: "info",
  };
}

/* Pages that render their own hero / heading block, so the layout heading band
   would duplicate them. Keyed by the exact pathname or a prefix. */
const PAGES_WITH_OWN_HEADING = [
  "/dashboard",
  "/dashboard/endpoint",
  "/dashboard/basic-chat",
  "/dashboard/providers/new",
  "/dashboard/plugins",
  "/dashboard/api-key-usage",
  "/dashboard/translator",
  "/dashboard/model-bypass",
  "/dashboard/arena",
  "/dashboard/model-editor",
  "/dashboard/prompt-injection",
  "/dashboard/combos",
  "/dashboard/prd-builder",
  "/dashboard/connection",
];

function hasOwnHeading(pathname) {
  if (!pathname) return false;
  if (pathname === "/dashboard" || pathname === "/dashboard/endpoint") return true;
  if (PAGES_WITH_OWN_HEADING.some((p) => p !== "/dashboard" && (pathname === p || pathname.startsWith(`${p}/`)))) return true;
  // Provider / media-provider detail pages draw their own hero with the logo.
  if (/^\/dashboard\/providers\/[^/]+$/.test(pathname)) return true;
  if (/^\/dashboard\/media-providers\/[^/]+\/[^/]+$/.test(pathname)) return true;
  if (/^\/dashboard\/cli-tools\/[^/]+$/.test(pathname)) return true;
  return false;
}

export default function DashboardLayout({ children, railPinned = true }) {
  const pathname = usePathname();
  const notifications = useNotificationStore((state) => state.notifications);
  const removeNotification = useNotificationStore((state) => state.removeNotification);

  // Security: Inactivity / Idle Timeout (OWASP protection for public PCs / warnet)
  // Auto-logs out if no keyboard/mouse/touch interaction for 60 minutes.
  useEffect(() => {
    const IDLE_TIMEOUT_MS = 60 * 60 * 1000; // 60 minutes of true inactivity
    let lastActive = Date.now();

    const recordActivity = () => {
      lastActive = Date.now();
    };

    const events = ["mousemove", "keydown", "click", "touchstart", "scroll"];
    events.forEach((evt) => window.addEventListener(evt, recordActivity, { passive: true }));

    const idleTimer = setInterval(() => {
      if (Date.now() - lastActive > IDLE_TIMEOUT_MS) {
        // True inactivity reached: terminate session to protect warnet / shared devices
        fetch("/api/auth/logout", { method: "POST" })
          .finally(() => {
            window.location.href = "/login?reason=idle_timeout";
          });
      }
    }, 30000); // Check every 30 seconds

    return () => {
      events.forEach((evt) => window.removeEventListener(evt, recordActivity));
      clearInterval(idleTimer);
    };
  }, []);

  const isChat = pathname === "/dashboard/basic-chat";
  const showHeading = !hasOwnHeading(pathname);

  return (
    // The root owns the stacking context (relative + z-0). It must NOT be
    // `isolate` on the inner <main>: an `isolate` there opens a nested
    // stacking context, which traps every `fixed inset-0 z-50` overlay
    // declared inside the page below the `z-30` mobile app header — the
    // header then paints over the modal and hides its title row.
    <div className="relative z-0 flex h-dvh w-full overflow-hidden bg-bg">
      {/* Toasts - High z-index to stay above charts, cards, and page headings */}
      <div className="fixed top-4 right-4 z-[9999] flex w-[min(92vw,380px)] flex-col gap-2 pointer-events-none">
        {notifications.map((n) => {
          const style = getToastStyle(n.type);
          return (
            <div
              key={n.id}
              role="status"
              className={`slide-in-right pointer-events-auto relative overflow-hidden rounded-xl shell-popover px-3.5 py-2.5 text-text-main ${style.border}`}
            >
              <span
                className={`pointer-events-none absolute inset-y-2 left-0 w-[3px] rounded-r-full ${style.accent}`}
                aria-hidden="true"
              />
              <div className="flex items-start gap-2.5 pl-1">
                <span className={`material-symbols-outlined mt-0.5 shrink-0 text-[18px] leading-5 ${style.iconWrap}`}>
                  {style.icon}
                </span>
                <div className="min-w-0 flex-1">
                  {n.title ? <p className="mb-0.5 text-xs font-semibold text-text-main">{n.title}</p> : null}
                  <p className="whitespace-pre-wrap break-words text-xs text-text-muted">{n.message}</p>
                </div>
                {n.dismissible ? (
                  <button
                    type="button"
                    onClick={() => removeNotification(n.id)}
                    className="flex size-6 shrink-0 items-center justify-center rounded-md text-text-subtle transition-colors hover:bg-surface-2 hover:text-text-main cursor-pointer"
                    aria-label="Dismiss notification"
                  >
                    <span className="material-symbols-outlined text-[15px]">close</span>
                  </button>
                ) : null}
              </div>
            </div>
          );
        })}
      </div>
      <WelcomeModal />

      {/* Left-docked navigation (desktop) */}
      <LeftRail initialPinned={railPinned} />

      {/* Content column */}
      <div className="flex min-w-0 flex-1 flex-col lg:my-2.5 lg:mr-2.5 lg:ml-2 lg:rounded-2xl lg:border lg:border-border/75 lg:bg-surface/45 lg:shadow-[var(--shadow-soft)] lg:overflow-hidden">
        <TaskDock />

        {/* Docked mobile top bar */}
        <div className="relative z-50">
          <MobileNav />
        </div>

        <div className="relative z-40">
          <UpdateBanner />
        </div>

        {/* No `isolate` here on purpose. `isolate` opens a stacking context,
            which traps every `fixed inset-0 z-50` overlay declared inside the
            page below the `z-30` mobile app header — the header then paints
            over the modal and hides its title row. The veil is the only
            element that needs to sit behind the content, and it already uses
            `-z-10` against the root stacking context. */}
        <main className="relative flex min-h-0 flex-1 flex-col overflow-hidden">
          <div className="ember-veil pointer-events-none absolute inset-0 -z-10" aria-hidden="true" />
          {showHeading ? <PageHeading className="shrink-0" showBell={false} /> : null}
          <div
            className={
              isChat
                ? "flex min-h-0 flex-1 flex-col overflow-hidden"
                : "min-h-0 flex-1 overflow-y-auto custom-scrollbar"
            }
          >
            <div className={isChat ? "flex h-full w-full flex-col" : "w-full px-4 py-5 lg:px-8 lg:py-6"}>
              {children}
            </div>
          </div>
        </main>

        {/* Instrument status strip */}
        <StatusBar />
      </div>
    </div>
  );
}
