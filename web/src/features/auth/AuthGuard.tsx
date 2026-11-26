import { AlertTriangle, ArrowRight, Loader2, RefreshCw } from "lucide-react";
import { type ReactNode, useEffect } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { Card } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { useAuthSession } from "./hooks";
import { authSessionToRouteState, resolvePostAuthRedirect } from "./schemas";

interface AuthGuardProps {
  children: ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";
  const authSessionQuery = useAuthSession();

  const redirectTarget = authSessionQuery.data
    ? resolvePostAuthRedirect(authSessionToRouteState(authSessionQuery.data), pathname)
    : null;
  const awaitingFreshRedirectCheck = Boolean(redirectTarget && authSessionQuery.isFetching);

  useEffect(() => {
    if (redirectTarget && !authSessionQuery.isFetching) {
      void navigate(redirectTarget, { replace: true });
    }
  }, [authSessionQuery.isFetching, navigate, redirectTarget]);

  if (authSessionQuery.isPending || awaitingFreshRedirectCheck) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md rounded-xl border-[var(--border-strong)]/80 p-6 shadow-sm sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                {isId ? "Verifikasi Sesi" : "Session Verification"}
              </p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="flex items-center gap-3 py-1" aria-busy="true" aria-live="polite">
            <Loader2
              className="h-4 w-4 shrink-0 animate-spin text-[var(--accent-primary)]"
              aria-hidden="true"
            />
            <div className="space-y-0.5">
              <p className="text-sm font-medium text-[var(--text-primary)]">
                {isId ? "Memverifikasi sesi operator" : "Verifying operator session"}
              </p>
              <p className="text-xs text-[var(--text-muted)]">
                {isId
                  ? "Memeriksa kebijakan akses control plane…"
                  : "Checking control plane access policy…"}
              </p>
            </div>
          </div>
        </Card>
      </main>
    );
  }

  if (authSessionQuery.isError) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md rounded-xl border-[var(--border-strong)]/80 p-6 shadow-sm sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                {isId ? "Verifikasi Sesi" : "Session Verification"}
              </p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="space-y-4">
            <Alert variant="destructive" className="rounded-lg">
              <AlertTriangle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle className="text-xs font-semibold">
                  {isId ? "Tidak dapat memverifikasi sesi" : "Unable to verify operator session"}
                </AlertTitle>
                <AlertDescription className="font-mono">
                  {authSessionQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>

            <div className="flex flex-wrap items-center gap-2.5">
              <Button
                type="button"
                size="sm"
                onClick={() => void authSessionQuery.refetch()}
                disabled={authSessionQuery.isFetching}
              >
                <RefreshCw
                  className={cn("h-3.5 w-3.5", authSessionQuery.isFetching && "animate-spin")}
                  aria-hidden="true"
                />
                <span>
                  {authSessionQuery.isFetching
                    ? isId
                      ? "Menghubungkan…"
                      : "Retrying…"
                    : isId
                      ? "Coba Lagi"
                      : "Retry"}
                </span>
              </Button>
              <Button asChild variant="outline" size="sm">
                <Link to="/login/">{isId ? "Ke Halaman Masuk" : "Go to Sign In"}</Link>
              </Button>
            </div>
          </div>
        </Card>
      </main>
    );
  }

  if (redirectTarget) {
    const destinationLabel =
      redirectTarget === "/init/"
        ? isId
          ? "pengaturan awal"
          : "first-run initialization"
        : isId
          ? "halaman masuk operator"
          : "operator sign in";

    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md rounded-xl border-[var(--border-strong)]/80 p-6 shadow-sm sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                {isId ? "Pengalihan Akses" : "Access Redirect"}
              </p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="space-y-4">
            <p className="text-sm text-[var(--text-secondary)]">
              {isId ? `Mengarahkan ke ${destinationLabel}…` : `Redirecting to ${destinationLabel}…`}
            </p>
            <Button asChild variant="outline" size="sm">
              <Link to={redirectTarget}>
                <span>
                  {isId ? `Lanjut ke ${redirectTarget}` : `Continue to ${redirectTarget}`}
                </span>
                <ArrowRight className="h-3.5 w-3.5" aria-hidden="true" />
              </Link>
            </Button>
          </div>
        </Card>
      </main>
    );
  }

  return <>{children}</>;
}
