"use client";

import Link from "next/link";
import {
  AlertTriangle,
  ArrowRight,
  KeyRound,
  RefreshCw,
  ShieldCheck,
} from "lucide-react";
import { useInitStatus, useSystemVersion } from "@/features/system/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

export default function LoginPage() {
  const initStatusQuery = useInitStatus();
  const versionQuery = useSystemVersion();

  const isRetrying = initStatusQuery.isFetching || versionQuery.isFetching;

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

          <ThemeToggle />
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
        <Card className="w-full max-w-md p-6 sm:p-8">
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
              <Alert variant="destructive">
                <AlertTriangle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Unable to reach control plane</AlertTitle>
                  <AlertDescription className="font-mono">
                    {initStatusQuery.error.message}
                  </AlertDescription>
                </div>
              </Alert>

              <div className="flex flex-col gap-2.5">
                <Button
                  type="button"
                  onClick={handleRetry}
                  disabled={isRetrying}
                  className="w-full"
                >
                  <RefreshCw className="h-4 w-4" aria-hidden="true" />
                  <span>{isRetrying ? "Retrying..." : "Retry Connection"}</span>
                </Button>
                <Button asChild variant="outline" className="w-full">
                  <Link href="/dashboard/">Open Dashboard</Link>
                </Button>
              </div>
            </div>
          ) : !initStatus?.initialized ? (
            <div className="space-y-4">
              <Alert variant="warning">
                First-run initialization is required before signing in.
              </Alert>
              <Button asChild className="w-full">
                <Link href="/init/">
                  <span>Complete First-Run Setup</span>
                  <ArrowRight className="h-4 w-4" aria-hidden="true" />
                </Link>
              </Button>
            </div>
          ) : (
            <div className="space-y-5">
              <Alert>
                <ShieldCheck
                  className="mt-0.5 h-4 w-4 shrink-0 text-[#E56A4A]"
                  aria-hidden="true"
                />
                <div className="text-xs leading-relaxed text-[var(--text-secondary)]">
                  <AlertTitle>Wave 0 Kernel Posture</AlertTitle>
                  <AlertDescription className="mt-1">
                    Cookie-backed operator login (
                    <span className="font-mono">/api/auth/login</span>) is
                    scheduled for Wave 1. Loopback control plane access is
                    currently direct.
                  </AlertDescription>
                </div>
              </Alert>

              <div>
                <Label
                  htmlFor="login-password-standby"
                  className="text-[var(--text-muted)]"
                >
                  Administrator Password
                </Label>
                <Input
                  id="login-password-standby"
                  type="password"
                  disabled
                  placeholder="Enabled when Wave 1 auth module mounts"
                  className="mt-1.5"
                />
              </div>

              <div className="flex flex-col gap-2.5 pt-1">
                <Button asChild className="w-full">
                  <Link href="/dashboard/">
                    <span>Continue to Dashboard</span>
                    <ArrowRight className="h-4 w-4" aria-hidden="true" />
                  </Link>
                </Button>
                <Button asChild variant="outline" className="w-full">
                  <Link href="/init/">Inspect Initialization Status</Link>
                </Button>
              </div>
            </div>
          )}
        </Card>
      </main>
    </div>
  );
}
