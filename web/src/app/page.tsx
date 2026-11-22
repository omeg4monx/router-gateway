"use client";

import { AlertTriangle, RefreshCw } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useInitStatus } from "@/features/system/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { Card } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";

export default function HomePage() {
  const router = useRouter();
  const initStatusQuery = useInitStatus();

  useEffect(() => {
    if (!initStatusQuery.data) {
      return;
    }
    if (!initStatusQuery.data.initialized) {
      router.replace("/init/");
    } else {
      router.replace("/dashboard/");
    }
  }, [initStatusQuery.data, router]);

  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12">
      <Card className="w-full max-w-md p-6 sm:p-8">
        <div className="flex items-center gap-3">
          <CaduceusLogo size={30} />
          <div>
            <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON ROUTER
            </p>
            <p className="font-mono text-xs text-[var(--text-muted)]">CONTROL PLANE BOOTSTRAP</p>
          </div>
        </div>

        <EmberRule className="my-5" />

        {initStatusQuery.isError ? (
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

            <div className="flex flex-wrap items-center gap-3">
              <Button type="button" size="sm" onClick={() => void initStatusQuery.refetch()}>
                <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Retry</span>
              </Button>
              <Button asChild variant="outline" size="sm">
                <Link href="/dashboard/">Open Dashboard</Link>
              </Button>
              <Button asChild variant="outline" size="sm">
                <Link href="/init/">Setup</Link>
              </Button>
            </div>
          </div>
        ) : (
          <div className="space-y-2">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              Inspecting initialization state
            </p>
            <p className="font-mono text-xs text-[var(--text-muted)]">
              GET /api/init &rarr; routing to setup or monitor view
            </p>
          </div>
        )}
      </Card>
    </main>
  );
}
