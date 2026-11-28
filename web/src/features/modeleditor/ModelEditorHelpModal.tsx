import { ArrowDownRight, ArrowUpRight, BookOpen, X } from "lucide-react";
import { useEffect } from "react";
import { cn } from "@/shared/lib/utils";
import type { UILocale } from "@/shared/store/uiStore";
import { Button } from "@/shared/ui/button";
import {
  getParameterGuideEntry,
  MODEL_EDITOR_COPY,
  PARAMETER_GUIDE_ENTRIES,
  type ParameterGuideId,
} from "./parameterGuide";

interface ModelEditorHelpModalProps {
  open: boolean;
  activeTopic: ParameterGuideId;
  locale: UILocale;
  onSelectTopic: (topic: ParameterGuideId) => void;
  onClose: () => void;
}

export function ModelEditorHelpModal({
  open,
  activeTopic,
  locale,
  onSelectTopic,
  onClose,
}: ModelEditorHelpModalProps) {
  const copy = MODEL_EDITOR_COPY[locale];
  const entry = getParameterGuideEntry(activeTopic);

  useEffect(() => {
    if (!open) return;
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-end justify-center p-2 sm:items-center sm:p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="model-editor-guide-title"
    >
      <button
        type="button"
        className="fixed inset-0 bg-black/55"
        aria-label={copy.closeLabel}
        onClick={onClose}
      />

      <div className="relative z-10 flex max-h-[85dvh] w-full max-w-xl flex-col overflow-hidden rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
        <div className="flex shrink-0 items-center justify-between gap-3 border-b border-[var(--border-subtle)] px-4 py-3">
          <div className="flex min-w-0 items-center gap-2">
            <BookOpen
              className="h-4 w-4 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
            <div className="min-w-0">
              <h2
                id="model-editor-guide-title"
                className="truncate text-sm font-semibold text-[var(--text-primary)]"
              >
                {copy.guideTitle}
              </h2>
              <p className="truncate text-xs text-[var(--text-muted)]">{copy.guideSubtitle}</p>
            </div>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={copy.closeLabel}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        <div
          className="flex shrink-0 flex-wrap gap-1.5 border-b border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]/60 px-4 py-2.5"
          role="tablist"
          aria-label={copy.guideTitle}
        >
          {PARAMETER_GUIDE_ENTRIES.map((item) => {
            const isSelected = item.id === entry.id;
            return (
              <button
                key={item.id}
                type="button"
                role="tab"
                aria-selected={isSelected}
                onClick={() => onSelectTopic(item.id)}
                className={cn(
                  "cursor-pointer rounded-md px-2.5 py-1 text-xs font-medium transition-colors",
                  isSelected
                    ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                    : "border border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]",
                )}
              >
                {item.shortLabel[locale]}
              </button>
            );
          })}
        </div>

        <div className="min-h-0 flex-1 space-y-3 overflow-y-auto p-4 text-xs">
          <div>
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">
              {entry.title[locale]}
            </h3>
            <p className="mt-1 leading-relaxed text-[var(--text-secondary)]">
              {entry.purpose[locale]}
            </p>
          </div>

          <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2">
            <p className="font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
              {copy.guideDefault}
            </p>
            <p className="mt-0.5 text-[var(--text-primary)]">{entry.defaultValue[locale]}</p>
          </div>

          <div className="grid gap-2.5 sm:grid-cols-2">
            <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3">
              <div className="flex items-center gap-1.5 font-semibold text-[var(--accent-text)]">
                <ArrowUpRight className="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                <span>{entry.whenIncreased.label[locale]}</span>
              </div>
              <p className="mt-1 leading-relaxed text-[var(--text-secondary)]">
                {entry.whenIncreased.effect[locale]}
              </p>
            </div>

            <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3">
              <div className="flex items-center gap-1.5 font-semibold text-[var(--status-ok)]">
                <ArrowDownRight className="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
                <span>{entry.whenDecreased.label[locale]}</span>
              </div>
              <p className="mt-1 leading-relaxed text-[var(--text-secondary)]">
                {entry.whenDecreased.effect[locale]}
              </p>
            </div>
          </div>

          {entry.tip ? (
            <div className="rounded-md border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-3 py-2.5">
              <p className="font-mono text-[10px] font-semibold uppercase tracking-wider text-[var(--status-warn)]">
                {copy.guideTip}
              </p>
              <p className="mt-1 leading-relaxed text-[var(--text-primary)]">{entry.tip[locale]}</p>
            </div>
          ) : (
            <div className="flex items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-2">
              <span className="font-mono text-[10px] uppercase tracking-wider text-[var(--text-muted)]">
                {copy.guideExample}
              </span>
              <code className="truncate font-mono text-xs text-[var(--text-primary)]">
                {entry.example}
              </code>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
