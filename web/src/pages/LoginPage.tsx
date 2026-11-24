"use client";

import {
  Activity,
  AlertTriangle,
  ArrowRight,
  CheckCircle2,
  Database,
  KeyRound,
  RefreshCw,
  ShieldCheck,
  Sparkles,
} from "lucide-react";
import { useEffect } from "react";
import { Link, useNavigate } from "react-router";
import { useAuthSession } from "@/features/auth/hooks";
import { describeAuthPosture, formatUptime, getLoginCopy } from "@/features/auth/i18n";
import { LoginForm } from "@/features/auth/LoginForm";
import { authSessionToRouteState, resolvePostAuthRedirect } from "@/features/auth/schemas";
import { useSystemHealth, useSystemVersion } from "@/features/system/hooks";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { Card } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { LocaleToggle } from "@/shared/ui/LocaleToggle";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

export function LoginPage() {
  const navigate = useNavigate();
  const locale = useUIStore((state) => state.locale);
  const copy = getLoginCopy(locale);

  const authSessionQuery = useAuthSession();
  const versionQuery = useSystemVersion();
  const healthQuery = useSystemHealth();

  const session = authSessionQuery.data;
  const version = versionQuery.data;
  const health = healthQuery.data;

  const redirectTarget = session
    ? resolvePostAuthRedirect(authSessionToRouteState(session), "/login/")
    : null;

  useEffect(() => {
    if (redirectTarget && !authSessionQuery.isFetching) {
      void navigate(redirectTarget, { replace: true });
    }
  }, [authSessionQuery.isFetching, navigate, redirectTarget]);

  const isRetrying =
    authSessionQuery.isFetching || versionQuery.isFetching || healthQuery.isFetching;

  const handleRetry = () => {
    void Promise.all([authSessionQuery.refetch(), versionQuery.refetch(), healthQuery.refetch()]);
  };

  const buildVersionLabel = versionQuery.isPending
    ? copy.buildVersionLoading
    : versionQuery.isError || !version
      ? copy.buildVersionUnavailable
      : version.version;

  const posture = describeAuthPosture(
    {
      isPending: authSessionQuery.isPending,
      isError: authSessionQuery.isError,
      session,
    },
    locale,
  );

  const gatewayState = healthQuery.isPending
    ? {
        label: copy.gatewayChecking,
        dotClass: "bg-[var(--text-muted)]",
        pillClass: "border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[var(--text-muted)]",
        pulse: false,
      }
    : healthQuery.isError || !health
      ? {
          label: copy.gatewayOffline,
          dotClass: "bg-[var(--status-error)]",
          pillClass:
            "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]",
          pulse: false,
        }
      : health.status === "ok"
        ? {
            label: copy.gatewayOnline,
            dotClass: "bg-[var(--status-ok)]",
            pillClass:
              "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]",
            pulse: true,
          }
        : {
            label: copy.gatewayDegraded,
            dotClass: "bg-[var(--status-warn)]",
            pillClass:
              "border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] text-[var(--status-warn)]",
            pulse: true,
          };

  const databaseStatusLabel =
    health?.db === "ok"
      ? copy.databaseOk
      : health?.db === "error"
        ? copy.databaseError
        : copy.databaseUnknown;

  const uptimeLabel = formatUptime(health?.uptimeSeconds, locale);

  return (
    <div className="relative flex min-h-screen flex-col overflow-x-clip bg-[var(--bg-canvas)] text-[var(--text-primary)] lg:grid lg:grid-cols-12">
      {/* Subtle architectural grid & ambient radial lighting */}
      <div aria-hidden="true" className="pointer-events-none fixed inset-0 overflow-hidden">
        <div
          className="absolute inset-0 opacity-45"
          style={{
            backgroundImage:
              "linear-gradient(to right, color-mix(in srgb, var(--border-strong) 34%, transparent) 1px, transparent 1px), linear-gradient(to bottom, color-mix(in srgb, var(--border-strong) 34%, transparent) 1px, transparent 1px)",
            backgroundSize: "40px 40px",
            maskImage: "radial-gradient(ellipse 80% 70% at 50% 35%, black 25%, transparent 100%)",
          }}
        />
        <div
          className="absolute -left-24 -top-24 h-[26rem] w-[26rem] rounded-full blur-3xl"
          style={{
            background:
              "radial-gradient(circle, color-mix(in srgb, var(--accent-primary) 16%, transparent) 0%, transparent 70%)",
          }}
        />
        <div
          className="absolute -bottom-28 right-0 h-[24rem] w-[24rem] rounded-full blur-3xl"
          style={{
            background:
              "radial-gradient(circle, color-mix(in srgb, var(--color-highlight) 10%, transparent) 0%, transparent 70%)",
          }}
        />
      </div>

      {/* Mobile compact header (visible below lg) */}
      <header className="relative z-10 flex items-center justify-between gap-3 border-b border-[var(--border-subtle)] bg-[var(--bg-surface)]/85 px-4 py-3.5 backdrop-blur-md sm:px-6 lg:hidden">
        <div className="flex items-center gap-2.5">
          <CaduceusLogo size={28} />
          <div className="leading-tight">
            <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON
            </p>
            <p className="font-mono text-[10px] tracking-[0.2em] text-[var(--text-muted)]">
              ROUTER
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <LocaleToggle />
          <ThemeToggle />
        </div>
      </header>

      {/* Left split brand & live telemetry panel (stacks below login card on mobile) */}
      <section className="relative z-10 order-2 flex flex-col justify-between border-t border-[var(--border-subtle)] bg-[var(--bg-surface)]/80 p-6 backdrop-blur-sm sm:p-8 lg:order-1 lg:col-span-5 lg:border-r lg:border-t-0 lg:p-10 xl:p-12">
        {/* Desktop brand lockup */}
        <div className="hidden items-center justify-between gap-3 lg:flex">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={32} />
            <div className="leading-tight">
              <p className="text-base font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON
              </p>
              <p className="font-mono text-[11px] tracking-[0.22em] text-[var(--text-muted)]">
                ROUTER
              </p>
            </div>
          </div>
        </div>

        {/* Hero & live gateway telemetry */}
        <div className="space-y-6 lg:my-auto lg:py-8">
          <div className="flex flex-wrap items-center gap-2">
            <span
              className={cn(
                "inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium",
                gatewayState.pillClass,
              )}
            >
              <span className="relative flex h-2 w-2" aria-hidden="true">
                {gatewayState.pulse ? (
                  <span
                    className={cn(
                      "absolute inline-flex h-full w-full animate-ping rounded-full opacity-60",
                      gatewayState.dotClass,
                    )}
                  />
                ) : null}
                <span
                  className={cn("relative inline-flex h-2 w-2 rounded-full", gatewayState.dotClass)}
                />
              </span>
              <span>
                {copy.gatewayStatusLabel} {gatewayState.label}
              </span>
            </span>

            <span className="inline-flex items-center rounded-full border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/70 px-2.5 py-1 font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
              {copy.buildVersionLabel}: {buildVersionLabel}
            </span>
          </div>

          <div className="space-y-2.5">
            <p className="font-mono text-[11px] font-medium uppercase tracking-[0.16em] text-[var(--accent-text)]">
              {copy.brandEyebrow}
            </p>
            <h1 className="text-2xl font-semibold tracking-tight text-[var(--text-primary)] sm:text-3xl">
              {copy.heroTitle}
            </h1>
            <p className="max-w-md text-sm leading-relaxed text-[var(--text-secondary)]">
              {copy.heroDescription}
            </p>
          </div>

          <EmberRule />

          {/* Standardized 2-column telemetry tiles */}
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div className="flex flex-col justify-between gap-2.5 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/75 p-4">
              <div className="flex items-center justify-between gap-2">
                <span className="font-mono text-[11px] uppercase tracking-[0.12em] text-[var(--text-muted)]">
                  {copy.authPostureLabel}
                </span>
                <ShieldCheck
                  className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
              </div>
              <div className="space-y-2">
                <Badge variant={posture.badgeVariant}>{posture.code}</Badge>
                <p className="text-xs leading-relaxed text-[var(--text-secondary)]">
                  {posture.summary}
                </p>
              </div>
            </div>

            <div className="flex flex-col justify-between gap-2.5 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/75 p-4">
              <div className="flex items-center justify-between gap-2">
                <span className="font-mono text-[11px] uppercase tracking-[0.12em] text-[var(--text-muted)]">
                  {copy.runtimeLabel}
                </span>
                <Activity
                  className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
                  aria-hidden="true"
                />
              </div>
              <dl className="space-y-2 text-xs">
                <div className="flex items-center justify-between gap-2">
                  <dt className="inline-flex items-center gap-1.5 text-[var(--text-muted)]">
                    <Database className="h-3 w-3" aria-hidden="true" />
                    <span>{copy.databaseLabel}</span>
                  </dt>
                  <dd
                    className={cn(
                      "font-mono font-medium",
                      health?.db === "ok"
                        ? "text-[var(--status-ok)]"
                        : health?.db === "error"
                          ? "text-[var(--status-error)]"
                          : "text-[var(--text-muted)]",
                    )}
                  >
                    {databaseStatusLabel}
                  </dd>
                </div>
                <div className="flex items-center justify-between gap-2">
                  <dt className="text-[var(--text-muted)]">{copy.uptimeLabel}</dt>
                  <dd className="font-mono font-medium tabular-nums text-[var(--text-primary)]">
                    {uptimeLabel}
                  </dd>
                </div>
              </dl>
            </div>
          </div>

          {/* Signed session reassurance strip */}
          <div className="flex items-start gap-3 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]/55 p-3.5">
            <ShieldCheck
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--accent-text)]"
              aria-hidden="true"
            />
            <div className="space-y-0.5">
              <p className="text-xs font-medium text-[var(--text-primary)]">
                {copy.securityNoteTitle}
              </p>
              <p className="text-xs leading-relaxed text-[var(--text-secondary)]">
                {copy.securityNoteBody}
              </p>
            </div>
          </div>
        </div>

        <p className="mt-8 font-mono text-[11px] tracking-[0.14em] text-[var(--text-muted)] lg:mt-0">
          {copy.footerContract}
        </p>
      </section>

      {/* Right split operator sign-in canvas (renders above telemetry on mobile) */}
      <main className="relative z-10 order-1 flex flex-1 flex-col justify-between p-5 sm:p-8 lg:order-2 lg:col-span-7 lg:p-10 xl:p-12">
        {/* Desktop utility bar */}
        <div className="hidden items-center justify-end gap-2 lg:flex">
          <LocaleToggle />
          <ThemeToggle />
        </div>

        {/* Centered elevated login card */}
        <div className="flex flex-1 items-center justify-center py-4 sm:py-8 lg:py-0">
          <Card className="w-full max-w-[440px] rounded-xl border border-[var(--border-strong)]/80 bg-[var(--bg-surface)]/95 p-6 shadow-[0_24px_60px_-28px_rgba(0,0,0,0.45)] backdrop-blur-md sm:p-8">
            <div className="flex items-start gap-3.5">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] text-[var(--accent-text)]">
                <KeyRound className="h-4 w-4" aria-hidden="true" />
              </div>
              <div className="min-w-0 space-y-1">
                <p className="font-mono text-[11px] font-medium uppercase tracking-[0.14em] text-[var(--accent-text)]">
                  {copy.cardEyebrow}
                </p>
                <h2 className="text-xl font-semibold tracking-tight text-[var(--text-primary)]">
                  {copy.signInHeading}
                </h2>
              </div>
            </div>

            <p className="mt-3 text-sm leading-relaxed text-[var(--text-secondary)]">
              {copy.signInSubtitle}
            </p>

            <EmberRule className="my-5" />

            {authSessionQuery.isPending ? (
              <div className="space-y-4 py-1" aria-busy="true" aria-live="polite">
                <div className="space-y-2">
                  <div className="h-3.5 w-36 animate-pulse rounded bg-[var(--bg-surface-raised)]" />
                  <div className="h-11 w-full animate-pulse rounded-md bg-[var(--bg-surface-raised)]" />
                </div>
                <div className="h-11 w-full animate-pulse rounded-md bg-[var(--bg-surface-raised)]" />
                <div className="space-y-1 pt-1">
                  <p className="text-xs font-medium text-[var(--text-primary)]">
                    {copy.pendingTitle}
                  </p>
                  <p className="text-xs text-[var(--text-muted)]">{copy.pendingSubtitle}</p>
                </div>
              </div>
            ) : authSessionQuery.isError ? (
              <div className="space-y-4">
                <Alert variant="destructive" className="rounded-lg">
                  <AlertTriangle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 space-y-1">
                    <AlertTitle className="text-xs font-semibold">
                      {copy.unreachableTitle}
                    </AlertTitle>
                    <AlertDescription className="text-xs leading-relaxed text-[var(--text-primary)]">
                      {copy.unreachableHint}
                    </AlertDescription>
                    <p className="truncate font-mono text-[11px] text-[var(--text-muted)]">
                      {authSessionQuery.error.message}
                    </p>
                  </div>
                </Alert>

                <ul className="space-y-1.5 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/75 p-3.5 text-xs leading-relaxed text-[var(--text-secondary)]">
                  {copy.unreachableTips.map((tip) => (
                    <li key={tip} className="flex items-start gap-2">
                      <span
                        className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-[var(--accent-primary)]"
                        aria-hidden="true"
                      />
                      <span>{tip}</span>
                    </li>
                  ))}
                </ul>

                <div className="flex flex-col gap-2.5">
                  <Button
                    type="button"
                    onClick={handleRetry}
                    disabled={isRetrying}
                    className="h-11 w-full"
                  >
                    <RefreshCw
                      className={cn("h-4 w-4", isRetrying && "animate-spin")}
                      aria-hidden="true"
                    />
                    <span>{isRetrying ? copy.retryingLabel : copy.retryLabel}</span>
                  </Button>
                  <Button asChild variant="outline" className="h-10 w-full">
                    <Link to="/dashboard/">{copy.openDashboardLabel}</Link>
                  </Button>
                </div>
              </div>
            ) : !session?.initialized ? (
              <div className="space-y-4">
                <Alert variant="warning" className="rounded-lg">
                  <Sparkles
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                    aria-hidden="true"
                  />
                  <div className="space-y-1">
                    <AlertTitle className="text-xs font-semibold">
                      {copy.uninitializedTitle}
                    </AlertTitle>
                    <AlertDescription className="text-xs leading-relaxed text-[var(--text-primary)]">
                      {copy.uninitializedBody}
                    </AlertDescription>
                  </div>
                </Alert>

                <Button asChild className="h-11 w-full">
                  <Link to="/init/">
                    <span>{copy.uninitializedCta}</span>
                    <ArrowRight className="h-4 w-4" aria-hidden="true" />
                  </Link>
                </Button>
              </div>
            ) : redirectTarget === "/dashboard/" ? (
              <div className="space-y-4">
                <div className="flex items-start gap-3 rounded-lg border border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] p-4">
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div className="space-y-1">
                    <p className="text-xs font-semibold text-[var(--text-primary)]">
                      {session.authenticated
                        ? copy.redirectAuthenticatedTitle
                        : copy.redirectLoopbackTitle}
                    </p>
                    <p className="text-xs leading-relaxed text-[var(--text-secondary)]">
                      {copy.redirectSubtitle}
                    </p>
                  </div>
                </div>

                <Button asChild className="h-11 w-full">
                  <Link to="/dashboard/">
                    <span>{copy.redirectCta}</span>
                    <ArrowRight className="h-4 w-4" aria-hidden="true" />
                  </Link>
                </Button>
              </div>
            ) : (
              <LoginForm remoteEnforced={session.remote_enforced} />
            )}
          </Card>
        </div>

        {/* Optical balance spacer matching desktop utility bar height */}
        <div className="hidden h-9 lg:block" aria-hidden="true" />
      </main>
    </div>
  );
}
