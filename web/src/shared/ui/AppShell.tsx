"use client";

import {
  Activity,
  ArrowRightLeft,
  Bell,
  Cpu,
  Image,
  KeyRound,
  Layers,
  Lock,
  LogIn,
  LogOut,
  Menu,
  MessageSquare,
  Network,
  PanelLeftClose,
  PanelLeftOpen,
  Radio,
  Send,
  Server,
  Shield,
  ShieldAlert,
  ShieldCheck,
  Sliders,
  SlidersHorizontal,
  Swords,
  Terminal,
  Ticket,
  Users,
  Workflow,
  X,
  Zap,
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { type ReactNode, useEffect } from "react";
import { useAuthSession, useLogoutMutation } from "@/features/auth/hooks";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Button } from "./button";
import { CaduceusLogo } from "./CaduceusLogo";
import { EmberRule } from "./EmberRule";
import { ThemeToggle } from "./ThemeToggle";

interface AppShellProps {
  children: ReactNode;
}

interface ActiveNavItem {
  label: string;
  href: string;
  icon: typeof Activity;
}

const ACTIVE_NAV_ITEMS: ActiveNavItem[] = [
  { label: "System Overview", href: "/dashboard/", icon: Activity },
  { label: "Providers", href: "/dashboard/providers/", icon: Server },
  { label: "Models & Routing", href: "/dashboard/models/", icon: Layers },
  { label: "Proxy Pools", href: "/dashboard/proxy-pools/", icon: Network },
  { label: "API Keys", href: "/dashboard/keys/", icon: KeyRound },
  { label: "Teams", href: "/dashboard/teams/", icon: Users },
  { label: "Vouchers", href: "/dashboard/vouchers/", icon: Ticket },
  { label: "Media Gateway", href: "/dashboard/media/", icon: Image },
  { label: "Tunnel", href: "/dashboard/tunnel/", icon: Radio },
  { label: "MITM Proxy", href: "/dashboard/mitm/", icon: Shield },
  { label: "PXPIPE", href: "/dashboard/pxpipe/", icon: Workflow },
  { label: "Telegram Bot", href: "/dashboard/telegram/", icon: Send },
  { label: "Model Arena", href: "/dashboard/arena/", icon: Swords },
  { label: "Basic Chat", href: "/dashboard/chat/", icon: MessageSquare },
  { label: "Usage & Telemetry", href: "/dashboard/usage/", icon: Cpu },
  { label: "Prompt Guard", href: "/dashboard/prompt-guard/", icon: ShieldAlert },
  { label: "Model Editor", href: "/dashboard/model-editor/", icon: Sliders },
  { label: "Token Saver", href: "/dashboard/token-saver/", icon: Zap },
  { label: "Translator", href: "/dashboard/translator/", icon: ArrowRightLeft },
  { label: "Security", href: "/dashboard/security/", icon: Lock },
  { label: "Notifications & Jobs", href: "/dashboard/notifications/", icon: Bell },
  { label: "CLI Tools", href: "/dashboard/cli-tools/", icon: Terminal },
  { label: "Settings", href: "/dashboard/settings/", icon: SlidersHorizontal },
];

const ACCESS_NAV_ITEMS: ActiveNavItem[] = [
  { label: "Sign In", href: "/login/", icon: LogIn },
  { label: "Initialization", href: "/init/", icon: ShieldCheck },
];

function normalizeRoute(path: string): string {
  return path.endsWith("/") ? path : `${path}/`;
}

function isNavItemActive(normalizedPathname: string, itemHref: string): boolean {
  const normalizedItem = normalizeRoute(itemHref);
  if (normalizedItem === "/dashboard/") {
    return normalizedPathname === "/dashboard/";
  }
  return normalizedPathname === normalizedItem || normalizedPathname.startsWith(normalizedItem);
}

