import { AlertTriangle, ArrowUpCircle, ExternalLink, Loader2, RefreshCw, X } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import { useUIStore } from "@/shared/store/uiStore";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { useForceCheckSystemUpdate, useSystemUpdateStatus, useTriggerSystemUpdate } from "./hooks";
import {
  getUnavailableReasonMessage,
  getUpdateCopy,
  getUpdateStateLabel,
  isUpdateInProgress,
  shouldShowUpdateButton,
} from "./updateCopy";

const RESTART_POLL_TIMEOUT_MS = 60_000;

export function SystemUpdateControl() {
  const locale = useUIStore((state) => state.locale);
  const copy = getUpdateCopy(locale);

  const updateQuery = useSystemUpdateStatus();
  const forceCheckMutation = useForceCheckSystemUpdate();
  const triggerMutation = useTriggerSystemUpdate();

  const [modalOpen, setModalOpen] = useState(false);
  const [waitingForRestart, setWaitingForRestart] = useState(false);
  const [restartTimedOut, setRestartTimedOut] = useState(false);
  const initialCommitRef = useRef<string>("");

  const status = updateQuery.data;
  const inProgress = (isUpdateInProgress(status?.state) || waitingForRestart) && !restartTimedOut;
  const visible = shouldShowUpdateButton(status) || waitingForRestart || restartTimedOut;

  const closeModal = useCallback(() => {
    if (inProgress || triggerMutation.isPending) {
      return;
    }
    setModalOpen(false);
  }, [inProgress, triggerMutation.isPending]);

  useEffect(() => {
    if (!modalOpen) {
      return;
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.preventDefault();
        closeModal();
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [closeModal, modalOpen]);

  useEffect(() => {
    if (status?.state === "restarting" && !waitingForRestart) {
      setWaitingForRestart(true);
    }
  }, [status?.state, waitingForRestart]);

  useEffect(() => {
    if (!waitingForRestart) {
      return;
    }
    let cancelled = false;
    let sawDowntime = false;
    const startedAt = Date.now();
    const targetShort = status?.latestShortCommit?.toLowerCase() ?? "";
    const previousCommit = initialCommitRef.current.toLowerCase();

    const timer = window.setInterval(async () => {
      if (cancelled) {
        return;
      }
      if (Date.now() - startedAt >= RESTART_POLL_TIMEOUT_MS) {
        window.clearInterval(timer);
        setWaitingForRestart(false);
        setRestartTimedOut(true);
        return;
      }
      try {
        const { data } = await apiClient.GET("/api/version");
        if (!data || cancelled) {
          return;
        }
        const liveCommit = data.commit.toLowerCase();
        const commitChanged =
          (targetShort.length > 0 && liveCommit.startsWith(targetShort)) ||
          (previousCommit.length > 0 && liveCommit !== previousCommit);
        if (commitChanged || sawDowntime) {
          window.clearInterval(timer);
          window.location.reload();
        }
      } catch {
        sawDowntime = true;
      }
    }, 1_500);

    return () => {
      cancelled = true;
      window.clearInterval(timer);
    };
  }, [status?.latestShortCommit, waitingForRestart]);

  if (!visible || !status) {
    return null;
  }

  const modeLabel =
    status.mode === "source"
      ? copy.modeSource
      : status.mode === "binary"
        ? copy.modeBinary
        : copy.modeNone;

  const handleOpenModal = () => {
    if (triggerMutation.error) {
      triggerMutation.reset();
    }
    setModalOpen(true);
  };

  const handleConfirmUpdate = () => {
    initialCommitRef.current = status.currentCommit;
    setRestartTimedOut(false);
    triggerMutation.mutate({ mode: "auto" });
  };

  const activeError = restartTimedOut
    ? copy.restartTimeoutNotice
    : triggerMutation.error
      ? extractErrorMessage(triggerMutation.error)
      : status.state === "failed" && status.error
        ? status.error
        : "";

  return (
    <>
      <Button
        type="button"
        variant="default"
        size="sm"
        onClick={handleOpenModal}
        aria-label={copy.buttonAriaLabel}
        title={`${copy.modalTitle} (${status.latestShortCommit})`}
        className="relative h-9 w-9 px-0 sm:w-auto sm:px-3.5"
      >
        {inProgress ? (
          <Loader2 className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
        ) : (
          <ArrowUpCircle className="h-3.5 w-3.5" aria-hidden="true" />
        )}
        <span className="hidden sm:inline">
          {inProgress ? copy.updatingLabel : copy.buttonLabel}
        </span>
      </Button>

      {modalOpen ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="system-update-title"
          aria-describedby="system-update-description"
        >
          <button
            type="button"
            className="fixed inset-0 bg-black/55"
            aria-label={copy.closeAction}
            disabled={inProgress || triggerMutation.isPending}
            onClick={closeModal}
          />
          <div className="relative z-10 w-full max-w-md overflow-hidden rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
            <div className="flex items-center justify-between border-b border-[var(--border-subtle)] px-5 py-4">
              <div className="flex items-center gap-2.5">
                <ArrowUpCircle
                  className="h-4 w-4 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
                <h2
                  id="system-update-title"
                  className="text-sm font-semibold text-[var(--text-primary)]"
                >
                  {copy.modalTitle}
                </h2>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={closeModal}
                disabled={inProgress || triggerMutation.isPending}
                aria-label={copy.closeAction}
              >
                <X className="h-4 w-4" aria-hidden="true" />
              </Button>
            </div>

            <div className="space-y-4 px-5 py-4">
              <p
                id="system-update-description"
                className="text-xs leading-relaxed text-[var(--text-secondary)]"
              >
                {copy.modalDescription}
              </p>

              <div className="space-y-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3 text-xs">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-[var(--text-muted)]">{copy.branchLabel}</span>
                  <span className="font-mono font-medium text-[var(--text-primary)]">
                    {status.branch}
                  </span>
                </div>
                <div className="flex items-center justify-between gap-2">
                  <span className="text-[var(--text-muted)]">
                    {copy.currentCommitLabel} → {copy.latestCommitLabel}
                  </span>
                  <span className="font-mono text-[var(--text-primary)]">
                    {status.currentCommit || "—"} →{" "}
                    <strong className="text-[var(--accent-text)]">
                      {status.latestShortCommit || "—"}
                    </strong>
                  </span>
                </div>
                <div className="flex items-center justify-between gap-2">
                  <span className="text-[var(--text-muted)]">{copy.modeLabel}</span>
                  <Badge variant="outline">{modeLabel}</Badge>
                </div>
                {status.latestMessage ? (
                  <div className="border-t border-[var(--border-subtle)] pt-2">
                    <div className="mb-1 flex items-center justify-between gap-2">
                      <span className="text-[11px] text-[var(--text-muted)]">
                        {copy.latestMessageLabel}
                      </span>
                      {status.latestUrl ? (
                        <a
                          href={status.latestUrl}
                          target="_blank"
                          rel="noreferrer"
                          className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--accent-text)] hover:underline"
                        >
                          <span>{copy.viewCommitLabel}</span>
                          <ExternalLink className="h-3 w-3" aria-hidden="true" />
                        </a>
                      ) : null}
                    </div>
                    <p className="break-words font-mono text-[11px] text-[var(--text-primary)]">
                      {status.latestMessage}
                    </p>
                  </div>
                ) : null}
              </div>

              {inProgress ? (
                <div
                  role="status"
                  className="flex items-start gap-2.5 rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-3 py-2.5 text-xs text-[var(--text-primary)]"
                >
                  <Loader2
                    className="mt-0.5 h-4 w-4 shrink-0 animate-spin text-[var(--accent-primary)]"
                    aria-hidden="true"
                  />
                  <div className="min-w-0 space-y-1">
                    <p className="font-medium">
                      {waitingForRestart
                        ? getUpdateStateLabel(locale, "restarting")
                        : getUpdateStateLabel(locale, status.state)}
                    </p>
                    {waitingForRestart ? (
                      <p className="text-[var(--text-secondary)]">{copy.restartingNotice}</p>
                    ) : status.message ? (
                      <p className="text-[var(--text-secondary)]">{status.message}</p>
                    ) : null}
                  </div>
                </div>
              ) : null}

              {!inProgress && !status.canApply && status.unavailableReason ? (
                <div
                  role="status"
                  className="flex items-start gap-2.5 rounded-md border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-3 py-2.5 text-xs text-[var(--status-warn)]"
                >
                  <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
                  <p className="leading-relaxed">
                    {getUnavailableReasonMessage(locale, status.unavailableReason)}
                  </p>
                </div>
              ) : null}

              {activeError ? (
                <p
                  role="alert"
                  className="rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-3 py-2 text-xs text-[var(--status-error)]"
                >
                  {activeError}
                </p>
              ) : null}
            </div>

            <div className="flex items-center justify-end gap-2 border-t border-[var(--border-subtle)] px-5 py-3.5">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={closeModal}
                disabled={inProgress || triggerMutation.isPending}
              >
                {copy.cancelAction}
              </Button>
              {restartTimedOut ? (
                <Button
                  type="button"
                  variant="default"
                  size="sm"
                  onClick={() => window.location.reload()}
                >
                  <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{copy.reloadAction}</span>
                </Button>
              ) : !status.canApply && !inProgress ? (
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => forceCheckMutation.mutate()}
                  disabled={forceCheckMutation.isPending}
                >
                  <RefreshCw
                    className={
                      forceCheckMutation.isPending ? "h-3.5 w-3.5 animate-spin" : "h-3.5 w-3.5"
                    }
                    aria-hidden="true"
                  />
                  <span>
                    {forceCheckMutation.isPending ? copy.recheckingLabel : copy.recheckAction}
                  </span>
                </Button>
              ) : (
                <Button
                  type="button"
                  variant="default"
                  size="sm"
                  onClick={handleConfirmUpdate}
                  disabled={inProgress || triggerMutation.isPending || !status.canApply}
                >
                  {inProgress || triggerMutation.isPending ? (
                    <Loader2 className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                  ) : (
                    <ArrowUpCircle className="h-3.5 w-3.5" aria-hidden="true" />
                  )}
                  <span>{status.state === "failed" ? copy.retryAction : copy.confirmAction}</span>
                </Button>
              )}
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
