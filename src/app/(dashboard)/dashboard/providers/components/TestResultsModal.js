"use client";

import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Button, Modal } from "@/shared/components";
import { cn } from "@/shared/utils/cn";

const MODE_LABELS = {
  oauth: "OAuth",
  free: "Free OAuth",
  apikey: "API Key",
  compatible: "Custom Endpoints",
  provider: "Provider",
  all: "All Active",
};

export default function TestResultsModal({
  results,
  onClose,
  onRerun,
  rerunning = false,
}) {
  const [failedOnly, setFailedOnly] = useState(false);

  useEffect(() => {
    setFailedOnly(false);
  }, [results]);

  const items = results?.results || [];
  const summary = results?.summary;
  const modeLabel = MODE_LABELS[results?.mode] || results?.mode || "Batch";
  const failedCount =
    summary?.failed ?? items.filter((r) => !r.valid).length;
  const visibleItems = failedOnly ? items.filter((r) => !r.valid) : items;

  const footer = (
    <div className="flex items-center justify-between gap-2">
      <div className="text-[11px] text-text-muted">
        {results?.testedAt
          ? `Completed ${new Date(results.testedAt).toLocaleTimeString()}`
          : null}
      </div>
      <div className="flex items-center gap-2">
        {onRerun && (
          <Button
            size="sm"
            variant="secondary"
            icon="refresh"
            onClick={onRerun}
            loading={rerunning}
            disabled={rerunning}
          >
            Re-run Test
          </Button>
        )}
        <Button size="sm" variant="primary" onClick={onClose}>
          Close
        </Button>
      </div>
    </div>
  );

  return (
    <Modal
      isOpen={Boolean(results)}
      onClose={onClose}
      title="Connection Test Results"
      size="lg"
      footer={footer}
    >
      {results?.error && !results?.results ? (
        <div className="flex flex-col items-center justify-center py-6 text-center gap-2">
          <span className="material-symbols-outlined text-danger text-[32px]">
            error
          </span>
          <p className="text-sm font-medium text-danger">{results.error}</p>
        </div>
      ) : (
        <div className="flex min-w-0 flex-col gap-3">
          {summary && (
            <div className="flex flex-wrap items-center justify-between gap-2 pb-2.5 border-b border-border-subtle">
              <div className="flex flex-wrap items-center gap-2 text-xs">
                <span className="font-medium text-text-main">
                  {modeLabel} Test
                </span>
                <span className="px-2 py-0.5 rounded-md bg-success-bg text-success font-medium font-mono tabular-nums">
                  {summary.passed} passed
                </span>
                {summary.failed > 0 && (
                  <span className="px-2 py-0.5 rounded-md bg-danger-bg text-danger font-medium font-mono tabular-nums">
                    {summary.failed} failed
                  </span>
                )}
                <span className="text-text-muted font-mono tabular-nums">
                  ({summary.total} tested)
                </span>
              </div>

              {failedCount > 0 && (
                <div className="inline-flex items-center gap-1 rounded-lg bg-surface-2 p-0.5 border border-border-subtle text-[11px]">
                  <button
                    type="button"
                    onClick={() => setFailedOnly(false)}
                    className={cn(
                      "px-2 py-0.5 rounded-md font-medium transition-colors cursor-pointer",
                      !failedOnly
                        ? "bg-surface text-text-main shadow-xs"
                        : "text-text-muted hover:text-text-main",
                    )}
                  >
                    All ({items.length})
                  </button>
                  <button
                    type="button"
                    onClick={() => setFailedOnly(true)}
                    className={cn(
                      "px-2 py-0.5 rounded-md font-medium transition-colors cursor-pointer",
                      failedOnly
                        ? "bg-danger/15 text-danger"
                        : "text-text-muted hover:text-text-main",
                    )}
                  >
                    Failed ({failedCount})
                  </button>
                </div>
              )}
            </div>
          )}

          {visibleItems.length > 0 ? (
            <div className="flex flex-col gap-2">
              {visibleItems.map((r, i) => {
                const errorDetail =
                  !r.valid && (r.error || r.diagnosis?.message);
                return (
                  <div
                    key={r.connectionId || i}
                    className={cn(
                      "flex flex-col gap-1.5 rounded-xl border px-3 py-2.5 text-xs",
                      r.valid
                        ? "border-border bg-surface-2/60"
                        : "border-danger/30 bg-danger/5",
                    )}
                  >
                    <div className="flex min-w-0 flex-wrap items-center gap-2 sm:flex-nowrap">
                      <span
                        className={cn(
                          "material-symbols-outlined text-[16px] shrink-0",
                          r.valid ? "text-success" : "text-danger",
                        )}
                      >
                        {r.valid ? "check_circle" : "error"}
                      </span>
                      <div className="min-w-0 flex-1">
                        <span className="font-medium text-text-main truncate block sm:inline">
                          {r.connectionName}
                        </span>
                        <span className="text-text-muted font-mono text-[11px] sm:ml-1.5 block sm:inline truncate">
                          ({r.provider})
                        </span>
                      </div>
                      {r.latencyMs !== undefined && (
                        <span className="shrink-0 text-text-muted font-mono tabular-nums text-[11px]">
                          {r.latencyMs}ms
                        </span>
                      )}
                      <span
                        className={cn(
                          "shrink-0 text-[10px] uppercase font-bold px-1.5 py-0.5 rounded-md",
                          r.valid
                            ? "bg-success-bg text-success"
                            : "bg-danger-bg text-danger",
                        )}
                      >
                        {r.valid ? "OK" : r.diagnosis?.type || "ERROR"}
                      </span>
                    </div>
                    {errorDetail && (
                      <p className="pl-6 font-mono text-[11px] text-danger break-words">
                        {errorDetail}
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="text-center py-6 text-text-muted text-xs">
              {items.length === 0
                ? "No active connections found for this group."
                : "No failed connections in this run."}
            </div>
          )}
        </div>
      )}
    </Modal>
  );
}

TestResultsModal.propTypes = {
  results: PropTypes.shape({
    mode: PropTypes.string,
    results: PropTypes.array,
    testedAt: PropTypes.string,
    summary: PropTypes.shape({
      total: PropTypes.number,
      passed: PropTypes.number,
      failed: PropTypes.number,
    }),
    error: PropTypes.string,
  }),
  onClose: PropTypes.func.isRequired,
  onRerun: PropTypes.func,
  rerunning: PropTypes.bool,
};
