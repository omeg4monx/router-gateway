import {
  Check,
  CheckCircle2,
  Copy,
  CopyPlus,
  FlaskConical,
  KeyRound,
  Pencil,
  Power,
  Sparkles,
  Trash2,
  XCircle,
} from "lucide-react";
import { memo } from "react";
import type { CatalogTestResult } from "@/features/models/types";
import type { ProviderConnectionSummary } from "@/features/providers/schemas";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { formatStepAccountLabel } from "./comboHelpers";
import type { Combo, ComboStep, ComboStrategy } from "./schemas";
import { getStrategyMeta } from "./strategies";

interface ComboCardProps {
  combo: Combo;
  connectionsById: ReadonlyMap<string, ProviderConnectionSummary>;
  testResult: CatalogTestResult | undefined;
  isTesting: boolean;
  isBusy: boolean;
  isCopied: boolean;
  isConfirmingDelete: boolean;
  isCurrentlyEditing: boolean;
  isId: boolean;
  onCopyName: (comboName: string) => void;
  onOpenStrategyGuide: (strategy: ComboStrategy) => void;
  onTestCombo: (combo: Combo) => void;
  onEditCombo: (combo: Combo) => void;
  onDuplicateCombo: (combo: Combo) => void;
  onToggleEnabled: (combo: Combo) => void;
  onRequestDelete: (comboId: string) => void;
  onCancelDelete: () => void;
  onConfirmDelete: (combo: Combo) => void;
}