export function AppShell({ children }: AppShellProps) {
  const pathname = usePathname();
  const router = useRouter();
  const authSessionQuery = useAuthSession();
  const logoutMutation = useLogoutMutation();
  const mobileNavOpen = useUIStore((state) => state.mobileNavOpen);
  const setMobileNavOpen = useUIStore((state) => state.setMobileNavOpen);
  const closeMobileNav = useUIStore((state) => state.closeMobileNav);
  const sidebarCollapsed = useUIStore((state) => state.sidebarCollapsed);
  const toggleSidebarCollapsed = useUIStore((state) => state.toggleSidebarCollapsed);

  const isAuthenticated = Boolean(authSessionQuery.data?.authenticated);

  const handleSignOut = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        closeMobileNav();
        router.replace("/login/");
      },
    });
  };

  useEffect(() => {
    if (!mobileNavOpen) {
      return;
    }
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        closeMobileNav();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [mobileNavOpen, closeMobileNav]);

  const normalizedPathname = normalizeRoute(pathname);
  const activeRouteLabel =
    ACTIVE_NAV_ITEMS.find((item) => isNavItemActive(normalizedPathname, item.href))?.label ??
    "Gateway Monitor";

  const renderSidebar = (compact: boolean) => (
    <div className="flex h-full flex-col bg-[var(--bg-surface)] text-[var(--text-primary)]">
      <div
        className={cn(
          "flex h-16 items-center justify-between",
          compact ? "justify-center px-3" : "px-5",
        )}
      >
        <Link
          href="/dashboard/"
          className="flex items-center gap-3 rounded-sm"
          onClick={closeMobileNav}
          aria-label="OMEG4MON ROUTER Dashboard"
        >
          <CaduceusLogo size={26} decorative />
          {!compact ? (
            <div className="flex flex-col leading-tight">
              <span className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON
              </span>
              <span className="font-mono text-[11px] tracking-widest text-[var(--text-muted)]">
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
            aria-label="Close navigation menu"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        ) : null}
      </div>

      <EmberRule />

      <nav className="flex-1 overflow-y-auto px-3 py-5" aria-label="Control plane navigation">
        <div className="space-y-1">
          {!compact ? (
            <p className="px-2.5 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Monitor
            </p>
          ) : null}
          {ACTIVE_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            const isActive = isNavItemActive(normalizedPathname, item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={closeMobileNav}
                title={compact ? item.label : undefined}
                className={cn(
                  "flex items-center rounded-md border px-2.5 py-2 text-sm font-medium transition-colors",
                  compact ? "justify-center" : "justify-between",
                  isActive
                    ? "border-[#E56A4A] bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                    : "border-transparent text-[var(--text-secondary)] hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]",
                )}
                aria-current={isActive ? "page" : undefined}
              >
                <span className="flex items-center gap-2.5">
                  <Icon
                    className={cn(
                      "h-4 w-4 shrink-0",
                      isActive ? "text-[#E56A4A]" : "text-[var(--text-muted)]",
                    )}
                    aria-hidden="true"
                  />
                  {!compact ? <span>{item.label}</span> : null}
                </span>
              </Link>
            );
          })}
        </div>

        <div className="mt-6 space-y-1">
          {!compact ? (
            <p className="px-2.5 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Access
            </p>
          ) : null}
          {ACCESS_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            const isActive = normalizedPathname === normalizeRoute(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={closeMobileNav}
                title={compact ? item.label : undefined}
                aria-current={isActive ? "page" : undefined}
                className={cn(
                  "flex items-center gap-2.5 rounded-md border px-2.5 py-2 text-sm font-medium transition-colors",
                  compact ? "justify-center" : "",
                  isActive
                    ? "border-[#E56A4A] bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                    : "border-transparent text-[var(--text-secondary)] hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]",
                )}
              >
                <Icon className="h-4 w-4 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
                {!compact ? <span>{item.label}</span> : null}
              </Link>
            );
          })}
        </div>
      </nav>

      <div
        className={cn(
          "border-t border-[var(--border-subtle)] py-3.5",
          compact ? "px-2 text-center" : "px-5",
        )}
      >
        <div
          className={cn(
            "flex items-center font-mono text-xs text-[var(--text-muted)]",
            compact ? "justify-center" : "justify-between",
          )}
        >
          {compact ? (
            <span className="text-[var(--text-secondary)]">W3D</span>
          ) : (
            <>
              <span>KERNEL</span>
              <span className="text-[var(--text-secondary)]">WAVE 3D</span>
            </>
          )}
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex min-h-screen bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      {/* Desktop sidebar */}
      <aside
        className={cn(
          "hidden shrink-0 border-r border-[var(--border-subtle)] transition-[width] lg:block",
          sidebarCollapsed ? "w-[72px]" : "w-64",
        )}
      >
        <div className="sticky top-0 h-screen">{renderSidebar(sidebarCollapsed)}</div>
      </aside>

      {/* Mobile navigation drawer (solid opaque surface) */}
      {mobileNavOpen ? (
        <div
          className="fixed inset-0 z-40 flex lg:hidden"
          role="dialog"
          aria-modal="true"
          aria-label="Navigation"
        >
          <div className="w-72 max-w-[85vw] border-r border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-xl">
            {renderSidebar(false)}
          </div>
          <button
            type="button"
            className="flex-1 bg-[var(--bg-canvas)]"
            aria-label="Dismiss navigation menu"
            onClick={closeMobileNav}
          />
        </div>
      ) : null}

      {/* Main content area */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-20 bg-[var(--bg-surface)]">
          <div className="flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
            <div className="flex items-center gap-3">
              <Button
                type="button"
                variant="secondary"
                size="icon"
                onClick={() => setMobileNavOpen(true)}
                className="lg:hidden"
                aria-label="Open navigation menu"
              >
                <Menu className="h-4 w-4" aria-hidden="true" />
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="icon"
                onClick={toggleSidebarCollapsed}
                className="hidden lg:inline-flex"
                aria-label={sidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
              >
                {sidebarCollapsed ? (
                  <PanelLeftOpen className="h-4 w-4" aria-hidden="true" />
                ) : (
                  <PanelLeftClose className="h-4 w-4" aria-hidden="true" />
                )}
              </Button>
              <div>
                <p className="font-mono text-[11px] uppercase tracking-widest text-[var(--text-muted)]">
                  Control Plane
                </p>
                <h1 className="text-base font-semibold tracking-tight text-[var(--text-primary)]">
                  {activeRouteLabel}
                </h1>
              </div>
            </div>

            <div className="flex items-center gap-2.5">
              {isAuthenticated ? (
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={handleSignOut}
                  disabled={logoutMutation.isPending}
                  title={logoutMutation.error?.message}
                >
                  <LogOut className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{logoutMutation.isPending ? "Signing out..." : "Sign Out"}</span>
                </Button>
              ) : null}
              <ThemeToggle />
            </div>
          </div>
          <EmberRule />
        </header>

        <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-6 sm:px-6 lg:px-8">
          {children}
        </main>
      </div>
    </div>
  );
}
