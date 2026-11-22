"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import {
  AlertTriangle,
  ArrowRight,
  KeyRound,
  Moon,
  RefreshCw,
  ShieldCheck,
  Sun,
} from "lucide-react";
import { useInitStatus, useSystemVersion } from "@/features/system/hooks";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { EmberRule } from "@/shared/ui/EmberRule";

export default function LoginPage() {
  const initStatusQuery = useInitStatus();
  const versionQuery = useSystemVersion();
  const [theme, setTheme] = useState<"dark" | "light">("dark");

  useEffect(() => {
    const current = document.documentElement.getAttribute("data-theme");
    if (current === "light" || current === "dark") {
      setTheme(current);
    }
  }, []);

  const toggleTheme = () => {
    const next = theme === "dark" ? "light" : "dark";
    setTheme(next);
    document.documentElement.setAttribute("data-theme", next);
    try {
      window.localStorage.setItem("omeg4mon-theme", next);
    } catch {
      // Ignore storage failures.
    }
  };

  const handleRetry = () => {
    void Promise.all([initStatusQuery.refetch(), versionQuery.refetch()]);
  };

  const initStatus = initStatusQuery.data;
  const version = versionQuery.data;

  const buildVersionLabel = versionQuery.isPending
    ? "Loading..."
    : versionQuery.isError || !version
      ? "Unavailable"
      : version.version;

  const authModuleLabel = initStatusQuery.isPending
    ? "CHECKING"
    : initStatusQuery.isError || !initStatus
      ? "UNAVAILABLE"
      : initStatus.requireLogin
        ? "ENFORCED"
        : "WAVE 1 STANDBY";

  return (
    <div className="grid min-h-screen grid-cols-1 bg-[var(--bg-canvas)] text-[var(--text-primary)] lg:grid-cols-12">
      {/* Left split brand panel (Configure archetype, RHYTHM 2 split composition) */}
      <section className="flex flex-col justify-between border-b border-[var(--border-subtle)] bg-[var(--bg-surface)] p-8 lg:col-span-5 lg:border-b-0 lg:border-r lg:p-12">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={32} />
            <div className="leading-tight">
              <p className="text-base font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON
              </p>
              <p className="font-mono text-xs tracking-widest text-[var(--text-muted)]">
                ROUTER
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={toggleTheme}
            className="inline-flex h-9 items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 text-xs font-medium text-[var(--text-secondary)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]"
            aria-label={theme === "dark" ? "Switch to light theme" : "Switch to dark theme"}
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

        <div className="my-10 space-y-5 lg:my-0">
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Self-Hosted LLM Gateway
          </p>
          <h1 className="text-3xl font-semibold tracking-tight text-[var(--text-primary)]">
            Control Plane Access
          </h1>
          <p className="max-w-md text-sm leading-relaxed text-[var(--text-secondary)]">
            Unified OpenAI and Anthropic compatible routing kernel. Single binary
            runtime with embedded SQLite state and deterministic provider failover.
          </p>
          <EmberRule />
          <dl className="grid grid-cols-2 gap-4 pt-1 text-xs">
            <div>
              <dt className="uppercase tracking-wider text-[var(--text-muted)]">
                Build Version
              </dt>
              <dd className="mt-1 font-mono text-sm text-[var(--text-primary)]">
                {buildVersionLabel}
              </dd>
            </div>
            <div>
              <dt className="uppercase tracking-wider text-[var(--text-muted)]">
                Auth Module
              </dt>
              <dd className="mt-1 font-mono text-sm text-[var(--text-primary)]">
                {authModuleLabel}
              </dd>
            </div>
          </dl>
        </div>

        <p className="font-mono text-xs text-[var(--text-muted)]">
          API CONTRACT &middot; OPENAPI 3.1
        </p>
      </section>

      {/* Right split configure panel */}
      <main className="flex items-center justify-center p-6 sm:p-10 lg:col-span-7">
        <div className="w-full max-w-md rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-6 sm:p-8">
          <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-widest text-[var(--accent-text)]">
            <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Operator Authentication</span>
          </div>
          <h2 className="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
            Sign in to Gateway
          </h2>

          <EmberRule className="my-5" />

          {initStatusQuery.isPending ? (
            <div className="space-y-2 py-2">
              <p className="text-sm font-medium text-[var(--text-primary)]">
                Checking authentication posture
              </p>
              <p className="font-mono text-xs text-[var(--text-muted)]">
                Querying GET /api/init...
              </p>
            </div>
          ) : initStatusQuery.isError ? (
            <div className="space-y-4">
              <div
                role="alert"
                className="flex items-start gap-2.5 rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-3.5 py-3 text-sm"
              >
                <AlertTriangle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <p className="font-medium text-[var(--text-primary)]">
                    Unable to reach control plane
                  </p>
                  <p className="mt-0.5 font-mono text-xs text-[var(--text-secondary)]">
                    {initStatusQuery.error.message}
                  </p>
                </div>
              </div>

              <div className="flex flex-col gap-2.5">
                <button
                  type="button"
                  onClick={handleRetry}
                  disabled={initStatusQuery.isFetching || versionQuery.isFetching}
                  className="inline-flex h-10 w-full items-center justify-center gap-2 rounded-md bg-[#E56A4A] px-4 text-sm font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)] disabled:opacity-60"
                >
                  <RefreshCw className="h-4 w-4" aria-hidden="true" />
                  <span>
                    {initStatusQuery.isFetching ? "Retrying..." : "Retry Connection"}
                  </span>
                </button>
                <Link
                  href="/dashboard/"
                  className="inline-flex h-10 w-full items-center justify-center rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-4 text-sm font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
                >
                  Open Dashboard
                </Link>
              </div>
            </div>
          ) : !initStatus?.initialized ? (
            <div className="space-y-4">
              <div className="rounded-md border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-4 py-3 text-sm text-[var(--text-primary)]">
                First-run initialization is required before signing in.
              </div>
              <Link
                href="/init/"
                className="inline-flex h-10 w-full items-center justify-center gap-2 rounded-md bg-[#E56A4A] px-4 text-sm font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)]"
              >
                <span>Complete First-Run Setup</span>
                <ArrowRight className="h-4 w-4" aria-hidden="true" />
              </Link>
            </div>
          ) : (
            <div className="space-y-5">
              <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
                <div className="flex items-start gap-2.5">
                  <ShieldCheck
                    className="mt-0.5 h-4 w-4 shrink-0 text-[#E56A4A]"
                    aria-hidden="true"
                  />
                  <div className="text-xs leading-relaxed text-[var(--text-secondary)]">
                    <p className="font-medium text-[var(--text-primary)]">
                      Wave 0 Kernel Posture
                    </p>
                    <p className="mt-1">
                      Cookie-backed operator login (<span className="font-mono">/api/auth/login</span>)
                      is scheduled for Wave 1. Loopback control plane access is
                      currently direct.
                    </p>
                  </div>
                </div>
              </div>

              <div>
                <label
                  htmlFor="login-password-standby"
                  className="block text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]"
                >
                  Administrator Password
                </label>
                <input
                  id="login-password-standby"
                  type="password"
                  disabled
                  placeholder="Enabled when Wave 1 auth module mounts"
                  className="mt-1.5 h-10 w-full cursor-not-allowed rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 text-sm text-[var(--text-muted)] placeholder:text-[var(--text-muted)]"
                />
              </div>

              <div className="flex flex-col gap-2.5 pt-1">
                <Link
                  href="/dashboard/"
                  className="inline-flex h-10 w-full items-center justify-center gap-2 rounded-md bg-[#E56A4A] px-4 text-sm font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)]"
                >
                  <span>Continue to Dashboard</span>
                  <ArrowRight className="h-4 w-4" aria-hidden="true" />
                </Link>
                <Link
                  href="/init/"
                  className="inline-flex h-10 w-full items-center justify-center rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-4 text-sm font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
                >
                  Inspect Initialization Status
                </Link>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
