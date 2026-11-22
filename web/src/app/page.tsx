"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { AlertTriangle, RefreshCw } from "lucide-react";
import { useInitStatus } from "@/features/system/hooks";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
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
      <div className="w-full max-w-md rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-6 sm:p-8">
        <div className="flex items-center gap-3">
          <CaduceusLogo size={30} />
          <div>
            <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON ROUTER
            </p>
            <p className="font-mono text-xs text-[var(--text-muted)]">
              CONTROL PLANE BOOTSTRAP
            </p>
          </div>
        </div>

        <EmberRule className="my-5" />

        {initStatusQuery.isError ? (
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

            <div className="flex flex-wrap items-center gap-3">
              <button
                type="button"
                onClick={() => void initStatusQuery.refetch()}
                className="inline-flex h-9 items-center gap-2 rounded-md bg-[#E56A4A] px-4 text-xs font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)]"
              >
                <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Retry</span>
              </button>
              <Link
                href="/dashboard/"
                className="inline-flex h-9 items-center rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-3.5 text-xs font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
              >
                Open Dashboard
              </Link>
              <Link
                href="/init/"
                className="inline-flex h-9 items-center rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-3.5 text-xs font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
              >
                Setup
              </Link>
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
      </div>
    </main>
  );
}
