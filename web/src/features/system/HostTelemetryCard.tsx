import { Cpu, HardDrive } from "lucide-react";
import { memo } from "react";
import { cn } from "@/shared/lib/utils";
import { Skeleton } from "@/shared/ui/skeleton";
import { useSystemTelemetry } from "./hooks";

function formatBytesGB(bytes: number | undefined): string {
  if (bytes === undefined || !Number.isFinite(bytes) || bytes < 0) {
    return "—";
  }
  const gb = bytes / 1024 ** 3;
  return `${gb.toFixed(1)} GB`;
}

function formatHostUptimeShort(sec: number | undefined, isId: boolean): string {
  if (sec === undefined || !Number.isFinite(sec) || sec < 0) {
    return "—";
  }
  const d = Math.floor(sec / 86_400);
  const h = Math.floor((sec % 86_400) / 3_600);
  const m = Math.floor((sec % 3_600) / 60);
  if (d > 0) return isId ? `${d} hari ${h} jam` : `${d}d ${h}h`;
  if (h > 0) return isId ? `${h} jam ${m} mnt` : `${h}h ${m}m`;
  return isId ? `${m} mnt` : `${m}m`;
}

export interface HostTelemetryCardProps {
  isId: boolean;
  healthDb?: string;
  healthUptimeSeconds?: number;
}

