import { memo } from "react";
import { cn } from "@/shared/lib/utils";
import { useUsageStream } from "./hooks";

interface UsageStreamIndicatorProps {
  isId: boolean;
}

export const UsageStreamIndicator = memo(function UsageStreamIndicator({
  isId,
}: UsageStreamIndicatorProps) {
  const { status } = useUsageStream(true);

  const label =
    status === "live"
      ? isId
        ? "Langsung"
        : "Live"
      : status === "error"
        ? isId
          ? "Menyambung ulang"
          : "Reconnecting"
        : isId
          ? "Menghubungkan"
          : "Connecting";

  const dotClass =
    status === "live"
      ? "bg-[var(--status-ok)] motion-safe:animate-pulse"
      : status === "error"
        ? "bg-[var(--status-error)]"
        : "bg-[var(--status-warn)]";

  const pillClass =
    status === "live"
      ? "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]"
      : status === "error"
        ? "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]"
        : "border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] text-[var(--status-warn)]";

  return (
    <span
      role="status"
      aria-live="polite"
      className={cn(
        "inline-flex items-center gap-1.5 rounded border px-2 py-0.5 font-mono text-[10px] font-medium uppercase tracking-wider",
        pillClass,
      )}
    >
      <span className={cn("h-1.5 w-1.5 shrink-0 rounded-full", dotClass)} aria-hidden="true" />
      <span>{label}</span>
    </span>
  );
});
