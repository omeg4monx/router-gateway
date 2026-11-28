import { ArrowDown, ArrowUp, KeyRound, Layers, Scale, Sparkles, X } from "lucide-react";
import { memo } from "react";
import type { ProviderConnectionSummary } from "@/features/providers/schemas";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { formatStepAccountLabel } from "./comboHelpers";
import type { ComboStep } from "./schemas";

interface ComboStepListProps {
  steps: readonly ComboStep[];
  connectionsById: ReadonlyMap<string, ProviderConnectionSummary>;
  usesWeights: boolean;
  totalWeight: number;
  isId: boolean;
  onAutoBalanceWeights: () => void;
  onMoveStep: (index: number, delta: -1 | 1) => void;
  onUpdateStep: (index: number, patch: Partial<ComboStep>) => void;
  onRemoveStep: (index: number) => void;
}

export const ComboStepList = memo(function ComboStepList({
  steps,
  connectionsById,
  usesWeights,
  totalWeight,
  isId,
  onAutoBalanceWeights,
  onMoveStep,
  onUpdateStep,
  onRemoveStep,
}: ComboStepListProps) {
  return (
    <div className="space-y-2 border-t border-[var(--border-subtle)] pt-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <span className="text-xs font-medium text-[var(--text-primary)]">
            {isId ? "Daftar Model Terpilih" : "Selected Models"}
          </span>
          <Badge variant="outline" className="text-[10px] tabular-nums">
            {steps.length}
          </Badge>
          {usesWeights && steps.length > 0 ? (
            <Badge
              variant={totalWeight === 100 ? "success" : "warning"}
              className="text-[10px] tabular-nums"
            >
              {isId ? `Bobot: ${totalWeight}%` : `Weight: ${totalWeight}%`}
            </Badge>
          ) : null}
        </div>

        {usesWeights && steps.length > 0 ? (
          <Button type="button" size="xs" variant="outline" onClick={onAutoBalanceWeights}>
            <Scale className="h-3 w-3" aria-hidden="true" />
            <span>{isId ? "Seimbangkan" : "Auto-Balance"}</span>
          </Button>
        ) : null}
      </div>

      {steps.length === 0 ? (
        <div className="rounded-md border border-dashed border-[var(--border-strong)] bg-[var(--bg-surface)]/50 px-3 py-5 text-center">
          <Layers className="mx-auto h-5 w-5 text-[var(--text-muted)]" aria-hidden="true" />
          <p className="mt-1.5 text-xs text-[var(--text-muted)]">
            {isId
              ? "Belum ada model ditambahkan. Pilih provider & model lalu klik Tambah Model."
              : "No models added yet. Pick a provider & model above and click Add Model."}
          </p>
        </div>
      ) : (
        <ol className="space-y-1.5">
          {steps.map((step, idx) => {
            const isRef = step.kind === "combo-ref";
            const ident = isRef ? (step.comboName ?? "") : (step.model ?? "");
            const accountLabel = isRef
              ? null
              : formatStepAccountLabel(step.connectionId, connectionsById, isId);
            const isPinned = Boolean(step.connectionId?.trim());

            return (
              <li
                key={`${step.kind}-${ident}-${step.connectionId ?? "dynamic"}`}
                className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-2.5"
              >
                <div className="flex items-center justify-between gap-2">
                  <div className="min-w-0 flex-1 space-y-1">
                    <div className="flex flex-wrap items-center gap-1.5">
                      <span className="font-mono text-[10px] font-semibold text-[var(--text-muted)]">
                        #{idx + 1}
                      </span>
                      {isRef ? (
                        <Badge variant="default" className="px-1.5 py-0 text-[9px]">
                          combo-ref
                        </Badge>
                      ) : null}
                      <span className="truncate font-mono text-xs font-semibold text-[var(--text-primary)]">
                        {ident}
                      </span>
                    </div>

                    {accountLabel ? (
                      <div className="flex items-center gap-1 text-[11px] text-[var(--text-secondary)]">
                        {isPinned ? (
                          <KeyRound
                            className="h-3 w-3 shrink-0 text-[var(--accent-primary)]"
                            aria-hidden="true"
                          />
                        ) : (
                          <Sparkles
                            className="h-3 w-3 shrink-0 text-[var(--status-ok)]"
                            aria-hidden="true"
                          />
                        )}
                        <span className="truncate">{accountLabel}</span>
                      </div>
                    ) : null}
                  </div>

                  <div className="flex shrink-0 items-center gap-1">
                    {usesWeights ? (
                      <div className="mr-1 flex items-center gap-1">
                        <Input
                          type="number"
                          min={0}
                          max={100}
                          value={step.weight ?? 0}
                          onChange={(e) =>
                            onUpdateStep(idx, {
                              weight: Math.max(0, Math.min(100, Number(e.target.value) || 0)),
                            })
                          }
                          className="h-7 w-14 px-1.5 text-right font-mono text-xs"
                          aria-label={isId ? "Bobot persen" : "Weight percentage"}
                        />
                        <span className="text-[11px] text-[var(--text-muted)]">%</span>
                      </div>
                    ) : null}

                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      disabled={idx === 0}
                      onClick={() => onMoveStep(idx, -1)}
                      aria-label={isId ? "Naikkan urutan" : "Move up"}
                    >
                      <ArrowUp className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      disabled={idx === steps.length - 1}
                      onClick={() => onMoveStep(idx, 1)}
                      aria-label={isId ? "Turunkan urutan" : "Move down"}
                    >
                      <ArrowDown className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      onClick={() => onRemoveStep(idx)}
                      aria-label={isId ? "Hapus langkah" : "Remove step"}
                    >
                      <X className="h-3.5 w-3.5 text-[var(--status-error)]" />
                    </Button>
                  </div>
                </div>

                {isRef ? (
                  <label className="mt-2 flex cursor-pointer items-center gap-1.5 border-t border-[var(--border-subtle)] pt-2 text-[11px] text-[var(--text-secondary)]">
                    <input
                      type="checkbox"
                      checked={Boolean(step.fallbackOnlyOnQuotaExhaustion)}
                      onChange={(e) =>
                        onUpdateStep(idx, {
                          fallbackOnlyOnQuotaExhaustion: e.target.checked,
                        })
                      }
                      className="h-3.5 w-3.5 accent-[var(--accent-primary)]"
                    />
                    <span>
                      {isId
                        ? "Failover hanya saat kuota habis"
                        : "Fallback only on quota exhaustion"}
                    </span>
                  </label>
                ) : null}
              </li>
            );
          })}
        </ol>
      )}
    </div>
  );
});