export const HostTelemetryCard = memo(function HostTelemetryCard({
  isId,
  healthDb,
  healthUptimeSeconds,
}: HostTelemetryCardProps) {
  const telemetryQuery = useSystemTelemetry();
  const sys = telemetryQuery.data;

  const cpuPct =
    sys && Number.isFinite(sys.cpuUsagePercent)
      ? Math.min(100, Math.max(0, sys.cpuUsagePercent))
      : null;
  const ramPct =
    sys && sys.memTotalBytes > 0
      ? Math.min(100, Math.max(0, Math.round((sys.memUsedBytes / sys.memTotalBytes) * 100)))
      : null;
  const diskPct =
    sys && sys.diskTotalBytes > 0
      ? Math.min(100, Math.max(0, Math.round((sys.diskUsedBytes / sys.diskTotalBytes) * 100)))
      : null;
  const hasLoadAvg = Boolean(sys && sys.loadAvg.length >= 3 && sys.loadAvg.some((v) => v > 0));
  const isTelemetryLive = Boolean(sys && !telemetryQuery.isError);

  return (
    <section className="flex min-w-0 flex-col justify-between gap-4 rounded-xl border border-[var(--color-border)]/80 bg-[var(--color-surface)] p-4 shadow-[0_1px_2px_rgba(11,13,18,0.03)] sm:p-5">
      <div className="space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex flex-wrap items-center gap-2">
            <Cpu className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <span className="text-sm font-semibold text-[var(--color-text-main)]">
              {isId ? "Spesifikasi Mesin & OS" : "Host & OS Specifications"}
            </span>
            {sys?.hostname ? (
              <span className="rounded-full border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] text-[var(--color-text-muted)]">
                {sys.hostname}
              </span>
            ) : null}
            {sys?.os ? (
              <span className="rounded-full border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[11px] text-[var(--color-text-muted)]">
                {sys.os}
              </span>
            ) : null}
          </div>
          {isTelemetryLive ? (
            <span
              className="inline-flex items-center gap-1.5 rounded-full border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[10px] font-semibold tracking-wider text-[var(--status-ok)]"
              title={isId ? "Diperbarui otomatis setiap 2 detik" : "Auto-updating every 2s"}
            >
              <span className="relative flex h-1.5 w-1.5" aria-hidden="true">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[var(--status-ok)] opacity-75 motion-reduce:animate-none" />
                <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-[var(--status-ok)]" />
              </span>
              LIVE
            </span>
          ) : null}
        </div>

        {telemetryQuery.isPending ? (
          <div className="grid grid-cols-1 gap-3 md:grid-cols-3" role="status" aria-busy="true">
            <Skeleton className="h-12 rounded-lg" />
            <Skeleton className="h-12 rounded-lg" />
            <Skeleton className="h-12 rounded-lg" />
          </div>
        ) : !sys ? (
          <p className="text-xs text-[var(--color-text-muted)]">
            {isId
              ? "Data telemetri mesin belum tersedia. Pastikan backend Go terbaru sedang berjalan."
              : "Host telemetry unavailable. Restart the Go backend to enable /api/system/telemetry."}
          </p>
        ) : (
          <div className="space-y-4">
            <div className="flex flex-wrap items-baseline justify-between gap-2">
              <p className="text-sm font-semibold text-[var(--color-text-main)]">{sys.cpuModel}</p>
              <p className="font-mono text-xs tabular-nums text-[var(--color-text-muted)]">
                {sys.cpuCores} {isId ? "core" : "cores"}
                {hasLoadAvg ? ` · Load ${sys.loadAvg.map((v) => v.toFixed(2)).join(" / ")}` : ""}
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              {/* CPU Usage Bar */}
              <div className="flex flex-col gap-1.5">
                <div className="flex items-center justify-between text-xs text-[var(--color-text-muted)]">
                  <span className="font-medium text-[var(--color-text-main)]">CPU</span>
                  <span className="font-mono tabular-nums">
                    {cpuPct !== null ? `${cpuPct.toFixed(1)}%` : "—"}
                  </span>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                  <div
                    className="h-full rounded-full bg-[var(--color-primary)] transition-[width] duration-700 ease-[var(--ease-spring)] will-change-[width]"
                    style={{ width: `${cpuPct ?? 0}%` }}
                  />
                </div>
              </div>

              {/* RAM Usage Bar */}
              <div className="flex flex-col gap-1.5">
                <div className="flex items-center justify-between text-xs text-[var(--color-text-muted)]">
                  <span className="font-medium text-[var(--color-text-main)]">RAM</span>
                  <span className="font-mono tabular-nums">
                    {formatBytesGB(sys.memUsedBytes)} / {formatBytesGB(sys.memTotalBytes)}
                    {ramPct !== null ? ` (${ramPct}%)` : ""}
                  </span>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                  <div
                    className="h-full rounded-full bg-[var(--color-primary)] transition-[width] duration-700 ease-[var(--ease-spring)] will-change-[width]"
                    style={{ width: `${ramPct ?? 0}%` }}
                  />
                </div>
              </div>

              {/* Disk Usage Bar */}
              <div className="flex flex-col gap-1.5">
                <div className="flex items-center justify-between text-xs text-[var(--color-text-muted)]">
                  <span className="inline-flex items-center gap-1.5 font-medium text-[var(--color-text-main)]">
                    <HardDrive
                      className="h-3.5 w-3.5 text-[var(--color-text-muted)]"
                      aria-hidden="true"
                    />
                    <span>Disk /</span>
                  </span>
                  <span className="font-mono tabular-nums">
                    {sys.diskTotalBytes > 0
                      ? `${formatBytesGB(sys.diskUsedBytes)} / ${formatBytesGB(sys.diskTotalBytes)}${diskPct !== null ? ` (${diskPct}%)` : ""}`
                      : "—"}
                  </span>
                </div>
                <div className="h-2 w-full overflow-hidden rounded-full bg-[var(--color-surface-2)]">
                  <div
                    className="h-full rounded-full bg-[var(--color-primary)] transition-[width] duration-700 ease-[var(--ease-spring)] will-change-[width]"
                    style={{ width: `${diskPct ?? 0}%` }}
                  />
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Host Footer Metadata */}
      <div className="flex flex-wrap items-center gap-x-5 gap-y-1.5 border-t border-[var(--color-border-subtle)] pt-3 text-xs text-[var(--color-text-muted)]">
        <span>
          {isId ? "Uptime OS" : "Host Uptime"}:{" "}
          <strong className="font-mono font-semibold text-[var(--color-text-main)]">
            {formatHostUptimeShort(sys?.hostUptimeSeconds, isId)}
          </strong>
        </span>
        <span>
          {isId ? "Uptime Router" : "App Uptime"}:{" "}
          <strong className="font-mono font-semibold text-[var(--color-text-main)]">
            {formatHostUptimeShort(sys?.appUptimeSeconds ?? healthUptimeSeconds, isId)}
          </strong>
        </span>
        <span>
          Database:{" "}
          <strong
            className={cn(
              "font-mono font-semibold",
              (sys?.dbStatus ?? healthDb) === "ok"
                ? "text-[var(--status-ok)]"
                : "text-[var(--status-error)]",
            )}
          >
            {(sys?.dbStatus ?? healthDb) === "ok" ? "● OK" : "● Error"}
          </strong>
        </span>
        {sys?.goVersion ? (
          <span className="font-mono text-[11px] text-[var(--color-text-subtle)]">
            {sys.goVersion}
          </span>
        ) : null}
      </div>
    </section>
  );
});
