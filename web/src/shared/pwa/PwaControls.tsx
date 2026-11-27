import { Download, RefreshCw, WifiOff } from "lucide-react";
import { useUIStore } from "@/shared/store/uiStore";
import { Button } from "@/shared/ui/button";
import { getPwaCopy, type PwaStatus } from "./usePwaStatus";

interface PwaInstallButtonProps {
  canInstall: boolean;
  isInstalling: boolean;
  onInstall: () => void;
}

export function PwaInstallButton({ canInstall, isInstalling, onInstall }: PwaInstallButtonProps) {
  const locale = useUIStore((state) => state.locale);
  if (!canInstall) {
    return null;
  }
  const copy = getPwaCopy(locale);

  return (
    <Button
      type="button"
      variant="secondary"
      size="sm"
      onClick={onInstall}
      disabled={isInstalling}
      aria-label={copy.installAriaLabel}
      className="h-9"
    >
      <Download className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
      <span className="hidden sm:inline">
        {isInstalling ? copy.installingLabel : copy.installLabel}
      </span>
    </Button>
  );
}

interface PwaStatusBannerProps {
  status: Pick<PwaStatus, "isOffline" | "updateAvailable" | "applyUpdate" | "dismissUpdate">;
}

export function PwaStatusBanner({ status }: PwaStatusBannerProps) {
  const locale = useUIStore((state) => state.locale);
  if (!status.isOffline && !status.updateAvailable) {
    return null;
  }
  const copy = getPwaCopy(locale);

  return (
    <div className="mb-4 flex flex-col gap-2">
      {status.isOffline ? (
        <div
          role="status"
          className="flex flex-wrap items-center gap-2.5 rounded-lg border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-3.5 py-2.5 text-xs text-[var(--status-warn)]"
        >
          <WifiOff className="h-4 w-4 shrink-0" aria-hidden="true" />
          <div className="min-w-0 flex-1">
            <span className="font-semibold">{copy.offlineTitle}</span>
            <span className="mx-1.5 opacity-60">·</span>
            <span>{copy.offlineDescription}</span>
          </div>
        </div>
      ) : null}

      {status.updateAvailable ? (
        <div
          role="status"
          className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] px-3.5 py-2.5 text-xs text-[var(--text-primary)]"
        >
          <div className="flex min-w-0 items-center gap-2.5">
            <RefreshCw
              className="h-4 w-4 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
            <div className="min-w-0">
              <span className="font-semibold">{copy.updateTitle}</span>
              <span className="mx-1.5 text-[var(--text-muted)]">·</span>
              <span className="text-[var(--text-secondary)]">{copy.updateDescription}</span>
            </div>
          </div>
          <div className="flex shrink-0 items-center gap-2">
            <Button type="button" variant="secondary" size="sm" onClick={status.dismissUpdate}>
              {copy.dismissAction}
            </Button>
            <Button type="button" variant="default" size="sm" onClick={status.applyUpdate}>
              {copy.updateAction}
            </Button>
          </div>
        </div>
      ) : null}
    </div>
  );
}
