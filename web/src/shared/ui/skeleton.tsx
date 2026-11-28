import type * as React from "react";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { CaduceusLogo } from "./CaduceusLogo";
import { Card } from "./card";
import { EmberRule } from "./EmberRule";
import { Spinner } from "./spinner";
import { useTopProgressHold } from "./TopProgressBar";

const ROUTE_KPI_KEYS = ["route-kpi-1", "route-kpi-2", "route-kpi-3", "route-kpi-4"] as const;
const ROUTE_ROW_KEYS = ["route-row-1", "route-row-2", "route-row-3"] as const;

export function Skeleton({ className, ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      aria-hidden="true"
      className={cn("loading-shimmer rounded-md bg-[var(--color-surface-2)]", className)}
      {...props}
    />
  );
}

export function RouteLoadingSkeleton() {
  useTopProgressHold();
  const locale = useUIStore((state) => state.locale);
  const label = locale === "id" ? "Memuat halaman" : "Loading view";

  return (
    <div
      role="status"
      aria-busy="true"
      aria-live="polite"
      aria-label={label}
      className="flex min-w-0 flex-col gap-5"
    >
      <span className="sr-only">{label}</span>

      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-2">
          <Skeleton className="h-6 w-44" />
          <Skeleton className="h-3.5 w-72 max-w-full" />
        </div>
        <div className="flex items-center gap-2">
          <Skeleton className="h-9 w-24 rounded-md" />
          <Skeleton className="h-9 w-28 rounded-md" />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
        {ROUTE_KPI_KEYS.map((key) => (
          <div
            key={key}
            className="space-y-2.5 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface)] p-4"
          >
            <Skeleton className="h-3 w-20" />
            <Skeleton className="h-7 w-24" />
          </div>
        ))}
      </div>

      <div className="space-y-3 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface)] p-5">
        <Skeleton className="h-4 w-36" />
        {ROUTE_ROW_KEYS.map((key) => (
          <Skeleton key={key} className="h-12 w-full rounded-lg" />
        ))}
      </div>
    </div>
  );
}

export function PageLoadingScreen() {
  useTopProgressHold();
  const locale = useUIStore((state) => state.locale);
  const isId = locale === "id";
  const statusText = isId ? "Memuat OMEG4MON ROUTER…" : "Loading OMEG4MON ROUTER…";

  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
      <Card
        role="status"
        aria-busy="true"
        aria-live="polite"
        className="w-full max-w-md rounded-xl border-[var(--border-strong)]/80 p-6 shadow-sm sm:p-8"
      >
        <div className="flex items-center gap-3">
          <CaduceusLogo size={30} decorative />
          <div>
            <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON ROUTER
            </p>
            <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Memulai Control Plane" : "Bootstrapping Control Plane"}
            </p>
          </div>
        </div>

        <EmberRule className="my-5" />

        <div className="flex items-center gap-3 py-1">
          <Spinner size="sm" decorative className="text-[var(--accent-primary)]" />
          <p className="text-sm font-medium text-[var(--text-secondary)]">{statusText}</p>
        </div>
      </Card>
    </main>
  );
}