export const ComboCard = memo(function ComboCard({
  combo,
  connectionsById,
  testResult,
  isTesting,
  isBusy,
  isCopied,
  isConfirmingDelete,
  isCurrentlyEditing,
  isId,
  onCopyName,
  onOpenStrategyGuide,
  onTestCombo,
  onEditCombo,
  onDuplicateCombo,
  onToggleEnabled,
  onRequestDelete,
  onCancelDelete,
  onConfirmDelete,
}: ComboCardProps) {
  const stratMeta = getStrategyMeta(combo.strategy);
  const steps: ComboStep[] =
    combo.steps && combo.steps.length > 0
      ? combo.steps
      : combo.targets.map((t): ComboStep => ({ kind: "model", model: t }));

  return (
    <Card
      className={`overflow-hidden transition-colors ${
        isCurrentlyEditing ? "border-[var(--accent-primary)]" : ""
      }`}
    >
      <div className="space-y-2 p-3">
        {/* Top Row: Name + Badges + Compact Actions */}
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-1.5">
              <span
                className={`h-2 w-2 shrink-0 rounded-full ${
                  combo.enabled ? "bg-[var(--status-ok)]" : "bg-[var(--status-warn)]"
                }`}
                aria-hidden="true"
              />
              <span className="truncate font-mono text-sm font-semibold text-[var(--text-primary)]">
                {combo.name}
              </span>
              <button
                type="button"
                onClick={() => onCopyName(combo.name)}
                className="rounded p-0.5 text-[var(--text-muted)] hover:bg-[var(--bg-canvas)] hover:text-[var(--text-primary)]"
                aria-label={isId ? "Salin nama combo" : "Copy combo name"}
                title={isId ? "Salin nama combo" : "Copy combo name"}
              >
                {isCopied ? (
                  <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" />
                ) : (
                  <Copy className="h-3.5 w-3.5" />
                )}
              </button>
              <button
                type="button"
                onClick={() => onOpenStrategyGuide(combo.strategy)}
                title={
                  isId
                    ? `${stratMeta.labelId} — klik untuk panduan`
                    : `${stratMeta.labelEn} — click for guide`
                }
                className="inline-flex items-center gap-1 rounded-full border border-[var(--border-strong)] bg-[var(--bg-canvas)] px-2 py-0.5 font-mono text-[10px] text-[var(--text-primary)] transition-colors hover:border-[var(--accent-primary)] hover:text-[var(--accent-text)]"
              >
                <span>{combo.strategy}</span>
              </button>
            </div>

            {combo.description ? (
              <p className="mt-0.5 truncate text-xs text-[var(--text-secondary)]">
                {combo.description}
              </p>
            ) : null}
          </div>

          {/* Compact Actions */}
          <div className="flex shrink-0 flex-wrap items-center gap-1">
            {isConfirmingDelete ? (
              <>
                <Button
                  type="button"
                  variant="destructive"
                  size="xs"
                  disabled={isBusy}
                  onClick={() => onConfirmDelete(combo)}
                >
                  {isBusy ? (isId ? "Menghapus..." : "Deleting...") : isId ? "Hapus" : "Confirm"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  size="xs"
                  disabled={isBusy}
                  onClick={onCancelDelete}
                >
                  {isId ? "Batal" : "Cancel"}
                </Button>
              </>
            ) : (
              <>
                <Button
                  type="button"
                  variant="outline"
                  size="xs"
                  disabled={isTesting}
                  onClick={() => onTestCombo(combo)}
                  title={isId ? "Uji combo" : "Test combo"}
                >
                  <FlaskConical className="h-3 w-3" />
                  <span>
                    {isTesting ? (isId ? "Menguji..." : "Testing...") : isId ? "Uji" : "Test"}
                  </span>
                </Button>
                <Button
                  type="button"
                  variant={isCurrentlyEditing ? "default" : "outline"}
                  size="xs"
                  disabled={isBusy}
                  onClick={() => onEditCombo(combo)}
                  title="Edit"
                >
                  <Pencil className="h-3 w-3" />
                  <span className="hidden sm:inline">Edit</span>
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="xs"
                  disabled={isBusy}
                  onClick={() => onDuplicateCombo(combo)}
                  aria-label={isId ? "Duplikat ke form" : "Duplicate into form"}
                  title={isId ? "Duplikat ke form" : "Duplicate into form"}
                >
                  <CopyPlus className="h-3.5 w-3.5" />
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  size="xs"
                  disabled={isBusy}
                  onClick={() => onToggleEnabled(combo)}
                  aria-label={
                    combo.enabled
                      ? isId
                        ? "Jeda combo"
                        : "Pause combo"
                      : isId
                        ? "Aktifkan combo"
                        : "Enable combo"
                  }
                  title={
                    combo.enabled
                      ? isId
                        ? "Jeda combo"
                        : "Pause combo"
                      : isId
                        ? "Aktifkan combo"
                        : "Enable combo"
                  }
                >
                  <Power className="h-3 w-3" />
                </Button>
                <Button
                  type="button"
                  variant="destructive"
                  size="xs"
                  disabled={isBusy}
                  onClick={() => onRequestDelete(combo.id)}
                  aria-label={isId ? "Hapus combo" : "Delete combo"}
                  title={isId ? "Hapus combo" : "Delete combo"}
                >
                  <Trash2 className="h-3 w-3" />
                </Button>
              </>
            )}
          </div>
        </div>

        {/* Compact Steps Chain */}
        <div className="flex flex-wrap items-center gap-1">
          {steps.map((step, idx) => {
            const isRef = step.kind === "combo-ref";
            const ident = isRef ? (step.comboName ?? "") : (step.model ?? "");
            const isPinned = Boolean(step.connectionId?.trim());
            const accountLabel = isRef
              ? null
              : formatStepAccountLabel(step.connectionId, connectionsById, isId);

            return (
              <div
                key={`${step.kind}-${ident}-${step.connectionId ?? "dynamic"}`}
                className="inline-flex max-w-full items-center gap-1 rounded border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 py-0.5 text-[11px]"
              >
                <span className="font-mono text-[10px] text-[var(--text-muted)]">#{idx + 1}</span>
                {isRef ? (
                  <Badge variant="default" className="px-1 py-0 text-[9px]">
                    ref
                  </Badge>
                ) : null}
                <span className="truncate font-mono font-medium text-[var(--text-primary)]">
                  {ident}
                </span>
                {accountLabel ? (
                  <span
                    className={`inline-flex shrink-0 items-center gap-0.5 rounded px-1 py-0.25 text-[10px] ${
                      isPinned
                        ? "bg-[var(--bg-surface-raised)] text-[var(--accent-text)]"
                        : "bg-[var(--bg-surface)] text-[var(--text-muted)]"
                    }`}
                  >
                    {isPinned ? (
                      <KeyRound className="h-2.5 w-2.5" aria-hidden="true" />
                    ) : (
                      <Sparkles className="h-2.5 w-2.5" aria-hidden="true" />
                    )}
                    <span>{isPinned ? accountLabel : "Dynamic"}</span>
                  </span>
                ) : null}
                {step.weight && step.weight > 0 ? (
                  <span className="font-mono text-[10px] font-semibold text-[var(--accent-text)]">
                    {step.weight}%
                  </span>
                ) : null}
              </div>
            );
          })}
        </div>

        {/* Inline Live Test Result */}
        {testResult ? (
          <div className="flex flex-wrap items-center gap-2 border-t border-[var(--border-subtle)] pt-1.5 text-xs">
            {testResult.status === "ok" || testResult.status === "slow" ? (
              <Badge variant={testResult.status === "ok" ? "success" : "warning"} className="gap-1">
                <CheckCircle2 className="h-3 w-3" />
                <span>
                  {testResult.status.toUpperCase()}
                  {typeof testResult.latencyMs === "number" ? ` · ${testResult.latencyMs}ms` : ""}
                </span>
              </Badge>
            ) : (
              <Badge variant="destructive" className="gap-1">
                <XCircle className="h-3 w-3" />
                <span>ERROR</span>
              </Badge>
            )}
            {testResult.resolvedBy ? (
              <span className="font-mono text-[11px] text-[var(--text-secondary)]">
                {testResult.resolvedBy}
              </span>
            ) : null}
            {testResult.error ? (
              <span className="font-mono text-[11px] text-[var(--status-error)]">
                {testResult.error}
              </span>
            ) : null}
          </div>
        ) : null}
      </div>
    </Card>
  );
});
