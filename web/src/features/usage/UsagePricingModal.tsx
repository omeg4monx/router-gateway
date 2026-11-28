import { AlertCircle, DollarSign, RotateCcw, Search, X } from "lucide-react";
import { type FormEvent, useEffect, useMemo, useState } from "react";
import { Alert, AlertDescription } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  useDeleteModelPricingMutation,
  useModelPricing,
  useUpsertModelPricingMutation,
} from "./hooks";
import type { ModelPricing } from "./schemas";

interface UsagePricingModalProps {
  isId: boolean;
  onClose: () => void;
}

export function UsagePricingModal({ isId, onClose }: UsagePricingModalProps) {
  const pricingQuery = useModelPricing(true);
  const upsertMutation = useUpsertModelPricingMutation();
  const deleteMutation = useDeleteModelPricingMutation();

  const [filterText, setFilterText] = useState("");
  const [modelInput, setModelInput] = useState("");
  const [inputRate, setInputRate] = useState("");
  const [outputRate, setOutputRate] = useState("");
  const [cachedRate, setCachedRate] = useState("0");
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        event.preventDefault();
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [onClose]);

  const pricingItems = pricingQuery.data?.pricing ?? [];
  const filteredPricing = useMemo(() => {
    const needle = filterText.trim().toLowerCase();
    if (!needle) {
      return pricingItems;
    }
    return pricingItems.filter((item) => item.model.toLowerCase().includes(needle));
  }, [pricingItems, filterText]);

  const handleSelectItem = (item: ModelPricing) => {
    setModelInput(item.model);
    setInputRate(String(item.input_per_million));
    setOutputRate(String(item.output_per_million));
    setCachedRate(String(item.cached_per_million));
    setFormError(null);
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    const model = modelInput.trim();
    if (!model || model === "*") {
      setFormError(
        isId ? "ID model wajib diisi dan bukan '*'." : "Model ID is required and cannot be '*'.",
      );
      return;
    }

    const parsedIn = Number(inputRate);
    const parsedOut = Number(outputRate);
    const parsedCached = cachedRate.trim() === "" ? 0 : Number(cachedRate);
    if (
      !Number.isFinite(parsedIn) ||
      parsedIn < 0 ||
      !Number.isFinite(parsedOut) ||
      parsedOut < 0 ||
      !Number.isFinite(parsedCached) ||
      parsedCached < 0
    ) {
      setFormError(
        isId
          ? "Tarif harus berupa angka non-negatif yang valid."
          : "Rates must be valid non-negative numbers.",
      );
      return;
    }

    setFormError(null);
    upsertMutation.mutate(
      {
        model,
        body: {
          input_per_million: parsedIn,
          output_per_million: parsedOut,
          cached_per_million: parsedCached,
        },
      },
      {
        onSuccess: () => {
          setModelInput("");
          setInputRate("");
          setOutputRate("");
          setCachedRate("0");
        },
      },
    );
  };

  const activeError = formError ?? upsertMutation.error?.message ?? deleteMutation.error?.message;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4"
      role="dialog"
      aria-modal="true"
      aria-labelledby="usage-pricing-modal-title"
    >
      <button
        type="button"
        className="fixed inset-0 bg-black/55"
        aria-label={isId ? "Tutup modal tarif" : "Close pricing modal"}
        onClick={onClose}
      />

      <div className="relative z-10 flex max-h-[85dvh] w-full max-w-2xl flex-col overflow-hidden rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
        <div className="flex shrink-0 items-center justify-between gap-2 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-3">
          <div className="flex min-w-0 items-center gap-2">
            <DollarSign
              className="h-4 w-4 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
            <h2
              id="usage-pricing-modal-title"
              className="truncate text-sm font-semibold text-[var(--text-primary)]"
            >
              {isId ? "Tarif Model ($ / 1M Token)" : "Model Pricing ($ / 1M Tokens)"}
            </h2>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={isId ? "Tutup" : "Close"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        {/* Override Form */}
        <form
          onSubmit={handleSubmit}
          className="shrink-0 space-y-2.5 border-b border-[var(--border-subtle)] px-4 py-3"
        >
          <div className="grid grid-cols-1 gap-2 sm:grid-cols-12">
            <div className="sm:col-span-5">
              <Label htmlFor="pricing-model-id" className="text-[11px]">
                Model
              </Label>
              <Input
                id="pricing-model-id"
                type="text"
                value={modelInput}
                onChange={(e) => setModelInput(e.target.value)}
                placeholder="gpt-4o / openai/gpt-4o"
                className="mt-1 h-8 font-mono text-xs"
              />
            </div>
            <div className="sm:col-span-2">
              <Label htmlFor="pricing-input-rate" className="text-[11px]">
                Input
              </Label>
              <Input
                id="pricing-input-rate"
                type="number"
                step="any"
                min="0"
                value={inputRate}
                onChange={(e) => setInputRate(e.target.value)}
                placeholder="2.5"
                className="mt-1 h-8 font-mono text-xs tabular-nums"
              />
            </div>
            <div className="sm:col-span-2">
              <Label htmlFor="pricing-output-rate" className="text-[11px]">
                Output
              </Label>
              <Input
                id="pricing-output-rate"
                type="number"
                step="any"
                min="0"
                value={outputRate}
                onChange={(e) => setOutputRate(e.target.value)}
                placeholder="10"
                className="mt-1 h-8 font-mono text-xs tabular-nums"
              />
            </div>
            <div className="sm:col-span-3">
              <Label htmlFor="pricing-cached-rate" className="text-[11px]">
                Cached
              </Label>
              <div className="mt-1 flex items-center gap-1.5">
                <Input
                  id="pricing-cached-rate"
                  type="number"
                  step="any"
                  min="0"
                  value={cachedRate}
                  onChange={(e) => setCachedRate(e.target.value)}
                  placeholder="0"
                  className="h-8 font-mono text-xs tabular-nums"
                />
                <Button
                  type="submit"
                  size="sm"
                  className="h-8 shrink-0 px-2.5 text-xs"
                  disabled={upsertMutation.isPending}
                >
                  {upsertMutation.isPending ? "..." : isId ? "Simpan" : "Save"}
                </Button>
              </div>
            </div>
          </div>

          {activeError ? (
            <Alert variant="destructive" className="py-2">
              <AlertCircle
                className="h-3.5 w-3.5 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <AlertDescription className="font-mono text-xs">{activeError}</AlertDescription>
            </Alert>
          ) : null}
        </form>

        {/* Search & Catalog List */}
        <div className="shrink-0 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-4 py-2">
          <div className="relative">
            <Search
              className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
              aria-hidden="true"
            />
            <Input
              aria-label={isId ? "Cari tarif model" : "Filter model pricing"}
              type="text"
              value={filterText}
              onChange={(e) => setFilterText(e.target.value)}
              placeholder={
                isId
                  ? "Filter katalog model (klik baris untuk edit)..."
                  : "Filter model catalog (click a row to edit)..."
              }
              className="h-8 pl-8 text-xs"
            />
          </div>
        </div>

        <div className="min-h-0 flex-1 overflow-y-auto">
          {pricingQuery.isPending ? (
            <div className="space-y-2 p-4">
              {[0, 1, 2, 3].map((idx) => (
                <Skeleton key={idx} className="h-9 w-full rounded" />
              ))}
            </div>
          ) : pricingQuery.isError ? (
            <div className="p-4">
              <Alert variant="destructive">
                <AlertCircle
                  className="h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <AlertDescription className="font-mono text-xs">
                  {pricingQuery.error.message}
                </AlertDescription>
              </Alert>
            </div>
          ) : filteredPricing.length === 0 ? (
            <div className="p-6 text-center text-xs text-[var(--text-muted)]">
              {isId ? "Tidak ada model yang cocok." : "No matching models."}
            </div>
          ) : (
            <table className="w-full border-collapse text-left text-xs">
              <thead className="sticky top-0 z-10 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[10px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                <tr>
                  <th scope="col" className="px-4 py-2">
                    Model
                  </th>
                  <th scope="col" className="px-2 py-2 text-right">
                    In
                  </th>
                  <th scope="col" className="px-2 py-2 text-right">
                    Out
                  </th>
                  <th scope="col" className="px-2 py-2 text-right">
                    Cached
                  </th>
                  <th scope="col" className="px-4 py-2 text-right">
                    {isId ? "Sumber" : "Source"}
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {filteredPricing.map((item) => (
                  <tr
                    key={item.model}
                    onClick={() => handleSelectItem(item)}
                    className="cursor-pointer transition-colors hover:bg-[var(--bg-surface-raised)]"
                  >
                    <td className="max-w-[220px] truncate px-4 py-2 font-mono text-xs text-[var(--text-primary)]">
                      {item.model}
                    </td>
                    <td className="px-2 py-2 text-right font-mono text-[11px] text-[var(--text-secondary)] tabular-nums">
                      ${item.input_per_million}
                    </td>
                    <td className="px-2 py-2 text-right font-mono text-[11px] text-[var(--text-secondary)] tabular-nums">
                      ${item.output_per_million}
                    </td>
                    <td className="px-2 py-2 text-right font-mono text-[11px] text-[var(--text-muted)] tabular-nums">
                      ${item.cached_per_million}
                    </td>
                    <td className="px-4 py-2 text-right">
                      <div className="inline-flex items-center gap-1.5">
                        <Badge variant={item.source === "custom" ? "default" : "secondary"}>
                          {item.source}
                        </Badge>
                        {item.source === "custom" ? (
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-sm"
                            className="h-6 w-6"
                            disabled={deleteMutation.isPending}
                            aria-label={
                              isId
                                ? `Reset tarif ${item.model} ke bawaan`
                                : `Reset ${item.model} pricing to default`
                            }
                            title={isId ? "Reset ke bawaan" : "Reset to default"}
                            onClick={(e) => {
                              e.stopPropagation();
                              deleteMutation.mutate(item.model);
                            }}
                          >
                            <RotateCcw className="h-3 w-3" aria-hidden="true" />
                          </Button>
                        ) : null}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
