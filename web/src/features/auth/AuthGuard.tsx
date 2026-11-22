"use client";

import { AlertTriangle, RefreshCw } from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { type ReactNode, useEffect } from "react";
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
  const pathname = usePathname();
  const router = useRouter();
  const authSessionQuery = useAuthSession();

  const redirectTarget = authSessionQuery.data
    ? resolvePostAuthRedirect(authSessionToRouteState(authSessionQuery.data), pathname)
    : null;
  const awaitingFreshRedirectCheck = Boolean(redirectTarget && authSessionQuery.isFetching);

  useEffect(() => {
    if (redirectTarget && !authSessionQuery.isFetching) {
      router.replace(redirectTarget);
    }
  }, [authSessionQuery.isFetching, redirectTarget, router]);

  if (authSessionQuery.isPending || awaitingFreshRedirectCheck) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md p-6 sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-xs text-[var(--text-muted)]">SESSION VERIFICATION</p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="space-y-2">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              Verifying operator session
            </p>
            <p className="font-mono text-xs text-[var(--text-muted)]">
              GET /api/auth/me &rarr; checking access policy
            </p>
          </div>
        </Card>
      </main>
    );
  }

  if (authSessionQuery.isError) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md p-6 sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-xs text-[var(--text-muted)]">SESSION VERIFICATION</p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="space-y-4">
            <Alert variant="destructive">
              <AlertTriangle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to verify operator session</AlertTitle>
                <AlertDescription className="font-mono">
                  {authSessionQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>

            <div className="flex flex-wrap items-center gap-3">
              <Button
                type="button"
                size="sm"
                onClick={() => void authSessionQuery.refetch()}
                disabled={authSessionQuery.isFetching}
              >
                <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{authSessionQuery.isFetching ? "Retrying..." : "Retry"}</span>
              </Button>
              <Button asChild variant="outline" size="sm">
                <Link href="/login/">Go to Sign In</Link>
              </Button>
            </div>
          </div>
        </Card>
      </main>
    );
  }

  if (redirectTarget) {
    const destinationLabel =
      redirectTarget === "/init/" ? "first-run initialization" : "operator sign in";

    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        <Card className="w-full max-w-md p-6 sm:p-8">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={30} />
            <div>
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-xs text-[var(--text-muted)]">ACCESS REDIRECT</p>
            </div>
          </div>

          <EmberRule className="my-5" />

          <div className="space-y-4">
            <p className="text-sm text-[var(--text-secondary)]">
              Redirecting to {destinationLabel}...
            </p>
            <Button asChild variant="outline" size="sm">
              <Link href={redirectTarget}>Continue to {redirectTarget}</Link>
            </Button>
          </div>
        </Card>
      </main>
    );
  }

  return <>{children}</>;
}
