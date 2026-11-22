"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState, type ReactNode } from "react";
import {
  Activity,
  Cpu,
  KeyRound,
  Layers,
  LogIn,
  Menu,
  Moon,
  Server,
  ShieldCheck,
  SlidersHorizontal,
  Sun,
  X,
} from "lucide-react";
import { CaduceusLogo } from "./CaduceusLogo";
import { EmberRule } from "./EmberRule";

interface AppShellProps {
  children: ReactNode;
}

interface ActiveNavItem {
  label: string;
  href: string;
  icon: typeof Activity;
}

interface PlannedNavItem {
  label: string;
  icon: typeof Activity;
  phase: string;
}

const ACTIVE_NAV_ITEMS: ActiveNavItem[] = [
  { label: "System Overview", href: "/dashboard/", icon: Activity },
];

const PLANNED_NAV_ITEMS: PlannedNavItem[] = [
  { label: "Providers", icon: Server, phase: "Wave 1" },
  { label: "Models & Routing", icon: Layers, phase: "Wave 1" },
  { label: "API Keys", icon: KeyRound, phase: "Wave 1" },
  { label: "Usage & Telemetry", icon: Cpu, phase: "Wave 1" },
  { label: "Settings", icon: SlidersHorizontal, phase: "Wave 1" },
];

const ACCESS_NAV_ITEMS: ActiveNavItem[] = [
  { label: "Sign In", href: "/login/", icon: LogIn },
  { label: "Initialization", href: "/init/", icon: ShieldCheck },
];

type ThemeMode = "dark" | "light";

function normalizeRoute(path: string): string {
  return path.endsWith("/") ? path : `${path}/`;
}

