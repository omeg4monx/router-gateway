import { Check, RefreshCw, RotateCcw, Shuffle, X } from "lucide-react";
import { useEffect, useId, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { Label } from "@/shared/ui/label";
import { useProviderStrategy, useUpdateProviderStrategyMutation } from "./hooks";
import {
  DEFAULT_PROVIDER_STRATEGY,
  getProviderStrategyMeta,
  isProviderAccountStrategy,
  PROVIDER_STRATEGY_GROUPS,
  type ProviderAccountStrategy,
  QUICK_PROVIDER_STRATEGIES,
} from "./providerStrategy";

interface ProviderStrategyCardProps {
  providerId: string;
  enabledConnectionCount: number;
  isId: boolean;
}

export function ProviderStrategyCard({
  providerId,
  enabledConnectionCount,
  isId,
}: ProviderStrategyCardProps) {
  const selectId = useId();
  const strategyQuery = useProviderStrategy(providerId);
  const updateStrategyMutation = useUpdateProviderStrategyMutation();

  const [saveError, setSaveError] = useState<string | null>(null);
  const [savedNotice, setSavedNotice] = useState<string | null>(null);

  useEffect(() => {
    if (!savedNotice) return;
    const timer = window.setTimeout(() => setSavedNotice(null), 2400);
    return () => window.clearTimeout(timer);
  }, [savedNotice]);

  const activeStrategy: ProviderAccountStrategy =
    strategyQuery.data?.strategy ?? DEFAULT_PROVIDER_STRATEGY;
  const isDefault = strategyQuery.data?.is_default ?? activeStrategy === DEFAULT_PROVIDER_STRATEGY;
  const activeMeta = getProviderStrategyMeta(activeStrategy);
  const isSaving = updateStrategyMutation.isPending;

  const handleSelectStrategy = (nextStrategy: ProviderAccountStrategy) => {
    if (isSaving || nextStrategy === activeStrategy) {
      return;
    }
    setSaveError(null);
    updateStrategyMutation.mutate(
      { providerId, strategy: nextStrategy },
      {
        onSuccess: (updated) => {
          const meta = getProviderStrategyMeta(updated.strategy);
          setSavedNotice(
            isId
              ? `Strategi aktif diubah ke ${meta.labelId}.`
              : `Active strategy updated to ${meta.labelEn}.`,
          );
        },
        onError: (err) => {
          setSaveError(err.message);
        },
      },
    );
  };

  return (
    <Card>
      <CardHeader className="flex-col items-start gap-2.5 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0 space-y-0.5">
          <div className="flex flex-wrap items-center gap-2">
            <Shuffle className="h-4 w-4 shrink-0 text-[var(--color-primary)]" aria-hidden="true" />
            <CardTitle className="text-base">
              {isId ? "Strategi Routing Akun" : "Account Routing Strategy"}
            </CardTitle>
            <Badge variant="success">● {isId ? activeMeta.labelId : activeMeta.labelEn}</Badge>
            <Badge variant="secondary">
              {isDefault ? (isId ? "Default Aktif" : "Active Default") : isId ? "Kustom" : "Custom"}
            </Badge>
          </div>
          <CardDescription>
            {isId
              ? "Menentukan cara memilih koneksi/akun saat request diarahkan langsung ke provider ini."
              : "Controls how requests are distributed across connections when routing directly to this provider."}
          </CardDescription>
        </div>

        {activeStrategy !== DEFAULT_PROVIDER_STRATEGY ? (
          <Button
            type="button"
            size="xs"
            variant="secondary"
            className="w-full shrink-0 sm:w-auto"
            disabled={isSaving}
            onClick={() => handleSelectStrategy(DEFAULT_PROVIDER_STRATEGY)}
          >
            <RotateCcw className="h-3.5 w-3.5" aria-hidden="true" />
            <span>{isId ? "Reset ke Default" : "Reset to Default"}</span>
          </Button>
        ) : null}
      </CardHeader>

      <CardContent className="space-y-3.5">
        {saveError ? (
          <Alert variant="destructive">
            <div className="flex items-center justify-between gap-2">
              <AlertDescription className="break-words">{saveError}</AlertDescription>
              <button
                type="button"
                onClick={() => setSaveError(null)}
                aria-label={isId ? "Tutup pesan error" : "Dismiss error"}
                className="cursor-pointer rounded p-1 hover:bg-black/10"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          </Alert>
        ) : null}

        {strategyQuery.error ? (
          <Alert variant="warning">
            <AlertDescription className="break-words">
              {isId
                ? `Gagal memuat strategi tersimpan (${strategyQuery.error.message}). Menggunakan default aktif: ${activeMeta.labelId}.`
                : `Could not load saved strategy (${strategyQuery.error.message}). Using active default: ${activeMeta.labelEn}.`}
            </AlertDescription>
          </Alert>
        ) : null}

        {/* Quick-select pills for common account strategies */}
        <div className="flex flex-wrap items-center gap-1.5">
          {QUICK_PROVIDER_STRATEGIES.map((strategyId) => {
            const meta = getProviderStrategyMeta(strategyId);
            const selected = activeStrategy === strategyId;
            return (
              <button
                key={strategyId}
                type="button"
                disabled={isSaving}
                aria-pressed={selected}
                onClick={() => handleSelectStrategy(strategyId)}
                className={cn(
                  "inline-flex cursor-pointer items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-xs font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50",
                  selected
                    ? "border-[var(--color-primary)] bg-[var(--color-primary)]/12 text-[var(--color-primary)]"
                    : "border-[var(--color-border)] bg-[var(--color-surface)] text-[var(--color-text-muted)] hover:border-[var(--border-strong)] hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]",
                )}
              >
                {selected ? <Check className="h-3 w-3 shrink-0" aria-hidden="true" /> : null}
                <span>{isId ? meta.labelId : meta.labelEn}</span>
              </button>
            );
          })}
        </div>

        {/* Full categorized selector + live active description */}
        <div className="grid grid-cols-1 gap-3 md:grid-cols-[minmax(0,280px)_1fr] md:items-start">
          <div className="space-y-1.5">
            <Label htmlFor={selectId} className="text-xs text-[var(--color-text-muted)]">
              {isId ? "Semua Strategi Akun" : "All Account Strategies"}
            </Label>
            <div className="relative">
              <select
                id={selectId}
                value={activeStrategy}
                disabled={isSaving}
                onChange={(event) => {
                  const next = event.target.value;
                  if (isProviderAccountStrategy(next)) {
                    handleSelectStrategy(next);
                  }
                }}
                className="flex h-9 w-full cursor-pointer rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-1.5 text-xs font-medium text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)] disabled:cursor-not-allowed disabled:opacity-60"
              >
                {PROVIDER_STRATEGY_GROUPS.map((group) => (
                  <optgroup key={group.id} label={isId ? group.labelId : group.labelEn}>
                    {group.options.map((opt) => (
                      <option key={opt.value} value={opt.value}>
                        {isId ? opt.labelId : opt.labelEn}
                        {opt.value === DEFAULT_PROVIDER_STRATEGY
                          ? isId
                            ? " (Default)"
                            : " (Default)"
                          : ""}
                      </option>
                    ))}
                  </optgroup>
                ))}
              </select>
            </div>
          </div>

          <div className="rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)]/50 px-3 py-2.5 text-xs">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="font-semibold text-[var(--color-text-main)]">
                {isId ? activeMeta.labelId : activeMeta.labelEn}
                <span className="ml-1.5 font-mono text-[10.5px] font-normal text-[var(--color-text-muted)]">
                  ({activeMeta.value})
                </span>
              </span>
              {isSaving ? (
                <span className="inline-flex items-center gap-1 text-[11px] text-[var(--color-text-muted)]">
                  <RefreshCw className="h-3 w-3 animate-spin" aria-hidden="true" />
                  <span>{isId ? "Menyimpan..." : "Saving..."}</span>
                </span>
              ) : savedNotice ? (
                <span className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--status-ok)]">
                  <Check className="h-3 w-3" aria-hidden="true" />
                  <span>{savedNotice}</span>
                </span>
              ) : null}
            </div>
            <p className="mt-1 leading-relaxed text-[var(--color-text-muted)]">
              {isId ? activeMeta.descId : activeMeta.descEn}
            </p>
            {enabledConnectionCount < 2 ? (
              <p className="mt-1.5 text-[11px] text-[var(--color-text-subtle)]">
                {isId
                  ? "Strategi sudah aktif dan otomatis berlaku penuh saat terdapat 2 atau lebih koneksi aktif."
                  : "Strategy is active now and automatically distributes traffic once 2 or more connections are enabled."}
              </p>
            ) : null}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
