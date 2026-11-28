import { BookOpen, Check, Clock, Cpu, Search, Sparkles, Target, Workflow, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import type { ComboStrategy } from "./schemas";
import { COMBO_STRATEGY_OPTIONS, STRATEGY_CATEGORIES, type StrategyCategory } from "./strategies";
import { STRATEGY_GUIDE_CONTENT } from "./strategyGuideContent";

interface StrategyGuideModalProps {
  open: boolean;
  selectedStrategy: ComboStrategy;
  onSelectStrategy?: (strategy: ComboStrategy) => void;
  onClose: () => void;
  isId: boolean;
}

export function StrategyGuideModal({
  open,
  selectedStrategy,
  onSelectStrategy,
  onClose,
  isId,
}: StrategyGuideModalProps) {
  const [search, setSearch] = useState("");
  const [activeCategory, setActiveCategory] = useState<"all" | StrategyCategory>("all");

  useEffect(() => {
    if (!open) return;
    setSearch("");
    setActiveCategory("all");

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

  const categoryLabels = useMemo(() => {
    const map = new Map<StrategyCategory, string>();
    for (const cat of STRATEGY_CATEGORIES) {
      map.set(cat.id, isId ? cat.labelId : cat.labelEn);
    }
    return map;
  }, [isId]);

  const filteredStrategies = useMemo(() => {
    const q = search.trim().toLowerCase();
    return COMBO_STRATEGY_OPTIONS.filter((item) => {
      if (activeCategory !== "all" && item.category !== activeCategory) {
        return false;
      }
      if (!q) return true;
      const guide = STRATEGY_GUIDE_CONTENT[item.value];
      const haystack = [
        item.value,
        item.labelEn,
        item.labelId,
        item.descEn,
        item.descId,
        guide.howItWorksEn,
        guide.howItWorksId,
        guide.bestForEn,
        guide.bestForId,
        guide.whenToUseEn,
        guide.whenToUseId,
      ]
        .join(" ")
        .toLowerCase();
      return haystack.includes(q);
    });
  }, [search, activeCategory]);

  if (!open) {
    return null;
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5"
      role="dialog"
      aria-modal="true"
      aria-labelledby="strategy-guide-modal-title"
      data-testid="strategy-guide-modal"
    >
      <button
        type="button"
        aria-label={isId ? "Tutup panduan strategi" : "Close strategy guide"}
        onClick={onClose}
        className="fixed inset-0 bg-black/60"
      />
      <div className="relative z-10 flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-xl border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
        {/* Header */}
        <div className="flex items-start justify-between gap-3 border-b border-[var(--border-subtle)] px-5 py-4">
          <div className="flex items-start gap-3">
            <div className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] text-[var(--accent-primary)]">
              <BookOpen className="h-4 w-4" aria-hidden="true" />
            </div>
            <div>
              <h3
                id="strategy-guide-modal-title"
                className="text-base font-semibold text-[var(--text-primary)]"
              >
                {isId ? "Panduan Strategi Routing Combo" : "Combo Routing Strategy Guide"}
              </h3>
              <p className="mt-0.5 text-xs text-[var(--text-secondary)]">
                {isId
                  ? "Pahami cara kerja, kecocokan workload, dan waktu penggunaan terbaik untuk semua 22 strategi routing."
                  : "Understand how each strategy works, what workloads it fits best, and when to use it across all 22 routing strategies."}
              </p>
            </div>
          </div>

          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={isId ? "Tutup panduan strategi" : "Close strategy guide"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        {/* Search & Category Filter */}
        <div className="space-y-3 border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3.5">
          <div className="relative">
            <Search
              className="pointer-events-none absolute top-1/2 left-3 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
              aria-hidden="true"
            />
            <Input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={
                isId
                  ? "Cari strategi (mis. failover, hemat biaya, cache, oauth, paralel)..."
                  : "Search strategies (e.g. failover, cost, cache, oauth, parallel)..."
              }
              className="h-9 pl-9 text-xs"
            />
          </div>

          <div className="flex flex-wrap items-center gap-1.5">
            <button
              type="button"
              onClick={() => setActiveCategory("all")}
              className={`rounded-full border px-3 py-1 text-xs font-medium transition-colors ${
                activeCategory === "all"
                  ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                  : "border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
              }`}
            >
              {isId
                ? `Semua (${COMBO_STRATEGY_OPTIONS.length})`
                : `All (${COMBO_STRATEGY_OPTIONS.length})`}
            </button>
            {STRATEGY_CATEGORIES.map((cat) => {
              const count = COMBO_STRATEGY_OPTIONS.filter((s) => s.category === cat.id).length;
              const isActive = activeCategory === cat.id;
              return (
                <button
                  key={cat.id}
                  type="button"
                  onClick={() => setActiveCategory(cat.id)}
                  className={`rounded-full border px-3 py-1 text-xs font-medium transition-colors ${
                    isActive
                      ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                      : "border-[var(--border-strong)] bg-[var(--bg-surface)] text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
                  }`}
                >
                  {isId ? cat.labelId : cat.labelEn} ({count})
                </button>
              );
            })}
          </div>
        </div>

        {/* Strategy Cards List */}
        <div className="custom-scrollbar flex-1 space-y-3 overflow-y-auto px-5 py-4">
          {filteredStrategies.length === 0 ? (
            <div className="rounded-lg border border-dashed border-[var(--border-strong)] p-10 text-center">
              <Sparkles className="mx-auto h-6 w-6 text-[var(--text-muted)]" aria-hidden="true" />
              <p className="mt-2 text-sm font-medium text-[var(--text-primary)]">
                {isId ? "Strategi tidak ditemukan" : "No matching strategies"}
              </p>
              <p className="mt-1 text-xs text-[var(--text-muted)]">
                {isId
                  ? "Coba kata kunci lain atau pilih kategori Semua."
                  : "Try another keyword or switch back to All categories."}
              </p>
            </div>
          ) : (
            filteredStrategies.map((item) => {
              const isSelected = item.value === selectedStrategy;
              const categoryLabel = categoryLabels.get(item.category) ?? item.category;
              const guide = STRATEGY_GUIDE_CONTENT[item.value];

              return (
                <article
                  key={item.value}
                  data-testid={`strategy-guide-item-${item.value}`}
                  className={`rounded-lg border p-4 transition-colors ${
                    isSelected
                      ? "border-[var(--accent-primary)] bg-[var(--bg-canvas)]"
                      : "border-[var(--border-subtle)] bg-[var(--bg-canvas)]/60 hover:border-[var(--border-strong)]"
                  }`}
                >
                  {/* Card Header */}
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <div className="space-y-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <h4 className="text-sm font-semibold text-[var(--text-primary)]">
                          {isId ? item.labelId : item.labelEn}
                        </h4>
                        <Badge variant="outline" className="font-mono text-[11px]">
                          {item.value}
                        </Badge>
                        <Badge variant="secondary" className="text-[10px]">
                          {categoryLabel}
                        </Badge>
                        {isSelected ? (
                          <Badge variant="default" className="gap-1 text-[10px]">
                            <Check className="h-3 w-3" aria-hidden="true" />
                            <span>{isId ? "Dipilih" : "Selected"}</span>
                          </Badge>
                        ) : null}
                      </div>

                      {/* Flow Preview */}
                      <div className="inline-flex items-center gap-1.5 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-2.5 py-1 font-mono text-[11px] text-[var(--accent-text)]">
                        <Workflow className="h-3 w-3 shrink-0" aria-hidden="true" />
                        <span>{item.flowPreview}</span>
                      </div>
                    </div>

                    {onSelectStrategy ? (
                      <Button
                        type="button"
                        size="xs"
                        variant={isSelected ? "secondary" : "default"}
                        onClick={() => {
                          onSelectStrategy(item.value);
                          onClose();
                        }}
                      >
                        {isSelected
                          ? isId
                            ? "Sedang Dipilih"
                            : "Selected"
                          : isId
                            ? "Pakai Strategi Ini"
                            : "Use This Strategy"}
                      </Button>
                    ) : null}
                  </div>

                  {/* Structured 3-Part Explanation */}
                  <div className="mt-3 grid grid-cols-1 gap-2.5 md:grid-cols-3">
                    <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-3">
                      <div className="flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide text-[var(--text-primary)]">
                        <Cpu
                          className="h-3.5 w-3.5 text-[var(--accent-primary)]"
                          aria-hidden="true"
                        />
                        <span>{isId ? "Cara Kerja" : "How It Works"}</span>
                      </div>
                      <p className="mt-1.5 text-xs leading-relaxed text-[var(--text-secondary)]">
                        {isId ? guide.howItWorksId : guide.howItWorksEn}
                      </p>
                    </div>

                    <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-3">
                      <div className="flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide text-[var(--text-primary)]">
                        <Target
                          className="h-3.5 w-3.5 text-[var(--status-ok)]"
                          aria-hidden="true"
                        />
                        <span>{isId ? "Cocok Untuk" : "Best For"}</span>
                      </div>
                      <p className="mt-1.5 text-xs leading-relaxed text-[var(--text-secondary)]">
                        {isId ? guide.bestForId : guide.bestForEn}
                      </p>
                    </div>

                    <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-3">
                      <div className="flex items-center gap-1.5 text-[11px] font-semibold uppercase tracking-wide text-[var(--text-primary)]">
                        <Clock
                          className="h-3.5 w-3.5 text-[var(--status-warn)]"
                          aria-hidden="true"
                        />
                        <span>{isId ? "Kapan Digunakan" : "When To Use"}</span>
                      </div>
                      <p className="mt-1.5 text-xs leading-relaxed text-[var(--text-secondary)]">
                        {isId ? guide.whenToUseId : guide.whenToUseEn}
                      </p>
                    </div>
                  </div>
                </article>
              );
            })
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3">
          <p className="text-xs text-[var(--text-muted)]">
            {isId
              ? `Menampilkan ${filteredStrategies.length} dari ${COMBO_STRATEGY_OPTIONS.length} strategi.`
              : `Showing ${filteredStrategies.length} of ${COMBO_STRATEGY_OPTIONS.length} strategies.`}
          </p>
          <Button type="button" variant="outline" size="sm" onClick={onClose}>
            {isId ? "Tutup" : "Close"}
          </Button>
        </div>
      </div>
    </div>
  );
}