export function AppShell({ children }: AppShellProps) {
  const pathname = usePathname();
  const [theme, setTheme] = useState<ThemeMode>("dark");
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  useEffect(() => {
    const current = document.documentElement.getAttribute("data-theme");
    if (current === "light" || current === "dark") {
      setTheme(current);
    }
  }, []);

  const toggleTheme = () => {
    const nextTheme: ThemeMode = theme === "dark" ? "light" : "dark";
    setTheme(nextTheme);
    document.documentElement.setAttribute("data-theme", nextTheme);
    try {
      window.localStorage.setItem("omeg4mon-theme", nextTheme);
    } catch {
      // Ignore storage errors in restricted browsing contexts.
    }
  };

  const normalizedPathname = normalizeRoute(pathname);

  const sidebarContent = (
    <div className="flex h-full flex-col bg-[var(--bg-surface)] text-[var(--text-primary)]">
      <div className="flex h-16 items-center justify-between px-5">
        <Link
          href="/dashboard/"
          className="flex items-center gap-3 focus-visible:outline-none"
          onClick={() => setMobileMenuOpen(false)}
        >
          <CaduceusLogo size={26} decorative />
          <div className="flex flex-col leading-tight">
            <span className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON
            </span>
            <span className="font-mono text-[11px] tracking-widest text-[var(--text-muted)]">
              ROUTER
            </span>
          </div>
        </Link>
        <button
          type="button"
          onClick={() => setMobileMenuOpen(false)}
          className="inline-flex h-8 w-8 items-center justify-center rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-[var(--text-secondary)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)] lg:hidden"
          aria-label="Close navigation menu"
        >
          <X className="h-4 w-4" aria-hidden="true" />
        </button>
      </div>

      <EmberRule />

      <nav className="flex-1 overflow-y-auto px-3 py-5" aria-label="Control plane navigation">
        <div className="space-y-1">
          <p className="px-2.5 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Monitor
          </p>
          {ACTIVE_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            const isActive = normalizedPathname === normalizeRoute(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setMobileMenuOpen(false)}
                className={`flex items-center justify-between rounded-md border px-2.5 py-2 text-sm font-medium transition-colors ${
                  isActive
                    ? "border-[#E56A4A] bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                    : "border-transparent text-[var(--text-secondary)] hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                }`}
                aria-current={isActive ? "page" : undefined}
              >
                <span className="flex items-center gap-2.5">
                  <Icon
                    className={`h-4 w-4 ${isActive ? "text-[#E56A4A]" : "text-[var(--text-muted)]"}`}
                    aria-hidden="true"
                  />
                  <span>{item.label}</span>
                </span>
              </Link>
            );
          })}
        </div>

        <div className="mt-6 space-y-1">
          <p className="px-2.5 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Control Plane Modules
          </p>
          {PLANNED_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            return (
              <div
                key={item.label}
                aria-disabled="true"
                className="flex cursor-default items-center justify-between rounded-md border border-transparent px-2.5 py-2 text-sm text-[var(--text-muted)]"
              >
                <span className="flex items-center gap-2.5">
                  <Icon className="h-4 w-4 opacity-75" aria-hidden="true" />
                  <span>{item.label}</span>
                </span>
                <span className="rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
                  {item.phase}
                </span>
              </div>
            );
          })}
        </div>

        <div className="mt-6 space-y-1">
          <p className="px-2.5 pb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
            Access
          </p>
          {ACCESS_NAV_ITEMS.map((item) => {
            const Icon = item.icon;
            const isActive = normalizedPathname === normalizeRoute(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setMobileMenuOpen(false)}
                className={`flex items-center gap-2.5 rounded-md border px-2.5 py-2 text-sm font-medium transition-colors ${
                  isActive
                    ? "border-[#E56A4A] bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                    : "border-transparent text-[var(--text-secondary)] hover:border-[var(--border-subtle)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                }`}
              >
                <Icon className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </div>
      </nav>

      <div className="border-t border-[var(--border-subtle)] px-5 py-3.5">
        <div className="flex items-center justify-between font-mono text-xs text-[var(--text-muted)]">
          <span>KERNEL</span>
          <span className="text-[var(--text-secondary)]">WAVE 0</span>
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex min-h-screen bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      {/* Desktop sidebar */}
      <aside className="hidden w-64 shrink-0 border-r border-[var(--border-subtle)] lg:block">
        <div className="sticky top-0 h-screen">{sidebarContent}</div>
      </aside>

      {/* Mobile navigation drawer (solid opaque surface) */}
      {mobileMenuOpen ? (
        <div className="fixed inset-0 z-40 flex lg:hidden" role="dialog" aria-modal="true" aria-label="Navigation">
          <div className="w-72 max-w-[85vw] border-r border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-xl">
            {sidebarContent}
          </div>
          <button
            type="button"
            className="flex-1 bg-[var(--bg-canvas)]"
            aria-label="Dismiss navigation menu"
            onClick={() => setMobileMenuOpen(false)}
          />
        </div>
      ) : null}

      {/* Main content area */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-20 bg-[var(--bg-surface)]">
          <div className="flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
            <div className="flex items-center gap-3">
              <button
                type="button"
                onClick={() => setMobileMenuOpen(true)}
                className="inline-flex h-9 w-9 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-[var(--text-secondary)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)] lg:hidden"
                aria-label="Open navigation menu"
              >
                <Menu className="h-4 w-4" aria-hidden="true" />
              </button>
              <div>
                <p className="font-mono text-[11px] uppercase tracking-widest text-[var(--text-muted)]">
                  Control Plane
                </p>
                <h1 className="text-base font-semibold tracking-tight text-[var(--text-primary)]">
                  Gateway Monitor
                </h1>
              </div>
            </div>

            <div className="flex items-center gap-2.5">
              <button
                type="button"
                onClick={toggleTheme}
                className="inline-flex h-9 items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 text-xs font-medium text-[var(--text-secondary)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]"
                aria-label={theme === "dark" ? "Switch to warm light theme" : "Switch to charcoal dark theme"}
              >
                {theme === "dark" ? (
                  <>
                    <Sun className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                    <span>Light</span>
                  </>
                ) : (
                  <>
                    <Moon className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                    <span>Dark</span>
                  </>
                )}
              </button>
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
