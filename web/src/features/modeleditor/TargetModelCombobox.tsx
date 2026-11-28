import { Check, ChevronsUpDown, Search, X } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import type { CatalogModelRow } from "@/features/models/types";
import { cn } from "@/shared/lib/utils";
import type { UILocale } from "@/shared/store/uiStore";
import { MODEL_EDITOR_COPY } from "./parameterGuide";

export interface TargetModelOption {
  id: string;
  name: string;
  provider: string;
  maxOutputTokens?: number;
}

export function buildTargetModelOptions(rows: readonly CatalogModelRow[]): TargetModelOption[] {
  const seen = new Set<string>();
  const options: TargetModelOption[] = [];
  for (const row of rows) {
    if (row.hidden) continue;
    const modelId = row.id.trim();
    if (!modelId || seen.has(modelId)) continue;
    seen.add(modelId);
    options.push({
      id: modelId,
      name: row.name?.trim() || modelId,
      provider: row.provider?.trim() || row.providerId,
      maxOutputTokens:
        typeof row.max_output_tokens === "number" && row.max_output_tokens > 0
          ? row.max_output_tokens
          : undefined,
    });
  }
  return options.sort((a, b) => a.id.localeCompare(b.id));
}

export function filterTargetModelOptions(
  options: readonly TargetModelOption[],
  query: string,
): TargetModelOption[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return options.slice(0, 80);
  return options
    .filter(
      (opt) =>
        opt.id.toLowerCase().includes(normalized) ||
        opt.name.toLowerCase().includes(normalized) ||
        opt.provider.toLowerCase().includes(normalized),
    )
    .slice(0, 80);
}

interface TargetModelComboboxProps {
  id: string;
  value: string;
  options: readonly TargetModelOption[];
  isLoading: boolean;
  locale: UILocale;
  onChange: (value: string) => void;
}

export function TargetModelCombobox({
  id,
  value,
  options,
  isLoading,
  locale,
  onChange,
}: TargetModelComboboxProps) {
  const copy = MODEL_EDITOR_COPY[locale];
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const filtered = useMemo(() => filterTargetModelOptions(options, search), [options, search]);

  useEffect(() => {
    if (!open) return;
    const frame = window.requestAnimationFrame(() => {
      searchInputRef.current?.focus();
    });
    const handlePointerDown = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handlePointerDown);
    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.cancelAnimationFrame(frame);
      document.removeEventListener("mousedown", handlePointerDown);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [open]);

  const handleSelect = (nextValue: string) => {
    onChange(nextValue);
    setSearch("");
    setOpen(false);
  };

  const handleSearchKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (filtered.length === 0) return;
    if (event.key === "ArrowDown") {
      event.preventDefault();
      setActiveIndex((prev) => (prev + 1) % filtered.length);
    } else if (event.key === "ArrowUp") {
      event.preventDefault();
      setActiveIndex((prev) => (prev - 1 + filtered.length) % filtered.length);
    } else if (event.key === "Enter") {
      event.preventDefault();
      const chosen = filtered[activeIndex];
      if (chosen) {
        handleSelect(chosen.id);
      }
    }
  };

  const activeOptionId = open && filtered[activeIndex] ? `${id}-opt-${activeIndex}` : undefined;

  return (
    <div ref={containerRef} className="relative">
      <div className="relative flex items-center">
        <button
          id={id}
          type="button"
          role="combobox"
          aria-expanded={open}
          aria-haspopup="listbox"
          aria-controls={`${id}-listbox`}
          onClick={() => {
            setActiveIndex(0);
            setOpen((prev) => !prev);
          }}
          className={cn(
            "flex h-9 w-full cursor-pointer items-center justify-between gap-1.5 rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 text-left text-xs transition-colors focus:border-[var(--accent-primary)] focus:outline-none",
            value ? "pr-12 font-mono text-[var(--text-primary)]" : "text-[var(--text-muted)]",
          )}
        >
          <span className="min-w-0 flex-1 truncate">
            {value || (isLoading ? copy.targetLoading : copy.targetPlaceholder)}
          </span>
          <ChevronsUpDown
            className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
            aria-hidden="true"
          />
        </button>

        {value ? (
          <button
            type="button"
            onClick={() => handleSelect("")}
            aria-label={copy.targetClearLabel}
            className="absolute right-7 inline-flex h-5 w-5 cursor-pointer items-center justify-center rounded text-[var(--text-muted)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
          >
            <X className="h-3 w-3" aria-hidden="true" />
          </button>
        ) : null}
      </div>

      {open ? (
        <div className="absolute left-0 right-0 top-full z-40 mt-1 flex max-h-60 flex-col overflow-hidden rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-xl">
          <div className="flex shrink-0 items-center gap-2 border-b border-[var(--border-subtle)] px-2.5 py-2">
            <Search className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
            <input
              ref={searchInputRef}
              type="text"
              value={search}
              onChange={(event) => {
                setSearch(event.target.value);
                setActiveIndex(0);
              }}
              onKeyDown={handleSearchKeyDown}
              placeholder={copy.targetSearchPlaceholder}
              aria-label={copy.targetSearchPlaceholder}
              aria-activedescendant={activeOptionId}
              className="w-full bg-transparent text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none"
            />
          </div>

          <div id={`${id}-listbox`} role="listbox" className="min-h-0 flex-1 overflow-y-auto p-1">
            {isLoading ? (
              <p className="px-2.5 py-3 text-center text-xs text-[var(--text-muted)]">
                {copy.targetLoading}
              </p>
            ) : filtered.length === 0 ? (
              <p className="px-2.5 py-3 text-center text-xs text-[var(--text-muted)]">
                {copy.targetEmpty}
              </p>
            ) : (
              filtered.map((opt, idx) => {
                const selected = opt.id === value;
                const isActive = activeIndex === idx;
                return (
                  <button
                    key={opt.id}
                    id={`${id}-opt-${idx}`}
                    type="button"
                    role="option"
                    aria-selected={selected}
                    onMouseEnter={() => setActiveIndex(idx)}
                    onClick={() => handleSelect(opt.id)}
                    className={cn(
                      "flex w-full cursor-pointer items-center justify-between gap-2 rounded px-2 py-1.5 text-left text-xs transition-colors",
                      selected
                        ? "bg-[var(--accent-primary)]/15 text-[var(--accent-text)]"
                        : isActive
                          ? "bg-[var(--bg-surface-raised)] text-[var(--text-primary)]"
                          : "text-[var(--text-primary)] hover:bg-[var(--bg-surface-raised)]",
                    )}
                  >
                    <div className="min-w-0 flex-1">
                      <p className="truncate font-mono text-xs">{opt.id}</p>
                      <p className="truncate text-[10px] text-[var(--text-muted)]">
                        {opt.provider}
                      </p>
                    </div>
                    {selected ? (
                      <Check
                        className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
                        aria-hidden="true"
                      />
                    ) : null}
                  </button>
                );
              })
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}
