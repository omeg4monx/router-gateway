import { Check, ChevronsUpDown, FileText, Search, Sparkles } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import type { UILocale } from "@/shared/store/uiStore";
import type { BuiltinPromptTemplate } from "./builtinPromptTemplates";
import { estimatePromptTokens, MODEL_EDITOR_COPY } from "./parameterGuide";

export function filterPromptTemplates(
  templates: readonly BuiltinPromptTemplate[],
  query: string,
): BuiltinPromptTemplate[] {
  const normalized = query.trim().toLowerCase();
  if (!normalized) return [...templates];
  return templates.filter(
    (tpl) =>
      tpl.name.toLowerCase().includes(normalized) ||
      tpl.content.slice(0, 240).toLowerCase().includes(normalized),
  );
}

function templatePreviewSnippet(content: string): string {
  const compact = content.replace(/\s+/g, " ").trim();
  return compact.length > 72 ? `${compact.slice(0, 72)}…` : compact;
}

interface PromptTemplateComboboxProps {
  selectedId: string;
  templates: readonly BuiltinPromptTemplate[];
  locale: UILocale;
  onSelectCustom: () => void;
  onSelectTemplate: (template: BuiltinPromptTemplate) => void;
}

export function PromptTemplateCombobox({
  selectedId,
  templates,
  locale,
  onSelectCustom,
  onSelectTemplate,
}: PromptTemplateComboboxProps) {
  const copy = MODEL_EDITOR_COPY[locale];
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const selectedTemplate = useMemo(
    () => templates.find((tpl) => tpl.id === selectedId) ?? null,
    [templates, selectedId],
  );

  const filtered = useMemo(() => filterPromptTemplates(templates, search), [templates, search]);

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

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        role="combobox"
        aria-expanded={open}
        aria-haspopup="listbox"
        aria-label={copy.promptTemplatePickerLabel}
        onClick={() => setOpen((prev) => !prev)}
        className="flex h-8 w-full cursor-pointer items-center justify-between gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-2.5 text-left text-xs transition-colors hover:border-[var(--border-strong)] focus:border-[var(--accent-primary)] focus:outline-none"
      >
        <div className="flex min-w-0 flex-1 items-center gap-1.5">
          {selectedTemplate ? (
            <Sparkles
              className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
              aria-hidden="true"
            />
          ) : (
            <FileText
              className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
              aria-hidden="true"
            />
          )}
          <span className="truncate font-medium text-[var(--text-primary)]">
            {selectedTemplate ? selectedTemplate.name : copy.promptFreeTextLabel}
          </span>
          <span className="shrink-0 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface)] px-1.5 py-0.2 font-mono text-[10px] text-[var(--text-muted)]">
            {templates.length} {copy.promptTemplatesUnit}
          </span>
        </div>
        <ChevronsUpDown
          className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]"
          aria-hidden="true"
        />
      </button>

      {open ? (
        <div className="absolute left-0 right-0 top-full z-40 mt-1 flex max-h-64 flex-col overflow-hidden rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-xl">
          <div className="flex shrink-0 items-center gap-2 border-b border-[var(--border-subtle)] px-2.5 py-2">
            <Search className="h-3.5 w-3.5 shrink-0 text-[var(--text-muted)]" aria-hidden="true" />
            <input
              ref={searchInputRef}
              type="text"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder={copy.promptTemplateSearchPlaceholder}
              aria-label={copy.promptTemplateSearchPlaceholder}
              className="w-full bg-transparent text-xs text-[var(--text-primary)] placeholder:text-[var(--text-muted)] focus:outline-none"
            />
          </div>

          <div role="listbox" className="min-h-0 flex-1 overflow-y-auto p-1">
            {/* Free Text option */}
            <button
              type="button"
              role="option"
              aria-selected={selectedId === "custom"}
              onClick={() => {
                onSelectCustom();
                setSearch("");
                setOpen(false);
              }}
              className={cn(
                "mb-1 flex w-full cursor-pointer items-center justify-between gap-2 rounded border-b border-[var(--border-subtle)] px-2 py-1.5 text-left text-xs transition-colors",
                selectedId === "custom"
                  ? "bg-[var(--accent-primary)]/15 text-[var(--accent-text)]"
                  : "text-[var(--text-primary)] hover:bg-[var(--bg-surface-raised)]",
              )}
            >
              <div className="min-w-0 flex-1">
                <p className="font-medium">{copy.promptFreeTextLabel}</p>
                <p className="truncate text-[10px] text-[var(--text-muted)]">
                  {copy.promptFreeTextHint}
                </p>
              </div>
              {selectedId === "custom" ? (
                <Check
                  className="h-3.5 w-3.5 shrink-0 text-[var(--accent-primary)]"
                  aria-hidden="true"
                />
              ) : null}
            </button>

            {filtered.length === 0 ? (
              <p className="px-2.5 py-3 text-center text-xs text-[var(--text-muted)]">
                {copy.promptTemplateEmpty}
              </p>
            ) : (
              filtered.map((tpl) => {
                const isSelected = tpl.id === selectedId;
                const tokens = estimatePromptTokens(tpl.content);
                return (
                  <button
                    key={tpl.id}
                    type="button"
                    role="option"
                    aria-selected={isSelected}
                    onClick={() => {
                      onSelectTemplate(tpl);
                      setSearch("");
                      setOpen(false);
                    }}
                    className={cn(
                      "flex w-full cursor-pointer items-center justify-between gap-2 rounded px-2 py-1.5 text-left text-xs transition-colors",
                      isSelected
                        ? "bg-[var(--accent-primary)]/15 text-[var(--accent-text)]"
                        : "text-[var(--text-primary)] hover:bg-[var(--bg-surface-raised)]",
                    )}
                  >
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-1.5">
                        <span className="truncate font-medium">{tpl.name}</span>
                        <span className="shrink-0 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-1 font-mono text-[9px] text-[var(--text-muted)]">
                          {tpl.position === "suffix"
                            ? copy.promptSuffixLabel
                            : copy.promptPrefixLabel}
                        </span>
                        <span className="shrink-0 font-mono text-[10px] text-[var(--text-muted)]">
                          ~{tokens} tok
                        </span>
                      </div>
                      <p className="mt-0.5 truncate text-[10px] text-[var(--text-muted)]">
                        {templatePreviewSnippet(tpl.content)}
                      </p>
                    </div>
                    {isSelected ? (
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
