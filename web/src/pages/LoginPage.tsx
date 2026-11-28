import { AlertTriangle, ArrowRight, CheckCircle2, RefreshCw } from "lucide-react";
import { useEffect } from "react";
import { Link, useNavigate } from "react-router";
import { useAuthSession } from "@/features/auth/hooks";
import { getLoginCopy } from "@/features/auth/i18n";
import { LoginForm } from "@/features/auth/LoginForm";
import { authSessionToRouteState, resolvePostAuthRedirect } from "@/features/auth/schemas";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { Card } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { LocaleToggle } from "@/shared/ui/LocaleToggle";
import { Skeleton } from "@/shared/ui/skeleton";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

export function LoginPage() {
  const navigate = useNavigate();
  const locale = useUIStore((state) => state.locale);
  const copy = getLoginCopy(locale);

  const authSessionQuery = useAuthSession();
  const session = authSessionQuery.data;

  const redirectTarget = session
    ? resolvePostAuthRedirect(authSessionToRouteState(session), "/login/")
    : null;

  useEffect(() => {
    if (redirectTarget && !authSessionQuery.isFetching) {
      void navigate(redirectTarget, { replace: true });
    }
  }, [authSessionQuery.isFetching, navigate, redirectTarget]);

  const isRetrying = authSessionQuery.isFetching;

  const handleRetry = () => {
    void authSessionQuery.refetch();
  };

  return (
    <div className="flex min-h-screen flex-col bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      <header className="flex items-center justify-end gap-2 px-4 py-4 sm:px-6">
        <LocaleToggle />
        <ThemeToggle />
      </header>

      <main className="flex flex-1 items-center justify-center px-4 pb-12 sm:px-6">
        <Card className="w-full max-w-[380px] rounded-xl border border-[var(--border-strong)] bg-[var(--bg-surface)] p-6 sm:p-7">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={28} />
            <div className="min-w-0">
              <p className="text-xs font-semibold tracking-wider text-[var(--text-muted)]">
                OMEG4MON ROUTER
              </p>
              <h1 className="text-lg font-semibold tracking-tight text-[var(--text-primary)]">
                {copy.signInHeading}
              </h1>
            </div>
          </div>

          <EmberRule className="my-5" />

          {authSessionQuery.isPending ? (
            <div className="space-y-4" role="status" aria-busy="true" aria-live="polite">
              <div className="space-y-2">
                <Skeleton className="h-3.5 w-24" />
                <Skeleton className="h-11 w-full rounded-md" />
              </div>
              <Skeleton className="h-11 w-full rounded-md" />
              <p className="sr-only">{copy.pendingLabel}</p>
            </div>
          ) : authSessionQuery.isError ? (
            <div className="space-y-4">
              <Alert variant="destructive" className="rounded-lg">
                <AlertTriangle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div className="min-w-0 space-y-1">
                  <AlertTitle className="text-xs font-semibold">{copy.unreachableTitle}</AlertTitle>
                  <AlertDescription className="truncate font-mono text-[11px] text-[var(--text-muted)]">
                    {authSessionQuery.error.message}
                  </AlertDescription>
                </div>
              </Alert>

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
            </div>
          ) : !session?.initialized ? (
            <div className="space-y-4">
              <Alert variant="warning" className="rounded-lg">
                <AlertTriangle
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
              <div className="flex items-center gap-2.5 rounded-lg border border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] px-3.5 py-3 text-xs font-medium text-[var(--text-primary)]">
                <CheckCircle2
                  className="h-4 w-4 shrink-0 text-[var(--status-ok)]"
                  aria-hidden="true"
                />
                <span>{copy.redirectTitle}</span>
              </div>

              <Button asChild className="h-11 w-full">
                <Link to="/dashboard/">
                  <span>{copy.redirectCta}</span>
                  <ArrowRight className="h-4 w-4" aria-hidden="true" />
                </Link>
              </Button>
            </div>
          ) : (
            <LoginForm />
          )}
        </Card>
      </main>
    </div>
  );
}
